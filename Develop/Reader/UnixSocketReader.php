<?php
namespace Develop\Reader;

use Exception;
class UnixSocketReader {
    protected $pool = [];
    protected $socketFile;
    protected $poolSize;
    private $container;
    protected $activePool = [];
    protected $maxPoolSize = 20;
    protected $idleTimeout = 300; // 秒
    protected $idlePool = []; // 格式改为[[socket资源, lastUsed时间戳], ...]

    public function __construct($container) {
        $this->socketFile = $container->get('socketMainFile');
        $this->poolSize = $container->get('socketPoolSize') ?? 5; // 移除has检查
        $this->maxPoolSize = $container->get('maxPoolSize') ?? 20;
        
        for ($i = 0; $i < $this->poolSize; $i++) {
            $this->idlePool[] = $this->createConnection();
        }
    }

    public function getConnection() {
        // 优先从空闲池获取
        while (!empty($this->idlePool)) {
            $socket = array_pop($this->idlePool);
            if ($this->isConnectionAlive($socket)) {
                $this->activePool[] = $socket;
                return $socket;
            }
            socket_close($socket);
        }
        
        // 创建新连接（不超过最大限制）
        if (count($this->activePool) + count($this->idlePool) < $this->maxPoolSize) {
            $newSocket = $this->createConnection();
            $this->activePool[] = $newSocket;
            return $newSocket;
        }
        
        throw new \RuntimeException("连接池已达最大限制 ({$this->maxPoolSize})");
    }

    public function release($socket) {
        $index = array_search($socket, array_column($this->activePool, 0), true);
        if ($index !== false) {
            list($socket) = array_splice($this->activePool, $index, 1);
            
            if ($this->isConnectionAlive($socket)) {
                // 存储连接时记录当前时间戳
                $this->idlePool[] = [$socket, time()];
            } else {
                socket_close($socket);
            }
        }
    }

    private function cleanupIdleConnections() {
        $now = time();
        foreach ($this->idlePool as $i => $item) {
            list($socket, $lastUsed) = $item;
            if ($now - $lastUsed > $this->idleTimeout) {
                socket_close($socket);
                unset($this->idlePool[$i]);
            }
        }
        $this->idlePool = array_values($this->idlePool);
    }

    protected function createConnection() {
        $socket = socket_create(AF_UNIX, SOCK_STREAM, 0);
        if (!$socket) {
            throw new \RuntimeException("Unable to create socket: " . socket_strerror(socket_last_error()));
        }
        if (!socket_connect($socket, $this->socketFile)) {
            socket_close($socket);
            throw new \RuntimeException("Unable to connect: " . socket_strerror(socket_last_error($socket)));
        }
        return $socket;
    }


    public function sendAndReceive($message) {
        try {
            $socket = $this->getConnection();
            
            // 修改协议头版本为v1.1（原为0x0100）
            $version = pack('n', 0x0101); // 协议版本v1.1
            $msgType = pack('C', 0x01);  // 同步消息类型
            $payload = json_encode($message);
            $payloadLen = pack('N', strlen($payload));
            
            $fullMessage = $version . $msgType . $payloadLen . $payload;
            
            if (socket_write($socket, $fullMessage, strlen($fullMessage)) === false) {
                throw new \RuntimeException("消息发送失败: " . socket_strerror(socket_last_error($socket)));
            }
            
            // 读取响应头时保持v1.1版本校验
            $header = socket_read($socket, 7);
            if (strlen($header) !== 7) {
                throw new \RuntimeException("无效响应头");
            }
            
            // 更新版本号解析注释
            list($version, $msgType, $payloadLen) = array_values(unpack('nversion/CmsgType/NpayloadLen', $header));
            
            // 读取响应体
            $response = socket_read($socket, $payloadLen);
            
            if ($response === false) {
                throw new \RuntimeException("响应读取失败: " . socket_strerror(socket_last_error($socket)));
            }
            
            $this->release($socket);
            return $response;
        } catch (\Exception $e) {
            error_log("Socket通信错误: " . $e->getMessage());
            throw $e;
        }
    }

    public function __destruct() {
        // 增加泄漏告警
        $leakCount = count($this->pool);
        if ($leakCount > 0) {
            error_log("连接池泄漏警告: 剩余{$leakCount}个未释放连接");
        }
        foreach ($this->pool as $socket) {
            if (is_resource($socket)) {
                socket_close($socket);
            }
        }
    }

    private function isConnectionAlive($socket) {
        // 使用v1.1心跳协议（与Go端同步）
        $heartbeat = pack('nCN', 0x0101, 0x02, 0); // 版本v1.1，类型0x02（心跳），负载长度0
        if (!@socket_write($socket, $heartbeat, 7)) {
            return false;
        }
        // 增加心跳响应验证（可选，根据Go端是否返回心跳确认）
        $response = @socket_read($socket, 7); // 读取心跳响应头
        return $response !== false && is_resource($socket);
    }
    
    protected function asyncReadLoop($socket) {
        while (true) {
            // ...
            
            // 每轮循环清理过期回调（每5秒执行一次）
            if (time() % 5 === 0) {
                $now = time();
                foreach ($this->asyncCallbacks as $id => $callbackData) {
                    if ($now > $callbackData['expire']) {
                        unset($this->asyncCallbacks[$id]);
                        error_log("异步回调ID={$id}超时未响应");
                    }
                }
            }
        }
    }
}
