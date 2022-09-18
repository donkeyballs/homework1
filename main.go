package main

import (
	"RPC/src/homework"
	"sync"
	"time"
)

var wg sync.WaitGroup

func main() {

	//启动一个线程 保持w监听
	//go w.Server()
	go func() {
		wg.Add(1)
		c := homework.MakeCoordinator()
		c.Test()
		wg.Done()
	}()

	time.Sleep(1 * time.Second)
	wg.Wait()

	return
}
