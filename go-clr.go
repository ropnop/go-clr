// +build windows

package clr

import (
	"fmt"
	"strings"
	"syscall"
)

var (
	Metahost        *ICLRMetaHost
	RuntimeInfo     *ICLRRuntimeInfo
	IsLoadable      bool
	LegacyV2Runtime bool
	CorRuntimeHost  *ICORRuntimeHost
	ClrRuntimeHost  *ICLRRuntimeHost
	Iu              *IUnknown
	LoadedAppDomain *AppDomain
	LoadedAssembly  *Assembly
	EntryPoint      *MethodInfo
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
	hr := runtimeInfo.GetInterface(&CLSID_CorRuntimeHost, &IID_ICLRRuntimeHost, &pRuntimeHost)
	err := checkOK(hr, "runtimeInfo.GetInterface")
	if err != nil {
		return nil, err
	}
	return NewICLRRuntimeHost(pRuntimeHost), nil
}

// ExecuteDLL is a wrapper function that will automatically load the latest installed CLR into the current process
// and execute a DLL on disk in the default app domain
func ExecuteDLL(dllpath, typeName, methodName, argument string) (retCode int, err error) {
	retCode = -1
	metahost, err := GetMetaHost()
	if err != nil {
		return
	}
	defer metahost.Release()

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
	defer runtimeInfo.Release()
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
	defer runtimeHost.Release()
	hr = runtimeHost.Start()
	err = checkOK(hr, "runtimeHost.Start")
	if err != nil {
		return
	}

	pDLLPath, _ := syscall.UTF16PtrFromString(dllpath)
	pTypeName, _ := syscall.UTF16PtrFromString(typeName)
	pMethodName, _ := syscall.UTF16PtrFromString(methodName)
	pArgument, _ := syscall.UTF16PtrFromString(argument)
	var pReturnVal *uint16
	hr = runtimeHost.ExecuteInDefaultAppDomain(pDLLPath, pTypeName, pMethodName, pArgument, pReturnVal)
	err = checkOK(hr, "runtimeHost.ExecuteInDefaultAppDomain")
	if err != nil {
		return int(*pReturnVal), err
	}
	return int(*pReturnVal), nil

}
