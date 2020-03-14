// +build windows

package clr

import (
	"fmt"
	"log"
)

const S_OK = 0x0

func checkOK(hr uintptr, caller string) error {
	if hr != S_OK {
		return fmt.Errorf("%s returned 0x%08x", caller, hr)
	} else {
		return nil
	}
}

func must(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
