package model

// OpenQuestions returns the authored question records belonging to an element.
func OpenQuestions(r Record) []Record {
	return records(r["open_questions"])
}
