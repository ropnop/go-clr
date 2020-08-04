// +build windows

package clr

import (
	"syscall"
	"unsafe"

	"golang.org/x/sys/windows"
)

// from mscorlib.tlh

type MethodInfo struct {
	vtbl *MethodInfoVtbl
}

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

func NewMethodInfoFromPtr(ppv uintptr) *MethodInfo {
	return (*MethodInfo)(unsafe.Pointer(ppv))
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

func (obj *MethodInfo) GetType(pRetVal *uintptr) uintptr {
	ret, _, _ := syscall.Syscall(
		obj.vtbl.GetType,
		2,
		uintptr(unsafe.Pointer(obj)),
		uintptr(unsafe.Pointer(pRetVal)),
		0)
	return ret
}

func (obj *MethodInfo) Invoke_3(variantObj Variant, parameters uintptr, pRetVal *uintptr) uintptr {
	ret, _, _ := syscall.Syscall6(
		obj.vtbl.Invoke_3,
		4,
		uintptr(unsafe.Pointer(obj)),
		uintptr(unsafe.Pointer(&variantObj)),
		parameters,
		uintptr(unsafe.Pointer(pRetVal)),
		0,
		0,
	)
	return ret
}

// GetString returns a string version of the method's signature
func (obj *MethodInfo) GetString(addr *uintptr) error {
	ret, _, _ := syscall.Syscall(
		obj.vtbl.get_ToString,
		2,
		uintptr(unsafe.Pointer(obj)),
		uintptr(unsafe.Pointer(addr)),
		0,
	)
	return checkOK(ret, "get_ToString")
}
