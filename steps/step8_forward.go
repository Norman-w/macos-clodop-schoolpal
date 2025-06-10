package steps

import (
	"context"
	"fmt"
	"os/exec"
	"time"

	"macos-clodop-schoolpal/config"
)

// StartPortForward 启动端口转发服务
func StartPortForward(cfg *config.Config) error {
	localPort := cfg.Network.LocalPort
	remoteHost := cfg.Network.RemoteHost
	remotePort := cfg.Network.RemotePort

	// 检查socat是否可用
	cmd := exec.Command("which", "socat")
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("socat未安装或不可用")
	}

	// 检查端口是否已经被占用
	if isPortInUse(localPort) {
		// 如果端口被占用，尝试停止现有的端口转发
		stopExistingPortForward(localPort)
	}

	// 启动端口转发
	forwardCmd := fmt.Sprintf("TCP-LISTEN:%s,fork", localPort)
	targetCmd := fmt.Sprintf("TCP:%s:%s", remoteHost, remotePort)

	ctx := context.Background()
	cmd = exec.CommandContext(ctx, "socat", forwardCmd, targetCmd)

	// 在后台启动端口转发
	err := cmd.Start()
	if err != nil {
		return fmt.Errorf("启动端口转发失败: %v", err)
	}

	// 等待一段时间确保端口转发启动成功
	time.Sleep(2 * time.Second)

	// 验证端口转发是否正常工作
	if !isPortInUse(localPort) {
		return fmt.Errorf("端口转发启动后端口仍不可用")
	}

	return nil
}

// isPortInUse 检查端口是否被占用
func isPortInUse(port string) bool {
	cmd := exec.Command("lsof", "-i", ":"+port)
	err := cmd.Run()
	return err == nil
}

// stopExistingPortForward 停止现有的端口转发
func stopExistingPortForward(port string) {
	// 查找占用端口的进程
	cmd := exec.Command("lsof", "-t", "-i", ":"+port)
	output, err := cmd.Output()
	if err != nil {
		return
	}

	// 终止进程
	if len(output) > 0 {
		pid := string(output)
		exec.Command("kill", "-9", pid).Run()
	}
}
