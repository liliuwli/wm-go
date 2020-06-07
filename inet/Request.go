package inet

import "github.com/study/iface"

type Request struct {
	conn iface.IConnection
	data []byte
}

func(r *Request) GetConnection() iface.IConnection{
	return r.conn
}

func(r *Request) GetData() []byte {
	return r.data
}

func (r *Request) GetDataLen() int {
	return  len(r.data)
}

func (r *Request) GetMsgId() uint32{
	return 0
}