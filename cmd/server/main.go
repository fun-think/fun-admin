package main

import (
	"context"
	"flag"
	"fmt"
	"fun-admin/cmd/bootstrap"
	"fun-admin/pkg/app"
	"fun-admin/pkg/logger"
	"fun-admin/pkg/server"

	"github.com/casbin/casbin/v2"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

// @title           Nunu Example API
// @version         1.0.0
// @description     This is a sample server celler server.
// @termsOfService  http://swagger.io/terms/
// @contact.name   API Support
// @contact.url    http://www.swagger.io/support
// @contact.email  support@swagger.io
// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html
// @host      localhost:8000
// @securityDefinitions.apiKey Bearer
// @in header
// @name Authorization
// @externalDocs.description  OpenAPI
// @externalDocs.url          https://swagger.io/resources/open-api/
func main() {
	var envConf = flag.String("conf", "config/local.yml", "config path, eg: -conf ./config/local.yml")
	flag.Parse()

	conf := bootstrap.LoadConfig(*envConf)
	log := logger.NewLogger(conf)

	bootstrap.EnsureStorage(log)
	bootstrap.InitI18n()

	containerManager := bootstrap.InitContainer(*envConf)
	cacheManager := bootstrap.InitCache(conf)

	db := containerManager.MustGet("database").(*gorm.DB)
	syncedEnforcer := containerManager.MustGet("enforcer").(*casbin.SyncedEnforcer)

	bootstrap.InitAdmin(containerManager, cacheManager, log, db, syncedEnforcer)

	httpServer := containerManager.MustGet("http_server").(server.Server)
	jobServer := containerManager.MustGet("job_server").(server.Server)

	servers := []server.Server{httpServer, jobServer}

	appManager := app.NewApp(
		app.WithServer(servers...),
		app.WithName("fun-server"),
	)

	log.Info("server start", zap.String("host", fmt.Sprintf("http://%s:%d", conf.GetString("http.host"), conf.GetInt("http.port"))))
	log.Info("docs addr", zap.String("addr", fmt.Sprintf("http://%s:%d/swagger/index.html", conf.GetString("http.host"), conf.GetInt("http.port"))))
	if err := appManager.Run(context.Background()); err != nil {
		panic(err)
	}
}
