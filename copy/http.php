<?php

// 创建监听套接字
$server = stream_socket_server("tcp://localhost:9999", $errno, $errstr);

if (!$server) {
    die("Error starting server: $errstr");
}

echo "Server started. Listening on port 9999...\n";

while (true) {
    // 接受客户端连接
    $client = stream_socket_accept($server);

    // 读取客户端发送的数据
    $data = fread($client, 1024);

    // 解析 HTTP 请求
    $request = parse_http_request($data);

    // 处理请求
    $response = handle_request($request);

    // 发送响应到客户端
    fwrite($client, $response);

    // 关闭客户端连接
    fclose($client);
}

// 关闭服务器
fclose($server);

// 解析 HTTP 请求
function parse_http_request($data) {
    $request = [];

    list($request["method"], $request["uri"], $request["protocol"]) = explode(" ", $data);
    $request["headers"] = [];
    $request["body"] = "";

    $headersAndBody = explode("\r\n\r\n", $data, 2);
    if (count($headersAndBody) == 2) {
        list($headers, $body) = $headersAndBody;
        $headerLines = explode("\r\n", $headers);
        foreach ($headerLines as $line) {
            list($key, $value) = explode(": ", $line);
            $request["headers"][$key] = $value;
        }
        $request["body"] = $body;
    }

    return $request;
}

// 处理 HTTP 请求
function handle_request($request) {
    $response = "HTTP/1.1 200 OK\r\n";
    $response .= "Content-Type: text/plain\r\n";
    $response .= "\r\n";
    $response .= "Hello, client! You requested: " . $request["uri"];

    return $response;
}
