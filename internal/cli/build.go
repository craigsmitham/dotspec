package cli

import (
	"errors"
	"fmt"
	"io"
	"path/filepath"
	"strings"

	"github.com/dotspec-md/dotspec/internal/check"
	"github.com/dotspec-md/dotspec/internal/diagnostic"
	"github.com/dotspec-md/dotspec/internal/model"
	"github.com/dotspec-md/dotspec/internal/specdoc"
	"github.com/dotspec-md/dotspec/internal/specfile"
	"github.com/dotspec-md/dotspec/internal/speckinds"
	"github.com/spf13/cobra"
)

func newBuildCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use: "build", Short: "Build the generated specification",
		Long:    "Build the generated specification for the enclosing project declared with dotspec.yaml: SPEC.md and the specification/ tree of pages that let a reader navigate the project's specifications as one. Nothing is written while a source has error findings; the findings are reported instead. " + findingsExit + "\n\nLearn more: dotspec docs docs/reference/inspection.md#check-the-generated-requirements-specification",
		Example: "dotspec build", Args: operands(0, 0),
		RunE: func(cmd *cobra.Command, _ []string) error {
			path, err := projectScope.dir("")
			if err != nil {
				return err
			}
			project, err := specfile.FindProject(path)
			if err != nil {
				return err
			}
			if project == nil {
				return fmt.Errorf("no dotspec.yaml declares a project here; run dotspec init <system-name> to start one")
			}
			if !project.Build.Enabled {
				return fmt.Errorf("build.enabled is false in %s; remove that setting to build the generated specification", displayPath(project.Manifest()))
			}
			plan, report := planSpec(project.Root)
			if plan == nil {
				// The lead line states the build's outcome, so the report
				// does not repeat the generated specification's state.
				report.GeneratedSpec = ""
				if _, err := fmt.Fprintf(cmd.OutOrStdout(), "SPEC.md not written: %s.\n\n", count(report.Summary.Errors, "source error", "source errors")); err != nil {
					return &exitError{code: 2, err: err}
				}
				return emitCheckReport(cmd.OutOrStdout(), report, false)
			}
			changes, err := specdoc.Write(project.Root, plan)
			if err != nil {
				return err
			}
			var out strings.Builder
			for _, path := range changes.Written {
				fmt.Fprintf(&out, "Wrote %s\n", path)
			}
			for _, path := range changes.Removed {
				fmt.Fprintf(&out, "Removed %s\n", path)
			}
			fmt.Fprintf(&out, "Built %s (%d written, %d removed).\n", displayPath(filepath.Join(project.Root, "SPEC.md")), len(changes.Written), len(changes.Removed))
			_, err = io.WriteString(cmd.OutOrStdout(), out.String())
			return err
		},
	}
	return cmd
}

func inspectProject(path string) check.Report {
	dir, err := specfile.Directory(path)
	if err != nil {
		return check.Project(path, speckinds.All()...)
	}
	project, err := specfile.FindProject(dir)
	if err != nil || project == nil || !project.Build.Enabled || dir != project.Root {
		return check.Project(path, speckinds.All()...)
	}
	plan, report := planSpec(dir)
	if plan == nil {
		return report
	}
	findings, err := specdoc.Check(dir, plan)
	if err != nil {
		report.Complete = false
		report.GeneratedSpec = "incomplete"
		report = report.WithFindings(*diagnostic.Error("spec.inspect", err.Error(), diagnostic.Location{Path: dir}))
	} else {
		report.GeneratedSpec = "checked"
	}
	return report.WithFindings(findings...)
}

func planSpec(root string) (*specdoc.Plan, check.Report) {
	inventory := check.InspectSources(root, speckinds.All()...)
	report := inventory.Report
	report.GeneratedSpec = "blocked-by-source"
	if report.ExitCode() != 0 {
		return nil, report
	}

	exported := model.ExportInspection(inventory)
	if exported.ExitCode != 0 {
		report.Complete = exported.ExitCode != 2
		return nil, report.WithFindings(modelErrors(exported.Snapshot, root)...)
	}
	plan, err := specdoc.Generate(root, exported.Snapshot)
	if err != nil {
		var finding *diagnostic.Diagnostic
		if errors.As(err, &finding) {
			return nil, report.WithFindings(*finding)
		}
		return nil, report.WithFindings(*diagnostic.Error("spec.structure", err.Error(), diagnostic.Location{Path: root}))
	}
	return plan, report
}

// Source inspection findings have already been reported. If mapping fails,
// preserve its source path and position in the existing check report.
func modelErrors(s *model.Snapshot, root string) []diagnostic.Diagnostic {
	paths := make(map[string]string, len(s.Sources))
	for _, source := range s.Sources {
		paths[model.String(source, "id")] = model.String(source, "locator")
	}
	var findings []diagnostic.Diagnostic
	for _, r := range s.Diagnostics {
		if model.String(r, "severity") != "error" {
			continue
		}
		path := paths[model.String(r, "source")]
		if path == "" {
			path = root
		}
		findings = append(findings, *diagnostic.Error(model.String(r, "code"), model.String(r, "message"), diagnostic.Location{Path: path, Line: model.Number(r["line"]), Column: model.Number(r["column"]), Element: model.String(r, "element")}))
	}
	return findings
}
