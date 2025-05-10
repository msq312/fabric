package chaincode

import (
	"encoding/json"
	"fmt"
	"time"
	"strconv"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
	//"github.com/google/uuid"
)

// SmartContract 定义合约结构体
type SmartContract struct {
	contractapi.Contract
}

// RegisterUser 注册用户
func (s *SmartContract) RegisterUser(ctx contractapi.TransactionContextInterface, userId,username string) error {
	user := &User{
		UserID:            userId,
		UserName:          username,
		Balance:           0,
        IsSeller:          NotApplied,
        IsBuyer:           NotApplied,
		Offers:         []string{},
		Contracts:      []string{},
		BalanceHistory: []string{},
		OfferHistory:   []string{},
		OfferDone:      []string{},
        CreditRating:      0,
        TradeCount:        0,
	}
	userJSON, err := json.Marshal(user)
	if err != nil {
		return err
	}
	return ctx.GetStub().PutState(userId, userJSON)
}
// getUser 获取用户信息
func (s *SmartContract) getUser(ctx contractapi.TransactionContextInterface, userId string) (*User, error) {
	userJSON, err := ctx.GetStub().GetState(userId)
	if err != nil {
		return nil, fmt.Errorf("获取用户信息时出错: %w", err)
	}
	if userJSON == nil {
		return nil, fmt.Errorf("用户 %s 不存在", userId)
	}
	var user User
	if err := json.Unmarshal(userJSON, &user); err != nil {
		return nil, err
	}
	return &user, nil
}

// getAdmin 获取管理员信息
func (s *SmartContract) getAdmin(ctx contractapi.TransactionContextInterface, adminId string) (*PlatformAdmin, error) {
	adminJSON, err := ctx.GetStub().GetState(adminId)
	if err != nil {
		return nil, fmt.Errorf("获取管理员信息时出错: %v", err)
	}
	if adminJSON == nil {
		return nil, fmt.Errorf("管理员 %s 不存在", adminId)
	}
	var admin PlatformAdmin
	if err := json.Unmarshal(adminJSON, &admin); err != nil {
		return nil, err
	}
	return &admin, nil
}

// updateUser 更新用户信息
func (s *SmartContract) updateUser(ctx contractapi.TransactionContextInterface, user *User) error {
	userJSON, err := json.Marshal(user)
	if err != nil {
		return err
	}
	return ctx.GetStub().PutState(user.UserID, userJSON)
}

// updateAdmin 更新管理员信息
func (s *SmartContract) updateAdmin(ctx contractapi.TransactionContextInterface, admin *PlatformAdmin) error {
	adminJSON, err := json.Marshal(admin)
	if err != nil {
		return err
	}
	return ctx.GetStub().PutState(admin.AdminID, adminJSON)
}
// UserApproveAs 用户购电/售电资质申请
func (s *SmartContract) UserApproveAs(ctx contractapi.TransactionContextInterface, adminId, userId, appId,applyType string) error {
	fmt.Printf("into\n")
	user, err := s.getUser(ctx, userId)
	if err != nil {
		fmt.Printf("Error: 找不到用户 %s\n", userId)
		return fmt.Errorf("找不到用户 %s", userId)
	}

	admin, err := s.getAdmin(ctx, adminId)
	if err != nil {
		fmt.Printf("Error: 找不到管理员 %s\n", adminId)
		return fmt.Errorf("找不到管理员 %s", adminId)
	}
    // 检查现有状态
    currentStatus := user.getCurrentStatus(applyType)
    if currentStatus != NotApplied && currentStatus != Rejected {
		fmt.Printf("Error: 存在进行中的申请，无法重复提交\n")
        return fmt.Errorf("存在进行中的申请，无法重复提交")
    }
    // 创建申请记录
    application := &Application{
        ApplicationID: appId,
        UserID:        userId,
		UserName:          user.UserName,
        ApplyType:     applyType,
        ApplyTime:     time.Now().Format(customFormat),
        AuditStatus:   AuditPending,
        AuditTime:   time.Now().Format(customFormat),
    }
    admin.Applications = append(admin.Applications, application)
	// 更新用户状态
    user.updateStatus(applyType, Pending)
	if err := s.updateUser(ctx, user); err != nil {
		fmt.Printf("Error: 无法更新用户 %s 的信息: %v\n", userId, err)
		return fmt.Errorf("无法更新用户 %s 的信息: %v", userId, err)
	}
	if err := s.updateAdmin(ctx, admin); err != nil {
		fmt.Printf("Error: 无法更新gly %s 的信息: %v\n", adminId, err)
		return fmt.Errorf("无法更新用户 %s 的信息: %v", adminId, err)
	}
	fmt.Printf("well done\n")
	//log.Println("SomeFunction 函数执行完毕")
	return nil//s.updateAdmin(ctx, admin)
}
// 获取用户当前状态
func (u *User) getCurrentStatus(applyType string) ApprovalStatus {
    if applyType == "sell" {
        return u.IsSeller
    }
    return u.IsBuyer
}
// 更新用户状态
func (u *User) updateStatus(applyType string, status ApprovalStatus) {
    if applyType == "sell" {
        u.IsSeller = status
    } else {
        u.IsBuyer = status
		u.Balance=INIT_BALANCE
    }
}
// RegisterPlatformAdmin 注册平台管理员
func (s *SmartContract) RegisterPlatformAdmin(ctx contractapi.TransactionContextInterface, adminId string) error {
	admin := &PlatformAdmin{
		AdminID:          adminId,
		Balance:          0,
		BalanceHistory:   []string{},
		AdminActionHistory: []string{},
        Applications: []*Application{},
		Contracts:        []string{},
		Contractnumber:   0,
	}
	adminJSON, err := json.Marshal(admin)
	if err != nil {
		return err
	}
	ctx.GetStub().PutState("MatchFrequency", []byte(strconv.Itoa(MatchFrequency)))
	ctx.GetStub().PutState("DepositRate", []byte(strconv.FormatFloat(DepositRate, 'f', 2, 64)))
	ctx.GetStub().PutState("FeeRate", []byte(strconv.FormatFloat(FeeRate, 'f', 2, 64)))
	return ctx.GetStub().PutState(adminId, adminJSON)
}



// ApproveUserAs 管理员审核通过用户成为买/卖方
func (s *SmartContract) ApproveUserAs(ctx contractapi.TransactionContextInterface, adminId string, applicationId string, approve bool) error {
	admin, err := s.getAdmin(ctx, adminId)
	if err != nil {
		return err
	}
    // 查找申请记录
    app, _ := findApplication(admin.Applications, applicationId)
    if app == nil {
        return fmt.Errorf("申请记录不存在")
    }
	
    // 更新申请记录
    if approve {
        app.AuditStatus = AuditPassed
    } else {
        app.AuditStatus = AuditRejected
    }
    app.AuditTime = time.Now().Format(customFormat)
    // 更新用户状态
    user, _ := s.getUser(ctx, app.UserID)
    newStatus := Approved
    if !approve {
        newStatus = Rejected
    }
    user.updateStatus(app.ApplyType, newStatus)
	if err := s.updateUser(ctx, user); err != nil {
		return err
	}
	//admin.AdminActionHistory = append(admin.AdminActionHistory, adminAction)
	return s.updateAdmin(ctx, admin)
}
// 查找申请记录
func findApplication(apps []*Application, id string) (*Application, int) {
    for i, app := range apps {
        if app.ApplicationID == id {
            return app, i
        }
    }
    return nil, -1
}



// AdminModify 管理员修改配置参数
func (s *SmartContract) AdminModify(ctx contractapi.TransactionContextInterface, adminId string, name string, newConfig string) error {
    // 获取管理员信息
	fmt.Println("into AdminModify ")
    admin, err := s.getAdmin(ctx, adminId)
    if err != nil {
		fmt.Printf("Error: 获取管理员信息失败: %v\n", err)
        return fmt.Errorf("获取管理员信息失败: %w", err)
    }

    var actionRecord *AdminActionRecord
    var updateSuccess bool

    switch name {
    case "MatchFrequency":
        i, err := strconv.Atoi(newConfig)
        if err != nil {
			fmt.Printf("Error: 转换撮合频率配置值失败: %v\n", err)
            return fmt.Errorf("转换撮合频率配置值失败: %w", err)
        }
        actionRecord = &AdminActionRecord{
            Action:    "修改撮合频率",
            Timestamp: time.Now().Format(customFormat),
            Details:   fmt.Sprintf("将撮合频率从 %d 修改为 %d", MatchFrequency, i),
        }
        MatchFrequency = i
        updateSuccess = true
    case "DepositRate":
        f, err := strconv.ParseFloat(newConfig, 64)
        if err != nil {
            return fmt.Errorf("转换保证金率配置值失败: %w", err)
        }
        actionRecord = &AdminActionRecord{
            Action:    "修改保证金率",
            Timestamp: time.Now().Format(customFormat),
            Details:   fmt.Sprintf("将保证金率从 %.2f 修改为 %.2f", DepositRate, f),
        }
        DepositRate = f
        updateSuccess = true
    case "FeeRate":
        f, err := strconv.ParseFloat(newConfig, 64)
        if err != nil {
            return fmt.Errorf("转换手续费率配置值失败: %w", err)
        }
        actionRecord = &AdminActionRecord{
            Action:    "修改手续费率",
            Timestamp: time.Now().Format(customFormat),
            Details:   fmt.Sprintf("将手续费率从 %.2f 修改为 %.2f", FeeRate, f),
        }
        FeeRate = f
        updateSuccess = true
    default:
		fmt.Printf("Error: 不支持的配置参数名: %s\n", name)
        return fmt.Errorf("不支持的配置参数名: %s", name)
    }

    if updateSuccess {
        // 将配置更新到账本
        if err := ctx.GetStub().PutState(name, []byte(newConfig)); err != nil {
			fmt.Printf("Error: 更新配置到账本失败: %v\n", err)
            return fmt.Errorf("更新配置到账本失败: %w", err)
        }

        // 创建复合键，使用自增 ID
        id, err := s.GetNextID(ctx, AdminActionID)
        if err != nil {
            return err
        }
		uniqueID:=strconv.Itoa(id)
		fmt.Printf("update uniqueID=%s\n",uniqueID)
        compositeKey, err := ctx.GetStub().CreateCompositeKey(AdminActionPrefix, []string{uniqueID})
        if err != nil {
			fmt.Printf("Error: 创建复合键失败: %v\n", err)
            return fmt.Errorf("创建复合键失败: %w", err)
        }

        // 序列化管理员操作记录
        actionRecordBytes, err := json.Marshal(actionRecord)
        if err != nil {
			fmt.Printf("Error: 序列化管理员操作记录失败: %v\n", err)
            return fmt.Errorf("序列化管理员操作记录失败: %w", err)
        }

        // 存储管理员操作记录到账本
        if err := ctx.GetStub().PutState(string(compositeKey), actionRecordBytes); err != nil {
			fmt.Printf("Error: 存储管理员操作记录到账本失败: %v\n", err)
            return fmt.Errorf("存储管理员操作记录到账本失败: %w", err)
        }
		fmt.Printf("将复合键添加到管理员操作历史\n")
        // 将复合键添加到管理员操作历史
        admin.AdminActionHistory = append(admin.AdminActionHistory, compositeKey)

        // 更新管理员信息到账本
        if err := s.updateAdmin(ctx, admin); err != nil {
			fmt.Printf("Error: 更新管理员信息失败: %v\n", err)
            return fmt.Errorf("更新管理员信息失败: %w", err)
        }
    }
	fmt.Println("AdminModify done")
    return nil
}
// // AdminModify 管理员修改配置参数
// func (s *SmartContract) AdminModify(ctx contractapi.TransactionContextInterface, adminId string, name string, newConfig string) error {
// 	admin, err := s.getAdmin(ctx, adminId)
// 	if err != nil {
// 		return err
// 	}

// 	var actionRecord *AdminActionRecord
// 	switch name {
// 	case "MatchFrequency":
// 		i, _ := strconv.Atoi(newConfig)
// 		actionRecord = &AdminActionRecord{
// 			Action:    "修改撮合频率",
// 			Timestamp: time.Now().Format(customFormat),
// 			Details:   fmt.Sprintf("将撮合频率从 %d 修改为 %d", MatchFrequency, i),
// 		}
// 		MatchFrequency = i
// 		ctx.GetStub().PutState(name, []byte(newConfig))
// 	case "DepositRate":
// 		f, _ := strconv.ParseFloat(newConfig, 64)
// 		actionRecord = &AdminActionRecord{
// 			Action:    "修改保证金率",
// 			Timestamp: time.Now().Format(customFormat),
// 			Details:   fmt.Sprintf("将保证金率从 %.2f 修改为 %.2f", DepositRate, f),
// 		}
// 		DepositRate = f
// 		ctx.GetStub().PutState(name, []byte(newConfig))
// 	case "FeeRate":
// 		f, _ := strconv.ParseFloat(newConfig, 64)
// 		actionRecord = &AdminActionRecord{
// 			Action:    "修改手续费率",
// 			Timestamp: time.Now().Format(customFormat),
// 			Details:   fmt.Sprintf("将手续费率从 %.2f 修改为 %.2f", FeeRate, f),
// 		}
// 		FeeRate = f
// 		ctx.GetStub().PutState(name, []byte(newConfig))
// 	}

// 	admin.AdminActionHistory = append(admin.AdminActionHistory, actionRecord)
// 	return s.updateAdmin(ctx, admin)
// }

// GetConfig 获取配置参数
func (s *SmartContract) GetConfig() *Config {
	return &Config{
		MatchFrequency: MatchFrequency,
		DepositRate:    DepositRate,
		FeeRate:        FeeRate,
	}
}

// GetNextID 获取下一个自增 ID
func (s *SmartContract) GetNextID(ctx contractapi.TransactionContextInterface, key string) (int, error) {
    idBytes, err := ctx.GetStub().GetState(key)
    if err != nil {
        return 0, fmt.Errorf("获取自增 ID 失败: %w", err)
    }
    var id int
    if idBytes == nil {
        id = 1
    } else {
        id, err = strconv.Atoi(string(idBytes))
        if err != nil {
            return 0, fmt.Errorf("解析自增 ID 失败: %w", err)
        }
        id++
    }
    newIDBytes := []byte(strconv.Itoa(id))
    if err := ctx.GetStub().PutState(key, newIDBytes); err != nil {
        return 0, fmt.Errorf("更新自增 ID 失败: %w", err)
    }
    return id, nil
}