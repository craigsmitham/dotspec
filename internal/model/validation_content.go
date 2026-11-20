package model

import (
	"encoding/json"
	"fmt"
	"math/big"
	"sort"
	"strconv"
	"unicode/utf8"
)

func natural(value any) (int, bool) {
	var text string
	switch value := value.(type) {
	case json.Number:
		text = string(value)
	case int:
		return value, value >= 0
	case float64:
		text = strconv.FormatFloat(value, 'f', -1, 64)
	default:
		return 0, false
	}
	if n, err := strconv.ParseUint(text, 10, strconv.IntSize-1); err == nil {
		return int(n), true
	}
	number, ok := new(big.Rat).SetString(text)
	if !ok || !number.IsInt() || !number.Num().IsInt64() || number.Sign() < 0 {
		return 0, false
	}
	n := number.Num().Int64()
	if uint64(n) > uint64(^uint(0)>>1) {
		return 0, false
	}
	return int(n), true
}

func (v *graphValidator) checkValue(object Record, class, path string) {
	switch class {
	case "SourceDocument":
		_, hasText := object["markdown"]
		if hasText != (String(object, "read_status") == "read" && String(object, "media_type") == "text/markdown") {
			v.add("source.markdown", path+"/markdown", "a successfully read Markdown source requires its text")
		}
	case "ByteSpan":
		start, validStart := natural(object["start"])
		end, validEnd := natural(object["end"])
		source := v.records[String(object, "source")].value
		markdown, hasText := source["markdown"].(string)
		if !validStart || !validEnd || !hasText || start > end || end > len(markdown) || start < len(markdown) && !utf8.RuneStart(markdown[start]) || end < len(markdown) && !utf8.RuneStart(markdown[end]) {
			v.add("source.span", path, "span must identify an ordered interval within acquired source bytes")
		}
	case "RichContent":
		last := 0
		owner, ok := v.containingRecord(path)
		for i, span := range records(object["spans"]) {
			if validationNumber(span["start"]) < last {
				v.add("content.order", fmt.Sprintf("%s/spans/%d", path, i), "field ranges must be ordered and disjoint")
			}
			if ok && String(span, "source") != v.sourceOf(owner.value) {
				v.add("content.ownership", path, "field ranges must belong to their owner's source")
			}
			last = validationNumber(span["end"])
		}
	case "Resolution":
		if String(object, "status") != "resolved" && object["target"] != nil {
			v.add("resolution.state", path+"/target", "only a resolved occurrence can claim a target")
		}
	case "SourceFragment", "Diagnostic", "SourceLocation":
		if span := Object(object["span"]); span != nil && String(object, "source") != "" && String(span, "source") != String(object, "source") {
			v.add("source.ownership", path+"/span", "source span disagrees with its location's source")
		}

	}
}

func (v *graphValidator) sortedRecords() []locatedRecord {
	if v.orderedRecords != nil {
		return v.orderedRecords
	}
	result := make([]locatedRecord, 0, len(v.records))
	for _, record := range v.records {
		result = append(result, record)
	}
	sort.Slice(result, func(i, j int) bool { return result[i].path < result[j].path })
	v.orderedRecords = result
	return result
}

func (v *graphValidator) checkCycles(code, kind string, edges func(Record) []string) {
	state := map[string]int{}
	var visit func(string)
	visit = func(id string) {
		r := v.records[id]
		if state[id] == 1 {
			v.add(code, r.path, "relationship hierarchy contains a cycle")
			return
		}
		if state[id] == 2 || !v.isA(String(r.value, "kind"), kind) {
			return
		}
		state[id] = 1
		for _, target := range edges(r.value) {
			visit(target)
		}
		state[id] = 2
	}
	for _, record := range v.sortedRecords() {
		visit(String(record.value, "id"))
	}
}

func validationNumber(value any) int { number, _ := natural(value); return number }
