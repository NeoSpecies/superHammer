package ipc

import (
	"bigHammer/internal/plugin"
	"bigHammer/internal/shared"
	"encoding/binary"
	"encoding/json"
	"errors"
	"io"
	"log"
	"net"
)

const (
	HeaderSize = 7 // 2+1+4 bytes
	MaxPayloadSize = 4 * 1024 * 1024 // 4MB
)

type ProtocolHeader struct {
	Version uint16
	MsgType byte
	Length  uint32
}

func HandleSocket(conn net.Conn) {
	defer conn.Close()
	buf := make([]byte, HeaderSize)

	// 读取协议头
	if _, err := io.ReadFull(conn, buf); err != nil {
		if !errors.Is(err, io.EOF) {
			log.Println("读取协议头错误:", err)
		}
		return
	}

	header := ProtocolHeader{
		Version: binary.BigEndian.Uint16(buf[:2]),
		MsgType: buf[2],
		Length:  binary.BigEndian.Uint32(buf[3:7]),
	}

	// 验证负载长度
	if header.Length > MaxPayloadSize {
		log.Println("负载大小超出限制:", header.Length)
		return
	}

	// 读取负载
	payload := make([]byte, header.Length)
	if _, err := io.ReadFull(conn, payload); err != nil {
		log.Println("读取负载错误:", err)
		return
	}

	// 处理请求
	pluginInterface, err := shared.GlobalContainer.Resolve("plugin")
	if err != nil {
		log.Println("解析插件错误:", err)
		return
	}

	var req plugin.Request
	if err := json.Unmarshal(payload, &req); err != nil {
		log.Println("解析JSON错误:", err)
		return
	}

	pluginInstance, ok := pluginInterface.(plugin.ServicePlugin)
	if !ok {
		log.Println("插件接口不匹配")
		return
	}

	// 生成响应
	response := pluginInstance.HandleRequest(req)
	responseData, err := json.Marshal(response)
	if err != nil {
		log.Println("序列化响应错误:", err)
		return
	}

	// 封装响应协议头
	respHeader := make([]byte, HeaderSize)
	binary.BigEndian.PutUint16(respHeader[:2], header.Version)
	respHeader[2] = 0x05 // 响应类型
	binary.BigEndian.PutUint32(respHeader[3:7], uint32(len(responseData)))

	// 发送响应
	if _, err := conn.Write(append(respHeader, responseData...)); err != nil {
		log.Println("发送响应错误:", err)
	}
}
