package system

import (
	"fmt"
	"os"
	"path"
	"runtime"
	"time"

	"github.com/hantbk/vbackup/internal/consts/global"
	"github.com/hantbk/vbackup/internal/consts/system_status"
)

func Upgrade(version string) error {
	if server.IsDocker() {
		return fmt.Errorf("Currently running in a Docker environment. Please update the container image version manually")
	}
	if version == "" {
		return fmt.Errorf("Version number cannot be empty")
	}
	timeStr := time.Now().Format("200601021504")
	upgradeDir := path.Join(server.Config().Data.CacheDir, fmt.Sprintf("upgrade/upgrade_%s/downloads", timeStr))
	if err := os.MkdirAll(upgradeDir, os.ModePerm); err != nil {
		return err
	}
	downloadUrl := fmt.Sprintf("%s/%s/vbackup_server_%s_%s_%s", global.DownloadUrl, version, version, runtime.GOOS, runtime.GOARCH)
	server.UpdateSystemStatus(system_status.Upgrade)
	go func() {
		err := http.DownloadFile(downloadUrl, path.Join(upgradeDir, "vbackup_server"))
		if err != nil {
			server.Logger().Errorf("Failed to download vbackup_server file, error: %v", err)
			server.UpdateSystemStatus(system_status.Normal)
			return
		}
		if runtime.GOOS == "linux" {
			err = http.DownloadFile(global.ServiceFileUrl, path.Join(upgradeDir, "vbackup.service"))
			if err != nil {
				server.Logger().Errorf("Failed to download vbackup_server file, error: %v", err)
				server.UpdateSystemStatus(system_status.Normal)
				return
			}
		}
		server.Logger().Println("All files downloaded successfully")
		defer func() {
			_ = os.Remove(upgradeDir)
		}()
		newFile := path.Join(upgradeDir, "vbackup_server")
		oldFile := "/usr/local/bin/vbackup_server"
		err = os.Remove(oldFile)
		if err != nil {
			server.Logger().Errorf("Failed to update vbackup_server, error: %v", err)
			server.UpdateSystemStatus(system_status.Normal)
			return
		}
		err = os.Rename(newFile, oldFile)
		if err != nil {
			server.Logger().Errorf("Failed to update vbackup_server, error: %v", err)
			server.UpdateSystemStatus(system_status.Normal)
			return
		}
		if runtime.GOOS == "linux" {
			_ = os.Remove("/etc/systemd/system/vbackup.service")
			_ = os.Rename(path.Join(upgradeDir, "vbackup.service"), "/etc/systemd/system/vbackup.service")
		}
		mode := os.FileMode(0755)
		err = os.Chmod(oldFile, mode)
		if err != nil {
			server.Logger().Errorf("Failed to update vbackup_server, error: %v", err)
			server.UpdateSystemStatus(system_status.Normal)
			return
		}
		if runtime.GOOS == "linux" {
			_, _ = cmd.ExecWithTimeOut("systemctl daemon-reload && systemctl restart vbackup.service", 2*time.Minute)
		}
		server.Logger().Println("Update successful")
		server.UpdateSystemStatus(system_status.Normal)
	}()
	return nil
}
