package xxl_job

import "github.com/liangboceo/yuanboot/cli/yuanbootctl/generate/projects"

var Project = projects.NewEmptyProject("xxl-job", "xxl-job Console Application").With(func(root *projects.ProjectItem) {
	root.AddFileWithContent("demojob.go", Demo_Job_Tel)
	root.AddFileWithContent("config_dev.yml", Config_Tel)
	root.AddFileWithContent("go.mod", Mod_Tel)
	root.AddFileWithContent("main.go", Main_Tel)
})
