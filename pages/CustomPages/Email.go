package CustomPages

import (
	"admin/common"
	"admin/fusion"
	"fmt"
	"github.com/GoAdminGroup/go-admin/context"
	"github.com/GoAdminGroup/go-admin/modules/config"
	"github.com/GoAdminGroup/go-admin/modules/db"
	"github.com/GoAdminGroup/go-admin/modules/language"
	form2 "github.com/GoAdminGroup/go-admin/plugins/admin/modules/form"
	template2 "github.com/GoAdminGroup/go-admin/template"
	"github.com/GoAdminGroup/go-admin/template/icon"
	"github.com/GoAdminGroup/go-admin/template/types"
	"github.com/GoAdminGroup/go-admin/template/types/form"
	"github.com/GoAdminGroup/themes/adminlte/components/infobox"
	"github.com/leonelquinteros/gotext"
	"html/template"
	"strings"
)

func GetEmailSettingPage(ctx *context.Context) (types.Panel, error) {
	components := template2.Get(config.GetTheme())
	col1 := components.Col().GetContent()
	btn1 := components.Button().SetType("submit").
		SetContent(language.GetFromHtml("Save")).
		SetThemePrimary().
		SetOrientationRight().
		SetLoadingText(icon.Icon("fa-spinner fa-spin", 2) + `Save`).
		GetContent()
	btn2 := components.Button().SetType("reset").
		SetContent(language.GetFromHtml("Reset")).
		SetThemeWarning().
		SetOrientationLeft().
		GetContent()
	col2 := components.Col().SetSize(types.SizeMD(8)).
		SetContent(btn1 + btn2).GetContent()

	var cfg = common.MyCfg
	info := types.NewFormPanel()
	info.AddField(gotext.Get("邮箱账号"), "EmailAccount", db.Varchar, form.Email).
		FieldMust().
		FieldValue(cfg.Email.EmailAccount).
		FieldHelpMsg(template.HTML(gotext.Get("用于发送紧急邮件的邮箱")))
	info.AddField(gotext.Get("授权码"), "EmailAuthorizationCode", db.Varchar, form.Text).
		FieldMust().
		FieldValue(cfg.Email.EmailAuthorizationCode).
		FieldHelpMsg(template.HTML(gotext.Get("邮箱的授权码")))
	info.AddField(gotext.Get("Host"), "EmailServerHost", db.Varchar, form.Url).
		FieldMust().
		FieldValue(cfg.Email.EmailServerHost).
		FieldHelpMsg(template.HTML(gotext.Get("对应邮箱的STMP服务器Host")))
	info.AddField(gotext.Get("Port"), "EmailServerPort", db.Varchar, form.Text).
		FieldMust().
		FieldValue(cfg.Email.EmailServerPort).
		FieldHelpMsg(template.HTML(gotext.Get("对应邮箱的STMP服务器Port")))

	aform := components.Form().
		SetHeader(info.Header).
		SetContent(info.FieldList).
		SetPrefix(config.PrefixFixSlash()).
		SetUrl("/admin/UpdateEmailSetting").
		SetHiddenFields(map[string]string{
			form2.PreviousKey: "/admin",
		}).
		SetMethod("GET").
		SetOperationFooter(col1 + col2)

	return types.Panel{
		Content: components.Box().
			SetHeader(aform.GetDefaultBoxHeader(true)).
			WithHeadBorder().
			SetBody(aform.GetContent() + info.FooterHtml).
			GetContent(),
		Title:     template.HTML(gotext.Get("电子邮件配置")),
		Callbacks: nil,
	}, nil
}

func UpdateEmailSetting(ctx *context.Context) (types.Panel, error) {
	EmailAccount := ctx.Request.FormValue("EmailAccount")
	EmailAuthorizationCode := ctx.Request.FormValue("EmailAuthorizationCode")
	EmailServerHost := ctx.Request.FormValue("EmailServerHost")
	EmailServerPort := ctx.Request.FormValue("EmailServerPort")

	go_manager := fusion.GetBaseGormDB("go_manager")
	if EmailAccount != "" {
		common.MyCfg.Email.EmailAccount = EmailAccount
		go_manager.Table("go_config").Where("key", "EmailAccount").Update("value", EmailAccount)
	}
	if EmailAuthorizationCode != "" {
		common.MyCfg.Email.EmailAuthorizationCode = EmailAuthorizationCode
		go_manager.Table("go_config").Where("key", "EmailAuthorizationCode").Update("value", EmailAuthorizationCode)
	}
	if EmailServerHost != "" {
		common.MyCfg.Email.EmailServerHost = EmailServerHost
		go_manager.Table("go_config").Where("key", "EmailServerHost").Update("value", EmailServerHost)
	}
	if EmailServerPort != "" {
		common.MyCfg.Email.EmailServerPort = EmailServerPort
		go_manager.Table("go_config").Where("key", "EmailServerPort").Update("value", EmailServerPort)
	}

	msg := gotext.Get("操作成功")
	if go_manager.Error != nil {
		msg = go_manager.Error.Error()
	}
	components := template2.Get(config.GetTheme())
	infobox1 := infobox.New().SetText(template.HTML(msg))
	str := gotext.Get("返回上级")
	preUrl := ctx.Request.Header.Get("Referer")
	index := strings.Index(preUrl, "admin")
	preUrl = preUrl[index:]
	col1 := components.Col()
	var size = map[string]string{"md": "3", "sm": "6", "xs": "12"}
	infoboxCol1 := col1.SetSize(size).SetContent(infobox1.GetContent()).GetContent()
	return types.Panel{
		Content:     infoboxCol1,
		Title:       template.HTML(gotext.Get("更新邮件发送配置")),
		Description: template.HTML(fmt.Sprintf(`<a href="/%s">%s</a>`, preUrl, str)),
		Callbacks:   nil,
	}, nil
}
