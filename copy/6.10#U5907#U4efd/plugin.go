package plugin

import "fmt"

// Request 定义了插件需要处理的请求结构
type Request struct {
	Service string            `json:"service"`
	Method  string            `json:"method"`
	Params  map[string]string `json:"params"`
}

// Response 定义了插件返回的响应结构
type Response struct {
	Status  int         `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

// ServicePlugin 定义了所有插件都必须实现的接口
type ServicePlugin interface {
	HandleRequest(req Request) Response
}

// 插件注册表，用于存储所有注册的插件
var Plugins = make(map[string]ServicePlugin)

// RegisterPlugin 允许插件注册自己，以便他们可以被调度
func RegisterPlugin(name string, plugin ServicePlugin) {
	Plugins[name] = plugin
}

// DispatchRequest 根据请求中的 Service 字段决定调用哪个插件，调用插件的处理函数并返回响应
func DispatchRequest(req Request) Response {
	fmt.Println(re)
	plugin, exists := Plugins[req.Service]
	if !exists {
		return Response{
			Status:  404,
			Message: "Service not found",
		}
	}
	return plugin.HandleRequest(req)
}
