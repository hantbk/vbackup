package rclone

import (
	"testing"

	"github.com/hantbk/vbackup/pkg/restic_source/rinternal/backend/test"
)

var configTests = []test.ConfigTestData[Config]{
	{
		S: "rclone:local:foo:/bar",
		Cfg: Config{
			Remote:      "local:foo:/bar",
			Program:     defaultConfig.Program,
			Args:        defaultConfig.Args,
			Connections: defaultConfig.Connections,
			Timeout:     defaultConfig.Timeout,
		},
	},
}

func TestParseConfig(t *testing.T) {
	test.ParseConfigTester(t, ParseConfig, configTests)
}
