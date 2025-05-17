package plugin

import (
	"bigHammer/internal/interface/database"
	"bigHammer/internal/plugin"
	"bigHammer/internal/shared"
	"fmt"
)

// InputPlugin 输入插件
// 用于处理输入请求，从数据库获取值
type InputPlugin struct{}

// HandleRequest 实现 ServicePlugin 接口的请求处理方法
// 功能：
// 1. 从全局容器获取数据库实例
// 2. 验证请求参数
// 3. 从数据库查询值
// 4. 返回查询结果
// 参数：
//   - req plugin.Request: 要处理的请求
// 返回值：
//   - plugin.Response: 处理结果，包含查询到的值或错误信息
func (p *InputPlugin) HandleRequest(req plugin.Request) plugin.Response {
	// 从全局容器获取数据库实例
	dbInterface, err := shared.GlobalContainer.Resolve("database")
	if err != nil {
		return plugin.Response{
			Status:  500,
			Message: fmt.Sprintf("Error resolving 'database': %v", err),
		}
	}

	// 断言数据库实例到正确的类型，这取决于你的数据库实现
	db, ok := dbInterface.(database.IDatabase) // 替换YourDatabaseType为您的数据库类型
	if !ok {
		return plugin.Response{
			Status:  500,
			Message: "Error asserting database instance to the correct type",
		}
	}

	// 获取请求中的key参数
	keyStr, exists := req.Params["key"]
	if !exists {
		return plugin.Response{
			Status:  400,
			Message: "Parameter 'key' is required",
		}
	}
	// 从数据库查询值
	value, exists := db.Get(keyStr) // 假设存在Get方法
	if !exists {
		return plugin.Response{
			Status:  404,
			Message: "Key not found in database",
		}
	}

	// 返回成功响应
	return plugin.Response{
		Status: 200,
		Data:   value,
	}
}

// init 包初始化函数
// 功能：
// 1. 打印初始化信息
// 2. 注册 InputPlugin 插件到插件系统
// 参数：无
// 返回值：无
func init() {
	fmt.Println("input package initialized")
	plugin.RegisterPlugin("input", &InputPlugin{})
}
