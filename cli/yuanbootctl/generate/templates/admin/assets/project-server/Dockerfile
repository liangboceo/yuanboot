
# 第一阶段：构建阶段
# 使用官方的 Go 基础镜像，用于编译 Go 代码
FROM golang:1.26-alpine AS builder

# 设置工作目录
WORKDIR /app
# 更换 Go 模块代理为国内常用代理
ENV GOPROXY=https://goproxy.cn,direct

# 将当前目录下的所有文件复制到工作目录
COPY . .
ARG ENVIROMENT=fat
RUN echo $ENVIROMENT
RUN chmod 777 /app/build.sh \
    && cd /app\
    && sh build.sh 1.0.0-SNAPSHOT linux amd64 ${ENVIROMENT}

FROM alpine:latest

# 安装 Go 运行环境（用于技能包的 go build 编译）
# gcc + musl-dev 用于 cgo 支持
# chromium 用于 HTML 导出 PDF（headless 模式）
# font-noto-cjk 用于 PDF 中文渲染，font-noto-emoji 用于报告图标渲染
RUN apk add --no-cache go gcc musl-dev chromium font-noto-cjk font-noto-emoji \
    && fc-cache -f

# 配置 Go 环境变量
ENV GOPROXY=https://goproxy.cn,direct
ENV GOSUMDB=off
ENV GOPATH=/root/go
ENV PATH=$PATH:/root/go/bin

# 设置工作目录
WORKDIR /app

# 从构建阶段的镜像中复制编译好的可执行文件到当前镜像
COPY --from=builder /app/build/sendex-server-linux-amd64 .

# 定义容器启动时执行的命令
CMD ["./sendex-server-1.0.0-SNAPSHOT"]

