## Docker版部署fabric-trace步骤
> 相较于主要的部署方式，不需要安装nvm，且使用docker部署项目更加方便管理





1. 安装docker 

    ```bash
    #下载docker 
    # 官方脚本当前已无法下载，使用gitee备份的脚本:
    curl -fsSL https://gitee.com/real__cool/fabric_install/raw/main/docker_install.sh | bash -s docker --mirror Aliyun
    #添加当前用户到docker用户组 
    sudo usermod -aG docker $USER 
    newgrp docker 
    sudo mkdir -p /etc/docker
    #配置docker镜像加速，近期非常不稳定，如果以下源不好用可以再找下其他源
    #下边的源2024.8.29日测试可用
    sudo tee /etc/docker/daemon.json <<-'EOF'
    {
        "registry-mirrors": [
            "https://docker.m.daocloud.io",
            "https://docker.1panel.live",
            "https://hub.rat.dev"
        ]
    }
    EOF

    #重启docker 
    sudo systemctl daemon-reload
    sudo systemctl restart docker
    ```

2. 安装开发使用的go、node、jq

    ```bash
    #下载二进制包
    wget https://golang.google.cn/dl/go1.19.linux-amd64.tar.gz
    #将下载的二进制包解压至 /usr/local目录
    sudo tar -C /usr/local -xzf go1.19.linux-amd64.tar.gz
    mkdir $HOME/go
    #将以下内容添加至环境变量 ~/.bashrc
    export GOPATH=$HOME/go
    export GOROOT=/usr/local/go
    export PATH=$GOROOT/bin:$PATH
    export PATH=$GOPATH/bin:$PATH
    #更新环境变量
    source  ~/.bashrc 
    #设置代理
    go env -w GO111MODULE=on
    go env -w GOPROXY=https://goproxy.cn,direct
    # 更新环境变量
    source  ~/.bashrc
    #安装jq 
    sudo apt install jq
    ```



3. 克隆本项目 

    ```bash
    git clone https://gitee.com/real__cool/fabric-trace
    ```

4. 启动区块链部分。在fabric-trace/blockchain/network目录下:

    ```bash
    # 仅在首次使用执行：下载Fabric Docker镜像。如果拉取速度过慢或失败请检查是否完成docker换源，或者更换一个其他的镜像源再试。
    ./install-fabric.sh -f 2.5.6 d 
    ```
    ```bash
    # 启动区块链网络
    ./start.sh
    ```	
     **如果在启动区块链网络时遇到报错可以尝试:**
    ```bash
    # 执行清理所有的容器指令：
    docker rm -f $(docker ps -aq)
    ```
    **然后再重新启动区块链网络**

5. 修改后端IP，将以下文件中的IP：`119.45.247.29`，换成自己云服务器的IP。
    ```bash
    fabric-trace/application/web/.env.production
    fabric-trace/application/web/.env.development
    fabric-trace/application/web/src/router/index.js
    ```
    或使用application/replaceip.sh脚本根据指引修改IP，在fabric-trace/application目录下
    ```bash
    ./replaceip.sh
    ```

6. 启动app，在fabric-trace/application目录执行： 

    ```bash
    ./stop_docker.sh
    ```


7. 关闭项目


    ```bash
    cd application
    #关闭前后端
    ./stop_docker.sh
    #关闭区块链
    cd blockchain/network
    ./stop.sh
    ```


