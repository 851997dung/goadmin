package AutoPages

import (
	"admin/fusion"
	"encoding/json"
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
)

func GetRankList(ctx *context.Context) table.Table {
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

	rankListType, _ := strconv.Atoi(ctx.Request.FormValue("RankListType"))
	myop := fusion.GetServerList()

	info := LevelRankList.GetInfo().HideNewButton().HideDeleteButton().HideDetailButton().HideEditButton().HideRowSelector()
	info.AddField(gotext.Get("排行榜类型"), "RankListType", db.Int).
		FieldFilterable(types.FilterType{FormType: form.SelectSingle}).
		FieldFilterOptions(types.FieldOptions{
			{Text: gotext.Get("等级排行"), Value: "1"},
			{Text: gotext.Get("充值排行"), Value: "2"},
			{Text: gotext.Get("网页充值排行"), Value: "3"},
		}).FieldHide()
	info.AddField(gotext.Get("服务器ID"), "ServerID", db.Int).
		FieldHide().FieldFilterable(types.FilterType{Options: myop, FormType: form.SelectSingle}).
		FieldFilterOptions(myop)
	info.AddField(gotext.Get("排名"), "Id", db.Int)
	info.AddField(gotext.Get("角色ID"), "PlayerId", db.Int)
	info.AddField(gotext.Get("角色名"), "PlayerName", db.Varchar)
	if rankListType == 1 {
		info.AddField(gotext.Get("账号ID"), "AcctId", db.Int)
		info.AddField(gotext.Get("角色等级"), "PlayerLevel", db.Int)
	} else {
		info.AddField(gotext.Get("充值金额"), "TotalPrice", db.Int)
	}
	info.AddField(gotext.Get("日期"), "Date", db.Date).
		FieldFilterable(types.FilterType{FormType: form.Date}).
		FieldHide()

	info.SetGetDataFn(
		func(param parameter.Parameters) (data []map[string]interface{}, size int) {
			rankListType, _ := strconv.Atoi(ctx.Request.FormValue("RankListType"))
			serverId, _ := strconv.Atoi(ctx.Request.FormValue("ServerID"))
			date := ctx.Request.FormValue("Date")
			if date == "" || serverId == 0 || rankListType == 0 {
				return nil, 0
			}
			var filterField string
			switch rankListType {
			case 1:
				filterField = "levelRankDaily"
				break
			case 2:
				filterField = "rechargeRankDaily"
				break
			case 3:
				filterField = "webRechargeRankDaily"
				break
			case 4:
				filterField = "rechargeRankDaily,webRechargeRankDaily"
				break
			}
			go_manager := fusion.GetBaseGormDB("go_manager")
			if rankListType != 4 {
				var levelRankString string
				go_manager.Table("player_data_count").
					Select(filterField).
					Where("gsId = ?", serverId).
					Where("logTime like ?", "%"+date+"%").
					Scan(&levelRankString)
				temp := fusion.String2Bytes(levelRankString)
				err := json.Unmarshal(temp, &data)
				if err != nil {
					return nil, 0
				}
				var playerIds = make([]interface{}, 0)
				for _, v := range data {
					playerIds = append(playerIds, v["PlayerId"])
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

				for k, v := range data {
					data[k]["ServerID"] = serverId
					data[k]["Id"] = k + 1
					tempStr := fusion.Strval(v["PlayerId"])
					if nameData[tempStr] != nil {
						data[k]["PlayerName"] = nameData[tempStr]
					}
				}
			} else {

			}
			min := math.Min(float64(param.PageInt*param.PageSizeInt), float64(len(data)))
			return data[((param.PageInt - 1) * param.PageSizeInt):int(min)], len(data)
		})

	return LevelRankList
}

func GetRankListPro(ctx *context.Context) table.Table {
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

	rankListType, _ := strconv.Atoi(ctx.Request.FormValue("RankListType"))
	myop := fusion.GetServerList()

	info := LevelRankList.GetInfo().HideNewButton().HideDeleteButton().HideDetailButton().HideEditButton().HideRowSelector()
	info.SetTitle(gotext.Get("数据排行榜"))

	info.AddField(gotext.Get("排行榜类型"), "RankListType", db.Int).
		FieldFilterable(types.FilterType{FormType: form.SelectSingle}).
		FieldFilterOptions(types.FieldOptions{
			{Text: gotext.Get("钻石排行"), Value: "1"},
			{Text: gotext.Get("单服累充排行"), Value: "2"},
			{Text: gotext.Get("实力排行"), Value: "3"},
		}).FieldHide()
	info.AddField(gotext.Get("服务器ID"), "ServerID", db.Int).
		FieldHide().FieldFilterable(types.FilterType{Options: myop, FormType: form.SelectSingle}).
		FieldFilterOptions(myop)
	info.AddField(gotext.Get("排名"), "Id", db.Int)

	switch rankListType {
	case 1:
		info.AddField(gotext.Get("钻石数量"), "TotalDiamond", db.Int)
		break
	case 2:
		info.AddField(gotext.Get("充值额度"), "TotalPrice", db.Int)
		break
	case 3:
		info.AddField(gotext.Get("实力"), "TotalPower", db.Int)
		break
	}
	info.AddField(gotext.Get("角色ID"), "PlayerId", db.Int)
	info.AddField(gotext.Get("角色名"), "PlayerName", db.Varchar)
	switch rankListType {
	case 2:
		info.AddField(gotext.Get("充值次数"), "TotalRechargeTimes", db.Int)
		break
	default:
		info.AddField(gotext.Get("充值额度"), "TotalPrice", db.Int)
	}

	info.SetGetDataFn(
		func(param parameter.Parameters) (data []map[string]interface{}, size int) {
			rankListType, _ := strconv.Atoi(ctx.Request.FormValue("RankListType"))
			serverId, _ := strconv.Atoi(ctx.Request.FormValue("ServerID"))
			if serverId == 0 || rankListType == 0 {
				return nil, 0
			}

			type cut1 struct {
				PlayerId uint32 `gorm:"column:playerId"`
				BuyPrice int64  `gorm:"column:buyPrice"`
				PayId    uint32 `gorm:"column:payId"`
			}
			type cut2 struct {
				PlayerId      uint32 `gorm:"column:playerId"`
				RechargeTimes uint32 `gorm:"column:rechargeTimes"`
			}

			var purchase = make([]map[string]interface{}, 0)

			db_world := fusion.GetBaseGormDB("db_world")
			db_world.Table("auto_purchase").
				Select("ID,platformID").Scan(&purchase)
			PlatForm4ID := map[int]int{}
			for _, v := range purchase {
				PlatFormId, _ := strconv.Atoi(fusion.Strval(v["platformID"]))
				id, _ := strconv.Atoi(fusion.Strval(v["ID"]))
				PlatForm4ID[id] = PlatFormId
			}

			if rankListType == 1 {
				type playerCut struct {
					PlayerId    uint32 `gorm:"column:playerId"`
					NickName    string `gorm:"column:nickName"`
					DiamondsNum int64  `gorm:"column:diamondsNum"`
				}
				playerCuts := make([]playerCut, 0)
				db_char := fusion.GetServerGormDB("db_char", int64(serverId))

				db_char.Table("inst_player_char").
					Select("ipcInstID AS playerId,ipcNickName AS NickName,ipcCurrencies->'$[2]' AS DiamondsNum ").
					Order("DiamondsNum DESC").Limit(50).Scan(&playerCuts)
				playerIds := make([]uint32, 0)
				for _, v := range playerCuts {
					playerIds = append(playerIds, v.PlayerId)
				}

				db_global := fusion.GetBaseGormDB("db_global")
				sql1, sql2 := fusion.CreateUnionTableSql("t_web_orders_", db_global, "where playerId in ? and gsId = ? ", " group by playerId,gsId,payId")

				cut1s := make([]cut1, 0)
				db_global.Raw(sql1, playerIds, serverId).Scan(&cut1s)
				playerId2Price := make(map[uint32]int64, 0)
				for _, v := range cut1s {
					if PlatForm4ID[int(v.PayId)] == 3 {
						playerId2Price[v.PlayerId] = v.BuyPrice * 1000
					} else {
						playerId2Price[v.PlayerId] = v.BuyPrice
					}
				}

				cut2s := make([]cut2, 0)
				db_global.Raw(sql2, playerIds, serverId).Scan(&cut2s)
				playerId2Times := make(map[uint32]uint32, 0)
				for _, v := range cut2s {
					playerId2Times[v.PlayerId] = v.RechargeTimes
				}

				for k, v := range playerCuts {
					temp := make(map[string]interface{}, 0)
					temp["Id"] = k + 1
					temp["PlayerId"] = v.PlayerId
					temp["PlayerName"] = v.NickName
					temp["ServerID"] = serverId
					temp["TotalDiamond"] = v.DiamondsNum

					price, err1 := playerId2Price[v.PlayerId]
					if !err1 {
						temp["TotalPrice"] = 0
					} else {
						temp["TotalPrice"] = price
					}
					times, err1 := playerId2Times[v.PlayerId]
					if !err1 {
						temp["TotalRechargeTimes"] = 0
					} else {
						temp["TotalRechargeTimes"] = times
					}
					data = append(data, temp)
				}
			} else if rankListType == 2 {
				db_global := fusion.GetBaseGormDB("db_global")
				sql1, sql2 := fusion.CreateUnionTableSql("t_web_orders_", db_global)

				sql1 = sql1 + "where gsId = ? " + " group by playerId,gsId " + " order by buyPrice desc limit 50"
				cut1s := make([]cut1, 0)
				db_global.Raw(sql1, serverId).Scan(&cut1s)
				playerId2Price := make(map[uint32]int64, 0)
				playerIds := make([]uint32, 0)
				for _, v := range cut1s {
					if PlatForm4ID[int(v.PayId)] == 3 {
						playerId2Price[v.PlayerId] = v.BuyPrice * 1000
					} else {
						playerId2Price[v.PlayerId] = v.BuyPrice
					}
				}

				sql2 = sql2 + "where playerId in ? and gsId = ? group by playerId,gsId"
				cut2s := make([]cut2, 0)
				db_global.Raw(sql2, playerIds, serverId).Scan(&cut2s)
				playerId2Times := make(map[uint32]uint32, 0)
				for _, v := range cut2s {
					playerId2Times[v.PlayerId] = v.RechargeTimes
				}

				type cut3 struct {
					PlayerId uint32 `gorm:"column:playerId"`
					NickName string `gorm:"column:nickName"`
				}
				cut3s := make([]cut3, 0)
				db_char := fusion.GetServerGormDB("db_char", int64(serverId))
				db_char.Table("inst_player_char").Select("ipcInstID as playerId,ipcNickName as nickName").
					Where("ipcInstID in ?", playerIds).Scan(&cut3s)
				playerId2Name := make(map[uint32]string, 0)
				for _, v := range cut3s {
					playerId2Name[v.PlayerId] = v.NickName
				}

				for k, v := range cut1s {
					temp := make(map[string]interface{}, 0)
					temp["Id"] = k + 1
					temp["PlayerId"] = v.PlayerId
					temp["ServerID"] = serverId
					temp["TotalPrice"] = v.BuyPrice
					name, err1 := playerId2Name[v.PlayerId]
					if !err1 {
						temp["PlayerName"] = ""
					} else {
						temp["PlayerName"] = name
					}
					times, err1 := playerId2Times[v.PlayerId]
					if !err1 {
						temp["TotalRechargeTimes"] = 0
					} else {
						temp["TotalRechargeTimes"] = times
					}
					data = append(data, temp)
				}
			} else if rankListType == 3 {
				type playerCut struct {
					PlayerId uint32 `gorm:"column:playerId"`
					NickName string `gorm:"column:nickName"`
					Power    int64  `gorm:"column:power"`
				}
				playerCuts := make([]playerCut, 0)
				db_char := fusion.GetServerGormDB("db_char", int64(serverId))

				db_char.Table("inst_player_char").
					Select("ipcInstID AS playerId,ipcNickName AS NickName,ipcS64Values->'$[7]' AS power ").
					Order("power DESC").Limit(50).Scan(&playerCuts)
				playerIds := make([]uint32, 0)
				for _, v := range playerCuts {
					playerIds = append(playerIds, v.PlayerId)
				}

				db_global := fusion.GetBaseGormDB("db_global")
				sql1, sql2 := fusion.CreateUnionTableSql("t_web_orders_", db_global, "where playerId in ? and gsId = ? ", " group by playerId,gsId,payId")

				cut1s := make([]cut1, 0)
				db_global.Raw(sql1, playerIds, serverId).Scan(&cut1s)
				playerId2Price := make(map[uint32]int64, 0)
				for _, v := range cut1s {
					if PlatForm4ID[int(v.PayId)] == 3 {
						playerId2Price[v.PlayerId] = v.BuyPrice * 1000
					} else {
						playerId2Price[v.PlayerId] = v.BuyPrice
					}
				}

				cut2s := make([]cut2, 0)
				db_global.Raw(sql2, playerIds, serverId).Scan(&cut2s)
				playerId2Times := make(map[uint32]uint32, 0)
				for _, v := range cut2s {
					playerId2Times[v.PlayerId] = v.RechargeTimes
				}

				for k, v := range playerCuts {
					temp := make(map[string]interface{}, 0)
					temp["Id"] = k + 1
					temp["PlayerId"] = v.PlayerId
					temp["PlayerName"] = v.NickName
					temp["ServerID"] = serverId
					temp["TotalPower"] = v.Power

					price, err1 := playerId2Price[v.PlayerId]
					if !err1 {
						temp["TotalPrice"] = 0
					} else {
						temp["TotalPrice"] = price
					}
					times, err1 := playerId2Times[v.PlayerId]
					if !err1 {
						temp["TotalRechargeTimes"] = 0
					} else {
						temp["TotalRechargeTimes"] = times
					}
					data = append(data, temp)
				}
			}
			min := math.Min(float64(param.PageInt*param.PageSizeInt), float64(len(data)))
			return data[((param.PageInt - 1) * param.PageSizeInt):int(min)], len(data)
		})
	return LevelRankList
}

func GetRankListAllServer(ctx *context.Context) table.Table {
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

	str := fusion.GetServerNameStrings()
	info := LevelRankList.GetInfo().HideNewButton().HideDeleteButton().HideDetailButton().HideEditButton().HideRowSelector()

	info.SetTitle(gotext.Get("全服累充排行"))

	info.AddField(gotext.Get("排名"), "Id", db.Int)
	info.AddField(gotext.Get("充值额度"), "TotalPrice", db.Int)
	info.AddField(gotext.Get("角色ID"), "PlayerId", db.Int)
	info.AddField(gotext.Get("角色名"), "PlayerName", db.Varchar)
	info.AddField(gotext.Get("充值次数"), "TotalRechargeTimes", db.Int)
	info.AddField(gotext.Get("所在服务器"), "ServerID", db.Int).
		FieldDisplay(func(model types.FieldModel) interface{} {
			temp, _ := strconv.Atoi(model.Value)
			return str[temp]
		})

	info.SetGetDataFn(
		func(param parameter.Parameters) (data []map[string]interface{}, size int) {
			var purchase = make([]map[string]interface{}, 0)

			db_world := fusion.GetBaseGormDB("db_world")
			db_world.Table("auto_purchase").
				Select("ID,platformID").Scan(&purchase)
			PlatForm4ID := map[int]int{}
			for _, v := range purchase {
				PlatFormId, _ := strconv.Atoi(fusion.Strval(v["platformID"]))
				id, _ := strconv.Atoi(fusion.Strval(v["ID"]))
				PlatForm4ID[id] = PlatFormId
			}

			type cut1 struct {
				PlayerId uint32 `gorm:"column:playerId"`
				GsId     uint32 `gorm:"column:gsId"`
				BuyPrice int64  `gorm:"column:buyPrice"`
				PayId    uint32 `gorm:"column:payId"`
			}
			type cut2 struct {
				PlayerId      uint32 `gorm:"column:playerId"`
				GsId          uint32 `gorm:"column:gsId"`
				RechargeTimes uint32 `gorm:"column:rechargeTimes"`
			}
			db_global := fusion.GetBaseGormDB("db_global")
			sql1, sql2 := fusion.CreateUnionTableSql("t_web_orders_", db_global)
			sql1 = sql1 + " group by playerId,gsId,payId order by buyPrice desc limit 50"
			cut1s := make([]cut1, 0)
			db_global.Raw(sql1).Scan(&cut1s)

			type tempCut struct {
				sort     uint32
				playerId uint32
			}

			playerId2Price := make(map[uint32]map[uint32]int64, 0)
			playerIdsAndSort := make(map[uint32][]tempCut, 0)
			playerIds := make(map[uint32][]uint32, 0)

			for k, v := range cut1s {
				_, err1 := playerIds[v.GsId]
				if !err1 {
					playerIdsAndSort[v.GsId] = make([]tempCut, 0)
				}
				playerIdsAndSort[v.GsId] = append(playerIdsAndSort[v.GsId], tempCut{sort: uint32(k), playerId: v.PlayerId})
				_, err2 := playerId2Price[v.GsId]
				if !err2 {
					playerId2Price[v.GsId] = make(map[uint32]int64, 0)
				}
				if PlatForm4ID[int(v.PayId)] == 3 {
					playerId2Price[v.GsId][v.PlayerId] = v.BuyPrice * 1000
				} else {
					playerId2Price[v.GsId][v.PlayerId] = v.BuyPrice
				}
				_, err3 := playerIds[v.GsId]
				if !err3 {
					playerIds[v.GsId] = make([]uint32, 0)
				}
				playerIds[v.GsId] = append(playerIds[v.GsId], v.PlayerId)
			}
			sql2 = sql2 + "where playerId in ? and gsId = ?  group by playerId,gsId"
			for k, v := range playerIdsAndSort {
				cut2s := make([]cut2, 0)
				db_global.Raw(sql2, playerIds[k], k).Scan(&cut2s)

				playerId2Times := make(map[uint32]uint32, 0)
				for _, v := range cut2s {
					playerId2Times[v.PlayerId] = v.RechargeTimes
				}
				type cut3 struct {
					PlayerId uint32 `gorm:"column:playerId"`
					NickName string `gorm:"column:nickName"`
				}
				cut3s := make([]cut3, 0)
				db_char := fusion.GetServerGormDB("db_char", int64(k))
				db_char.Table("inst_player_char").Select("ipcInstID as playerId,ipcNickName as nickName").
					Where("ipcInstID in ?", playerIds[k]).Scan(&cut3s)
				playerId2Name := make(map[uint32]string, 0)
				for _, _v := range cut3s {
					playerId2Name[_v.PlayerId] = _v.NickName
				}

				for _, _v := range v {
					temp := make(map[string]interface{}, 0)
					temp["ServerID"] = k
					temp["Id"] = _v.sort + 1
					temp["PlayerId"] = _v.playerId
					name, err1 := playerId2Name[_v.playerId]
					if !err1 {
						temp["PlayerName"] = ""
					} else {
						temp["PlayerName"] = name
					}
					times, err1 := playerId2Times[_v.playerId]
					if !err1 {
						temp["TotalRechargeTimes"] = 0
					} else {
						temp["TotalRechargeTimes"] = times
					}

					price4player, err1 := playerId2Price[k]
					if !err1 {
						temp["TotalPrice"] = 0
					} else {
						price, err2 := price4player[_v.playerId]
						if !err2 {
							temp["TotalPrice"] = 0
						} else {
							temp["TotalPrice"] = price
						}
					}
					data = append(data, temp)
				}
			}

			sort.Slice(data, func(i, j int) bool {
				temp1, _ := strconv.Atoi(fusion.Strval(data[i]["Id"]))
				temp2, _ := strconv.Atoi(fusion.Strval(data[j]["Id"]))
				return temp1 < temp2
			})

			min := math.Min(float64(param.PageInt*param.PageSizeInt), float64(len(data)))
			return data[((param.PageInt - 1) * param.PageSizeInt):int(min)], len(data)
		})
	return LevelRankList
}
