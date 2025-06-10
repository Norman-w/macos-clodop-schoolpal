package steps

import (
	"fmt"
	"os/exec"
	"strings"

	"macos-clodop-schoolpal/config"
)

// ConfigureCUPS 配置CUPS打印服务
func ConfigureCUPS(cfg *config.Config) error {
	// 检查是否已经配置过
	if isCUPSConfigured() {
		fmt.Println("CUPS已经配置完成，跳过此步骤")
		return nil
	}

	// 如果CUPS未运行，先启动它
	if !isCUPSRunning() {
		fmt.Println("CUPS服务未运行，正在启动...")
		err := startCUPS()
		if err != nil {
			return fmt.Errorf("启动CUPS服务失败: %v", err)
		}
		fmt.Println("CUPS服务启动成功")
	}

	// 使用osascript执行需要管理员权限的命令
	script := `
	do shell script "cupsctl WebInterface=yes" with administrator privileges
	do shell script "cupsctl --remote-admin --remote-any --share-printers" with administrator privileges
	do shell script "launchctl stop org.cups.cupsd; launchctl start org.cups.cupsd" with administrator privileges
	`

	cmd := exec.Command("osascript", "-e", script)
	output, err := cmd.CombinedOutput()
	if err != nil {
		if strings.Contains(string(output), "User canceled") {
			return fmt.Errorf("用户取消了权限授权")
		}
		return fmt.Errorf("配置CUPS失败: %v\n输出: %s", err, string(output))
	}

	fmt.Println("CUPS配置完成")
	return nil
}

// isCUPSRunning 检查CUPS服务是否运行
func isCUPSRunning() bool {
	cmd := exec.Command("launchctl", "list", "org.cups.cupsd")
	err := cmd.Run()
	return err == nil
}

// startCUPS 启动CUPS服务
func startCUPS() error {
	script := `do shell script "launchctl start org.cups.cupsd" with administrator privileges`
	cmd := exec.Command("osascript", "-e", script)
	output, err := cmd.CombinedOutput()
	if err != nil {
		if strings.Contains(string(output), "User canceled") {
			return fmt.Errorf("用户取消了权限授权")
		}
		return fmt.Errorf("启动失败: %v", err)
	}
	return nil
}

// isCUPSConfigured 检查CUPS是否已经配置
func isCUPSConfigured() bool {
	// 检查WebInterface是否启用
	cmd := exec.Command("cupsctl")
	output, err := cmd.Output()
	if err != nil {
		return false
	}

	configOutput := string(output)
	return strings.Contains(configOutput, "WebInterface=yes")
}

// restartCUPS 重启CUPS服务
func restartCUPS() error {
	// 停止CUPS服务
	cmd := exec.Command("sudo", "launchctl", "stop", "org.cups.cupsd")
	cmd.Run() // 忽略错误，可能服务已经停止

	// 启动CUPS服务
	cmd = exec.Command("sudo", "launchctl", "start", "org.cups.cupsd")
	return cmd.Run()
}
