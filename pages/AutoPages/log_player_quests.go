package AutoPages

import (
	"admin/common/def"
	"admin/common/def/questDef"
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

func GetLogPlayerQuestsTable(ctx *context.Context) table.Table {

	logPlayerQuests := table.NewDefaultTable(table.Config{
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

	info := logPlayerQuests.GetInfo()
	info.HideNewButton().HideEditButton().HideDeleteButton().HideDetailButton().HideRowSelector().HideExportButton()

	myOP := fusion.GetServerList()
	questOp := questDef.GetQuestOptions()
	questSubOp := questDef.GetQuestSubOptions()
	whenTypeOp := questDef.GetQuestWhenTypeOptions()

	questStr := questDef.GetQuestString()
	questSubStr := questDef.GetQuestSubString()
	whenTypeStr := questDef.GetQuestWhenTypeString()

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
	info.AddField("PlayerFightValue", "PlayerFightValue", db.Bigint).
		FieldHide()
	info.AddField("PlayerRichValue", "PlayerRichValue", db.Double).
		FieldHide()
	info.AddField(gotext.Get("任务ID"), "QuestTypeID", db.Int).
		FieldFilterable()
	info.AddField(gotext.Get("任务名"), "QuestName", db.Varchar).
		FieldFilterable()
	info.AddField(gotext.Get("任务类型"), "QuestClass", db.Int).
		FieldFilterable(types.FilterType{FormType: form.SelectSingle}).
		FieldFilterOptions(questOp).
		FieldDisplay(func(model types.FieldModel) interface{} {
			temp, _ := strconv.Atoi(model.Value)
			if questStr[temp] != "" {
				return questStr[temp]
			}
			return model.Value
		})
	info.AddField(gotext.Get("任务子类型"), "QuestSubClass", db.Int).
		FieldFilterable(types.FilterType{FormType: form.SelectSingle}).
		FieldFilterOptions(questSubOp).
		FieldDisplay(func(model types.FieldModel) interface{} {
			temp, _ := strconv.Atoi(model.Value)
			if questStr[temp] != "" {
				return questSubStr[temp]
			}
			return model.Value
		})
	info.AddField(gotext.Get("任务进度"), "WhenType", db.Int).
		FieldFilterable(types.FilterType{FormType: form.SelectSingle}).
		FieldFilterOptions(whenTypeOp).
		FieldDisplay(func(model types.FieldModel) interface{} {
			temp, _ := strconv.Atoi(model.Value)
			if questStr[temp] != "" {
				return whenTypeStr[temp]
			}
			return model.Value
		})
	info.AddField(gotext.Get("记录时间"), "LogTime", db.Datetime).
		FieldFilterable(types.FilterType{FormType: form.DatetimeRange}).
		FieldSortable()

	info.SetTable("log_player_quests").SetTitle(gotext.Get("任务日志")).
		SetDescription(gotext.Get("必须选择服务器与记录时间后方可查询"))

	info.SetGetDataFn(
		func(param parameter.Parameters) (data []map[string]interface{}, size int) {
			serverId, _ := strconv.Atoi(ctx.Request.FormValue("IpcServerID"))

			acctId := ctx.Request.FormValue("AcctId")
			playerId := ctx.Request.FormValue("PlayerId")
			playerName := ctx.Request.FormValue("PlayerName")

			questTypeID := ctx.Request.FormValue("QuestTypeID")
			QuestName := ctx.Request.FormValue("QuestName")
			QuestClass := ctx.Request.FormValue("QuestClass")
			QuestSubClass := ctx.Request.FormValue("QuestSubClass")
			WhenType := ctx.Request.FormValue("WhenType")

			playerLevelStart := ctx.Request.FormValue("PlayerLevel_start__goadmin")
			playerLevelEnd := ctx.Request.FormValue("PlayerLevel_end__goadmin")
			logTimeStart := ctx.Request.FormValue("LogTime_start__goadmin")
			logTimeEnd := ctx.Request.FormValue("LogTime_end__goadmin")

			fields := fusion.GetFields2Struct(def.QuestLog{})
			sql1 := "select " + fields + " from " + def.QuestLog{}.TableName()
			sql2 := "select count(*) from " + def.QuestLog{}.TableName()
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
			if questTypeID != "" {
				where += " questTypeID = '" + questTypeID + "' and "
			}
			if QuestName != "" {
				where += " questName like " + "'%" + QuestName + "%'" + " and "
			}
			if QuestClass != "" {
				where += " questClass = " + QuestClass + " and "
			}
			if QuestSubClass != "" {
				where += " questSubClass = " + QuestSubClass + " and "
			}
			if WhenType != "" {
				where += " whenType = " + WhenType + " and "
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

			sql1 = sql1 + where
			if !param.IsAll() {
				if param.SortField != "" {
					sql1 += " ORDER BY " + param.SortField + " " + param.SortType
				}
				sql1 += "  LIMIT " + strconv.Itoa(param.PageSizeInt) + " OFFSET " + strconv.Itoa((param.PageInt-1)*param.PageSizeInt)
			}
			sql2 = sql2 + where

			questLog := []def.QuestLog{}
			var count int64
			temp := []map[string]interface{}{}

			serverNames := fusion.GetServerNameStrings()
			serverIDS := ctx.Request.Form["ServerIDs"]
			if strings.Index(ctx.Request.RequestURI, "/export/") != -1 && len(serverIDS) != 0 {
				for _, v := range serverIDS {
					var tempCount int64
					tempId, _ := strconv.Atoi(v)
					gormDB1 := fusion.GetServerGormDB("db_log", int64(tempId))
					if gormDB1 == nil {
						continue
					}
					gormDB1.Raw(sql1).Scan(&questLog)
					for i, _ := range questLog {
						questLog[i].LogTime = strings.Replace(questLog[i].LogTime, "T", " ", -1)
						questLog[i].LogTime = strings.Replace(questLog[i].LogTime, "Z", " ", -1)
					}
					for _, v := range questLog {
						m3 := structs.Map(&v)
						m3["ServerName"] = serverNames[tempId]
						temp = append(temp, m3)
					}
					gormDB1.Raw(sql2).Scan(&tempCount)
					count += tempCount
				}
			} else {
				gormDB := fusion.GetServerGormDB("db_log", int64(serverId))
				if gormDB == nil {
					return nil, 0
				}
				gormDB.Raw(sql1).Scan(&questLog)
				for i, _ := range questLog {
					questLog[i].LogTime = strings.Replace(questLog[i].LogTime, "T", " ", -1)
					questLog[i].LogTime = strings.Replace(questLog[i].LogTime, "Z", " ", -1)
				}
				gormDB.Raw(sql2).Scan(&count)
				for _, v := range questLog {
					m3 := structs.Map(&v)
					temp = append(temp, m3)
				}
			}
			return temp, int(count)
		})

	info.SetExportProcessFn(func(param parameter.Parameters) (types.PanelInfo, error) {
		panelInfo, _ := logPlayerQuests.GetData(param.WithIsAll(param.IsAll()))
		return types.PanelInfo{Thead: panelInfo.Thead, InfoList: panelInfo.InfoList}, nil
	})
	return logPlayerQuests
}
