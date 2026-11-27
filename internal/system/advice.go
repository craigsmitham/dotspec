package system

import (
	"github.com/dotspec-md/dotspec/internal/diagnostic"
	"github.com/dotspec-md/dotspec/internal/specfile"
)

// boundaryAdvice prompts for the support a reviewer needs to judge a named
// assumption or interface. Locations omit the path, which the reader adds.
func boundaryAdvice(doc *specfile.Outline) []diagnostic.Diagnostic {
	var findings []diagnostic.Diagnostic
	for _, section := range doc.Sections {
		switch section.Title {
		case "Assumptions":
			for _, assumption := range section.Children {
				findings = append(findings, assumptionAdvice(doc, section, assumption)...)
			}
		case "Required interfaces", "Provided interfaces":
			for _, item := range section.Children {
				if !hasContract(doc, item) {
					location := doc.Location(item.Start, section.Title+" / "+item.Title)
					findings = append(findings, diagnostic.Advice("interface.contract", "Link to the native contract that defines this interface under a Contracts heading.", location))
				}
			}
		}
	}
	return findings
}

func assumptionAdvice(doc *specfile.Outline, section, assumption *specfile.Section) []diagnostic.Diagnostic {
	prompts := []struct{ rule, title, message string }{
		{"assumption.basis", "Basis", "Explain under a Basis heading why it is reasonable to rely on this assumption."},
		{"assumption.applicability", "Applicability", "State under an Applicability heading where this assumption holds."},
	}
	var findings []diagnostic.Diagnostic
	for _, prompt := range prompts {
		location := doc.Location(assumption.Start, section.Title+" / "+assumption.Title)
		present := false
		for _, detail := range assumption.Children {
			if detail.Title == prompt.title {
				location = doc.Location(detail.Start, section.Title+" / "+assumption.Title+" / "+detail.Title)
				present = present || doc.Body(detail) != ""
			}
		}
		if !present {
			findings = append(findings, diagnostic.Advice(prompt.rule, prompt.message, location))
		}
	}
	return findings
}

func hasContract(doc *specfile.Outline, item *specfile.Section) bool {
	for _, detail := range item.Children {
		if detail.Title == "Contracts" && len(doc.LinksWithin(detail)) > 0 {
			return true
		}
	}
	return false
}

func withPath(findings []diagnostic.Diagnostic, path string) []diagnostic.Diagnostic {
	located := make([]diagnostic.Diagnostic, len(findings))
	for i, finding := range findings {
		finding.Location.Path = path
		located[i] = finding
	}
	return located
}
