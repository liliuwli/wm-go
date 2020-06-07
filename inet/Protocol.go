package inet

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/study/iface"
	"github.com/study/utils"
)


//配合telnet调试
type BaseProtocol struct {}

func NewBase() *BaseProtocol {
	return &BaseProtocol{}
}

//返回包长
func (b *BaseProtocol) GetHeadLen(recv []byte,c iface.IConnection) (int,error) {

	if utils.GlobalObject.MaxPackageSize > 0 && len(recv) > utils.GlobalObject.MaxPackageSize{
		fmt.Printf("server input to large MaxPackageSize: %d %d %d",utils.GlobalObject.MaxPackageSize,len(recv),int(utils.GlobalObject.MaxPackageSize))
		c.Stop()
		PACKSIZEERROR := errors.New("PACK SIZE TOO LARGE")
		return 0,PACKSIZEERROR
	}

	pos := utils.BytesFind(recv,"\n")

	if pos != -1{
		return pos+1,nil
	}else{
		PACKEOFERROR := errors.New("undefind pack eof")
		return 0,PACKEOFERROR
	}
}

//返回加密后的包内容
func (b *BaseProtocol) Pack(data []byte) ([]byte,error){
	return utils.BytesCombine(data,[]byte("\r\n")),nil
}

//返回解密后的[]byte
func (b *BaseProtocol) Unpack(data []byte) ([]byte,error){
	return bytes.TrimRight(data,"\n"),nil
}




