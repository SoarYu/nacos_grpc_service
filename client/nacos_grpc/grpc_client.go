package nacos_grpc

import (
	"fmt"
	"github.com/nacos-group/nacos-sdk-go/v2/model"

	"github.com/nacos-group/nacos-sdk-go/v2/clients"
	"github.com/nacos-group/nacos-sdk-go/v2/clients/naming_client"
	"github.com/nacos-group/nacos-sdk-go/v2/common/constant"
	"github.com/nacos-group/nacos-sdk-go/v2/vo"
)

type GrpcClient struct {
	clientConfig  constant.ClientConfig       //nacos-coredns客户端配置
	serverConfigs []constant.ServerConfig     //nacos服务器集群配置
	grpcClient    naming_client.INamingClient //nacos-coredns与nacos服务器的grpc连接
}

func NewGrpcClient() *GrpcClient {
	var nacosGrpcClient GrpcClient

	nacosGrpcClient.clientConfig = *constant.NewClientConfig(
		constant.WithNamespaceId(""), //When namespace is public, fill in the blank string here.
		constant.WithTimeoutMs(5000),
		constant.WithNotLoadCacheAtStart(true),
		constant.WithLogDir("/tmp/nacos/log"),
		constant.WithCacheDir("/tmp/nacos/cache"),
		constant.WithLogLevel("debug"),
	)
	//Another way of create serverConfigs
	nacosGrpcClient.serverConfigs = []constant.ServerConfig{
		*constant.NewServerConfig(
			"192.168.66.146",
			8848,
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

//func (ngc *GrpcClient) setGrpcClient(gc *naming_client.INamingClient) {
//	&(ngc.grpcClient) = gc
//}

func (ngc *GrpcClient) GetService() (model.Service, error) {
	param := vo.GetServiceParam{
		ServiceName: "demo.go",
		GroupName:   "group-a",
		Clusters:    []string{"cluster-a"},
	}
	service, err := ngc.getGrpcClient().GetService(param)
	if err != nil {
		panic("GetService failed!")
	}
	fmt.Printf("GetService,param:%+v, result:%+v \n\n", param, service)

	return service, err
}
