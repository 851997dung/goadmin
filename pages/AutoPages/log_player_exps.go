package AutoPages

import (
	"admin/common/def"
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

func GetLogPlayerExpsTable(ctx *context.Context) table.Table {

	logPlayerExps := table.NewDefaultTable(table.Config{
		Driver:     "mysql",
		Connection: "db_log",
		CanAdd:     false,
		Editable:   false,
		Deletable:  false,
		Exportable: true,
		PrimaryKey: table.PrimaryKey{
			Type: db.Int,
			Name: "Id",
		},
	})

	myops := pageMgr.GetFlowTypeOptions()
	flowstr := flowDef.GetFlowTypeStrings()

	info := logPlayerExps.GetInfo()
	info.HideNewButton().HideEditButton().HideDeleteButton().HideDetailButton().HideRowSelector().HideExportButton()

	myOP := fusion.GetServerList()
	if strings.Index(ctx.Request.RequestURI, "/export/") != -1 {
		info.AddField(gotext.Get("服务器名"), "ServerName", db.Varchar)
	}
	info.AddField(gotext.Get("服务器ID"), "IpcServerID", db.Int).
		FieldFilterable(types.FilterType{FormType: form.SelectSingle}).
		FieldFilterOptions(myOP).
		FieldHide()
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
	info.AddField(gotext.Get("获取经验"), "Exp", db.Bigint).
		FieldSortable()
	info.AddField(gotext.Get("次数"), "Count", db.Int)
	info.AddField(gotext.Get("操作来源"), "FlowType", db.Int).
		FieldFilterable(types.FilterType{FormType: form.SelectSingle}).
		FieldFilterOptions(myops).
		FieldDisplay(func(model types.FieldModel) interface{} {
			temp, _ := strconv.Atoi(model.Value)
			if flowstr[temp] == "" {
				return model.Value
			}
			return flowstr[temp]
		})
	info.AddField("FlowParams", "FlowParams", db.JSON).FieldHide()

	//info.AddField("LogStartTime", "LogStartTime", db.Datetime).
	//	FieldHide()
	//FieldFilterable(types.FilterType{FormType: form.DatetimeRange})
	info.AddField(gotext.Get("记录时间"), "LogTime", db.Datetime).
		FieldFilterable(types.FilterType{FormType: form.DatetimeRange}).
		FieldSortable()

	info.SetTable("log_player_exps").SetTitle(gotext.Get("经验获取日志")).
		SetDescription(gotext.Get("必须选择服务器与记录时间后方可查询"))

	info.SetGetDataFn(
		func(param parameter.Parameters) (data []map[string]interface{}, size int) {
			serverId, _ := strconv.Atoi(ctx.Request.FormValue("IpcServerID"))
			acctId := ctx.Request.FormValue("AcctId")
			playerId := ctx.Request.FormValue("PlayerId")
			playerName := ctx.Request.FormValue("PlayerName")
			playerLevelStart := ctx.Request.FormValue("PlayerLevel_start__goadmin")
			playerLevelEnd := ctx.Request.FormValue("PlayerLevel_end__goadmin")

			flowType := ctx.Request.FormValue("FlowType")
			//logStartTimeStart := ctx.Request.FormValue("LogStartTime_start__goadmin")
			//logStartTimeEnd := ctx.Request.FormValue("LogStartTime_end__goadmin")
			logStopTimeStart := ctx.Request.FormValue("LogTime_start__goadmin")
			logStopTimeEnd := ctx.Request.FormValue("LogTime_end__goadmin")
			if logStopTimeStart == "" || logStopTimeEnd == "" {
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
			if logStopTimeStart != "" {
				where += " logStopTime >= '" + logStopTimeStart + "' and "
			}
			if logStopTimeEnd != "" {
				where += " logStopTime <= '" + logStopTimeEnd + "' and "
			}
			//if logStartTimeStart != "" {
			//	where += " logStartTime >= '" + logStartTimeStart + "' and "
			//}
			//if logStartTimeEnd != "" {
			//	where += " logStartTime <= '" + logStartTimeEnd + "' and "
			//}
			if where != "" {
				where = " where " + where
				where = strings.TrimSuffix(where, " and ")
			} else {
				return nil, 0
			}

			expLog := []def.ExpLog{}

			serverNames := fusion.GetServerNameStrings()
			var count int64
			temp := []map[string]interface{}{}
			serverIDS := ctx.Request.Form["ServerIDs"]
			if strings.Index(ctx.Request.RequestURI, "/export/") != -1 && len(serverIDS) != 0 {
				for _, v := range serverIDS {
					var tempCount []int64
					tempId, _ := strconv.Atoi(v)
					gormDB1 := fusion.GetServerGormDB("db_log", int64(tempId))
					if gormDB1 == nil {
						continue
					}
					sql1, sql2 := fusion.CreateUnionTableSqlByDays(logStopTimeStart, logStopTimeEnd, def.ExpLog{}.TableName(),
						gormDB1, def.ExpLog{}, where, false)
					gormDB1.Raw(sql1).Scan(&expLog)
					for i, _ := range expLog {
						expLog[i].LogTime = strings.Replace(expLog[i].LogTime, "T", " ", -1)
						expLog[i].LogTime = strings.Replace(expLog[i].LogTime, "Z", " ", -1)
					}
					for _, v := range expLog {
						m3 := structs.Map(&v)
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
				sql1, sql2 := fusion.CreateUnionTableSqlByDays(logStopTimeStart, logStopTimeEnd, def.ExpLog{}.TableName(),
					gormDB, def.ExpLog{}, where, false)
				if !param.IsAll() {
					if param.SortField != "" {
						if param.SortField == "LogTime" {
							param.SortField = "logStopTime"
						}
						sql1 += " ORDER BY " + param.SortField + " " + param.SortType
					}
					sql1 += "  LIMIT " + strconv.Itoa(param.PageSizeInt) + " OFFSET " + strconv.Itoa((param.PageInt-1)*param.PageSizeInt)
				}
				gormDB.Raw(sql1).Scan(&expLog)
				for i, _ := range expLog {
					expLog[i].LogTime = strings.Replace(expLog[i].LogTime, "T", " ", -1)
					expLog[i].LogTime = strings.Replace(expLog[i].LogTime, "Z", " ", -1)
				}
				var tempCount []int64
				gormDB.Raw(sql2).Scan(&tempCount)
				for _, val := range tempCount {
					count += val
				}
				for _, v := range expLog {
					m3 := structs.Map(&v)
					temp = append(temp, m3)
				}
			}
			return temp, int(count)
		})

	info.SetExportProcessFn(func(param parameter.Parameters) (types.PanelInfo, error) {
		panelInfo, _ := logPlayerExps.GetData(param.WithIsAll(param.IsAll()))
		return types.PanelInfo{Thead: panelInfo.Thead, InfoList: panelInfo.InfoList}, nil
	})
	return logPlayerExps
}
