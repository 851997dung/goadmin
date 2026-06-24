package AutoPages

import (
	"admin/common/def"
	"admin/fusion"
	"log"
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

func GetLogCharacterLoginTable(ctx *context.Context) table.Table {

	logCharacterlogin := table.NewDefaultTable(table.Config{
		Driver:     "mysql",
		Connection: "db_login_log",
		CanAdd:     false,
		Editable:   false,
		Deletable:  false,
		Exportable: true,
		PrimaryKey: table.PrimaryKey{
			Type: db.Char,
			Name: "DeviceUniqueIdentifier",
		},
	})

	info := logCharacterlogin.GetInfo().HideRowSelector()
	//myOP := fusion.GetServerList()
	//if strings.Index(ctx.Request.RequestURI, "/export/") != -1 {
	//	info.AddField(gotext.Get("服务器名"), "ServerName", db.Varchar)
	//}
	info.AddField(gotext.Get("服务器ID"), "IpcServerID", db.Int).
		FieldHide()
	//info.AddField("Id", "Id", db.Bigint)
	info.AddField(gotext.Get("账号"), "AcctId", db.Int).FieldFilterable()
	info.AddField(gotext.Get("设备码"), "DeviceUniqueIdentifier", db.Char)
	info.AddField(gotext.Get("设备型号"), "DeviceModel", db.Varchar)
	info.AddField(gotext.Get("最后时间"), "LastTime", db.Char)
	info.AddField(gotext.Get("最后ip"), "LastIP", db.Char)
	info.AddField(gotext.Get("当前步骤"), "CurStep", db.Char)
	info.AddField(gotext.Get("最后步骤"), "FinishStep", db.Char)
	info.AddField(gotext.Get("完美运行"), "IsPerfectPlay", db.Int).
		FieldDisplay(func(model types.FieldModel) interface{} {
			temp, _ := strconv.Atoi(model.Value)
			return def.LoginLog{}.GetStatusStrings()[temp]
		})
	info.AddField(gotext.Get("开启时间"), "DayFirstLoginTime", db.Datetime).
		FieldFilterable(types.FilterType{FormType: form.DatetimeRange}).
		FieldSortable()
	info.AddField(gotext.Get("更新"), "UpdateState", db.Int).
		FieldDisplay(func(model types.FieldModel) interface{} {
			temp, _ := strconv.Atoi(model.Value)
			return def.LoginLog{}.GetStatusStrings()[temp]
		})
	info.AddField(gotext.Get("资源加载"), "ResourceLoad", db.Int).
		FieldDisplay(func(model types.FieldModel) interface{} {
			temp, _ := strconv.Atoi(model.Value)
			return def.LoginLog{}.GetStatusStrings()[temp]
		})
	info.AddField(gotext.Get("SDK账号"), "SDKAcct", db.Int).
		FieldDisplay(func(model types.FieldModel) interface{} {
			temp, _ := strconv.Atoi(model.Value)
			return def.LoginLog{}.GetStatusStrings()[temp]
		})
	info.AddField(gotext.Get("列表失败次数"), "GetSerListFaildNum", db.Int)
	info.AddField(gotext.Get("链接服务器次数"), "ConnSerFaildNum", db.Int)
	info.AddField(gotext.Get("进入游戏状态"), "EnterGameStatus", db.Int).
		FieldDisplay(func(model types.FieldModel) interface{} {
			temp, _ := strconv.Atoi(model.Value)
			return def.LoginLog{}.GetStatusStrings()[temp]
		})
	info.AddField(gotext.Get("包状态"), "PackageStatus", db.Char)
	info.AddField(gotext.Get("服务器列表状态"), "ServerListStatus", db.Int).
		FieldDisplay(func(model types.FieldModel) interface{} {
			temp, _ := strconv.Atoi(model.Value)
			return def.LoginLog{}.GetStatusStrings()[temp]
		})
	info.AddField(gotext.Get("服务器连接状态"), "ServerConnectStatus", db.Int).
		FieldDisplay(func(model types.FieldModel) interface{} {
			temp, _ := strconv.Atoi(model.Value)
			return def.LoginLog{}.GetStatusStrings()[temp]
		})
	info.AddField(gotext.Get("sdk登录失败次数"), "SdkLoginFailCount", db.Int)

	info.SetGetDataFn(
		func(param parameter.Parameters) (data []map[string]interface{}, size int) {
			//serverId, _ := strconv.Atoi(ctx.Request.FormValue("IpcServerID"))
			acctId := ctx.Request.FormValue("AcctId")
			logTimeStart := ctx.Request.FormValue("DayFirstLoginTime_start__goadmin")
			logTimeEnd := ctx.Request.FormValue("DayFirstLoginTime_end__goadmin")

			//fields := fusion.GetFields2Struct(def.LoginLog{})
			//sql1 := "select " + fields + " from " + def.LoginLog{}.TableName()
			//sql2 := "select count(*) from " + def.LoginLog{}.TableName()
			where := ""
			if acctId != "" {
				where += " accountId = " + acctId + " and "
			}
			if logTimeStart != "" {
				where += " firstStartTime >= '" + logTimeStart + "' and "
			}
			if logTimeEnd != "" {
				where += " firstStartTime <= '" + logTimeEnd + "' and "
			}
			if where != "" {
				where = " where " + where
				where = strings.TrimSuffix(where, " and ")
			} else {
				return nil, 0
			}
			gormDB := fusion.GetBaseGormDB("db_login_log")
			if gormDB == nil {
				return nil, 0
			}
			sql1, sql2 := fusion.CreateUnionTableSqlByDays(logTimeStart, logTimeEnd, def.LoginLog{}.TableName(), gormDB, def.LoginLog{}, where, false)
			if !param.IsAll() {
				if param.SortField != "" {
					sql1 += " ORDER BY " + param.SortField + " " + param.SortType
				}
				sql1 += "  LIMIT " + strconv.Itoa(param.PageSizeInt) + " OFFSET " + strconv.Itoa((param.PageInt-1)*param.PageSizeInt)
			}
			var count int64
			temp := []map[string]interface{}{}
			LoginLog := []def.LoginLog{}

			log.Printf("%s", sql1)
			gormDB.Raw(sql1).Scan(&LoginLog)
			for i, _ := range LoginLog {
				//log.Printf("%d", LoginLog[i].Id)
				sgss := def.LoginLog{}.GetOpStrings()[LoginLog[i].CurStep]
				if sgss != "" {
					LoginLog[i].CurStep = def.LoginLog{}.GetOpStrings()[LoginLog[i].CurStep]
				}
				sgss2 := def.LoginLog{}.GetOpStrings()[LoginLog[i].FinishStep]
				if sgss2 != "" {
					LoginLog[i].FinishStep = def.LoginLog{}.GetOpStrings()[LoginLog[i].FinishStep]
				}
				LoginLog[i].FinishStep = def.LoginLog{}.GetOpStrings()[LoginLog[i].FinishStep]
				LoginLog[i].DayFirstLoginTime = strings.Replace(LoginLog[i].DayFirstLoginTime, "T", " ", -1)
				LoginLog[i].DayFirstLoginTime = strings.Replace(LoginLog[i].DayFirstLoginTime, "Z", " ", -1)
			}
			gormDB.Raw(sql2).Scan(&count)
			for _, v := range LoginLog {
				m3 := structs.Map(&v)
				temp = append(temp, m3)
			}

			return temp, int(count)
		})

	info.SetExportProcessFn(func(param parameter.Parameters) (types.PanelInfo, error) {
		panelInfo, _ := logCharacterlogin.GetData(param.WithIsAll(param.IsAll()))
		return types.PanelInfo{Thead: panelInfo.Thead, InfoList: panelInfo.InfoList}, nil
	})

	return logCharacterlogin
}
