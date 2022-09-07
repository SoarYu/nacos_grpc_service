package main

import (
	nacosService "nacos-grpc-service/server"
	"sync"
)

func main() {
	var wg sync.WaitGroup

	wg.Add(1) // 因为有两个动作，所以增加2个计数
	go func() {
		nacosService.Serve("demo.go", []string{"192.168.66.146"}, "192.168.66.146", "106.52.77.111", 8461, "")
		wg.Done() // 操作完成，减少一个计数
	}()


	wg.Wait()

}
