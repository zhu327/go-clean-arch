# ---- Build Stage ----
FROM golang:1.21-alpine AS builder

WORKDIR /app

# 复制 go.mod 和 go.sum 并下载依赖，利用 Docker 的层缓存
COPY go.mod go.sum ./
RUN go mod download

# 复制所有源代码
COPY . .

# 执行 DI 和 build
# wire 生成的程序码也应该被 .dockerignore 排除，或者在这里生成
RUN go generate ./...
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o /app/server ./cmd/api/main.go


# ---- Final Stage ----
FROM alpine:latest

WORKDIR /app

# 从 builder 阶段复制编译好的二进制文件和配置
COPY --from=builder /app/server .
# 如果有配置文件，也一并复制
# COPY --from=builder /app/configs ./configs
# COPY .env .

# 暴露端口
EXPOSE 8000

# 运行二进制文件
CMD ["./server"]