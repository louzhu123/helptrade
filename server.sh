#!/bin/bash

# 应用程序的名称或路径
APP_NAME="helptrade"  # 你可以设置成你的应用程序名称或完整路径
#APP_PATH="go run ."  # 如果APP_NAME是路径，这一行可以删除
APP_PATH="./helptrade"

# 获取传递的参数
ACTION=$1

# 获取应用程序的PID
get_pid() {
    echo $(pgrep -f $APP_NAME)
}

# 启动应用程序
start_app() {
    PID=$(get_pid)
    if [ -z "$PID" ]; then
        echo "Starting $APP_NAME..."
        nohup $APP_PATH > system.log 2>&1 &
        echo "$APP_NAME started with PID $(get_pid)"
    else
        echo "$APP_NAME is already running with PID $PID"
    fi
}

# 停止应用程序
stop_app() {
    PID=$(get_pid)
    if [ -z "$PID" ]; then
        echo "$APP_NAME is not running."
    else
        echo "Stopping $APP_NAME..."
        kill $PID
        echo "$APP_NAME stopped."
    fi
}

# 重启应用程序
restart_app() {
    echo "Restarting $APP_NAME..."
    stop_app
    start_app
}

# 检查应用程序状态
status_app() {
    PID=$(get_pid)
    if [ -z "$PID" ]; then
        echo "$APP_NAME is not running."
    else
        echo "$APP_NAME is running with PID $PID."
    fi
}

# 根据传入的参数执行操作
case $ACTION in
    start)
        start_app
        ;;
    stop)
        stop_app
        ;;
    restart)
        restart_app
        ;;
    status)
        status_app
        ;;
    *)
        echo "Usage: $0 {start|stop|restart|status}"
        exit 1
        ;;
esac
