package gate

import (
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
	UpdateAuth (bool2 bool)
	GetAuth() interface{}
	setPingTime(time int64)
	getPingTime() int64
	getOptions() *Options
	setAuth(auth *Auth)
}
