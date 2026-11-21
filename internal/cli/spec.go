package cli

import (
	"errors"
	"fmt"
	"io"
	"path/filepath"
	"strings"

	"github.com/craigsmitham/dotspec/internal/check"
	"github.com/craigsmitham/dotspec/internal/diagnostic"
	"github.com/craigsmitham/dotspec/internal/model"
	"github.com/craigsmitham/dotspec/internal/specdoc"
	"github.com/craigsmitham/dotspec/internal/specfile"
	"github.com/spf13/cobra"
)

func newSpecCommand() *cobra.Command {
	cmd := group("spec", "Generate the project's readable requirements specification")
	build := &cobra.Command{
		Use: "build", Short: "Build SPEC.md and the supporting Markdown tree",
		Long:    "Build the requirements specification for one whole project. Reports are enabled by default in projects declared with dotspec.yaml. Validate the candidate before replacing generator-owned output. Run dotspec check to verify the resulting documents.",
		Example: "dotspec spec build", Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			path, err := projectPath(cmd)
			if err != nil {
				return err
			}
			project, err := specfile.FindProject(path)
			if err != nil {
				return err
			}
			if project == nil {
				return fmt.Errorf("spec build requires a project declared in dotspec.yaml")
			}
			if !project.Report.Enabled {
				return fmt.Errorf("report.enabled is false in %s. Remove that setting to build the specification", project.Manifest())
			}
			if path != project.Root {
				return fmt.Errorf("spec build requires the whole project; select --path %s or --project %s", project.Root, project.Name)
			}
			plan, report := planSpec(project.Root)
			if plan == nil {
				return emitCheckReport(cmd.OutOrStdout(), report, "text")
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
	addPath(build)
	cmd.AddCommand(build)
	return cmd
}

func inspectProject(path string) check.Report {
	dir, err := specfile.Directory(path)
	if err != nil {
		return check.Project(path, supportedTypes()...)
	}
	project, err := specfile.FindProject(dir)
	if err != nil || project == nil || !project.Report.Enabled || dir != project.Root {
		return check.Project(path, supportedTypes()...)
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
	types := supportedTypes()
	inventory := check.InspectSources(root, types...)
	report := inventory.Report
	report.GeneratedSpec = "blocked-by-source"
	if report.ExitCode() != 0 {
		return nil, report
	}

	exported := model.ExportInspection(inventory, types...)
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
