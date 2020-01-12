package network

import "github.com/threewe/cool/network/json"

type Processor interface {
	// must goroutine safe
	Route(msg interface{}, userData interface{}) error
	// must goroutine safe
	Unmarshal(data []byte) (interface{}, error)
	// must goroutine safe
	Marshal(msg interface{}) ([][]byte, error)
	Register(msg interface{}) string
	SetHandler(msg interface{}, msgHandler json.MsgHandler)
}
