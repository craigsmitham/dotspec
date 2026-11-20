// Package assessment exposes the assessment vocabulary generated from the metamodel.
package assessment

type Purpose string

type Definition struct {
	Purpose  Purpose
	Label    string
	Question string
}

func (d Definition) Heading() string { return d.Label + " plan" }

func FromHeading(heading string) (Definition, bool) {
	for _, d := range Definitions {
		if d.Heading() == heading {
			return d, true
		}
	}
	return Definition{}, false
}

func (p Purpose) Definition() Definition {
	for _, d := range Definitions {
		if p == d.Purpose {
			return d
		}
	}
	return Definition{}
}
