#!/bin/bash

# 测试脚本：生成表结构并插入数据库
# 使用方式: ./scripts/init_db.sh [options]
# 选项:
#   -h, --help     显示帮助信息
#   -v, --verbose  显示详细输出

set -euo pipefail  # 设置严格的错误处理

# 默认变量
VERBOSE=false
MYSQL_HOST="${MYSQL_HOST:-127.0.0.1}"
MYSQL_PORT="${MYSQL_PORT:-4000}"
MYSQL_USER="${MYSQL_USER:-root}"
MYSQL_PASS="${MYSQL_PASS:-test}"
MYSQL_DB="${MYSQL_DB:-test}"

# 显示帮助信息
show_help() {
    echo "测试脚本：生成表结构并插入数据库"
    echo
    echo "使用方式: $0 [options]"
    echo
    echo "选项:"
    echo "  -h, --help     显示帮助信息"
    echo "  -v, --verbose  显示详细输出"
    echo
    echo "环境变量:"
    echo "  MYSQL_HOST  MySQL主机地址 (默认: 127.0.0.1)"
    echo "  MYSQL_PORT  MySQL端口 (默认: 3306)"
    echo "  MYSQL_USER  MySQL用户名 (默认: root)"
    echo "  MYSQL_PASS  MySQL密码 (默认: test)"
    echo "  MYSQL_DB    MySQL数据库名 (默认: test)"
}

# 解析命令行参数
while [[ $# -gt 0 ]]; do
    case $1 in
        -h|--help)
            show_help
            exit 0
            ;;
        -v|--verbose)
            VERBOSE=true
            shift
            ;;
        *)
            echo "未知选项: $1"
            show_help
            exit 1
            ;;
    esac
done

# 输出函数
log() {
    echo "[$(date +%Y-%m-%d\ %H:%M:%S)] $1"
}

verbose() {
    if [[ "$VERBOSE" = true ]]; then
        log "VERBOSE: $1"
    fi
}

# 构建DSN供Go程序使用
build_dsn() {
    echo "${MYSQL_USER}:${MYSQL_PASS}@tcp(${MYSQL_HOST}:${MYSQL_PORT})/${MYSQL_DB}?charset=utf8mb4&parseTime=True&loc=Local"
}

# 构建MySQL命令行参数
mysql_args() {
    echo "-h" "$MYSQL_HOST" "-P" "$MYSQL_PORT" "-u" "$MYSQL_USER" "-p$MYSQL_PASS"
}

# 检查MySQL连接
check_mysql_connection() {
    log "检查MySQL连接..."
    if ! mysql $(mysql_args) -e "SELECT 1" "$MYSQL_DB" >/dev/null 2>&1; then
        log "无法连接到MySQL数据库，请确保MySQL服务正在运行且配置正确"
        exit 1
    fi
    log "MySQL连接成功"
}

# 执行SQL语句
execute_sql() {
    local sql="$1"
    verbose "执行SQL: $sql"
    if ! mysql $(mysql_args) -e "$sql" "$MYSQL_DB" >/dev/null 2>&1; then
        log "SQL执行失败: $sql"
        exit 1
    fi
}

# 执行SQL文件
execute_sql_file() {
    local file="$1"
    verbose "执行SQL文件: $file"
    if ! mysql $(mysql_args) "$MYSQL_DB" < "$file" >/dev/null 2>&1; then
        log "SQL文件执行失败: $file"
        exit 1
    fi
}

# 运行Go程序
run_go_program() {
    local program="$1"
    local description="$2"
    
    log "运行${description}..."
    # 导出DSN供Go程序使用
    export MYSQL_DSN=$(build_dsn)
    if ! go run "$program"; then
        log "${description}运行失败"
        exit 1
    fi
    log "${description}运行成功"
}

# 主程序流程
main() {
    log "开始测试流程..."
    
    # 检查MySQL连接
    check_mysql_connection
    
    # 创建数据库（如果不存在）
    log "创建数据库..."
    execute_sql "CREATE DATABASE IF NOT EXISTS $MYSQL_DB"
    
    # 创建表结构
    log "创建表结构..."
    execute_sql_file "scripts/schema.sql"
    
    # 运行Go程序插入表结构数据
    run_go_program "cmd/populate_tables/main.go" "populate_tables程序"
    
    log "所有测试完成！"
}

# 执行主程序
main "$@"