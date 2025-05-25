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
	asyncPending = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "ipc_async_pending",
		Help: "当前未完成的异步请求数量",
	})
)

type AsyncRequest struct {
	ID     string      `json:"id"`   // 全局唯一ID（UUIDv4）
	Method string      `json:"method"` // 目标方法（如JSON-RPC规范）
	Params interface{} `json:"params"` // 业务参数
}

func HandleSocket(conn net.Conn) {
	defer conn.Close()
	buf := make([]byte, HeaderSize)
	requestMap := make(map[string]chan []byte)
	var requestMapMu sync.RWMutex

	for {
		// 读取协议头（固定7字节）
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

		if header.Version < 0x0101 {
			log.Printf("不支持的协议版本: %d", header.Version)
			return
		}

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

		// 处理异步请求（MsgType=0x04）
		if header.MsgType == 0x04 {
			// 新增：解析异步请求体
			var asyncReq AsyncRequest
			if err := json.Unmarshal(payload, &asyncReq); err != nil {
				log.Printf("解析异步请求失败: %v", err)
				return
			}

			respChan := make(chan []byte, 1)
			requestMapMu.Lock()
			requestMap[asyncReq.ID] = respChan
			requestMapMu.Unlock()
			asyncPending.Inc()

			// 启动goroutine处理异步逻辑
			go func(id string, respChan chan []byte) {
				defer func() {
					requestMapMu.Lock()
					delete(requestMap, id)
					requestMapMu.Unlock()
					asyncPending.Dec()
				}()

				select {
				case respData := <-respChan:
					// 构造响应头
					versionBuf := make([]byte, 2)
					binary.BigEndian.PutUint16(versionBuf, header.Version)
					msgTypeBuf := []byte{0x05} // 响应类型
					payloadLenBuf := make([]byte, 4)
					binary.BigEndian.PutUint32(payloadLenBuf, uint32(len(respData)))

					fullResp := append(append(append(versionBuf, msgTypeBuf...), payloadLenBuf...), respData...)
					if _, err := conn.Write(fullResp); err != nil {
						log.Printf("异步响应ID=%s发送失败: %v", id, err)
					}
				case <-time.After(30 * time.Second):
					log.Printf("异步请求ID=%s超时", id)
				}
			}(asyncReq.ID, respChan)
			continue
		}

		// 恢复同步请求处理逻辑（确保导入包被使用）
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

		// 生成响应（示例逻辑，根据实际需求调整）
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

		// 发送响应（使用responseData）
		if _, err := conn.Write(append(respHeader, responseData...)); err != nil {
			log.Println("发送响应错误:", err)
		}
	}
}
