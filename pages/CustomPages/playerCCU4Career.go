package CustomPages

import (
	"admin/common"
	"admin/common/def"
	"admin/common/def/careerDef"
	"admin/fusion"
	"encoding/json"
	"github.com/GoAdminGroup/go-admin/context"
	"github.com/GoAdminGroup/go-admin/modules/db"
	"github.com/GoAdminGroup/go-admin/modules/logger"
	"github.com/GoAdminGroup/go-admin/plugins/admin/modules/parameter"
	"github.com/GoAdminGroup/go-admin/plugins/admin/modules/table"
	"github.com/GoAdminGroup/go-admin/template/types"
	"github.com/GoAdminGroup/go-admin/template/types/form"
	"github.com/leonelquinteros/gotext"
	"math"
	"net/url"
	"strconv"
	"time"
)

func GetPlayerCCU4Career(ctx *context.Context) table.Table {
	Diamond4RebornRec := table.NewDefaultTable(table.Config{
		Driver:     "mysql",
		Connection: "db_log",
		CanAdd:     true,
		Editable:   true,
		Deletable:  true,
		Exportable: true,
		PrimaryKey: table.PrimaryKey{
			Type: db.Int,
			Name: "ServerId",
		},
	})

	info := Diamond4RebornRec.GetInfo().
		HideEditButton().
		HideNewButton().
		HideRowSelector().
		HideDetailButton().
		HideDeleteButton().
		HideRowSelector().
		SetDefaultPageSize(20).
		SetPageSizeList([]int{20})

	serverList := fusion.GetServerList()

	info.AddField(gotext.Get("服务器ID"), "ServerId", db.Int).
		FieldFilterable(types.FilterType{FormType: form.SelectSingle}).
		FieldFilterOptions(serverList)
	info.AddField(gotext.Get("职业id"), "Career", db.Int)
	info.AddField(gotext.Get("服务器开启时间"), "SeverOpenTime", db.Varchar)
	info.AddField(gotext.Get("创建角色数量"), "CharCount", db.Int)

	info.AddField(gotext.Get("最大CCU"), "MaxCCU", db.Int)
	info.AddField(gotext.Get("平均CCU"), "AvrCCU", db.Int)
	info.AddField(gotext.Get("当前CCU"), "CurCCU", db.Int)

	info.AddField(gotext.Get("转生1-5次"), "Reset1", db.Int)
	info.AddField(gotext.Get("转生6-10次"), "Reset2", db.Int)
	info.AddField(gotext.Get("转生11-50次"), "Reset3", db.Int)
	info.AddField(gotext.Get("转生51-100次"), "Reset4", db.Int)
	info.AddField(gotext.Get("转生101-200次"), "Reset5", db.Int)
	info.AddField(gotext.Get("转生201-500次"), "Reset6", db.Int)
	info.AddField(gotext.Get("转生500次以上"), "Reset7", db.Int)

	info.SetTable("PlayerCCU4Career").SetTitle(gotext.Get("职业ccu"))
	info.SetGetDataFn(
		func(param parameter.Parameters) (data []map[string]interface{}, size int) {
			var ServerId = ctx.Request.FormValue("ServerId")
			var serverIdInt, _ = strconv.Atoi(ServerId)
			if ServerId == "" {
				return nil, 0
			}

			var OnlineNumRec4CareerInfos = make([]def.OnlineNumRec4Career, 0)
			var go_mananger = fusion.GetBaseGormDB("go_manager")
			go_mananger.Select("*").
				Table("player_online_num_record_4_Career").
				Scan(&OnlineNumRec4CareerInfos)

			var OnlineData = make(map[int]int)
			var OnlineMaxData = make(map[int]int)
			var countTimes = len(OnlineNumRec4CareerInfos)

			for _, v := range OnlineNumRec4CareerInfos {
				var data4Career = map[int]map[int]int{}
				var temp1 = fusion.String2Bytes(v.Record)
				json.Unmarshal(temp1, &data4Career)
				careerData := data4Career[serverIdInt]

				for k, online := range careerData {
					_, isOk := OnlineData[k]
					if !isOk {
						OnlineData[k] = 0
					}
					OnlineData[k] += online

					_, isOk = OnlineMaxData[k]
					if !isOk {
						OnlineMaxData[k] = 0
					}
					if OnlineMaxData[k] < online {
						OnlineMaxData[k] = online
					}
				}
			}

			db_global := fusion.GetBaseGormDB("db_global")
			var openTime time.Time
			db_global.Table("t_game_servers").Select("logicOpenTime").Where("Id=?", serverIdInt).Scan(&openTime)

			type HTTPRespBase struct {
				Err int    `json:"error"`
				Msg string `json:"message"`
			}
			type httpOnlineNum4CareerResp struct {
				HTTPRespBase
				OnlineNum4Career map[int]int `json:",omitempty"`
			}

			ctx1 := &context.Context{}
			vals := url.Values{}
			vals.Add("gsId", ServerId)
			err, res := fusion.CallToCenter(&vals, common.MyCfg.Api["GetGSOnlineNumber4Career"], "GET", ctx1)
			if err != nil {
				logger.Info("CountPlayerOnlineNum Failed, Because CenterServer has some error")
			}
			temp := fusion.String2Bytes(res)
			resp := httpOnlineNum4CareerResp{}
			json.Unmarshal(temp, &resp)

			type charCut struct {
				IpcCareer    int `gorm:"column:ipcCareer"`
				Num          int `gorm:"column:num"`
				IpcRebornNum int `gorm:"column:ipcRebornNum"`
			}
			var charData = make([]charCut, 0)
			var charCount = make(map[int]int)
			db_char := fusion.GetServerGormDB("db_char", int64(serverIdInt))
			db_char.Table("inst_player_char").
				Select("ipcCareer,count(*)as num,ipcRebornNum").
				Group("ipcCareer,ipcRebornNum").
				Scan(&charData)

			var num4Reborn = make(map[int]map[int]int)
			for _, temp := range charData {
				_, isOk := charCount[temp.IpcCareer]
				if !isOk {
					charCount[temp.IpcCareer] = 0
				}
				charCount[temp.IpcCareer] += temp.Num

				if num4Reborn[temp.IpcCareer] == nil {
					num4Reborn[temp.IpcCareer] = make(map[int]int)
				}
				if 1 <= temp.IpcRebornNum && temp.IpcRebornNum <= 5 {
					_, isOk := num4Reborn[temp.IpcCareer][0]
					if !isOk {
						num4Reborn[temp.IpcCareer][0] = 0
					}
					num4Reborn[temp.IpcCareer][0] += temp.Num
				}
				if 6 <= temp.IpcRebornNum && temp.IpcRebornNum <= 10 {
					_, isOk := num4Reborn[temp.IpcCareer][1]
					if !isOk {
						num4Reborn[temp.IpcCareer][1] = 0
					}
					num4Reborn[temp.IpcCareer][1] += temp.Num
				}
				if 11 <= temp.IpcRebornNum && temp.IpcRebornNum <= 50 {
					_, isOk := num4Reborn[temp.IpcCareer][2]
					if !isOk {
						num4Reborn[temp.IpcCareer][2] = 0
					}
					num4Reborn[temp.IpcCareer][2] += temp.Num
				}
				if 51 <= temp.IpcRebornNum && temp.IpcRebornNum <= 100 {
					_, isOk := num4Reborn[temp.IpcCareer][3]
					if !isOk {
						num4Reborn[temp.IpcCareer][3] = 0
					}
					num4Reborn[temp.IpcCareer][3] += temp.Num
				}
				if 101 <= temp.IpcRebornNum && temp.IpcRebornNum <= 200 {
					_, isOk := num4Reborn[temp.IpcCareer][4]
					if !isOk {
						num4Reborn[temp.IpcCareer][4] = 0
					}
					num4Reborn[temp.IpcCareer][4] += temp.Num
				}
				if 201 <= temp.IpcRebornNum && temp.IpcRebornNum <= 500 {
					_, isOk := num4Reborn[temp.IpcCareer][5]
					if !isOk {
						num4Reborn[temp.IpcCareer][5] = 0
					}
					num4Reborn[temp.IpcCareer][5] += temp.Num
				}
				if 501 <= temp.IpcRebornNum {
					_, isOk := num4Reborn[temp.IpcCareer][6]
					if !isOk {
						num4Reborn[temp.IpcCareer][6] = 0
					}
					num4Reborn[temp.IpcCareer][6] += temp.Num
				}
			}

			for _, value := range careerDef.Careers {
				var tempMap = make(map[string]interface{})
				tempMap["ServerId"] = serverIdInt
				tempMap["Career"] = value
				tempMap["SeverOpenTime"] = openTime.Format("2006-01-02 15:04:05")
				tempMap["CharCount"] = charCount[value]
				tempMap["MaxCCU"] = OnlineMaxData[value]
				tempMap["AvrCCU"] = OnlineData[value] / int(math.Max(1, float64(countTimes)))
				if resp.OnlineNum4Career[value] != 0 {
					tempMap["CurCCU"] = resp.OnlineNum4Career[value]
				} else {
					tempMap["CurCCU"] = 0
				}
				if num4Reborn[value] != nil {
					_, isOk := num4Reborn[value][0]
					if isOk {
						tempMap["Reset1"] = num4Reborn[value][0]
					} else {
						tempMap["Reset1"] = 0
					}

					_, isOk = num4Reborn[value][1]
					if isOk {
						tempMap["Reset2"] = num4Reborn[value][1]
					} else {
						tempMap["Reset2"] = 0
					}

					_, isOk = num4Reborn[value][2]
					if isOk {
						tempMap["Reset3"] = num4Reborn[value][2]
					} else {
						tempMap["Reset3"] = 0
					}

					_, isOk = num4Reborn[value][3]
					if isOk {
						tempMap["Reset4"] = num4Reborn[value][3]
					} else {
						tempMap["Reset4"] = 0
					}

					_, isOk = num4Reborn[value][4]
					if isOk {
						tempMap["Reset5"] = num4Reborn[value][4]
					} else {
						tempMap["Reset5"] = 0
					}

					_, isOk = num4Reborn[value][5]
					if isOk {
						tempMap["Reset6"] = num4Reborn[value][5]
					} else {
						tempMap["Reset6"] = 0
					}

					_, isOk = num4Reborn[value][6]
					if isOk {
						tempMap["Reset7"] = num4Reborn[value][6]
					} else {
						tempMap["Reset7"] = 0
					}
				} else {
					tempMap["Reset1"] = 0
					tempMap["Reset2"] = 0
					tempMap["Reset3"] = 0
					tempMap["Reset4"] = 0
					tempMap["Reset5"] = 0
					tempMap["Reset6"] = 0
					tempMap["Reset7"] = 0
				}
				data = append(data, tempMap)
			}
			return data, len(data)
		})

	info.SetExportProcessFn(func(param parameter.Parameters) (types.PanelInfo, error) {
		panelInfo, _ := Diamond4RebornRec.GetData(param.WithIsAll(param.IsAll()))
		return types.PanelInfo{Thead: panelInfo.Thead, InfoList: panelInfo.InfoList}, nil
	})
	return Diamond4RebornRec
}
