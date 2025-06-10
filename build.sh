#!/bin/bash

# HPRT打印机配置工具 - 构建脚本
# 用于在开发环境编译和打包发布版本

set -e

echo "🔨 开始构建HPRT打印机配置工具..."

# 检查Go环境
if ! command -v go &> /dev/null; then
    echo "❌ 错误：未找到Go环境"
    exit 1
fi

echo "✅ Go环境检查通过"

# 清理之前的构建
echo "🧹 清理构建目录..."
rm -rf build/
mkdir -p build

# 下载依赖
echo "📦 下载项目依赖..."
go mod tidy

# 编译程序 (针对macOS，支持中文)
echo "🔨 编译程序..."
export LC_ALL=zh_CN.UTF-8
export LANG=zh_CN.UTF-8
CGO_ENABLED=1 GOOS=darwin go build -ldflags="-s -w" -o build/printer-setup .

if [ $? -ne 0 ]; then
    echo "❌ 编译失败"
    exit 1
fi

echo "✅ 编译完成"

# 复制必要文件到构建目录
echo "📋 复制配置文件..."
cp config.yaml build/
cp README.md build/
cp PROJECT_SUMMARY.md build/
cp run_chinese.sh build/

# 设置可执行文件权限
echo "📝 设置程序权限..."
chmod +x build/printer-setup
chmod +x build/run_chinese.sh

# 创建使用说明
echo "📖 创建使用说明..."
cat > build/使用说明.txt << 'EOF'
HPRT打印机配置工具 - 使用说明
=====================================

📋 文件说明：
- printer-setup         程序主文件（双击运行）
- config.yaml          配置文件（需要修改）
- README.md            详细文档
- 使用说明.txt          本文件

🔧 使用步骤：

1. 修改配置文件
   编辑 config.yaml 文件：
   - 修改VPN名称为你的实际VPN连接名称
   - 修改Windows电脑IP地址

2. 放入驱动文件
   将 hprt-pos-printer-driver-v1.2.16.pkg 文件放入此目录

3. 运行程序
   双击 printer-setup 文件即可启动

⚠️  注意事项：
- 程序需要管理员权限
- 确保网络连接正常
- 确保VPN可以正常连接
- 如有问题请查看详细文档 README.md

🆘 常见问题：
- 如果提示权限不足，请右键选择"以管理员身份运行"
- 如果VPN连接失败，请检查VPN名称是否正确
- 如果打印机无法识别，请检查USB连接和驱动安装

EOF

# 创建发布包
echo "📦 创建发布包..."
cd build
tar -czf "../hprt-printer-setup-v1.0.tar.gz" .
cd ..

echo ""
echo "🎉 构建完成！"
echo ""
echo "📁 构建文件位置："
echo "   - 构建目录: $(pwd)/build/"
echo "   - 发布包: $(pwd)/hprt-printer-setup-v1.0.tar.gz"
echo ""
echo "📋 发布步骤："
echo "1. 将 hprt-printer-setup-v1.0.tar.gz 复制到目标电脑"
echo "2. 解压文件包"
echo "3. 将 hprt-pos-printer-driver-v1.2.16.pkg 放入解压目录"
echo "4. 修改 config.yaml 配置文件"
echo "5. 双击 printer-setup 开始配置"
echo "" 