#!/usr/bin/env bash
# 一键更新并运行（kratos run）——在项目根目录内执行即可
# 用法：
#   bash scripts/server_update_run.sh [branch] [--daemon]
# 例子：
#   bash scripts/server_update_run.sh master
#   bash scripts/server_update_run.sh master --daemon

set -euo pipefail

BRANCH="${1:-master}"
DAEMON="${2:-}"

# 1) 校验当前目录
REPO_DIR="$(pwd)"
if [[ ! -d .git ]]; then
  echo "[ERROR] 当前目录不是 Git 仓库，请先 git clone 项目后在仓库根目录执行本脚本" >&2
  exit 1
fi
if [[ ! -f Makefile ]]; then
  echo "[ERROR] 未找到 Makefile，请在项目根目录执行本脚本" >&2
  exit 1
fi

# 2) 拉取最新代码
echo "[INFO] Git update: branch=$BRANCH ..."
git fetch --all --prune
git checkout "$BRANCH"
git pull --ff-only origin "$BRANCH"

# 3) 准备工具链（kratos 在 make init 里安装）
export PATH="$(go env GOPATH)/bin:/usr/local/go/bin:$PATH"
if ! command -v kratos >/dev/null 2>&1; then
  echo "[INFO] 安装工具链（kratos/protoc-plugins/wire/openapi）..."
  make init || true
fi

# 4) 生成代码并整理依赖
if command -v protoc >/dev/null 2>&1; then
  echo "[INFO] 执行完整生成链路：make all (api/service/openapi/wire/generate)"
  make all || true
else
  echo "[INFO] 未检测到 protoc，执行最小链路：wire + generate + go mod tidy"
  make wire || true
  make generate || true
  go mod tidy || true
fi

# 5) 运行（kratos run）
CMD=(kratos run -conf ./configs)
echo "[INFO] Starting: ${CMD[*]}"
if [[ "$DAEMON" == "--daemon" ]]; then
  # 后台运行到 nohup 日志
  mkdir -p ./logs || true
  nohup "${CMD[@]}" > ./logs/pet-angel.out 2>&1 & disown
  echo "[OK] started in background. Logs: $(pwd)/logs/pet-angel.out"
else
  exec "${CMD[@]}"
fi

