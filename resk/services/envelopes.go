package services

import (
	"github.com/shopspring/decimal"
	"joeytest.com/resk/infra/base"
	"time"
)

var IRedEnvelopeService RedEnvelopeService

func GetRedEnvelopeService() RedEnvelopeService {
	base.Check(IRedEnvelopeService)
	return IRedEnvelopeService
}

type RedEnvelopeService interface {
	// 发红包
	SendOut(dto RedEnvelopeService) (activity *RedEnvelopeActivity, err error)
	// 收红包
	Receive(dto RedEnvelopeReceiveDTO) (item *RedEnvelopeItemDTO, err error)
	// 退款
	Refund(envelopeNo string) (order *RedEnvelopeGoodsDTO)
	// 查询红包订单
	Get(envelopeNo string) (order *RedEnvelopeGoodsDTO)
}

type RedEnvelopeSendingDTO struct {
	EnvelopeType int             `json:"envelopeType" validate:"required"`     // 红包类型：普通红包、碰运气红包
	Username     string          `json:"username" validate:"required"`         // 用户名称
	UserId       string          `json:"userId" validate:"required"`           // 用户编号，红包所属用户
	Blessing     string          `json:"blessing"`                             // 祝福语
	Amount       decimal.Decimal `json:"amount" validate:"required,numeric"`   // 红包金额：普通红包指单个红包金额，碰运气红包指总金额
	Quantity     int             `json:"quantity" validate:"required,numeric"` // 红包总数量
}

type RedEnvelopeReceiveDTO struct {
	EnvelopeNo   string `json:"envelopeNo" validate:"required"`   // 红包编号，红包唯一标识
	RecvUsername string `json:"recvUsername" validate:"required"` // 红包接收者用户名称
	RecvUserId   string `json:"recvUserId" validate:"required"`   // 红包接收者用户编号
	AccountNo    string `json:"accountNo"`
}

type RedEnvelopeActivity struct {
	RedEnvelopeGoodsDTO
	Link string `json:"link"`
}

type RedEnvelopeGoodsDTO struct {
	EnvelopeNo       string          `json:"envelopeNo"`                           //红包编号,红包唯一标识
	EnvelopeType     int             `json:"envelopeType" validate:"required"`     //红包类型：普通红包，碰运气红包
	Username         string          `json:"username" validate:"required"`         //用户名称
	UserId           string          `json:"userId" validate:"required"`           //用户编号, 红包所属用户
	Blessing         string          `json:"blessing"`                             //祝福语
	Amount           decimal.Decimal `json:"amount" validate:"required,numeric"`   //红包总金额
	AmountOne        decimal.Decimal `json:"amountOne"`                            //单个红包金额，碰运气红包无效
	Quantity         int             `json:"quantity" validate:"required,numeric"` //红包总数量
	RemainAmount     decimal.Decimal `json:"remainAmount"`                         //红包剩余金额额
	RemainQuantity   int             `json:"remainQuantity"`                       //红包剩余数量
	ExpiredAt        time.Time       `json:"expiredAt" `                           //过期时间
	Status           OrderStatus     `json:"status"`                               //红包状态：0红包初始化，1启用，2失效
	OrderType        OrderType       `json:"orderType"`                            //订单类型：发布单、退款单
	PayStatus        PayStatus       `json:"payStatus"`                            //支付状态：未支付，支付中，已支付，支付失败
	CreatedAt        time.Time       `json:"createdAt"`                            //创建时间
	UpdatedAt        time.Time       `json:"updatedAt"`                            //更新时间
	AccountNo        string          `json:"accountNo"`
	OriginEnvelopeNo string          `json:"originEnvelopeNo"`
}

type RedEnvelopeItemDTO struct {
	ItemNo       string          `json:"itemNo"`       //红包订单详情编号
	EnvelopeNo   string          `json:"envelopeNo"`   //订单编号 红包编号,红包唯一标识
	RecvUsername string          `json:"recvUsername"` //红包接收者用户名称
	RecvUserId   string          `json:"recvUserId"`   //红包接收者用户编号
	Amount       decimal.Decimal `json:"amount"`       //收到金额
	Quantity     int             `json:"quantity"`     //收到数量：对于收红包来说是1
	RemainAmount decimal.Decimal `json:"remainAmount"` //收到后红包剩余金额
	AccountNo    string          `json:"accountNo"`    //红包接收者账户ID
	PayStatus    int             `json:"payStatus"`    //支付状态：未支付，支付中，已支付，支付失败
	CreatedAt    time.Time       `json:"createdAt"`    //创建时间
	UpdatedAt    time.Time       `json:"updatedAt"`    //更新时间
	IsLuckiest   bool            `json:"isLuckiest"`
	Desc         string          `json:"desc"`
}
