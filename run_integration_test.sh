#!/bin/bash
set -e

echo "=========================================="
echo "  IDaaS 集成测试运行脚本"
echo "=========================================="
echo ""

# 检查阿里云凭证
if [ -z "$ALIBABA_CLOUD_ACCESS_KEY_ID" ] || [ -z "$ALIBABA_CLOUD_ACCESS_KEY_SECRET" ]; then
    echo "错误: 阿里云凭证未设置"
    echo ""
    echo "请先设置环境变量:"
    echo "  export ALIBABA_CLOUD_ACCESS_KEY_ID='你的AK'"
    echo "  export ALIBABA_CLOUD_ACCESS_KEY_SECRET='你的SK'"
    echo ""
    exit 1
fi

echo "✓ 阿里云凭证已设置"

# 配置文件路径
CONFIG_FILE=${1:-"config/integration-config.json"}

if [ ! -f "$CONFIG_FILE" ]; then
    echo "错误: 配置文件不存在: $CONFIG_FILE"
    echo ""
    echo "请先创建配置文件:"
    echo "  1. 复制模板: cp config/integration-config.json.template config/integration-config.json"
    echo "  2. 编辑配置: vim config/integration-config.json"
    echo ""
    exit 1
fi

echo "✓ 配置文件: $CONFIG_FILE"
echo ""

# 显示配置信息
echo "配置信息:"
if command -v jq &> /dev/null; then
    echo "  Instance ID: $(jq -r '.idaasInstanceId' "$CONFIG_FILE")"
    echo "  Client ID: $(jq -r '.clientId' "$CONFIG_FILE")"
    echo "  Scope: $(jq -r '.scope' "$CONFIG_FILE")"
    echo "  Endpoint: $(jq -r '.openApiEndpoint' "$CONFIG_FILE")"
else
    echo "  (安装 jq 可以显示格式化的配置信息)"
fi
echo ""

echo "=========================================="
echo "  开始运行集成测试"
echo "=========================================="
echo ""

go test -v -tags=integration -run TestIntegration

echo ""
echo "=========================================="
echo "  集成测试完成"
echo "=========================================="
