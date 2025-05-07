package chaincode

const (
	OfferPrefix = "Offer"
	ContractPrefix = "Contract"
)
const (
	INIT_BALANCE = 50
    customFormat = "2006-01-02 15:04:05"
)

var (
	MatchFrequency = 10 // 撮合频率，单位为分钟
	DepositRate    = 0.1 // 保证金率
	FeeRate        = 0.02 // 手续费率
)

type Config struct {
	MatchFrequency int     `json:"matchFrequency"` // 撮合频率，单位为分钟
	DepositRate    float64 `json:"depositRate"`    // 保证金率
	FeeRate        float64 `json:"feeRate"`       // 手续费率
}
type ApprovalStatus string
const (
    NotApplied   ApprovalStatus = "未申请"
    Pending      ApprovalStatus = "申请中" 
    Approved     ApprovalStatus = "已通过"
    Rejected     ApprovalStatus = "未通过"
)

type AuditStatus string
const (
    AuditPending  AuditStatus = "审核中"
    AuditPassed   AuditStatus = "审核通过"
    AuditRejected AuditStatus = "审核拒绝"
)
// User 表示用户信息
type User struct {
	UserID             string                `json:"userId"`
	Balance            float64               `json:"balance"`
	IsSeller           ApprovalStatus                 `json:"isSeller"`
	IsBuyer            ApprovalStatus            `json:"isBuyer"`
	//ApproveUserAsSeller bool                  `json:"approveUserAsSeller"`
	//ApproveUserAsBuyer  bool                  `json:"approveUserAsBuyer"`
	Offers             []*Offer             `json:"offers"`          // 用户的报价列表
	Contracts          []*Contract          `json:"contracts"`       // 用户的购电合同 ID 列表
	BalanceHistory     []*BalanceRecord     `json:"balanceHistory"` // 余额变动历史记录
	OfferHistory       []*OfferHistoryRecord `json:"offerHistory"`   // 报价历史记录
	OfferDone          []*Offer             `json:"offerDone"`
    CreditRating int     `json:"creditRating"` // 信用评级（0-100）
    TradeCount   int     `json:"tradeCount"`   // 累计交易次数
}
type OfferStatus string
const (
    OfferPending   OfferStatus = "待撮合"
    OfferMatched   OfferStatus = "已撮合"
    OfferCancelled OfferStatus = "已取消"
)
// Offer 表示报价信息
type Offer struct {
	OfferID   string  `json:"offerId"`
	UserID    string  `json:"userId"`
	Price     float64 `json:"price"`
	Quantity  int     `json:"quantity"`
	Deposit   float64 `json:"deposit"`   // 保证金
	IsSeller  bool    `json:"isSeller"`
	Timestamp string  `json:"timestamp"` // 报价时间
    UpdatedTime string      `json:"updatedTime"` // 更新时间
	Status    OfferStatus  `json:"status"`    // 报价状态
	Round     int     `json:"round"`     // 撮合轮数
}
type ContractStatus string
const (
    ContractCreated  ContractStatus = "created"
    ContractSettled  ContractStatus = "settled"
    ContractCanceled ContractStatus = "canceled"
)
// Contract 表示购电合同
type Contract struct {
	ContractID string  `json:"contractId"`
	SellerID   string  `json:"sellerId"`
	BuyerID    string  `json:"buyerId"`
	Price      float64 `json:"price"`
	Quantity   int     `json:"quantity"`
	Timestamp  string  `json:"timestamp"` // 重命名原Timestamp
    SettledTime  string  `json:"settledTime"` // 新增结算时间
    TradeAmount  float64    `json:"tradeAmount"` // 实际交易金额
    Status       ContractStatus  `json:"status"`      // 状态：created/settled/cancelled
}

// BalanceRecord 余额变动记录
type BalanceRecord struct {
	Timestamp string  `json:"timestamp"`
	Amount    float64 `json:"amount"`
	Rest      float64 `json:"rest"`
	Reason    string  `json:"reason"`
}

// OfferHistoryRecord 报价历史记录
type OfferHistoryRecord struct {
	Offer     *Offer  `json:"offer"`
	Timestamp string  `json:"timestamp"`
	Action    string  `json:"action"` // 如 "提交", "修改", "完成"
}

// PlatformAdmin 表示平台管理员
type PlatformAdmin struct {
	AdminID            string                `json:"adminId"`
	Balance            float64               `json:"balance"`          // 存储收到的保证金和手续费
	BalanceHistory     []*BalanceRecord      `json:"balanceHistory"`   // 余额变动历史记录
	AdminActionHistory []*AdminActionRecord  `json:"adminActionHistory"` // 管理员操作历史记录
    Applications []*Application `json:"applications"` // 待审核列表
	//SellList           []string              `json:"sellList"`         // 审核列表
	//BuyList            []string              `json:"buyList"`          // 审核列表
	Contracts          []string          `json:"contracts"`        // 用户的购电合同 ID 列表
}

type Application struct {
    ApplicationID string      `json:"applicationId"`
    UserID        string      `json:"userId"`
    ApplyType     string      `json:"applyType"` // "buy" or "sell"
    ApplyTime     string      `json:"applyTime"`
    AuditStatus   AuditStatus `json:"auditStatus"`
    AuditTime     string      `json:"auditTime"`
    //Reason        string      `json:"reason"`      // 拒绝原因
}

// AdminActionRecord 管理员操作记录
type AdminActionRecord struct {
	Action    string `json:"action"`
	Timestamp string `json:"timestamp"`
	Details   string `json:"details"`
}




