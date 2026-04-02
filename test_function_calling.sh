#!/bin/bash

# 测试 Function Calling

echo "=== 测试 1: 普通对话（不调用工具） ==="
echo "你好" | ./myopenclaw

echo ""
echo "=== 测试 2: 请求调用工具 ==="
echo "帮我 echo 一下 hello world" | ./myopenclaw

echo ""
echo "=== 测试 3: 明确要求使用 echo 工具 ==="
echo "使用 echo 工具输出：测试成功" | ./myopenclaw
