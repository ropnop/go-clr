// +build windows

package clr

import (
	"fmt"
	"io"
	"os"

	"golang.org/x/sys/windows"
)

// origSTDOUT is a Windows Handle to the program's original STDOUT
var origSTDOUT = windows.Stdout

// origSTDERR is a Windows Handle to the program's original STDERR
var origSTDERR = windows.Stderr

// rSTDOUT is an io.Reader for STDOUT
var rSTDOUT *os.File

// wSTDOUT is an io.Writer for STDOUT
var wSTDOUT *os.File

// rSTDERR is an io.Reader for STDERR
var rSTDERR *os.File

// wSTDERR is an io.Writer for STDERR
var wSTDERR *os.File

// RedirectStdoutStderr redirects the program's STDOUT/STDERR to an *os.File that can be read from this Go program
// The CLR executes assemblies outside of Go and therefore STDOUT/STDERR can't be captured using normal functions
// Intended to be used with a Command & Control framework so STDOUT/STDERR can be captured and returned
func RedirectStdoutStderr() (err error) {
	// Create a new reader and writer for STDOUT
	rSTDOUT, wSTDOUT, err = os.Pipe()
	if err != nil {
		err = fmt.Errorf("there was an error calling the os.Pipe() function to create a new STDOUT:\n%s", err)
		return
	}

	// Createa new reader and writer for STDERR
	rSTDERR, wSTDERR, err = os.Pipe()
	if err != nil {
		err = fmt.Errorf("there was an error calling the os.Pipe() function to create a new STDERR:\n%s", err)
		return
	}

	// Set STDOUT/STDERR to the new files from os.Pipe()
	// https://docs.microsoft.com/en-us/windows/console/setstdhandle
	if err = windows.SetStdHandle(windows.STD_OUTPUT_HANDLE, windows.Handle(wSTDOUT.Fd())); err != nil {
		err = fmt.Errorf("there was an error calling the windows.SetStdHandle function for STDOUT:\n%s", err)
		return
	}

	if err = windows.SetStdHandle(windows.STD_ERROR_HANDLE, windows.Handle(wSTDERR.Fd())); err != nil {
		err = fmt.Errorf("there was an error calling the windows.SetStdHandle function for STDERR:\n%s", err)
		return
	}

	return
}

// RestoreStdoutStderr returns the program's original STDOUT/STDERR handles before they were redirected an *os.File
// Previously instantiated CLRs will continue to use the REDIRECTED STDOUT/STDERR handles and will not resume
// using the restored handles
func RestoreStdoutStderr() error {
	if err := windows.SetStdHandle(windows.STD_OUTPUT_HANDLE, origSTDOUT); err != nil {
		return fmt.Errorf("there was an error calling the windows.SetStdHandle function to restore the original STDOUT handle:\n%s", err)
	}
	if err := windows.SetStdHandle(windows.STD_ERROR_HANDLE, origSTDERR); err != nil {
		return fmt.Errorf("there was an error calling the windows.SetStdHandle function to restore the original STDERR handle:\n%s", err)
	}
	return nil
}

// ReadStdoutStderr reads from the REDIRECTED STDOUT/STDERR
// Only use when RedirectStdoutStderr was previously called
func ReadStdoutStderr() (stdout string, stderr string, err error) {
	debugPrint("Entering io.ReadStdoutStderr()...")

	// If nothing was written to STDOUT then Read() will block
	// Can't call Close() because the pipe needs to remain open for the duration of the top-level program
	// A "workaround" is to write in a null byte so that way it can be read and won't block
	_, err = wSTDOUT.Write([]byte{0x00})
	if err != nil {
		err = fmt.Errorf("there was an error writing a null-byte into STDOUT Writer:\n%s", err)
		return
	}

	// TODO Update to use io.ReadAll(), requires GO 1.16
	// https://golang.org/pkg/io/#ReadAll
	bStdout := make([]byte, 500000)
	c, err := rSTDOUT.Read(bStdout)
	// Will return EOF if there is no data to be read
	if err != nil && err != io.EOF {
		err = fmt.Errorf("there was an error reading from the STDOUT Reader:\n%s", err)
		return
	}
	// If STDOUT is contains more than the null byte we wrote into it, then capture it
	if c > 1 && bStdout[1] != 0x00 {
		stdout = string(bStdout[:])
	}

	// If nothing was written to STDERR then Read() will block
	// Can't call Close() because the pipe needs to remain open for the duration of the top-level program
	// A "workaround" is to write in a null byte so that way it can be read and won't block
	_, err = wSTDERR.Write([]byte{0x00})
	if err != nil {
		err = fmt.Errorf("there was an error writing a null-byte into STDERR Writer:\n%s", err)
		return
	}
	bStderr := make([]byte, 500000)
	c, err = rSTDERR.Read(bStderr)
	// Will return EOF when nothing was written to it if Close() was called first
	if err != nil && err != io.EOF {
		err = fmt.Errorf("there was an error reading from the STDERR Reader:\n%s", err)
		return
	}
	err = nil
	// If STDERR is contains more than the null byte we wrote into it, then capture it
	if c > 1 && bStderr[1] != 0x00 {
		stderr = string(bStderr[:])
	}

	return
}

// CloseSTdoutStderr closes the Reader/Writer for the prviously redirected STDOUT/STDERR
// that was changed to an *os.File
func CloseStdoutStderr() (err error) {
	err = rSTDOUT.Close()
	if err != nil {
		err = fmt.Errorf("there was an error closing the STDOUT Reader:\n%s", err)
		return
	}

	err = wSTDOUT.Close()
	if err != nil {
		err = fmt.Errorf("there was an error closing the STDOUT Writer:\n%s", err)
		return
	}

	err = rSTDERR.Close()
	if err != nil {
		err = fmt.Errorf("there was an error closing the STDERR Reader:\n%s", err)
		return
	}

	err = wSTDERR.Close()
	if err != nil {
		err = fmt.Errorf("there was an error closing the STDERR Writer:\n%s", err)
		return
	}
	return nil
}
