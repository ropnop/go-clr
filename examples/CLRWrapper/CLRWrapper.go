// +build windows

package main

import (
	clr "github.com/ropnop/go-clr"
	"log"
	"fmt"
	"io/ioutil"
	"runtime"
)

func main() {
	fmt.Println("[+] Loading DLL from Disk")
	ret, err := clr.ExecuteDLL(
		"TestDLL.dll",
		"TestDLL.HelloWorld",
		"SayHello",
		"foobar")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("[+] DLL Return Code: %d\n", ret)

	
	fmt.Println("[+] Executing EXE from memory")
	exebytes, err := ioutil.ReadFile("helloworld.exe")
	if err != nil {
		log.Fatal(err)
	}
	runtime.KeepAlive(exebytes)

	ret2, err := clr.ExecuteByteArray(exebytes)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("[+] EXE Return Code: %d\n", ret2)
}
