package main

import (
	"fmt"
	"github.com/lgc202/mall-micro/common/registry"
	"github.com/lgc202/mall-micro/service/user/rpc/global"
	"github.com/lgc202/mall-micro/service/user/rpc/handler"
	"github.com/lgc202/mall-micro/service/user/rpc/initialize"
	"github.com/lgc202/mall-micro/service/user/rpc/proto"
	uuid "github.com/satori/go.uuid"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
	"net"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	// 初始化
	initialize.InitLogger()
	initialize.InitConfig()
	initialize.InitDb()

	zap.S().Info("host: ", global.ServerConfig.Host)
	zap.S().Info("port: ", global.ServerConfig.Port)

	server := grpc.NewServer()
	proto.RegisterUserServer(server, &handler.UserServer{})
	listen, err := net.Listen("tcp", fmt.Sprintf("%s:%d", global.ServerConfig.Host, global.ServerConfig.Port))
	if err != nil {
		panic("failed to listen:" + err.Error())
	}

	// 注册grpc服务健康检查
	grpc_health_v1.RegisterHealthServer(server, health.NewServer())

	// 服务注册
	serviceID := fmt.Sprintf("%s", uuid.NewV4())
	consulClient, err := registry.NewRegisry(global.ServerConfig.ConsulInfo.Host, global.ServerConfig.ConsulInfo.Port)
	if err != nil {
		panic(err)
	}

	if err := consulClient.Register(
		global.ServerConfig.Host, serviceID,
		global.ServerConfig.Name, []string{global.ServerConfig.Name},
		global.ServerConfig.Port, registry.GRPC); err != nil {
		panic(err)
	}

	go func() {
		err = server.Serve(listen)
		if err != nil {
			panic("failed to start grpc:" + err.Error())
		}
	}()

	// 接收终止信号
	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	if err = consulClient.DeRegister(serviceID); err != nil {
		zap.S().Info("注销失败")
		panic(err)
	}
	zap.S().Infof("注销成功")
}
