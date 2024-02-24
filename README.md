# SOCKS5

该项目是一个使用 Go 语言编写的 SOCKS5 服务器模块，它允许用户搭建一个支持 SOCKS5 协议的代理服务器，用于转发 TCP 流量。

## 功能特点

- 支持 SOCKS5 协议版本，可以与任何兼容 SOCKS5 协议的客户端进行通信。
- 支持多种认证方式，包括无需认证和用户名/密码认证，使用户能够根据需求选择合适的认证方式，也可以自己实现认证接口。
- 可以指定服务器监听的主机地址和端口号，方便用户根据实际情况进行配置。
- 支持日志记录功能，用户可以选择是否开启日志记录，并可以自定义日志输出流。

## 安装

使用 `go get` 命令安装该项目：

```bash
go get github.com/xaitx/socks5
```

## 使用用例
```go
package main

import "github.com/xaitx/socks5"

func main() {
	config := socks5.NewConfig("127.0.0.1", 1080)
	socks5.StartServer(config)
}

```

也可以通过结构体的方式配置：

```go
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

```


## 配置选项

- `Host`：监听的主机地址
- `Port`：监听的端口号
- `Auth`：认证方式，可以是无需认证 (`NoAuthenticator{}`) 或者用户名/密码认证 (`UsernamePasswordAuthenticator{}`), 默认是无密码认证，也可以实现 `Authenticator` 接口，自定义认证方式
- `LogEnabled`：是否开启日志记录
- `LogOutput`：日志输出流，可以是文件流或者标准输出流，默认输出到终端

## 贡献

- 欢迎提出问题、报告 bug 或者贡献代码！请提交 issue 或者 PR。