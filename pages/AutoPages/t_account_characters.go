package AutoPages

import (
	"admin/fusion"
	"github.com/GoAdminGroup/go-admin/context"
	"github.com/GoAdminGroup/go-admin/modules/db"
	"github.com/GoAdminGroup/go-admin/plugins/admin/modules/table"
	"github.com/GoAdminGroup/go-admin/template/types"
	"github.com/GoAdminGroup/go-admin/template/types/form"
	table2 "github.com/GoAdminGroup/go-admin/template/types/table"
	"github.com/leonelquinteros/gotext"
	"strconv"
)

func GetTAccountCharactersTable(ctx *context.Context) table.Table {

	tAccountCharacters := table.NewDefaultTable(table.Config{
		Driver:     "mysql",
		Connection: "db_global",
		CanAdd:     false,
		Editable:   true,
		Deletable:  false,
		Exportable: false,
		PrimaryKey: table.PrimaryKey{
			Type: db.Int,
			Name: "characterId",
		},
	})

	info := tAccountCharacters.GetInfo()

	myOP := fusion.GetServerList()
	info.AddField(gotext.Get("服务器Id"), "serverId", db.Int).
		FieldFilterable(types.FilterType{FormType: form.SelectSingle}).
		FieldFilterOptions(myOP)
	info.AddField(gotext.Get("角色ID"), "characterId", db.Int).FieldFilterable()
	info.AddField(gotext.Get("角色名字"), "characterName", db.Varchar).
		FieldFilterable(types.FilterType{Operator: types.FilterOperatorLike})
	info.AddField(gotext.Get("是否启用GM命令"), "allowGM", db.Tinyint).
		FieldEditAble(table2.Select).FieldEditOptions(types.FieldOptions{
		{Value: "0", Text: gotext.Get("关闭")},
		{Value: "1", Text: gotext.Get("开启")},
	}).FieldDisplay(func(model types.FieldModel) interface{} {
		temp, _ := strconv.Atoi(model.Value)
		if temp == 1 {
			return gotext.Get("打开")
		}
		return gotext.Get("关闭")
	})
	info.SetTable("t_account_characters").SetTitle(gotext.Get("设置角色白名单")).
		HideEditButton().
		HideDetailButton().
		HideDeleteButton().
		HideExportButton().
		HideNewButton()

	formList := tAccountCharacters.GetForm()
	formList.AddField("ServerId", "serverId", db.Int, form.Number)
	formList.AddField("CharacterId", "characterId", db.Int, form.Number)
	formList.AddField("CharacterName", "characterName", db.Int, form.Number)
	formList.AddField("AllowGM", "allowGM", db.Tinyint, form.Number)
	formList.SetTable("t_account_characters").SetTitle("TAccountCharacters").SetDescription("TAccountCharacters")

	return tAccountCharacters
}
