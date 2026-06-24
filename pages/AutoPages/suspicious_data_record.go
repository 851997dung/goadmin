package AutoPages

import (
	"admin/common/def"
	"admin/common/def/currencyDef"
	"admin/fusion"
	"strconv"

	"github.com/GoAdminGroup/go-admin/context"
	"github.com/GoAdminGroup/go-admin/modules/db"
	"github.com/GoAdminGroup/go-admin/plugins/admin/modules/table"
	"github.com/GoAdminGroup/go-admin/template"
	"github.com/GoAdminGroup/go-admin/template/types"
	"github.com/GoAdminGroup/go-admin/template/types/action"
	"github.com/GoAdminGroup/go-admin/template/types/form"
	"github.com/leonelquinteros/gotext"
)

func GetSuspiciousDataRecordTable(ctx *context.Context) table.Table {

	suspiciousDataRecord := table.NewDefaultTable(table.Config{
		Driver:     "mysql",
		Connection: "go_manager",
		CanAdd:     true,
		Editable:   true,
		Deletable:  true,
		Exportable: true,
		PrimaryKey: table.PrimaryKey{
			Type: db.Int,
			Name: "Id",
		},
	})

	info := suspiciousDataRecord.GetInfo()

	serverStr := fusion.GetServerNameStrings()
	cfgOps := def.GetGlobalCfgOps()
	cfgStrS := def.GetGlobalCfgStr()
	currencyStr := currencyDef.GetCurrencyStrings()
	itemNamStr := fusion.GetItemNameStrings()

	cheatTypeS := def.GetCheatTypeStr()
	//cheatOps := def.GetCheatTypeOps()

	info.AddField("Id", "Id", db.Int)
	info.AddField(gotext.Get("服务器ID"), "gsId", db.Int).
		FieldDisplay(func(model types.FieldModel) interface{} {
			temp, _ := strconv.Atoi(model.Value)
			if serverStr[temp] != "" {
				return serverStr[temp]
			}
			return model.Value
		}).FieldFilterable(types.FilterType{})
	info.AddField(gotext.Get("账号ID"), "accountId", db.Int).FieldFilterable(types.FilterType{})
	info.AddField(gotext.Get("角色ID"), "playerId", db.Int).FieldFilterable(types.FilterType{})
	info.AddField(gotext.Get("角色名"), "playerName", db.Varchar).FieldFilterable(types.FilterType{})
	info.AddField(gotext.Get("记录类型"), "type", db.Int).
		FieldDisplay(func(model types.FieldModel) interface{} {
			temp, _ := strconv.Atoi(model.Value)
			return cfgStrS[temp]
		}).FieldFilterable(types.FilterType{FormType: form.SelectSingle}).FieldFilterOptions(cfgOps)
	info.AddField("Key", "key", db.Int).
		FieldDisplay(func(model types.FieldModel) interface{} {
			Type, _ := strconv.Atoi(fusion.Strval(model.Row["type"]))
			tempStr := make(map[int]string)
			switch Type {
			case 0:
			case 2:
				tempStr = currencyStr
				break
			case 1:
			case 3:
				tempStr = itemNamStr
				break
			case 5:
				tempStr = cheatTypeS
				break
			}
			temp, _ := strconv.Atoi(model.Value)
			if tempStr[temp] != "" {
				return tempStr[temp] + "(" + model.Value + ")"
			}
			return model.Value
		}).FieldFilterable(types.FilterType{})
	info.AddField("Value", "value", db.Int)
	info.AddField(gotext.Get("记录时间"), "logTime", db.Datetime).
		FieldFilterable(types.FilterType{FormType: form.DatetimeRange})

	info.SetTable("suspicious_data_record").
		SetTitle(gotext.Get("监测日志")).
		HideNewButton().
		HideDetailButton().
		HideEditButton().
		HideRowSelector()

	info.AddActionButton(template.HTML(gotext.Get("玩家详情")), action.Jump("/admin/userMgr/PlayerDetail?IpcInstID={{(index .Value \"PlayerId\").Value}}&"+
		"IpcServerID={{(index .Value \"gsId\").Value}}&IpcAcctID={{(index .Value \"accountId\").Value}}"))

	formList := suspiciousDataRecord.GetForm()
	formList.AddField("Id", "Id", db.Int, form.Number).
		FieldDisableWhenCreate()
	formList.AddField("GsId", "gsId", db.Int, form.Number)
	formList.AddField("AccountId", "accountId", db.Int, form.Number)
	formList.AddField("PlayerId", "playerId", db.Int, form.Number)
	formList.AddField("PlayerName", "playerName", db.Varchar, form.Text)
	formList.AddField("Type", "type", db.Int, form.Number)
	formList.AddField("Key", "key", db.Int, form.Number)
	formList.AddField("Value", "value", db.Int, form.Number)
	formList.AddField("LogTime", "logTime", db.Datetime, form.Datetime)

	formList.SetTable("suspicious_data_record").SetTitle(gotext.Get("监测日志"))

	return suspiciousDataRecord
}
