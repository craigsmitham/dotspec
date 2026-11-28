package specdoc

import (
	"fmt"
	"slices"
	"strings"

	"github.com/dotspec-md/dotspec/internal/model"
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

func (p *planner) renderAssessment(w *writer, s *subject, record model.Record, heading string) {
	a := p.assessments[id(record)]
	if a == nil {
		return
	}
	if len(a.Procedures) == 0 && len(a.Observations) == 0 && (record["criterion"] == nil || p.evidence.Selection.ProductRevision == "") {
		return
	}
	w.text(heading)
	if len(a.Procedures) > 0 {
		for _, key := range model.Strings(record["verification"]) {
			w.text("- " + p.catalogLink(w, s, p.elements[key]))
		}
		if a.UsedAsProcedure {
			w.text("Procedure: " + p.catalogLink(w, s, record))
		}
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
