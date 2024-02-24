package socks5

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"net"
	"sync"
)

// SOCKS5版本号
const socks5Version = 0x05

// 命令类型
const (
	cmdConnect = 0x01
	cmdBind    = 0x02
	cmdUDP     = 0x03
)

type server struct {
	config *Config
}

func StartServer(cfg *Config) error {
	// 设置config
	s := &server{config: cfg}

	// 初始化日志
	initLog(cfg)

	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", s.config.Host, s.config.Port))
	if err != nil {
		logger.Error(err.Error())
		return fmt.Errorf("failed to start listening: %w", err)
	}
	defer listener.Close()

	logger.Info("SOCKS5 server listening on %s:%d", s.config.Host, s.config.Port)

	for {
		conn, err := listener.Accept()
		if err != nil {
			logger.Error(err.Error())
			return fmt.Errorf("failed to accept connection: %w", err)
		}
		logger.Info("Accepted connection from %s", conn.RemoteAddr())
		go s.handleSocks5Conn(conn)
	}
}

func (s *server) handleSocks5Conn(clientConn net.Conn) {
	// SOCKS5判断协议并读取方法
	authMethods, err := readMethods(clientConn)
	if err != nil {
		clientConn.Close()
		return
	}
	// 遍历config的auth,通过GetCode返回值判断是否在authMethods中，选择执行对用的Authenticate方法
	for _, auth := range s.config.Auth {
		if bytes.Contains(authMethods, []byte{auth.GetCode()}) {
			// 认证
			if ok, err := auth.Authenticate(clientConn); err != nil {
				logger.Info("Authenticate error: %v", err)
				// 认证出现错误
				clientConn.Close()
				return
			} else if !ok {
				// 认证失败
				clientConn.Close()
				return
			}
			break
		}
	}
	// 调用request函数实现请求阶段
	var destConn net.Conn
	if destConn, err = s.request(clientConn); err != nil {
		logger.Info("Request error: %v", err)
		clientConn.Close()
		return
	}
	// 开始处理流量转发
	copyData(destConn, clientConn)
}

// readMethods函数从net.Conn连接中读取SOCKS方法列表
func readMethods(conn net.Conn) ([]byte, error) {
	buf := make([]byte, 2)
	if _, err := io.ReadFull(conn, buf); err != nil {
		return nil, err
	}
	if buf[0] != socks5Version {
		return nil, errors.New("invalid SOCKS version")
	}
	nMethods := int(buf[1])
	methods := make([]byte, nMethods)
	if _, err := io.ReadFull(conn, methods); err != nil {
		return nil, err
	}
	return methods, nil
}

func writeMethodSelection(conn net.Conn, selectedMethod []byte) error {
	if _, err := conn.Write([]byte{socks5Version, selectedMethod[0]}); err != nil {
		return err
	}
	return nil
}

// 实现SOCKS5的请求阶段
// 先处理客户端发送过来的请求，根据CMD进行不同的处理，并连接到目标地址
// 然后向客户端发送一个响应
func (s *server) request(conn net.Conn) (net.Conn, error) {
	// 读取数据
	buf := make([]byte, 256)
	_, err := conn.Read(buf)
	if err != nil {
		return nil, err
	}
	// 解析客户端请求
	version := buf[0]  // SOCKS协议版本
	cmd := buf[1]      // 请求的命令类型
	addrType := buf[3] // 地址类型

	// 判断协议版本是否支持
	if version != 5 {
		// 发送不支持的协议版本响应给客户端
		response := []byte{5, 0xFF}
		conn.Write(response)
		return nil, errors.New("Unsupported SOCKS version")
	}
	// 判断请求的命令类型
	switch cmd {
	case 1: // CONNECT请求
		// 解析请求的目标地址
		var destAddr string
		switch addrType {
		case 1: // IPv4地址类型
			destAddr = net.IPv4(buf[4], buf[5], buf[6], buf[7]).String()
			port := binary.BigEndian.Uint16(buf[8:10])
			destAddr = fmt.Sprintf("%s:%d", destAddr, port)
		case 3: // 域名地址类型
			destAddrLen := int(buf[4])
			destAddr = string(buf[5 : 5+destAddrLen])
			port := binary.BigEndian.Uint16(buf[5+destAddrLen : 5+destAddrLen+2])
			destAddr = fmt.Sprintf("%s:%d", destAddr, port)
		default:
			// 不支持的地址类型
			response := []byte{5, 0x08}
			conn.Write(response)
			return nil, errors.New("Unsupported address type")
		}

		// 连接到目标地址
		destConn, err := net.Dial("tcp", destAddr)
		if err != nil {
			// 连接失败，发送连接失败响应给客户端
			response := []byte{5, 0x01, 0x00, 0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}
			conn.Write(response)
			return nil, err
		}

		// 连接成功，发送连接成功响应给客户端
		response := []byte{5, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}
		conn.Write(response)

		return destConn, nil
	default:
		// 不支持的命令类型
		response := []byte{5, 0x07}
		conn.Write(response)
		return nil, errors.New("Unsupported command")
	}
}

// 实现流的复制
func copyData(dst net.Conn, src net.Conn) {
	defer dst.Close()
	defer src.Close()

	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()
		io.Copy(dst, src)
	}()

	go func() {
		defer wg.Done()
		io.Copy(src, dst)
	}()

	wg.Wait()
}
