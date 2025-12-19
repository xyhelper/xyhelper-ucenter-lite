#!/bin/bash
set -e
# 检测是否存在目录 ./backend/resource/public/xyhelper
if [ ! -d "./backend/resource/public/xyhelper" ]; then
    echo "Create directory ./backend/resource/public/xyhelper"
    mkdir -p "./backend/resource/public/xyhelper"

    cd frontend
    yarn build
    cd ..
fi

cd backend
gf build main.go -a amd64 -s linux -p ./temp
gf docker main.go -p -t ghcr.io/xyhelper/xyhelper-ucenter-lite:dev