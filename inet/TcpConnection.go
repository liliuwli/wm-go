package inet

import (
	"errors"
	"fmt"
	"github.com/study/iface"
	"github.com/study/utils"
	"io"
	"net"
	"sync"
)

var EMPTYATTR = errors.New("this attr is not set")

type TcpConnection struct {
	Id uint32

	Conn *net.TCPConn

	isClosed bool

	ExitChan chan bool

	MsgChan chan []byte

	Handler iface.IHandler

	MyServer iface.IServer

	//读取缓冲区 后续加缓冲区满的操作
	recvBuffer []byte

	//累计读取
	bytesRead int

	//当前包长
	CurrentPackageLength int

	//当前协议
	Parser iface.IProtocol

	//链接属性
	property map[string]interface{}
	propertyLock sync.RWMutex
}

func NewConnection(server iface.IServer,conn *net.TCPConn,connID uint32,handler iface.IHandler) *TcpConnection {
	c := &TcpConnection{
		Id		:  connID,
		Conn	:  conn,
		isClosed:  false,
		Handler	: handler,
		ExitChan:  make(chan bool,1),
		MsgChan: make(chan []byte),
		MyServer:server,
		property:make(map[string]interface{}),
		//缓冲区初始化
		recvBuffer:[]byte(""),
		bytesRead:0,
		CurrentPackageLength:0,
		Parser : NewBase(),
	}

	//添加conn加入到connectionmanger
	c.MyServer.GetConnManger().Add(c)

	return c
}

// 链接得读业务方法
func (c *TcpConnection) StartReader() {
	fmt.Println("Reader Gorountine is running ... ")

	defer c.Stop()
	defer fmt.Println("connID = ",c.Id, "Reader is exit")

	for {
		//单次读
		buffer := make([]byte,utils.GlobalObject.MaxPackageSize)

		read_int,err := c.Conn.Read(buffer)
		real_buffer,_ := utils.BytesSlice(buffer,read_int)
		fmt.Printf("telnet readlen : %d \n\n\n", read_int)
		if err != nil {
			//连接已关闭
			if err == io.EOF{
				break
			}else{
				fmt.Println("read stream error")
				break
			}
		}

		//加载数据到缓冲区
		c.bytesRead += read_int
		c.recvBuffer = utils.BytesCombine(c.recvBuffer,real_buffer)

		fmt.Printf(" read length : %d recv data:%s \n",c.bytesRead,c.recvBuffer)

		for{
			recvlen := len(c.recvBuffer)
			if recvlen > 0 {
				// 缓冲区有数据
				if c.CurrentPackageLength > 0 {
					// 当前包数据不足
					if c.CurrentPackageLength > recvlen {
						break
					}
				} else {
					//获取包长
					c.CurrentPackageLength,err = c.Parser.GetHeadLen(c.recvBuffer,c)
					if err != nil {
						fmt.Println("uncatch a head")
						break
					}

					//当前数据无包长
					if c.CurrentPackageLength == 0 {
						break
					}

					if c.CurrentPackageLength > 0 && c.CurrentPackageLength <= utils.GlobalObject.MaxPackageSize{
						//缓冲区数据 不够解析数据包
						if c.CurrentPackageLength > recvlen {
							break
						}
					} else {
						//读到的包长过大
						fmt.Printf("package_length : %d",c.CurrentPackageLength)
						c.Stop()
						break
					}

				}
			}else{
				break
			}




			//当前数据足够解析出n个包
			reqbuf := make([]byte,c.CurrentPackageLength)
			if recvlen == c.CurrentPackageLength {
				reqbuf = c.recvBuffer
				c.recvBuffer = []byte("")
			} else {
				reqbuf,c.recvBuffer = utils.BytesSlice(c.recvBuffer,c.CurrentPackageLength)
			}
			c.CurrentPackageLength = 0

			msgbuf,err := c.Parser.Unpack(reqbuf)

			if err != nil {
				fmt.Println("ParserDecodeMsgError")
				break
			}

			req := Request{
				conn: c,
				data: msgbuf,
			}

			if utils.GlobalObject.WorkerPoolSize > 0{
				c.Handler.SendMsgToTask(&req)
			}else{
				go c.Handler.DoHandler(&req)
			}

		}
	}
}

func (c *TcpConnection) StartWriter(){
	fmt.Println("Writer Gorountine is running ... ")

	defer fmt.Println(c.GetRemoteAddr().String(),"[conn writer exit!]")

	//wait for channel msg
	for{
		select {
		case data := <-c.MsgChan:
			if _,err := c.Conn.Write(data);err != nil {
				fmt.Println("send data error,",err," conn writer exit!")
				return
			}
		case <-c.ExitChan:
			return
		}
	}
}

func (c *TcpConnection) Start (){
	fmt.Println("connection start : id = ",c.Id)

	//read handle
	go c.StartReader()
	//writer handle
	go c.StartWriter()

	//触发start回调
	c.MyServer.CallOnConnectionStart(c)
}

func (c *TcpConnection) Stop(){
	//log
	fmt.Println("connection stop : id = ",c.Id)
	if c.isClosed == true{
		return
	}

	c.isClosed = true

	c.MyServer.CallOnConnectionStop(c)

	c.Conn.Close()

	c.ExitChan <- true

	c.MyServer.GetConnManger().Remove(c)

	close(c.ExitChan)
	close(c.MsgChan)
}

func (c *TcpConnection) GetTcpConnection() *net.TCPConn{
	return c.Conn
}

func (c *TcpConnection) GetConnID() uint32{
	return c.Id
}

func (c *TcpConnection) GetRemoteAddr() net.Addr{
	return c.Conn.RemoteAddr()
}

func (c *TcpConnection) Send(data []byte) error{
	if c.isClosed {
		return errors.New("conn is closed when send msg")
	}

	send,err := c.Parser.Pack(data)
	if err != nil {
		return errors.New("send msg pack error")
	}

	c.MsgChan <- send

	return nil
}

func (c *TcpConnection) SetProperty(key string,value interface{}){
	c.propertyLock.Lock()
	defer c.propertyLock.Unlock()

	c.property[key] = value
}
func (c *TcpConnection) GetProperty(key string) (interface{},error){
	c.propertyLock.RLock()
	defer c.propertyLock.RUnlock()

	if value,isset := c.property[key]; isset{
		return value,nil
	} else {
		return nil,EMPTYATTR
	}

}
func (c *TcpConnection) RemoveProperty(key string){
	c.propertyLock.Lock()
	defer c.propertyLock.Unlock()

	delete(c.property,key)
}