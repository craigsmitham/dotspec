package cli

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"slices"
	"strings"

	"github.com/dotspec-md/dotspec/internal/check"
	"github.com/dotspec-md/dotspec/internal/specfile"
	"github.com/dotspec-md/dotspec/internal/speckinds"
	"github.com/spf13/cobra"
)

func newSourcesCommand() *cobra.Command {
	cmd := &cobra.Command{Use: "sources [<path>]", Short: "List specifications, verification procedures, reference documents, system context, and referenced files",
		Long:    "List the files check reads: specifications, reference documents, enclosing system context and directly referenced files, each with the references that reach it. A named specification source shows what it defines. A verification procedure identifies the specification it assesses. With no <path>, sources covers the enclosing project; a directory narrows it to that subtree. A file selects its incoming references from the project inventory.\n\nA listed reference means its target exists and can be addressed, not that the target is authoritative. " + findingsExit + "\n\nLearn more: dotspec docs docs/reference/inspection.md#discover-specification-sources",
		Example: "dotspec sources\ndotspec sources internal --json\ndotspec sources tests/acceptance.py", Args: operands(0, 1),
		RunE: func(cmd *cobra.Command, args []string) error {
			r, err := sourceInventory(optionalPath(args))
			if err != nil {
				return err
			}
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

// A reverse lookup scans the project that can contain incoming declarations.
// For shared files outside a project, use the project of the invoking directory.
func sourceInventory(operand string) (check.SourceReport, error) {
	if info, err := os.Lstat(operand); err == nil && info.Mode().IsRegular() {
		dir, err := specfile.Directory(filepath.Dir(operand))
		if err != nil {
			return check.SourceReport{}, err
		}
		file := filepath.Join(dir, filepath.Base(operand))
		project, err := specfile.FindProject(dir)
		if err != nil {
			return check.SourceReport{}, err
		}
		if project == nil {
			project, err = specfile.FindProject(".")
			if err != nil {
				return check.SourceReport{}, err
			}
		}
		if project != nil {
			dir = project.Root
		} else if repo, err := specfile.RepositoryRoot(dir); err != nil {
			return check.SourceReport{}, err
		} else if repo != "" {
			dir = repo
		}
		result := check.Sources(dir, speckinds.All()...)
		result.SelectedFile = file
		result.Files = slices.DeleteFunc(result.Files, func(f check.SourceFile) bool { return f.Path != file })
		return result, nil
	}
	dir, err := projectScope.dir(operand)
	if err != nil {
		return check.SourceReport{}, err
	}
	return check.Sources(dir, speckinds.All()...), nil
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
				declaration := slices.IndexFunc(file.Declarations, func(d check.SourceDeclaration) bool { return d.Location == location })
				if declaration >= 0 {
					d := file.Declarations[declaration]
					description := "Verification procedure"
					if d.Defines != "" {
						description = "Defines " + d.Defines
					} else if d.VerifiedSpecification != "" {
						description += " for " + d.VerifiedSpecification
					}
					fmt.Fprintf(&text, "  %s: %s · %s (%s)\n", description, d.Name, diagnosticPath(location), d.Source)
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
