// +build windows

package clr

import (
	"fmt"
	"syscall"
	"unsafe"

	"golang.org/x/sys/windows"
)

type IUnknown struct {
	vtbl *IUnknownVtbl
}

type IUnknownVtbl struct {
	QueryInterface uintptr
	AddRef         uintptr
	Release        uintptr
}

// QueryInterface queries a COM object for a pointer to one of its interface;
// identifying the interface by a reference to its interface identifier (IID).
// If the COM object implements the interface, then it returns a pointer to that interface after calling IUnknown::AddRef on it.
// HRESULT QueryInterface(
//   REFIID riid,
//   void   **ppvObject
// );
// https://docs.microsoft.com/en-us/windows/win32/api/unknwn/nf-unknwn-iunknown-queryinterface(refiid_void)
func (obj *IUnknown) QueryInterface(riid windows.GUID, ppvObject unsafe.Pointer) error {
	debugPrint("Entering into iunknown.QueryInterface()...")
	hr, _, err := syscall.Syscall(
		obj.vtbl.QueryInterface,
		3,
		uintptr(unsafe.Pointer(obj)),
		uintptr(unsafe.Pointer(&riid)), // A reference to the interface identifier (IID) of the interface being queried for.
		uintptr(ppvObject),
	)
	if err != syscall.Errno(0) {
		return fmt.Errorf("the ICORRuntimeHost::GetDefaultDomain method returned an error:\r\n%s", err)
	}
	if hr != S_OK {
		return fmt.Errorf("the ICORRuntimeHost::GetDefaultDomain method method returned a non-zero HRESULT: 0x%x", hr)
	}
	return nil
}

func (obj *IUnknown) AddRef() uintptr {
	ret, _, _ := syscall.Syscall(
		obj.vtbl.AddRef,
		1,
		uintptr(unsafe.Pointer(obj)),
		0,
		0)
	return ret
}

func (obj *IUnknown) Release() uintptr {
	ret, _, _ := syscall.Syscall(
		obj.vtbl.Release,
		1,
		uintptr(unsafe.Pointer(obj)),
		0,
		0)
	return ret
}
