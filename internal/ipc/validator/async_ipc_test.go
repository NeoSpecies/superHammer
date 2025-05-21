package main

import (
	"encoding/json"
	"log"
	"net"
	"os"
	"time"
)

const socketPath = "/tmp/async_ipc.sock"

type AsyncMessage struct {
	ID       string      `json:"id"`
	Protocol string      `json:"protocol"`
	Data     interface{} `json:"data"`
}

func main() {
	// 清理旧的socket文件
	os.Remove(socketPath)

	// 启动服务端
	go startAsyncServer()
	time.Sleep(100 * time.Millisecond)

	// 测试客户端
	testAsyncClient()
}

func startAsyncServer() {
	listener, err := net.Listen("unix", socketPath)
	if err != nil {
		log.Fatal("Server error:", err)
	}
	defer listener.Close()

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Println("Accept error:", err)
			continue
		}
		go handleAsyncConnection(conn)
	}
}

func handleAsyncConnection(conn net.Conn) {
	defer conn.Close()

	decoder := json.NewDecoder(conn)
	encoder := json.NewEncoder(conn)

	for {
		var msg AsyncMessage
		if err := decoder.Decode(&msg); err != nil {
			log.Println("Decode error:", err)
			return
		}

		// 异步处理模拟
		go func(m AsyncMessage) {
			time.Sleep(100 * time.Millisecond) // 模拟处理耗时
			response := AsyncMessage{
				ID:       m.ID,
				Protocol: m.Protocol,
				Data:     "Processed: " + m.ID,
			}
			encoder.Encode(response)
		}(msg)
	}
}

func testAsyncClient() {
	conn, err := net.Dial("unix", socketPath)
	if err != nil {
		log.Fatal("Client dial error:", err)
	}
	defer conn.Close()

	// 发送测试请求
	msg := AsyncMessage{
		ID:       "test_123",
		Protocol: "http",
		Data:     "test data",
	}

	if err := json.NewEncoder(conn).Encode(msg); err != nil {
		log.Fatal("Client encode error:", err)
	}

	// 等待响应
	var resp AsyncMessage
	if err := json.NewDecoder(conn).Decode(&resp); err != nil {
		log.Fatal("Client decode error:", err)
	}

	log.Printf("Received response: %+v\n", resp)
}
