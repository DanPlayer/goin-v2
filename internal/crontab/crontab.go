package crontab

import (
	"github.com/robfig/cron/v3"
)

func InitCron() *cron.Cron {
	c := cron.New(cron.WithSeconds())

	c.Start()
	return c
}
