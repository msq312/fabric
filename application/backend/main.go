package main

import (
	"backend/pkg"
	"backend/router"
	//"backend/pkg/mysql"
	"backend/model"
	con "backend/controller"
	setting "backend/settings"
	"fmt"
	"time"
	//"strconv"
	"github.com/spf13/viper"
	//chain "fabric-trace/blockchain/chaincode/chaincode"
)
/////////////////////////////////////////////////////////////
func periodicTask() {
	fmt.Println("Periodic task executed at", time.Now())
	fmt.Println("mathch offers result:", con.ExecutePowerMatching())
}
func main() {
	fmt.Printf("main.go启动\n")
	// 加载配置文件
	if err := setting.Init(); err != nil {
		fmt.Printf("init settings failed,err:%v\n", err)
	}
	fmt.Printf("加载配置文件\n")
	// 初始化数据库
	if err := pkg.MysqlInit(); err != nil {
		fmt.Printf("init mysql failed,err:%v\n", err)
	}
	fmt.Printf("初始化数据库\n")
	
	// 注册路由
	r := router.SetupRouter()
	fmt.Printf("注册路由\n")

	// 定义初始化账号信息
    initUser := &model.MysqlUser{
        UserID:   "",
        Username: "admin",
        Password: "123",
        RealInfo: "",
    }

    // 调用注册函数
    err := con.RegisterAdmin(initUser)
    if err != nil {
        fmt.Printf("Failed to register init user: %v\n", err)
    } else {
        fmt.Println("Init user registered successfully")
    }
	//启动定时器
	// 创建一个每半小时触发一次的 Ticker
	MatchFrequency:=con.GetMatchFre()
	//m,err:=strconv.Atoi(MatchFrequency)
	if err != nil {
        fmt.Println("转换错误:", err)
        //return
    }
	fmt.Printf("MatchFrequency=%d\n",MatchFrequency)
	ticker := time.NewTicker(time.Duration(MatchFrequency)  * time.Minute)
	defer ticker.Stop()
	fmt.Printf("启动计时器\n")
	// 启动一个 Goroutine 来处理定时任务
	go func() {
		for {
			fmt.Printf("另一个线程\n")
			select {
			case <-ticker.C: // 等待 Ticker 的时间信号
				periodicTask() // 执行定时任务
			}
		}
	}()

	// 启动服务
	fmt.Printf("Server is listening on port %d \n", viper.GetInt("app.port"))
	r.Run(fmt.Sprintf(":%d", viper.GetInt("app.port")))
}
