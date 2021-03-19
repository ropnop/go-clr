// +build windows

package clr

import (
	"fmt"
	"syscall"
	"unsafe"

	"golang.org/x/sys/windows"
)

// from mscorlib.tlh

type MethodInfo struct {
	vtbl *MethodInfoVtbl
}

// MethodInfoVtbl Discovers the attributes of a method and provides access to method metadata.
// Inheritance: Object -> MemberInfo -> MethodBase -> MethodInfo
// MethodInfo Class: https://docs.microsoft.com/en-us/dotnet/api/system.reflection.methodinfo?view=net-5.0
// MethodBase Class: https://docs.microsoft.com/en-us/dotnet/api/system.reflection.methodbase?view=net-5.0
// MemberInfo Class: https://docs.microsoft.com/en-us/dotnet/api/system.reflection.memberinfo?view=net-5.0
// Object Class: https://docs.microsoft.com/en-us/dotnet/api/system.object?view=net-5.0
type MethodInfoVtbl struct {
	QueryInterface                 uintptr
	AddRef                         uintptr
	Release                        uintptr
	GetTypeInfoCount               uintptr
	GetTypeInfo                    uintptr
	GetIDsOfNames                  uintptr
	Invoke                         uintptr
	get_ToString                   uintptr
	Equals                         uintptr
	GetHashCode                    uintptr
	GetType                        uintptr
	get_MemberType                 uintptr
	get_name                       uintptr
	get_DeclaringType              uintptr
	get_ReflectedType              uintptr
	GetCustomAttributes            uintptr
	GetCustomAttributes_2          uintptr
	IsDefined                      uintptr
	GetParameters                  uintptr
	GetMethodImplementationFlags   uintptr
	get_MethodHandle               uintptr
	get_Attributes                 uintptr
	get_CallingConvention          uintptr
	Invoke_2                       uintptr
	get_IsPublic                   uintptr
	get_IsPrivate                  uintptr
	get_IsFamily                   uintptr
	get_IsAssembly                 uintptr
	get_IsFamilyAndAssembly        uintptr
	get_IsFamilyOrAssembly         uintptr
	get_IsStatic                   uintptr
	get_IsFinal                    uintptr
	get_IsVirtual                  uintptr
	get_IsHideBySig                uintptr
	get_IsAbstract                 uintptr
	get_IsSpecialName              uintptr
	get_IsConstructor              uintptr
	Invoke_3                       uintptr
	get_returnType                 uintptr
	get_ReturnTypeCustomAttributes uintptr
	GetBaseDefinition              uintptr
}

func (obj *MethodInfo) QueryInterface(riid *windows.GUID, ppvObject *uintptr) uintptr {
	ret, _, _ := syscall.Syscall(
		obj.vtbl.QueryInterface,
		3,
		uintptr(unsafe.Pointer(obj)),
		uintptr(unsafe.Pointer(riid)),
		uintptr(unsafe.Pointer(ppvObject)))
	return ret
}

func (obj *MethodInfo) AddRef() uintptr {
	ret, _, _ := syscall.Syscall(
		obj.vtbl.AddRef,
		1,
		uintptr(unsafe.Pointer(obj)),
		0,
		0)
	return ret
}

func (obj *MethodInfo) Release() uintptr {
	ret, _, _ := syscall.Syscall(
		obj.vtbl.Release,
		1,
		uintptr(unsafe.Pointer(obj)),
		0,
		0)
	return ret
}

// Invoke_3 Invokes the method or constructor reflected by this MethodInfo instance.
//      virtual HRESULT __stdcall Invoke_3 (
//      /*[in]*/ VARIANT obj,
//      /*[in]*/ SAFEARRAY * parameters,
//      /*[out,retval]*/ VARIANT * pRetVal ) = 0;
// https://docs.microsoft.com/en-us/dotnet/api/system.reflection.methodbase.invoke?view=net-5.0
func (obj *MethodInfo) Invoke_3(variantObj Variant, parameters *SafeArray) (err error) {
	debugPrint("Entering into methodinfo.Invoke_3()...")
	var pRetVal *Variant
	hr, _, err := syscall.Syscall6(
		obj.vtbl.Invoke_3,
		4,
		uintptr(unsafe.Pointer(obj)),
		uintptr(unsafe.Pointer(&variantObj)),
		uintptr(unsafe.Pointer(parameters)),
		uintptr(unsafe.Pointer(pRetVal)),
		0,
		0,
	)
	if err != syscall.Errno(0) {
		err = fmt.Errorf("the MethodInfo::Invoke_3 method returned an error:\r\n%s", err)
		return
	}
	if hr != S_OK {
		err = fmt.Errorf("the Assembly::Invoke_3 method returned a non-zero HRESULT: 0x%x", hr)
		return
	}
	if pRetVal != nil {
		err = fmt.Errorf("the Assembly::Invoke_3 method returned a non-zero pRetVal: %+v", pRetVal)
		return
	}
	err = nil
	return
}

// GetString returns a string that represents the current object
// a string version of the method's signature
// public virtual string ToString ();
// https://docs.microsoft.com/en-us/dotnet/api/system.object.tostring?view=net-5.0#System_Object_ToString
func (obj *MethodInfo) GetString() (str string, err error) {
	debugPrint("Entering into methodinfo.GetString()...")
	var object *string
	hr, _, err := syscall.Syscall(
		obj.vtbl.get_ToString,
		2,
		uintptr(unsafe.Pointer(obj)),
		uintptr(unsafe.Pointer(&object)),
		0,
	)
	if err != syscall.Errno(0) {
		err = fmt.Errorf("the MethodInfo::ToString method returned an error:\r\n%s", err)
		return
	}
	if hr != S_OK {
		err = fmt.Errorf("the Assembly::ToString method returned a non-zero HRESULT: 0x%x", hr)
		return
	}
	err = nil
	str = ReadUnicodeStr(unsafe.Pointer(object))
	return
}
