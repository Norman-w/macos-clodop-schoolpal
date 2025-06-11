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
		fmt.Printf("✅ 检测到同目录socat: %s\n", bundledPath)
		return nil
	} else if bundledErr == nil {
		fmt.Printf("⚠️ 找到同目录socat文件但不可执行: %s\n", bundledPath)
		fmt.Println("   正在检查权限...")
		// 尝试给socat添加执行权限
		exec.Command("chmod", "+x", bundledPath).Run()
		if isSocatExecutable(bundledPath) {
			fmt.Println("✅ 修复权限成功，同目录socat现在可用")
			return nil
		}
		fmt.Println("❌ 无法修复同目录socat的执行权限")
	} else {
		fmt.Printf("ℹ️ 同目录未找到socat文件 (路径: %s)\n", bundledPath)
	}

	// 检查系统中的socat是否已经安装
	if isSocatInstalled() {
		systemPath, _ := getSystemSocatPath()
		fmt.Printf("✅ 系统中已安装socat: %s\n", systemPath)
		return nil
	}

	fmt.Println("⚠️ 既没有同目录socat，也没有系统安装的socat")

	// 检查Homebrew是否安装
	if !isHomebrewInstalled() {
		fmt.Println("❌ 未安装Homebrew，无法自动安装socat")
		fmt.Println("💡 解决方案:")
		fmt.Println("   1. 访问 https://brew.sh 安装Homebrew")
		fmt.Println("   2. 或者将socat文件放在程序同目录下")
		fmt.Println("   3. 或者手动安装socat: brew install socat")
		return fmt.Errorf("需要先安装Homebrew或提供bundled socat")
	}

	// 使用Homebrew安装socat
	fmt.Println("📦 正在通过Homebrew安装socat...")
	cmd := exec.Command("brew", "install", "socat")
	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Printf("❌ Homebrew安装socat失败: %v\n", err)
		fmt.Printf("   输出: %s\n", string(output))
		fmt.Println("💡 请尝试手动安装: brew install socat")
		return fmt.Errorf("安装socat失败: %v", err)
	}

	// 验证安装是否成功
	if !isSocatInstalled() {
		fmt.Println("❌ socat安装后仍无法找到")
		fmt.Println("💡 请检查Homebrew配置或手动安装socat")
		return fmt.Errorf("socat安装后仍无法找到，请手动安装")
	}

	systemPath, _ := getSystemSocatPath()
	fmt.Printf("✅ socat安装完成: %s\n", systemPath)
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
			fmt.Printf("✅ 使用同目录socat: %s\n", bundledPath)
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
		fmt.Println("   1. 确保socat文件存在于程序同目录")
		fmt.Println("   2. 检查socat文件权限 (chmod +x socat)")
		fmt.Println("   3. 或安装系统版本: brew install socat")
		return "", fmt.Errorf("找不到可用的socat")
	}

	fmt.Printf("✅ 使用系统socat: %s\n", systemPath)
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
