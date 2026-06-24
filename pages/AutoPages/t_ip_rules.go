package AutoPages

import (
	"admin/common"
	"admin/fusion"
	"html/template"
	"net/url"
	"strconv"

	"github.com/GoAdminGroup/go-admin/context"
	"github.com/GoAdminGroup/go-admin/modules/db"
	"github.com/GoAdminGroup/go-admin/plugins/admin/modules/table"
	"github.com/GoAdminGroup/go-admin/template/types"
	"github.com/GoAdminGroup/go-admin/template/types/action"
	"github.com/GoAdminGroup/go-admin/template/types/form"
	"github.com/leonelquinteros/gotext"
)

func GetTIpRulesTable(ctx *context.Context) table.Table {

	tIpRules := table.NewDefaultTable(table.Config{
		Driver:     "mysql",
		Connection: "db_global",
		CanAdd:     true,
		Editable:   true,
		Deletable:  true,
		Exportable: true,
		PrimaryKey: table.PrimaryKey{
			Type: db.Char,
			Name: "IP",
		},
	})

	info := tIpRules.GetInfo().HideRowSelector()

	info.AddField("IP", "IP", db.Char)
	info.AddField("Rule", "rule", db.Int).
		FieldDisplay(func(model types.FieldModel) interface{} {
			temp, _ := strconv.Atoi(model.Value)
			if temp == -1 {
				return gotext.Get("黑名单")
			}
			return gotext.Get("白名单")
		})

	info.SetTable("t_ip_rules").SetTitle("IpRules").HideDetailButton()

	formList := tIpRules.GetForm()
	formList.AddField("IP", "IP", db.Char, form.Text).FieldMust()
	formList.AddField("Rule", "rule", db.Int, form.SelectSingle).FieldOptions(types.FieldOptions{
		{Text: gotext.Get("黑名单"), Value: "-1"},
		{Text: gotext.Get("白名单"), Value: "1"},
	}).FieldMust()

	formList.SetTable("t_ip_rules").SetTitle("IpRules").HideContinueNewCheckBox().HideContinueEditCheckBox()

	info.AddButton(template.HTML(gotext.Get("通知到服务器")), "",
		action.PopUp("/admin/notifyIpRule", gotext.Get("操作结果"), func(ctx *context.Context) (success bool, msg string, data interface{}) {
			faild := []string{}
			vars := &url.Values{}
			ctx1 := &context.Context{}
			err, res := fusion.CallToCenter(vars, common.MyCfg.Api["ReloadIPRules"], "GET", ctx1)
			if err != nil {
				faild = append(faild, "<h5>"+err.Error()+"</h5>")
			}
			if res != "" {
				faild = append(faild, "<h5>"+res+"</h5>")
			}
			return true, "", faild
		}))

	return tIpRules
}
