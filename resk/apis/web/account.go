package web

import (
	"github.com/kataras/iris/v12"
	"github.com/sirupsen/logrus"
	"joeytest.com/resk/infra/base"
	"joeytest.com/resk/services"
)

func init() {
	accountApi := new(AccountApi)
	logrus.Info("accountApi:", accountApi)
	base.RegisterApi(accountApi)
}

type AccountApi struct {
	service services.AccountService
}

func (a *AccountApi) Init() {
	a.service = services.GetAccountService()
	logrus.Info("GetAccountService:", services.GetAccountService())
	groupRouter := base.Iris().Party("/v1/account")
	logrus.Info("groupRouter:", groupRouter)
	groupRouter.Post("/create", a.creatHandler)
}

//账户创建的接口: /v1/account/create
func (a *AccountApi) creatHandler(ctx iris.Context) {
	//获取请求参数
	account := services.AccountCreatedDTO{}
	err := ctx.ReadJSON(&account)
	r := base.Res{
		Code: base.ResCodeOk,
	}
	if err != nil {
		r.Code = base.ResCodeRequestParamsError
		r.Message = err.Error()
		ctx.JSON(&r)
		logrus.Error(err)
		return
	}
	logrus.Info("执行创建账户的密码")
	logrus.Info("service:", a.service)
	// 执行创建账户的代码
	dto, err := a.service.CreateAccount(account)
	if err != nil {
		r.Code = base.ResCodeInnerServerError
		r.Message = err.Error()
		logrus.Error(err)
		ctx.JSON(&r)
		return
	}
	r.Data = dto
	ctx.JSON(&r)
}

//转账的接口 :/v1/account/transfer
func (a *AccountApi) transferHandler(ctx iris.Context) {
	account := services.AccountTransferDTO{}
	err := ctx.ReadJSON(&account)
	r := base.Res{
		Code: base.ResCodeOk,
	}
	if err != nil {
		r.Code = base.ResCodeRequestParamsError
		r.Message = err.Error()
		ctx.JSON(&r)
		logrus.Error(err)
		return
	}

	// 执行转账逻辑
	status, err := a.service.Transfer(account)
	if err != nil {
		r.Code = base.ResCodeInnerServerError
		r.Message = err.Error()
		logrus.Error(err)
		ctx.JSON(r)
		return
	}
	if status != services.TransferStatusSuccess {
		r.Code = base.ResCodeBizError
		r.Message = err.Error()
	}
	r.Data = status
	ctx.JSON(&r)
}

func (a *AccountApi) getEnvelopeAccountHandler(ctx iris.Context) {
	userId := ctx.URLParam("userId")
	r := base.Res{
		Code: base.ResCodeOk,
	}
	if userId == "" {
		r.Code = base.ResCodeRequestParamsError
		r.Message = "用户ID不能为空"
		ctx.JSON(&r)
		return
	}
	account := a.service.GetEnvelopeAccountByUserId(userId)
	r.Data = account
	ctx.JSON(&r)
}

//查询账户信息的web接口：/v1/account/get
func (a *AccountApi) getAccountHandler(ctx iris.Context) {
	accountNo := ctx.URLParam("accountNo")
	r := base.Res{
		Code: base.ResCodeOk,
	}
	if accountNo == "" {
		r.Code = base.ResCodeRequestParamsError
		r.Message = "账户编号不能为空"
		ctx.JSON(&r)
		return
	}
	account := a.service.GetAccount(accountNo)
	r.Data = account
	ctx.JSON(&r)
}
