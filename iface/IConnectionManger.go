package iface

type IConnectionManger interface {
	//添加链接
	Add(connection IConnection)
	//删除链接
	Remove(connection IConnection)
	//根据connid获取链接
	Get(cid uint32)(IConnection,error)
	//得到当前连接总数
	Len() int
	//清除并终止所有连接
	ClearConn()
}
