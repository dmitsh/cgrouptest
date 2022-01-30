package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"syscall"

	"github.com/dmitsh/cgrouptest/pkg/utils"
)

func main() {
	flag.Parse()

	// create cgroup directory
	cgroupMemDir := filepath.Join(utils.CGroupRoot, "memory", utils.CGroupName)
	if err := os.MkdirAll(cgroupMemDir, 0755); err != nil {
		log.Panic(err)
	}

	// set memory limit
	mPath := filepath.Join(cgroupMemDir, utils.MemoryLimitFile)
	utils.WriteFile(mPath, utils.RssLimit*1024*1024)

	// set swap memory limit to zero
	sPath := filepath.Join(cgroupMemDir, utils.SwapLimitFile)
	utils.WriteFile(sPath, 0)

	cmd := exec.Command("./app")
	cmd.Stdout = os.Stdout

	// start app
	if err := cmd.Start(); err != nil {
		log.Panic(err)
	}

	fmt.Println("add pid", cmd.Process.Pid, "to file cgroup.procs")

	// set cgroup procs id
	pPath := filepath.Join(cgroupMemDir, utils.ProcsFile)
	utils.WriteFile(pPath, cmd.Process.Pid)

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
