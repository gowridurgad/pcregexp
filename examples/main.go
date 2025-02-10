package main

import (
	"fmt"

	"github.com/dwisiswant0/pcregexp"
)

func main() {
	// Compile a pattern. (Panics on error)
	re := pcregexp.MustCompile("p([a-z]+)ch")
	defer re.Close()

	// Check if the string matches the pattern.
	fmt.Println("MatchString(\"peach\"):", re.MatchString("peach"))

	// Find the leftmost match.
	fmt.Println("FindString(\"peach punch\"):", re.FindString("peach punch"))

	// Get the start and end indexes of the match.
	fmt.Println("FindStringIndex(\"peach punch\"):", re.FindStringIndex("peach punch"))

	// Retrieve the match along with its captured submatch.
	fmt.Println("FindStringSubmatch(\"peach punch\"):", re.FindStringSubmatch("peach punch"))

	// Perform global replacement (naively replaces all non-overlapping matches).
	src := "peach punch pinch"
	repl := "<fruit>"
	fmt.Println("ReplaceAllString:", re.ReplaceAllString(src, repl))
}
