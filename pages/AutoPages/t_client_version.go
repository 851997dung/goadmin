package AutoPages

import (
	"admin/common"
	"admin/fusion"
	"html/template"
	"net/url"
	"strconv"
	"strings"

	"github.com/GoAdminGroup/go-admin/context"
	"github.com/GoAdminGroup/go-admin/modules/db"
	"github.com/GoAdminGroup/go-admin/plugins/admin/modules/table"
	"github.com/GoAdminGroup/go-admin/template/icon"
	"github.com/GoAdminGroup/go-admin/template/types/action"
	"github.com/leonelquinteros/gotext"
)

func GetTClientVersionTable(ctx *context.Context) table.Table {

	tClientVersion := table.NewDefaultTable(table.Config{
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

	info := tClientVersion.GetInfo()
	info.AddField("Id", "Id", db.Int)
	info.AddField("Version", "version", db.Int)

	info.AddButton(template.HTML(gotext.Get("更新选中版本号(不踢人)")), icon.Modx, action.PopUp("/test/ajax222", gotext.Get("结果"),
		func(ctx *context.Context) (success bool, msg string, data interface{}) {
			idStr := ctx.Request.FormValue("ids")
			ids := strings.Split(idStr, ",")
			var idsNum = make([]int, 0)

			for _, v := range ids {
				if v == "" {
					continue
				}
				id, _ := strconv.Atoi(v)
				idsNum = append(idsNum, id)
			}
			results := []string{}

			db_global := fusion.GetBaseGormDB("db_global")
			var temp int
			db_global.Raw("Update t_client_version set `version` = `version`+1 where id in ?", idsNum).Scan(&temp)

			vars := &url.Values{}
			ctx1 := &context.Context{}
			err, res := fusion.CallToCenter(vars, common.MyCfg.Api["ReloadClientVersion"], "GET", ctx1)
			if err != nil {
				return false, "", err.Error()
			}
			if res != "" {
				if res != "All Is OK!" {
					return true, "", res
				} else {
					results = append(results, gotext.Get("<h5>已通知到中心服！</h5>"))
				}
			}
			return true, "", results
		}))

	info.AddButton(template.HTML(gotext.Get("更新选中版本号(包含踢人)")), icon.Modx, action.PopUp("/test/ajax", gotext.Get("结果"),
		func(ctx *context.Context) (success bool, msg string, data interface{}) {
			idStr := ctx.Request.FormValue("ids")
			ids := strings.Split(idStr, ",")
			var idsNum = make([]int, 0)

			for _, v := range ids {
				if v == "" {
					continue
				}
				id, _ := strconv.Atoi(v)
				idsNum = append(idsNum, id)
			}
			results := []string{}

			db_global := fusion.GetBaseGormDB("db_global")
			var temp int
			db_global.Raw("Update t_client_version set `version` = `version`+1 where id in ?", idsNum).Scan(&temp)

			vars := &url.Values{}
			ctx1 := &context.Context{}
			err, res := fusion.CallToCenter(vars, common.MyCfg.Api["ReloadClientVersion"], "GET", ctx1)
			if err != nil {
				return false, "", err.Error()
			}
			if res != "" {
				if res != "All Is OK!" {
					return true, "", res
				} else {
					results = append(results, gotext.Get("<h5>已通知到中心服！</h5>"))
				}
			}

			ctx1 = &context.Context{}
			vars = &url.Values{}
			vars.Add("cmd", "ForceExitClient")
			err, res = fusion.CallToCenter(vars, common.MyCfg.Api["GM2GS"], "GET", ctx1)
			if err != nil {
				return false, "", err.Error()
			}
			if res != "" {
				results = append(results, "<h5>"+res+"</h5>")
			}
			results = append(results, "<h5>操作结束，请刷新页面!</h5>")
			return true, "", results
		}))

	info.SetTable("t_client_version").SetTitle(gotext.Get("客户端版本控制")).
		HideNewButton().HideEditButton().HideDeleteButton().HideDetailButton().HideRowSelector().HideExportButton()

	//formList := tClientVersion.GetForm()
	//formList.AddField("Id", "id", db.Int, form.Default)
	//formList.AddField("Version", "version", db.Int, form.Number)
	//
	//formList.SetTable("t_client_version").SetTitle("TClientVersion").SetDescription("TClientVersion")

	return tClientVersion
}
