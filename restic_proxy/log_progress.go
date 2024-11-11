package resticProxy

import (
	"fmt"
	"time"

	wsTaskInfo "github.com/hantbk/vbackup/internal/store/ws_task_info"
	"github.com/hantbk/vbackup/pkg/restic_source/rinternal/ui/progress"
	"github.com/hantbk/vbackup/pkg/utils"
)

// newProgressMax returns a progress.Counter that prints to stdout.
func newProgressMax(show bool, max uint64, description string, spr *wsTaskInfo.Sprintf) *progress.Counter {
	if !show {
		return nil
	}
	return progress.NewCounter(spr.MinUpdatePause, max, func(v uint64, max uint64, d time.Duration, final bool) {
		var status string
		if max == 0 {
			status = fmt.Sprintf("[%s]          %d %s", utils.FormatDuration(d), v, description)
		} else {
			status = fmt.Sprintf("[%s] %s  %d / %d %s",
				utils.FormatDuration(d), utils.FormatPercent(v, max), v, max, description)
		}
		spr.AppendForClear(wsTaskInfo.Info, status, final)
	})
}
