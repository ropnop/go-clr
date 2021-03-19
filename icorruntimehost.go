// +build windows

package clr

import (
	"fmt"
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
	debugPrint("Entering into icorruntimehost.GetICORRuntimeHost()...")
	var runtimeHost *ICORRuntimeHost
	err := runtimeInfo.GetInterface(CLSID_CorRuntimeHost, IID_ICorRuntimeHost, unsafe.Pointer(&runtimeHost))
	if err != nil {
		return nil, err
	}

	err = runtimeHost.Start()
	return runtimeHost, err
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

// Start starts the common language runtime (CLR).
// HRESULT Start ();
// https://docs.microsoft.com/en-us/dotnet/framework/unmanaged-api/hosting/icorruntimehost-start-method
func (obj *ICORRuntimeHost) Start() error {
	debugPrint("Entering into icorruntimehost.Start()...")
	hr, _, err := syscall.Syscall(
		obj.vtbl.Start,
		1,
		uintptr(unsafe.Pointer(obj)),
		0,
		0,
	)
	if err != syscall.Errno(0) {
		// The system could not find the environment option that was entered.
		// TODO Why is this error message returned?
		fmt.Printf("the ICORRuntimeHost::Start method returned an error:\r\n%s\n", err)
	}
	if hr != S_OK {
		return fmt.Errorf("the ICORRuntimeHost::Start method method returned a non-zero HRESULT: 0x%x", hr)
	}
	return nil
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
