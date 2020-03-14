// +build windows

package main

import (
	"fmt"
	"log"
	"syscall"

	clr "github.com/ropnop/go-clr"
)

func must(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func checkOK(hr uintptr, caller string) {
	if hr != 0x0 {
		log.Fatalf("%s returned 0x%08x", hr)
	}
}

func main() {
	metahost, err := clr.GetICLRMetaHost()
	must(err)
	fmt.Println("[+] Got metahost")

	installedRuntimes, err := clr.GetInstalledRuntimes(metahost)
	must(err)
	fmt.Printf("[+] Found installed runtimes: %s\n", installedRuntimes)
	versionString := "v4.0.30319"
	pwzVersion, _ := syscall.UTF16PtrFromString(versionString)
	var pRuntimeInfo uintptr

	hr := metahost.GetRuntime(pwzVersion, &clr.IID_ICLRRuntimeInfo, &pRuntimeInfo)
	checkOK(hr, "metaHost.GetRuntime")
	runtimeInfo := clr.NewICLRRuntimeInfo(pRuntimeInfo)
	fmt.Printf("[+] Using runtime: %s\n", versionString)

	var isLoadable bool
	hr = runtimeInfo.IsLoadable(&isLoadable)
	checkOK(hr, "runtimeInfo.IsLoadable")
	if !isLoadable {
		log.Fatal("[!] IsLoadable returned false. Bailing...")
	}

	var pRuntimeHost uintptr
	hr = runtimeInfo.GetInterface(&clr.CLSID_CLRRuntimeHost, &clr.IID_ICLRRuntimeHost, &pRuntimeHost)
	checkOK(hr, "runtimeInfo.GetInterface")
	runtimeHost := clr.NewICLRRuntimeHost(pRuntimeHost)

	hr = runtimeHost.Start()
	checkOK(hr, "runtimeHost.Start")
	fmt.Println("[+] Loaded CLR into this process")

	fmt.Println("[+] Executing assembly...")
	pDLLPath, _ := syscall.UTF16PtrFromString("TestDLL.dll")
	pTypeName, _ := syscall.UTF16PtrFromString("TestDLL.HelloWorld")
	pMethodName, _ := syscall.UTF16PtrFromString("SayHello")
	pArgument, _ := syscall.UTF16PtrFromString("foobar")
	var pReturnVal *uint16
	hr = runtimeHost.ExecuteInDefaultAppDomain(
		pDLLPath,
		pTypeName,
		pMethodName,
		pArgument,
		pReturnVal)

	checkOK(hr, "runtimeHost.ExecuteInDefaultAppDomain")
	fmt.Printf("[+] Assembly returned: 0x%x\n", pReturnVal)
}
