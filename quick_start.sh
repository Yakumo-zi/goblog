#!/bin/bash

# GoBlog 快速启动入口脚本
# 这个脚本位于项目根目录，方便用户直接运行

echo "🚀 启动 GoBlog PostgreSQL 版本..."
echo ""

# 检查scripts目录中的完整脚本是否存在
if [ -f "./scripts/quick_start.sh" ]; then
    # 运行完整的快速启动脚本
    exec ./scripts/quick_start.sh "$@"
else
    echo "❌ 错误：找不到 ./scripts/quick_start.sh 脚本"
    echo ""
    echo "请确保您在正确的项目根目录中运行此脚本。"
    echo ""
    echo "项目结构应该是："
    echo "goblog/"
    echo "├── quick_start.sh       # 这个脚本"
    echo "├── scripts/"
    echo "│   └── quick_start.sh   # 完整的启动脚本"
    echo "├── docker-compose.yml"
    echo "└── ..."
    echo ""
    exit 1
fi 