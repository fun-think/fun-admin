package main

import (
	"context"
	"flag"
	"fun-admin/internal/repository"
	"fun-admin/internal/server"
	"fun-admin/internal/task"
	"fun-admin/pkg/app"
	"fun-admin/pkg/config"
	"fun-admin/pkg/logger"
	"fun-admin/pkg/sid"

	"go.uber.org/zap"
)

func main() {
	var envConf = flag.String("conf", "config/local.yml", "config path, eg: -conf ./config/local.yml")
	flag.Parse()
	conf := config.NewConfig(*envConf)

	log := logger.NewLogger(conf)
	log.Info("start task")

	// 创建依赖与仓库
	db := repository.NewDB(conf, log)
	repo := repository.NewRepository(log, db, nil)
	transaction := repository.NewTransaction(repo)
	userRepo := repository.NewUserRepository(log, db, nil)

	// 创建任务组件
	baseSid := sid.NewSid()
	baseTask := task.NewTask(transaction, log, baseSid)
	userTask := task.NewUserTask(baseTask, userRepo)
	taskServer := server.NewTaskServer(log, userTask)

	// 创建应用并运行
	taskApp := app.NewApp(
		app.WithServer(taskServer),
		app.WithName("fun-task"),
	)

	if err := taskApp.Run(context.Background()); err != nil {
		log.Error("task run error", zap.Error(err))
	}
}
