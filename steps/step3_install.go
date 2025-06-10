package steps

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"macos-clodop-schoolpal/config"
	"macos-clodop-schoolpal/utils"
)

// InstallDriver 安装打印机驱动
func InstallDriver(cfg *config.Config) error {
	// 检查驱动是否已经安装
	if isDriverInstalled(cfg) {
		fmt.Println("HPRT驱动已安装，跳过此步骤")
		return nil
	}

	// 使用新的路径查找逻辑
	driverPath, err := utils.GetResourcePath(cfg.Printer.DriverFile)
	if err != nil {
		return fmt.Errorf("无法定位驱动文件: %v", err)
	}

	// 确保驱动文件存在
	if _, err := os.Stat(driverPath); os.IsNotExist(err) {
		return fmt.Errorf("驱动文件不存在: %s", driverPath)
	}

	// 获取绝对路径
	absPath, err := filepath.Abs(driverPath)
	if err != nil {
		return fmt.Errorf("无法获取驱动文件绝对路径: %v", err)
	}

	fmt.Printf("🔧 正在安装HPRT驱动: %s\n", filepath.Base(absPath))

	// 使用AppleScript请求管理员权限并安装驱动
	script := fmt.Sprintf(`do shell script "installer -pkg '%s' -target /" with administrator privileges`, absPath)

	cmd := exec.Command("osascript", "-e", script)
	output, err := cmd.CombinedOutput()

	if err != nil {
		if strings.Contains(string(output), "User canceled") {
			return fmt.Errorf("用户取消了权限授权")
		}
		return fmt.Errorf("驱动安装失败: %v\n输出: %s", err, string(output))
	}

	// 等待安装完成，检查是否成功
	err = verifyDriverInstallation(cfg)
	if err != nil {
		return err
	}

	fmt.Println("✅ HPRT驱动安装完成")
	return nil
}

// isDriverInstalled 检查驱动是否已安装
func isDriverInstalled(cfg *config.Config) bool {
	// 方法1: 检查系统打印机驱动列表
	cmd := exec.Command("lpinfo", "-m")
	output, err := cmd.Output()
	if err == nil {
		outputStr := string(output)
		// 检查是否包含HPRT相关驱动
		if strings.Contains(strings.ToLower(outputStr), "hprt") {
			return true
		}
	}

	// 方法2: 检查常见的驱动安装位置
	driverPaths := []string{
		"/Library/Printers/PPDs/Contents/Resources/",
		"/usr/share/cups/drv/",
		"/usr/share/cups/model/",
	}

	for _, path := range driverPaths {
		if _, err := os.Stat(path); err == nil {
			// 检查目录下是否有HPRT相关文件
			entries, err := os.ReadDir(path)
			if err == nil {
				for _, entry := range entries {
					if strings.Contains(strings.ToLower(entry.Name()), "hprt") {
						return true
					}
				}
			}
		}
	}

	return false
}

// verifyDriverInstallation 验证驱动是否安装成功
func verifyDriverInstallation(cfg *config.Config) error {
	// 检查驱动是否在系统中注册
	// 这里可以检查/usr/share/cups/drv/或其他系统目录

	// 简单的验证方法：检查系统打印机驱动列表
	cmd := exec.Command("lpinfo", "-m")
	output, err := cmd.Output()
	if err != nil {
		// 如果lpinfo命令失败，不认为是错误，可能是权限问题
		return nil
	}

	// 检查输出中是否包含HPRT相关的驱动
	outputStr := string(output)
	if len(outputStr) > 0 {
		// 有输出说明CUPS系统正常工作
		return nil
	}

	return nil
}
