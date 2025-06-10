package steps

import (
	"fmt"
	"os/exec"
	"strings"

	"macos-clodop-schoolpal/config"
)

// InstallSocat 安装socat网络工具
func InstallSocat(cfg *config.Config) error {
	// 首先检查socat是否已经安装
	if isSocatInstalled() {
		return nil
	}

	// 检查Homebrew是否安装
	if !isHomebrewInstalled() {
		return fmt.Errorf("需要先安装Homebrew。请访问 https://brew.sh 安装Homebrew后重新运行")
	}

	// 使用Homebrew安装socat
	cmd := exec.Command("brew", "install", "socat")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("安装socat失败: %v\n输出: %s", err, string(output))
	}

	// 验证安装是否成功
	if !isSocatInstalled() {
		return fmt.Errorf("socat安装后仍无法找到，请手动安装")
	}

	return nil
}

// isSocatInstalled 检查socat是否已安装
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
