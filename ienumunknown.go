// +build windows

package clr

import (
	"fmt"
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

// Next retrieves the specified number of items in the enumeration sequence.
// HRESULT Next(
//   ULONG    celt,
//   IUnknown **rgelt,
//   ULONG    *pceltFetched
// );
// https://docs.microsoft.com/en-us/windows/win32/api/objidl/nf-objidl-ienumunknown-next
func (obj *IEnumUnknown) Next(celt uint32, pEnumRuntime unsafe.Pointer, pceltFetched *uint32) (err error) {
	debugPrint("Entering into ienumunknown.Next()...")
	hr, _, err := syscall.Syscall6(
		obj.vtbl.Next,
		4,
		uintptr(unsafe.Pointer(obj)),
		uintptr(celt),
		uintptr(pEnumRuntime),
		uintptr(unsafe.Pointer(pceltFetched)),
		0,
		0,
	)
	if err != syscall.Errno(0) {
		err = fmt.Errorf("there was an error calling the IEnumUnknown::Next method:\r\n%s", err)
		return
	}
	if hr != S_OK {
		err = fmt.Errorf("the IEnumUnknown::Next method method returned a non-zero HRESULT: 0x%x", hr)
		return
	}
	err = nil
	return
}
