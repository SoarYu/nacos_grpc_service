package main

import (
	nacosService "nacos-grpc-server/server"
	"sync"
)

func main() {
	var wg sync.WaitGroup

	wg.Add(2) // 因为有两个动作，所以增加2个计数
	go func() {
		nacosService.Serve([]string{"192.168.66.146", "192.168.66.148"}, "192.168.66.146", "106.52.77.111", 8461, "")
		wg.Done() // 操作完成，减少一个计数
	}()
	go func() {
		nacosService.Serve([]string{"192.168.66.146", "192.168.66.148"}, "192.168.66.146", "106.52.77.111", 8462, "cluster-146")
		wg.Done() // 操作完成，减少一个计数
	}()

	wg.Wait()

}
