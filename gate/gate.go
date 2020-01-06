package gate

import (
	"github.com/name5566/leaf/chanrpc"
	"github.com/name5566/leaf/log"
	"gitee.com/webkit/vuitl/network"
	"net"
	"reflect"
	"time"
)

type Gate struct {
	MaxConnNum      int
	PendingWriteNum int
	MaxMsgLen       uint32
	Processor       network.Processor
	AgentChanRPC    *chanrpc.Server

	// websocket
	WSAddr      string
	HTTPTimeout time.Duration
	CertFile    string
	KeyFile     string

	// tcp
	TCPAddr      string
	LenMsgLen    int
	LittleEndian bool

	Options *Options
}

type Options struct {
	PingTimeOut int
	PongTimeOut int
	AuthTimeOut time.Duration
}
func (gate *Gate) Run(closeSig chan bool) {
	var wsServer *network.WSServer

	if gate.Options.PingTimeOut <= 0 {
		gate.Options.PingTimeOut = 10
	} else {
		gate.Options.PingTimeOut = gate.Options.PingTimeOut
	}
	if gate.Options.PongTimeOut >= gate.Options.PingTimeOut - 2 || gate.Options.PongTimeOut <= 0{
		if gate.Options.PingTimeOut - 2 <= 0 {
			gate.Options.PongTimeOut =gate.Options.PingTimeOut
			gate.Options.PingTimeOut += 2
		}
	} else {
		gate.Options.PingTimeOut = gate.Options.PongTimeOut
	}
	if gate.WSAddr != "" {
		wsServer = new(network.WSServer)
		wsServer.Addr = gate.WSAddr
		wsServer.MaxConnNum = gate.MaxConnNum
		wsServer.PendingWriteNum = gate.PendingWriteNum
		wsServer.MaxMsgLen = gate.MaxMsgLen
		wsServer.HTTPTimeout = gate.HTTPTimeout
		wsServer.CertFile = gate.CertFile
		wsServer.KeyFile = gate.KeyFile
		wsServer.NewAgent = func(conn *network.WSConn) network.Agent {
			a := &agent{conn: conn, gate: gate}
			a.SetPongHandler(gate.Options)
			a.SetAuth() // 设置验证
			if gate.AgentChanRPC != nil {
				gate.AgentChanRPC.Go("NewAgent", a)
			}
			return a
		}
	}

	var tcpServer *network.TCPServer
	if gate.TCPAddr != "" {
		tcpServer = new(network.TCPServer)
		tcpServer.Addr = gate.TCPAddr
		tcpServer.MaxConnNum = gate.MaxConnNum
		tcpServer.PendingWriteNum = gate.PendingWriteNum
		tcpServer.LenMsgLen = gate.LenMsgLen
		tcpServer.MaxMsgLen = gate.MaxMsgLen
		tcpServer.LittleEndian = gate.LittleEndian
		tcpServer.NewAgent = func(conn *network.TCPConn) network.Agent {
			a := &agent{conn: conn, gate: gate}
			a.SetPongHandler(gate.Options)
			a.SetAuth() // 设置验证
			if gate.AgentChanRPC != nil {
				gate.AgentChanRPC.Go("NewAgent", a)
			}
			return a
		}
	}

	if wsServer != nil {
		wsServer.Start()
	}
	if tcpServer != nil {
		tcpServer.Start()
	}
	<-closeSig
	if wsServer != nil {
		wsServer.Close()
	}
	if tcpServer != nil {
		tcpServer.Close()
	}
}

func (gate *Gate) OnDestroy() {}

type agent struct {
	conn     network.Conn
	gate     *Gate
	userData interface{}
	pongTime int
	isAuth bool // 是否通过验证
	options *Options
}

func (a *agent) Run() {
	for {
		data, err := a.conn.ReadMsg()
		if err != nil {
			log.Debug("read message: %v", err)
			break
		}

		if a.gate.Processor != nil {
			msg, err := a.gate.Processor.Unmarshal(data)
			if err != nil {
				log.Debug("unmarshal message error: %v", err)
				break
			}
			err = a.gate.Processor.Route(msg, a)
			if err != nil {
				log.Debug("route message error: %v", err)
				break
			}
		}
	}
}

func (a *agent) OnClose() {
	if a.gate.AgentChanRPC != nil {
		err := a.gate.AgentChanRPC.Call0("CloseAgent", a)
		if err != nil {
			log.Error("chanrpc error: %v", err)
		}
	}
}

func (a *agent) WriteMsg(msg interface{}) {
	if a.gate.Processor != nil {
		data, err := a.gate.Processor.Marshal(msg)
		if err != nil {
			log.Error("marshal message %v error: %v", reflect.TypeOf(msg), err)
			return
		}
		err = a.conn.WriteMsg(data...)
		if err != nil {
			log.Error("write message %v error: %v", reflect.TypeOf(msg), err)
		}
	}
}

func (a *agent) LocalAddr() net.Addr {
	return a.conn.LocalAddr()
}

func (a *agent) RemoteAddr() net.Addr {
	return a.conn.RemoteAddr()
}

func (a *agent) Close() {
	a.conn.Close()
}

func (a *agent) Destroy() {
	a.conn.Destroy()
}

func (a *agent) UserData() interface{} {
	return a.userData
}

func (a *agent) SetUserData(data interface{}) {
	a.userData = data
}

func (a *agent) Auth(bool2 bool) {
	if bool2 {
		a.isAuth = true
	} else {
		a.isAuth = false
	}
}
// 设置心跳
func (a *agent) SetPongHandler(options *Options) {
	a.options = options
	// 发送参数
	a.WriteMsg(&Options{
		PingTimeOut: a.options.PingTimeOut,
		PongTimeOut: a.options.PongTimeOut,
		AuthTimeOut: a.options.AuthTimeOut,
	})
}
// 设置验证
func (a *agent) SetAuth() {


	// 开启验证
	if a.options.AuthTimeOut > 0 {
		go func(agent *agent) {
			time.Sleep(time.Second * a.options.AuthTimeOut)
			if !a.isAuth {a.Close()} // 验证不通过则关闭连接
		}(a)
	} else {

	}
}
