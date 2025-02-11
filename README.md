# pcregexp

`pcregexp` is a drop‑in replacement for Go's standard [`regexp`](https://pkg.go.dev/regexp) package that uses the full capabilities of [PCRE2](https://github.com/PCRE2Project/pcre2) by loading the shared library dynamically at runtime, which enables cross‑compilation without the need for a C compiler (**no Cgo required!**). The API closely mirrors that of the standard library's `regexp` package while supporting advanced regex features like lookarounds and backreferences that PCRE2 provides.

> [!WARNING]
> PCRE2 supports features that can lead to exponential runtime in some cases. Use `pcregexp` only with *trusted* regex patterns to avoid potential regular expression denial-of-service (ReDoS) issues ([CWE-1333](https://cwe.mitre.org/data/definitions/1333.html)).

## Requirements

* **go1.18** or later.
* **PCRE2 10.x shared library** must be [installed](https://github.com/PCRE2Project/pcre2#quickstart) on your system.
* Supported platforms:
  * [PCRE2](https://github.com/PCRE2Project/pcre2#platforms)
  * [purego](https://github.com/ebitengine/purego#supported-platforms)

## Install

```bash
go install -v github.com/dwisiswant0/pcregexp@latest
```

## Usage

```go
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
```

## Benchmark

```console
$ go test -run - -bench=. -benchmem
goos: linux
goarch: amd64
pkg: github.com/dwisiswant0/pcregexp
cpu: 11th Gen Intel(R) Core(TM) i9-11900H @ 2.50GHz
BenchmarkCompile/impl=pcregexp-16         	  816678	      1530 ns/op	     728 B/op	      13 allocs/op
BenchmarkCompile/impl=stdlib-16           	  534636	      2201 ns/op	    3272 B/op	      33 allocs/op
BenchmarkMatchString/impl=pcregexp-16     	  979143	      1070 ns/op	     736 B/op	      15 allocs/op
BenchmarkMatchString/impl=stdlib-16       	10747933	       150.1 ns/op	       0 B/op	       0 allocs/op
BenchmarkFindString/impl=pcregexp-16      	 1128944	      1168 ns/op	     736 B/op	      15 allocs/op
BenchmarkFindString/impl=stdlib-16        	10892875	       171.6 ns/op	       0 B/op	       0 allocs/op
BenchmarkReplaceAllString/impl=pcregexp-16         	  310188	      3455 ns/op	    2232 B/op	      46 allocs/op
BenchmarkReplaceAllString/impl=stdlib-16           	 1949434	       658.9 ns/op	      96 B/op	       5 allocs/op
PASS
ok  	github.com/dwisiswant0/pcregexp	13.346s
```

## TODO

* [ ] Implement PCRE2 JIT compilation support
  * Use native PCRE2 API JIT functions for improved performance
  * Add JIT compilation options and configurations
  * Implement memory management for JIT-compiled patterns

## Status

> [!CAUTION]
> `pcregexp` has NOT reached 1.0 yet. Therefore, this library is currently not supported and does not offer a stable API; use at your own risk.

There are no guarantees of stability for the APIs in this library, and while they are not expected to change dramatically. API tweaks and bug fixes may occur.

## License

`pcregexp` is released by [**@dwisiswant0**](https://github.com/dwisiswant0) under the Apache 2.0 license. See [LICENSE](/LICENSE).