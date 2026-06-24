package AutoPages

import (
	"admin/common"
	"admin/common/def"
	"admin/fusion"
	"html/template"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/GoAdminGroup/go-admin/context"
	"github.com/GoAdminGroup/go-admin/modules/db"
	form2 "github.com/GoAdminGroup/go-admin/plugins/admin/modules/form"
	"github.com/GoAdminGroup/go-admin/plugins/admin/modules/table"
	"github.com/GoAdminGroup/go-admin/template/icon"
	"github.com/GoAdminGroup/go-admin/template/types"
	"github.com/GoAdminGroup/go-admin/template/types/action"
	"github.com/GoAdminGroup/go-admin/template/types/form"
	table2 "github.com/GoAdminGroup/go-admin/template/types/table"
	"github.com/leonelquinteros/gotext"
)

func GetTGameServersTable(ctx *context.Context) table.Table {

	tGameServers := table.NewDefaultTable(table.Config{
		Driver:     "mysql",
		Connection: "db_global",
		CanAdd:     true,
		Editable:   true,
		Deletable:  true,
		Exportable: true,
		PrimaryKey: table.PrimaryKey{
			Type: db.Int,
			Name: "Id",
		},
	})
	open := def.GetServerStatusOP(1)
	flag := def.GetServerStatusOP(2)
	openString := def.GetServerOpenStatusOpString()
	flagString := def.GetServerSpecialOpString()
	info := tGameServers.GetInfo()
	myOP := fusion.GetServerList()
	// info.AddField(gotext.Get("序号"), "customOrder", db.Int).
	// 	FieldSortable().FieldEditAble()

	info.AddField("Id", "Id", db.Int).
		FieldSortable().
		FieldFilterable(types.FilterType{Options: myOP, FormType: form.SelectSingle}).
		FieldFilterOptions(myOP)
	info.AddField("ExternalIP", "externalIP", db.Char).
		FieldHide()
	info.AddField("ExternalPort", "externalPort", db.Smallint).
		FieldHide()
	info.AddField("InternalIP", "internalIP", db.Char).
		FieldHide()
	info.AddField("InternalName", "internalName", db.Char).
		FieldSortable().FieldHide()
	info.AddField("LogicId", "logicId", db.Int).
		FieldSortable().FieldHide()
	info.AddField(gotext.Get("服务器名"), "logicName", db.Char).
		FieldSortable().FieldEditAble(table2.Text)
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

	info.AddButton(template.HTML(gotext.Get("中心服数据库备份")), "", action.PopUp("/server/backupGDB", gotext.Get("操作结果"),
		func(ctx *context.Context) (success bool, msg string, data interface{}) {
			faild := []string{}
			ctx1 := &context.Context{}
			vars := url.Values{}
			vars.Add("serverName", "Global")
			err, res := fusion.CallToBackup(&vars, common.MyCfg.Api["RequestBackUp"], "GET", ctx1,
				common.MyCfg.BackUp.Host, common.MyCfg.BackUp.Port)
			if err != nil {
				faild = append(faild, "<h5>"+err.Error()+"</h5>")
			}
			if res != "" {
				if res == "All Is OK!" {
					faild = append(faild, "<h5>"+"Database backup start"+"</h5>")
				} else {
					faild = append(faild, "<h5>"+res+"</h5>")
				}
			}
			return true, "", faild
		}))

	info.AddButton(template.HTML(gotext.Get("游戏服数据库备份")), "", action.PopUp("/server/backupGSDB", gotext.Get("操作结果"),
		func(ctx *context.Context) (success bool, msg string, data interface{}) {
			faild := []string{}
			ctx1 := &context.Context{}
			vars := url.Values{}
			vars.Add("gsIds", ctx.Request.FormValue("ids"))
			vars.Add("serverName", "GameServer")
			err, res := fusion.CallToBackup(&vars, common.MyCfg.Api["RequestBackUp"], "GET", ctx1,
				common.MyCfg.BackUp.Host, common.MyCfg.BackUp.Port)
			if err != nil {
				faild = append(faild, "<h5>"+err.Error()+"</h5>")
			}
			if res != "" {
				if res == "All Is OK!" {
					faild = append(faild, "<h5>"+"Database backup start"+"</h5>")
				} else {
					faild = append(faild, "<h5>"+res+"</h5>")
				}
			}
			return true, "", faild
		}))

	info.AddButton(template.HTML(gotext.Get("查看进程守护状态")), "", action.PopUp("/server/CheckWatchMode", gotext.Get("操作结果"),
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
						faild = append(faild, "<h5>"+v+":"+gotext.Get("进程守护状态为")+res+"</h5>")
					}
				}
			}

			if ifEmpty {
				return false, gotext.Get("请选择服务器"), ""
			}
			return true, "", faild
		}))

	info.AddButton(template.HTML(gotext.Get("批量设置服务器标签")), "", action.PopUpWithForm(action.PopUpData{
		Id:     "/admin/popup/setServerSymbol",
		Title:  gotext.Get("批量设置服务器标签"),
		Width:  "900px",
		Height: "430px",
	}, func(panel *types.FormPanel) *types.FormPanel {
		panel.AddField("type", "OPType", db.Int, form.Text).
			FieldDefault("1").FieldHide()
		panel.AddField(gotext.Get("服务器ID"), "Id", db.Int, form.SelectBox).FieldFoot(template.HTML(ChangeSelectBoxHeight)).
			FieldOptions(fusion.GetServerList())
		panel.AddField(gotext.Get("服务器标签"), "logicSpecialFlags", db.Int, form.SelectSingle).
			FieldOptions(flag)
		panel.EnableAjax(gotext.Get("操作成功"), gotext.Get("操作失败，请联系管理员"), "/admin/info/t_game_servers")
		return panel
	}, "/admin/popup/updateServerSetting"))

	info.AddButton(template.HTML(gotext.Get("批量设置服务器状态")), "", action.PopUpWithForm(action.PopUpData{
		Id:     "/admin/popup/setServerStatus",
		Title:  gotext.Get("批量设置服务器状态"),
		Width:  "900px",
		Height: "430px",
	}, func(panel *types.FormPanel) *types.FormPanel {
		panel.AddField("type", "OPType", db.Int, form.Text).
			FieldDefault("2").FieldHide()
		panel.AddField(gotext.Get("服务器ID"), "Id", db.Int, form.SelectBox).FieldFoot(template.HTML(ChangeSelectBoxHeight)).
			FieldOptions(fusion.GetServerList())
		panel.AddField(gotext.Get("服务器开启状态"), "logicOpenStatus", db.Int, form.SelectSingle).
			FieldOptions(open)
		panel.EnableAjax(gotext.Get("操作成功"), gotext.Get("操作失败，请联系管理员"), "/admin/info/t_game_servers")
		return panel
	}, "/admin/popup/updateServerSetting"))

	info.AddActionButton(template.HTML(gotext.Get("服务器详情")), action.Jump("/admin/ServerDetail?Id={{.Id}}&logicOpenStatus={{(index .Value \"logicOpenStatus\").Value}}"+
		"&logicName={{(index .Value \"logicName\").Value}}"))

	info.SetTable("t_game_servers").SetTitle(gotext.Get("全服信息")).
		HideEditButton().HideNewButton().HideFilterArea().HideRowSelector().HideDetailButton().HideDeleteButton()

	info.AddButton(template.HTML(gotext.Get("通知到服务器")), icon.Save,
		action.PopUp("/admin/NotifyServer", gotext.Get("操作结果"), func(ctx *context.Context) (success bool, msg string, data interface{}) {
			vars := url.Values{}
			ctx1 := &context.Context{}
			err, res := fusion.CallToCenter(&vars, common.MyCfg.Api["reloadGameServerInfos"], "GET", ctx1)
			if err != nil {
				str := gotext.Get("操作失败")
				return false, "", "<h2>" + str + "</h2>"
			}
			if res != "" {
				if res != "All Is OK!" {
					return false, "", "<h2>" + res + "</h2>"
				}
			}
			str := gotext.Get("操作成功")
			return true, "", "<h2>" + str + "</h2>"
		}))

	//info.SetGetDataFn(
	//	func(param parameter.Parameters) (data []map[string]interface{}, size int) {
	//		db_global := fusion.GetBaseGormDB("db_global")
	//		db_global.Raw("SELECT *FROM t_game_servers WHERE id IN (SELECT MIN(id) FROM t_game_servers GROUP BY externalIP,externalPort)").
	//			Scan(&data)
	//		min := math.Min(float64(param.PageInt*param.PageSizeInt), float64(len(data)))
	//		return data[((param.PageInt - 1) * param.PageSizeInt):int(min)], len(data)
	//	})

	formList := tGameServers.GetForm()
	//formList.AddField(gotext.Get("序号"), "customOrder", db.Int, form.Number)
	formList.AddField("Id", "Id", db.Int, form.Number).
		FieldDisableWhenCreate()
	formList.AddField("ExternalIP", "externalIP", db.Char, form.Text)
	formList.AddField("ExternalPort", "externalPort", db.Smallint, form.Number)
	formList.AddField("InternalIP", "internalIP", db.Char, form.Text)
	formList.AddField("InternalName", "internalName", db.Char, form.Text)
	formList.AddField("LogicId", "logicId", db.Int, form.Number)
	formList.AddField("LogicName", "logicName", db.Char, form.Text)
	formList.AddField("LogicOpenTime", "logicOpenTime", db.Datetime, form.Datetime)
	formList.AddField("LogicOpenStatus", "logicOpenStatus", db.Tinyint, form.Number)
	formList.AddField("LogicSpecialFlags", "logicSpecialFlags", db.Int, form.Number)

	formList.SetTable("t_game_servers").SetTitle("TGameServers").SetDescription("TGameServers")

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

		//vars := url.Values{}
		//ctx1 := &context.Context{}
		//err, res := fusion.CallToCenter(&vars, common.MyCfg.Api["reloadGameServerInfos"], "GET", ctx1)
		//if err != nil {
		//	str := gotext.Get("操作失败")
		//	return errors.New(str)
		//}
		//if res != "" {
		//	if res != "All Is OK!" {
		//		return errors.New(res)
		//	}
		//}
		return nil
	})

	return tGameServers
}
