package cli

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/craigsmitham/dotspec/internal/businessrule"
	"github.com/craigsmitham/dotspec/internal/concept"
	"github.com/craigsmitham/dotspec/internal/job"
	"github.com/craigsmitham/dotspec/internal/qualitymodel"
	"github.com/craigsmitham/dotspec/internal/requirement"
	"github.com/craigsmitham/dotspec/internal/speccollection"
	"github.com/craigsmitham/dotspec/internal/specfile"
	"github.com/craigsmitham/dotspec/internal/system"
	"github.com/craigsmitham/dotspec/internal/usecase"
	"github.com/spf13/cobra"
)

// kind is one specification kind as the CLI presents it. new and list
// generate one subcommand per entry, in table order, so adding a kind adds an
// entry here and no command wiring.
type kind struct {
	layout specfile.Layout
	short  string // what the kind specifies, in domain terms
	plural string // lowercase plural for Shorts and list text
	// operand is new's usage operand: <name> for named subjects, <slug> for
	// collection kinds and <id>... for adopted quality attributes.
	operand  string
	about    string   // what new writes, for new's Long
	examples []string // new's example operands and flags
	guidance string   // what new prints after the created path
	flags    func(*cobra.Command)
	// create writes one specification from the working directory and returns
	// its path. adopt replaces it for kinds whose creation reports decisions.
	create  func(cmd *cobra.Command, dir, operand string) (string, error)
	adopt   func(cmd *cobra.Command, dir string, operands []string) error
	members func(dir string) ([]listEntry, error)
	// tree kinds print their structure instead of grouping by system.
	tree      func(entries []listEntry) string
	listShort string
}

// kinds lists every specification kind the CLI creates and lists.
func kinds() []kind {
	return []kind{
		{layout: system.Layout, short: "Describe a system's purpose, scope, and context", plural: "systems", operand: "<name>",
			about:    "The <name> becomes the heading, and the filename is derived from it. --part-of names the parent system, relative to the working directory; a project's root system comes from dotspec init.",
			examples: []string{"Billing -C internal --part-of ../specs/new-system.system.spec.md"},
			guidance: "Edit Description, Purpose, and Scope. Keep the filename in sync if you rename the system.",
			flags:    partOfFlag("parent system source file, relative to the working directory"),
			create: func(cmd *cobra.Command, dir, name string) (string, error) {
				parent := flagValue(cmd, "part-of")
				return system.InitSubsystem(dir, name, parent)
			},
			members: systemEntries, tree: systemTree, listShort: "List systems by directory"},
		{layout: system.ComponentLayout, short: "Describe a part of a system and its responsibilities", plural: "components", operand: "<name>",
			about:    "The <name> becomes the heading, and the filename is derived from it. --part-of names the parent system or component, relative to the working directory; Part of is never inferred from the file's location.",
			examples: []string{"Intake -C internal --part-of ../specs/new-system.system.spec.md"},
			flags:    partOfFlag("parent system or component source file, relative to the working directory"),
			create: func(cmd *cobra.Command, dir, name string) (string, error) {
				parent := flagValue(cmd, "part-of")
				return system.InitComponent(dir, name, parent)
			},
			members: componentEntries, tree: componentTree, listShort: "List components by Part of"},
		product(system.FeatureLayout, "Group stakeholder-visible behavior", "features"),
		product(system.CapabilityLayout, "Describe an ability and the parts that provide it", "capabilities"),
		collection(job.Layout, "Define jobs to be done (JTBDs)", "jobs",
			"Describe what the performer wants to accomplish and how they judge success.", job.Init,
			heading(job.Read, func(s *job.Job) string { return s.Name })),
		{layout: usecase.Layout, short: "Describe actor goals and interactions with a system", plural: "use cases", operand: "<slug>",
			examples: []string{usecase.Layout.SlugExample, "review-change --template expanded"},
			guidance: "Edit the actor goal and system boundary; add success criteria, scenarios, and extensions as they become known.",
			flags: func(cmd *cobra.Command) {
				stringFlag(cmd.Flags(), "template", "", "brief|expanded", "brief", "writing prompts: brief or expanded (same document format)")
			},
			create: func(cmd *cobra.Command, dir, slug string) (string, error) {
				template := flagValue(cmd, "template")
				return usecase.InitTemplate(dir, slug, template)
			},
			members: collected(usecase.Layout, "use cases", heading(usecase.Read, func(s *usecase.UseCase) string { return s.Name }))},
		collection(concept.Layout, "Describe purposeful product behavior through state and actions", "concepts",
			"Describe the purpose, state, actions, and operational principle as they become known.", concept.Init,
			heading(concept.Read, func(s *concept.Concept) string { return s.Name })),
		collection(businessrule.Layout, "Describe business constraints independently of their enforcement", "business rules",
			"Describe the business constraint, applicability, rationale, authority, and exceptions as they become known.", businessrule.Init,
			heading(businessrule.Read, func(s *businessrule.BusinessRule) string { return s.Name })),
		{layout: requirement.Layout, short: "Specify required results, rules, properties, and restrictions", plural: "requirements", operand: "<slug>",
			examples: []string{requirement.Layout.SlugExample},
			guidance: "Describe the scope, statement, any criterion, rationale, and verification as they become known.",
			flags: func(cmd *cobra.Command) {
				stringFlag(cmd.Flags(), "attribute", "", "<id>", "", "classify the requirement by this quality attribute of its governing system")
			},
			create: func(cmd *cobra.Command, dir, slug string) (string, error) {
				attribute := flagValue(cmd, "attribute")
				if cmd.Flags().Changed("attribute") && attribute == "" {
					return "", usageError(cmd, "--attribute requires a quality attribute ID")
				}
				return requirement.InitWithAttribute(dir, slug, attribute)
			},
			members: collected(requirement.Layout, "requirements", heading(requirement.Read, func(s *requirement.Requirement) string { return s.Name }))},
		{layout: system.AttributeLayout, short: "Adopt a quality attribute for a system, blank or from a quality model", plural: "quality attributes", operand: "<id>...",
			about:    "Each <id> becomes the attribute's ID and slug in the nearest enclosing system. Without --from, each attribute starts with a blank definition in specs/quality/<id>/. With --from, the definition is copied from the model's term and cited, and the term's missing parents are adopted too; parents the system already declares are reused. --dry-run previews the decisions and new file contents without writing.",
			examples: []string{"recoverability --from iso-25010-2023 --dry-run"},
			flags: func(cmd *cobra.Command) {
				stringFlag(cmd.Flags(), "from", "", "<model>", "", "copy from a bundled model, local model slug, or model source path")
				cmd.Flags().Bool("dry-run", false, "preview the decisions and new file contents without writing")
			},
			adopt:   adoptQualityAttributes,
			members: collected(system.AttributeLayout, "quality attributes", heading(system.ReadQualityAttribute, func(s *system.AttributeDocument) string { return s.Name }))},
		{layout: qualitymodel.Layout, short: "Define reusable quality attributes for systems to adopt", plural: "quality models", operand: "<slug>",
			about:    "With --from, the new model starts as a copy of a bundled model.",
			examples: []string{qualitymodel.Layout.SlugExample, "team-quality --from iso-25010-2023"},
			flags: func(cmd *cobra.Command) {
				stringFlag(cmd.Flags(), "from", "", "<model>", "", "start from a copy of this bundled model")
			},
			create: func(cmd *cobra.Command, dir, slug string) (string, error) {
				from := flagValue(cmd, "from")
				if cmd.Flags().Changed("from") && from == "" {
					return "", usageError(cmd, "--from requires a bundled model")
				}
				return qualitymodel.Init(dir, slug, from)
			},
			members: collected(qualitymodel.Layout, "quality models", heading(qualitymodel.Read, func(s *qualitymodel.QualityModel) string { return s.Name }))},
	}
}

// collection describes a kind created from its layout's template.
func collection(layout specfile.Layout, short, plural, guidance string, create func(path, slug string) (string, error), name func(string) (string, error)) kind {
	return kind{layout: layout, short: short, plural: plural, operand: "<slug>", examples: []string{layout.SlugExample}, guidance: guidance,
		create:  func(_ *cobra.Command, dir, slug string) (string, error) { return create(dir, slug) },
		members: collected(layout, plural, name)}
}

func product(layout specfile.Layout, short, plural string) kind {
	return kind{layout: layout, short: short, plural: plural, operand: "<slug>", examples: []string{layout.SlugExample},
		create: func(_ *cobra.Command, dir, slug string) (string, error) { return system.InitProduct(layout, dir, slug) },
		members: collected(layout, plural, func(path string) (string, error) {
			spec, err := system.ReadProductSpec(layout, path)
			if err != nil {
				return "", err
			}
			return spec.Name, nil
		})}
}

func partOfFlag(usage string) func(*cobra.Command) {
	return func(cmd *cobra.Command) {
		stringFlag(cmd.Flags(), "part-of", "", "<path>", "", usage)
		_ = cmd.MarkFlagRequired("part-of")
	}
}

// heading reads a specification's heading through its kind's reader.
func heading[T any](read func(string) (*T, error), name func(*T) string) func(string) (string, error) {
	return func(path string) (string, error) {
		spec, err := read(path)
		if err != nil {
			return "", err
		}
		return name(spec), nil
	}
}

// listEntry is one specification in a list result. Paths are absolute.
type listEntry struct {
	Kind   string     `json:"kind"`
	Slug   string     `json:"slug"`
	Name   string     `json:"name"`
	Path   string     `json:"path"`
	System *systemRef `json:"system"`
	PartOf string     `json:"part_of,omitempty"`
}

// systemRef names the system an entry belongs to: its associated system, a
// component's governing system or a system's enclosing system.
type systemRef struct {
	Name string `json:"name"`
	Path string `json:"path"`
}

func refer(s *system.System) *systemRef {
	if s == nil {
		return nil
	}
	return &systemRef{s.Name, s.Path}
}

func entry(layout specfile.Layout, name, path string, owner *systemRef) listEntry {
	slug := strings.TrimSuffix(filepath.Base(path), layout.Suffix)
	if slug == filepath.Base(path) || "."+slug == layout.Suffix {
		slug = ""
	}
	return listEntry{Kind: layout.Kind, Slug: slug, Name: name, Path: path, System: owner}
}

// collected lists a kind's files with their associated systems.
func collected(layout specfile.Layout, plural string, name func(string) (string, error)) func(string) ([]listEntry, error) {
	return func(dir string) ([]listEntry, error) {
		members, err := speccollection.Members(dir, layout, plural, name)
		if err != nil {
			return nil, err
		}
		entries := []listEntry{}
		for _, m := range members {
			entries = append(entries, entry(layout, m.Name, m.Path, refer(m.System)))
		}
		return entries, nil
	}
}

func systemEntries(dir string) ([]listEntry, error) {
	nodes, err := system.List(dir)
	if err != nil {
		return nil, err
	}
	entries := []listEntry{}
	var visit func([]*system.Node, *system.System)
	visit = func(nodes []*system.Node, parent *system.System) {
		for _, node := range nodes {
			entries = append(entries, entry(system.Layout, node.System.Name, node.System.Path, refer(parent)))
			visit(node.Children, node.System)
		}
	}
	visit(nodes, nil)
	return entries, nil
}

func componentEntries(dir string) ([]listEntry, error) {
	members, err := system.Components(dir)
	if err != nil {
		return nil, err
	}
	entries := []listEntry{}
	for _, m := range members {
		e := entry(system.ComponentLayout, m.Name, m.Path, refer(m.System))
		e.PartOf = m.Parent
		entries = append(entries, e)
	}
	return entries, nil
}

// systemTree nests each system under its enclosing system's directory.
func systemTree(entries []listEntry) string {
	children := map[string][]listEntry{}
	for _, e := range entries {
		parent := ""
		if e.System != nil {
			parent = e.System.Path
		}
		children[parent] = append(children[parent], e)
	}
	var out strings.Builder
	out.WriteString("Systems by directory\n")
	writeBranches(&out, children, "", "")
	return out.String()
}

// componentTree nests each component under the parent its Part of link names,
// beneath the governing system.
func componentTree(entries []listEntry) string {
	children := map[string][]listEntry{}
	var roots []*systemRef
	seen := map[string]bool{}
	for _, e := range entries {
		children[e.PartOf] = append(children[e.PartOf], e)
		if !seen[e.System.Path] {
			seen[e.System.Path] = true
			roots = append(roots, e.System)
		}
	}
	var out strings.Builder
	out.WriteString("Components by Part of\n")
	for _, root := range roots {
		fmt.Fprintf(&out, "\n%s  (%s)\n", root.Name, displayPath(root.Path))
		writeBranches(&out, children, root.Path, "")
	}
	return out.String()
}

func writeBranches(out *strings.Builder, children map[string][]listEntry, parent, prefix string) {
	nodes := children[parent]
	for i, node := range nodes {
		branch, next := "├── ", "│   "
		if i == len(nodes)-1 {
			branch, next = "└── ", "    "
		}
		fmt.Fprintf(out, "%s%s%s  (%s)\n", prefix, branch, node.Name, displayPath(node.Path))
		writeBranches(out, children, node.Path, prefix+next)
	}
}
