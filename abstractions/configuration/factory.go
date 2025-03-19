package configuration

import (
	"embed"
	"github.com/liangboceo/yuanboot/abstractions"
)

// YAML config by yaml or yml file
func YAML(configPath string) *abstractions.Configuration {
	return abstractions.NewConfigurationBuilder().AddEnvironment().AddYamlFile(configPath).Build()
}
func YAML_EMBED(configPath string, fs embed.FS) *abstractions.Configuration {
	return abstractions.NewConfigurationBuilder().AddEnvironment().AddEmbedFs(fs).AddYamlFile(configPath).Build()
}

// JSON config by json file
func JSON(configPath string) *abstractions.Configuration {
	return abstractions.NewConfigurationBuilder().AddEnvironment().AddJsonFile(configPath).Build()
}

// JSON config by json file
func JSON_EMBED(configPath string, fs embed.FS) *abstractions.Configuration {
	return abstractions.NewConfigurationBuilder().AddEmbedFs(fs).AddEnvironment().AddJsonFile(configPath).Build()
}

// PROPERITES config by properties file
func PROPERITES(configPath string) *abstractions.Configuration {
	return abstractions.NewConfigurationBuilder().AddEnvironment().AddPropertiesFile(configPath).Build()
}

func PROPERITES_EMBED(configPath string, fs embed.FS) *abstractions.Configuration {
	return abstractions.NewConfigurationBuilder().AddEmbedFs(fs).AddEnvironment().AddPropertiesFile(configPath).Build()
}
