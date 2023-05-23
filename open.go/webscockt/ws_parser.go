package webscockt

import (
	"encoding/json"
	"go-fast/open.go/language"
	"go-fast/open.go/net"
)

// wsParser 默认消息解析器
type wsParser struct {
}

/*Receive 收到消息
* 1.解析数据
* 2.调用默认的路由器
* 注意: 该解析器只解析json字符串的消息/返回json字符串消息
 */
func (a *wsParser) Receive(request []byte, data net.ActionModel) {

	err := a.Decode(request, &data)

	data.Data = request

	if err != nil {
		data.Conn.Send(net.RespModel{Msg: language.L("收到无法解析的请求")})
		return
	}

	a.toRouter(data)
}

// Decode 解析收到的消息
func (a *wsParser) Decode(data []byte, v *net.ActionModel) error {
	return json.Unmarshal(data, &v.Mess)
}

// Encode 发送给客户端消息预处理
// 注意: 只是返回消息,并没有发送消息
func (a *wsParser) Encode(data interface{}) (interface{}, error) {

	switch data.(type) {
	case string:
		return data, nil
	case []byte:
		return string(data.([]byte)), nil
	}

	_bye, err := json.Marshal(data)

	return string(_bye), err
}

// toRouter websocket 路由
func (a *wsParser) toRouter(data net.ActionModel) {

	_router, _has := WebSocketRouters.GetRouter(data.Mess.Action)
	if !_has {
		data.Conn.Send(net.RespModel{Msg: language.L("未知处理器")})
		return
	}

	switch _router.Type {
	case net.CheckTypeNone:
		a.authNone(_router, data)
		return
	case net.CheckTypePlatform:
		a.authPlatform(_router, data)
		return
	case net.CheckTypeUI:

	}

	a.authUI(_router, data)

}

// authPlatform 验证平台
func (a *wsParser) authPlatform(r *Router, data net.ActionModel) {
	a.authNone(r, data)
}

// authUI 验证用户
func (a *wsParser) authUI(r *Router, data net.ActionModel) {
	if !data.Conn.Handshake() {
		data.Conn.Send(net.RespModel{Action: data.Mess.Action, Flag: 998, Msg: language.L("请先执行握手操作")})
		return
	}
	a.authNone(r, data)
}

// authNone 不验证
func (a *wsParser) authNone(r *Router, data net.ActionModel) {
	if r.Handler != nil {
		v := r.Handler(data)
		if v != nil {
			data.Conn.Send(v)
		}
		return
	}
	var b net.Action

	if r.Action != nil {
		b = r.Action
	} else {
		b = r.Builder()
	}

	if b == nil {
		data.Conn.Send(net.RespModel{Msg: language.L("系统错误!!"), Action: data.Mess.Action})
		return
	}

	v := b.Process(data)
	if v != nil {
		data.Conn.Send(v)
	}
}

var WsParser = wsParser{}
