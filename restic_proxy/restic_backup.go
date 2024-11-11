package resticProxy

import (
	ser "github.com/hantbk/vbackup/internal/service/v1/task"
	"github.com/hantbk/vbackup/pkg/restic_source/rinternal/errors"
	"github.com/hantbk/vbackup/pkg/restic_source/rinternal/restic"
)

var taskHistoryService ser.Service

func init() {
	taskHistoryService = ser.GetService()
}

var ErrInvalidSourceData = errors.New("at least one source file could not be read")

// BackupOptions bundles all options for the backup command.
type BackupOptions struct {
	excludePatternOptions
	Parent            string
	GroupBy           restic.SnapshotGroupByOptions
	Force             bool
	ExcludeOtherFS    bool
	ExcludeIfPresent  []string
	ExcludeCaches     bool
	ExcludeLargerThan string
	Tags              restic.TagLists
	Host              string
	TimeStamp         string // `time` of the backup (ex. '2012-11-01 22:08:41') (default: now)
	WithAtime         bool
	IgnoreInode       bool
	IgnoreCtime       bool
	UseFsSnapshot     bool
	DryRun            bool
	ReadConcurrency   uint // // Read concurrency count, default is 2
}
