package AutoPages

import (
	"admin/fusion"
	"strconv"

	"github.com/GoAdminGroup/go-admin/context"
	"github.com/GoAdminGroup/go-admin/modules/db"
	"github.com/GoAdminGroup/go-admin/plugins/admin/modules/parameter"
	"github.com/GoAdminGroup/go-admin/plugins/admin/modules/table"
	"github.com/GoAdminGroup/go-admin/template/types"
	"github.com/GoAdminGroup/go-admin/template/types/form"
	"github.com/leonelquinteros/gotext"
)

func GetAccountDataCountTable(ctx *context.Context) table.Table {

	accountDataCount := table.NewDefaultTable(table.DefaultConfigWithDriverAndConnection("mysql", "go_manager"))

	info := accountDataCount.GetInfo().HideEditButton().HideNewButton().HideRowSelector().HideDetailButton().HideDeleteButton().HideRowSelector()
	serverList := fusion.GetServerList()

	info.AddField("Id", "id", db.Int).FieldHide()
	info.AddField(gotext.Get("服务器ID"), "serverId", db.Int).
		FieldFilterable(types.FilterType{FormType: form.SelectBox}).
		FieldFilterOptions(serverList)
	info.AddField(gotext.Get("离线七天以上账号数量"), "offline7Day", db.Int)
	info.AddField(gotext.Get("总在线超过48小时的账号数量"), "online48Hour", db.Int)
	info.AddField(gotext.Get("总在线超过48小时的付费账号数量"), "online48HourPaid", db.Int)
	info.AddField(gotext.Get("付费账号/账号"), "paidPer", db.Float).
		FieldDisplay(func(model types.FieldModel) interface{} {
			temp, _ := strconv.ParseFloat(model.Value, 32)
			return strconv.FormatFloat(temp*100, 'f', 2, 32) + "%"
		})
	info.AddField(gotext.Get("日期"), "logTime", db.Datetime).
		FieldFilterable(types.FilterType{FormType: form.DateRange})

	info.SetTable("account_data_count").SetTitle(gotext.Get("用户数据汇总")).SetPageSizeList([]int{50, 100, 9999})

	//formList := accountDataCount.GetForm()
	//formList.AddField("Id", "id", db.Int, form.Default)
	//formList.AddField("ServerId", "serverId", db.Int, form.Number)
	//formList.AddField("Offline7Day", "offline7Day", db.Int, form.Number)
	//formList.AddField("Online48Hour", "online48Hour", db.Int, form.Number)
	//formList.AddField("Online48HourPaid", "online48HourPaid", db.Int, form.Number)
	//formList.AddField("PaidPer", "paidPer", db.Float, form.Text)
	//formList.AddField("LogTime", "logTime", db.Datetime, form.Datetime)
	//
	//formList.SetTable("account_data_count").SetTitle("AccountDataCount").SetDescription("AccountDataCount")

	info.SetExportProcessFn(func(param parameter.Parameters) (types.PanelInfo, error) {
		panelInfo, _ := accountDataCount.GetData(param.WithIsAll(param.IsAll()))
		return types.PanelInfo{Thead: panelInfo.Thead, InfoList: panelInfo.InfoList}, nil
	})
	return accountDataCount
}
