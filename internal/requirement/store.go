package requirement

import (
	"fmt"
	"strings"

	"github.com/craigsmitham/dotspec/internal/specfile"
	"github.com/craigsmitham/dotspec/internal/system"
)

// InitWithAttribute creates a requirement, classified by attribute when one is given.
func InitWithAttribute(path, slug, attribute string) (string, error) {
	if err := Layout.ValidateSlug(slug); err != nil {
		return "", err
	}
	if attribute == "" {
		return Layout.Create(path, slug, Template)
	}
	dir, err := specfile.Directory(path)
	if err != nil {
		return "", err
	}
	chain, err := system.SubjectChain(dir)
	if err != nil {
		return "", err
	}
	if len(chain) == 0 {
		return "", fmt.Errorf("no governing system for quality attribute %q in %s", attribute, dir)
	}
	sys, err := system.Read(chain[len(chain)-1])
	if err != nil {
		return "", err
	}
	if sys.FindQualityAttribute(attribute) == nil {
		return "", fmt.Errorf("unknown quality attribute %q in %s", attribute, sys.Path)
	}
	template := strings.Replace(Template, "type: requirement-spec\n", fmt.Sprintf("type: requirement-spec\nquality_attributes: [%q]\n", attribute), 1)
	return Layout.Create(dir, slug, template)
}
