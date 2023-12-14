package main

import (
	"kerdo.dev/taavi/app"
	"kerdo.dev/taavi/data"
	"kerdo.dev/taavi/external"
	"kerdo.dev/taavi/logger"
	"kerdo.dev/taavi/scheduler"
	"kerdo.dev/taavi/zulip"
)

func main() {
	app.InitEnv()
	logger.Init()

	data.Init()
	zulip.Init()
	external.Init()

	scheduler.Init()
	scheduler.ScheduleMaster()
	scheduler.ScheduleAll()

	logger.Infow("starting Taavi", nil)

	listenAddr := app.GetListenAddr()
	app := app.GetApp()

	if err := app.Listen(listenAddr); err != nil {
		logger.Fatalw("error starting Taavi", logger.M{
			"err": err.Error(),
		})
	}
}
