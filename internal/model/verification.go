package model

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"strings"
)

// VerificationEvidence is an optional project-side interchange document. The
// selection describes the report being requested, independently of individual runs.
type VerificationEvidence struct {
	Selection VerificationContext  `json:"selection"`
	Results   []VerificationResult `json:"results"`
}
type VerificationContext struct {
	ProductRevision       string `json:"product_revision"`
	Conditions            string `json:"conditions"`
	SpecificationRevision string `json:"specification_revision,omitempty"`
}
type VerificationResult struct {
	VerificationContext
	Subject     string `json:"subject"`
	Procedure   string `json:"procedure"`
	Outcome     string `json:"outcome"`
	Observation string `json:"observation"`
}

type ProcedureAssessment struct {
	Subject      string               `json:"subject"`
	Procedures   []string             `json:"procedures"`
	Observations []VerificationResult `json:"observations,omitempty"`
}

const VerificationFile = "verification-results.json"

// ReadVerificationEvidence reads the optional observation source independently
// of the specification snapshot. A missing file is an empty evidence set.
func ReadVerificationEvidence(root string) (VerificationEvidence, error) {
	path := filepath.Join(root, VerificationFile)
	data, err := os.ReadFile(path)
	if errors.Is(err, os.ErrNotExist) {
		return VerificationEvidence{}, nil
	}
	if err != nil {
		return VerificationEvidence{}, err
	}
	evidence, err := ParseVerificationEvidence(data)
	if err != nil {
		return evidence, fmt.Errorf("%s: %w", path, err)
	}
	return evidence, nil
}

func (c VerificationContext) Matches(observed VerificationContext) bool {
	return c.ProductRevision != "" && c.ProductRevision == observed.ProductRevision &&
		c.Conditions != "" && c.Conditions == observed.Conditions &&
		(c.SpecificationRevision == "" || c.SpecificationRevision == observed.SpecificationRevision)
}

func ParseVerificationEvidence(input []byte) (VerificationEvidence, error) {
	var evidence VerificationEvidence
	d := json.NewDecoder(bytes.NewReader(input))
	d.DisallowUnknownFields()
	if err := d.Decode(&evidence); err != nil {
		return evidence, err
	}
	if err := d.Decode(new(any)); err != io.EOF {
		return evidence, fmt.Errorf("expected one verification evidence document")
	}
	if strings.TrimSpace(evidence.Selection.ProductRevision) == "" || strings.TrimSpace(evidence.Selection.Conditions) == "" || evidence.Results == nil {
		return evidence, fmt.Errorf("supply selection.product_revision, selection.conditions and a results array")
	}
	for i, r := range evidence.Results {
		if r.Outcome != "pass" && r.Outcome != "fail" && r.Outcome != "inconclusive" {
			return evidence, fmt.Errorf("results[%d]: expected pass, fail or inconclusive", i)
		}
		for _, field := range []string{r.Subject, r.Procedure, r.Observation} {
			if u, err := url.Parse(field); err != nil || strings.TrimSpace(field) == "" || u.Opaque != "" || strings.ContainsAny(field, "\r\n<>") {
				return evidence, fmt.Errorf("results[%d]: invalid subject, procedure or observation locator %q", i, field)
			}
		}
	}
	return evidence, nil
}

// Assessments indexes the original observations by their stated subject: a
// requirement's project-relative URI. Its procedures are the companions it
// lists with role verifies; a procedure's portable identity is its URI
// resolved relative to the project.
func (s *Snapshot) Assessments(root string, evidence VerificationEvidence) (map[string][]*ProcedureAssessment, error) {
	index := map[string]Record{}
	for _, collection := range [][]Record{s.Sources, s.References, s.Elements} {
		for _, r := range collection {
			index[String(r, "id")] = r
		}
	}
	bySubject := map[string]*ProcedureAssessment{}
	byRequirement := map[string][]*ProcedureAssessment{}
	for _, r := range s.Elements {
		source := index[String(r, "source_document")]
		if r["kind"] != "Requirement" || source["selected"] != true {
			continue
		}
		rel, err := filepath.Rel(root, filepath.FromSlash(String(source, "locator")))
		if err != nil {
			return nil, err
		}
		subject := (&url.URL{Path: filepath.ToSlash(rel)}).String()
		assessment := &ProcedureAssessment{Subject: subject}
		for _, id := range Strings(r["companions"]) {
			companion := index[id]
			if String(companion, "role") != "verifies" {
				continue
			}
			raw := String(index[String(companion, "target")], "authored_target")
			u, err := url.Parse(raw)
			if err != nil {
				return nil, err
			}
			if !u.IsAbs() && u.Host == "" && !strings.HasPrefix(u.Path, "/") {
				u.Path = path.Join(path.Dir(filepath.ToSlash(rel)), u.Path)
				u.RawPath = ""
			}
			assessment.Procedures = appendUnique(assessment.Procedures, u.String())
		}
		bySubject[subject] = assessment
		id := String(r, "id")
		byRequirement[id] = append(byRequirement[id], assessment)
	}
	roots := Strings(s.Scope["roots"])
	wholeProject := len(roots) == 1 && roots[0] == filepath.ToSlash(root) && s.Scope["inspection_complete"] == true
	for i, result := range evidence.Results {
		target := bySubject[result.Subject]
		if target == nil {
			if !wholeProject {
				continue
			}
			return nil, fmt.Errorf("%s: results[%d]: unknown requirement %q", VerificationFile, i, result.Subject)
		}
		target.Observations = append(target.Observations, result)
	}
	return byRequirement, nil
}
