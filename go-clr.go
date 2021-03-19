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
	enumICLRRuntimeInfo, err := metahost.EnumerateInstalledRuntimes()
	if err != nil {
		return runtimes, err
	}

	var hr int
	for hr != S_FALSE {
		var runtimeInfo *ICLRRuntimeInfo
		var fetched = uint32(0)
		hr, err = enumICLRRuntimeInfo.Next(1, unsafe.Pointer(&runtimeInfo), &fetched)
		if err != nil {
			return runtimes, fmt.Errorf("InstalledRuntimes Next Error:\r\n%s\n", err)
		}
		if hr == S_FALSE {
			break
		}
		// Only release if an interface pointer was returned
		runtimeInfo.Release()

		version, err := runtimeInfo.GetVersionString()
		if err != nil {
			return runtimes, err
		}
		runtimes = append(runtimes, version)
	}
	if len(runtimes) == 0 {
		return runtimes, fmt.Errorf("Could not find any installed runtimes")
	}
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
	metahost, err := CLRCreateInstance(CLSID_CLRMetaHost, IID_ICLRMetaHost)
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

	isLoadable, err := runtimeInfo.IsLoadable()
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

	pDLLPath, err := syscall.UTF16PtrFromString(dllpath)
	must(err)
	pTypeName, err := syscall.UTF16PtrFromString(typeName)
	must(err)
	pMethodName, err := syscall.UTF16PtrFromString(methodName)
	must(err)
	pArgument, err := syscall.UTF16PtrFromString(argument)
	must(err)

	ret, err := runtimeHost.ExecuteInDefaultAppDomain(pDLLPath, pTypeName, pMethodName, pArgument)
	must(err)
	if *ret != 0 {
		return int16(*ret), fmt.Errorf("the ICLRRuntimeHost::ExecuteInDefaultAppDomain method returned a non-zero return value: %d", *ret)
	}

	runtimeHost.Release()
	runtimeInfo.Release()
	metahost.Release()
	return 0, nil
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
	metahost, err := CLRCreateInstance(CLSID_CLRMetaHost, IID_ICLRMetaHost)
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

	isLoadable, err := runtimeInfo.IsLoadable()
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

	var paramSafeArray *SafeArray
	methodSignature, err := methodInfo.GetString()
	if err != nil {
		return
	}

	if expectsParams(methodSignature) {
		if paramSafeArray, err = PrepareParameters(params); err != nil {
			return
		}
	}

	nullVariant := Variant{
		VT:  1,
		Val: uintptr(0),
	}
	err = methodInfo.Invoke_3(nullVariant, paramSafeArray)
	if err != nil {
		return
	}
	appDomain.Release()
	runtimeHost.Release()
	runtimeInfo.Release()
	metahost.Release()
	return 0, nil
}
