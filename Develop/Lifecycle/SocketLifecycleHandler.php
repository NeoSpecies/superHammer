<?php
namespace Develop\Lifecycle;


class SocketLifecycleHandler
{
    private $router;
    private $container; // 用于存储容器实例
    private $input;

    public function __construct($container)
    {
        $this->router = $container->get('router');
        $this->input = $container->get('input');
    }
    
    public function onServerStart()
    {
        echo "Server started.\n";
    }

    public function onServerStop()
    {
        echo "Server stopped.\n";
    }

    public function onDataReceived(&$data) {
        // echo "Data received: " . $data . "\n";
        $data = substr($data, 7);
        // 解析原始数据为数组（假设数据是JSON格式）
        $this->receivedData = json_decode($data, true);
        
        if (json_last_error() !== JSON_ERROR_NONE) {
            error_log("数据解析失败: " . json_last_error_msg());
            $this->receivedData = [];
        }
        
        return true;
    }
    

    public function onCtrl($data) {
        
        // 直接使用 onDataReceived 中存储的解析后数据
        if (empty($this->receivedData)) {
            error_log("未接收到有效请求数据");
            return "无效请求";
        }
        
        // 从已解析的数据中提取URI
        $uri = $this->receivedData['params']['uri'] ?? '';
        $request = $this->receivedData;
        
        // error_log("Trying to match route for URI: " . $uri);
        $route = $this->router->match($uri);
        
        if ($route) {
            // error_log("Route matched: " . print_r($route, true));
            return $this->router->dispatch($route, $request);
        }
        
        error_log("No route found for URI: " . $uri);
        return "No matching route found.";
    }

    public function onDataSent(&$data)
    {
        // 构造协议头（版本v1.1，响应类型0x05）
        $version = pack('n', 0x0101);
        $msgType = pack('C', 0x05); // 响应类型
        $payloadLen = pack('N', strlen($data));
        
        // 组合协议头和数据
        $data = $version . $msgType . $payloadLen . $data;
        
        // 添加调试日志（正式环境可移除）
        error_log("完整响应数据长度：" . strlen($data));
    }
}
