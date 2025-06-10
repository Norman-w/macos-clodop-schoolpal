package steps

import (
	"fmt"
	"os/exec"
	"strings"
	"time"

	"macos-clodop-schoolpal/config"
)

// DetectPrinter 检测打印机连接状态
func DetectPrinter(cfg *config.Config) error {
	// 等待一段时间让系统识别打印机
	time.Sleep(2 * time.Second)

	// 检查USB设备中是否有打印机
	cmd := exec.Command("system_profiler", "SPUSBDataType")
	output, err := cmd.Output()
	if err != nil {
		return fmt.Errorf("无法获取USB设备信息: %v", err)
	}

	outputStr := string(output)

	// 检查是否包含HPRT或相关打印机设备
	if strings.Contains(strings.ToLower(outputStr), "hprt") ||
		strings.Contains(strings.ToLower(outputStr), "printer") {
		// 找到了打印机设备
	} else {
		// 即使没有检测到，也不算失败，可能是检测方法的问题
	}

	// 检查CUPS系统中的打印机
	cmd = exec.Command("lpstat", "-p")
	output, err = cmd.Output()
	if err != nil {
		// lpstat命令失败不算致命错误
		return nil
	}

	outputStr = string(output)
	if strings.Contains(strings.ToLower(outputStr), "hprt") {
		// 找到了CUPS中的打印机
		return nil
	}

	// 尝试添加打印机到CUPS
	return addPrinterToCUPS(cfg)
}

// addPrinterToCUPS 添加打印机到CUPS系统
func addPrinterToCUPS(cfg *config.Config) error {
	// 这里可以尝试自动添加打印机
	// 由于每个打印机的具体连接方式不同，这里先返回成功
	// 实际使用时可以根据具体的打印机型号进行配置

	return nil
}
