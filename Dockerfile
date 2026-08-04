# 构建阶段
FROM golang:1.21-alpine AS builder

WORKDIR /build

# 安装构建依赖
RUN apk add --no-cache git gcc musl-dev

# 复制 go mod 文件
COPY go.mod go.sum ./

# 下载依赖
RUN go mod download

# 复制源代码
COPY . .

# 编译
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags="-s -w" -o iscsi-web-panel main.go

# 运行阶段
FROM alpine:latest

# 安装运行时依赖
RUN apk add --no-cache \
    tgt \
    sqlite \
    bash \
    curl \
    && mkdir -p /app/data \
    && mkdir -p /sys/kernel/config

WORKDIR /app

# 从构建阶段复制二进制文件
COPY --from=builder /build/iscsi-web-panel /app/iscsi-web-panel
COPY --from=builder /build/frontend /app/frontend

# 设置环境变量
ENV LISTEN_ADDR=:3005
ENV DATA_DIR=/app/data
ENV DB_PATH=/app/data/iscsi.db

# 暴露端口
EXPOSE 3005

# 健康检查
HEALTHCHECK --interval=30s --timeout=10s --start-period=5s --retries=3 \
    CMD curl -f http://localhost:3005/api/v1/api-doc || exit 1

# 启动命令
CMD ["/app/iscsi-web-panel"]
