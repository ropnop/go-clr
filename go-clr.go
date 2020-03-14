// +build windows

package clr

import (
	"fmt"
	"runtime"
	"strings"
	"syscall"
	"unsafe"
)

func GetMetaHost() (*ICLRMetaHost, error) {
	var pMetaHost uintptr
	hr := CLRCreateInstance(&CLSID_CLRMetaHost, &IID_ICLRMetaHost, &pMetaHost)
	err := checkOK(hr, "CLRCreateInstance")
	if err != nil {
		return nil, err
	}
	return NewICLRMetaHost(pMetaHost), nil
}

func GetInstalledRuntimes(metahost *ICLRMetaHost) ([]string, error) {
	var runtimes []string
	var pInstalledRuntimes uintptr
	hr := metahost.EnumerateInstalledRuntimes(&pInstalledRuntimes)
	err := checkOK(hr, "EnumerateInstalledRuntimes")
	if err != nil {
		return runtimes, err
	}
	installedRuntimes := NewIEnumUnknown(pInstalledRuntimes)
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
		runtimeInfo = NewICLRRuntimeInfo(pRuntimeInfo)
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

func GetRuntimeInfo(metahost *ICLRMetaHost, version string) (*ICLRRuntimeInfo, error) {
	pwzVersion, err := syscall.UTF16PtrFromString(version)
	if err != nil {
		return nil, err
	}
	var pRuntimeInfo uintptr
	hr := metahost.GetRuntime(pwzVersion, &IID_ICLRRuntimeInfo, &pRuntimeInfo)
	err = checkOK(hr, "metahost.GetRuntime")
	if err != nil {
		return nil, err
	}
	return NewICLRRuntimeInfo(pRuntimeInfo), nil
}

func GetICLRRuntimeHost(runtimeInfo *ICLRRuntimeInfo) (*ICLRRuntimeHost, error) {
	var pRuntimeHost uintptr
	hr := runtimeInfo.GetInterface(&CLSID_CLRRuntimeHost, &IID_ICLRRuntimeHost, &pRuntimeHost)
	err := checkOK(hr, "runtimeInfo.GetInterface")
	if err != nil {
		return nil, err
	}
	runtimeHost := NewICLRRuntimeHost(pRuntimeHost)
	hr = runtimeHost.Start()
	err = checkOK(hr, "runtimeHost.Start")
	return runtimeHost, err
}

func GetICORRuntimeHost(runtimeInfo *ICLRRuntimeInfo) (*ICORRuntimeHost, error) {
	var pRuntimeHost uintptr
	hr := runtimeInfo.GetInterface(&CLSID_CorRuntimeHost, &IID_ICorRuntimeHost, &pRuntimeHost)
	err := checkOK(hr, "runtimeInfo.GetInterface")
	if err != nil {
		return nil, err
	}
	runtimeHost := NewICORRuntimeHost(pRuntimeHost)
	hr = runtimeHost.Start()
	err = checkOK(hr, "runtimeHost.Start")
	return runtimeHost, err
}

// ExecuteDLL is a wrapper function that will automatically load the latest installed CLR into the current process
// and execute a DLL on disk in the default app domain
func ExecuteDLL(dllpath, typeName, methodName, argument string) (retCode int16, err error) {
	retCode = -1
	metahost, err := GetMetaHost()
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

func ExecuteByteArray(rawBytes []byte) (retCode int32, err error) {
	retCode = -1
	metahost, err := GetMetaHost()
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
	assembly := NewAssembly(pAssembly)
	var pEntryPointInfo uintptr
	hr = assembly.GetEntryPoint(&pEntryPointInfo)
	err = checkOK(hr, "assembly.GetEntryPoint")
	if err != nil {
		return
	}
	methodInfo := NewMethodInfo(pEntryPointInfo)
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
