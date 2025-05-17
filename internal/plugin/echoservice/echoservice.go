package echoservice

import (
	"bigHammer/internal/plugin"
	"fmt"
)

// EchoService 回显服务插件
// 用于测试和调试，将请求参数原样返回
type EchoService struct{}

// HandleRequest 实现 ServicePlugin 接口的请求处理方法
// 功能：
// 1. 接收请求
// 2. 将请求参数作为响应数据返回
// 3. 返回成功状态
// 参数：
//   - req plugin.Request: 要处理的请求
// 返回值：
//   - plugin.Response: 处理结果，包含原始请求参数
func (e EchoService) HandleRequest(req plugin.Request) plugin.Response {
	return plugin.Response{
		Status:  200,
		Message: "Echo successful hhh",
		Data:    req.Params, // 简单地回显请求参数
	}
}

// init 包初始化函数
// 功能：
// 1. 打印初始化信息
// 2. 注册 EchoService 插件到插件系统
// 参数：无
// 返回值：无
func init() {
	fmt.Println("echoservice package initialized")
	plugin.RegisterPlugin("echo", EchoService{})
}
