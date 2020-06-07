package iface

//定义协议 处理粘包
type IProtocol interface {
	GetHeadLen( []byte, IConnection) (int,error)
	Pack([]byte) ([]byte,error)
	Unpack(data []byte) ([]byte,error)

}
