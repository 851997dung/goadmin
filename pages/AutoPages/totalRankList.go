package AutoPages

import (
	"admin/common/def"
	"admin/fusion"
	"github.com/GoAdminGroup/go-admin/context"
	"github.com/GoAdminGroup/go-admin/modules/db"
	"github.com/GoAdminGroup/go-admin/plugins/admin/modules/parameter"
	"github.com/GoAdminGroup/go-admin/plugins/admin/modules/table"
	"github.com/GoAdminGroup/go-admin/template/types"
	"github.com/GoAdminGroup/go-admin/template/types/form"
	"github.com/leonelquinteros/gotext"
	"math"
	"sort"
	"strconv"
	"strings"
)

func GetTotalRechargeRankList(ctx *context.Context) table.Table {
	LevelRankList := table.NewDefaultTable(table.Config{
		Driver:     "mysql",
		Connection: "",
		CanAdd:     false,
		Editable:   false,
		Deletable:  false,
		Exportable: false,
		PrimaryKey: table.PrimaryKey{
			Type: db.Int,
			Name: "Id",
		},
	})

	myop := fusion.GetServerList()
	info := LevelRankList.GetInfo().HideNewButton().HideDeleteButton().HideDetailButton().HideEditButton().HideRowSelector()
	info.AddField(gotext.Get("服务器ID"), "ServerID", db.Int).
		FieldHide().
		FieldFilterable(types.FilterType{FormType: form.SelectSingle}).
		FieldFilterOptions(myop)
	info.AddField(gotext.Get("排名"), "Id", db.Int)
	info.AddField(gotext.Get("角色ID"), "PlayerId", db.Int).FieldHide()
	info.AddField(gotext.Get("角色名"), "PlayerName", db.Varchar)
	info.AddField(gotext.Get("充值金额"), "TotalPrice", db.Int)

	info.SetGetDataFn(
		func(param parameter.Parameters) (data []map[string]interface{}, size int) {
			serverId, _ := strconv.Atoi(ctx.Request.FormValue("ServerID"))
			if serverId == 0 {
				return nil, 0
			}
			db_global := fusion.GetBaseGormDB("db_global")
			tableNames := []string{}
			db_global.Raw("SELECT table_name FROM information_schema.TABLES WHERE TABLE_SCHEMA='mmorpg_global' ").Scan(&tableNames)
			PriceCut := make([]def.TotalPriceCut, 0)
			var playerId2Price = make(map[int]int)
			for _, v := range tableNames {
				if strings.Index(v, "t_orders_") != -1 || -1 != strings.Index(v, "t_web_orders_") {
					sql := "select playerId,sum(buyPrice) AS TotalPrice from " + v + " WHERE gsId = ? And payTime != NULL group by playerId"
					db_global.Raw(sql, serverId).Scan(&PriceCut)
					for _, cut := range PriceCut {
						playerId2Price[cut.PlayerId] = cut.TotalPrice
					}
				}
			}

			var playerIds = make([]interface{}, 0)
			for _, v := range PriceCut {
				playerIds = append(playerIds, v.PlayerId)
			}
			var tempNameData = make([]map[string]interface{}, 0)
			db_char := fusion.GetServerGormDB("db_char", int64(serverId))
			db_char.Table("inst_player_char").
				Select("ipcInstID,ipcNickName").
				Where("ipcInstID in ?", playerIds).
				Scan(&tempNameData)

			var nameData = make(map[string]interface{})
			for _, v := range tempNameData {
				tempStr := fusion.Strval(v["ipcInstID"])
				nameData[tempStr] = v["ipcNickName"]
			}

			for k, v := range playerId2Price {
				var temp = make(map[string]interface{})
				//temp["Id"] =
				temp["PlayerId"] = k
				temp["PlayerName"] = nameData[strconv.Itoa(k)]
				temp["TotalPrice"] = v
				data = append(data, temp)
			}
			sort.Slice(data, func(i, j int) bool {
				temp1, _ := strconv.Atoi(fusion.Strval(data[i]["TotalPrice"]))
				temp2, _ := strconv.Atoi(fusion.Strval(data[j]["TotalPrice"]))
				return temp1 > temp2
			})
			for k, _ := range data {
				data[k]["Id"] = k + 1
			}
			min := math.Min(float64(param.PageInt*param.PageSizeInt), float64(len(data)))
			return data[((param.PageInt - 1) * param.PageSizeInt):int(min)], len(data)
		})

	return LevelRankList
}
