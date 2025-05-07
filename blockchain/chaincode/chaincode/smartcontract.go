package chaincode

import (
	"encoding/json"
	"fmt"
	"sort"
	"strconv"
	"time"

	"github.com/bwmarrin/snowflake"
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

// SmartContract 定义合约结构体
type SmartContract struct {
	contractapi.Contract
}

// GenerateID 使用 snowflake 库生成唯一的 ID
func GenerateID() string {
	node, _ := snowflake.NewNode(1)
	return node.Generate().String()
}

// RegisterUser 注册用户
func (s *SmartContract) RegisterUser(ctx contractapi.TransactionContextInterface, userId string) error {
	user := &User{
		UserID:            userId,
		Balance:           0,
        IsSeller:          NotApplied,
        IsBuyer:           NotApplied,
		Offers:            []*Offer{},
		Contracts:         []*Contract{},
		BalanceHistory:    []*BalanceRecord{},
		OfferHistory:      []*OfferHistoryRecord{},
		OfferDone:         []*Offer{},
        CreditRating:      0,
        TradeCount:        0,
	}
	userJSON, err := json.Marshal(user)
	if err != nil {
		return err
	}
	return ctx.GetStub().PutState(userId, userJSON)
}

// UserApproveAs 用户购电/售电资质申请
func (s *SmartContract) UserApproveAs(ctx contractapi.TransactionContextInterface, adminId, userId, applyType string) error {
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
        ApplicationID: GenerateID(),
        UserID:        userId,
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
    }
}
// RegisterPlatformAdmin 注册平台管理员
func (s *SmartContract) RegisterPlatformAdmin(ctx contractapi.TransactionContextInterface, adminId string) error {
	admin := &PlatformAdmin{
		AdminID:          adminId,
		Balance:          0,
		BalanceHistory:   []*BalanceRecord{},
		AdminActionHistory: []*AdminActionRecord{},
        Applications: []*Application{},
		//SellList:         []string{},
		//BuyList:          []string{},
		Contracts:        []string{},
	}
	adminJSON, err := json.Marshal(admin)
	if err != nil {
		return err
	}
	ctx.GetStub().PutState("MatchFrequency", []byte("10"))
	ctx.GetStub().PutState("DepositRate", []byte("0.1"))
	ctx.GetStub().PutState("FeeRate", []byte("0.02"))
	return ctx.GetStub().PutState(adminId, adminJSON)
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
	// adminAction := &AdminActionRecord{
	// 	Action:    "用户资质审核",
	// 	Timestamp: time.Now().Format(customFormat),
	// 	Details:   "",
	// }
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
// SubmitOffer 提交报价
func (s *SmartContract) SubmitOffer(ctx contractapi.TransactionContextInterface, userId string, price float64, quantity int, isSeller bool, adminId string) (string, error) {
	user, err := s.getUser(ctx, userId)
	if err != nil {
		return "", err
	}

	var depositAmount float64
	offerId:=GenerateID()
	if isSeller {
		if user.IsSeller != Approved {
			return "", fmt.Errorf("用户 %s 尚未通过售电资质申请，无法提交卖方报价", userId)
		}
	} else {
		if user.IsBuyer != Approved {
			return "", fmt.Errorf("用户 %s 尚未通过购电资质申请，无法提交买方报价", userId)
		}
		depositAmount = float64(quantity) * price * DepositRate
		if user.Balance < depositAmount {
			return "", fmt.Errorf("用户 %s 余额不足，无法提交保证金", userId)
		}
        
		user.Balance -= depositAmount
		balanceRecord := &BalanceRecord{
			Timestamp: time.Now().Format(customFormat),
			Amount:    -depositAmount,
			Rest:      user.Balance,
			Reason:    fmt.Sprintf("买方:%s提交报价:%s保证金", user.UserID,offerId),
		}
		user.BalanceHistory = append(user.BalanceHistory, balanceRecord)

		admin, err := s.getAdmin(ctx, adminId)
		if err != nil {
			return "", err
		}
		admin.Balance += depositAmount
		adminBalanceRecord := &BalanceRecord{
			Timestamp: time.Now().Format(customFormat),
			Amount:    depositAmount,
			Rest:      admin.Balance,
			Reason:    fmt.Sprintf("收到买方:%s提交报价:%s的保证金", user.UserID,offerId),
		}
		admin.BalanceHistory = append(admin.BalanceHistory, adminBalanceRecord)
		if err := s.updateAdmin(ctx, admin); err != nil {
			return "", err
		}
	}

	offer := &Offer{
		OfferID:    offerId,
		UserID:     userId,
		Price:      price,
		Quantity:   quantity,
		Deposit:    depositAmount,
		IsSeller:   isSeller,
		Timestamp:  time.Now().Format(customFormat),
        UpdatedTime:time.Now().Format(customFormat),
		Status:     OfferPending,
		Round:      1,
	}
	user.Offers = append(user.Offers, offer)
	offerRecord := &OfferHistoryRecord{
		Offer:     offer,
		Timestamp: time.Now().Format(customFormat),
		Action:    "提交",
	}
	// 新增：将 Offer 单独存储一份
	offerKey, _ := ctx.GetStub().CreateCompositeKey(OfferPrefix, []string{offerId})
	offerJSON, _ := json.Marshal(offer)
	ctx.GetStub().PutState(offerKey, offerJSON)
	
	user.OfferHistory = append(user.OfferHistory, offerRecord)
	if err := s.updateUser(ctx, user); err != nil {
		return "", err
	}
	return offerId, nil
}

// // createContract 创建合同，增加错误处理
// func (s *SmartContract) createContract(sellerOffer, buyerOffer *Offer) (*Contract, error) {
// 	if sellerOffer == nil || buyerOffer == nil {
// 		return nil, fmt.Errorf("无效的报价信息")
// 	}

// 	quantity := sellerOffer.Quantity
// 	if buyerOffer.Quantity < sellerOffer.Quantity {
// 		quantity = buyerOffer.Quantity
// 	}

// 	if quantity <= 0 {
// 		return nil, fmt.Errorf("无效的交易数量")
// 	}

// 	return &Contract{
// 		ContractID: GenerateID(),
// 		SellerID:   sellerOffer.UserID,
// 		BuyerID:    buyerOffer.UserID,
// 		Price:      sellerOffer.Price,
// 		Quantity:   quantity,
// 		Timestamp:  time.Now().Format(customFormat),
// 		Status:     "待结算",
// 	}, nil
// }

// 优化后的 MatchOffers 函数
func (s *SmartContract) MatchOffers(ctx contractapi.TransactionContextInterface, adminId string) error {
    // 获取所有状态为“待撮合”的报价
    sellerOffers, buyerOffers, err := s.getPendingOffers(ctx)
    if err != nil {
        return fmt.Errorf("获取待撮合报价失败: %v", err)
    }

    // 对卖方报价按价格升序排序，买方报价按价格降序排序
    sort.Slice(sellerOffers, func(i, j int) bool {
        return sellerOffers[i].Price < sellerOffers[j].Price
    })
    sort.Slice(buyerOffers, func(i, j int) bool {
        return buyerOffers[i].Price > buyerOffers[j].Price
    })

    admin, err := s.getAdmin(ctx, adminId)
    if err != nil {
        return fmt.Errorf("获取管理员信息失败: %v", err)
    }

    for len(sellerOffers) > 0 && len(buyerOffers) > 0 {
        sellerOffer := sellerOffers[0]
        buyerOffer := buyerOffers[0]

        // 验证买卖方不能是同一用户
        if sellerOffer.UserID == buyerOffer.UserID {
            // 跳过当前买方报价，尝试下一个买方
            buyerOffers = buyerOffers[1:]
            continue
        }

        // 价格匹配检查
        if buyerOffer.Price >= sellerOffer.Price {
            // 创建合同
            contract, err := s.createContract(ctx,sellerOffer, buyerOffer)
            if err != nil {
                return fmt.Errorf("创建合同失败: %v", err)
            }

            // 更新报价和合同信息
            if err := s.updateOfferAndContract(ctx, sellerOffer, buyerOffer, contract, adminId); err != nil {
                return fmt.Errorf("更新报价和合同信息失败: %v", err)
            }
			// 在 MatchOffers 撮合成功后
			admin.Contracts = append(admin.Contracts, contract.ContractID) // 只存ID
			if err := s.updateAdmin(ctx, admin); err != nil {
				return fmt.Errorf("更新管理员合同列表失败: %v", err)
			}
            // 更新报价列表状态
            sellerOffers, buyerOffers = s.updateOfferLists(ctx, sellerOffers, buyerOffers, sellerOffer, buyerOffer, admin)
        } else {
            // 价格不匹配，终止撮合
            break
        }
    }

    return nil
}
// func (s *SmartContract) getPendingOffers(ctx contractapi.TransactionContextInterface) ([]*Offer, []*Offer, error) {
//     resultsIterator, err := ctx.GetStub().GetStateByPartialCompositeKey("Offer~Status", []string{OfferPending})
//     if err != nil {
//         return nil, nil, err
//     }
//     defer resultsIterator.Close()

//     var sellerOffers []*Offer
//     var buyerOffers []*Offer

//     for resultsIterator.HasNext() {
//         queryResponse, err := resultsIterator.Next()
//         if err != nil {
//             return nil, nil, err
//         }

//         var offer Offer
//         err = json.Unmarshal(queryResponse.Value, &offer)
//         if err != nil {
//             continue
//         }

//         if offer.IsSeller {
//             sellerOffers = append(sellerOffers, &offer)
//         } else {
//             buyerOffers = append(buyerOffers, &offer)
//         }
//     }

//     return sellerOffers, buyerOffers, nil
// }
func (s *SmartContract) getPendingOffers(ctx contractapi.TransactionContextInterface) ([]*Offer, []*Offer, error) {
    resultsIterator, err := ctx.GetStub().GetStateByPartialCompositeKey("Offer", []string{})
    if err != nil {
        return nil, nil, err
    }
    defer resultsIterator.Close()

    var sellerOffers []*Offer
    var buyerOffers []*Offer

    for resultsIterator.HasNext() {
        queryResponse, err := resultsIterator.Next()
        if err != nil {
            return nil, nil, err
        }

        var offer Offer
        err = json.Unmarshal(queryResponse.Value, &offer)
        if err != nil {
            continue
        }
		if offer.Status!=OfferPending {
            continue
        }
        if offer.IsSeller  {
            sellerOffers = append(sellerOffers, &offer)
        } else {
            buyerOffers = append(buyerOffers, &offer)
        }
    }

    return sellerOffers, buyerOffers, nil
}
func (s *SmartContract) createContract(ctx contractapi.TransactionContextInterface,sellerOffer, buyerOffer *Offer) (*Contract, error) {
    if sellerOffer == nil || buyerOffer == nil {
        return nil, fmt.Errorf("无效的报价信息")
    }
    quantity := sellerOffer.Quantity
    if buyerOffer.Quantity < sellerOffer.Quantity {
        quantity = buyerOffer.Quantity
    }
    if quantity <= 0 {
        return nil, fmt.Errorf("无效的交易数量")
    }

    contract:=& Contract{
        ContractID: GenerateID(),
        SellerID:   sellerOffer.UserID,
        BuyerID:    buyerOffer.UserID,
        Price:      sellerOffer.Price,
        Quantity:   quantity,
        Timestamp:  time.Now().Format(customFormat),
        Status:     ContractCreated,//"待结算",
    }
	// 存储到账本
    contractKey, _ := ctx.GetStub().CreateCompositeKey(ContractPrefix, []string{contract.ContractID})
    contractJSON, _ := json.Marshal(contract)
    if err := ctx.GetStub().PutState(contractKey, contractJSON); err != nil {
        return nil, fmt.Errorf("存储合同失败: %v", err)
    }
	return contract, nil
}
// 新增获取合同的函数
func (s *SmartContract) getContract(ctx contractapi.TransactionContextInterface, contractID string) (*Contract, error) {
    contractKey, _ := ctx.GetStub().CreateCompositeKey(ContractPrefix, []string{contractID})
    contractJSON, err := ctx.GetStub().GetState(contractKey)
    if err != nil {
        return nil, fmt.Errorf("获取合同失败: %v", err)
    }
    var contract Contract
    if err := json.Unmarshal(contractJSON, &contract); err != nil {
        return nil, err
    }
    return &contract, nil
}
func (s *SmartContract) updateOfferAndContract(ctx contractapi.TransactionContextInterface, sellerOffer, buyerOffer *Offer, contract *Contract, adminId string) error {
    // 更新卖方信息
    sellerUser, err := s.getUser(ctx, sellerOffer.UserID)
    if err != nil {
        return err
    }
    s.updateUserOffer(sellerUser, sellerOffer, contract.Quantity)
    sellerUser.Contracts = append(sellerUser.Contracts, contract)
    if err := s.updateUser(ctx, sellerUser); err != nil {
        return err
    }

    // 更新买方信息
    buyerUser, err := s.getUser(ctx, buyerOffer.UserID)
    if err != nil {
        return err
    }
    s.updateUserOffer(buyerUser, buyerOffer, contract.Quantity)
    buyerUser.Contracts = append(buyerUser.Contracts, contract)
    if err := s.updateUser(ctx, buyerUser); err != nil {
        return err
    }

    return nil
}
func (s *SmartContract) updateUserOffer(user *User, offer *Offer, matchedQuantity int) {
    for _, o := range user.Offers {
        if o.OfferID == offer.OfferID {
            o.Quantity -= matchedQuantity
            if o.Quantity == 0 {
                o.Status = OfferMatched//"已撮合"
                user.OfferDone = append(user.OfferDone, o)
            }
            break
        }
    }
}
func (s *SmartContract) updateOfferLists(ctx contractapi.TransactionContextInterface, sellerOffers, buyerOffers []*Offer, sellerOffer, buyerOffer *Offer, admin *PlatformAdmin) ([]*Offer, []*Offer) {
    // 更新卖方报价
    if sellerOffer.Quantity <= 0 {
        sellerOffers = sellerOffers[1:]
        s.updateOfferHistory(ctx, sellerOffer, "完成")
    }

    // 更新买方报价
    if buyerOffer.Quantity <= 0 {
        buyerOffers = buyerOffers[1:]
        s.refundBuyerDeposit(ctx, buyerOffer, admin)
    }

    // 移除数量为0的报价
    sellerOffers = s.removeEmptyOffers(sellerOffers)
    buyerOffers = s.removeEmptyOffers(buyerOffers)

    return sellerOffers, buyerOffers
}
// updateOfferHistory 更新报价历史记录
func (s *SmartContract) updateOfferHistory(ctx contractapi.TransactionContextInterface, offer *Offer, action string) error {
	user, err := s.getUser(ctx, offer.UserID)
	if err != nil {
		return err
	}

	historyRecord := &OfferHistoryRecord{
		Offer:     offer,
		Timestamp: time.Now().Format(customFormat),
		Action:    action,
	}
	user.OfferHistory = append(user.OfferHistory, historyRecord)

	return s.updateUser(ctx, user)
}

// refundBuyerDeposit 退还买方保证金
func (s *SmartContract) refundBuyerDeposit(ctx contractapi.TransactionContextInterface, offer *Offer, admin *PlatformAdmin) error {
	user, err := s.getUser(ctx, offer.UserID)
	if err != nil {
		return err
	}

	admin.Balance -= offer.Deposit
	user.Balance += offer.Deposit

	// 记录余额变动
	balanceRecord := &BalanceRecord{
		Timestamp: time.Now().Format(customFormat),
		Amount:    offer.Deposit,
		Rest:      user.Balance,
		Reason:    fmt.Sprintf("退还保证金：%s", offer.OfferID),
	}
	user.BalanceHistory = append(user.BalanceHistory, balanceRecord)

	// 管理员操作记录
	adminAction := &AdminActionRecord{
		Action:    "保证金退还",
		Timestamp: time.Now().Format(customFormat),
		Details:   fmt.Sprintf("退还买方 %s 的报价 %s 保证金 %.2f", offer.UserID, offer.OfferID, offer.Deposit),
	}
	admin.AdminActionHistory = append(admin.AdminActionHistory, adminAction)

	// 更新用户和管理员状态
	if err := s.updateUser(ctx, user); err != nil {
		return err
	}
	if err := s.updateAdmin(ctx, admin); err != nil {
		return err
	}

	return nil
}
func (s *SmartContract) removeEmptyOffers(offers []*Offer) []*Offer {
    validOffers := make([]*Offer, 0, len(offers))
    for _, offer := range offers {
        if offer.Quantity > 0 {
            validOffers = append(validOffers, offer)
        }
    }
    return validOffers
}


// // GetAllUsers 获取所有用户
// func (s *SmartContract) GetAllUsers(ctx contractapi.TransactionContextInterface) ([]*User, error) {
// 	var users []*User
// 	resultsIterator, err := ctx.GetStub().GetStateByRange("", "")
// 	if err != nil {
// 		return nil, err
// 	}
// 	defer resultsIterator.Close()

// 	for resultsIterator.HasNext() {
// 		queryResponse, err := resultsIterator.Next()
// 		if err != nil {
// 			return nil, err
// 		}

// 		var user User
// 		err = json.Unmarshal(queryResponse.Value, &user)
// 		if err != nil {
// 			continue
// 		}
// 		users = append(users, &user)
// 	}
// 	return users, nil
// }

// // UpdateOfferAndContract 更新报价和合同信息
// func (s *SmartContract) UpdateOfferAndContract(ctx contractapi.TransactionContextInterface, sellerOffer *Offer, buyerOffer *Offer, contract *Contract, adminId string) error {
// 	sellerUser, err := s.getUser(ctx, sellerOffer.UserID)
// 	if err != nil {
// 		return err
// 	}
// 	for _, offer := range sellerUser.Offers {
// 		if offer.OfferID == sellerOffer.OfferID {
// 			offer.Quantity -= contract.Quantity
// 			if offer.Quantity == 0 {
// 				offer.Status = "已撮合"
// 				sellerUser.OfferDone = append(sellerUser.OfferDone, offer)
// 			} else {
// 				sellerUser.OfferHistory = append(sellerUser.OfferHistory, &OfferHistoryRecord{
// 					Offer:     offer,
// 					Timestamp: time.Now().Format(customFormat),
// 					Action:    "撮合后更新",
// 				})
// 			}
// 			break
// 		}
// 	}
// 	sellerUser.Contracts = append(sellerUser.Contracts, contract)
// 	if err := s.updateUser(ctx, sellerUser); err != nil {
// 		return err
// 	}

// 	buyerUser, err := s.getUser(ctx, buyerOffer.UserID)
// 	if err != nil {
// 		return err
// 	}
// 	for _, offer := range buyerUser.Offers {
// 		if offer.OfferID == buyerOffer.OfferID {
// 			offer.Quantity -= contract.Quantity
// 			if offer.Quantity == 0 {
// 				offer.Status = "已撮合"
// 				buyerUser.OfferDone = append(buyerUser.OfferDone, offer)
// 			} else {
// 				buyerUser.OfferHistory = append(buyerUser.OfferHistory, &OfferHistoryRecord{
// 					Offer:     offer,
// 					Timestamp: time.Now().Format(customFormat),
// 					Action:    "撮合后更新",
// 				})
// 			}
// 			break
// 		}
// 	}
// 	buyerUser.Contracts = append(buyerUser.Contracts, contract)
// 	return s.updateUser(ctx, buyerUser)
// }

// SettleContract 结算合同
func (s *SmartContract) SettleContract(ctx contractapi.TransactionContextInterface, adminId string) error {
	admin, err := s.getAdmin(ctx, adminId)
	if err != nil {
		return err
	}
	for _, contractID := range admin.Contracts {
		contract, err := s.getContract(ctx, contractID)
		if err != nil || contract.Status != ContractCreated {
			continue // 跳过已处理或无效合同
		}
		sellerUser, err := s.getUser(ctx, contract.SellerID)
		if err != nil {
			return err
		}

		buyerUser, err := s.getUser(ctx, contract.BuyerID)
		if err != nil {
			return err
		}

		totalAmount := float64(contract.Quantity) * contract.Price
		fee := totalAmount * FeeRate

		sellerUser.Balance += totalAmount - fee
		sellerUser.BalanceHistory = append(sellerUser.BalanceHistory, &BalanceRecord{
			Timestamp: time.Now().Format(customFormat),
			Amount:    totalAmount - fee,
			Rest:      sellerUser.Balance,
			Reason:    "合同结算收入（扣除手续费）",
		})

		buyerUser.Balance -= totalAmount
		buyerUser.BalanceHistory = append(buyerUser.BalanceHistory, &BalanceRecord{
			Timestamp: time.Now().Format(customFormat),
			Amount:    -totalAmount,
			Rest:      buyerUser.Balance,
			Reason:    "合同结算支出",
		})

		admin.Balance += fee
		admin.BalanceHistory = append(admin.BalanceHistory, &BalanceRecord{
			Timestamp: time.Now().Format(customFormat),
			Amount:    fee,
			Rest:      admin.Balance,
			Reason:    "收取结算手续费",
		})

		contract.Status = ContractSettled
		// 更新合同到账本
		contractKey, _ := ctx.GetStub().CreateCompositeKey(ContractPrefix, []string{contractID})
		contractJSON, _ := json.Marshal(contract)
		
		if err:=ctx.GetStub().PutState(contractKey, contractJSON);err!=nil{
			return err
		}	
		if err := s.updateAdmin(ctx, admin); err != nil {
			return err
		}
		if err := s.updateUser(ctx, sellerUser); err != nil {
			return err
		}
		if err := s.updateUser(ctx, buyerUser); err != nil {
			return err
		}
	}
	return nil
}

// ModifyOffer （用户）修改报价信息
func (s *SmartContract) ModifyOffer(ctx contractapi.TransactionContextInterface, userId string, offerId string, newPrice float64, newQuantity int, adminId string) error {
	user, err := s.getUser(ctx, userId)
	if err != nil {
		return err
	}

	for i, o := range user.Offers {
		if o.OfferID == offerId {
			if !o.IsSeller {
				oldDeposit := float64(o.Quantity) * o.Price * DepositRate
				newDeposit := float64(newQuantity) * newPrice * DepositRate

				if newDeposit > oldDeposit {
					additionalDeposit := newDeposit - oldDeposit
					if user.Balance < additionalDeposit {
						return fmt.Errorf("用户 %s 余额不足，无法增加保证金", user.UserID)
					}
					user.Balance -= additionalDeposit
					o.Deposit += additionalDeposit
					user.BalanceHistory = append(user.BalanceHistory, &BalanceRecord{
						Timestamp: time.Now().Format(customFormat),
						Amount:    -additionalDeposit,
						Rest:      user.Balance,
						Reason:    "修改报价增加保证金",
					})

					admin, err := s.getAdmin(ctx, adminId)
					if err != nil {
						return err
					}
					admin.Balance += additionalDeposit
					admin.BalanceHistory = append(admin.BalanceHistory, &BalanceRecord{
						Timestamp: time.Now().Format(customFormat),
						Amount:    additionalDeposit,
						Rest:      admin.Balance,
						Reason:    "用户修改报价增加的保证金",
					})
					if err := s.updateAdmin(ctx, admin); err != nil {
						return err
					}
				} else if newDeposit < oldDeposit {
					refundDeposit := oldDeposit - newDeposit
					user.Balance += refundDeposit
					user.BalanceHistory = append(user.BalanceHistory, &BalanceRecord{
						Timestamp: time.Now().Format(customFormat),
						Amount:    refundDeposit,
						Rest:      user.Balance,
						Reason:    "修改报价退还保证金",
					})

					admin, err := s.getAdmin(ctx, adminId)
					if err != nil {
						return err
					}
					admin.Balance -= refundDeposit
					admin.BalanceHistory = append(admin.BalanceHistory, &BalanceRecord{
						Timestamp: time.Now().Format(customFormat),
						Amount:    -refundDeposit,
						Rest:      admin.Balance,
						Reason:    "用户修改报价退还的保证金",
					})
					if err := s.updateAdmin(ctx, admin); err != nil {
						return err
					}
				}
			}

			oldOffer := o
			o.Price = newPrice
			o.Quantity = newQuantity
			o.Status = OfferPending//"待撮合"
			o.Round = 0
			user.Offers[i] = o
			user.OfferHistory = append(user.OfferHistory, &OfferHistoryRecord{
				Offer:     oldOffer,
				Timestamp: time.Now().Format(customFormat),
				Action:    "修改",
			})

			if err := s.updateUser(ctx, user); err != nil {
				return err
			}
			return nil
		}
	}

	return s.updateUser(ctx, user)
}

// // GetAllOffers 获取所有报价（分页查询）
// func (s *SmartContract) GetAllOffers(ctx contractapi.TransactionContextInterface, pageSize int, bookmark string) ([]*Offer, string, error) {
//     resultsIterator, responseMetadata, err := ctx.GetStub().GetStateByPartialCompositeKeyWithPagination(OfferPrefix, []string{}, int32(pageSize), bookmark)
//     if err != nil {
//         return nil, "", err
//     }
//     defer resultsIterator.Close()

//     var allOffers []*Offer

//     for resultsIterator.HasNext() {
//         queryResponse, err := resultsIterator.Next()
//         if err != nil {
//             return nil, "", err
//         }

//         var offer Offer
//         if err := json.Unmarshal(queryResponse.Value, &offer); err != nil {
//             continue
//         }

//         allOffers = append(allOffers, &offer)
//     }

//     return allOffers, responseMetadata.Bookmark, nil
// }

// // GetAllUsers 获取所有用户（分页查询）
// func (s *SmartContract) GetAllUsers(ctx contractapi.TransactionContextInterface, pageSize int, bookmark string) ([]*User, string, error) {
//     resultsIterator, responseMetadata, err := ctx.GetStub().GetStateByRangeWithPagination("", "", int32(pageSize), bookmark)
//     if err != nil {
//         return nil, "", err
//     }
//     defer resultsIterator.Close()

//     var users []*User

//     for resultsIterator.HasNext() {
//         queryResponse, err := resultsIterator.Next()
//         if err != nil {
//             return nil, "", err
//         }

//         var user User
//         err = json.Unmarshal(queryResponse.Value, &user)
//         if err != nil {
//             continue
//         }
//         users = append(users, &user)
//     }

//     return users, responseMetadata.Bookmark, nil
// }

// GetUserInfo 获取用户信息
func (s *SmartContract) GetUserInfo(ctx contractapi.TransactionContextInterface, userId string) (*User, error) {
	return s.getUser(ctx, userId)
}

// GetAdminInfo 获取管理员信息
func (s *SmartContract) GetAdminInfo(ctx contractapi.TransactionContextInterface, adminId string) (*PlatformAdmin, error) {
	return s.getAdmin(ctx, adminId)
}

// GetAllOffer 用户查询所有报价信息
func (s *SmartContract) GetAllOffer(ctx contractapi.TransactionContextInterface, userId string) ([]*Offer, error) {
	user, err := s.getUser(ctx, userId)
	if err != nil {
		return nil, err
	}
	return append(user.Offers, user.OfferDone...), nil
}

// GetOfferHistory 用户查询所有报价历史信息
func (s *SmartContract) GetOfferHistory(ctx contractapi.TransactionContextInterface, userId string) ([]*OfferHistoryRecord, error) {
	user, err := s.getUser(ctx, userId)
	if err != nil {
		return nil, err
	}
	return user.OfferHistory, nil
}

// GetBalanceHistory 查询账户余额历史
func (s *SmartContract) GetBalanceHistory(ctx contractapi.TransactionContextInterface, userId string) ([]*BalanceRecord, error) {
	user, err := s.getUser(ctx, userId)
	if err != nil {
		return nil, err
	}
	return user.BalanceHistory, nil
}

// GetUserContracts 查询用户参与的购电合同
func (s *SmartContract) GetUserContracts(ctx contractapi.TransactionContextInterface, userId string) ([]*Contract, error) {
	user, err := s.getUser(ctx, userId)
	if err != nil {
		return nil, err
	}
	return user.Contracts, nil
}

// AdminModify 管理员修改配置参数
func (s *SmartContract) AdminModify(ctx contractapi.TransactionContextInterface, adminId string, name string, newConfig string) error {
	admin, err := s.getAdmin(ctx, adminId)
	if err != nil {
		return err
	}

	var actionRecord *AdminActionRecord
	switch name {
	case "MatchFrequency":
		i, _ := strconv.Atoi(newConfig)
		actionRecord = &AdminActionRecord{
			Action:    "修改撮合频率",
			Timestamp: time.Now().Format(customFormat),
			Details:   fmt.Sprintf("将撮合频率从 %d 修改为 %d", MatchFrequency, i),
		}
		MatchFrequency = i
		ctx.GetStub().PutState(name, []byte(newConfig))
	case "DepositRate":
		f, _ := strconv.ParseFloat(newConfig, 64)
		actionRecord = &AdminActionRecord{
			Action:    "修改保证金率",
			Timestamp: time.Now().Format(customFormat),
			Details:   fmt.Sprintf("将保证金率从 %.2f 修改为 %.2f", DepositRate, f),
		}
		DepositRate = f
		ctx.GetStub().PutState(name, []byte(newConfig))
	case "FeeRate":
		f, _ := strconv.ParseFloat(newConfig, 64)
		actionRecord = &AdminActionRecord{
			Action:    "修改手续费率",
			Timestamp: time.Now().Format(customFormat),
			Details:   fmt.Sprintf("将手续费率从 %.2f 修改为 %.2f", FeeRate, f),
		}
		FeeRate = f
		ctx.GetStub().PutState(name, []byte(newConfig))
	}

	admin.AdminActionHistory = append(admin.AdminActionHistory, actionRecord)
	return s.updateAdmin(ctx, admin)
}

// GetConfig 获取配置参数
func (s *SmartContract) GetConfig() *Config {
	return &Config{
		MatchFrequency: MatchFrequency,
		DepositRate:    DepositRate,
		FeeRate:        FeeRate,
	}
}

// GetAdminActionHistory 获取管理员操作历史
func (s *SmartContract) GetAdminActionHistory(ctx contractapi.TransactionContextInterface, adminId string) ([]*AdminActionRecord, error) {
	admin, err := s.getAdmin(ctx, adminId)
	if err != nil {
		return nil, err
	}
	return admin.AdminActionHistory, nil
}

// GetAdminMoneyHistory 获取管理员账户余额历史
func (s *SmartContract) GetAdminMoneyHistory(ctx contractapi.TransactionContextInterface, adminId string) ([]*BalanceRecord, error) {
	admin, err := s.getAdmin(ctx, adminId)
	if err != nil {
		return nil, err
	}
	return admin.BalanceHistory, nil
}

// AdminGetAllOffer 查询系统内所有报价
func (s *SmartContract) AdminGetAllOffer(ctx contractapi.TransactionContextInterface) ([]*Offer, error) {
	resultsIterator, err := ctx.GetStub().GetStateByRange("", "")
	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()

	var offerlist []*Offer
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return nil, err
		}

		var user User
		err = json.Unmarshal(queryResponse.Value, &user)
		if err != nil {
			continue
		}
		list, err := s.GetAllOffer(ctx, user.UserID)
		offerlist = append(offerlist, list...)
	}
	return offerlist, nil
}

// GetAllContract 查询系统内所有合同
func (s *SmartContract) GetAllContract(ctx contractapi.TransactionContextInterface) ([]*Contract, error) {
	resultsIterator, err := ctx.GetStub().GetStateByRange("", "")
	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()

	var contractlist []*Contract
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return nil, err
		}

		var user User
		err = json.Unmarshal(queryResponse.Value, &user)
		if err != nil {
			continue
		}
		list, err := s.GetUserContracts(ctx, user.UserID)
		contractlist = append(contractlist, list...)
	}
	return contractlist, nil
}

// // GetAllContract 查询系统内所有合同（分页查询）
// func (s *SmartContract) GetAllContract(ctx contractapi.TransactionContextInterface, pageSize int, bookmark string) ([]*Contract, string, error) {
//     resultsIterator, responseMetadata, err := ctx.GetStub().GetStateByRangeWithPagination("", "", int32(pageSize), bookmark)
//     if err != nil {
//         return nil, "", err
//     }
//     defer resultsIterator.Close()

//     var contractlist []*Contract

//     for resultsIterator.HasNext() {
//         queryResponse, err := resultsIterator.Next()
//         if err != nil {
//             return nil, "", err
//         }

//         var user User
//         err = json.Unmarshal(queryResponse.Value, &user)
//         if err != nil {
//             continue
//         }
//         list, err := s.GetUserContracts(ctx, user.UserID)
//         contractlist = append(contractlist, list...)
//     }

//     return contractlist, responseMetadata.Bookmark, nil
// }

