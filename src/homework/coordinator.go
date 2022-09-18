package homework

import (
	"fmt"
	"log"
	"math/rand"
	"net/rpc"
	"sync"
	"time"
)

type Coordinator struct {
	lock       sync.RWMutex
	stage      string     //运行阶段
	toDo       chan NUM   //需要更新的数据
	Deadline   time.Time  //Deadline
	WorkersAdd [3]string  //workers列表
	Data       [10000]int //存储的数据，默认值为0
}

var DataPre [10000]int

func MakeCoordinator() *Coordinator {
	var i int
	c := Coordinator{
		toDo:       make(chan NUM, 3),
		WorkersAdd: [3]string{":8000", ":8001", ":8002"},
	}
	for i := 0; i < 10000; i++ {
		c.Data[i] = rand.Intn(10)
	}
	DataPre = c.Data
	for {
		fmt.Printf("第%d次Check\n", i)
		if c.stage != DONE {
			c.Check()
			fmt.Printf("所有都连接成功，并传输完数据，开始让Worker Commit\n")
			//让Worker提交并且Updata
			if c.DoCommit() {
				c.UpdataData()
			} else {
				c.RollBack()
			}
		} else {
			break
		}
		i++
		//因为访问过快，来不及释放
		if i%500 == 0 {
			time.Sleep(10 * time.Second)
		}
	}
	for i, j := range c.Data {
		fmt.Printf("下标：%d,  原先：%d，现在：%d\n", i, DataPre[i], j)
	}
	return &c
}

func (c *Coordinator) Test() {
	fmt.Printf("测试")
}

func (c *Coordinator) UpdataData() {
	for i := 0; i < 3; i++ {
		msg := <-c.toDo
		go func() {
			c.lock.Lock()
			if &msg.index == nil {
				goto end
			}
			if msg.index >= 10000 {
				msg.index = msg.index % 10000
			}
			c.Data[msg.index] = msg.Number
			fmt.Printf("下标：%d ；更新数字：%d\n", msg.index, msg.Number)
			c.lock.Unlock()
		end:
		}()
	}
	fmt.Printf("更新数据完成\n")
}

//Check 所有服务器
func (c *Coordinator) Check() bool {
	args := OrderInfo{
		Stage: PREPARE,
	}
	for _, j := range c.WorkersAdd {
		conn, err := rpc.DialHTTP("tcp", j)
		if err != nil {
			fmt.Printf("1 Check连接失败\n")
			log.Fatal(err)
		}
		reply := Reply{}
		//得到回应，并且知道所需要的数据
		err1 := conn.Call("Worker.DoPrePare", args, &reply)
		if reply.Stage == DONE {
			c.stage = DONE
			break
		}
		if err1 != nil {
			continue
		}
		defer conn.Close()
		//读取数据 和 传输数据
		c.ReadData(reply.IndexI, j)
	}
	return true
}

func (c *Coordinator) ReadData(IndexI int, add string) OrderInfo {
	c.lock.RLock()
	arg := OrderInfo{}
	arg.Args[0] = c.Data[(IndexI)%10000]
	arg.Args[1] = c.Data[(IndexI+1)%10000]
	arg.Args[2] = c.Data[(IndexI+2)%10000]
	c.lock.RUnlock()
	conn, err := rpc.DialHTTP("tcp", add)
	defer conn.Close()
	if err != nil {
		log.Fatal(err)
	}
	reply := Reply{}
	err1 := conn.Call("Worker.GetData", arg, &reply)
	if err1 != nil {
		log.Fatal(err1)
	}
	/*
		if reply.Stage == Checked {
			fmt.Printf("数据传输成功\n")
		}
	*/
	return arg
}

func (c *Coordinator) DoCommit() bool {
	for _, j := range c.WorkersAdd {
		conn, err := rpc.DialHTTP("tcp", j)
		defer conn.Close()
		if err != nil {
			return false
			log.Fatal(err)
		}
		arg := OrderInfo{
			Stage: COMMIT,
		}
		reply := NUM{}
		//得到回应，并且知道所需要的数据
		//不知道为什么INDEX J 传送不过来
		err1 := conn.Call("Worker.DoCommit", arg, &reply)
		if err1 != nil {
			continue
		}
		var index int
		err2 := conn.Call("Worker.DoIndex", arg, &index)
		if err2 != nil {
			continue
		}
		reply.index = index
		c.toDo <- reply
	}
	return true
}

func (c *Coordinator) RollBack() {

}

//一个函数访问多个方法
func (c *Coordinator) Call(add string, rpcname string, Args interface{}, Rly interface{}) bool {
	// 拨号服务
	coon, err := rpc.DialHTTP("tcp", add)
	if err != nil {
		log.Fatal("dialing:", err)
	}
	err = coon.Call(rpcname, Args, Rly)
	defer coon.Close()
	if err == nil {
		return true
	}
	fmt.Println(err)
	return false
}
