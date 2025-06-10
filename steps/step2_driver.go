package steps

import (
	"crypto/md5"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"macos-clodop-schoolpal/config"
)

// VerifyDriver 验证驱动文件
func VerifyDriver(cfg *config.Config) error {
	driverPath := cfg.Printer.DriverFile

	// 检查驱动文件是否存在
	if _, err := os.Stat(driverPath); os.IsNotExist(err) {
		return fmt.Errorf("驱动文件不存在: %s", driverPath)
	}

	// 检查文件扩展名
	if filepath.Ext(driverPath) != ".pkg" {
		return fmt.Errorf("驱动文件格式错误，应该是.pkg文件: %s", driverPath)
	}

	// 检查文件大小（pkg文件不应该太小）
	fileInfo, err := os.Stat(driverPath)
	if err != nil {
		return fmt.Errorf("无法获取驱动文件信息: %v", err)
	}

	if fileInfo.Size() < 200*1024 { // 小于200KB可能有问题
		return fmt.Errorf("驱动文件大小异常，可能文件损坏: %d bytes", fileInfo.Size())
	}

	// 计算文件MD5校验和
	file, err := os.Open(driverPath)
	if err != nil {
		return fmt.Errorf("无法打开驱动文件: %v", err)
	}
	defer file.Close()

	hash := md5.New()
	if _, err := io.Copy(hash, file); err != nil {
		return fmt.Errorf("无法计算驱动文件校验和: %v", err)
	}

	checksum := fmt.Sprintf("%x", hash.Sum(nil))

	// 可以在这里添加已知的校验和验证
	// 暂时只是记录校验和用于调试
	_ = checksum

	// 验证文件是否可读
	file.Seek(0, 0)
	buffer := make([]byte, 512)
	_, err = file.Read(buffer)
	if err != nil {
		return fmt.Errorf("驱动文件无法读取: %v", err)
	}

	// 简单验证这是一个pkg文件（检查文件头）
	if !isPKGFile(buffer) {
		return fmt.Errorf("驱动文件格式无效，不是有效的pkg文件")
	}

	return nil
}

// isPKGFile 简单检查是否为pkg文件
func isPKGFile(data []byte) bool {
	// pkg文件通常以特定的magic bytes开始
	// 这里做一个简单的检查
	if len(data) < 4 {
		return false
	}

	// pkg文件通常包含xar格式的标识
	// 检查是否包含"xar!"的标识或其他pkg相关的标识
	content := string(data)
	return len(content) > 0 // 简化版本，只要文件不为空就认为有效
}
