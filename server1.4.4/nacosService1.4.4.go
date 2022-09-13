package nacosService1_4_4

import (
	"context"
	"fmt"
	"github.com/nacos-group/nacos-sdk-go/clients"
	"github.com/nacos-group/nacos-sdk-go/clients/naming_client"
	"github.com/nacos-group/nacos-sdk-go/common/constant"
	"github.com/nacos-group/nacos-sdk-go/vo"
	"google.golang.org/grpc"
	"gopkg.in/natefinch/lumberjack.v2"
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
	// DeRegis(serviceName, serverAddr[0], servicePort)
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

	// time.Sleep(1 * time.Second)

	//RegisterNacos(serviceName, serverAddr[0], servicePort)
}

var NacosClient naming_client.INamingClient

func InitServerClient() naming_client.INamingClient {
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
	newNacosClient, err := clients.NewNamingClient(
		vo.NacosClientParam{
			ClientConfig:  &cc,
			ServerConfigs: sc,
		},
	)
	if err != nil {
		panic(err)
	}
	return newNacosClient

}

func DeRegis(serviceName string, serverAddr string, servicePort uint64) {
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

//
//func ExampleServiceClient_GetService(client naming_client.INamingClient, param vo.GetServiceParam) {
//	service, _ := client.GetService(param)
//	fmt.Printf("GetService,param:%+v, result:%+v \n\n", param, service)
//}
//
//func ExampleServiceClient_SelectAllInstances(client naming_client.INamingClient, param vo.SelectAllInstancesParam) {
//	instances, _ := client.SelectAllInstances(param)
//	fmt.Printf("SelectAllInstance,param:%+v, result:%+v \n\n", param, instances)
//}
//
//func ExampleServiceClient_SelectInstances(client naming_client.INamingClient, param vo.SelectInstancesParam) {
//	instances, _ := client.SelectInstances(param)
//	fmt.Printf("SelectInstances,param:%+v, result:%+v \n\n", param, instances)
//}
//
//func ExampleServiceClient_SelectOneHealthyInstance(client naming_client.INamingClient, param vo.SelectOneHealthInstanceParam) {
//	instances, _ := client.SelectOneHealthyInstance(param)
//	fmt.Printf("SelectInstances,param:%+v, result:%+v \n\n", param, instances)
//}
//
//func ExampleServiceClient_Subscribe(client naming_client.INamingClient, param *vo.SubscribeParam) {
//	client.Subscribe(param)
//}
//
//func ExampleServiceClient_UnSubscribe(client naming_client.INamingClient, param *vo.SubscribeParam) {
//	client.Unsubscribe(param)
//}
//
//func ExampleServiceClient_GetAllService(client naming_client.INamingClient, param vo.GetAllServiceInfoParam) {
//	service, _ := client.GetAllServicesInfo(param)
//	fmt.Printf("GetAllService,param:%+v, result:%+v \n\n", param, service)
//}
