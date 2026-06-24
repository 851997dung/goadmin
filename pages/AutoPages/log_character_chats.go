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

func GetLogCharacterChatsTable(ctx *context.Context) table.Table {

	logCharacterChats := table.NewDefaultTable(table.Config{
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

	info := logCharacterChats.GetInfo()

	info.HideNewButton().HideEditButton().HideDeleteButton().HideDetailButton().HideRowSelector().HideExportButton()
	myOP := fusion.GetServerList()
	channelOp := def.GetChannelOptions()
	if strings.Index(ctx.Request.RequestURI, "/export/") != -1 {
		info.AddField(gotext.Get("服务器名"), "ServerName", db.Varchar)
	}
	info.AddField(gotext.Get("服务器ID"), "IpcServerID", db.Int).
		FieldFilterable(types.FilterType{FormType: form.SelectSingle}).
		FieldFilterOptions(myOP).
		FieldHide()
	info.AddField("ID", "Id", db.Bigint)
	info.AddField(gotext.Get("账号ID"), "AcctId", db.Int).FieldHide()
	info.AddField(gotext.Get("角色ID"), "PlayerId", db.Int).FieldHide()
	info.AddField(gotext.Get("角色名"), "PlayerName", db.Varchar).
		FieldFilterable()
	info.AddField(gotext.Get("频道"), "ChannelType", db.Tinyint).
		FieldDisplay(func(model types.FieldModel) interface{} {
			temp, _ := strconv.Atoi(model.Value)
			return def.GetOpStringsChannel()[temp]
		}).
		FieldFilterable(types.FilterType{FormType: form.SelectSingle}).
		FieldFilterOptions(channelOp)
	info.AddField(gotext.Get("接收玩家"), "ToPlayerId", db.Int).
		FieldFilterable()
	info.AddField(gotext.Get("内容"), "StrMsg", db.Text).
		FieldFilterable()
	info.AddField(gotext.Get("记录时间"), "LogTime", db.Datetime).
		FieldFilterable(types.FilterType{FormType: form.DatetimeRange}).
		FieldSortable()
	info.SetExportProcessFn(func(param parameter.Parameters) (types.PanelInfo, error) {
		temp := types.PanelInfo{}

		return temp, nil
	})

	info.SetTable("log_character_chats").SetTitle(gotext.Get("聊天日志")).
		SetDescription(gotext.Get("必须选择服务器与记录时间后方可查询"))

	info.SetGetDataFn(
		func(param parameter.Parameters) (data []map[string]interface{}, size int) {
			serverId, _ := strconv.Atoi(ctx.Request.FormValue("IpcServerID"))
			acctId := ctx.Request.FormValue("AcctId")
			characterId := ctx.Request.FormValue("PlayerId")
			characterName := ctx.Request.FormValue("PlayerName")
			channelType := ctx.Request.FormValue("ChannelType")
			strMsg := ctx.Request.FormValue("StrMsg")
			logTimeStart := ctx.Request.FormValue("LogTime_start__goadmin")
			logTimeEnd := ctx.Request.FormValue("LogTime_end__goadmin")

			fields := fusion.GetFields2Struct(def.ChatLog{})
			sql1 := "select " + fields + " from " + def.ChatLog{}.TableName()
			sql2 := "select count(*) from " + def.ChatLog{}.TableName()
			where := ""
			if acctId != "" {
				where += " acctId = " + acctId + " and "
			}
			if characterId != "" {
				where += " characterId = " + characterId + " and "
			}
			if characterName != "" {
				where += " characterName like " + "'%" + characterName + "%'" + " and "
			}
			if channelType != "" {
				where += " channelType = " + "'" + channelType + "'" + " and "
			}
			if strMsg != "" {
				where += " strMsg like " + "'%" + strMsg + "%'" + " and "
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

			chatLogs := []def.ChatLog{}
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
					gormDB1.Raw(sql1).Scan(&chatLogs)
					for i, _ := range chatLogs {
						chatLogs[i].LogTime = strings.Replace(chatLogs[i].LogTime, "T", " ", -1)
						chatLogs[i].LogTime = strings.Replace(chatLogs[i].LogTime, "Z", " ", -1)
					}
					for _, v := range chatLogs {
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
				gormDB.Raw(sql1).Scan(&chatLogs)
				for i, _ := range chatLogs {
					chatLogs[i].LogTime = strings.Replace(chatLogs[i].LogTime, "T", " ", -1)
					chatLogs[i].LogTime = strings.Replace(chatLogs[i].LogTime, "Z", " ", -1)
				}
				gormDB.Raw(sql2).Scan(&count)
				for _, v := range chatLogs {
					m3 := structs.Map(&v)
					temp = append(temp, m3)
				}
			}
			return temp, int(count)
		})

	info.SetExportProcessFn(func(param parameter.Parameters) (types.PanelInfo, error) {
		panelInfo, _ := logCharacterChats.GetData(param.WithIsAll(true))
		return types.PanelInfo{Thead: panelInfo.Thead, InfoList: panelInfo.InfoList}, nil
	})

	return logCharacterChats
}
