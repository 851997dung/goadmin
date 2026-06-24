package AutoPages

import (
	"admin/common"
	"admin/common/def"
	"admin/fusion"
	"admin/hotfix"
	"errors"
	"fmt"
	"html/template"
	"log"
	"math"
	"net/url"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/GoAdminGroup/go-admin/context"
	"github.com/GoAdminGroup/go-admin/modules/db"
	form2 "github.com/GoAdminGroup/go-admin/plugins/admin/modules/form"
	"github.com/GoAdminGroup/go-admin/plugins/admin/modules/parameter"
	"github.com/GoAdminGroup/go-admin/plugins/admin/modules/response"
	"github.com/GoAdminGroup/go-admin/plugins/admin/modules/table"
	"github.com/GoAdminGroup/go-admin/template/types"
	"github.com/GoAdminGroup/go-admin/template/types/action"
	"github.com/GoAdminGroup/go-admin/template/types/form"
	table2 "github.com/GoAdminGroup/go-admin/template/types/table"
	"github.com/leonelquinteros/gotext"
)

func GetOperationServer(ctx *context.Context) table.Table {

	tGameServers := table.NewDefaultTable(table.Config{
		Driver:     "mysql",
		Connection: "db_global",
		CanAdd:     true,
		Editable:   true,
		Deletable:  false,
		Exportable: false,
		PrimaryKey: table.PrimaryKey{
			Type: db.Int,
			Name: "Id",
		},
	})
	open := def.GetServerStatusOP(1)
	flag := def.GetServerStatusOP(2)
	openString := def.GetServerOpenStatusOpString()
	flagString := def.GetServerSpecialOpString()
	info := tGameServers.GetInfo().HideDetailButton().HideDeleteButton()
	myOP := fusion.GetServerList()

	// info.AddField(gotext.Get("序号"), "customOrder", db.Int).
	// 	FieldSortable().FieldEditAble()
	info.AddField("Id", "Id", db.Int).
		FieldSortable().
		FieldFilterable(types.FilterType{Options: myOP, FormType: form.SelectSingle}).
		FieldFilterOptions(myOP)
	info.AddField(gotext.Get("服务器名"), "logicName", db.Char).
		FieldSortable()
	info.AddField(gotext.Get("服务器开启状态"), "logicOpenStatus", db.Tinyint).
		FieldFilterable(types.FilterType{
			FormType: form.Select,
			Options:  open}).
		FieldSortable().
		FieldDisplay(func(model types.FieldModel) interface{} {
			temp, _ := strconv.Atoi(model.Value)
			return openString[temp]
		})
	info.AddField(gotext.Get("服务器标签"), "logicSpecialFlags", db.Int).
		FieldFilterable(types.FilterType{
			FormType: form.Select,
			Options:  flag}).
		FieldSortable().
		FieldDisplay(func(model types.FieldModel) interface{} {
			temp, _ := strconv.Atoi(model.Value)
			return flagString[temp]
		})
	info.AddField(gotext.Get("开服时间"), "logicOpenTime", db.Datetime).
		FieldEditAble(table2.Datetime).SetPageSizeList([]int{5, 10, 20, 100, 9999})

	info.AddActionButton(template.HTML(gotext.Get("修改后台服务器读取数据库配置")), action.Jump("/admin/EditServerDatabaseCfg?Id={{.Id}}"+
		"&logicName={{(index .Value \"logicName\").Value}}"))

	info.AddActionButton(template.HTML(gotext.Get("修改备份服务器读取数据库配置")), action.Jump("/admin/EditServerDatabaseBackupCfg?Id={{.Id}}"+
		"&logicName={{(index .Value \"logicName\").Value}}"))

	info.SetGetDataFn(
		func(param parameter.Parameters) (data []map[string]interface{}, size int) {
			db_global := fusion.GetBaseGormDB("db_global")
			if param.Fields["__pk"] != nil && param.Fields["__pk"][0] != "" {
				db_global.Raw("SELECT *FROM t_game_servers WHERE id = ?", param.Fields["__pk"][0]).
					Scan(&data)
				return data, 1
			} else {
				db_global.Raw("SELECT *FROM t_game_servers WHERE id IN (SELECT MIN(id) FROM t_game_servers GROUP BY externalIP,externalPort)").
					Scan(&data)
				min := math.Min(float64(param.PageInt*param.PageSizeInt), float64(len(data)))
				return data[((param.PageInt - 1) * param.PageSizeInt):int(min)], len(data)
			}
		})

	//info.AddButton(template.HTML(gotext.Get("清理数据库")), icon.Save,
	//	action.PopUp("/admin/ClearDB", gotext.Get("操作结果"), func(ctx *context.Context) (success bool, msg string, data interface{}) {
	//		faild := make([]string, 0)
	//		ids := ctx.Request.FormValue("ids")
	//		gsids := strings.Split(ids, ",")
	//		ifEmpty := true
	//		sqlSource := "SELECT CONCAT('TRUNCATE ', table_name, ';')\n" +
	//			"FROM information_schema.tables\n" +
	//			"WHERE table_schema = '%s';\n"
	//		deleteSqlSource := "SELECT CONCAT('DROP TABLE IF EXISTS ', table_name, ';')\n" +
	//			"FROM information_schema.tables\n" +
	//			"WHERE table_schema = '%s' and table_name like \"%%_2%%\";\n"
	//
	//		for _, v := range gsids {
	//			if v != "" {
	//				ifEmpty = false
	//				id, _ := strconv.Atoi(v)
	//				var dbName string
	//				var sqls = make([]string, 0)
	//
	//				db_char := fusion.GetServerGormDB("db_char", int64(id))
	//				if db_char == nil {
	//					faild = append(faild, "<h5>"+v+": db_char is nil, please check database config"+"</h5>")
	//				} else {
	//					db_char.Raw("select database ()").Scan(&dbName)
	//					sql := fmt.Sprintf(sqlSource, dbName)
	//					db_char.Raw(sql).Scan(&sqls)
	//					for _, truncateSql := range sqls {
	//						db_char.Exec(truncateSql)
	//					}
	//					faild = append(faild, "<h5>"+v+":db_char Clear Success"+"</h5>")
	//				}
	//
	//				db_log := fusion.GetServerGormDB("db_log", int64(id))
	//				if db_log == nil {
	//					faild = append(faild, "<h5>"+v+": db_log is nil, please check database config"+"</h5>")
	//				} else {
	//					dbName = ""
	//					db_log.Raw("select database ()").Scan(&dbName)
	//					deleteSql := fmt.Sprintf(deleteSqlSource, dbName)
	//					sqls = make([]string, 0)
	//					db_log.Raw(deleteSql).Scan(&sqls)
	//					for _, temp := range sqls {
	//						db_log.Exec(temp)
	//					}
	//
	//					sql := fmt.Sprintf(sqlSource, dbName)
	//					sqls = make([]string, 0)
	//					db_log.Raw(sql).Scan(&sqls)
	//					for _, truncateSql := range sqls {
	//						db_log.Exec(truncateSql)
	//					}
	//					faild = append(faild, "<h5>"+v+":db_log Clear Success"+"</h5>")
	//				}
	//
	//				db_global := fusion.GetBaseGormDB("db_global")
	//				if db_global == nil {
	//					faild = append(faild, "<h5>"+v+": db_global is nil, please check database config"+"</h5>")
	//				} else {
	//					db_global.Table("t_account_characters").Where("groupId = ?", v).Delete("t_account_characters")
	//					db_global.Table("t_guilds").Where("groupId = ?", v).Delete("t_guilds")
	//					faild = append(faild, "<h5>"+v+":db_global :t_account_characters,t_guilds Clear Success"+"</h5>")
	//				}
	//				var datas = make(map[string]interface{})
	//				go_manager := fusion.GetBaseGormDB("go_manager")
	//				go_manager.Table("databasecfg").Select("redisSn,deployHost,deployPort").Where("serverId = ?", v).Scan(&datas)
	//
	//				vars := url.Values{}
	//				ctx1 := &context.Context{}
	//				vars.Add("redisSn", fusion.Strval(datas["redisSn"]))
	//				err, res := fusion.CallToDeploy(&vars, common.MyCfg.Api["CleanRedis"], "GET", ctx1, fusion.Strval(datas["deployHost"]), fusion.Strval(datas["deployPort"]))
	//				if err != nil {
	//					faild = append(faild, "<h5>"+v+err.Error()+"</h5>")
	//				}
	//				if res != "" {
	//					if res == "All Is OK!" {
	//						faild = append(faild, "<h5>"+v+":Redis Clear Success"+"</h5>")
	//					}
	//					faild = append(faild, "<h5>"+v+":"+res+"</h5>")
	//				}
	//
	//			}
	//		}
	//		if ifEmpty {
	//			return false, gotext.Get("请选择服务器"), ""
	//		}
	//		return true, "", faild
	//	}))

	//info.AddButton(template.HTML(gotext.Get("开启进程守护")), "", action.PopUp("/server/deployResFiles", gotext.Get("操作结果"),
	//	func(ctx *context.Context) (success bool, msg string, data interface{}) {
	//		faild := []string{}
	//		ids := ctx.Request.FormValue("ids")
	//		gsids := strings.Split(ids, ",")
	//		vars := url.Values{}
	//		vars.Add("mode", "0")
	//
	//		ifEmpty := true
	//		for _, v := range gsids {
	//			if v != "" {
	//				ifEmpty = false
	//				ctx1 := &context.Context{}
	//				var datas = make(map[string]interface{})
	//				go_manager := fusion.GetBaseGormDB("go_manager")
	//				go_manager.Table("databasecfg").Select("redisSn,deployHost,deployPort").Where("serverId = ?", v).Scan(&datas)
	//
	//				err, res := fusion.CallToDeploy(&vars, common.MyCfg.Api["ChangeWatchMode"], "GET", ctx1,
	//					fusion.Strval(datas["deployHost"]), fusion.Strval(datas["deployPort"]))
	//				if err != nil {
	//					faild = append(faild, "<h5>"+v+err.Error()+"</h5>")
	//				}
	//				if res != "" {
	//					if res == "All Is OK!" {
	//						faild = append(faild, "<h5>"+v+":process watch start"+"</h5>")
	//					}
	//					faild = append(faild, "<h5>"+v+":"+res+"</h5>")
	//				}
	//			}
	//		}
	//
	//		if ifEmpty {
	//			return false, gotext.Get("请选择服务器"), ""
	//		}
	//		return true, "", faild
	//	}))

	info.AddButton(template.HTML(gotext.Get("向玩家开放")), "",
		action.PopUp("/admin/OpenServerToPlayer", gotext.Get("操作结果"), func(ctx *context.Context) (success bool, msg string, data interface{}) {
			faild := []string{}
			ids := ctx.Request.FormValue("ids")
			gsids := strings.Split(ids, ",")
			ifEmpty := true
			for _, v := range gsids {
				if v != "" {
					ifEmpty = false
					db_global := fusion.GetBaseGormDB("db_global")
					db_global.Table("t_game_servers").Where("Id", v).Update("logicOpenStatus", def.ServerOpenStatus_Normal)
				}
			}
			ctx1 := &context.Context{}
			vars := url.Values{}
			err, res := fusion.CallToCenter(&vars, common.MyCfg.Api["reloadGameServerInfos"], "GET", ctx1)
			if err != nil {
				faild = append(faild, "<h5>"+err.Error()+"</h5>")
			}
			if res != "" {
				if res == "All Is OK!" {
					faild = append(faild, "<h5>"+gotext.Get("已将选中服务器设为开放状态")+"</h5>")
				} else {
					faild = append(faild, "<h5>"+":"+res+"</h5>")
				}
			}
			if ifEmpty {
				return false, gotext.Get("请选择服务器"), ""
			}
			return true, "", faild
		}))

	// info.AddButton(template.HTML(gotext.Get("热更新服务器")), "", action.PopUpWithForm(action.PopUpData{
	// 	Id:     "/admin/popup/reqHotfixServer",
	// 	Title:  gotext.Get("热更新服务器"),
	// 	Width:  "1400px",
	// 	Height: "830px",
	// }, func(panel *types.FormPanel) *types.FormPanel {
	// 	panel.AddField(gotext.Get("服务器ID"), "gsIdX", db.Int, form.SelectBox).
	// 		FieldOptions(fusion.GetServerList())
	// 	panel.AddField(gotext.Get("热更新文件名"), "fileName", db.Varchar, form.Text)
	// 	panel.AddField(gotext.Get("热更新参数"), "hotfixArgs", db.Varchar, form.Text)
	// 	panel.EnableAjax()
	// 	return panel
	// }, "/admin/popup/HotfixServer"))

	//info.AddButton(template.HTML(gotext.Get("开启服务器(清空数据后)")), icon.Warning,
	//	action.PopUp("/admin/startServer4ClearDb", gotext.Get("操作结果"), func(ctx *context.Context) (success bool, msg string, data interface{}) {
	//		faild := make([]string, 0)
	//		ids := ctx.Request.FormValue("ids")
	//		gsids := strings.Split(ids, ",")
	//		ifEmpty := true
	//		sqlSource := "SELECT CONCAT('TRUNCATE ', table_name, ';')\n" +
	//			"FROM information_schema.tables\n" +
	//			"WHERE table_schema = '%s';\n"
	//		deleteSqlSource := "SELECT CONCAT('DROP TABLE IF EXISTS ', table_name, ';')\n" +
	//			"FROM information_schema.tables\n" +
	//			"WHERE table_schema = '%s' and table_name like \"%%_2%%\";\n"
	//
	//		for _, v := range gsids {
	//			if v != "" {
	//				ifEmpty = false
	//				id, _ := strconv.Atoi(v)
	//				var dbName string
	//				var sqls = make([]string, 0)
	//
	//				db_char := fusion.GetServerGormDB("db_char", int64(id))
	//				if db_char == nil {
	//					faild = append(faild, "<h5>"+v+": db_char is nil, please check database config"+"</h5>")
	//				} else {
	//					db_char.Raw("select database ()").Scan(&dbName)
	//					sql := fmt.Sprintf(sqlSource, dbName)
	//					db_char.Raw(sql).Scan(&sqls)
	//					for _, truncateSql := range sqls {
	//						db_char.Exec(truncateSql)
	//					}
	//					faild = append(faild, "<h5>"+v+":db_char Clear Success"+"</h5>")
	//				}
	//
	//				db_log := fusion.GetServerGormDB("db_log", int64(id))
	//				if db_log == nil {
	//					faild = append(faild, "<h5>"+v+": db_log is nil, please check database config"+"</h5>")
	//				} else {
	//					dbName = ""
	//					db_log.Raw("select database ()").Scan(&dbName)
	//					deleteSql := fmt.Sprintf(deleteSqlSource, dbName)
	//					sqls = make([]string, 0)
	//					db_log.Raw(deleteSql).Scan(&sqls)
	//					for _, temp := range sqls {
	//						db_log.Exec(temp)
	//					}
	//
	//					sql := fmt.Sprintf(sqlSource, dbName)
	//					sqls = make([]string, 0)
	//					db_log.Raw(sql).Scan(&sqls)
	//					for _, truncateSql := range sqls {
	//						db_log.Exec(truncateSql)
	//					}
	//					faild = append(faild, "<h5>"+v+":db_log Clear Success"+"</h5>")
	//				}
	//
	//				db_global := fusion.GetBaseGormDB("db_global")
	//				if db_global == nil {
	//					faild = append(faild, "<h5>"+v+": db_global is nil, please check database config"+"</h5>")
	//				} else {
	//					db_global.Table("t_account_characters").Where("groupId = ?", v).Delete("t_account_characters")
	//					db_global.Table("t_guilds").Where("groupId = ?", v).Delete("t_guilds")
	//					faild = append(faild, "<h5>"+v+":db_global :t_account_characters,t_guilds Clear Success"+"</h5>")
	//				}
	//				var datas = make(map[string]interface{})
	//				go_manager := fusion.GetBaseGormDB("go_manager")
	//				go_manager.Table("databasecfg").Select("redisSn,deployHost,deployPort").Where("serverId = ?", v).Scan(&datas)
	//
	//				vars := url.Values{}
	//				ctx1 := &context.Context{}
	//				vars.Add("redisSn", fusion.Strval(datas["redisSn"]))
	//				err, res := fusion.CallToDeploy(&vars, common.MyCfg.Api["CleanRedis"], "GET", ctx1, fusion.Strval(datas["deployHost"]), fusion.Strval(datas["deployPort"]))
	//				if err != nil {
	//					faild = append(faild, "<h5>"+v+err.Error()+"</h5>")
	//				}
	//				if res != "" {
	//					if res == "All Is OK!" {
	//						faild = append(faild, "<h5>"+v+":Redis Clear Success"+"</h5>")
	//					} else {
	//						faild = append(faild, "<h5>"+v+":"+res+"</h5>")
	//					}
	//				}
	//
	//				vars = url.Values{}
	//				ctx1 = &context.Context{}
	//				vars.Add("mode", "0")
	//				err, res = fusion.CallToDeploy(&vars, common.MyCfg.Api["ChangeWatchMode"], "GET", ctx1,
	//					fusion.Strval(datas["deployHost"]), fusion.Strval(datas["deployPort"]))
	//				if err != nil {
	//					faild = append(faild, "<h5>"+v+err.Error()+"</h5>")
	//				}
	//				if res != "" {
	//					if res == "All Is OK!" {
	//						faild = append(faild, "<h5>"+v+":process watch start"+"</h5>")
	//					} else {
	//						faild = append(faild, "<h5>"+v+":"+res+"</h5>")
	//					}
	//				}
	//			}
	//		}
	//		if ifEmpty {
	//			return false, gotext.Get("请选择服务器"), ""
	//		}
	//		return true, "", faild
	//	}))

	info.AddButton(template.HTML(gotext.Get("开启服务器(更新后)")), "", action.PopUp("/server/startServer4Update", gotext.Get("操作结果"),
		func(ctx *context.Context) (success bool, msg string, data interface{}) {
			faild := []string{}
			ids := ctx.Request.FormValue("ids")
			gsids := strings.Split(ids, ",")
			vars := url.Values{}
			vars.Add("timeout", "10")

			ifEmpty := true
			for _, v := range gsids {
				if v != "" {
					ifEmpty = false
					ctx1 := &context.Context{}
					var datas = make(map[string]interface{})
					go_manager := fusion.GetBaseGormDB("go_manager")
					go_manager.Table("databasecfg").Select("redisSn,deployHost,deployPort").Where("serverId = ?", v).Scan(&datas)

					err, res := fusion.CallToDeploy(&vars, common.MyCfg.Api["RestartServer"], "GET", ctx1,
						fusion.Strval(datas["deployHost"]), fusion.Strval(datas["deployPort"]))
					if err != nil {
						faild = append(faild, "<h5>"+v+err.Error()+"</h5>")
					}
					if res != "" {
						if res == "All Is OK!" {
							faild = append(faild, "<h5>"+v+":process watch start"+"</h5>")
						} else {
							faild = append(faild, "<h5>"+v+":"+res+"</h5>")
						}
					}
				}
			}

			if ifEmpty {
				return false, gotext.Get("请选择服务器"), ""
			}
			return true, "", faild
		}))

	//info.AddButton(template.HTML(gotext.Get("部署资源文件")), "", action.PopUp("/server/deployResFiles", gotext.Get("操作结果"),
	//	func(ctx *context.Context) (success bool, msg string, data interface{}) {
	//		faild := []string{}
	//		ids := ctx.Request.FormValue("ids")
	//		gsids := strings.Split(ids, ",")
	//		vars := url.Values{}
	//
	//		ifEmpty := true
	//		for _, v := range gsids {
	//			if v != "" {
	//				ifEmpty = false
	//				ctx1 := &context.Context{}
	//				var datas = make(map[string]interface{})
	//				go_manager := fusion.GetBaseGormDB("go_manager")
	//				go_manager.Table("databasecfg").Select("redisSn,deployHost,deployPort").Where("serverId = ?", v).Scan(&datas)
	//
	//				err, res := fusion.CallToDeploy(&vars, common.MyCfg.Api["deployResFiles"], "GET", ctx1,
	//					fusion.Strval(datas["deployHost"]), fusion.Strval(datas["deployPort"]))
	//				if err != nil {
	//					faild = append(faild, "<h5>"+v+err.Error()+"</h5>")
	//				}
	//				if res != "" {
	//					if res == "All Is OK!" {
	//						faild = append(faild, "<h5>"+v+":ResFile Deploy Success"+"</h5>")
	//					}
	//					faild = append(faild, "<h5>"+v+":"+res+"</h5>")
	//				}
	//			}
	//		}
	//
	//		if ifEmpty {
	//			return false, gotext.Get("请选择服务器"), ""
	//		}
	//		return true, "", faild
	//	}))
	//
	//info.AddButton(template.HTML(gotext.Get("下载资源文件")), "", action.PopUp("/server/downResFiles", gotext.Get("操作结果"),
	//	func(ctx *context.Context) (success bool, msg string, data interface{}) {
	//		faild := []string{}
	//		ids := ctx.Request.FormValue("ids")
	//		gsids := strings.Split(ids, ",")
	//		vars := url.Values{}
	//
	//		ifEmpty := true
	//		for _, v := range gsids {
	//			if v != "" {
	//				ifEmpty = false
	//				ctx1 := &context.Context{}
	//				var datas = make(map[string]interface{})
	//				go_manager := fusion.GetBaseGormDB("go_manager")
	//				go_manager.Table("databasecfg").Select("redisSn,deployHost,deployPort").Where("serverId = ?", v).Scan(&datas)
	//
	//				err, res := fusion.CallToDeploy(&vars, common.MyCfg.Api["downResFiles"], "GET", ctx1,
	//					fusion.Strval(datas["deployHost"]), fusion.Strval(datas["deployPort"]))
	//				if err != nil {
	//					faild = append(faild, "<h5>"+v+err.Error()+"</h5>")
	//				}
	//				if res != "" {
	//					if res == "All Is OK!" {
	//						faild = append(faild, "<h5>"+v+":ResFile download Success"+"</h5>")
	//					}
	//					faild = append(faild, "<h5>"+v+":"+res+"</h5>")
	//				}
	//			}
	//		}
	//
	//		if ifEmpty {
	//			return false, gotext.Get("请选择服务器"), ""
	//		}
	//		return true, "", faild
	//	}))
	//
	//info.AddButton(template.HTML(gotext.Get("清除资源文件")), "", action.PopUp("/server/clearResFile", gotext.Get("操作结果"),
	//	func(ctx *context.Context) (success bool, msg string, data interface{}) {
	//		faild := []string{}
	//		ids := ctx.Request.FormValue("ids")
	//		gsids := strings.Split(ids, ",")
	//		vars := url.Values{}
	//
	//		ifEmpty := true
	//		for _, v := range gsids {
	//			if v != "" {
	//				ifEmpty = false
	//				ctx1 := &context.Context{}
	//				var datas = make(map[string]interface{})
	//				go_manager := fusion.GetBaseGormDB("go_manager")
	//				go_manager.Table("databasecfg").Select("redisSn,deployHost,deployPort").Where("serverId = ?", v).Scan(&datas)
	//
	//				err, res := fusion.CallToDeploy(&vars, common.MyCfg.Api["clearResFiles"], "GET", ctx1,
	//					fusion.Strval(datas["deployHost"]), fusion.Strval(datas["deployPort"]))
	//				if err != nil {
	//					faild = append(faild, "<h5>"+v+err.Error()+"</h5>")
	//				}
	//				if res != "" {
	//					if res == "All Is OK!" {
	//						faild = append(faild, "<h5>"+v+":ResFile Clear Success"+"</h5>")
	//					}
	//					faild = append(faild, "<h5>"+v+":"+res+"</h5>")
	//				}
	//			}
	//		}
	//
	//		if ifEmpty {
	//			return false, gotext.Get("请选择服务器"), ""
	//		}
	//		return true, "", faild
	//	}))

	// info.AddButton(template.HTML(gotext.Get("清理旧资源->更新新资源->解压部署资源")), "", action.PopUp("/server/clear_update_deploy", gotext.Get("操作结果"),
	// 	func(ctx *context.Context) (success bool, msg string, data interface{}) {
	// 		faild := []string{}
	// 		ids := ctx.Request.FormValue("ids")
	// 		gsids := strings.Split(ids, ",")
	// 		vars := url.Values{}

	// 		ifEmpty := true
	// 		for _, v := range gsids {
	// 			if v != "" {
	// 				ifEmpty = false
	// 				ctx1 := &context.Context{}
	// 				var datas = make(map[string]interface{})
	// 				go_manager := fusion.GetBaseGormDB("go_manager")
	// 				go_manager.Table("databasecfg").Select("redisSn,deployHost,deployPort").Where("serverId = ?", v).Scan(&datas)

	// 				err, res := fusion.CallToDeploy(&vars, common.MyCfg.Api["clearResFiles"], "GET", ctx1,
	// 					fusion.Strval(datas["deployHost"]), fusion.Strval(datas["deployPort"]))
	// 				if err != nil {
	// 					faild = append(faild, "<h5>"+v+err.Error()+"</h5>")
	// 				}
	// 				if res != "" {
	// 					if res == "All Is OK!" {
	// 						faild = append(faild, "<h5>"+v+":ResFile Clear Success"+"</h5>")
	// 					} else {
	// 						faild = append(faild, "<h5>"+v+":"+res+"</h5>")
	// 					}
	// 				}

	// 				err, res = fusion.CallToDeploy(&vars, common.MyCfg.Api["downResFiles"], "GET", ctx1,
	// 					fusion.Strval(datas["deployHost"]), fusion.Strval(datas["deployPort"]))
	// 				if err != nil {
	// 					faild = append(faild, "<h5>"+v+err.Error()+"</h5>")
	// 				}
	// 				if res != "" {
	// 					if res == "All Is OK!" {
	// 						faild = append(faild, "<h5>"+v+":ResFile download Success"+"</h5>")
	// 					} else {
	// 						faild = append(faild, "<h5>"+v+":"+res+"</h5>")
	// 					}
	// 				}

	// 				err, res = fusion.CallToDeploy(&vars, common.MyCfg.Api["deployResFiles"], "GET", ctx1,
	// 					fusion.Strval(datas["deployHost"]), fusion.Strval(datas["deployPort"]))
	// 				if err != nil {
	// 					faild = append(faild, "<h5>"+v+err.Error()+"</h5>")
	// 				}
	// 				if res != "" {
	// 					if res == "All Is OK!" {
	// 						faild = append(faild, "<h5>"+v+":ResFile Deploy Success"+"</h5>")
	// 					} else {
	// 						faild = append(faild, "<h5>"+v+":"+res+"</h5>")
	// 					}
	// 				}

	// 			}
	// 		}
	// 		if ifEmpty {
	// 			return false, gotext.Get("请选择服务器"), ""
	// 		}
	// 		return true, "", faild
	// 	}))

	//info.AddButton(template.HTML(gotext.Get("关闭选中服务器")), "",
	//	action.PopUp("/admin/CloseServer", gotext.Get("操作结果"), func(ctx *context.Context) (success bool, msg string, data interface{}) {
	//		faild := []string{}
	//		ids := ctx.Request.FormValue("ids")
	//		gsids := strings.Split(ids, ",")
	//		ifEmpty := true
	//		for _, v := range gsids {
	//			if v != "" {
	//				ifEmpty = false
	//				vars := &url.Values{}
	//				vars.Add("name", "AdminServer")
	//				vars.Add("cmd", "ShutdownServer")
	//				vars.Add("gsId", v)
	//				ctx1 := &context.Context{}
	//				err, res := fusion.CallToCenter(vars, common.MyCfg.Api["GM2GS"], "GET", ctx1)
	//				if err != nil {
	//					faild = append(faild, "<h5>"+v+":"+err.Error()+"</h5>")
	//				}
	//				if res != "" {
	//					faild = append(faild, "<h5>"+res+"</h5>")
	//				}
	//			}
	//		}
	//		if ifEmpty {
	//			return false, gotext.Get("请选择服务器"), ""
	//		}
	//		return true, "", faild
	//	}))
	//
	//info.AddButton(template.HTML(gotext.Get("关闭进程守护")), "", action.PopUp("/server/CloseProcessWatch", gotext.Get("操作结果"),
	//	func(ctx *context.Context) (success bool, msg string, data interface{}) {
	//		faild := []string{}
	//		ids := ctx.Request.FormValue("ids")
	//		gsids := strings.Split(ids, ",")
	//		vars := url.Values{}
	//		vars.Add("mode", "1")
	//
	//		ifEmpty := true
	//		for _, v := range gsids {
	//			if v != "" {
	//				ifEmpty = false
	//				ctx1 := &context.Context{}
	//				var datas = make(map[string]interface{})
	//				go_manager := fusion.GetBaseGormDB("go_manager")
	//				go_manager.Table("databasecfg").Select("redisSn,deployHost,deployPort").Where("serverId = ?", v).Scan(&datas)
	//
	//				err, res := fusion.CallToDeploy(&vars, common.MyCfg.Api["ChangeWatchMode"], "GET", ctx1,
	//					fusion.Strval(datas["deployHost"]), fusion.Strval(datas["deployPort"]))
	//				if err != nil {
	//					faild = append(faild, "<h5>"+v+err.Error()+"</h5>")
	//				}
	//				if res != "" {
	//					if res == "All Is OK!" {
	//						faild = append(faild, "<h5>"+v+":process watch closed"+"</h5>")
	//					}
	//					faild = append(faild, "<h5>"+v+":"+res+"</h5>")
	//				}
	//			}
	//		}
	//
	//		if ifEmpty {
	//			return false, gotext.Get("请选择服务器"), ""
	//		}
	//		return true, "", faild
	//	}))

	info.AddButton(template.HTML(gotext.Get("关服")), "", action.PopUp("/server/closeServerProcess", gotext.Get("操作结果"),
		func(ctx *context.Context) (success bool, msg string, data interface{}) {
			faild := []string{}
			ids := ctx.Request.FormValue("ids")
			gsids := strings.Split(ids, ",")

			ifEmpty := true
			for _, v := range gsids {
				if v != "" {

					ifEmpty = false
					ctx1 := &context.Context{}

					var datas = make(map[string]interface{})
					go_manager := fusion.GetBaseGormDB("go_manager")
					go_manager.Table("databasecfg").Select("redisSn,deployHost,deployPort").Where("serverId = ?", v).Scan(&datas)

					vars := url.Values{}
					vars.Add("mode", "1")
					err, res := fusion.CallToDeploy(&vars, common.MyCfg.Api["ChangeWatchMode"], "GET", ctx1,
						fusion.Strval(datas["deployHost"]), fusion.Strval(datas["deployPort"]))
					if err != nil {
						faild = append(faild, "<h5>"+v+err.Error()+"</h5>")
					}
					if res != "" {
						if res == "All Is OK!" {
							faild = append(faild, "<h5>"+v+":process watch closed"+"</h5>")
						} else {
							faild = append(faild, "<h5>"+v+":"+res+"</h5>")
						}
					}

					vars = url.Values{}
					vars.Add("name", "AdminServer")
					vars.Add("cmd", "ShutdownServer")
					vars.Add("gsId", v)
					var api string
					if find := strings.Contains(vars.Get("name"), "MapServer"); find {
						api = common.MyCfg.Api["GM2MS"]
					} else if find = strings.Contains(vars.Get("name"), "GameServer"); find {
						api = common.MyCfg.Api["GM2GS"]
					} else if find = strings.Contains(vars.Get("name"), "GateServer"); find {
						api = common.MyCfg.Api["GM2GATE"]
					} else if find = strings.Contains(vars.Get("name"), "SocialServer"); find {
						api = common.MyCfg.Api["GM2SOCIAL"]
					} else if find = strings.Contains(vars.Get("name"), "DBPServer"); find {
						api = common.MyCfg.Api["GM2DBP"]
					} else {
						api = common.MyCfg.Api["GM2S"]
					}
					ctx1 = &context.Context{}
					err, res = fusion.CallToCenter(&vars, api, "GET", ctx1)
					if err != nil {
						faild = append(faild, "<h5>"+v+":"+err.Error()+"</h5>")
					}
					if res != "" {
						faild = append(faild, "<h5>"+res+"</h5>")
					}

					db_global := fusion.GetBaseGormDB("db_global")
					temp := make(map[string]interface{})
					db_global.Table("t_game_servers").Where("Id", v).Select("externalIP", "externalPort").Scan(&temp)
					db_global.Table("t_game_servers").Where("externalIP", temp["externalIP"]).
						Where("externalPort", temp["externalPort"]).Where("logicOpenStatus", def.ServerOpenStatus_Normal).Update("logicOpenStatus", def.ServerOpenStatus_Maintain)
				}
			}
			ctx1 := &context.Context{}
			vars := url.Values{}
			err, res := fusion.CallToCenter(&vars, common.MyCfg.Api["reloadGameServerInfos"], "GET", ctx1)
			if err != nil {
				faild = append(faild, "<h5>"+err.Error()+"</h5>")
			}
			if res != "" {
				if res == "All Is OK!" {
					faild = append(faild, "<h5>"+gotext.Get("已将选中服务器设为维护状态")+"</h5>")
				} else {
					faild = append(faild, "<h5>"+":"+res+"</h5>")
				}
			}

			if ifEmpty {
				return false, gotext.Get("请选择服务器"), ""
			}
			return true, "", faild
		}))

	info.AddButton(template.HTML(gotext.Get("开启进程守护")), "", action.PopUp("/server/startProcessPro", gotext.Get("操作结果"),
		func(ctx *context.Context) (success bool, msg string, data interface{}) {
			faild := []string{}
			ids := ctx.Request.FormValue("ids")
			gsids := strings.Split(ids, ",")
			vars := url.Values{}
			vars.Add("mode", "0")

			ifEmpty := true
			for _, v := range gsids {
				if v != "" {
					ifEmpty = false
					ctx1 := &context.Context{}
					deployHost, deployPort, deployName := fusion.GetDeployTargetByServerId(v)
					vars.Set("name", deployName)
					err, res := fusion.CallToDeploy(&vars, common.MyCfg.Api["ChangeWatchMode"], "GET", ctx1,
						deployHost, deployPort)
					if err != nil {
						faild = append(faild, "<h5>"+v+err.Error()+"</h5>")
					}
					if res != "" {
						if res == "All Is OK!" {
							faild = append(faild, "<h5>"+v+":process watch start"+"</h5>")
						}
						faild = append(faild, "<h5>"+v+":"+res+"</h5>")
					}
				}
			}

			if ifEmpty {
				return false, gotext.Get("请选择服务器"), ""
			}
			return true, "", faild
		}))

	info.AddButton(template.HTML(gotext.Get("检查是否为跨服")), "", action.PopUpWithForm(action.PopUpData{
		Id:     "/server/isCrossServerPro",
		Title:  gotext.Get("检查是否为跨服"),
		Width:  "1400px",
		Height: "830px",
	}, func(panel *types.FormPanel) *types.FormPanel {
		panel.AddField(gotext.Get(""), "operation", db.Int, form.Text).FieldHide().FieldDefault("9")
		panel.AddField(gotext.Get("服务器ID"), "gsIdX", db.Int, form.SelectBox).FieldFoot(template.HTML(ChangeSelectBoxHeight)).
			FieldOptions(fusion.GetServerList())
		panel.AddField(gotext.Get("分组名"), "groupName", db.Varchar, form.Text)
		panel.EnableAjax()
		return panel
	}, "/admin/popup/OperationResFilesPro"))

	info.AddButton(template.HTML(gotext.Get("检查进程数量")), "", action.PopUpWithForm(action.PopUpData{
		Id:     "/server/CheckProcessPro",
		Title:  gotext.Get("检查进程数量"),
		Width:  "1400px",
		Height: "830px",
	}, func(panel *types.FormPanel) *types.FormPanel {
		panel.AddField(gotext.Get(""), "operation", db.Int, form.Text).FieldHide().FieldDefault("10")
		panel.AddField(gotext.Get("服务器ID"), "gsIdX", db.Int, form.SelectBox).FieldFoot(template.HTML(ChangeSelectBoxHeight)).
			FieldOptions(fusion.GetServerList())
		panel.AddField(gotext.Get("分组名"), "groupName", db.Varchar, form.Text)
		panel.EnableAjax()
		return panel
	}, "/admin/popup/OperationResFilesPro"))

	info.AddButton(template.HTML(gotext.Get("重新启动服务器")), "", action.PopUpWithForm(action.PopUpData{
		Id:     "/server/RestartServerPro",
		Title:  gotext.Get("重新启动服务器"),
		Width:  "1400px",
		Height: "830px",
	}, func(panel *types.FormPanel) *types.FormPanel {
		panel.AddField(gotext.Get(""), "operation", db.Int, form.Text).FieldHide().FieldDefault("4")
		panel.AddField(gotext.Get("服务器ID"), "gsIdX", db.Int, form.SelectBox).FieldFoot(template.HTML(ChangeSelectBoxHeight)).
			FieldOptions(fusion.GetServerList())
		panel.AddField(gotext.Get("分组名"), "groupName", db.Varchar, form.Text)
		panel.EnableAjax()
		return panel
	}, "/admin/popup/OperationResFilesPro"))

	info.AddButton(template.HTML(gotext.Get("强制关闭服务器（kill）")), "", action.PopUpWithForm(action.PopUpData{
		Id:     "/server/StopServerPro",
		Title:  gotext.Get("强制关闭服务器（kill）"),
		Width:  "1400px",
		Height: "830px",
	}, func(panel *types.FormPanel) *types.FormPanel {
		panel.AddField(gotext.Get(""), "operation", db.Int, form.Text).FieldHide().FieldDefault("11")
		panel.AddField(gotext.Get("服务器ID"), "gsIdX", db.Int, form.SelectBox).FieldFoot(template.HTML(ChangeSelectBoxHeight)).
			FieldOptions(fusion.GetServerList())
		panel.AddField(gotext.Get("分组名"), "groupName", db.Varchar, form.Text)
		panel.EnableAjax()
		return panel
	}, "/admin/popup/OperationResFilesPro"))

	info.AddButton(template.HTML(gotext.Get("运行命令")), "", action.PopUpWithForm(action.PopUpData{
		Id:     "/server/RunBashPro",
		Title:  gotext.Get("运行命令"),
		Width:  "1400px",
		Height: "630px",
	}, func(panel *types.FormPanel) *types.FormPanel {
		panel.AddField(gotext.Get(""), "operation", db.Int, form.Text).FieldHide().FieldDefault("12")
		panel.AddField(gotext.Get("服务器ID"), "gsIdX", db.Int, form.SelectBox).FieldFoot(template.HTML(ChangeSelectBoxHeight)).
			FieldOptions(fusion.GetServerList())
		panel.AddField(gotext.Get("分组名"), "groupName", db.Varchar, form.Text)
		panel.AddField(gotext.Get("执行命令"), "cmdline", db.Varchar, form.Text)
		panel.EnableAjax()

		return panel
	}, "/admin/popup/OperationResFilesPro"))

	info.AddButton(template.HTML(gotext.Get("重启deploy")), "", action.PopUpWithForm(action.PopUpData{
		Id:     "/server/RebootSelfPro",
		Title:  gotext.Get("重启deploy"),
		Width:  "1400px",
		Height: "830px",
	}, func(panel *types.FormPanel) *types.FormPanel {
		panel.AddField(gotext.Get(""), "operation", db.Int, form.Text).FieldHide().FieldDefault("13")
		panel.AddField(gotext.Get("服务器ID"), "gsIdX", db.Int, form.SelectBox).FieldFoot(template.HTML(ChangeSelectBoxHeight)).
			FieldOptions(fusion.GetServerList())
		panel.AddField(gotext.Get("分组名"), "groupName", db.Varchar, form.Text)
		panel.EnableAjax()
		return panel
	}, "/admin/popup/OperationResFilesPro"))

	info.AddButton(template.HTML(gotext.Get("deploy重载配置文件")), "", action.PopUpWithForm(action.PopUpData{
		Id:     "/server/ReloadConfigPro",
		Title:  gotext.Get("deploy重载配置文件"),
		Width:  "1400px",
		Height: "830px",
	}, func(panel *types.FormPanel) *types.FormPanel {
		panel.AddField(gotext.Get(""), "operation", db.Int, form.Text).FieldHide().FieldDefault("14")
		panel.AddField(gotext.Get("服务器ID"), "gsIdX", db.Int, form.SelectBox).
			FieldOptions(fusion.GetServerList())
		panel.AddField(gotext.Get("分组名"), "groupName", db.Varchar, form.Text)
		panel.EnableAjax()
		return panel
	}, "/admin/popup/OperationResFilesPro"))

	info.AddButton(template.HTML(gotext.Get("热更数据表")), "", action.PopUpWithForm(action.PopUpData{
		Id:     "/admin/popup/hotfixDataTables",
		Title:  gotext.Get("热更数据表"),
		Width:  "1400px",
		Height: "830px",
	}, func(panel *types.FormPanel) *types.FormPanel {
		panel.AddField(gotext.Get("服务器ID"), "gsIdX", db.Int, form.SelectBox).FieldFoot(template.HTML(ChangeSelectBoxHeight)).
			FieldOptions(fusion.GetServerList())
		panel.AddField(gotext.Get("数据表名"), "tableNames", db.Varchar, form.Text).
			FieldDefault("")
		panel.EnableAjax()
		return panel
	}, "/admin/popup/HotfixDBTables"))

	info.AddButton(template.HTML(gotext.Get("热更3rd")), "", action.PopUp("/server/hotfix3rd", gotext.Get("操作结果"),
		func(ctx *context.Context) (success bool, msg string, data interface{}) {
			faild := []string{}
			ids := ctx.Request.FormValue("ids")
			gsids := strings.Split(ids, ",")

			ifEmpty := true
			for _, v := range gsids {
				if v != "" {
					ifEmpty = false
					ctx1 := &context.Context{}
					vars := &url.Values{}
					vars.Add("name", "MapServer")
					vars.Add("cmd", "HotfixSimulateScripts")
					vars.Add("gsId", v)
					var api string
					if find := strings.Contains(vars.Get("name"), "MapServer"); find {
						api = common.MyCfg.Api["GM2MS"]
					} else if find = strings.Contains(vars.Get("name"), "GameServer"); find {
						api = common.MyCfg.Api["GM2GS"]
					} else if find = strings.Contains(vars.Get("name"), "GateServer"); find {
						api = common.MyCfg.Api["GM2GATE"]
					} else if find = strings.Contains(vars.Get("name"), "SocialServer"); find {
						api = common.MyCfg.Api["GM2SOCIAL"]
					} else if find = strings.Contains(vars.Get("name"), "DBPServer"); find {
						api = common.MyCfg.Api["GM2DBP"]
					} else {
						api = common.MyCfg.Api["GM2S"]
					}
					err, res := fusion.CallToCenter(vars, api, "GET", ctx1)
					if err != nil {
						faild = append(faild, "<h5>"+v+":"+err.Error()+"</h5>")
					}
					if res != "" {
						if res == "All Is OK!" {
							faild = append(faild, "<h5>"+v+":热更3rd已通知"+"</h5>")
						} else {
							faild = append(faild, "<h5>"+v+":"+res+"</h5>")
						}
					}
				}
			}
			if ifEmpty {
				return false, gotext.Get("请选择服务器"), ""
			}

			return true, "", faild
		}))

	info.AddButton(template.HTML(gotext.Get("手动保存日志文件")), "", action.PopUp("/server/persistLogFiles", gotext.Get("操作结果"),
		func(ctx *context.Context) (success bool, msg string, data interface{}) {
			faild := []string{}
			ids := ctx.Request.FormValue("ids")
			gsids := strings.Split(ids, ",")
			vars := url.Values{}

			ifEmpty := true
			for _, v := range gsids {
				if v != "" {
					ifEmpty = false
					ctx1 := &context.Context{}
					deployHost, deployPort, deployName := fusion.GetDeployTargetByServerId(v)
					vars.Set("name", deployName)
					err, res := fusion.CallToDeploy(&vars, common.MyCfg.Api["persistLogFiles"], "GET", ctx1,
						deployHost, deployPort)
					if err != nil {
						faild = append(faild, "<h5>"+v+err.Error()+"</h5>")
					}
					if res != "" {
						if res == "All Is OK!" {
							faild = append(faild, "<h5>"+v+":日志文件已保存"+"</h5>")
						} else {
							faild = append(faild, "<h5>"+v+":"+res+"</h5>")
						}
					}
				}
			}
			if ifEmpty {
				return false, gotext.Get("请选择服务器"), ""
			}
			return true, "", faild
		}))

	info.AddButton(template.HTML(gotext.Get("关闭选中服务器")), "",
		action.PopUp("/admin/CloseServerPro", gotext.Get("操作结果"), func(ctx *context.Context) (success bool, msg string, data interface{}) {
			faild := []string{}
			ids := ctx.Request.FormValue("ids")
			gsids := strings.Split(ids, ",")
			ifEmpty := true
			for _, v := range gsids {
				if v != "" {
					ifEmpty = false
					vars := &url.Values{}
					vars.Add("name", "AdminServer")
					vars.Add("cmd", "ShutdownServer")
					vars.Add("gsId", v)
					vars.Add("name", "AdminServer")
					var api string
					if find := strings.Contains(vars.Get("name"), "MapServer"); find {
						api = common.MyCfg.Api["GM2MS"]
					} else if find = strings.Contains(vars.Get("name"), "GameServer"); find {
						api = common.MyCfg.Api["GM2GS"]
					} else if find = strings.Contains(vars.Get("name"), "GateServer"); find {
						api = common.MyCfg.Api["GM2GATE"]
					} else if find = strings.Contains(vars.Get("name"), "SocialServer"); find {
						api = common.MyCfg.Api["GM2SOCIAL"]
					} else if find = strings.Contains(vars.Get("name"), "DBPServer"); find {
						api = common.MyCfg.Api["GM2DBP"]
					} else {
						api = common.MyCfg.Api["GM2S"]
					}
					ctx1 := &context.Context{}
					err, res := fusion.CallToCenter(vars, api, "GET", ctx1)
					if err != nil {
						faild = append(faild, "<h5>"+v+":"+err.Error()+"</h5>")
					}
					if res != "" {
						faild = append(faild, "<h5>"+res+"</h5>")
					}
				}
			}
			if ifEmpty {
				return false, gotext.Get("请选择服务器"), ""
			}
			return true, "", faild
		}))

	info.AddButton(template.HTML(gotext.Get("关闭进程守护")), "", action.PopUp("/server/CloseProcessWatchPro", gotext.Get("操作结果"),
		func(ctx *context.Context) (success bool, msg string, data interface{}) {
			faild := []string{}
			ids := ctx.Request.FormValue("ids")
			gsids := strings.Split(ids, ",")
			vars := url.Values{}
			vars.Add("mode", "1")

			ifEmpty := true
			for _, v := range gsids {
				if v != "" {
					ifEmpty = false
					ctx1 := &context.Context{}
					deployHost, deployPort, deployName := fusion.GetDeployTargetByServerId(v)
					vars.Set("name", deployName)
					err, res := fusion.CallToDeploy(&vars, common.MyCfg.Api["ChangeWatchMode"], "GET", ctx1,
						deployHost, deployPort)
					if err != nil {
						faild = append(faild, "<h5>"+v+err.Error()+"</h5>")
					}
					if res != "" {
						if res == "All Is OK!" {
							faild = append(faild, "<h5>"+v+":process watch closed"+"</h5>")
						}
						faild = append(faild, "<h5>"+v+":"+res+"</h5>")
					}
				}
			}

			if ifEmpty {
				return false, gotext.Get("请选择服务器"), ""
			}
			return true, "", faild
		}))

	info.AddButton(template.HTML(gotext.Get("查看进程守护状态")), "", action.PopUp("/server/CheckWatchModePro", gotext.Get("操作结果"),
		func(ctx *context.Context) (success bool, msg string, data interface{}) {
			faild := []string{}
			ids := ctx.Request.FormValue("ids")
			gsids := strings.Split(ids, ",")
			vars := url.Values{}

			ifEmpty := true
			for _, v := range gsids {
				if v != "" {
					ifEmpty = false
					ctx1 := &context.Context{}
					deployHost, deployPort, deployName := fusion.GetDeployTargetByServerId(v)
					vars.Set("name", deployName)
					err, res := fusion.CallToDeploy(&vars, common.MyCfg.Api["CheckWatchMode"], "GET", ctx1,
						deployHost, deployPort)
					if err != nil {
						faild = append(faild, "<h5>"+v+err.Error()+"</h5>")
					}
					if res != "" {
						faild = append(faild, "<h5>"+v+":"+"进程守护状态为"+res+"</h5>")
					}
				}
			}

			if ifEmpty {
				return false, gotext.Get("请选择服务器"), ""
			}
			return true, "", faild
		}))

	info.AddButton(template.HTML(gotext.Get("批量设置服务器标签")), "", action.PopUpWithForm(action.PopUpData{
		Id:     "/admin/popup/setServerSymbolPro",
		Title:  gotext.Get("批量设置服务器标签"),
		Width:  "1400px",
		Height: "830px",
	}, func(panel *types.FormPanel) *types.FormPanel {
		panel.AddField("type", "OPType", db.Int, form.Text).
			FieldDefault("1").FieldHide()
		panel.AddField(gotext.Get("服务器ID"), "Id", db.Int, form.SelectBox).FieldFoot(template.HTML(ChangeSelectBoxHeight)).
			FieldOptions(fusion.GetServerList())
		panel.AddField(gotext.Get("服务器标签"), "logicSpecialFlags", db.Int, form.SelectSingle).
			FieldOptions(flag)
		panel.EnableAjax("操作成功", "操作失败，请联系管理员", "/admin/info/t_game_servers")
		return panel
	}, "/admin/popup/UpdateServerSettingPro"))

	info.AddButton(template.HTML(gotext.Get("批量设置服务器状态")), "", action.PopUpWithForm(action.PopUpData{
		Id:     "/admin/popup/setServerStatusPro",
		Title:  gotext.Get("批量设置服务器状态"),
		Width:  "1400px",
		Height: "830px",
	}, func(panel *types.FormPanel) *types.FormPanel {
		panel.AddField("type", "OPType", db.Int, form.Text).
			FieldDefault("2").FieldHide()
		panel.AddField(gotext.Get("服务器ID"), "Id", db.Int, form.SelectBox).FieldFoot(template.HTML(ChangeSelectBoxHeight)).
			FieldOptions(fusion.GetServerList())
		panel.AddField(gotext.Get("服务器开启状态"), "logicOpenStatus", db.Int, form.SelectSingle).
			FieldOptions(open)
		panel.EnableAjax("操作成功", "操作失败，请联系管理员", "/admin/info/t_game_servers")
		return panel
	}, "/admin/popup/UpdateServerSettingPro"))

		info.AddButton(template.HTML(gotext.Get("踢出所有玩家")), "", action.PopUp("/server/KickAllPlayerPro", gotext.Get("操作结果"),
			func(ctx *context.Context) (success bool, msg string, data interface{}) {
				faild := []string{}
				ids := ctx.Request.FormValue("ids")
				gsids := strings.Split(ids, ",")

				ifEmpty := true
				for _, v := range gsids {
					if v != "" {
						ifEmpty = false
						ctx1 := &context.Context{}
						vars := url.Values{}
						vars.Add("cmd", "KickAllPlayer")
						vars.Add("gsId", v)
						err, res := fusion.CallToCenter(&vars, common.MyCfg.Api["GM2GS"], "GET", ctx1)
						if err != nil {
							faild = append(faild, "<h5>"+err.Error()+"</h5>")
						}
						if res != "" {
							if res == "All Is OK!" {
								faild = append(faild, "<h5>"+gotext.Get("已经踢出所有玩家")+"</h5>")
							} else {
								faild = append(faild, "<h5>"+":"+res+"</h5>")
							}
						}
					}
				}
				if ifEmpty {
					return false, gotext.Get("请选择服务器"), ""
				}
				return true, "", faild
			}))

	formList := tGameServers.GetForm()
	formList.AddField(gotext.Get("服务器Id"), "Id", db.Int, form.Number).
		FieldDisableWhenCreate().FieldDisplayButCanNotEditWhenUpdate().FieldMust()
	formList.AddField(gotext.Get("外部IP"), "externalIP", db.Char, form.Text).FieldMust()
	formList.AddField(gotext.Get("外部端口"), "externalPort", db.Smallint, form.Number).FieldMust()
	formList.AddField(gotext.Get("内部IP"), "internalIP", db.Char, form.Text).FieldMust()
	formList.AddField("InternalName", "internalName", db.Char, form.Text)
	//formList.AddField("LogicId", "logicId", db.Int, form.Number).FieldDefault("0")
	formList.AddField(gotext.Get("服务器名"), "logicName", db.Char, form.Text).FieldMust()
	formList.AddField(gotext.Get("服务器开启时间"), "logicOpenTime", db.Datetime, form.Datetime).
		FieldDisableWhenUpdate()
	formList.AddField(gotext.Get("服务器状态"), "logicOpenStatus", db.Tinyint, form.Number).
		FieldDisableWhenCreate().FieldMust()
	formList.AddField(gotext.Get("服务器标签"), "logicSpecialFlags", db.Int, form.Number).
		FieldDisableWhenCreate().FieldMust()

	formList.SetTable("t_game_servers").SetTitle("TGameServers").SetDescription("TGameServers")

	formList.SetInsertFn(func(values form2.Values) error {
		serverInfo := def.GameServerInfo{}
		serverInfo.ExternalIP = values.Get("externalIP")
		temp, _ := strconv.Atoi(values.Get("externalPort"))
		serverInfo.ExternalPort = uint16(temp)
		serverInfo.InternalIP = values.Get("internalIP")
		serverInfo.InternalName = values.Get("internalName")
		if serverInfo.InternalName == "" {
			serverInfo.InternalName = "s1"
		}
		serverInfo.LogicID = 0
		serverInfo.LogicName = values.Get("logicName")
		serverInfo.LogicOpenTime, _ = time.ParseInLocation("2006-01-02 15:04:05", values.Get("logicOpenTime"), time.Local)
		serverInfo.LogicOpenStatus = def.ServerOpenStatus_Hide
		serverInfo.LogicSpecialFlags = def.ServerSpecialFlag_New

		db_global := fusion.GetBaseGormDB("db_global")
		db_global.Create(&serverInfo)

		vars := url.Values{}
		ctx1 := &context.Context{}
		err, res := fusion.CallToCenter(&vars, common.MyCfg.Api["reloadGameServerInfos"], "GET", ctx1)
		if err != nil {
			str := gotext.Get("操作失败")
			return errors.New(str)
		}
		if res != "" {
			if res != "All Is OK!" {
				return errors.New(res)
			}
		}
		return nil
	})

	formList.SetUpdateFn(func(values form2.Values) error {
		serverInfo := def.GameServerInfo{}
		temp, _ := strconv.Atoi(values.Get("Id"))
		serverInfo.ID = uint32(temp)
		serverInfo.ExternalIP = values.Get("externalIP")
		temp, _ = strconv.Atoi(values.Get("externalPort"))
		serverInfo.ExternalPort = uint16(temp)
		serverInfo.InternalIP = values.Get("internalIP")
		serverInfo.InternalName = values.Get("internalName")
		if serverInfo.InternalName == "" {
			serverInfo.InternalName = "s1"
		}
		serverInfo.LogicID = 0
		serverInfo.LogicName = values.Get("logicName")
		serverInfo.LogicOpenTime, _ = time.ParseInLocation("2006-01-02 15:04:05", values.Get("logicOpenTime"), time.Local)
		temp, _ = strconv.Atoi(values.Get("logicOpenStatus"))
		serverInfo.LogicOpenStatus = uint8(temp)
		temp, _ = strconv.Atoi(values.Get("logicSpecialFlags"))
		serverInfo.LogicSpecialFlags = uint32(temp)

		db_global := fusion.GetBaseGormDB("db_global")
		db_global.UpdateColumns(&serverInfo)

		vars := url.Values{}
		ctx1 := &context.Context{}
		err, res := fusion.CallToCenter(&vars, common.MyCfg.Api["reloadGameServerInfos"], "GET", ctx1)
		if err != nil {
			str := gotext.Get("操作失败")
			return errors.New(str)
		}
		if res != "" {
			if res != "All Is OK!" {
				return errors.New(res)
			}
		}
		return nil
	})

	return tGameServers
}

func HotfixServer(ctx *context.Context) {
	var faild string
	ids := ctx.Request.Form["gsIdX[]"]
	fileName := ctx.Request.FormValue("fileName")
	hotfixArgs := ctx.Request.FormValue("hotfixArgs")

	ifEmpty := true
	for _, v := range ids {
		if v != "" {
			ifEmpty = false
			ctx1 := &context.Context{}
			vars := url.Values{}
			vars.Add("fileName", fileName)
			var datas = make(map[string]interface{})
			go_manager := fusion.GetBaseGormDB("go_manager")
			go_manager.Table("databasecfg").Select("redisSn,deployHost,deployPort").Where("serverId = ?", v).Scan(&datas)

			err, res := fusion.CallToDeploy(&vars, common.MyCfg.Api["downHotfixFiles"], "GET", ctx1,
				fusion.Strval(datas["deployHost"]), fusion.Strval(datas["deployPort"]))
			if err != nil {
				faild += v + err.Error() + "/n"
			}
			if res != "" {
				if res == "All Is OK!" {
					faild += v + ":热更新文件已下载" + "/n"
				} else {
					faild += v + ":" + res + "/n"
				}
			}

			ctx1 = &context.Context{}
			vars = url.Values{}
			vars.Add("cmd", "HotfixServer")
			vars.Add("args", hotfixArgs)
			vars.Add("gsId", v)
			var api string

			if find := strings.Contains(hotfixArgs, "MapServer"); find {
				api = common.MyCfg.Api["GM2MS"]
			} else if find = strings.Contains(hotfixArgs, "GameServer"); find {
				api = common.MyCfg.Api["GM2GS"]
			} else if find = strings.Contains(hotfixArgs, "GateServer"); find {
				api = common.MyCfg.Api["GM2GATE"]
			} else if find = strings.Contains(hotfixArgs, "SocialServer"); find {
				api = common.MyCfg.Api["GM2SOCIAL"]
			} else if find = strings.Contains(hotfixArgs, "DBPServer"); find {
				api = common.MyCfg.Api["GM2DBP"]
			}

			if api == "" {
				response.Error(ctx, gotext.Get("参数错误"))
			}
			err, res = fusion.CallToCenter(&vars, api, "GET", ctx1)
			if err != nil {
				faild += v + err.Error() + "/n"
			}
			if res != "" {
				if res == "All Is OK!" {
					faild += gotext.Get("该服务器已经通知热更") + "/n"
				} else {
					faild += v + ":" + res + "/n"
				}
			}
		}
	}

	if ifEmpty {
		response.Error(ctx, gotext.Get("请选择服务器"))
		return
	}
	response.OkWithMsg(ctx, faild)
}

var ChangeSelectBoxHeight = "<script>a= document.querySelectorAll('select.form-control');a[3].style.height='400px';a[4].style.height='400px';</script>"

func GetOperationServerPro(ctx *context.Context) table.Table {

	tGameServers := table.NewDefaultTable(table.Config{
		Driver:     "mysql",
		Connection: "db_global",
		CanAdd:     true,
		Editable:   true,
		Deletable:  false,
		Exportable: false,
		PrimaryKey: table.PrimaryKey{
			Type: db.Int,
			Name: "Id",
		},
	})
	open := def.GetServerStatusOP(1)
	flag := def.GetServerStatusOP(2)
	openString := def.GetServerOpenStatusOpString()
	flagString := def.GetServerSpecialOpString()
	info := tGameServers.GetInfo().HideDetailButton().HideDeleteButton()
	myOP := fusion.GetServerList()
	info.AddField("Id", "Id", db.Int).
		FieldSortable().
		FieldFilterable(types.FilterType{Options: myOP, FormType: form.SelectSingle}).
		FieldFilterOptions(myOP)
	info.AddField("服务器名", "logicName", db.Char).
		FieldSortable()
	info.AddField("服务器开启状态", "logicOpenStatus", db.Tinyint).
		FieldFilterable(types.FilterType{
			FormType: form.Select,
			Options:  open}).
		FieldSortable().
		FieldDisplay(func(model types.FieldModel) interface{} {
			temp, _ := strconv.Atoi(model.Value)
			return openString[temp]
		})
	info.AddField("服务器标签", "logicSpecialFlags", db.Int).
		FieldFilterable(types.FilterType{
			FormType: form.Select,
			Options:  flag}).
		FieldSortable().
		FieldDisplay(func(model types.FieldModel) interface{} {
			temp, _ := strconv.Atoi(model.Value)
			return flagString[temp]
		})
	info.AddField("开服时间", "logicOpenTime", db.Datetime).
		FieldEditAble(table2.Datetime).SetPageSizeList([]int{5, 10, 20, 100, 9999})

	info.AddActionButton(template.HTML("修改后台服务器读取数据库配置"), action.Jump("/admin/EditServerDatabaseCfg?Id={{.Id}}"+
		"&logicName={{(index .Value \"logicName\").Value}}"))

	info.AddActionButton(template.HTML("修改备份服务器读取数据库配置"), action.Jump("/admin/EditServerDatabaseBackupCfg?Id={{.Id}}"+
		"&logicName={{(index .Value \"logicName\").Value}}"))

	info.SetGetDataFn(
		func(param parameter.Parameters) (data []map[string]interface{}, size int) {
			db_global := fusion.GetBaseGormDB("db_global")
			if param.Fields["__pk"] != nil && param.Fields["__pk"][0] != "" {
				db_global.Raw("SELECT *FROM t_game_servers WHERE id = ?", param.Fields["__pk"][0]).
					Scan(&data)
				return data, 1
			} else {
				db_global.Raw("SELECT *FROM t_game_servers WHERE id IN (SELECT MIN(id) FROM t_game_servers GROUP BY externalIP,externalPort)").
					Scan(&data)
				min := math.Min(float64(param.PageInt*param.PageSizeInt), float64(len(data)))
				return data[((param.PageInt - 1) * param.PageSizeInt):int(min)], len(data)
			}
		})

	//info.AddButton(template.HTML("清理数据库"), icon.Save,
	//	action.PopUp("/admin/ClearDB", "操作结果", func(ctx *context.Context) (success bool, msg string, data interface{}) {
	//		faild := make([]string, 0)
	//		ids := ctx.Request.FormValue("ids")
	//		gsids := strings.Split(ids, ",")
	//		ifEmpty := true
	//		sqlSource := "SELECT CONCAT('TRUNCATE ', table_name, ';')\n" +
	//			"FROM information_schema.tables\n" +
	//			"WHERE table_schema = '%s';\n"
	//		deleteSqlSource := "SELECT CONCAT('DROP TABLE IF EXISTS ', table_name, ';')\n" +
	//			"FROM information_schema.tables\n" +
	//			"WHERE table_schema = '%s' and table_name like \"%%_2%%\";\n"
	//
	//		for _, v := range gsids {
	//			if v != "" {
	//				ifEmpty = false
	//				id, _ := strconv.Atoi(v)
	//				var dbName string
	//				var sqls = make([]string, 0)
	//
	//				db_char := fusion.GetServerGormDB("db_char", int64(id))
	//				if db_char == nil {
	//					faild = append(faild, "<h5>"+v+": db_char is nil, please check database config"+"</h5>")
	//				} else {
	//					db_char.Raw("select database ()").Scan(&dbName)
	//					sql := fmt.Sprintf(sqlSource, dbName)
	//					db_char.Raw(sql).Scan(&sqls)
	//					for _, truncateSql := range sqls {
	//						db_char.Exec(truncateSql)
	//					}
	//					faild = append(faild, "<h5>"+v+":db_char Clear Success"+"</h5>")
	//				}
	//
	//				db_log := fusion.GetServerGormDB("db_log", int64(id))
	//				if db_log == nil {
	//					faild = append(faild, "<h5>"+v+": db_log is nil, please check database config"+"</h5>")
	//				} else {
	//					dbName = ""
	//					db_log.Raw("select database ()").Scan(&dbName)
	//					deleteSql := fmt.Sprintf(deleteSqlSource, dbName)
	//					sqls = make([]string, 0)
	//					db_log.Raw(deleteSql).Scan(&sqls)
	//					for _, temp := range sqls {
	//						db_log.Exec(temp)
	//					}
	//
	//					sql := fmt.Sprintf(sqlSource, dbName)
	//					sqls = make([]string, 0)
	//					db_log.Raw(sql).Scan(&sqls)
	//					for _, truncateSql := range sqls {
	//						db_log.Exec(truncateSql)
	//					}
	//					faild = append(faild, "<h5>"+v+":db_log Clear Success"+"</h5>")
	//				}
	//
	//				db_global := fusion.GetBaseGormDB("db_global")
	//				if db_global == nil {
	//					faild = append(faild, "<h5>"+v+": db_global is nil, please check database config"+"</h5>")
	//				} else {
	//					db_global.Table("t_account_characters").Where("groupId = ?", v).Delete("t_account_characters")
	//					db_global.Table("t_guilds").Where("groupId = ?", v).Delete("t_guilds")
	//					faild = append(faild, "<h5>"+v+":db_global :t_account_characters,t_guilds Clear Success"+"</h5>")
	//				}
	//				var datas = make(map[string]interface{})
	//				go_manager := fusion.GetBaseGormDB("go_manager")
	//				go_manager.Table("databasecfg").Select("redisSn,deployHost,deployPort").Where("serverId = ?", v).Scan(&datas)
	//
	//				vars := url.Values{}
	//				ctx1 := &context.Context{}
	//				vars.Add("redisSn", fusion.Strval(datas["redisSn"]))
	//				err, res := fusion.CallToDeploy(&vars, common.MyCfg.Api["CleanRedis"], "GET", ctx1, fusion.Strval(datas["deployHost"]), fusion.Strval(datas["deployPort"]))
	//				if err != nil {
	//					faild = append(faild, "<h5>"+v+err.Error()+"</h5>")
	//				}
	//				if res != "" {
	//					if res == "All Is OK!" {
	//						faild = append(faild, "<h5>"+v+":Redis Clear Success"+"</h5>")
	//					}
	//					faild = append(faild, "<h5>"+v+":"+res+"</h5>")
	//				}
	//
	//			}
	//		}
	//		if ifEmpty {
	//			return false, "请选择服务器", ""
	//		}
	//		return true, "", faild
	//	}))

	info.AddButton(template.HTML("开启进程守护"), "", action.PopUp("/server/startProcessPro", "操作结果",
		func(ctx *context.Context) (success bool, msg string, data interface{}) {
			faild := []string{}
			ids := ctx.Request.FormValue("ids")
			gsids := strings.Split(ids, ",")
			vars := url.Values{}
			vars.Add("mode", "0")

			ifEmpty := true
			for _, v := range gsids {
				if v != "" {
					ifEmpty = false
					ctx1 := &context.Context{}
					var datas = make(map[string]interface{})
					go_manager := fusion.GetBaseGormDB("go_manager")
					go_manager.Table("databasecfg").Select("redisSn,deployHost,deployPort").Where("serverId = ?", v).Scan(&datas)

					err, res := fusion.CallToDeploy(&vars, common.MyCfg.Api["ChangeWatchMode"], "GET", ctx1,
						fusion.Strval(datas["deployHost"]), fusion.Strval(datas["deployPort"]))
					if err != nil {
						faild = append(faild, "<h5>"+v+err.Error()+"</h5>")
					}
					if res != "" {
						if res == "All Is OK!" {
							faild = append(faild, "<h5>"+v+":process watch start"+"</h5>")
						}
						faild = append(faild, "<h5>"+v+":"+res+"</h5>")
					}
				}
			}

			if ifEmpty {
				return false, "请选择服务器", ""
			}
			return true, "", faild
		}))

	info.AddButton(template.HTML("向玩家开放"), "",
		action.PopUp("/admin/CloseServer", "操作结果", func(ctx *context.Context) (success bool, msg string, data interface{}) {
			faild := []string{}
			ids := ctx.Request.FormValue("ids")
			gsids := strings.Split(ids, ",")
			ifEmpty := true
			for _, v := range gsids {
				if v != "" {
					ifEmpty = false
					db_global := fusion.GetBaseGormDB("db_global")
					db_global.Table("t_game_servers").Where("Id", v).Update("logicOpenStatus", def.ServerOpenStatus_Normal)
				}
			}
			ctx1 := &context.Context{}
			vars := url.Values{}
			err, res := fusion.CallToCenter(&vars, common.MyCfg.Api["reloadGameServerInfos"], "GET", ctx1)
			if err != nil {
				faild = append(faild, "<h5>"+err.Error()+"</h5>")
			}
			if res != "" {
				if res == "All Is OK!" {
					faild = append(faild, "<h5>"+"已将选中服务器设为开放状态"+"</h5>")
				} else {
					faild = append(faild, "<h5>"+":"+res+"</h5>")
				}
			}
			if ifEmpty {
				return false, "请选择服务器", ""
			}
			return true, "", faild
		}))

	info.AddButton(template.HTML("部署资源文件"), "", action.PopUpWithForm(action.PopUpData{
		Id:     "/server/deployResFilesPro",
		Title:  "部署资源文件",
		Width:  "1400px",
		Height: "830px",
	}, func(panel *types.FormPanel) *types.FormPanel {
		panel.AddField(gotext.Get(""), "operation", db.Int, form.Text).FieldHide().FieldDefault("3")
		panel.AddField("服务器ID", "gsIdX", db.Int, form.SelectBox).
			FieldOptions(myOP).FieldFoot(template.HTML(ChangeSelectBoxHeight))
		panel.AddField("分组名", "groupName", db.Varchar, form.Text)
		panel.EnableAjax()
		return panel
	}, "/admin/popup/OperationResFilesPro"))

	info.AddButton(template.HTML("下载资源文件"), "", action.PopUpWithForm(action.PopUpData{
		Id:     "/server/downResFilesPro",
		Title:  "下载资源文件",
		Width:  "1400px",
		Height: "830px",
	}, func(panel *types.FormPanel) *types.FormPanel {
		panel.AddField(gotext.Get(""), "operation", db.Int, form.Text).FieldHide().FieldDefault("2")
		panel.AddField("服务器ID", "gsIdX", db.Int, form.SelectBox).
			FieldOptions(myOP).FieldFoot(template.HTML(ChangeSelectBoxHeight))
		panel.AddField("分组名", "groupName", db.Varchar, form.Text)
		panel.EnableAjax()
		return panel
	}, "/admin/popup/OperationResFilesPro"))

	info.AddButton(template.HTML("清除资源文件"), "", action.PopUpWithForm(action.PopUpData{
		Id:     "/server/clearResFilePro",
		Title:  "清除资源文件",
		Width:  "1400px",
		Height: "830px",
	}, func(panel *types.FormPanel) *types.FormPanel {
		panel.AddField(gotext.Get(""), "operation", db.Int, form.Text).FieldHide().FieldDefault("1")
		panel.AddField("服务器ID", "gsIdX", db.Int, form.SelectBox).
			FieldOptions(myOP).FieldFoot(template.HTML(ChangeSelectBoxHeight))
		panel.AddField("分组名", "groupName", db.Varchar, form.Text)
		panel.EnableAjax()
		return panel
	}, "/admin/popup/OperationResFilesPro"))

	info.AddButton(template.HTML("下载热更新文件"), "", action.PopUpWithForm(action.PopUpData{
		Id:     "/server/downHotfixFilesPro",
		Title:  "下载热更新文件",
		Width:  "1400px",
		Height: "830px",
	}, func(panel *types.FormPanel) *types.FormPanel {
		panel.AddField(gotext.Get(""), "operation", db.Int, form.Text).FieldHide().FieldDefault("5")
		panel.AddField("服务器ID", "gsIdX", db.Int, form.SelectBox).
			FieldOptions(myOP).FieldFoot(template.HTML(ChangeSelectBoxHeight))
		panel.AddField("分组名", "groupName", db.Varchar, form.Text)
		panel.AddField("文件名", "fileName", db.Varchar, form.Text)
		panel.EnableAjax()
		return panel
	}, "/admin/popup/OperationResFilesPro"))

	info.AddButton(template.HTML("部署热更新文件"), "", action.PopUpWithForm(action.PopUpData{
		Id:     "/server/deployHotfixFilesPro",
		Title:  "部署热更新文件",
		Width:  "1400px",
		Height: "830px",
	}, func(panel *types.FormPanel) *types.FormPanel {
		panel.AddField(gotext.Get(""), "operation", db.Int, form.Text).FieldHide().FieldDefault("6")
		panel.AddField("服务器ID", "gsIdX", db.Int, form.SelectBox).
			FieldOptions(myOP).FieldFoot(template.HTML(ChangeSelectBoxHeight))
		panel.AddField("分组名", "groupName", db.Varchar, form.Text)
		panel.AddField("文件名", "fileName", db.Varchar, form.Text)
		panel.EnableAjax()
		return panel
	}, "/admin/popup/OperationResFilesPro"))

	info.AddButton(template.HTML("校验热更新文件Hash"), "", action.PopUpWithForm(action.PopUpData{
		Id:     "/server/checkHotfixFilesPro",
		Title:  "下载热更新文件",
		Width:  "1400px",
		Height: "830px",
	}, func(panel *types.FormPanel) *types.FormPanel {
		panel.AddField(gotext.Get(""), "operation", db.Int, form.Text).FieldHide().FieldDefault("7")
		panel.AddField("服务器ID", "gsIdX", db.Int, form.SelectBox).FieldFoot(template.HTML(ChangeSelectBoxHeight)).
			FieldOptions(myOP)
		panel.AddField("分组名", "groupName", db.Varchar, form.Text)
		panel.AddField("文件名", "fileName", db.Varchar, form.Text)
		panel.AddField("文件hash值", "fileHash", db.Varchar, form.Text)
		panel.EnableAjax()
		return panel
	}, "/admin/popup/OperationResFilesPro"))

	info.AddButton(template.HTML("导出配置文件"), "", action.PopUpWithForm(action.PopUpData{
		Id:     "/server/exportConfigFilePro",
		Title:  "下载热更新文件",
		Width:  "1400px",
		Height: "830px",
	}, func(panel *types.FormPanel) *types.FormPanel {
		panel.AddField(gotext.Get(""), "operation", db.Int, form.Text).FieldHide().FieldDefault("8")
		panel.AddField("服务器ID", "gsIdX", db.Int, form.SelectBox).FieldFoot(template.HTML(ChangeSelectBoxHeight)).
			FieldOptions(myOP)
		panel.AddField("分组名", "groupName", db.Varchar, form.Text)
		panel.EnableAjax()
		return panel
	}, "/admin/popup/OperationResFilesPro"))

	info.AddButton(template.HTML("检查是否为跨服"), "", action.PopUpWithForm(action.PopUpData{
		Id:     "/server/isCrossServerPro",
		Title:  "检查是否为跨服",
		Width:  "1400px",
		Height: "830px",
	}, func(panel *types.FormPanel) *types.FormPanel {
		panel.AddField(gotext.Get(""), "operation", db.Int, form.Text).FieldHide().FieldDefault("9")
		panel.AddField("服务器ID", "gsIdX", db.Int, form.SelectBox).
			FieldOptions(myOP).FieldFoot(template.HTML(ChangeSelectBoxHeight))
		panel.AddField("分组名", "groupName", db.Varchar, form.Text)
		panel.EnableAjax()
		return panel
	}, "/admin/popup/OperationResFilesPro"))

	info.AddButton(template.HTML("检查进程数量"), "", action.PopUpWithForm(action.PopUpData{
		Id:     "/server/CheckProcessPro",
		Title:  "检查进程数量",
		Width:  "1400px",
		Height: "830px",
	}, func(panel *types.FormPanel) *types.FormPanel {
		panel.AddField(gotext.Get(""), "operation", db.Int, form.Text).FieldHide().FieldDefault("10")
		panel.AddField("服务器ID", "gsIdX", db.Int, form.SelectBox).
			FieldOptions(myOP).FieldFoot(template.HTML(ChangeSelectBoxHeight))
		panel.AddField("分组名", "groupName", db.Varchar, form.Text)
		panel.EnableAjax()
		return panel
	}, "/admin/popup/OperationResFilesPro"))

	info.AddButton(template.HTML("重新启动服务器"), "", action.PopUpWithForm(action.PopUpData{
		Id:     "/server/RestartServerPro",
		Title:  "重新启动服务器",
		Width:  "1400px",
		Height: "830px",
	}, func(panel *types.FormPanel) *types.FormPanel {
		panel.AddField(gotext.Get(""), "operation", db.Int, form.Text).FieldHide().FieldDefault("4")
		panel.AddField("服务器ID", "gsIdX", db.Int, form.SelectBox).
			FieldOptions(myOP).FieldFoot(template.HTML(ChangeSelectBoxHeight))
		panel.AddField("分组名", "groupName", db.Varchar, form.Text)
		panel.EnableAjax()
		return panel
	}, "/admin/popup/OperationResFilesPro"))

	info.AddButton(template.HTML("强制闭服务器（kill）"), "", action.PopUpWithForm(action.PopUpData{
		Id:     "/server/StopServerPro",
		Title:  "强制闭服务器（kill）",
		Width:  "1400px",
		Height: "830px",
	}, func(panel *types.FormPanel) *types.FormPanel {
		panel.AddField(gotext.Get(""), "operation", db.Int, form.Text).FieldHide().FieldDefault("11")
		panel.AddField("服务器ID", "gsIdX", db.Int, form.SelectBox).
			FieldOptions(myOP).FieldFoot(template.HTML(ChangeSelectBoxHeight))
		panel.AddField("分组名", "groupName", db.Varchar, form.Text)
		panel.EnableAjax()
		return panel
	}, "/admin/popup/OperationResFilesPro"))

	info.AddButton(template.HTML("运行命令"), "", action.PopUpWithForm(action.PopUpData{
		Id:     "/server/RunBashPro",
		Title:  "运行命令",
		Width:  "1400px",
		Height: "830px",
	}, func(panel *types.FormPanel) *types.FormPanel {
		myOP := fusion.GetServerList()
		panel.AddField(gotext.Get(""), "operation", db.Int, form.Text).FieldHide().FieldDefault("12")
		panel.AddField("服务器ID", "gsIdX", db.Int, form.SelectBox).
			FieldOptions(myOP).FieldFoot(template.HTML(ChangeSelectBoxHeight))
		panel.AddField("分组名", "groupName", db.Varchar, form.Text)
		panel.AddField("执行命令", "cmdline", db.Varchar, form.Text)
		panel.EnableAjax()
		return panel
	}, "/admin/popup/OperationResFilesPro"))

	info.AddButton(template.HTML("重启deploy"), "", action.PopUpWithForm(action.PopUpData{
		Id:     "/server/RebootSelfPro",
		Title:  "重启deploy",
		Width:  "1400px",
		Height: "830px",
	}, func(panel *types.FormPanel) *types.FormPanel {
		panel.AddField(gotext.Get(""), "operation", db.Int, form.Text).FieldHide().FieldDefault("13")
		panel.AddField("服务器ID", "gsIdX", db.Int, form.SelectBox).FieldFoot(template.HTML(ChangeSelectBoxHeight)).
			FieldOptions(myOP)
		panel.AddField("分组名", "groupName", db.Varchar, form.Text)
		panel.EnableAjax()
		return panel
	}, "/admin/popup/OperationResFilesPro"))

	info.AddButton(template.HTML("deploy重载配置文件"), "", action.PopUpWithForm(action.PopUpData{
		Id:     "/server/ReloadConfigPro",
		Title:  "deploy重载配置文件",
		Width:  "1400px",
		Height: "830px",
	}, func(panel *types.FormPanel) *types.FormPanel {
		panel.AddField(gotext.Get(""), "operation", db.Int, form.Text).FieldHide().FieldDefault("14")
		panel.AddField("服务器ID", "gsIdX", db.Int, form.SelectBox).
			FieldOptions(myOP).FieldFoot(template.HTML(ChangeSelectBoxHeight))
		panel.AddField("分组名", "groupName", db.Varchar, form.Text)
		panel.EnableAjax()
		return panel
	}, "/admin/popup/OperationResFilesPro"))

	info.AddButton(template.HTML("手动保存日志文件"), "", action.PopUp("/server/persistLogFiles", "操作结果",
		func(ctx *context.Context) (success bool, msg string, data interface{}) {
			faild := []string{}
			ids := ctx.Request.FormValue("ids")
			gsids := strings.Split(ids, ",")
			vars := url.Values{}

			ifEmpty := true
			for _, v := range gsids {
				if v != "" {
					ifEmpty = false
					ctx1 := &context.Context{}
					var datas = make(map[string]interface{})
					go_manager := fusion.GetBaseGormDB("go_manager")
					go_manager.Table("databasecfg").Select("redisSn,deployHost,deployPort").Where("serverId = ?", v).Scan(&datas)

					err, res := fusion.CallToDeploy(&vars, common.MyCfg.Api["persistLogFiles"], "GET", ctx1,
						fusion.Strval(datas["deployHost"]), fusion.Strval(datas["deployPort"]))
					if err != nil {
						faild = append(faild, "<h5>"+v+err.Error()+"</h5>")
					}
					if res != "" {
						if res == "All Is OK!" {
							faild = append(faild, "<h5>"+v+":日志文件已保存"+"</h5>")
						} else {
							faild = append(faild, "<h5>"+v+":"+res+"</h5>")
						}
					}
				}
			}
			if ifEmpty {
				return false, "请选择服务器", ""
			}
			return true, "", faild
		}))

	// info.AddButton(template.HTML("开启服务器(更新后)"), "", action.PopUp("/server/startServer4UpdatePro", "操作结果",
	// 	func(ctx *context.Context) (success bool, msg string, data interface{}) {
	// 		faild := []string{}
	// 		ids := ctx.Request.FormValue("ids")
	// 		gsids := strings.Split(ids, ",")
	// 		vars := url.Values{}
	// 		vars.Add("timeout", "10")

	// 		ifEmpty := true
	// 		for _, v := range gsids {
	// 			if v != "" {
	// 				ifEmpty = false
	// 				ctx1 := &context.Context{}
	// 				var datas = make(map[string]interface{})
	// 				go_manager := fusion.GetBaseGormDB("go_manager")
	// 				go_manager.Table("databasecfg").Select("redisSn,deployHost,deployPort").Where("serverId = ?", v).Scan(&datas)

	// 				err, res := fusion.CallToDeploy(&vars, common.MyCfg.Api["RestartServer"], "GET", ctx1,
	// 					fusion.Strval(datas["deployHost"]), fusion.Strval(datas["deployPort"]))
	// 				if err != nil {
	// 					faild = append(faild, "<h5>"+v+err.Error()+"</h5>")
	// 				}
	// 				if res != "" {
	// 					if res == "All Is OK!" {
	// 						faild = append(faild, "<h5>"+v+":process watch start"+"</h5>")
	// 					} else {
	// 						faild = append(faild, "<h5>"+v+":"+res+"</h5>")
	// 					}
	// 				}
	// 			}
	// 		}

	// 		if ifEmpty {
	// 			return false, "请选择服务器", ""
	// 		}
	// 		return true, "", faild
	// 	}))

	info.AddButton(template.HTML("关闭选中服务器"), "",
		action.PopUp("/admin/CloseServerPro", "操作结果", func(ctx *context.Context) (success bool, msg string, data interface{}) {
			faild := []string{}
			ids := ctx.Request.FormValue("ids")
			gsids := strings.Split(ids, ",")
			ifEmpty := true
			for _, v := range gsids {
				if v != "" {
					ifEmpty = false
					vars := &url.Values{}
					vars.Add("name", "AdminServer")
					vars.Add("cmd", "ShutdownServer")
					vars.Add("gsId", v)
					ctx1 := &context.Context{}
					err, res := fusion.CallToCenter(vars, common.MyCfg.Api["GM2GS"], "GET", ctx1)
					if err != nil {
						faild = append(faild, "<h5>"+v+":"+err.Error()+"</h5>")
					}
					if res != "" {
						faild = append(faild, "<h5>"+res+"</h5>")
					}
				}
			}
			if ifEmpty {
				return false, "请选择服务器", ""
			}
			return true, "", faild
		}))

	info.AddButton(template.HTML("关闭进程守护"), "", action.PopUp("/server/CloseProcessWatchPro", "操作结果",
		func(ctx *context.Context) (success bool, msg string, data interface{}) {
			faild := []string{}
			ids := ctx.Request.FormValue("ids")
			gsids := strings.Split(ids, ",")
			vars := url.Values{}
			vars.Add("mode", "1")

			ifEmpty := true
			for _, v := range gsids {
				if v != "" {
					ifEmpty = false
					ctx1 := &context.Context{}
					var datas = make(map[string]interface{})
					go_manager := fusion.GetBaseGormDB("go_manager")
					go_manager.Table("databasecfg").Select("redisSn,deployHost,deployPort").Where("serverId = ?", v).Scan(&datas)

					err, res := fusion.CallToDeploy(&vars, common.MyCfg.Api["ChangeWatchMode"], "GET", ctx1,
						fusion.Strval(datas["deployHost"]), fusion.Strval(datas["deployPort"]))
					if err != nil {
						faild = append(faild, "<h5>"+v+err.Error()+"</h5>")
					}
					if res != "" {
						if res == "All Is OK!" {
							faild = append(faild, "<h5>"+v+":process watch closed"+"</h5>")
						}
						faild = append(faild, "<h5>"+v+":"+res+"</h5>")
					}
				}
			}

			if ifEmpty {
				return false, "请选择服务器", ""
			}
			return true, "", faild
		}))

	info.AddButton(template.HTML("查看进程守护状态"), "", action.PopUp("/server/CheckWatchModePro", "操作结果",
		func(ctx *context.Context) (success bool, msg string, data interface{}) {
			faild := []string{}
			ids := ctx.Request.FormValue("ids")
			gsids := strings.Split(ids, ",")
			vars := url.Values{}

			ifEmpty := true
			for _, v := range gsids {
				if v != "" {
					ifEmpty = false
					ctx1 := &context.Context{}
					var datas = make(map[string]interface{})
					go_manager := fusion.GetBaseGormDB("go_manager")
					go_manager.Table("databasecfg").Select("redisSn,deployHost,deployPort").Where("serverId = ?", v).Scan(&datas)

					err, res := fusion.CallToDeploy(&vars, common.MyCfg.Api["CheckWatchMode"], "GET", ctx1,
						fusion.Strval(datas["deployHost"]), fusion.Strval(datas["deployPort"]))
					if err != nil {
						faild = append(faild, "<h5>"+v+err.Error()+"</h5>")
					}
					if res != "" {
						faild = append(faild, "<h5>"+v+":"+"进程守护状态为"+res+"</h5>")
					}
				}
			}

			if ifEmpty {
				return false, "请选择服务器", ""
			}
			return true, "", faild
		}))

	info.AddButton(template.HTML("批量设置服务器标签"), "", action.PopUpWithForm(action.PopUpData{
		Id:     "/admin/popup/setServerSymbolPro",
		Title:  "批量设置服务器标签",
		Width:  "1400px",
		Height: "830px",
	}, func(panel *types.FormPanel) *types.FormPanel {
		panel.AddField("type", "OPType", db.Int, form.Text).
			FieldDefault("1").FieldHide()
		panel.AddField("服务器ID", "Id", db.Int, form.SelectBox).FieldFoot(template.HTML(ChangeSelectBoxHeight)).
			FieldOptions(myOP)
		panel.AddField("服务器标签", "logicSpecialFlags", db.Int, form.SelectSingle).
			FieldOptions(flag)
		panel.EnableAjax("操作成功", "操作失败，请联系管理员", "/admin/info/t_game_servers")
		return panel
	}, "/admin/popup/UpdateServerSettingPro"))

	info.AddButton(template.HTML("批量设置服务器状态"), "", action.PopUpWithForm(action.PopUpData{
		Id:     "/admin/popup/setServerStatusPro",
		Title:  "批量设置服务器状态",
		Width:  "1400px",
		Height: "830px",
	}, func(panel *types.FormPanel) *types.FormPanel {
		panel.AddField("type", "OPType", db.Int, form.Text).
			FieldDefault("2").FieldHide()
		panel.AddField("服务器ID", "Id", db.Int, form.SelectBox).
			FieldOptions(myOP).FieldFoot(template.HTML(ChangeSelectBoxHeight))
		panel.AddField("服务器开启状态", "logicOpenStatus", db.Int, form.SelectSingle).
			FieldOptions(open)
		panel.EnableAjax("操作成功", "操作失败，请联系管理员", "/admin/info/t_game_servers")
		return panel
	}, "/admin/popup/UpdateServerSettingPro"))

	info.AddButton(template.HTML("踢出所有玩家"), "", action.PopUp("/server/KickAllPlayerPro", "操作结果",
		func(ctx *context.Context) (success bool, msg string, data interface{}) {
			faild := []string{}
			ids := ctx.Request.FormValue("ids")
			gsids := strings.Split(ids, ",")

			ifEmpty := true
			for _, v := range gsids {
				if v != "" {
					ifEmpty = false
					ctx1 := &context.Context{}
					vars := url.Values{}
					vars.Add("cmd", "KickAllPlayer")
					vars.Add("gsId", v)
					err, res := fusion.CallToCenter(&vars, common.MyCfg.Api["GM2GS"], "GET", ctx1)
					if err != nil {
						faild = append(faild, "<h5>"+err.Error()+"</h5>")
					}
					if res != "" {
						if res == "All Is OK!" {
							faild = append(faild, "<h5>"+"已经踢出所有玩家"+"</h5>")
						} else {
							faild = append(faild, "<h5>"+":"+res+"</h5>")
						}
					}
				}
			}
			if ifEmpty {
				return false, "请选择服务器", ""
			}
			return true, "", faild
		}))

	formList := tGameServers.GetForm()
	// formList.AddField("序号", "customOrder", db.Int, form.Number)
	formList.AddField("服务器Id", "Id", db.Int, form.Number).
		FieldDisableWhenCreate().FieldDisplayButCanNotEditWhenUpdate().FieldMust()
	formList.AddField("外部IP", "externalIP", db.Char, form.Text).FieldMust()
	formList.AddField("外部端口", "externalPort", db.Smallint, form.Number).FieldMust()
	formList.AddField("内部IP", "internalIP", db.Char, form.Text).FieldMust()
	formList.AddField("InternalName", "internalName", db.Char, form.Text)
	//formList.AddField("LogicId", "logicId", db.Int, form.Number).FieldDefault("0")
	formList.AddField("服务器名", "logicName", db.Char, form.Text).FieldMust()
	formList.AddField("服务器开启时间", "logicOpenTime", db.Datetime, form.Datetime).
		FieldDisableWhenUpdate()
	formList.AddField("服务器状态", "logicOpenStatus", db.Tinyint, form.Number).
		FieldDisableWhenCreate().FieldMust()
	formList.AddField("服务器标签", "logicSpecialFlags", db.Int, form.Number).
		FieldDisableWhenCreate().FieldMust()

	formList.SetTable("t_game_servers").SetTitle("TGameServers").SetDescription("TGameServers")

	formList.SetInsertFn(func(values form2.Values) error {
		serverInfo := def.GameServerInfo{}
		serverInfo.ExternalIP = values.Get("externalIP")
		temp, _ := strconv.Atoi(values.Get("externalPort"))
		serverInfo.ExternalPort = uint16(temp)
		temp, _ = strconv.Atoi(values.Get("customOrder"))
		serverInfo.CustomOrder = uint32(temp)
		serverInfo.InternalIP = values.Get("internalIP")
		serverInfo.InternalName = values.Get("internalName")
		if serverInfo.InternalName == "" {
			serverInfo.InternalName = "s1"
		}
		serverInfo.LogicID = 0
		serverInfo.LogicName = values.Get("logicName")
		serverInfo.LogicOpenTime, _ = time.ParseInLocation("2006-01-02 15:04:05", values.Get("logicOpenTime"), time.Local)
		serverInfo.LogicOpenStatus = def.ServerOpenStatus_Hide
		serverInfo.LogicSpecialFlags = def.ServerSpecialFlag_New

		db_global := fusion.GetBaseGormDB("db_global")
		db_global.Create(&serverInfo)

		vars := url.Values{}
		ctx1 := &context.Context{}
		err, res := fusion.CallToCenter(&vars, common.MyCfg.Api["reloadGameServerInfos"], "GET", ctx1)
		if err != nil {
			str := "操作失败"
			return errors.New(str)
		}
		if res != "" {
			if res != "All Is OK!" {
				return errors.New(res)
			}
		}
		return nil
	})

	formList.SetUpdateFn(func(values form2.Values) error {
		serverInfo := def.GameServerInfo{}
		temp, _ := strconv.Atoi(values.Get("Id"))
		serverInfo.ID = uint32(temp)
		temp, _ = strconv.Atoi(values.Get("customOrder"))
		serverInfo.CustomOrder = uint32(temp)
		serverInfo.ExternalIP = values.Get("externalIP")
		temp, _ = strconv.Atoi(values.Get("externalPort"))
		serverInfo.ExternalPort = uint16(temp)
		serverInfo.InternalIP = values.Get("internalIP")
		serverInfo.InternalName = values.Get("internalName")
		if serverInfo.InternalName == "" {
			serverInfo.InternalName = "s1"
		}
		serverInfo.LogicID = 0
		serverInfo.LogicName = values.Get("logicName")
		serverInfo.LogicOpenTime, _ = time.ParseInLocation("2006-01-02 15:04:05", values.Get("logicOpenTime"), time.Local)
		temp, _ = strconv.Atoi(values.Get("logicOpenStatus"))
		serverInfo.LogicOpenStatus = uint8(temp)
		temp, _ = strconv.Atoi(values.Get("logicSpecialFlags"))
		serverInfo.LogicSpecialFlags = uint32(temp)

		db_global := fusion.GetBaseGormDB("db_global")
		db_global.UpdateColumns(&serverInfo)

		vars := url.Values{}
		ctx1 := &context.Context{}
		err, res := fusion.CallToCenter(&vars, common.MyCfg.Api["reloadGameServerInfos"], "GET", ctx1)
		if err != nil {
			str := "操作失败"
			return errors.New(str)
		}
		if res != "" {
			if res != "All Is OK!" {
				return errors.New(res)
			}
		}
		return nil
	})

	return tGameServers
}

func UpdateServerSettingPro(ctx *context.Context) {
	OPType, _ := strconv.Atoi(ctx.Request.FormValue("OPType"))
	formValue := ctx.Request.Form["Id[]"]
	var LogicIds []int
	for _, v := range formValue {
		id, err := strconv.Atoi(v)
		if err != nil {
			continue
		}
		LogicIds = append(LogicIds, id)
	}
	db_global := fusion.GetBaseGormDB("db_global")
	// 合服修复：将选中服务器的同IP同端口（合服）服务器一并纳入更新范围
	if len(LogicIds) > 0 {
		var relatedIds []int
		db_global.Raw(`SELECT Id FROM t_game_servers WHERE CONCAT(externalIP, ':', externalPort) IN (SELECT CONCAT(externalIP, ':', externalPort) FROM t_game_servers WHERE Id IN ?)`, LogicIds).Scan(&relatedIds)
		LogicIds = relatedIds
	}
	if OPType == 1 {
		status, _ := strconv.Atoi(ctx.Request.FormValue("logicSpecialFlags"))
		db_global.Table("t_game_servers").
			Where("Id in ?", LogicIds).
			Update("logicSpecialFlags", status)
	} else if OPType == 2 {
		symbol, _ := strconv.Atoi(ctx.Request.FormValue("logicOpenStatus"))
		db_global.Table("t_game_servers").
			Where("Id in ?", LogicIds).
			Update("logicOpenStatus", symbol)
	}

	ctx1 := &context.Context{}
	vars := url.Values{}
	err, res := fusion.CallToCenter(&vars, common.MyCfg.Api["reloadGameServerInfos"], "GET", ctx1)
	if err != nil {
		response.Error(ctx, err.Error())
		return
	}
	if res != "" {
		if res == "All Is OK!" {
			response.OkWithMsg(ctx, "已同步标签/状态到服务器")
			return
		} else {
			response.Error(ctx, res)
			return
		}
	}
	response.Ok(ctx)
}
func OperationResFilesPro(ctx *context.Context) {
	var failed string = "<div style='height: 500px; overflow-y: scroll;'><table  class='table table-hover '  style='  max-width: 100%;font-size: 13px;' table-layout: auto;><tbody'><tr><th>serverId</th><th>command</th><th>result</th></tr>"
	ids := ctx.Request.Form["gsIdX[]"]
	operation, _ := strconv.Atoi(ctx.Request.FormValue("operation"))
	var api = ""
	var fileHash = ""
	var fileName = ""
	var fileNames = ""
	var cmdline = ""
	var Method = "GET"
	var isDown = false
	switch operation {
	case 1:
		api = "clearResFiles"
		break
	case 2:
		api = "downResFiles"
		isDown = true
		break
	case 3:
		api = "deployResFiles"
		break
	case 4:
		api = "RestartServer"
		break
	case 5:
		api = "downResFiles"
		fileNames = ctx.Request.FormValue("fileName")
		isDown = true
		break
	case 6:
		api = "deployHotfixFiles"
		fileNames = ctx.Request.FormValue("fileName")
		break
	case 7:
		api = "checkHotfixFiles"
		fileHash = ctx.Request.FormValue("fileHash")
		fileName = ctx.Request.FormValue("fileName")
		break
	case 8:
		api = "exportConfigFile"
		break
	case 9:
		api = "isCrossServer"
		break
	case 10:
		api = "CheckProcess"
		break
	case 11:
		api = "StopServer"
		break
	case 12:
		api = "RunBash"
		cmdline = ctx.Request.FormValue("cmdline")
		Method = "POST"
		break
	case 13:
		api = "RebootSelf"
		break
	case 14:
		api = "ReloadConfig"
		break
	default:
		response.Error(ctx, gotext.Get("操作类型有误"))
		return
	}
	ifEmpty := true
	var wg sync.WaitGroup
	for _, v := range ids {
		if v != "" {
				deployHost, deployPort, deployName := fusion.GetDeployTargetByServerId(v)
				ifEmpty = false
				ctx1 := context.Context{}
			vars := url.Values{}
			vars.Add("name", deployName)
			if fileName != "" {
				vars.Add("fileName", fileName)
			}
			if fileNames != "" {
				vars.Add("fileNames", fileNames)
			}

			if fileHash != "" {
				vars.Add("fileHash", fileHash)
			}

			if cmdline != "" {
				vars.Add("cmdline", cmdline)
			}
			wg.Add(1)
			id := v
			if !isDown {
				fusion.GoPool.Submit(func() {
					defer wg.Done()
					if deployHost == "" {
					failed += " <tr><th>" + id + "</th><th>" + api + "</th><th>deployHost 未配置</th></tr>"
					return
					}
					err, res := fusion.CallToDeploy1(vars, common.MyCfg.Api[api], Method, ctx1, deployHost, deployPort)
					if err != nil {
						failed += " <tr><th>" + id + "</th><th>" + api + "</th><th>" + err.Error() + "</th></tr>"
					} else {
						if res != "" {
							if res == "All Is OK!" {
								failed += "<tr><th> " + id + "</th><th>" + api + "</th><th>操作已完成</th></tr>"
							} else {
								failed += "<tr><th>" + id + "</th><th>" + api + " </th><th>" + res + "</th></tr>"
							}
						}
					}
				})
			} else {
				fusion.GoPool10.Submit(func() {
					defer wg.Done()
					err, res := fusion.CallToDeploy1(vars, common.MyCfg.Api[api], Method, ctx1,
						deployHost, deployPort)
					if err != nil {
						failed += " <tr><th>" + id + "</th><th>" + api + "</th><th>" + err.Error() + "</th></tr>"
					} else {
						if res != "" {
							if res == "All Is OK!" {
								failed += "<tr><th> " + id + "</th><th>" + api + "</th><th>操作已完成</th></tr>"
							} else {
								failed += "<tr><th>" + id + "</th><th>" + api + " </th><th>" + res + "</th></tr>"
							}
						}
					}
				})
			}

		}
	}
	wg.Wait()
	failed += "</tbody></table></div>"
	if ifEmpty {
		response.Error(ctx, gotext.Get("请选择服务器"))
		return
	}
	response.OkWithMsg(ctx, failed)
}

func HotfixDBTables(ctx *context.Context) {
	idStrs := ctx.Request.Form["gsIdX[]"]
	if len(idStrs) == 0 || idStrs[0] == "" {
		response.Error(ctx, gotext.Get("请选择服务器"))
		return
	}

	var gsIDs []uint32
	for _, s := range idStrs {
		n, err := strconv.ParseUint(strings.TrimSpace(s), 10, 32)
		if err != nil {
			continue
		}
		gsIDs = append(gsIDs, uint32(n))
	}
	if len(gsIDs) == 0 {
		response.Error(ctx, gotext.Get("请选择服务器"))
		return
	}

	tableNamesStr := ctx.Request.FormValue("tableNames")
	if tableNamesStr == "" {
		response.Error(ctx, gotext.Get("请输入数据表名"))
		return
	}
	tableNames := strings.FieldsFunc(tableNamesStr, func(r rune) bool { return r == ',' || r == ' ' || r == '\t' || r == '\n' || r == '\r' })

	// Resolve special table commands
	tableNamesForGM, specialCmds := resolveSpecialTableCmds(tableNames)

	var faild string
	for _, gsID := range gsIDs {
		faild += fmt.Sprintf("%d:\n", gsID)

		// Step 1: AdminServer HotfixDBFile
		gsIDStr := fmt.Sprint(gsID)
			deployHost, deployPort, deployName := fusion.GetDeployTargetByServerId(gsIDStr)
			cmdline := "./runXML.app xml2dbfile"
			vars := url.Values{}
			vars.Add("name", deployName)
			vars.Add("cmdline", cmdline)
			vars.Add("mode", "silent")
			log.Printf("[HotfixDBTables] gsId=%d: CallToDeploy %s:%s %s?%s",
				gsID, deployHost, deployPort, common.MyCfg.Api["RunBash"], vars.Encode())
			ctx1 := &context.Context{}
			err, res := fusion.CallToDeploy(&vars, common.MyCfg.Api["RunBash"], "GET", ctx1, deployHost, deployPort)
			if err != nil {
				faild += fmt.Sprintf("  xml2dbfile failed, %s\n", err.Error())
			} else if res != "" && res != "All Is OK!" {
				faild += fmt.Sprintf("  xml2dbfile, %s\n", res)
			}

			faild += sendGMCommandLine(gsID, "AdminServer", "HotfixDBFile", "")

		// Step 2: AdminServer HotfixDBTables
		faild += sendGMCommandLine(gsID, "AdminServer", "HotfixDBTables", strings.Join(tableNamesForGM, ","))

		// Steps 4-10: Conditional special table commands
		for _, cmd := range specialCmds {
			faild += sendGMCommandLine(gsID, cmd.server, cmd.cmd, "")
		}

		faild += fmt.Sprintf("  %s\n", gotext.Get("操作完成"))
	}

	response.OkWithMsg(ctx, faild)
}

type specialCmd struct {
	server string
	cmd    string
}

func resolveSpecialTableCmds(tableNames []string) ([]string, []specialCmd) {
	var cmds []specialCmd
	remaining := make([]string, 0, len(tableNames))

	hasSpell := false
	hasLoot := false
	for _, name := range tableNames {
		if name == "$Spell" || hotfix.IsSpellTable(name) {
			hasSpell = true
		} else if name == "$Loot" || hotfix.IsLootTable(name) {
			hasLoot = true
		} else {
			remaining = append(remaining, name)
		}
	}
	if hasSpell {
		cmds = append(cmds, specialCmd{"MapServer", "HotfixSpellRelation"})
		cmds = append(cmds, specialCmd{"MapServer", "ClearSpellEffectArgs"})
	}
	if hasLoot {
		cmds = append(cmds, specialCmd{"MapServer", "HotfixLootRelation"})
	}

	specialMap := map[string]specialCmd{
		"QuestPrototype":  {"MapServer", "HotfixQuests"},
		"ItemPrototype":   {"MapServer", "HotfixItemPrototypes"},
		"MapSpecial":      {"MapServer", "HotfixMapSpecial"},
		"operating_activities": {"SocialServer", "ReloadActivity"},
	}
	for i := 0; i < len(remaining); i++ {
		if cmd, ok := specialMap[remaining[i]]; ok {
			cmds = append(cmds, cmd)
			if remaining[i] == "operating_activities" {
				cmds = append(cmds, specialCmd{"MapServer", "ReloadActivity"})
			}
			remaining = append(remaining[:i], remaining[i+1:]...)
			i--
		}
	}
	return remaining, cmds
}

func sendGMCommandLine(gsID uint32, server, cmd, args string) string {
	apis := []string{
		common.MyCfg.Api["GM2S"],
		common.MyCfg.Api["GM2MS"],
		common.MyCfg.Api["GM2GS"],
		common.MyCfg.Api["GM2GATE"],
		common.MyCfg.Api["GM2SOCIAL"],
		common.MyCfg.Api["GM2DBP"],
	}
	var result string
	for _, api := range apis {
		vars := url.Values{}
		vars.Add("name", server)
		vars.Add("cmd", cmd)
		vars.Add("gsId", fmt.Sprint(gsID))
		if args != "" {
			vars.Add("args", args)
		}
		log.Printf("[HotfixDBTables] gsId=%d: CallToCenter %s:%s/%s?%s",
			gsID, common.MyCfg.Center.Host, common.MyCfg.Center.Port, api, vars.Encode())
		ctx1 := &context.Context{}
		err, res := fusion.CallToCenter(&vars, api, "GET", ctx1)
		if err != nil {
			result += fmt.Sprintf("  %s(%s) failed, %s\n", cmd, api, err.Error())
		} else if res != "" && res != "All Is OK!" {
			result += fmt.Sprintf("  %s(%s), %s\n", cmd, api, res)
		}
	}
	return result
}
