// +build windows

// go-clr is a PoC package that wraps Windows syscalls necessary to load and the CLR into the current process and
// execute a managed DLL from disk or a managed EXE from memory
package clr

import (
	"fmt"
	"strings"
	"syscall"
	"unsafe"
)

// GetInstallRuntimes is a wrapper function that returns an array of installed runtimes. Requires an existing ICLRMetaHost
func GetInstalledRuntimes(metahost *ICLRMetaHost) ([]string, error) {
	var runtimes []string
	installedRuntimes, err := metahost.EnumerateInstalledRuntimes()
	if err != nil {
		return runtimes, err
	}

	var fetched = uint32(0)
	var versionString string
	versionStringBytes := make([]uint16, 20)
	versionStringSize := uint32(len(versionStringBytes))
	var runtimeInfo *ICLRRuntimeInfo
	for {
		err := installedRuntimes.Next(1, unsafe.Pointer(runtimeInfo), &fetched)
		if err != nil {
			break
		}
		if ret := runtimeInfo.GetVersionString(&versionStringBytes[0], &versionStringSize); ret != S_OK {
			return runtimes, fmt.Errorf("GetVersionString returned 0x%08x", ret)
		}
		versionString = syscall.UTF16ToString(versionStringBytes)
		runtimes = append(runtimes, versionString)
	}
	if len(runtimes) == 0 {
		return runtimes, fmt.Errorf("Could not find any installed runtimes")
	}
	runtimeInfo.Release()
	return runtimes, err
}

// ExecuteDLLFromDisk is a wrapper function that will automatically load the latest installed CLR into the current process
// and execute a DLL on disk in the default app domain. It takes in the target runtime, DLLPath, TypeName, MethodName
// and Argument to use as strings. It returns the return code from the assembly
func ExecuteDLLFromDisk(targetRuntime, dllpath, typeName, methodName, argument string) (retCode int16, err error) {
	retCode = -1
	if targetRuntime == "" {
		targetRuntime = "v4"
	}
	metahost, err := GetICLRMetaHost()
	if err != nil {
		return
	}

	runtimes, err := GetInstalledRuntimes(metahost)
	if err != nil {
		return
	}
	var latestRuntime string
	for _, r := range runtimes {
		if strings.Contains(r, targetRuntime) {
			latestRuntime = r
			break
		} else {
			latestRuntime = r
		}
	}
	runtimeInfo, err := GetRuntimeInfo(metahost, latestRuntime)
	if err != nil {
		return
	}
	var isLoadable bool
	err = runtimeInfo.IsLoadable(&isLoadable)
	if err != nil {
		return
	}
	if !isLoadable {
		return -1, fmt.Errorf("%s is not loadable for some reason", latestRuntime)
	}
	runtimeHost, err := GetICLRRuntimeHost(runtimeInfo)
	if err != nil {
		return
	}

	pDLLPath, _ := syscall.UTF16PtrFromString(dllpath)
	pTypeName, _ := syscall.UTF16PtrFromString(typeName)
	pMethodName, _ := syscall.UTF16PtrFromString(methodName)
	pArgument, _ := syscall.UTF16PtrFromString(argument)
	var pReturnVal uint16
	hr := runtimeHost.ExecuteInDefaultAppDomain(pDLLPath, pTypeName, pMethodName, pArgument, &pReturnVal)
	err = checkOK(hr, "runtimeHost.ExecuteInDefaultAppDomain")
	if err != nil {
		return int16(pReturnVal), err
	}
	runtimeHost.Release()
	runtimeInfo.Release()
	metahost.Release()
	return int16(pReturnVal), nil

}

// ExecuteByteArray is a wrapper function that will automatically loads the supplied target framework into the current
// process using the legacy APIs, then load and execute an executable from memory. If no targetRuntime is specified, it
// will default to latest. It takes in a byte array of the executable to load and run and returns the return code.
// You can supply an array of strings as command line arguments.
func ExecuteByteArray(targetRuntime string, rawBytes []byte, params []string) (retCode int32, err error) {
	retCode = -1
	if targetRuntime == "" {
		targetRuntime = "v4"
	}
	metahost, err := GetICLRMetaHost()
	if err != nil {
		return
	}

	runtimes, err := GetInstalledRuntimes(metahost)
	if err != nil {
		return
	}
	var latestRuntime string
	for _, r := range runtimes {
		if strings.Contains(r, targetRuntime) {
			latestRuntime = r
			break
		} else {
			latestRuntime = r
		}
	}
	runtimeInfo, err := GetRuntimeInfo(metahost, latestRuntime)
	if err != nil {
		return
	}
	var isLoadable bool
	err = runtimeInfo.IsLoadable(&isLoadable)
	if err != nil {
		return
	}
	if !isLoadable {
		return -1, fmt.Errorf("%s is not loadable for some reason", latestRuntime)
	}
	runtimeHost, err := GetICORRuntimeHost(runtimeInfo)
	if err != nil {
		return
	}
	appDomain, err := GetAppDomain(runtimeHost)
	if err != nil {
		return
	}
	safeArrayPtr, err := CreateSafeArray(rawBytes)
	if err != nil {
		return
	}

	assembly, err := appDomain.Load_3(safeArrayPtr)
	if err != nil {
		return
	}

	methodInfo, err := assembly.GetEntryPoint()
	if err != nil {
		return
	}

	var paramPtr uintptr
	methodSignature, err := methodInfo.GetString()
	if err != nil {
		return
	}

	if expectsParams(methodSignature) {
		if paramPtr, err = PrepareParameters(params); err != nil {
			return
		}
	}

	var pRetCode uintptr
	nullVariant := Variant{
		VT:  1,
		Val: uintptr(0),
	}
	hr := methodInfo.Invoke_3(
		nullVariant,
		paramPtr,
		&pRetCode)
	err = checkOK(hr, "methodInfo.Invoke_3")
	if err != nil {
		return
	}
	appDomain.Release()
	runtimeHost.Release()
	runtimeInfo.Release()
	metahost.Release()
	return int32(pRetCode), nil

}

// PrepareParameters creates a safe array of strings (arguments) nested inside a Variant object, which is itself
// appended to the final safe array
func PrepareParameters(params []string) (uintptr, error) {
	sab := SafeArrayBound{
		cElements: uint32(len(params)),
		lLbound:   0,
	}
	listStrSafeArrayPtr, err := SafeArrayCreate(VT_BSTR, 1, &sab) // VT_BSTR
	if err != nil {
		return 0, err
	}
	for i, p := range params {
		bstr, _ := SysAllocString(p)
		SafeArrayPutElement(listStrSafeArrayPtr, int32(i), bstr)
	}

	paramVariant := Variant{
		VT:  VT_BSTR | VT_ARRAY, // VT_BSTR | VT_ARRAY
		Val: uintptr(unsafe.Pointer(listStrSafeArrayPtr)),
	}

	sab2 := SafeArrayBound{
		cElements: uint32(1),
		lLbound:   0,
	}
	paramsSafeArrayPtr, err := SafeArrayCreate(VT_VARIANT, 1, &sab2) // VT_VARIANT
	if err != nil {
		return 0, err
	}
	err = SafeArrayPutElement(paramsSafeArrayPtr, int32(0), unsafe.Pointer(&paramVariant))
	if err != nil {
		return 0, err
	}
	return uintptr(unsafe.Pointer(paramsSafeArrayPtr)), nil
}
