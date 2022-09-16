package client

import (
	"fmt"
	"github.com/nacos-group/nacos-sdk-go/clients"
	"github.com/nacos-group/nacos-sdk-go/clients/naming_client"
	"github.com/nacos-group/nacos-sdk-go/common/constant"
	"github.com/nacos-group/nacos-sdk-go/vo"
	"gopkg.in/natefinch/lumberjack.v2"
)

var NacosClient naming_client.INamingClient

func InitNacosClient(serverHosts []string) {
	//or a more graceful way to c/reate ServerConfig
	sc := []constant.ServerConfig{
		*constant.NewServerConfig("106.52.77.111", 8848),
	}

	//or a more graceful way to create ClientConfig
	cc := *constant.NewClientConfig(
		constant.WithNamespaceId(""),
		constant.WithTimeoutMs(5000),
		constant.WithNotLoadCacheAtStart(true),
		constant.WithLogDir("/tmp/nacos/log"),
		constant.WithCacheDir("/tmp/nacos/cache"),
		constant.WithLogLevel("debug"),
		constant.WithLogRollingConfig(&lumberjack.Logger{MaxSize: 10}),
		constant.WithLogStdout(true),
	)

	// a more graceful way to create naming client
	var err error
	NacosClient, err = clients.NewNamingClient(
		vo.NacosClientParam{
			ClientConfig:  &cc,
			ServerConfigs: sc,
		},
	)

	if err != nil {
		panic(err)
	}
}

func DeRegisterNacos(serviceName string, serviceHost string, servicePort uint64) {
	success, _ := NacosClient.DeregisterInstance(vo.DeregisterInstanceParam{
		Ip:          "10.0.0.10",
		Port:        8848,
		ServiceName: "demo.go",
		Ephemeral:   true, //it must be true
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
	})
	if !success {
		fmt.Println("服务： " + serviceName + " 注册失败！")
		return
	} else {
		fmt.Println("服务： " + serviceName + " 注册成功！")
	}
}

//func SelectAllInstances(serviceName string) ([]model.Instance, error) {
//	return NacosClient.SelectAllInstances(vo.SelectAllInstancesParam{
//		ServiceName: serviceName,
//	})
//	//return instances
//}
//func SelectOneHealthyInstance(serviceName string) (*model.Instance, error) {
//	return NacosClient.SelectOneHealthyInstance(vo.SelectOneHealthInstanceParam{
//		ServiceName: serviceName,
//	})
//}
