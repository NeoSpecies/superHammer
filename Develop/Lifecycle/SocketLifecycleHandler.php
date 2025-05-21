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

    public function onDataReceived(&$data)
    {
        echo "Data received: " . $data . "\n";
        // 使用已经注入的 reader 实例
        // $this->input
        // 进行数据处理...
        return true;
    }

    public function onCtrl($data)
    {
        $uri = $this->input->getUri($data);
        $request = $this->input->getRequest($data);
        // 添加调试日志
        error_log("Trying to match route for URI: " . $uri);
        
        $route = $this->router->match($uri);
        if ($route) {
            error_log("Route matched: " . print_r($route, true));
            return $this->router->dispatch($route, $request);
        }
        
        error_log("No route found for URI: " . $uri);
        return "No matching route found.";
    }

    public function onDataSent(&$data)
    {
        echo "Data sent: " . $data . "\n";
        // 对控制器的处理结果进行数据加工或者直接放行
    }
}
