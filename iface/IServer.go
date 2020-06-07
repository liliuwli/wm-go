package iface

//定义服务接口
type IServer interface{
    //启动
    Start()
    //停止
    Stop()
    //运行
    Serve()
    //路由注册
    AddRouter(callbackid uint32,router IRouter)
    //获取连接管理器
    GetConnManger () IConnectionManger

    //设置回调
    SetOnConnectionStart(func(connection IConnection))
    SetOnConnectionStop(func(connection IConnection))
    SetOnWorkerStart(func(workerid int))

    //触发回调
    CallOnConnectionStart(connection IConnection)
    CallOnConnectionStop(connection IConnection)
    CallOnWorkerStart(workerid int)
}