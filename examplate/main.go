package main

import "github.com/xaitx/socks5"

func main() {
	config := &socks5.Config{
		Host: "127.0.0.1",
		Port: 10800,
		Auth: []socks5.Authenticator{
			socks5.NoAuthenticator{},
		},
	}
	socks5.StartServer(config)
}
