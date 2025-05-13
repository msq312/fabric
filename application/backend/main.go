package main

import (
	"backend/pkg"
	"backend/router"
	"backend/model"
	con "backend/controller"
	setting "backend/settings"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
	"github.com/spf13/viper"
)

// WebSocket 配置
var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize:  1024,
	CheckOrigin: func(r *http.Request) bool { return true },
}

var clients = make(map[*websocket.Conn]bool)
var broadcast = make(chan string)

// 处理 WebSocket 连接
func handleConnections(w http.ResponseWriter, r *http.Request) {
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Fatal("WebSocket升级失败:", err)
		return
	}
	defer ws.Close()

	clients[ws] = true
	log.Println("新客户端连接")

	// 发送当前状态和下一次执行时间
	sendCurrentStatus(ws)

	// 保持连接
	for {
		_, _, err := ws.ReadMessage()
		if err != nil {
			log.Println("客户端断开连接:", err)
			delete(clients, ws)
			break
		}
	}
}

// 向所有客户端广播消息
func handleMessages() {
	for msg := range broadcast {
		for client := range clients {
			err := client.WriteMessage(websocket.TextMessage, []byte(msg))
			if err != nil {
				log.Printf("发送消息失败: %v", err)
				client.Close()
				delete(clients, client)
			}
		}
	}
}

// 发送当前状态和下一次执行时间
func sendCurrentStatus(ws *websocket.Conn) {
	nextTime := calculateNextExecutionTime()
	msg := fmt.Sprintf("系统正常运行中。下次撮合将于 %s 执行", nextTime.Format("15:04:05"))
	ws.WriteMessage(websocket.TextMessage, []byte(msg))
}

// 计算下次执行时间
func calculateNextExecutionTime() time.Time {
	now := time.Now()
	matchFreq := con.GetMatchFre()
	if matchFreq <= 0 {
		matchFreq = 30 // 默认30分钟
	}
	
	currentMinutes := now.Hour()*60 + now.Minute()
	remainingMinutes := matchFreq - (currentMinutes % matchFreq)
	if remainingMinutes == matchFreq {
		remainingMinutes = 0
	}
	
	nextTime := now.Add(time.Duration(remainingMinutes) * time.Minute)
	if nextTime.Before(now) {
		nextTime = nextTime.Add(24 * time.Hour)
	}
	return nextTime.Truncate(time.Minute)
}

// func periodicTask() {
// 	// 发送开始消息和下次执行时间
// 	nextTime := calculateNextExecutionTime()
// 	startMsg := fmt.Sprintf("撮合报价开始 (下次将于 %s 执行)", nextTime.Format("15:04:05"))
// 	broadcast <- startMsg
// 	log.Println("Periodic task started at", time.Now())
	
// 	// 执行核心逻辑
// 	result := con.ExecutePowerMatching()
// 	log.Println("Matching result:", result)
	
// 	// 发送结束消息和下次执行时间
// 	nextTime = calculateNextExecutionTime() // 重新计算
// 	endMsg := fmt.Sprintf("撮合报价完成 ✅ 结果: %v (下次将于 %s 执行)", result, nextTime.Format("15:04:05"))
// 	broadcast <- endMsg
// }

func periodicTask() {
    // 发送开始消息和下次执行时间
    nextTime := calculateNextExecutionTime()
    startMsg := map[string]interface{}{
        "type":    "info",
        "message": fmt.Sprintf("撮合报价开始 (下次将于 %s 执行)", nextTime.Format("15:04:05")),
    }
    broadcast <- startMsg

    log.Println("Periodic task started at", time.Now())

    // 执行核心逻辑
    result := con.ExecutePowerMatching()
    log.Println("Matching result:", result)

    // 发送结束消息和下次执行时间
    nextTime = calculateNextExecutionTime() // 重新计算
    endMsg := map[string]interface{}{
        "type":    "success",
        "message": fmt.Sprintf("撮合报价完成 ✅ 结果: %v (下次将于 %s 执行)", result, nextTime.Format("15:04:05")),
    }
    broadcast <- endMsg
}

func main() {
	fmt.Println("main.go启动")
	// 初始化配置和服务...
	if err := setting.Init(); err != nil {
		log.Fatal("初始化配置失败:", err)
	}
	if err := pkg.MysqlInit(); err != nil {
		log.Fatal("初始化数据库失败:", err)
	}
	
	r := router.SetupRouter()
	
	// 初始化管理员账户
	if err := con.RegisterAdmin(&model.MysqlUser{
		Username: "admin",
		Password: "123",
        UserType: "管理员",
	}); err != nil {
		log.Println("注册管理员账户失败:", err)
	} else {
		log.Println("管理员账户注册成功")
	}

	// 启动WebSocket服务
	go handleMessages()
	http.HandleFunc("/ws", handleConnections)

	// 启动定时任务
	go startScheduledTask()

	// 启动HTTP服务器
	port := viper.GetInt("app.port")
	log.Printf("服务器启动，监听端口 %d", port)
	if err := r.Run(fmt.Sprintf(":%d", port)); err != nil {
		log.Fatal("启动服务器失败:", err)
	}
}

// 启动定时任务逻辑
func startScheduledTask() {
	log.Println("定时任务协程启动")
	var ticker *time.Ticker
	defer func() {
		if ticker != nil {
			ticker.Stop()
		}
		log.Println("定时任务已停止")
	}()

	// 初始等待，计算首次执行时间
	now := time.Now()
	matchFreq := con.GetMatchFre()
	if matchFreq <= 0 {
		matchFreq = 30 // 默认30分钟
	}
	
	currentMinutes := now.Hour()*60 + now.Minute()
	remainingMinutes := matchFreq - (currentMinutes % matchFreq)
	if remainingMinutes == matchFreq {
		remainingMinutes = 0
	}
	
	nextTime := now.Add(time.Duration(remainingMinutes) * time.Minute)
	if nextTime.Before(now) {
		nextTime = nextTime.Add(24 * time.Hour)
	}
	nextTime = nextTime.Truncate(time.Minute)
	
	log.Printf("首次执行时间: %s", nextTime.Format("15:04:05"))
	time.Sleep(nextTime.Sub(now))

	for {
		// 获取最新频率并更新ticker
		newFreq := con.GetMatchFre()
		if newFreq > 0 && newFreq != matchFreq {
			matchFreq = newFreq
			if ticker != nil {
				ticker.Stop()
			}
			ticker = time.NewTicker(time.Duration(matchFreq) * time.Minute)
			log.Printf("撮合频率已更新为 %d 分钟", matchFreq)
		}

		// 执行任务
		periodicTask()

		// 等待下一个周期
		if ticker == nil {
			ticker = time.NewTicker(time.Duration(matchFreq) * time.Minute)
		}
		<-ticker.C
	}
}