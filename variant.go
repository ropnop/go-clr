// +build windows

package clr

import "unsafe"

// from https://github.com/go-ole/go-ole/blob/master/variant_amd64.go

type Variant struct {
	VT         uint16 // VARTYPE
	wReserved1 uint16
	wReserved2 uint16
	wReserved3 uint16
	Val        uintptr
	_          [8]byte
}

func NewVariantFromPtr(ppv uintptr) *Variant {
	return (*Variant)(unsafe.Pointer(ppv))
}
