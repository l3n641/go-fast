package websocketApi

import "go-fast/open.go/webscockt"

func init() {
	var sendMessageRouter = webscockt.Router{Action: SendMessage{}}
	var routers = map[string]*webscockt.Router{
		"sendMessage": &sendMessageRouter,
	}
	webscockt.WebSocketRouters.SetRouters(routers)
}
