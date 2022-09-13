package nacosService

import (
	"context"
	"fmt"
	"github.com/nacos-group/nacos-sdk-go/v2/clients"
	"github.com/nacos-group/nacos-sdk-go/v2/clients/naming_client"
	"github.com/nacos-group/nacos-sdk-go/v2/common/constant"
	"github.com/nacos-group/nacos-sdk-go/v2/vo"
	"google.golang.org/grpc"
	"nacos-grpc-service/pb/person"
	"net"
	"strconv"
)

var NacosClient naming_client.INamingClient

// 定义类
type Children struct {
}

// 绑定类方法, 实现借口
func (this *Children) SayHello(ctx context.Context, p *person.Person) (*person.Person, error) {
	p.Name = "hello  " + p.Name
	return p, nil
}

func Serve(localAddr string, servicePort uint64) {
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

func InitNacosClient() {
	// 创建clientConfig
	clientConfig := constant.ClientConfig{
		NamespaceId:         "", // 如果需要支持多namespace，我们可以场景多个client,它们有不同的NamespaceId。当namespace是public时，此处填空字符串。
		TimeoutMs:           5000,
		NotLoadCacheAtStart: true,
		LogDir:              "/tmp/nacos/log",
		CacheDir:            "/tmp/nacos/cache",
		//RotateTime:          "1h",
		//MaxAge:              3,
		LogLevel: "debug",
	}

	// 至少一个ServerConfig
	serverConfigs := make([]constant.ServerConfig, 1)

	serverConfigs[0] = *constant.NewServerConfig(
		"106.52.77.111",
		8848,
		constant.WithScheme("http"),
		constant.WithContextPath("/nacos"),
	)

	var err error
	// 创建服务发现客户端 (推荐)
	NacosClient, err = clients.NewNamingClient(
		vo.NacosClientParam{
			ClientConfig:  &clientConfig,
			ServerConfigs: serverConfigs,
		},
	)
	if err != nil {
		fmt.Println("clients.NewNamingClient err,", err)
	}

}

func DeRegisterNacos(serviceName string, serviceHost string, servicePort uint64) {
	success, _ := NacosClient.DeregisterInstance(vo.DeregisterInstanceParam{
		Ip:          serviceHost,
		Port:        servicePort,
		ServiceName: serviceName,
	})
	if !success {
		return
	} else {
		fmt.Println("namingClient.DeRegisterInstance Success")
	}
}

func RegisterNacos(serviceName string, serviceHost string, servicePort uint64) {
	success, _ := NacosClient.RegisterInstance(vo.RegisterInstanceParam{
		Ip:          serviceHost,
		Port:        servicePort,
		ServiceName: serviceName,
		Weight:      10,
		Enable:      true,
		Healthy:     true,
		Ephemeral:   true,
		Metadata:    map[string]string{"idc": "shanghai"},
		GroupName:   "", // 默认值DEFAULT_GROUP11
	})
	if !success {
		return
	} else {
		fmt.Println("namingClient.RegisterInstance Success")
	}
}
