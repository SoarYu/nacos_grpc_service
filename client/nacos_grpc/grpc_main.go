package nacos_grpc

type Domain struct {
	Name          string `json:"dom"`
	Clusters      string
	CacheMillis   int64
	LastRefMillis int64
	Instances     []Instance `json:"hosts"`
	Env           string
	TTL           int
}

type Instance struct {
	IP         string
	Port       int
	Weight     float64
	Valid      bool
	Unit       string
	AppUseType string
	Site       string
}

//var (
//	nacosClient naming_client.INamingClient
//)

//
//func initNacos() {
//	clientConfig := *constant.NewClientConfig(
//		constant.WithNamespaceId(""), //When namespace is public, fill in the blank string here.
//		constant.WithTimeoutMs(5000),
//		constant.WithNotLoadCacheAtStart(true),
//		constant.WithLogDir("/tmp/nacos/log"),
//		constant.WithCacheDir("/tmp/nacos/cache"),
//		constant.WithLogLevel("debug"),
//	)
//	//Another way of create serverConfigs
//	serverConfigs := []constant.ServerConfig{
//		*constant.NewServerConfig(
//			"192.168.66.146",
//			8848,
//			constant.WithScheme("http"),
//			constant.WithContextPath("/nacos"),
//		),
//	}
//	// Another way of create naming client for service discovery (recommend)
//	nc, err := clients.NewNamingClient(
//		vo.NacosClientParam{
//			ClientConfig:  &clientConfig,
//			ServerConfigs: serverConfigs,
//		},
//	)
//	if err != nil {
//		fmt.Println("init nacos-client error")
//	}
//	nacosClient = nc
//}
//
//func getAllService(c *gin.Context) {
//	param := vo.GetAllServiceInfoParam{
//		NameSpace: "",        //optional,default:public
//		GroupName: "group-a", //optional,default:DEFAULT_GROUP
//		PageNo:    1,
//		PageSize:  10,
//	}
//	service, err :=
//	//service, err := nacosClient.GetAllServicesInfo(param)
//	if err != nil {
//		panic("GetAllService failed!")
//	}
//	fmt.Printf("GetAllService,param:%+v, result:%+v \n\n", param, service)
//
//	c.JSON(200, service)
//
//}
//
//func GetService(c *gin.Context) {
//	param := vo.GetServiceParam{
//		ServiceName: "demo.go",
//		GroupName:   "group-a",
//		Clusters:    []string{"cluster-a"},
//	}
//	service, err := nacosClient.GetService(param)
//	if err != nil {
//		panic("GetService failed!")
//	}
//	fmt.Printf("GetService,param:%+v, result:%+v \n\n", param, service)
//
//	c.JSON(200, service)
//
//}
//
//func DeregisterInstance(c *gin.Context) {
//	param := vo.DeregisterInstanceParam{
//		Ip:          "192.168.66.148",
//		Port:        8801,
//		Cluster:     "cluster-a",
//		GroupName:   "group-a",
//		ServiceName: "demo.go",
//		Ephemeral:   true,
//	}
//	service, err := nacosClient.DeregisterInstance(param)
//	if err != nil {
//		panic("DeregisterInstance failed!")
//	}
//	fmt.Printf("DeregisterInstance,param:%+v, result:%+v \n\n", param, service)
//
//	c.JSON(200, service)
//
//}
//
//func main() {
//
//	initNacos()
//
//	InitGprcClient()
//
//	router := gin.Default()
//
//	router.GET("/srvIPXT", srvIPXT)
//	router.GET("/getAllService", getAllService)
//	router.GET("/getService", GetService)
//	router.GET("/deregisterInstance", DeregisterInstance)
//
//	router.Run(":8001")
//}
//
//func srvIPXT(c *gin.Context) {
//	var domain Domain
//	var ins1 Instance
//	var ins2 Instance
//
//	domain.Name = "demo.go"
//	domain.Clusters = "cluster-a"
//	domain.CacheMillis = 3000
//	domain.Env = ""
//
//	ins1.IP = "192.168.66.146"
//	ins1.Port = 8000
//	ins1.Weight = 10
//	ins1.Valid = true
//
//	ins2.IP = "192.168.66.146"
//	ins2.Port = 8001
//	ins2.Weight = 10
//	ins2.Valid = true
//
//	domain.Instances = append(domain.Instances, ins1, ins2)
//
//	data, _ := json.Marshal(domain)
//
//	fmt.Println(string(data))
//	c.JSON(200, domain)
//}
var nacosGrpcClient *GrpcClient

func main() {
	nacosGrpcClient = NewGrpcClient()
	//

	_, err := nacosGrpcClient.GetService()

	if err != nil {
		return
	}
}
