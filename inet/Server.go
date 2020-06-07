package inet

import (
    "errors"
    "fmt"
    "github.com/study/utils"
    "net"
    "github.com/study/iface"
)

type Server struct {
    Name        string
    IPVersion   string
    IP          string
    Port        int
    //当前work的回调管理器
    Handler     iface.IHandler
    //当前server连接管理器
    ConnManger  iface.IConnectionManger

    //回调函数
    OnConnectionStart func(conn iface.IConnection)
    OnConnectionStop func(conn iface.IConnection)
    OnWorkerStart func(workerid int)
}

func CallBackToClient(conn *net.TCPConn,data []byte,cnt int)error{
    fmt.Printf("[handler] start %s %d \n",data,cnt)
    if _,err := conn.Write(data[:cnt]); err != nil {
        fmt.Println("write back buf err",err)
        return errors.New("callbacktoclient error")
    }
    return nil
}

func (s *Server) Start() {
    fmt.Printf("Server Name : %s , listenner at Ip: %s,Port: %d is starting \n",
        utils.GlobalObject.Name,utils.GlobalObject.Host,utils.GlobalObject.TcpPort)
    go func() {
        //开启workerlist
        s.Handler.StartWorkerPool()

        // 获取一个tcp得addr
        addr,err := net.ResolveTCPAddr(s.IPVersion,fmt.Sprintf("%s:%d",s.IP,s.Port))
        
        if err != nil {
            fmt.Println("resolve tcp addt error",err)
            return 
        }
        
        // 监听地址端口
        listenner,err := net.ListenTCP(s.IPVersion,addr)
        
        if err != nil {
            fmt.Println("listen ",s.IPVersion," err ",err)
            return 
        }
        
        fmt.Println(" start  server success : ",s.Name," success , listenning")
        // 处理客户端链接
        var cid uint32
        cid = 0
        for{
            conn,err := listenner.AcceptTCP()
            
            if err != nil {
                fmt.Println("Accept err",err)
                continue
            }

            //最大连接数
            if s.ConnManger.Len() >= utils.GlobalObject.MaxConn {
                fmt.Println("connection too many")
                conn.Close()
                continue;
            }

            Connection := NewConnection(s,conn,cid,s.Handler)
            cid++
            Connection.Start()
        }
    }()
}

func (s *Server) Stop() {
    fmt.Println("server stop")
    //回收资源
    s.ConnManger.ClearConn()
}


func (s *Server) Serve() {
    s.Start()
    
    select{}
}

func (s *Server) GetConnManger () iface.IConnectionManger {
    return s.ConnManger
}

func (s *Server) AddRouter(callbackid uint32,router iface.IRouter) {
    s.Handler.AddRouter(callbackid,router)
    fmt.Println("add router success \n")
}

/*
    初始化的server模块的方法
*/

func NewServer (name string) iface.IServer{
    s := &Server{
        Name        :utils.GlobalObject.Name,
        IPVersion   :"tcp4",
        IP          :utils.GlobalObject.Host,
        Port        :utils.GlobalObject.TcpPort,
        Handler     :nil,
        ConnManger  :NewConnManger(),
    }

    s.Handler = NewHandle(s)
    return s
}



//设置回调
func (s *Server) SetOnConnectionStart(callback func(connection iface.IConnection)){
    s.OnConnectionStart = callback
}
func (s *Server) SetOnConnectionStop(callback func(connection iface.IConnection)){
    s.OnConnectionStop = callback
}
func (s *Server) SetOnWorkerStart(callback func(workerid int)){
    s.OnWorkerStart = callback
}

//触发回调
func (s *Server) CallOnConnectionStart(connection iface.IConnection){
    if s.OnConnectionStart != nil{
        fmt.Println("connction start callback ")
        s.OnConnectionStart(connection)
        return
    }
}
func (s *Server) CallOnConnectionStop(connection iface.IConnection){
    if s.OnConnectionStop != nil {
        fmt.Println("connction stop callback  111")
        s.OnConnectionStop(connection)
        return
    }
}
func (s *Server) CallOnWorkerStart(workerid int){
    if s.OnWorkerStart != nil{
        fmt.Println("worker start callback")
        s.OnWorkerStart(workerid)
        return
    }
}