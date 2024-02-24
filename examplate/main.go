package main

import (
	"os"

	"github.com/xaitx/socks5"
)

func main() {
	// config := socks5.NewConfig("127.0.0.1", 1080)  //快速启动一个简单的例子
	config := &socks5.Config{
		Host:       "127.0.0.1",                                                                                          //监听的接口
		Port:       1080,                                                                                                 //监听的端口
		Auth:       []socks5.Authenticator{&socks5.UsernamePasswordAuthenticator{Username: "admin", Password: "123456"}}, //认证方式，也可以选择无密码认证，或者自己实现认证的接口，并传入
		LogEnabled: true,                                                                                                 // 是否开启日志
		LogOutput:  os.Stdout,                                                                                            //日志输出
	}
	socks5.StartServer(config)
}
