package ipc

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"net"
	"time"

	"github.com/google/uuid" // 新增UUID生成库
)

// 协议常量定义（与设计文档v1.1一致）
const (
	ProtocolVersion = 0x0101           // v1.1协议版本
	MsgTypeSync     = 0x01             // 同步消息类型
	MsgTypeAsyncReq = 0x04             // 异步请求类型
	AsyncTimeout    = 30 * time.Second // 异步超时时间
	// 移除重复声明的 MaxPayloadSize，直接使用 socket_receive.go 中已定义的常量
)

// 同步请求结构体（用于构造负载）
type SyncRequest struct {
	Method string      `json:"method"`
	Params interface{} `json:"params"`
}

// 通用发送函数（支持同步/异步）
func TransmitIPC(isAsync bool, method string, params interface{}, socketPath string) ([]byte, string, error) {
	// 创建Unix Socket连接
	conn, err := net.Dial("unix", socketPath)
	if err != nil {
		return nil, "", fmt.Errorf("连接PHP Socket失败: %v", err)
	}
	defer conn.Close()

	var payload []byte
	var msgType byte
	var requestID string

	// 构造请求负载和消息类型
	if isAsync {
		// 生成UUIDv4作为异步请求ID
		requestID = uuid.New().String()
		asyncReq := AsyncRequest{
			ID:     requestID,
			Method: method,
			Params: params,
		}
		payload, err = json.Marshal(asyncReq)
		msgType = MsgTypeAsyncReq
	} else {
		syncReq := SyncRequest{
			Method: method,
			Params: params,
		}
		payload, err = json.Marshal(syncReq)
		msgType = MsgTypeSync
	}
	if err != nil {
		return nil, "", fmt.Errorf("序列化请求失败: %v", err)
	}

	// 校验负载大小（直接使用 socket_receive.go 中已定义的 MaxPayloadSize）
	if len(payload) > MaxPayloadSize {
		return nil, "", fmt.Errorf("负载大小超出限制（最大4MB）")
	}

	// 封装协议头（2字节版本 + 1字节类型 + 4字节负载长度）
	header := make([]byte, 7)
	binary.BigEndian.PutUint16(header[:2], ProtocolVersion)
	header[2] = msgType
	binary.BigEndian.PutUint32(header[3:7], uint32(len(payload)))
	fmt.Println(string(append(header, payload...)))
	// 发送完整消息（头+负载）
	_, err = conn.Write(append(header, payload...))
	if err != nil {
		return nil, "", fmt.Errorf("发送请求失败: %v", err)
	}

	// 同步请求直接等待响应
	if !isAsync {
		return readSyncResponse(conn)
	}

	// 异步请求返回ID（响应通过asyncReadLoop处理）
	return nil, requestID, nil
}

// 读取同步响应
func readSyncResponse(conn net.Conn) ([]byte, string, error) {
	// 读取响应头（固定7字节）
	headerBuf := make([]byte, 7)
	_, err := conn.Read(headerBuf)
	if err != nil {
		return nil, "", fmt.Errorf("读取响应头失败: %v", err)
	}

	// 解析响应头
	version := binary.BigEndian.Uint16(headerBuf[:2])
	msgType := headerBuf[2]
	payloadLen := binary.BigEndian.Uint32(headerBuf[3:7])

	// 校验协议版本
	if version != ProtocolVersion {
		return nil, "", fmt.Errorf("不支持的响应协议版本: %d", version)
	}

	// 读取响应负载
	payload := make([]byte, payloadLen)
	_, err = conn.Read(payload)
	if err != nil {
		return nil, "", fmt.Errorf("读取响应负载失败: %v", err)
	}

	// 同步响应类型应为0x05（与设计文档一致）
	if msgType != 0x05 {
		return nil, "", fmt.Errorf("无效的同步响应类型: 0x%x", msgType)
	}

	return payload, "", nil
}
