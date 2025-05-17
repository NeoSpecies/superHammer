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
        $uri =$this->input->getUri($data);
        // // 获取 headers, body, route 等...
        // $headers = $input->getHeaders();
        $request = $this->input->getRequest($data);
        // $route = $input->getRoute();
        // $timestamp = $input->getTimestamp();
        // $clientIp = $input->getClientIp();
        // $host = $input->getHost();
        // $uri = $input->getUri();
        // $url = $input->getUrl();
        // var_dump($request);
        $route = $this->router->match($uri);
        if ($route) {
            return $this->router->dispatch($route,$request);
        }
        return "No matching route found.";
    }

    public function onDataSent(&$data)
    {
        echo "Data sent: " . $data . "\n";
        // 对控制器的处理结果进行数据加工或者直接放行
    }
}
