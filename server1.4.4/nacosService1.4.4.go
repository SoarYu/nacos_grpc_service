package nacosService1_4_4

import (
	"context"
	"fmt"
	"github.com/nacos-group/nacos-sdk-go/clients"
	"github.com/nacos-group/nacos-sdk-go/clients/naming_client"
	"github.com/nacos-group/nacos-sdk-go/common/constant"
	"github.com/nacos-group/nacos-sdk-go/vo"
	"google.golang.org/grpc"
	"nacos-grpc-service/pb/person"
	"net"
	"strconv"
	// "time"
)

// 定义类
type Children struct {
}

// 绑定类方法, 实现借口
func (this *Children) SayHello(ctx context.Context, p *person.Person) (*person.Person, error) {
	p.Name = "hello  " + p.Name
	return p, nil
}

func Serve(serviceName string, serverAddr []string, localAddr string, npsAddr string, servicePort uint64, cluster string) {
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

var NacosClient naming_client.INamingClient

func InitServerClient() {
	sc := []constant.ServerConfig{
		{
			IpAddr: "106.52.77.111",
			Port:   8848,
		},
	}

	cc := constant.ClientConfig{
		NamespaceId:         "", //namespace id
		TimeoutMs:           5000,
		NotLoadCacheAtStart: true,
		LogDir:              "/tmp/nacos/log",
		CacheDir:            "/tmp/nacos/cache",
		LogRollingConfig:    &lumberjack.Logger{MaxSize: 10},
		LogLevel:            "debug",
		AppendToStdout:      true,
	}

	// a more graceful way to create naming client
	//var err error
	NacosClient, err := clients.NewNamingClient(
		vo.NacosClientParam{
			ClientConfig:  &cc,
			ServerConfigs: sc,
		},
	)
	if err != nil {
		panic(err)
	}

}

func DeRegisterNacos(serviceName string, serverAddr string, servicePort uint64) {
	deparam := vo.DeregisterInstanceParam{
		Ip:          serverAddr,
		Port:        servicePort,
		ServiceName: serviceName,
		Ephemeral:   true,
	}
	deRegisSUCC, _ := NacosClient.DeregisterInstance(deparam)
	fmt.Printf("RegisterServiceInstance,param:%+v,result:%+v \n\n", deparam, deRegisSUCC)
}

//func RegisterNacos(serviceName string, serverAddr []string, npsAddr string, servicePort uint64, cluster string) {
func RegisterNacos(serviceName string, serverAddr string, servicePort uint64) {

	param := vo.RegisterInstanceParam{
		Ip:          serverAddr,
		Port:        servicePort,
		ServiceName: serviceName,
		Weight:      10,
		Enable:      true,
		Healthy:     true,
		Ephemeral:   true,
		Metadata:    map[string]string{"idc": "shanghai"},
	}
	success, _ := NacosClient.RegisterInstance(param)
	fmt.Printf("RegisterServiceInstance,param:%+v,result:%+v \n\n", param, success)
}
