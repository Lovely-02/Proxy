package main

import (
	"flag"
	"fmt"
	"os"
	"proxy/config"
	"proxy/proxy"
	"proxy/utils"
)

func main() {
	// 解析命令行参数
	target := flag.String("target", "", "目标服务器地址")
	listen := flag.String("listen", "", "监听地址")
	flag.Parse()

	// 加载配置
	cfg := config.LoadConfig(*target, *listen)

	// 初始化日志
	utils.InitLogger()

	// 启动代理服务器
	server := proxy.NewProxyServer(cfg)
	utils.Info(fmt.Sprintf("代理服务器启动在 %s，转发到 %s", cfg.Listen, cfg.Target))

	if err := server.Start(); err != nil {
		utils.Error(fmt.Sprintf("启动代理服务器失败: %v", err))
		os.Exit(1)
	}
}
