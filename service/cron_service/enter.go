package cron_service

import (
	"Blog_server/service/redis_service/redis_count"
	"github.com/robfig/cron/v3"
	"time"
)

func CronInit() {
	timezone, _ := time.LoadLocation("Asia/Shanghai")
	Cron := cron.New(cron.WithSeconds(), cron.WithLocation(timezone))
	Cron.AddFunc("*/10 * * * * *", redis_count.Update)
	Cron.AddFunc("*/10 * * * * *", redis_count.UpdateToDB)
	Cron.Start()
}
