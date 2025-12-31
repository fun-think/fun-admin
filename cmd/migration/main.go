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
	"os"

	"go.uber.org/zap"
)

func main() {
	var envConf = flag.String("conf", "config/local.yml", "config path, eg: -conf ./config/local.yml")
	var allowDrop = flag.Bool("allow-drop", false, "confirm dropping existing tables before migration")
	var adminPasswordFlag = flag.String("admin-password", "", "initial admin password (if empty, a random password will be generated)")
	flag.Parse()
	conf := config.NewConfig(*envConf)

	adminPassword := *adminPasswordFlag
	if adminPassword == "" {
		adminPassword = os.Getenv("ADMIN_INITIAL_PASSWORD")
	}

	log := logger.NewLogger(conf)
	db := repository.NewDB(conf, log)
	enforcer := repository.NewCasbinEnforcer(conf, log, db)
	baseSid, err := sid.NewSid()
	if err != nil {
		log.Fatal("create sid error", zap.Error(err))
	}

	migrateServer := server.NewMigrateServer(db, log, baseSid, enforcer, *allowDrop, adminPassword)

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
