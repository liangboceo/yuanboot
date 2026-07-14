package console

import (
	"github.com/liangboceo/yuanboot/cli/yuanbootctl/docker"
	"github.com/liangboceo/yuanboot/cli/yuanbootctl/generate/projects"
	"github.com/liangboceo/yuanboot/cli/yuanbootctl/spec"
)

var Project = projects.NewEmptyProject("console", "Console Application").With(func(root *projects.ProjectItem) {
	root.AddFileWithContent("main.go", ProjectItem_Main_go)
	root.AddFileWithContent("startup.go", ProjectItem_startup_go)
	root.AddFileWithContent("hostservice.go", ProjectItem_hostservice_go)
	root.AddFileWithContent("config.yml", ProjectItem_conf_yml)
	root.AddFileWithContent("go.mod", ProjectItem_go_mod)
	root.AddDir("spec").AddFileWithContent("yuanboot.md", spec.YuanbootSpec)
	root.AddDir("spec").AddFileWithContent("backend-module-development.md", spec.BackendModuleDevelopment)
	root.AddFileWithContent("README.md", spec.ReadMe)
	root.AddDir("version").AddFileWithContent("version.go", docker.Version_Tel)
	root.AddFileWithContent("Dockerfile", docker.DockerFile_Tel)
	root.AddFileWithContent("docker-compose.yml", docker.DockerCompose_Tel)
	root.AddFileWithContent("build.sh", docker.BuildSh_Tel)
})
