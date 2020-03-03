package main

import (
	"syscall"
	"unsafe"
)

type ICORRuntimeHost struct {
	vtbl *ICORRuntimeHostVtbl
}

type ICORRuntimeHostVtbl struct {
	QueryInterface                uintptr
	AddRef                        uintptr
	Release                       uintptr
	CreateLogicalThreadState      uintptr
	DeleteLogicalThreadState      uintptr
	SwitchInLogicalThreadState    uintptr
	SwitchOutLogicalThreadState   uintptr
	LocksHeldByLogicalThreadState uintptr
	MapFile                       uintptr
	GetConfiguration              uintptr
	Start                         uintptr
	Stop                          uintptr
	CreateDomain                  uintptr
	GetDefaultDomain              uintptr
	EnumDomains                   uintptr
	NextDomain                    uintptr
	CloseEnum                     uintptr
	CreateDomainEx                uintptr
	CreateDomainSetup             uintptr
	CreateEvidence                uintptr
	UnloadDomain                  uintptr
	CurrentDomain                 uintptr
}

func newICORRuntimeHost(ppv uintptr) *ICORRuntimeHost {
	return (*ICORRuntimeHost)(unsafe.Pointer(ppv))
}

func (obj *ICORRuntimeHost) AddRef() uintptr {
	ret, _, _ := syscall.Syscall(
		obj.vtbl.AddRef,
		1,
		uintptr(unsafe.Pointer(obj)),
		0,
		0)
	return ret
}

func (obj *ICORRuntimeHost) Release() uintptr {
	ret, _, _ := syscall.Syscall(
		obj.vtbl.Release,
		1,
		uintptr(unsafe.Pointer(obj)),
		0,
		0)
	return ret
}

func (obj *ICORRuntimeHost) Start() uintptr {
	ret, _, _ := syscall.Syscall(
		obj.vtbl.Start,
		1,
		uintptr(unsafe.Pointer(obj)),
		0,
		0)
	return ret
}

func (obj *ICORRuntimeHost) GetDefaultDomain(pAppDomain *uintptr) uintptr {
	ret, _, _ := syscall.Syscall(
		obj.vtbl.GetDefaultDomain,
		2,
		uintptr(unsafe.Pointer(obj)),
		uintptr(unsafe.Pointer(pAppDomain)),
		0)
	return ret
}
