package AutoPages

import (
	"admin/common/def"
	"admin/fusion"
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

func GetLogPlayerStoneTable(ctx *context.Context) table.Table {

	logPlayerStone := table.NewDefaultTable(table.Config{
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

	info := logPlayerStone.GetInfo()

	info.HideNewButton().HideEditButton().HideDeleteButton().HideDetailButton().HideRowSelector().HideExportButton()

	myOP := fusion.GetServerList()
	if strings.Index(ctx.Request.RequestURI, "/export/") != -1 {
		info.AddField(gotext.Get("服务器名"), "ServerName", db.Varchar)
	}
	info.AddField(gotext.Get("服务器ID"), "IpcServerID", db.Int).
		FieldFilterable(types.FilterType{FormType: form.SelectSingle}).
		FieldFilterOptions(myOP).
		FieldHide()
	info.AddField("Id", "Id", db.Bigint)
	info.AddField(gotext.Get("账号ID"), "AcctId", db.Int).FieldHide()
	info.AddField(gotext.Get("角色ID"), "PlayerId", db.Int).FieldHide()
	info.AddField(gotext.Get("角色昵称"), "PlayerName", db.Varchar).
		FieldFilterable()
	info.AddField(gotext.Get("角色等级"), "PlayerLevel", db.Int).
		FieldFilterable(types.FilterType{FormType: form.NumberRange}).
		FieldSortable()
	info.AddField(gotext.Get("VIP等级"), "PlayerVipLevel", db.Int).
		FieldHide()
	info.AddField("PlayerFightValue", "PlayerFightValue", db.Int).
		FieldHide()
	info.AddField("PlayerRichValue", "PlayerRichValue", db.Double).
		FieldHide()
	info.AddField(gotext.Get("记录类型"), "LogType", db.Int).
		FieldFilterable(types.FilterType{FormType: form.SelectSingle}).
		FieldFilterOptions(def.StoneLog{}.GetOptions()).
		FieldDisplay(func(model types.FieldModel) interface{} {
			temp, _ := strconv.Atoi(model.Value)
			return def.StoneLog{}.GetOpStrings()[temp]
		})
	logType := 999
	temp := ctx.Request.FormValue("LogType")
	if temp != "" {
		logType, _ = strconv.Atoi(temp)
	}
	if logType == def.Inlay || logType == def.Remove || strings.Index(ctx.Request.RequestURI, "/export/") != -1 {
		info.AddField(gotext.Get("孔位"), "HoleID", db.Int)
		info.AddField(gotext.Get("宝石类型ID"), "GemItemTypeID", db.Int)
	}
	if logType == def.Inlay || strings.Index(ctx.Request.RequestURI, "/export/") != -1 {
		info.AddField(gotext.Get("幸运属性"), "GemsLucky", db.Mediumtext)
		info.AddField(gotext.Get("宝石类型ID"), "GemItemTypeID", db.Int)
	}
	if logType == def.Enhancement || strings.Index(ctx.Request.RequestURI, "/export/") != -1 {
		info.AddField(gotext.Get("宝石等级"), "GemLevel", db.Int)
		info.AddField(gotext.Get("宝石类型ID"), "GemItemTypeID", db.Int)
	}
	info.AddField(gotext.Get("记录时间"), "LogTime", db.Datetime).
		FieldFilterable(types.FilterType{FormType: form.DatetimeRange}).
		FieldSortable()

	info.SetTable("log_player_stone").SetTitle(gotext.Get("荧之石日志")).
		SetDescription(gotext.Get("必须选择服务器与记录时间后方可查询"))

	info.SetGetDataFn(
		func(param parameter.Parameters) (data []map[string]interface{}, size int) {
			serverId, _ := strconv.Atoi(ctx.Request.FormValue("IpcServerID"))
			acctId := ctx.Request.FormValue("AcctId")
			playerId := ctx.Request.FormValue("PlayerId")
			playerName := ctx.Request.FormValue("PlayerName")
			playerLevelStart := ctx.Request.FormValue("PlayerLevel_start__goadmin")
			playerLevelEnd := ctx.Request.FormValue("PlayerLevel_end__goadmin")
			logTimeStart := ctx.Request.FormValue("LogTime_start__goadmin")
			logTimeEnd := ctx.Request.FormValue("LogTime_end__goadmin")

			logType := ctx.Request.FormValue("LogType")

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
			if logType != "" {
				where += " logType = " + logType + " and "
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

			serverNames := fusion.GetServerNameStrings()
			stoneLog := []def.StoneLog{}
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
					sql1, sql2 := fusion.CreateUnionTableSqlByDays(logTimeStart, logTimeEnd, def.StoneLog{}.TableName(),
						gormDB1, def.StoneLog{}, where, false)
					gormDB1.Raw(sql1).Scan(&stoneLog)
					for i, _ := range stoneLog {
						stoneLog[i].LogTime = strings.Replace(stoneLog[i].LogTime, "T", " ", -1)
						stoneLog[i].LogTime = strings.Replace(stoneLog[i].LogTime, "Z", " ", -1)
					}
					for _, v := range stoneLog {
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
				sql1, sql2 := fusion.CreateUnionTableSqlByDays(logTimeStart, logTimeEnd, def.StoneLog{}.TableName(),
					gormDB, def.StoneLog{}, where, false)
				if !param.IsAll() {
					if param.SortField != "" {
						sql1 += " ORDER BY " + param.SortField + " " + param.SortType
					}
					sql1 += "  LIMIT " + strconv.Itoa(param.PageSizeInt) + " OFFSET " + strconv.Itoa((param.PageInt-1)*param.PageSizeInt)
				}
				gormDB.Raw(sql1).Scan(&stoneLog)
				for i, _ := range stoneLog {
					stoneLog[i].LogTime = strings.Replace(stoneLog[i].LogTime, "T", " ", -1)
					stoneLog[i].LogTime = strings.Replace(stoneLog[i].LogTime, "Z", " ", -1)
				}
				var tempCount []int64
				gormDB.Raw(sql2).Scan(&tempCount)
				for _, val := range tempCount {
					count += val
				}
				for _, v := range stoneLog {
					m3 := structs.Map(&v)
					temp = append(temp, m3)
				}
			}
			return temp, int(count)
		})

	info.SetExportProcessFn(func(param parameter.Parameters) (types.PanelInfo, error) {
		panelInfo, _ := logPlayerStone.GetData(param.WithIsAll(param.IsAll()))
		return types.PanelInfo{Thead: panelInfo.Thead, InfoList: panelInfo.InfoList}, nil
	})
	return logPlayerStone
}
