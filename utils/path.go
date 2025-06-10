package utils

import (
	"os"
	"path/filepath"
)

// GetExecutableDir 获取可执行文件所在的目录
func GetExecutableDir() (string, error) {
	executable, err := os.Executable()
	if err != nil {
		return "", err
	}

	// 解析符号链接
	executable, err = filepath.EvalSymlinks(executable)
	if err != nil {
		return "", err
	}

	return filepath.Dir(executable), nil
}

// GetResourcePath 获取资源文件的绝对路径
// 优先查找可执行文件目录，如果不存在则查找当前工作目录
func GetResourcePath(filename string) (string, error) {
	// 首先尝试可执行文件目录
	execDir, err := GetExecutableDir()
	if err == nil {
		execPath := filepath.Join(execDir, filename)
		if _, err := os.Stat(execPath); err == nil {
			return execPath, nil
		}
	}

	// 如果可执行文件目录没有，尝试当前工作目录
	wd, err := os.Getwd()
	if err != nil {
		return "", err
	}

	wdPath := filepath.Join(wd, filename)
	if _, err := os.Stat(wdPath); err == nil {
		return wdPath, nil
	}

	// 如果都没找到，返回可执行文件目录的路径（即使不存在）
	if execDir != "" {
		return filepath.Join(execDir, filename), nil
	}

	// 最后返回工作目录路径
	return wdPath, nil
}
