package gate

import (
	"gitee.com/webkit/vuitl/network/json"
	"net"
)

type Agent interface {
	WriteMsg(msg interface{})
	LocalAddr() net.Addr
	RemoteAddr() net.Addr
	Close()
	Destroy()
	UserData() interface{}
	SetUserData(data interface{})
	SetPong(ping *json.Ping)
	Auth (bool2 bool)
	SetPingTime(time int64)
	GetPingTime() int64
	GetOptions() *Options
}
