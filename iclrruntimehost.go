// +build windows

package clr

import (
	"syscall"
	"unsafe"
)

type ICLRRuntimeHost struct {
	vtbl *ICLRRuntimeHostVtbl
}

type ICLRRuntimeHostVtbl struct {
	QueryInterface            uintptr
	AddRef                    uintptr
	Release                   uintptr
	Start                     uintptr
	Stop                      uintptr
	SetHostControl            uintptr
	GetCLRControl             uintptr
	UnloadAppDomain           uintptr
	ExecuteInAppDomain        uintptr
	GetCurrentAppDomainId     uintptr
	ExecuteApplication        uintptr
	ExecuteInDefaultAppDomain uintptr
}
// GetICLRRuntimeHost is a wrapper function that takes an ICLRRuntimeInfo object and
// returns an ICLRRuntimeHost and loads it into the current process
func GetICLRRuntimeHost(runtimeInfo *ICLRRuntimeInfo) (*ICLRRuntimeHost, error) {
	var pRuntimeHost uintptr
	hr := runtimeInfo.GetInterface(&CLSID_CLRRuntimeHost, &IID_ICLRRuntimeHost, &pRuntimeHost)
	err := checkOK(hr, "runtimeInfo.GetInterface")
	if err != nil {
		return nil, err
	}
	runtimeHost := NewICLRRuntimeHostFromPtr(pRuntimeHost)
	hr = runtimeHost.Start()
	err = checkOK(hr, "runtimeHost.Start")
	return runtimeHost, err
}

func NewICLRRuntimeHostFromPtr(ppv uintptr) *ICLRRuntimeHost {
	return (*ICLRRuntimeHost)(unsafe.Pointer(ppv))
}

func (obj *ICLRRuntimeHost) AddRef() uintptr {
	ret, _, _ := syscall.Syscall(
		obj.vtbl.AddRef,
		1,
		uintptr(unsafe.Pointer(obj)),
		0,
		0)
	return ret
}

func (obj *ICLRRuntimeHost) Release() uintptr {
	ret, _, _ := syscall.Syscall(
		obj.vtbl.Release,
		1,
		uintptr(unsafe.Pointer(obj)),
		0,
		0)
	return ret
}

func (obj *ICLRRuntimeHost) Start() uintptr {
	ret, _, _ := syscall.Syscall(
		obj.vtbl.Start,
		1,
		uintptr(unsafe.Pointer(obj)),
		0,
		0)
	return ret
}

func (obj *ICLRRuntimeHost) ExecuteInDefaultAppDomain(pwzAssemblyPath, pwzTypeName, pwzMethodName, pwzArgument, pReturnValue *uint16) uintptr {
	ret, _, _ := syscall.Syscall9(
		obj.vtbl.ExecuteInDefaultAppDomain,
		6,
		uintptr(unsafe.Pointer(obj)),
		uintptr(unsafe.Pointer(pwzAssemblyPath)),
		uintptr(unsafe.Pointer(pwzTypeName)),
		uintptr(unsafe.Pointer(pwzMethodName)),
		uintptr(unsafe.Pointer(pwzArgument)),
		uintptr(unsafe.Pointer(pReturnValue)),
		0,
		0,
		0)
	return ret
}

func (obj *ICLRRuntimeHost) GetCurrentAppDomainID(pdwAppDomainId *uint16) uintptr {
	ret, _, _ := syscall.Syscall(
		obj.vtbl.GetCurrentAppDomainId,
		2,
		uintptr(unsafe.Pointer(obj)),
		uintptr(unsafe.Pointer(pdwAppDomainId)),
		0)
	return ret
}
