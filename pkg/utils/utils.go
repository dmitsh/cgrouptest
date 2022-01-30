package utils

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
)

const (
	CGroupRoot = "/sys/fs/cgroup"

	ProcsFile       = "cgroup.procs"
	MemoryLimitFile = "memory.limit_in_bytes"
	SwapLimitFile   = "memory.swappiness"
)

var (
	RssLimit   int
	CGroupName string
)

func init() {
	flag.IntVar(&RssLimit, "memory", 10, "memory limit with MB.")
	flag.StringVar(&CGroupName, "cgroup", "cgtest", "cgroup name")
}

func WriteFile(path string, value int) {
	if err := ioutil.WriteFile(path, []byte(fmt.Sprintf("%d", value)), 0755); err != nil {
		log.Panic(err)
	}
}
