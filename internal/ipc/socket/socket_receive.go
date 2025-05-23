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
	"sync"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"

	"time" // 补全缺失的time包
)

const (
	HeaderSize     = 7               // 2+1+4 bytes
	MaxPayloadSize = 4 * 1024 * 1024 // 4MB
)

type ProtocolHeader struct {
	Version uint16
	MsgType byte
	Length  uint32
}

var (
	requestMap   = make(map[string]net.Conn) // 新增请求映射表
	mapMutex     sync.RWMutex                // 新增互斥锁
	asyncPending = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "ipc_async_pending",
		Help: "当前未完成的异步请求数量",
	})
)

// 新增异步请求结构体
type AsyncRequest struct {
	ID      string
	Payload json.RawMessage
}

func handleAsyncRequest(version uint16, payload []byte, conn net.Conn) {
	// 修复顺序：先解析请求体
	var asyncReq AsyncRequest
	if err := json.Unmarshal(payload, &asyncReq); err != nil {
		log.Printf("解析异步请求失败: %v", err)
		return
	}

	// 合并互斥锁操作
	mapMutex.Lock()
	defer mapMutex.Unlock()

	// 设置超时必须放在锁操作之后
	conn.SetDeadline(time.Now().Add(30 * time.Second))

	if _, exists := requestMap[asyncReq.ID]; exists {
		log.Printf("重复的异步请求ID: %s", asyncReq.ID)
		return
	}

	// 正确注册连接
	requestMap[asyncReq.ID] = conn
	asyncPending.Inc()

	defer func() {
		delete(requestMap, asyncReq.ID)
		asyncPending.Dec()
		conn.Close()
	}()
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
