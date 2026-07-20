package templates

import (
	"github.com/liangboceo/yuanboot/cli/yuanbootctl/generate/templates/admin"
	"github.com/liangboceo/yuanboot/cli/yuanbootctl/generate/templates/console"
	"github.com/liangboceo/yuanboot/cli/yuanbootctl/generate/templates/grpc"
	"github.com/liangboceo/yuanboot/cli/yuanbootctl/generate/templates/iot"
	"github.com/liangboceo/yuanboot/cli/yuanbootctl/generate/templates/microservice"
	"github.com/liangboceo/yuanboot/cli/yuanbootctl/generate/templates/mvc"
	"github.com/liangboceo/yuanboot/cli/yuanbootctl/generate/templates/simpleweb"
	"github.com/liangboceo/yuanboot/cli/yuanbootctl/generate/templates/webapi"
	"github.com/liangboceo/yuanboot/cli/yuanbootctl/generate/templates/xxl_job"
)

func init() {
	registerProject("admin", admin.Project)
	registerProject("console", console.Project)
	registerProject("webapi", webapi.Project)
	registerProject("mvc", mvc.Project)
	registerProject("grpc", grpc.Project)
	registerProject("xxl-job", xxl_job.Project)
	registerProject("iot", iot.Project)
	registerProject("simpleweb", simpleweb.Project)
	registerProject("microservice", microservice.Project)
}
