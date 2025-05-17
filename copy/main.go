package main

import (
	"fmt"
	"log"
	"net"
	"net/http"
	_ "net/http/pprof"
	"os"
	_ "os/exec"
	"os/signal"
	"syscall"

	_ "github.com/mkevac/debugcharts" // 可选，添加后可以查看几个实时图表数据
)

type Route struct {
	Path     string `json:"path"`
	Language string `json:"language"`
	Command  string `json:"command"`
}

type Router struct {
	Routes []Route `json:"routes"`
}

func (r *Router) handleHTTP(w http.ResponseWriter, req *http.Request) {
	log.Println("进入请求")
	// 获取请求路径
	path := req.URL.Path
	log.Println(path)
	// 在路由表中查找对应的配置
	var route Route
	for _, r := range r.Routes {
		if r.Path == path {
			route = r
			break
		}
	}

	if route.Path == "" {
		http.NotFound(w, req)
		return
	}

	// 执行PHP命令通过Socket
	socketPath := "/tmp/mysocket.sock" // 替换为您的PHP解释器Socket路径
	output, err := executePHPViaSocket(route.Command, socketPath)
	if err != nil {
		log.Println("执行PHP命令时出错:", err)
		http.Error(w, "内部服务器错误", http.StatusInternalServerError)
		return
	}

	// 将输出结果返回给客户端
	w.Header().Set("Content-Type", "text/plain")
	w.Write(output)
}

func executePHPViaSocket(command string, socketPath string) ([]byte, error) {
	// 创建一个TCP连接到PHP解释器的Socket
	conn, err := net.Dial("unix", socketPath)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to PHP socket: %v", err)
	}
	defer conn.Close()

	// 将命令写入Socket
	_, err = fmt.Fprintf(conn, "%s\n", command)
	if err != nil {
		return nil, fmt.Errorf("failed to write command to PHP socket: %v", err)
	}

	// 从Socket读取输出
	buffer := make([]byte, 1024)
	n, err := conn.Read(buffer)
	if err != nil {
		return nil, fmt.Errorf("failed to read output from PHP socket: %v", err)
	}

	return buffer[:n], nil
}

func (r *Router) handleSocket(conn net.Conn) {
	// 在此处添加处理 Socket 请求的逻辑
	data := make([]byte, 1024)
	_, err := conn.Read(data)
	if err != nil {
		log.Println("Error reading data from socket:", err)
		return
	}

	// 处理接收到的数据
	// ...

	// 发送响应到 Socket
	response := []byte("Hello from Socket!")
	_, err = conn.Write(response)
	if err != nil {
		log.Println("Error writing data to socket:", err)
		return
	}
}
func cleanup() {
	fmt.Println("执行清理操作...")
	// 在这里执行你的清理操作，例如关闭文件、释放资源等。
	fmt.Println("清理完成，程序已退出。")
}

// WritePidToFile 将当前进程的PID写入指定的文件
func WritePidToFile(pidFile string) {
	// 获取当前进程的PID
	mainPid := os.Getpid()
	fmt.Printf("当前进程的PID为：%d\n", mainPid)

	// 创建或打开PID文件
	file, err := os.Create(pidFile)
	if err != nil {
		log.Fatalf("无法创建PID文件: %s", err)
	}
	defer file.Close()

	// 将PID写入文件
	_, err = file.WriteString(fmt.Sprintf("%d", mainPid))
	if err != nil {
		log.Fatalf("无法写入PID到文件: %s", err)
	}
}

var router Router

func main() {

	// 调用函数，将PID写入文件
	WritePidToFile("duangMain.pid")

	// 创建一个接收信号的通道
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGTERM)

	fmt.Println("程序已启动。按下Ctrl+C或发送SIGTERM信号以终止程序。")

	// 启动 HTTP 服务器
	go func() {
		http.HandleFunc("/", router.handleHTTP)
		fmt.Println("HTTP server listening on port 80...")
		err := http.ListenAndServe(":80", nil)
		if err != nil {
			log.Fatal("Error starting HTTP server:", err)
		}
	}()

	// 启动 Socket 服务器
	go func() {
		listener, err := net.Listen("unix", "/tmp/mysocket.sock")
		if err != nil {
			log.Fatal("Error starting listener:", err)
		}
		defer listener.Close()

		fmt.Println("Listening on Unix socket /tmp/mysocket.sock...")

		for {
			conn, err := listener.Accept()
			if err != nil {
				log.Println("Error accepting connection:", err)
				continue
			}

			go router.handleSocket(conn)
		}
	}()
	// 阻塞等待信号
	sig := <-sigChan
	fmt.Printf("接收到信号：%s\n", sig)
	// 执行清理操作或退出程序
	cleanup()
	//select {} // 阻塞主goroutine，等待信号
}
