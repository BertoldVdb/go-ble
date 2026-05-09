package bleadvertiser

import "testing"

func TestTruncateUTF8(t *testing.T) {
	cases := []struct {
		name    string
		in      string
		max     int
		want    string
	}{
		{"ascii under limit", "hello", 27, "hello"},
		{"ascii at limit", "0123456789012345678901234567", 27, "012345678901234567890123456"},
		{"empty", "", 27, ""},
		{"max=0", "abc", 0, ""},
		{"multibyte cut at boundary", "αβγδε", 4, "αβ"},     // 2-byte runes; 4 bytes = 2 runes
		{"multibyte cut mid-rune", "αβγδε", 5, "αβ"},        // would split at byte 5; back off to byte 4
		{"emoji", "x🚀x", 4, "x"},                            // 🚀 is 4 bytes; cut at byte 4 falls inside it → "x" only
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got := truncateUTF8(c.in, c.max)
			if got != c.want {
				t.Errorf("truncateUTF8(%q, %d) = %q, want %q", c.in, c.max, got, c.want)
			}
		})
	}
}
