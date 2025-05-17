#!/bin/bash

# 定义路径
SOCKET_PATH="/tmp/mysocket.sock"
MAIN_GO_PATH="./main.go"  # 替换为你的main.go文件路径
TEST_PHP_PATH="./test.php"  # 替换为你的test.php文件路径

# 存储PID的文件
MAIN_PID_FILE="duangMain.pid"
BUSSINESS_PID_FILE="duangBussiness.pid"

# 帮助信息
usage() {
    echo "Usage: $0 {start|stop|reload|start -d|help}"
    exit 1
}

# 检查并删除已存在的socket文件
cleanup_socket() {
    if [ -S "$SOCKET_PATH" ]; then
        echo "Removing existing socket file: $SOCKET_PATH"
        rm "$SOCKET_PATH"
    fi
}

# 检查端口是否被占用
check_and_free_port() {
    echo "Checking if port 80 is in use..."
    if lsof -i :80 | grep -q LISTEN; then
        echo "Port 80 is in use. Stopping the process..."
        PORT_PID=$(lsof -t -i :80)
        kill -9 $PORT_PID
        echo "Process on port 80 stopped."
    else
        echo "Port 80 is free."
    fi
}

# 启动服务
start_services() {
    check_and_free_port
    cleanup_socket

    echo "Starting main.go..."
    #nohup go run "$MAIN_GO_PATH" > /dev/null 2>&1 &
    nohup go run "$MAIN_GO_PATH" > main_output.log 2>&1 &
    sleep 5
    MAIN_PID=$(cat $MAIN_PID_FILE)
    echo "Waiting for socket to be created..."
    while [ ! -S "$SOCKET_PATH" ]; do
        sleep 1
    done
    echo "Socket created, starting test.php..."
    nohup php "$TEST_PHP_PATH" > /dev/null 2>&1 &
    sleep 1
    PHP_PID=$(pgrep -f "php $TEST_PHP_PATH" | head -n 1)
    # 记录main.go和test.php的PID到文件中
    echo "$MAIN_PID" > "$MAIN_PID_FILE"
    echo "$PHP_PID" >> "$BUSSINESS_PID_FILE"
    echo "Services started. Main.go PID: $MAIN_PID, PHP script PID: $PHP_PID"
}

# 停止服务
stop_services() {
    if [ -f "$MAIN_PID_FILE" ] && [ -f "$BUSSINESS_PID_FILE" ]; then
        MAIN_PID=$(cat "$MAIN_PID_FILE" | cut -d' ' -f2)  # 提取PID部分
        PHP_PID=$(cat "$BUSSINESS_PID_FILE" | cut -d' ' -f2)  # 提取PID部分
        echo "Stopping main.go with PID: $MAIN_PID"
        kill -s TERM $MAIN_PID
        sleep 5
        if ps -p $MAIN_PID > /dev/null; then
            echo "main.go did not stop gracefully, sending SIGKILL"
            kill -s KILL $MAIN_PID
        else
            echo "main.go stopped"
        fi
        rm "$MAIN_PID_FILE"
        echo "Stopping test.php with PID: $PHP_PID"
        kill -s TERM $PHP_PID
        sleep 5
        if pgrep -f "$TEST_PHP_PATH" > /dev/null; then
            echo "test.php did not stop gracefully, sending SIGKILL"
            kill -s KILL $PHP_PID
        else
            echo "test.php stopped"
        fi
        rm "$BUSSINESS_PID_FILE"
        cleanup_socket
        echo "Services stopped and socket file removed."
    else
        echo "Services are not running."
    fi
}

# 主逻辑
case $1 in
    start)
        start_services
        ;;
    stop)
        stop_services
        ;;
    reload)
        stop_services
        start_services
        ;;
    start)
        shift
        if [ "$1" = "-d" ]; then
            start_services
            tail -f /dev/null &
        else
            usage
        fi
        ;;
    help)
        usage
        ;;
    *)
        usage
        ;;
esac

exit 0
