package main

import (
	"nacos_person/pb/person"
	"context"
	"fmt"
	"strconv"

	"github.com/nacos-group/nacos-sdk-go/clients"
	"github.com/nacos-group/nacos-sdk-go/common/constant"
	"github.com/nacos-group/nacos-sdk-go/vo"
	"google.golang.org/grpc"
)

func main() {
	// 创建clientConfig
	clientConfig := constant.ClientConfig{
		//NamespaceId:         "e525eafa-f7d7-4029-83d9-008937f9d468", // 如果需要支持多namespace，我们可以场景多个client,它们有不同的NamespaceId。当namespace是public时，此处填空字符串。
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

	// 创建服务发现客户端的另一种方式 (推荐)
	namingClient, err := clients.NewNamingClient(
		vo.NacosClientParam{
			ClientConfig:  &clientConfig,
			ServerConfigs: serverConfigs,
		},
	)
	// SelectAllInstance可以返回全部实例列表,包括healthy=false,enable=false,weight<=0
	instances, err := namingClient.SelectAllInstances(vo.SelectAllInstancesParam{
		ServiceName: "demo.go",
		GroupName:   "group-a",             // 默认值DEFAULT_GROUP
		Clusters:    []string{"cluster-a"}, // 默认值DEFAULT
	})

	fmt.Println(instances)
	// SelectOneHealthyInstance将会按加权随机轮询的负载均衡策略返回一个健康的实例
	// 实例必须满足的条件：health=true,enable=true and weight>0
	instance, err := namingClient.SelectOneHealthyInstance(vo.SelectOneHealthInstanceParam{
		ServiceName: "demo.go",
		GroupName:   "group-a",             // 默认值DEFAULT_GROUP
		Clusters:    []string{"cluster-a"}, // 默认值DEFAULT
	})

	addr := instance.Ip + ":" + strconv.Itoa(int(instance.Port))
	fmt.Println(addr)
	//////////////////////以下为 grpc 服务远程调用//////////////////////////////
	// 1. 链接服务
	//grpcConn, _ := grpc.Dial("127.0.0.1:8800", grpc.WithInsecure())

	// 使用 服务发现consul 上的 IP/port 来与服务建立链接
	grpcConn, _ := grpc.Dial(addr, grpc.WithInsecure())

	// 2. 初始化 grpc 客户端
	grpcClient := person.NewHelloClient(grpcConn)

	var person person.Person
	person.Name = "Andy"
	person.Age = 18

	// 3. 调用远程函数
	p, err := grpcClient.SayHello(context.TODO(), &person)

	fmt.Println(p, err)
}