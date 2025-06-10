package steps

import (
	"fmt"
	"net"
	"net/http"
	"os/exec"
	"strings"
	"time"

	"macos-clodop-schoolpal/config"
)

// ConfigureCUPS 配置CUPS打印服务
func ConfigureCUPS(cfg *config.Config) error {
	fmt.Println("🖨️ ========== CUPS打印服务配置 ==========")

	// 检查是否已经配置过
	if isCUPSConfigured() {
		fmt.Println("✅ CUPS已经配置完成")

		// 显示详细状态信息
		err := showCUPSStatus()
		if err != nil {
			fmt.Printf("⚠️ 获取CUPS状态时出错: %v\n", err)
		}

		return nil
	}

	// 如果CUPS未运行，先启动它
	if !isCUPSRunning() {
		fmt.Println("🔄 CUPS服务未运行，正在启动...")
		err := startCUPS()
		if err != nil {
			return fmt.Errorf("启动CUPS服务失败: %v", err)
		}
		fmt.Println("✅ CUPS服务启动成功")
	}

	fmt.Println("🔧 配置CUPS共享设置...")

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

	fmt.Println("✅ CUPS配置完成")

	// 等待服务重启
	fmt.Println("⏳ 等待CUPS服务重启...")
	time.Sleep(3 * time.Second)

	// 显示详细状态信息
	err = showCUPSStatus()
	if err != nil {
		fmt.Printf("⚠️ 获取CUPS状态时出错: %v\n", err)
	}

	return nil
}

// showCUPSStatus 显示CUPS详细状态信息
func showCUPSStatus() error {
	fmt.Println("\n📊 ========== CUPS状态信息 ==========")

	// 1. 获取本机IP地址
	localIP, err := getLocalIP()
	if err != nil {
		fmt.Printf("⚠️ 无法获取本机IP: %v\n", err)
		localIP = "localhost"
	} else {
		fmt.Printf("🌐 本机IP地址: %s\n", localIP)
	}

	// 2. CUPS管理界面
	cupsAdminURL := fmt.Sprintf("http://%s:631", localIP)
	fmt.Printf("🖥️ CUPS管理界面: %s\n", cupsAdminURL)

	// 3. 测试CUPS管理界面是否可访问
	fmt.Printf("🔍 测试CUPS管理界面访问性...")
	if testCUPSAccess(localIP) {
		fmt.Println(" ✅ 可访问")
	} else {
		fmt.Println(" ❌ 无法访问")
	}

	// 4. 获取已安装的打印机
	printers, err := getInstalledPrinters()
	if err != nil {
		fmt.Printf("⚠️ 获取打印机列表失败: %v\n", err)
	} else if len(printers) > 0 {
		fmt.Println("🖨️ 已安装的打印机:")
		for _, printer := range printers {
			printerURL := fmt.Sprintf("ipp://%s:631/printers/%s", localIP, printer)
			fmt.Printf("   • %s\n", printer)
			fmt.Printf("     📡 共享地址: %s\n", printerURL)
			fmt.Printf("     🪟 Windows添加: http://%s:631/printers/%s\n", localIP, printer)
		}
	} else {
		fmt.Println("ℹ️ 暂无已安装的打印机")
	}

	// 5. 提供操作提示
	fmt.Println("\n💡 ========== 使用提示 ==========")
	fmt.Printf("1. 在浏览器中打开: %s\n", cupsAdminURL)
	fmt.Println("2. 在CUPS管理界面中添加和管理打印机")
	fmt.Println("3. Windows电脑添加网络打印机时使用上述共享地址")
	fmt.Println("4. 确保防火墙允许631端口访问")

	return nil
}

// getLocalIP 获取本机局域网IP地址
func getLocalIP() (string, error) {
	interfaces, err := net.Interfaces()
	if err != nil {
		return "", err
	}

	for _, iface := range interfaces {
		// 跳过回环接口
		if iface.Flags&net.FlagLoopback != 0 {
			continue
		}
		// 跳过未启用的接口
		if iface.Flags&net.FlagUp == 0 {
			continue
		}

		addrs, err := iface.Addrs()
		if err != nil {
			continue
		}

		for _, addr := range addrs {
			var ip net.IP
			switch v := addr.(type) {
			case *net.IPNet:
				ip = v.IP
			case *net.IPAddr:
				ip = v.IP
			}

			// 只返回IPv4地址，且不是回环地址
			if ip == nil || ip.IsLoopback() {
				continue
			}
			if ip.To4() == nil {
				continue // 跳过IPv6
			}

			// 优先返回局域网地址
			if ip.IsPrivate() {
				return ip.String(), nil
			}
		}
	}

	return "", fmt.Errorf("未找到有效的局域网IP地址")
}

// testCUPSAccess 测试CUPS服务是否可访问
func testCUPSAccess(ip string) bool {
	client := &http.Client{
		Timeout: 5 * time.Second,
	}

	url := fmt.Sprintf("http://%s:631", ip)
	resp, err := client.Get(url)
	if err != nil {
		return false
	}
	defer resp.Body.Close()

	return resp.StatusCode == 200
}

// getInstalledPrinters 获取已安装的打印机列表
func getInstalledPrinters() ([]string, error) {
	cmd := exec.Command("lpstat", "-p")
	output, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	var printers []string
	lines := strings.Split(string(output), "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "printer ") {
			// 格式: printer PrinterName is idle. ...
			parts := strings.Fields(line)
			if len(parts) >= 2 {
				printers = append(printers, parts[1])
			}
		}
	}

	return printers, nil
}

// OpenCUPSAdmin 打开CUPS管理界面
func OpenCUPSAdmin() error {
	localIP, err := getLocalIP()
	if err != nil {
		localIP = "localhost"
	}

	url := fmt.Sprintf("http://%s:631", localIP)
	fmt.Printf("🌐 正在打开CUPS管理界面: %s\n", url)

	cmd := exec.Command("open", url)
	return cmd.Run()
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
