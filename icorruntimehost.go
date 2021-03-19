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

// ICORRuntimeHos Provides methods that enable the host to start and stop the common language runtime (CLR)
// explicitly, to create and configure application domains, to access the default domain, and to enumerate all
// domains running in the process.
type ICORRuntimeHostVtbl struct {
	QueryInterface uintptr
	AddRef         uintptr
	Release        uintptr
	// CreateLogicalThreadState Do not use.
	CreateLogicalThreadState uintptr
	// DeleteLogicalThreadSate Do not use.
	DeleteLogicalThreadState uintptr
	// SwitchInLogicalThreadState Do not use.
	SwitchInLogicalThreadState uintptr
	// SwitchOutLogicalThreadState Do not use.
	SwitchOutLogicalThreadState uintptr
	// LocksHeldByLogicalThreadState Do not use.
	LocksHeldByLogicalThreadState uintptr
	// MapFile Maps the specified file into memory. This method is obsolete.
	MapFile uintptr
	// GetConfiguration Gets an object that allows the host to specify the callback configuration of the CLR.
	GetConfiguration uintptr
	// Start Starts the CLR.
	Start uintptr
	// Stop Stops the execution of code in the runtime for the current process.
	Stop uintptr
	// CreateDomain Creates an application domain. The caller receives an interface pointer of
	// type _AppDomain to an instance of type System.AppDomain.
	CreateDomain uintptr
	// GetDefaultDomain Gets an interface pointer of type _AppDomain that represents the default domain for the current process.
	GetDefaultDomain uintptr
	// EnumDomains Gets an enumerator for the domains in the current process.
	EnumDomains uintptr
	// NextDomain Gets an interface pointer to the next domain in the enumeration.
	NextDomain uintptr
	// CloseEnum Resets a domain enumerator back to the beginning of the domain list.
	CloseEnum uintptr
	// CreateDomainEx Creates an application domain. This method allows the caller to pass an
	// IAppDomainSetup instance to configure additional features of the returned _AppDomain instance.
	CreateDomainEx uintptr
	// CreateDomainSetup Gets an interface pointer of type IAppDomainSetup to an AppDomainSetup instance.
	// IAppDomainSetup provides methods to configure aspects of an application domain before it is created.
	CreateDomainSetup uintptr
	// CreateEvidence Gets an interface pointer of type IIdentity, which allows the host to create security
	// evidence to pass to CreateDomain or CreateDomainEx.
	CreateEvidence uintptr
	// UnloadDomain Unloads the specified application domain from the current process.
	UnloadDomain uintptr
	// CurrentDomain Gets an interface pointer of type _AppDomain that represents the domain loaded on the current thread.
	CurrentDomain uintptr
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
		debugPrint(fmt.Sprintf("the ICORRuntimeHost::Start method returned an error:\r\n%s", err))
	}
	if hr != S_OK {
		return fmt.Errorf("the ICORRuntimeHost::Start method method returned a non-zero HRESULT: 0x%x", hr)
	}
	return nil
}

// GetDefaultDomain gets an interface pointer of type System._AppDomain that represents the default domain for the current process.
// HRESULT GetDefaultDomain (
//   [out] IUnknown** pAppDomain
// );
// https://docs.microsoft.com/en-us/dotnet/framework/unmanaged-api/hosting/icorruntimehost-getdefaultdomain-method
func (obj *ICORRuntimeHost) GetDefaultDomain() (IUnknown *IUnknown, err error) {
	debugPrint("Entering into icorruntimehost.GetDefaultDomain()...")
	hr, _, err := syscall.Syscall(
		obj.vtbl.GetDefaultDomain,
		2,
		uintptr(unsafe.Pointer(obj)),
		uintptr(unsafe.Pointer(&IUnknown)),
		0,
	)
	if err != syscall.Errno(0) {
		// The specified procedure could not be found.
		// TODO Why is this error message returned?
		debugPrint(fmt.Sprintf("the ICORRuntimeHost::GetDefaultDomain method returned an error:\r\n%s", err))
	}
	if hr != S_OK {
		err = fmt.Errorf("the ICORRuntimeHost::GetDefaultDomain method method returned a non-zero HRESULT: 0x%x", hr)
		return
	}
	err = nil
	return
}
