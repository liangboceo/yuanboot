package simpleweb

import "github.com/liangboceo/yuanboot/cli/yuanbootctl/generate/projects"

var Project = projects.NewEmptyProject("simpleweb", "Simple Web Application").With(func(root *projects.ProjectItem) {
	root.AddFileWithContent("main.go", Main_Tel)
	root.AddFileWithContent("config_dev.yml", Config_Tel)
	root.AddFileWithContent("go.mod", Mod_Tel)
	root.AddDir("static").AddDir("templates").AddFileWithContent("index.html", IndexHtml_Tel)
	root.AddFileWithContent("README.md", Readme_Tel)
})
