package resticProxy

import (
	"context"
	"fmt"
	"log"
	"os"
	"runtime"
	"time"

	"github.com/hantbk/vbackup/internal/consts/globalcontext"
	"github.com/hantbk/vbackup/internal/server"
	"github.com/hantbk/vbackup/internal/service/v1/common"
	ser "github.com/hantbk/vbackup/internal/service/v1/task"
	"github.com/hantbk/vbackup/internal/store/task"
	"github.com/hantbk/vbackup/pkg/restic_source/rinternal/archiver"
	"github.com/hantbk/vbackup/pkg/restic_source/rinternal/errors"
	"github.com/hantbk/vbackup/pkg/restic_source/rinternal/fs"
	"github.com/hantbk/vbackup/pkg/restic_source/rinternal/repository"
	"github.com/hantbk/vbackup/pkg/restic_source/rinternal/restic"
	"github.com/hantbk/vbackup/pkg/restic_source/rinternal/ui/backup"
	"gopkg.in/tomb.v2"
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
	ReadConcurrency   uint // Read concurrency count, default is 2
}

func RunBackup(opts BackupOptions, repoid int, taskinfo task.TaskInfo) error {

	// fmt.Println("RunBackup called")

	if opts.Host == "" {
		hostname, err := os.Hostname()
		if err != nil {
			return fmt.Errorf("os.Hostname() returned err: %v", err)
		}
		opts.Host = hostname
	}
	repoHandler, err := GetRepository(repoid)
	if err != nil {
		return err
	}
	repo := repoHandler.repo

	ctx, cancel := context.WithCancel(context.Background())
	clean := NewCleanCtx()
	clean.AddCleanCtx(func() {
		cancel()
	})

	targets := []string{taskinfo.Path}
	timeStamp := time.Now()

	var t tomb.Tomb
	progressPrinter := NewTaskProgress(&taskinfo, time.Second)
	progressReporter := backup.NewProgress(progressPrinter, time.Second)
	clean.AddCleanCtx(func() {
		progressReporter.Done()
	})
	if opts.DryRun {
		repo.SetDryRun()
	}

	lock, err := lockRepo(ctx, repo)
	if err != nil {
		clean.Cleanup()
		return err
	}
	clean.AddCleanCtx(func() {
		unlockRepo(lock)
	})

	rejectByNameFuncs, err := collectRejectByNameFuncs(opts, repo)
	if err != nil {
		clean.Cleanup()
		return err
	}
	rejectFuncs, err := collectRejectFuncs(opts, targets)
	if err != nil {
		clean.Cleanup()
		return err
	}
	parentSnapshot, err := findParentSnapshot(ctx, repo, opts, targets, timeStamp)
	if err != nil {
		clean.Cleanup()
		return err
	}
	if parentSnapshot != nil {
		err = taskHistoryService.UpdateField(taskinfo.GetId(), "ParentId", parentSnapshot.ID().Str(), common.DBOptions{})
		if err != nil {
			clean.Cleanup()
			return err
		}
	}

	selectByNameFilter := func(item string) bool {
		for _, reject := range rejectByNameFuncs {
			if reject(item) {
				return false
			}
		}
		return true
	}

	selectFilter := func(item string, fi os.FileInfo) bool {
		for _, reject := range rejectFuncs {
			if reject(item, fi) {
				return false
			}
		}
		return true
	}
	var targetFS fs.FS = fs.Local{}
	if runtime.GOOS == "windows" && opts.UseFsSnapshot {
		if err = fs.HasSufficientPrivilegesForVSS(); err != nil {
			return err
		}

		errorHandler := func(item string, err error) error {
			return progressReporter.Error(item, err)
		}

		messageHandler := func(msg string, args ...interface{}) {
			progressPrinter.P(msg, args...)
		}

		localVss := fs.NewLocalVss(errorHandler, messageHandler)
		defer localVss.DeleteSnapshots()
		targetFS = localVss
	}
	sc := archiver.NewScanner(targetFS)
	sc.SelectByName = selectByNameFilter
	sc.Select = selectFilter
	sc.Error = progressPrinter.ScannerError
	sc.Result = progressReporter.ReportTotal

	t.Go(func() error { return sc.Scan(t.Context(ctx), targets) })

	arch := archiver.New(repo, targetFS, archiver.Options{
		ReadConcurrency: opts.ReadConcurrency,
	})
	arch.SelectByName = selectByNameFilter
	arch.Select = selectFilter
	arch.WithAtime = opts.WithAtime
	success := true
	arch.Error = func(item string, err error) error {
		success = false
		return progressReporter.Error(item, err)
	}
	arch.CompleteItem = progressReporter.CompleteItem
	arch.StartFile = progressReporter.StartFile
	arch.CompleteBlob = progressReporter.CompleteBlob

	if opts.IgnoreInode {
		// --ignore-inode implies --ignore-ctime: on FUSE, the ctime is not
		// reliable either.
		arch.ChangeIgnoreFlags |= archiver.ChangeIgnoreCtime | archiver.ChangeIgnoreInode
	}
	if opts.IgnoreCtime {
		arch.ChangeIgnoreFlags |= archiver.ChangeIgnoreCtime
	}

	snapshotOpts := archiver.SnapshotOptions{
		Excludes:       opts.Excludes,
		Tags:           opts.Tags.Flatten(),
		Time:           timeStamp,
		Hostname:       opts.Host,
		ParentSnapshot: parentSnapshot,
		ProgramVersion: "restic " + version,
	}
	bound := make(chan string)
	taskinfo.SetBound(bound)
	task.TaskInfos.Set(taskinfo.GetId(), &taskinfo)
	t.Go(func() error {
		for {
			select {
			case <-t.Context(ctx).Done():
				return nil
			case <-task.TaskInfos.Get(taskinfo.GetId()).GetBound():
				info := task.TaskInfos.Get(taskinfo.GetId())
				progressPrinter.UpdateTaskInfo(info)
			}
		}
	})
	if !BackupLock(repoid, taskinfo.Path) {
		clean.Cleanup()
		return fmt.Errorf("Repository \"%d\" is being backed up: %s", repoid, taskinfo.Path)
	}
	clean.AddCleanCtx(func() {
		BackupUnLock(repoid, taskinfo.Path)
	})

	go func() {
		defer clean.Cleanup()

		// var errorMessages []string

		err = taskHistoryService.UpdateField(taskinfo.GetId(), "Status", task.StatusRunning, common.DBOptions{})
		if err != nil {
			server.Logger().Error(err)
			// errorMessages = append(errorMessages, fmt.Sprintf("Failed to update task status to running: %v", err))
		}

		_, id, err := arch.Snapshot(ctx, targets, snapshotOpts)
		if err != nil {
			progressPrinter.E(fmt.Errorf("unable to save snapshot: %v", err).Error())
			// errorMessages = append(errorMessages, fmt.Sprintf("Unable to save snapshot: %v", err))
		}

		t.Kill(nil)
		werr := t.Wait()
		if werr != nil {
			server.Logger().Error(werr)
			// errorMessages = append(errorMessages, fmt.Sprintf("Error waiting for task to complete: %v", werr))
		}

		if !success {
			progressPrinter.E(ErrInvalidSourceData.Error())
			// errorMessages = append(errorMessages, fmt.Sprintf("Invalid source data: %v", ErrInvalidSourceData.Error()))
		}

		// // Nếu có lỗi, gửi email tổng hợp
		// if len(errorMessages) > 0 {
		// 	sendCombinedErrorEmail(errorMessages, taskinfo)
		// }

		progressReporter.Finish(id, opts.DryRun)
	}()

	return nil
}

func sendCombinedErrorEmail(errorMessages []string, taskinfo task.TaskInfo) {

	// Get the current user
	currentUser, _ := globalcontext.GetCurrentUser()

	// Get the current time
	currentTime := time.Now().Format("2006-01-02 15:04:05")

	// Tạo nội dung email từ danh sách lỗi
	errorDetails := ""
	for i, msg := range errorMessages {
		errorDetails += fmt.Sprintf("<p><strong>Error %d:</strong> %s</p>", i+1, msg)
	}

	emailBody := fmt.Sprintf(`
	<html>
		<body style="font-family: Arial, sans-serif; color: #333; background-color: #f4f4f4; padding: 20px;">
		<div style="max-width: 600px; margin: 0 auto; background-color: #ffffff; padding: 20px; border-radius: 8px; box-shadow: 0 4px 10px rgba(0, 0, 0, 0.1);">
			<img src="https://github.com/hantbk/vbackup/blob/main/web/dashboard/src/assets/logo/vbackup-bar.png?raw=true" 
			     alt="vBackup Logo" 
			     style="max-width: 200px; margin: 0 auto 20px auto; display: block;" />
			<h2 style="color: #d9534f; font-size: 24px; margin-bottom: 20px;">Task Error Notification</h2>
			<p style="font-size: 16px; margin-bottom: 10px;">
				<strong>Task ID:</strong> %d
			</p>
			<p style="font-size: 16px; margin-bottom: 10px;">
				<strong>Task Name:</strong> %s
			</p>
			<p style="font-size: 16px; margin-bottom: 10px;">
				<strong>Time:</strong> %s
			</p>
			<h3 style="color: #d9534f;">Error Details:</h3>
			%s
			<p style="font-style: italic; color: gray; font-size: 14px;">
				This is an automated message. Please do not reply.
			</p>
		</div>
	</body>
	</html>`, taskinfo.GetId(), taskinfo.Name, currentTime, errorDetails)

	// Gửi email
	if err := SendEmail(context.Background(), currentUser.Email, "[vBackup] Task Error Notification", emailBody); err != nil {
		log.Printf("Failed to send error email: %v", err)
	}
}

func findParentSnapshot(ctx context.Context, repo restic.Repository, opts BackupOptions, targets []string, timeStampLimit time.Time) (*restic.Snapshot, error) {
	if opts.Force {
		return nil, nil
	}

	snName := opts.Parent
	if snName == "" {
		snName = "latest"
	}
	f := restic.SnapshotFilter{TimestampLimit: timeStampLimit}
	if opts.GroupBy.Host {
		f.Hosts = []string{opts.Host}
	}
	if opts.GroupBy.Path {
		f.Paths = targets
	}
	if opts.GroupBy.Tag {
		f.Tags = []restic.TagList{opts.Tags.Flatten()}
	}

	sn, _, err := f.FindLatest(ctx, repo.Backend(), repo, snName)
	// Snapshot not found is ok if no explicit parent was set
	if opts.Parent == "" && errors.Is(err, restic.ErrNoSnapshotFound) {
		err = nil
	}
	return sn, err
}
func collectRejectFuncs(opts BackupOptions, targets []string) (fs []RejectFunc, err error) {
	// allowed devices
	if opts.ExcludeOtherFS {
		f, err := rejectByDevice(targets)
		if err != nil {
			return nil, err
		}
		fs = append(fs, f)
	}

	if len(opts.ExcludeLargerThan) != 0 {
		f, err := rejectBySize(opts.ExcludeLargerThan)
		if err != nil {
			return nil, err
		}
		fs = append(fs, f)
	}

	return fs, nil
}

func collectRejectByNameFuncs(opts BackupOptions, repo *repository.Repository) (fs []RejectByNameFunc, err error) {
	// exclude restic cache
	if repo.Cache != nil {
		f, err := rejectResticCache(repo)
		if err != nil {
			return nil, err
		}

		fs = append(fs, f)
	}

	fsPatterns, err := opts.excludePatternOptions.CollectPatterns()
	if err != nil {
		return nil, err
	}
	fs = append(fs, fsPatterns...)

	if opts.ExcludeCaches {
		opts.ExcludeIfPresent = append(opts.ExcludeIfPresent, "CACHEDIR.TAG:Signature: 8a477f597d28d172789f06886806bc55")
	}

	for _, spec := range opts.ExcludeIfPresent {
		f, err := rejectIfPresent(spec)
		if err != nil {
			return nil, err
		}

		fs = append(fs, f)
	}

	return fs, nil
}
