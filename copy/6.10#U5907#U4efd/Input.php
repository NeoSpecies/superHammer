<?php
namespace Develop\Tool;

use Exception;

class Input
{
    private $requestData;
    private $container;

    public function __construct($container)
    {
        $this->container = $container;
    }

    private function connectAndReadData($stringToSend)
    {
        // 获取UnixSocketReader实例，该方法假设容器已正确配置
        $reader = $this->container->get('unixSocketReader');
        try {
            // 使用sendAndReceive方法发送数据并等待响应
            $response = $reader->sendAndReceive($stringToSend);
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

    public function getHeaders($stringToSend)
    {
        $this->connectAndReadData($stringToSend);
        return $this->requestData['headers'] ?? [];
    }

    public function getBody($stringToSend)
    {
        $this->connectAndReadData($stringToSend);
        return $this->requestData['body'] ?? null;
    }

    public function getRoute($stringToSend)
    {
        $this->connectAndReadData($stringToSend);
        return $this->requestData['route'] ?? null;
    }

    public function getTimestamp($stringToSend)
    {
        $this->connectAndReadData($stringToSend);
        return $this->requestData['timestamp'] ?? null;
    }

    public function getClientIp($stringToSend)
    {
        $this->connectAndReadData($stringToSend);
        return $this->requestData['client_ip'] ?? null;
    }

    public function getHost($stringToSend)
    {
        $this->connectAndReadData($stringToSend);
        return $this->requestData['host'] ?? null;
    }

    public function getUri($stringToSend)
    {
        $this->connectAndReadData($stringToSend);
        return $this->requestData['uri'] ?? null;
    }

    public function getUrl($stringToSend)
    {
        $this->connectAndReadData($stringToSend);
        return $this->requestData['url'] ?? null;
    }
    public function getRequest($stringToSend)
    {
        $this->connectAndReadData($stringToSend);
        return $this->requestData ?? null;
    }
}
