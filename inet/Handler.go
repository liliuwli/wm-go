package inet

import (
	"fmt"
	"github.com/study/iface"
	"github.com/study/utils"
	"strconv"
)

type Handler struct {
	//every callbackid for a callback func
	Apis map[uint32] iface.IRouter
	//worker pool number
	WorkerPoolSize uint32
	//a message queue from connection to worker
	TaskQueue []chan iface.IRequest

	MyServer iface.IServer
}

func NewHandle(server iface.IServer) *Handler{
	return &Handler{
		Apis:make(map[uint32] iface.IRouter),
		WorkerPoolSize:utils.GlobalObject.WorkerPoolSize,
		TaskQueue:make([]chan iface.IRequest,utils.GlobalObject.WorkerPoolSize),
		MyServer:server,
	}
}

func (h* Handler) DoHandler(request iface.IRequest)  {
	handler,isset := h.Apis[request.GetMsgId()]

	if !isset {
		panic(" api callback id :" + strconv.Itoa(int(request.GetMsgId())) + "is not found ")
	}
	handler.PreHandle(request)
	handler.Handle(request)
	handler.PostHandle(request)
}

func (h* Handler) AddRouter(cid uint32,router iface.IRouter){
	if _,isset := h.Apis[cid]; isset {
		h.Apis[cid] = router
		panic(" repeat router id :" + strconv.Itoa(int(cid)))
	}

	h.Apis[cid] = router
}

//all worker start
func (h* Handler) StartWorkerPool(){
	for i := 0;i < int(h.WorkerPoolSize);i++{
		h.TaskQueue[i] = make(chan iface.IRequest,utils.GlobalObject.MaxWorkerTaskLen)
		go h.StartOneWorker(i,h.TaskQueue[i])
	}
}

//a worker start
func (h* Handler) StartOneWorker(workerid int,taskqueen chan iface.IRequest){
	fmt.Println("worker id = ",workerid, " starting ")
	//触发回调
	h.MyServer.CallOnWorkerStart(workerid)
	for{
		select {
		//when a request come
		case req := <- taskqueen:
			h.DoHandler(req)
		}
	}
}

func (h* Handler) SendMsgToTask(req iface.IRequest)  {
	//create a array , item is queue length
	queuesize := make([]int,int(h.WorkerPoolSize))
	for i:=0;i<int(h.WorkerPoolSize);i++{
		queuesize[i] = len(h.TaskQueue[i])
	}

	// free work do task first
	workerid := utils.ArrMinKeyFind(queuesize)
	h.TaskQueue[workerid] <- req
}