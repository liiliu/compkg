package ai

import (
	"fmt"
	"os"
	"weihu_server/library/cache"
	"weihu_server/library/config"
	"weihu_server/library/cos"
	"weihu_server/library/database"
	"weihu_server/library/language"
	"weihu_server/library/logger"
	"weihu_server/library/oss"
	"weihu_server/library/refresh"
)

func InitTest() {
	// 读取环境变量
	projectDir := os.Getenv("WEIHU_PROJECT_DIR")
	if projectDir == "" {
		fmt.Println("Environment variable WEIHU_PROJECT_DIR is not set")
		return
	}

	// 设置工作目录为项目根目录
	err := os.Chdir(projectDir)
	if err != nil {
		fmt.Printf("Failed to change working directory to project root: %v", err)
		os.Exit(1)
	}

	config.Init("test")
	logger.Init()
	database.InitialPostgresql()
	oss.Initial()
	cache.Initial()
	refresh.BaseData()
	cos.Initial()
	language.Initial()
}
