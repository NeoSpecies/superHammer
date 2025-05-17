package ipc

import (
	"bigHammer/internal/interface/database"
	"encoding/json"
	"log"
	"net"
	"strings"
	"time"
)

var keepAliveInterval = 30 * time.Second // 设置心跳间隔
// 将请求信息组装到map中
type Response struct {
	Status  int         `json:"status"`  // 状态码
	Message string      `json:"message"` // 信息描述
	Data    interface{} `json:"data"`    // 数据载体，可以是任何类型
}

func HandleSocket(conn net.Conn, db database.IDatabase) {
	defer func() {
		// 延迟关闭连接
		conn.Close()
	}()

	// 读取数据
	data := make([]byte, 1024)
	n, err := conn.Read(data)
	if err != nil {
		log.Println("Error reading data from socket:", err)
		return
	}
	log.Printf("Received data: %s\n", data)
	// 处理接收到的数据
	// 这里可以添加具体的处理逻辑
	// 创建 Response 结构体实例
	response := Response{
		Status:  200,                                 // OK 状态码
		Message: "Hello from Socket!",                // 自定义消息
		Data:    map[string]string{"uri": "/route1"}, // 数据部分，可以自定义
	}
	// 在尝试获取值之前，打印键的值来确认其格式
	keyStr := string(data[:n])
	keyStr = strings.TrimSpace(keyStr) // 如果有疑虑，去除可能的首尾空格
	// 然后使用修正后的键来获取数据
	value, exists := db.Get(keyStr)
	if exists {
		response.Data = value
	} else {
		response.Data = map[string]string{}
	}

	// 将 Response 结构体转换成 JSON 字符串
	responseJSON, _ := json.Marshal(response)
	// 发送响应到 Socket
	// response := []byte("Hello from Socket!")
	_, err = conn.Write(responseJSON)
	if err != nil {
		log.Println("Error writing data to socket:", err)
		return
	}

	// 启动心跳机制
	go sendHeartbeat(conn)

	// 主循环，持续读取数据
	for {
		_, err := conn.Read(data)
		if err != nil {
			log.Println("Connection closed by peer")
			return
		}

		// 根据接收到的数据或处理逻辑来决定是否保持连接打开
		// 例如，如果接收到的数据中包含特定的命令，可以设置 keepConnectionOpen 为 true
		// keepConnectionOpen := true
	}
}

func sendHeartbeat(conn net.Conn) {
	ticker := time.NewTicker(keepAliveInterval)
	defer ticker.Stop()

	for range ticker.C {
		// 检查连接是否仍然打开
		if conn == nil {
			log.Println("Connection is nil, stopping heartbeat.")
			return
		}

		// 发送心跳包
		_, err := conn.Write([]byte("heartbeat"))
		if err != nil {
			log.Println("Error sending heartbeat:", err)
			// 如果发送失败，检查是否是连接已关闭
			if err, ok := err.(*net.OpError); ok && err.Op == "write" && err.Err.Error() == "use of closed network connection" {
				log.Println("Connection is closed, stopping heartbeat.")
				return
			}
		}
	}
}
