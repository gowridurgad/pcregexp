package pcregexp_test

import (
	"regexp"
	"testing"

	"github.com/dwisiswant0/pcregexp"
)

func BenchmarkCompile(b *testing.B) {
	pattern := `\b\w+@\w+\.\w+\b`

	b.Run("impl=pcregexp", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			re, _ := pcregexp.Compile(pattern)
			re.Close()
		}
	})

	b.Run("impl=stdlib", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			regexp.Compile(pattern)
		}
	})
}

func BenchmarkMatchString(b *testing.B) {
	pattern := `p([a-z]+)ch`
	text := "peach punch pinch"

	pcre := pcregexp.MustCompile(pattern)
	defer pcre.Close()
	re := regexp.MustCompile(pattern)

	b.Run("impl=pcregexp", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			pcre.MatchString(text)
		}
	})

	b.Run("impl=stdlib", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			re.MatchString(text)
		}
	})
}

func BenchmarkFindString(b *testing.B) {
	pattern := `p([a-z]+)ch`
	text := "peach punch pinch"

	pcre := pcregexp.MustCompile(pattern)
	defer pcre.Close()
	re := regexp.MustCompile(pattern)

	b.Run("impl=pcregexp", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			pcre.FindString(text)
		}
	})

	b.Run("impl=stdlib", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			re.FindString(text)
		}
	})
}

func BenchmarkReplaceAllString(b *testing.B) {
	pattern := `p([a-z]+)ch`
	text := "peach punch pinch"
	repl := "FRUIT"

	pcre := pcregexp.MustCompile(pattern)
	defer pcre.Close()
	re := regexp.MustCompile(pattern)

	b.Run("impl=pcregexp", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			pcre.ReplaceAllString(text, repl)
		}
	})

	b.Run("impl=stdlib", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			re.ReplaceAllString(text, repl)
		}
	})
}
