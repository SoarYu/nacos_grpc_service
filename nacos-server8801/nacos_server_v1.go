package main

import (
	"context"
	"fmt"
	"github.com/nacos-group/nacos-sdk-go/clients"
	"github.com/nacos-group/nacos-sdk-go/common/constant"
	"github.com/nacos-group/nacos-sdk-go/vo"
	"nacos_person/pb/person"
	"net"

	"google.golang.org/grpc"
)

// 定义类
type Children struct {
}

// 绑定类方法, 实现借口
func (this *Children) SayHello(ctx context.Context, p *person.Person) (*person.Person, error) {
	p.Name = "hello  " + p.Name
	return p, nil
}

func main() {
	RegisterNacos()
	//////////////////////以下为 grpc 服务远程调用//////////////////////////////

	// 1.初始化 grpc 对象,
	grpcServer := grpc.NewServer()

	// 2.注册服务
	person.RegisterHelloServer(grpcServer, new(Children))

	// 3.设置监听, 指定 IP/port
	listener, err := net.Listen("tcp", "192.168.66.146:8801")
	if err != nil {
		fmt.Println("Listen err:", err)
		return
	}
	defer listener.Close()

	fmt.Println("服务启动... ")

	// 4. 启动服务
	grpcServer.Serve(listener)

}

func RegisterNacos() {
	// 创建clientConfig
	clientConfig := constant.ClientConfig{
		NamespaceId:         "", // 如果需要支持多namespace，我们可以场景多个client,它们有不同的NamespaceId。当namespace是public时，此处填空字符串。
		TimeoutMs:           5000,
		NotLoadCacheAtStart: true,
		LogDir:              "/tmp/nacos/log",
		CacheDir:            "/tmp/nacos/cache",
		//RotateTime:          "1h",
		//MaxAge:              3,
		LogLevel:            "debug",
	}

	// 至少一个ServerConfig
	serverConfigs := []constant.ServerConfig{
		{
			IpAddr:      "192.168.66.148",
			ContextPath: "/nacos",
			Port:        8848,
			Scheme:      "http",
		},
	}

	// 创建服务发现客户端 (推荐)
	namingClient, err := clients.NewNamingClient(
		vo.NacosClientParam{
			ClientConfig:  &clientConfig,
			ServerConfigs: serverConfigs,
		},
	)
	if err != nil {
		fmt.Println("clients.NewNamingClient err,", err)
	}
	success, err := namingClient.RegisterInstance(vo.RegisterInstanceParam{
		Ip:          "192.168.66.146",
		Port:        8801,
		ServiceName: "demo.go",
		Weight:      10,
		Enable:      true,
		Healthy:     true,
		Ephemeral:   true,
		Metadata:    map[string]string{"idc": "shanghai"},
		ClusterName: "cluster-a", // 默认值DEFAULT
		GroupName:   "group-a",   // 默认值DEFAULT_GROUP
	})
	if !success {
		return
	} else {
		fmt.Println("namingClient.RegisterInstance Success")
	}
}
