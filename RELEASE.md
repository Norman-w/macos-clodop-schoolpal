# 发布说明

## 自动构建和发布流程

本项目使用 GitHub Actions 进行自动构建，确保为 Intel 和 ARM (Apple Silicon) 架构提供真正的原生版本。

### 🚀 如何发布新版本

#### 方法一：创建版本标签（推荐）

1. **本地创建标签并推送**
   ```bash
   # 创建新的版本标签
   git tag v1.0.0
   
   # 推送标签到 GitHub
   git push origin v1.0.0
   ```

2. **GitHub Actions 自动执行**
   - 自动在两个架构上构建：Intel (x86_64) 和 ARM (arm64)
   - 自动获取和打包 socat 组件
   - 自动创建 GitHub Release
   - 自动上传构建产物

#### 方法二：手动触发构建

1. 访问 GitHub 仓库的 Actions 页面
2. 选择 "构建和发布 MacOS校宝打印组件" 工作流
3. 点击 "Run workflow" 按钮
4. 选择分支并手动触发

### 📦 构建产物

每次构建会产生以下文件：

- `MacOS校宝打印组件-Intel-YYYYMMDD-HHMM.zip` - Intel 版本
- `MacOS校宝打印组件-ARM-YYYYMMDD-HHMM.zip` - ARM 版本

每个压缩包包含：
- **MacOS校宝打印组件** - 主应用程序
- **config.yaml** - 配置文件
- **hprt-pos-printer-driver-v1.2.16.pkg** - HPRT 驱动程序
- **socat** - 网络代理工具（预装，无需 brew）
- **README.txt** - 详细说明文档

### 🔧 构建环境

- **Intel 版本**: 在 `macos-13` (Ventura) 上构建
- **ARM 版本**: 在 `macos-14` (Sonoma) 上构建  
- **Go 版本**: 1.21
- **包含组件**: socat 通过 Homebrew 自动安装和打包

### 🌟 优势

1. **真正的跨平台支持** - 每个架构都在对应的原生环境中构建
2. **socat 预装** - 解决 GFW 和网络环境问题，用户无需安装 brew
3. **自动化发布** - 无需手动操作，推送标签即可发布
4. **版本管理** - 自动创建 GitHub Release 和版本说明
5. **构建一致性** - 统一的构建环境确保产物质量

### 📋 版本命名规范

- 使用语义化版本号：`v主版本.次版本.修订版本`
- 例如：`v1.0.0`, `v1.1.0`, `v1.0.1`

### 🐛 故障排除

如果构建失败：

1. 检查 GitHub Actions 日志
2. 确保所有依赖文件存在（config.yaml, 驱动程序等）
3. 检查 Go 代码是否可以正常编译
4. 验证版本标签格式是否正确

### 🔍 本地构建

对于本地开发和测试，可以使用 `build.sh` 脚本：

```bash
./build.sh
```

注意：本地构建由于 Fyne GUI 框架的 CGO 依赖限制，无法进行真正的交叉编译。生产环境请使用 GitHub Actions。 