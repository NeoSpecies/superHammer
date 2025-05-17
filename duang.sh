#!/bin/bash

# 配置项（根据实际环境调整）
PROJECT_DIR="/www/wwwroot"          # 项目根目录
GO_BIN="/usr/local/go/bin/go"       # Go 可执行文件路径
PHP_BIN="/usr/bin/php"              # PHP 可执行文件路径
GO_PID_FILE="${PROJECT_DIR}/go.pid" # Go 进程 PID 文件
PHP_PID_FILE="${PROJECT_DIR}/php.pid" # PHP 进程 PID 文件
LOG_FILE="${PROJECT_DIR}/duang.log"  # 日志文件

# 日志函数
log() {
    echo "[$(date '+%Y-%m-%d %H:%M:%S')] $1" >> "${LOG_FILE}"
}

# 检查进程是否存在
check_process() {
    local pid_file=$1
    if [ -f "${pid_file}" ]; then
        local pid=$(cat "${pid_file}")
        if ps -p "${pid}" > /dev/null; then
            return 0  # 进程存在
        else
            rm -f "${pid_file}"
            return 1  # 进程不存在但 PID 文件残留
        fi
    else
        return 1  # 无 PID 文件
    fi
}

# 启动 Go 进程
start_go() {
    if check_process "${GO_PID_FILE}"; then
        log "Go 进程已运行（PID: $(cat ${GO_PID_FILE})）"
        return
    fi
    cd "${PROJECT_DIR}" && ${GO_BIN} run . &
    local go_pid=$!
    echo "${go_pid}" > "${GO_PID_FILE}"
    log "启动 Go 进程（PID: ${go_pid}）"
}

# 启动 PHP 进程（修改后）
start_php() {
    if check_process "${PHP_PID_FILE}"; then
        log "PHP 进程已运行（PID: $(cat ${PHP_PID_FILE})）"
        return
    fi
    # 调整路径为 Develop 目录，并指定入口文件 bootstrap.php
    cd "${PROJECT_DIR}/Develop" && ${PHP_BIN} -S 0.0.0.0:8000 bootstrap.php &
    local php_pid=$!
    echo "${php_pid}" > "${PHP_PID_FILE}"
    log "启动 PHP 进程（PID: ${php_pid}）"
}

# 停止 Go 进程
stop_go() {
    if check_process "${GO_PID_FILE}"; then
        local pid=$(cat "${GO_PID_FILE}")
        kill "${pid}"
        rm -f "${GO_PID_FILE}"
        log "停止 Go 进程（PID: ${pid}）"
    else
        log "Go 进程未运行"
    fi
}

# 停止 PHP 进程
stop_php() {
    if check_process "${PHP_PID_FILE}"; then
        local pid=$(cat "${PHP_PID_FILE}")
        kill "${pid}"
        rm -f "${PHP_PID_FILE}"
        log "停止 PHP 进程（PID: ${pid}）"
    else
        log "PHP 进程未运行"
    fi
}

# 主逻辑
case "$1" in
    start)
        start_go
        start_php
        ;;
    stop)
        stop_go
        stop_php
        ;;
    restart)
        stop_go
        stop_php
        start_go
        start_php
        ;;
    status)
        check_process "${GO_PID_FILE}" && echo "Go 进程运行中（PID: $(cat ${GO_PID_FILE})）" || echo "Go 进程未运行"
        check_process "${PHP_PID_FILE}" && echo "PHP 进程运行中（PID: $(cat ${PHP_PID_FILE})）" || echo "PHP 进程未运行"
        ;;
    *)
        echo "用法: $0 {start|stop|restart|status}"
        exit 1
        ;;
esac
