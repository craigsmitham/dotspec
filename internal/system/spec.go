// Package system reads and discovers system specifications.
package system

import (
	"fmt"
	"path/filepath"

	"github.com/craigsmitham/dotspec/internal/diagnostic"
	"github.com/craigsmitham/dotspec/internal/specfile"
)

const Filename = "system.spec.md"

var Layout = specfile.Layout{Kind: "system", Filename: Filename}

type metadata struct {
	specfile.Metadata `yaml:",inline"`
}

const Template = `---
type: system
---

# New system

## Purpose

Describe why this system exists and whom it serves.

## Scope

Describe the system's boundary and what it includes.
`

type System struct {
	Name    string `json:"name"`
	Purpose string `json:"purpose"`
	Scope   string `json:"scope"`
	Path    string `json:"path"`
}

func Read(path string) (*System, error) {
	if filepath.Base(path) != Filename {
		return nil, diagnostic.Error("layout.system-location", fmt.Sprintf("expected a %s filename", Filename), diagnostic.Location{Path: path})
	}
	if _, err := Layout.Owner(path); err != nil {
		return nil, err
	}
	doc, err := specfile.Read(path, Layout.Kind, &metadata{}, "Purpose", "Scope")
	if err != nil {
		return nil, err
	}
	return fromDocument(doc), nil
}

func Parse(data []byte) (*System, error) {
	doc, err := specfile.Parse(data, Layout.Kind, &metadata{}, "Purpose", "Scope")
	if err != nil {
		return nil, err
	}
	return fromDocument(doc), nil
}

func fromDocument(doc *specfile.Document) *System {
	return &System{Name: doc.Name, Purpose: doc.Sections["Purpose"], Scope: doc.Sections["Scope"], Path: doc.Path}
}
