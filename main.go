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

	nacosService.InitNacosClient([]string{"console.nacos.io:8848"}) //设置Nacos服务端的Ip地址和端口

	count := 10 //开启服务数

	wg.Add(count + 1) // 增加计数

	for i := 1; i <= count; i++ {
		serviceName := "demo" + strconv.Itoa(i) + ".go"
		serviceHost := "127.0.0.1"      //注册服务的ip地址
		servicePort := uint64(8560 + i) //服务绑定的端口
		nacosService.DeRegisterNacos(serviceName, serviceHost, servicePort)
		go nacosService.Serve(serviceHost, servicePort)
		nacosService.RegisterNacos(serviceName, serviceHost, servicePort)
		wg.Done() // 操作完成，减少一个计数
	}

	wg.Wait()

}
