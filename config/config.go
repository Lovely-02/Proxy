package config

import (
	"bufio"
	"os"
	"strings"
)

// Config 代理服务器配置
type Config struct {
	Target string
	Listen string
}

// LoadConfig 加载配置，优先使用命令行参数，其次使用环境变量
func LoadConfig(cmdTarget, cmdListen string) *Config {
	// 加载.env文件
	loadEnvFile()

	// 从命令行参数获取配置
	target := cmdTarget
	listen := cmdListen

	// 如果命令行参数未提供，则从环境变量获取
	if target == "" {
		target = os.Getenv("target")
	}

	if listen == "" {
		listen = os.Getenv("listen")
		if listen == "" {
			listen = "0.0.0.0:7860" // 默认监听地址
		}
	}

	return &Config{
		Target: target,
		Listen: listen,
	}
}

// loadEnvFile 加载.env文件
func loadEnvFile() {
	// 获取当前工作目录
	cwd, err := os.Getwd()
	if err != nil {
		return
	}

	// 打开.env文件
	file, err := os.Open(cwd + "/.env")
	if err != nil {
		// .env文件不存在，忽略错误
		return
	}
	defer file.Close()

	// 读取文件内容
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		// 跳过空行和注释
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		// 解析环境变量
		parts := strings.SplitN(line, "=", 2)
		if len(parts) == 2 {
			key := strings.TrimSpace(parts[0])
			value := strings.TrimSpace(parts[1])
			// 设置环境变量
			os.Setenv(key, value)
		}
	}

	// 检查是否有扫描错误
	if err := scanner.Err(); err != nil {
		return
	}
}
