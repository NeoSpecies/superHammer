<?php
namespace Develop\Reader;

use Exception;
class UnixSocketReader {
    protected $pool = [];
    protected $socketFile;
    protected $poolSize;
    private $container;

    public function __construct($container) {
        $this->socketFile = $container->get('socketMainFile');
        $this->poolSize = 5;

        for ($i = 0; $i < $this->poolSize; $i++) {
            $this->pool[$i] = $this->createConnection();
        }
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

    public function getConnection() {
        foreach ($this->pool as $key => $socket) {
            if($this->isConnectionAlive($socket)) {
                return $socket;
            }
        }
        // If all sockets are busy, create a new one
        $newSocket = $this->createConnection();
        $this->pool[] = $newSocket;
        return $newSocket;
    }

    public function release($socket) {
        // For now, do nothing. Add back to pool or recreate if needed.
    }

    public function sendAndReceive($message) {
        try {
            $socket = $this->getConnection();
            
            // 添加二进制协议头封装（已正确使用v1.1）
            $version = pack('n', 0x0101); // 协议版本v1.1
            $msgType = pack('C', 0x01);  // 同步消息类型
            $payload = json_encode($message);
            $payloadLen = pack('N', strlen($payload));
            
            $fullMessage = $version . $msgType . $payloadLen . $payload;
            
            if (socket_write($socket, $fullMessage, strlen($fullMessage)) === false) {
                throw new \RuntimeException("消息发送失败: " . socket_strerror(socket_last_error($socket)));
            }
            
            // 读取响应头
            $header = socket_read($socket, 7);
            if (strlen($header) !== 7) {
                throw new \RuntimeException("无效响应头");
            }
            
            // 解析响应头
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
        // 增加心跳发送逻辑
        $heartbeat = pack('nCN', 0x0101, 0x02, 0); // v1.1心跳协议
        if (!@socket_write($socket, $heartbeat, 7)) {
            return false;
        }
        return is_resource($socket) && 
               @socket_get_option($socket, SOL_SOCKET, SO_ERROR) === 0;
    }
}
