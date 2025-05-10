package chaincode

import (
	"encoding/json"
	"fmt"
	//"strconv"
	//"time"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

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
		return nil, fmt.Errorf("获取用户信息失败: %w", err)
	}

	var offers []*Offer
	for _, offerID := range user.Offers {
		offer, err := s.getOffer(ctx, offerID)
		if err != nil {
			// 记录错误但继续处理其他报价
			fmt.Printf("获取报价 %s 信息失败: %v\n", offerID, err)
			continue
		}
		offers = append(offers, offer)
	}

	for _, offerID := range user.OfferDone {
		offer, err := s.getOffer(ctx, offerID)
		if err != nil {
			// 记录错误但继续处理其他报价
			fmt.Printf("获取已完成报价 %s 信息失败: %v\n", offerID, err)
			continue
		}
		offers = append(offers, offer)
	}

	return offers, nil
}

// GetOfferHistory 用户查询所有报价历史信息
func (s *SmartContract) GetOfferHistory(ctx contractapi.TransactionContextInterface, userId string) ([]*OfferHistoryRecord, error) {
	user, err := s.getUser(ctx, userId)
	if err != nil {
		return nil, fmt.Errorf("获取用户信息失败: %w", err)
	}

	var historyRecords []*OfferHistoryRecord
	for _, historyKey := range user.OfferHistory {
		historyJSON, err := ctx.GetStub().GetState(historyKey)
		if err != nil {
			// 记录错误但继续处理其他历史记录
			fmt.Printf("获取报价历史记录 %s 信息失败: %v\n", historyKey, err)
			continue
		}

		var historyRecord OfferHistoryRecord
		if err := json.Unmarshal(historyJSON, &historyRecord); err != nil {
			// 记录错误但继续处理其他历史记录
			fmt.Printf("解析报价历史记录 %s 信息失败: %v\n", historyKey, err)
			continue
		}
		historyRecords = append(historyRecords, &historyRecord)
	}

	return historyRecords, nil
}

// GetBalanceHistory 用户查询账户余额历史
func (s *SmartContract) GetBalanceHistory(ctx contractapi.TransactionContextInterface, userId string) ([]*BalanceRecord, error) {
	user, err := s.getUser(ctx, userId)
	if err != nil {
		return nil, fmt.Errorf("获取用户信息失败: %w", err)
	}

	var balanceRecords []*BalanceRecord
	for _, balanceKey := range user.BalanceHistory {
		balanceJSON, err := ctx.GetStub().GetState(balanceKey)
		if err != nil {
			// 记录错误但继续处理其他余额记录
			fmt.Printf("获取余额记录 %s 信息失败: %v\n", balanceKey, err)
			continue
		}

		var balanceRecord BalanceRecord
		if err := json.Unmarshal(balanceJSON, &balanceRecord); err != nil {
			// 记录错误但继续处理其他余额记录
			fmt.Printf("解析余额记录 %s 信息失败: %v\n", balanceKey, err)
			continue
		}
		balanceRecords = append(balanceRecords, &balanceRecord)
	}

	return balanceRecords, nil
}

// GetUserContracts 查询用户参与的购电合同
func (s *SmartContract) GetUserContracts(ctx contractapi.TransactionContextInterface, userId string) ([]*Contract, error) {
	user, err := s.getUser(ctx, userId)
	if err != nil {
		return nil, fmt.Errorf("获取用户信息失败: %w", err)
	}

	var contracts []*Contract
	for _, contractID := range user.Contracts {
		contract, err := s.getContract(ctx, contractID)
		if err != nil {
			// 记录错误但继续处理其他合同
			fmt.Printf("获取合同 %s 信息失败: %v\n", contractID, err)
			continue
		}
		contracts = append(contracts, contract)
	}

	return contracts, nil
}

// GetAdminActionHistory 获取管理员操作历史
func (s *SmartContract) GetAdminActionHistory(ctx contractapi.TransactionContextInterface, adminId string) ([]*AdminActionRecord, error) {
	admin, err := s.getAdmin(ctx, adminId)
	if err != nil {
		return nil, fmt.Errorf("获取管理员信息失败: %w", err)
	}

	var actionRecords []*AdminActionRecord
	for _, actionKey := range admin.AdminActionHistory {
		actionJSON, err := ctx.GetStub().GetState(actionKey)
		if err != nil {
			// 记录错误但继续处理其他操作记录
			fmt.Printf("获取管理员操作记录 %s 信息失败: %v\n", actionKey, err)
			continue
		}

		var actionRecord AdminActionRecord
		if err := json.Unmarshal(actionJSON, &actionRecord); err != nil {
			// 记录错误但继续处理其他操作记录
			fmt.Printf("解析管理员操作记录 %s 信息失败: %v\n", actionKey, err)
			continue
		}
		actionRecords = append(actionRecords, &actionRecord)
	}

	return actionRecords, nil
}

// GetAdminMoneyHistory 获取管理员账户余额历史
func (s *SmartContract) GetAdminMoneyHistory(ctx contractapi.TransactionContextInterface, adminId string) ([]*BalanceRecord, error) {
	admin, err := s.getAdmin(ctx, adminId)
	if err != nil {
		return nil, fmt.Errorf("获取管理员信息失败: %w", err)
	}

	var balanceRecords []*BalanceRecord
	for _, balanceKey := range admin.BalanceHistory {
		balanceJSON, err := ctx.GetStub().GetState(balanceKey)
		if err != nil {
			// 记录错误但继续处理其他余额记录
			fmt.Printf("获取管理员余额记录 %s 信息失败: %v\n", balanceKey, err)
			continue
		}

		var balanceRecord BalanceRecord
		if err := json.Unmarshal(balanceJSON, &balanceRecord); err != nil {
			// 记录错误但继续处理其他余额记录
			fmt.Printf("解析管理员余额记录 %s 信息失败: %v\n", balanceKey, err)
			continue
		}
		balanceRecords = append(balanceRecords, &balanceRecord)
	}

	return balanceRecords, nil
}

// AdminGetAllOffer 查询系统内所有报价
func (s *SmartContract) AdminGetAllOffer(ctx contractapi.TransactionContextInterface) ([]*Offer, error) {
    // 使用 Offer 复合键前缀直接查询所有报价
    resultsIterator, err := ctx.GetStub().GetStateByPartialCompositeKey(OfferPrefix, []string{})
    if err != nil {
        return nil, fmt.Errorf("获取所有报价信息失败: %w", err)
    }
    defer resultsIterator.Close()

    var offerList []*Offer
    for resultsIterator.HasNext() {
        queryResponse, err := resultsIterator.Next()
        if err != nil {
            return nil, fmt.Errorf("迭代报价信息失败: %w", err)
        }

        var offer Offer
        if err := json.Unmarshal(queryResponse.Value, &offer); err != nil {
            fmt.Printf("解析报价信息失败: %v\n", err)
            continue
        }
        offerList = append(offerList, &offer)
    }

    return offerList, nil
}
// GetAllContract 查询系统内所有合同
func (s *SmartContract) GetAllContract(ctx contractapi.TransactionContextInterface) ([]*Contract, error) {
	resultsIterator, err := ctx.GetStub().GetStateByPartialCompositeKey(ContractPrefix, []string{})
	if err != nil {
		return nil, fmt.Errorf("获取所有合同信息失败: %w", err)
	}
	defer resultsIterator.Close()

	var contractlist []*Contract
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return nil, fmt.Errorf("迭代合同信息失败: %w", err)
		}

		var contract Contract
		if err := json.Unmarshal(queryResponse.Value, &contract); err != nil {
			// 不是合同信息，跳过
			continue
		}
		contractlist = append(contractlist, &contract)
	}

	return contractlist, nil
}