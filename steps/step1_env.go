package steps

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"

	"macos-clodop-schoolpal/config"
)

// CheckEnvironment 检查系统环境
func CheckEnvironment(cfg *config.Config) error {
	// 检查操作系统
	if runtime.GOOS != "darwin" {
		return fmt.Errorf("此程序仅支持macOS系统，当前系统: %s", runtime.GOOS)
	}

	// 检查macOS版本
	cmd := exec.Command("sw_vers", "-productVersion")
	output, err := cmd.Output()
	if err != nil {
		return fmt.Errorf("无法获取macOS版本: %v", err)
	}

	version := strings.TrimSpace(string(output))
	if !isValidMacOSVersion(version) {
		return fmt.Errorf("macOS版本过低，需要10.13.6或更高版本，当前版本: %s", version)
	}

	// 检查当前用户是否为管理员组成员
	cmd = exec.Command("id", "-Gn")
	output, err = cmd.Output()
	if err != nil {
		return fmt.Errorf("无法检查用户权限: %v", err)
	}

	groups := strings.Fields(string(output))
	isAdmin := false
	for _, group := range groups {
		if group == "admin" {
			isAdmin = true
			break
		}
	}

	if !isAdmin {
		return fmt.Errorf("当前用户不是管理员，无法执行系统配置")
	}

	// 检查网络连接
	cmd = exec.Command("ping", "-c", "1", "8.8.8.8")
	err = cmd.Run()
	if err != nil {
		return fmt.Errorf("网络连接检查失败，请确保网络正常")
	}

	// 检查工作目录权限
	wd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("无法获取当前工作目录: %v", err)
	}

	// 检查是否可以在当前目录创建文件
	testFile := wd + "/test_permission.tmp"
	file, err := os.Create(testFile)
	if err != nil {
		return fmt.Errorf("当前目录没有写入权限: %v", err)
	}
	file.Close()
	os.Remove(testFile)

	return nil
}

// isValidMacOSVersion 检查macOS版本是否满足要求
func isValidMacOSVersion(version string) bool {
	// 简单的版本检查，支持10.13.6及以上版本
	parts := strings.Split(version, ".")
	if len(parts) < 2 {
		return false
	}

	// 对于macOS 11及以上版本（Big Sur及以后）
	if parts[0] == "11" || parts[0] == "12" || parts[0] == "13" || parts[0] == "14" {
		return true
	}

	// 对于macOS 10.x版本
	if parts[0] == "10" && len(parts) >= 2 {
		major := parts[1]
		if major == "15" || major == "14" || major == "13" {
			return true
		}
	}

	return false
}
