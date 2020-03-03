package main

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

func newICLRRuntimeHost(ppv uintptr) *ICLRRuntimeHost {
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

func (obj *ICLRRuntimeHost) GetCurrentappDomainID(pdwAppDomainId *uint16) uintptr {
	ret, _, _ := syscall.Syscall(
		obj.vtbl.GetCurrentAppDomainId,
		2,
		uintptr(unsafe.Pointer(obj)),
		uintptr(unsafe.Pointer(pdwAppDomainId)),
		0)
	return ret
}
