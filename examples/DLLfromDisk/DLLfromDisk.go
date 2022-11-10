// +build windows

package main

import (
	"fmt"
	"log"
	"syscall"
	"unsafe"

	clr "github.com/Ne0nd0g/go-clr"
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
	metahost, err := clr.CLRCreateInstance(clr.CLSID_CLRMetaHost, clr.IID_ICLRMetaHost)
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

	isLoadable, err := runtimeInfo.IsLoadable()
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
	pDLLPath, err := syscall.UTF16PtrFromString("TestDLL.dll")
	must(err)
	pTypeName, err := syscall.UTF16PtrFromString("TestDLL.HelloWorld")
	must(err)
	pMethodName, err := syscall.UTF16PtrFromString("SayHello")
	must(err)
	pArgument, err := syscall.UTF16PtrFromString("foobar")
	must(err)
	ret, err := runtimeHost.ExecuteInDefaultAppDomain(
		pDLLPath,
		pTypeName,
		pMethodName,
		pArgument,
	)
	if *ret != 0 {
		err = fmt.Errorf("the ICLRRuntimeHost::ExecuteInDefaultAppDomain method returned a non-zero return value: %d", *ret)
	}
	must(err)
}
