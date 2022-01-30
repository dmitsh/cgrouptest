package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"syscall"

	"github.com/dmitsh/cgrouptest/pkg/utils"
)

func main() {
	flag.Parse()

	switch os.Args[1] {
	case "start":
		start()
	case "cgr":
		cgr()
	default:
		panic("invalid input")
	}
}

func start() {
	fmt.Printf("Running %v as %d\n", os.Args[2:], os.Getpid())

	cmd := exec.Command("/proc/self/exe", append([]string{"cgr"}, os.Args[2:]...)...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	cmd.SysProcAttr = &syscall.SysProcAttr{
		Cloneflags: syscall.CLONE_NEWPID,
	}
	cmd.Run()
}

func cgr() {
	fmt.Printf("Setting cgroup for %v as %d\n", os.Args[2:], os.Getpid())
	// create cgroup directory

	pidsPath := filepath.Join(utils.CGroupRoot, "pids", utils.CGroupName)
	os.MkdirAll(pidsPath, 0755)
	utils.WriteFile(pidsPath+"/pids.max", 10)

	cgroupMemRoot := filepath.Join(utils.CGroupRoot, "memory", utils.CGroupName)

	// set memory limit
	mPath := filepath.Join(cgroupMemRoot, utils.MemoryLimitFile)
	utils.WriteFile(mPath, utils.RssLimit*1024*1024)

	// set swap memory limit to zero
	sPath := filepath.Join(cgroupMemRoot, utils.SwapLimitFile)
	utils.WriteFile(sPath, 0)

	// set cgroup procs id
	pPath := filepath.Join(cgroupMemRoot, utils.ProcsFile)
	utils.WriteFile(pPath, os.Getpid())

	cmd := exec.Command(os.Args[2], os.Args[3:]...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	cmd.Run()

	if err := cmd.Wait(); err != nil {
		fmt.Println("cmd return with error:", err)
	}

	status := cmd.ProcessState.Sys().(syscall.WaitStatus)
	exitStatus := status.ExitStatus()

	var sig syscall.Signal
	if status.Signaled() {
		sig = status.Signal()
	}

	cmd.Process.Kill()

	switch sig {
	case os.Kill:
		fmt.Println("app is killed by system")
	default:
		fmt.Println("app exit with code:", exitStatus)
	}
}
