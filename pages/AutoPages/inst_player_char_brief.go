package AutoPages

import (
	"admin/common/def"
	"admin/fusion"
	"html/template"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/GoAdminGroup/go-admin/context"
	"github.com/GoAdminGroup/go-admin/modules/db"
	"github.com/GoAdminGroup/go-admin/plugins/admin/modules/parameter"
	"github.com/GoAdminGroup/go-admin/plugins/admin/modules/table"
	"github.com/GoAdminGroup/go-admin/template/types"
	"github.com/GoAdminGroup/go-admin/template/types/action"
	"github.com/GoAdminGroup/go-admin/template/types/form"
	"github.com/fatih/structs"
	"github.com/leonelquinteros/gotext"
)

func GetPlayerBriefInfo(ctx *context.Context) table.Table {
	PlayerBriefInfo := table.NewDefaultTable(table.Config{
		Driver:     "mysql",
		Connection: "",
		CanAdd:     false,
		Editable:   false,
		Deletable:  false,
		Exportable: false,
		PrimaryKey: table.PrimaryKey{
			Type: db.Int,
			Name: "IpcInstID",
		},
	})
	panelList := PlayerBriefInfo.GetInfo()
	myOP := fusion.GetServerList()
	panelList.AddField(gotext.Get("服务器ID"), "IpcServerID", db.Int).
		FieldFilterable(types.FilterType{FormType: form.SelectSingle}).
		FieldFilterOptions(myOP)
	panelList.AddField(gotext.Get("玩家ID"), "IpcInstID", db.Int)
	panelList.AddField(gotext.Get("账号ID"), "IpcAcctID", db.Int).
		FieldFilterable()
	panelList.AddField(gotext.Get("玩家昵称"), "IpcNickName", db.Text).
		FieldFilterable()
	panelList.AddField(gotext.Get("玩家等级"), "IpcLevel", db.Int)
	panelList.AddField(gotext.Get("角色创建时间"), "IpcCreateTime", db.Int).
		FieldDisplay(func(model types.FieldModel) interface{} {
			cTime, _ := strconv.Atoi(model.Value)
			displayTime := time.Unix(int64(cTime), 0)
			return displayTime.Format("2006-01-02 15:04:05")
		}).
		FieldFilterable(types.FilterType{FormType: form.DatetimeRange})
	panelList.AddField(gotext.Get("角色最后登录时间"), "IpcLastLoginTime", db.Int).
		FieldSortable().FieldDisplay(func(model types.FieldModel) interface{} {
		cTime, _ := strconv.Atoi(model.Value)
		displayTime := time.Unix(int64(cTime), 0)
		return displayTime.Format("2006-01-02 15:04:05")
	}).
		FieldFilterable(types.FilterType{FormType: form.DatetimeRange})

	panelList.SetTable("inst_player_char").SetTitle(gotext.Get("查找玩家")).
		SetDescription(gotext.Get("必须选择服务器后方可查询"))

	panelList.HideNewButton().HideEditButton().HideDeleteButton().HideDetailButton().HideRowSelector().HideExportButton()
	panelList.AddActionButton(template.HTML(gotext.Get("玩家详情")), action.Jump("/admin/userMgr/PlayerDetail?IpcInstID={{(index .Value \"IpcInstID\").Value}}&"+
		"IpcServerID={{(index .Value \"IpcServerID\").Value}}&IpcAcctID={{(index .Value \"IpcAcctID\").Value}}"))

	panelList.SetGetDataFn(
		func(param parameter.Parameters) (data []map[string]interface{}, size int) {
			tStart := time.Now()
			defer func() {
				log.Printf("[PlayerSearch] total elapsed=%v server=%s player=%s acct=%s nick=%s",
					time.Since(tStart), ctx.Request.FormValue("IpcServerID"),
					ctx.Request.FormValue("IpcInstID"), ctx.Request.FormValue("IpcAcctID"),
					ctx.Request.FormValue("IpcNickName"))
			}()

			ServerID := ctx.Request.FormValue("IpcServerID")
			AccountID := ctx.Request.FormValue("IpcAcctID")
			PlayerID := ctx.Request.FormValue("IpcInstID")
			NickName := ctx.Request.FormValue("IpcNickName")
			CreateTimeStart := ctx.Request.FormValue("IpcCreateTime_start__goadmin")
			CreateTimeEnd := ctx.Request.FormValue("IpcCreateTime_end__goadmin")
			LastLoginTimeStart := ctx.Request.FormValue("IpcLastLoginTime_start__goadmin")
			LastLoginTimeEnd := ctx.Request.FormValue("IpcLastLoginTime_end__goadmin")

			temp1, _ := time.ParseInLocation("2006-01-02 15:04:05", CreateTimeStart, time.Local)
			CreateTimeStartUnix := temp1.Unix()
			temp1, _ = time.ParseInLocation("2006-01-02 15:04:05", CreateTimeEnd, time.Local)
			CreateTimeEndUnix := temp1.Unix()
			temp1, _ = time.ParseInLocation("2006-01-02 15:04:05", LastLoginTimeStart, time.Local)
			LastLoginTimeStartUnix := temp1.Unix()
			temp1, _ = time.ParseInLocation("2006-01-02 15:04:05", LastLoginTimeEnd, time.Local)
			LastLoginTimeEndUnix := temp1.Unix()
			log.Print(ServerID)
			ServerIDint, _ := strconv.Atoi(ServerID)
			log.Print(ServerIDint)

			tDB := time.Now()
			gormDB := fusion.GetServerGormDB("db_char", int64(ServerIDint))
			log.Printf("[PlayerSearch] GetServerGormDB elapsed=%v serverID=%d", time.Since(tDB), ServerIDint)
			if gormDB == nil {
				return nil, 0
			}
			baseSQL := "select ipcInstID,ipcLevel,ipcAcctID,ipcNickName,ipcServerID,ipcCreateTime," +
				"ipcLastLoginTime from inst_player_char"
			countSQL := "select count(*) from inst_player_char"
			conditions := []string{}
			args := []interface{}{}
			if AccountID != "" {
				conditions = append(conditions, "ipcAcctID = ?")
				args = append(args, AccountID)
			}
			if PlayerID != "" {
				conditions = append(conditions, "ipcInstID = ?")
				args = append(args, PlayerID)
			}
			if NickName != "" {
				conditions = append(conditions, "ipcNickName LIKE ?")
				args = append(args, "%"+NickName+"%")
			}
			if CreateTimeStart != "" {
				conditions = append(conditions, "ipcCreateTime >= ?")
				args = append(args, CreateTimeStartUnix)
			}
			if CreateTimeEnd != "" {
				conditions = append(conditions, "ipcCreateTime <= ?")
				args = append(args, CreateTimeEndUnix)
			}
			if LastLoginTimeStart != "" {
				conditions = append(conditions, "ipcLastLoginTime >= ?")
				args = append(args, LastLoginTimeStartUnix)
			}
			if LastLoginTimeEnd != "" {
				conditions = append(conditions, "ipcLastLoginTime <= ?")
				args = append(args, LastLoginTimeEndUnix)
			}
			whereClause := " where ipcDeleteTime = 0"
			if len(conditions) > 0 {
				whereClause = " where " + strings.Join(conditions, " and ") + " and ipcDeleteTime = 0"
			}
			sql1 := baseSQL + whereClause + " LIMIT ? OFFSET ?"
			queryArgs := make([]interface{}, len(args), len(args)+2)
			copy(queryArgs, args)
			queryArgs = append(queryArgs, param.PageSizeInt, (param.PageInt-1)*param.PageSizeInt)
			playerInfo := []def.PlayerBriefInfo{}
			tData := time.Now()
			gormDB.Raw(sql1, queryArgs...).Scan(&playerInfo)
			log.Printf("[PlayerSearch] dataQuery elapsed=%v rows=%d", time.Since(tData), len(playerInfo))
			sql2 := countSQL + whereClause
			countArgs := make([]interface{}, len(args))
			copy(countArgs, args)
			var count int64
			tCount := time.Now()
			gormDB.Raw(sql2, countArgs...).Scan(&count)
			log.Printf("[PlayerSearch] countQuery elapsed=%v count=%d", time.Since(tCount), count)
			temp := []map[string]interface{}{}
			for _, v := range playerInfo {
				m3 := structs.Map(&v)
				temp = append(temp, m3)
			}
			return temp, int(count)
		})
	return PlayerBriefInfo
}
