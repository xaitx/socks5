package socks5

import (
	"os"
)

type Config struct {
	Host       string          // host
	Port       int             //	port
	Auth       []Authenticator // 认证接口的实现
	LogEnabled bool            // 是否开启日志
	LogOutput  *os.File        //日志输入流
}

// NewConfig 创建一个配置对象，并简单化操纵，设置部分默认值
// 返回一个Config的指针
func NewConfig(host string, port int) *Config {
	return &Config{
		Host:       host,
		Port:       port,
		Auth:       []Authenticator{NoAuthenticator{}},
		LogEnabled: true,
		LogOutput:  os.Stdout,
	}
}
