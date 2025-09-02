# 构建阶段
FROM golang:1.24-alpine AS builder

# 安装必要的工具和证书
RUN apk add --no-cache git ca-certificates tzdata

# 设置工作目录
WORKDIR /app

# 设置 Go 环境变量以优化构建
ENV CGO_ENABLED=0
ENV GOOS=linux
ENV GOARCH=amd64

# 复制 go mod 文件并下载依赖（利用 Docker 层缓存）
COPY go.mod go.sum ./
RUN go mod download && go mod verify

# 复制源代码
COPY . .

# 构建应用，使用优化参数
RUN go build -a -installsuffix cgo \
    -ldflags '-w -s -extldflags "-static"' \
    -o main .

# 验证构建的二进制文件
RUN ls -la main && file main

# 运行阶段
FROM alpine:latest

# 安装运行时依赖
RUN apk --no-cache add ca-certificates tzdata curl

# 创建非 root 用户
RUN addgroup -g 1001 -S appgroup && \
    adduser -u 1001 -S appuser -G appgroup

# 设置工作目录
WORKDIR /app

# 复制可执行文件和配置
COPY --from=builder --chown=appuser:appgroup /app/main .
COPY --from=builder --chown=appuser:appgroup /app/config ./config
COPY --from=builder --chown=appuser:appgroup /app/assets ./assets

# 设置可执行权限
RUN chmod +x ./main

# 切换到非 root 用户
USER appuser

# 暴露端口
EXPOSE 8080

# 健康检查
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD curl -f http://localhost:8080/health || exit 1

# 启动应用
CMD ["./main"]