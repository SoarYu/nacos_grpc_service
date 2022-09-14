package main

import (
	//nacosService "nacos-grpc-service/server"
	nacosService1_3_1 "nacos-grpc-service/server1.4.4"
	"strconv"
	"sync"
)

var (
	wg167 sync.WaitGroup
)

func main() {

	nacosService1_3_1.InitServerClient()

	count := 100

	wg167.Add(count + 1) // 因为有两个动作，所以增加2个计数

	for i := 1; i <= count; i++ {
		serviceName := "demo" + strconv.Itoa(i) + ".go"
		servers := []string{"106.52.77.111"}
		localAddr := "106.52.77.111"
		port := uint64(8560 + i)
		nacosService1_3_1.DeRegisterNacos(serviceName, servers[0], port)
		go nacosService1_3_1.Serve(serviceName, servers, localAddr, "", port, "")
		nacosService1_3_1.RegisterNacos(serviceName, servers[0], port)
		wg.Done() // 操作完成，减少一个计数
	}

	wg167.Wait()

}
