package check_test

import (
	"errors"
	"testing"

	"github.com/craigsmitham/dotspec/internal/check"
	"github.com/craigsmitham/dotspec/internal/diagnostic"
	"github.com/craigsmitham/dotspec/internal/specfile"
)

type member struct{ references *specfile.ReferenceDocument }

func TestReadTypeAdvisesOnlyAfterReading(t *testing.T) {
	refused := errors.New("refused")
	read := func(_ *specfile.SourceSet, path string) (*member, error) {
		if path == "bad" {
			return nil, refused
		}
		return &member{references: &specfile.ReferenceDocument{}}, nil
	}
	project := func(m *member) *specfile.ReferenceDocument { return m.references }
	advised := 0
	kind := check.ReadType(specfile.Layout{Kind: "member"}, read, project, func(*member) []diagnostic.Diagnostic {
		advised++
		return []diagnostic.Diagnostic{diagnostic.Advice("member.advice", "Advise.", diagnostic.Location{})}
	})
	files := &specfile.SourceSet{}
	if findings, err := kind.Inspect(files, "bad"); !errors.Is(err, refused) || findings != nil || advised != 0 {
		t.Fatalf("failed read: %v %v, advised %d times", findings, err, advised)
	}
	if references, err := kind.References(files, "bad"); !errors.Is(err, refused) || references != nil {
		t.Fatalf("failed references: %v %v", references, err)
	}
	if findings, err := kind.Inspect(files, "good"); err != nil || len(findings) != 1 || advised != 1 {
		t.Fatalf("advice: %v %v", findings, err)
	}
	if references, err := kind.References(files, "good"); err != nil || references == nil {
		t.Fatalf("references: %v %v", references, err)
	}
	silent := check.ReadType(specfile.Layout{Kind: "member"}, read, project, nil)
	if findings, err := silent.Inspect(files, "good"); err != nil || findings != nil {
		t.Fatalf("without advice: %v %v", findings, err)
	}
	if kind.Context != nil {
		t.Fatal("ReadType must leave context to the caller")
	}
}
