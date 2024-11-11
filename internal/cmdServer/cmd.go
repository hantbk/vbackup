package cmdServer

import (
	"fmt"
	"path"
	"path/filepath"
	"runtime"
	"time"

	"github.com/asdine/storm/v3"
	cf "github.com/hantbk/vbackup/internal/config"
	"github.com/hantbk/vbackup/internal/entity/v1/config"
	"github.com/hantbk/vbackup/internal/service/v1/common"
	"github.com/hantbk/vbackup/internal/service/v1/user"
	shell "github.com/hantbk/vbackup/pkg/utils/cmd"
	fileutil "github.com/hantbk/vbackup/pkg/file"
)

type Cmd struct {
	config *config.Config
	db     *storm.DB
}

func Instance(path string) *Cmd {
	cmd := &Cmd{}
	cmd.setConfig(path)
	cmd.setUpDB()
	return cmd
}

func (cmd *Cmd) setConfig(path string) {
	c, err := cf.ReadConfig(path)
	if err != nil {
		panic(err)
	}
	cmd.config = c
}

func (cmd *Cmd) setUpDB() {
	dbpath := fileutil.ReplaceHomeDir(cmd.config.Data.DbDir)
	if !fileutil.Exist(dbpath) {
		fmt.Println("Database directory does not exist")
		return
	}
	d, err := storm.Open(path.Join(dbpath, string(filepath.Separator), "vbackup.db"))
	if err != nil {
		d = nil
	}
	cmd.db = d
}

// ClearOtp - Disable two-factor authentication for a user
// username - Username
// mode - Whether to start the service. If the service is running, it needs to be temporarily stopped and restarted after the task is completed.
func (cmd *Cmd) ClearOtp(username string, mode int) {
	userService := user.GetService()
	if cmd.db == nil {
		if runtime.GOOS == "linux" {
			_, _ = shell.ExecWithTimeOut("systemctl stop vbackup.service", 1*time.Minute)
			cmd.setUpDB()
			cmd.ClearOtp(username, 1)
		} else {
			fmt.Println("Database is busy. Please manually stop vbackup_server and try again.")
		}
		return
	}
	err := userService.ClearOtp(username, common.DBOptions{DB: cmd.db})
	if mode == 1 && runtime.GOOS == "linux" {
		_, _ = shell.ExecWithTimeOut("systemctl start vbackup.service", 1*time.Minute)
	}
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("Disabled successfully.")
}

// ClearPwd - Reset user password
// username - Username
// mode - Whether to start the service. If the service is running, it needs to be temporarily stopped and restarted after the task is completed.
func (cmd *Cmd) ClearPwd(username string, mode int) {
	userService := user.GetService()
	if cmd.db == nil {
		if runtime.GOOS == "linux" {
			_, _ = shell.ExecWithTimeOut("systemctl stop vbackup.service", 1*time.Minute)
			cmd.setUpDB()
			cmd.ClearPwd(username, 1)
		} else {
			fmt.Println("Database is busy. Please manually stop vbackup_server and try again.")
		}
		return
	}
	err := userService.ClearPwd(username, common.DBOptions{DB: cmd.db})
	if mode == 1 && runtime.GOOS == "linux" {
		_, _ = shell.ExecWithTimeOut("systemctl start vbackup.service", 1*time.Minute)
	}
	if err != nil {
		fmt.Println(err)
		return
	}
}

