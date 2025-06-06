#!/bin/bash

# 定义测试线程数
thread_counts=(10 30 50 70)

# 定义测试时长（秒）
test_duration=30

# 循环遍历每个线程数
for thread_count in "${thread_counts[@]}"
do
    echo "Starting stress test with $thread_count threads for $test_duration seconds..."

    # 备份原始配置文件
    cp /home/msq/fabric-trace/blockchain/tape/config_register.yaml /home/msq/fabric-trace/blockchain/tape/config_register.yaml.bak

    # 修改配置文件中的线程数
    sed -i "s/num_of_conn: [0-9]*/num_of_conn: $thread_count/" /home/msq/fabric-trace/blockchain/tape/config_register.yaml
    sed -i "s/client_per_conn: [0-9]*/client_per_conn: 1/" /home/msq/fabric-trace/blockchain/tape/config_register.yaml

    # 计算在 30 秒内需要执行的请求次数
    # 假设每个请求平均耗时 0.1 秒，可根据实际情况调整
    request_count=$((thread_count * test_duration * 10))

    # 执行压力测试
    timeout $test_duration ./tape --config /home/msq/fabric-trace/blockchain/tape/config_register.yaml -n $request_count

    # 恢复原始配置文件
    mv /home/msq/fabric-trace/blockchain/tape/config_register.yaml.bak /home/msq/fabric-trace/blockchain/tape/config_register.yaml

    echo "Stress test with $thread_count threads completed."
    echo "------------------------"
done