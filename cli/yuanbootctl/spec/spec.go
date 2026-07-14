package spec

import _ "embed"

//go:embed spec/yuanboot.md
var YuanbootSpec string

//go:embed README.md
var ReadMe string

//go:embed spec/backend-module-development.md
var BackendModuleDevelopment string
