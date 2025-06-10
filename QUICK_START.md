# 快速开始指南

## 🎯 目标
在您的开发电脑上构建程序，然后发布到目标电脑运行。

## 📋 准备工作

### 开发电脑 (当前macOS)
- [x] Go 1.19+ 已安装
- [x] 项目代码已就绪
- [ ] 需要构建发布版本

### 目标电脑 (macOS 10.13.6)
- [ ] 需要接收发布包
- [ ] 需要配置和运行

## 🔨 构建步骤

### 1. 在开发电脑上构建
```bash
cd /Users/norman/GolandProjects/macos-clodop-schoolpal
chmod +x build.sh
./build.sh
```

构建完成后会生成：
- `build/` 目录 - 包含所有发布文件
- `hprt-printer-setup-v1.0.tar.gz` - 发布包

### 2. 传输到目标电脑
将 `hprt-printer-setup-v1.0.tar.gz` 复制到目标电脑

### 3. 在目标电脑上部署
```bash
# 解压发布包
tar -xzf hprt-printer-setup-v1.0.tar.gz
cd hprt-printer-setup-v1.0

# 放入驱动文件
# 将 hprt-pos-printer-driver-v1.2.16.pkg 复制到当前目录

# 修改配置
# 编辑 config.yaml 文件，设置VPN名称和Windows IP

# 运行程序
双击 printer-setup 文件
```

## ⚙️ 配置说明

在目标电脑上修改 `config.yaml`：

```yaml
vpn:
  name: "实际的VPN连接名称"  # 查看系统偏好设置->网络

network:
  local_port: "8443"
  remote_host: "192.168.1.200"  # Windows电脑的实际IP
  remote_port: "8443"

printer:
  model: "HPRT_TP80B"
  driver_file: "hprt-pos-printer-driver-v1.2.16.pkg"
```

## 🚀 运行效果

程序启动后会显示图形界面，自动执行9个配置步骤：

1. ✅ 环境检查
2. ✅ 驱动验证  
3. ⏳ 安装驱动 (需要用户授权)
4. ✅ 检测打印机
5. ✅ 安装socat
6. ✅ 配置CUPS
7. ✅ 连接VPN
8. ✅ 端口转发
9. ✅ 测试连接

## 🆘 常见问题

### 权限问题
- 程序会自动请求管理员权限
- 如果权限被拒绝，请手动运行 `sudo ./printer-setup`

### VPN连接失败
- 检查VPN名称是否正确
- 确保VPN配置文件存在于系统中

### 网络连接问题
- 确保Windows电脑IP地址正确
- 检查防火墙设置

## 📞 技术支持

如果遇到问题：
1. 查看程序界面中的详细日志
2. 检查 `README.md` 中的详细说明
3. 验证配置文件设置是否正确 