# IPC通讯架构设计文档
## 一、设计目标
1. 协议标准化 ：定义跨语言消息协议，解决粘包/拆包问题，支持协议版本管理
2. 连接池优化 ：PHP端实现高效连接复用，Go端增强连接生命周期管理
3. 稳定性增强 ：添加心跳机制、自动重连、错误重试、goroutine泄漏防护
4. 可观测性 ：集成日志与监控指标（连接数/吞吐量/错误率）
5. 扩展性 ：预留协议扩展字段，支持未来添加新功能（如加密、压缩）
6. 异步通讯支持 ：实现读写分离的异步通信模式，通过非阻塞IO和事件驱动解耦请求与响应，提升并发性能
7. 消息可追踪性 ：参考JSON-RPC规范，在消息中增加唯一ID字段，实现请求-响应的精准匹配与链路追踪
## 二、核心协议设计（PHP/Go同步）
### 2.1 消息格式规范（二进制）
```
[2字节版本号][1字节消息类型][4字节负载长度][N字节负载]
```
- 版本号 ：大端序2字节（如0x0100表示v1.0，0x0101表示v1.1），支持协议升级兼容
- 消息类型 ：0x01=同步业务消息，0x02=心跳包，0x03=错误通知，0x04=异步请求，0x05=异步响应
- 负载长度 ：大端序4字节整数（最大支持4GB负载）
- 负载 ：使用JSON/Protobuf等序列化后的数据（v1.1及以上版本需包含 id 字段）
### 2.2 PHP端协议实现（/www/wwwroot/Develop/Reader/UnixSocketReader.php）
```
public function sendAndReceive($message) {
    try {
        $socket = $this->getConnection();
        
        // 封装同步请求协议头（v1.0）
        $version = pack('n', 0x0100);
        $msgType = pack('C', 0x01);
        $payload = json_encode($message);
        $payloadLen = pack('N', strlen($payload));
        
        $fullMessage = $version . $msgType . $payloadLen . 
        $payload;
        
        if (socket_write($socket, $fullMessage, strlen
        ($fullMessage)) === false) {
            throw new \RuntimeException("同步消息发送失败: " . 
            socket_strerror(socket_last_error($socket)));
        }

        // 读取同步响应
        $header = socket_read($socket, 7);
        if (strlen($header) !== 7) {
            throw new \RuntimeException("无效响应头");
        }
        list($version, $msgType, $payloadLen) = unpack('nversion/
        CmsgType/NpayloadLen', $header);
        
        $response = socket_read($socket, $payloadLen);
        if ($response === false) {
            throw new \RuntimeException("同步响应读取失败: " . 
            socket_strerror(socket_last_error($socket)));
        }

        $this->release($socket);
        return json_decode($response, true);
    } catch (\Exception $e) {
        static $retry = 0;
        if ($retry++ < 3) {
            return $this->sendAndReceive($message);
        }
        error_log($e->getMessage());
        throw $e;
    }
}

/**
 * 异步消息发送（支持非阻塞）
 * @param mixed $message 业务数据（需包含method/params字段）
 * @param callable $callback 响应回调函数（参数：响应数据）
 * @return string 请求ID（UUIDv4）
 */
public function sendAsync($message, callable $callback) {
    $requestId = bin2hex(random_bytes(16)); // 生成UUIDv4
    $message['id'] = $requestId;
    
    try {
        $socket = $this->getConnection();
        $this->asyncCallbacks[$requestId] = [
            'socket' => $socket,
            'callback' => $callback,
            'expire' => time() + 30 // 30秒超时
        ];

        // 封装异步请求协议头（v1.1）
        $version = pack('n', 0x0101);
        $msgType = pack('C', 0x04);
        $payload = json_encode($message);
        $payloadLen = pack('N', strlen($payload));
        
        $fullMessage = $version . $msgType . $payloadLen . 
        $payload;
        if (socket_write($socket, $fullMessage, strlen
        ($fullMessage)) === false) {
            throw new \RuntimeException("异步消息发送失败: " . 
            socket_strerror(socket_last_error($socket)));
        }

        // 启动非阻塞读循环
        stream_set_blocking($socket, false);
        $this->asyncReadLoop($socket);
        
        return $requestId;
    } catch (\Exception $e) {
        error_log("异步发送异常: " . $e->getMessage());
        throw $e;
    }
}

protected function asyncReadLoop($socket) {
    $read = [$socket];
    $write = null;
    $except = null;
    
    while (true) {
        $numChanged = socket_select($read, $write, $except, 0, 
        100000); // 100ms轮询
        if ($numChanged === false) break;
        
        if ($numChanged > 0) {
            $header = socket_read($socket, 7);
            if (strlen($header) === 7) {
                list($version, $msgType, $payloadLen) = unpack
                ('nversion/CmsgType/NpayloadLen', $header);
                if ($msgType == 0x05) { // 异步响应类型
                    $response = socket_read($socket, $payloadLen);
                    $responseData = json_decode($response, true);
                    $callback = $this->asyncCallbacks[$responseData
                    ['id']] ?? null;
                    if ($callback) {
                        ($callback['callback'])($responseData); // 
                        触发回调
                        $this->release($callback['socket']); // 释
                        放连接
                        unset($this->asyncCallbacks[$responseData
                        ['id']]);
                    }
                }
            }
        }
    }
}
// ... existing code ...
```
### 2.3 Go端协议实现（/www/wwwroot/internal/ipc/socket/socket_receive.go）
```
func HandleSocket(conn net.Conn) {
    defer conn.Close()
    buf := make([]byte, 4096)
    requestMap := make(map[string]chan []byte) // ID到响应通道的映射

    for {
        // 读取协议头（7字节）
        _, err := io.ReadFull(conn, buf[:7])
        if err != nil {
            log.Println("Read header failed:", err)
            return
        }

        version := binary.BigEndian.Uint16(buf[:2])
        msgType := buf[2]
        payloadLen := binary.BigEndian.Uint32(buf[3:7])

        // 读取负载
        payload := make([]byte, payloadLen)
        _, err = io.ReadFull(conn, payload)
        if err != nil {
            log.Println("Read payload failed:", err)
            return
        }

        // 处理不同消息类型
        switch msgType {
        case 0x01: // 同步业务消息
            handleSyncMessage(version, payload, conn)
        case 0x02: // 心跳包
            sendHeartbeatAck(conn)
        case 0x04: // 异步请求
            handleAsyncRequest(version, payload, conn, requestMap)
        }
    }
}

func handleAsyncRequest(version uint16, payload []byte, conn net.
Conn, requestMap map[string]chan []byte) {
    var asyncReq struct {
        ID     string      `json:"id"`
        Method string      `json:"method"`
        Params interface{} `json:"params"`
    }
    if err := json.Unmarshal(payload, &asyncReq); err != nil {
        log.Println("异步请求解析失败:", err)
        return
    }

    // 创建响应通道并启动异步处理
    respChan := make(chan []byte, 1)
    requestMap[asyncReq.ID] = respChan
    go func() {
        result := processAsyncMethod(asyncReq.Method, asyncReq.
        Params) // 实际业务处理
        resp := struct {
            ID     string      `json:"id"`
            Result interface{} `json:"result"`
        }{
            ID:     asyncReq.ID,
            Result: result,
        }
        respData, _ := json.Marshal(resp)
        respChan <- respData
    }()

    // 独立goroutine等待响应并回传
    go func(id string) {
        select {
        case respData := <-respChan:
            // 封装异步响应协议头（v1.1）
            versionBuf := make([]byte, 2)
            binary.BigEndian.PutUint16(versionBuf, 0x0101)
            msgTypeBuf := []byte{0x05}
            payloadLenBuf := make([]byte, 4)
            binary.BigEndian.PutUint32(payloadLenBuf, uint32(len
            (respData)))
            
            fullResp := append(append(append(versionBuf, 
            msgTypeBuf...), payloadLenBuf...), respData...)
            conn.Write(fullResp)
        case <-time.After(30 * time.Second):
            log.Printf("异步请求ID=%s超时", id)
        }
        delete(requestMap, id)
    }(asyncReq.ID)
}
// ... existing code ...
```
### 2.4 异步消息扩展规范（兼容v1.0）
为支持异步通讯和消息追踪，v1.1版本协议在负载中增加 id 字段（JSON格式示例）：

```
{
  "id": "550e8400-e29b-41d4-a716-446655440000", // UUIDv4全局唯一ID
  "method": "device.status.update",            // 目标服务方法（类似
  JSON-RPC）
  "params": {"deviceId": "dev_123", "status": "online"}, // 业务参数
  "timestamp": 1717324800                       // 消息时间戳（可选）
}
```
- 版本兼容 ：v1.0消息不包含 id 字段时，默认视为同步请求；v1.1消息必须包含 id 字段用于异步追踪
- 消息类型扩展 ：0x04=异步请求（需携带 id ），0x05=异步响应（需携带相同 id ）
## 三、连接池优化（PHP端）
### 3.1 核心改进点
- 新增 idlePool 空闲连接池，优先复用健康连接
- 添加连接健康检查（心跳检测）
- 限制最大连接数（默认20）防止资源耗尽
- 支持异步场景下的连接复用（异步请求与响应共享连接）
### 3.2 关键代码实现（/www/wwwroot/Develop/Reader/UnixSocketReader.php）
```
protected $idlePool = []; // 空闲连接池
protected $maxPoolSize = 20; // 最大连接数
protected $asyncCallbacks = []; // 异步回调映射（id => [socket, 
callback, expire]）

public function getConnection() {
    // 优先从空闲池获取
    if (!empty($this->idlePool)) {
        $socket = array_pop($this->idlePool);
        if ($this->checkConnectionHealth($socket)) {
            return $socket;
        }
        socket_close($socket);
    }

    // 从活跃池获取或新建
    foreach ($this->pool as $key => $socket) {
        if (is_resource($socket) && $this->checkConnectionHealth
        ($socket)) {
            unset($this->pool[$key]);
            return $socket;
        }
    }

    if (count($this->pool) < $this->maxPoolSize) {
        $newSocket = $this->createConnection();
        $this->pool[] = $newSocket;
        return $newSocket;
    }

    throw new \RuntimeException("连接池已满");
}

public function release($socket) {
    if (is_resource($socket) && $this->checkConnectionHealth
    ($socket)) {
        $this->idlePool[] = $socket; // 健康连接放回空闲池
    } else {
        socket_close($socket); // 异常连接直接关闭
    }
}

protected function checkConnectionHealth($socket) {
    // 发送心跳包检测连接状态（v1.0心跳）
    $heartbeat = pack('n', 0x0100) . pack('C', 0x02) . pack('N', 
    0);
    return socket_write($socket, $heartbeat, strlen($heartbeat)) 
    !== false;
}
// ... existing code ...
```
## 四、Go端稳定性增强
### 4.1 关键改进点
- 使用 context 管理goroutine生命周期
- 添加心跳响应逻辑
- 捕获goroutine panic防止进程崩溃
- 支持异步请求的超时控制与资源回收
### 4.2 关键代码实现（/www/wwwroot/internal/service/socket/socket_server.go）
```
func StartSocketServer(ctx context.Context) {
    err := config.LoadConfig()
    if err != nil {
        log.Fatal("加载配置失败:", err)
    }

    socketPath, _ := utils.ResolvePath(config.GlobalConfig.
    SocketPath)
    listener, err := net.Listen("unix", socketPath)
    if err != nil {
        log.Fatal("监听失败:", err)
    }
    defer listener.Close()

    go func() {
        <-ctx.Done()
        listener.Close() // 上下文取消时关闭监听器
    }()

    for {
        conn, err := listener.Accept()
        if err != nil {
            if ctx.Err() != nil {
                return // 上下文已取消
            }
            log.Println("接收连接错误:", err)
            continue
        }

        // 使用带超时的context管理单个连接（30秒超时）
        connCtx, cancel := context.WithTimeout(ctx, 30*time.Second)
        go func() {
            defer cancel()
            defer conn.Close()
            defer func() { // 捕获panic
                if r := recover(); r != nil {
                    log.Println("HandleSocket panic:", r)
                }
            }()
            HandleSocket(connCtx, conn) // 处理连接
        }()
    }
}
// ... existing code ...
```
## 五、配套措施
1. 日志与监控 ：
   
   - Go端：在 HandleSocket 中记录 id 与goroutine的关联关系，监控指标增加 ipc_async_pending （未完成异步请求数）、 ipc_async_timeout （异步超时数）
   - PHP端：在 UnixSocketReader 的 sendAsync 方法中记录 async_request_id 到日志，使用Monolog记录异步链路耗时
2. 测试验证 ：
   
   - 单元测试：验证消息ID的生成与匹配（Go使用 testing 的 t.Run 并发测试，PHP使用 PHPUnit 的异步测试扩展）
   - 集成测试：模拟1000+并发异步请求，验证连接池在异步场景下的复用效率（观察 idlePool 的回收延迟）
   - 压力测试：验证异步请求超时机制（手动构造30秒未响应的请求，检查是否触发超时日志）
## 六、下一步计划
1. 完成协议v1.1扩展开发（本周）
2. 实现PHP异步通讯模块（下周）
3. 集成Go端异步监控指标（下下周）
4. 编写异步场景测试用例（改造完成后）


### 一、协议版本升级（v1.0→v1.1）
目标 ：实现协议版本兼容，支持异步消息追踪（id字段）与类型扩展（0x04/0x05）。 涉及文件 ：

- PHP端： /www/wwwroot/Develop/Reader/UnixSocketReader.php （ sendAsync 、 asyncReadLoop 方法）
- Go端： /www/wwwroot/internal/ipc/socket/socket_receive.go （ HandleSocket 、 handleAsyncRequest 函数） 步骤1：PHP端协议调整
1. 验证 sendAsync 方法是否正确封装v1.1协议头：
   
   - 确认 $version = pack('n', 0x0101); （v1.1版本号）
   - 确认 $msgType = pack('C', 0x04); （异步请求类型）
   - 检查 $message['id'] 是否为UUIDv4（当前代码使用 bin2hex(random_bytes(16)) 生成，需确保唯一性）。
2. 验证 asyncReadLoop 方法是否正确解析异步响应（0x05类型）：
   
   - 检查 $msgType == 0x05 的判断逻辑
   - 确认从响应中提取 id 并匹配 $this->asyncCallbacks 中的回调 步骤2：Go端协议调整
1. 验证 HandleSocket 函数是否正确处理v1.1消息类型：
   
   - 检查 switch msgType 分支是否包含 0x04 （异步请求）和 0x05 （异步响应）的处理逻辑。
2. 验证 handleAsyncRequest 函数是否正确解析 id 字段并回传响应：
   
   - 检查 asyncReq 结构体是否包含 ID string 字段（JSON标签 json:"id" ）
   - 确认异步响应封装时使用 versionBuf = 0x0101 、 msgTypeBuf = 0x05
   - 检查超时控制（ select 中 time.After(30 * time.Second) ）是否生效
交叉检查 ：使用Wireshark或自定义抓包工具，验证PHP发送的异步请求与Go返回的响应是否包含正确的v1.1协议头（版本号、消息类型、负载长度）。

### 二、PHP端连接池优化落地
目标 ：实现 idlePool 空闲连接复用、健康检查、最大连接数限制。 涉及文件 ： /www/wwwroot/Develop/Reader/UnixSocketReader.php （ getConnection 、 release 、 checkConnectionHealth 方法）
 步骤1：连接池逻辑验证
1. 验证 getConnection 方法优先级：
   
   - 优先从 idlePool 获取健康连接（ checkConnectionHealth 返回 true ）
   - 若 idlePool 无可用连接，从活跃池（ $this->pool ）获取或新建（不超过 $maxPoolSize=20 ）
2. 验证 release 方法：
   
   - 健康连接放回 idlePool ，异常连接直接关闭（ socket_close ）
3. 验证 checkConnectionHealth 方法：
   
   - 发送v1.0心跳包（ version=0x0100 、 msgType=0x02 ）
   - 检查 socket_write 是否成功（避免将已断开的连接放回池）
交叉检查 ：模拟高并发场景（如100+异步请求），观察 idlePool 的回收延迟（设计文档要求异步场景下连接复用），确保无连接泄漏。

### 三、Go端稳定性增强
目标 ：通过context管理goroutine生命周期、心跳响应、panic捕获。 涉及文件 ： /www/wwwroot/internal/service/socket/socket_server.go （ StartSocketServer 函数）
 步骤1：context与超时控制
1. 验证 StartSocketServer 是否使用带超时的context管理单个连接：
   
   - 检查 connCtx, cancel := context.WithTimeout(ctx, 30*time.Second)
   - 确认 defer cancel() 在goroutine中执行（避免资源泄漏）
2. 验证panic捕获逻辑：
   
   - 检查 defer func() { if r := recover(); r != nil { ... } }() 是否存在 步骤2：心跳响应实现
1. 验证 HandleSocket 函数是否处理心跳包（ msgType=0x02 ）：
   - 检查 case 0x02: sendHeartbeatAck(conn) 分支是否存在
   - 确认 sendHeartbeatAck 函数返回心跳响应（协议头与PHP端 checkConnectionHealth 发送的心跳匹配）
交叉检查 ：手动断开Go服务，观察PHP端是否触发自动重连（设计文档要求错误重试≤3次）。

### 四、异步通讯功能联调
目标 ：验证异步请求-响应链路的完整性（ID匹配、回调触发、连接复用）。 验证步骤 ：

1. 启动Go服务端（ StartSocketServer ）和PHP客户端（ UnixSocketReader ）。
2. 从PHP端发送异步请求（ sendAsync ），携带 method 和 params 。
3. 检查Go端 handleAsyncRequest 是否正确解析 id ，并通过独立goroutine处理业务逻辑。
4. 验证PHP端 asyncReadLoop 是否接收到Go返回的异步响应（ msgType=0x05 ），并触发回调函数。
5. 检查连接是否被正确释放（ release 方法将健康连接放回 idlePool ）。
异常场景验证 ：

- 手动让Go端异步处理超时（超过30秒），检查是否触发 log.Printf("异步请求ID=%s超时", id) 。
- 模拟PHP回调函数异常，检查是否影响其他异步请求（应仅当前连接被关闭，不影响连接池）。
### 五、测试用例编写与执行
目标 ：通过测试确保升级后的功能符合设计文档要求。 测试类型与执行方式 ：

测试类型 执行方法 验证点 单元测试 Go使用 testing 包（ t.Run 并发测试），PHP使用PHPUnit异步测试扩展 消息ID生成的唯一性、协议头编码/解码的正确性、异步回调的触发时机 集成测试 模拟1000+并发异步请求（使用 ab 或自定义压测工具） 连接池复用效率（ idlePool 回收延迟≤100ms）、吞吐量（设计目标需达到5000+QPS） 压力测试 构造30秒未响应的请求（Go端 processAsyncMethod 中添加 time.Sleep(35*time.Second) ） 超时日志是否触发（ ipc_async_timeout 指标增加）、资源是否回收（无goroutine泄漏）

### 六、全链路交叉检查清单
为避免协议/类型变更导致的上下文文件异常，需重点检查以下关联点：

关联文件/模块 检查项 PHP UnixSocketReader.php asyncCallbacks 的 expire 字段是否正确（30秒超时），避免内存泄漏 Go socket_receive.go requestMap 是否在超时后删除 id （防止内存溢出） Go socket_server.go listener.Close() 是否在 ctx.Done() 时触发（确保服务优雅退出） 日志与监控（设计文档第五章） Go端 HandleSocket 是否记录 id 与goroutine的关联，PHP端 sendAsync 是否记录 async_request_id