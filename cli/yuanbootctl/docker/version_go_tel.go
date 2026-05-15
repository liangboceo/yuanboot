package docker

const Version_Tel = `
package {{.CurrentModelName}}

const version = "1.0.0"

var (
	env = "dev"
)

// Version return the version string
func Version() string {
	return version
}

// Env return the env string
func Env() string {
	return env
}

`
