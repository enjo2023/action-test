#!/bin/bash

# 拉取最新代码
echo "拉取最新代码..."
git pull origin main

# 停止并移除旧容器
echo "停止并移除旧容器..."
docker-compose down

# 重新构建并启动容器
echo "重新构建并启动容器..."
docker-compose up -d --build

echo "部署完成！应用运行在 http://localhost:8080"
