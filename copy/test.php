<?php
// 创建Socket连接
$socket = socket_create(AF_UNIX, SOCK_STREAM, 0);
if ($socket === false) {
    echo "Failed to create socket: " . socket_strerror(socket_last_error()) . "\n";
    exit(1);
}

// 连接到Go的Socket服务
$result = socket_connect($socket, '../config/mysocket.sock');
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

// 持续监听Go的响应并回复
while (true) {
    // 读取Go的响应
    $response = socket_read($socket, 1024);

    // 检查是否成功读取了数据
    if ($response === false) {
        // 处理错误，例如记录日志或重试
        echo "Error reading from socket: " . socket_strerror(socket_last_error($socket)) . "\n";
        // 可能需要添加重试逻辑或退出循环
        continue;
    }

    // 检查是否有数据可读
    if ($response === '') {
        // 没有数据可读，可能是连接已关闭
        // echo "No data received, possibly the connection is closed.\n";
        // 可能需要退出循环或进行其他处理
        continue;
    }

    echo "Received from Go: " . $response . "\n";

    // 构造回复的JSON数据
    $reply = json_encode(array(
        "status" => "success",
        "message" => "Received your message: " . $response
    ));

    // 发送回复给Go
    if (socket_write($socket, $reply, strlen($reply)) === false) {
        // 处理写入错误
        echo "Error writing to socket: " . socket_strerror(socket_last_error($socket)) . "\n";
        // 可能需要添加重试逻辑或退出循环
        continue;
    }
}

// 关闭Socket连接
socket_close($socket);
?>
