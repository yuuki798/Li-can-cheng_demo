# 使用官方的golang镜像作为基础
FROM golang:1.20.6 AS builder
# 使用官方Go镜像作为基础
# 设置工作目录
WORKDIR /app

# 将Go模块和源代码复制到容器中
COPY . .

# 设置GO111MODULE环境变量，并构建Go应用程序
RUN go build -o main main.go

# 指定容器启动时要运行的命令
CMD ["/app/main"]
