package specfile

import (
	"bytes"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"slices"
	"strings"

	"github.com/dotspec-md/dotspec/internal/diagnostic"
	"go.yaml.in/yaml/v3"
)

// CompanionsField is the frontmatter key that lists a specification's companions.
const CompanionsField = "companions"

// A companion carries part of a specification's intent in another form, such
// as a formal model, an executable example, a test or an explanatory document.
// The role states how it relates to the specification; a link alone confers no
// authority.
type Companion struct {
	Target   string
	Role     string
	Location diagnostic.Location
}

// CompanionRoles lists the closed role vocabulary in presentation order.
var CompanionRoles = []string{"formalizes", "illustrates", "verifies", "explains"}

// CompanionRoleLabel is the reader-facing phrase for a role, read from the
// specification toward its companion.
func CompanionRoleLabel(role string) string {
	switch role {
	case "formalizes":
		return "Formalized by"
	case "illustrates":
		return "Illustrated by"
	case "verifies":
		return "Verified by"
	case "explains":
		return "Explained by"
	}
	return role
}

var companionKeys = []string{"target", "role", "note"}

// Companions validates the frontmatter list of a specification of this kind
// and returns its entries with the location of each target. A missing key
// yields no entries.
func (md *Markdown) Companions(kind string) ([]Companion, error) {
	list := md.Value(CompanionsField)
	if list == nil {
		return nil, nil
	}
	at := func(n *yaml.Node) diagnostic.Location {
		return diagnostic.Location{Line: n.Line + 1, Column: n.Column, Element: CompanionsField}
	}
	if kind == "quality-model" {
		return nil, diagnostic.Error("companion.invalid", "A quality model is reference material. Declare companions on the specifications that use it.", at(list))
	}
	if list.Kind != yaml.SequenceNode {
		return nil, diagnostic.Error("companion.invalid", "List companions as a YAML sequence of entries with target and role.", at(list))
	}
	var result []Companion
	seen := map[[2]string]diagnostic.Location{}
	for _, item := range list.Content {
		if item.Kind != yaml.MappingNode {
			return nil, diagnostic.Error("companion.invalid", "Give each companion a mapping with target and role.", at(item))
		}
		entry := Companion{Location: at(item)}
		values := map[string]*yaml.Node{}
		for i := 0; i+1 < len(item.Content); i += 2 {
			key, value := item.Content[i], item.Content[i+1]
			if !slices.Contains(companionKeys, key.Value) {
				return nil, diagnostic.Error("companion.invalid", fmt.Sprintf("Unknown companion key %q. Use target, role and optionally note.", key.Value), at(key))
			}
			if values[key.Value] != nil {
				return nil, diagnostic.Error("companion.invalid", fmt.Sprintf("Duplicate companion key %q.", key.Value), at(key))
			}
			if value.Kind != yaml.ScalarNode || value.ShortTag() != "!!str" {
				return nil, diagnostic.Error("companion.invalid", fmt.Sprintf("Give companion %s a single-line text value.", key.Value), at(value))
			}
			if strings.ContainsAny(value.Value, "\r\n") || strings.TrimSpace(value.Value) != value.Value {
				return nil, diagnostic.Error("companion.invalid", fmt.Sprintf("Give companion %s a single-line value without surrounding spaces.", key.Value), at(value))
			}
			values[key.Value] = value
		}
		target, role := values["target"], values["role"]
		if target == nil || target.Value == "" {
			return nil, diagnostic.Error("companion.invalid", "Give each companion a target: a relative path or URL, optionally with a #definition fragment.", at(item))
		}
		if role == nil || role.Value == "" {
			return nil, diagnostic.Error("companion.role", "Give each companion a role: "+strings.Join(CompanionRoles, ", ")+".", at(item))
		}
		if !slices.Contains(CompanionRoles, role.Value) {
			return nil, diagnostic.Error("companion.role", fmt.Sprintf("Unknown companion role %q. Use one of %s.", role.Value, strings.Join(CompanionRoles, ", ")), at(role))
		}
		if strings.HasPrefix(target.Value, "#") {
			return nil, diagnostic.Error("companion.invalid", "A companion is another file or resource. Link to this specification's own passages from its body.", at(target))
		}
		entry.Target, entry.Role, entry.Location = target.Value, role.Value, at(target)
		if note := values["note"]; note != nil {
			if note.Value == "" {
				return nil, diagnostic.Error("companion.invalid", "Omit an empty note.", at(note))
			}
		}
		key := [2]string{entry.Target, entry.Role}
		if previous, ok := seen[key]; ok {
			return nil, diagnostic.Error("companion.duplicate", fmt.Sprintf("%s already lists %q with role %s.", CompanionsField, entry.Target, entry.Role), entry.Location, previous)
		}
		seen[key] = entry.Location
		result = append(result, entry)
	}
	return result, nil
}

// NativeDefinition is a named declaration in a non-Markdown companion.
type NativeDefinition struct {
	Name string
	Line int
}

type definitionRule struct {
	pattern *regexp.Regexp
	list    bool // the captured group is a comma-separated list of names
	slug    bool // the captured title becomes a heading-style fragment
	last    bool // the captured path names its final segment
	// continues reports whether a matched line leaves a list open on the next line.
	continues func(captured string) bool
	// opens starts a grouped block whose indented lines each declare a name.
	opens bool
}

var (
	definitionName = regexp.MustCompile(`^[A-Za-z_][A-Za-z0-9_]*$`)
	blockMember    = regexp.MustCompile(`^\s+([A-Za-z_][A-Za-z0-9_]*)\b`)
	blockEnd       = regexp.MustCompile(`^\)`)
)

// A TLA+ CONSTANTS or VARIABLES list continues after a bare keyword or a trailing comma.
func openList(captured string) bool {
	captured = strings.TrimSpace(captured)
	return captured == "" || strings.HasSuffix(captured, ",")
}

// Definition lookup supports addressability only: it recognizes declaration
// forms by their leading keywords and never interprets semantics.
var nativeDefinitionRules = map[string][]definitionRule{
	".qnt": {
		{pattern: regexp.MustCompile(`^\s*(?:pure\s+)?(?:module|const|var|type|val|def|action|run|temporal|assume|nondet)\s+([A-Za-z_][A-Za-z0-9_]*)`)},
	},
	".als": {
		{pattern: regexp.MustCompile(`^\s*(?:(?:abstract|one|lone|some|private|var)\s+)*sig\s+([A-Za-z_][\w'"]*(?:\s*,\s*[A-Za-z_][\w'"]*)*)`), list: true},
		{pattern: regexp.MustCompile(`^\s*module\s+([A-Za-z_][\w'"/]*)`), last: true},
		{pattern: regexp.MustCompile(`^\s*(?:private\s+)?(?:enum|fact|pred|fun|assert|check|run)\s+(?:[A-Za-z_][\w'"]*\.)?([A-Za-z_][\w'"]*)`)},
	},
	".tla": {
		{pattern: regexp.MustCompile(`^-*\s*MODULE\s+([A-Za-z_][A-Za-z0-9_]*)`)},
		{pattern: regexp.MustCompile(`^\s*(?:CONSTANTS?|VARIABLES?)\b(.*)$`), list: true, continues: openList},
		{pattern: regexp.MustCompile(`^\s*(?:THEOREM|LEMMA|COROLLARY|PROPOSITION|AXIOM|ASSUME)\s+([A-Za-z_][A-Za-z0-9_]*)\s*==`)},
		{pattern: regexp.MustCompile(`^(?:LOCAL\s+)?([A-Za-z_][A-Za-z0-9_]*)\s*(?:\((?:[^()]|\([^()]*\))*\)|\[[^\]]*\])?\s*==`)},
	},
	".feature": {
		{pattern: regexp.MustCompile(`^\s*(?:Feature|Rule|Background|Scenario Outline|Scenario Template|Scenario|Example):\s*(.+?)\s*$`), slug: true},
	},
	".go": {
		{pattern: regexp.MustCompile(`^type\s*\(\s*$`), opens: true},
		{pattern: regexp.MustCompile(`^func\s+(?:\([^)]*\)\s*)?([A-Za-z_][A-Za-z0-9_]*)`)},
		{pattern: regexp.MustCompile(`^type\s+([A-Za-z_][A-Za-z0-9_]*)`)},
	},
	".py": {
		{pattern: regexp.MustCompile(`^\s*(?:async\s+)?(?:def|class)\s+([A-Za-z_][A-Za-z0-9_]*)`)},
	},
	".txtar": {
		{pattern: regexp.MustCompile(`^-- (.+?) --\s*$`)},
	},
}

// NativeDefinitionFormat reports whether dotspec can look up named definitions
// in files with this name.
func NativeDefinitionFormat(name string) bool {
	return nativeDefinitionRules[strings.ToLower(filepath.Ext(name))] != nil
}

// NativeDefinitions lists the named declarations in a supported file, in
// source order. The first declaration of a repeated name wins.
func NativeDefinitions(name string, data []byte) []NativeDefinition {
	rules := nativeDefinitionRules[strings.ToLower(filepath.Ext(name))]
	if rules == nil {
		return nil
	}
	var result []NativeDefinition
	seen := map[string]bool{}
	slugs := RenderedHeadingIDs{}
	add := func(name string, line int) {
		name = strings.TrimSpace(name)
		if name != "" && !seen[name] {
			seen[name] = true
			result = append(result, NativeDefinition{name, line})
		}
	}
	addList := func(list string, line int) {
		for name := range strings.SplitSeq(list, ",") {
			if name = strings.TrimSpace(name); definitionName.MatchString(name) {
				add(name, line)
			}
		}
	}
	var continued *definitionRule // a list left open by the previous line
	inBlock := false
	for i, line := range bytes.Split(Normalize(data), []byte("\n")) {
		text := string(line)
		if inBlock {
			if blockEnd.MatchString(text) {
				inBlock = false
			} else if match := blockMember.FindStringSubmatch(text); match != nil {
				add(match[1], i+1)
			}
			continue
		}
		if continued != nil {
			if strings.TrimSpace(text) == "" {
				continue
			}
			addList(text, i+1)
			if !continued.continues(text) {
				continued = nil
			}
			continue
		}
		for r := range rules {
			rule := &rules[r]
			match := rule.pattern.FindStringSubmatch(text)
			if match == nil {
				continue
			}
			switch {
			case rule.opens:
				inBlock = true
			case rule.slug:
				add(slugs.Next(match[1]), i+1)
			case rule.list:
				addList(match[1], i+1)
				if rule.continues != nil && rule.continues(match[1]) {
					continued = rule
				}
			case rule.last:
				add(path.Base(match[1]), i+1)
			default:
				add(match[1], i+1)
			}
			break
		}
	}
	return result
}

// nativeDefinitionLine returns the line declaring name, or zero.
func nativeDefinitionLine(file string, data []byte, name string) int {
	for _, d := range NativeDefinitions(file, data) {
		if d.Name == name {
			return d.Line
		}
	}
	return 0
}

// NativeDefinitionAt reads a supported native file and returns the line
// declaring name, or zero.
func NativeDefinitionAt(path, name string) (int, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return 0, err
	}
	return nativeDefinitionLine(path, data, name), nil
}
