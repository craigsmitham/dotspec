package cli

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"strings"

	"github.com/craigsmitham/dotspec/internal/job"
	"github.com/craigsmitham/dotspec/internal/system"
)

const jobUsage = `Usage: dotspec job <command> [slug] [--path DIR]

Commands:
  init <slug>          Create a job template
  show <slug>          Show all authored sections and the associated system
  list                 Discover jobs, grouped by associated system
  step list <slug>     List the job's steps in document order
  outcome list <slug>  List job and step outcomes; optionally --step ID
  map <slug>           Generate a job map; --format text|mermaid, --outcomes
  check <slug>         Check structure and suggest writing improvements

Slugs use lowercase letters or digits separated by single hyphens.
--path defaults to the current directory and can precede or follow the slug.
Use --help for command help.
`

type jobOptions struct {
	path, format, step string
	outcomes           bool
}

func runJobArgs(args []string, out io.Writer) error {
	if len(args) == 0 || isHelp(args[0]) {
		_, err := fmt.Fprint(out, jobUsage)
		return err
	}
	command := args[0]
	args = args[1:]
	if command == "step" || command == "outcome" {
		if len(args) == 0 || isHelp(args[0]) {
			_, err := fmt.Fprintf(out, "Usage: dotspec job %s list <slug> [--path DIR]", command)
			if err != nil {
				return err
			}
			if command == "outcome" {
				_, err = fmt.Fprint(out, " [--step ID]")
			}
			if err != nil {
				return err
			}
			_, err = fmt.Fprintln(out)
			return err
		}
		if args[0] != "list" {
			return fmt.Errorf("unknown job %s command %q; use dotspec job --help", command, args[0])
		}
		command += " list"
		args = args[1:]
	}
	switch command {
	case "init", "show", "list", "map", "check", "step list", "outcome list":
	default:
		return fmt.Errorf("unknown job command %q; use dotspec job --help", command)
	}
	flags := flag.NewFlagSet("dotspec job "+command, flag.ContinueOnError)
	flags.SetOutput(io.Discard)
	options := jobOptions{format: "text"}
	flags.StringVar(&options.path, "path", ".", "directory containing the job")
	if command == "map" {
		flags.StringVar(&options.format, "format", "text", "text or mermaid")
		flags.BoolVar(&options.outcomes, "outcomes", false, "include job and step outcomes")
	}
	if command == "outcome list" {
		flags.StringVar(&options.step, "step", "", "select a step by its ID")
	}
	operands, err := parseFlagArgs(flags, args)
	if errors.Is(err, flag.ErrHelp) {
		slug := " <slug>"
		if command == "list" {
			slug = ""
		}
		if _, err := fmt.Fprintf(out, "Usage: dotspec job %s%s [options]\n\n", command, slug); err != nil {
			return err
		}
		flags.SetOutput(out)
		flags.PrintDefaults()
		return nil
	}
	if err != nil {
		return err
	}
	if strings.TrimSpace(options.path) == "" {
		return fmt.Errorf("--path must not be empty")
	}
	if options.format != "text" && options.format != "mermaid" {
		return fmt.Errorf("unknown map format %q; use text or mermaid", options.format)
	}
	if command == "list" {
		if len(operands) != 0 {
			return fmt.Errorf("job list does not take a slug; select a directory with --path DIR")
		}
		return listJobs(options.path, out)
	}
	if len(operands) != 1 {
		return fmt.Errorf("dotspec job %s requires exactly one slug; use dotspec job --help", command)
	}
	slug := operands[0]
	if command == "init" {
		target, err := job.Init(options.path, slug)
		if err != nil {
			return err
		}
		_, err = fmt.Fprintf(out, "Created %s\nEdit the job statement and performer; add context, steps, and outcomes as they become known.\n", displayPath(target))
		return err
	}
	spec, err := job.Show(options.path, slug)
	if err != nil {
		return err
	}
	switch command {
	case "show":
		_, err = fmt.Fprintf(out, "%s\n\nSystem: %s\nSource: %s\n", spec.Markdown, systemLabel(spec.System), displayPath(spec.Path))
	case "step list":
		err = listSteps(spec, out)
	case "outcome list":
		err = listOutcomes(spec, options.step, out)
	case "map":
		err = writeJobMap(spec, options, out)
	case "check":
		suggestions := job.Check(spec)
		if _, err = fmt.Fprintf(out, "No structural errors in %s.\n", displayPath(spec.Path)); err != nil {
			return err
		}
		if len(suggestions) == 0 {
			_, err = fmt.Fprintln(out, "No writing suggestions from the current checks.")
		} else {
			if _, err = fmt.Fprintf(out, "%d writing suggestions (advisory):\n", len(suggestions)); err != nil {
				return err
			}
			for _, suggestion := range suggestions {
				if _, err = fmt.Fprintf(out, "- %s: %s\n", suggestion.Location, suggestion.Message); err != nil {
					return err
				}
			}
		}
	}
	return err
}

func listJobs(path string, out io.Writer) error {
	groups, err := job.List(path)
	if err != nil {
		return err
	}
	if len(groups) == 0 {
		_, err := fmt.Fprintln(out, "No jobs found in the selected subtree.")
		return err
	}
	var output strings.Builder
	output.WriteString("Jobs\n")
	for _, group := range groups {
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
