package cli

import (
	"fmt"
	"io"
	"slices"
	"strings"

	"github.com/dotspec-md/dotspec/internal/check"
	"github.com/dotspec-md/dotspec/internal/speckinds"
	"github.com/spf13/cobra"
)

func newSourcesCommand() *cobra.Command {
	cmd := &cobra.Command{Use: "sources [<path>]", Short: "List specifications, their companions, reference documents, system context, and referenced files",
		Long:    "List the files check reads: specifications, reference documents, enclosing system context and directly referenced files, each with the references that reach it. A file a specification lists as a companion shows the role that specification states for it. With no <path>, sources covers the enclosing project; a directory narrows it to that subtree.\n\nA listed reference means its target exists and can be addressed, not that the target is authoritative; a companion's role is the author's declaration, not a verified relationship. " + findingsExit + "\n\nLearn more: dotspec docs docs/reference/inspection.md#discover-specification-sources",
		Example: "dotspec sources\ndotspec sources internal --json", Args: operands(0, 1),
		RunE: func(cmd *cobra.Command, args []string) error {
			path, err := projectScope.dir(optionalPath(args))
			if err != nil {
				return err
			}
			r := check.Sources(path, speckinds.All()...)
			if err := writeSources(cmd.OutOrStdout(), r, jsonOutput(cmd)); err != nil {
				return &exitError{code: 2, err: err}
			}
			if code := r.ExitCode(); code != 0 {
				return &exitError{code: code}
			}
			return nil
		}}
	addJSON(cmd)
	return cmd
}

const outsideProjects = "Outside any project"

func writeSources(out io.Writer, r check.SourceReport, asJSON bool) error {
	if asJSON {
		return writeJSON(out, r)
	}
	var text strings.Builder
	// Files are grouped by project, in the order each project first appears,
	// so the project is printed once per group.
	var groups []string
	members := make(map[string][]check.SourceFile)
	for _, file := range r.Files {
		group := outsideProjects
		if file.Project != nil {
			group = "Project: " + file.Project.Name
		}
		if _, seen := members[group]; !seen {
			groups = append(groups, group)
		}
		members[group] = append(members[group], file)
	}
	// Without any project, a group heading would say nothing.
	headed := len(groups) > 1 || len(groups) == 1 && groups[0] != outsideProjects
	for _, group := range groups {
		if headed {
			fmt.Fprintf(&text, "%s\n\n", group)
		}
		for _, file := range members[group] {
			fmt.Fprintf(&text, "%s\n", displayPath(file.Path))
			for _, location := range file.ReferencedBy {
				// A frontmatter companion entry is printed with its stated role.
				companion := slices.IndexFunc(file.CompanionOf, func(c check.SourceCompanion) bool { return c.Location == location })
				if companion >= 0 {
					fmt.Fprintf(&text, "  Companion that %s %s (%s)\n", file.CompanionOf[companion].Role, diagnosticPath(location), file.CompanionOf[companion].Target)
					continue
				}
				fmt.Fprintf(&text, "  Referenced by %s", diagnosticPath(location))
				if location.Element != "" {
					fmt.Fprintf(&text, " (%s)", location.Element)
				}
				text.WriteByte('\n')
			}
			text.WriteByte('\n')
		}
	}
	if r.Complete && len(r.Files) == 0 {
		text.WriteString("No source files found.\n\n")
	}
	for _, d := range r.Diagnostics {
		fmt.Fprintf(&text, "%s\n  %s [%s] %s\n", diagnosticPath(d.Location), d.Severity, d.Rule, d.Message)
		for _, related := range d.Related {
			fmt.Fprintf(&text, "  Related: %s\n", diagnosticPath(related))
		}
		text.WriteByte('\n')
	}
	completeness := "scan complete"
	if !r.Complete {
		completeness = "scan incomplete"
	}
	fmt.Fprintf(&text, "Listed %s in %s (%s).\n", count(len(r.Files), "source file", "source files"), scopeName(r.Scope), completeness)
	_, err := io.WriteString(out, text.String())
	return err
}
