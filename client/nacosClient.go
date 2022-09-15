package client

import (
	"fmt"
	"github.com/nacos-group/nacos-sdk-go/v2/clients"
	"github.com/nacos-group/nacos-sdk-go/v2/clients/naming_client"
	"github.com/nacos-group/nacos-sdk-go/v2/common/constant"
	"github.com/nacos-group/nacos-sdk-go/v2/model"
	"github.com/nacos-group/nacos-sdk-go/v2/vo"
	"strconv"
	"strings"
)

var NacosClient naming_client.INamingClient

func InitNacosClient(serverHosts []string) {
	// 创建clientConfig
	clientConfig := constant.ClientConfig{
		NamespaceId:         "", // 如果需要支持多namespace，我们可以场景多个client,它们有不同的NamespaceId。当namespace是public时，此处填空字符串。
		TimeoutMs:           5000,
		NotLoadCacheAtStart: true,
		LogDir:              "/tmp/nacos/log",
		CacheDir:            "/tmp/nacos/cache",
		LogLevel:            "debug",
	}

	// 至少一个ServerConfig
	serverConfigs := make([]constant.ServerConfig, len(serverHosts))

	for i, serverHost := range serverHosts {
		serverIp := strings.Split(serverHost, ":")[0]
		serverPort, err := strconv.Atoi(strings.Split(serverHost, ":")[1])
		if err != nil {
			fmt.Errorf("nacos server host config error!", err)
		}
		serverConfigs[i] = *constant.NewServerConfig(
			serverIp,
			uint64(serverPort),
			constant.WithScheme("http"),
			constant.WithContextPath("/nacos"),
		)
	}

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
		fmt.Println("服务： " + serviceName + " 注册失败！")
		return
	} else {
		fmt.Println("服务： " + serviceName + " 注册成功！")
	}
}
func SelectAllInstances(serviceName string) ([]model.Instance, error) {
	return NacosClient.SelectAllInstances(vo.SelectAllInstancesParam{
		ServiceName: serviceName,
	})
	//return instances
}
func SelectOneHealthyInstance(serviceName string) (*model.Instance, error) {
	return NacosClient.SelectOneHealthyInstance(vo.SelectOneHealthInstanceParam{
		ServiceName: serviceName,
	})
}
