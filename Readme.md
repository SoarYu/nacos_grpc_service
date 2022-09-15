# gRPC+protobuf搭建的简单服务端注册到nacos中,并进行负载均衡测试

## protobuf生成的简单服务端
```
    #生成protobuf文件
    cd pb
    protoc --go_out=plugins=grpc:./ *.proto

    #运行命令安装包
    go mod init nacos-grpc-service
    go mod tidy

```

## 开启 ProtoBuf 生成的 gRPC服务 并将服务注册到 Nacos 服务器中。
```
    go run main.go
```

## 查看nacos服务注册的状态
![](img/nacos_grpc_service_list.png)

## 查看“服务详情” sdk获取服务信息
![](img/nacos_grpc_service_detail.png)

[代码参考](https://gitee.com/xiaonqedu/nacos-grpc-gotest?_from=gitee_search)




