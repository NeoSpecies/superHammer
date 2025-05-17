package ipc

import (
	"fmt"
	"net"
)

// executeSocket 通过Socket执行PHP命令
func ExecuteSocket(command string, socketPath string, uuid string) ([]byte, error) {
	// 创建一个TCP连接到PHP解释器的Socket
	conn, err := net.Dial("unix", socketPath)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to PHP socket: %v", err)
	}
	defer conn.Close()

	// 将命令和uuid写入Socket，这里使用空格分隔
	// fullCommand := fmt.Sprintf("%s", uuid)
	// fmt.Println(fullCommand)
	_, err = fmt.Fprintf(conn, "%s\n", uuid)
	if err != nil {
		return nil, fmt.Errorf("failed to write command and uuid to PHP socket: %v", err)
	}

	// 从Socket读取输出
	buffer := make([]byte, 1024)
	n, err := conn.Read(buffer)
	if err != nil {
		return nil, fmt.Errorf("failed to read output from PHP socket: %v", err)
	}

	return buffer[:n], nil
}
