package inet

import (
	"errors"
	"fmt"
	"github.com/study/iface"
	"sync"
)

var NOTFOUNDCONN = errors.New("connction not found")

/*
	管理连接
*/

type ConnectionManger struct {
	connections map[uint32] iface.IConnection
	connlock sync.RWMutex						//保护连接集合的读写锁
}

func NewConnManger() *ConnectionManger{
	return &ConnectionManger{
		connections: make(map[uint32] iface.IConnection),
	}
}

//添加链接
func(c *ConnectionManger)Add(connection iface.IConnection){
	//保护共享资源 加写锁
	c.connlock.Lock()
	defer  c.connlock.Unlock()

	//add
	c.connections[connection.GetConnID()] = connection

	fmt.Println("add conn success  now count:" ,c.Len())

}

//删除链接
func(c *ConnectionManger)Remove(connection iface.IConnection){
	c.connlock.Lock()
	defer  c.connlock.Unlock()

	delete(c.connections,connection.GetConnID())
	fmt.Println("del conn success  now count:" ,c.Len())
}

//根据connid获取链接
func(c *ConnectionManger)Get(cid uint32)(iface.IConnection,error){
	//read lock
	c.connlock.RLock()
	defer c.connlock.RUnlock()

	if conn,isset := c.connections[cid]; isset{
		return conn,nil
	}else{
		return nil,NOTFOUNDCONN
	}
}
//得到当前连接总数
func(c *ConnectionManger)Len() int{
	return len(c.connections)
}
//清除并终止所有连接
func(c *ConnectionManger) ClearConn(){
	//write lock
	c.connlock.Lock()
	defer  c.connlock.Unlock()

	//删除conn并停止conn
	for connID,conn:= range c.connections{
		conn.Stop()
		delete(c.connections,connID)
	}

	fmt.Println("Clear ALL connections success!")
}