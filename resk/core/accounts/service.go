package accounts

import (
	"errors"
	"github.com/shopspring/decimal"
	"joeytest.com/resk/infra/base"
	"joeytest.com/resk/services"
)

type accountService struct {
}

func (a *accountService) CreateAccount(dto services.AccountCreatedDTO) (*services.AccountDTO, error) {
	domain := accountDomain{}
	// 验证输入参数
	err := base.ValidateStruct(&dto)
	if err != nil {
		return nil, err
	}
	// 执行账户业务逻辑
	amount, err := decimal.NewFromString(dto.Amount)
	if err != nil {
		return nil, err
	}
	account := services.AccountDTO{
		AccountName:  dto.AccountName,
		AccountType:  dto.AccountType,
		CurrencyCode: dto.CurrencyCode,
		UserId:       dto.UserId,
		Username:     dto.Username,
		Balance:      amount,
		Status:       1,
	}
	rdto, err := domain.Create(account)
	return rdto, err
}

func (a *accountService) Transfer(dto services.AccountTransferDTO) (services.TransferStatus, error) {
	domain := accountDomain{}
	// 验证输入参数
	err := base.ValidateStruct(&dto)
	if err != nil {
		return services.TransferStatusFailure, err
	}
	// 执行账户业务逻辑
	amount, err := decimal.NewFromString(dto.AmountStr)
	if err != nil {
		return services.TransferStatusFailure, err
	}
	dto.Amount = amount
	if dto.ChangeFlag == services.FlagTransferOut {
		if dto.ChangeType > 0 {
			return services.TransferStatusFailure, errors.New("如果changeFlag为支出，那么changeType必须小于0")
		}
	} else {
		if dto.ChangeType < 0 {
			return services.TransferStatusFailure, errors.New("如果changeFlag为收入，那么changeType必须大于0")
		}
	}
	status, err := domain.Transfer(dto)
	return status, err
}

func (a *accountService) StoreValue(dto services.AccountTransferDTO) (services.TransferStatus, error) {
	dto.TradeTarget = dto.TradeBody
	dto.ChangeFlag = services.FlagTransferIn
	dto.ChangeType = services.AccountStoreValue
	return a.Transfer(dto)
}

func (a *accountService) GetEnvelopeAccountByUserId(userId string) *services.AccountDTO {
	panic("implement me")
}
