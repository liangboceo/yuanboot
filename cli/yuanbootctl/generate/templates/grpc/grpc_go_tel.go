package grpc

const Main_Tel = `
package main

import (
	"github.com/liangboceo/yuanboot/abstractions"
	"github.com/liangboceo/dependencyinjection"
	yrpc "github.com/liangboceo/yuanboot/grpc"
	"github.com/liangboceo/yuanboot/pkg/servicediscovery/nacos"
	"google.golang.org/grpc"
	pb "{{.ModelName}}/proto/helloworld"
	"{{.ModelName}}/services"
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

`

const Mod_Tel = `
module {{.ModelName}}

go 1.16

require (
	github.com/liangboceo/dependencyinjection v1.0.0
	github.com/liangboceo/yuanboot {{.Version}}
	github.com/golang/protobuf v1.5.2
	github.com/grpc-ecosystem/go-grpc-middleware v1.3.0
	google.golang.org/grpc v1.38.0
	google.golang.org/protobuf v1.26.0
)
`
