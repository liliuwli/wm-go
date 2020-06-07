package main

import (
	"fmt"
	"github.com/study/iface"
	"github.com/study/inet"
)


//php trait
type SelfRouter struct {
	inet.BaseRouter
}

func (br *SelfRouter) Handle(request iface.IRequest){

	fmt.Println("call router Handle ")

	fmt.Printf("recv msg : %s  \n",request.GetData())
	//读取到的数据
	err := request.GetConnection().Send([]byte("pong \n"))
	if err != nil{
		fmt.Println(err)
	}

}

func OnConnectionStart(conn iface.IConnection){
	fmt.Println("a connection start")
	conn.Send([]byte("hello"))
}

func OnConnectionStop(conn iface.IConnection){
	fmt.Println("a connection will be close id = ",conn.GetConnID())
}


func main() {
	//创建一个server句柄，使用api
	s := inet.NewServer("test1.0")
	//onmessage
	s.AddRouter(0,&SelfRouter{})

	s.SetOnConnectionStart(OnConnectionStart)
	s.SetOnConnectionStop(OnConnectionStop)
	//启动server
	s.Serve()
}
