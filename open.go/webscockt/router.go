package webscockt

import (
	"go-fast/open.go/net"
	"sync"
)

/*
Router 路由模型
路由可以是一个方法
或者一个Action接口
*/
type Router struct {
	Handler net.Handler       //调用方法(优先级最高)
	Action  net.Action        //调用接口
	Builder func() net.Action //通过构造器获取调用接口
	Type    net.CheckType     //接口验证类型
}

// Routers 路由列表管理结构体(程序启动后就初始化完成,运行中不需要添加)
type Routers struct {
	lock sync.RWMutex
	data map[string]*Router
}

// SetRouters 设置路由列表(注意一次性设置)
func (a *Routers) SetRouters(routers map[string]*Router) {
	a.lock.Lock()
	defer a.lock.Unlock()
	a.data = routers
}

// GetRouter 获取路由
func (a *Routers) GetRouter(action string) (*Router, bool) {
	a.lock.RLock()
	defer a.lock.RUnlock()
	value, has := a.data[action]
	return value, has
}

var WebSocketRouters = Routers{}
