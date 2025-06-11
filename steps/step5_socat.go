package steps

import (
	"fmt"
	"os/exec"
	"strings"

	"macos-clodop-schoolpal/config"
	"macos-clodop-schoolpal/utils"
)

// InstallSocat 安装socat网络工具
func InstallSocat(cfg *config.Config) error {
	fmt.Println("🔧 ========== Socat网络工具检查 ==========")

	// 首先检查是否有预装的socat（与应用程序同目录）
	bundledPath, bundledErr := utils.GetResourcePath("socat")
	if bundledErr == nil && isSocatExecutable(bundledPath) {
		fmt.Printf("✅ 检测到同目录静态socat: %s\n", bundledPath)
		fmt.Println("💡 使用内置静态编译版本，无需任何系统依赖！")
		return nil
	} else if bundledErr == nil {
		fmt.Printf("⚠️ 找到同目录socat文件但不可执行: %s\n", bundledPath)
		fmt.Println("   正在检查权限...")
		// 尝试给socat添加执行权限
		exec.Command("chmod", "+x", bundledPath).Run()
		if isSocatExecutable(bundledPath) {
			fmt.Println("✅ 修复权限成功，同目录静态socat现在可用")
			fmt.Println("💡 使用内置静态编译版本，无需任何系统依赖！")
			return nil
		}
		fmt.Println("❌ 无法修复同目录socat的执行权限")
	} else {
		fmt.Printf("ℹ️ 同目录未找到socat文件 (路径: %s)\n", bundledPath)
		fmt.Println("💡 推荐使用官方发布版本，内置静态编译的socat")
	}

	// 检查系统中的socat是否已经安装
	if isSocatInstalled() {
		systemPath, _ := getSystemSocatPath()
		fmt.Printf("⚠️ 发现系统socat: %s\n", systemPath)
		fmt.Println("   注意：系统版本可能有动态库依赖问题")
		fmt.Println("   建议使用官方发布版本的内置静态socat")
		return nil
	}

	fmt.Println("⚠️ 既没有同目录socat，也没有系统安装的socat")
	fmt.Println("")
	fmt.Println("🎯 推荐解决方案（按优先级排序）:")
	fmt.Println("   1. ⭐ 下载官方发布版本 - 内置静态编译socat，无依赖")
	fmt.Println("      GitHub Releases: https://github.com/Norman-w/macos-clodop-schoolpal/releases")
	fmt.Println("   2. 📁 手动放置socat - 将socat文件放在程序同目录")
	fmt.Println("   3. 🍺 使用Homebrew - brew install socat（可能有依赖问题）")
	fmt.Println("")

	// 检查Homebrew是否安装
	if !isHomebrewInstalled() {
		fmt.Println("❌ 未安装Homebrew，无法自动安装socat")
		fmt.Println("💡 强烈建议下载官方发布版本，避免复杂的安装过程")
		return fmt.Errorf("需要socat支持，请下载官方发布版本或手动安装")
	}

	fmt.Println("🤔 检测到Homebrew，是否尝试安装系统版socat？")
	fmt.Println("⚠️ 警告：Homebrew安装的socat可能在目标机器上有依赖问题")
	fmt.Println("📦 正在通过Homebrew安装socat（不推荐用于生产）...")

	cmd := exec.Command("brew", "install", "socat")
	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Printf("❌ Homebrew安装socat失败: %v\n", err)
		fmt.Printf("   输出: %s\n", string(output))
		fmt.Println("💡 建议下载官方发布版本，包含静态编译的socat")
		return fmt.Errorf("安装socat失败，建议使用官方发布版本")
	}

	// 验证安装是否成功
	if !isSocatInstalled() {
		fmt.Println("❌ socat安装后仍无法找到")
		fmt.Println("💡 建议下载官方发布版本，避免安装问题")
		return fmt.Errorf("socat安装失败，建议使用官方发布版本")
	}

	systemPath, _ := getSystemSocatPath()
	fmt.Printf("✅ socat安装完成: %s\n", systemPath)
	fmt.Println("⚠️ 注意：当前使用的是动态链接版本，在其他机器上可能有依赖问题")
	fmt.Println("💡 建议在生产环境使用官方发布版本的静态编译socat")
	return nil
}

// GetSocatPath 获取socat的路径，优先返回预装版本
func GetSocatPath() (string, error) {
	fmt.Println("🔍 查找socat路径...")

	// 优先使用预装的socat（同目录）
	bundledPath, err := utils.GetResourcePath("socat")
	if err == nil {
		fmt.Printf("   检查同目录socat: %s\n", bundledPath)
		if isSocatExecutable(bundledPath) {
			fmt.Printf("✅ 使用同目录静态socat: %s\n", bundledPath)
			fmt.Println("💡 静态编译版本，无外部依赖，推荐！")
			return bundledPath, nil
		} else {
			fmt.Printf("⚠️ 同目录socat不可执行: %s\n", bundledPath)
		}
	} else {
		fmt.Printf("   同目录未找到socat: %v\n", err)
	}

	// 如果没有预装版本，使用系统安装的版本
	systemPath, err := getSystemSocatPath()
	if err != nil {
		fmt.Println("❌ 也未找到系统安装的socat")
		fmt.Println("💡 故障排除:")
		fmt.Println("   1. ⭐ 推荐：下载官方发布版本（内置静态socat）")
		fmt.Println("   2. 确保socat文件存在于程序同目录")
		fmt.Println("   3. 检查socat文件权限 (chmod +x socat)")
		fmt.Println("   4. 或安装系统版本: brew install socat")
		return "", fmt.Errorf("找不到可用的socat")
	}

	fmt.Printf("⚠️ 使用系统socat: %s\n", systemPath)
	fmt.Println("   注意：系统版本可能有动态库依赖，在其他机器上可能无法运行")

	if isSocatExecutable(systemPath) {
		return systemPath, nil
	}

	fmt.Printf("❌ 系统socat不可执行: %s\n", systemPath)
	return "", fmt.Errorf("找不到可用的socat")
}

// getSystemSocatPath 获取系统安装的socat路径
func getSystemSocatPath() (string, error) {
	cmd := exec.Command("which", "socat")
	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("socat未安装或不在PATH中")
	}
	return strings.TrimSpace(string(output)), nil
}

// isBundledSocatAvailable 检查是否有预装的socat
func isBundledSocatAvailable() bool {
	bundledPath, err := utils.GetResourcePath("socat")
	if err != nil {
		return false
	}
	return isSocatExecutable(bundledPath)
}

// isSocatExecutable 检查指定路径的socat是否可执行
func isSocatExecutable(path string) bool {
	cmd := exec.Command(path, "-V")
	err := cmd.Run()
	return err == nil
}

// isSocatInstalled 检查socat是否已安装（系统版本）
func isSocatInstalled() bool {
	cmd := exec.Command("which", "socat")
	err := cmd.Run()
	return err == nil
}

// isHomebrewInstalled 检查Homebrew是否已安装
func isHomebrewInstalled() bool {
	cmd := exec.Command("which", "brew")
	err := cmd.Run()
	if err != nil {
		return false
	}

	// 进一步验证brew命令是否可用
	cmd = exec.Command("brew", "--version")
	output, err := cmd.Output()
	if err != nil {
		return false
	}

	return strings.Contains(string(output), "Homebrew")
}
