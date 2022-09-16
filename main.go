package main

import (
	NacosClient "nacos-grpc-service/client"
	NacosService "nacos-grpc-service/server"
	"strconv"
	"sync"
)

var (
	wg sync.WaitGroup
)

//开启 ProtoBuf 生成的 gRPC服务 并将服务注册到 Nacos 服务器中。
func main() {

	NacosClient.InitNacosClient([]string{"192.168.66.146:8848"}) //设置Nacos服务端的Ip地址和端口

	count := 10000 //选择开启的服务数

	wg.Add(count + 1) // 增加计数

	serviceName := "person.go"      //设置服务名
	serviceHost := "192.168.66.146" //注册服务的ip地址
	for i := 1; i <= count; i++ {
		servicePort := uint64(8560 + i) //服务绑定的端口
		go NacosService.Serve(serviceName+strconv.Itoa(i%1000), serviceHost, servicePort)
		wg.Done() // 操作完成，减少一个计数
	}

	wg.Wait()
}
