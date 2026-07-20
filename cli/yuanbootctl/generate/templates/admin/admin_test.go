package admin

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestProjectGenerate(t *testing.T) {
	target := t.TempDir()
	Project.Generate(target, "example")

	checks := []string{
		"example/example-server/main.go",
		"example/example-server/go.mod",
		"example/example-server/conf/config_dev.yml",
		"example/example-server/internal/controller/sys_controller.go",
		"example/example-server/internal/service/system_menu_seed.go",
		"example/example-web/package.json",
		"example/example-web/src/layout/index.vue",
		"example/example-web/src/views/login/index.vue",
		"example/example-web/src/views/system/user/index.vue",
		"example/example-web/src/views/system/role/index.vue",
		"example/example-web/src/views/system/menu/index.vue",
		"example/example-web/src/views/system/system-config/index.vue",
		"example/example-web/.gitignore",
	}
	for _, name := range checks {
		if _, err := os.Stat(filepath.Join(target, name)); err != nil {
			t.Fatalf("expected generated file %s: %v", name, err)
		}
	}

	goMod, err := os.ReadFile(filepath.Join(target, "example/example-server/go.mod"))
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(goMod), "module example-server") {
		t.Fatalf("module name was not rendered: %s", goMod)
	}

	controller, err := os.ReadFile(filepath.Join(target, "example/example-server/internal/controller/sys_controller.go"))
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(string(controller), "sendex-server") {
		t.Fatal("generated server still contains the source module name")
	}
}
