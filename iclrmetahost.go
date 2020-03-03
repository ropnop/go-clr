//+build windows

package main

import (
	"github.com/Microsoft/go-winio/pkg/guid"
	"syscall"
	"unsafe"
)

var (
	modMSCoree            = syscall.MustLoadDLL("mscoree.dll")
	procCLRCreateInstance = modMSCoree.MustFindProc("CLRCreateInstance")
)

func CLRCreateInstance(clsid, riid *guid.GUID, ppInterface *uintptr) uintptr {
	ret, _, _ := procCLRCreateInstance.Call(
		uintptr(unsafe.Pointer(clsid)),
		uintptr(unsafe.Pointer(riid)),
		uintptr(unsafe.Pointer(ppInterface)))
	return ret
}

// GetICLRMetaHost
//func GetICLRMetaHost() (*ICLRMetaHost, error) {
//	var ppInterface uintptr
//	if hr := getICLRMetaHostPtr(&ppInterface); hr != S_OK {
//		return nil, fmt.Errorf("Could not get pointer to ICLRMetaHost. HRESULT: 0x%x", hr)
//	}
//	return newICLRMetaHost(ppInterface), nil
//}

//ICLRMetaHost Interface from metahost.h
// https://stackoverflow.com/questions/37781676/how-to-use-com-component-object-model-in-golang
type ICLRMetaHost struct {
	vtbl *ICLRMetaHostVtbl
}
type ICLRMetaHostVtbl struct {
	QueryInterface                   uintptr
	AddRef                           uintptr
	Release                          uintptr
	GetRuntime                       uintptr
	GetVersionFromFile               uintptr
	EnumerateInstalledRuntimes       uintptr
	EnumerateLoadedRuntimes          uintptr
	RequestRuntimeLoadedNotification uintptr
	QueryLegacyV2RuntimeBinding      uintptr
	ExitProcess                      uintptr
}

// newICLRMetaHost takes a uintptr to ICLRMetahost, which must come from the syscall CLRCreateInstance
func newICLRMetaHost(ppv uintptr) *ICLRMetaHost {
	return (*ICLRMetaHost)(unsafe.Pointer(ppv))
}

func (obj *ICLRMetaHost) AddRef() uintptr {
	ret, _, _ := syscall.Syscall(
		obj.vtbl.AddRef,
		1,
		uintptr(unsafe.Pointer(obj)),
		0,
		0)
	return ret
}

func (obj *ICLRMetaHost) Release() uintptr {
	ret, _, _ := syscall.Syscall(
		obj.vtbl.Release,
		1,
		uintptr(unsafe.Pointer(obj)),
		0,
		0)
	return ret
}

func (obj *ICLRMetaHost) EnumerateInstalledRuntimes(pInstalledRuntimes *uintptr) uintptr {
	ret, _, _ := syscall.Syscall(
		obj.vtbl.EnumerateInstalledRuntimes,
		2,
		uintptr(unsafe.Pointer(obj)),
		uintptr(unsafe.Pointer(pInstalledRuntimes)),
		0)
	return ret
}

func (obj *ICLRMetaHost) GetRuntime(pwzVersion *uint16, riid *guid.GUID, pRuntimeHost *uintptr) uintptr {
	v4Ptr, err := syscall.UTF16PtrFromString("v4.0.30319")
	if err != nil {
		panic(err)
	}
	ret, _, _ := syscall.Syscall6(
		obj.vtbl.GetRuntime,
		4,
		uintptr(unsafe.Pointer(obj)),
		uintptr(unsafe.Pointer(v4Ptr)),
		uintptr(unsafe.Pointer(&IID_ICLRRuntimeInfo)),
		uintptr(unsafe.Pointer(pRuntimeHost)),
		0,
		0)
	return ret
}
