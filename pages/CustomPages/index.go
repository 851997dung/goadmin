package CustomPages

import (
	"admin/common"
	"admin/common/def"
	"admin/fusion"
	"fmt"
	"html/template"
	"strings"

	"github.com/GoAdminGroup/go-admin/context"
	"github.com/GoAdminGroup/go-admin/modules/auth"
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
)

func GetIndex(ctx *context.Context) (types.Panel, error) {
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
	info.AddField(gotext.Get("切换语言"), "Language", db.Varchar, form.SelectSingle).
		FieldOptions(types.FieldOptions{
			{Text: gotext.Get("中文"), Value: "cn"},
			{Text: gotext.Get("越南文"), Value: "vi"},
			//{Text: gotext.Get("英文"), Value: "en"},
		}).FieldValue(cfg.Language)

	user := auth.Auth(ctx)
	if !user.CheckRole("playerManager") {
		info.AddField(gotext.Get("中心服host"), "host", db.Varchar, form.Ip).
			//FieldMust().
			FieldValue(cfg.Center.Host)
		info.AddField(gotext.Get("中心服端口"), "port", db.Varchar, form.Text).
			//FieldMust().
			FieldValue(cfg.Center.Port)

		info.AddField(gotext.Get("备份服务器host"), "backupHost", db.Varchar, form.Ip).
			//FieldMust().
			FieldValue(cfg.BackUp.Host)
		info.AddField(gotext.Get("备份服务器端口"), "backupPort", db.Varchar, form.Select).
			//FieldMust().
			FieldValue(cfg.BackUp.Port)
		info.AddField(gotext.Get("国战跨服主服务器ID"), "MainServer", db.Varchar, form.Text).
			//FieldMust().
			FieldValue(cfg.Api["MainServer"])

		var jumpUrl string
		db_global := fusion.GetBaseGormDB("db_global")
		if db_global != nil {
			db_global.Table("web_address").
				Select("addr").
				Where("ID = 1").
				Scan(&jumpUrl)
		}
		info.AddField(gotext.Get("跳转网页"), "jumpUrl", db.Varchar, form.Text).
			//FieldMust().
			FieldValue(jumpUrl)
	}
	aform := components.Form().
		SetHeader(info.Header).
		SetContent(info.FieldList).
		SetPrefix(config.PrefixFixSlash()).
		SetUrl("/admin/SaveConfig").
		SetHiddenFields(map[string]string{
			form2.PreviousKey: "/admin",
		}).
		SetMethod("GET").
		SetOperationFooter(col1 + col2)

	//popup := components.Popup().SetTitle("test").SetID("/popup/test").SetBody("wdnmd")

	return types.Panel{
		Content: components.Box().
			SetHeader(aform.GetDefaultBoxHeader(true)).
			WithHeadBorder().
			SetBody(aform.GetContent() + info.FooterHtml /*+ popup.GetContent()*/).
			GetContent(),
		Title:     template.HTML(gotext.Get("配置管理")),
		Callbacks: nil,
		//JS:        template.JS(GetJs()),
	}, nil
}
func SaveConfig(ctx *context.Context) (types.Panel, error) {
	Language := ctx.Request.FormValue("Language")
	centerHost := ctx.Request.FormValue("host")
	centerPort := ctx.Request.FormValue("port")

	backupHost := ctx.Request.FormValue("backupHost")
	backupPort := ctx.Request.FormValue("backupPort")
	MainServer := ctx.Request.FormValue("MainServer")

	jumpUrl := ctx.FormValue("jumpUrl")

	go_manager := fusion.GetBaseGormDB("go_manager")
	if centerHost != "" {
		common.MyCfg.Center.Host = centerHost
		go_manager.Table("go_config").Where("key", "centerHost").Update("value", centerHost)
	}
	if centerPort != "" {
		common.MyCfg.Center.Port = centerPort
		go_manager.Table("go_config").Where("key", "centerPort").Update("value", centerPort)
	}

	if backupHost != "" {
		common.MyCfg.BackUp.Host = backupHost
		go_manager.Table("go_config").Where("key", "backupHost").Update("value", backupHost)
	}
	if backupPort != "" {
		common.MyCfg.BackUp.Port = backupPort
		go_manager.Table("go_config").Where("key", "backupPort").Update("value", backupPort)
	}
	if MainServer != "" {
		common.MyCfg.Api["MainServer"] = MainServer
		go_manager.Table("go_config").Where("key", "MainServer").Update("value", MainServer)
	}
	if Language != "" {
		common.MyCfg.Language = Language
		go_manager.Table("go_config").Where("key", "language").Update("value", Language)
		gotext.Configure("path/to/locales", common.MyCfg.Language, common.MyCfg.Language)
		defaultData := common.AdminEngine.DefaultConnection().GetConfig("default")
		go_admin := fusion.OpenDB(&def.DataBase{
			Host:         defaultData.Host,
			Port:         defaultData.Port,
			User:         defaultData.User,
			Pwd:          defaultData.Pwd,
			DataBaseName: defaultData.Name,
		})
		go_admin.Table("goadmin_site").Where("key", "language").Update("value", Language)
		var m = make(map[string]string)
		if Language == "vi" {
			Language = "en"
		}
		m["language"] = Language
		common.AdminEngine.UpdateConfig(m)
	}

	//if jumpUrl != "" {
	db_global := fusion.GetBaseGormDB("db_global")
	if db_global != nil {
		db_global.Table("web_address").Where("ID", 1).Update("addr", jumpUrl)
	}
	//}

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
		Title:       template.HTML(gotext.Get("更新配置")),
		Description: template.HTML(fmt.Sprintf(`<a href="/%s">%s</a>`, preUrl, str)),
		Callbacks:   nil,
	}, nil
}

//func GetJs() string {
//	msg := "fuck you admin"
//	return "window.setInterval(function () {\n" +
//		"    var refreshHours = new Date().getHours();\n" +
//		"    var refreshMin = new Date().getMinutes();\n" +
//		"    var refreshSec = new Date().getSeconds();\n" +
//		"    if (refreshMin == '13') {\n" +
//		"        $.ajax({\n" +
//		"            type: \"GET\",\n" +
//		"            url: '/admin/SendComplaint',\n" +
//		"            success: function (msg) {\n" +
//		"                alert(\"" + msg + "\");\n" +
//		"            },\n" +
//		"        })\n" +
//		"    }\n" +
//		"}, 1000);"
//}
