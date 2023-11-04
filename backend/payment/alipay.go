package payment

import (
	"github.com/go-pay/gopay/alipay"
	"github.com/go-pay/gopay/pkg/xlog"
)

func InitAlipay() {
	client, err := alipay.NewClient("2021004124640823", privateKey, false)
	if err != nil {
		xlog.Error(err)
		return
	}
}
