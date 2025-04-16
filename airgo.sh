#!/bin/bash
cd server
# 判断www/dist/index.html是否存在
if [ ! -f www/dist/index.html ]; then
  # 创建dist目录
  mkdir -p www/dist
  # 创建index.html文件
    echo "ok" > www/dist/index.html
fi
# 判断www/dist/static是否存在
if [ ! -d www/dist/static ]; then
  # 创建static目录
  mkdir -p www/dist/static
  # 创建index.html文件
    echo "ok" > www/dist/static/index.html
fi
# 判断tmp/main文件是否存在
if [ ! -f tmp/main ]; then
    # 创建tmp目录
    mkdir -p tmp
    # 编译
    go build -o tmp/main ./server.go
fi
air server