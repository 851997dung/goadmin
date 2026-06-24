package AutoPages

import (
	"admin/common/def"
	"admin/fusion"
	"strconv"

	"github.com/GoAdminGroup/go-admin/context"
	"github.com/GoAdminGroup/go-admin/modules/auth"
	"github.com/GoAdminGroup/go-admin/modules/db"
	form2 "github.com/GoAdminGroup/go-admin/plugins/admin/modules/form"
	"github.com/GoAdminGroup/go-admin/plugins/admin/modules/table"
	"github.com/GoAdminGroup/go-admin/template/types"
	"github.com/GoAdminGroup/go-admin/template/types/form"
	table2 "github.com/GoAdminGroup/go-admin/template/types/table"
	"github.com/leonelquinteros/gotext"
)

func GetComplaintManagementTable(ctx *context.Context) table.Table {

	complaintManagement := table.NewDefaultTable(table.Config{
		Driver:     "mysql",
		Connection: "go_manager",
		CanAdd:     true,
		Editable:   true,
		Deletable:  true,
		Exportable: true,
		PrimaryKey: table.PrimaryKey{
			Type: db.Int,
			Name: "Id",
		},
	})

	info := complaintManagement.GetInfo().HideEditButton().HideNewButton().HideRowSelector().HideDetailButton().HideDeleteButton().HideRowSelector()

	typeStr := def.GetComplaintTypeString()
	fromStr := def.GetComplaintFromString()
	statusStr := def.GetComplaintStatusString()

	statusOP := def.GetComplaintOps(2)
	serverStr := fusion.GetServerNameStrings()
	info.AddField("Id", "Id", db.Int)
	info.AddField(gotext.Get("问题分类"), "cType", db.Int).
		FieldDisplay(func(model types.FieldModel) interface{} {
			temp, _ := strconv.Atoi(model.Value)
			return typeStr[temp]
		})
	info.AddField(gotext.Get("问题来源"), "cFrom", db.Int).
		FieldDisplay(func(model types.FieldModel) interface{} {
			temp, _ := strconv.Atoi(model.Value)
			return fromStr[temp]
		})
	info.AddField(gotext.Get("服务器ID"), "gsId", db.Int).
		FieldDisplay(func(model types.FieldModel) interface{} {
			temp, _ := strconv.Atoi(model.Value)
			return serverStr[temp]
		})
	info.AddField(gotext.Get("角色名"), "playerName", db.Varchar)
	info.AddField(gotext.Get("角色Id"), "uid", db.Int)
	info.AddField(gotext.Get("记录时间"), "logTime", db.Datetime)
	info.AddField(gotext.Get("内容"), "content", db.Mediumtext)
	info.AddField(gotext.Get("状态"), "status", db.Int).
		FieldEditOptions(statusOP).
		FieldEditAble(table2.Select).
		FieldDisplay(func(model types.FieldModel) interface{} {
			temp, _ := strconv.Atoi(model.Value)
			return statusStr[temp]
		})
	info.AddField(gotext.Get("操作者"), "operatorName", db.Varchar)

	info.SetTable("complaint_management").SetTitle(gotext.Get("投诉建议管理"))

	formList := complaintManagement.GetForm()
	formList.AddField("Id", "Id", db.Int, form.Number).
		FieldDisableWhenCreate()
	formList.AddField(gotext.Get("问题分类"), "cType", db.Int, form.SelectSingle).FieldOptions(def.GetComplaintOps(0))
	formList.AddField(gotext.Get("问题来源"), "cFrom", db.Int, form.SelectSingle).FieldOptions(def.GetComplaintOps(1))
	formList.AddField(gotext.Get("内容"), "content", db.Mediumtext, form.RichText)
	formList.AddField(gotext.Get("服务器ID"), "gsId", db.Int, form.SelectSingle).FieldOptions(fusion.GetServerList())
	formList.AddField(gotext.Get("角色名"), "playerName", db.Varchar, form.Text)
	formList.AddField(gotext.Get("角色ID"), "uid", db.Int, form.Text)
	formList.AddField(gotext.Get("记录时间"), "logTime", db.Datetime, form.Datetime).
		FieldDisableWhenCreate()
	formList.AddField(gotext.Get("状态"), "status", db.Int, form.Number).
		FieldDisableWhenCreate()
	formList.AddField(gotext.Get("操作者"), "operatorName", db.Varchar, form.Text).
		FieldDisableWhenCreate()

	formList.SetTable("complaint_management").SetTitle(gotext.Get("投诉建议管理"))

	formList.SetUpdateFn(func(values form2.Values) error {
		user := auth.Auth(ctx)
		Id := values.Get("Id")
		status := values.Get("status")
		operatorName := user.Name

		go_manager := fusion.GetBaseGormDB("go_manager")
		go_manager.Table("complaint_management").
			Where("Id = ?", Id).
			Update("status", status).
			Update("operatorName", operatorName)

		return nil
	})

	return complaintManagement
}
