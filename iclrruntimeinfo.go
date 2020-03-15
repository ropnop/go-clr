// +build windows

package clr

import (
	"golang.org/x/sys/windows"
	"syscall"
	"unsafe"
)

type ICLRRuntimeInfo struct {
	vtbl *ICLRRuntimeInfoVtbl
}

type ICLRRuntimeInfoVtbl struct {
	QueryInterface         uintptr
	AddRef                 uintptr
	Release                uintptr
	GetVersionString       uintptr
	GetRuntimeDirectory    uintptr
	IsLoaded               uintptr
	LoadErrorString        uintptr
	LoadLibrary            uintptr
	GetProcAddress         uintptr
	GetInterface           uintptr
	IsLoadable             uintptr
	SetDefaultStartupFlags uintptr
	GetDefaultStartupFlags uintptr
	BindAsLegacyV2Runtime  uintptr
	IsStarted              uintptr
}

// GetRuntimeInfo is a wrapper function to return an ICLRRuntimeInfo from a standard version string
func GetRuntimeInfo(metahost *ICLRMetaHost, version string) (*ICLRRuntimeInfo, error) {
	pwzVersion, err := syscall.UTF16PtrFromString(version)
	if err != nil {
		return nil, err
	}
	var pRuntimeInfo uintptr
	hr := metahost.GetRuntime(pwzVersion, &IID_ICLRRuntimeInfo, &pRuntimeInfo)
	err = checkOK(hr, "metahost.GetRuntime")
	if err != nil {
		return nil, err
	}
	return NewICLRRuntimeInfoFromPtr(pRuntimeInfo), nil
}

func NewICLRRuntimeInfoFromPtr(ppv uintptr) *ICLRRuntimeInfo {
	return (*ICLRRuntimeInfo)(unsafe.Pointer(ppv))
}

func (obj *ICLRRuntimeInfo) AddRef() uintptr {
	ret, _, _ := syscall.Syscall(
		obj.vtbl.AddRef,
		1,
		uintptr(unsafe.Pointer(obj)),
		0,
		0)
	return ret
}

func (obj *ICLRRuntimeInfo) Release() uintptr {
	ret, _, _ := syscall.Syscall(
		obj.vtbl.Release,
		1,
		uintptr(unsafe.Pointer(obj)),
		0,
		0)
	return ret
}

func (obj *ICLRRuntimeInfo) GetVersionString(pcchBuffer *uint16, pVersionstringSize *uint32) uintptr {
	ret, _, _ := syscall.Syscall(
		obj.vtbl.GetVersionString,
		3,
		uintptr(unsafe.Pointer(obj)),
		uintptr(unsafe.Pointer(pcchBuffer)),
		uintptr(unsafe.Pointer(&pVersionstringSize)))
	return ret
}

func (obj *ICLRRuntimeInfo) GetInterface(rclsid *windows.GUID, riid *windows.GUID, ppUnk *uintptr) uintptr {
	ret, _, _ := syscall.Syscall6(
		obj.vtbl.GetInterface,
		4,
		uintptr(unsafe.Pointer(obj)),
		uintptr(unsafe.Pointer(rclsid)),
		uintptr(unsafe.Pointer(riid)),
		uintptr(unsafe.Pointer(ppUnk)),
		0,
		0)
	return ret
}

func (obj *ICLRRuntimeInfo) BindAsLegacyV2Runtime() uintptr {
	ret, _, _ := syscall.Syscall(
		obj.vtbl.BindAsLegacyV2Runtime,
		1,
		uintptr(unsafe.Pointer(obj)),
		0,
		0,
	)
	return ret
}

func (obj *ICLRRuntimeInfo) IsLoadable(pbLoadable *bool) uintptr {
	ret, _, _ := syscall.Syscall(
		obj.vtbl.IsLoadable,
		2,
		uintptr(unsafe.Pointer(obj)),
		uintptr(unsafe.Pointer(pbLoadable)),
		0)
	return ret
}
