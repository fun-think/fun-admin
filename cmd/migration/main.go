package main

import (
	"context"
	"flag"
	"fun-admin/internal/repository"
	"fun-admin/internal/server"
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
	db := repository.NewDB(conf, log)
	enforcer := repository.NewCasbinEnforcer(conf, log, db)
	baseSid := sid.NewSid()

	migrateServer := server.NewMigrateServer(db, log, baseSid, enforcer)

	// 创建应用
	migrationApp := app.NewApp(
		app.WithServer(migrateServer),
		app.WithName("fun-migration"),
	)

	// 运行应用
	if err := migrationApp.Run(context.Background()); err != nil {
		log.Error("migration run error", zap.Error(err))
	}
}
