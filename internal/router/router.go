package router

import (
	"bigHammer/internal/config"
	"bigHammer/internal/interface/database"
	"bigHammer/pkg/utils"
	"encoding/json"
	"fmt"
	"os"
)

type Route struct {
	Path     string `json:"path"`
	Language string `json:"language"`
	Command  string `json:"command"`
}

type Router struct {
	Routes []Route `json:"routes"`
	DB     database.IDatabase
}

func NewRouter(db database.IDatabase) *Router {
	return &Router{DB: db}
}

// LoadRouterConfig 从文件加载并解析路由配置
func LoadRouterConfig() (Router, error) {
	var router Router
	err := config.LoadConfig()
	if err != nil {
		return router, fmt.Errorf("error loading config: %v", err)
	}
	config_file_path, err := utils.ResolvePath(config.GlobalConfig.RouterPath)
	if err != nil {
		return router, fmt.Errorf("error resolving config path: %v", err)
	}
	file, err := os.Open(config_file_path)
	if err != nil {
		return router, fmt.Errorf("failed to open router configuration file: %w", err)
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	err = decoder.Decode(&router)
	if err != nil {
		return router, fmt.Errorf("failed to decode router configuration: %w", err)
	}

	return router, nil
}
