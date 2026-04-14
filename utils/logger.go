package utils

import (
	"fmt"
	"log"
	"os"
	"time"
)

// Logger 日志记录器
type Logger struct {
	infoLogger  *log.Logger
	errorLogger *log.Logger
}

// 全局日志实例
var logger *Logger

// InitLogger 初始化日志记录器
func InitLogger() {
	logger = &Logger{
		infoLogger:  log.New(os.Stdout, "", 0),
		errorLogger: log.New(os.Stderr, "", 0),
	}
}

// Info 记录信息级日志
func Info(message string) {
	logger.infoLogger.Println(formatLog("INFO", message))
}

// Error 记录错误级日志
func Error(message string) {
	logger.errorLogger.Println(formatLog("ERROR", message))
}

// formatLog 格式化日志输出
func formatLog(level, message string) string {
	return fmt.Sprintf("[%s] [%s] %s", time.Now().Format("2006-01-02 15:04:05"), level, message)
}
