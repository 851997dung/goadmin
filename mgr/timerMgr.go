package mgr

import (
	"admin/SendEMail"
	"admin/common"
	"admin/common/def"
	"admin/common/def/flowDef"
	"admin/fusion"
	"admin/fusion/timer"
	"database/sql"
	"encoding/json"
	"math"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/GoAdminGroup/go-admin/context"
	"github.com/GoAdminGroup/go-admin/modules/logger"
	"github.com/leonelquinteros/gotext"
)

var DaySecond = 24 * 60 * 60

type StatisticsTimer struct {
	timer.WheelTriggerOwner
}

func NewTimeMgrTimer() *StatisticsTimer {
	recorder := new(StatisticsTimer)
	recorder.InitVT(recorder)
	return recorder
}

func (recorder *StatisticsTimer) Impl_GetWheelTimerMgr() *timer.WheelTimerMgr {
	return fusion.MainService.TimerMgr
}

func CountPlayerOnlineNum() {
	gsIdx := fusion.GetActiveServerIds()
	//fusion.GoPool.Submit(func() {
	logger.Info("Start CountPlayerOnlineNum")
	onlineRecords := []def.OnlineRecord{}
	for _, v := range gsIdx {
		ctx1 := &context.Context{}
		vals := url.Values{}
		vals.Add("gsId", strconv.Itoa(v))
		err, res := fusion.CallToCenter(&vals, common.MyCfg.Api["getGSOnlineNumber"], "GET", ctx1)
		if err != nil {
			logger.Info("CountPlayerOnlineNum Failed, Because CenterServer has some error")
			continue
		}
		var gsOnlineNum int
		gsOnlineNum, err = strconv.Atoi(res)
		if err != nil {
			logger.Infof("CountPlayerOnlineNum Failed,serverID = %d,With %s", v, res)
			continue
		}
		onlineRecords = append(onlineRecords, def.OnlineRecord{GsId: v, OnlineNum: gsOnlineNum})
	}
	temp := map[int]int{}
	for _, o := range onlineRecords {
		temp[o.GsId] = o.OnlineNum
	}
	record, _ := json.Marshal(temp)
	onlineNumRec := def.OnlineNumRec{Record: string(record)}
	goManager := fusion.GetBaseGormDB("go_manager")
	goManager.Create(&onlineNumRec)
	//})
	logger.Info("End CountPlayerOnlineNum")
}

func CountWastageRate() {
	gsIdx := fusion.GetActiveServerIds()
	logger.Info("Start CountWastageRate")
	now := time.Now().Unix()
	levelRecords := []def.LevelWastageRecord{}
	questRecords := []def.QuestWastageRecord{}
	for _, v := range gsIdx {
		lRecords := []def.LevelRecord{}
		db_char := fusion.GetServerGormDB("db_char", int64(v))
		if db_char == nil {
			logger.Infof("CountWastageRate() get db_char failed server id = %d", v)
			continue
		}
		logTimeMax := now - int64(3*DaySecond)
		logTimeMin := now - int64(4*DaySecond)
		//temp := []map[string]interface{}{}
		db_char.Table("inst_player_char").Select("ipcInstID", "ipcLevel").
			Where("ipcLastLoginTime >= ?", logTimeMin).
			Where("ipcLastLoginTime <= ?", logTimeMax).Scan(&lRecords)
		temp := map[int]int{}
		for _, l := range lRecords {
			temp[l.Level]++
		}
		record, _ := json.Marshal(temp)
		levelRecord := def.LevelWastageRecord{Record: string(record), GsId: v}
		levelRecords = append(levelRecords, levelRecord)

		userIds := []int{}
		for _, l := range lRecords {
			userIds = append(userIds, l.PlayerID)
		}
		qRecords := []def.QuestRecord{}

		db_log := fusion.GetServerGormDB("db_log", int64(v))
		if db_log == nil {
			logger.Infof("CountWastageRate() get db_log failed server id = %d", v)
			continue
		}
		db_log.Raw("select a.questTypeID,any_value(a.playerId)AS playerId\n"+
			"from\n"+
			"(select questTypeID,playerId,whenType from log_player_quests where (whenType = 0 )and QuestClass = 1 AND playerId IN ?) as a\n"+
			"inner join \n"+
			"(select questTypeID,playerId,whenType from log_player_quests where (whenType = 4 )and QuestClass = 1) as b\n"+
			"on a.playerId = b.playerId and a.questTypeID = b.questTypeID and a.whenType != b.whenType\n"+
			"group by a.questTypeID", userIds).Scan(&qRecords)

		temp = map[int]int{}
		for _, l := range qRecords {
			temp[l.QuestTypeId]++
		}

		record, _ = json.Marshal(temp)
		questRecord := def.QuestWastageRecord{Record: string(record), GsId: v}
		questRecords = append(questRecords, questRecord)
	}
	goManager := fusion.GetBaseGormDB("go_manager")
	goManager.Create(&levelRecords)
	goManager.Create(&questRecords)
	logger.Info("End CountWastageRate")
}

func CountPlayerData() {
	gsIdx := fusion.GetActiveServerIds()
	logger.Info("Start CountPlayerData")
	countPlayerDatas := []def.CountPlayerData{}
	for _, v := range gsIdx {
		cut := []def.PlayerInfoCut{}
		db_char := fusion.GetServerGormDB("db_char", int64(v))
		if db_char == nil {
			logger.Infof("CountPlayerData() get db_char failed server id = %d", v)
			continue
		}
		db_char.Table("inst_player_char").
			Select("ipcInstID", "ipcLevel", "ipcCamp", "ipcS64Values").
			Scan(&cut)
		levelCount := map[int]int{}
		campCount := map[int]int{}
		for _, data := range cut {
			levelCount[data.Level]++
			campCount[data.Camp]++
		}
		record1, _ := json.Marshal(levelCount)
		record2, _ := json.Marshal(campCount)

		start, end := fusion.GetAimDayStartEnd(1)
		tempTime, _ := time.ParseInLocation("2006-01-02 15:04:05", start, time.Local)
		activeCut := make([]def.ActiveCut, 0)
		db_char.Table("inst_player_char").
			Select("ipcAcctID,ipcS64Values").
			Where("ipcLastLoginTime >= ?", tempTime.Unix()).
			Scan(&activeCut)
		var onlineCount int64
		var activeCount int64

		var AcctId4OnlineCount = make(map[int]int)

		for _, data := range activeCut {
			temp := fusion.String2Bytes(fusion.Strval(data.S64Values))
			temp1 := make([]int, 0)
			json.Unmarshal(temp, &temp1)
			if len(temp1) >= 4 {
				if temp1[3] != 0 {
					onlineCount += int64(temp1[3])
				}
				if AcctId4OnlineCount[data.AcctId] == 0 {
					AcctId4OnlineCount[data.AcctId] = temp1[3]
				} else {
					AcctId4OnlineCount[data.AcctId] += temp1[3]
				}
			}
		}
		for _, data := range AcctId4OnlineCount {
			if data > 60*60 {
				activeCount++
			}
		}

		var count int64
		db_char.Table("guild_member").Count(&count)
		playerIds := []int{}
		db_log := fusion.GetServerGormDB("db_log", int64(v))
		if db_log == nil {
			logger.Infof("CountPlayerData() get db_log failed server id = %d", v)
			continue
		}
		db_log.Table("log_player_quests").
			Select("playerId").
			Where("logTime >= ?", start).
			Where("logTime <= ?", end).
			Where("whenType =?", def.Family).
			Scan(&playerIds)
		familyQuestCount := map[int]int{}
		for _, id := range playerIds {
			familyQuestCount[id]++
		}
		record3, _ := json.Marshal(familyQuestCount)
		levelUpsCut := make([]def.LevelUpsCut, 0)
		db_log.Raw("select acctId,playerId,MAX(playerLevel)as playerLevel,min(logTime) as logTime " +
			"from  " +
			"log_player_levelups GROUP BY playerId,acctId " +
			"order by playerLevel desc, logTime ASC limit 200").
			Scan(&levelUpsCut)

		record4, _ := json.Marshal(levelUpsCut)

		timeStr := tempTime.Format("200601")
		totalPriceCut := make([]def.TotalPriceCut, 0)
		db_global := fusion.GetBaseGormDB("db_global")
		db_global.Table("t_web_orders_"+timeStr).
			Select("playerId,sum(buyPrice)as totalPrice  ").
			Where("createTime >= ?", start).
			Where("createTime <= ?", end).
			Group("playerId").Order("totalPrice DESC").
			Scan(&totalPriceCut)
		record5, _ := json.Marshal(totalPriceCut)

		totalPriceCut = make([]def.TotalPriceCut, 0)

		db_global.Table("t_web_orders_"+timeStr).
			Select("playerId,sum(buyPrice)as totalPrice  ").
			Where("createTime >= ?", start).
			Where("createTime <= ?", end).
			Group("playerId").
			Order("totalPrice DESC").
			Scan(&totalPriceCut)
		record6, _ := json.Marshal(totalPriceCut)

		var totalAcctNum int64
		db_global.Table("t_accounts").Count(&totalAcctNum)

		countPlayerDatas = append(countPlayerDatas,
			def.CountPlayerData{
				LevelCount:           string(record1),
				FamilyCount:          string(record2),
				GsId:                 v,
				OnlineDurationCount:  onlineCount,
				FamilyQuestCount:     string(record3),
				GuildCount:           count,
				OnlinePlayerCount:    activeCount,
				LevelRankDaily:       string(record4),
				RechargeRankDaily:    string(record5),
				WebRechargeRankDaily: string(record6),
				TotalAcctNum:         totalAcctNum,
			})
	}
	goManager := fusion.GetBaseGormDB("go_manager")
	goManager.Create(&countPlayerDatas)
	logger.Info("End CountPlayerData")
}

func CountShopBuyItem() {
	gsIdx := fusion.GetActiveServerIds()
	logger.Info("Start CountShopBuyItem")
	buyShopItemRecords := []def.BuyShopItemRecord{}
	now := time.Now().Unix()
	now = now - int64(DaySecond)
	tm := time.Unix(now, 0)
	date := tm.Format("20060102")
	for _, v := range gsIdx {
		itemCut := []def.ItemCut{}
		db_log := fusion.GetServerGormDB("db_log", int64(v))
		if db_log == nil {
			logger.Infof("CountShopBuyItem() get db_log failed server id = %d", v)
			continue
		}
		db_log.Table("log_player_items_"+date).
			Where("flowType", flowDef.LFT_SHOP).Scan(&itemCut)
		temp := map[int]int{}
		for _, l := range itemCut {
			temp[l.ItemTypeId] += l.ItemNum
		}
		record, _ := json.Marshal(temp)
		buyShopItemRecord := def.BuyShopItemRecord{Record: string(record), GsId: v}
		buyShopItemRecords = append(buyShopItemRecords, buyShopItemRecord)
	}
	goManager := fusion.GetBaseGormDB("go_manager")
	goManager.Create(&buyShopItemRecords)
	logger.Info("End CountShopBuyItem")
}

func CountOrder() {
	gsIdx := fusion.GetActiveServerIds()
	logger.Info("Start CountOrder")
	orderRecords := []def.OrderRecord{}
	now := time.Now().Unix()

	now = now - int64(DaySecond)
	tm := time.Unix(now, 0)
	date := tm.Format("200601")
	start, end := fusion.GetAimDayStartEnd(1)
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
	for _, v := range gsIdx {
		orderCut := []def.OrderCut{}
		db_global := fusion.GetBaseGormDB("db_global")
		db_global.Table("t_orders_"+date).
			Where("gsId = ?", v).
			Where("payTime >= ?", start).
			Where("payTime <= ?", end).
			Scan(&orderCut)
		var priceCount int
		var priceCountWeb int
		tempIOS := map[int]int{}
		tempAD := map[int]int{}
		tempWeb := map[int]int{}
		playerIds := make([]int, 0)
		for _, data := range orderCut {
			if PlatForm4ID[data.PayId] == 1 {
				tempAD[data.PayId]++
				priceCount += data.BuyPrice
			}
			if PlatForm4ID[data.PayId] == 2 {
				tempIOS[data.PayId]++
				priceCount += data.BuyPrice
			}
			if PlatForm4ID[data.PayId] == 3 {
				tempWeb[data.PayId]++
				priceCountWeb += data.BuyPrice
			}
			playerIds = append(playerIds, data.PlayerID)

		}
		recordIOS, _ := json.Marshal(tempIOS)
		recordAD, _ := json.Marshal(tempAD)
		orderCut = []def.OrderCut{}
		db_global.Table("t_web_orders_"+date).
			Where("createTime >= ?", start).
			Where("createTime <= ?", end).
			Where("gsId = ?", v).
			Scan(&orderCut)
		temp := map[int]int{}
		for _, data := range orderCut {
			if PlatForm4ID[data.PayId] == 1 {
				tempAD[data.PayId]++
				priceCount += data.BuyPrice
			}
			if PlatForm4ID[data.PayId] == 2 {
				tempIOS[data.PayId]++
				priceCount += data.BuyPrice
			}
			if PlatForm4ID[data.PayId] == 3 {
				tempWeb[data.PayId]++
				priceCountWeb += data.BuyPrice
			}
			playerIds = append(playerIds, data.PlayerID)
		}
		var payAccountId int64
		db_global.Table("t_account_characters").
			Select(" COUNT(DISTINCT accountId)").
			Where("groupId = ? and characterId in ?", v, playerIds).
			Scan(&payAccountId)

		var TotalRechargeAcct int64

		db_global.Table("t_account_characters").
			Select(" COUNT(DISTINCT accountId)").
			Where("groupId = ? and characterId in (?)", v,
				db_global.Table("t_web_orders_"+date).Select("playerId").Where("gsId = ?", v)).
			Scan(&TotalRechargeAcct)

		webRecord, _ := json.Marshal(temp)
		orderRecord := def.OrderRecord{
			IOSRecord:         string(recordIOS),
			AndroidRecord:     string(recordAD),
			WebRecord:         string(webRecord),
			GsId:              v,
			TotalPayNum:       int(payAccountId),
			TotalPrice:        priceCount,
			TotalPriceWeb:     priceCountWeb,
			TotalRechargeAcct: int(TotalRechargeAcct),
		}
		orderRecords = append(orderRecords, orderRecord)
	}
	goManager := fusion.GetBaseGormDB("go_manager")
	goManager.Create(&orderRecords)
	logger.Info("End CountOrder")
}

func CountPlayDungeon() {
	gsIdx := fusion.GetActiveServerIds()
	logger.Info("Start CountPlayDungeon")
	dungeonsRecords := []def.DungeonsRecord{}
	for _, v := range gsIdx {
		dRecords := []def.DungeonsCut{}
		db_log := fusion.GetServerGormDB("db_log", int64(v))
		if db_log == nil {
			logger.Infof("CountPlayDungeon() get db_log failed server id = %d", v)
			continue
		}
		start, end := fusion.GetAimDayStartEnd(1)
		db_log.Table("log_player_dungeons").
			Select("playerID", "dungeonId", "progress").
			Where("logTime >= ?", start).
			Where("logTime <= ?", end).
			Scan(&dRecords)
		temp := map[int]map[string]int{}

		for _, l := range dRecords {
			if temp[l.DungeonId] == nil {
				temp[l.DungeonId] = map[string]int{}
			}
			if l.Progress == 0 {
				temp[l.DungeonId]["join"]++
			} else if l.Progress == 1 {
				temp[l.DungeonId]["clear"]++
			}
		}
		record, _ := json.Marshal(temp)

		dungeonsRecord := def.DungeonsRecord{
			Record: string(record),
			GsId:   v,
		}
		dungeonsRecords = append(dungeonsRecords, dungeonsRecord)
	}
	goManager := fusion.GetBaseGormDB("go_manager")
	goManager.Create(&dungeonsRecords)
	logger.Info("End CountPlayDungeon")
}

func GetPlayerRetainedData() {
	gsIdx := fusion.GetActiveServerIds()
	logger.Info("Start GetPlayerRetainedData")
	for _, v := range gsIdx {
		gormDB := fusion.GetServerGormDB("db_char", int64(v))
		if gormDB == nil {
			logger.Infof("GetPlayerRetainedData() get db_char failed server id = %d", v)
			continue
		}
		playerRetainedRec := []def.PlayerRetainedAccount{}
		go_manager := fusion.GetBaseGormDB("go_manager")

		now := time.Now()
		str := now.Format("20060102")
		go_manager.Model(def.PlayerRetainedAccount{}).
			Where("DATEDIFF(? , dayNo) <= 8", str).
			Where("DATEDIFF(? , dayNo) >= 0", str).
			Where("gsId = ?", v).
			Scan(&playerRetainedRec)

		playerRetained := []def.PlayerRetainedAccount{}
		gormDB.Raw("SELECT\n"+
			"    logMin.day1 AS dayNo,\n"+
			"	 ? AS gsId,\n"+
			"    COUNT(DISTINCT logMin.acctId) AS new,\n"+
			"    COUNT(\n"+
			"        DISTINCT CASE\n"+
			"            WHEN DATEDIFF(day2, day1) = 1 THEN logMin.acctId\n"+
			"        END\n"+
			"    ) AS retained2,\n"+
			"    COUNT(\n"+
			"        DISTINCT CASE\n"+
			"            WHEN DATEDIFF(day2, day1) = 2 THEN logMin.acctId\n"+
			"        END\n"+
			"    ) AS retained3,\n"+
			"    COUNT(\n"+
			"        DISTINCT CASE\n"+
			"            WHEN DATEDIFF(day2, day1) = 6 THEN logMin.acctId\n"+
			"        END\n"+
			"    ) AS retained7,\n"+
			"    CONCAT(\n"+
			"        COUNT(\n"+
			"            DISTINCT CASE\n"+
			"                WHEN DATEDIFF(day2, day1) = 1 THEN logMin.acctId\n"+
			"            END\n"+
			"        ) / COUNT(DISTINCT logMin.acctId) * 100\n"+
			"    ) AS retainedPer2,\n"+
			"    CONCAT(\n"+
			"        COUNT(\n"+
			"            DISTINCT CASE\n"+
			"                WHEN DATEDIFF(day2, day1) = 2 THEN logMin.acctId\n"+
			"            END\n"+
			"        ) / COUNT(DISTINCT logMin.acctId) * 100\n"+
			"    ) AS retainedPer3,\n"+
			"    CONCAT(\n"+
			"        COUNT(\n"+
			"            DISTINCT CASE\n"+
			"                WHEN DATEDIFF(day2, day1) = 6 THEN logMin.acctId\n"+
			"            END\n"+
			"        ) / COUNT(DISTINCT logMin.acctId) * 100\n"+
			"    ) AS retainedPer7,\n"+
			"	NOW() AS createTime,\n"+
			"	NOW() AS updateTime\n"+
			"FROM\n"+
			"    (\n"+
			"        SELECT\n"+
			"            ipcAcctID AS acctId,\n"+
			"            MAX(FROM_UNIXTIME(ipcLastLoginTime, '%Y%m%d')) AS day2\n"+
			"        FROM\n"+
			"            inst_player_char\n"+
			"        WHERE\n"+
			"            ipcAcctID IN (\n"+
			"                SELECT\n"+
			"                    acctId\n"+
			"                FROM\n"+
			"                    (\n"+
			"                        SELECT\n"+
			"                            ipcAcctID AS acctId,\n"+
			"                            MIN(FROM_UNIXTIME(ipcCreateTime, '%Y%m%d')) AS creatTime\n"+
			"                        FROM\n"+
			"                            inst_player_char\n"+
			"                        WHERE\n"+
			"                            ipcServerID = ?\n"+
			"                        GROUP BY\n"+
			"                            ipcAcctID\n"+
			"                    ) AS a\n"+
			"            )\n"+
			"            AND ipcServerID = ?\n"+
			"        GROUP BY\n"+
			"            ipcAcctID\n"+
			"    ) AS logMax\n"+
			"    INNER JOIN (\n"+
			"        SELECT\n"+
			"            ipcAcctID AS acctId,\n"+
			"            MIN(FROM_UNIXTIME(ipcCreateTime, '%Y%m%d')) AS day1\n"+
			"        FROM\n"+
			"            inst_player_char\n"+
			"        WHERE\n"+
			"            ipcServerID = ?\n"+
			"        GROUP BY\n"+
			"            ipcAcctID\n"+
			"    ) AS logMin ON logMax.acctId = logMin.acctId WHERE DATEDIFF(now(),logMin.day1) <=8\n"+
			"GROUP BY\n"+
			"    logMin.day1", v, v, v, v).Scan(&playerRetained)
		newRec := []def.PlayerRetainedAccount{}

		for _, temp := range playerRetained {
			flag := false
			for key, v := range playerRetainedRec {
				vDayNo, _ := time.ParseInLocation("20060102", v.DayNo, time.Local)
				tempDayNo, _ := time.ParseInLocation("20060102", temp.DayNo, time.Local)
				diff := fusion.DiffNatureDays(vDayNo.Unix(), tempDayNo.Unix())
				if diff == -1 || diff == 0 {
					//if temp.DayNo == v.DayNo {
					flag = true
					tempDay, err := time.ParseInLocation("20060102", v.DayNo, time.Local)
					if err != nil {
						return
					}
					if err != nil {
						logger.Infof("GetPlayerRetainedData() convert Failed err:%s", err)
						continue
					}
					interVal := fusion.DiffNatureDays(now.Unix(), tempDay.Unix())
					//interVal := math.Floor(math.Abs(float64(now.Sub(tempDay))) / 60 * 60 * 24 * 1000)
					if interVal == 2 {
						playerRetainedRec[key].Retained2 = temp.Retained2
						playerRetainedRec[key].RetainedPer2 = temp.RetainedPer2
					}
					if interVal == 3 {
						playerRetainedRec[key].Retained3 = temp.Retained3
						playerRetainedRec[key].RetainedPer3 = temp.RetainedPer3
					}
					if interVal == 7 {
						playerRetainedRec[key].Retained7 = temp.Retained7
						playerRetainedRec[key].RetainedPer7 = temp.RetainedPer7
					}
					playerRetainedRec[key].UpdateTime = temp.UpdateTime
					//}
				}
			}
			if !flag {
				tempDayNo, _ := time.ParseInLocation("20060102", temp.DayNo, time.Local)
				diff := fusion.DiffNatureDays(now.Unix(), tempDayNo.Unix())
				if diff == -1 || diff == 0 {
					continue
				}
				newRec = append(newRec, temp)
			}
		}

		if len(playerRetainedRec) != 0 {
			go_manager.Save(&playerRetainedRec)
		}

		if len(newRec) != 0 {
			go_manager.Create(&newRec)
		}

		//playerRetainedRecAccount := []def.PlayerRetainedAccount{}
		//go_manager.Model(def.PlayerRetainedAccount{}).
		//	Where("DATEDIFF(? , dayNo) <= 7", str).
		//	Where("DATEDIFF(? , dayNo) >= 0", str).
		//	Where("gsId = ?", v).
		//	Scan(&playerRetainedRecAccount)
		//
		//playerRetainedAccount := []def.PlayerRetainedAccount{}
		//
		//db_global := fusion.GetBaseGormDB("db_global")
		//db_global.Raw("SELECT *, ? AS gsId,\n"+
		//	"100*Retained2/new AS  RetainedPer2,\n"+
		//	"100*Retained3/new AS RetainedPer3,\n"+
		//	"100*Retained7/new AS RetainedPer7,\n"+
		//	"NOW() AS createTime,\n"+
		//	"NOW() AS updateTime\n"+
		//	"FROM \n"+
		//	"(\n"+
		//	"        SELECT\n"+
		//	"        FROM_UNIXTIME(c.createTime,'%Y%m%d') dayNo,\n"+
		//	"        COUNT(DISTINCT c.username)  new,\n"+
		//	"        COUNT(DISTINCT d.username)  Retained2,\n"+
		//	"        COUNT(DISTINCT e.username)  Retained3,\n"+
		//	"        COUNT(DISTINCT f.username)  Retained7\n"+
		//	"        FROM\n"+
		//	"        (\n"+
		//	"                SELECT username,lastLoginTime,createTime FROM t_accounts\n"+
		//	"                WHERE DATEDIFF(NOW(),FROM_UNIXTIME(createTime,'%Y-%m-%d')) >0\n"+
		//	"                AND  DATEDIFF(NOW(),FROM_UNIXTIME(createTime,'%Y-%m-%d')) <=7 AND lastLoginGsId = ?\n"+
		//	") c\n"+
		//	"        LEFT JOIN t_accounts d ON c.username = d.username  AND  DATEDIFF(FROM_UNIXTIME(d.lastLoginTime,'%Y-%m-%d'), \n"+
		//	"FROM_UNIXTIME(c.createTime,'%Y-%m-%d')) = 1\n"+
		//	"        LEFT JOIN t_accounts e ON c.username = e.username  AND  DATEDIFF(FROM_UNIXTIME(e.lastLoginTime,'%Y-%m-%d'),\n"+
		//	"FROM_UNIXTIME(c.createTime,'%Y-%m-%d')) = 3\n"+
		//	"        LEFT JOIN t_accounts f ON c.username = f.username  AND  DATEDIFF(FROM_UNIXTIME(f.lastLoginTime,'%Y-%m-%d'),\n"+
		//	"FROM_UNIXTIME(c.createTime,'%Y-%m-%d')) = 7\n"+
		//	"        GROUP BY dayNo\n"+
		//	") p;", v, v).Scan(&playerRetainedAccount)
		//newRecAccount := []def.PlayerRetainedAccount{}
		////nowday, _ = strconv.Atoi(str)
		//
		//for _, temp := range playerRetainedAccount {
		//	flag := false
		//	for key, v := range playerRetainedRecAccount {
		//		vDayNo, _ := time.ParseInLocation("20060102", v.DayNo, time.Local)
		//		tempDayNo, _ := time.ParseInLocation("20060102", temp.DayNo, time.Local)
		//		diff := fusion.DiffNatureDays(vDayNo.Unix(), tempDayNo.Unix())
		//		if diff == -1 || diff == 0 {
		//			flag = true
		//			tempDay, err := time.ParseInLocation("20060102", v.DayNo, time.Local)
		//			if err != nil {
		//				logger.Infof("GetPlayerRetainedData() convert Failed err:%s", err)
		//				continue
		//			}
		//			interVal := fusion.DiffNatureDays(now.Unix(), tempDay.Unix())
		//			//interVal := math.Floor(math.Abs(float64(now.Sub(tempDay))) / 60 * 60 * 24 * 1000)
		//			if interVal == 1 {
		//				playerRetainedRecAccount[key].Retained2 = temp.Retained2
		//				playerRetainedRecAccount[key].RetainedPer2 = temp.RetainedPer2
		//			}
		//			if interVal == 3 {
		//				playerRetainedRecAccount[key].Retained3 = temp.Retained3
		//				playerRetainedRecAccount[key].RetainedPer3 = temp.RetainedPer3
		//			}
		//			if interVal == 7 {
		//				playerRetainedRecAccount[key].Retained7 = temp.Retained7
		//				playerRetainedRecAccount[key].RetainedPer7 = temp.RetainedPer7
		//			}
		//			playerRetainedRecAccount[key].UpdateTime = temp.UpdateTime
		//		}
		//	}
		//	if !flag {
		//		newRecAccount = append(newRecAccount, temp)
		//	}
		//}
		//if len(playerRetainedRecAccount) != 0 {
		//	go_manager.Save(&playerRetainedRecAccount)
		//}
		//if len(newRecAccount) != 0 {
		//	go_manager.Create(&newRecAccount)
		//}
		//
		//playerRetainedRecIP := []def.PlayerRetainedIP{}
		//go_manager.Model(def.PlayerRetainedIP{}).
		//	Where("DATEDIFF(? , dayNo) <= 7", str).
		//	Where("DATEDIFF(? , dayNo) >= 0", str).
		//	Where("gsId = ?", v).
		//	Scan(&playerRetainedRecIP)
		//
		//playerRetainedIP := []def.PlayerRetainedIP{}
		//db_global.Raw("SELECT *, ? AS gsId,\n"+
		//	"100*Retained2/new AS  RetainedPer2,\n"+
		//	"100*Retained3/new AS RetainedPer3,\n"+
		//	"100*Retained7/new AS RetainedPer7,\n"+
		//	"NOW() AS createTime,\n"+
		//	"NOW() AS updateTime\n"+
		//	"FROM \n"+
		//	"(\n"+
		//	"        SELECT\n"+
		//	"        FROM_UNIXTIME(c.createTime,'%Y%m%d') dayNo,\n"+
		//	"        COUNT(DISTINCT c.lastLoginIP)  new,\n"+
		//	"        COUNT(DISTINCT d.lastLoginIP)  Retained2,\n"+
		//	"        COUNT(DISTINCT e.lastLoginIP)  Retained3,\n"+
		//	"        COUNT(DISTINCT f.lastLoginIP)  Retained7\n"+
		//	"        FROM\n"+
		//	"        (\n"+
		//	"                SELECT lastLoginIP,lastLoginTime,createTime FROM t_accounts\n"+
		//	"                WHERE DATEDIFF(NOW(),FROM_UNIXTIME(createTime,'%Y-%m-%d')) >0\n"+
		//	"                AND  DATEDIFF(NOW(),FROM_UNIXTIME(createTime,'%Y-%m-%d')) <=7 AND lastLoginGsId = ?\n"+
		//	") c\n"+
		//	"        LEFT JOIN t_accounts d ON c.lastLoginIP = d.lastLoginIP  AND  DATEDIFF(FROM_UNIXTIME(d.lastLoginTime,'%Y-%m-%d'), \n"+
		//	"FROM_UNIXTIME(c.createTime,'%Y-%m-%d')) = 1\n"+
		//	"        LEFT JOIN t_accounts e ON c.lastLoginIP = e.lastLoginIP  AND  DATEDIFF(FROM_UNIXTIME(e.lastLoginTime,'%Y-%m-%d'),\n"+
		//	"FROM_UNIXTIME(c.createTime,'%Y-%m-%d')) = 3\n"+
		//	"        LEFT JOIN t_accounts f ON c.lastLoginIP = f.lastLoginIP  AND  DATEDIFF(FROM_UNIXTIME(f.lastLoginTime,'%Y-%m-%d'),\n"+
		//	"FROM_UNIXTIME(c.createTime,'%Y-%m-%d')) = 7\n"+
		//	"        GROUP BY dayNo\n"+
		//	") p;", v, v).Scan(&playerRetainedIP)
		//newRecIP := []def.PlayerRetainedIP{}
		////nowday, _ = strconv.Atoi(str)
		//
		//for _, temp := range playerRetainedIP {
		//	flag := false
		//	for key, v := range playerRetainedRecIP {
		//		vDayNo, _ := time.ParseInLocation("20060102", v.DayNo, time.Local)
		//		tempDayNo, _ := time.ParseInLocation("20060102", temp.DayNo, time.Local)
		//		diff := fusion.DiffNatureDays(vDayNo.Unix(), tempDayNo.Unix())
		//		if diff == -1 || diff == 0 {
		//			flag = true
		//			tempDay, err := time.ParseInLocation("20060102", v.DayNo, time.Local)
		//			if err != nil {
		//				logger.Infof("GetPlayerRetainedData() convert Failed err:%s", err)
		//				continue
		//			}
		//			interVal := fusion.DiffNatureDays(now.Unix(), tempDay.Unix())
		//			//interVal := math.Floor(math.Abs(float64(now.Sub(tempDay))) / 60 * 60 * 24 * 1000)
		//			if interVal == 1 {
		//				playerRetainedRecIP[key].Retained2 = temp.Retained2
		//				playerRetainedRecIP[key].RetainedPer2 = temp.RetainedPer2
		//			}
		//			if interVal == 3 {
		//				playerRetainedRecIP[key].Retained3 = temp.Retained3
		//				playerRetainedRecIP[key].RetainedPer3 = temp.RetainedPer3
		//			}
		//			if interVal == 7 {
		//				playerRetainedRecIP[key].Retained7 = temp.Retained7
		//				playerRetainedRecIP[key].RetainedPer7 = temp.RetainedPer7
		//			}
		//			playerRetainedRecIP[key].UpdateTime = temp.UpdateTime
		//		}
		//	}
		//	if !flag {
		//		newRecIP = append(newRecIP, temp)
		//	}
		//}
		//if len(playerRetainedRecIP) != 0 {
		//	go_manager.Save(&playerRetainedRecIP)
		//}
		//if len(newRecIP) != 0 {
		//	go_manager.Create(&newRecIP)
		//}
	}
	logger.Info("End GetPlayerRetainedData")
}
func AddUniqueUint32(slice *[]uint32, newIds []uint32) {
	if len(newIds) == 0 {
		return
	}
	existing := make(map[uint32]bool)
	for _, id := range *slice {
		existing[id] = true
	}
	for _, id := range newIds {
		if !existing[id] {
			*slice = append(*slice, id)
			existing[id] = true
		}
	}
}

func CountPlayerActiveData() {
	gsIdx := fusion.GetActiveServerIds()
	totalAcctIDs := make([]uint32, 0)
	logger.Info("Start CountPlayerActiveData")
	for i := len(gsIdx) - 1; i >= 0; i-- {
		v := gsIdx[i]
		initLen := len(totalAcctIDs)
		gormDB := fusion.GetServerGormDB("db_char", int64(v))
		if gormDB == nil {
			logger.Infof("CountPlayerActiveData() get db_char failed server id = %d", v)
			continue
		}
		ipcAcctIDs := make([]uint32, 0)
		gormDB.Raw("select \n" +
			"ipcAcctID\n" +
			"from\n" +
			"inst_player_char\n" +
			"where\n" +
			"FROM_UNIXTIME(ipcLastLoginTime,'%Y-%m-%d') \n" +
			">= DATE_FORMAT((DATE_SUB(NOW(),INTERVAL 1 DAY)),'%Y-%m-%d') GROUP BY ipcAcctID").
			Scan(&ipcAcctIDs)
		AddUniqueUint32(&totalAcctIDs, ipcAcctIDs)
		activeA := len(totalAcctIDs) - initLen
		ipcAcctIDs = make([]uint32, 0)

		gormDB.Raw("select \n" +
			"ipcAcctID\n" +
			"from\n" +
			"inst_player_char\n" +
			"where\n" +
			"FROM_UNIXTIME(ipcCreateTime,'%Y-%m-%d') \n" +
			">= DATE_FORMAT((DATE_SUB(NOW(),INTERVAL 1 DAY)),'%Y-%m-%d') GROUP BY ipcAcctID").
			Scan(&ipcAcctIDs)
		newP := len(ipcAcctIDs)
		logCount := 0
		db_log := fusion.GetServerGormDB("db_log", int64(v))
		db_log.Raw("SELECT \n" +
			"COUNT(id) \n" +
			"FROM \n" +
			"log_character_onlines \n" +
			"WHERE \n" +
			"DATE_FORMAT((DATE_SUB(NOW(),INTERVAL 1 DAY)),'%Y-%m-%d') \n" +
			">= DATE_FORMAT(logTime,'%Y-%m-%d')").Scan(&logCount)

		activeCount := def.PlayerActiveCount{
			ActiveA:       activeA,
			ActiveB:       0,
			NewActivePer:  float64(newP),
			LoginManTimes: logCount,
			GsId:          v,
		}
		go_manager := fusion.GetBaseGormDB("go_manager")
		go_manager.Create(&activeCount)
	}
	logger.Info("End CountPlayerActiveData")
}

func CountSingleServer() {
	gsIdx := fusion.GetActiveServerIds()
	logger.Info("Start countSingleServer")
	for _, v := range gsIdx {
		db_log := fusion.GetServerGormDB("db_log", int64(v))
		if db_log == nil {
			logger.Infof("countSingleServer() get db_char failed server id = %d", v)
			continue
		}
		type numCountCut struct {
			Max     int `gorm:"column:max" db:"max"`
			Min     int `gorm:"column:min" db:"min"`
			Average int `gorm:"column:average" db:"average"`
		}
		numCount := make([]numCountCut, 0)
		db_log.Raw("SELECT MAX(num) as max ,MIN(num) as min, ROUND(SUM(num)/24) as average FROM\n" +
			"(\n" +
			"SELECT COUNT(a.characterId)AS num,DATE_FORMAT(logTime,'%y-%m-%d %h:%i:00')AS logM\n" +
			"FROM \n" +
			"(SELECT \n" +
			"characterId,logTime\n" +
			"FROM \n" +
			"log_character_onlines\n" +
			"WHERE\n" +
			"DATE_FORMAT(NOW(),'%y%m%d') = DATE_FORMAT(logTime,'%y%m%d')\n" +
			") AS a GROUP BY logM) AS b").Scan(&numCount)
		type accountCut struct {
			NewC int `gorm:"column:newC" db:"newC"`
			AllC int `gorm:"column:allC" db:"allC"`
		}
		account := make([]accountCut, 0)
		db_global := fusion.GetBaseGormDB("db_global")
		db_global.Raw("SELECT \n" +
			"COUNT(DISTINCT CASE WHEN DATE_FORMAT(NOW(),'%y%m%d') = FROM_UNIXTIME(a.createTime,'%y%m%d') THEN a.Id END)AS newC,\n" +
			"COUNT(Id)AS allC FROM\n" +
			"(SELECT Id, FROM_UNIXTIME(createTime,'%Y%-m-%d') AS createTime FROM t_accounts)AS a").Scan(&account)

		date := time.Now().Format("200601")
		type orderCut struct {
			CurDayR         int `gorm:"column:curDayR" db:"curDayR"`
			CountR          int `gorm:"column:countR" db:"countR"`
			CurDayManTimesR int `gorm:"column:curDayManTimesR" db:"curDayManTimesR"`
		}
		order := make([]orderCut, 0)
		sql2 := "select\n" +
			"COUNT(DISTINCT CASE WHEN date_format(Now(),'%Y-%m-%d') = payDay THEN a.playerId END)AS curDayR,\n" +
			"count(DISTINCT a.playerId) countR,\n" +
			"count(if (date_format(Now(),'%Y-%m-%d') = payDay,true,NULL)) as curDayManTimesR\n" +
			"from \n" +
			"(select playerId,buyPrice, DATE_FORMAT(payTime,'%Y-%m-%d')as payDay from t_orders_" + date + " )as a"
		db_global.Raw(sql2).Scan(&order)

		type orderCut1 struct {
			PayTime  sql.NullTime `gorm:"column:payTime" db:"payTime"`
			PlayerId int          `gorm:"column:playerId" db:"playerId"`
		}
		order1 := make([]orderCut1, 0)
		sql1, _ := fusion.CreateUnionTableSqlByDays("", "", "t_orders_", db_global,
			orderCut1{}, "", true)
		sql1 = "SELECT\n" +
			"playerId,\n" +
			"MIN(a.payTime)AS payTime\n" +
			"FROM\n" +
			"(" + sql1 + ")AS a\n" +
			"WHERE DATE_FORMAT(NOW(),'%Y-%m-%d') = DATE_FORMAT(a.payTime,'%Y-%m-%d') \n" +
			"GROUP BY playerId "
		db_global.Raw(sql1).Scan(&order1)

		dataCount := def.SingleServerDataCount{
			GsId:            v,
			MaxOnline:       numCount[0].Max,
			MinOnline:       numCount[0].Min,
			AverageOnline:   numCount[0].Average,
			NewAccountNum:   account[0].NewC,
			AllAccountNum:   account[0].AllC,
			NewRNum:         len(order1),
			CurDayRNum:      order[0].CurDayR,
			CountRNum:       order[0].CurDayR,
			CurDayManTimesR: order[0].CurDayManTimesR,
			LogTime:         time.Now().Format("2006-01-02 15:04:05"),
		}
		go_manager := fusion.GetBaseGormDB("go_manager")
		go_manager.Create(&dataCount)
	}
	logger.Info("End countSingleServer")
}

func countCurDaySuspiciousData() {
	gsIdx := fusion.GetActiveServerIds()
	logger.Info("Start countCurDaySuspiciousData")
	go_manager := fusion.GetBaseGormDB("go_manager")

	type2kv := make(map[int]map[int]int)
	type2keys := make(map[int][]int)

	temp := make([]map[string]interface{}, 0)

	go_manager.Table("global_cfg_check").
		Select("`type`,`key`,min(`value`) as `value`").
		Group("`type`,`key`").
		Scan(&temp)

	for _, v := range temp {
		gType, _ := strconv.Atoi(fusion.Strval(v["type"]))
		key, _ := strconv.Atoi(fusion.Strval(v["key"]))
		value, _ := strconv.Atoi(fusion.Strval(v["value"]))
		if type2kv[gType] == nil {
			type2kv[gType] = make(map[int]int)
		}
		type2kv[gType][key] = value
		type2keys[gType] = append(type2keys[gType], key)
	}
	startTime, _ := fusion.GetAimDayStartEndTime(1)
	start := startTime.Format("20060102")

	type tempMoneyCut struct {
		MoneyType  int `gorm:"column:moneyType"`
		MoneyValue int `gorm:"column:moneyValue"`
		PlayerId   int `gorm:"column:playerId"`
	}

	type tempItemCut struct {
		ItemTypeId int `gorm:"column:itemTypeId"`
		ItemNum    int `gorm:"column:itemNum"`
		LogType    int `gorm:"column:logType"`
		PlayerId   int `gorm:"column:playerId"`
	}

	type tempPlayerCut struct {
		PlayerId        int    `gorm:"column:ipcInstID"`
		IpcCurrencies   string `gorm:"column:ipcCurrencies"`
		IpcStorageItems string `gorm:"column:ipcStorageItems"`
	}

	type tempAccountCut struct {
		PlayerId   int    `gorm:"column:characterId"`
		AccountId  int    `gorm:"column:accountId"`
		ServerId   int    `gorm:"column:serverId"`
		PlayerName string `gorm:"column:characterName"`
	}

	vaildRecord := make([]def.SuspiciousDataRecord, 0)
	for _, v := range gsIdx {
		db_log := fusion.GetServerGormDB("db_log", int64(v))
		if db_log == nil {
			continue
		}
		records := make([]def.SuspiciousDataRecord, 0)
		cut1 := make([]tempMoneyCut, 0)
		cut2 := make([]tempItemCut, 0)
		cut3 := make([]tempPlayerCut, 0)
		cut4 := make([]tempAccountCut, 0)
		PlayerIds := make([]int, 0)

		kv := type2kv[def.Cfg_type_chequeDaily]
		name := "log_player_moneys_" + start
		db_log.Table(name).
			Select("moneyType ,SUM(moneyValue)AS moneyValue,playerId").
			Where("moneyValue >0").
			Where("moneyType in ?", type2keys[def.Cfg_type_chequeDaily]).
			Where("flowType not in ?",
				[]int{flowDef.LFT_ITEM_ARRANGE, flowDef.LFT_ITEM_SPLIT, flowDef.LFT_ITEM_SWAP, flowDef.LFT_ITEM_EQUIP, flowDef.LFT_ITEM_UNEQUIP}).
			Group("playerId,moneyType").
			Scan(&cut1)
		for _, tempCut := range cut1 {
			val, err := kv[tempCut.MoneyType]
			if !err {
				continue
			}
			if tempCut.MoneyValue >= val {
				records = append(records, def.SuspiciousDataRecord{
					GsId:     v,
					PlayerId: tempCut.PlayerId,
					Type:     def.Cfg_type_chequeDaily,
					Key:      tempCut.MoneyType,
					Value:    tempCut.MoneyValue,
					LogTime:  def.MyTime(time.Now()),
				})
				PlayerIds = append(PlayerIds, tempCut.PlayerId)
			}
		}

		kv = type2kv[def.Cfg_type_ItemDaily]
		name = "log_player_items_" + start
		db_log.Table(name).
			Select("itemTypeId ,SUM(itemNum)AS itemNum,playerId,logType").
			Where("logType IN ?", []int{0, 2}).
			Where("itemTypeId IN ?", type2keys[def.Cfg_type_ItemDaily]).
			Where("flowType not in ?",
				[]int{flowDef.LFT_ITEM_ARRANGE, flowDef.LFT_ITEM_SPLIT, flowDef.LFT_ITEM_SWAP, flowDef.LFT_ITEM_EQUIP, flowDef.LFT_ITEM_UNEQUIP}).
			Group("playerId,itemTypeId,logType").
			Scan(&cut2)
		for _, tempCut := range cut2 {
			val, err := kv[tempCut.ItemTypeId]
			if !err {
				continue
			}
			if tempCut.ItemNum >= val {
				records = append(records, def.SuspiciousDataRecord{
					GsId:     v,
					PlayerId: tempCut.PlayerId,
					Type:     def.Cfg_type_ItemDaily,
					Key:      tempCut.ItemTypeId,
					Value:    tempCut.ItemNum,
					LogTime:  def.MyTime(time.Now()),
				})
				PlayerIds = append(PlayerIds, tempCut.PlayerId)
			}
		}

		startChar, _ := fusion.GetAimDayStartEndTime(3)
		db_char := fusion.GetServerGormDB("db_char", int64(v))
		if db_char == nil {
			continue
		}
		db_char.Table("inst_player_char").
			Select("ipcInstID,ipcCurrencies,ipcStorageItems").
			Where("ipcLastOnlineTime >= ?", startChar.Unix()).
			Scan(&cut3)

		kv = type2kv[def.Cfg_type_Cheque]
		for _, player := range cut3 {
			tempCurrency := fusion.String2Bytes(fusion.Strval(player.IpcCurrencies))
			Currencies := []int{}
			json.Unmarshal(tempCurrency, &Currencies)
			for id, num := range Currencies {
				val, err := kv[id+1]
				if !err {
					continue
				}
				if num >= val {
					records = append(records, def.SuspiciousDataRecord{
						GsId:     v,
						PlayerId: player.PlayerId,
						Type:     def.Cfg_type_Cheque,
						Key:      id + 1,
						Value:    num,
						LogTime:  def.MyTime(time.Now()),
					})
					PlayerIds = append(PlayerIds, player.PlayerId)
				}
			}
			kv = type2kv[def.Cfg_type_Item]
			arr := strings.Split(player.IpcStorageItems, ";")
			var list = fusion.ReadItemListFromByte(arr, true)
			for _, item := range list {
				val, err := kv[int(item.ItemTypeID)]
				if !err {
					continue
				}
				if int(item.ItemCount) >= val {
					records = append(records, def.SuspiciousDataRecord{
						GsId:     v,
						PlayerId: player.PlayerId,
						Type:     def.Cfg_type_Item,
						Key:      int(item.ItemTypeID),
						Value:    int(item.ItemCount),
						LogTime:  def.MyTime(time.Now()),
					})
					PlayerIds = append(PlayerIds, player.PlayerId)
				}
			}
		}
		db_global := fusion.GetBaseGormDB("db_global")
		db_global.Table("t_account_characters").
			Select("accountId,characterId,serverId,characterName").
			Where("serverId = ?", v).
			Where("characterId IN ?", PlayerIds).
			Scan(&cut4)
		for _, accountCut := range cut4 {
			for _, record := range records {
				if accountCut.PlayerId == record.PlayerId {
					record.AccountId = accountCut.AccountId
					record.PlayerName = accountCut.PlayerName
					vaildRecord = append(vaildRecord, record)
				}
			}
		}
	}
	go_manager.Create(&vaildRecord)
	if len(vaildRecord) != 0 {
		SendEMail.SendEmail(gotext.Get("有") + strconv.Itoa(len(vaildRecord)) + gotext.Get("条监测信息待处理!"))
	}
	logger.Info("End countCurDaySuspiciousData")
}

func CountPlayerOnline4Min() {
	gsIdx := fusion.GetActiveServerIds()
	logger.Info("Start CountPlayerOnlineNum")
	onlineRecords := []def.OnlineRecordMin{}
	for _, v := range gsIdx {
		ctx1 := &context.Context{}
		vals := url.Values{}
		vals.Add("gsId", strconv.Itoa(v))
		err, res := fusion.CallToCenter(&vals, common.MyCfg.Api["getGSOnlineNumber"], "GET", ctx1)
		if err != nil {
			logger.Info("CountPlayerOnlineNum Failed, Because CenterServer has some error")
			continue
		}
		var gsOnlineNum int
		gsOnlineNum, err = strconv.Atoi(res)
		if err != nil {
			logger.Infof("CountPlayerOnlineNum Failed,serverID = %d,With %s", v, res)
			continue
		}
		onlineRecords = append(onlineRecords, def.OnlineRecordMin{GsId: v, OnlineNum: gsOnlineNum})
	}
	temp := map[int]int{}
	for _, o := range onlineRecords {
		temp[o.GsId] = o.OnlineNum
	}
	record, _ := json.Marshal(temp)
	onlineNumRec := def.OnlineNumRecMin{Record: string(record)}
	goManager := fusion.GetBaseGormDB("go_manager")
	goManager.Create(&onlineNumRec)
	//})
	date := time.Now().Add(-(time.Hour * 6))
	goManager.Where("logTime <=?", def.MyTime(date)).Delete(def.OnlineNumRecMin{})
	logger.Info("End CountPlayerOnlineNum")
}

func CountAccountData() {
	gsIdx := fusion.GetActiveServerIds()
	logger.Info("Start CountAccountData")
	AccountCountData := make([]def.AccountCountData, 0)
	for _, v := range gsIdx {
		db_char := fusion.GetServerGormDB("db_char", int64(v))
		if db_char == nil {
			continue
		}
		var offline7Day int
		db_char.Table("inst_player_char").
			Where("DATEDIFF(NOW(),FROM_UNIXTIME(ipcLastLoginTime, '%Y%m%d')) >=7 and ipcServerID = ?", v).
			Select("count(DISTINCT ipcAcctID) as offline7Day").
			Scan(&offline7Day)

		var online48Hour int
		db_char.Table("inst_player_char").
			Where("ipcS64Values->'$[6]' > 172800 and ipcServerID = ?", v).
			Select("count(DISTINCT ipcAcctID) as online48Hour").
			Scan(&online48Hour)

		var online48HourAccount = make([]int, 0)
		db_char.Table("inst_player_char").
			Where("ipcS64Values->'$[6]' > 172800 and ipcServerID = ?", v).
			Select("ipcAcctID").
			Scan(&online48HourAccount)

		var paidPlayerIds = make([]int, 0)

		db_global := fusion.GetBaseGormDB("db_global")
		sql1 := fusion.GetSql4AccountCount("t_web_orders_", db_global, "where gsId = ? ")
		db_global.Raw(sql1, v).Scan(&paidPlayerIds)

		var online48HourPaid int
		for _, id := range paidPlayerIds {
			if fusion.Contain(online48HourAccount, id) {
				online48HourPaid++
			}
		}

		var paidAccountNum int
		db_global.Table("t_account_characters").
			Select("count(DISTINCT accountId) as paidAccountNum").
			Where("serverId = ? and characterId in ?", v, paidPlayerIds).
			Scan(&paidAccountNum)

		var allAccountNum int
		db_global.Table("t_account_characters").
			Select("count(DISTINCT accountId) as allAccountNum").
			Where("serverId = ?", v).
			Scan(&allAccountNum)

		PaidPer, _ := strconv.ParseFloat(strconv.FormatFloat(float64(paidAccountNum)/math.Max(float64(allAccountNum), 1), 'f', 2, 32), 32)
		AccountCountData = append(AccountCountData, def.AccountCountData{
			ServerId:         v,
			Offline7Day:      offline7Day,
			Online48Hour:     online48Hour,
			Online48HourPaid: online48HourPaid,
			PaidPer:          PaidPer,
			LogTime:          def.MyTime(time.Now()),
		})
	}
	db_manager := fusion.GetBaseGormDB("go_manager")
	db_manager.Create(&AccountCountData)
	logger.Info("end CountAccountData")
}

func CountPlayerOnlineNum4Career() {
	type HTTPRespBase struct {
		Err int    `json:"error"`
		Msg string `json:"message"`
	}

	type httpOnlineNum4CareerResp struct {
		HTTPRespBase
		OnlineNum4Career map[int]int `json:",omitempty"`
	}

	logger.Info("start CountPlayerOnlineNum4Career")
	gsIdx := fusion.GetActiveServerIds()
	onlineRecords := []def.OnlineRecord4Career{}
	for _, v := range gsIdx {
		ctx1 := &context.Context{}
		vals := url.Values{}
		vals.Add("gsId", strconv.Itoa(v))
		err, res := fusion.CallToCenter(&vals, common.MyCfg.Api["GetGSOnlineNumber4Career"], "GET", ctx1)
		if err != nil {
			logger.Info("CountPlayerOnlineNum Failed, Because CenterServer has some error")
			continue
		}
		temp := fusion.String2Bytes(res)
		resp := httpOnlineNum4CareerResp{}
		json.Unmarshal(temp, &resp)

		if err != nil {
			logger.Infof("CountPlayerOnlineNum Failed,serverID = %d,With %s", v, res)
			continue
		}
		onlineRecords = append(onlineRecords, def.OnlineRecord4Career{GsId: v, OnlineData: resp.OnlineNum4Career})
	}
	temp := map[int]map[int]int{}
	for _, o := range onlineRecords {
		temp[o.GsId] = o.OnlineData
	}
	record, _ := json.Marshal(temp)
	onlineNumRec := def.OnlineNumRec4Career{Record: string(record)}
	goManager := fusion.GetBaseGormDB("go_manager")
	goManager.Create(&onlineNumRec)
	logger.Info("end CountPlayerOnlineNum4Career")
}

func StartAllTimer() {
	recorder := NewTimeMgrTimer()
	recorder.CreateTriggerX(timer.ByHourlyTrigger,
		timer.TriggerPoint{Hour: 0, Min: 0, Sec: 0}, func() {
			CountPlayerOnlineNum()
		}, 0)
	recorder.CreateTriggerX(timer.ByDailyTrigger,
		timer.TriggerPoint{Hour: 0, Min: 0, Sec: 0}, func() {
			CountWastageRate()
		}, 0)
	recorder.CreateTriggerX(timer.ByDailyTrigger,
		timer.TriggerPoint{Hour: 0, Min: 0, Sec: 0}, func() {
			CountShopBuyItem()
		}, 0)
	recorder.CreateTriggerX(timer.ByDailyTrigger,
		timer.TriggerPoint{Hour: 0, Min: 0, Sec: 0}, func() {
			CountOrder()
		}, 0)
	recorder.CreateTriggerX(timer.ByDailyTrigger,
		timer.TriggerPoint{Hour: 0, Min: 0, Sec: 0}, func() {
			CountPlayDungeon()
		}, 0)
	recorder.CreateTriggerX(timer.ByDailyTrigger,
		timer.TriggerPoint{Hour: 0, Min: 0, Sec: 0}, func() {
			CountPlayerData()
		}, 0)
	recorder.CreateTriggerX(timer.ByDailyTrigger,
		timer.TriggerPoint{Hour: 23, Min: 59, Sec: 30}, func() {
			CountPlayerActiveData()
		}, 0)
	recorder.CreateTriggerX(timer.ByDailyTrigger,
		timer.TriggerPoint{Hour: 23, Min: 59, Sec: 30}, func() {
			CountSingleServer()
		}, 0)
	recorder.CreateTriggerX(timer.ByDailyTrigger,
		timer.TriggerPoint{Hour: 2, Min: 0, Sec: 0}, func() {
			countCurDaySuspiciousData()
		}, 0)
	recorder.CreateTriggerX(timer.ByDailyTrigger,
		timer.TriggerPoint{Hour: 0, Min: 0, Sec: 0}, func() {
			GetPlayerRetainedData()
		}, 0)
	//recorder.CreateTriggerX(timer.ByDailyTrigger,
	//	timer.TriggerPoint{Hour: 12, Min: 0, Sec: 0}, func() {
	//		GetPlayerRetainedData(gsIdx)
	//	}, 0)
	recorder.CreateTriggerX(timer.By5MinutelyTrigger,
		timer.TriggerPoint{Hour: 0, Min: 0, Sec: 0}, func() {
			CountPlayerOnline4Min()
		}, 0)
	recorder.CreateTriggerX(timer.ByDailyTrigger,
		timer.TriggerPoint{Hour: 0, Min: 0, Sec: 0}, func() {
			CountAccountData()
		}, 0)
	recorder.CreateTriggerX(timer.ByHourlyTrigger,
		timer.TriggerPoint{Hour: 0, Min: 0, Sec: 0}, func() {
			CountPlayerOnlineNum4Career()
		}, 0)
}
