package controller

import (
	"backend/model"
	"backend/pkg"
	"fmt"
	//"time"
	//"strconv"
	"github.com/gin-gonic/gin"
)
var AdminCount=0
var ADMINID string
// func periodicTask() {
// 	fmt.Println("Periodic task executed at", time.Now())
// 	ExecutePowerMatching()
// }
// func startPeriodicTask() {
// 	// 创建一个每半小时触发一次的 Ticker
// 	MatchFrequency:=GetMatchFre()
// 	m,err:=strconv.Atoi(MatchFrequency)
// 	if err != nil {
//         fmt.Println("转换错误:", err)
//         return
//     }
// 	fmt.Printf("MatchFrequency=%d\n",m)
// 	ticker := time.NewTicker(time.Duration(m)  * time.Minute)
// 	defer ticker.Stop()
// 	fmt.Printf("启动计时器\n")
// 	// 启动一个 Goroutine 来处理定时任务
// 	go func() {
// 		for {
// 			fmt.Printf("另一个线程\n")
// 			select {
// 			case <-ticker.C: // 等待 Ticker 的时间信号
// 				periodicTask() // 执行定时任务
// 			}
// 		}
// 	}()
// }
func RegisterAdmin(user *model.MysqlUser)error{
	user.UserID = "1917140260131704832"//pkg.GenerateID()//"1916567760880537600"
	user.RealInfo = pkg.EncryptByMD5(user.Username)
	user.UserType="管理员"
	if user.UserType=="管理员"{
		AdminCount+=1
		if AdminCount==1{
			ADMINID=user.UserID
			//startPeriodicTask()
		}
	}
	
	err := pkg.InsertUser(user)
	if err != nil {
		fmt.Printf("admin id=%s\n",ADMINID)
		fmt.Println("管理员注册插入mysql失败")
		return err
	}
	var args []string
	args = append(args, user.UserID)
	_,err = pkg.ChaincodeInvoke("RegisterPlatformAdmin", args)
	if err != nil {
		fmt.Println("区块链注册管理员失败")
		return err
	}
	return err
}
func Register(c *gin.Context) {
	// 将用户信息存入mysql数据库
	var user model.MysqlUser
	user.UserID = pkg.GenerateID()
	user.Username = c.PostForm("username")
	user.Password = c.PostForm("password")
	user.RealInfo = pkg.EncryptByMD5(c.PostForm("username"))
	user.UserType=c.PostForm("userType")
	if user.UserType=="管理员"{
		AdminCount+=1
		if AdminCount==1{
			ADMINID=user.UserID
			//startPeriodicTask()
		}
		if AdminCount>1 {
			c.JSON(200, gin.H{
				"message": "register failed：管理员已存在！请登录" ,
			})
			return
		}
	}
	err := pkg.InsertUser(&user)
	if err != nil {
		c.JSON(200, gin.H{
			"message": "register failed：" + err.Error(),
		})
		return
	}
	// 将用户信息存入区块链
	// userID string, userType string, realInfoHash string
	// 将post请求的参数封装成一个数组args
	var args []string
	args = append(args, user.UserID,user.Username)
	//args = append(args, c.PostForm("userType"))
	//args = append(args, user.RealInfo)
	var res string 
	if user.UserType!="管理员"{
		_,err := pkg.ChaincodeInvoke("RegisterUser", args)
		if err != nil {
			c.JSON(200, gin.H{
				"message": "register failed：" + err.Error(),
			})
			return
		}
	}else{
		_,err := pkg.ChaincodeInvoke("RegisterPlatformAdmin", args)
		if err != nil {
			c.JSON(200, gin.H{
				"message": "register failed：" + err.Error(),
			})
			return
		}
		if AdminCount==1{
			//startPeriodicTask()
		}
	}
	
	c.JSON(200, gin.H{
		"code":    200,
		"message": "register success",
		"txid":    res,
	})
}

func Login(c *gin.Context) {
	var user model.MysqlUser
	user.Username = c.PostForm("username")
	user.Password = c.PostForm("password")
	// 获取用户ID
	var err error
	user.UserID, err = pkg.GetUserID(user.Username)
	if err != nil {
		c.JSON(200, gin.H{
			"message": "没有找到该用户",
		})
		return
	}
	user.UserType, err = pkg.GetUserType(user.Username)
	if err != nil {
		c.JSON(200, gin.H{
			"message": "没有找到用户类型" ,
		})
		return
	}
	err = pkg.Login(&user)
	if err != nil {
		c.JSON(200, gin.H{
			"message": "login failed:" + err.Error(),
		})
		return
	}

	// 生成jwt, userType
	jwt, err := pkg.GenToken(user.UserID,user.UserType)
	if err != nil {
		c.JSON(200, gin.H{
			"message": "login failed:" + err.Error(),
		})
		return
	}
	c.JSON(200, gin.H{
		"code":    200,
		"message": "login success",
		"jwt":     jwt,
		"userType": user.UserType, 
	})
}

func Logout(c *gin.Context) {
	c.JSON(200, gin.H{
		"code":    200,
		"message": "logout success",
	})
}

//获取用户所有信息
func GetUserInfo(c *gin.Context) {
	userID, _ := c.Get("userID")
	fmt.Println("into getuserinfo")
	userInfo, err := pkg.ChaincodeQuery("GetUserInfo", userID.(string))
	if err != nil {
		fmt.Println("getuserinfo wrong")
		c.JSON(200, gin.H{
			"message": "get userinfo failed" + err.Error(),
		})
	}
	fmt.Println("getuserinfo ok")
	c.JSON(200, gin.H{
		"code":     200,
		"message":  "get userinfo success",
		"data": userInfo,
	})
}
//获取管理员所有信息
func GetAdminInfo(c *gin.Context) {
	fmt.Println("into getadmininfo")
	userInfo, err := pkg.ChaincodeQuery("GetAdminInfo", ADMINID)
	if err != nil {
		fmt.Println("getadmininfo wrong")
		c.JSON(200, gin.H{
			"message": "get admininfo failed" + err.Error(),
		})
	}
	fmt.Println("getadmininfo ok")
	c.JSON(200, gin.H{
		"code":     200,
		"message":  "get admininfo success",
		"data": userInfo,
	})
}
func GetName(c *gin.Context){
	username,err:=pkg.GetUsername(c.PostForm("id"))
	if err != nil {
		c.JSON(200, gin.H{
			"message": "get user name failed" + err.Error(),
		})
	}
	c.JSON(200, gin.H{
		"code":     200,
		"message":  "get user type success",
		"data": username,
	})
}
// 获取用户信息
func GetInfo(c *gin.Context) {
	userID, _ := c.Get("userID")
	// if !exist {
	// 	c.JSON(200, gin.H{
	// 		"message": "get user ID failed",
	// 	})
	// }

	// userType, err := pkg.GetUserType(userID.(string))
	// if err != nil {
	// 	c.JSON(200, gin.H{
	// 		"message": "get user type failed" + err.Error(),
	// 	})
	// }

	username, err := pkg.GetUsername(userID.(string))
	userType,err:=pkg.GetUserType(username)
	if err != nil {
		c.JSON(200, gin.H{
			"message": "get user name failed" + err.Error(),
		})
	}

	c.JSON(200, gin.H{
		"code":     200,
		"message":  "get user type success",
		"userType": userType,
		"username": username,
	})
}
