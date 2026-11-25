package specdoc

import (
	"fmt"
	"slices"
	"strings"

	"github.com/craigsmitham/dotspec/internal/model"
)

func (p *planner) readEvidence() error {
	evidence, err := model.ReadVerificationEvidence(p.root)
	if err != nil {
		return err
	}
	p.evidence = evidence
	p.assessments, err = p.snapshot.Assessments(p.root, p.evidence)
	return err
}

func (p *planner) renderVerification(w *writer, s *subject) {
	for _, a := range p.assessments[id(s.r)] {
		criterion := strings.HasSuffix(a.Subject, "#criterion")
		if len(a.Procedures) == 0 && len(a.Observations) == 0 && (!criterion || p.evidence.Selection.ProductRevision == "") {
			continue
		}
		title := "Requirement verification"
		if criterion {
			title = "Criterion verification"
		}
		w.heading(title)
		if len(a.Procedures) > 0 {
			w.text("Procedure linked.")
		} else {
			w.text("No procedure linked.")
		}
		selection := p.evidence.Selection
		matched := false
		var observations []string
		for _, r := range a.Observations {
			current := selection.Matches(r.VerificationContext) && slices.Contains(a.Procedures, r.Procedure)
			matched = matched || current
			scope := "Other recorded observation"
			if current {
				scope = "Matching observation"
			}
			dest := p.documentLink(w.file, model.VerificationFile, r.Observation)

			line := fmt.Sprintf("- %s · %s · %s · %s · %s", scope, inlineText(r.Outcome), inlineText(r.ProductRevision), inlineText(strings.Join(strings.Fields(r.Conditions), " ")), link("Observation", dest))
			line += " · " + inlineText(r.Procedure)
			if r.SpecificationRevision != "" {
				line += " · specification " + inlineText(r.SpecificationRevision)
			}
			observations = append(observations, line)
		}
		if selection.ProductRevision != "" {
			caption := "Selected revision " + inlineText(selection.ProductRevision) + " · " + inlineText(strings.Join(strings.Fields(selection.Conditions), " "))
			if selection.SpecificationRevision != "" {
				caption += " · specification " + inlineText(selection.SpecificationRevision)
			}
			w.text(caption)
		}
		if !matched {
			w.text("No matching result recorded.")
		}
		w.text(strings.Join(observations, "\n"))
	}
}
