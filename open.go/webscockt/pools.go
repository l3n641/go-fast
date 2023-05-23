package webscockt

import (
	"errors"
	"fmt"
	"sync"
)

/*pools 连接池结构体*/
type pools struct {
	lock  sync.RWMutex
	data  map[int64]*Conn //连接列表
	total int64           //连接总数
}

// Get 根据 sno 获取连接
func (p *pools) Get(sno int64) (*Conn, bool) {
	p.lock.RLock()
	defer p.lock.RUnlock()
	v, b := p.data[sno]
	return v, b
}

func (p *pools) GetID(id int64) *Conn {
	return p.Foreach(func(sno int64, conn *Conn) bool {
		return conn.id == id
	})
}

// Set 添加连接到连接池中
func (p *pools) Set(v *Conn) error {
	p.lock.Lock()
	defer p.lock.Unlock()
	sno := v.GetSno()
	if sno <= 0 {
		return errors.New("sno 值无效")
	}
	if p.data == nil {
		p.data = make(map[int64]*Conn)
	}
	p.data[sno] = v
	p.total++
	p.printTotal()
	return nil
}

/*Foreach 遍历连接池并返回 f() == true 的连接
 * 注意: 回调返回 true 时将中断遍历
 */
func (p *pools) Foreach(f func(sno int64, conn *Conn) bool) *Conn {
	p.lock.Lock()
	defer p.lock.Unlock()

	for i := range p.data {
		conn := p.data[i]
		if f(i, conn) {
			return conn
		}
	}

	return nil
}

// Del 删除并返回 sno 的连接
// 注意: 这里连接还未关闭
func (p *pools) Del(sno int64) (*Conn, bool) {
	p.lock.Lock()
	defer p.lock.Unlock()
	v, b := p.data[sno]
	if b {
		delete(p.data, sno)
	}

	p.total--
	p.printTotal()
	return v, b
}

/*Remove 遍历删除连接
 * single == true 时只删除一个
 * 注意: 连接在这里并没有关闭
 */
func (p *pools) Remove(f func(sno int64, conn *Conn) bool, single bool) {
	p.lock.Lock()
	defer p.lock.Unlock()
	for i := range p.data {
		conn := p.data[i]
		if f(i, conn) {
			delete(p.data, i)
			p.total--
			if single {
				break
			}
		}
	}

	p.printTotal()
}

/*Clean 关闭并移除所有连接 */
func (p *pools) Clean() {
	p.lock.Lock()
	defer p.lock.Unlock()
	for i := range p.data {
		p.data[i].ws.Close()
	}
	p.data = make(map[int64]*Conn)
	p.total = 0
}

/*Send 指定sno的连接推送数据客户端*/
func (p *pools) Send(sno int64, msg []byte) error {
	v, b := p.Get(sno)
	if b {
		v.ws.Write(msg)
		return nil
	}
	return fmt.Errorf("SNO:%v 连接不存在", sno)
}

/*SendExcept 除了指定sno的连接不推送,其他的连接全部推送数据*/
func (p *pools) SendExcept(sno int64, msg []byte) {
	p.lock.Lock()
	defer p.lock.Unlock()
	for i := range p.data {
		if i != sno {
			p.data[i].ws.Write(msg)
		}
	}
}

/*SendToID 先指定id的连接推送数据*/
func (p *pools) SendToID(id int64, data interface{}) error {
	conn := p.Foreach(func(sno int64, conn *Conn) bool {
		return conn.GetID() == id
	})
	if conn != nil {
		conn.Send(data)
		return nil
	}
	return fmt.Errorf("ID:%v 连接不存在", id)
}

/*SendExceptID 除了指定的ID不推送数据,其他连接全部推送数据*/
func (p *pools) SendExceptID(id int64, data interface{}) {
	p.lock.Lock()
	defer p.lock.Unlock()
	for i := range p.data {
		if p.data[i].GetID() != id {
			p.data[i].Send(data)
		}
	}
}

/*SendToGroup 向指定分组的连接推送数据*/
func (p *pools) SendToGroup(group int, data interface{}) {
	p.Foreach(func(sno int64, conn *Conn) bool {
		if conn.GetGroup() == group {
			conn.Send(data)
		}
		return false
	})
}

/*SendAll 向所有连接推送消息*/
func (p *pools) SendAll(data interface{}) {
	p.Foreach(func(sno int64, conn *Conn) bool {
		conn.Send(data)
		return false
	})
}

/*Close 关闭指定sno的连接*/
func (p *pools) Close(sno int64) error {
	v, b := p.Del(sno)
	if b {
		v.Close()
		return nil
	}

	return fmt.Errorf("SNO:%v 连接不存在", sno)
}

/*CloseID 关闭指定ID的连接*/
func (p *pools) CloseID(id int64) error {
	var b bool
	p.Remove(func(sno int64, conn *Conn) bool {
		if conn.GetID() == id {
			conn.Close()
			return true
		}
		return false
	}, true)

	if b {
		return nil
	}
	return fmt.Errorf("ID:%v 连接不存在", id)
}

func (p *pools) CloseOther(sno int64, id int64) {
	p.Remove(func(s int64, conn *Conn) bool {
		if conn.GetID() == id && sno != s {
			conn.Close()
			return true
		}
		return false
	}, true)
}

/*CloseGroup 关闭指定分组的连接*/
func (p *pools) CloseGroup(group int) {
	p.Remove(func(sno int64, conn *Conn) bool {
		if conn.GetGroup() == group {
			conn.Close()
			return true
		}
		return false
	}, false)
}

/*CloseAll 关闭所有连接*/
func (p *pools) CloseAll() {
	p.Clean()
}

// 打印当前连接的数量
func (p *pools) printTotal() {
	fmt.Printf("[Pool] 当前连接数: %v", len(p.data))
}

// Pools 连接池管理
var Pools = pools{data: make(map[int64]*Conn)}
