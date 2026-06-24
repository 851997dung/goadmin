package CustomPages

import (
	"admin/common/def"
	"admin/common/def/flowDef"
	"admin/fusion"
	"github.com/GoAdminGroup/go-admin/context"
	"github.com/GoAdminGroup/go-admin/modules/db"
	"github.com/GoAdminGroup/go-admin/plugins/admin/modules/parameter"
	"github.com/GoAdminGroup/go-admin/plugins/admin/modules/table"
	"github.com/GoAdminGroup/go-admin/template/types"
	"github.com/GoAdminGroup/go-admin/template/types/form"
	"github.com/leonelquinteros/gotext"
	"strconv"
	"strings"
)

func GetBuyStoneFromShopTable(ctx *context.Context) table.Table {
	BuyStoneFromShop := table.NewDefaultTable(table.Config{
		Driver:     "mysql",
		Connection: "db_log",
		CanAdd:     true,
		Editable:   true,
		Deletable:  true,
		Exportable: true,
		PrimaryKey: table.PrimaryKey{
			Type: db.Int,
			Name: "serverId",
		},
	})

	info := BuyStoneFromShop.GetInfo().HideEditButton().HideNewButton().HideRowSelector().HideDetailButton().HideDeleteButton().HideRowSelector()
	serverList := fusion.GetServerList()

	info.AddField(gotext.Get("服务器ID"), "serverId", db.Int).
		FieldFilterable(types.FilterType{FormType: form.SelectSingle}).
		FieldFilterOptions(serverList)
	info.AddField(gotext.Get("道具ID"), "stoneType", db.Varchar).
		FieldFilterable()
	info.AddField(gotext.Get("道具数量"), "stoneNum", db.Int)
	info.AddField(gotext.Get("记录时间"), "timeRange", db.Int).
		FieldFilterable(types.FilterType{FormType: form.DatetimeRange}).
		FieldHide()

	info.SetTable("Diamond4RebornRec").SetTitle(gotext.Get("单服购买商店物品查询")).SetPageSizeList([]int{50, 100, 9999})

	info.SetGetDataFn(
		func(param parameter.Parameters) (data []map[string]interface{}, size int) {
			serverId, _ := strconv.Atoi(ctx.Request.FormValue("serverId"))
			stoneType := ctx.Request.FormValue("stoneType")

			timeStart := ctx.Request.FormValue("timeRange_start__goadmin")
			timeEnd := ctx.Request.FormValue("timeRange_end__goadmin")

			if stoneType == "" {
				return nil, 0
			}

			where := ""
			where += " flowType = " + strconv.Itoa(flowDef.LFT_SHOP) + " and "
			where += " itemNum >=0 and "
			if timeStart != "" {
				where += " logTime >= '" + timeStart + "' and "
			}
			if timeEnd != "" {
				where += " logTime <= '" + timeEnd + "' and "
			}
			if stoneType != "" {
				where += " itemTypeId = '" + stoneType + "' and "
			}

			if where != "" {
				where = " where " + where
				where = strings.TrimSuffix(where, " and ")
			} else {
				return nil, 0
			}

			gormDB := fusion.GetServerGormDB("db_log", int64(serverId))
			if gormDB == nil {
				return nil, 0
			}
			sql1, _ := fusion.CreateUnionTableSqlByDays(timeStart, timeEnd, def.ItemLog{}.TableName(),
				gormDB, def.ItemLog{}, where, false)
			var ItemSum int
			sql1 = "select sum(itemNum) from (" + sql1 + ") as A"
			gormDB.Raw(sql1).Scan(&ItemSum)
			temp := make(map[string]interface{}, 0)
			temp["serverId"] = serverId
			temp["stoneType"] = stoneType
			temp["stoneNum"] = ItemSum
			data = append(data, temp)

			return data, len(data)
		})

	info.SetExportProcessFn(func(param parameter.Parameters) (types.PanelInfo, error) {
		panelInfo, _ := BuyStoneFromShop.GetData(param.WithIsAll(param.IsAll()))
		return types.PanelInfo{Thead: panelInfo.Thead, InfoList: panelInfo.InfoList}, nil
	})
	return BuyStoneFromShop
}
