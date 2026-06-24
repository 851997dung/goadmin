package AutoPages

import (
	"admin/common/def"
	"admin/common/def/currencyDef"
	"admin/common/def/flowDef"
	"admin/fusion"
	"admin/mgr/pageMgr"
	"github.com/GoAdminGroup/go-admin/context"
	"github.com/GoAdminGroup/go-admin/modules/db"
	"github.com/GoAdminGroup/go-admin/plugins/admin/modules/parameter"
	"github.com/GoAdminGroup/go-admin/plugins/admin/modules/table"
	"github.com/GoAdminGroup/go-admin/template/types"
	"github.com/GoAdminGroup/go-admin/template/types/form"
	"github.com/fatih/structs"
	"github.com/leonelquinteros/gotext"
	"strconv"
	"strings"
)

func GetLogPlayerMoneysTable(ctx *context.Context) table.Table {

	logPlayerMoneys := table.NewDefaultTable(table.Config{
		Driver:     "mysql",
		Connection: "db_log",
		CanAdd:     true,
		Editable:   true,
		Deletable:  true,
		Exportable: true,
		PrimaryKey: table.PrimaryKey{
			Type: db.Int,
			Name: "Id",
		},
	})

	info := logPlayerMoneys.GetInfo()
	info.HideNewButton().HideEditButton().HideDeleteButton().HideDetailButton().HideRowSelector().HideExportButton()
	myops := pageMgr.GetFlowTypeOptions()
	myOP := fusion.GetServerList()
	currencyOPs := currencyDef.GetCurrencyOps()
	flowStr := flowDef.GetFlowTypeStrings()
	if strings.Index(ctx.Request.RequestURI, "/export/") != -1 {
		info.AddField(gotext.Get("服务器名"), "ServerName", db.Varchar)
	}
	info.AddField(gotext.Get("服务器ID"), "IpcServerID", db.Int).
		FieldFilterable(types.FilterType{FormType: form.SelectSingle}).
		FieldFilterOptions(myOP)
	info.AddField("Id", "Id", db.Bigint).
		FieldHide()
	info.AddField(gotext.Get("账号ID"), "AcctId", db.Int).FieldHide()
	info.AddField(gotext.Get("角色ID"), "PlayerId", db.Int).FieldHide()
	info.AddField(gotext.Get("角色昵称"), "PlayerName", db.Varchar).
		FieldFilterable()
	info.AddField(gotext.Get("角色等级"), "PlayerLevel", db.Int).
		FieldFilterable(types.FilterType{FormType: form.NumberRange}).
		FieldSortable()
	info.AddField(gotext.Get("VIP等级"), "PlayerVipLevel", db.Int).
		FieldHide()
	info.AddField("PlayerFightValue", "PlayerFightValue", db.Bigint).
		FieldHide()
	info.AddField("PlayerRichValue", "PlayerRichValue", db.Double).
		FieldHide()
	info.AddField(gotext.Get("货币类型"), "MoneyType", db.Int).
		FieldFilterable(types.FilterType{FormType: form.SelectSingle}).
		FieldFilterOptions(currencyOPs).
		FieldDisplay(func(model types.FieldModel) interface{} {
			temp, _ := strconv.Atoi(model.Value)
			if currencyDef.GetCurrencyStrings()[temp] != "" {
				return currencyDef.GetCurrencyStrings()[temp]
			}
			return model.Value
		})
	info.AddField(gotext.Get("货币名称"), "MoneyName", db.Varchar).
		FieldHide()
	info.AddField(gotext.Get("货币值"), "MoneyValue", db.Bigint).
		FieldSortable()
	info.AddField(gotext.Get("当前货币值"), "MoneyCurValue", db.Bigint)
	info.AddField(gotext.Get("收入/支出"), "OpType", db.Bigint).
		FieldFilterable(types.FilterType{FormType: form.SelectSingle}).
		FieldFilterOptions(types.FieldOptions{
			{Value: "0", Text: gotext.Get("支出")},
			{Value: "1", Text: gotext.Get("收入")},
		}).FieldFilterOptionExt(map[string]interface{}{"AllowClear": true}).
		FieldHide()
	info.AddField(gotext.Get("操作来源"), "FlowType", db.Int).
		FieldFilterable(types.FilterType{FormType: form.SelectSingle}).
		FieldFilterOptions(myops).
		FieldDisplay(func(model types.FieldModel) interface{} {
			temp, _ := strconv.Atoi(model.Value)
			if flowStr[temp] == "" {
				return model.Value
			}
			return flowStr[temp]
		})
	info.AddField("FlowParams", "FlowParams", db.JSON).
		FieldDisplay(func(model types.FieldModel) interface{} {
			flowType, _ := strconv.Atoi(fusion.Strval(model.Row["FlowType"]))
			IpcServerID, _ := strconv.Atoi(fusion.Strval(model.Row["IpcServerID"]))
			return fusion.ParseFlowParams4FlowType(uint32(flowType), model.Value, int64(IpcServerID))
		})

	info.AddField(gotext.Get("记录时间"), "LogTime", db.Datetime).
		FieldFilterable(types.FilterType{FormType: form.DatetimeRange}).
		FieldSortable()
	info.AddField(gotext.Get("是否将货币按来源汇总"), "ifCount", db.Int).FieldHide().
		FieldFilterable(types.FilterType{FormType: form.Switch}).FieldFilterOptions(types.FieldOptions{
		{Text: gotext.Get("是"), Value: "true"},
		{Text: gotext.Get("否"), Value: "false"},
	})

	info.SetTable("log_player_moneys").SetTitle(gotext.Get("货币日志")).
		SetDescription(gotext.Get("必须选择服务器与记录时间后方可查询"))

	info.SetGetDataFn(
		func(param parameter.Parameters) (data []map[string]interface{}, size int) {
			serverId, _ := strconv.Atoi(ctx.Request.FormValue("IpcServerID"))
			acctId := ctx.Request.FormValue("AcctId")
			playerId := ctx.Request.FormValue("PlayerId")
			playerName := ctx.Request.FormValue("PlayerName")
			playerLevelStart := ctx.Request.FormValue("PlayerLevel_start__goadmin")
			playerLevelEnd := ctx.Request.FormValue("PlayerLevel_end__goadmin")
			ifCount, _ := strconv.ParseBool(ctx.Request.FormValue("ifCount"))

			moneyType := ctx.Request.FormValue("MoneyType")
			moneyName := ctx.Request.FormValue("MoneyName")
			//moneyValue := ctx.Request.FormValue("MoneyValue")
			flowType := ctx.Request.FormValue("FlowType")
			opType := ctx.Request.FormValue("OpType")

			logTimeStart := ctx.Request.FormValue("LogTime_start__goadmin")
			logTimeEnd := ctx.Request.FormValue("LogTime_end__goadmin")

			if logTimeStart == "" || logTimeEnd == "" {
				return nil, 0
			}

			where := ""
			if acctId != "" {
				where += " acctId = " + acctId + " and "
			}
			if playerId != "" {
				where += " playerId = " + playerId + " and "
			}
			if playerName != "" {
				where += " playerName like " + "'%" + playerName + "%'" + " and "
			}
			if playerLevelStart != "" {
				where += " playerLevel >= '" + playerLevelStart + "' and "
			}
			if playerLevelEnd != "" {
				where += " playerLevel <= '" + playerLevelEnd + "' and "
			}
			if flowType != "" {
				where += " flowType = " + flowType + " and "
			}
			if moneyType != "" {
				where += " moneyType = " + moneyType + " and "
			} else {
				where += " moneyType in (1,3,4) and"
			}
			if moneyName != "" {
				where += " moneyName like " + "'%" + moneyName + "%'" + " and "
			}
			if opType != "" {
				temp, _ := strconv.Atoi(opType)
				if temp == 0 {
					where += " moneyValue < 0 and "
				} else {
					where += " moneyValue >= 0 and "
				}
			}
			if logTimeStart != "" {
				where += " logTime >= '" + logTimeStart + "' and "
			}
			if logTimeEnd != "" {
				where += " logTime <= '" + logTimeEnd + "' and "
			}

			if where != "" {
				where = " where " + where
				where = strings.TrimSuffix(where, " and ")
			} else {
				return nil, 0
			}

			var count int64
			temp := []map[string]interface{}{}
			moneyLog := []def.MoneyLog{}
			serverNames := fusion.GetServerNameStrings()
			serverIDS := ctx.Request.Form["ServerIDs"]
			if strings.Index(ctx.Request.RequestURI, "/export/") != -1 && len(serverIDS) != 0 {
				for _, v := range serverIDS {
					var tempCount []int64
					tempId, _ := strconv.Atoi(v)
					gormDB1 := fusion.GetServerGormDB("db_log", int64(tempId))
					if gormDB1 == nil {
						continue
					}

					sql1, sql2 := fusion.CreateUnionTableSqlByDays(logTimeStart, logTimeEnd, def.MoneyLog{}.TableName(),
						gormDB1, def.MoneyLog{}, where, false)

					gormDB1.Raw(sql1).Scan(&moneyLog)
					for i, _ := range moneyLog {
						moneyLog[i].LogTime = strings.Replace(moneyLog[i].LogTime, "T", " ", -1)
						moneyLog[i].LogTime = strings.Replace(moneyLog[i].LogTime, "Z", " ", -1)
					}
					for _, v := range moneyLog {
						m3 := structs.Map(&v)
						m3["IpcServerID"] = tempId
						m3["ServerName"] = serverNames[tempId]
						temp = append(temp, m3)
					}
					gormDB1.Raw(sql2).Scan(&tempCount)
					for _, val := range tempCount {
						count += val
					}
				}
			} else {
				gormDB := fusion.GetServerGormDB("db_log", int64(serverId))
				if gormDB == nil {
					return nil, 0
				}
				var sql1, sql2 string
				if ifCount {
					sql1, sql2 = fusion.CreateUnionTableSqlByDays4Money(logTimeStart, logTimeEnd, def.MoneyLog{}.TableName(),
						gormDB, def.MoneyLog{}, where, false)
				} else {
					sql1, sql2 = fusion.CreateUnionTableSqlByDays(logTimeStart, logTimeEnd, def.MoneyLog{}.TableName(),
						gormDB, def.MoneyLog{}, where, false)
				}
				if !param.IsAll() {
					if param.SortField != "" {
						sql1 += " ORDER BY " + param.SortField + " " + param.SortType
					}
					sql1 += "  LIMIT " + strconv.Itoa(param.PageSizeInt) + " OFFSET " + strconv.Itoa((param.PageInt-1)*param.PageSizeInt)
				}
				gormDB.Raw(sql1).Scan(&moneyLog)
				for i, _ := range moneyLog {
					moneyLog[i].LogTime = strings.Replace(moneyLog[i].LogTime, "T", " ", -1)
					moneyLog[i].LogTime = strings.Replace(moneyLog[i].LogTime, "Z", " ", -1)
				}
				var tempCount []int64
				gormDB.Raw(sql2).Scan(&tempCount)
				for _, val := range tempCount {
					count += val
				}
				for _, v := range moneyLog {
					m3 := structs.Map(&v)
					m3["IpcServerID"] = serverId
					temp = append(temp, m3)
				}
			}
			return temp, int(count)
		})

	info.SetExportProcessFn(func(param parameter.Parameters) (types.PanelInfo, error) {
		panelInfo, _ := logPlayerMoneys.GetData(param.WithIsAll(param.IsAll()))
		return types.PanelInfo{Thead: panelInfo.Thead, InfoList: panelInfo.InfoList}, nil
	})
	return logPlayerMoneys
}
