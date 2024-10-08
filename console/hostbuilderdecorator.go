package console

import (
	"github.com/liangboceo/dependencyinjection"
	"github.com/liangboceo/yuanboot/abstractions"
)

type HostBuilderDecorator struct {
}

func (h HostBuilderDecorator) OverrideConfigure(configureFunc interface{}, builder abstractions.IApplicationBuilder) {
	configureWebAppFunc := configureFunc.(func(applicationBuilder *ApplicationBuilder))
	configureWebAppFunc(builder.(*ApplicationBuilder))
}

func (h HostBuilderDecorator) OverrideNewApplicationBuilder(context *abstractions.HostBuilderContext) abstractions.IApplicationBuilder {
	applicationBuilder := NewApplicationBuilder()
	applicationBuilder.SetHostBuildContext(context)
	return applicationBuilder
}

func (h HostBuilderDecorator) OverrideNewHost(server abstractions.IServer, context *abstractions.HostBuilderContext) abstractions.IServiceHost {
	return NewHost(NewServer(), context)
}

func (h HostBuilderDecorator) OverrideIOCInnerConfigures(serviceCollection *dependencyinjection.ServiceCollection) {
	//panic("implement me")
}

func NewHostBuilderDecorator() HostBuilderDecorator {
	return HostBuilderDecorator{}
}
