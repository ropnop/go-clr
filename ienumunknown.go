// +build windows

package clr

import (
	"syscall"
	"unsafe"
)

type IEnumUnknown struct {
	vtbl *IEnumUnknownVtbl
}

type IEnumUnknownVtbl struct {
	QueryInterface uintptr
	AddRef         uintptr
	Release        uintptr
	Next           uintptr
	Skip           uintptr
	Reset          uintptr
	Clone          uintptr
}

func NewIEnumUnknownFromPtr(ppv uintptr) *IEnumUnknown {
	return (*IEnumUnknown)(unsafe.Pointer(ppv))
}

func (obj *IEnumUnknown) AddRef() uintptr {
	ret, _, _ := syscall.Syscall(
		obj.vtbl.AddRef,
		1,
		uintptr(unsafe.Pointer(obj)),
		0,
		0)
	return ret
}

func (obj *IEnumUnknown) Release() uintptr {
	ret, _, _ := syscall.Syscall(
		obj.vtbl.Release,
		1,
		uintptr(unsafe.Pointer(obj)),
		0,
		0)
	return ret
}

func (obj *IEnumUnknown) Next(celt uint32, pEnumRuntime *uintptr, pCeltFetched *uint32) uintptr {
	ret, _, _ := syscall.Syscall6(
		obj.vtbl.Next,
		4,
		uintptr(unsafe.Pointer(obj)),
		uintptr(celt),
		uintptr(unsafe.Pointer(pEnumRuntime)),
		uintptr(unsafe.Pointer(pCeltFetched)),
		0,
		0)
	return ret
}
