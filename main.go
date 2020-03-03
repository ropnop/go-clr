//+build windows
package main

import (
	"fmt"
	"github.com/Microsoft/go-winio/pkg/guid"
	"io/ioutil"
	"log"
	"syscall"
	"unsafe"
	//ole "github.com/go-ole/go-ole"
	//"github.com/go-ole/go-ole/oleutil"
)

const S_OK = 0
const E_POINTER = 0x80004003

var (
	CLSID_CLRMetaHost   = guid.GUID{0x9280188d, 0xe8e, 0x4867, [8]byte{0xb3, 0xc, 0x7f, 0xa8, 0x38, 0x84, 0xe8, 0xde}}
	IID_ICLRMetaHost    = guid.GUID{0xD332DB9E, 0xB9B3, 0x4125, [8]byte{0x82, 0x07, 0xA1, 0x48, 0x84, 0xF5, 0x32, 0x16}}
	IID_ICLRRuntimeInfo = guid.GUID{0xBD39D1D2, 0xBA2F, 0x486a, [8]byte{0x89, 0xB0, 0xB4, 0xB0, 0xCB, 0x46, 0x68, 0x91}}

	//EXTERN_GUID(CLSID_CLRRuntimeHost, 0x90F1A06E, 0x7712, 0x4762, 0x86, 0xB5, 0x7A, 0x5E, 0xBA, 0x6B, 0xDB, 0x02);
	CLSID_CLRRuntimeHost = guid.GUID{0x90F1A06E, 0x7712, 0x4762, [8]byte{0x86, 0xB5, 0x7A, 0x5E, 0xBA, 0x6B, 0xDB, 0x02}}

	//EXTERN_GUID(IID_ICLRRuntimeHost, 0x90F1A06C, 0x7712, 0x4762, 0x86, 0xB5, 0x7A, 0x5E, 0xBA, 0x6B, 0xDB, 0x02);
	IID_ICLRRuntimeHost = guid.GUID{0x90F1A06C, 0x7712, 0x4762, [8]byte{0x86, 0xB5, 0x7A, 0x5E, 0xBA, 0x6B, 0xDB, 0x02}}

	//EXTERN_GUID(IID_ICorRuntimeHost, 0xcb2f6722, 0xab3a, 0x11d2, 0x9c, 0x40, 0x00, 0xc0, 0x4f, 0xa3, 0x0a, 0x3e);
	IID_ICorRuntimeHost = guid.GUID{0xcb2f6722, 0xab3a, 0x11d2, [8]byte{0x9c, 0x40, 0x00, 0xc0, 0x4f, 0xa3, 0x0a, 0x3e}}

	//EXTERN_GUID(CLSID_CorRuntimeHost, 0xcb2f6723, 0xab3a, 0x11d2, 0x9c, 0x40, 0x00, 0xc0, 0x4f, 0xa3, 0x0a, 0x3e);
	CLSID_CorRuntimeHost = guid.GUID{0xcb2f6723, 0xab3a, 0x11d2, [8]byte{0x9c, 0x40, 0x00, 0xc0, 0x4f, 0xa3, 0x0a, 0x3e}}

	//https://docs.microsoft.com/en-us/dotnet/api/system._appdomain?view=netframework-4.8
	// _AppDomain Interface GUID
	IID_AppDomain = guid.GUID{0x5f696dc, 0x2b29, 0x3663, [8]uint8{0xad, 0x8b, 0xc4, 0x38, 0x9c, 0xf2, 0xa7, 0x13}}
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
	var pMetaHost uintptr
	var metaHost *ICLRMetaHost
	hr := CLRCreateInstance(&CLSID_CLRMetaHost, &IID_ICLRMetaHost, &pMetaHost)
	checkOK(hr)
	metaHost = newICLRMetaHost(pMetaHost)

	//metaHost, err := GetICLRMetaHost()
	//must(err)
	//var pInstalledRuntimes uintptr
	//
	//hr = metaHost.EnumerateInstalledRuntimes(&pInstalledRuntimes)
	//
	//checkOK(hr)
	//installedRuntimes := newIEnumUnknown(pInstalledRuntimes)
	//
	//var pRuntimeInfo uintptr
	//var fetched = uint32(0)
	//var versionString string
	//versionStringBytes := make([]uint16, 20)
	//versionStringSize := uint32(20)
	//var runtimeInfo *ICLRRuntimeInfo
	//fmt.Println("[+] Enumerating installed runtimes...")
	//for {
	//	hr = installedRuntimes.Next(1, &pRuntimeInfo, &fetched )
	//	if hr != S_OK {break}
	//	runtimeInfo = newICLRRuntimeInfo(pRuntimeInfo)
	//	if ret := runtimeInfo.GetVersionString(&versionStringBytes[0], &versionStringSize); ret != S_OK {
	//		log.Fatalf("[+] Error getting version string: 0x%x", ret)
	//	}
	//	versionString = syscall.UTF16ToString(versionStringBytes)
	//	fmt.Printf("\t[+] Found: %s\n", versionString)
	//}
	//fmt.Printf("[+] Using latest supported framework: %s\n", versionString)

	//So I'll use the older API for now
	var pRuntimeInfo uintptr
	var runtimeInfo *ICLRRuntimeInfo
	versionString := "v4.0.30319"
	pwzVersion, _ := syscall.UTF16PtrFromString(versionString)
	hr = metaHost.GetRuntime(pwzVersion, &IID_ICLRRuntimeInfo, &pRuntimeInfo)
	checkOK(hr)
	runtimeInfo = newICLRRuntimeInfo(pRuntimeInfo)
	fmt.Println("[+] Got RuntimeHost")
	var isLoadable bool
	hr = runtimeInfo.IsLoadable(&isLoadable)
	checkOK(hr)
	fmt.Printf("[+] isLoadable: %t\n", isLoadable)

	if !isLoadable {
		log.Fatal("[!] IsLoadable returned false. Won't load CLR")
	}

	fmt.Println("[+] BindAsLegacyV2Runtime...")
	hr = runtimeInfo.BindAsLegacyV2Runtime()
	checkOK(hr)

	var pRuntimeHost uintptr
	// This is the "new" API, but I can't figure out how to load assemblies from memory in unmanaged code :/
	//hr = runtimeInfo.GetInterface(&CLSID_CLRRuntimeHost, &IID_ICLRRuntimeHost, &pRuntimeHost)

	// So I'll use the older API for now
	//pwzVersion, _ := syscall.UTF16PtrFromString(versionString)
	//hr = metaHost.GetRuntime(pwzVersion, &IID_ICLRRuntimeInfo, &pRuntimeInfo)
	//checkOK(hr)
	//runtimeInfo = newICLRRuntimeInfo(pRuntimeInfo)
	fmt.Println("[+] Getting interface...")
	hr = runtimeInfo.GetInterface(&CLSID_CorRuntimeHost, &IID_ICorRuntimeHost, &pRuntimeHost)
	checkOK(hr)
	runTimeHost := newICORRuntimeHost(pRuntimeHost)
	fmt.Println("[+] Got interface. Loading CLR...")
	hr = runTimeHost.Start()
	checkOK(hr)
	fmt.Println("[+] Loaded CLR into this process")

	fmt.Println("[+] Getting Default AppDomain")

	//ole.CoInitialize(0)
	//var pIUnknown uintptr
	//https://docs.microsoft.com/en-us/dotnet/api/system._appdomain?view=netframework-4.8
	//unknown, err := oleutil.CreateObject("{05F696DC-2B29-3663-AD8B-C4389CF2A713}")
	//if err != nil {
	//	fmt.Printf("createobject failed: %s", err)
	//}
	//must(err)
	//iid, _ := oleutil.ClassIDFrom("{05F696DC-2B29-3663-AD8B-C4389CF2A713}")
	//oleappdomain, _ := unknown.QueryInterface(ole.IID_IDispatch)

	var pAppDomain uintptr
	var pIUnknown uintptr
	hr = runTimeHost.GetDefaultDomain(&pIUnknown)
	checkOK(hr)
	iu := newIUnknown(pIUnknown)
	//appDomainGuid, err := guid.FromString("05F696DC-2B29-3663-AD8B-C4389CF2A713")
	//must(err)
	hr = iu.QueryInterface(&IID_AppDomain, &pAppDomain)
	appDomain := newAppDomain(pAppDomain)
	fmt.Println("[+] Got default appdomain")

	//box := packr.New("EXEs", "./static")
	//testEXEBytes, err := box.Find("TestEXE.exe")
	//must(err)
	testEXEBytes, err := ioutil.ReadFile("./static/TestEXE.exe")
	must(err)
	//testEXEBytes := []byte{0x41, 0x41, 0x41, 0x41, 0x41, 0x41, 0x41, 0x41}
	fmt.Printf("[+] Loaded %d bytes from TestExe.exe\n", len(testEXEBytes))

	safeArray, err := createSafeArray(testEXEBytes)
	must(err)

	var pAssembly uintptr
	hr = appDomain.Load_3(uintptr(unsafe.Pointer(&safeArray)), &pAssembly)
	checkOK(hr)

	fmt.Printf("[+] Assembly loaded into memory at 0x%08x\n", pAssembly)
	assembly := newAssembly(pAssembly)
	var pEntryPointInfo uintptr
	fmt.Printf("entrypoint: 0x%08x\n", pEntryPointInfo)
	hr = assembly.GetEntryPoint(&pEntryPointInfo)
	checkOK(hr)
	fmt.Printf("entrypoint: 0x%08x\n", pEntryPointInfo)
	fmt.Println("[+] Found assembly entrypoint")

	var pRetCode uintptr
	nullVariant := Variant{
		VT:  1,
		Val: 0,
	}


	//printRawData(uintptr(unsafe.Pointer(&nullVariant)), unsafe.Sizeof(nullVariant))
	methodInfo := newMethodInfo(pEntryPointInfo)
	//methodGUID, err := guid.FromString("ffcc1b5d-ecb8-38dd-9b01-3dc8abc2aa5f")
	//must(err)
	//iu = newIUnknown(pEntryPointInfo)
	//hr = iu.QueryInterface(&methodGUID, &pRetCode)
	//checkOK(hr)
	//fmt.Printf("[+] MethodInfo at 0x%08x\n", pRetCode)
	////
	////methodInfo := newMethodInfo(pRetCode)
	//fmt.Println("testing gettype")
	//hr = methodInfo.GetType(&pRetCode)
	//checkOK(hr)
	//

	fmt.Println("[+] Calling default entry point with no args")
	hr = methodInfo.Invoke_3(
		uintptr(unsafe.Pointer(&nullVariant)),
		uintptr(0),
		&pRetCode)
	checkOK(hr)

	//var appDomainID uint16
	//runTimeHost.GetCurrentappDomainID(&appDomainID)
	//checkOK(hr)
	//fmt.Printf("[+] Current AppDomain ID: %d\n", appDomainID)
	//
	//fmt.Println("[+] Executing assembly...")
	//pDLLPath, _ := syscall.UTF16PtrFromString("TestDLL.dll")
	//pTypeName, _ := syscall.UTF16PtrFromString("TestDLL.HelloWorld")
	//pMethodName, _ := syscall.UTF16PtrFromString("SayHello")
	//pArgument, _ := syscall.UTF16PtrFromString("foobar")
	//var pReturnVal *uint16
	//hr = runTimeHost.ExecuteInDefaultAppDomain(
	//	pDLLPath,
	//	pTypeName,
	//	pMethodName,
	//	pArgument,
	//	pReturnVal)
	//
	//checkOK(hr)
	//fmt.Printf("[+] Assembly returned: 0x%x\n", pReturnVal)

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
	for i < ptr+size {
		if offset%16 == 0 {
			fmt.Printf("%016x : ", i)
		}

		fmt.Printf("%02x", *(*byte)(unsafe.Pointer(i)))

		i++
		offset++
		if offset%16 == 0 || offset == size {
			fmt.Print("\n")
		} else if offset%8 == 0 {
			fmt.Print("  ")
		} else {
			fmt.Print(" ")
		}
	}
}
