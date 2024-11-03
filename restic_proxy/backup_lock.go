package resticProxy

import (
	"strconv"
	"time"

	"github.com/fanjindong/go-cache"
	"github.com/hantbk/vbackup/internal/consts"
)

func BackupLock(repo int, path string) bool {
	key := consts.Key("BackupIsRun", strconv.Itoa(repo), path)
	c := server.Cache()
	res, ok := c.Get(key)
	if ok && res == 1 {
		return false
	} else {
		c.Set(key, 1, cache.WithEx(24*time.Hour))
		return true
	}
}

func BackupUnlock(repo int, path string) {
	key := consts.Key("BackupIsRun", strconv.Itoa(repo), path)
	c := server.Cache()
	c.Del(key)
}
