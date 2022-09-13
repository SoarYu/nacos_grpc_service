package main

import (
	nacosService "nacos-grpc-service/server"
	"strconv"
	"sync"
)

var (
	wg sync.WaitGroup
)

func main() {

	nacosService.InitNacosClient()

	count := 100

	wg.Add(count + 1) // 因为有两个动作，所以增加2个计数

	for i := 1; i <= count; i++ {
		serviceName := "demo" + strconv.Itoa(i) + ".go"
		serviceHost := "106.52.77.111"
		localAddr := "106.52.77.111"
		servicePort := uint64(8560 + i)
		nacosService.DeRegisterNacos(serviceName, serviceHost, servicePort)
		go nacosService.Serve(localAddr, servicePort)
		nacosService.RegisterNacos(serviceName, serviceHost, servicePort)
		wg.Done() // 操作完成，减少一个计数

	}

	wg.Wait()

}
