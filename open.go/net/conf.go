package net

import (
	"time"
)

// ------------------------------------------------------------
// 服务配置等等
// ------------------------------------------------------------

// ProtocolType 服务(协议)类型
type ProtocolType int

const (
	ProtocolTypeHttp         ProtocolType = 1 //http 协议
	ProtocolTypeHttps        ProtocolType = 2 //https 协议
	ProtocolTypeWebSocket    ProtocolType = 3 //ws  协议
	ProtocolTypeWebSocketSSL ProtocolType = 4 //wss 协议
)

// DataType 数据类型
type DataType int

const (
	DataTypeKeyValue          = 1 //键值对 (http/https post,get...)
	DataTypeJson              = 2 //json字符串(http/https post) (ws/wss)
	DataTypeByteArray         = 3 //字节流(ws/wss)
	DataTypeCustomerByteArray = 4 //自定义字节流(ws/wss)
	DataTypeBlob              = 5 //文件流(http/https 上传文件用)
	DataTypeString            = 6 //字符串(http/https post) (ws/wss) .. 不常用
)

// ServerItem ..
type ServerItem struct {
	Name          string            `json:"name"`      //服务名称
	Router        string            `json:"router"`    //路由名称
	Type          ProtocolType      `json:"type"`      //服务类型
	DataType      DataType          `json:"data_type"` //数据类型
	Enabled       bool              `json:"enabled"`   //是否启用
	Pattern       string            `json:"pattern"`   //监听路由
	Timeout       time.Duration     `json:"timeout"`   //超时时间(秒)
	Header        map[string]string `json:"header"`    //请求头
	RequestHeader map[string]string `json:"request_header"`
}

// Other 其他配置
type Other struct {
	Debug       bool   `json:"debug"`        //是否开启debug
	LogEnabled  bool   `json:"log_enabled"`  //是否开启日志
	LogPath     string `json:"log_path"`     //日志路径
	Welcome     string `json:"welcome"`      //启动欢迎文本
	PrintConfig bool   `json:"print_config"` //启动时是否打印配置
}
