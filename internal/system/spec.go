// Package system reads and discovers system specifications.
package system

import (
	"fmt"
	"path/filepath"

	"github.com/craigsmitham/dotspec/internal/specfile"
)

const Filename = "system.spec.md"

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
	Name    string
	Purpose string
	Scope   string
	Path    string
}

func Read(path string) (*System, error) {
	if filepath.Base(path) != Filename {
		return nil, fmt.Errorf("%s: expected a %s filename", path, Filename)
	}
	if _, err := specfile.Owner(path, ""); err != nil {
		return nil, err
	}
	doc, err := specfile.Read(path, "system", "Purpose", "Scope")
	if err != nil {
		return nil, err
	}
	return fromDocument(doc), nil
}

func Parse(data []byte) (*System, error) {
	doc, err := specfile.Parse(data, "system", "Purpose", "Scope")
	if err != nil {
		return nil, err
	}
	return fromDocument(doc), nil
}

func fromDocument(doc *specfile.Document) *System {
	return &System{Name: doc.Name, Purpose: doc.Sections["Purpose"], Scope: doc.Sections["Scope"], Path: doc.Path}
}
