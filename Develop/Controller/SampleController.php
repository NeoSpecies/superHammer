<?php
namespace Develop\Controller;

/**
 * 示例控制器
 * 用于演示控制器的基本结构和功能
 * 提供简单的问候服务
 */
class SampleController
{
    /**
     * 问候方法
     * 处理问候请求并返回响应数据
     * 
     * 可选参数：
     * - headers: 请求头信息
     * - body: 请求体数据
     * - route: 路由信息
     * - timestamp: 时间戳
     * - client_ip: 客户端IP
     * - host: 主机信息
     * - uri: 请求URI
     * - url: 完整URL
     * 
     * @param string $body 请求体数据，JSON格式
     * @param string $uri 请求URI
     * @return string JSON格式的响应数据
     */
    public function greet($body, $uri)
    {
        $data = ["body"=>json_decode($body,true),"uri"=>$uri,"msg"=>"Hello, welcome to our service!"];
        return json_encode($data);
    }
}