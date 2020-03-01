//+build windows
package main

import (
	"fmt"
	"github.com/Microsoft/go-winio/pkg/guid"
	"log"
	"syscall"
	"unsafe"
)

const S_OK = 0
const E_POINTER = 0x80004003

var (
	CLSID_CLRMetaHost = guid.GUID{0x9280188d,0xe8e,  0x4867, [8]byte{0xb3, 0xc, 0x7f, 0xa8, 0x38, 0x84, 0xe8, 0xde}}
	IID_ICLRMetaHost = guid.GUID{0xD332DB9E, 0xB9B3, 0x4125, [8]byte{0x82, 0x07, 0xA1, 0x48, 0x84, 0xF5, 0x32, 0x16}}
	IID_ICLRRuntimeInfo = guid.GUID{0xBD39D1D2, 0xBA2F, 0x486a, [8]byte{0x89, 0xB0, 0xB4, 0xB0, 0xCB, 0x46, 0x68, 0x91}}

	//EXTERN_GUID(CLSID_CLRRuntimeHost, 0x90F1A06E, 0x7712, 0x4762, 0x86, 0xB5, 0x7A, 0x5E, 0xBA, 0x6B, 0xDB, 0x02);
	CLSID_CLRRuntimeHost = guid.GUID{ 0x90F1A06E, 0x7712, 0x4762, [8]byte{0x86, 0xB5, 0x7A, 0x5E, 0xBA, 0x6B, 0xDB, 0x02}}

	//EXTERN_GUID(IID_ICLRRuntimeHost, 0x90F1A06C, 0x7712, 0x4762, 0x86, 0xB5, 0x7A, 0x5E, 0xBA, 0x6B, 0xDB, 0x02);
	IID_ICLRRuntimeHost = guid.GUID{0x90F1A06C, 0x7712, 0x4762, [8]byte{0x86, 0xB5, 0x7A, 0x5E, 0xBA, 0x6B, 0xDB, 0x02}}

	)

func must(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func checkOK(hr uintptr) {
	if hr == S_OK {
		return
	} else if hr == E_POINTER {
		log.Fatalf("Error getting pointer: 0x%x\n", hr)
	} else {
		log.Fatalf("Unknown return result: 0x%x\n", hr)
	}
}

func main() {
	metaHost, err := GetICLRMetaHost()
	must(err)
	var pInstalledRuntimes uintptr

	hr := metaHost.EnumerateInstalledRuntimes(&pInstalledRuntimes)

	checkOK(hr)
	installedRuntimes := newIEnumUnknown(pInstalledRuntimes)

	var pRuntimeInfo uintptr
	var fetched = uint32(0)
	var versionString string
	versionStringBytes := make([]uint16, 20)
	versionStringSize := uint32(20)
	var runtimeInfo *ICLRRuntimeInfo
	fmt.Println("[+] Enumerating installed runtimes...")
	for {
		hr = installedRuntimes.Next(1, &pRuntimeInfo, &fetched )
		if hr != S_OK {break}
		runtimeInfo = newICLRRuntimeInfo(pRuntimeInfo)
		if ret := runtimeInfo.GetVersionString(&versionStringBytes[0], &versionStringSize); ret != S_OK {
			log.Fatalf("[+] Error getting version string: 0x%x", ret)
		}
		versionString = syscall.UTF16ToString(versionStringBytes)
		fmt.Printf("\t[+] Found: %s\n", versionString)
	}
	fmt.Printf("[+] Using latest supported framework: %s\n", versionString)


	var isLoadable bool
	hr = runtimeInfo.IsLoadable(&isLoadable)
	checkOK(hr)
	fmt.Printf("isLoadable: %t\n", isLoadable)

	if (!isLoadable) {
		log.Fatal("[!] IsLoadable returned false. Won't load CLR")
	}

	fmt.Println("[+] BindAsLegacyV2Runtime...")
	hr = runtimeInfo.BindAsLegacyV2Runtime()

	var pRuntimeHost uintptr
	hr = runtimeInfo.GetInterface(&CLSID_CLRRuntimeHost, &IID_ICLRRuntimeHost, &pRuntimeHost)
	checkOK(hr)
	runTimeHost := newICLRRuntimeHost(pRuntimeHost)
	hr = runTimeHost.Start()
	checkOK(hr)
	fmt.Println("[+] Loaded CLR into this process")

	var appDomainID uint16
	runTimeHost.GetCurrentappDomainID(&appDomainID)
	checkOK(hr)
	fmt.Printf("[+] Current AppDomain ID: %d\n", appDomainID)

	fmt.Println("[+] Executing assembly...")
	pDLLPath, _ := syscall.UTF16PtrFromString("TestDLL.dll")
	pTypeName, _ := syscall.UTF16PtrFromString("TestDLL.HelloWorld")
	pMethodName, _ := syscall.UTF16PtrFromString("SayHello")
	pArgument, _ := syscall.UTF16PtrFromString("foobar")
	var pReturnVal *uint16
	hr = runTimeHost.ExecuteInDefaultAppDomain(
		pDLLPath,
		pTypeName,
		pMethodName,
		pArgument,
		pReturnVal)

	checkOK(hr)
	fmt.Printf("[+] Assembly returned: 0x%x\n", pReturnVal)





	//if hr == S_OK {
	//	log.Println("[+] Enumerating runtimes")
	//} else if hr == E_POINTER {
	//	log.Fatal("[!] Couldn't get pointer to enumerate runtimes")
	//} else {
	//	log.Fatalf("[!] Unknown error: %x\n", hr)
	//}
	//
	//
	//
	//var pRuntimeInfo uintptr
	//
	//hr = metaHost.GetRuntime("v4.0.30319", pRuntimeInfo)
	//if hr == S_OK {
	//	log.Println("[+] Got Runtime!")
	//} else if hr == E_POINTER {
	//	log.Fatalf("[!] Error getting runtime pointer: %x\n", hr)
	//} else {
	//	log.Fatalf("[!] Unknown error: %x\n", hr)
	//}


}





func printRawData(ptr uintptr, size uintptr) {
	fmt.Printf("Printing ptr %016x size %d\n", ptr, size)
	i := ptr
	var offset uintptr
	for i < ptr + size {
		if offset % 16 == 0 {
			fmt.Printf("%016x : ", i)
		}

		fmt.Printf("%02x", *(*byte)(unsafe.Pointer(i)))

		i++
		offset++
		if offset % 16 == 0 || offset == size {
			fmt.Print("\n")
		} else if offset % 8 == 0 {
			fmt.Print("  ")
		} else {
			fmt.Print(" ")
		}
	}
}
