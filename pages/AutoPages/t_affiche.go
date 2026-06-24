package AutoPages

import (
	"github.com/GoAdminGroup/go-admin/context"
	"github.com/GoAdminGroup/go-admin/modules/db"
	"github.com/GoAdminGroup/go-admin/plugins/admin/modules/table"
	"github.com/GoAdminGroup/go-admin/template/types"
	"github.com/GoAdminGroup/go-admin/template/types/form"
	"github.com/leonelquinteros/gotext"
	"time"
)

func GetTAfficheTable(ctx *context.Context) table.Table {

	tAffiche := table.NewDefaultTable(table.Config{
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

	info := tAffiche.GetInfo()

	info.AddField("Id", "Id", db.Int)
	info.AddField(gotext.Get("内容"), "content", db.Text).
		FieldHide()
	info.AddField(gotext.Get("更新时间"), "updateTime", db.Time)
	info.AddField(gotext.Get("开始时间"), "startTime", db.Time)
	info.AddField(gotext.Get("结束时间"), "endTime", db.Time)
	info.AddField(gotext.Get("是否激活"), "isActive", db.Bool)

	info.SetTable("t_affiche").SetTitle(gotext.Get("公告管理")).HideRowSelector().HideExportButton()

	now := time.Now().Format("2006-01-02 15:04:05")
	formList := tAffiche.GetForm()
	formList.AddField("Id", "Id", db.Int, form.Number).
		FieldDisableWhenCreate().FieldDisplayButCanNotEditWhenUpdate()
	formList.AddField(gotext.Get("内容"), "content", db.Mediumtext, form.TextArea)
	formList.AddField(gotext.Get("更新时间"), "updateTime", db.Time, form.Datetime).FieldDefault(now).FieldHide()
	formList.AddField(gotext.Get("开始时间"), "startTime", db.Time, form.Datetime)
	formList.AddField(gotext.Get("结束时间"), "endTime", db.Time, form.Datetime)
	formList.AddField(gotext.Get("是否激活"), "isActive", db.Tinyint, form.Switch).
		FieldOptions(types.FieldOptions{
			{Value: "1"},
			{Value: "0"},
		})

	formList.SetTable("t_affiche").SetTitle(gotext.Get("公告管理"))

	return tAffiche
}
