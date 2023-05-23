package webscockt

import (
	"fmt"
	"github.com/olahol/melody"
	"go-fast/open.go/net"
	"net/http"
	"sync"
)

func NewWebsocket(config *melody.Config) *WS {
	return &WS{melodyConfig: config, pools: pools{data: make(map[int64]*Conn)}}

}

type WS struct {
	instance     *melody.Melody
	melodyConfig *melody.Config
	once         sync.Once
	lock         sync.RWMutex
	pools        pools
	sno          int64 //连接编码

}

func (w *WS) getSno() int64 {
	w.lock.RLock()
	defer w.lock.RUnlock()
	w.sno++
	return w.sno
}

func (w *WS) GetWebsocket() *WS {
	w.once.Do(func() {
		instance := melody.New()
		if w.melodyConfig != nil {
			instance.Config = w.melodyConfig
		}
		instance.HandleConnect(w.handleConnect)
		instance.HandleMessage(w.handleMessage)
		instance.HandleDisconnect(w.handleDisconnect)
		instance.HandlePong(w.handlePong)
		w.instance = instance
	})
	return w
}

func (w *WS) HandleRequest(write http.ResponseWriter, r *http.Request) error {
	return w.instance.HandleRequest(write, r)
}

func (w *WS) handleConnect(s *melody.Session) {
	sno := w.getSno()
	coon := Conn{
		sno: sno,
		ws:  s,
	}
	w.pools.Set(&coon)
	s.Set("sno", sno)
}

// 注册 WebSocket 消息处理函数
func (w *WS) handleMessage(s *melody.Session, request []byte) {
	snoStr, exist := s.Get("sno")
	if !exist {
		return
	}
	sno := snoStr.(int64)
	conn, exist := w.pools.Get(sno)
	if !exist {
		return
	}
	var _data = net.ActionModel{Conn: conn}

	WsParser.Receive(request, _data)

}

// 注册 WebSocket 连接关闭处理函数

func (w *WS) handleDisconnect(s *melody.Session) {
	snoStr, exist := s.Get("sno")
	if !exist {
		return
	}
	sno := snoStr.(int64)
	w.pools.Del(sno)
}

func (w *WS) handlePong(s *melody.Session) {
	fmt.Println("Received pong from session:", s)
}
