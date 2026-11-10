package cli

import (
	"fmt"
	"io"
	"strings"

	"github.com/craigsmitham/dotspec/internal/check"
	"github.com/craigsmitham/dotspec/internal/job"
	"github.com/craigsmitham/dotspec/internal/system"
	"github.com/spf13/cobra"
)

type jobOptions struct {
	path, format, step string
	outcomes           bool
}
type jobResult struct {
	SchemaVersion int      `json:"schema_version"`
	Job           *job.Job `json:"job"`
}
type jobsResult struct {
	SchemaVersion int          `json:"schema_version"`
	Groups        []*job.Group `json:"groups"`
}
type stepsResult struct {
	SchemaVersion int        `json:"schema_version"`
	Name          string     `json:"name"`
	Path          string     `json:"path"`
	Steps         []job.Step `json:"steps"`
}
type outcomesResult struct {
	SchemaVersion int           `json:"schema_version"`
	Name          string        `json:"name"`
	Path          string        `json:"path"`
	Overall       []job.Outcome `json:"overall"`
	Steps         []job.Step    `json:"steps"`
}

func newJobCommand() *cobra.Command {
	root := group("job", "Define and discover jobs to be done (JTBDs)")
	for _, action := range []struct{ name, short string }{
		{"init", "Create a job specification without overwriting an existing file"},
		{"show", "Show authored Markdown and the associated system"},
		{"list", "Discover jobs, grouped by associated system"},
		{"map", "Map authored step order (text, JSON, or Mermaid)"},
		{"check", "Check structure and suggest writing improvements"},
		{"step", "List the job's steps in authored order"},
		{"outcome", "List overall and step outcomes"},
	} {
		use := action.name + " <slug>"
		if action.name == "list" {
			use = "list"
		}
		cmd := &cobra.Command{Use: use, Short: action.short,
			Example: "dotspec job " + action.name + " develop-requirements --path ."}
		command := action.name
		if command == "step" || command == "outcome" {
			parent := group(command, action.short)
			cmd.Use = "list <slug>"
			command += " list"
			cmd.Example = "dotspec job " + command + " develop-requirements --path ."
			parent.AddCommand(cmd)
			root.AddCommand(parent)
		} else {
			root.AddCommand(cmd)
		}
		cmd.Args = cobra.ExactArgs(1)
		if command == "list" {
			cmd.Args = cobra.NoArgs
			cmd.Example = "dotspec job list --path ."
		}
		if command == "check" {
			cmd.Long = action.short + ".\n\n" + exitPolicy
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
			options := jobOptions{path: path, format: format}
			if command == "map" {
				options.outcomes, _ = cmd.Flags().GetBool("outcomes")
			}
			if command == "outcome list" {
				options.step, _ = cmd.Flags().GetString("step")
			}
			if command == "list" {
				return listJobs(path, format, cmd.OutOrStdout())
			}
			if err := job.ValidateSlug(args[0]); err != nil {
				return err
			}
			return runJob(command, args[0], options, cmd.OutOrStdout())
		}
		addPath(cmd)
		addOutput(cmd, command == "map")
		if command == "map" {
			cmd.Flags().Bool("outcomes", false, "include job and step outcomes")
		}
		if command == "outcome list" {
			cmd.Flags().String("step", "", "select a step by its explicit ID")
		}
	}
	return root
}

func runJob(command, slug string, options jobOptions, out io.Writer) error {
	if command == "init" {
		target, err := job.Init(options.path, slug)
		if err != nil {
			return err
		}
		if options.format == "json" {
			return writeJSON(out, createdResult{1, target})
		}
		_, err = fmt.Fprintf(out, "Created %s\nEdit the job statement and performer; add context, steps, and outcomes as they become known.\n", displayPath(target))
		return err
	}
	if command == "check" {
		return emitCheckReport(out, check.FocusedJob(options.path, slug), options.format)
	}
	spec, err := job.Show(options.path, slug)
	if err != nil {
		return err
	}
	switch command {
	case "show":
		result := jobResult{1, spec}
		if options.format == "json" {
			return writeJSON(out, result)
		}
		_, err = fmt.Fprintf(out, "%s\n\nSystem: %s\nSource: %s\n", result.Job.Markdown, systemLabel(result.Job.System), displayPath(result.Job.Path))
	case "step list":
		if options.format == "json" {
			return writeJSON(out, stepsResult{1, spec.Name, spec.Path, nonnil(spec.Steps)})
		}
		err = listSteps(spec, out)
	case "outcome list":
		result, selectionErr := selectOutcomes(spec, options.step)
		if selectionErr != nil {
			return selectionErr
		}
		if options.format == "json" {
			return writeJSON(out, result)
		}
		err = renderOutcomes(result, out)
	case "map":
		if options.format == "json" {
			steps := append([]job.Step{}, spec.Steps...)
			overall := []job.Outcome{}
			if options.outcomes {
				overall = nonnil(spec.Outcomes)
			} else {
				for i := range steps {
					steps[i].Outcomes = []job.Outcome{}
				}
			}
			return writeJSON(out, outcomesResult{1, spec.Name, spec.Path, overall, steps})
		}
		err = writeJobMap(spec, options, out)
	}
	return err
}

func nonnil[T any](values []T) []T {
	if values == nil {
		return []T{}
	}
	return values
}

func listJobs(path, format string, out io.Writer) error {
	groups, err := job.List(path)
	if err != nil {
		return err
	}
	result := jobsResult{1, nonnil(groups)}
	if format == "json" {
		return writeJSON(out, result)
	}
	if len(result.Groups) == 0 {
		_, err = fmt.Fprintln(out, "No jobs found in the selected subtree.")
		return err
	}
	var output strings.Builder
	output.WriteString("Jobs\n")
	for _, group := range result.Groups {
		fmt.Fprintf(&output, "\nSystem: %s\n", systemLabel(group.System))
		for _, spec := range group.Jobs {
			fmt.Fprintf(&output, "  %s  (%s)\n", spec.Name, displayPath(spec.Path))
		}
	}
	_, err = io.WriteString(out, output.String())
	return err
}

func systemLabel(spec *system.System) string {
	if spec == nil {
		return "None"
	}
	return fmt.Sprintf("%s (%s)", spec.Name, displayPath(spec.Path))
}
