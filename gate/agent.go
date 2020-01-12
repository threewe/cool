package gate

import (
	"net"
	"time"
)

type Agent interface {
	WriteMsg(msg interface{})
	LocalAddr() net.Addr
	RemoteAddr() net.Addr
	Close()
	Destroy()
	UserData() interface{}
	SetUserData(data interface{})
	Timer(d time.Duration, cb func())
	Auth (bool2 bool)
	setPingTime(time int64)
	getPingTime() int64
	getOptions() *Options
}
