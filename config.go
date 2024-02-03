package socks5

type Config struct {
	Host string          // host
	Port int             //	port
	Auth []Authenticator // 认证接口的实现
}
