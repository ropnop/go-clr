// +build windows

package clr

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

// GetICORRuntimeHost is a wrapper function that takes in an ICLRRuntimeInfo and returns an ICORRuntimeHost object
// and loads it into the current process. This is the "deprecated" API, but the only way currently to load an assembly
// from memory (afaict)
func GetICORRuntimeHost(runtimeInfo *ICLRRuntimeInfo) (*ICORRuntimeHost, error) {
	var pRuntimeHost uintptr
	hr := runtimeInfo.GetInterface(&CLSID_CorRuntimeHost, &IID_ICorRuntimeHost, &pRuntimeHost)
	err := checkOK(hr, "runtimeInfo.GetInterface")
	if err != nil {
		return nil, err
	}
	runtimeHost := NewICORRuntimeHostFromPtr(pRuntimeHost)
	hr = runtimeHost.Start()
	err = checkOK(hr, "runtimeHost.Start")
	return runtimeHost, err
}

func NewICORRuntimeHostFromPtr(ppv uintptr) *ICORRuntimeHost {
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
