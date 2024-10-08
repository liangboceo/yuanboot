package console

const ProjectItem_Main_go = `
package {{.CurrentModelName}}

import (
	"github.com/liangboceo/yuanboot/abstractions"
	"github.com/liangboceo/yuanboot/console"
)

func main() {
	// -f ./conf/test_conf.yml 指定配置文件 , 默认读取 config_{profile}.yml , -profile [dev,test,prod]
	config := abstractions.NewConfigurationBuilder().
		AddEnvironment().
		AddYamlFile("config").Build()

	console.NewHostBuilder().
		UseConfiguration(config).
		UseStartup(Startup).
		Build().
		Run()
}
`
