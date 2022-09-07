package main

import (
	//nacosService "nacos-grpc-service/server"
	nacosService1_3_1 "nacos-grpc-service/server1.4.4"
	"sync"
)

func main() {
	var wg sync.WaitGroup

	wg.Add(1) // 因为有两个动作，所以增加2个计数
	go func() {
		//nacosService1_3_1.Serve("demo.go", []string{"192.168.66.146"}, "192.168.66.146", "106.52.77.111", 8461, "")
		nacosService1_3_1.Serve("demo.go", []string{"192.168.66.146"}, "106.52.77.111", "106.52.77.111", 8461, "")
		wg.Done() // 操作完成，减少一个计数
	}()

	wg.Wait()

}
