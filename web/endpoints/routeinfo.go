package endpoints

import (
	"fmt"
	"github.com/liangboceo/yuanboot/abstractions"
	"github.com/liangboceo/yuanboot/abstractions/xlog"
	"github.com/liangboceo/yuanboot/web/context"
	"github.com/liangboceo/yuanboot/web/mvc"
	"github.com/liangboceo/yuanboot/web/router"
	"strings"
	"sync"
)

var (
	once           sync.Once
	routerInfoList []router.Info
)

func UseRouteInfo(route router.IRouterBuilder) {
	xlog.GetXLogger("Endpoint").Debug("loaded router information endpoint.")

	route.GET("/actuator/routers", func(ctx *context.HttpContext) {
		once.Do(func() {
			var builder *mvc.ControllerBuilder
			var env *abstractions.HostEnvironment
			_ = ctx.RequiredServices.GetService(&builder)
			_ = ctx.RequiredServices.GetService(&env)
			serverPath := env.MetaData["server.path"]
			mvcTemplate := env.MetaData["mvc.template"]
			// route
			routerInfoList = make([]router.Info, len(route.GetRouteInfo()))
			copy(routerInfoList, route.GetRouteInfo())
			for idx, _ := range routerInfoList {
				if serverPath == "" {
					routerInfoList[idx].Path = fmt.Sprintf("/%s", strings.TrimLeft(routerInfoList[idx].Path, "/"))
				} else {
					routerInfoList[idx].Path = fmt.Sprintf("/%s/%s", strings.TrimRight(serverPath, "/"), strings.TrimLeft(routerInfoList[idx].Path, "/"))
				}
			}
			// mvc
			mvcTemplate = strings.ReplaceAll(mvcTemplate, "{controller}", "%s")
			mvcTemplate = strings.ReplaceAll(mvcTemplate, "{action}", "%s")
			if serverPath == "" {
				mvcTemplate = fmt.Sprintf("/%s", strings.TrimLeft(mvcTemplate, "/"))
			} else {
				mvcTemplate = fmt.Sprintf("/%s/%s", strings.TrimRight(serverPath, "/"), strings.TrimLeft(mvcTemplate, "/"))
			}
			descriptorList := builder.GetControllerDescriptorList()
			for _, desc := range descriptorList {
				for _, action := range desc.GetActionDescriptors() {
					colName := strings.ReplaceAll(desc.ControllerName, "controller", "")
					actionName := getActionPath(action.ActionName)

					routerInfoList = append(routerInfoList, router.Info{Method: strings.ToUpper(action.ActionMethod), Path: fmt.Sprintf(mvcTemplate, colName, actionName), Type: "mvc"})
				}
			}
		})

		ctx.JSON(200, routerInfoList)
	})
}

func getActionPath(actionName string) string {
	name := strings.ToLower(actionName)
	name = strings.Replace(name, "get", "", 5)
	name = strings.Replace(name, "post", "", 5)
	name = strings.Replace(name, "put", "", 5)
	name = strings.Replace(name, "delete", "", 5)
	return name
}
