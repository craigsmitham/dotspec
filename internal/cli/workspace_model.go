package cli

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"

	"github.com/craigsmitham/dotspec/internal/model"
	modelschema "github.com/craigsmitham/dotspec/model"
)

// Workspace exports retain separate snapshots, including their acquisition scope.
// Validate each result without merging record identities or reading the workspace.
func validateExportJSON(data []byte) ([]model.Finding, error) {
	var probe map[string]json.RawMessage
	if json.Unmarshal(data, &probe) != nil || probe["projects"] == nil {
		return model.ValidateJSON(data)
	}
	fail := func(message string) ([]model.Finding, error) {
		return []model.Finding{{Code: "workspace.shape", Path: "/", Message: message}}, nil
	}
	var envelope workspaceResult
	decoder := json.NewDecoder(bytes.NewReader(data))
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&envelope); err != nil {
		return fail(err.Error())
	}
	var extra any
	if decoder.Decode(&extra) != io.EOF {
		return fail("expected one workspace export")
	}
	if envelope.SchemaVersion != 1 || envelope.Workspace == "" || len(envelope.Projects) == 0 {
		return fail("workspace export needs schema_version 1, workspace and projects")
	}
	findings := []model.Finding{}
	names, roots := map[string]bool{}, map[string]bool{}
	for i, p := range envelope.Projects {
		prefix := fmt.Sprintf("/projects/%d", i)
		if p.Name == "" || p.Root == "" || names[p.Name] || roots[p.Root] || p.ExitCode < 0 || p.ExitCode > 2 {
			return fail("project identities must be nonempty and unique, with an exit_code from 0 to 2")
		}
		names[p.Name], roots[p.Root] = true, true
		if len(p.Result) == 0 {
			findings = append(findings, model.Finding{Code: "workspace.result", Path: prefix, Message: "project has no exported snapshot"})
			continue
		}
		current, err := model.ValidateJSON(p.Result)
		if err != nil {
			return nil, err
		}
		for _, f := range current {
			f.Path = prefix + "/result" + f.Path
			findings = append(findings, f)
		}
		if len(current) == 0 {
			var snapshot model.Snapshot
			if err := json.Unmarshal(p.Result, &snapshot); err != nil {
				return nil, err
			}
			project := model.Object(snapshot.Scope["project"])
			if model.String(project, "name") != p.Name || model.String(project, "root") != p.Root || model.String(project, "workspace") != envelope.Workspace {
				findings = append(findings, model.Finding{Code: "workspace.scope", Path: prefix + "/result/scope/project", Message: "snapshot project must match its workspace entry"})
			}
		}
	}
	return findings, nil
}

func workspaceExportSchema() (map[string]any, error) {
	var snapshot map[string]any
	if err := json.Unmarshal(modelschema.SchemaJSON, &snapshot); err != nil {
		return nil, err
	}
	str := map[string]any{"type": "string", "minLength": 1}
	return map[string]any{
		"$schema": snapshot["$schema"], "$id": "https://github.com/craigsmitham/dotspec/workspace-export.schema.json",
		"$defs": snapshot["$defs"], "type": "object", "additionalProperties": false,
		"required": []string{"schema_version", "workspace", "projects"},
		"properties": map[string]any{
			"schema_version": map[string]any{"const": 1}, "workspace": str,
			"projects": map[string]any{"type": "array", "minItems": 1, "items": map[string]any{
				"type": "object", "additionalProperties": false, "required": []string{"name", "root", "exit_code", "result"},
				"properties": map[string]any{"name": str, "root": str, "exit_code": map[string]any{"type": "integer", "minimum": 0, "maximum": 2}, "error": str, "result": map[string]any{"$ref": "#/$defs/Snapshot"}},
			}},
		},
	}, nil
}
