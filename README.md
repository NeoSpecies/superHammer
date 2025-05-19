### BigHammer 微服务网关框架：开启 AIoT 新时代的通信枢纽

在 AIoT（人工智能物联网）的浪潮下，传统单一语言架构在应对复杂场景时显得力不从心。BigHammer 微服务网关框架应运而生，它以 Golang 为底层通信核心，通过进程间通信（IPC）开放应用逻辑层给多语言，为 AIoT 领域带来了高性能、易扩展、高稳定的通信解决方案。

#### 适用场景广泛，满足多样需求
- **IoT 设备接入**：支持 TCP、UDP、MQTT 等多种协议的设备接入，轻松实现海量 IoT 设备的互联互通。无论是智能家居中的传感器、工业生产中的监控设备，还是物流运输中的追踪器，BigHammer 都能稳定高效地处理设备连接和数据传输。
- **实时数据交互**：借助 WebSocket 协议，实现设备与服务器之间的实时数据交互。在金融交易、能源管理、交通监控等领域，实时数据的及时处理和反馈至关重要，BigHammer 能够确保数据的实时性和准确性。
- **边缘 AI 推理**：框架内集成边缘 AI（最大不超过 4B 的模型），主要用于格式化文件处理等工作。在智能安防、医疗影像分析、农业监测等场景中，边缘 AI 能够在本地进行数据处理和分析，减少数据传输延迟，提高系统的响应速度和智能化水平。
- **多技术栈微服务集成**：作为 API 网关，支持多语言服务集成（PHP、Python 等），实现多技术栈微服务的无缝协作。在大型企业级应用、互联网平台等项目中，不同团队可能使用不同的技术栈进行开发，BigHammer 能够将这些微服务整合在一起，提高开发效率和系统的可维护性。
- **硬件联动控制**：通过 AIoT 自定义协议（MCP），实现硬件之间的联动控制。在智能工厂、智能建筑、智能交通等领域，硬件设备之间的协同工作需要精确的控制和协调，BigHammer 能够根据设备状态和业务规则，自动下发控制指令，实现硬件的智能化联动。

#### 架构设计精妙，保障系统性能
- **分层架构清晰**：采用分层架构设计，分为 Golang 核心通信引擎、多语言应用层和平台级功能层。Golang 核心通信引擎提供高性能、高并发的网络通信基础能力，屏蔽底层协议复杂性；多语言应用层降低业务开发门槛，允许开发者使用熟悉的语言实现逻辑；平台级功能层提供 AIoT 场景的垂直能力支持，如边缘 AI 集成、设备管理、API 网关扩展等。
- **核心模块强大**：Golang 核心通信引擎包含多协议监听器、连接管理器、协议分发器、IPC 管理器和序列化引擎等核心模块。多协议监听器支持 HTTP/HTTPS、WebSocket、TCP、UDP、MQTT 等多种协议；连接管理器能够管理百万级并发连接，支持心跳检测、超时断开、重连触发；协议分发器基于元数据（端口 / 消息头）动态路由至应用服务；IPC 管理器实现 Golang 与多语言应用服务的可靠通信；序列化引擎支持跨语言消息序列化与反序列化。
- **稳定性保障完善**：具备完善的错误处理和资源监控机制，确保系统在高并发下稳定运行。每个 goroutine 都添加了 recover 机制，能够捕获并处理异常；通过 Prometheus 暴露资源监控指标，实时监控系统的运行状态。同时，还提供了健康检查和自动重连机制，提升系统的容错能力；采用 Docker 容器化部署，实现资源隔离和快速部署。

#### 插件化架构灵活，支持快速扩展
- **插件系统丰富**：采用插件化架构设计，支持开发者通过插件化方式开发 MCP 程序（Clients/Hosts/Servers）。框架提供了插件间能力调用、插件与边缘大模型的集成调用以及硬件链接支持等功能。开发者可以根据业务需求，灵活开发和部署插件，实现系统的功能扩展和定制化。
- **多语言 SDK 支持**：提供多语言 SDK，降低开发门槛，支持不同技术栈团队协作。Python、Java 等语言的开发者可以使用相应的 SDK 快速开发应用服务，与 Golang 核心通信引擎进行交互。同时，SDK 还提供了自动重连、心跳发送等通用功能，提高开发效率和系统的稳定性。

#### 开发体验优良，助力高效开发
- **完善的文档与工具链**：提供详细的架构设计文档、插件开发指南和多语言示例代码，帮助开发者快速上手。同时，还提供了服务管理脚本（duang.sh），实现服务的自动化管理，包括启动、停止、状态检查、重启等功能。
- **灵活的配置系统**：支持灵活的配置系统，通过配置文件（config.json、plugins.json、router.json 等）可以方便地进行系统配置和参数调整。同时，还支持热重载功能，在不重启服务的情况下，实现配置的动态更新。

BigHammer 微服务网关框架是一款功能强大、性能卓越、易于扩展的 AIoT 通信框架，能够为企业和开发者提供一站式的解决方案，助力 AIoT 应用的快速开发和部署。无论是小型创业公司还是大型企业，都可以借助 BigHammer 的优势，在 AIoT 领域取得更大的成功。

## 项目特点

- 高性能Go语言网关层
- 支持多语言服务集成（PHP、Python等）
- 插件化架构设计
- 完善的进程管理和监控
- 灵活的配置系统
- 支持热重载
- 多种通信协议支持（HTTP、WebSocket、TCP、Unix Domain Socket）

## 系统架构

### 目录结构

```
/bigHammer
├── /cmd                    # 命令行工具
│   └── /gateway           # 网关服务入口
├── /config                # 配置文件目录
│   ├── config.json        # 主配置文件
│   ├── plugins.json       # 插件配置
│   ├── router.json        # 路由配置
│   └── data.json          # 内存数据库数据
├── /internal              # 内部包
│   ├── /config           # 配置管理
│   ├── /di               # 依赖注入容器
│   ├── /ipc              # 进程间通信
│   ├── /plugin           # 插件系统
│   ├── /router           # 路由管理
│   ├── /service          # 服务实现
│   │   ├── /http        # HTTP服务
│   │   └── /socket      # Socket服务
│   └── /shared           # 共享资源
├── /pkg                   # 公共包
│   └── /utils            # 工具函数
├── /runtime              # 运行时文件
│   ├── duangMain.pid     # 主进程PID
│   ├── duangBussiness.pid # 业务进程PID
│   ├── mainSocket.sock   # 主Socket文件
│   └── phpSocket.sock    # PHP服务Socket文件
├── /scripts              # 脚本文件
├── /test                 # 测试文件
├── /Develop              # 开发目录
├── go.mod                # Go模块定义
├── go.sum                # Go模块依赖
└── README.md             # 项目文档
```
## 服务管理脚本 (duang.sh)

`duang.sh` 是一个用于管理 BigHammer 服务的脚本，提供了完整的服务生命周期管理功能。

### 功能特点

1. **完整的服务管理**
   - 自动启动 PHP 和 Go 服务
   - 优雅停止服务
   - 服务状态检查
   - 自动重启功能

2. **健壮的错误处理**
   - 自动检测和清理 socket 和 pid 文件
   - 服务启动超时检测
   - 优雅的服务停止（先尝试 SIGTERM，再使用 SIGKILL）
   - 启动失败时自动清理

3. **完善的日志系统**
   - 彩色控制台输出
   - 时间戳记录
   - 日志级别（INFO/WARN/ERROR）
   - 同时输出到控制台和日志文件

4. **服务状态监控**
   - 检查进程是否存在
   - 检查 socket 文件是否存在
   - 检查服务是否正常运行

### 使用方法

```bash
# 启动所有服务
./duang.sh start

# 停止所有服务
./duang.sh stop

# 检查服务状态
./duang.sh status

# 重启所有服务
./duang.sh restart

# 显示帮助信息
./duang.sh help
```
## 系统流程图

### 1. 整体架构图

```
+-------------+     +----------------+     +-----------------+
|   客户端    | --> |   Go网关层     | --> |   业务服务层    |
+-------------+     +----------------+     +-----------------+
                           |                       |
                           v                       v
                    +----------------+     +-----------------+
                    |   后端服务     | <-- | PHP/Python等    |
                    +----------------+     +-----------------+

Go网关层内部结构：
+----------------+
|   HTTP服务     |
+----------------+
       |
       v
+----------------+
|   WebSocket    |
+----------------+
       |
       v
+----------------+
|   TCP服务      |
+----------------+
       |
       v
+----------------+
|   路由系统     |
+----------------+
       |
       v
+----------------+
|   插件系统     |
+----------------+
       |
       v
+----------------+
|  Socket服务    |
+----------------+
```

### 2. 请求处理流程

```
客户端 --> Go网关 --> 插件系统 --> Socket服务 --> 业务服务
   ^          |          |            |            |
   |          v          v            v            v
   +---- 响应处理 <-- 响应插件 <-- 返回结果 <-- 处理结果
```

### 3. 插件系统流程

```
请求 --> 初始化插件 --> 预处理插件 --> 处理插件 --> 后处理插件 --> 响应
```

### 4. 热重载机制

```
文件变化 --> 文件监控器
   |
   +--> 配置文件变化 --> 重载配置 --> 更新配置
   |
   +--> PHP文件变化 --> 重启服务 --> 停止服务 --> 启动服务
   |
   +--> 插件文件变化 --> 重载插件 --> 卸载插件 --> 加载插件
```

### 5. 进程管理流程

```
启动阶段：
主进程 --> 写入PID --> 启动Socket服务 --> 启动业务进程

运行阶段：
主进程 <--> 业务进程 (心跳检测)

关闭阶段：
主进程 --> 发送关闭信号 --> 业务进程确认关闭
   |
   +--> 关闭Socket服务
   |
   +--> 删除PID文件
```

### 6. 依赖注入流程

```
注册服务 --> 容器
   |
   +--> 单例服务
   |
   +--> 工厂服务

解析服务 --> 容器
   |
   +--> 缓存实例
   |
   +--> 新建实例
```

### 7. 配置管理流程

```
加载配置 --> 读取文件 --> 解析JSON --> 验证配置 --> 更新配置

配置热重载：
监控文件 --> 文件变化 --> 重新加载 --> 加载配置
```

### 8. Socket通信流程

```
客户端 --> 建立连接 --> Socket服务器
   |                        |
   |                        v
   |                    接受连接
   |                        |
   |                        v
   |                    读取请求数据
   |                        |
   |                        v
   |                    处理请求
   |                        |
   |                        v
   |                    发送响应
   |                        |
   +------------------- 返回数据
   |
   +------------------- 关闭连接
```

## 请求处理流程

1. **请求进入**
   - 客户端请求通过HTTP/WebSocket/TCP进入网关
   - 请求经过路由系统进行分发
   - 根据配置决定是否需要插件处理

2. **插件处理**
   - 请求经过已注册的插件链
   - 插件可以修改请求、响应或中断处理
   - 支持认证、日志、限流等插件

3. **业务处理**
   - 请求通过Socket转发到对应的业务服务
   - 业务服务处理请求并返回结果
   - 结果通过Socket返回给网关

4. **响应处理**
   - 网关接收业务服务响应
   - 响应经过插件链处理
   - 最终响应返回给客户端

## 生命周期

1. **启动阶段**
   - 加载配置文件
   - 初始化依赖注入容器
   - 注册插件系统
   - 启动Socket服务器
   - 启动HTTP服务器
   - 启动业务服务

2. **运行阶段**
   - 处理客户端请求
   - 管理插件生命周期
   - 监控系统状态
   - 处理热重载

3. **关闭阶段**
   - 优雅关闭HTTP服务器
   - 关闭Socket连接
   - 停止业务服务
   - 清理资源

## 插件开发指南

### 插件接口

```go
type Plugin interface {
    // 插件初始化
    Init() error
    
    // 请求处理
    HandleRequest(ctx *Context) error
    
    // 响应处理
    HandleResponse(ctx *Context) error
    
    // 插件关闭
    Close() error
}
```

### 插件开发步骤

1. 创建插件目录和文件
2. 实现Plugin接口
3. 在plugins.json中注册插件
4. 配置插件参数

### 示例插件

```go
package example

type ExamplePlugin struct {
    // 插件配置
    config map[string]interface{}
}

func (p *ExamplePlugin) Init() error {
    // 初始化逻辑
    return nil
}

func (p *ExamplePlugin) HandleRequest(ctx *Context) error {
    // 请求处理逻辑
    return nil
}

func (p *ExamplePlugin) HandleResponse(ctx *Context) error {
    // 响应处理逻辑
    return nil
}

func (p *ExamplePlugin) Close() error {
    // 清理逻辑
    return nil
}
```

## 项目启动

### 环境要求

- Go 1.16+
- 支持Unix Domain Socket的操作系统
- 足够的文件描述符限制

### 启动步骤

1. 安装依赖
```bash
go mod download
```

2. 编译项目
```bash
go build -o bigHammer
```

3. 启动服务
```bash
./bigHammer
```

### 配置文件

主要配置文件位于`config/config.json`：

```json
{
    "database_url": "postgres://user:password@localhost/dbname",
    "memory_db_path": "/config/data.json",
    "mainPid_path": "/runtime/duangMain.pid",
    "bussinessPid_path": "/runtime/duangBussiness.pid",
    "socket_path": "/runtime/mainSocket.sock",
    "bussiness_socket_path": "/runtime/phpSocket.sock",
    "router_path": "/config/router.json",
    "plugins_path": "/config/plugins.json",
    "bussiness_main_path": "/Develop/test.php",
    "ports": {
        "ipc_port": "8000",
        "http_port": "8080",
        "websocket_port": "8081",
        "tcp_port": "8082"
    }
}
```

## 热重载

系统支持以下热重载功能：

1. **配置热重载**
   - 修改配置文件后自动重新加载
   - 支持动态更新路由配置
   - 支持动态更新插件配置

2. **插件热重载**
   - 支持动态加载新插件
   - 支持更新现有插件
   - 支持禁用/启用插件

3. **业务服务热重载**
   - 监控业务服务文件变化
   - 自动重启业务服务
   - 保持连接状态

## 开发注意事项

1. 使用`go run .`启动项目，不要使用`go run main.go`
2. 确保runtime目录具有正确的权限
3. 注意文件描述符限制
4. 遵循插件开发规范
5. 保持配置文件格式正确

## 常见问题

1. **Socket连接失败**
   - 检查Socket文件权限
   - 确认进程是否正在运行
   - 检查文件描述符限制

2. **插件加载失败**
   - 检查插件配置格式
   - 确认插件接口实现
   - 查看错误日志

3. **热重载不生效**
   - 检查文件监控配置
   - 确认文件权限
   - 查看进程日志

## 贡献指南

1. Fork项目
2. 创建特性分支
3. 提交更改
4. 推送到分支
5. 创建Pull Request

## 许可证

[许可证类型]

## 核心机制实现详解

### 1. Socket通信机制

#### 1.1 Socket服务器实现
```go
// internal/service/socket/socket_server.go

// StartSocketServer 启动Socket服务器
// 参数:
//   - ctx: 上下文对象，用于控制服务器的生命周期
// 功能:
//   1. 加载配置文件
//   2. 创建Unix Domain Socket监听器
//   3. 启动优雅关闭处理
//   4. 进入主循环处理连接
func StartSocketServer(ctx context.Context) {
    // 加载配置文件
    err := config.LoadConfig()
    if err != nil {
        fmt.Println("Error loading config:", err)
        return
    }
    
    // 创建Unix Domain Socket监听器
    // 使用Unix Domain Socket而不是TCP Socket，提供更好的性能和安全性
    configFilePath, _ := utils.ResolvePath(config.GlobalConfig.SocketPath)
    listener, err := net.Listen("unix", configFilePath)
    if err != nil {
        log.Fatal("Error starting listener:", err)
    }

    // 优雅关闭处理
    // 当收到取消信号时，关闭监听器
    go func() {
        <-ctx.Done()
        listener.Close()
    }()

    // 主循环：接受并处理连接
    // 每个连接在独立的goroutine中处理，实现并发
    for {
        conn, err := listener.Accept()
        if err != nil {
            log.Println("Listener accept error:", err)
            break
        }
        go ipc.HandleSocket(conn)
    }
}

// StopSocketServer 停止Socket服务器
// 功能:
//   1. 关闭监听器
//   2. 删除Socket文件
//   3. 清理资源
func StopSocketServer() {
    // 加载配置
    err := config.LoadConfig()
    if err != nil {
        fmt.Println("Error loading config:", err)
        return
    }
    
    // 获取Socket文件路径
    configFilePath, _ := utils.ResolvePath(config.GlobalConfig.SocketPath)

    // 尝试关闭监听器
    _, err = net.Dial("unix", configFilePath)
    if err != nil {
        log.Println("Socket server is not running")
        return
    }

    // 关闭监听器
    listener, err := net.FileListener(os.NewFile(3, configFilePath))
    if err != nil {
        log.Println("Error getting listener:", err)
        return
    }
    
    // 关闭监听器并删除Socket文件
    err = listener.Close()
    if err != nil {
        log.Println("Error closing listener:", err)
    }
    
    err = os.Remove(configFilePath)
    if err != nil {
        log.Println("Error removing socket file:", err)
    }
}
```

#### 1.2 Socket连接处理
```go
// internal/ipc/socket/socket_receive.go

// Response 定义响应结构
// 字段:
//   - Status: HTTP状态码
//   - Message: 响应消息
//   - Data: 响应数据
type Response struct {
    Status  int         `json:"status"`
    Message string      `json:"message"`
    Data    interface{} `json:"data"`
}

// HandleSocket 处理Socket连接
// 参数:
//   - conn: 网络连接对象
// 功能:
//   1. 读取请求数据
//   2. 解析请求
//   3. 获取插件实例
//   4. 处理请求
//   5. 发送响应
func HandleSocket(conn net.Conn) {
    // 确保连接最终被关闭
    defer conn.Close()

    // 创建缓冲区并读取数据
    data := make([]byte, 1024)
    n, err := conn.Read(data)
    if err != nil {
        if errors.Is(err, io.EOF) {
            // 连接被客户端关闭
            return
        }
        log.Println("Error reading data from socket:", err)
        return
    }

    // 从容器中获取插件实例
    pluginInterface, err := shared.GlobalContainer.Resolve("plugin")
    if err != nil {
        log.Println("Error resolving plugin from container:", err)
        return
    }

    // 解析请求数据
    var req plugin.Request
    if err := json.Unmarshal(data[:n], &req); err != nil {
        log.Println("Error parsing request:", err)
        return
    }

    // 类型断言，确保插件实现了正确的接口
    pluginInstance, ok := pluginInterface.(plugin.ServicePlugin)
    if !ok {
        log.Println("The provided plugin service does not implement the ServicePlugin interface.")
        return
    }

    // 处理请求并发送响应
    response := pluginInstance.HandleRequest(req)
    responseJSON, err := json.Marshal(response)
    if err != nil {
        log.Println("Error marshalling response:", err)
        return
    }
    
    // 发送响应数据
    _, err = conn.Write(responseJSON)
    if err != nil {
        log.Println("Error writing data to socket:", err)
        return
    }
}
```

### 2. 插件系统实现

#### 2.1 插件接口定义
```go
// internal/plugin/plugin.go

// Plugin 定义插件接口
// 所有插件必须实现这个接口
type Plugin interface {
    // Init 初始化插件
    // 返回值:
    //   - error: 初始化过程中的错误
    Init() error
    
    // HandleRequest 处理请求
    // 参数:
    //   - ctx: 请求上下文，包含请求和配置信息
    // 返回值:
    //   - error: 处理过程中的错误
    HandleRequest(ctx *Context) error
    
    // HandleResponse 处理响应
    // 参数:
    //   - ctx: 响应上下文，包含响应和配置信息
    // 返回值:
    //   - error: 处理过程中的错误
    HandleResponse(ctx *Context) error
    
    // Close 关闭插件
    // 返回值:
    //   - error: 关闭过程中的错误
    Close() error
}

// Context 插件上下文
// 包含请求处理所需的所有信息
type Context struct {
    Request  *Request           // 请求对象
    Response *Response          // 响应对象
    Config   map[string]interface{} // 插件配置
}
```

#### 2.2 插件注册机制
```go
// internal/plugin/dispatcher.go

// PluginDispatcher 插件调度器
// 负责管理和执行插件链
type PluginDispatcher struct {
    plugins []Plugin                    // 插件列表
    config  map[string]interface{}      // 配置信息
}

// RegisterPlugin 注册插件
// 参数:
//   - plugin: 要注册的插件实例
// 返回值:
//   - error: 注册过程中的错误
// 功能:
//   1. 初始化插件
//   2. 将插件添加到插件链
func (d *PluginDispatcher) RegisterPlugin(plugin Plugin) error {
    // 初始化插件
    if err := plugin.Init(); err != nil {
        return err
    }
    
    // 添加到插件链
    d.plugins = append(d.plugins, plugin)
    return nil
}

// HandleRequest 处理请求
// 参数:
//   - ctx: 请求上下文
// 返回值:
//   - error: 处理过程中的错误
// 功能:
//   按顺序执行插件链中的每个插件
func (d *PluginDispatcher) HandleRequest(ctx *Context) error {
    // 按顺序执行插件链
    for _, plugin := range d.plugins {
        if err := plugin.HandleRequest(ctx); err != nil {
            return err
        }
    }
    return nil
}
```

### 3. 依赖注入容器

#### 3.1 容器实现
```go
// internal/di/container.go

// Container 依赖注入容器
// 负责管理服务的生命周期和依赖关系
type Container struct {
    services map[string]interface{}     // 已注册的服务实例
    factory  map[string]func() interface{} // 服务工厂函数
}

// NewContainer 创建新的容器实例
// 返回值:
//   - *Container: 新创建的容器实例
func NewContainer() *Container {
    return &Container{
        services: make(map[string]interface{}),
        factory:  make(map[string]func() interface{}),
    }
}

// Register 注册服务
// 参数:
//   - name: 服务名称
//   - factory: 服务工厂函数
//   - scope: 服务作用域（单例/工厂）
// 返回值:
//   - error: 注册过程中的错误
func (c *Container) Register(name string, factory func() interface{}, scope Scope) error {
    c.factory[name] = factory
    return nil
}

// Resolve 解析服务
// 参数:
//   - name: 服务名称
// 返回值:
//   - interface{}: 服务实例
//   - error: 解析过程中的错误
// 功能:
//   1. 检查缓存中是否存在服务实例
//   2. 如果不存在，使用工厂函数创建新实例
//   3. 将新实例缓存并返回
func (c *Container) Resolve(name string) (interface{}, error) {
    // 检查缓存
    if service, exists := c.services[name]; exists {
        return service, nil
    }
    
    // 使用工厂函数创建实例
    if factory, exists := c.factory[name]; exists {
        service := factory()
        c.services[name] = service
        return service, nil
    }
    
    return nil, fmt.Errorf("service %s not found", name)
}
```

### 4. 配置管理

#### 4.1 配置加载
```go
// internal/config/config.go

// Config 配置结构
// 包含系统运行所需的所有配置项
type Config struct {
    DatabaseURL        string            `json:"database_url"`         // 数据库连接URL
    MemoryDBPath      string            `json:"memory_db_path"`       // 内存数据库路径
    MainPIDPath       string            `json:"mainPid_path"`         // 主进程PID文件路径
    BusinessPIDPath   string            `json:"bussinessPid_path"`    // 业务进程PID文件路径
    SocketPath        string            `json:"socket_path"`          // Socket文件路径
    BusinessSocketPath string           `json:"bussiness_socket_path"` // 业务Socket文件路径
    RouterPath        string            `json:"router_path"`          // 路由配置路径
    PluginsPath       string            `json:"plugins_path"`         // 插件配置路径
    BusinessMainPath  string            `json:"bussiness_main_path"`  // 业务主程序路径
    Ports             map[string]string `json:"ports"`                // 端口配置
}

// LoadConfig 加载配置
// 返回值:
//   - error: 加载过程中的错误
// 功能:
//   1. 读取配置文件
//   2. 解析JSON数据
//   3. 更新全局配置
func LoadConfig() error {
    // 读取配置文件
    data, err := os.ReadFile("config/config.json")
    if err != nil {
        return err
    }
    
    // 解析配置
    if err := json.Unmarshal(data, &GlobalConfig); err != nil {
        return err
    }
    
    return nil
}
```

### 5. 热重载机制

#### 5.1 文件监控
```go
// pkg/utils/watcher.go
type FileWatcher struct {
    paths []string
    pid   string
    main  string
}

func (w *FileWatcher) Start(ctx context.Context) {
    // 1. 创建文件监控器
    watcher, err := fsnotify.NewWatcher()
    if err != nil {
        log.Fatal(err)
    }
    defer watcher.Close()

    // 2. 添加监控路径
    for _, path := range w.paths {
        if err := watcher.Add(path); err != nil {
            log.Fatal(err)
        }
    }

    // 3. 监控文件变化
    for {
        select {
        case event := <-watcher.Events:
            if event.Op&fsnotify.Write == fsnotify.Write {
                w.handleFileChange(event.Name)
            }
        case err := <-watcher.Errors:
            log.Println("Error:", err)
        case <-ctx.Done():
            return
        }
    }
}

func (w *FileWatcher) handleFileChange(path string) {
    // 1. 检查文件类型
    if strings.HasSuffix(path, ".php") {
        // 2. 重启PHP服务
        w.restartPHPService()
    } else if strings.HasSuffix(path, ".json") {
        // 3. 重新加载配置
        w.reloadConfig()
    }
}
```

### 6. 进程管理

#### 6.1 PID文件管理
```go
// main.go
func WritePidToFile() {
    mainPid := os.Getpid()
    config_file_path, err := utils.ResolvePath(globalConfig.MainPIDPath)
    if err != nil {
        return
    }
    
    file, err := os.Create(config_file_path)
    if err != nil {
        return
    }
    defer file.Close()
    
    file.WriteString(fmt.Sprintf("%d", mainPid))
}

func TerminateProcessByPIDFile(pidFilePath string) error {
    // 1. 读取PID
    pidBytes, err := os.ReadFile(pidFilePath)
    if err != nil {
        return err
    }
    
    // 2. 终止进程
    pid := string(pidBytes)
    output, errOutput, err := cmd.ExecCommand("kill", pid)
    if err != nil {
        return err
    }
    
    return nil
}
```

## 开发指南

### 1. 添加新插件

1. 创建插件文件：
```go
// internal/plugin/example/example_plugin.go
package example

type ExamplePlugin struct {
    config map[string]interface{}
}

func (p *ExamplePlugin) Init() error {
    // 初始化逻辑
    return nil
}

func (p *ExamplePlugin) HandleRequest(ctx *plugin.Context) error {
    // 请求处理逻辑
    return nil
}

func (p *ExamplePlugin) HandleResponse(ctx *plugin.Context) error {
    // 响应处理逻辑
    return nil
}

func (p *ExamplePlugin) Close() error {
    // 清理逻辑
    return nil
}
```

2. 注册插件：
```json
// config/plugins.json
{
    "plugins": [
        {
            "name": "example",
            "enabled": true,
            "config": {
                "key": "value"
            }
        }
    ]
}
```

### 2. 添加新路由

1. 更新路由配置：
```json
// config/router.json
{
    "routes": [
        {
            "path": "/api/example",
            "method": "GET",
            "handler": "example_handler",
            "plugins": ["auth", "logger"]
        }
    ]
}
```

### 3. 开发新服务

1. 创建服务文件：
```php
// Develop/example_service.php
<?php
namespace Develop;

class ExampleService {
    public function handle($request) {
        // 处理请求
        return [
            'status' => 200,
            'data' => [
                'message' => 'Hello World'
            ]
        ];
    }
}
```

2. 注册服务：
```json
// config/services.json
{
    "services": [
        {
            "name": "example",
            "path": "/Develop/example_service.php",
            "enabled": true
        }
    ]
}
```