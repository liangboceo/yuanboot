package main

import (
	"github.com/liangboceo/dependencyinjection"
	"github.com/liangboceo/yuanboot/abstractions"
	yrpc "github.com/liangboceo/yuanboot/grpc"
	"github.com/liangboceo/yuanboot/pkg/servicediscovery/nacos"
	"google.golang.org/grpc"
	pb "grpc-demo/proto/helloworld"
	"grpc-demo/services"
)

func main() {
	configuration := abstractions.NewConfigurationBuilder().
		AddEnvironment().
		AddYamlFile("config").Build()

	hosting := yrpc.NewHostBuilder().
		UseConfiguration(configuration).
		Configure(func(app *yrpc.ApplicationBuilder) {
			//app.AddUnaryServerInterceptor( logger.UnaryServerInterceptor() )
			//app.AddStreamServerInterceptor( logger.StreamServerInterceptor() )
			app.AddGrpcService(func(server *grpc.Server, ctx *yrpc.ServiceContext) {
				ctx.Register(pb.RegisterGreeterServer) // register grpc service
			})

		}).
		ConfigureServices(func(collection *dependencyinjection.ServiceCollection) {
			collection.AddSingleton(services.NewGreeterServer) // add grpc service
			collection.AddSingleton(services.NewIOCDemo)
			nacos.UseServiceDiscovery(collection)
		}).Build()

	hosting.Run()
}
