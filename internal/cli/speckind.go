package cli

import (
	"encoding/json"
	"fmt"
	"io"
	"strings"

	"github.com/craigsmitham/dotspec/internal/check"
	"github.com/craigsmitham/dotspec/internal/specfile"
	"github.com/craigsmitham/dotspec/internal/speckinds"
	"github.com/craigsmitham/dotspec/internal/system"
	"github.com/spf13/cobra"
)

// Specifications kept in a system's collection share one command convention: init takes a
// strict slug, show and check also accept a source path, reads and lists use schema version
// 2, and lists group members by their associated system.
const readSchemaVersion = 2

// readResult encodes {"schema_version": 2, "<key>": value}, keeping the version first.
type readResult struct {
	key   string
	value any
}

func (r readResult) MarshalJSON() ([]byte, error) {
	key, err := json.Marshal(r.key)
	if err != nil {
		return nil, err
	}
	value, err := json.Marshal(r.value)
	if err != nil {
		return nil, err
	}
	return fmt.Appendf(nil, `{"schema_version":%d,%s:%s}`, readSchemaVersion, key, value), nil
}

type groupsResult[G any] struct {
	SchemaVersion int `json:"schema_version"`
	Groups        []G `json:"groups"`
}

// newSpecAction builds a collection command: init <slug>, show and check
// <slug-or-source-path>, list without arguments, and any other action <slug>. The
// example uses the layout's example slug.
func newSpecAction(layout specfile.Layout, action, short string) *cobra.Command {
	kind := layout.Kind
	cmd := &cobra.Command{Use: action + " <slug>", Short: short, Args: cobra.ExactArgs(1),
		Example: "dotspec " + kind + " " + action + " " + layout.SlugExample + " --path ."}
	switch action {
	case "list":
		cmd.Use, cmd.Example, cmd.Args = "list", "dotspec "+kind+" list --path .", cobra.NoArgs
	case "show", "check":
		cmd.Use = action + " <slug-or-source-path>"
	}
	if action == "check" {
		cmd.Long = short + ".\n\n" + exitPolicy
	}
	return cmd
}

// validateSelector requires init to name a valid slug; other commands also accept a
// source path in the kind's collection.
func validateSelector(layout specfile.Layout, action, selector string) error {
	if err := layout.ValidateSlug(selector); err != nil && (action == "init" || !layout.Match(selector)) {
		return err
	}
	return nil
}

// writeCreated reports a created specification; guidance follows the created path in text.
func writeCreated(out io.Writer, format, target, guidance string) error {
	if format == "json" {
		return writeJSON(out, createdResult{1, target})
	}
	_, err := fmt.Fprintf(out, "Created %s\n%s\n", displayPath(target), guidance)
	return err
}

func emitFocusedCheck(out io.Writer, format, path, selector string, kind check.Type) error {
	return emitCheckReport(out, check.Focused(path, selector, kind, speckinds.All()...), format)
}

// writeShown prints authored Markdown, any details lines, and the associated system and source.
func writeShown(out io.Writer, markdown, details string, owner *system.System, source string) error {
	_, err := fmt.Fprintf(out, "%s\n\n%sSystem: %s\nSource: %s\n", markdown, details, systemLabel(owner), displayPath(source))
	return err
}

// writeGroupedList prints a heading, then each associated system followed by the lines
// writeMembers adds for that group.
func writeGroupedList[G any](out io.Writer, plural, heading string, groups []G, owner func(G) *system.System, writeMembers func(*strings.Builder, G)) error {
	if len(groups) == 0 {
		_, err := fmt.Fprintf(out, "No %s found in the selected subtree.\n", plural)
		return err
	}
	var output strings.Builder
	output.WriteString(heading + "\n")
	for _, group := range groups {
		fmt.Fprintf(&output, "\nSystem: %s\n", systemLabel(owner(group)))
		writeMembers(&output, group)
	}
	_, err := io.WriteString(out, output.String())
	return err
}

func writeListEntry(out *strings.Builder, name, path string) {
	fmt.Fprintf(out, "  %s  (%s)\n", name, displayPath(path))
}

// specSource is what a collection command prints about one authored specification.
type specSource struct {
	name, markdown, path string
	system               *system.System
}

// specKind describes a collection kind whose init, show, list and check commands need no
// options beyond those that flags adds.
type specKind[S, G any] struct {
	short      string
	plural     string // lowercase plural used in list text
	checkShort string
	key        string // JSON key of a shown specification
	guidance   string // text init prints after the created path
	layout     specfile.Layout
	checkType  func() check.Type
	init       func(path, slug string) (string, error)
	show       func(path, selector string) (*S, error)
	list       func(path string) ([]G, error)
	source     func(*S) specSource
	members    func(G) (*system.System, []*S)
	flags      func(action string, cmd *cobra.Command) // optional
	writeEntry func(*strings.Builder, *S)              // optional list line
}

func (k specKind[S, G]) command() *cobra.Command {
	root := group(k.layout.Kind, k.short)
	for _, action := range []struct{ name, short string }{
		{"init", "Create a " + k.layout.Kind + " specification without overwriting an existing file"},
		{"show", "Show authored Markdown and the associated system"},
		{"list", "Discover " + k.plural + ", grouped by associated system"},
		{"check", k.checkShort},
	} {
		cmd := newSpecAction(k.layout, action.name, action.short)
		action := action.name
		if k.flags != nil {
			k.flags(action, cmd)
		}
		cmd.RunE = func(cmd *cobra.Command, args []string) error {
			path, err := projectPath(cmd)
			if err != nil {
				return err
			}
			format, err := outputFormat(cmd)
			if err != nil {
				return err
			}
			out := cmd.OutOrStdout()
			if action == "list" {
				return k.writeList(out, path, format)
			}
			if err := validateSelector(k.layout, action, args[0]); err != nil {
				return err
			}
			return k.run(out, action, path, args[0], format)
		}
		addPath(cmd)
		addOutput(cmd, false)
		root.AddCommand(cmd)
	}
	return root
}

func (k specKind[S, G]) run(out io.Writer, action, path, selector, format string) error {
	switch action {
	case "init":
		target, err := k.init(path, selector)
		if err != nil {
			return err
		}
		return writeCreated(out, format, target, k.guidance)
	case "check":
		return emitFocusedCheck(out, format, path, selector, k.checkType())
	}
	spec, err := k.show(path, selector)
	if err != nil {
		return err
	}
	if format == "json" {
		return writeJSON(out, readResult{k.key, spec})
	}
	source := k.source(spec)
	return writeShown(out, source.markdown, "", source.system, source.path)
}

func (k specKind[S, G]) writeList(out io.Writer, path, format string) error {
	groups, err := k.list(path)
	if err != nil {
		return err
	}
	groups = nonnil(groups)
	if format == "json" {
		return writeJSON(out, groupsResult[G]{readSchemaVersion, groups})
	}
	owner := func(group G) *system.System {
		associated, _ := k.members(group)
		return associated
	}
	heading := strings.ToUpper(k.plural[:1]) + k.plural[1:]
	return writeGroupedList(out, k.plural, heading, groups, owner, func(output *strings.Builder, group G) {
		_, specs := k.members(group)
		for _, spec := range specs {
			if k.writeEntry != nil {
				k.writeEntry(output, spec)
				continue
			}
			source := k.source(spec)
			writeListEntry(output, source.name, source.path)
		}
	})
}
