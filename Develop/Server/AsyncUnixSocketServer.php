<?php
namespace Develop\Server;

class AsyncUnixSocketServer
{
    private $container; // 用于存储容器实例
    private $socketFile;
    private $pidFile;
    private $base;
    private $event;
    private $socket;
    private $lifecycleHandler;

    public function __construct($socketFile, $pidFile, $container)
    {
        $this->socketFile = $socketFile;
        $this->pidFile = $pidFile;
        $this->container = $container; // 保存容器实例
    }

    public function registerLifecycleHandler($handler)
    {
        // 检查handler是否实现了所有必要的生命周期方法
        if (method_exists($handler, 'onServerStart') && 
            method_exists($handler, 'onServerStop') &&
            method_exists($handler, 'onDataReceived') &&
            method_exists($handler, 'onCtrl') &&
            method_exists($handler, 'onDataSent')) {
            $this->lifecycleHandler = $handler;
        } else {
            throw new \Exception("Handler does not implement required lifecycle methods.");
        }
    }

    public function start()
    {
        $this->createSocket();
        $this->setupEventHandler();
        $this->writePidToFile($this->pidFile);
        if ($this->lifecycleHandler) {
            $this->lifecycleHandler->onServerStart();
        }
        $this->base->loop();
    }

    public function stop()
    {
        if ($this->lifecycleHandler) {
            $this->lifecycleHandler->onServerStop();
        }
        if (is_resource($this->socket)) {
            socket_close($this->socket);
        }
        if ($this->event) {
            $this->event->free();
        }
        if ($this->base) {
            $this->base->exit();
        }
        if (file_exists($this->socketFile)) {
            unlink($this->socketFile);
        }
        if (file_exists($this->pidFile)) {
            unlink($this->pidFile);
        }
    }

    private function createSocket()
    {
        $this->socket = socket_create(AF_UNIX, SOCK_STREAM, 0);
        if ($this->socket === false) {
            throw new \Exception("socket_create() failed: reason: " . socket_strerror(socket_last_error()));
        }
        if (file_exists($this->socketFile)) {
            unlink($this->socketFile);
        }
        if (!socket_bind($this->socket, $this->socketFile)) {
            throw new \Exception("socket_bind() failed: reason: " . socket_strerror(socket_last_error($this->socket)));
        }
        if (!socket_listen($this->socket, 5)) {
            throw new \Exception("socket_listen() failed: reason: " . socket_strerror(socket_last_error($this->socket)));
        }
        if (!socket_set_nonblock($this->socket)) {
            throw new \Exception("socket_set_nonblock() failed: reason: " . socket_strerror(socket_last_error($this->socket)));
        }
    }

    private function setupEventHandler()
    {
        $this->base = new \EventBase();
        $this->event = new \Event($this->base, $this->socket, \Event::READ | \Event::PERSIST, function() {
            $clientSocket = socket_accept($this->socket);
            if ($clientSocket === false) {
                echo "socket_accept() failed: reason: " . socket_strerror(socket_last_error($this->socket)) . "\n";
                return;
            }
            socket_set_nonblock($clientSocket);
            $data = '';
            while (true) {
                $read = socket_read($clientSocket, 1024);
                if ($read === false) {
                    if (socket_last_error($clientSocket) == SOCKET_EAGAIN) {
                        break;
                    }
                    echo "socket_read() failed: reason: " . socket_strerror(socket_last_error($clientSocket)) . "\n";
                    socket_close($clientSocket);
                    return;
                }
                $data .= $read;
                if (strlen($read) < 1024) {
                    break;
                }
            }
            if (!empty($data)) {
                $data = trim($data);
                // 触发onDataReceived周期
                if ($this->lifecycleHandler && $this->lifecycleHandler->onDataReceived($data)) {
                    // 在生命周期处理类中处理路由和控制逻辑
                    $response = $this->lifecycleHandler->onCtrl($data);
                    // 触发onDataSent周期
                    if ($this->lifecycleHandler) {
                        $this->lifecycleHandler->onDataSent($response);
                    }
                    socket_write($clientSocket, $response, strlen($response));
                }
            }
            socket_close($clientSocket);
        });
        $this->event->add();
    }

    private function writePidToFile($filePath)
    {
        $pid = getmypid();
        file_put_contents($filePath, $pid);
        echo "PID {$pid} 已写入文件 {$filePath}\n";
    }
}
