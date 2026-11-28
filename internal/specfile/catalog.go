package specfile

import (
	"bytes"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"slices"
	"strings"

	"github.com/dotspec-md/dotspec/internal/diagnostic"
	"go.yaml.in/yaml/v3"
)

const SpecificationsField = "specifications"
const VerificationField = "verification"

// CatalogEntry identifies a named source and the exact authored scalars for export.
// Defines is empty for a verification procedure.
type CatalogEntry struct {
	Name, Source *yaml.Node
	Defines      string
	Verification []CatalogEntry
	Flow         bool
	Location     diagnostic.Location
}

func SpecificationAnchor(name string) string {
	return "specification." + RenderedHeadingIDs{}.Next(name)
}

func (md *Markdown) Catalog(kind string) ([]CatalogEntry, error) {
	var result []CatalogEntry
	for _, field := range []string{SpecificationsField, VerificationField} {
		list := md.Value(field)
		if list == nil {
			continue
		}
		if kind == "quality-model" {
			return nil, diagnostic.Error("catalog.invalid", "A quality model is reference material. Declare specification sources and verification on a specification.", md.MetadataLocation(field))
		}
		entries, err := catalogList(list, field)
		if err != nil {
			return nil, err
		}
		result = append(result, entries...)
	}
	return result, nil
}

func catalogList(list *yaml.Node, field string) ([]CatalogEntry, error) {
	at := func(n *yaml.Node) diagnostic.Location {
		return diagnostic.Location{Line: n.Line + 1, Column: n.Column, Element: field}
	}
	invalid := func(n *yaml.Node, message string) error { return diagnostic.Error("catalog.invalid", message, at(n)) }
	if list.Kind != yaml.SequenceNode {
		return nil, invalid(list, "Write "+field+" as a YAML sequence of named sources.")
	}
	var entries []CatalogEntry
	names := map[string]diagnostic.Location{}
	for _, item := range list.Content {
		if item.Kind != yaml.MappingNode {
			return nil, invalid(item, "Write each entry as a mapping with a name and source.")
		}
		values := map[string]*yaml.Node{}
		for i := 0; i < len(item.Content); i += 2 {
			k, v := item.Content[i], item.Content[i+1]
			if values[k.Value] != nil {
				return nil, invalid(k, "Repeated key "+k.Value+".")
			}
			switch k.Value {
			case "name", "source":
			case "defines", VerificationField:
				if field != SpecificationsField {
					return nil, invalid(k, "Unexpected verification key "+k.Value+".")
				}
			default:
				return nil, invalid(k, "Unknown "+field+" key "+k.Value+".")
			}
			if k.Value != VerificationField && (v.Kind != yaml.ScalarNode || v.ShortTag() != "!!str" || v.Value == "" || strings.TrimSpace(v.Value) != v.Value || strings.ContainsAny(v.Value, "\r\n")) {
				return nil, invalid(v, "Give "+k.Value+" a nonempty single-line text value without surrounding spaces.")
			}
			values[k.Value] = v
		}
		for _, key := range []string{"name", "source"} {
			if values[key] == nil {
				return nil, invalid(item, "Give each entry a "+key+".")
			}
		}
		entry := CatalogEntry{Name: values["name"], Source: values["source"], Flow: item.Style&yaml.FlowStyle != 0 || list.Style&yaml.FlowStyle != 0, Location: at(values["source"])}
		locator, err := url.Parse(entry.Source.Value)
		if err != nil || (!locator.IsAbs() && (locator.Path == "" || locator.Host != "" || strings.HasPrefix(locator.Path, "/") || strings.Contains(locator.Path, "\\"))) {
			return nil, invalid(entry.Source, "Give the source a path or URL, optionally with a fragment.")
		}
		anchor := SpecificationAnchor(entry.Name.Value)
		if anchor == "specification." {
			return nil, invalid(entry.Name, "Give the name at least one addressable character.")
		}
		if prior, ok := names[anchor]; ok {
			return nil, diagnostic.Error("catalog.duplicate", "Choose distinct names in this collection.", at(entry.Name), prior)
		}
		names[anchor] = at(entry.Name)
		if field == SpecificationsField {
			defines := values["defines"]
			if defines == nil {
				return nil, invalid(item, "Give each specification a defines category.")
			}
			categories := schemaCatalog().Enums["SpecificationCategory"]
			if !slices.Contains(categories, defines.Value) {
				return nil, invalid(defines, "Choose defines from "+strings.Join(categories, ", ")+".")
			}
			entry.Defines = defines.Value
			if nested := values[VerificationField]; nested != nil {
				procedures, err := catalogList(nested, VerificationField)
				if err != nil {
					return nil, err
				}
				entry.Verification = procedures
			}
		}
		entries = append(entries, entry)
	}
	return entries, nil
}

// NativeDefinition is a named declaration in a native source.
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
