package iface

type IHandler interface {
	DoHandler(request IRequest)
	AddRouter(callbackid uint32,router IRouter)
	StartWorkerPool()
	SendMsgToTask(request IRequest)
}
