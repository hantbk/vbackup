package docker

import (
	"io/ioutil"
	"os"
	"strings"
)

// IsDockerEnv returns true if the current process is running in a docker container.
func IsDockerEnv() bool {
	_, err := os.Stat("/.dockerenv")
	if err == nil {
		return true
	}
	content, err := ioutil.ReadFile("/proc/1/cgroup")
	if err != nil {
		return false
	}
	for _, line := range strings.Split(string(content), "\n") {
		if strings.Contains(line, "/docker/") || strings.Contains(line, "/kubepod/") {
			return true
		}
	}
	return false
}
