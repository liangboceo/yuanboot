package admin

import (
	"embed"
	"io/fs"
	"path"
	"sort"
	"strings"

	"github.com/liangboceo/yuanboot/cli/yuanbootctl/generate/projects"
)

//go:embed assets
var assets embed.FS

var Project = projects.NewEmptyProject("admin", "基于 Yuanboot 和 PureAdmin 的管理后台项目").With(
	func(root *projects.ProjectItem) {
		addEmbeddedDir(root.AddDir("{{.ModelName}}-server"), "assets/project-server")
		addEmbeddedDir(root.AddDir("{{.ModelName}}-web"), "assets/project-web")
		root.AddDir("{{.ModelName}}-server").AddFileWithContent("go.mod", Mod_Tel)

	},
)

func addEmbeddedDir(parent *projects.ProjectItem, source string) {
	entries, err := fs.ReadDir(assets, source)
	if err != nil {
		panic(err)
	}
	sort.Slice(entries, func(i, j int) bool { return entries[i].Name() < entries[j].Name() })
	for _, entry := range entries {
		name := outputName(entry.Name())
		assetPath := path.Join(source, entry.Name())
		if entry.IsDir() {
			addEmbeddedDir(parent.AddDir(name), assetPath)
			continue
		}
		content, err := assets.ReadFile(assetPath)
		if err != nil {
			panic(err)
		}
		parent.AddFileWithContent(name, normalizeContent(string(content)))
	}
}

func outputName(name string) string {
	switch name {
	case "go.mod.tpl":
		return "go.mod"
	case "gitignore":
		return ".gitignore"
	}
	if strings.HasPrefix(name, "dot") {
		return "." + strings.TrimPrefix(name, "dot")
	}
	if strings.HasPrefix(name, "env.") {
		return "." + name
	}
	return name
}

func normalizeContent(content string) string {
	content = strings.ReplaceAll(content, "{{", "{{\"{{\"}}")
	content = strings.ReplaceAll(content, "sendex-agent-web", "{{.ModelName}}-web")
	content = strings.ReplaceAll(content, "sendex", "{{.ModelName}}")
	content = strings.ReplaceAll(content, "Sendex", "{{.ModelName}}")
	content = strings.ReplaceAll(content, "SendEx", "{{.ModelName}}")
	content = strings.ReplaceAll(content, "vue-pure-admin", "{{.ModelName}}-web")
	return content
}
