#!/bin/bash

# PostgreSQL 数据库迁移和测试脚本
# 用于从SQLite迁移到PostgreSQL或设置新的PostgreSQL环境

set -e

# 颜色输出
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# 配置变量
POSTGRES_HOST=${POSTGRES_HOST:-localhost}
POSTGRES_PORT=${POSTGRES_PORT:-5432}
POSTGRES_DB=${POSTGRES_DB:-goblog}
POSTGRES_USER=${POSTGRES_USER:-goblog}
POSTGRES_PASSWORD=${POSTGRES_PASSWORD:-goblog123}

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

# 检查PostgreSQL连接
check_postgres_connection() {
    log_info "检查PostgreSQL连接..."
    
    if ! command -v psql &> /dev/null; then
        log_error "psql客户端未安装，请先安装PostgreSQL客户端"
        return 1
    fi
    
    export PGPASSWORD=$POSTGRES_PASSWORD
    if psql -h $POSTGRES_HOST -p $POSTGRES_PORT -U $POSTGRES_USER -d $POSTGRES_DB -c "SELECT 1;" > /dev/null 2>&1; then
        log_info "PostgreSQL连接成功"
        return 0
    else
        log_error "无法连接到PostgreSQL数据库"
        return 1
    fi
}

# 创建数据库（如果不存在）
create_database() {
    log_info "检查数据库是否存在..."
    
    export PGPASSWORD=$POSTGRES_PASSWORD
    if psql -h $POSTGRES_HOST -p $POSTGRES_PORT -U $POSTGRES_USER -lqt | cut -d \| -f 1 | grep -qw $POSTGRES_DB; then
        log_info "数据库 $POSTGRES_DB 已存在"
    else
        log_info "创建数据库 $POSTGRES_DB..."
        createdb -h $POSTGRES_HOST -p $POSTGRES_PORT -U $POSTGRES_USER $POSTGRES_DB
        log_info "数据库创建成功"
    fi
}

# 运行应用程序迁移
run_app_migration() {
    log_info "运行应用程序数据库迁移..."
    
    # 设置环境变量
    export DB_DRIVER=postgres
    export DB_DSN="host=$POSTGRES_HOST port=$POSTGRES_PORT user=$POSTGRES_USER password=$POSTGRES_PASSWORD dbname=$POSTGRES_DB sslmode=disable"
    
    # 构建并运行迁移
    if go run ./cmd/server/main.go --migrate-only 2>/dev/null; then
        log_info "应用程序迁移成功"
    else
        log_warn "应用程序迁移失败，将尝试手动启动应用让Ent自动创建表"
        # 创建一个临时的迁移程序
        cat > /tmp/migrate.go << 'EOF'
package main

import (
    "context"
    "log"
    "os"
    
    "goblog/ent"
    _ "github.com/lib/pq"
)

func main() {
    dsn := os.Getenv("DB_DSN")
    if dsn == "" {
        log.Fatal("DB_DSN environment variable is required")
    }
    
    client, err := ent.Open("postgres", dsn)
    if err != nil {
        log.Fatalf("failed opening connection to postgres: %v", err)
    }
    defer client.Close()
    
    // Run the auto migration tool.
    if err := client.Schema.Create(context.Background()); err != nil {
        log.Fatalf("failed creating schema resources: %v", err)
    }
    
    log.Println("Migration completed successfully")
}
EOF
        
        if go run /tmp/migrate.go; then
            log_info "数据库schema创建成功"
            rm -f /tmp/migrate.go
        else
            log_error "数据库迁移失败"
            rm -f /tmp/migrate.go
            return 1
        fi
    fi
}

# 测试数据库连接和基本操作
test_database() {
    log_info "测试数据库连接和基本操作..."
    
    export PGPASSWORD=$POSTGRES_PASSWORD
    
    # 检查表是否存在
    if psql -h $POSTGRES_HOST -p $POSTGRES_PORT -U $POSTGRES_USER -d $POSTGRES_DB -c "\dt" | grep -q "articles\|categories\|tags"; then
        log_info "数据库表存在，检查表结构..."
        psql -h $POSTGRES_HOST -p $POSTGRES_PORT -U $POSTGRES_USER -d $POSTGRES_DB -c "\d articles" > /dev/null
        psql -h $POSTGRES_HOST -p $POSTGRES_PORT -U $POSTGRES_USER -d $POSTGRES_DB -c "\d categories" > /dev/null
        psql -h $POSTGRES_HOST -p $POSTGRES_PORT -U $POSTGRES_USER -d $POSTGRES_DB -c "\d tags" > /dev/null
        log_info "所有表结构正常"
    else
        log_warn "数据库表不存在或不完整"
        return 1
    fi
}

# 显示数据库信息
show_db_info() {
    log_info "数据库连接信息："
    echo "  主机: $POSTGRES_HOST"
    echo "  端口: $POSTGRES_PORT"
    echo "  数据库: $POSTGRES_DB"
    echo "  用户: $POSTGRES_USER"
    echo ""
    
    export PGPASSWORD=$POSTGRES_PASSWORD
    log_info "数据库统计信息："
    psql -h $POSTGRES_HOST -p $POSTGRES_PORT -U $POSTGRES_USER -d $POSTGRES_DB -c "
        SELECT 
            schemaname,
            tablename,
            attname as column_name,
            typname as data_type
        FROM pg_tables 
        JOIN pg_attribute ON pg_tables.tablename = pg_attribute.attrelid::regclass::text
        JOIN pg_type ON pg_attribute.atttypid = pg_type.oid
        WHERE schemaname = 'public' 
        AND attnum > 0 
        AND NOT attisdropped
        ORDER BY tablename, attnum;
    " 2>/dev/null || echo "无法获取详细表信息"
}

# 主函数
main() {
    log_info "开始PostgreSQL数据库设置和迁移..."
    
    case "${1:-migrate}" in
        "check")
            check_postgres_connection
            ;;
        "create")
            create_database
            ;;
        "migrate")
            check_postgres_connection && \
            create_database && \
            run_app_migration && \
            test_database && \
            show_db_info
            ;;
        "test")
            check_postgres_connection && test_database
            ;;
        "info")
            show_db_info
            ;;
        *)
            echo "用法: $0 {check|create|migrate|test|info}"
            echo ""
            echo "  check   - 检查PostgreSQL连接"
            echo "  create  - 创建数据库"
            echo "  migrate - 完整迁移过程（推荐）"
            echo "  test    - 测试数据库连接和表"
            echo "  info    - 显示数据库信息"
            exit 1
            ;;
    esac
}

# 运行主函数
main "$@" 