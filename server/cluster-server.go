package main

import (
	"context"
	"fmt"
	"github.com/nacos-group/nacos-sdk-go/v2/clients"
	"github.com/nacos-group/nacos-sdk-go/v2/common/constant"
	"github.com/nacos-group/nacos-sdk-go/v2/vo"
	"nacos-grpc-gin/pb/person"
	"net"
	"strconv"
	"sync"
	"time"

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
	var wg sync.WaitGroup

	//go serve("47.115.216.190", "192.168.66.146", "47.115.216.190", 8461)
	//go serve("47.115.216.190", "192.168.66.146", "47.115.216.190", 8462)

	wg.Add(2) // 因为有两个动作，所以增加2个计数
	go func() {
		serve([]string{"192.168.66.146", "192.168.66.148"}, "192.168.66.146", "106.52.77.111", 8461, "")
		wg.Done() // 操作完成，减少一个计数
	}()
	go func() {
		serve([]string{"192.168.66.146", "192.168.66.148"}, "192.168.66.146", "106.52.77.111", 8462, "cluster-146")
		wg.Done() // 操作完成，减少一个计数
	}()

	//running()
	wg.Wait()
	//serve("47.115.216.190", "192.168.66.146", "47.115.216.190", 8462)

}

func running() {
	var times int
	// 构建一个无限循环
	for {
		times++
		//fmt.Println("tick", times)
		// 延时1秒
		time.Sleep(time.Second)
	}
}

func serve(serverAddr []string, localAddr string, npsAddr string, servicePort uint64, cluster string) {
	RegisterNacos(serverAddr, npsAddr, servicePort, cluster)
	//////////////////////以下为 grpc 服务远程调用//////////////////////////////

	// 1.初始化 grpc 对象,
	grpcServer := grpc.NewServer()

	// 2.注册服务
	person.RegisterHelloServer(grpcServer, new(Children))

	// 3.设置监听, 指定 IP/port
	listener, err := net.Listen("tcp", localAddr+":"+strconv.FormatUint(servicePort, 10))
	if err != nil {
		fmt.Println("Listen err:", err)
		return
	}
	defer listener.Close()

	fmt.Println(strconv.FormatUint(servicePort, 10) + "服务启动... ")

	// 4. 启动服务
	grpcServer.Serve(listener)
}

func RegisterNacos(serverAddr []string, npsAddr string, servicePort uint64, cluster string) {
	// 创建clientConfig
	clientConfig := constant.ClientConfig{
		//NamespaceId:         "e525eafa-f7d7-4029-83d9-008937f9d468", // 如果需要支持多namespace，我们可以场景多个client,它们有不同的NamespaceId。当namespace是public时，此处填空字符串。
		TimeoutMs:           5000,
		NotLoadCacheAtStart: true,
		LogDir:              "/tmp/nacos/log",
		CacheDir:            "/tmp/nacos/cache",
		//RotateTime:          "1h",
		//MaxAge:              3,
		LogLevel: "debug",
	}

	// 至少一个ServerConfig
	serverConfigs := []constant.ServerConfig{
		{
			IpAddr:      serverAddr[0],
			ContextPath: "/nacos",
			Port:        8848,
			Scheme:      "http",
		},
		{
			IpAddr:      serverAddr[0],
			ContextPath: "/nacos",
			Port:        8849,
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
		Ip:          npsAddr,
		Port:        servicePort,
		ServiceName: "demo.go",
		Weight:      10,
		Enable:      true,
		Healthy:     true,
		Ephemeral:   true,
		Metadata:    map[string]string{"idc": "shanghai"},
		ClusterName: cluster, // 默认值DEFAULT
		GroupName:   "",            // 默认值DEFAULT_GROUP11
	})
	if !success {
		return
	} else {
		fmt.Println("namingClient.RegisterInstance Success")
	}
}
