package chaincode

import (
	"encoding/json"
	"fmt"
	"sort"
	"time"
	//"strings"
	"strconv"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
	//"github.com/google/uuid"
)

// SubmitOffer 提交报价
func (s *SmartContract) SubmitOffer(ctx contractapi.TransactionContextInterface, userId string, offerId string, price float64, quantity int, isSeller bool, adminId string) (*Offer, error) {
	fmt.Println("into submitoffer")
	// 获取用户信息
	user, err := s.getUser(ctx, userId)
	if err != nil {
		fmt.Println("获取用户信息失败:", err)
		return nil, fmt.Errorf("获取用户信息失败: %w", err)
	}

	var depositAmount float64
	// 检查用户资质
	if isSeller {
		if user.IsSeller != Approved {
			return nil, fmt.Errorf("用户 %s 尚未通过售电资质申请，无法提交卖方报价", userId)
		}
	} else {
		if user.IsBuyer != Approved {
			return nil, fmt.Errorf("用户 %s 尚未通过购电资质申请，无法提交买方报价", userId)
		}
		fmt.Println("用户 %s 通过购电资质申请，计算保证金", userId)
		// 计算保证金
		depositAmount = float64(quantity) * price * DepositRate
		if user.Balance < depositAmount {
			fmt.Println("用户 %s 余额不足，无法提交保证金", userId)
			return nil, fmt.Errorf("用户 %s 余额不足，无法提交保证金", userId)
		}

		// 扣除用户保证金
		user.Balance -= depositAmount
		// 记录用户余额变动
		if err := s.recordBalanceChange(ctx, userId, -depositAmount, user.Balance, fmt.Sprintf("报价:%s保证金", offerId), false, AdminName, false); err != nil {
			fmt.Println("记录用户余额变动失败:", err)
			return nil, err
		}

		// 获取管理员信息
		admin, err := s.getAdmin(ctx, adminId)
		if err != nil {
			fmt.Println("获取管理员信息失败:", err)
			return nil, fmt.Errorf("获取管理员信息失败: %w", err)
		}
		// 增加管理员余额
		admin.Balance += depositAmount
		// 记录管理员余额变动
		if err := s.recordBalanceChange(ctx, adminId, depositAmount, admin.Balance, fmt.Sprintf("报价:%s的保证金", offerId), true, user.UserName, true); err != nil {
			fmt.Println("记录管理员余额变动失败:", err)
			return nil, err
		}
		// 更新管理员信息
		if err := s.updateAdmin(ctx, admin); err != nil {
			fmt.Println("更新管理员信息失败:", err)
			return nil, err
		}
	}

	// 创建报价对象
	offer := &Offer{
		OfferID:    offerId,
		UserID:     userId,
		UserName:   user.UserName,
		Price:      price,
		Quantity:   quantity,
		Deposit:    depositAmount,
		IsSeller:   isSeller,
		Timestamp:  time.Now().Format(customFormat),
		UpdatedTime: time.Now().Format(customFormat),
		Status:     OfferPending,
		Round:      1,
	}
	fmt.Println("创建offer")
	// 存储报价信息
	if err := s.saveOffer(ctx, offer); err != nil {
		fmt.Println("存储报价信息失败:", err)
		return nil, err
	}
	// 更新用户报价列表
	user.Offers = append(user.Offers, offer.OfferID)
	// 记录报价历史
	if err := s.recordOfferHistory(ctx, offer, "提交"); err != nil {
		fmt.Println("记录报价历史失败:", err)
		return nil, err
	}
	// 更新用户信息
	if err := s.updateUser(ctx, user); err != nil {
		return nil, err
	}
	fmt.Println("submit offer success")
	return offer, nil
}

// 优化后的 MatchOffers 函数
func (s *SmartContract) MatchOffers(ctx contractapi.TransactionContextInterface, adminId string) error {
	fmt.Println("into matchoffer")
	// 获取所有状态为“待撮合”的报价
	sellerOffers, buyerOffers, err := s.getPendingOffers(ctx)
	if err != nil {
		fmt.Println("获取待撮合报价失败:", err)
		return fmt.Errorf("获取待撮合报价失败: %w", err)
	}

	fmt.Println("sort offer")
	// 对卖方报价按价格升序排序，买方报价按价格降序排序
	sort.Slice(sellerOffers, func(i, j int) bool {
		return sellerOffers[i].Price < sellerOffers[j].Price
	})
	sort.Slice(buyerOffers, func(i, j int) bool {
		return buyerOffers[i].Price > buyerOffers[j].Price
	})

	// 打印获取到的待撮合报价
	fmt.Println("获取到的卖方待撮合报价：")
	for _, offer := range sellerOffers {
		fmt.Printf("报价ID: %s, 价格: %.2f, 数量: %d\n", offer.OfferID, offer.Price, offer.Quantity)
	}
	fmt.Println("获取到的买方待撮合报价：")
	for _, offer := range buyerOffers {
		fmt.Printf("报价ID: %s, 价格: %.2f, 数量: %d\n", offer.OfferID, offer.Price, offer.Quantity)
	}

	// 获取管理员信息
	admin, err := s.getAdmin(ctx, adminId)
	if err != nil {
		fmt.Println("获取管理员信息失败:", err)
		return fmt.Errorf("获取管理员信息失败: %w", err)
	}

	fmt.Println("start match")
	for len(sellerOffers) > 0 && len(buyerOffers) > 0 {
		sellerOffer := sellerOffers[0]
		buyerOffer := buyerOffers[0]

		// 验证买卖方不能是同一用户
		if sellerOffer.UserID == buyerOffer.UserID {
			// 跳过当前买方报价，尝试下一个买方
			buyerOffers = buyerOffers[1:]
			continue
		}

		// 检查报价是否有效
		if !isValidOffer(sellerOffer) || !isValidOffer(buyerOffer) {
			if !isValidOffer(sellerOffer) {
				fmt.Println("sell quantity<=0111")
				sellerOffers = sellerOffers[1:]
			}
			if !isValidOffer(buyerOffer) {
				fmt.Println("buy quantity<=0111")
				buyerOffers = buyerOffers[1:]
			}
			continue
		}

		// 价格匹配检查
		if buyerOffer.Price >= sellerOffer.Price {
			// 创建合同
			contract, err := s.createContract(ctx, sellerOffer, buyerOffer, strconv.Itoa(admin.Contractnumber))
			if err != nil {
				fmt.Println("创建合同失败:", err)
				return fmt.Errorf("创建合同失败: %w", err)
			}
			admin.Contractnumber++
			fmt.Println("id:", admin.Contractnumber)

			// 更新报价
			if err := s.updateOffer(ctx, sellerOffer, buyerOffer, contract, admin); err != nil {
				fmt.Println("更新报价和合同信息失败:", err)
				return fmt.Errorf("更新报价和合同信息失败: %w", err)
			}

			// 在 MatchOffers 撮合成功后
			admin.Contracts = append(admin.Contracts, contract.ContractID) // 只存ID
			if err := s.updateAdmin(ctx, admin); err != nil {
				fmt.Println("更新管理员合同列表失败:", err)
				return fmt.Errorf("更新管理员合同列表失败: %w", err)
			}
		} else {
			// 价格不匹配，终止撮合
			fmt.Println("价格不匹配，终止撮合")
			break
		}
	}

	fmt.Println("well done matchoffer")
	return nil
}

// getPendingOffers 获取所有待撮合的报价
func (s *SmartContract) getPendingOffers(ctx contractapi.TransactionContextInterface) ([]*Offer, []*Offer, error) {
	fmt.Println("into getPendingOffers")
	// 获取所有报价
	resultsIterator, err := ctx.GetStub().GetStateByPartialCompositeKey(OfferPrefix, []string{})
	if err != nil {
		return nil, nil, err
	}
	defer resultsIterator.Close()

	var sellerOffers []*Offer
	var buyerOffers []*Offer

	// 遍历所有报价，筛选出待撮合的报价
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
		if offer.Status != OfferPending {
			continue
		}
		if offer.IsSeller {
			sellerOffers = append(sellerOffers, &offer)
		} else {
			buyerOffers = append(buyerOffers, &offer)
		}
	}

	return sellerOffers, buyerOffers, nil
}

// createContract 创建合同
func (s *SmartContract) createContract(ctx contractapi.TransactionContextInterface, sellerOffer, buyerOffer *Offer, contractId string) (*Contract, error) {
	if sellerOffer == nil || buyerOffer == nil {
		return nil, fmt.Errorf("无效的报价信息")
	}
	// 计算交易数量
	quantity := sellerOffer.Quantity
	if buyerOffer.Quantity < sellerOffer.Quantity {
		quantity = buyerOffer.Quantity
	}
	if quantity <= 0 {
		return nil, fmt.Errorf("无效的交易数量")
	}

	// 获取用户信息以填充合同字段
	sellerUser, err := s.getUser(ctx, sellerOffer.UserID)
	if err != nil {
		return nil, err
	}

	buyerUser, err := s.getUser(ctx, buyerOffer.UserID)
	if err != nil {
		return nil, err
	}

	// 创建合同对象
	contract := &Contract{
		ContractID: contractId,
		SellerID:   sellerOffer.UserID,
		BuyerID:    buyerOffer.UserID,
		SellerName: sellerUser.UserName,
		BuyerName:  buyerUser.UserName,
		sellerOfferID:   sellerOffer.OfferID,
		BuyerOfferID:    buyerOffer.OfferID,
		Price:      sellerOffer.Price,
		Quantity:   quantity,
		Timestamp:  time.Now().Format(customFormat),
		SettledTime:     "", // 初始为空，结算时更新
		TradeAmount:     float64(quantity) * sellerOffer.Price,
		Status:     ContractCreated,
	}

	// 存储合同信息到账本
	contractKey, _ := ctx.GetStub().CreateCompositeKey(ContractPrefix, []string{contract.ContractID})
	contractJSON, _ := json.Marshal(contract)
	if err := ctx.GetStub().PutState(contractKey, contractJSON); err != nil {
		return nil, fmt.Errorf("存储合同失败: %w", err)
	}

	return contract, nil
}

// getContract 获取合同信息
func (s *SmartContract) getContract(ctx contractapi.TransactionContextInterface, contractID string) (*Contract, error) {
	contractKey, _ := ctx.GetStub().CreateCompositeKey(ContractPrefix, []string{contractID})
	contractJSON, err := ctx.GetStub().GetState(contractKey)
	if err != nil {
		return nil, fmt.Errorf("获取合同失败: %w", err)
	}
	var contract Contract
	if err := json.Unmarshal(contractJSON, &contract); err != nil {
		return nil, err
	}
	return &contract, nil
}

// updateOffer 更新报价信息
func (s *SmartContract) updateOffer(ctx contractapi.TransactionContextInterface, sellerOffer, buyerOffer *Offer, contract *Contract, admin *PlatformAdmin) error {
	// 更新卖方信息
	sellerUser, err := s.getUser(ctx, sellerOffer.UserID)
	if err != nil {
		return err
	}
	s.updateUserOffer(ctx, sellerUser, sellerOffer, contract.Quantity, admin)
	sellerUser.Contracts = append(sellerUser.Contracts, contract.ContractID)
	if err := s.updateUser(ctx, sellerUser); err != nil {
		return err
	}
	if err := s.saveOffer(ctx, sellerOffer); err != nil {
		return err
	}

	// 更新买方信息
	buyerUser, err := s.getUser(ctx, buyerOffer.UserID)
	if err != nil {
		return err
	}
	s.updateUserOffer(ctx, buyerUser, buyerOffer, contract.Quantity, admin)
	buyerUser.Contracts = append(buyerUser.Contracts, contract.ContractID)
	if err := s.updateUser(ctx, buyerUser); err != nil {
		return err
	}
	if err := s.saveOffer(ctx, buyerOffer); err != nil {
		return err
	}

	return nil
}

// updateUserOffer 更新用户报价信息
func (s *SmartContract) updateUserOffer(ctx contractapi.TransactionContextInterface, user *User, offer *Offer, matchedQuantity int, admin *PlatformAdmin) {
	offer.Quantity -= matchedQuantity
	if offer.Quantity == 0 {
		offer.Status = OfferMatched
		// 将已完成的OfferID添加到OfferDone列表
		user.OfferDone = append(user.OfferDone, offer.OfferID)
		if offer.IsSeller {
			s.updateOfferHistory(ctx, offer, "完成")
		} else {
			s.refundBuyerDeposit(ctx, offer, admin, offer.Deposit)
			s.updateOfferHistory(ctx, offer, "完成")
		}
	} else {
		if offer.IsSeller {
			s.updateOfferHistory(ctx, offer, "更新")
		} else {
			// 计算未撮合部分的保证金（按比例退还）
			oldDeposit := offer.Deposit
			originalQuantity := offer.Quantity + matchedQuantity
			remainingRatio := float64(offer.Quantity) / float64(originalQuantity)
			newDeposit := oldDeposit * remainingRatio
			refundDeposit := oldDeposit - newDeposit
			offer.Deposit = newDeposit
			s.refundBuyerDeposit(ctx, offer, admin, refundDeposit)
			s.updateOfferHistory(ctx, offer, "更新")
		}
	}
}

// updateOfferHistory 更新报价历史记录
func (s *SmartContract) updateOfferHistory(ctx contractapi.TransactionContextInterface, offer *Offer, action string) error {
	user, err := s.getUser(ctx, offer.UserID)
	if err != nil {
		return err
	}
	offer.UpdatedTime = time.Now().Format(customFormat)
	historyRecord := &OfferHistoryRecord{
		OfferID:     offer.OfferID,
		Action:    action,
		OfferSnapshot: *offer,
	}
	// 创建历史记录的复合键
	id, err := s.GetNextID(ctx, OfferHistoryID)
    if err != nil {
		fmt.Println("获取下一个ID失败: %w", err)
        return err
    }
	uniqueID := strconv.Itoa(id)
	historyRecordKey, _ := ctx.GetStub().CreateCompositeKey(OfferHistoryPrefix, []string{offer.OfferID, uniqueID})

	// 序列化并存储历史记录
	historyRecordJSON, _ := json.Marshal(historyRecord)
	if err := ctx.GetStub().PutState(historyRecordKey, historyRecordJSON); err != nil {
		return err
	}

	// 将历史记录键添加到用户的OfferHistory
	user.OfferHistory = append(user.OfferHistory, historyRecordKey)

	return s.updateUser(ctx, user)
}

// refundBuyerDeposit 退还买方保证金
func (s *SmartContract) refundBuyerDeposit(ctx contractapi.TransactionContextInterface, offer *Offer, admin *PlatformAdmin, refundAmount float64) error {
	user, err := s.getUser(ctx, offer.UserID)
	if err != nil {
		return err
	}

	// 生成唯一ID替代txID，避免复合键重复
	//uniqueID := uuid.New().String()

	admin.Balance -= refundAmount
	user.Balance += refundAmount

	// 记录余额变动
	if err := s.recordBalanceChange(ctx, user.UserID, refundAmount, user.Balance, fmt.Sprintf("退还保证金：%s", offer.OfferID), true, AdminName, false); err != nil {
		return err
	}

	// 记录管理员余额变动
	if err := s.recordBalanceChange(ctx, admin.AdminID, -refundAmount, admin.Balance, fmt.Sprintf("退还报价 %s 保证金", offer.OfferID), false, user.UserName, true); err != nil {
		return err
	}

	// 更新用户和管理员状态
	if err := s.updateUser(ctx, user); err != nil {
		return err
	}
	if err := s.updateAdmin(ctx, admin); err != nil {
		return err
	}

	return nil
}

// SettleContract 结算合同
func (s *SmartContract) SettleContract(ctx contractapi.TransactionContextInterface, adminId string) error {
	fmt.Println("into settle")
	// 获取管理员信息
	admin, err := s.getAdmin(ctx, adminId)
	if err != nil {
		return err
	}
	fmt.Println("admin contracts:", len(admin.Contracts))

	// 遍历管理员的所有合同
	for _, contractID := range admin.Contracts {
		contract, err := s.getContract(ctx, contractID)
		if err != nil {
			fmt.Println("获取合同失败:", err)
			continue
		}

		// 只处理已创建但未结算的合同
		if contract.Status != ContractCreated {
			continue
		}

		// 获取买卖方用户
		sellerUser, err := s.getUser(ctx, contract.SellerID)
		if err != nil {
			return err
		}

		buyerUser, err := s.getUser(ctx, contract.BuyerID)
		if err != nil {
			return err
		}

		// 计算交易金额和手续费
		totalAmount := float64(contract.Quantity) * contract.Price
		fee := totalAmount * FeeRate

		// 买方扣款
		buyerUser.Balance -= totalAmount
		if err := s.recordBalanceChange(ctx, buyerUser.UserID, -totalAmount, buyerUser.Balance, fmt.Sprintf("合同结算支出（合同ID：%s）", contract.ContractID), false, sellerUser.UserName, false); err != nil {
			return err
		}

		// 卖方收款（扣除手续费）
		sellerUser.Balance += totalAmount - fee
		if err := s.recordBalanceChange(ctx, sellerUser.UserID, totalAmount - fee, sellerUser.Balance, fmt.Sprintf("合同结算收入（合同ID：%s，扣除手续费）", contract.ContractID), true, buyerUser.UserName, false); err != nil {
			return err
		}

		// 平台收取手续费
		admin.Balance += fee
		if err := s.recordBalanceChange(ctx, admin.AdminID, fee, admin.Balance, fmt.Sprintf("收取合同（ID：%s）结算手续费", contract.ContractID), true, sellerUser.UserName, true); err != nil {
			return err
		}

		// 更新合同状态和结算时间
		contract.Status = ContractSettled
		contract.SettledTime = time.Now().Format(customFormat)

		// 更新合同到账本
		contractKey, err := ctx.GetStub().CreateCompositeKey(ContractPrefix, []string{contract.ContractID})
		if err != nil {
			return fmt.Errorf("创建合同复合键失败: %w", err)
		}
		contractJSON, err := json.Marshal(contract)
		if err != nil {
			return err
		}
		if err := ctx.GetStub().PutState(contractKey, contractJSON); err != nil {
			return err
		}

		// 更新用户和管理员状态
		if err := s.updateUser(ctx, sellerUser); err != nil {
			return err
		}
		if err := s.updateUser(ctx, buyerUser); err != nil {
			return err
		}
		if err := s.updateAdmin(ctx, admin); err != nil {
			return err
		}

		fmt.Printf("合同 %s 结算完成\n", contract.ContractID)
	}

	fmt.Println("well done settle")
	return nil
}

// ModifyOffer 修改报价信息
func (s *SmartContract) ModifyOffer(ctx contractapi.TransactionContextInterface, userId string, offerId string, newPrice float64, newQuantity int, adminId string) error {
	// 获取完整的Offer信息
	offer, err := s.getOffer(ctx, offerId)
	if err != nil {
		return fmt.Errorf("获取报价信息失败: %w", err)
	}

	// 验证用户身份
	if offer.UserID != userId {
		return fmt.Errorf("用户 %s 无权修改报价 %s", userId, offerId)
	}

	// 保存旧状态用于历史记录
	//oldOffer := *offer

	// 买方需要处理保证金变化
	if !offer.IsSeller {
		// 直接使用当前存储的保证金值
		oldDeposit := offer.Deposit

		// 计算新保证金（基于新的数量和价格，使用当前保证金率）
		newDeposit := float64(newQuantity) * newPrice * DepositRate

		if newDeposit > oldDeposit {
			// 需要追加保证金
			additionalDeposit := newDeposit - oldDeposit

			// 检查用户余额是否足够
			user, err := s.getUser(ctx, userId)
			if err != nil {
				return err
			}

			if user.Balance < additionalDeposit {
				return fmt.Errorf("用户 %s 余额不足，无法追加保证金", userId)
			}

			// 扣除用户余额并增加保证金
			user.Balance -= additionalDeposit
			offer.Deposit = newDeposit

			// 记录保证金增加
			if err := s.recordBalanceChange(ctx, userId, -additionalDeposit, user.Balance, fmt.Sprintf("修改报价（ID:%s）追加保证金", offerId), false, AdminName, false); err != nil {
				return err
			}

			// 更新管理员账户
			admin, err := s.getAdmin(ctx, adminId)
			if err != nil {
				return err
			}

			admin.Balance += additionalDeposit

			// 记录管理员余额变化
			if err := s.recordBalanceChange(ctx, adminId, additionalDeposit, admin.Balance, fmt.Sprintf("修改报价（ID:%s）追加保证金", offerId), true, user.UserName, true); err != nil {
				return err
			}

			// 更新用户和管理员状态
			if err := s.updateUser(ctx, user); err != nil {
				return err
			}

			if err := s.updateAdmin(ctx, admin); err != nil {
				return err
			}
		} else if newDeposit < oldDeposit {
			// 需要退还部分保证金
			refundAmount := oldDeposit - newDeposit

			// 获取用户信息
			user, err := s.getUser(ctx, userId)
			if err != nil {
				return err
			}

			// 增加用户余额并减少保证金
			user.Balance += refundAmount
			offer.Deposit = newDeposit

			// 记录保证金退还
			if err := s.recordBalanceChange(ctx, userId, refundAmount, user.Balance, fmt.Sprintf("修改报价（ID:%s）退还保证金", offerId), true, AdminName, false); err != nil {
				return err
			}

			// 更新管理员账户
			admin, err := s.getAdmin(ctx, adminId)
			if err != nil {
				return err
			}

			admin.Balance -= refundAmount

			// 记录管理员余额变化
			if err := s.recordBalanceChange(ctx, adminId, -refundAmount, admin.Balance, fmt.Sprintf("修改报价（ID:%s）退还保证金", offerId), false, user.UserName, true); err != nil {
				return err
			}

			// 更新用户和管理员状态
			if err := s.updateUser(ctx, user); err != nil {
				return err
			}

			if err := s.updateAdmin(ctx, admin); err != nil {
				return err
			}
		}
	}

	// 更新报价信息
	offer.Price = newPrice
	offer.Quantity = newQuantity
	offer.Status = OfferPending
	offer.Round += 1
	offer.UpdatedTime = time.Now().Format(customFormat)

	// 创建报价修改历史记录
	if err := s.recordOfferHistory(ctx, offer, "修改"); err != nil {
		return err
	}

	// 更新账本中的Offer
	if err := s.saveOffer(ctx, offer); err != nil {
		return err
	}

	return nil
}

// CancelOffer 撤销报价
func (s *SmartContract) CancelOffer(ctx contractapi.TransactionContextInterface, userId string, offerId string, adminId string) error {
	// 获取完整的Offer信息
	offer, err := s.getOffer(ctx, offerId)
	if err != nil {
		return fmt.Errorf("获取报价信息失败: %w", err)
	}

	// 验证用户身份
	if offer.UserID != userId {
		return fmt.Errorf("用户 %s 无权撤销报价 %s", userId, offerId)
	}

	// 验证报价状态是否可撤销
	if offer.Status != OfferPending {
		return fmt.Errorf("报价 %s 状态为 %s，不可撤销", offerId, offer.Status)
	}

	// 退还买方保证金（如果有）
	if!offer.IsSeller {
		admin, err := s.getAdmin(ctx, adminId)
		if err != nil {
			return err
		}
		if err := s.refundBuyerDeposit(ctx, offer, admin, offer.Deposit); err != nil {
			return err
		}
	}

	// 更新报价状态为已取消
	offer.Status = OfferCancelled
	offer.UpdatedTime = time.Now().Format(customFormat)

	// 创建报价撤销历史记录
	if err := s.recordOfferHistory(ctx, offer, "撤销"); err != nil {
		return err
	}

	// // 从用户的Offers列表中移除该报价
	// user, err := s.getUser(ctx, userId)
	// if err != nil {
	// 	return err
	// }
	// for i, o := range user.Offers {
	// 	if o == offerId {
	// 		user.Offers = append(user.Offers[:i], user.Offers[i+1:]...)
	// 		break
	// 	}
	// }

	// // 更新用户信息
	// if err := s.updateUser(ctx, user); err != nil {
	// 	return err
	// }

	// 更新账本中的Offer
	if err := s.saveOffer(ctx, offer); err != nil {
		return err
	}

	return nil
}

// isValidOffer 检查报价是否有效
func isValidOffer(offer *Offer) bool {
	return offer.Quantity > 0 && offer.Status == OfferPending
}

// getOffer 获取报价信息
func (s *SmartContract) getOffer(ctx contractapi.TransactionContextInterface, offerId string) (*Offer, error) {
	offerKey, _ := ctx.GetStub().CreateCompositeKey(OfferPrefix, []string{offerId})
	offerJSON, err := ctx.GetStub().GetState(offerKey)
	if err != nil {
		return nil, fmt.Errorf("获取报价信息失败: %w", err)
	}
	var offer Offer
	if err := json.Unmarshal(offerJSON, &offer); err != nil {
		return nil, err
	}
	return &offer, nil
}

// saveOffer 保存报价信息
func (s *SmartContract) saveOffer(ctx contractapi.TransactionContextInterface, offer *Offer) error {
	offerKey, _ := ctx.GetStub().CreateCompositeKey(OfferPrefix, []string{offer.OfferID})
	offerJSON, _ := json.Marshal(offer)
	if err := ctx.GetStub().PutState(offerKey, offerJSON); err != nil {
		return fmt.Errorf("保存报价信息失败: %w", err)
	}
	return nil
}

// recordBalanceChange 记录余额变动
func (s *SmartContract) recordBalanceChange(ctx contractapi.TransactionContextInterface, userId string, amount float64, rest float64, reason string, isIncome bool, userName string, isAdmin bool) error {
	// 生成唯一 ID，使用自增 ID
    id, err := s.GetNextID(ctx, BalanceRecordID)
    if err != nil {
        return err
    }
	uniqueID := strconv.Itoa(id)
	balanceRecord := &BalanceRecord{
		Timestamp: time.Now().Format(customFormat),
		Amount:    amount,
		Rest:      rest,
		Reason:    reason,
		IsIncome:  isIncome,
		UserName:  userName,
	}

	// 创建余额记录的复合键
	balanceRecordKey, err := ctx.GetStub().CreateCompositeKey(BalanceRecordPrefix, []string{userId, uniqueID})
	if err != nil {
		return fmt.Errorf("创建余额记录复合键失败: %w", err)
	}

	// 序列化并存储余额记录
	balanceRecordJSON, err := json.Marshal(balanceRecord)
	if err != nil {
		return fmt.Errorf("序列化余额记录失败: %w", err)
	}
	
	if err := ctx.GetStub().PutState(balanceRecordKey, balanceRecordJSON); err != nil {
		return fmt.Errorf("存储余额记录失败: %w", err)
	}

	// 根据用户类型更新相应的实体
	if isAdmin {
		// 管理员
		admin, err := s.getAdmin(ctx, userId)
		if err != nil {
			return fmt.Errorf("获取管理员信息失败: %w", err)
		}
		admin.BalanceHistory = append(admin.BalanceHistory, balanceRecordKey)
		if err := s.updateAdmin(ctx, admin); err != nil {
			return fmt.Errorf("更新管理员信息失败: %w", err)
		}
	} else {
		// 普通用户
		user, err := s.getUser(ctx, userId)
		if err != nil {
			return fmt.Errorf("获取用户信息失败: %w", err)
		}
		user.BalanceHistory = append(user.BalanceHistory, balanceRecordKey)
		if err := s.updateUser(ctx, user); err != nil {
			return fmt.Errorf("更新用户信息失败: %w", err)
		}
	}

	return nil
}

// recordOfferHistory 记录报价历史
func (s *SmartContract) recordOfferHistory(ctx contractapi.TransactionContextInterface, offer *Offer, action string) error {
	user, err := s.getUser(ctx, offer.UserID)
	if err != nil {
		return err
	}
	offer.UpdatedTime = time.Now().Format(customFormat)
	historyRecord := &OfferHistoryRecord{
		OfferID:     offer.OfferID,
		Action:    action,
		OfferSnapshot: *offer,
	}
	// 创建历史记录的复合键，使用自增 ID
    id, err := s.GetNextID(ctx, OfferHistoryID)
    if err != nil {
        return err
    }
	uniqueID := strconv.Itoa(id)
	historyRecordKey, _ := ctx.GetStub().CreateCompositeKey(OfferHistoryPrefix, []string{offer.OfferID, uniqueID})

	// 序列化并存储历史记录
	historyRecordJSON, _ := json.Marshal(historyRecord)
	if err := ctx.GetStub().PutState(historyRecordKey, historyRecordJSON); err != nil {
		return err
	}

	// 将历史记录键添加到用户的OfferHistory
	user.OfferHistory = append(user.OfferHistory, historyRecordKey)

	return s.updateUser(ctx, user)
}