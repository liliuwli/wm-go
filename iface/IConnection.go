package iface

import "net"

type IConnection interface {
	//启动链接
	Start()
	//停用链接 结束当前链接得工作
	Stop()
	//获取当前链接得绑定socket
	GetTcpConnection() *net.TCPConn
	//获取当前链接模块得id
	GetConnID() uint32
	//获取远程客户端 IP PORT
	GetRemoteAddr() net.Addr
	//发送数据
	Send(data []byte) error

	//设置连接属性
	SetProperty(key string,value interface{})
	GetProperty(key string) (interface{},error)
	RemoveProperty(key string)
}
