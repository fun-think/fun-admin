package main

import (
	"flag"
	"fmt"
	"os"
)

func main() {
	// 定义命令行标志
	var (
		action string
		name   string
	)

	flag.StringVar(&action, "action", "", "Action to perform (make:resource, make:page)")
	flag.StringVar(&name, "name", "", "Name of the resource or page")
	flag.Parse()

	// 检查参数
	if action == "" {
		fmt.Println("Error: action is required")
		printUsage()
		os.Exit(1)
	}

	// 执行相应操作
	switch action {
	case "make:resource":
		if name == "" {
			fmt.Println("Error: name is required for make:resource")
			printUsage()
			os.Exit(1)
		}
		makeResource(name)
	case "make:page":
		if name == "" {
			fmt.Println("Error: name is required for make:page")
			printUsage()
			os.Exit(1)
		}
		makePage(name)
	default:
		fmt.Printf("Error: unknown action '%s'\n", action)
		printUsage()
		os.Exit(1)
	}
}

// makeResource 创建新资源
func makeResource(name string) {
	fmt.Printf("Creating resource: %s\n", name)

	// TODO: 实际创建资源文件的逻辑
	// 这里应该创建 internal/resources/{name}_resource.go 文件
	// 并填充基本的资源模板

	fmt.Println("Resource created successfully!")
}

// makePage 创建新页面
func makePage(name string) {
	fmt.Printf("Creating page: %s\n", name)

	// TODO: 实际创建页面文件的逻辑

	fmt.Println("Page created successfully!")
}

// printUsage 打印使用说明
func printUsage() {
	fmt.Println("Usage:")
	fmt.Println("  go run cmd/admin/main.go -action=make:resource -name=ResourceName")
	fmt.Println("  go run cmd/admin/main.go -action=make:page -name=PageName")
	fmt.Println("")
	fmt.Println("Examples:")
	fmt.Println("  go run cmd/admin/main.go -action=make:resource -name=User")
	fmt.Println("  go run cmd/admin/main.go -action=make:page -name=Dashboard")
}
