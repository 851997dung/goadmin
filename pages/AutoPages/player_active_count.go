package AutoPages

import (
	"admin/fusion"
	"github.com/GoAdminGroup/go-admin/context"
	"github.com/GoAdminGroup/go-admin/modules/db"
	"github.com/GoAdminGroup/go-admin/plugins/admin/modules/table"
	"github.com/GoAdminGroup/go-admin/template/types"
	"github.com/GoAdminGroup/go-admin/template/types/form"
	"github.com/leonelquinteros/gotext"
	"strconv"
)

func GetPlayerActiveCountTable(ctx *context.Context) table.Table {

	playerActiveCount := table.NewDefaultTable(table.Config{
		Driver:     "mysql",
		Connection: "go_manager",
		CanAdd:     false,
		Editable:   false,
		Deletable:  false,
		Exportable: false,
		PrimaryKey: table.PrimaryKey{
			Type: db.Int,
			Name: "Id",
		},
	})

	info := playerActiveCount.GetInfo().HideNewButton().HideDeleteButton().HideDetailButton().HideEditButton().HideRowSelector()
	ops := fusion.GetServerNameStrings()
	myOP := fusion.GetServerList()
	info.AddField("Id", "Id", db.Int)
	info.AddField(gotext.Get("服务器ID"), "gsId", db.Int).
		FieldDisplay(func(model types.FieldModel) interface{} {
			temp, _ := strconv.Atoi(model.Value)
			str, err := ops[temp]
			if !err {
				return model.Value
			}
			return str
		}).FieldFilterable(types.FilterType{FormType: form.SelectBox, Options: myOP})
	info.AddField(gotext.Get("日活跃用户A"), "activeA", db.Bigint)
	info.AddField(gotext.Get("日活跃用户B"), "activeB", db.Bigint)
	info.AddField(gotext.Get("新注册用户活跃率"), "newActivePer", db.Double).
		FieldDisplay(func(model types.FieldModel) interface{} {
			temp, _ := strconv.ParseFloat(model.Value, 32)
			return strconv.FormatFloat(temp, 'f', 2, 32)
		})
	info.AddField(gotext.Get("日登录人次"), "loginManTimes", db.Bigint)
	info.AddField(gotext.Get("记录时间"), "logTime", db.Date).
		FieldFilterable(types.FilterType{FormType: form.DateRange})

	info.SetTable("player_active_count").SetTitle(gotext.Get("用户日活跃"))

	return playerActiveCount
}
