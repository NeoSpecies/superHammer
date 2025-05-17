<?php
namespace Develop\Tool;

use Exception;

/**
 * 输入处理工具类
 * 用于处理和管理HTTP请求的输入数据
 * 通过Unix Socket与主服务进行通信获取请求信息
 */
class Input
{
    /**
     * @var array|null 存储请求数据的数组
     */
    private $requestData;

    /**
     * @var object 依赖注入容器实例
     */
    private $container;

    /**
     * 构造函数
     * 
     * @param object $container 依赖注入容器实例
     */
    public function __construct($container)
    {
        $this->container = $container;
    }

    /**
     * 连接Unix Socket并读取数据
     * 
     * @param string $stringToSend 要发送的键值
     * @throws Exception 当Socket通信失败时抛出异常
     */
    private function connectAndReadData($stringToSend)
    {
        // 获取UnixSocketReader实例，该方法假设容器已正确配置
        $reader = $this->container->get('unixSocketReader');
        $data = json_encode([
            'service' => 'input',
            'method' => 'input',
            'params' => ['key' => $stringToSend]
        ]);
        try {
            // 使用sendAndReceive方法发送数据并等待响应
            $response = $reader->sendAndReceive($data);
            $response = json_decode($response, true);
            if($response['status']==200){
                $this->requestData = json_decode($response['data'],true);
            }else{
                $this->requestData = [];
            }
        } catch (Exception $e) {
            error_log("Failed to read data from Unix socket: " . $e->getMessage());
            // 设置一个状态或者标记以便外部检查是否成功
            $this->requestData = null;
        }
    }

    /**
     * 获取请求头信息
     * 
     * @param string $stringToSend 要发送的键值
     * @return array 请求头信息数组
     */
    public function getHeaders($stringToSend)
    {
        $this->connectAndReadData($stringToSend);
        return $this->requestData['headers'] ?? [];
    }

    /**
     * 获取请求体数据
     * 
     * @param string $stringToSend 要发送的键值
     * @return mixed|null 请求体数据
     */
    public function getBody($stringToSend)
    {
        $this->connectAndReadData($stringToSend);
        return $this->requestData['body'] ?? null;
    }

    /**
     * 获取路由信息
     * 
     * @param string $stringToSend 要发送的键值
     * @return string|null 路由信息
     */
    public function getRoute($stringToSend)
    {
        $this->connectAndReadData($stringToSend);
        return $this->requestData['route'] ?? null;
    }

    /**
     * 获取请求时间戳
     * 
     * @param string $stringToSend 要发送的键值
     * @return string|null 请求时间戳
     */
    public function getTimestamp($stringToSend)
    {
        $this->connectAndReadData($stringToSend);
        return $this->requestData['timestamp'] ?? null;
    }

    /**
     * 获取客户端IP地址
     * 
     * @param string $stringToSend 要发送的键值
     * @return string|null 客户端IP地址
     */
    public function getClientIp($stringToSend)
    {
        $this->connectAndReadData($stringToSend);
        return $this->requestData['client_ip'] ?? null;
    }

    /**
     * 获取请求主机名
     * 
     * @param string $stringToSend 要发送的键值
     * @return string|null 请求主机名
     */
    public function getHost($stringToSend)
    {
        $this->connectAndReadData($stringToSend);
        return $this->requestData['host'] ?? null;
    }

    /**
     * 获取请求URI
     * 
     * @param string $stringToSend 要发送的键值
     * @return string|null 请求URI
     */
    public function getUri($stringToSend)
    {
        $this->connectAndReadData($stringToSend);
        return $this->requestData['uri'] ?? null;
    }

    /**
     * 获取完整URL
     * 
     * @param string $stringToSend 要发送的键值
     * @return string|null 完整URL
     */
    public function getUrl($stringToSend)
    {
        $this->connectAndReadData($stringToSend);
        return $this->requestData['url'] ?? null;
    }

    /**
     * 获取完整的请求数据
     * 
     * @param string $stringToSend 要发送的键值
     * @return array|null 完整的请求数据数组
     */
    public function getRequest($stringToSend)
    {
        $this->connectAndReadData($stringToSend);
        return $this->requestData ?? null;
    }
}