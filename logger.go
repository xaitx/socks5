package socks5

import (
	"log"
	"os"
)

type customLogger struct {
	logEnabled bool     // 是否开启日志
	logOutput  *os.File //日志输入流
	currentLog *log.Logger
}

var logger *customLogger = &customLogger{}

// Init 初始化日志模块，设置默认的日志输出
func initLog(cfg *Config) {
	logger.logEnabled = cfg.LogEnabled
	logger.logOutput = cfg.LogOutput
	logger.SetOutput()
}

// SetOutput 自定义日志输出位置（如：文件）
func (l *customLogger) SetOutput() {
	l.currentLog = log.New(l.logOutput, "", log.LstdFlags|log.LUTC)
}

// writeIfEnabled 是一个内部辅助函数，用于在日志启用时执行实际的写入操作
func (l *customLogger) writeIfEnabled(level string, format string, v ...interface{}) {
	if l.logEnabled {
		l.currentLog.Printf("[%s] "+format, append([]interface{}{level}, v...)...)
	}
}

// Info 打印信息日志
func (l *customLogger) Info(format string, v ...interface{}) {
	l.writeIfEnabled("INFO", format, v...)
}

// Warn 打印警告日志
func (l *customLogger) Warn(format string, v ...interface{}) {
	l.writeIfEnabled("WARN", format, v...)
}

// Error 打印错误日志
func (l *customLogger) Error(format string, v ...interface{}) {
	l.writeIfEnabled("ERROR", format, v...)
}
