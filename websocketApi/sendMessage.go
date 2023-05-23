package websocketApi

import "go-fast/open.go/net"

type SendMessage struct {
}

func (s SendMessage) Process(data net.ActionModel) interface{} {
	return ""
}
