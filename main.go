package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"

	"macos-clodop-schoolpal/config"
	"macos-clodop-schoolpal/steps"
	"macos-clodop-schoolpal/utils"
)

// initChineseFont 初始化中文字体支持
func initChineseFont() {
	// 设置UTF-8环境
	os.Setenv("LC_ALL", "zh_CN.UTF-8")
	os.Setenv("LANG", "zh_CN.UTF-8")
	os.Setenv("LC_CTYPE", "zh_CN.UTF-8")

	// 完全不设置FYNE_FONT，让Fyne使用内置字体
	// 这样可以避免加载系统字体文件时的错误
}

// createChineseTheme 创建支持中文的主题
func createChineseTheme() fyne.Theme {
	return theme.DefaultTheme()
}

// preFlightCheck 程序启动前的检查
func preFlightCheck() error {
	// 检查系统版本
	cmd := exec.Command("sw_vers", "-productVersion")
	output, err := cmd.Output()
	if err == nil {
		version := strings.TrimSpace(string(output))
		if !isValidVersion(version) {
			return fmt.Errorf("系统版本过低，需要macOS 10.13.6或更高版本，当前版本: %s", version)
		}
	}

	// 检查必要文件 - 使用新的路径查找逻辑
	configPath, err := utils.GetResourcePath("config.yaml")
	if err != nil {
		return fmt.Errorf("无法定位配置文件: %v", err)
	}

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		return fmt.Errorf("配置文件不存在: %s", configPath)
	}

	// 给当前程序设置执行权限（如果需要的话）
	if err := os.Chmod(os.Args[0], 0755); err != nil {
		// 权限设置失败不算致命错误
	}

	return nil
}

// isValidVersion 检查macOS版本是否满足要求
func isValidVersion(version string) bool {
	parts := strings.Split(version, ".")
	if len(parts) < 2 {
		return false
	}

	// 对于macOS 11及以上版本
	if parts[0] == "11" || parts[0] == "12" || parts[0] == "13" || parts[0] == "14" || parts[0] == "15" {
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

func main() {
	// 初始化中文字体支持
	initChineseFont()

	// 前置检查
	if err := preFlightCheck(); err != nil {
		log.Fatalf("前置检查失败: %v", err)
	}

	// 确定配置文件路径 - 使用新的路径查找逻辑
	configPath, err := utils.GetResourcePath("config.yaml")
	if err != nil {
		log.Fatalf("无法定位配置文件: %v", err)
	}

	// 读取配置文件
	cfg, err := config.LoadConfig(configPath)
	if err != nil {
		log.Printf("配置文件错误: %v", err)
		// 如果配置文件有问题，仍然启动程序，但会在界面中显示错误
	}

	// 创建应用程序
	myApp := app.New()
	myApp.SetIcon(nil) // 可以后续添加图标

	// 设置中文主题
	myApp.Settings().SetTheme(createChineseTheme())

	// 创建主窗口
	window := myApp.NewWindow("HPRT打印机一键配置工具")
	window.Resize(fyne.NewSize(600, 500))
	window.CenterOnScreen()

	// 创建UI组件
	titleLabel := widget.NewLabel("HPRT打印机一键配置工具")
	titleLabel.Alignment = fyne.TextAlignCenter
	titleLabel.TextStyle = fyne.TextStyle{Bold: true}

	statusLabel := widget.NewLabel("准备开始配置...")
	progressBar := widget.NewProgressBar()
	progressBar.SetValue(0)

	// 创建日志文本区域
	logText := widget.NewEntry()
	logText.MultiLine = true
	logText.Wrapping = fyne.TextWrapWord
	logText.Disable() // 只读
	logContainer := container.NewScroll(logText)
	logContainer.SetMinSize(fyne.NewSize(580, 200))

	// 操作按钮
	cupsButton := widget.NewButton("打开CUPS管理", func() {
		go func() {
			err := steps.OpenCUPSAdmin()
			if err != nil {
				addLog(logText, fmt.Sprintf("❌ 打开CUPS管理界面失败: %v", err))
			} else {
				addLog(logText, "🌐 已打开CUPS管理界面")
			}
		}()
	})

	exitButton := widget.NewButton("退出程序", func() {
		myApp.Quit()
	})

	// 按钮容器
	buttonContainer := container.NewHBox(cupsButton, exitButton)

	// 布局
	content := container.NewVBox(
		titleLabel,
		widget.NewSeparator(),
		statusLabel,
		progressBar,
		widget.NewLabel("详细日志:"),
		logContainer,
		buttonContainer,
	)

	window.SetContent(content)

	// 如果配置文件正常，自动开始执行
	if cfg != nil {
		go runAllSteps(cfg, progressBar, statusLabel, logText, window)
	} else {
		statusLabel.SetText("❌ 配置文件错误，请检查config.yaml")
		addLog(logText, "❌ 配置文件加载失败: "+err.Error())
		addLog(logText, "💡 请检查并修改config.yaml文件后重新启动程序")
	}

	// 显示窗口并运行
	window.ShowAndRun()
}

// runAllSteps 执行所有配置步骤
func runAllSteps(cfg *config.Config, progressBar *widget.ProgressBar, statusLabel *widget.Label, logText *widget.Entry, window fyne.Window) {
	addLog := func(msg string) {
		addLog(logText, msg)
	}

	addLog("🚀 开始HPRT打印机自动配置")
	addLog(fmt.Sprintf("📋 配置信息: VPN=%s, 远程主机=%s:%s",
		cfg.VPN.Name, cfg.Network.RemoteHost, cfg.Network.RemotePort))

	// 定义所有步骤
	allSteps := []struct {
		Name        string
		Description string
		Execute     func(*config.Config) error
	}{
		{"环境检查", "检查系统版本和权限", steps.CheckEnvironment},
		{"验证驱动", "确认驱动文件完整性", steps.VerifyDriver},
		{"安装驱动", "安装HPRT打印机驱动", steps.InstallDriver},
		{"检测打印机", "检测打印机连接状态", steps.DetectPrinter},
		{"安装工具", "安装socat网络工具", steps.InstallSocat},
		{"配置CUPS", "配置CUPS打印服务", steps.ConfigureCUPS},
		{"连接VPN", "连接到指定VPN", steps.ConnectVPN},
		{"端口转发", "启动端口转发服务", steps.StartPortForward},
		{"测试连接", "测试打印机连接", steps.TestConnection},
	}

	totalSteps := len(allSteps)
	allSuccess := true

	for i, step := range allSteps {
		statusLabel.SetText(fmt.Sprintf("第%d步: %s", i+1, step.Description))
		addLog(fmt.Sprintf("🔄 第%d/%d步: %s", i+1, totalSteps, step.Name))

		err := step.Execute(cfg)
		if err != nil {
			// 在GUI上显示错误信息
			addLog(fmt.Sprintf("❌ %s 失败: %s", step.Name, err.Error()))
			statusLabel.SetText(fmt.Sprintf("❌ 配置失败: %s", step.Name))

			// 特别处理测试连接步骤的失败
			if step.Name == "测试连接" {
				addLog("⚠️ 打印测试失败！这可能导致打印功能无法正常工作")
				statusLabel.SetText("⚠️ 打印测试失败 - 请检查错误信息")
			}

			addLog("💡 配置失败，请查看错误信息后重新运行程序")
			addLog("🔍 请检查以下可能的问题:")

			// 根据失败的步骤提供具体建议
			switch step.Name {
			case "连接VPN":
				addLog("   - VPN配置是否正确（服务器地址、用户名、密码、共享密钥）")
				addLog("   - 网络连接是否正常")
				addLog("   - VPN服务器是否可访问")
			case "测试连接":
				addLog("   ⚠️ 以下问题可能导致打印测试失败:")
				addLog("   - 远程Windows电脑上Clodop服务未运行")
				addLog("   - 打印机未连接或未开机")
				addLog("   - VPN连接不稳定或已断开")
				addLog("   - 端口转发设置有问题")
				addLog("   - 防火墙阻止了HTTPS连接（端口8443）")
				addLog("   - SSL证书验证问题")
				addLog("💡 建议操作:")
				addLog("   1. 确认远程Windows电脑已安装并启动Clodop服务")
				addLog("   2. 检查打印机电源和USB连接")
				addLog("   3. 验证VPN连接状态")
				addLog("   4. 重新启动配置程序重试")
			default:
				addLog("   - 检查网络连接")
				addLog("   - 确认所需权限")
			}

			allSuccess = false
			break // 停止执行后续步骤
		}

		addLog(fmt.Sprintf("✅ %s 完成", step.Name))
		progressBar.SetValue(float64(i+1) / float64(totalSteps))

		// 添加短暂延迟，让用户看到进度
		time.Sleep(500 * time.Millisecond)
	}

	// 只有在所有步骤都成功时才隐藏窗口
	if allSuccess {
		statusLabel.SetText("🎉 配置完成！打印机已就绪")
		addLog("🎉 所有配置步骤完成！")
		addLog("✨ HPRT打印机现在可以通过Clodop正常使用了")
		addLog("📝 如果打印机已出纸，说明配置完全正常")
		addLog("🕒 请等待10秒确认打印结果...")

		// 延长等待时间，确保打印任务完成
		go func() {
			// 等待10秒，让用户确认打印结果
			for i := 10; i > 0; i-- {
				time.Sleep(1 * time.Second)
				if i <= 5 {
					addLog(fmt.Sprintf("💡 程序将在 %d 秒后隐藏窗口", i))
				}
			}
			addLog("🫥 程序已转入后台运行，可以关闭此窗口")
			window.Hide()
		}()
	} else {
		// 配置失败时，窗口保持显示，让用户查看错误信息
		addLog("🚫 配置未完成，窗口将保持显示以便查看错误信息")
		addLog("🔧 请根据上述建议修复问题后重新启动程序")
		addLog("📞 如需技术支持，请保存此日志信息")

		// 确保状态显示失败信息
		if !strings.Contains(statusLabel.Text, "❌") && !strings.Contains(statusLabel.Text, "⚠️") {
			statusLabel.SetText("❌ 配置过程中出现错误")
		}
	}
}

// addLog 添加日志信息
func addLog(logText *widget.Entry, msg string) {
	timestamp := time.Now().Format("15:04:05")
	newText := logText.Text + fmt.Sprintf("[%s] %s\n", timestamp, msg)
	logText.SetText(newText)

	// 自动滚动到底部
	logText.CursorRow = len(logText.Text)
}
