package SendEMail

import (
	"admin/common"
	"admin/fusion"
	"github.com/jordan-wright/email"
	"github.com/leonelquinteros/gotext"
	"log"
	"net/smtp"
)

func SendEmail(content string) {
	auth := smtp.PlainAuth("", common.MyCfg.Email.EmailAccount,
		common.MyCfg.Email.EmailAuthorizationCode, common.MyCfg.Email.EmailServerHost)
	var receivers = make([]string, 0)
	var db_default = fusion.GetBaseGormDB("default")
	db_default.Table("goadmin_users").Select("eMail").Scan(&receivers)

	for _, value := range receivers {
		if value == "" {
			continue
		}
		e := &email.Email{
			From:    gotext.Get("后台管理服务器") + " <" + common.MyCfg.Email.EmailAccount + ">",
			To:      []string{value},
			Subject: gotext.Get("后台管理服务器通知！"),
			Text:    []byte(content),
		}
		err := e.Send(common.MyCfg.Email.EmailServerHost+common.MyCfg.Email.EmailServerPort, auth)
		if err != nil {
			log.Fatal("cant send mail account:" + value + "   " + err.Error())
		}
	}
}
