// +build windows

package clr

import (
	"fmt"
	"syscall"
	"unsafe"

	"golang.org/x/sys/windows"
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
	return metahost.GetRuntime(pwzVersion, IID_ICLRRuntimeInfo)
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

// GetInterface loads the CLR into the current process and returns runtime interface pointers,
// such as ICLRRuntimeHost, ICLRStrongName, and IMetaDataDispenserEx.
// HRESULT GetInterface(
//   [in]  REFCLSID rclsid,
//   [in]  REFIID   riid,
//   [out, iid_is(riid), retval] LPVOID *ppUnk); unsafe pointer of a pointer to an object pointer
// https://docs.microsoft.com/en-us/dotnet/framework/unmanaged-api/hosting/iclrruntimeinfo-getinterface-method
func (obj *ICLRRuntimeInfo) GetInterface(rclsid windows.GUID, riid windows.GUID, ppUnk unsafe.Pointer) error {
	debugPrint("Entering into iclrruntimeinfo.GetInterface()...")
	hr, _, err := syscall.Syscall6(
		obj.vtbl.GetInterface,
		4,
		uintptr(unsafe.Pointer(obj)),
		uintptr(unsafe.Pointer(&rclsid)),
		uintptr(unsafe.Pointer(&riid)),
		uintptr(ppUnk),
		0,
		0,
	)
	// The syscall returns "The requested lookup key was not found in any active activation context." in the error position
	// TODO Why is this error message returned?
	if err.Error() != "The requested lookup key was not found in any active activation context." {
		return fmt.Errorf("the ICLRRuntimeInfo::GetInterface method returned an error:\r\n%s", err)
	}
	if hr != S_OK {
		return fmt.Errorf("the ICLRRuntimeInfo::GetInterface method returned a non-zero HRESULT: 0x%x", hr)
	}
	return nil
}

// BindAsLegacyV2Runtime binds the current runtime for all legacy common language runtime (CLR) version 2 activation policy decisions.
// HRESULT BindAsLegacyV2Runtime ();
// https://docs.microsoft.com/en-us/dotnet/framework/unmanaged-api/hosting/iclrruntimeinfo-bindaslegacyv2runtime-method
func (obj *ICLRRuntimeInfo) BindAsLegacyV2Runtime() error {
	debugPrint("Entering into iclrruntimeinfo.BindAsLegacyV2Runtime()...")
	hr, _, err := syscall.Syscall(
		obj.vtbl.BindAsLegacyV2Runtime,
		1,
		uintptr(unsafe.Pointer(obj)),
		0,
		0,
	)
	if err != syscall.Errno(0) {
		return fmt.Errorf("the ICLRRuntimeInfo::BindAsLegacyV2Runtime method returned an error:\r\n%s", err)
	}
	if hr != S_OK {
		return fmt.Errorf("the ICLRRuntimeInfo::BindAsLegacyV2Runtime method returned a non-zero HRESULT: 0x%x", hr)
	}
	return nil
}

// IsLoadable indicates whether the runtime associated with this interface can be loaded into the current process,
// taking into account other runtimes that might already be loaded into the process.
// HRESULT IsLoadable(
//   [out, retval] BOOL *pbLoadable);
// https://docs.microsoft.com/en-us/dotnet/framework/unmanaged-api/hosting/iclrruntimeinfo-isloadable-method
func (obj *ICLRRuntimeInfo) IsLoadable(pbLoadable *bool) error {
	debugPrint("Entering into iclrruntimeinfo.IsLoadable()...")
	hr, _, err := syscall.Syscall(
		obj.vtbl.IsLoadable,
		2,
		uintptr(unsafe.Pointer(obj)),
		uintptr(unsafe.Pointer(pbLoadable)),
		0)
	if err != syscall.Errno(0) {
		return fmt.Errorf("the ICLRRuntimeInfo::IsLoadable method returned an error:\r\n%s", err)
	}
	if hr != S_OK {
		return fmt.Errorf("the ICLRRuntimeInfo::IsLoadable method  returned a non-zero HRESULT: 0x%x", hr)
	}
	return nil
}
