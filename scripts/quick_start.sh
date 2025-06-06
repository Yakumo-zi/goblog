#!/bin/bash

# GoBlog PostgreSQL 快速启动脚本
# 这个脚本会自动设置和启动完整的博客系统

set -e

# 颜色输出
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# 日志函数
log_info() {
    echo -e "${GREEN}[INFO]${NC} $1"
}

log_warn() {
    echo -e "${YELLOW}[WARN]${NC} $1"
}

log_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

log_step() {
    echo -e "${BLUE}[STEP]${NC} $1"
}

# 检查依赖
check_dependencies() {
    log_step "检查系统依赖..."
    
    # 检查Docker
    if ! command -v docker &> /dev/null; then
        log_error "Docker未安装，请先安装Docker"
        exit 1
    fi
    
    # 检查Docker Compose
    if ! command -v docker-compose &> /dev/null; then
        log_error "Docker Compose未安装，请先安装Docker Compose"
        exit 1
    fi
    
    # 检查Go
    if ! command -v go &> /dev/null; then
        log_warn "Go未安装，某些开发功能可能不可用"
    else
        log_info "Go版本: $(go version)"
    fi
    
    log_info "所有依赖检查完成"
}

# 停止现有服务
stop_existing_services() {
    log_step "停止现有服务..."
    
    if docker-compose ps | grep -q "goblog"; then
        log_info "发现运行中的服务，正在停止..."
        docker-compose down
    fi
}

# 构建和启动服务
start_services() {
    log_step "启动PostgreSQL数据库和博客服务..."
    
    # 启动服务
    docker-compose up -d
    
    log_info "等待服务启动..."
    sleep 10
    
    # 检查服务状态
    if docker-compose ps | grep -q "Up"; then
        log_info "服务启动成功！"
    else
        log_error "服务启动失败，请检查日志"
        docker-compose logs
        exit 1
    fi
}

# 等待数据库就绪
wait_for_database() {
    log_step "等待PostgreSQL数据库就绪..."
    
    local max_attempts=30
    local attempt=1
    
    while [ $attempt -le $max_attempts ]; do
        if docker exec goblog-postgres pg_isready -U goblog -d goblog > /dev/null 2>&1; then
            log_info "数据库已就绪"
            return 0
        fi
        
        echo -n "."
        sleep 2
        attempt=$((attempt + 1))
    done
    
    log_error "数据库启动超时"
    return 1
}

# 运行数据库迁移
run_migration() {
    log_step "运行数据库迁移..."
    
    # 等待应用服务启动
    sleep 5
    
    # 检查应用是否自动完成了迁移
    if docker logs goblog-backend 2>&1 | grep -q "Migration completed\|Schema created"; then
        log_info "数据库迁移已自动完成"
    else
        log_info "手动运行数据库迁移..."
        if command -v go &> /dev/null; then
            # 如果有Go环境，使用本地脚本
            # 检测脚本位置
            if [ -f "./scripts/migrate_postgres.sh" ]; then
                ./scripts/migrate_postgres.sh migrate
            elif [ -f "./migrate_postgres.sh" ]; then
                ./migrate_postgres.sh migrate
            else
                log_warn "找不到迁移脚本，跳过手动迁移"
            fi
        else
            # 否则检查容器是否已自动完成迁移
            log_info "等待容器自动完成数据库迁移..."
            sleep 10
        fi
    fi
}

# 测试API连接
test_api() {
    log_step "测试API连接..."
    
    local max_attempts=10
    local attempt=1
    
    while [ $attempt -le $max_attempts ]; do
        if curl -s http://localhost:8080/api/categories > /dev/null 2>&1; then
            log_info "API服务正常运行"
            return 0
        fi
        
        echo -n "."
        sleep 3
        attempt=$((attempt + 1))
    done
    
    log_warn "API服务可能还未完全启动，请稍后再试"
}

# 显示访问信息
show_access_info() {
    log_step "服务访问信息"
    
    echo ""
    echo "🎉 GoBlog PostgreSQL版本启动成功！"
    echo ""
    echo "📡 服务地址："
    echo "  - API服务:        http://localhost:8080"
    echo "  - API文档:        http://localhost:8080/api"
    echo "  - pgAdmin管理:    http://localhost:8081"
    echo ""
    echo "🔑 默认登录信息："
    echo "  - 博客管理员:     admin / admin123"
    echo "  - pgAdmin:       admin@goblog.com / admin123"
    echo ""
    echo "🗄️ 数据库连接信息："
    echo "  - 主机:          localhost"
    echo "  - 端口:          5432"
    echo "  - 数据库:        goblog"
    echo "  - 用户名:        goblog"
    echo "  - 密码:          goblog123"
    echo ""
    echo "📋 常用命令："
    echo "  - 查看日志:      docker-compose logs -f"
    echo "  - 停止服务:      docker-compose down"
    echo "  - 重启服务:      docker-compose restart"
    echo "  - 数据库备份:    make db-backup"
    echo ""
    echo "🧪 测试API："
    echo "  # 获取分类列表"
    echo "  curl http://localhost:8080/api/categories"
    echo ""
    echo "  # 登录获取token"
    echo "  curl -X POST http://localhost:8080/auth/login \\"
    echo "    -H 'Content-Type: application/json' \\"
    echo "    -d '{\"username\":\"admin\",\"password\":\"admin123\"}'"
    echo ""
}

# 主函数
main() {
    echo "🚀 GoBlog PostgreSQL 快速启动脚本"
    echo "=================================="
    echo ""
    
    check_dependencies
    stop_existing_services
    start_services
    wait_for_database
    run_migration
    test_api
    show_access_info
    
    log_info "启动完成！享受使用GoBlog吧！"
}

# 错误处理
trap 'log_error "脚本执行失败，请检查错误信息"; exit 1' ERR

# 运行主函数
main "$@" 