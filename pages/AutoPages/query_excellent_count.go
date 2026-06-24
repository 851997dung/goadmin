package AutoPages

import (
	"admin/common/def/excellentCountDef"
	"admin/fusion"
	"github.com/GoAdminGroup/go-admin/context"
	"github.com/GoAdminGroup/go-admin/modules/db"
	"github.com/GoAdminGroup/go-admin/plugins/admin/modules/parameter"
	"github.com/GoAdminGroup/go-admin/plugins/admin/modules/table"
	"github.com/GoAdminGroup/go-admin/template"
	"github.com/GoAdminGroup/go-admin/template/types"
	"github.com/GoAdminGroup/go-admin/template/types/action"
	"github.com/GoAdminGroup/go-admin/template/types/form"
	"github.com/fatih/structs"
	"github.com/leonelquinteros/gotext"
	"math"
	"sort"
	"strconv"
	"strings"
	"time"
)

func QueryExcellentCount(ctx *context.Context) table.Table {

	queryExcellentCount := table.NewDefaultTable(table.Config{
		Driver:     "mysql",
		Connection: "db_global",
		CanAdd:     true,
		Editable:   true,
		Deletable:  true,
		Exportable: true,
		PrimaryKey: table.PrimaryKey{
			Type: db.Int,
			Name: "PlayerId",
		},
	})

	info := queryExcellentCount.GetInfo().HideRowSelector().HideEditButton().HideDeleteButton().HideNewButton()
	myOP := fusion.GetServerList()
	info.AddField(gotext.Get("服务器ID"), "IpcServerID", db.Int).
		FieldFilterable(types.FilterType{FormType: form.SelectSingle}).
		FieldFilterOptions(myOP)
	info.AddField(gotext.Get("账号ID"), "AcctId", db.Int).
		FieldFilterable()
	info.AddField(gotext.Get("账号"), "AcctName", db.Varchar)
	info.AddField(gotext.Get("角色ID"), "PlayerId", db.Int).
		FieldFilterable()
	info.AddField(gotext.Get("角色昵称"), "PlayerName", db.Varchar).
		FieldFilterable()
	info.AddField(gotext.Get("角色等级"), "PlayerLevel", db.Int)
	info.AddField(gotext.Get("需求的卓越属性条数"), "NeedExAttrCount", db.Int).
		FieldFilterable(types.FilterType{HelpMsg: "该条必须填写，否则无法查询到结果"}).FieldHide()
	info.AddField(gotext.Get("满足条件的装备数量"), "EligibleEquipCount", db.Int).
		FieldSortable()
	info.AddField(gotext.Get("累计充值数量"), "RechargeCount", db.Bigint)
	info.AddField(gotext.Get("最近充值时间"), "LastRechargeTime", db.Datetime)
	info.AddField(gotext.Get("最近登录IP"), "LastLoginIP", db.Varchar)

	info.SetTable("t_ip_rules").SetTitle(gotext.Get("以卓越属性条数查询玩家")).HideDetailButton()
	info.SetGetDataFn(
		func(param parameter.Parameters) (data []map[string]interface{}, size int) {
			serverIdStr := ctx.Request.FormValue("IpcServerID")
			serverId, _ := strconv.Atoi(serverIdStr)
			acctId := ctx.Request.FormValue("AcctId")
			playerId := ctx.Request.FormValue("PlayerId")
			playerName := ctx.Request.FormValue("PlayerName")
			needExAttrCount, _ := strconv.Atoi(ctx.Request.FormValue("NeedExAttrCount"))

			if serverId == 0 || needExAttrCount == 0 {
				return nil, 0
			}

			where := ""
			if acctId != "" {
				where += " ipcAcctID = " + acctId + " and "
			}
			if playerId != "" {
				where += " ipcInstID = " + playerId + " and "
			}
			if playerName != "" {
				where += " ipcNickName like " + "'%" + playerName + "%'" + " and "
			}
			if serverIdStr != "" {
				where += " ipcServerID = " + serverIdStr + " and "
			}
			if where != "" {
				where = " where " + where
				where = strings.TrimSuffix(where, " and ")
			}

			db_char := fusion.GetServerGormDB("db_char", int64(serverId))
			if db_char == nil {
				return nil, 0
			}
			var playerInfos = make([]map[string]interface{}, 0)
			sql1 := "select ipcInstID,ipcAcctID,ipcNickName,ipcLevel,ipcStorageItems from inst_player_char " + where
			db_char.Raw(sql1).Scan(&playerInfos)

			CountDatas := make([]excellentCountDef.ExcellentCount, 0)

			var AcctIds = make([]int, 0)
			var playerIds = make([]int, 0)
			for _, playerInfo := range playerInfos {
				itemStr := fusion.Strval(playerInfo["ipcStorageItems"])
				arr := strings.Split(itemStr, ";")
				count := fusion.GetEligibleEquipCount(needExAttrCount, arr)
				if count != 0 {
					countData := excellentCountDef.ExcellentCount{}
					countData.IpcServerID = int64(serverId)
					countData.PlayerId, _ = strconv.Atoi(fusion.Strval(playerInfo["ipcInstID"]))
					countData.AcctId, _ = strconv.Atoi(fusion.Strval(playerInfo["ipcAcctID"]))
					countData.PlayerName = fusion.Strval(playerInfo["ipcNickName"])
					countData.PlayerLevel, _ = strconv.Atoi(fusion.Strval(playerInfo["ipcLevel"]))
					countData.EligibleEquipCount = count
					countData.NeedExAttrCount = needExAttrCount
					playerIds = append(playerIds, countData.PlayerId)
					AcctIds = append(AcctIds, countData.AcctId)
					CountDatas = append(CountDatas, countData)
				}
			}
			var acctCut = make([]map[string]interface{}, 0)

			type cut3 struct {
				name string
				ip   string
			}
			var AcctId2Name = make(map[int]cut3)
			db_global := fusion.GetBaseGormDB("db_global")
			db_global.Table("t_accounts").Select("Id,username,lastLoginIP").
				Where("Id in ?", AcctIds).Scan(&acctCut)

			for _, v := range acctCut {
				acctId, _ := strconv.Atoi(fusion.Strval(v["Id"]))
				acctName := fusion.Strval(v["username"])
				lastLoginIP := fusion.Strval(v["lastLoginIP"])
				AcctId2Name[acctId] = cut3{name: acctName, ip: lastLoginIP}
			}

			type cut1 struct {
				PlayerId   uint32    `gorm:"column:playerId"`
				BuyPrice   int64     `gorm:"column:buyPrice"`
				CreateTime time.Time `gorm:"column:createTime"`
			}

			type cut2 struct {
				BuyPrice   int64
				CreateTime time.Time
			}

			sql := fusion.GetSql4ExcellentAttr("t_web_orders_", db_global, "where playerId in ? and gsId = ? ", " group by playerId,gsId")
			cut1s := make([]cut1, 0)
			db_global.Raw(sql, playerIds, serverId).Scan(&cut1s)
			playerId2Price := make(map[uint32]cut2, 0)
			for _, v := range cut1s {
				playerId2Price[v.PlayerId] = cut2{BuyPrice: v.BuyPrice, CreateTime: v.CreateTime}
			}

			for _, countData := range CountDatas {
				cut2, err1 := playerId2Price[uint32(countData.PlayerId)]
				if !err1 {
					countData.RechargeCount = 0
				} else {
					countData.RechargeCount = uint64(cut2.BuyPrice)
					countData.LastRechargeTime = cut2.CreateTime.Format("2006-01-02 15:04:05")
				}

				cut3, err1 := AcctId2Name[countData.AcctId]
				if !err1 {
					countData.AcctName = ""
					countData.LastLoginIP = ""
				} else {
					countData.AcctName = cut3.name
					countData.LastLoginIP = cut3.ip
				}
				m3 := structs.Map(&countData)
				data = append(data, m3)
			}
			sort.Slice(data, func(i, j int) bool {
				tempI := data[i][param.SortField]
				tempi, err1 := strconv.Atoi(fusion.Strval(tempI))

				tempJ := data[j][param.SortField]
				tempj, err2 := strconv.Atoi(fusion.Strval(tempJ))
				if err1 != nil || err2 != nil {
					return false
				}
				if param.SortType == "desc" {
					return tempi > tempj
				} else {
					return tempi < tempj
				}
			})
			if strings.Index(ctx.Request.RequestURI, "/export/") != -1 {
				return data, len(data)
			} else {
				min := math.Min(float64(param.PageInt*param.PageSizeInt), float64(len(data)))
				return data[((param.PageInt - 1) * param.PageSizeInt):int(min)], len(data)
			}
		})
	info.SetExportProcessFn(func(param parameter.Parameters) (types.PanelInfo, error) {
		panelInfo, _ := queryExcellentCount.GetData(param.WithIsAll(param.IsAll()))
		return types.PanelInfo{Thead: panelInfo.Thead, InfoList: panelInfo.InfoList}, nil
	})

	info.AddActionButton(template.HTML(gotext.Get("玩家详情")), action.Jump("/admin/userMgr/PlayerDetail?IpcInstID={{(index .Value \"PlayerId\").Value}}&"+
		"IpcServerID={{(index .Value \"IpcServerID\").Value}}&IpcAcctID={{(index .Value \"AcctId\").Value}}"))

	return queryExcellentCount
}
