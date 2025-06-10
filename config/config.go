package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

// Config 应用程序配置结构
type Config struct {
	VPN struct {
		Name string `yaml:"name"`
	} `yaml:"vpn"`

	Network struct {
		LocalPort  string `yaml:"local_port"`
		RemoteHost string `yaml:"remote_host"`
		RemotePort string `yaml:"remote_port"`
	} `yaml:"network"`

	Printer struct {
		Model      string `yaml:"model"`
		DriverFile string `yaml:"driver_file"`
	} `yaml:"printer"`
}

// LoadConfig 从YAML文件加载配置
func LoadConfig(filename string) (*Config, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("无法读取配置文件 %s: %v", filename, err)
	}

	var config Config
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		return nil, fmt.Errorf("无法解析配置文件: %v", err)
	}

	// 验证必要的配置项
	if config.VPN.Name == "" || config.VPN.Name == "请修改为你的VPN连接名称" {
		return nil, fmt.Errorf("请在config.yaml中设置正确的VPN名称")
	}

	if config.Network.RemoteHost == "" {
		return nil, fmt.Errorf("请在config.yaml中设置Windows电脑的IP地址")
	}

	return &config, nil
}

// Validate 验证配置是否完整
func (c *Config) Validate() error {
	if c.VPN.Name == "" {
		return fmt.Errorf("VPN名称不能为空")
	}

	if c.Network.LocalPort == "" {
		return fmt.Errorf("本地端口不能为空")
	}

	if c.Network.RemoteHost == "" {
		return fmt.Errorf("远程主机地址不能为空")
	}

	if c.Network.RemotePort == "" {
		return fmt.Errorf("远程端口不能为空")
	}

	if c.Printer.DriverFile == "" {
		return fmt.Errorf("打印机驱动文件名不能为空")
	}

	return nil
}
