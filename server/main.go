package main

import (
	"context"
	"etcdservice-grpc/config"
	"etcdservice-grpc/etcdservice"
	pb "etcdservice-grpc/protos"
	"fmt"
	"google.golang.org/grpc"
	"net"
	"os"
	"os/signal"
	"syscall"
)

// EtcdGRPCServer rpc服务绑定的struct
type EtcdGRPCServer struct {
}

// SayHi 简单rpc示例
func (s *EtcdGRPCServer) SayHi(ctx context.Context, req *pb.HelloRequest) (*pb.HelloResponse, error) {
	msg := "你好：" + req.Name
	return &pb.HelloResponse{
		Message: msg,
	}, nil
}

func main() {
	err := config.InitFromIni("conf/conf.ini")
	if err != nil {
		panic(err)
	}
	addr := fmt.Sprintf(":%d", config.Conf.Port)

	// 监听网络
	ln, err := net.Listen("tcp", addr)
	if err != nil {
		fmt.Println("网络异常：", err)
		return
	}
	defer ln.Close()

	// 创建grpc句柄
	srv := grpc.NewServer()
	defer srv.GracefulStop()

	// 将server结构体注册到grpc服务中
	pb.RegisterHelloServerServer(srv, &EtcdGRPCServer{})

	// etcd服务注册
	go etcdservice.Register(config.Conf.EtcdAddr, config.Conf.ServiceName, addr, config.Conf.Schema, config.Conf.TTL)

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGTERM, syscall.SIGINT, syscall.SIGKILL, syscall.SIGHUP, syscall.SIGQUIT)
	go func() {
		s := <-ch
		etcdservice.UnRegister(config.Conf.ServiceName, addr, config.Conf.Schema)

		if i, ok := s.(syscall.Signal); ok {
			os.Exit(int(i))
		} else {
			os.Exit(0)
		}
	}()

	// 监听服务
	err = srv.Serve(ln)
	if err != nil {
		fmt.Println("监听异常：", err)
		return
	}
}
