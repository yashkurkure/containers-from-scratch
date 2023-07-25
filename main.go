package main

import (
	"fmt"
	"os"
	"os/exec"
	"syscall"
)

// docker 			run		image	<cmd>	<params>
// go run main.go	run 			<cmd>	<params>


func main() {
	switch os.Args[1] {
	case "run":
		run()
	case "child":
		child()
	default:
		panic("bad command")
	}
}

// This creates the namespace
func run() {
	fmt.Printf("Running %v as %d\n", os.Args[2:], os.Getpid())

	cmd := exec.Command("/proc/self/exe", append([]string{"child"}, os.Args[2:]...)...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	// Let us containerize the process in which the command is run
	// Configuring the namespace
	// UTS (UNIX Time-Sharing) namespaces allow a single system to appear to have different host and domain names to different processes
	// PID (Process ID) namespaces isolate the process ID number space, meaning that processes in different PID namespaces can have the same PID
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Cloneflags:	syscall.CLONE_NEWUTS | syscall.CLONE_NEWPID,
	}

	cmd.Run()
}


func child() {
	fmt.Printf("Running %v as %d\n", os.Args[2:], os.Getpid())

	// Change the host name
	syscall.Sethostname([]byte("container"))

	// Change the root
	syscall.Chroot("/home/yash/arg/containers-from-scratch/fakeroot")
	os.Chdir("/")

	cmd := exec.Command(os.Args[2], os.Args[3:]...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	cmd.Run()
}


func must(err error) {
	if err != nil {
		panic(err)
	}
}
