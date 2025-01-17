package resticProxy

import (
	"math"
	"sort"
	"time"

	"github.com/hantbk/vbackup/internal/model"
	"github.com/hantbk/vbackup/internal/server"
	"github.com/hantbk/vbackup/internal/service/v1/common"
	"github.com/hantbk/vbackup/internal/store/task"
	wsTaskInfo "github.com/hantbk/vbackup/internal/store/ws_task_info"
	"github.com/hantbk/vbackup/pkg/restic_source/rinternal/archiver"
	"github.com/hantbk/vbackup/pkg/restic_source/rinternal/restic"
	"github.com/hantbk/vbackup/pkg/restic_source/rinternal/ui"
	"github.com/hantbk/vbackup/pkg/restic_source/rinternal/ui/backup"
	"github.com/hantbk/vbackup/pkg/utils"
)

type TaskProgress struct {
	*ui.StdioWrapper
	task           wsTaskInfo.WsTaskInfo
	weightCount    float64 // Count progress weight
	weightSize     float64 // Size progress weight
	lastUpdate     time.Time
	errors         []model.ErrorUpdate
	minUpdatePause time.Duration
}

func (t *TaskProgress) E(msg string, args ...interface{}) {
	errorUpdate := model.ErrorUpdate{
		MessageType: "error",
		Error:       msg,
		During:      "backup",
		Item:        "",
	}
	if len(t.errors) > 20 {
		return
	}
	t.print(errorUpdate, true)
	t.errors = append(t.errors, errorUpdate)
	_ = taskHistoryService.UpdateField(t.task.GetId(), "ArchivalError", t.errors, common.DBOptions{})
}

func (t *TaskProgress) P(msg string, args ...interface{}) {
	t.print(msg, true)
}

func (t *TaskProgress) V(msg string, args ...interface{}) {
	t.print(msg, true)
}

var _ backup.ProgressPrinter = &TaskProgress{}

func NewTaskProgress(task wsTaskInfo.WsTaskInfo, minUpdatePause time.Duration) *TaskProgress {
	return &TaskProgress{
		task:           task,
		errors:         make([]model.ErrorUpdate, 0),
		weightCount:    1,
		weightSize:     1,
		minUpdatePause: minUpdatePause,
	}
}

func (t *TaskProgress) UpdateTaskInfo(task wsTaskInfo.WsTaskInfo) {
	t.task = task
}
func (t *TaskProgress) SetWeight(weightCount, weightSize float64) {
	t.weightSize = weightSize
	t.weightCount = weightCount
}

func (t *TaskProgress) print(status interface{}, forceUpdate bool) {
	// limit update frequency
	if !forceUpdate && (time.Since(t.lastUpdate) < t.minUpdatePause || t.minUpdatePause == 0) {
		return
	}
	t.task.SendMsg(status)
	if !forceUpdate {
		t.lastUpdate = time.Now()
	}
}

func (t *TaskProgress) Update(total, processed backup.Counter, errors uint, currentFiles map[string]struct{}, start time.Time, secs uint64) {
	duration := time.Since(start)
	avg := ""
	if duration/time.Second > 0 {
		avg = utils.FormatBytesSpeed(processed.Bytes / uint64(duration/time.Second))
	}

	status := model.StatusUpdate{
		MessageType:      "status",
		SecondsElapsed:   utils.FormatDuration(duration),
		SecondsRemaining: utils.FormatSeconds(secs),
		TotalFiles:       total.Files,
		FilesDone:        processed.Files,
		TotalBytes:       utils.FormatBytes(total.Bytes),
		BytesDone:        utils.FormatBytes(processed.Bytes),
		ErrorCount:       errors,
		AvgSpeed:         avg,
	}

	if total.Bytes > 0 && total.Files > 0 {
		denominator := float64(total.Files)*t.weightCount + float64(total.Bytes)*t.weightSize
		numerator := float64(processed.Files)*t.weightCount + float64(processed.Bytes)*t.weightSize
		status.PercentDone = numerator / denominator
		status.PercentDone = math.Floor(status.PercentDone*100) / 100
	}

	for filename := range currentFiles {
		status.CurrentFiles = append(status.CurrentFiles, filename)
	}
	sort.Strings(status.CurrentFiles)
	t.task.(*task.TaskInfo).Progress = &status
	task.TaskInfos.Set(t.task.GetId(), t.task)
	// The sending frequency is controlled by the reporter, no control is done here
	t.print(&status, true)
}

func (t *TaskProgress) ScannerError(item string, err error) error {
	errorUpdate := &model.ErrorUpdate{
		MessageType: "error",
		Error:       err.Error(),
		During:      "scan",
		Item:        item,
	}
	t.print(errorUpdate, true)
	err1 := taskHistoryService.UpdateField(t.task.GetId(), "ScannerError", errorUpdate, common.DBOptions{})
	if err1 != nil {
		return err1
	}
	return err
}

func (t *TaskProgress) Error(item string, err error) error {
	errorUpdate := model.ErrorUpdate{
		MessageType: "error",
		Error:       err.Error(),
		During:      "archival",
		Item:        item,
	}
	if len(t.errors) > 20 {
		return err
	}
	t.print(&errorUpdate, true)
	t.errors = append(t.errors, errorUpdate)
	err1 := taskHistoryService.UpdateField(t.task.GetId(), "ArchivalError", t.errors, common.DBOptions{})
	if err1 != nil {
		return err1
	}
	return err
}

func (t *TaskProgress) CompleteItem(messageType string, item string, s archiver.ItemStats, d time.Duration) {
	var status model.VerboseUpdate
	switch messageType {
	case "dir new":
		status = model.VerboseUpdate{
			MessageType:  "verbose_status",
			Action:       "new",
			Item:         item,
			Duration:     utils.FormatDuration(d),
			DataSize:     utils.FormatBytes(s.DataSize),
			MetadataSize: utils.FormatBytes(s.TreeSize),
		}
	case "dir unchanged":
		status = model.VerboseUpdate{
			MessageType: "verbose_status",
			Action:      "unchanged",
			Item:        item,
		}
	case "dir modified":
		status = model.VerboseUpdate{
			MessageType:  "verbose_status",
			Action:       "modified",
			Item:         item,
			Duration:     utils.FormatDuration(d),
			DataSize:     utils.FormatBytes(s.DataSize),
			MetadataSize: utils.FormatBytes(s.TreeSize),
		}
	case "file new":
		status = model.VerboseUpdate{
			MessageType: "verbose_status",
			Action:      "new",
			Item:        item,
			Duration:    utils.FormatDuration(d),
			DataSize:    utils.FormatBytes(s.DataSize),
		}
	case "file unchanged":
		status = model.VerboseUpdate{
			MessageType: "verbose_status",
			Action:      "unchanged",
			Item:        item,
		}
	case "file modified":
		status = model.VerboseUpdate{
			MessageType: "verbose_status",
			Action:      "modified",
			Item:        item,
			Duration:    utils.FormatDuration(d),
			DataSize:    utils.FormatBytes(s.DataSize),
		}
	}
	t.print(&status, false)
}

func (t *TaskProgress) ReportTotal(start time.Time, s archiver.ScanStats) {
	ver := &model.VerboseUpdate{
		MessageType: "verbose_status",
		Action:      "scan_finished",
		Duration:    utils.FormatDuration(time.Since(start)),
		DataSize:    utils.FormatBytes(s.Bytes),
		TotalFiles:  s.Files,
	}
	t.print(ver, true)
	err := taskHistoryService.UpdateField(t.task.GetId(), "Scanner", ver, common.DBOptions{})
	if err != nil {
		return
	}
}

func (t *TaskProgress) Finish(snapshotID restic.ID, start time.Time, summary *backup.Summary, dryRun bool) {
	var summaryOut *model.SummaryOutput
	var p1 *model.StatusUpdate
	if summary != nil {
		summaryOut = &model.SummaryOutput{
			MessageType:         "summary",
			FilesNew:            summary.Files.New,
			FilesChanged:        summary.Files.Changed,
			FilesUnmodified:     summary.Files.Unchanged,
			DirsNew:             summary.Dirs.New,
			DirsChanged:         summary.Dirs.Changed,
			DirsUnmodified:      summary.Dirs.Unchanged,
			DataBlobs:           summary.ItemStats.DataBlobs,
			TreeBlobs:           summary.ItemStats.TreeBlobs,
			DataAdded:           utils.FormatBytes(summary.ItemStats.DataSize + summary.ItemStats.TreeSize),
			TotalFilesProcessed: summary.Files.New + summary.Files.Changed + summary.Files.Unchanged,
			TotalBytesProcessed: utils.FormatBytes(summary.ProcessedBytes),
			TotalDuration:       utils.FormatDuration(time.Since(start)),
			SnapshotID:          snapshotID.Str(),
			DryRun:              dryRun,
		}
		t.print(summaryOut, true)
		p1 = t.task.(*task.TaskInfo).Progress
		if p1 != nil {
			p1.BytesDone = p1.TotalBytes
			p1.PercentDone = 1
			p1.FilesDone = p1.TotalFiles
			p1.SecondsRemaining = "0"
			p1.SecondsElapsed = summaryOut.TotalDuration
			t.print(p1, true)
		}
	}
	taskhis, err3 := taskHistoryService.Get(t.task.GetId(), common.DBOptions{})
	if err3 != nil {
		server.Logger().Error(err3)
		return
	}
	status := task.StatusEnd
	if taskhis.ScannerError != nil || len(taskhis.ArchivalError) > 0 {
		status = task.StatusError
	}
	taskhis.Status = status
	taskhis.Summary = summaryOut
	taskhis.Progress = p1
	_ = taskHistoryService.Update(taskhis, common.DBOptions{})
	task.TaskInfos.Close(t.task.GetId(), "process end", 1)
	go GetAllRepoStats()
}

func (t *TaskProgress) Reset() {

}
