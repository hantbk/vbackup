package utils

import (
	"fmt"
	"runtime"

	"github.com/super-l/machine-code/machine"
)

func GetSN() string {
	sn, err := machine.GetBoardSerialNumber()
	if err != nil {
		return ""
	}
	return sn
}

func GetCpuCores() int {
	fmt.Println(runtime.GOMAXPROCS(0))
	return runtime.NumCPU()
}
