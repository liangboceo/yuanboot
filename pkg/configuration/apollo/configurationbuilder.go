package apollo

import (
	"embed"
	"github.com/liangboceo/yuanboot/abstractions"
)

func AddRemoteWithApollo(builder *abstractions.ConfigurationBuilder) *abstractions.ConfigurationBuilder {
	if builder.Context.ConfigType == "" {
		builder.Context.ConfigType = "yml"
	}
	builder.Context.EnableRemote = true
	builder.Context.RemoteProvider = NewRemoteProvider(builder.Context.ConfigType)
	return builder
}

func RemoteConfig(configPath string) *abstractions.Configuration {
	return AddRemoteWithApollo(abstractions.NewConfigurationBuilder().AddEnvironment().AddYamlFile(configPath)).Build()
}

func RemoteConfigEmbed(configPath string, fs embed.FS) *abstractions.Configuration {
	return AddRemoteWithApollo(abstractions.NewConfigurationBuilder().AddEnvironment().AddEmbedFs(fs).AddYamlFile(configPath)).Build()
}
