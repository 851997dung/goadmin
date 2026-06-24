package AutoPages

import (
	"admin/fusion"
	"github.com/GoAdminGroup/go-admin/context"
	"github.com/GoAdminGroup/go-admin/modules/db"
	"github.com/GoAdminGroup/go-admin/plugins/admin/modules/table"
	"github.com/GoAdminGroup/go-admin/template/types"
	"github.com/GoAdminGroup/go-admin/template/types/form"
	"github.com/leonelquinteros/gotext"
)

func GetTExpCfgTable(ctx *context.Context) table.Table {

	tExpCfg := table.NewDefaultTable(table.Config{
		Driver:     "mysql",
		Connection: "db_global",
		CanAdd:     true,
		Editable:   true,
		Deletable:  true,
		Exportable: true,
		PrimaryKey: table.PrimaryKey{
			Type: db.Int,
			Name: "CreatureId",
		},
	})

	info := tExpCfg.GetInfo().HideDetailButton()
	myOP := fusion.GetServerList()
	str := fusion.GetServerNameStrings()
	info.AddField(gotext.Get("服务器ID"), "gsId", db.Int).
		FieldFilterable(types.FilterType{FormType: form.SelectSingle}).
		FieldFilterOptions(myOP)
	info.AddColumn(gotext.Get("服务器名"), func(value types.FieldModel) interface{} {
		return str[int(value.Row["gsId"].(int64))]
	})
	info.AddField(gotext.Get("怪物id"), "CreatureId", db.Int).FieldEditAble().FieldFilterable()
	info.AddField(gotext.Get("怪物名"), "CreatureName", db.Varchar).FieldEditAble().
		FieldFilterable(types.FilterType{Operator: types.FilterOperatorLike})
	info.AddField(gotext.Get("经验系数（万分比）"), "Coefficient", db.Int).FieldEditAble()
	info.SetTable("t_exp_cfg").SetTitle(gotext.Get("怪物经验系数表"))

	formList := tExpCfg.GetForm()
	formList.AddField(gotext.Get("服务器ID"), "gsId", db.Int, form.Text).
		FieldDisplayButCanNotEditWhenUpdate()
	formList.AddField(gotext.Get("服务器名"), "serverName", db.Varchar, form.Text).
		FieldDisplay(func(value types.FieldModel) interface{} {
			return str[int(value.Row["gsId"].(int64))]
		}).
		FieldDisplayButCanNotEditWhenUpdate()
	formList.AddField(gotext.Get("怪物id"), "CreatureId", db.Int, form.Number)
	formList.AddField(gotext.Get("怪物名"), "CreatureName", db.Varchar, form.Text)
	formList.AddField(gotext.Get("经验系数（万分比）"), "Coefficient", db.Int, form.Number)
	formList.SetTable("t_exp_cfg").SetTitle(gotext.Get("怪物经验系数表"))
	return tExpCfg
}
