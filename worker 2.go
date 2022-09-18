package main

import (
	"RPC/src/homework"
	"fmt"
	"os"
)

func main() {
	addres := ":8001"
	var id int64
	id = int64(os.Getegid() + 2)
	j := homework.MakeWork(id, addres)
	fmt.Printf("Worker Id : %d启动服务\n", j.WorkerId)
}
