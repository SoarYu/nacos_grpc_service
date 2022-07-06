package main

import (
	"context"
	"fmt"
	"github.com/nacos-group/nacos-sdk-go/v2/clients"
	"github.com/nacos-group/nacos-sdk-go/v2/common/constant"
	"github.com/nacos-group/nacos-sdk-go/v2/vo"
	"google.golang.org/grpc"
	"nacos_person/pb/person"
	"strconv"
)

func main() {
	// 创建clientConfig
	clientConfig := *constant.NewClientConfig(
		constant.WithNamespaceId(""), //When namespace is public, fill in the blank string here.
		constant.WithTimeoutMs(5000),
		constant.WithNotLoadCacheAtStart(true),
		constant.WithLogDir("/tmp/nacos/log"),
		constant.WithCacheDir("/tmp/nacos/cache"),
		constant.WithLogLevel("debug"),
	)

	//Another way of create serverConfigs
	serverConfigs := []constant.ServerConfig{
		*constant.NewServerConfig(
			"192.168.66.146",
			8848,
			constant.WithScheme("http"),
			constant.WithContextPath("/nacos"),
		),
	}

	// Another way of create naming client for service discovery (recommend)
	namingClient, err := clients.NewNamingClient(
		vo.NacosClientParam{
			ClientConfig:  &clientConfig,
			ServerConfigs: serverConfigs,
		},
	)

	if err != nil {
		panic(err)
	}

	//SelectAllInstance可以返回全部实例列表,包括healthy=false,enable=false,weight<=0
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
	fmt.Println(p)
}
