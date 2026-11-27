package check_test

import (
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"github.com/dotspec-md/dotspec/internal/check"
	"github.com/dotspec-md/dotspec/internal/check/testdata/quality"
	"github.com/dotspec-md/dotspec/internal/diagnostic"
	"github.com/dotspec-md/dotspec/internal/job"
	"github.com/dotspec-md/dotspec/internal/specfile"
	"github.com/dotspec-md/dotspec/internal/system"
	"github.com/dotspec-md/dotspec/internal/usecase"
)

const qualitySource = "---\ntype: quality-spec\nunit: milliseconds\n---\n# Response time\n\n## Criterion\n95% of requests finish in 200 milliseconds.\n\n## Evidence\nRetain this **authored prose**.\n"

func testScope(t *testing.T) string {
	t.Helper()
	dir, err := specfile.Directory(t.TempDir())
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
	write(t, filepath.Join(dir, "specs/product.system.spec.md"), "broken")
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
	write(t, filepath.Join(dir, "specs/product.system.spec.md"), validSystem)
	write(t, filepath.Join(dir, "specs/jobs/unrelated.job.spec.md"), validJob)
	write(t, filepath.Join(dir, "specs/use-cases/work.use-case.spec.md"), usecase.Template)
	report := check.Project(dir, system.CheckType(), job.CheckType(), usecase.CheckType(), kind)
	if report.ExitCode() != 0 || report.Summary.FilesChecked != 4 || len(report.Scope.Types) != 4 || report.Summary.ContextFilesChecked != 0 {
		t.Fatalf("mixed type composition failed: %+v", report)
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
	if report := check.Focused(dir, "Bad_Slug", kind); report.ExitCode() != 2 || report.Summary.FilesChecked != 0 || len(report.Diagnostics) != 1 || report.Diagnostics[0].Rule != "scan.scope" || !strings.Contains(report.Diagnostics[0].Message, "invalid quality slug") {
		t.Fatalf("invalid slug contract: %+v", report)
	}
	misplaced := filepath.Join(dir, "specs", "misplaced"+quality.Layout.Suffix)
	write(t, misplaced, qualitySource)
	for _, report := range []check.Report{check.Project(dir, kind), check.Focused(dir, "response-time", kind)} {
		if report.ExitCode() != 0 || !report.Complete {
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
	kind.Context = &check.ContextProvider{Resolve: func(_ *specfile.SourceSet, owner string) ([]string, error) {
		calls++
		if owner != filepath.Dir(first) {
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
	other := check.Type{Layout: specfile.Layout{Kind: "evidence", Collection: "quality", Suffix: ".evidence.spec.md"},
		Inspect: func(*specfile.SourceSet, string) ([]diagnostic.Diagnostic, error) { return nil, nil },
		Context: &check.ContextProvider{Resolve: func(*specfile.SourceSet, string) ([]string, error) {
			otherCalls++
			return []string{evidence}, nil
		}},
	}
	write(t, other.Layout.Path(dir, "x"), "evidence")
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

func TestInventoryRecordsEachSourceContextResolution(t *testing.T) {
	dir := testScope(t)
	kind := quality.CheckType()
	first, second := quality.Layout.Path(dir, "first"), quality.Layout.Path(dir, "second")
	broken := quality.Layout.Path(dir, "broken")
	write(t, first, qualitySource)
	write(t, second, qualitySource)
	write(t, broken, "broken frontmatter")
	evidence := filepath.Join(dir, "measurements.csv")
	failure := &os.PathError{Op: "open", Path: evidence, Err: os.ErrNotExist}
	kind.Context = &check.ContextProvider{
		Resolve: func(*specfile.SourceSet, string) ([]string, error) {
			t.Fatal("source-aware provider used the directory default")
			return nil, nil
		},
		ResolveSource: func(_ *specfile.SourceSet, _, source string) ([]string, error) {
			if source == second {
				return []string{evidence}, failure
			}
			return []string{evidence}, nil
		},
	}
	inventory := check.InspectSources(dir, kind)
	// Each source keeps exactly what its provider returned; a source check
	// could not read has no resolution.
	want := map[string]check.ContextResult{first: {Paths: []string{evidence}}, second: {Paths: []string{evidence}, Err: failure}}
	if !reflect.DeepEqual(inventory.Contexts, want) {
		t.Fatalf("recorded contexts: %+v", inventory.Contexts)
	}
}

func TestContextFailuresKeepFileIdentityAndCompletion(t *testing.T) {
	dir := testScope(t)
	write(t, quality.Layout.Path(dir, "response-time"), qualitySource)
	contextPath := filepath.Join(dir, "measurement_plan.csv")
	for _, tt := range []struct {
		name     string
		err      error
		complete bool
		code     int
	}{
		{"invalid", diagnostic.Error("measurement_plan.invalid", "Invalid measurement.", diagnostic.Location{Path: contextPath, Line: 2, Column: 1}), true, 1},
		{"missing", &os.PathError{Op: "open", Path: contextPath, Err: os.ErrNotExist}, false, 2},
	} {
		t.Run(tt.name, func(t *testing.T) {
			kind := quality.CheckType()
			kind.Context = &check.ContextProvider{Resolve: func(*specfile.SourceSet, string) ([]string, error) { return []string{contextPath}, tt.err }}
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
	kind.Context = &check.ContextProvider{Resolve: func(*specfile.SourceSet, string) ([]string, error) {
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
		Resolve: func(*specfile.SourceSet, string) ([]string, error) { resolveCalls++; return []string{path}, nil },
		Inspect: func(*specfile.SourceSet, string) ([]diagnostic.Diagnostic, error) {
			inspectCalls++
			return nil, inspectionError
		},
	}
	kind.InspectContext = func(_ *specfile.SourceSet, source string, paths []string) ([]diagnostic.Diagnostic, error) {
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

func TestOptionalSpecsDirectoryMayBeAbsent(t *testing.T) {
	dir := testScope(t)
	git(t, dir, "init", "-q")
	child := filepath.Join(dir, "component")
	write(t, filepath.Join(child, "specs/jobs/work.job.spec.md"), validJob)
	write(t, filepath.Join(dir, "specs"), "a file instead of the metadata directory")
	for _, report := range []check.Report{checkProject(child), checkJob(child, "work")} {
		if !report.Complete || report.ExitCode() != 0 || report.Summary.Errors != 0 || len(report.ContextFiles) != 0 || report.Summary.ContextFilesChecked != 0 {
			t.Fatalf("context directory counted as a source file: %+v", report)
		}
	}
}
