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

func GetTActivationCodeUseTable(ctx *context.Context) table.Table {

	tActivationCodeUse := table.NewDefaultTable(table.Config{
		Driver:     "mysql",
		Connection: "db_global",
		CanAdd:     true,
		Editable:   true,
		Deletable:  true,
		Exportable: true,
		PrimaryKey: table.PrimaryKey{
			Type: db.Int,
			Name: "playerId",
		},
	})

	info := tActivationCodeUse.GetInfo().HideDeleteButton().HideDetailButton().HideNewButton().HideRowSelector()
	myOp := fusion.GetBatchNumOption()
	info.AddField(gotext.Get("玩家ID"), "playerId", db.Int).FieldFilterable()
	info.AddField(gotext.Get("批号"), "batchNum", db.Int).
		FieldFilterable(types.FilterType{FormType: form.SelectSingle}).
		FieldFilterOptions(myOp)
	info.AddField(gotext.Get("使用时间"), "useTime", db.Datetime)

	info.SetTable("t_activation_code_use").SetTitle(gotext.Get("礼包码使用记录"))

	return tActivationCodeUse
}
