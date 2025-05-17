package plugin

// Request 定义了插件需要处理的请求结构
// 包含服务名称、方法和参数信息
type Request struct {
	// Service 服务名称
	// 用于标识要调用的服务
	Service string            `json:"service"`
	// Method 方法名称
	// 用于标识要调用的方法
	Method  string            `json:"method"`
	// Params 请求参数
	// 存储请求的参数键值对
	Params  map[string]string `json:"params"`
}

// Response 定义了插件返回的响应结构
// 包含状态码、消息和数据
type Response struct {
	// Status HTTP状态码
	// 表示请求处理的结果状态
	Status  int         `json:"status"`
	// Message 响应消息
	// 用于描述处理结果或错误信息
	Message string      `json:"message"`
	// Data 响应数据
	// 存储实际的响应数据
	Data    interface{} `json:"data"`
}

// ServicePlugin 定义了所有插件都必须实现的接口
// 用于统一插件的处理方式
type ServicePlugin interface {
	// HandleRequest 处理请求的方法
	// 参数：
	//   - req Request: 要处理的请求
	// 返回值：
	//   - Response: 处理结果
	HandleRequest(req Request) Response
}

// Plugins 插件注册表
// 用于存储所有注册的插件
// key: 插件名称
// value: 插件实例
var Plugins = make(map[string]ServicePlugin)

// RegisterPlugin 注册插件到插件系统
// 功能：
// 1. 将插件实例存储到插件注册表中
// 参数：
//   - name string: 插件名称
//   - plugin ServicePlugin: 插件实例
// 返回值：无
func RegisterPlugin(name string, plugin ServicePlugin) {
	Plugins[name] = plugin
}

// DispatchRequest 分发请求到对应的插件
// 功能：
// 1. 根据请求中的服务名称查找对应的插件
// 2. 调用插件的处理方法
// 3. 返回处理结果
// 参数：
//   - req Request: 要处理的请求
// 返回值：
//   - Response: 处理结果
func DispatchRequest(req Request) Response {
	plugin, exists := Plugins[req.Service]
	if !exists {
		return Response{
			Status:  404,
			Message: "Service not found",
		}
	}
	return plugin.HandleRequest(req)
}
