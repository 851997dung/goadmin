package AutoPages

import (
	"admin/common/def"
	"admin/fusion"
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

func GetLogCharacterOnlinesTable(ctx *context.Context) table.Table {

	logCharacterOnlines := table.NewDefaultTable(table.Config{
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

	info := logCharacterOnlines.GetInfo().HideRowSelector()

	myOP := fusion.GetServerList()
	if strings.Index(ctx.Request.RequestURI, "/export/") != -1 {
		info.AddField(gotext.Get("服务器名"), "ServerName", db.Varchar)
	}
	info.AddField(gotext.Get("服务器ID"), "IpcServerID", db.Int).
		FieldFilterable(types.FilterType{FormType: form.SelectSingle}).
		FieldFilterOptions(myOP).
		FieldHide()
	info.AddField("Id", "Id", db.Bigint)
	info.AddField(gotext.Get("账号"), "AcctId", db.Int)
	info.AddField(gotext.Get("角色ID"), "PlayerId", db.Int)
	info.AddField(gotext.Get("角色名"), "PlayerName", db.Varchar).
		FieldFilterable()
	info.AddField(gotext.Get("角色等级"), "PlayerLevel", db.Int).
		FieldFilterable(types.FilterType{FormType: form.NumberRange}).
		FieldSortable()
	info.AddField(gotext.Get("VIP等级"), "PlayerVipLevel", db.Int).
		FieldHide()
	info.AddField("CharacterFightValue", "PlayerFightValue", db.Bigint).
		FieldHide()
	info.AddField("CharacterRichValue", "PlayerRichValue", db.Double).
		FieldHide()
	info.AddField(gotext.Get("上线/下线"), "IsOnline", db.Boolean).
		FieldDisplay(func(model types.FieldModel) interface{} {
			temp, _ := strconv.ParseBool(model.Value)
			if temp {
				return gotext.Get("上线")
			}
			return gotext.Get("下线")
		}).
		FieldFilterable(types.FilterType{FormType: form.SelectSingle}).
		FieldFilterOptions(types.FieldOptions{
			{Value: "0", Text: gotext.Get("上线")},
			{Value: "1", Text: gotext.Get("下线")},
		}).FieldFilterOptionExt(map[string]interface{}{"AllowClear": true})
	info.AddField(gotext.Get("登录IP"), "CharacterLoginIP", db.Varchar)
	info.AddField(gotext.Get("登录时间"), "LogTime", db.Datetime).
		FieldFilterable(types.FilterType{FormType: form.DatetimeRange}).
		FieldSortable()
	info.HideNewButton().HideEditButton().HideDeleteButton().HideDetailButton().HideRowSelector().HideExportButton()
	info.SetTable("log_character_onlines").SetTitle(gotext.Get("登录日志表格")).
		SetDescription(gotext.Get("必须选择服务器与记录时间后方可查询"))

	info.SetGetDataFn(
		func(param parameter.Parameters) (data []map[string]interface{}, size int) {
			serverId, _ := strconv.Atoi(ctx.Request.FormValue("IpcServerID"))
			acctId := ctx.Request.FormValue("AcctId")
			characterId := ctx.Request.FormValue("PlayerId")
			characterName := ctx.Request.FormValue("PlayerName")
			characterLevelStart := ctx.Request.FormValue("PlayerLevel_start__goadmin")
			characterLevelEnd := ctx.Request.FormValue("PlayerLevel_end__goadmin")
			isOnline := ctx.Request.FormValue("IsOnline")
			logTimeStart := ctx.Request.FormValue("LogTime_start__goadmin")
			logTimeEnd := ctx.Request.FormValue("LogTime_end__goadmin")

			fields := fusion.GetFields2Struct(def.OnlinesLog{})
			sql1 := "select " + fields + " from " + def.OnlinesLog{}.TableName()
			sql2 := "select count(*) from " + def.OnlinesLog{}.TableName()
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
			if isOnline != "" {
				where += " isOnline = " + isOnline + " and "
			}
			if characterLevelStart != "" {
				where += " CharacterLevel >= '" + characterLevelStart + "' and "
			}
			if characterLevelEnd != "" {
				where += " CharacterLevel <= '" + characterLevelEnd + "' and "
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
					if param.SortField == "PlayerLevel" {
						param.SortField = "CharacterLevel"
					}
					sql1 += " ORDER BY " + param.SortField + " " + param.SortType
				}
				sql1 += "  LIMIT " + strconv.Itoa(param.PageSizeInt) + " OFFSET " + strconv.Itoa((param.PageInt-1)*param.PageSizeInt)
			}
			sql2 = sql2 + where
			serverNames := fusion.GetServerNameStrings()
			var count int64
			temp := []map[string]interface{}{}
			onlinesLog := []def.OnlinesLog{}
			serverIDS := ctx.Request.Form["ServerIDs"]
			if strings.Index(ctx.Request.RequestURI, "/export/") != -1 && len(serverIDS) != 0 {
				for _, v := range serverIDS {
					var tempCount int64
					tempId, _ := strconv.Atoi(v)
					gormDB1 := fusion.GetServerGormDB("db_log", int64(tempId))
					if gormDB1 == nil {
						continue
					}
					gormDB1.Raw(sql1).Scan(&onlinesLog)
					for i, _ := range onlinesLog {
						onlinesLog[i].LogTime = strings.Replace(onlinesLog[i].LogTime, "T", " ", -1)
						onlinesLog[i].LogTime = strings.Replace(onlinesLog[i].LogTime, "Z", " ", -1)
					}
					for _, v := range onlinesLog {
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
				gormDB.Raw(sql1).Scan(&onlinesLog)
				for i, _ := range onlinesLog {
					onlinesLog[i].LogTime = strings.Replace(onlinesLog[i].LogTime, "T", " ", -1)
					onlinesLog[i].LogTime = strings.Replace(onlinesLog[i].LogTime, "Z", " ", -1)
				}
				gormDB.Raw(sql2).Scan(&count)
				for _, v := range onlinesLog {
					m3 := structs.Map(&v)
					temp = append(temp, m3)
				}
			}
			return temp, int(count)
		})

	info.SetExportProcessFn(func(param parameter.Parameters) (types.PanelInfo, error) {
		panelInfo, _ := logCharacterOnlines.GetData(param.WithIsAll(param.IsAll()))
		return types.PanelInfo{Thead: panelInfo.Thead, InfoList: panelInfo.InfoList}, nil
	})

	return logCharacterOnlines
}
