#!/bin/bash
#
# NodePilot 启动脚本
# 用法: ./start.sh [FRONTEND_IP] [FRONTEND_PORT] [DEBUG]
# 示例:
#   ./start.sh                 # 默认模式
#   ./start.sh 127.0.0.1 8080  # 指定IP和端口
#   ./start.sh 127.0.0.1 8080 debug  # 开启Debug日志
#

set -e

# 默认值
FRONTEND_IP="${1:-127.0.0.1}"
FRONTEND_PORT="${2:-8080}"
DEBUG_MODE="${3:-}"

# 获取脚本所在目录
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_DIR="$(dirname "$SCRIPT_DIR")"
DATA_DIR="$PROJECT_DIR/data"
LOG_FILE="$DATA_DIR/node-pilot.log"

# 确保数据目录存在
mkdir -p "$DATA_DIR"

echo "=========================================="
echo "  NodePilot 批量服务器管理平台"
echo "=========================================="
echo ""
echo "前端地址: http://${FRONTEND_IP}:${FRONTEND_PORT}"
echo "数据目录: ${DATA_DIR}"
echo "日志文件: ${LOG_FILE}"
if [ "$DEBUG_MODE" = "debug" ]; then
    echo "Debug模式: 开启"
else
    echo "Debug模式: 关闭"
fi
echo ""

BINARY="$PROJECT_DIR/backend/node-pilot"
if [ ! -f "$BINARY" ]; then
    echo "未找到二进制文件，开始自动构建..."
    cd "$PROJECT_DIR/backend"
    go build -o node-pilot ./cmd/server
    echo "构建完成!"
fi

# 启动服务
echo "启动服务..."
cd "$PROJECT_DIR/backend"

# 根据是否开启debug模式来设置参数
CMD="$BINARY --db $DATA_DIR/servers.db --listen :${FRONTEND_PORT} --log $LOG_FILE"

if [ "$DEBUG_MODE" = "debug" ]; then
    CMD="$CMD --debug"
fi

eval "$CMD" >> "$LOG_FILE" 2>&1 &

PID=$!
echo "服务已启动 (PID: $PID)"
echo ""

# 等待服务启动
sleep 2

# 检查服务是否成功启动
if kill -0 $PID 2>/dev/null; then
    echo "✅ 服务启动成功!"
    echo ""
    echo "请访问: http://${FRONTEND_IP}:${FRONTEND_PORT}"
    echo ""
    if [ "$DEBUG_MODE" = "debug" ]; then
        echo "🔍 Debug模式已开启，日志文件: $LOG_FILE"
        echo ""
        echo "查看实时日志: tail -f $LOG_FILE"
    fi
    echo ""
    echo "按 Ctrl+C 停止服务"
    
    # 捕获退出信号
    trap "echo '正在停止服务...'; kill $PID 2>/dev/null; exit 0" SIGINT SIGTERM
    
    # 保持运行
    wait $PID
else
    echo "❌ 服务启动失败，请检查日志: $LOG_FILE"
    exit 1
fi
