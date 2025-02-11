package pcregexp_test

import (
	"bytes"
	"reflect"
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

func TestRegexp_ByteMethods(t *testing.T) {
	re := pcregexp.MustCompile(`p([a-z]+)ch`)
	defer re.Close()

	t.Run("Match", func(t *testing.T) {
		tests := []struct {
			input []byte
			want  bool
		}{
			{[]byte("peach"), true},
			{[]byte("punch"), true},
			{[]byte("pch"), false},
		}

		for _, tt := range tests {
			if got := re.Match(tt.input); got != tt.want {
				t.Errorf("Match(%q) = %v, want %v", tt.input, got, tt.want)
			}
		}
	})

	t.Run("Find", func(t *testing.T) {
		tests := []struct {
			input []byte
			want  []byte
		}{
			{[]byte("peach punch"), []byte("peach")},
			{[]byte("no match"), nil},
		}

		for _, tt := range tests {
			got := re.Find(tt.input)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Find(%q) = %q, want %q", tt.input, got, tt.want)
			}
		}
	})

	t.Run("FindSubmatch", func(t *testing.T) {
		input := []byte("peach")
		want := [][]byte{[]byte("peach"), []byte("ea")}
		got := re.FindSubmatch(input)

		if len(got) != len(want) {
			t.Errorf("FindSubmatch(%q) returned %d submatches, want %d", input, len(got), len(want))
			return
		}

		for i := range got {
			if !bytes.Equal(got[i], want[i]) {
				t.Errorf("FindSubmatch(%q)[%d] = %q, want %q", input, i, got[i], want[i])
			}
		}
	})
}

func TestRegexp_FindAll(t *testing.T) {
	re := pcregexp.MustCompile(`p([a-z]+)ch`)
	defer re.Close()

	t.Run("FindAllString", func(t *testing.T) {
		tests := []struct {
			input string
			n     int
			want  []string
		}{
			{"peach punch pinch", -1, []string{"peach", "punch", "pinch"}},
			{"peach punch pinch", 2, []string{"peach", "punch"}},
			{"peach punch pinch", 0, nil},
			{"no matches", -1, nil},
		}

		for _, tt := range tests {
			got := re.FindAllString(tt.input, tt.n)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("FindAllString(%q, %d) = %v, want %v", tt.input, tt.n, got, tt.want)
			}
		}
	})

	t.Run("FindAllStringSubmatch", func(t *testing.T) {
		input := "peach punch"
		want := [][]string{
			{"peach", "ea"},
			{"punch", "un"},
		}
		got := re.FindAllStringSubmatch(input, -1)

		if !reflect.DeepEqual(got, want) {
			t.Errorf("FindAllStringSubmatch(%q, -1) = %v, want %v", input, got, want)
		}
	})

	t.Run("FindAllStringIndex", func(t *testing.T) {
		input := "peach punch"
		want := [][]int{{0, 5}, {6, 11}}
		got := re.FindAllStringIndex(input, -1)

		if !reflect.DeepEqual(got, want) {
			t.Errorf("FindAllStringIndex(%q, -1) = %v, want %v", input, got, want)
		}
	})
}

func TestRegexp_ReplaceAll(t *testing.T) {
	re := pcregexp.MustCompile(`a([a-z])e`)
	defer re.Close()

	t.Run("ReplaceAll", func(t *testing.T) {
		tests := []struct {
			src  []byte
			repl []byte
			want []byte
		}{
			{[]byte("age ace"), []byte("X"), []byte("X X")},
			{[]byte("no match"), []byte("X"), []byte("no match")},
		}

		for _, tt := range tests {
			got := re.ReplaceAll(tt.src, tt.repl)
			if !bytes.Equal(got, tt.want) {
				t.Errorf("ReplaceAll(%q, %q) = %q, want %q", tt.src, tt.repl, got, tt.want)
			}
		}
	})

	t.Run("ReplaceAllFunc", func(t *testing.T) {
		input := []byte("age ace")
		want := []byte("AGE ACE")
		got := re.ReplaceAllFunc(input, bytes.ToUpper)

		if !bytes.Equal(got, want) {
			t.Errorf("ReplaceAllFunc(%q, bytes.ToUpper) = %q, want %q", input, got, want)
		}
	})
}

func TestRegexp_Split(t *testing.T) {
	re := pcregexp.MustCompile(`\s+`)
	defer re.Close()

	tests := []struct {
		input string
		n     int
		want  []string
	}{
		{"foo bar baz", -1, []string{"foo", "bar", "baz"}},
		{"foo bar baz", 2, []string{"foo", "bar baz"}},
		{"foo", -1, []string{"foo"}},
		{"", -1, []string{""}},
	}

	for _, tt := range tests {
		got := re.Split(tt.input, tt.n)
		if !reflect.DeepEqual(got, tt.want) {
			t.Errorf("Split(%q, %d) = %v, want %v", tt.input, tt.n, got, tt.want)
		}
	}
}

func TestRegexp_Utility(t *testing.T) {
	pattern := `p([a-z]+)ch`
	re := pcregexp.MustCompile(pattern)
	defer re.Close()

	t.Run("String", func(t *testing.T) {
		if got := re.String(); got != pattern {
			t.Errorf("String() = %q, want %q", got, pattern)
		}
	})

	// t.Run("NumSubexp", func(t *testing.T) {
	// 	want := 1
	// 	if got := re.NumSubexp(); got != want {
	// 		t.Errorf("NumSubexp() = %d, want %d", got, want)
	// 	}
	// })
}
