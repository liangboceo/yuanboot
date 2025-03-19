package nacos

import (
	"embed"
	nacos_viper_remote "github.com/liangboceo/nacos-viper-remote"
	"github.com/liangboceo/yuanboot/abstractions"
)

func AddRemoteWithNacos(builder *abstractions.ConfigurationBuilder) *abstractions.ConfigurationBuilder {
	if builder.Context.ConfigType == "" {
		builder.Context.ConfigType = "yml"
	}
	builder.Context.EnableRemote = true
	builder.Context.RemoteProvider = nacos_viper_remote.NewRemoteProvider(builder.Context.ConfigType)
	return builder
}

func RemoteConfig(configPath string) *abstractions.Configuration {
	return AddRemoteWithNacos(abstractions.NewConfigurationBuilder().AddEnvironment().AddYamlFile(configPath)).Build()
}

func RemoteConfigEmbed(configPath string, fs embed.FS) *abstractions.Configuration {
	return AddRemoteWithNacos(abstractions.NewConfigurationBuilder().AddEnvironment().AddEmbedFs(fs).AddYamlFile(configPath)).Build()
}
