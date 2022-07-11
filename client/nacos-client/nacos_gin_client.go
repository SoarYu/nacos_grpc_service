package main

import (
	"nacos_person/client/nacos_grpc"
)

var nacosGrpcClient *nacos_grpc.GrpcClient

func main() {
	nacosGrpcClient = nacos_grpc.NewGrpcClient()
	_, err := nacosGrpcClient.GetService()

	if err != nil {
		return
	}
}
