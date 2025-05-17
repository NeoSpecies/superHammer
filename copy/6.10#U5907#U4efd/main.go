package main

import (
	"bigHammer/internal/config"
	"bigHammer/internal/di"
	"bigHammer/internal/ipc/cmd"
	"bigHammer/internal/plugins/agilitymemdb"
	"bigHammer/internal/service/http"
	"bigHammer/internal/service/socket"
	"bigHammer/pkg/utils"
	"bytes"
	"context"
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"sync"
	"syscall"
)

// WritePidToFile 将当前进程的PID写入指定的文件
func WritePidToFile() {
	err := config.LoadConfig()
	if err != nil {
		fmt.Println("Error loading config:", err)
		return
	}
	mainPid := os.Getpid()
	fmt.Printf("当前进程的PID为：%d\n", mainPid)
	config_file_path, _ := utils.ResolvePath(config.GlobalConfig.MainPIDPath)
	file, err := os.Create(config_file_path)
	if err != nil {
		fmt.Printf("无法创建PID文件: %s", err)
	}
	defer file.Close()
	_, err = file.WriteString(fmt.Sprintf("%d", mainPid))
	if err != nil {
		fmt.Printf("无法写入PID到文件: %s", err)
	}
}

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

func cleanup() {
	fmt.Println("执行清理操作...")
	err := config.LoadConfig()
	if err != nil {
		fmt.Println("Error loading config:", err)
		return
	}
	mainPid_path, _ := utils.ResolvePath(config.GlobalConfig.MainPIDPath)
	socket_path, _ := utils.ResolvePath(config.GlobalConfig.SocketPath)
	bussinessPid_path, _ := utils.ResolvePath(config.GlobalConfig.BussinessPIDPath)
	err = TerminateProcessByPIDFile(bussinessPid_path)
	if err != nil {
		fmt.Println("操作失败:", err)
		return
	}
	fmt.Println("成功终止进程")
	socket.StopSocketServer()
	filesToDelete := []string{mainPid_path, socket_path, bussinessPid_path}
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

func startup(ctx context.Context) error {
	bussinessMainPath, err := utils.ResolvePath(config.GlobalConfig.BussinessMainPath)
	if err != nil {
		return fmt.Errorf("解析路径时发生错误: %w", err)
	}
	cmd := exec.CommandContext(ctx, "php", bussinessMainPath)
	var output, errOutput bytes.Buffer
	cmd.Stdout = &output
	cmd.Stderr = &errOutput
	err = cmd.Start()
	if err != nil {
		return fmt.Errorf("启动进程时发生错误: %w, 输出: %s, 错误输出: %s", err, output.String(), errOutput.String())
	}
	fmt.Printf("PHP脚本的PID: %d\n", cmd.Process.Pid)
	bussinessPIDPath, err := utils.ResolvePath(config.GlobalConfig.BussinessPIDPath)
	if err != nil {
		return fmt.Errorf("解析业务PID文件路径时发生错误: %w", err)
	}
	file, err := os.Create(bussinessPIDPath)
	if err != nil {
		return fmt.Errorf("创建业务PID文件时发生错误: %w", err)
	}
	defer file.Close()
	_, err = file.WriteString(fmt.Sprintf("%d", cmd.Process.Pid))
	if err != nil {
		return fmt.Errorf("写入业务PID到文件时发生错误: %w", err)
	}
	select {
	case <-ctx.Done():
		if cmd.Process != nil {
			err = cmd.Process.Signal(os.Interrupt)
			if err != nil {
				fmt.Printf("无法发送中断信号: %v\n", err)
			}
		}
		return ctx.Err()
	default:
	}
	return nil
}

func main() {
	err := config.LoadConfig()
	if err != nil {
		fmt.Printf("Error loading config: %v", err)
	}
	mainPIDPath, _ := utils.ResolvePath(config.GlobalConfig.BussinessPIDPath)
	businessPath, _ := utils.ResolvePath(config.GlobalConfig.BussinessMainPath)
	memoryDBPath, _ := utils.ResolvePath(config.GlobalConfig.MemoryDBPath)
	watcherPath, _ := utils.ResolvePath("/Develop")
	container := di.NewContainer()
	err = container.Register("database", func() interface{} {
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
	dbService, err := container.Resolve("database")
	if err != nil {
		fmt.Printf("解析数据库服务失败: %v\n", err)
		os.Exit(1)
	}
	db, ok := dbService.(*agilitymemdb.AgilityMemDB)
	if !ok {
		fmt.Println("Failed to assert dbService as *agilitymemdb.AgilityMemDB")
		os.Exit(1)
	}
	WritePidToFile()
	ctx, cancel := context.WithCancel(context.Background())
	var wg sync.WaitGroup
	watcher := utils.New([]string{watcherPath}, mainPIDPath, businessPath)
	wg.Add(1)
	go func() {
		defer wg.Done()
		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
		sig := <-sigChan
		fmt.Printf("接收到信号：%s\n", sig)
		signal.Stop(sigChan)
		cancel()
	}()
	wg.Add(1)
	go func() {
		defer wg.Done()
		socket.StartSocketServer(ctx, db)
	}()
	wg.Add(1)
	go func() {
		defer wg.Done()
		http.StartHTTPServer(ctx, db)
	}()
	//启动phpsocket监听
	// wg.Add(1)
	// go func() {
	// 	defer wg.Done()
	// 	err := startup(ctx)
	// 	if err != nil {
	// 		fmt.Printf("Startup failed: %v\n", err)
	// 	}
	// }()
	wg.Add(1)
	go func() {
		defer wg.Done()
		watcher.Start(ctx)
	}()
	wg.Wait()
	cleanup()
}
