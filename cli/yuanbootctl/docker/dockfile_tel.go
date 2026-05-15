package docker

const DockerFile_Tel = `
# 第一阶段：构建阶段
# 使用官方的 Go 基础镜像，用于编译 Go 代码
FROM golang:1.24-alpine AS builder

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
# 设置工作目录
WORKDIR /app

# 从构建阶段的镜像中复制编译好的可执行文件到当前镜像
COPY --from=builder /app/build/{{.ModelName}}-linux-amd64 .

# 定义容器启动时执行的命令
CMD ["./{{.ModelName}}-1.0.0-SNAPSHOT"]

`

const DockerCompose_Tel = `
version: '3'
services:
  iot-link-bridge-server: # 服务名称 建议与容器名称一致
    image: registry.cn-shanghai.aliyuncs.com/{{.ModelName}}:fat
    container_name: {{.ModelName}} # 容器名称
    environment:
      - TZ=Asia/Shanghai # 设置容器时区 我这里通过下面挂载方式同步的宿主机时区和时间了,这里忽略
      - YUANBOOT_PROFILE=fat
    volumes:
      - /etc/hosts:/etc/hosts:ro
      - /mnt/data/log:/mnt/data/log
    network_mode: host
    restart: always # 容器随docker启动自启
`

const BuildSh_Tel = `
#!/bin/sh
out_file="{{.ModelName}}"
VersionPath="{{.ModelName}}/version"
[ $# -lt 4 ] && {
	echo "Usage: $0 1.0.0-SNAPSHOT linux amd64 dev"
	exit 1
}
build() {
  local version="$1"
  local os="$2"
  local arch="$3"
  local profile="$4"
  local dir="build/$out_file-$os-$arch"
  out_file="${out_file}-${version}"
  go env -w GOSUMDB=off
  [ "$os" = "windows" ] && {
  		out_file="${out_file}.exe"
  	}
  rm -rf $dir
  mkdir -p $dir
  GOOS=$os GOARCH=$arch CGO_ENABLED=0 go build -ldflags="-s -w -X $VersionPath.env=$profile " -o "${dir}/${out_file}" .

}

main() {
  echo "mod download"
  go get -t .
  go mod download
  go get $out_file
  build $1 $2 $3 $4
}

main $1 $2 $3 $4
`
