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
	// 首先检查是否有预装的socat（与应用程序同目录）
	if isBundledSocatAvailable() {
		fmt.Println("✅ 检测到预装的socat，跳过安装步骤")
		return nil
	}

	// 检查系统中的socat是否已经安装
	if isSocatInstalled() {
		fmt.Println("✅ 系统中已安装socat，跳过安装步骤")
		return nil
	}

	// 检查Homebrew是否安装
	if !isHomebrewInstalled() {
		return fmt.Errorf("需要先安装Homebrew。请访问 https://brew.sh 安装Homebrew后重新运行")
	}

	// 使用Homebrew安装socat
	fmt.Println("📦 正在通过Homebrew安装socat...")
	cmd := exec.Command("brew", "install", "socat")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("安装socat失败: %v\n输出: %s", err, string(output))
	}

	// 验证安装是否成功
	if !isSocatInstalled() {
		return fmt.Errorf("socat安装后仍无法找到，请手动安装")
	}

	fmt.Println("✅ socat安装完成")
	return nil
}

// GetSocatPath 获取socat的路径，优先返回预装版本
func GetSocatPath() (string, error) {
	// 优先使用预装的socat
	bundledPath, err := utils.GetResourcePath("socat")
	if err == nil {
		if isSocatExecutable(bundledPath) {
			return bundledPath, nil
		}
	}

	// 如果没有预装版本，使用系统安装的版本
	cmd := exec.Command("which", "socat")
	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("socat未安装或不可用")
	}

	systemPath := strings.TrimSpace(string(output))
	if isSocatExecutable(systemPath) {
		return systemPath, nil
	}

	return "", fmt.Errorf("找不到可用的socat")
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
