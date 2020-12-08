package services

// 转账状态
type TransferStatus int8

const (
	// 转账失败
	TransferStatusFailure TransferStatus = -1
	// 余额不足
	TransferStatusSufficientFunds TransferStatus = 0
	// 转账成功
	TransferStatusSuccess TransferStatus = 1
)

// 转账的类型： 0 = 创建账户  >= 1 进账  <= 1 支出
type ChangeType int8

const (
	// 账户创建
	AccountCreated ChangeType = 0
	// 储值
	AccountStoreValue ChangeType = 1
	// 红包资金的支出
	EnvelopeOutgoing ChangeType = -2
	// 红包资金的收入
	EnvelopeIncoming ChangeType = 2
	// 红包退款
	EnvelopeExpired ChangeType = 3
)

// 资金交易变化标识
type ChangeFlag int8

const (
	// 创建账户
	FlagAccountCreated ChangeFlag = 0
	// 支出
	FlagTransferOut ChangeFlag = -1
	// 收入
	FlagTransferIn ChangeFlag = 1
)

// 账户类型
type AccountType int8

const (
	EnvelopeAccountType       AccountType = 1
	SystemEnvelopeAccountType AccountType = 2
)
