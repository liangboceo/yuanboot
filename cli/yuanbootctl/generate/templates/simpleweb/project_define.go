package simpleweb

import (
	"github.com/liangboceo/yuanboot/cli/yuanbootctl/docker"
	"github.com/liangboceo/yuanboot/cli/yuanbootctl/generate/projects"
	"github.com/liangboceo/yuanboot/cli/yuanbootctl/spec"
)

var Project = projects.NewEmptyProject("simpleweb", "Simple Web Application").With(func(root *projects.ProjectItem) {
	root.AddFileWithContent("main.go", Main_Tel)
	root.AddFileWithContent("config_dev.yml", Config_Tel)
	root.AddFileWithContent("go.mod", Mod_Tel)
	root.AddDir("static").AddDir("templates").AddFileWithContent("index.html", IndexHtml_Tel)
	root.AddFileWithContent("README.md", Readme_Tel)
	root.AddDir("spec").AddFileWithContent("yuanboot.md", spec.YuanbootSpec)
	root.AddDir("spec").AddFileWithContent("backend-module-development.md", spec.BackendModuleDevelopment)
	root.AddFileWithContent("README.md", spec.ReadMe)
	root.AddDir("version").AddFileWithContent("version.go", docker.Version_Tel)
	root.AddFileWithContent("Dockerfile", docker.DockerFile_Tel)
	root.AddFileWithContent("docker-compose.yml", docker.DockerCompose_Tel)
	root.AddFileWithContent("build.sh", docker.BuildSh_Tel)
})
