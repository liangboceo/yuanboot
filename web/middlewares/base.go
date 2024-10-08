package middlewares

import "github.com/liangboceo/yuanboot/abstractions"

type IConfigurationMiddleware interface {
	SetConfiguration(config abstractions.IConfiguration)
}

type BaseMiddleware struct {
	// Configuration
	config abstractions.IConfiguration
}

func (mdw *BaseMiddleware) SetConfiguration(config abstractions.IConfiguration) {
	mdw.config = config
}
