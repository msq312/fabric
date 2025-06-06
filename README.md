安装**Ubuntu20.04**系统
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
	#下边的源2024.12.26日测试可用
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
	
	#下载nvm安装脚本
	wget https://gitee.com/real__cool/fabric_install/raw/main/nvminstall.sh
	#安装nvm；屏幕输出内容添加环境变量
	chmod +x nvminstall.sh
	./nvminstall.sh
	# 将环境变量写入.bashrc
	export NVM_DIR="$HOME/.nvm"
	[ -s "$NVM_DIR/nvm.sh" ] && \. "$NVM_DIR/nvm.sh"  # This loads nvm
	[ -s "$NVM_DIR/bash_completion" ] && \. "$NVM_DIR/bash_completion"  # This loads nvm bash_completion
	export NVM_NODEJS_ORG_MIRROR=http://npmmirror.com/mirrors/node/ #更换阿里云nvm node源
 
	# 更新环境变量
	source  ~/.bashrc
	# 安装node16
	nvm install 16
	#换源
	npm config set registry https://registry.npmmirror.com
	
	#安装jq 
	sudo apt install jq
	```

3. 克隆本项目 

	```bash
	git clone https://gitee.com/real__cool/fabric-trace
	```

4. 启动区块链部分。在fabric-trace/blockchain/network目录下:

	```bash
	# 仅在首次使用执行：下载Fabric Docker镜像。如果拉取速度过慢或失败请检查是否完成docker换源并执行了重启docker命令。
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

5. 启动后端 在fabric-trace/application/backend目录下： 执行： `go run main.go`

6. 修改后端IP，将以下文件中的IP：`119.45.247.29`，换成自己云服务器的IP。
	```bash
	fabric-trace/application/web/.env.production
	fabric-trace/application/web/.env.development
	fabric-trace/application/web/src/router/index.js
	```
	或使用application/replaceip.sh脚本根据指引修改IP，在fabric-trace/application目录下
	```bash
	./replaceip.sh
	```

7. 新开一个窗口，启动前端 在fabric-trace/application/web目录下： 执行： 

	```bash
	# 仅在首次运行执行：安装依赖
	npm install 
	```

	```bash
	# 启动前端
	npm run dev
	```


#### 关闭项目与重新运行步骤
##### 关闭项目：
1. 前端（`npm run dev`界面）与后端（`go run main.go`界面：

	使用键盘组合键：`ctrl+c`

2. 区块链部分：

	在`fabric-trace/blockchain/network`目录`./stop.sh`，此步骤会清理所有的区块链、Mysql中的数据。

##### 开发模式启动项目：
1. 在`fabric-trace/blockchain/network`目录
`./start.sh` 如果遇到报错可以执行以下命令后再试：
执行清理所有的容器指令：
`docker rm -f $(docker ps -aq)`
2. 在`fabric-trace/application/backend`目录下： 执行： `go run main.go`
3. 在`fabric-trace/application/web`目录下： 执行：
`npm run dev`
4. 在http://服务器IP:9528打开

##### 生产环境部署项目(后台运行，访问速度更快)
1. 在`fabric-trace/blockchain/network`目录
`./start.sh` 如果遇到报错可以执行以下命令后再试：
执行清理所有的容器指令：
`docker rm -f $(docker ps -aq)`
2. 在`fabric-trace/application`目录下： 执行： `./start_prod.sh`
3. 在http://服务器IP:9090打开

	注意：此方式部署项目会在后台运行，如果后续遇到端口号占用可以尝试关闭占用9090端口号的进程，可以参考：
	[解决端口占用 bind:address already in use](https://blog.csdn.net/qq_41575489/article/details/137434008?spm=1001.2014.3001.5501)

#### 特别提示：
1. 使用虚拟机时区块链浏览器有时候会出现无法访问的情况，可以尝试重启浏览器容器。
2. 为了减少用户运行本项目时的难度，区块链目录的start.sh脚本在启动区块链时同时会清理掉所有的历史数据！如果重启机器后不希望清理原来的数据启动区块链，可以使用指令：`docker start $(docker ps -aq)`启动所有节点

