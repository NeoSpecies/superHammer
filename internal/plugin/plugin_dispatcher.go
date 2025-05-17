package plugin

// PluginDispatcher 插件调度器
// 负责将请求转发到对应的插件进行处理
type PluginDispatcher struct{}

// HandleRequest 实现 ServicePlugin 接口的请求处理方法
// 功能：
// 1. 接收请求
// 2. 使用 DispatchRequest 函数处理请求
// 3. 返回处理结果
// 参数：
//   - req Request: 要处理的请求
// 返回值：
//   - Response: 处理结果
func (pd *PluginDispatcher) HandleRequest(req Request) Response {
	return DispatchRequest(req) // 使用 DispatchRequest 函数处理请求
}
