package router

// 路由文件
import (
	con "backend/controller"
	"backend/middleware"
	"time"

	"github.com/gin-contrib/cors"//Gin 框架的跨域资源共享（CORS）中间件，用于解决跨域请求的问题。
	"github.com/gin-gonic/gin"//Gin 是一个轻量级的 Web 框架，用于构建 Web 应用和 API。
)

func SetupRouter() *gin.Engine {
//*gin.Engine 实例是 Gin 框架的核心对象，用于处理 HTTP 请求。包含了一些默认的中间件，如日志和恢复中间件
	r := gin.Default()
	// 解决跨域问题
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},                                       // 允许的来源
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE"},            // 允许的请求方法
		AllowHeaders:     []string{"Origin", "Authorization", "Content-Type"}, // 允许的请求头
		ExposeHeaders:    []string{"Content-Length"},                          // 暴露的响应头
		AllowCredentials: true,                                                // 允许传递凭据（例如 Cookie）
		MaxAge:           12 * time.Hour,                                      // 预检请求的有效期
	}))
	// 设置静态文件目录
	r.Static("/static", "./dist/static")//设置静态文件目录，将 /static 路径映射到 ./dist/static 目录，使得客户端可以通过 /static 路径访问该目录下的静态文件。
	r.LoadHTMLGlob("dist/*.html")//加载 dist 目录下的所有 HTML 文件作为模板。
	r.GET("/", func(c *gin.Context) {
		c.HTML(200, "index.html", nil)
	})
	// 测试GET请求
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
	//注册
	r.POST("/register", con.Register)
	//登录
	r.POST("/login", con.Login)
	//登出
	r.POST("/logout", con.Logout)
	//查询用户的类型，当客户端访问 /getInfo 路径时，首先调用 middleware.JWTAuthMiddleware() 中间件进行 JWT 身份验证，验证通过后再调用 con.GetInfo 函数查询用户的类型。
	r.POST("/getUserInfo", middleware.JWTAuthMiddleware(), con.GetUserInfo)
	r.POST("/getAdminInfo", middleware.JWTAuthMiddleware(), con.GetAdminInfo)

	r.POST("/getInfo", middleware.JWTAuthMiddleware(), con.GetInfo)
	r.POST("/getName", middleware.JWTAuthMiddleware(), con.GetName)

	//  用户申请成为买/卖方
	r.POST("/userApproveAs", middleware.JWTAuthMiddleware(), con.UserApproveAs)
	// ApproveUserAs 审核用户申请成为买/卖方
	r.POST("/approveUserAs", middleware.JWTAuthMiddleware(), con.ApproveUserAs)
	//用户报价
	r.POST("/uplink", middleware.JWTAuthMiddleware(), con.Uplink)
	//用户修改报价
	r.POST("/usermodify", middleware.JWTAuthMiddleware(), con.ModifyOffer)
	r.POST("/userCancel", middleware.JWTAuthMiddleware(), con.CancelOffer)

	//用户获取所有的报价信息
	r.POST("/userGetAllOffer", middleware.JWTAuthMiddleware(), con.GetAllOffer)
	//（管理员）获取所有用户正在进行的报价列表
	r.POST("/adminGetAllOffers", middleware.JWTAuthMiddleware(), con.GetAllOffers)
	// 用户获取报价信息溯源
	r.POST("/getOfferHistory", middleware.JWTAuthMiddleware(),con.GetOfferHistory)
	// 查询账户余额历史
	r.POST("/getBalanceHistory", middleware.JWTAuthMiddleware(), con.GetBalanceHistory)
	//  查询用户参与的购电合同
	r.POST("/getUserContracts", middleware.JWTAuthMiddleware(), con.GetUserContracts)


	r.POST("/adminModify", middleware.JWTAuthMiddleware(), con.AdminModify)
	r.POST("/getConfig", middleware.JWTAuthMiddleware(), con.GetConfig)

	//  获取管理员操作历史
	r.POST("/getAdminActionHistory", middleware.JWTAuthMiddleware(), con.GetAdminActionHistory)
	//  获取管理员账户余额历史
	r.POST("/getAdminMoneyHistory", middleware.JWTAuthMiddleware(), con.GetAdminMoneyHistory)
	//查询系统内所有报价
	r.POST("/adminGetAllOffer", middleware.JWTAuthMiddleware(), con.AdminGetAllOffer)
	r.POST("/getAllContract", middleware.JWTAuthMiddleware(), con.GetAllContract)


	return r
}
