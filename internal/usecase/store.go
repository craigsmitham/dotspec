package usecase

import "fmt"

// InitTemplate selects writing prompts, not a separate format or maturity state.
func InitTemplate(path, slug, template string) (string, error) {
	content := Template
	switch template {
	case "brief":
	case "expanded":
		content = ExpandedTemplate
	default:
		return "", fmt.Errorf("unknown use-case template %q; use brief or expanded", template)
	}
	return Layout.Create(path, slug, content)
}
