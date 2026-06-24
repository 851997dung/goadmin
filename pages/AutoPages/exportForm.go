package AutoPages

import (
	"admin/common/def"
	"admin/fusion"
	"strconv"

	"github.com/GoAdminGroup/go-admin/context"
	"github.com/GoAdminGroup/go-admin/modules/db"
	"github.com/GoAdminGroup/go-admin/plugins/admin/modules/parameter"
	"github.com/GoAdminGroup/go-admin/plugins/admin/modules/table"
	"github.com/GoAdminGroup/go-admin/template/types"
	"github.com/GoAdminGroup/go-admin/template/types/form"
	"github.com/leonelquinteros/gotext"
)

func GetTestFilterForm(ctx *context.Context) table.Table {

	testTable := table.NewDefaultTable(table.Config{
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
	myOP := fusion.GetServerList()
	myTableOp := def.GetTableNameOp()
	info := testTable.GetInfo().HidePagination().HideRowSelector().HideNewButton().HideDeleteButton().HideDetailButton().HideEditButton()
	info.AddField(gotext.Get("请选择要导出的表"), "TableId", db.Text).
		FieldFilterable(types.FilterType{FormType: form.SelectSingle}).
		FieldFilterOptions(myTableOp).
		FieldHide()
	info.AddField(gotext.Get("选择要导出的服务器"), "ServerIDs", db.Int).
		FieldFilterable(types.FilterType{FormType: form.SelectBox}).
		FieldFilterOptions(myOP).
		FieldHide()
	info.AddField(gotext.Get("账号ID"), "AcctId", db.Int).
		FieldHide()
	info.AddField(gotext.Get("角色ID"), "PlayerId", db.Int).
		FieldHide()
	info.AddField(gotext.Get("角色名"), "PlayerName", db.Varchar).
		FieldFilterable().
		FieldHide()
	info.AddField(gotext.Get("记录时间"), "LogTime", db.Datetime).
		FieldFilterable(types.FilterType{FormType: form.DatetimeRange}).
		FieldHide()
	info.HideCheckBoxColumn().HideEditButton().HideDeleteButton().HideDetailButton()

	info.SetGetDataFn(func(param parameter.Parameters) (data []map[string]interface{}, size int) {
		return nil, 0
	})

	info.SetExportProcessFn(func(param parameter.Parameters) (types.PanelInfo, error) {
		tableId, _ := strconv.Atoi(ctx.Request.FormValue("TableId"))
		panelInfo := table.PanelInfo{}
		if tableId == def.Log_character_chats {
			panelInfo, _ = GetLogCharacterChatsTable(ctx).GetData(param.WithIsAll(true))
		}
		if tableId == def.Log_character_onlines {
			panelInfo, _ = GetLogCharacterOnlinesTable(ctx).GetData(param.WithIsAll(true))
		}
		if tableId == def.Log_player_amulet {
			panelInfo, _ = GetLogPlayerAmuletTable(ctx).GetData(param.WithIsAll(true))
		}
		if tableId == def.Log_player_dungeons {
			panelInfo, _ = GetLogPlayerDungeonsTable(ctx).GetData(param.WithIsAll(true))
		}
		if tableId == def.Log_player_equip {
			panelInfo, _ = GetLogPlayerEquipTable(ctx).GetData(param.WithIsAll(true))
		}
		if tableId == def.Log_player_ertl {
			panelInfo, _ = GetLogPlayerErtlTable(ctx).GetData(param.WithIsAll(true))
		}
		if tableId == def.Log_player_exps {
			panelInfo, _ = GetLogPlayerExpsTable(ctx).GetData(param.WithIsAll(true))
		}
		if tableId == def.Log_player_items {
			panelInfo, _ = GetLogPlayerItemsTable(ctx).GetData(param.WithIsAll(true))
		}
		if tableId == def.Log_player_levelups {
			panelInfo, _ = GetLogPlayerLevelupsTable(ctx).GetData(param.WithIsAll(true))
		}
		if tableId == def.Log_player_moneys {
			panelInfo, _ = GetLogPlayerMoneysTable(ctx).GetData(param.WithIsAll(true))
		}
		if tableId == def.Log_player_quests {
			panelInfo, _ = GetLogPlayerQuestsTable(ctx).GetData(param.WithIsAll(true))
		}
		if tableId == def.Log_player_stone {
			panelInfo, _ = GetLogPlayerStoneTable(ctx).GetData(param.WithIsAll(true))
		}
		if tableId == def.Log_player_login {
			panelInfo, _ = GetLogPlayerStoneTable(ctx).GetData(param.WithIsAll(true))
		}
		return types.PanelInfo{Thead: panelInfo.Thead, InfoList: panelInfo.InfoList}, nil
	})
	return testTable
}
