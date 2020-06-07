package utils

import (
	"encoding/json"
	"fmt"
	"github.com/study/iface"
	"io/ioutil"
)

type GlobalConf struct {
	TcpServer 	iface.IServer
	Host 		string
	TcpPort		int
	Name		string

	Version		string
	MaxConn		int
	MaxPackageSize int
	SendBufferMaxSize   int
	WorkerPoolSize uint32
	MaxWorkerTaskLen uint32
}

var GlobalObject *GlobalConf

func (g *GlobalConf) Reload()  {
	data,err := ioutil.ReadFile("conf/conf.json")
	if err != nil {
		fmt.Printf("%s", err)
		panic(err)
	}

	err = json.Unmarshal(data,&GlobalObject)
	if err != nil {
		panic(err)
	}
}

func init(){
	GlobalObject = &GlobalConf{
		//default
		TcpServer		:      nil,
		Host			:      "0.0.0.0",
		TcpPort			:      5200,
		Name			:      "DefaultServer",
		Version			:      "1.0",
		MaxConn			:      3,
		MaxPackageSize	: 	   65535,
		SendBufferMaxSize	:	   1048576,
		WorkerPoolSize	:		16,
		MaxWorkerTaskLen:		1024,
	}

	GlobalObject.Reload()
}