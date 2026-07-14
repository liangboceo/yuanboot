package grpc

import (
	"github.com/liangboceo/yuanboot/cli/yuanbootctl/docker"
	"github.com/liangboceo/yuanboot/cli/yuanbootctl/generate/projects"
	"github.com/liangboceo/yuanboot/cli/yuanbootctl/spec"
)

var Project = projects.NewEmptyProject("grpc", "Grpc Application").With(func(root *projects.ProjectItem) {
	clientDir := root.AddDir("client")
	clientDir.AddFileWithContent("api.go", Client_Api_Tel)
	clientDir.AddFileWithContent("clientservice.go", Client_Service_Tel)
	clientDir.AddFileWithContent("main.go", Client_Main_Tel)
	clientDir.AddFileWithContent("config.yml", CLient_Config_Tel)
	protoDir := root.AddDir("proto")
	protoDir.AddFileWithContent("helloworld.proto", Hello_World_Tel)
	protoDir.AddDir("helloworld").AddFileWithContent("helloworld.pd.go", Hello_World_PD_TEL)
	serviceDir := root.AddDir("services")
	serviceDir.AddFileWithContent("demo.go", Demo_Tel)
	serviceDir.AddFileWithContent("greeterservice.go", Greeter_Server_Tel)
	root.AddFileWithContent("config.yml", ServiceConfig_Tel)
	root.AddFileWithContent("go.mod", Mod_Tel)
	root.AddFileWithContent("main.go", Main_Tel)
	root.AddDir("spec").AddFileWithContent("yuanboot.md", spec.YuanbootSpec)
	root.AddDir("spec").AddFileWithContent("backend-module-development.md", spec.BackendModuleDevelopment)
	root.AddFileWithContent("README.md", spec.ReadMe)
	root.AddDir("version").AddFileWithContent("version.go", docker.Version_Tel)
	root.AddFileWithContent("Dockerfile", docker.DockerFile_Tel)
	root.AddFileWithContent("docker-compose.yml", docker.DockerCompose_Tel)
	root.AddFileWithContent("build.sh", docker.BuildSh_Tel)
})
