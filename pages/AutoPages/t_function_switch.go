package AutoPages

import (
	"admin/common"
	"admin/common/def"
	"admin/fusion"
	"net/url"
	"strings"
	"time"

	"github.com/GoAdminGroup/go-admin/context"
	"github.com/GoAdminGroup/go-admin/modules/db"
	"github.com/GoAdminGroup/go-admin/plugins/admin/modules/table"
	"github.com/GoAdminGroup/go-admin/template"
	"github.com/GoAdminGroup/go-admin/template/types"
	"github.com/GoAdminGroup/go-admin/template/types/action"
	"github.com/GoAdminGroup/go-admin/template/types/form"
	"github.com/leonelquinteros/gotext"
)

func GetTFunctionSwitchTable(ctx *context.Context) table.Table {

	tFunctionSwitch := table.NewDefaultTable(table.Config{
		Driver:     "mysql",
		Connection: "db_global",
		CanAdd:     true,
		Editable:   true,
		Deletable:  true,
		Exportable: true,
		PrimaryKey: table.PrimaryKey{
			Type: db.Int,
			Name: "Id",
		},
	})

	info := tFunctionSwitch.GetInfo()

	info.AddField("Id", "Id", db.Int)
	info.AddField(gotext.Get("开关名字"), "switchName", db.Varchar).FieldDisplay(func(model types.FieldModel) interface{} {
		return gotext.Get(model.Value)
	})
	info.AddField(gotext.Get("开启该功能的服务器"), "gsIds", db.Mediumtext)
	info.AddField(gotext.Get("关闭该功能的服务器"), "gsNotIds", db.Mediumtext)
	info.AddField(gotext.Get("更新时间"), "updateTime", db.Datetime)

	info.SetTable("t_function_switch").SetTitle(gotext.Get("功能开关")).HideDeleteButton().HideDetailButton().HideRowSelector().FieldDisplay(func(model types.FieldModel) interface{} {
		return gotext.Get(model.Value)
	})

	info.AddButton(template.HTML(gotext.Get("将选中条目同步到服务器")), "",
		action.PopUp("/functionSwitch/send", gotext.Get("操作结果"), func(ctx *context.Context) (success bool, msg string, data interface{}) {
			failed := []string{}
			ids := ctx.Request.FormValue("ids")
			noticeIds := strings.Split(ids, ",")
			ifEmpty := true
			for _, v := range noticeIds {
				if v != "" {
					ifEmpty = false
					vars := &url.Values{}
					vars.Add("switchID", v)
					ctx1 := &context.Context{}
					err, res := fusion.CallToCenter(vars, common.MyCfg.Api["pushFunctionSwitch"], "GET", ctx1)
					if err != nil {
						failed = append(failed, "<h5>"+v+":"+err.Error()+"</h5>")
					}
					if res != "" {
						db_global := fusion.GetBaseGormDB("db_global")
						db_global.Table("t_function_switch").Where("id = ?", v).Update("updateTime", def.MyTime(time.Now()))
						failed = append(failed, "<h5>"+res+"</h5>")
					}
				}
			}
			if ifEmpty {
				return false, gotext.Get("请选择条目"), ""
			}
			return true, "", failed
		}))
	myOP := fusion.GetServerList()
	formList := tFunctionSwitch.GetForm()
	formList.AddField("Id", "Id", db.Int, form.Number).
		FieldDisplayButCanNotEditWhenUpdate().FieldHide()
	formList.AddField(gotext.Get("开关名字"), "switchName", db.Mediumtext, form.Text).FieldDisplay(func(model types.FieldModel) interface{} {
		return gotext.Get(model.Value)
	})
	formList.AddField("UpdateTime", "updateTime", db.Datetime, form.Datetime).
		FieldValue(time.Now().Format("2006-01-02 15:04:05")).
		FieldDefault(time.Now().Format("2006-01-02 15:04:05")).
		FieldHide()

	formList.AddField(gotext.Get("发往指定服务器"), "gsIds", db.Text, form.SelectBox).
		FieldOptions(myOP).
		FieldHelpMsg(template.HTML(gotext.Get("在此项选择要发往的服务器,若不填则向所有服务器发送")))

	formList.AddField(gotext.Get("不发往指定服务器"), "gsNotIds", db.Text, form.SelectBox).
		FieldOptions(myOP).
		FieldHelpMsg(template.HTML(gotext.Get("在此项选择不发往的服务器")))

	formList.SetTable("t_function_switch").HideContinueNewCheckBox().HideContinueEditCheckBox()

	return tFunctionSwitch
}
