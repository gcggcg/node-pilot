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
BACKEND_DIR="$PROJECT_DIR/backend"
FRONTEND_DIR="$PROJECT_DIR/frontend"
WEB_DIR="$BACKEND_DIR/web"

# 检查后端二进制是否存在
if [ ! -f "$BINARY" ]; then
    echo "开始自动构建..."
    echo ""

    # 构建前端
    echo "[1/3] 构建前端..."
    cd "$FRONTEND_DIR"
    if [ ! -d "node_modules" ]; then
        echo "安装前端依赖..."
        npm install --silent 2>/dev/null || npm install
    fi
    npm run build
    if [ $? -ne 0 ]; then
        echo "❌ 前端构建失败"
        exit 1
    fi
    echo "前端构建完成"
    echo ""

    # 复制前端到后端web目录
    echo "[2/3] 复制前端到后端..."
    mkdir -p "$WEB_DIR"
    cp -r "$FRONTEND_DIR/dist/"* "$WEB_DIR/"
    echo "前端复制完成"
    echo ""

    # 构建后端
    echo "[3/3] 构建后端..."
    cd "$BACKEND_DIR"
    go build -o node-pilot ./cmd/server
    if [ $? -ne 0 ]; then
        echo "❌ 后端构建失败"
        exit 1
    fi
    echo "后端构建完成"
    echo ""
    echo "✅ 构建流程完成!"
fi

# 启动服务
echo ""
echo "启动服务..."
cd "$BACKEND_DIR"

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
