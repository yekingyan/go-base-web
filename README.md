#

## 一、Protobuf安装

### 下载并解压

<https://github.com/protocolbuffers/protobuf/releases>

### 环境变量

{protobuf}/bin目录 设置进环境变量

### 查看

>protoc --version

## 二、生成Protobuf代码

server>./genProto

## 三、 vscode 启动服务配置

>mkdir -p .vscode
>cp config/vscodeLaunch.json .vscode/launch.json

## 四 密钥配置

server> touch auth/private.key
server> touch share/auth/public.key

并填入密钥
