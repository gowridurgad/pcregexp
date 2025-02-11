package pcregexp

import (
	"fmt"
	"reflect"
	"runtime"
	"strings"
	"unicode/utf8"
	"unsafe"

	"github.com/ebitengine/purego"
)

func init() {
	var libPath string

	switch runtime.GOOS {
	case "darwin":
		libPath = "libpcre2-8.dylib"
	case "linux", "freebsd":
		libPath = "libpcre2-8.so"
	case "windows":
		libPath = "pcre2-8.dll"
	default:
		panic(fmt.Errorf("GOOS=%s is not supported", runtime.GOOS))
	}

	lib, err := openLibrary(libPath)
	if err != nil {
		panic(fmt.Errorf("failed to load %s: %w", libPath, err))
	}

	// Register the functions by their PCRE2 symbol names.
	// (For the 8-bit versions, the symbols are suffixed with "_8".)
	funcs := [][2]any{
		{&pcre2_compile, "pcre2_compile_8"},
		{&pcre2_code_free, "pcre2_code_free_8"},
		{&pcre2_pattern_info, "pcre2_pattern_info_8"},
		{&pcre2_match, "pcre2_match_8"},
		{&pcre2_match_data_create_from_pattern, "pcre2_match_data_create_from_pattern_8"},
		{&pcre2_match_data_free, "pcre2_match_data_free_8"},
		{&pcre2_get_ovector_pointer, "pcre2_get_ovector_pointer_8"},
	}

	for _, f := range funcs {
		purego.RegisterLibFunc(f[0], lib, f[1].(string))
	}
}

type PCREgexp struct {
	pattern   string  // original pattern
	buf       []int   // cached match offsets
	code      uintptr // pointer to compiled pcre2_code
	matchData uintptr // cached match data
}

// Compile compiles the given pattern and returns a [PCREgexp].
func Compile(pattern string) (*PCREgexp, error) {
	var patPtr *uint8
	var errcode int32
	var errOffset uint64

	if len(pattern) == 0 {
		var dummy byte = 0
		patPtr = &dummy
	} else {
		strHeader := (*reflect.StringHeader)(unsafe.Pointer(&pattern))
		patPtr = (*uint8)(unsafe.Pointer(strHeader.Data))
		// patPtr = (*uint8)(unsafe.StringData(pattern))
	}

	code := pcre2_compile(patPtr, uint64(len(pattern)), 0, &errcode, &errOffset, 0)
	if code == 0 {
		return nil, fmt.Errorf("pcre2_compile failed at offset %d, error code %d", errOffset, errcode)
	}

	return &PCREgexp{code: code, pattern: pattern}, nil
}

// MustCompile is like Compile but panics on error.
func MustCompile(pattern string) *PCREgexp {
	re, err := Compile(pattern)
	if err != nil {
		panic(err)
	}

	return re
}

// Close frees the resources associated with the compiled pattern.
func (re *PCREgexp) Close() {
	if re.matchData != 0 {
		pcre2_match_data_free(re.matchData)
		re.matchData = 0
	}

	if re.code != 0 {
		pcre2_code_free(re.code)
		re.code = 0
	}
}

// saveMatchData creates a new match data object if it doesn't exist yet.
//
// It returns the pointer to the match data object. The match data object is
// used to store the results of a match.
func (re *PCREgexp) saveMatchData() uintptr {
	if re.matchData == 0 {
		re.matchData = pcre2_match_data_create_from_pattern(re.code, 0)
	}

	return re.matchData
}

// match performs a PCRE2 match on the given subject.
//
// It returns a slice of start/end indexes as returned by PCRE2.
func (re *PCREgexp) match(subject []byte) []int {
	if re.code == 0 || len(subject) == 0 {
		return nil
	}

	md := re.saveMatchData()
	if md == 0 {
		return nil
	}

	var subjectPtr *uint8

	if len(subject) > 0 {
		subjectPtr = (*uint8)(ptr(&subject[0]))
	}

	ret := pcre2_match(re.code, subjectPtr, uint64(len(subject)), 0, 0, md, 0)
	if ret < 0 {
		return nil
	}

	n := int(ret)
	reqLen := n * 2

	if cap(re.buf) < reqLen {
		re.buf = make([]int, reqLen)
	} else {
		re.buf = re.buf[:reqLen]
	}

	ovector := pcre2_get_ovector_pointer(md)
	if ovector == nil {
		return nil
	}

	size := unsafe.Sizeof(uint64(0))
	for i := 0; i < reqLen; i++ {
		ptr := (*uint64)(ptr(uintptr(ptr(ovector)) + uintptr(i)*size))
		re.buf[i] = int(*ptr)
	}

	return re.buf
}

// MatchString reports whether the Regexp matches the given string.
func (re *PCREgexp) MatchString(s string) bool {
	return re.match(stringToBytesUnsafe(s)) != nil
}

// FindString returns the text of the leftmost match in s.
func (re *PCREgexp) FindString(s string) string {
	indexes := re.match(stringToBytesUnsafe(s))
	if indexes == nil || len(indexes) < 2 {
		return ""
	}

	return s[indexes[0]:indexes[1]]
}

// FindStringIndex returns a two-element slice of integers defining the start
// and end of the leftmost match in s.
func (re *PCREgexp) FindStringIndex(s string) []int {
	return re.match(stringToBytesUnsafe(s))
}

// FindStringSubmatch returns a slice holding the text of the leftmost match and
// its submatches. It uses the actual number of captured groups as returned by
// PCRE2.
func (re *PCREgexp) FindStringSubmatch(s string) []string {
	indexes := re.match(stringToBytesUnsafe(s))
	if indexes == nil || len(indexes) < 2 {
		return nil
	}

	n := len(indexes) / 2
	submatches := make([]string, n)

	for i := 0; i < n; i++ {
		start := indexes[2*i]
		end := indexes[2*i+1]
		if start < 0 || end < 0 {
			submatches[i] = ""
		} else {
			submatches[i] = s[start:end]
		}
	}

	return submatches
}

// ReplaceAllString returns a copy of src in which all matches of the [PCREgexp]
// have been replaced by repl.
//
// If an empty match is encountered, it advances one UTF-8 rune to avoid
// infinite loop.
func (re *PCREgexp) ReplaceAllString(src, repl string) string {
	if src == "" {
		return ""
	}

	var b strings.Builder
	b.Grow(len(src))

	remaining := src
	for {
		indexes := re.match(stringToBytesUnsafe(remaining))
		if indexes == nil || len(indexes) < 2 || indexes[0] < 0 {
			b.WriteString(remaining)
			break
		}

		b.WriteString(remaining[:indexes[0]])
		b.WriteString(repl)

		if indexes[0] == indexes[1] {
			if indexes[1] < len(remaining) {
				r, size := utf8.DecodeRuneInString(remaining[indexes[1]:])
				if r == utf8.RuneError || size == 0 {
					b.WriteString(remaining[indexes[1]:])
					break
				}

				b.WriteString(remaining[indexes[1] : indexes[1]+size])
				remaining = remaining[indexes[1]+size:]
			} else {
				remaining = ""
			}
		} else {
			remaining = remaining[indexes[1]:]
		}
	}

	return b.String()
}
