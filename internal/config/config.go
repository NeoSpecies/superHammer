package config

import (
	"bigHammer/pkg/utils"
	"encoding/json"
	"fmt"
	"os"
)

// Config 结构体定义了所有的配置项
// 包含系统运行所需的所有配置参数，如数据库连接、文件路径、端口等
type Config struct {
	// DatabaseURL 数据库连接URL
	// 用于连接主数据库的完整URL字符串
	DatabaseURL         string      `json:"database_url"`
	// MemoryDBPath 内存数据库路径
	// 用于存储内存数据库文件的路径
	MemoryDBPath        string      `json:"memory_db_path"`
	// MainPIDPath 主进程PID文件路径
	// 用于存储主进程PID的文件路径
	MainPIDPath         string      `json:"mainPid_path"`
	// BussinessPIDPath 业务进程PID文件路径
	// 用于存储业务进程PID的文件路径
	BussinessPIDPath    string      `json:"bussinessPid_path"`
	// AttachmentStorage 附件存储路径
	// 用于存储系统附件的目录路径
	AttachmentStorage   string      `json:"attachment_storage"`
	// SocketPath Socket文件路径
	// 用于Unix Domain Socket通信的文件路径
	SocketPath          string      `json:"socket_path"`
	// BussinessSocketPath 业务Socket文件路径
	// 用于业务进程Unix Domain Socket通信的文件路径
	BussinessSocketPath string      `json:"bussiness_socket_path"`
	// RouterPath 路由配置文件路径
	// 用于存储路由配置的文件路径
	RouterPath          string      `json:"router_path"`
	// PluginsPath 插件目录路径
	// 用于存储系统插件的目录路径
	PluginsPath         string      `json:"plugins_path"`
	// BussinessMainPath 业务主程序路径
	// 用于存储业务主程序的文件路径
	BussinessMainPath   string      `json:"bussiness_main_path"`
	// Ports 端口配置
	// 包含系统使用的所有端口配置
	Ports               PortsConfig `json:"ports"`
}

// PortsConfig 定义了端口配置
// 包含系统各个服务使用的端口号
type PortsConfig struct {
	// IPCPort IPC服务端口
	// 用于进程间通信的端口号
	IPCPort       string `json:"ipc_port"`
	// HTTPPort HTTP服务端口
	// 用于HTTP服务的端口号
	HTTPPort      string `json:"http_port"`
	// WebSocketPort WebSocket服务端口
	// 用于WebSocket服务的端口号
	WebSocketPort string `json:"websocket_port"`
	// TCPPort TCP服务端口
	// 用于TCP服务的端口号
	TCPPort       string `json:"tcp_port"`
}

// GlobalConfig 全局配置变量
// 存储系统运行时的配置信息
var GlobalConfig *Config

// LoadConfig 从文件加载配置到全局变量
// 功能：
// 1. 解析配置文件路径
// 2. 检查配置文件是否存在
// 3. 读取配置文件内容
// 4. 解析JSON配置到结构体
// 5. 更新全局配置变量
// 参数：无
// 返回值：
//   - error: 加载过程中的错误信息，如果成功则返回nil
func LoadConfig() error {
	config_file_path, err := utils.ResolvePath("/config/config.json")
	if err != nil {
		return fmt.Errorf("error resolving config path: %v", err)
	}

	if _, err := os.Stat(config_file_path); os.IsNotExist(err) {
		return fmt.Errorf("error loading config: %s does not exist", config_file_path)
	}

	content, err := os.ReadFile(config_file_path)
	if err != nil {
		return fmt.Errorf("error loading config: %v", err)
	}

	var config Config
	err = json.Unmarshal(content, &config)
	if err != nil {
		return fmt.Errorf("error parsing config: %v", err)
	}

	GlobalConfig = &config
	return nil
}

// SaveConfig 将全局配置保存到文件
// 功能：
// 1. 将配置结构体转换为JSON格式
// 2. 将JSON数据写入配置文件
// 参数：无
// 返回值：
//   - error: 保存过程中的错误信息，如果成功则返回nil
func SaveConfig() error {
	content, err := json.MarshalIndent(GlobalConfig, "", "    ")
	if err != nil {
		return fmt.Errorf("error marshalling config: %v", err)
	}
	return os.WriteFile("config/config.json", content, 0644)
}
