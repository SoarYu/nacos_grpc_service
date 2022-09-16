package main

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	NacosClient "nacos-grpc-service/client"
	"nacos-grpc-service/pb/person"
	NacosService "nacos-grpc-service/server"
	"strconv"
	"sync"
	"time"
)

var (
	wg sync.WaitGroup
)

//开启 ProtoBuf 生成的 gRPC服务 并将服务注册到 Nacos 服务器中。
func main() {

	NacosClient.InitNacosClient([]string{"127.0.0.1:8848"}) //设置Nacos服务端的Ip地址和端口

	count := 10 //选择开启的服务数

	wg.Add(count + 1) // 增加计数

	serviceName := "person.go" //设置服务名
	serviceHost := "127.0.0.1" //注册服务的ip地址
	for i := 1; i <= count; i++ {
		servicePort := uint64(8560 + i) //服务绑定的端口
		go NacosService.Serve(serviceName, serviceHost, servicePort)
		wg.Done() // 操作完成，减少一个计数
	}

	time.Sleep(3 * time.Second)

	////SelectAllInstance可以返回全部实例列表,包括healthy=false,enable=false,weight<=0
	instances, _ := NacosClient.SelectAllInstances(serviceName)
	//
	fmt.Println(instances)
	//// SelectOneHealthyInstance将会按加权随机轮询的负载均衡策略返回一个健康的实例
	//// 实例必须满足的条件：health=true,enable=true and weight>0
	instance, _ := NacosClient.SelectOneHealthyInstance(serviceName)

	for i := 0; i < count; i++ {
		addr := instance.Ip + ":" + strconv.Itoa(int(instance.Port))
		fmt.Println(addr)
		//////////////////////以下为 grpc 服务远程调用//////////////////////////////
		// 1. 链接服务
		// 使用 服务发现consul 上的 IP/port 来与服务建立链接
		grpcConn, _ := grpc.Dial(addr, grpc.WithInsecure())

		// 2. 初始化 grpc 客户端
		grpcClient := person.NewHelloClient(grpcConn)

		var person person.Person
		person.Name = "Andy"
		person.Age = 18

		// 3. 调用远程函数
		p, _ := grpcClient.SayHello(context.TODO(), &person)
		fmt.Println(p)
	}

	wg.Wait()
}
