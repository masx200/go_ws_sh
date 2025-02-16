#!/bin/bash
# 无限循环脚本示例

while true
do
    #echo "脚本正在运行... 按 [CTRL+C] 停止"
    # 这里可以放置你想要重复执行的命令
    # 例如：./your_program 或者其他命令
    ./main -config server-config-tls.json -mode server
    
    # 可选：让循环每秒执行一次（或任何其他时间间隔）
    #sleep 1
done