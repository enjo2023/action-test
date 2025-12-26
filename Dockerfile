# 第一阶段：构建阶段
FROM golang:1.24-alpine AS builder

# 安装构建依赖和sqlite3开发库
RUN apk add --no-cache build-base sqlite-dev

# 设置工作目录
WORKDIR /app

# 复制go.mod和go.sum文件
COPY go.mod go.sum ./

# 下载依赖
RUN go mod download

# 复制项目文件
COPY . .

# 编译项目，启用CGO
RUN CGO_ENABLED=1 go build -o drink_water_helper .

# 第二阶段：运行阶段
FROM alpine:latest

# 安装sqlite3运行时依赖
RUN apk add --no-cache sqlite-libs

# 设置工作目录
WORKDIR /app

# 复制编译后的可执行文件
COPY --from=builder /app/drink_water_helper .

# 复制模板文件
COPY --from=builder /app/templates ./templates

# 创建静态文件目录
RUN mkdir -p ./static

# 暴露端口
EXPOSE 8080

# 运行应用
CMD ["./drink_water_helper"]
