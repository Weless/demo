package accounts

import (
	"github.com/segmentio/ksuid"
	"github.com/shopspring/decimal"
	"github.com/sirupsen/logrus"
	. "github.com/smartystreets/goconvey/convey"
	"github.com/tietang/dbx"
	"joeytest.com/resk/infra/base"
	"joeytest.com/resk/services"
	"testing"
)

func TestAccountLogDao(t *testing.T) {
	err := base.Tx(func(runner *dbx.TxRunner) error {
		dao := AccountLogDao{
			runner: runner,
		}
		Convey("测试通过流水编号查询流水数据", t, func() {
			a := &AccountLog{
				LogNo:      ksuid.New().Next().String(),
				TradeNo:    ksuid.New().Next().String(),
				AccountNo:  ksuid.New().Next().String(),
				UserId:     ksuid.New().Next().String(),
				Username:   "测试用户",
				Amount:     decimal.NewFromFloat(1),
				Balance:    decimal.NewFromFloat(100),
				ChangeType: services.AccountCreated,
				ChangeFlag: services.FlagAccountCreated,
				Status:     1,
			}
			n, err := dao.Insert(a)
			So(err, ShouldBeNil)
			So(n, ShouldBeGreaterThan, 0)
			Convey("通过流水编号来查询流水数据", func() {
				na := dao.GetOne(a.LogNo)
				So(na, ShouldNotBeNil)
				So(na.Balance.String(), ShouldEqual, a.Balance.String())
				So(na.Amount.String(), ShouldEqual, a.Amount.String())
				So(na.CreatedAt, ShouldNotBeNil)
			})
			Convey("通过交易编号来查询流水记录", func() {
				na := dao.GetByTradeNo(a.TradeNo)
				So(na, ShouldNotBeNil)
				So(na.Balance.String(), ShouldEqual, a.Balance.String())
				So(na.Amount.String(), ShouldEqual, a.Amount.String())
				So(na.CreatedAt, ShouldNotBeNil)
			})
		})
		return nil
	})
	if err != nil {
		logrus.Error(err)
	}
}
