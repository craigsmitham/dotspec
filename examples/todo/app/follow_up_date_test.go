package main

import "testing"

// Each case describes a date supplied by the person waiting for a reply.
func TestFollowUpDateCalendar(t *testing.T) {
	for _, day := range []string{"0001-01-01", "2000-02-29", "2024-02-29", "10000-02-29"} {
		if !validFollowUpDate(day) {
			t.Errorf("expected to accept %q", day)
		}
	}
	for _, day := range []string{"0000-01-01", "1900-02-29", "10100-02-29", "2024-04-31"} {
		if validFollowUpDate(day) {
			t.Errorf("expected to reject %q", day)
		}
	}
}

func TestFollowUpDateSpelling(t *testing.T) {
	for _, day := range []string{"", "24-02-29", "2024-2-29", "2024-02-9", "tomorrow"} {
		if validFollowUpDate(day) {
			t.Errorf("expected to reject %q", day)
		}
	}
}
