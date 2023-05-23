package webscockt

import (
	"encoding/json"
	"fmt"
	"github.com/olahol/melody"
	"go-fast/open.go/net"
)

// Conn WebSocket连接对象
type Conn struct {
	//连接配置
	ServerItem *net.ServerItem

	//当前编号
	sno int64

	//status 状态
	status int

	//连接ID,由业务逻辑中生成一个唯一ID.如:用户ID
	id int64

	//连接分组,由业务逻辑中生成的分组.如: 房间ID
	group int

	//websocket连接
	ws *melody.Session

	//连接中保存的其他数据
	Data interface{}
}

// GetSno 获取编号
func (c *Conn) GetSno() int64 {
	return c.sno
}

// GetID 获取ID
func (c *Conn) GetID() int64 {
	return c.id
}

// SetID 设置ID
func (c *Conn) SetID(id int64) {
	c.id = id
}

func (c *Conn) GetData() interface{} {
	return c.Data
}
func (c *Conn) SetData(data interface{}) {
	c.Data = data
}

// GetGroup 设置分组
func (c *Conn) GetGroup() int {
	return c.group
}

// SetGroup .. 设置连接组
// group == 0 表示离开分组
func (c *Conn) SetGroup(group int) {
	c.group = group
}

// GetStatus 获取连接状态 (TODO:未实现,目前返回0)
func (c *Conn) GetStatus() int {
	return c.status
}

func (c *Conn) GetConfig() *net.ServerItem {
	return c.ServerItem
}

// Handshake 是否握手
// 注意 ID <=0 表示未握手
func (c *Conn) Handshake() bool {
	return c.id > 0
}

// Send 先客户端发送数据
// 这里回调解析的Encode
func (c *Conn) Send(v interface{}) {

	data, err := WsParser.Encode(v)
	if err != nil {
		c.onError(fmt.Errorf("[scocket.Encode] SNO: %v 数据编码错误\nError: %v", c.sno, err.Error()))
		return
	}
	bytes, err := json.Marshal(data)
	if err != nil {
		fmt.Println("转换失败:", err)
		return
	}

	if err = c.ws.Write(bytes); err != nil {
		c.onError(fmt.Errorf("[scocket.Send] SNO: %v 发送失败\nError: %v", c.sno, err.Error()))
	}
}

// Close 关闭连接
func (c *Conn) Close() {
	c.ws.Close()
}

// Error 连接错误()
func (c *Conn) Error(err error) {
	c.onError(err)
}

/*onError 连接中产生错误时调用*/
func (c *Conn) onError(err error) {

}
