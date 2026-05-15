package microservice

import "github.com/liangboceo/yuanboot/cli/yuanbootctl/generate/projects"

var Project = projects.NewEmptyProject("microservice", "Complete Microservice Application").With(func(root *projects.ProjectItem) {
	root.AddDir("cmd").AddFileWithContent("main.go", Main_Tel)
	root.AddDir("internal").AddDir("controller").AddFileWithContent("user_controller.go", UserController_Tel)
	root.AddDir("internal").AddDir("service").AddFileWithContent("user_service.go", UserService_Tel)
	root.AddDir("internal").AddDir("repository").AddFileWithContent("user_repository.go", UserRepository_Tel)
	root.AddDir("internal").AddDir("model").AddFileWithContent("user.go", UserModel_Tel)
	root.AddDir("internal").AddDir("middleware").AddFileWithContent("auth.go", Middleware_Tel)
	root.AddDir("config").AddFileWithContent("config_dev.yml", ConfigDev_Tel)
	root.AddDir("config").AddFileWithContent("config_prod.yml", ConfigProd_Tel)
	root.AddFileWithContent("go.mod", Mod_Tel)
	root.AddFileWithContent("Makefile", Makefile_Tel)
	root.AddFileWithContent("README.md", Readme_Tel)
})
