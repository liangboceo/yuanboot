package configuration

import (
	"embed"
	"github.com/liangboceo/yuanboot/abstractions"
)

func LocalConfig(configPath string) *abstractions.Configuration {
	return abstractions.NewConfigurationBuilder().AddEnvironment().AddYamlFile(configPath).Build()
}

func LocalConfigEmbed(configPath string, fs embed.FS) *abstractions.Configuration {
	return abstractions.NewConfigurationBuilder().AddEnvironment().AddEmbedFs(fs).AddYamlFile(configPath).Build()
}
