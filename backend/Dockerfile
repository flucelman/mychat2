# 构建阶段
FROM golang:1.24-alpine AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o main .

# 运行阶段
FROM alpine:latest

RUN apk --no-cache add ca-certificates tzdata
WORKDIR /root/

# 复制可执行文件和配置
COPY --from=builder /app/main .
COPY --from=builder /app/config ./config
COPY --from=builder /app/assets ./assets

EXPOSE 8080

CMD ["./main"]
