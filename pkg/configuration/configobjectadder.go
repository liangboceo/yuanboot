package configuration

import (
	"github.com/liangboceo/dependencyinjection"
	"github.com/liangboceo/reflectx"
	"github.com/liangboceo/yuanboot/abstractions"
	"reflect"
)

// AddConfiguration 注入函数 用户API
func AddConfiguration(sc *dependencyinjection.ServiceCollection, objType interface{}) {
	sc.AddTransient(objType)
}

func Configure[T any](services *dependencyinjection.ServiceCollection) {
	services.AddSingleton(configTypeFactory[T])
}

func configTypeFactory[T any](configuration abstractions.IConfiguration) OptionsSnapshot[T] {
	var config T
	_, objectType := reflectx.GetCtorFuncOutTypeName(config)
	configObject := reflect.New(objectType).Interface().(abstractions.IConfigurationProperties)
	sectionName := configObject.GetSection()
	return OptionsSnapshot[T]{config: configuration, sectionName: sectionName, value: config}
}
