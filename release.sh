#!/bin/bash
set -e
now=$(date +"%Y%m%d%H%M%S")
# 将 dev 版本推送到 latest
docker tag ghcr.io/xyhelper/xyhelper-ucenter-lite:dev ghcr.io/xyhelper/xyhelper-ucenter-lite:latest
docker push ghcr.io/xyhelper/xyhelper-ucenter-lite:latest
# 删除 dev 版本 防止重复提交
docker rmi ghcr.io/xyhelper/xyhelper-ucenter-lite:dev
# 以当前时间为版本
docker tag ghcr.io/xyhelper/xyhelper-ucenter-lite:latest ghcr.io/xyhelper/xyhelper-ucenter-lite:$now
docker push ghcr.io/xyhelper/xyhelper-ucenter-lite:$now
echo "release success" $now
# 写入发布日志 release.log
# echo $now >> ./release.log