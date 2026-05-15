package cmds

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/liangboceo/yuanboot/cli/yuanbootctl/utils"
	"github.com/spf13/cobra"
)

// 支持的操作系统和架构
var supportedPlatforms = []struct {
	GOOS   string
	GOARCH string
}{
	{"linux", "amd64"},
	{"linux", "arm64"},
	{"darwin", "amd64"},
	{"darwin", "arm64"},
	{"windows", "amd64"},
	{"windows", "386"},
}

var buildOS string
var buildArch string

func init() {
	BuildCmd.Flags().StringVarP(&buildOS, "os", "o", "", "Target OS (linux, darwin, windows, all)")
	BuildCmd.Flags().StringVarP(&buildArch, "arch", "a", "", "Target ARCH (amd64, arm64, 386, all)")
}

var BuildCmd = &cobra.Command{
	Use:   "build",
	Short: "build Project application of yuanboot fx",
	Long:  `build Project application of yuanboot fx. Use --os and --arch to specify target platforms, or use --os all to build for all supported platforms.`,
	Run: func(cmd *cobra.Command, args []string) {
		buildProject()
	},
}

func buildProject() {
	targets := resolveTargets()

	for _, target := range targets {
		buildForPlatform(target.GOOS, target.GOARCH)
	}

	// 移动静态文件
	utils.CopyPath("config"+utils.DirDot(), "build"+utils.DirDot()+"config")
	fmt.Println("build success")
}

// resolveTargets 根据入参确定编译目标平台
func resolveTargets() []struct {
	GOOS   string
	GOARCH string
} {
	// 都不传参数，编译当前环境
	if buildOS == "" && buildArch == "" {
		return []struct {
			GOOS   string
			GOARCH string
		}{{runtime.GOOS, runtime.GOARCH}}
	}

	// os=all，编译所有支持的环境
	if strings.ToLower(buildOS) == "all" {
		return supportedPlatforms
	}

	// 单独指定 os 或 arch
	goos := buildOS
	goarch := buildArch

	// 如果只指定了 os，使用该 os 的所有架构
	if goarch == "" || strings.ToLower(goarch) == "all" {
		var targets []struct {
			GOOS   string
			GOARCH string
		}
		for _, p := range supportedPlatforms {
			if p.GOOS == goos {
				targets = append(targets, p)
			}
		}
		if len(targets) > 0 {
			return targets
		}
	}

	// 如果只指定了 arch，使用该 arch 的所有操作系统
	if goos == "" || strings.ToLower(goos) == "all" {
		var targets []struct {
			GOOS   string
			GOARCH string
		}
		for _, p := range supportedPlatforms {
			if p.GOARCH == goarch {
				targets = append(targets, p)
			}
		}
		if len(targets) > 0 {
			return targets
		}
	}

	// 指定具体的 os 和 arch
	if goos != "" && goarch != "" {
		return []struct {
			GOOS   string
			GOARCH string
		}{{goos, goarch}}
	}

	// 默认返回当前环境
	return []struct {
		GOOS   string
		GOARCH string
	}{{runtime.GOOS, runtime.GOARCH}}
}

// buildForPlatform 为指定平台编译
func buildForPlatform(goos, goarch string) {
	pwd, _ := utils.ExecShell("pwd", "")
	pwd = strings.Replace(pwd, " ", "", -1)
	pwd = strings.Replace(pwd, "\n", "", -1)
	pwdArr := utils.Explode("/", pwd)
	if len(pwdArr) == 0 {
		return
	}
	projectName := pwdArr[len(pwdArr)-1]

	ext := ""
	if goos == "windows" {
		ext = ".exe"
	}

	outputName := projectName + ext

	// 只有在多平台或指定平台编译时才添加后缀

	// 指定了 os 或 arch 参数时，文件名带上平台后缀
	if buildOS != "" || buildArch != "" {
		outputName = fmt.Sprintf("%s-%s-%s%s", projectName, goos, goarch, ext)
	}

	buildDir := "build"
	if err := os.MkdirAll(buildDir, 0755); err != nil {
		fmt.Printf("Failed to create build directory: %v\n", err)
		return
	}

	outputPath := filepath.Join(buildDir, outputName)
	cmd := fmt.Sprintf("GOOS=%s GOARCH=%s go build -o %s", goos, goarch, outputPath)

	fmt.Printf("Building for %s/%s...\n", goos, goarch)
	utils.ExecShell(cmd, "")

	// 如果是 windows 且当前系统不是 windows，还需要编译无后缀版本
	if goos == "windows" && runtime.GOOS != "windows" {
		winextPath := filepath.Join(buildDir, projectName+".exe")
		if _, err := os.Stat(outputPath); err == nil {
			os.Rename(outputPath, winextPath)
		}
	}
}

// linux下编译打包（保留兼容）
func buildProjectWithLinux() {
	pwd, _ := utils.ExecShell("pwd", "")
	pwd = strings.Replace(pwd, " ", "", -1)
	pwd = strings.Replace(pwd, "\r\n", "", -1)
	pwdArr := utils.Explode("/", pwd)
	if len(pwdArr) == 0 {
		return
	}
	projectName := pwdArr[len(pwdArr)-1]
	utils.ExecShell(fmt.Sprintf("go build -o build/%s", projectName), "")
}

// windows下编译打包（保留兼容）
func buildProjectWithWindows() {
	pwd, _ := utils.ExecShell("cd", "")
	pwd = strings.Replace(pwd, " ", "", -1)
	pwd = strings.Replace(pwd, "\r\n", "", -1)
	pwdArr := utils.Explode("\\", pwd)
	if len(pwdArr) == 0 {
		return
	}
	projectName := pwdArr[len(pwdArr)-1]

	utils.ExecShell(fmt.Sprintf("go build -o build/%s.exe", projectName), "")
}
