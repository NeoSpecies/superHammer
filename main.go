package main

import (
	"bigHammer/internal/config"
	"bigHammer/internal/di"
	"bigHammer/internal/ipc/cmd"
	"bigHammer/internal/plugin"
	"bigHammer/internal/plugin/agilitymemdb"
	"bigHammer/internal/service/http"
	"bigHammer/internal/service/socket"
	"bigHammer/internal/shared"
	"bigHammer/pkg/utils"
	"context"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

// globalConfig 全局配置变量
// 存储系统运行所需的配置信息，包括：
// - 主进程PID文件路径
// - 业务进程PID文件路径
// - Socket文件路径
// - 业务主程序路径
// - 内存数据库路径
var globalConfig *config.Config

// init 初始化函数
// 在main函数执行前自动调用，负责：
// 1. 加载配置文件
// 2. 初始化全局配置变量
// 3. 如果配置加载失败，程序将直接退出
func init() {
	err := config.LoadConfig()
	if err != nil {
		fmt.Println("Error loading config:", err)
		os.Exit(1)
	}
	globalConfig = config.GlobalConfig
}

// WritePidToFile 将当前进程的PID写入指定的文件
// 功能：
// 1. 获取当前进程PID
// 2. 解析PID文件路径
// 3. 创建PID文件
// 4. 将PID写入文件
// 5. 确保文件正确关闭
// 参数：无
// 返回值：无
func WritePidToFile() {
	mainPid := os.Getpid()
	fmt.Printf("当前进程的PID为：%d\n", mainPid)
	config_file_path, err := utils.ResolvePath(globalConfig.MainPIDPath)
	if err != nil {
		fmt.Printf("无法解析PID文件路径: %s", err)
		return
	}
	file, err := os.Create(config_file_path)
	if err != nil {
		fmt.Printf("无法创建PID文件: %s", err)
		return
	}
	defer func() {
		if err := file.Close(); err != nil {
			fmt.Printf("无法关闭PID文件: %s", err)
		}
	}()
	_, err = file.WriteString(fmt.Sprintf("%d", mainPid))
	if err != nil {
		fmt.Printf("无法写入PID到文件: %s", err)
	}
}

// TerminateProcessByPIDFile 通过PID文件终止进程
// 功能：
// 1. 读取PID文件内容
// 2. 验证PID不为空
// 3. 使用kill命令终止进程
// 4. 检查终止结果
// 参数：
//   - pidFilePath string: PID文件的路径
// 返回值：
//   - error: 终止过程中的错误信息，如果成功则返回nil
func TerminateProcessByPIDFile(pidFilePath string) error {
	pidBytes, err := os.ReadFile(pidFilePath)
	if err != nil {
		return fmt.Errorf("读取 PID 文件时发生错误: %w", err)
	}
	pid := string(pidBytes)
	if pid == "" {
		return fmt.Errorf("PID 文件为空")
	}
	output, errOutput, err := cmd.ExecCommand("kill", pid)
	if err != nil {
		return fmt.Errorf("终止进程时发生错误: %w, 输出: %s, 错误输出: %s", err, output, errOutput)
	}
	if output == "" && errOutput == "" {
		return nil
	}
	return fmt.Errorf("进程终止时出现问题, 输出: %s, 错误输出: %s", output, errOutput)
}

// cleanup 清理函数
// 功能：
// 1. 终止业务进程
// 2. 停止Socket服务器
// 3. 删除临时文件（PID文件、Socket文件等）
// 4. 确保所有资源被正确释放
// 参数：无
// 返回值：无
func cleanup() {
	fmt.Println("执行清理操作...")
	socket_path, err := utils.ResolvePath(globalConfig.SocketPath)
	if err != nil {
		fmt.Printf("无法解析Socket路径: %s", err)
		return
	}
	bussinessPid_path, err := utils.ResolvePath(globalConfig.BussinessPIDPath)
	if err != nil {
		fmt.Printf("无法解析业务PID路径: %s", err)
		return
	}
	err = TerminateProcessByPIDFile(bussinessPid_path)
	if err != nil {
		fmt.Println("操作失败:", err)
		return
	}
	fmt.Println("成功终止进程")
	socket.StopSocketServer()
	filesToDelete := []string{globalConfig.MainPIDPath, socket_path, bussinessPid_path}
	for _, filePath := range filesToDelete {
		err := os.Remove(filePath)
		if err != nil {
			fmt.Printf("无法删除文件 %s: %s\n", filePath, err)
		} else {
			fmt.Printf("文件 %s 已成功删除\n", filePath)
		}
	}
	fmt.Println("清理完成，程序已退出。")
}

// main 程序入口函数
// 功能：
// 1. 初始化配置和路径
//    - 解析主PID路径
//    - 解析业务路径
//    - 解析内存数据库路径
//    - 解析监视器路径
// 2. 初始化依赖注入容器
//    - 创建全局容器实例
//    - 注册插件服务
//    - 注册数据库服务
// 3. 启动服务组件
//    - 写入PID文件
//    - 创建上下文和等待组
//    - 创建文件监视器
//    - 启动信号处理
//    - 启动Socket服务器
//    - 启动HTTP服务器
//    - 启动文件监视器
// 4. 处理信号和优雅退出
//    - 等待所有goroutine完成
//    - 执行清理操作
// 参数：无
// 返回值：无
func main() {
	// 解析各种路径配置
	mainPIDPath, err := utils.ResolvePath(globalConfig.BussinessPIDPath)
	if err != nil {
		fmt.Printf("无法解析主PID路径: %s", err)
		os.Exit(1)
	}
	businessPath, err := utils.ResolvePath(globalConfig.BussinessMainPath)
	if err != nil {
		fmt.Printf("无法解析业务路径: %s", err)
		os.Exit(1)
	}
	memoryDBPath, err := utils.ResolvePath(globalConfig.MemoryDBPath)
	if err != nil {
		fmt.Printf("无法解析内存数据库路径: %s", err)
		os.Exit(1)
	}
	watcherPath, err := utils.ResolvePath("/Develop")
	if err != nil {
		fmt.Printf("无法解析监视器路径: %s", err)
		os.Exit(1)
	}

	// 初始化依赖注入容器
	shared.GlobalContainer = di.NewContainer()
	
	// 注册插件服务
	err = shared.GlobalContainer.Register("plugin", func() interface{} {
		return &plugin.PluginDispatcher{}
	}, di.Singleton)
	if err != nil {
		fmt.Printf("注册插件服务失败: %v\n", err)
		os.Exit(1)
	}

	// 注册数据库服务
	err = shared.GlobalContainer.Register("database", func() interface{} {
		db := agilitymemdb.NewAgilityMemDB(memoryDBPath)
		if err := db.LoadData(); err != nil {
			fmt.Printf("加载数据失败: %v\n", err)
			os.Exit(1)
		}
		return db
	}, di.Singleton)
	if err != nil {
		fmt.Printf("注册数据库服务失败: %v\n", err)
		os.Exit(1)
	}

	// 写入PID文件
	WritePidToFile()

	// 创建上下文和等待组
	ctx, cancel := context.WithCancel(context.Background())
	var wg sync.WaitGroup

	// 创建文件监视器
	watcher := utils.New([]string{watcherPath}, mainPIDPath, businessPath)

	// 启动信号处理
	wg.Add(1)
	go func() {
		defer wg.Done()
		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
		<-sigChan
		fmt.Println("接收到退出信号，开始清理...")
		cancel()
	}()

	// 启动Socket服务器
	wg.Add(1)
	go func() {
		defer wg.Done()
		socket.StartSocketServer(ctx)
	}()

	// 启动HTTP服务器
	httpPort := globalConfig.Ports.HTTPPort  // 修正字段名并获取HTTP端口
	fmt.Println("HTTP server listening on port ：", httpPort)
	wg.Add(1)
	go func() {
		defer wg.Done()
		http.StartHTTPServer(ctx,httpPort)
	}()

	// 启动文件监视器
	wg.Add(1)
	go func() {
		defer wg.Done()
		watcher.Start(ctx)
	}()

	// 等待所有goroutine完成
	wg.Wait()

	// 执行清理操作
	cleanup()
}
