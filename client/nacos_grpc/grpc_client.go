package nacos_grpc

import (
	"fmt"
	"github.com/nacos-group/nacos-sdk-go/v2/clients"
	"github.com/nacos-group/nacos-sdk-go/v2/clients/naming_client"
	"github.com/nacos-group/nacos-sdk-go/v2/common/constant"
	"github.com/nacos-group/nacos-sdk-go/v2/model"
	"github.com/nacos-group/nacos-sdk-go/v2/util"
	"github.com/nacos-group/nacos-sdk-go/v2/vo"
)

type GrpcClient struct {
	clientConfig  constant.ClientConfig       //nacos-coredns客户端配置
	serverConfigs []constant.ServerConfig     //nacos服务器集群配置
	grpcClient    naming_client.INamingClient //nacos-coredns与nacos服务器的grpc连接
}

func NewGrpcClient(namespaceId string, serverAddr []string) *GrpcClient {
	var nacosGrpcClient GrpcClient

	nacosGrpcClient.clientConfig = *constant.NewClientConfig(
		constant.WithNamespaceId(namespaceId), //When namespace is public, fill in the blank string here.
		constant.WithTimeoutMs(5000),
		constant.WithNotLoadCacheAtStart(true),
		constant.WithLogDir("/tmp/nacos/log"),
		constant.WithCacheDir("/tmp/nacos/cache"),
		constant.WithLogLevel("debug"),
	)
	//Another way of create serverConfigs
	nacosGrpcClient.serverConfigs = []constant.ServerConfig{
		*constant.NewServerConfig(
			serverAddr[0],
			8848,
			constant.WithScheme("http"),
			constant.WithContextPath("/nacos"),
		),
		*constant.NewServerConfig(
			serverAddr[0],
			8849,
			constant.WithScheme("http"),
			constant.WithContextPath("/nacos"),
		),
	}
	// Another way of create naming client for service discovery (recommend)
	var err error
	nacosGrpcClient.grpcClient, err = clients.NewNamingClient(
		vo.NacosClientParam{
			ClientConfig:  &nacosGrpcClient.clientConfig,
			ServerConfigs: nacosGrpcClient.serverConfigs,
		},
	)
	if err != nil {
		fmt.Println("init nacos-client error")
	}

	return &nacosGrpcClient
}

func (ngc *GrpcClient) getGrpcClient() naming_client.INamingClient {
	//检查是否超时断连
	return ngc.grpcClient
}

func (ngc *GrpcClient) GetService(serviceName string, clusters []string) (model.Service, error) {
	param := vo.GetServiceParam{
		ServiceName: serviceName,
		GroupName:   "",
		Clusters:    clusters,
	}
	service, err := ngc.getGrpcClient().GetService(param)
	if err != nil {
		panic("GetService failed!")
	}
	fmt.Printf("GetService,param:%+v, result: \n\n", param)

	return service, err
}

func (ngc *GrpcClient) GetAllServicesInfo() (model.ServiceList, error) {
	param := vo.GetAllServiceInfoParam{
		NameSpace: "",  //optional,default:public
		GroupName: "group",  //optional,default:DEFAULT_GROUP
		PageNo:    1,   //optional,default:1
		PageSize:  100, //optional,default:10
	}
	services, err := ngc.getGrpcClient().GetAllServicesInfo(param)
	if err != nil {
		panic("GetAllServicesInfo failed!")
	}
	fmt.Printf("GetAllServicesInfo,param:%+v, result:%+v \n\n", param, services)

	return services, err
}

// Subscribe ...
func (ngc *GrpcClient) Subscribe(serviceName string, groupName string) {
	//Subscribe key=serviceName+groupName+cluster
	//Note:We call add multiple SubscribeCallback with the same key.
	param := &vo.SubscribeParam{
		ServiceName: serviceName,
		GroupName:   groupName,
		SubscribeCallback: func(services []model.Instance, err error) {
			fmt.Printf("SubscribeCallback return services:%s \n\n", util.ToJsonString(services))
		},
	}
	err := ngc.getGrpcClient().Subscribe(param)
	fmt.Println("Subscribe service ", serviceName, err)
}

func (ngc *GrpcClient) UpdateServiceInstance(param vo.UpdateInstanceParam) (bool, error) {

	success, err := ngc.getGrpcClient().UpdateInstance(param)
	if !success || err != nil {
		panic("UpdateInstance failed!")
	}
	fmt.Printf("UpdateServiceInstance,param:%+v,result:%+v \n\n", param, success)

	return success, err
}

func (ngc *GrpcClient) SelectOneHealthInstance(param vo.SelectOneHealthInstanceParam) (*model.Instance, error) {
	instance, err := ngc.getGrpcClient().SelectOneHealthyInstance(param)
	if err != nil {
		panic("SelectOneHealthyInstance failed!")
	}
	fmt.Printf("SelectOneHealthyInstance,param:%+v, result:%+v \n\n", param, instance)

	return instance, err
}

func (ngc *GrpcClient) SelectAllInstances(param vo.SelectAllInstancesParam) ([]model.Instance, error) {
	instances, err := ngc.getGrpcClient().SelectAllInstances(param)
	if err != nil {
		panic("SelectAllInstances failed!")
	}
	fmt.Printf("SelectAllInstance,param:%+v, result:%+v \n\n", param, instances)
	return instances, err
}

func (ngc *GrpcClient) SelectInstances(param vo.SelectInstancesParam) ([]model.Instance, error) {
	instances, err := ngc.getGrpcClient().SelectInstances(param)
	if err != nil {
		panic("SelectInstances failed!")
	}
	fmt.Printf("SelectInstances,param:%+v, result:%+v \n\n", param, instances)
	return instances, err
}
