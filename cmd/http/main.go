package main

import (
	"kerdo.dev/taavi/pkg/app"
	"kerdo.dev/taavi/pkg/data"
	"kerdo.dev/taavi/pkg/external"
	"kerdo.dev/taavi/pkg/logger"
	"kerdo.dev/taavi/pkg/scheduler"
	"kerdo.dev/taavi/pkg/zulip"
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
