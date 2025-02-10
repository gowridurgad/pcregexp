package pcregexp_test

import (
	"testing"

	"github.com/dwisiswant0/pcregexp"
)

func TestCompile(t *testing.T) {
	tests := []struct {
		name    string
		pattern string
		wantErr bool
	}{
		{"empty pattern", "", false},
		{"valid pattern", "a+b", false},
		{"invalid pattern", "a[", true},
		{"complex pattern", `\b\w+@\w+\.\w+\b`, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			re, err := pcregexp.Compile(tt.pattern)
			if (err != nil) != tt.wantErr {
				t.Errorf("Compile() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if err == nil {
				re.Close()
			}
		})
	}
}

func TestRegexp_Methods(t *testing.T) {
	re := pcregexp.MustCompile(`p([a-z]+)ch`)
	defer re.Close()

	t.Run("MatchString", func(t *testing.T) {
		tests := []struct {
			input string
			want  bool
		}{
			{"peach", true},
			{"punch", true},
			{"pinch", true},
			{"pch", false},
			{"each", false},
		}

		for _, tt := range tests {
			if got := re.MatchString(tt.input); got != tt.want {
				t.Errorf("MatchString(%q) = %v, want %v", tt.input, got, tt.want)
			}
		}
	})

	t.Run("FindString", func(t *testing.T) {
		tests := []struct {
			input string
			want  string
		}{
			{"peach punch", "peach"},
			{"no match", ""},
			{"pinch first", "pinch"},
		}

		for _, tt := range tests {
			if got := re.FindString(tt.input); got != tt.want {
				t.Errorf("FindString(%q) = %q, want %q", tt.input, got, tt.want)
			}
		}
	})

	t.Run("FindStringIndex", func(t *testing.T) {
		input := "peach punch"
		want := []int{0, 5}
		got := re.FindStringIndex(input)

		if got == nil || got[0] != want[0] || got[1] != want[1] {
			t.Errorf("FindStringIndex(%q) = %v, want %v", input, got, want)
		}
	})

	t.Run("FindStringSubmatch", func(t *testing.T) {
		tests := []struct {
			input string
			want  []string
		}{
			{"peach", []string{"peach", "ea"}},
			{"no match", nil},
		}

		for _, tt := range tests {
			got := re.FindStringSubmatch(tt.input)
			if (got == nil) != (tt.want == nil) {
				t.Errorf("FindStringSubmatch(%q) = %v, want %v", tt.input, got, tt.want)
				continue
			}

			if got != nil && (got[0] != tt.want[0] || got[1] != tt.want[1]) {
				t.Errorf("FindStringSubmatch(%q) = %v, want %v", tt.input, got, tt.want)
			}
		}
	})

	t.Run("ReplaceAllString", func(t *testing.T) {
		tests := []struct {
			input string
			repl  string
			want  string
		}{
			{"peach punch", "FRUIT", "FRUIT FRUIT"},
			{"no match here", "FRUIT", "no match here"},
			{"peach", "FRUIT", "FRUIT"},
		}

		for _, tt := range tests {
			if got := re.ReplaceAllString(tt.input, tt.repl); got != tt.want {
				t.Errorf("ReplaceAllString(%q, %q) = %q, want %q", tt.input, tt.repl, got, tt.want)
			}
		}
	})
}

func TestLookarounds(t *testing.T) {
	tests := []struct {
		name    string
		pattern string
		input   string
		want    bool
	}{
		{"positive lookahead", `foo(?=bar)`, "foobar", true},
		{"positive lookahead no match", `foo(?=bar)`, "foobaz", false},
		{"negative lookahead", `foo(?!bar)`, "foobaz", true},
		{"negative lookahead no match", `foo(?!bar)`, "foobar", false},
		{"positive lookbehind", `(?<=foo)bar`, "foobar", true},
		{"positive lookbehind no match", `(?<=foo)bar`, "bazbar", false},
		{"negative lookbehind", `(?<!foo)bar`, "bazbar", true},
		{"negative lookbehind no match", `(?<!foo)bar`, "foobar", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			re := pcregexp.MustCompile(tt.pattern)
			defer re.Close()

			if got := re.MatchString(tt.input); got != tt.want {
				t.Errorf("MatchString(%q) = %v, want %v", tt.input, got, tt.want)
			}
		})
	}
}

func TestBackreferences(t *testing.T) {
	tests := []struct {
		name    string
		pattern string
		input   string
		want    string
	}{
		{"simple backreference", `(\w+)\s+\1`, "hello hello world", "hello hello"},
		{"no match backreference", `(\w+)\s+\1`, "hello world", ""},
		{"nested groups", `((\w+)\s+\2)`, "hello hello world", "hello hello"},
		{"multiple backreferences", `(\w+)\s+(\w+)\s+\1\s+\2`, "cat dog cat dog", "cat dog cat dog"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			re := pcregexp.MustCompile(tt.pattern)
			defer re.Close()

			if got := re.FindString(tt.input); got != tt.want {
				t.Errorf("FindString(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}
