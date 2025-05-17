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

	// 生成UUID
	uuid, err := utils.GenerateUUID()
	if err != nil {
		fmt.Printf("Error generating UUID: %v\n", err)
		return
	}

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

	// 将请求信息组装到map中
	requestData := map[string]interface{}{
		"headers":   req.Header,
		"body":      string(bodyBytes),
		"route":     req.URL.Path,
		"timestamp": requestTimestamp.Format(time.RFC3339Nano),
		"client_ip": clientIP,
		"host":      host,
		"uri":       uri,
		"url":       fullURL,
	}

	// JSON序列化请求数据
	requestDataJSON, err := json.Marshal(requestData)
	if err != nil {
		log.Println("Error marshalling request data:", err)
		http.Error(w, "内部服务器错误", http.StatusInternalServerError)
		return
	}

	// 存储到内存数据库
	r.DB.Put(uuid, string(requestDataJSON))
	// 内存数据库持久化到文件，开启后大幅降低系统性能
	r.DB.Persist()

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

	// 记录请求到达时间
	// requestStartTime = time.Now()

	// 执行Socket通信
	configFilePath, _ := utils.ResolvePath(config.GlobalConfig.BussinessSocketPath)
	output, err := ipc.ExecuteSocket(route.Command, configFilePath, uuid)
	if err != nil {
		log.Println("执行PHP命令时出错:", err)
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
