package cli

import (
	"fmt"
	"io"
	"strings"

	"github.com/craigsmitham/dotspec/internal/job"
	"github.com/craigsmitham/dotspec/internal/system"
	"github.com/spf13/cobra"
)

type jobOptions struct {
	path, format, step, preset string
	outcomes                   bool
}
type outcomesResult struct {
	SchemaVersion int                  `json:"schema_version"`
	Name          string               `json:"name"`
	Path          string               `json:"path"`
	DeclaredByJob []job.DesiredOutcome `json:"declared_by_job"`
	Steps         []job.Step           `json:"steps"`
}

func newJobCommand() *cobra.Command {
	root := group("job", "Define and discover jobs to be done (JTBDs)")
	for _, action := range []struct{ name, short string }{
		{"init", "Create a job specification without overwriting an existing file"},
		{"show", "Show authored Markdown and the associated system"},
		{"list", "Discover jobs, grouped by associated system"},
		{"map", "Map authored step order (text, JSON, or Mermaid)"},
		{"check", "Check structure and suggest writing improvements"},
		{"outcome", "List outcomes by declaration owner"},
	} {
		cmd := newSpecAction(job.Layout, action.name, action.short)
		command := action.name
		if command == "outcome" {
			parent := group(command, action.short)
			cmd.Use = "list <slug>"
			command += " list"
			cmd.Example = "dotspec job " + command + " " + job.Layout.SlugExample + " --path ."
			parent.AddCommand(cmd)
			root.AddCommand(parent)
		} else {
			root.AddCommand(cmd)
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
			if command == "init" || command == "check" {
				options.preset, _ = cmd.Flags().GetString("preset")
				if err := job.ValidatePreset(options.preset); err != nil {
					return err
				}
			}
			if command == "map" {
				options.outcomes, _ = cmd.Flags().GetBool("outcomes")
			}
			if command == "outcome list" {
				options.step, _ = cmd.Flags().GetString("step")
			}
			if command == "list" {
				return listJobs(path, format, cmd.OutOrStdout())
			}
			if err := validateSelector(job.Layout, action.name, args[0]); err != nil {
				return err
			}
			return runJob(command, args[0], options, cmd.OutOrStdout())
		}
		if command == "init" || command == "check" {
			cmd.Flags().String("preset", "general", "JTBD authoring guidance: general, odi, or job-story")
		}
		addPath(cmd)
		addOutput(cmd, command == "map")
		if command == "map" {
			cmd.Flags().Bool("outcomes", false, "include job and step outcomes")
		}
		if command == "outcome list" {
			cmd.Flags().String("step", "", "select the declaring step by its explicit ID")
		}
	}
	return root
}

func runJob(command, slug string, options jobOptions, out io.Writer) error {
	if command == "init" {
		target, err := job.InitWithPreset(options.path, slug, options.preset)
		if err != nil {
			return err
		}
		return writeCreated(out, options.format, target, "Describe what the performer wants to accomplish and how they judge success.\nWrite a job from a concrete situation: "+authoringPage("job"))
	}
	if command == "check" {
		return emitFocusedCheck(out, options.format, options.path, slug, job.CheckTypeWithPreset(options.preset))
	}
	spec, err := job.Show(options.path, slug)
	if err != nil {
		return err
	}
	switch command {
	case "show":
		if options.format == "json" {
			return writeJSON(out, readResult{"job", spec})
		}
		err = writeShown(out, spec.Markdown, "", spec.System, spec.Path)
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
			overall := []job.DesiredOutcome{}
			if options.outcomes {
				overall = nonnil(spec.Outcomes)
			} else {
				for i := range steps {
					steps[i].Outcomes = []job.DesiredOutcome{}
				}
			}
			return writeJSON(out, outcomesResult{2, spec.Name, spec.Path, overall, steps})
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
	groups = nonnil(groups)
	if format == "json" {
		return writeJSON(out, groupsResult[*job.Group]{readSchemaVersion, groups})
	}
	owner := func(group *job.Group) *system.System { return group.System }
	return writeGroupedList(out, "jobs", "Jobs", groups, owner, func(output *strings.Builder, group *job.Group) {
		for _, spec := range group.Jobs {
			writeListEntry(output, spec.Name, spec.Path)
		}
	})
}

func systemLabel(spec *system.System) string {
	if spec == nil {
		return "None"
	}
	return fmt.Sprintf("%s (%s)", spec.Name, displayPath(spec.Path))
}
