package net

import (
	"net/http"
)

// ------------------------------------------------------------
// 类型定义
// ------------------------------------------------------------
type CheckType int

const (
	CheckTypeUI       CheckType = 0 //验证UI / 验证握手
	CheckTypePlatform CheckType = 1 //验证平台 / 签名
	CheckTypeNone     CheckType = 2 //不验证
)

// ------------------------------------------------------------
// 路由定义
// ------------------------------------------------------------

// Conn 连接的方法
type ConnData interface {
	GetSno() int64
	GetID() int64
	SetID(id int64)
	GetData() interface{}
	SetData(v interface{})
	GetGroup() int
	SetGroup(group int)
	Send(v interface{})
	Error(err error)
	Handshake() bool
}

/*
MessageModel 默认的消息结构
用于路由,和验证UI
*/
type MessageModel struct {
	Action string `json:"action"`
	UI     struct {
		UserID int64  `json:"userid"`
		Token  string `json:"token"`
	} `json:"ui"`
}

type Header struct {
	http.Header
	RemoteAddr string
}

/*ActionModel 请求数据模型*/
type ActionModel struct {
	Config *ServerItem  //当前连接配置
	Conn   ConnData     //连接数据 (ConnData/ConnHttp)
	Header Header       //请求头等参数
	Mess   MessageModel //基础消息结构
	Data   interface{}  //原始数据(TODO:类型还需要考虑考虑)
}

// Resp 内部处理的回复消息模型
type RespModel struct {
	Action string `json:"action"`
	Flag   int    `json:"flag"`
	Msg    string `json:"msg"`
}

/*Action 路由接口*/
type Action interface {
	Process(data ActionModel) interface{}
}

// Handler 路由方法类型
type Handler func(data ActionModel) interface{}

/*WebSocketEvent WebSocket 事件接口*/
type WebSocketEvent interface {
	Connected(data ConnData)
	Close(data ConnData)
}

// ------------------------------------------------------------
// 解析器定义
// ------------------------------------------------------------

type Parser interface {
	Receive(data []byte, v ActionModel)
	Decode(data []byte, v *ActionModel) error
	Encode(data interface{}) (interface{}, error)
}
