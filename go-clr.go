// +build windows

// go-clr is a PoC package that wraps Windows syscalls necessary to load and the CLR into the current process and
// execute a managed DLL from disk or a managed EXE from memory
package clr

import (
	"fmt"
	"runtime"
	"strings"
	"syscall"
	"unsafe"
)

// GetInstallRuntimes is a wrapper function that returns an array of installed runtimes. Requires an existing ICLRMetaHost
func GetInstalledRuntimes(metahost *ICLRMetaHost) ([]string, error) {
	var runtimes []string
	var pInstalledRuntimes uintptr
	hr := metahost.EnumerateInstalledRuntimes(&pInstalledRuntimes)
	err := checkOK(hr, "EnumerateInstalledRuntimes")
	if err != nil {
		return runtimes, err
	}
	installedRuntimes := NewIEnumUnknownFromPtr(pInstalledRuntimes)
	var pRuntimeInfo uintptr
	var fetched = uint32(0)
	var versionString string
	versionStringBytes := make([]uint16, 20)
	versionStringSize := uint32(len(versionStringBytes))
	var runtimeInfo *ICLRRuntimeInfo
	for {
		hr = installedRuntimes.Next(1, &pRuntimeInfo, &fetched)
		if hr != S_OK {
			break
		}
		runtimeInfo = NewICLRRuntimeInfoFromPtr(pRuntimeInfo)
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
// and execute a DLL on disk in the default app domain. It takes in the DLLPath, TypeName, MethodName and Argument to use
// as strings. It returns the return code from the assembly
func ExecuteDLLFromDisk(dllpath, typeName, methodName, argument string) (retCode int16, err error) {
	retCode = -1
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
		if strings.Contains(r, "v4") {
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
	hr := runtimeInfo.IsLoadable(&isLoadable)
	err = checkOK(hr, "runtimeInfo.IsLoadable")
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
	hr = runtimeHost.ExecuteInDefaultAppDomain(pDLLPath, pTypeName, pMethodName, pArgument, &pReturnVal)
	err = checkOK(hr, "runtimeHost.ExecuteInDefaultAppDomain")
	if err != nil {
		return int16(pReturnVal), err
	}
	runtimeHost.Release()
	runtimeInfo.Release()
	metahost.Release()
	return int16(pReturnVal), nil

}

// ExecuteByteArray is a wrapper function that will automatically load the latest supported framework into the current
// process using the legacy APIs, then load and execute an executable from memory. It takes in a byte array of the
// executable to load and run and returns the return code. It currently does not support any arguments on the entry point
func ExecuteByteArray(rawBytes []byte) (retCode int32, err error) {
	retCode = -1
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
		if strings.Contains(r, "v4") {
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
	hr := runtimeInfo.IsLoadable(&isLoadable)
	err = checkOK(hr, "runtimeInfo.IsLoadable")
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
	safeArray, err := CreateSafeArray(rawBytes)
	if err != nil {
		return
	}
	runtime.KeepAlive(&safeArray)
	var pAssembly uintptr
	hr = appDomain.Load_3(uintptr(unsafe.Pointer(&safeArray)), &pAssembly)
	err = checkOK(hr, "appDomain.Load_3")
	if err != nil {
		return
	}
	assembly := NewAssemblyFromPtr(pAssembly)
	var pEntryPointInfo uintptr
	hr = assembly.GetEntryPoint(&pEntryPointInfo)
	err = checkOK(hr, "assembly.GetEntryPoint")
	if err != nil {
		return
	}
	methodInfo := NewMethodInfoFromPtr(pEntryPointInfo)
	var pRetCode uintptr
	nullVariant := Variant{
		VT:  1,
		Val: uintptr(0),
	}
	hr = methodInfo.Invoke_3(
		nullVariant,
		uintptr(0),
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
