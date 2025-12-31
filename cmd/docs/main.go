package main

import (
	"flag"
	"fmt"
	"fun-admin/cmd/bootstrap"
	"fun-admin/pkg/admin"
	"fun-admin/pkg/docs"
	"fun-admin/pkg/logger"
	"os"

	"github.com/casbin/casbin/v2"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

func main() {
	var (
		confPath   = flag.String("conf", "config/local.yml", "config path, eg: -conf ./config/local.yml")
		outputPath = flag.String("output", "docs/swagger.json", "output swagger file path")
	)
	flag.Parse()

	conf := bootstrap.LoadConfig(*confPath)
	log := logger.NewLogger(conf)

	bootstrap.InitI18n()

	containerManager := bootstrap.InitContainer(*confPath)
	cacheManager := bootstrap.InitCache(conf)

	db := containerManager.MustGet("database").(*gorm.DB)
	enforcer := containerManager.MustGet("enforcer").(*casbin.SyncedEnforcer)

	bootstrap.InitAdmin(containerManager, cacheManager, log, db, enforcer)

	if err := generateDocs(*outputPath); err != nil {
		log.Error("生成文档失败", zap.Error(err))
		os.Exit(1)
	}

	log.Info("swagger 生成成功", zap.String("output", *outputPath))
}

func generateDocs(output string) error {
	resources := admin.GlobalResourceManager.GetResources()
	resourceMap := make(map[string]admin.Resource, len(resources))
	for _, resource := range resources {
		resourceMap[resource.GetSlug()] = resource
	}

	generator := docs.NewDocumentGenerator()
	if err := generator.Generate(nil, resourceMap); err != nil {
		return fmt.Errorf("生成文档失败: %w", err)
	}

	if err := generator.SaveToFile(output); err != nil {
		return fmt.Errorf("保存文档失败: %w", err)
	}

	return nil
}
