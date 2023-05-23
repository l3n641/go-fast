package parser

import (
	"go-fast/open.go/net"
	"sync"
)

type mgr struct {
	lock      sync.RWMutex
	webSocket map[string]net.Parser
}

func (a *mgr) SetWebSocketParsers(data map[string]net.Parser) {
	a.lock.Lock()
	defer a.lock.Unlock()
	a.webSocket = data
}

func (a *mgr) SetWebSocketParser(name string, parser net.Parser) {
	a.lock.Lock()
	defer a.lock.Unlock()
	a.webSocket[name] = parser
}

func (a *mgr) GetWebSocketParser(name string) (net.Parser, bool) {
	a.lock.RLock()
	defer a.lock.RUnlock()
	s, b := a.webSocket[name]
	return s, b
}

var Mgr = mgr{
	webSocket: make(map[string]net.Parser),
}
