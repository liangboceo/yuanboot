package admin

const Mod_Tel = `
module {{.ModelName}}-server

go 1.18

require (
	github.com/liangboceo/dependencyinjection v1.0.0
	github.com/liangboceo/yuanboot {{.Version}}
)

require (
	github.com/valyala/fasthttp v1.51.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)
`
