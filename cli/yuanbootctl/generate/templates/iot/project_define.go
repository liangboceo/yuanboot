package iot

import "github.com/liangboceo/yuanboot/cli/yuanbootctl/generate/projects"

var Project = projects.NewEmptyProject("iot", "Iot Application").With(func(root *projects.ProjectItem) {
	root.AddDir("controller").AddFileWithContent("demo.go", DemoController_Tel)
	root.AddDir("conf").AddFileWithContent("bootstrap_dev.yml", Config_Tel)
	root.AddDir("conf").AddFileWithContent("log.yml", Log_Config_Tel)
	root.AddDir("version").AddFileWithContent("version.go", Version_Tel)
	root.AddDir("pkg").AddDir("mq").AddFileWithContent("kafka_client.go", Kafka_Tel)
	root.AddDir("pkg").AddDir("service").AddFileWithContent("cache_service.go", Cache_Tel)
	root.AddDir("iothub").AddFileWithContent("iothub.go", IotHub_Tel)
	root.AddFileWithContent("go.mod", Mod_Tel)
	root.AddFileWithContent("main.go", Main_Tel)
})
