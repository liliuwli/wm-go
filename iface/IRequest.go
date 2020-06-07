package iface

type IRequest interface {
	GetConnection() IConnection
	GetData() []byte
	GetDataLen() int
	GetMsgId() uint32
}
