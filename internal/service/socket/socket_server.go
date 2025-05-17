package socket

import (
	"bigHammer/internal/config"
	ipc "bigHammer/internal/ipc/socket"
	"bigHammer/pkg/utils"
	"context"
	"fmt"
	"log"
	"net"
	"os"
)

// StartSocketServer 启动socket服务器
func StartSocketServer(ctx context.Context) {

	err := config.LoadConfig()
	if err != nil {
		fmt.Println("Error loading config:", err)
		return
	}
	configFilePath, _ := utils.ResolvePath(config.GlobalConfig.SocketPath)
	listener, err := net.Listen("unix", configFilePath)
	if err != nil {
		log.Fatal("Error starting listener:", err)
	}

	// 添加日志输出以确保goroutine正在执行
	go func() {
		log.Println("Starting goroutine to listen for ctx.Done()")
		<-ctx.Done()
		log.Println("Received ctx.Done() signal, closing listener")
		listener.Close()
	}()

	log.Println("Listening on Unix socket ...", configFilePath)

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Println("Listener accept error:", err)
			break
		}
		go ipc.HandleSocket(conn)
	}

	log.Println("Shutting down socket server...")
}

// StopSocketServer 停止socket服务器
func StopSocketServer() {
	err := config.LoadConfig()
	if err != nil {
		fmt.Println("Error loading config:", err)
		return
	}
	configFilePath, _ := utils.ResolvePath(config.GlobalConfig.SocketPath)

	// 尝试关闭监听器
	_, err = net.Dial("unix", configFilePath)
	if err != nil {
		log.Println("Socket server is not running")
		return
	}

	// 关闭监听器
	listener, err := net.FileListener(os.NewFile(3, configFilePath))
	if err != nil {
		log.Println("Error getting listener:", err)
		return
	}
	err = listener.Close()
	if err != nil {
		log.Println("Error closing listener:", err)
	} else {
		log.Println("Listener closed successfully")
	}

	// 删除socket文件
	err = os.Remove(configFilePath)
	if err != nil {
		log.Println("Error removing socket file:", err)
	} else {
		log.Println("Socket file removed successfully")
	}
}
