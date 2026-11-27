package job

import "testing"

func TestJTBDPresetsDoNotImposeODIOnOtherJobs(t *testing.T) {
	source := "---\ntype: job-spec\n---\n# Travel\n## Job\nArrive confidently.\n## Performer\nA commuter.\n## Outcomes\n### Feel confident about keeping commitments\n\nFeel confident about keeping commitments.\n"
	j, err := Parse([]byte(source))
	if err != nil {
		t.Fatal(err)
	}
	if got := CheckWithPreset(j, "general"); len(got) != 0 {
		t.Fatal("general check imposed a method", got)
	}
	odi := CheckWithPreset(j, "odi")
	found := false
	for _, d := range odi {
		found = found || d.Rule == "job.outcome-measure"
	}
	if !found {
		t.Fatal("ODI guidance missing", odi)
	}
	if len(CheckWithPreset(j, "job-story")) == 0 {
		t.Fatal("preset guidance missing")
	}
	for _, preset := range []string{"general", "odi", "job-story"} {
		template, err := TemplateFor(preset)
		if err != nil {
			t.Fatal(err)
		}
		parsed, err := Parse([]byte(template))
		if err != nil {
			t.Fatalf("%s: %v", preset, err)
		}
		if parsed.Statement == "" {
			t.Fatal(preset, parsed)
		}
		if preset == "odi" && (len(parsed.Outcomes) != 1 || parsed.Outcomes[0].AppliesTo == "") {
			t.Fatal("outcome scope missing", parsed)
		}
	}
	for _, preset := range []string{"unknown", ""} {
		if _, err := TemplateFor(preset); err == nil {
			t.Fatalf("preset %q accepted", preset)
		}
	}
}

func TestNamedOutcomesPreserveFullStatementAndIdentity(t *testing.T) {
	source := "---\ntype: job-spec\n---\n# Job\n## Steps\n### Plan\n**ID:** plan\n### Monitor\n**ID:** monitor\n## Outcomes\n### Confidence\n**ID:** confidence\nFeel confident about arrival.\n#### Applies to\n[Plan](#plan) and [Monitor](#monitor)\n#### Supporting content\nA heading without an ID remains supporting prose.\n"
	j, err := Parse([]byte(source))
	if err != nil {
		t.Fatal(err)
	}
	if len(j.Outcomes) != 1 || j.Outcomes[0].Name != "Confidence" || j.Outcomes[0].Statement != "Feel confident about arrival." {
		t.Fatal(j.Outcomes)
	}
	if !j.References.Anchors["confidence"] {
		t.Fatal("named outcome ID not referenceable")
	}
	_, err = Parse([]byte(source + "### Duplicate\n**ID:** confidence\nAnother result.\n"))
	if err == nil {
		t.Fatal("duplicate outcome ID accepted")
	}
}
