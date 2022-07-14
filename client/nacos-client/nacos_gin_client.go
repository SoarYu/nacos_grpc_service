package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/nacos-group/nacos-sdk-go/v2/vo"
	"nacos-grpc-gin/client/nacos_grpc"
	"net/http"
)

var nacosGrpcClient *nacos_grpc.GrpcClient

func GetService(c *gin.Context) {
	serviceName, _ := c.GetQuery("serviceName")
	clusters, _ := c.GetQueryArray("clusters")
	service, err := nacosGrpcClient.GetService(serviceName, clusters)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status": "Error",
		})
		return
	}
	c.JSON(http.StatusOK, service)
}

func GetAllServicesInfo(c *gin.Context) {
	service, err := nacosGrpcClient.GetAllServicesInfo()

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status": "Error",
		})
		return
	}
	c.JSON(http.StatusOK, service)
}

func Subscribe(c *gin.Context) {
	serviceName, ok := c.GetQuery("serviceName")
	//err := nacosGrpcClient.Subscribe("", "")
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{
			//"serviceName": serviceName,
			"ok": ok,
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"serviceName": serviceName,
		"ok":          ok,
	})
}

func UpdateServiceInstance(c *gin.Context) {
	//serviceName, _ := c.GetQuery("serviceName")
	ip, _ := c.GetQuery("ip")
	//metadata[idc]=shanghai
	metadata, _ := c.GetQueryMap("metadata")

	fmt.Println("getQuery:", ip, metadata)
	param := vo.UpdateInstanceParam{
		Ip:          ip, //update ip
		Port:        8461,
		ServiceName: "demo.go",
		GroupName:   "group-demo",
		ClusterName: "cluster-146",
		Weight:      10,
		Enable:      true,
		Healthy:     true,
		Ephemeral:   true,
		Metadata:    metadata, //update metadata
	}
	success, err := nacosGrpcClient.UpdateServiceInstance(param)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status": "Error",
		})
		return
	}
	c.JSON(http.StatusOK, success)
}

//SelectAllInstance
//GroupName=DEFAULT_GROUP
func SelectAllInstances(c *gin.Context) {
	serviceName, _ := c.GetQuery("serviceName")
	groupName, _ := c.GetQuery("groupName")
	//clusters=cluster-146,cluster-148
	//clusters, _ := c.GetQueryArray("clusters")
	param := vo.SelectAllInstancesParam{
		ServiceName: serviceName,
		GroupName:   groupName,
		//Clusters:    clusters,
	}
	instances, err := nacosGrpcClient.SelectAllInstances(param)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status": "Error",
		})
		return
	}
	c.JSON(http.StatusOK, instances)
}

//SelectInstances only return the instances of healthy=${HealthyOnly},enable=true and weight>0
//ClusterName=DEFAULT,GroupName=DEFAULT_GROUP
func SelectInstances(c *gin.Context) {
	serviceName, _ := c.GetQuery("serviceName")
	groupName, _ := c.GetQuery("groupName")
	//clusters=cluster-146,cluster-148
	//clusters, _ := c.GetQueryArray("clusters")
	param := vo.SelectInstancesParam{
		ServiceName: serviceName,
		GroupName:   groupName,
		//Clusters:    clusters,
	}
	instances, err := nacosGrpcClient.SelectInstances(param)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status": "Error",
		})
		return
	}
	c.JSON(http.StatusOK, instances)
}

//SelectOneHealthyInstance return one instance by WRR strategy for load balance
//And the instance should be health=true,enable=true and weight>0
//ClusterName=DEFAULT,GroupName=DEFAULT_GROUP
func SelectOneHealthInstance(c *gin.Context) {
	serviceName, _ := c.GetQuery("serviceName")
	groupName, ok := c.GetQuery("groupName")
	if !ok {
		groupName = ""
	}
	//clusters=cluster-146,cluster-148
	//clusters, _ := c.GetQueryArray("clusters")
	param := vo.SelectOneHealthInstanceParam{
		ServiceName: serviceName,
		GroupName:   groupName,
		//Clusters:    clusters,
	}
	instance, err := nacosGrpcClient.SelectOneHealthInstance(param)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status": "Error",
		})
		return
	}
	c.JSON(http.StatusOK, instance)
}

func main() {
	nacosGrpcClient = nacos_grpc.NewGrpcClient("", "47.115.216.190")

	router := gin.Default()

	router.GET("/getService", GetService)
	router.GET("/getAllServicesInfo", GetAllServicesInfo)
	router.GET("/subscribe", Subscribe)
	router.GET("/updateServiceInstance", UpdateServiceInstance)
	router.GET("/selectOneHealthInstance", SelectOneHealthInstance)
	router.GET("/selectInstances", SelectInstances)
	router.GET("/selectAllInstances", SelectAllInstances)
	router.GET("/allDomNames", AllDomNames)
	router.GET("/srvIPXT", SrvIPXT)

	router.Run(":8000")
}

func AllDomNames(c *gin.Context) {
	services, err := nacosGrpcClient.GetAllServicesInfo()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status": "Error",
		})
		return
	}
	c.JSON(http.StatusOK, services)
}

func SrvIPXT(c *gin.Context) {
	serviceName, _ := c.GetQuery("serviceName")
	clusters, _ := c.GetQueryArray("clusters")
	service, err := nacosGrpcClient.GetService(serviceName, clusters)
	if err != nil {
		c.JSON(http.StatusBadRequest, service)
	}
	c.JSON(http.StatusOK, service)
}
