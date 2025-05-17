<?php
// 创建Socket连接
$socket = socket_create(AF_UNIX, SOCK_STREAM, 0);
if ($socket === false) {
    echo "Failed to create socket: " . socket_strerror(socket_last_error()) . "\n";
    exit(1);
}

// 连接到Go的Socket服务
$result = socket_connect($socket, '/tmp/mysocket.sock');
if ($result === false) {
    echo "Failed to connect to socket: " . socket_strerror(socket_last_error($socket)) . "\n";
    exit(1);
}

// 向Go发送请求
$request = json_encode(array(
    "method" => "exec",
    "params" => array("php", "echo 'Hello from Go!'")
));

socket_write($socket, $request, strlen($request));

// 读取Go的响应
$response = socket_read($socket, 1024);
echo "Response from Go: " . $response . "\n";

// 关闭Socket连接
socket_close($socket);
?>