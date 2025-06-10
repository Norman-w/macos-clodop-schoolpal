#!/bin/bash

# MacOS校宝打印组件 本地构建脚本
# 注意：这个脚本主要用于本地开发和测试
# 生产环境的构建请使用 GitHub Actions，它会创建真正的跨平台版本
# 支持Intel和ARM版本，包含配置文件和驱动程序

set -e

# 颜色输出
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[0;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# 应用程序名称
APP_NAME="MacOS校宝打印组件"
BUILD_DIR="build"

echo -e "${BLUE}开始构建 ${APP_NAME}...${NC}"

# 1. 清空build目录
echo -e "${YELLOW}清空构建目录...${NC}"
if [ -d "$BUILD_DIR" ]; then
    rm -rf "$BUILD_DIR"
fi
mkdir -p "$BUILD_DIR"

# 2. 创建平台目录
INTEL_DIR="$BUILD_DIR/intel"
ARM_DIR="$BUILD_DIR/arm"
mkdir -p "$INTEL_DIR"
mkdir -p "$ARM_DIR"

# 检测当前系统架构
CURRENT_ARCH=$(uname -m)
echo -e "${BLUE}当前系统架构: ${CURRENT_ARCH}${NC}"

# 3. 构建当前架构版本
if [[ "$CURRENT_ARCH" == "x86_64" ]]; then
    echo -e "${YELLOW}构建Intel版本 (x86_64)...${NC}"
    go build -ldflags="-s -w" -o "$INTEL_DIR/$APP_NAME" .
    if [ $? -eq 0 ]; then
        echo -e "${GREEN}✓ Intel版本构建成功${NC}"
    else
        echo -e "${RED}✗ Intel版本构建失败${NC}"
        exit 1
    fi
    
    echo -e "${YELLOW}⚠ 由于Fyne GUI框架使用CGO，无法交叉编译ARM版本${NC}"
    echo -e "${YELLOW}复制Intel版本到ARM目录（用户需要在ARM Mac上重新编译以获得最佳性能）${NC}"
    # 复制Intel版本作为备用
    cp "$INTEL_DIR/$APP_NAME" "$ARM_DIR/$APP_NAME"
    echo -e "${GREEN}✓ 已复制Intel版本到ARM目录${NC}"
    
elif [[ "$CURRENT_ARCH" == "arm64" ]]; then
    echo -e "${YELLOW}构建ARM版本 (Apple Silicon)...${NC}"
    go build -ldflags="-s -w" -o "$ARM_DIR/$APP_NAME" .
    if [ $? -eq 0 ]; then
        echo -e "${GREEN}✓ ARM版本构建成功${NC}"
    else
        echo -e "${RED}✗ ARM版本构建失败${NC}"
        exit 1
    fi
    
    echo -e "${YELLOW}⚠ 由于Fyne GUI框架使用CGO，无法交叉编译Intel版本${NC}"
    echo -e "${YELLOW}复制ARM版本到Intel目录（用户需要在Intel Mac上重新编译以获得最佳性能）${NC}"
    # 复制ARM版本作为备用
    cp "$ARM_DIR/$APP_NAME" "$INTEL_DIR/$APP_NAME"
    echo -e "${GREEN}✓ 已复制ARM版本到Intel目录${NC}"
    
else
    echo -e "${RED}不支持的架构: ${CURRENT_ARCH}${NC}"
    exit 1
fi

# 5. 复制必需文件到Intel目录
echo -e "${YELLOW}复制必需文件到Intel目录...${NC}"
if [ -f "config.yaml" ]; then
    cp "config.yaml" "$INTEL_DIR/"
    echo -e "${GREEN}✓ 复制 config.yaml${NC}"
else
    echo -e "${RED}⚠ config.yaml 文件不存在${NC}"
fi

if [ -f "hprt-pos-printer-driver-v1.2.16.pkg" ]; then
    cp "hprt-pos-printer-driver-v1.2.16.pkg" "$INTEL_DIR/"
    echo -e "${GREEN}✓ 复制 HPRT驱动程序${NC}"
else
    echo -e "${RED}⚠ HPRT驱动程序文件不存在${NC}"
fi

# 6. 复制必需文件到ARM目录
echo -e "${YELLOW}复制必需文件到ARM目录...${NC}"
if [ -f "config.yaml" ]; then
    cp "config.yaml" "$ARM_DIR/"
    echo -e "${GREEN}✓ 复制 config.yaml${NC}"
else
    echo -e "${RED}⚠ config.yaml 文件不存在${NC}"
fi

if [ -f "hprt-pos-printer-driver-v1.2.16.pkg" ]; then
    cp "hprt-pos-printer-driver-v1.2.16.pkg" "$ARM_DIR/"
    echo -e "${GREEN}✓ 复制 HPRT驱动程序${NC}"
else
    echo -e "${RED}⚠ HPRT驱动程序文件不存在${NC}"
fi

# 7. 创建README文件
echo -e "${YELLOW}创建README文件...${NC}"
cat > "$INTEL_DIR/README.txt" << EOF
MacOS校宝打印组件 - Intel版本 (x86_64)
=======================================

安装说明：
1. 运行 "${APP_NAME}" 启动配置工具
2. 按照步骤完成HPRT打印机配置
3. 如需手动安装驱动，双击 "hprt-pos-printer-driver-v1.2.16.pkg"

系统要求：
- macOS 10.15 或更高版本
- Intel处理器Mac

配置文件：
- config.yaml: 应用程序配置文件

联系方式：
如有问题，请访问：https://github.com/Norman-w/macos-clodop-schoolpal
EOF

cat > "$ARM_DIR/README.txt" << EOF
MacOS校宝打印组件 - ARM版本 (Apple Silicon)
============================================

安装说明：
1. 运行 "${APP_NAME}" 启动配置工具
2. 按照步骤完成HPRT打印机配置
3. 如需手动安装驱动，双击 "hprt-pos-printer-driver-v1.2.16.pkg"

系统要求：
- macOS 11.0 或更高版本
- Apple Silicon处理器Mac (M1/M2/M3)

配置文件：
- config.yaml: 应用程序配置文件

联系方式：
如有问题，请访问：https://github.com/Norman-w/macos-clodop-schoolpal
EOF

# 8. 打包Intel版本
echo -e "${YELLOW}打包Intel版本...${NC}"
cd "$BUILD_DIR"
ZIP_INTEL="${APP_NAME}-Intel-$(date +%Y%m%d).zip"
zip -r "$ZIP_INTEL" intel/ > /dev/null
if [ $? -eq 0 ]; then
    echo -e "${GREEN}✓ Intel版本打包完成: ${ZIP_INTEL}${NC}"
else
    echo -e "${RED}✗ Intel版本打包失败${NC}"
    exit 1
fi

# 9. 打包ARM版本
echo -e "${YELLOW}打包ARM版本...${NC}"
ZIP_ARM="${APP_NAME}-ARM-$(date +%Y%m%d).zip"
zip -r "$ZIP_ARM" arm/ > /dev/null
if [ $? -eq 0 ]; then
    echo -e "${GREEN}✓ ARM版本打包完成: ${ZIP_ARM}${NC}"
else
    echo -e "${RED}✗ ARM版本打包失败${NC}"
    exit 1
fi

cd ..

# 10. 显示构建结果
echo -e "${BLUE}构建完成！${NC}"
echo -e "${GREEN}================================${NC}"
echo -e "${GREEN}构建产物位置：${NC}"
echo -e "  Intel目录: ${BUILD_DIR}/intel/"
echo -e "  ARM目录:   ${BUILD_DIR}/arm/"
echo -e "  Intel包:   ${BUILD_DIR}/${ZIP_INTEL}"
echo -e "  ARM包:     ${BUILD_DIR}/${ZIP_ARM}"
echo -e "${GREEN}================================${NC}"

# 11. 显示文件大小信息
echo -e "${BLUE}文件大小信息：${NC}"
if [ -f "$BUILD_DIR/$ZIP_INTEL" ]; then
    INTEL_SIZE=$(du -h "$BUILD_DIR/$ZIP_INTEL" | cut -f1)
    echo -e "  Intel包大小: ${INTEL_SIZE}"
fi
if [ -f "$BUILD_DIR/$ZIP_ARM" ]; then
    ARM_SIZE=$(du -h "$BUILD_DIR/$ZIP_ARM" | cut -f1)
    echo -e "  ARM包大小:   ${ARM_SIZE}"
fi

echo -e "${GREEN}🎉 ${APP_NAME} 构建成功完成！${NC}" 