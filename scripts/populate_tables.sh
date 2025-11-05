#!/bin/bash

# 脚本用于生成15个关联表并将其写入数据库

echo "开始生成表结构..."

# 运行Go程序插入表结构数据
echo "运行populate_tables程序..."
go run cmd/populate_tables/main.go

if [ $? -eq 0 ]; then
    echo "表结构数据插入成功"
else
    echo "表结构数据插入失败"
    exit 1
fi

echo "表结构生成完成！"