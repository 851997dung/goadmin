package CustomPages

import (
	"admin/common"
	"admin/common/def"
	"admin/fusion"
	"encoding/json"
	"math"
	"net/url"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/GoAdminGroup/go-admin/context"
	"github.com/GoAdminGroup/go-admin/modules/db"
	"github.com/GoAdminGroup/go-admin/modules/logger"
	"github.com/GoAdminGroup/go-admin/plugins/admin/modules/parameter"
	"github.com/GoAdminGroup/go-admin/plugins/admin/modules/table"
	"github.com/GoAdminGroup/go-admin/template/types"
	"github.com/GoAdminGroup/go-admin/template/types/form"
	"github.com/fatih/structs"
	"github.com/leonelquinteros/gotext"
)

func GetRealTimeOnline(ctx *context.Context) table.Table {
	tRealTimeOnline := table.NewDefaultTable(table.Config{
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
	info := tRealTimeOnline.GetInfo()
	str := fusion.GetServerNameStrings()
	info.AddField(gotext.Get("服务器ID"), "Id", db.Int).
		FieldDisplay(func(model types.FieldModel) interface{} {
			temp, _ := strconv.Atoi(model.Value)
			return str[temp]
		})
	info.AddField(gotext.Get("实时在线人数"), "OnlineNum", db.Int)

	info.SetTitle(gotext.Get("实时在线人数")).
		HideEditButton().HideNewButton().HideRowSelector().HideDetailButton().HideDeleteButton()

	info.SetGetDataFn(
		func(param parameter.Parameters) (data []map[string]interface{}, size int) {
			serverIds := fusion.GetActiveServerIds()
			totalOnline := 0
			for _, v := range serverIds {
				ctx1 := &context.Context{}
				vals := url.Values{}
				vals.Add("gsId", strconv.Itoa(v))
				err, res := fusion.CallToCenter(&vals, common.MyCfg.Api["getGSOnlineNumber"], "GET", ctx1)
				if err != nil {
					logger.Info("GetRealTimeOnline Failed, Because CenterServer has some error")
				}
				var gsOnlineNum int
				gsOnlineNum, err = strconv.Atoi(res)
				if err != nil {
					logger.Infof("GetRealTimeOnline Failed,serverID = %d,With %s", v, res)
				}
				temp := make(map[string]interface{})
				temp["Id"] = v
				temp["OnlineNum"] = gsOnlineNum
				data = append(data, temp)
				totalOnline += gsOnlineNum
			}
			temp := make(map[string]interface{})
			temp["Id"] = 0
			temp["OnlineNum"] = totalOnline
			data = append(data, temp)
			sort.Slice(data, func(i, j int) bool {
				idI := fusion.Strval(data[i]["Id"])
				idJ := fusion.Strval(data[j]["Id"])
				return idI < idJ
			})

			if param.IsAll() {
				return data, len(data)
			} else {
				max := math.Min(float64(len(data)), float64((param.PageInt)*param.PageSizeInt))
				return data[(param.PageInt-1)*param.PageSizeInt : int(max)], len(data)
			}
		})
	info.SetExportProcessFn(func(param parameter.Parameters) (types.PanelInfo, error) {
		panelInfo, _ := tRealTimeOnline.GetData(param.WithIsAll(param.IsAll()))
		return types.PanelInfo{Thead: panelInfo.Thead, InfoList: panelInfo.InfoList}, nil
	})
	return tRealTimeOnline
}

func GetOnlineNumByHour(ctx *context.Context) table.Table {
	OnlineNumByHour := table.NewDefaultTable(table.Config{
		Driver:     "mysql",
		Connection: "db_global",
		CanAdd:     true,
		Editable:   true,
		Deletable:  true,
		Exportable: true,
		PrimaryKey: table.PrimaryKey{
			Type: db.Varchar,
			Name: "date",
		},
	})
	info := OnlineNumByHour.GetInfo().HideEditButton().HideNewButton().HideRowSelector().HideDetailButton().HideDeleteButton()

	myOp := fusion.GetServerList()
	info.AddField(gotext.Get("服务器ID"), "serverId", db.Int).
		FieldFilterable(types.FilterType{FormType: form.SelectSingle}).
		FieldFilterOptions(myOp)
	info.AddField(gotext.Get("日期"), "date", db.Date).
		FieldFilterable(types.FilterType{FormType: form.DateRange})
	info.AddField(gotext.Get("当日最大在线"), "maxNumDay", db.Int)
	for i := 0; i <= 23; i++ {
		info.AddField(strconv.Itoa(i), strconv.Itoa(i), db.Int).
			FieldDisplay(func(model types.FieldModel) interface{} {
				_, err := strconv.Atoi(model.Value)
				if err != nil {
					return "0"
				}
				return model.Value
			})
	}
	info.SetTitle(gotext.Get("分时在线"))
	info.SetGetDataFn(
		func(param parameter.Parameters) (data []map[string]interface{}, size int) {
			serverId, _ := strconv.Atoi(ctx.Request.FormValue("serverId"))
			dateStart := ctx.Request.FormValue("date_start__goadmin")
			dateEnd := ctx.Request.FormValue("date_end__goadmin")
			if serverId == 0 || dateStart == "" || dateEnd == "" {
				return nil, 0
			}

			go_manager := fusion.GetBaseGormDB("go_manager")
			onlineNumRec := []def.OnlineNumRec{}
			go_manager.Table("player_online_num_record").
				Where("logTime >= ?", dateStart).
				Where("logTime <= ?", dateEnd).
				Order("logTime ASC").
				Scan(&onlineNumRec)

			startTime, _ := time.ParseInLocation("2006-01-02", dateStart, time.Local)
			endTime, _ := time.ParseInLocation("2006-01-02", dateEnd, time.Local)

			for i := startTime; i.Before(endTime); i = i.AddDate(0, 0, 1) {
				aimDate := i.Format("2006-01-02")
				jsonData := map[int]int{}
				temp := make(map[string]interface{})
				temp["serverId"] = serverId
				temp["date"] = aimDate
				var max = 0

				for _, v := range onlineNumRec {
					logDate := time.Time(v.LogTime).Format("2006-01-02")
					if logDate == aimDate {
						b := fusion.String2Bytes(v.Record)
						err := json.Unmarshal(b, &jsonData)
						if err != nil {
							continue
						}
						hour, _ := strconv.Atoi(v.LogTime.Hour())
						temp[strconv.Itoa(hour)] = jsonData[serverId]
						if max < jsonData[serverId] {
							max = jsonData[serverId]
						}
					}
				}
				temp["maxNumDay"] = max
				data = append(data, temp)
			}

			if param.IsAll() {
				return data, len(data)
			} else {
				max := math.Min(float64(len(data)), float64((param.PageInt)*param.PageSizeInt))
				return data[(param.PageInt-1)*param.PageSizeInt : int(max)], len(data)
			}
		})

	info.SetExportProcessFn(func(param parameter.Parameters) (types.PanelInfo, error) {
		panelInfo, _ := OnlineNumByHour.GetData(param.WithIsAll(param.IsAll()))
		return types.PanelInfo{Thead: panelInfo.Thead, InfoList: panelInfo.InfoList}, nil
	})

	return OnlineNumByHour
}

func GetSingleSeverCount(ctx *context.Context) table.Table {
	SingleSeverCount := table.NewDefaultTable(table.Config{
		Driver:     "mysql",
		Connection: "go_manager",
		CanAdd:     true,
		Editable:   true,
		Deletable:  true,
		Exportable: true,
		PrimaryKey: table.PrimaryKey{
			Type: db.Int,
			Name: "Id",
		},
	})
	info := SingleSeverCount.GetInfo().HideEditButton().HideNewButton().HideRowSelector().HideDetailButton().HideDeleteButton()

	myOp := fusion.GetServerList()
	info.AddField("Id", "Id", db.Int)
	info.AddField(gotext.Get("服务器ID"), "GsId", db.Int).
		FieldFilterable(types.FilterType{FormType: form.SelectSingle}).
		FieldFilterOptions(myOp).
		FieldHide()
	info.AddField(gotext.Get("在线峰值"), "MaxOnline", db.Int)
	info.AddField(gotext.Get("在线谷值"), "MinOnline", db.Int)
	info.AddField(gotext.Get("在线均值"), "AverageOnline", db.Int)
	info.AddField(gotext.Get("新增账号数"), "NewAccountNum", db.Int)
	info.AddField(gotext.Get("总账号数"), "AllAccountNum", db.Int)
	info.AddField(gotext.Get("新增充值玩家"), "NewRNum", db.Int)
	info.AddField(gotext.Get("当天充值玩家数"), "CurDayRNum", db.Int)
	info.AddField(gotext.Get("总充值玩家数"), "CountRNum", db.Int)
	info.AddField(gotext.Get("当天总充值人次"), "CurDayManTimesR", db.Int)
	info.AddField(gotext.Get("记录时间"), "LogTime", db.Datetime).
		FieldFilterable(types.FilterType{FormType: form.DateRange})

	info.SetTitle(gotext.Get("单服数据"))

	info.SetGetDataFn(func(param parameter.Parameters) (data []map[string]interface{}, size int) {
		serverId := ctx.Request.FormValue("GsId")
		logTimeStart := ctx.Request.FormValue("LogTime_start__goadmin")
		logTimeEnd := ctx.Request.FormValue("LogTime_end__goadmin")
		go_manager := fusion.GetBaseGormDB("go_manager")

		fields := fusion.GetFields2Struct(def.SingleServerDataCount{})
		sql1 := "select " + fields + " from " + def.SingleServerDataCount{}.TableName()
		sql2 := "select count(*) from " + def.SingleServerDataCount{}.TableName()
		where := ""
		if serverId != "" {
			where += " gsId = " + serverId + " and "
		}
		if logTimeStart != "" {
			where += " logTime >= '" + logTimeStart + "' and "
		}
		if logTimeEnd != "" {
			where += " logTime <= '" + logTimeEnd + "' and "
		}
		if where != "" {
			where = " where " + where
			where = strings.TrimSuffix(where, " and ")
		} else {
			return nil, 0
		}

		sql1 = sql1 + where
		if !param.IsAll() {
			if param.SortField != "" {
				if param.SortField == "PlayerLevel" {
					param.SortField = "CharacterLevel"
				}
				sql1 += " ORDER BY " + param.SortField + " " + param.SortType
			}
			sql1 += "  LIMIT " + strconv.Itoa(param.PageSizeInt) + " OFFSET " + strconv.Itoa((param.PageInt-1)*param.PageSizeInt)
		}
		sql2 = sql2 + where

		var tempCount []int64
		singleData := make([]def.SingleServerDataCount, 0)

		go_manager.Raw(sql1).Scan(&singleData)

		go_manager.Raw(sql2).Scan(&tempCount)

		for i, _ := range singleData {
			singleData[i].LogTime = strings.Replace(singleData[i].LogTime, "T", " ", -1)
			singleData[i].LogTime = strings.Replace(singleData[i].LogTime, "Z", " ", -1)
		}
		temp := []map[string]interface{}{}
		for _, v := range singleData {
			m3 := structs.Map(&v)
			temp = append(temp, m3)
		}
		var count int64
		for _, val := range tempCount {
			count += val
		}
		return temp, int(count)
	})

	info.SetExportProcessFn(func(param parameter.Parameters) (types.PanelInfo, error) {
		panelInfo, _ := SingleSeverCount.GetData(param.WithIsAll(param.IsAll()))
		return types.PanelInfo{Thead: panelInfo.Thead, InfoList: panelInfo.InfoList}, nil
	})

	return SingleSeverCount
}

func GetDeviceInfo(ctx *context.Context) table.Table {
	DeviceInfo := table.NewDefaultTable(table.Config{
		Driver:     "mysql",
		Connection: "go_manager",
		CanAdd:     true,
		Editable:   true,
		Deletable:  true,
		Exportable: true,
		PrimaryKey: table.PrimaryKey{
			Type: db.Varchar,
			Name: "deviceModel",
		},
	})
	info := DeviceInfo.GetInfo().HideEditButton().HideNewButton().HideRowSelector().HideDetailButton().HideDeleteButton().HideRowSelector()
	info.AddField(gotext.Get("设备类型"), "deviceModel", db.Varchar)
	info.AddField(gotext.Get("人数"), "num", db.Int)

	info.SetTitle(gotext.Get("设备信息"))

	info.SetGetDataFn(func(param parameter.Parameters) (data []map[string]interface{}, size int) {
		db_global := fusion.GetBaseGormDB("db_global")

		db_global.Raw("SELECT deviceModel,COUNT(deviceModel)AS num FROM  t_accounts GROUP BY deviceModel").Scan(&data)
		if param.IsAll() {
			return data, len(data)
		} else {
			max := math.Min(float64(len(data)), float64((param.PageInt)*param.PageSizeInt))
			return data[(param.PageInt-1)*param.PageSizeInt : int(max)], len(data)
		}
	})

	info.SetExportProcessFn(func(param parameter.Parameters) (types.PanelInfo, error) {
		panelInfo, _ := DeviceInfo.GetData(param.WithIsAll(param.IsAll()))
		return types.PanelInfo{Thead: panelInfo.Thead, InfoList: panelInfo.InfoList}, nil
	})

	return DeviceInfo
}

func GetCCUByDate(ctx *context.Context) table.Table {
	CCUByDate := table.NewDefaultTable(table.Config{
		Driver:     "mysql",
		Connection: "db_global",
		CanAdd:     true,
		Editable:   true,
		Deletable:  true,
		Exportable: true,
		PrimaryKey: table.PrimaryKey{
			Type: db.Int,
			Name: "serverId",
		},
	})

	serverList := fusion.GetServerListWithAll()
	info := CCUByDate.GetInfo().HideEditButton().HideNewButton().HideRowSelector().HideDetailButton().HideDeleteButton().HideRowSelector()
	str := fusion.GetServerNameStrings()

	info.SetTitle(gotext.Get("单服ccu"))
	info.AddField(gotext.Get("服务器Id"), "serverId", db.Int).
		FieldFilterable(types.FilterType{FormType: form.SelectSingle, Options: serverList}).
		FieldDisplay(func(model types.FieldModel) interface{} {
			temp, _ := strconv.Atoi(model.Value)
			return str[temp]
		})

	info.AddField(gotext.Get("查询类型"), "queryType", db.Int).FieldHide().
		FieldFilterable(types.FilterType{FormType: form.SelectSingle, Operator: types.FilterOperatorEqual}).
		FieldFilterOptions(types.FieldOptions{
			{Text: "按月查询", Value: "0"},
			{Text: "按年查询", Value: "1"},
		})

	info.AddField(gotext.Get("日期"), "date", db.Date).FieldFilterable(types.FilterType{FormType: form.Date})

	queryType, _ := strconv.Atoi(ctx.FormValue("queryType"))
	switch queryType {
	case 0:
		info.AddField(gotext.Get("月最大在线"), "maxNumMonth", db.Int)
		info.AddField(gotext.Get("第一周最大在线"), "maxNumWeek1", db.Int)
		info.AddField(gotext.Get("第二周最大在线"), "maxNumWeek2", db.Int)
		info.AddField(gotext.Get("第三周最大在线"), "maxNumWeek3", db.Int)
		info.AddField(gotext.Get("第四周最大在线"), "maxNumWeek4", db.Int)
		info.AddField(gotext.Get("第五周最大在线"), "maxNumWeek5", db.Int)
		break
	case 1:
		info.AddField(gotext.Get("第一季度最大在线"), "maxNumQuarter1", db.Int)
		info.AddField(gotext.Get("第二季度最大在线"), "maxNumQuarter2", db.Int)
		info.AddField(gotext.Get("第三季度最大在线"), "maxNumQuarter3", db.Int)
		info.AddField(gotext.Get("第四季度最大在线"), "maxNumQuarter4", db.Int)
		for i := 1; i <= 12; i++ {
			info.AddField(strconv.Itoa(i)+gotext.Get("月最大在线"), "maxNumMonth"+strconv.Itoa(i), db.Int)
		}
		break
	}

	info.SetGetDataFn(
		func(param parameter.Parameters) (data []map[string]interface{}, size int) {
			serverId, _ := strconv.Atoi(ctx.Request.FormValue("serverId"))
			queryType, _ := strconv.Atoi(ctx.FormValue("queryType"))
			date := ctx.Request.FormValue("date")

			if date == "" {
				return nil, 0
			}

			go_manager := fusion.GetBaseGormDB("go_manager")
			onlineNumRec := []def.OnlineNumRec{}

			gsIdx := fusion.GetActiveServerIds()

			switch queryType {
			case 0:
				startTime, _ := time.ParseInLocation("2006-01-02", date, time.Local)
				endTime := startTime.AddDate(0, 1, 0)
				start := startTime.Format("2006-01") + "-01"
				end := endTime.Format("2006-01") + "-01"

				go_manager.Table("player_online_num_record").
					Where("logTime >= ?", start).
					Where("logTime < ?", end).
					Order("logTime ASC").
					Scan(&onlineNumRec)

				if serverId != 0 {
					jsonData := map[int]int{}
					temp := make(map[string]interface{})
					temp["serverId"] = serverId
					temp["date"] = startTime.Format("2006-01")
					var monthMax = 0
					var weekMax = 0
					var day = 0
					var lastDay = 0
					var week = 1

					temp["maxNumWeek1"] = 0
					temp["maxNumWeek2"] = 0
					temp["maxNumWeek3"] = 0
					temp["maxNumWeek4"] = 0
					temp["maxNumWeek5"] = 0

					for _, v := range onlineNumRec {

						b := fusion.String2Bytes(v.Record)
						err := json.Unmarshal(b, &jsonData)
						if err != nil {
							continue
						}

						curDay, _ := strconv.Atoi(time.Time(v.LogTime).Format("02"))
						if curDay != lastDay {
							lastDay = curDay
							day++
							if day == 7 {
								week++
								weekMax = 0
								day = 0
							}
						}

						if weekMax < jsonData[serverId] {
							weekMax = jsonData[serverId]
							temp["maxNumWeek"+strconv.Itoa(week)] = weekMax
						}
						if monthMax < jsonData[serverId] {
							monthMax = jsonData[serverId]
						}

					}
					temp["maxNumMonth"] = monthMax
					data = append(data, temp)
				} else {
					for _, gsid := range gsIdx {
						jsonData := map[int]int{}
						temp := make(map[string]interface{})
						temp["serverId"] = gsid
						temp["date"] = startTime.Format("2006-01")
						var monthMax = 0
						var weekMax = 0
						var day = 0
						var lastDay = 0
						var week = 1

						temp["maxNumWeek1"] = 0
						temp["maxNumWeek2"] = 0
						temp["maxNumWeek3"] = 0
						temp["maxNumWeek4"] = 0
						temp["maxNumWeek5"] = 0

						for _, v := range onlineNumRec {

							b := fusion.String2Bytes(v.Record)
							err := json.Unmarshal(b, &jsonData)
							if err != nil {
								continue
							}

							curDay, _ := strconv.Atoi(time.Time(v.LogTime).Format("02"))
							if curDay != lastDay {
								lastDay = curDay
								day++
								if day == 7 {
									week++
									weekMax = 0
									day = 0
								}
							}

							if weekMax < jsonData[gsid] {
								weekMax = jsonData[gsid]
								temp["maxNumWeek"+strconv.Itoa(week)] = weekMax
							}
							if monthMax < jsonData[gsid] {
								monthMax = jsonData[gsid]
							}
						}
						temp["maxNumMonth"] = monthMax
						data = append(data, temp)
					}
				}
				break
			case 1:
				startTime, _ := time.ParseInLocation("2006-01-02", date, time.Local)
				endTime := startTime.AddDate(1, 0, 0)
				start := startTime.Format("2006") + "-01-01"
				end := endTime.Format("2006") + "-01-01"

				go_manager.Table("player_online_num_record").
					Where("logTime >= ?", start).
					Where("logTime < ?", end).
					Order("logTime ASC").
					Scan(&onlineNumRec)
				if serverId != 0 {
					jsonData := map[int]int{}
					temp := make(map[string]interface{})
					temp["serverId"] = serverId
					temp["date"] = startTime.Format("2006")
					var quarterMax = 0
					var startMonth = 1
					var quarter = 1
					var monthMax = 0

					temp["maxNumQuarter1"] = 0
					temp["maxNumQuarter2"] = 0
					temp["maxNumQuarter3"] = 0
					temp["maxNumQuarter4"] = 0
					for i := 1; i <= 12; i++ {
						temp["maxNumMonth"+strconv.Itoa(i)] = 0
					}

					for _, v := range onlineNumRec {

						b := fusion.String2Bytes(v.Record)
						err := json.Unmarshal(b, &jsonData)
						if err != nil {
							continue
						}
						tempQuarter := fusion.GetDataQuarter(time.Time(v.LogTime))
						if tempQuarter != quarter {
							quarter = tempQuarter
							quarterMax = 0
						}

						if quarterMax < jsonData[serverId] {
							quarterMax = jsonData[serverId]
							temp["maxNumQuarter"+strconv.Itoa(quarter)] = quarterMax
						}

						monthNum, _ := strconv.Atoi(time.Time(v.LogTime).Format("01"))
						if monthNum != startMonth {
							monthMax = 0
							startMonth = monthNum
						}

						if monthMax < jsonData[serverId] {
							monthMax = jsonData[serverId]
							temp["maxNumMonth"+strconv.Itoa(monthNum)] = monthMax
						}
					}
					data = append(data, temp)
				} else {
					for _, gsid := range gsIdx {
						jsonData := map[int]int{}
						temp := make(map[string]interface{})
						temp["serverId"] = gsid
						temp["date"] = startTime.Format("2006")
						var quarterMax = 0
						var startMonth = 1
						var quarter = 1
						var monthMax = 0

						temp["maxNumQuarter1"] = 0
						temp["maxNumQuarter2"] = 0
						temp["maxNumQuarter3"] = 0
						temp["maxNumQuarter4"] = 0
						for i := 1; i <= 12; i++ {
							temp["maxNumMonth"+strconv.Itoa(i)] = 0
						}

						for _, v := range onlineNumRec {

							b := fusion.String2Bytes(v.Record)
							err := json.Unmarshal(b, &jsonData)
							if err != nil {
								continue
							}
							tempQuarter := fusion.GetDataQuarter(time.Time(v.LogTime))
							if tempQuarter != quarter {
								quarter = tempQuarter
								quarterMax = 0
							}

							if quarterMax < jsonData[gsid] {
								quarterMax = jsonData[gsid]
								temp["maxNumQuarter"+strconv.Itoa(quarter)] = quarterMax
							}

							monthNum, _ := strconv.Atoi(time.Time(v.LogTime).Format("01"))
							if monthNum != startMonth {
								monthMax = 0
								startMonth = monthNum
							}

							if monthMax < jsonData[gsid] {
								monthMax = jsonData[gsid]
								temp["maxNumMonth"+strconv.Itoa(monthNum)] = monthMax
							}
						}
						data = append(data, temp)
					}
				}
				break
			}

			if param.IsAll() {
				return data, len(data)
			} else {
				max := math.Min(float64(len(data)), float64((param.PageInt)*param.PageSizeInt))
				return data[(param.PageInt-1)*param.PageSizeInt : int(max)], len(data)
			}
			return nil, 0
		})

	info.SetExportProcessFn(func(param parameter.Parameters) (types.PanelInfo, error) {
		panelInfo, _ := CCUByDate.GetData(param.WithIsAll(param.IsAll()))
		return types.PanelInfo{Thead: panelInfo.Thead, InfoList: panelInfo.InfoList}, nil
	})

	return CCUByDate
}

func GetAllServerCCUOneDay(ctx *context.Context) table.Table {
	AllServerCCUOneDay := table.NewDefaultTable(table.Config{
		Driver:     "mysql",
		Connection: "db_global",
		CanAdd:     true,
		Editable:   true,
		Deletable:  true,
		Exportable: true,
		PrimaryKey: table.PrimaryKey{
			Type: db.Int,
			Name: "serverId",
		},
	})

	info := AllServerCCUOneDay.GetInfo().HideEditButton().HideNewButton().HideRowSelector().HideDetailButton().HideDeleteButton().HideRowSelector()
	str := fusion.GetServerNameStrings()

	info.SetTitle(gotext.Get("全服ccu"))
	info.AddField(gotext.Get("服务器Id"), "serverId", db.Int).FieldDisplay(func(model types.FieldModel) interface{} {
		temp, _ := strconv.Atoi(model.Value)
		return str[temp]
	})

	info.AddField(gotext.Get("日期"), "date", db.Date).FieldFilterable(types.FilterType{FormType: form.Date})
	info.AddField(gotext.Get("CCU"), "CCU", db.Int)

	info.SetGetDataFn(
		func(param parameter.Parameters) (data []map[string]interface{}, size int) {
			date := ctx.Request.FormValue("date")
			if date == "" {
				return nil, 0
			}

			go_manager := fusion.GetBaseGormDB("go_manager")
			onlineNumRec := []def.OnlineNumRec{}
			go_manager.Table("player_online_num_record").
				Where("logTime like ?", date+"%").
				Order("logTime ASC").
				Scan(&onlineNumRec)

			gsIdxs := fusion.GetActiveServerIds()
			for _, gsId := range gsIdxs {
				var serverMax = 0
				jsonData := map[int]int{}
				temp := make(map[string]interface{})

				temp["serverId"] = gsId
				temp["date"] = date

				for _, v := range onlineNumRec {
					b := fusion.String2Bytes(v.Record)
					err := json.Unmarshal(b, &jsonData)
					if err != nil {
						continue
					}

					for id, val := range jsonData {
						if gsId == id {
							if serverMax <= val {
								serverMax = val
							}
						}
					}

				}
				temp["CCU"] = serverMax
				data = append(data, temp)
			}

			if param.IsAll() {
				return data, len(data)
			} else {
				max := math.Min(float64(len(data)), float64((param.PageInt)*param.PageSizeInt))
				return data[(param.PageInt-1)*param.PageSizeInt : int(max)], len(data)
			}
			return nil, 0
		})

	info.SetExportProcessFn(func(param parameter.Parameters) (types.PanelInfo, error) {
		panelInfo, _ := AllServerCCUOneDay.GetData(param.WithIsAll(param.IsAll()))
		return types.PanelInfo{Thead: panelInfo.Thead, InfoList: panelInfo.InfoList}, nil
	})
	return AllServerCCUOneDay
}
