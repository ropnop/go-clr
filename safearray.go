// +build windows

package clr

import (
	"fmt"
	"runtime"
	"syscall"
	"unsafe"
)

// VARTYPE uint16
// UINT uint32
// VT_UI1 = 0x0011
// ULONG uint32
// LONG int32
// USHORT uint16

// from OAld.h

type SafeArray struct {
	cDims      uint16
	fFeatures  uint16
	cbElements uint32
	cLocks     uint32
	pvData     uintptr
	rgsabound  [1]SafeArrayBound
}

type SafeArrayBound struct {
	cElements uint32
	lLbound   int32
}

func CreateSafeArray(rawBytes []byte) (SafeArray, error) {
	modOleAuto, err := syscall.LoadDLL("OleAut32.dll")
	if err != nil {
		return SafeArray{}, err
	}
	procSafeArrayCreate, err := modOleAuto.FindProc("SafeArrayCreate")
	if err != nil {
		return SafeArray{}, err
	}

	size := len(rawBytes)
	sab := SafeArrayBound{
		cElements: uint32(size),
		lLbound:   0,
	}
	vt := uint16(0x11) // VT_UI1
	ret, _, _ := procSafeArrayCreate.Call(
		uintptr(vt),
		uintptr(1),
		uintptr(unsafe.Pointer(&sab)))

	if ret == 0 {
		return SafeArray{}, fmt.Errorf("Error creating SafeArray")
	}

	sa := (*SafeArray)(unsafe.Pointer(ret))
	runtime.KeepAlive(sa)
	// now we need to use RtlCopyMemory to copy our bytes to the SafeArray
	modNtDll, err := syscall.LoadDLL("ntdll.dll")
	if err != nil {
		return SafeArray{}, err
	}
	procRtlCopyMemory, err := modNtDll.FindProc("RtlCopyMemory")
	if err != nil {
		return SafeArray{}, err
	}

	ret, _, err = procRtlCopyMemory.Call(
		sa.pvData,
		uintptr(unsafe.Pointer(&rawBytes[0])),
		uintptr(size))
	if err != syscall.Errno(0) {
		return SafeArray{}, err
	}

	return *sa, nil

}
