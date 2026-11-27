package specfile

import "testing"

func TestLinkLabelEnd(t *testing.T) {
	for _, c := range []struct {
		raw   string
		start int
		want  int
	}{
		{"[label](x.md)", 0, 6},
		{"![alt](x.png)", 1, 5},
		{"[`]`](x.md)", 0, 4},
		{"[``a]`b``](x.md)", 0, 9},
		{"[a `` ` ``](x.md)", 0, 10},
		{"[`a](x.md)", 0, 3},
		{"[``a`](x.md)", 0, 5},
		{"[a [b] c](x.md)", 0, 8},
		{"[![img](i.png)](x.md)", 0, 14},
		{`[a \] b](x.md)`, 0, 7},
		{`[a \[ b](x.md)`, 0, 7},
		{"[\\`]`](x.md)", 0, 3},
		{"[unterminated", 0, -1},
		{"[a [b](x.md)", 0, -1},
		{"label](x.md)", 0, -1},
		{"[label]", -1, -1},
		{"[label]", 7, -1},
	} {
		if got := LinkLabelEnd([]byte(c.raw), c.start); got != c.want {
			t.Errorf("LinkLabelEnd(%q, %d) = %d, want %d", c.raw, c.start, got, c.want)
		}
	}
}

func TestFormatDestination(t *testing.T) {
	for _, c := range []struct{ dest, want string }{
		{"guide.md#topic", "<guide.md#topic>"},
		{"", "<>"},
		{"a b (c).md", "<a b (c).md>"},
		{"a<b>c", "<a%3Cb%3Ec>"},
		{"a\nb\r", "<a%0Ab%0D>"},
		{`a\_b\)`, `<a\_b\)>`},
		{`a\>b`, `<a\>b>`},
		{`dir\`, `<dir\\>`},
		{`C:\dir`, `<C:\dir>`},
		{"a\\\nb", `<a\\%0Ab>`},
	} {
		if got := FormatDestination(c.dest); got != c.want {
			t.Errorf("FormatDestination(%q) = %q, want %q", c.dest, got, c.want)
		}
	}
}

func TestFormatTitle(t *testing.T) {
	for _, c := range []struct{ title, want string }{
		{"Guide", `"Guide"`},
		{`say "hi"`, `"say \"hi\""`},
		{`say \"hi\"`, `"say \"hi\""`},
		{`ends in \\`, `"ends in \\"`},
		{`ends in \`, `"ends in \\"`},
		{`it\'s (a\)`, `"it\'s (a\)"`},
		{`\a`, `"\a"`},
	} {
		if got := FormatTitle(c.title); got != c.want {
			t.Errorf("FormatTitle(%q) = %q, want %q", c.title, got, c.want)
		}
	}
}
