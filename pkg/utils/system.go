package utils

import (
	"fmt"
	"os"
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

// GetUserHomePath returns the user's home directory path based on the OS
func GetUserHomePath() string {
	homeDir, err := os.UserHomeDir() // Fetches the home directory
	if err != nil {
		fmt.Println("Error fetching user home directory:", err)
		return ""
	}
	return homeDir
}

