#!/usr/bin/env pwsh

# 拉取最新代码
Write-Output "拉取最新代码..."
git pull origin main

# 停止并移除旧容器
Write-Output "停止并移除旧容器..."
docker-compose down

# 重新构建并启动容器
Write-Output "重新构建并启动容器..."
docker-compose up -d --build

Write-Output "部署完成！应用运行在 http://localhost:8080"
