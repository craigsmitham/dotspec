package check_test

import (
	"errors"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"github.com/craigsmitham/dotspec/internal/check"
	"github.com/craigsmitham/dotspec/internal/check/testdata/quality"
	"github.com/craigsmitham/dotspec/internal/diagnostic"
	"github.com/craigsmitham/dotspec/internal/job"
	"github.com/craigsmitham/dotspec/internal/specfile"
	"github.com/craigsmitham/dotspec/internal/system"
	"github.com/craigsmitham/dotspec/internal/usecase"
)

const qualitySource = "---\ntype: quality\nid: QR-1\nunit: milliseconds\n---\n# Response time\n\n## Criterion\n95% of requests finish in 200 milliseconds. See [requirement](#QR-1).\n\n## Evidence\nRetain this **authored prose** and [criterion](#criterion).\n"

func testScope(t *testing.T) string {
	t.Helper()
	dir, err := specfile.Scope(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	return dir
}

func TestIndependentTypeIntegration(t *testing.T) {
	dir := testScope(t)
	source, err := quality.Layout.Create(dir, "response-time", qualitySource)
	if err != nil {
		t.Fatal(err)
	}
	if source != filepath.Join(dir, "specs/quality/response-time.quality.spec.md") {
		t.Fatalf("unexpected created path: %s", source)
	}
	if _, err := quality.Layout.Create(dir, "response-time", "replacement"); err == nil {
		t.Fatal("creation overwrote an existing requirement")
	}
	// An unrelated malformed system must not become implicit context.
	write(t, filepath.Join(dir, "specs/system.spec.md"), "broken")
	write(t, filepath.Join(dir, "specs/jobs/unrelated.job.spec.md"), "broken")
	kind := quality.CheckType()
	project, focused := check.Project(dir, kind), check.Focused(dir, "response-time", kind)
	for _, report := range []check.Report{project, focused} {
		if !report.Complete || report.ExitCode() != 0 || report.Summary.FilesChecked != 1 || report.Summary.ContextFilesChecked != 0 || len(report.Diagnostics) != 0 {
			t.Fatalf("type acquired unrelated dependencies: %+v", report)
		}
		if !reflect.DeepEqual(report.Scope.Types, []string{"quality"}) {
			t.Fatalf("wrong registered scope: %+v", report.Scope)
		}
	}
	data, err := os.ReadFile(source)
	if err != nil || string(data) != qualitySource {
		t.Fatalf("authored document changed: %q (%v)", data, err)
	}
	// Adding a fourth type requires only composition, not another engine branch.
	write(t, filepath.Join(dir, "specs/system.spec.md"), validSystem)
	write(t, filepath.Join(dir, "specs/jobs/unrelated.job.spec.md"), validJob)
	write(t, filepath.Join(dir, "specs/use-cases/work.use-case.spec.md"), usecase.Template)
	report := check.Project(dir, system.CheckType(), job.CheckType(), usecase.CheckType(), kind)
	if report.ExitCode() != 0 || report.Summary.FilesChecked != 4 || len(report.Scope.Types) != 4 || report.Summary.ContextFilesChecked != 0 {
		t.Fatalf("mixed type composition failed: %+v", report)
	}
}

func TestTypeOwnsMetadataSchemaAndIdentity(t *testing.T) {
	requirement, err := quality.Parse([]byte(qualitySource))
	if err != nil {
		t.Fatal(err)
	}
	if requirement.ID != "QR-1" || requirement.Unit != "milliseconds" || !strings.Contains(requirement.Markdown, "Retain this **authored prose**") || requirement.Location.Line != 6 {
		t.Fatalf("metadata, source, or locations lost: %+v", requirement)
	}
	for _, tt := range []struct {
		name, source, rule string
		line               int
	}{
		{"unknown field", strings.Replace(qualitySource, "unit:", "unknown:", 1), "metadata.invalid", 4},
		{"wrong scalar", strings.Replace(qualitySource, "unit: milliseconds", "unit: [milliseconds]", 1), "metadata.invalid", 4},
		{"wrong kind", strings.Replace(qualitySource, "type: quality", "type: another", 1), "metadata.type", 2},
		{"invalid ID", strings.Replace(qualitySource, "id: QR-1", "id: QR 1", 1), "identity.invalid", 3},
		{"unresolved ID", strings.Replace(qualitySource, "(#QR-1)", "(#QR-2)", 1), "reference.unresolved", 9},
	} {
		t.Run(tt.name, func(t *testing.T) {
			for _, source := range []string{tt.source, "\ufeff" + strings.ReplaceAll(tt.source, "\n", "\r\n")} {
				_, err := quality.Parse([]byte(source))
				var finding *diagnostic.Diagnostic
				if !errors.As(err, &finding) || finding.Rule != tt.rule || finding.Location.Line != tt.line {
					t.Fatalf("unexpected diagnostic: %+v (%v)", finding, err)
				}
			}
		})
	}
}

func TestIndependentTypeLayoutAndAdvice(t *testing.T) {
	dir := testScope(t)
	kind := quality.CheckType()
	write(t, quality.Layout.Path(dir, "response-time"), strings.Replace(qualitySource, "unit: milliseconds\n", "", 1))
	project, focused := check.Project(dir, kind), check.Focused(dir, "response-time", kind)
	if project.ExitCode() != 0 || project.Summary.Advice != 1 || !reflect.DeepEqual(project.Diagnostics, focused.Diagnostics) {
		t.Fatalf("advice contract: %+v; %+v", project, focused)
	}
	if report := check.Focused(dir, "missing", kind); report.ExitCode() != 2 || report.Complete || report.Summary.FilesChecked != 1 {
		t.Fatalf("missing file contract: %+v", report)
	}
	misplaced := filepath.Join(dir, "specs", "misplaced"+quality.Layout.Suffix)
	write(t, misplaced, qualitySource)
	for _, report := range []check.Report{check.Project(dir, kind), check.Focused(dir, "response-time", kind)} {
		if report.ExitCode() != 1 || !report.Complete || report.Diagnostics[0].Rule != "layout.quality-location" || report.Diagnostics[0].Location.Path != misplaced {
			t.Fatalf("layout inferred another type's identity: %+v", report)
		}
	}
}

func TestContextProvidersAreExplicitAndShared(t *testing.T) {
	dir := testScope(t)
	kind := quality.CheckType()
	first := quality.Layout.Path(dir, "first")
	write(t, first, qualitySource)
	write(t, quality.Layout.Path(dir, "second"), qualitySource)
	evidence := filepath.Join(dir, "measurements.csv")
	write(t, evidence, "latency\n150\n")
	calls := 0
	kind.Context = &check.ContextProvider{Resolve: func(owner string) ([]string, error) {
		calls++
		if owner != dir {
			t.Fatalf("wrong context owner: %s", owner)
		}
		// A provider may consult an already-discovered file as well as another
		// source format. Repeated paths still count only once.
		return []string{first, evidence, evidence}, nil
	}}
	report := check.Project(dir, kind)
	if calls != 1 || report.ExitCode() != 0 || report.Summary.FilesChecked != 2 || !reflect.DeepEqual(report.ContextFiles, []string{evidence}) {
		t.Fatalf("context sharing/accounting failed: calls=%d %+v", calls, report)
	}
	calls = 0
	report = check.Focused(dir, "first", kind)
	if calls != 1 || !reflect.DeepEqual(report.ContextFiles, []string{evidence}) || report.Summary.ContextFilesChecked != 1 {
		t.Fatalf("focused context differs: calls=%d %+v", calls, report)
	}
	// Distinct providers for the same owner both run, while shared paths are
	// still deduplicated. Sharing is by provider identity, not just directory.
	otherCalls := 0
	other := check.Type{Layout: specfile.Layout{Kind: "evidence", Filename: "evidence.spec.md"},
		Inspect: func(string) ([]diagnostic.Diagnostic, error) { return nil, nil },
		Context: &check.ContextProvider{Resolve: func(string) ([]string, error) {
			otherCalls++
			return []string{evidence}, nil
		}},
	}
	write(t, other.Layout.Path(dir, ""), "evidence")
	calls = 0
	report = check.Project(dir, kind, other)
	if calls != 1 || otherCalls != 1 || report.Summary.ContextFilesChecked != 1 {
		t.Fatalf("distinct providers merged: %d %d %+v", calls, otherCalls, report)
	}
	other.Context = kind.Context
	calls = 0
	report = check.Project(dir, kind, other)
	if calls != 1 || report.Summary.ContextFilesChecked != 1 {
		t.Fatalf("shared provider repeated across types: %d %+v", calls, report)
	}
}

func TestContextFailuresKeepFileIdentityAndCompletion(t *testing.T) {
	dir := testScope(t)
	write(t, quality.Layout.Path(dir, "response-time"), qualitySource)
	contextPath := filepath.Join(dir, "measurement.csv")
	for _, tt := range []struct {
		name     string
		err      error
		complete bool
		code     int
	}{
		{"invalid", diagnostic.Error("measurement.invalid", "Invalid measurement.", diagnostic.Location{Path: contextPath, Line: 2, Column: 1}), true, 1},
		{"missing", &os.PathError{Op: "open", Path: contextPath, Err: os.ErrNotExist}, false, 2},
	} {
		t.Run(tt.name, func(t *testing.T) {
			kind := quality.CheckType()
			kind.Context = &check.ContextProvider{Resolve: func(string) ([]string, error) { return []string{contextPath}, tt.err }}
			for _, report := range []check.Report{check.Project(dir, kind), check.Focused(dir, "response-time", kind)} {
				if report.Complete != tt.complete || report.ExitCode() != tt.code || report.Summary.Errors != 1 || report.Summary.ContextFilesChecked != 1 || report.Diagnostics[0].Location.Path != contextPath {
					t.Fatalf("context failure contract: %+v", report)
				}
			}
		})
	}
}

func TestMalformedTypeSkipsDependentContext(t *testing.T) {
	dir := testScope(t)
	write(t, quality.Layout.Path(dir, "broken"), "broken frontmatter")
	kind := quality.CheckType()
	kind.Context = &check.ContextProvider{Resolve: func(string) ([]string, error) {
		t.Fatal("resolved context after structural failure")
		return nil, nil
	}}
	for _, report := range []check.Report{check.Project(dir, kind), check.Focused(dir, "broken", kind)} {
		if !report.Complete || report.ExitCode() != 1 || report.Summary.ContextFilesChecked != 0 || report.Summary.Advice != 0 {
			t.Fatalf("dependent checks were not skipped: %+v", report)
		}
	}
}

func TestContextInspectionAndDependentChecks(t *testing.T) {
	dir := testScope(t)
	kind := quality.CheckType()
	write(t, quality.Layout.Path(dir, "first"), qualitySource)
	write(t, quality.Layout.Path(dir, "second"), qualitySource)
	path := filepath.Join(dir, "model.txt")
	resolveCalls, inspectCalls, dependentCalls := 0, 0, 0
	var inspectionError error
	kind.Context = &check.ContextProvider{
		Resolve: func(string) ([]string, error) { resolveCalls++; return []string{path}, nil },
		Inspect: func(string) ([]diagnostic.Diagnostic, error) { inspectCalls++; return nil, inspectionError },
	}
	kind.InspectContext = func(source string, paths []string) ([]diagnostic.Diagnostic, error) {
		dependentCalls++
		if !reflect.DeepEqual(paths, []string{path}) {
			t.Fatalf("resolved paths: %v", paths)
		}
		return []diagnostic.Diagnostic{diagnostic.Advice("test.context", "Resolved.", diagnostic.Location{Path: source})}, nil
	}
	report := check.Project(dir, kind)
	if resolveCalls != 1 || inspectCalls != 1 || dependentCalls != 2 || report.Summary.Advice != 2 {
		t.Fatalf("context calls: %d %d %d; %+v", resolveCalls, inspectCalls, dependentCalls, report)
	}
	resolveCalls, inspectCalls, dependentCalls = 0, 0, 0
	inspectionError = &os.PathError{Op: "open", Path: path, Err: os.ErrNotExist}
	report = check.Project(dir, kind)
	if report.ExitCode() != 2 || inspectCalls != 1 || dependentCalls != 0 || report.Summary.Errors != 1 {
		t.Fatalf("failed context leaked to dependents: %d %d %+v", inspectCalls, dependentCalls, report)
	}
}

func TestContextDirectoryFailureIsNotCountedAsAFile(t *testing.T) {
	dir := testScope(t)
	git(t, dir, "init", "-q")
	child := filepath.Join(dir, "component")
	write(t, filepath.Join(child, "specs/jobs/work.job.spec.md"), validJob)
	write(t, filepath.Join(dir, "specs"), "a file instead of the metadata directory")
	for _, report := range []check.Report{checkProject(child), checkJob(child, "work")} {
		if !report.Complete || report.ExitCode() != 1 || report.Summary.Errors != 1 || report.Diagnostics[0].Rule != "layout.directory" || len(report.ContextFiles) != 0 || report.Summary.ContextFilesChecked != 0 {
			t.Fatalf("context directory counted as a source file: %+v", report)
		}
	}
}
