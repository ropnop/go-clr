// +build windows

package main

import (
	"fmt"
	"log"
	"syscall"
	"unsafe"

	clr "github.com/ropnop/go-clr"
)

func must(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func checkOK(hr uintptr, caller string) {
	if hr != 0x0 {
		log.Fatalf("%s returned 0x%x", caller, hr)
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
	pwzVersion, err := syscall.UTF16PtrFromString(versionString)
	must(err)

	runtimeInfo, err := metahost.GetRuntime(pwzVersion, clr.IID_ICLRRuntimeInfo)
	must(err)
	fmt.Printf("[+] Using runtime: %s\n", versionString)

	var isLoadable bool
	err = runtimeInfo.IsLoadable(&isLoadable)
	must(err)
	if !isLoadable {
		log.Fatal("[!] IsLoadable returned false. Bailing...")
	}

	var runtimeHost *clr.ICLRRuntimeHost
	err = runtimeInfo.GetInterface(clr.CLSID_CLRRuntimeHost, clr.IID_ICLRRuntimeHost, unsafe.Pointer(&runtimeHost))
	must(err)

	err = runtimeHost.Start()
	must(err)
	fmt.Println("[+] Loaded CLR into this process")

	fmt.Println("[+] Executing assembly...")
	pDLLPath, _ := syscall.UTF16PtrFromString("TestDLL.dll")
	pTypeName, _ := syscall.UTF16PtrFromString("TestDLL.HelloWorld")
	pMethodName, _ := syscall.UTF16PtrFromString("SayHello")
	pArgument, _ := syscall.UTF16PtrFromString("foobar")
	var pReturnVal *uint16
	hr := runtimeHost.ExecuteInDefaultAppDomain(
		pDLLPath,
		pTypeName,
		pMethodName,
		pArgument,
		pReturnVal)

	checkOK(hr, "runtimeHost.ExecuteInDefaultAppDomain")
	fmt.Printf("[+] Assembly returned: 0x%x\n", pReturnVal)
}
