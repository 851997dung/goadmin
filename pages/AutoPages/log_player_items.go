package AutoPages

import (
	"admin/common/def"
	"admin/common/def/flowDef"
	"admin/fusion"
	"admin/mgr/pageMgr"
	"strconv"
	"strings"

	"github.com/GoAdminGroup/go-admin/context"
	"github.com/GoAdminGroup/go-admin/modules/db"
	"github.com/GoAdminGroup/go-admin/plugins/admin/modules/parameter"
	"github.com/GoAdminGroup/go-admin/plugins/admin/modules/table"
	"github.com/GoAdminGroup/go-admin/template/types"
	"github.com/GoAdminGroup/go-admin/template/types/form"
	"github.com/fatih/structs"
	"github.com/leonelquinteros/gotext"
)

func GetLogPlayerItemsTable(ctx *context.Context) table.Table {

	logPlayerItems := table.NewDefaultTable(table.Config{
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

	info := logPlayerItems.GetInfo()

	info.HideNewButton().HideEditButton().HideDeleteButton().HideDetailButton().HideRowSelector().HideExportButton()
	myOP := fusion.GetServerList()

	if strings.Index(ctx.Request.RequestURI, "/export/") != -1 {
		info.AddField(gotext.Get("服务器名"), "ServerName", db.Varchar)
	}
	info.AddField(gotext.Get("服务器ID"), "IpcServerID", db.Int).
		FieldFilterable(types.FilterType{FormType: form.SelectSingle}).
		FieldFilterOptions(myOP)
	info.AddField("Id", "Id", db.Bigint).
		FieldHide()
	info.AddField(gotext.Get("账号ID"), "AcctId", db.Int).FieldHide()
	info.AddField(gotext.Get("角色ID"), "PlayerId", db.Int).FieldHide().
		FieldFilterable()
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
	info.AddField("itemSlot", "ItemSlot", db.Int).
		FieldHide()
	info.AddField("ItemGuid", "ItemGuid", db.Int).
		FieldHide()

	info.AddField(gotext.Get("道具唯一ID"), "ItemConstGuid", db.Int).
		FieldFilterable()
	info.AddField(gotext.Get("道具装载位置"), "ItemSlotType", db.Int).
		FieldFilterable()
	info.AddField(gotext.Get("道具ID"), "ItemTypeId", db.Int).
		FieldFilterable()
	info.AddField(gotext.Get("道具名称"), "ItemName", db.Varchar).
		FieldFilterable()
	info.AddField(gotext.Get("道具数量"), "ItemNum", db.Int).
		FieldSortable()
	info.AddField(gotext.Get("当前道具数量"), "ItemCurNum", db.Int)
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
	info.AddField("FlowParams", "FlowParams", db.JSON).
		FieldDisplay(func(model types.FieldModel) interface{} {
			flowType, _ := strconv.Atoi(fusion.Strval(model.Row["FlowType"]))
			IpcServerID, _ := strconv.Atoi(fusion.Strval(model.Row["IpcServerID"]))
			return fusion.ParseFlowParams4FlowType(uint32(flowType), model.Value, int64(IpcServerID))
		})

	info.AddField(gotext.Get("记录类型"), "LogType", db.Int).
		FieldFilterable(types.FilterType{FormType: form.SelectSingle}).
		FieldFilterOptions(def.ItemLog{}.GetOptions()).
		FieldDisplay(func(model types.FieldModel) interface{} {
			temp, _ := strconv.Atoi(model.Value)
			var ops = def.ItemLog{}.GetOpStrings()
			return ops[temp]
		})
	info.AddField(gotext.Get("记录时间"), "LogTime", db.Datetime).
		FieldFilterable(types.FilterType{FormType: form.DatetimeRange}).
		FieldSortable()
	info.SetTable("log_player_items").SetTitle(gotext.Get("道具日志")).
		SetDescription(gotext.Get("必须选择服务器与记录时间后方可查询"))

	info.SetGetDataFn(
		func(param parameter.Parameters) (data []map[string]interface{}, size int) {
			serverId, _ := strconv.Atoi(ctx.Request.FormValue("IpcServerID"))
			acctId := ctx.Request.FormValue("AcctId")
			playerId := ctx.Request.FormValue("PlayerId")
			playerName := ctx.Request.FormValue("PlayerName")
			playerLevelStart := ctx.Request.FormValue("PlayerLevel_start__goadmin")
			playerLevelEnd := ctx.Request.FormValue("PlayerLevel_end__goadmin")

			itemConstGuid := ctx.Request.FormValue("ItemConstGuid")
			itemSlotType := ctx.Request.FormValue("ItemSlotType")
			itemTypeId := ctx.Request.FormValue("ItemTypeId")
			itemName := ctx.Request.FormValue("ItemName")
			logType := ctx.Request.FormValue("LogType")

			flowType := ctx.Request.FormValue("FlowType")
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
			if itemSlotType != "" {
				where += " itemSlotType = " + itemSlotType + " and "
			}
			if itemTypeId != "" {
				where += " itemTypeId = " + itemTypeId + " and "
			}
			if itemName != "" {
				where += " itemName like " + "'%" + itemName + "%'" + " and "
			}
			if logType != "" {
				where += " logType = " + logType + " and "
			}
			if logTimeStart != "" {
				where += " logTime >= '" + logTimeStart + "' and "
			}
			if logTimeEnd != "" {
				where += " logTime <= '" + logTimeEnd + "' and "
			}
			if itemConstGuid != "" {
				where += " itemConstGuid = " + itemConstGuid + " and "
			}

			if where != "" {
				where = " where " + where
				where = strings.TrimSuffix(where, " and ")
			} else {
				return nil, 0
			}

			serverNames := fusion.GetServerNameStrings()
			itemLog := []def.ItemLog{}
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
					sql1, sql2 := fusion.CreateUnionTableSqlByDays(logTimeStart, logTimeEnd, def.ItemLog{}.TableName(),
						gormDB1, def.ItemLog{}, where, false)
					gormDB1.Raw(sql1).Scan(&itemLog)
					for i, _ := range itemLog {
						itemLog[i].LogTime = strings.Replace(itemLog[i].LogTime, "T", " ", -1)
						itemLog[i].LogTime = strings.Replace(itemLog[i].LogTime, "Z", " ", -1)
					}
					for _, v := range itemLog {
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
				sql1, sql2 := fusion.CreateUnionTableSqlByDays(logTimeStart, logTimeEnd, def.ItemLog{}.TableName(),
					gormDB, def.ItemLog{}, where, false)
				if !param.IsAll() {
					if param.SortField != "" {
						sql1 += " ORDER BY " + param.SortField + " " + param.SortType
					}
					sql1 += "  LIMIT " + strconv.Itoa(param.PageSizeInt) + " OFFSET " + strconv.Itoa((param.PageInt-1)*param.PageSizeInt)
				}
				gormDB.Raw(sql1).Scan(&itemLog)
				for i, _ := range itemLog {
					itemLog[i].LogTime = strings.Replace(itemLog[i].LogTime, "T", " ", -1)
					itemLog[i].LogTime = strings.Replace(itemLog[i].LogTime, "Z", " ", -1)
				}
				var tempCount []int64
				gormDB.Raw(sql2).Scan(&tempCount)
				for _, val := range tempCount {
					count += val
				}
				for _, v := range itemLog {
					m3 := structs.Map(&v)
					m3["IpcServerID"] = serverId
					temp = append(temp, m3)
				}
			}
			return temp, int(count)
		})

	info.SetExportProcessFn(func(param parameter.Parameters) (types.PanelInfo, error) {
		panelInfo, _ := logPlayerItems.GetData(param.WithIsAll(param.IsAll()))
		return types.PanelInfo{Thead: panelInfo.Thead, InfoList: panelInfo.InfoList}, nil
	})
	return logPlayerItems
}
