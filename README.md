# HPRT打印机一键配置工具

## 项目背景
基于您的文章《在MacOS使用Clodop打印的实现方法》，为macOS 10.13.6系统开发的一键配置工具。

**目标**：开机启动程序后自动完成所有配置，确保HPRT打印机可以通过Clodop正常工作。

## 技术架构
- **语言**: Go
- **GUI**: Fyne (原生跨平台GUI)
- **配置**: YAML格式 (用户友好)
- **执行方式**: 线性流程，一键自动化

## 9步配置流程
1. **环境检查** - 验证macOS版本和系统权限
2. **驱动验证** - 检查驱动文件完整性
3. **安装驱动** - 安装HPRT打印机驱动(.pkg文件)
4. **检测打印机** - 确认打印机连接状态
5. **安装工具** - 安装socat网络转发工具
6. **配置CUPS** - 设置CUPS打印服务
7. **连接VPN** - 连接到指定VPN
8. **端口转发** - 启动8443端口转发服务
9. **测试连接** - 验证打印机连接

## 项目结构
```
macos-clodop-schoolpal/
├── main.go                    # 程序入口
├── go.mod                     # Go模块文件
├── config.yaml                # 配置文件
├── hprt-pos-printer-driver-v1.2.16.pkg  # 打印机驱动(需要放入)
├── config/
│   └── config.go             # 配置文件读取
├── gui/
│   ├── window.go             # 主窗口界面
│   └── progress.go           # 进度显示组件
├── steps/
│   ├── step1_env.go          # 第1步：环境检查
│   ├── step2_driver.go       # 第2步：驱动验证
│   ├── step3_install.go      # 第3步：安装驱动
│   ├── step4_detect.go       # 第4步：检测打印机
│   ├── step5_socat.go        # 第5步：安装socat
│   ├── step6_cups.go         # 第6步：配置CUPS
│   ├── step7_vpn.go          # 第7步：连接VPN
│   ├── step8_forward.go      # 第8步：端口转发
│   └── step9_test.go         # 第9步：测试连接
└── utils/
    ├── system.go             # 系统操作工具
    ├── printer.go            # 打印机操作工具
    └── network.go            # 网络操作工具
```

## 开发和发布

### 开发环境要求
- macOS (开发机)
- Go 1.19+ 
- Fyne GUI依赖

### 构建发布版本
```bash
cd /Users/norman/GolandProjects/macos-clodop-schoolpal
chmod +x build.sh
./build.sh
```

### 发布到目标电脑
1. 将生成的 `hprt-printer-setup-v1.0.tar.gz` 复制到目标电脑
2. 解压文件包
3. 将 `hprt-pos-printer-driver-v1.2.16.pkg` 放入解压目录
4. 修改 `config.yaml` 配置文件
5. 运行 `./run.sh` 或直接运行 `./printer-setup`

## 目标电脑使用

### 环境要求
- macOS 10.13.6 或更高版本
- 管理员权限
- 网络连接

### 使用步骤
解压后的文件夹中包含详细的使用说明文件

## 配置说明

### config.yaml 配置项：
```yaml
# VPN配置
vpn:
  name: "你的VPN连接名称"  # 在系统偏好设置->网络中查看

# 网络配置  
network:
  local_port: "8443"
  remote_host: "192.168.1.200"  # Windows电脑IP地址
  remote_port: "8443"

# 打印机配置
printer:
  model: "HPRT_TP80B"
  driver_file: "hprt-pos-printer-driver-v1.2.16.pkg"
```

## 错误排查

### 常见问题：
1. **权限不足** - 程序会自动请求管理员权限
2. **VPN连接失败** - 检查VPN名称是否正确
3. **打印机未识别** - 确认USB连接和驱动安装
4. **端口转发失败** - 检查网络配置和防火墙设置

### 日志查看：
程序运行时会在界面显示详细的执行日志，包括每一步的成功/失败状态。

## 开发说明

### 项目背景来源
基于文章《在MacOS使用Clodop打印的实现方法》中的手动配置步骤，自动化实现：
- MacOS发起Clodop打印请求
- 通过socat转发到Windows打印机
- 再转发到连接在MacOS上的热敏打印机

### 技术决策
- **选择Go**: 系统集成能力强，编译成单一可执行文件
- **选择Fyne**: 原生GUI，无需web协议复杂性
- **选择YAML**: 比JSON更适合非技术用户配置
- **线性流程**: 简单直接，符合"一键配置"的需求

## 许可证
MIT License 