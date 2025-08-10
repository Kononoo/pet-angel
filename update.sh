#!/usr/bin/env bash
# 一键更新项目脚本
# 用法：./update.sh [--daemon]
# 默认拉取 master 分支，编译并运行

set -euo pipefail

# 设置 PATH
export PATH="/usr/local/go/bin:$HOME/go/bin:/usr/local/bin:$PATH"

# 检查是否在项目根目录
if [[ ! -d .git ]] || [[ ! -f Makefile ]]; then
    echo "[ERROR] 请在项目根目录执行此脚本" >&2
    exit 1
fi

echo "[INFO] 开始更新项目..."

# 拉取最新代码
echo "[INFO] 拉取 master 分支最新代码..."
git fetch --all --prune
git checkout master
git pull --ff-only origin master

# 检查并安装依赖
if ! command -v go >/dev/null 2>&1; then
    echo "[ERROR] Go 未安装，请先安装 Go" >&2
    exit 1
fi

# 生成代码
echo "[INFO] 生成代码..."
if command -v protoc >/dev/null 2>&1; then
    make all
else
    make wire
    make generate
    go mod tidy
fi

# 编译项目
echo "[INFO] 编译项目..."
go build -o ./bin/pet-angel ./cmd/pet-angel

# 运行项目
if [[ "${1:-}" == "--daemon" ]]; then
    echo "[INFO] 后台启动项目..."
    mkdir -p ./logs
    chmod +x ./bin/pet-angel
    
    # 停止旧进程（如果存在）
    if [[ -f ./logs/pet-angel.pid ]]; then
        OLD_PID=$(cat ./logs/pet-angel.pid)
        if kill -0 "$OLD_PID" 2>/dev/null; then
            echo "[INFO] 停止旧进程 PID: $OLD_PID"
            kill "$OLD_PID"
            sleep 2
        fi
    fi
    
    # 启动新进程
    nohup ./bin/pet-angel -conf ./configs > ./logs/pet-angel.out 2>&1 & echo $! > ./logs/pet-angel.pid
    NEW_PID=$(cat ./logs/pet-angel.pid)
    echo "[OK] 项目已后台启动，PID: $NEW_PID"
    echo "[INFO] 日志文件: $(pwd)/logs/pet-angel.out"
    echo "[INFO] 查看日志: tail -f ./logs/pet-angel.out"
else
    echo "[INFO] 前台启动项目..."
    exec ./bin/pet-angel -conf ./configs
fi 