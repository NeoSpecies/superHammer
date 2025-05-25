package router

import (
	"bigHammer/internal/config"
	ipc "bigHammer/internal/ipc/socket"
	"bigHammer/pkg/utils"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
)

func (r *Router) HandleHTTP(w http.ResponseWriter, req *http.Request) {
	log.Printf("开始处理请求: %s %s", req.Method, req.URL.Path)
	requestStartTime := time.Now()

	// 加载配置
	err := config.LoadConfig()
	if err != nil {
		fmt.Println("Error loading config:", err)
		return
	}
	log.Println("进入请求")
	// 读取请求体
	bodyBytes, err := io.ReadAll(req.Body)
	if err != nil {
		log.Println("Error reading request body:", err)
		http.Error(w, "内部服务器错误", http.StatusInternalServerError)
		return
	}
	defer req.Body.Close()

	// 获取请求时间戳
	requestTimestamp := time.Now()

	// 获取客户端的 IP 地址
	clientIP := req.RemoteAddr

	// 获取Host和URI
	host := req.Host
	uri := req.RequestURI

	// 获取完整的URL
	scheme := "http"
	if req.TLS != nil {
		scheme += "s"
	}
	fullURL := fmt.Sprintf("%s://%s%s", scheme, host, req.RequestURI)

	// 尝试解析请求体为JSON
	var bodyData interface{}
	if len(bodyBytes) > 0 {
		// 尝试解析为JSON
		if err := json.Unmarshal(bodyBytes, &bodyData); err != nil {
			// 解析失败，作为普通字符串处理
			bodyData = string(bodyBytes)
		}
	} else {
		bodyData = nil
	}

	// 将请求信息组装到map中
	requestData := map[string]interface{}{
		"headers":   req.Header,
		"body":      bodyData, // 这里存储解析后的JSON对象或原始字符串
		"route":     req.URL.Path,
		"timestamp": requestTimestamp.Format(time.RFC3339Nano),
		"client_ip": clientIP,
		"host":      host,
		"uri":       uri,
		"url":       fullURL,
	}

	// 获取请求路径
	path := req.URL.Path

	// 查找路由表中的对应配置
	var route Route
	for _, r := range r.Routes {
		if r.Path == path {
			route = r
			break
		}
	}

	// 如果未找到路由，则返回404
	if route.Path == "" {
		http.NotFound(w, req)
		return
	}

	// 执行Socket通信
	configFilePath, _ := utils.ResolvePath(config.GlobalConfig.BussinessSocketPath)
	output, _, err := ipc.TransmitIPC(false, route.Command, requestData, configFilePath)
	if err != nil {
		log.Println("执行Socket通信失败:", err)
		http.Error(w, "内部服务器错误", http.StatusInternalServerError)
		return
	}

	// 计算请求处理时间
	requestProcessTime := time.Since(requestStartTime)

	// 返回输出结果给客户端
	w.Header().Set("Content-Type", "text/plain")
	w.Write(output)

	// 打印请求处理时间
	log.Printf("请求处理时间: %s", requestProcessTime)

	log.Printf("结束处理请求: %s %s", req.Method, req.URL.Path)
}
