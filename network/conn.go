package network

import (
	"net"
)
type Options struct {
	PingTimeOut int
	PongTimeOut int
	AuthTimeOut int
}
type Conn interface {
	ReadMsg() ([]byte, error)
	WriteMsg(args ...[]byte) error
	LocalAddr() net.Addr
	RemoteAddr() net.Addr
	Close()
	Destroy()
}
