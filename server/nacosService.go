package server

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	NacosClient "nacos-grpc-service/client"
	"nacos-grpc-service/pb/person"
	"net"
	"strconv"
)

// 定义类
type Children struct {
}

// 绑定类方法, 实现借口
func (this *Children) SayHello(ctx context.Context, p *person.Person) (*person.Person, error) {
	p.Name = "hello  " + p.Name
	return p, nil
}

func Serve(serviceName string, serviceHost string, servicePort uint64) {
	//////////////////////以下为 grpc 服务远程调用//////////////////////////////
	NacosClient.DeRegisterNacos(serviceName, serviceHost, servicePort)
	// 1.初始化 grpc 对象,
	grpcServer := grpc.NewServer()

	// 2.注册服务
	person.RegisterHelloServer(grpcServer, new(Children))

	// 3.设置监听, 指定 IP/port
	listener, err := net.Listen("tcp", serviceHost+":"+strconv.FormatUint(servicePort, 10))
	if err != nil {
		fmt.Println("Listen err:", err)
		return
	}
	defer listener.Close()

	fmt.Println(serviceName + "服务启动... ")

	NacosClient.RegisterNacos(serviceName, serviceHost, servicePort)

	// 4. 启动服务
	grpcServer.Serve(listener)

}
