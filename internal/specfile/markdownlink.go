package specfile

import "strings"

// LinkLabelEnd returns the index of the bracket that closes the link or image
// label opened at raw[start], or -1 when raw[start] is not '[' or the label is
// unterminated. Following CommonMark link text, backslash escapes, matched
// bracket pairs and backtick code spans bind more tightly than the closing
// bracket, so [`]`](x.md) closes after the code span.
func LinkLabelEnd(raw []byte, start int) int {
	if start < 0 || start >= len(raw) || raw[start] != '[' {
		return -1
	}
	depth := 0
	for i := start; i < len(raw); i++ {
		switch raw[i] {
		case '\\':
			i++
		case '`':
			i = codeSpanEnd(raw, i)
		case '[':
			depth++
		case ']':
			depth--
			if depth == 0 {
				return i
			}
		}
	}
	return -1
}

// codeSpanEnd returns the last index of the code span opened by the backtick
// run at raw[i], or of the run itself when no run of equal length closes it.
func codeSpanEnd(raw []byte, i int) int {
	open := runEnd(raw, i)
	for j := open; j < len(raw); {
		if raw[j] != '`' {
			j++
			continue
		}
		close := runEnd(raw, j)
		if close-j == open-i {
			return close - 1
		}
		j = close
	}
	return open - 1
}

func runEnd(raw []byte, i int) int {
	for i < len(raw) && raw[i] == '`' {
		i++
	}
	return i
}

// The formatters take a destination or title as written in the source,
// without its delimiters, which is how goldmark records it whichever form the
// author used. A computed path is written the same way. Backslash escapes are
// kept as written, since every delimited form decodes them alike.

// FormatDestination writes a link destination in its pointy-bracket form,
// which accepts spaces and parentheses without further escaping. Unescaped
// angle brackets and line endings cannot appear in that form, so they are
// percent-encoded as a URL encodes them. A trailing literal backslash is
// doubled so that it cannot escape the closing bracket.
func FormatDestination(dest string) string {
	return "<" + escapeDelimited(dest, destinationEscapes) + ">"
}

var destinationEscapes = map[byte]string{'<': "%3C", '>': "%3E", '\n': "%0A", '\r': "%0D"}

// FormatTitle writes a link title as a double-quoted title. Unescaped double
// quotes are escaped, and a trailing literal backslash is doubled so that it
// cannot escape the closing quote.
func FormatTitle(title string) string {
	return `"` + escapeDelimited(title, titleEscapes) + `"`
}

var titleEscapes = map[byte]string{'"': `\"`}

// escapeDelimited replaces the unescaped bytes in replace and keeps backslash
// escapes intact. A literal backslash is doubled where it would otherwise
// escape what follows it in the output: a replacement or the closing delimiter.
func escapeDelimited(raw string, replace map[byte]string) string {
	var out strings.Builder
	for i := 0; i < len(raw); i++ {
		c := raw[i]
		if c == '\\' {
			next := i + 1
			switch {
			case next < len(raw) && strings.IndexByte(asciiPunctuation, raw[next]) >= 0:
				out.WriteString(raw[i : next+1])
				i = next
			case next == len(raw) || replace[raw[next]] != "":
				out.WriteString(`\\`)
			default:
				out.WriteByte(c)
			}
			continue
		}
		if r, ok := replace[c]; ok {
			out.WriteString(r)
		} else {
			out.WriteByte(c)
		}
	}
	return out.String()
}

// asciiPunctuation lists the characters a CommonMark backslash can escape.
const asciiPunctuation = "!\"#$%&'()*+,-./:;<=>?@[\\]^_`{|}~"
