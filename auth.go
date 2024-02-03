package socks5

import (
	"errors"
	"fmt"
	"io"
	"net"
)

// Authenticator 是SOCKS5协议中的认证接口
type Authenticator interface {
	//获取认证方法的Code值
	GetCode() byte
	// Authenticate 用于认证客户端的连接
	Authenticate(conn net.Conn) (bool, error)
}

// NoAuthenticator 实现了无需认证的方法
type NoAuthenticator struct{}

func (n NoAuthenticator) GetCode() byte {
	return 0x00
}

func (n NoAuthenticator) Authenticate(conn net.Conn) (bool, error) {
	// 发送支持无认证的消息
	if _, err := conn.Write([]byte{socks5Version, n.GetCode()}); err != nil {
		return false, fmt.Errorf("failed to send authentication methods: %w", err)
	}
	return true, nil
}

// UsernamePasswordAuthenticator 实现了用户名/密码认证的方法
type UsernamePasswordAuthenticator struct {
	Username string
	Password string
}

func (u UsernamePasswordAuthenticator) GetCode() byte {
	return 0x02
}

func (u UsernamePasswordAuthenticator) Authenticate(conn net.Conn) (bool, error) {
	// 发送支持用户名/密码认证的消息
	if _, err := conn.Write([]byte{socks5Version, 0x02}); err != nil {
		return false, fmt.Errorf("failed to send authentication methods: %w", err)
	}

	// 读取并验证客户端发送的用户名/密码
	usernameLength := make([]byte, 1)
	if _, err := io.ReadFull(conn, usernameLength); err != nil {
		return false, fmt.Errorf("failed to read username length: %w", err)
	}

	username := make([]byte, int(usernameLength[0]))
	if _, err := io.ReadFull(conn, username); err != nil {
		return false, fmt.Errorf("failed to read username: %w", err)
	}

	passwordLength := make([]byte, 1)
	if _, err := io.ReadFull(conn, passwordLength); err != nil {
		return false, fmt.Errorf("failed to read password length: %w", err)
	}

	password := make([]byte, int(passwordLength[0]))
	if _, err := io.ReadFull(conn, password); err != nil {
		return false, fmt.Errorf("failed to read password: %w", err)
	}

	// 这里仅做简单比较，实际应用中应加密存储和比对密码
	if u.Username == string(username) && u.Password == string(password) {
		// 发送成功响应
		successResponse := []byte{socks5Version, 0x00}
		if _, err := conn.Write(successResponse); err != nil {
			return false, fmt.Errorf("failed to send successful authentication response: %w", err)
		}
		return true, nil
	}

	// 发送失败响应
	failureResponse := []byte{socks5Version, 0x01}
	if _, err := conn.Write(failureResponse); err != nil {
		return false, fmt.Errorf("failed to send failed authentication response: %w", err)
	}
	return false, errors.New("invalid username or password")
}
