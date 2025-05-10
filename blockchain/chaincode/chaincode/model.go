package chaincode

const (
	OfferPrefix = "Offer"//前端提交生成
	ContractPrefix = "Contract"//admin自增
	BalanceRecordPrefix = "BalanceRecord"//账本存储自增
	OfferHistoryPrefix = "OfferHistory"
	AdminActionPrefix = "AdminAction"
)
const (
	INIT_BALANCE = 50
    customFormat = "2006-01-02 15:04:05"
	AdminName = "admin"
	OfferHistoryID= "OfferHistoryID"
	BalanceRecordID= "BalanceRecordID"
	AdminActionID= "AdminActionID"
)

var (
	MatchFrequency = 2 // 撮合频率，单位为分钟
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
	UserName           string                `json:"userName"`
	Balance            float64               `json:"balance"`
	IsSeller           ApprovalStatus                 `json:"isSeller"`
	IsBuyer            ApprovalStatus            `json:"isBuyer"`
	Offers             []string               `json:"offers"`          // 用户的报价 ID 列表
	Contracts          []string          `json:"contracts"`       // 用户的购电合同 ID 列表
	BalanceHistory     []string     `json:"balanceHistory"` // 存储 BalanceRecord 的键
	OfferHistory       []string      `json:"offerHistory"`   // 存储 OfferHistoryRecord 的键
	OfferDone          []string             `json:"offerDone"`
    CreditRating int     `json:"creditRating"` // 信用评级（0-100）
    TradeCount   int     `json:"tradeCount"`   // 累计交易次数
}
type OfferStatus string
const (
    OfferPending   OfferStatus = "待撮合"
    OfferMatched   OfferStatus = "已撮合"
    OfferCancelled OfferStatus = "已取消"   //提交，修改，更新，完成，撤销
)
// Offer 表示报价信息
type Offer struct {
	OfferID   string  `json:"offerId"`
	UserID    string  `json:"userId"`
	UserName  string  `json:"userName"`
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
	SellerName string  `json:"sellerName"`
	BuyerName  string  `json:"buyerName"`
	sellerOfferID string  `json:"sellerOfferId"`
	BuyerOfferID string  `json:"buyerOfferId"`
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
	IsIncome bool `json:"isIncome"` // true 表示收入，false 表示支出
	UserName  string  `json:"userName"`
}

// OfferHistoryRecord 报价历史记录
type OfferHistoryRecord struct {
	OfferID     string  `json:"offerId"`
	//Timestamp   string  `json:"timestamp"`
	Action      string  `json:"action"` // 提交，修改，更新，完成，撤销
	OfferSnapshot Offer  `json:"offerSnapshot"` // 存储操作时的报价快照
}

// PlatformAdmin 表示平台管理员
type PlatformAdmin struct {
	AdminID            string   `json:"adminId"`
	Balance            float64  `json:"balance"`
	BalanceHistory     []string `json:"balanceHistory"`     // 改为存储 BalanceRecord 的键
	AdminActionHistory []string `json:"adminActionHistory"` // 改为存储 AdminActionRecord 的键
	Applications       []*Application `json:"applications"`       //
	Contracts          []string `json:"contracts"`          // ID
	Contractnumber     int      `json:"contractnumber"`
}

type Application struct {
    ApplicationID string      `json:"applicationId"`
    UserID        string      `json:"userId"`
	UserName      string      `json:"userName"`
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




