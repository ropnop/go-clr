// +build windows

package clr

import (
	"syscall"
	"unsafe"

	"golang.org/x/sys/windows"
)

var (
	modMSCoree            = syscall.NewLazyDLL("mscoree.dll")
	procCLRCreateInstance = modMSCoree.NewProc("CLRCreateInstance")
)

func CLRCreateInstance(clsid, riid *windows.GUID, ppInterface *uintptr) uintptr {
	ret, _, _ := procCLRCreateInstance.Call(
		uintptr(unsafe.Pointer(clsid)),
		uintptr(unsafe.Pointer(riid)),
		uintptr(unsafe.Pointer(ppInterface)))
	return ret
}

func GetICLRMetaHost() (metahost *ICLRMetaHost, err error) {
	var pMetaHost uintptr
	hr := CLRCreateInstance(&CLSID_CLRMetaHost, &IID_ICLRMetaHost, &pMetaHost)
	err = checkOK(hr, "CLRCreateInstance")
	if err != nil {
		return
	}
	metahost = NewICLRMetaHost(pMetaHost)
	return
}

//ICLRMetaHost Interface from metahost.h
// Couldnt have done any of this without this SO answer I stumbled on:
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
func NewICLRMetaHost(ppv uintptr) *ICLRMetaHost {
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

func (obj *ICLRMetaHost) GetRuntime(pwzVersion *uint16, riid *windows.GUID, pRuntimeHost *uintptr) uintptr {
	// v4Ptr, err := syscall.UTF16PtrFromString("v4.0.30319")
	// if err != nil {
	// 	panic(err)
	// }
	ret, _, _ := syscall.Syscall6(
		obj.vtbl.GetRuntime,
		4,
		uintptr(unsafe.Pointer(obj)),
		uintptr(unsafe.Pointer(pwzVersion)),
		uintptr(unsafe.Pointer(&IID_ICLRRuntimeInfo)),
		uintptr(unsafe.Pointer(pRuntimeHost)),
		0,
		0)
	return ret
}
