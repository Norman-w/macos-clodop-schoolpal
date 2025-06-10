package steps

import (
	"fmt"
	"os/exec"
	"strings"
	"time"

	"macos-clodop-schoolpal/config"
)

// ConnectVPN 连接到指定VPN
func ConnectVPN(cfg *config.Config) error {
	vpnName := cfg.VPN.Name

	if vpnName == "" {
		return fmt.Errorf("配置文件中未指定VPN名称")
	}

	// 获取所有可用的VPN列表
	availableVPNs, err := getAvailableVPNs()
	if err != nil {
		return fmt.Errorf("无法获取VPN列表: %v", err)
	}

	if len(availableVPNs) == 0 {
		return fmt.Errorf("系统中没有配置任何VPN连接")
	}

	// 尝试找到匹配的VPN名称
	actualVPNName := findMatchingVPN(vpnName, availableVPNs)
	if actualVPNName == "" {
		fmt.Printf("❌ 找不到VPN '%s'\n", vpnName)
		fmt.Println("📋 系统中可用的VPN列表:")
		for i, vpn := range availableVPNs {
			fmt.Printf("  %d. %s\n", i+1, vpn)
		}
		return fmt.Errorf("VPN '%s' 不存在，请检查配置文件中的VPN名称", vpnName)
	}

	if actualVPNName != vpnName {
		fmt.Printf("💡 找到匹配VPN: '%s' -> '%s'\n", vpnName, actualVPNName)
	}

	// 检查VPN是否已连接
	if isVPNConnected(actualVPNName) {
		fmt.Printf("✅ VPN '%s' 已连接，跳过此步骤\n", actualVPNName)
		return nil
	}

	fmt.Printf("🔗 正在连接VPN '%s'...\n", actualVPNName)

	// 使用networksetup连接VPN（可以正确访问keychain）
	cmd := exec.Command("networksetup", "-connectpppoeservice", actualVPNName)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("无法连接VPN '%s': %v\n输出: %s", actualVPNName, err, string(output))
	}

	// 等待连接成功
	fmt.Print("⏳ 等待VPN连接")
	for i := 0; i < 30; i++ {
		time.Sleep(time.Second)
		fmt.Print(".")

		// 检查连接状态
		status := getVPNStatus(actualVPNName)
		if strings.Contains(status, "Connected") {
			fmt.Println()
			fmt.Printf("✅ VPN '%s' 连接成功\n", actualVPNName)
			return nil
		}

		// 检查是否有连接错误
		if strings.Contains(status, "Disconnected") && i > 5 {
			fmt.Println()
			return fmt.Errorf("VPN连接失败，请检查VPN配置和网络状况")
		}
	}
	fmt.Println()

	return fmt.Errorf("VPN连接超时，请检查VPN配置和网络状况")
}

// getAvailableVPNs 获取所有可用的VPN连接
func getAvailableVPNs() ([]string, error) {
	// 使用networksetup获取VPN列表
	cmd := exec.Command("networksetup", "-listallnetworkservices")
	output, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	lines := strings.Split(string(output), "\n")
	var vpns []string

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "An asterisk") {
			continue
		}

		// 跳过Wi-Fi等非VPN服务，只保留VPN相关的
		if !strings.Contains(strings.ToLower(line), "wi-fi") &&
			!strings.Contains(strings.ToLower(line), "ethernet") &&
			!strings.Contains(strings.ToLower(line), "ax88179a") &&
			!strings.Contains(strings.ToLower(line), "xreal") &&
			line != "" {
			vpns = append(vpns, line)
		}
	}

	return vpns, nil
}

// findMatchingVPN 查找匹配的VPN名称（支持模糊匹配）
func findMatchingVPN(targetName string, availableVPNs []string) string {
	// 1. 精确匹配
	for _, vpn := range availableVPNs {
		if vpn == targetName {
			return vpn
		}
	}

	// 2. 忽略大小写匹配
	targetLower := strings.ToLower(targetName)
	for _, vpn := range availableVPNs {
		if strings.ToLower(vpn) == targetLower {
			return vpn
		}
	}

	// 3. 移除空格后匹配
	targetNoSpaces := strings.ReplaceAll(targetName, " ", "")
	for _, vpn := range availableVPNs {
		vpnNoSpaces := strings.ReplaceAll(vpn, " ", "")
		if strings.ToLower(vpnNoSpaces) == strings.ToLower(targetNoSpaces) {
			return vpn
		}
	}

	// 4. 包含匹配
	for _, vpn := range availableVPNs {
		if strings.Contains(strings.ToLower(vpn), targetLower) {
			return vpn
		}
	}

	return ""
}

// isVPNConnected 检查VPN是否已连接
func isVPNConnected(vpnName string) bool {
	cmd := exec.Command("scutil", "--nc", "status", vpnName)
	output, err := cmd.Output()
	if err != nil {
		return false
	}

	return strings.Contains(string(output), "Connected")
}

// getVPNStatus 获取VPN详细状态
func getVPNStatus(vpnName string) string {
	cmd := exec.Command("scutil", "--nc", "status", vpnName)
	output, err := cmd.Output()
	if err != nil {
		return "Unknown"
	}
	return string(output)
}

// DisconnectVPN 断开VPN连接（新增功能）
func DisconnectVPN(vpnName string) error {
	cmd := exec.Command("networksetup", "-disconnectpppoeservice", vpnName)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("断开VPN失败: %v\n输出: %s", err, string(output))
	}
	return nil
}
