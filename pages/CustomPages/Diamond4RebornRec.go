package CustomPages

import (
	"admin/common/def"
	"admin/common/def/currencyDef"
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

func GetDiamond4RebornRecTable(ctx *context.Context) table.Table {
	Diamond4RebornRec := table.NewDefaultTable(table.Config{
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

	info := Diamond4RebornRec.GetInfo().HideEditButton().HideNewButton().HideRowSelector().HideDetailButton().HideDeleteButton().HideRowSelector()
	serverList := fusion.GetServerList()

	info.AddField(gotext.Get("服务器ID"), "serverId", db.Int).
		FieldFilterable(types.FilterType{FormType: form.SelectSingle}).
		FieldFilterOptions(serverList)
	info.AddField(gotext.Get("转生消耗钻石总量"), "totalUsedDiamond", db.Int)
	info.AddField(gotext.Get("记录时间"), "timeRange", db.Int).
		FieldFilterable(types.FilterType{FormType: form.DatetimeRange}).
		FieldHide()

	info.SetTable("Diamond4RebornRec").SetTitle(gotext.Get("单服转生消耗钻石总量")).SetPageSizeList([]int{50, 100, 9999})

	info.SetGetDataFn(
		func(param parameter.Parameters) (data []map[string]interface{}, size int) {
			serverId, _ := strconv.Atoi(ctx.Request.FormValue("serverId"))
			timeStart := ctx.Request.FormValue("timeRange_start__goadmin")
			timeEnd := ctx.Request.FormValue("timeRange_end__goadmin")
			if timeStart == "" || timeEnd == "" {
				return nil, 0
			}

			where := ""
			where += " flowType = " + strconv.Itoa(flowDef.LFT_REBORN_EX) + " and "
			where += " moneyType = " + strconv.Itoa(currencyDef.Diamond) + " and "
			if timeStart != "" {
				where += " logTime >= '" + timeStart + "' and "
			}
			if timeEnd != "" {
				where += " logTime <= '" + timeEnd + "' and "
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
			sql1, _ := fusion.CreateUnionTableSqlByDays(timeStart, timeEnd, def.MoneyLog{}.TableName(),
				gormDB, def.MoneyLog{}, where, false)
			var moneyValueSum int
			sql1 = "select sum(moneyValue) from (" + sql1 + ") as A"
			gormDB.Raw(sql1).Scan(&moneyValueSum)
			temp := make(map[string]interface{}, 0)
			temp["serverId"] = serverId
			temp["totalUsedDiamond"] = moneyValueSum
			data = append(data, temp)

			return data, len(data)
		})

	info.SetExportProcessFn(func(param parameter.Parameters) (types.PanelInfo, error) {
		panelInfo, _ := Diamond4RebornRec.GetData(param.WithIsAll(param.IsAll()))
		return types.PanelInfo{Thead: panelInfo.Thead, InfoList: panelInfo.InfoList}, nil
	})
	return Diamond4RebornRec
}
