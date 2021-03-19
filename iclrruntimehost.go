// +build windows

package clr

import (
	"fmt"
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
	debugPrint("Entering into iclrruntimehost.GetICLRRuntimeHost()...")
	var runtimeHost *ICLRRuntimeHost
	err := runtimeInfo.GetInterface(CLSID_CLRRuntimeHost, IID_ICLRRuntimeHost, unsafe.Pointer(&runtimeHost))
	if err != nil {
		return nil, err
	}

	err = runtimeHost.Start()
	return runtimeHost, err
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

// Start Initializes the common language runtime (CLR) into a process.
// HRESULT Start();
// https://docs.microsoft.com/en-us/dotnet/framework/unmanaged-api/hosting/iclrruntimehost-start-method
func (obj *ICLRRuntimeHost) Start() error {
	debugPrint("Entering into iclrruntimehost.Start()...")
	hr, _, err := syscall.Syscall(
		obj.vtbl.Start,
		1,
		uintptr(unsafe.Pointer(obj)),
		0,
		0,
	)
	if err != syscall.Errno(0) {
		return fmt.Errorf("the ICLRRuntimeHost::Start method returned an error:\r\n%s", err)
	}
	if hr != S_OK {
		return fmt.Errorf("the ICLRRuntimeHost::Start method method returned a non-zero HRESULT: 0x%x", hr)
	}
	return nil
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
