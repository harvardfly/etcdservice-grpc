package main

import (
	"etcdservice-grpc/config"
	"etcdservice-grpc/etcdservice"
	pb "etcdservice-grpc/protos"
	"fmt"
	"google.golang.org/grpc/balancer/roundrobin"
	"google.golang.org/grpc/resolver"
	"time"

	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

func main() {
	err := config.InitFromIni("conf/conf.ini")
	if err != nil {
		panic(err)
	}

	// 解析etcd服务地址
	r := etcdservice.NewResolver(config.Conf.EtcdAddr, config.Conf.Schema)
	resolver.Register(r)

	// 客户端连接服务器
	conn, err := grpc.Dial(config.Conf.Schema+"://"+config.Conf.Caller+"/"+config.Conf.Callee, grpc.WithDefaultServiceConfig(fmt.Sprintf(`{"LoadBalancingPolicy": "%s"}`, roundrobin.Name)), grpc.WithInsecure())
	if err != nil {
		fmt.Println("连接服务器失败", err)
	}
	defer conn.Close()

	// 获得grpc句柄
	c := pb.NewHelloServerClient(conn)
	ticker := time.NewTicker(1 * time.Second)
	for t := range ticker.C {
		// 远程单调用 SayHi 接口
		r1, err := c.SayHi(context.Background(), &pb.HelloRequest{Name: "aaa"})
		if err != nil {
			fmt.Println("Can not get SayHi:", err)
			return
		}
		fmt.Printf("%v: SayHi 响应：%s\n", t, r1.Message)
	}
}
