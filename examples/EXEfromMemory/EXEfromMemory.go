// +build windows

package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"runtime"
	"strings"
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
		log.Fatalf("%s returned 0x%08x", caller, hr)
	}
}

func init() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: EXEfromMemory.exe <exe_file> <exe_args>")
		os.Exit(1)
	}
}

func main() {
	filename := os.Args[1]
	exebytes, err := ioutil.ReadFile(filename)
	must(err)
	runtime.KeepAlive(exebytes)

	var params []string
	if len(os.Args) > 2 {
		params = os.Args[2:]
	}

	metaHost, err := clr.CLRCreateInstance(clr.CLSID_CLRMetaHost, clr.IID_ICLRMetaHost)
	must(err)

	versionString := "v4.0.30319"
	pwzVersion, err := syscall.UTF16PtrFromString(versionString)
	must(err)
	runtimeInfo, err := metaHost.GetRuntime(pwzVersion, clr.IID_ICLRRuntimeInfo)
	must(err)

	isLoadable, err := runtimeInfo.IsLoadable()
	must(err)
	if !isLoadable {
		log.Fatal("[!] IsLoadable returned false. Bailing...")
	}

	err = runtimeInfo.BindAsLegacyV2Runtime()
	must(err)

	var runtimeHost *clr.ICORRuntimeHost
	err = runtimeInfo.GetInterface(clr.CLSID_CorRuntimeHost, clr.IID_ICorRuntimeHost, unsafe.Pointer(&runtimeHost))
	must(err)
	err = runtimeHost.Start()
	must(err)
	fmt.Println("[+] Loaded CLR into this process")

	iu, err := runtimeHost.GetDefaultDomain()
	must(err)

	var appDomain *clr.AppDomain
	err = iu.QueryInterface(clr.IID_AppDomain, unsafe.Pointer(&appDomain))
	must(err)
	fmt.Println("[+] Got default AppDomain")

	safeArray, err := clr.CreateSafeArray(exebytes)
	must(err)
	runtime.KeepAlive(safeArray)
	fmt.Println("[+] Created SafeArray from byte array")

	assembly, err := appDomain.Load_3(safeArray)
	must(err)
	fmt.Printf("[+] Loaded %d bytes into memory from %s\n", len(exebytes), filename)
	fmt.Printf("[+] Executable loaded into memory at %p\n", assembly)

	methodInfo, err := assembly.GetEntryPoint()
	must(err)
	fmt.Printf("[+] Executable entrypoint found at 0x%x\n", uintptr(unsafe.Pointer(methodInfo)))

	var paramSafeArray *clr.SafeArray
	methodSignature, err := methodInfo.GetString()
	if err != nil {
		return
	}

	fmt.Println("[+] Checking if the assembly requires arguments...")
	if !strings.Contains(methodSignature, "Void Main()") {
		if len(params) < 1 {
			log.Fatal("the assembly requires arguments but none were provided\nUsage: EXEfromMemory.exe <exe_file> <exe_args>")
		}
		if paramSafeArray, err = clr.PrepareParameters(params); err != nil {
			log.Fatal(fmt.Sprintf("there was an error preparing the assembly arguments:\r\n%s", err))
		}
	}

	nullVariant := clr.Variant{
		VT:  1,
		Val: uintptr(0),
	}
	fmt.Println("[+] Invoking...")
	err = methodInfo.Invoke_3(nullVariant, paramSafeArray)
	must(err)

	appDomain.Release()
	runtimeHost.Release()
	runtimeInfo.Release()
	metaHost.Release()
}
