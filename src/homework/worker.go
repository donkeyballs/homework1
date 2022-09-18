package homework

import (
	"fmt"
	"log"
	"math"
	"math/rand"
	"net/http"
	"net/rpc"
)

type Worker struct {
	WorkerId    int64
	WorkerTimes int    //工作次数
	Data        [3]int //存储的数据
	IndexI      int    //需要的下标I
	IndexJ      int    //需要的下标J
	Up          NUM
}

func MakeWork(Id int64, add string) *Worker {
	w := Worker{
		WorkerTimes: 0,
		IndexI:      0,
		IndexJ:      0,
		WorkerId:    Id,
	}

	rand.Seed(w.WorkerId)
	//生成数据
	//启动监听
	fmt.Printf("Worker Id：%d 生成！\n", w.WorkerId)
	w.Server(add)
	return &w
}

//随机生成新数据

//做好准备，回应coordinator，并且传输自己所需要的数据
func (w *Worker) DoPrePare(args OrderInfo, reply *Reply) error {
	if w.WorkerTimes != 10000 {
		reply.IndexI = w.IndexI
		reply.Stage = Checked
	} else {
		reply.Stage = DONE

	}
	return nil
}

//收到自己需要的数据
func (w *Worker) GetData(args OrderInfo, reply *Reply) error {
	w.Data = args.Args
	reply.Stage = Checked
	return nil
}

func (w *Worker) DoCommit(args OrderInfo, reply *NUM) error {
	if args.Stage == COMMIT {
		w.Compute()
		*reply = w.Up
	}
	fmt.Printf("提交的Reply ：%d  工作次数 %d\n", reply, w.WorkerTimes)
	w.New()
	return nil
}

func (w *Worker) DoIndex(args OrderInfo, reply *int) error {
	if args.Stage == COMMIT {
		*reply = w.IndexJ
	}
	return nil
}

func (w *Worker) DoRollBack(args OrderInfo, reply *Reply) error {
	w.WorkerTimes--
	reply.Stage = ACK
	return nil
}

func (w *Worker) Compute() {
	w.WorkerTimes++
	w.Up.index = w.IndexJ
	w.Up.Number = (w.Data[0] + w.Data[1] + w.Data[2]) % math.MaxInt
}

func (w *Worker) Rollback() {
	w.WorkerTimes--
}

func (w *Worker) Check(request int, answer *int) error {
	if request != 0 {
		*answer = 1
	} else {
		*answer = 0
	}
	return nil
}

func (w *Worker) New() {
	for {
		w.IndexI = rand.Int() % 10000
		w.IndexJ = rand.Int() % 10000
		if w.IndexJ-w.IndexI > 2 || w.IndexJ-w.IndexI < 0 {
			break
		}
	}
	if w.IndexJ >= 10000 {
		w.IndexJ = w.IndexJ % 10000
	}
	if w.IndexI >= 10000 {
		w.IndexI = w.IndexI % 10000
	}
	w.Up.index = w.IndexJ
}

//服务启动
func (w *Worker) Server(add string) {
	//1. 注册w
	fmt.Printf("开始注册服务器\n")
	rpc.Register(w)
	//2. 服务处理绑定到HTTP协议上
	rpc.HandleHTTP()
	//3. 启动监听服务，端口add
	fmt.Printf("服务器开始工作\n")
	fmt.Printf("%s \n", add)
	err := http.ListenAndServe(add, nil)
	if err != nil {
		log.Fatal(err)
	}
}
