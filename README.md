# Proxy

Proxy是一个轻量级的HTTP代理服务，用于处理请求转发，支持HTTP和HTTPS目标服务器。

## 功能特点

- 支持HTTP和HTTPS目标服务器
- 自动处理重定向响应的Location头
- 支持WebSocket连接转发
- 简洁的配置方式
- 详细的日志记录

## 安装方法

### 前提条件

- Go 1.20或更高版本

### 编译

```bash
go build -o Proxy.exe main.go
```

## 配置说明

### 配置文件

创建`.env`文件，配置以下参数：

```
target=https://enka.network  # 目标服务器地址（包含协议）
listen=0.0.0.0:7860            # 监听地址
```

### 命令行参数

也可以通过命令行参数覆盖配置文件中的设置：

```bash
./Proxy -target=https://enka.network -listen=0.0.0.0:7860
```

## 使用方法

### 启动代理服务器

```bash
./Proxy
```

### 测试代理

使用浏览器或curl等工具访问代理服务器：

```bash
curl http://127.0.0.1:7860
```

## 日志说明

代理服务器会输出以下类型的日志：

- INFO: 记录请求信息和服务器启动状态
- ERROR: 记录代理过程中的错误

## 代码结构

```
Proxy/
├── config/           # 配置加载
├── Proxy/            # 代理核心实现
├── utils/            # 工具函数
├── .env              # 配置文件
├── README.md         # 说明文档
├── go.mod            # Go模块文件
└── main.go           # 主入口
```
