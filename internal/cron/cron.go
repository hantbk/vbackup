package cron

import (
	"fmt"
	"strings"
	"time"

	"github.com/hantbk/vbackup/internal/api/v1/policy"
	"github.com/hantbk/vbackup/internal/api/v1/task"
	"github.com/hantbk/vbackup/internal/consts"
	"github.com/hantbk/vbackup/internal/server"
	resticProxy "github.com/hantbk/vbackup/restic_proxy"
	"github.com/robfig/cron/v3"
)

var c *cron.Cron

func InitCron() {
	c = cron.New(cron.WithSeconds())
	initSystemCronJob()
	c.Start()
	defer c.Stop()
	select {}
}

// initSystemCronJob - Initialize system cron jobs
func initSystemCronJob() {
	// Prepare homepage data
	_, err := c.AddJob("0 0 0 * * *", SystemJob(func() {
		server.Logger().Info("Preparing homepage data")
		go resticProxy.GetAllRepoStats()
	}))
	if err != nil {
		fmt.Println(fmt.Errorf("GetAllRepoStats cron job failed to start: %s", err))
	}
	// Clean up running tasks
	_, err = c.AddJob("0 */10 * * * *", SystemJob(func() {
		go task.ClearTaskRunning()
	}))
	if err != nil {
		fmt.Println(fmt.Errorf("ClearTaskRunning cron job failed to start: %s", err))
	}
	// Execute cleanup policy
	_, err = c.AddJob("0 0 6 * * *", SystemJob(func() {
		server.Logger().Info("Executing cleanup policy")
		go policy.DoPolicy()
	}))
	if err != nil {
		fmt.Println(fmt.Errorf("DoPolicy cron job failed to start: %s", err))
	}
}

func AddJob(cronStr string, job BaseJob) error {
	cronStr = CheckCron(cronStr)
	_, err := c.AddJob(cronStr, job)
	if err != nil {
		return err
	}
	return nil
}

func ClearJob() {
	entries := c.Entries()
	for _, entry := range entries {
		if entry.Job.(BaseJob).GetType() == JOB_TYPE_SYSTEM {
			continue
		}
		c.Remove(entry.ID)
	}
}

func CheckCron(cronStr string) string {
	ts := strings.Fields(cronStr)
	if len(ts) > 6 {
		var res string
		for _, s := range ts[:6] {
			res += s + " "
		}
		return res
	} else {
		return cronStr
	}
}

// GetNextTimes - Generate the list of next execution times
func GetNextTimes(cronStr string) ([]string, error) {
	cronStr = CheckCron(cronStr)
	res := make([]string, 0)
	tmpcron := cron.New(cron.WithSeconds())
	entryID, err := tmpcron.AddFunc(cronStr, func() {

	})
	if err != nil {
		return nil, err
	}
	entry := tmpcron.Entry(entryID)
	nexttime := time.Now()
	for i := 0; i < 5; i++ {
		nexttime = entry.Schedule.Next(nexttime)
		res = append(res, nexttime.Format(consts.Custom))
	}
	return res, nil
}
