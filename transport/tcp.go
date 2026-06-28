package transport

import (
	"crypto/tls"
	"fmt"
	"net"

	"golang.org/x/net/proxy"
)

type TcpTransport struct {
	host   string
	port   int
	prx    string
	useSSL bool
	conn   net.Conn
}

func NewTcpTransport(host string, port int, proxyStr string, useSSL bool) *TcpTransport {
	return &TcpTransport{host: host, port: port, prx: proxyStr, useSSL: useSSL}
}

func (t *TcpTransport) Connect() error {
	addr := net.JoinHostPort(t.host, fmt.Sprintf("%d", t.port))
	var conn net.Conn
	var err error

	if t.prx != "" {
		dialer, err := proxy.SOCKS5("tcp", t.prx, nil, proxy.Direct)
		if err != nil {
			return err
		}
		conn, err = dialer.Dial("tcp", addr)
	} else {
		conn, err = net.Dial("tcp", addr)
	}
	if err != nil {
		return err
	}

	if t.useSSL {
		tlsConn := tls.Client(conn, &tls.Config{ServerName: t.host})
		if err := tlsConn.Handshake(); err != nil {
			conn.Close()
			return err
		}
		conn = tlsConn
	}

	t.conn = conn
	return nil
}

func (t *TcpTransport) Close() error {
	if t.conn != nil {
		return t.conn.Close()
	}
	return nil
}

func (t *TcpTransport) Send(data []byte) error {
	if t.conn == nil {
		return fmt.Errorf("tcp: not connected")
	}
	_, err := t.conn.Write(data)
	return err
}

func (t *TcpTransport) Recv() ([]byte, error) {
	if t.conn == nil {
		return nil, fmt.Errorf("tcp: not connected")
	}
	buf := make([]byte, 1024)
	n, err := t.conn.Read(buf)
	if err != nil {
		return nil, err
	}
	return buf[:n], nil
}

func (t *TcpTransport) Connected() bool {
	return t.conn != nil
}
