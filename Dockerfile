# 多阶段构建 - 构建阶段
FROM golang:1.24-alpine AS builder

# 设置工作目录
WORKDIR /app

# 安装必要的系统依赖
RUN apk add --no-cache git ca-certificates tzdata && \
    update-ca-certificates

# 复制Go模块文件
COPY go.mod go.sum ./
ENV GOPROXY https://goproxy.cn
# 下载依赖（利用Docker缓存层）
RUN go mod download && go mod verify

# 复制源代码
COPY . .

# 构建应用程序
# 使用-ldflags="-w -s"去除调试信息，减小二进制文件大小
# 使用-a强制重新构建所有包
RUN GOOS=linux GOARCH=amd64 go build \
    -ldflags="-w -s -X main.version=$(date -u +%Y%m%d%H%M%S)" \
    -o goblog \
    ./cmd/server/main.go

# 多阶段构建 - 运行阶段
FROM scratch

# 从builder阶段复制必要文件
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo
COPY --from=builder /app/goblog /app/goblog

# 设置环境变量
ENV TZ=Asia/Shanghai \
    GIN_MODE=release \
    PORT=8080

# 暴露端口
EXPOSE 8080

# 设置用户（安全性）
# 由于使用scratch镜像，我们需要使用nobody用户的UID/GID
USER 65534:65534

# 健康检查
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD ["/app/goblog", "--health-check"] || exit 1

# 入口点
ENTRYPOINT ["/app/goblog"]
