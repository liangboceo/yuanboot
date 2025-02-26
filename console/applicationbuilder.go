package console

import "github.com/liangboceo/yuanboot/abstractions"

type ApplicationBuilder struct {
	hostBuilderContext *abstractions.HostBuilderContext
	extendConfigures   []func(context *abstractions.HostBuilderContext)
}

func (builder *ApplicationBuilder) Build() interface{} {
	builder.buildExtends()
	return builder
}

func (builder *ApplicationBuilder) SetHostBuildContext(context *abstractions.HostBuilderContext) {
	builder.hostBuilderContext = context
}

func NewApplicationBuilder() *ApplicationBuilder {
	return &ApplicationBuilder{}
}

func (builder *ApplicationBuilder) UseExtends(configure func(context *abstractions.HostBuilderContext)) *ApplicationBuilder {
	builder.extendConfigures = append(builder.extendConfigures, configure)
	return builder
}
func (builder *ApplicationBuilder) buildExtends() {
	for _, configure := range builder.extendConfigures {
		configure(builder.hostBuilderContext)
	}
}
