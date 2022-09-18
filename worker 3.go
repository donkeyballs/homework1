package main

import (
	"RPC/src/homework"
	"fmt"
	"os"
)

func main() {
	addres := ":8002"
	var id int64
	id = int64(os.Getegid())
	w := homework.MakeWork(id, addres)

	fmt.Printf("Worker Id : %d启动服务\n", w.WorkerId)

}
