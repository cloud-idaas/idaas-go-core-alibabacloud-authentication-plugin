#!/bin/bash
set -e

echo "=========================================="
echo "  IDaaS Integration Test Runner"
echo "=========================================="
echo ""

# Check Alibaba Cloud credentials
if [ -z "$ALIBABA_CLOUD_ACCESS_KEY_ID" ] || [ -z "$ALIBABA_CLOUD_ACCESS_KEY_SECRET" ]; then
    echo "Error: Alibaba Cloud credentials are not set"
    echo ""
    echo "Export them first:"
    echo "  export ALIBABA_CLOUD_ACCESS_KEY_ID='<your AK>'"
    echo "  export ALIBABA_CLOUD_ACCESS_KEY_SECRET='<your SK>'"
    echo ""
    exit 1
fi

echo "OK: Alibaba Cloud credentials are set"

# Configuration file path
CONFIG_FILE=${1:-"config/integration-config.json"}

if [ ! -f "$CONFIG_FILE" ]; then
    echo "Error: config file not found: $CONFIG_FILE"
    echo ""
    echo "Create it first:"
    echo "  1. Copy the template: cp config/integration-config.json.template config/integration-config.json"
    echo "  2. Edit the values:   vim config/integration-config.json"
    echo ""
    exit 1
fi

echo "OK: config file: $CONFIG_FILE"
echo ""

# Show non-sensitive configuration summary
echo "Configuration summary:"
if command -v jq &> /dev/null; then
    echo "  Instance ID: $(jq -r '.idaasInstanceId' "$CONFIG_FILE")"
    echo "  Client ID:   $(jq -r '.clientId' "$CONFIG_FILE")"
    echo "  Scope:       $(jq -r '.scope' "$CONFIG_FILE")"
    echo "  Endpoint:    $(jq -r '.openApiEndpoint' "$CONFIG_FILE")"
else
    echo "  (install jq for a formatted summary; falling back to raw fields)"
    grep -E '"(idaasInstanceId|clientId|scope|openApiEndpoint)"' "$CONFIG_FILE" || true
fi
echo ""

echo "=========================================="
echo "  Running integration tests"
echo "=========================================="
echo ""

go test -v -tags=integration -run TestIntegration

echo ""
echo "=========================================="
echo "  Integration tests finished"
echo "=========================================="
