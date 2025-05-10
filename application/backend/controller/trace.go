package controller

import (
	"backend/pkg"
	"fmt"
	"encoding/json"
	"github.com/gin-gonic/gin"
)
type Config struct{
    MatchFrequency   int   `json:"matchFrequency"` // 撮合频率，单位为分钟
    DepositRate      float64           `json:"depositRate"`       // 保证金率
    FeeRate    float64           `json:"feeRate"`       // 手续费率
}
// UserApproveAs 用户申请成为买/卖方
func UserApproveAs(c *gin.Context) {
	userID, _ := c.Get("userID")
	var args []string
	fmt.Println("into UserApproveAs")
	appId:= pkg.GenerateID()[1:]
	args=append(args,ADMINID,userID.(string),appId,c.PostForm("status"))
	fmt.Println("args:",args)
	res, err := pkg.ChaincodeInvoke("UserApproveAs", args)
	fmt.Println("ok UserApproveAs")
	if err != nil {
		c.JSON(200, gin.H{
			"message": "查询失败：" + err.Error(),
		})
		return
	}
	fmt.Println("query success")
	c.JSON(200, gin.H{
		"code":    200,
		"message": "query success",
		"data":    res,
	})
}
// ApproveUserAs 审核用户申请成为买/卖方
func ApproveUserAs(c *gin.Context) {
	//userID,_ := c.Get("userID")
	var args []string
	fmt.Println("into ApproveUserAs")
	
	args=append(args,ADMINID,c.PostForm("arg1"),c.PostForm("arg2"))

	fmt.Println("ok ApproveUserAs")
	res, err := pkg.ChaincodeInvoke("ApproveUserAs", args)
	if err != nil {
		c.JSON(200, gin.H{
			"message": "查询失败：" + err.Error(),
		})
		return
	}
	c.JSON(200, gin.H{
		"code":    200,
		"message": "query success",
		"data":    res,
	})
}
//ExecutePowerMatching执行电力撮合
func ExecutePowerMatching() string{	
	_,err := pkg.ChaincodeInvoke("MatchOffers", []string{ADMINID})
	if err != nil {
		// c.JSON(200, gin.H{
		// 	"message": "MatchOffers failed" + err.Error(),
		// })
		return "MatchOffers failed" + err.Error()
	}
	fmt.Println("exe well done matchoffers")
	_,err=pkg.ChaincodeInvoke("SettleContract", []string{ADMINID})
	if err != nil {
		// c.JSON(200, gin.H{
		// 	"message": "SettleContract failed" + err.Error(),
		// })
		return "SettleContract failed" + err.Error()
	}
	return "exe well done settlecontract"
}
// 用户报价
func Uplink(c *gin.Context) {
	// 与userID不一样，取ID从第二位作为溯源码，即18位，雪花ID的结构如下（转化为10进制共19位）：
	// +--------------------------------------------------------------------------+
	// | 1 Bit Unused | 41 Bit Timestamp |  10 Bit NodeID  |   12 Bit Sequence ID |
	// +--------------------------------------------------------------------------+
	offerId := pkg.GenerateID()[1:]
	args := buildArgs(c,offerId,true)
	if args == nil {
		return
	}
	res, err := pkg.ChaincodeInvoke("SubmitOffer", args)
	if err != nil {
		c.JSON(200, gin.H{
			"message": "uplink failed" + err.Error(),
		})
		return
	}
	c.JSON(200, gin.H{
		"code":              200,
		"message":           "uplink success",
		"data":              res,
		"offerId":			 args[1],
	})
}

//用户修改报价信息
func ModifyOffer(c *gin.Context) {
	//offerId string, userId
	args := buildArgs(c,c.PostForm("offerId"),false)
	if args == nil {
		return
	}
	_,err := pkg.ChaincodeInvoke("ModifyOffer", args)
	if err != nil {
		c.JSON(200, gin.H{
			"message": "修改报价失败：" + err.Error(),
		})
		return
	}
	c.JSON(200, gin.H{
		"code":    200,
		"message": "query success",
		//"data":    res,
	})

}

// （管理员）获取所有用户正在进行的报价列表
func GetAllOffers(c *gin.Context) {
	//userID, _ := c.Get("userID")
	res, err := pkg.ChaincodeQuery("GetAllOffers","")
	if err != nil {
		c.JSON(200, gin.H{
			"message": "query failed" + err.Error(),
		})
	}
	c.JSON(200, gin.H{
		"code":    200,
		"message": "query success",
		"data":    res,
	})
}

// 用户获取所有的报价信息
func GetAllOffer(c *gin.Context) {
	userID, _ := c.Get("userID")
	res, err := pkg.ChaincodeQuery("GetAllOffer", userID.(string))
	fmt.Print("res", res)
	if err != nil {
		c.JSON(200, gin.H{
			"message": "query failed" + err.Error(),
		})
	}
	c.JSON(200, gin.H{
		"code":    200,
		"message": "query success",
		"data":    res,
	})
}

// GetOfferHistory 用户查询所有报价历史信息
func GetOfferHistory(c *gin.Context) {
	userID, _ := c.Get("userID")
	res, err := pkg.ChaincodeQuery("GetOfferHistory", userID.(string))
	if err != nil {
		c.JSON(200, gin.H{
			"message": "query failed" + err.Error(),
		})
	}
	c.JSON(200, gin.H{
		"code":    200,
		"message": "query success",
		"data":    res,
	})
}
//GetBalanceHistory 查询账户余额历史
func GetBalanceHistory(c *gin.Context) {
	userID, _ := c.Get("userID")
	res, err := pkg.ChaincodeQuery("GetBalanceHistory", userID.(string))
	if err != nil {
		c.JSON(200, gin.H{
			"message": "query failed" + err.Error(),
		})
	}
	c.JSON(200, gin.H{
		"code":    200,
		"message": "query success",
		"data":    res,
	})
}
// GetUserContracts 查询用户参与的购电合同
func GetUserContracts(c *gin.Context) {
	userID, _ := c.Get("userID")
	res, err := pkg.ChaincodeQuery("GetUserContracts", userID.(string))
	if err != nil {
		c.JSON(200, gin.H{
			"message": "query failed" + err.Error(),
		})
	}
	c.JSON(200, gin.H{
		"code":    200,
		"message": "query success",
		"data":    res,
	})
}
func GetMatchFre()int{
	res, err := pkg.ChaincodeQuery("GetConfig", "")
	if err != nil {
		fmt.Printf(err.Error())
		return -1
	}
	var config Config
	err = json.Unmarshal([]byte(res), &config)
	if err != nil {
		fmt.Printf("Error unmarshaling response: %s\n", err.Error())
		return -1
	}
	return config.MatchFrequency
}
func GetConfig(c *gin.Context){
	res, err := pkg.ChaincodeQuery("GetConfig", "")
	if err != nil {
		c.JSON(200, gin.H{
			"message": "获取失败：" + err.Error(),
		})
	}
	c.JSON(200, gin.H{
		"code":    200,
		"message": "获取成功",
		"data":    res,
	})
}
//// AdminModify 管理员修改配置参数
func AdminModify(c *gin.Context) {
	var args []string
	args=append(args, ADMINID,c.PostForm("name"),c.PostForm("newConfig"))
	fmt.Println(c.PostForm("name"),c.PostForm("newConfig"))
	res, err := pkg.ChaincodeInvoke("AdminModify",args)
	if err != nil {
		c.JSON(200, gin.H{
			"message": "修改失败：" + err.Error(),
		})
	}
	c.JSON(200, gin.H{
		"code":    200,
		"message": "修改成功",
		"data":    res,
	})
}
// GetAdminActionHistory 获取管理员操作历史
func GetAdminActionHistory(c *gin.Context) {
	res, err := pkg.ChaincodeQuery("GetAdminActionHistory", ADMINID)
	if err != nil {
		c.JSON(200, gin.H{
			"message": "query failed" + err.Error(),
		})
	}
	c.JSON(200, gin.H{
		"code":    200,
		"message": "query success",
		"data":    res,
	})
}
// GetAdminMoneyHistory 获取管理员账户余额历史
func GetAdminMoneyHistory(c *gin.Context) {
	res, err := pkg.ChaincodeQuery("GetAdminMoneyHistory", ADMINID)
	if err != nil {
		c.JSON(200, gin.H{
			"message": "query failed" + err.Error(),
		})
	}
	c.JSON(200, gin.H{
		"code":    200,
		"message": "query success",
		"data":    res,
	})
}
//查询系统内所有报价AdminGetAllOffer
func AdminGetAllOffer(c *gin.Context) {
	res, err := pkg.ChaincodeQuery("AdminGetAllOffer", ADMINID)
	if err != nil {
		c.JSON(200, gin.H{
			"message": "query failed" + err.Error(),
		})
	}
	c.JSON(200, gin.H{
		"code":    200,
		"message": "query success",
		"data":    res,
	})
}
//查询系统内所有合同GetAllContract
func GetAllContract(c *gin.Context) {
	res, err := pkg.ChaincodeQuery("GetAllContract", ADMINID)
	if err != nil {
		c.JSON(200, gin.H{
			"message": "query failed" + err.Error(),
		})
	}
	c.JSON(200, gin.H{
		"code":    200,
		"message": "query success",
		"data":    res,
	})
}
func buildArgs(c *gin.Context,offerId string,create bool) []string {
	var args []string
	userID, _ := c.Get("userID")
	//userType, _ := pkg.ChaincodeQuery("GetUserType", userID.(string))
	args = append(args, userID.(string))
	fmt.Println(userID)
	
	// // 检查溯源码是否正确
	// res1, res2,err := pkg.ChaincodeQuery("GetAllOffers", offerId)
	// if res == "" || err != nil || len(offerId) != 18 {
	// 	c.JSON(200, gin.H{
	// 		"message": "请检查溯源码是否正确!!",
	// 	})
	// 	return nil
	// } else {
	// 	args = append(args, offerId)
	// }
	args = append(args, offerId)	
	if create{
		args = append(args, c.PostForm("arg1"))
		args = append(args, c.PostForm("arg2"))
		args = append(args, c.PostForm("arg3"))
	}else{
		args = append(args, c.PostForm("arg1"))
		args = append(args, c.PostForm("arg2"))
	}
	
	//args = append(args, c.PostForm("arg4"))
	//adminid 
	args = append(args, ADMINID)
	return args
}
