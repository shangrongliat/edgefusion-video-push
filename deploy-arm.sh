#!/bin/bash

# 设置编译参数，如果需要的话
export GOOS=linux
export GOARCH=arm64

PROJECT_NAME="video-push"

sleep 1
# 编译Go项目
echo "build $PROJECT_NAME file......."
go build -o build/arm64/$PROJECT_NAME main.go

# 检查编译是否成功
if [ $? -eq 0 ]; then
    echo "Compilation successful."
else
    echo "Compilation failed."
    exit 1
fi

echo "Deployment completed successfully."