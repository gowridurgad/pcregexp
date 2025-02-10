//go:build windows
// +build windows

package pcregexp

import "syscall"

func openLibrary(name string) (uintptr, error) {
	handle, err := syscall.LoadLibrary(name)
	return uintptr(handle), err
}
