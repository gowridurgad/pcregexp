package pcregexp

import "unsafe"

// stringToBytesUnsafe returns a byte slice header that points to the string's
// data. This conversion is safe only if the receiver does not modify the
// returned slice.
func stringToBytesUnsafe(s string) []byte {
	return *(*[]byte)(unsafe.Pointer(&s))
}

// ptr aliases [unsafe.Pointer].
type ptr = unsafe.Pointer
