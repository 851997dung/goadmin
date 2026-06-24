package pageMgr

import (
	"admin/common/def"
	"admin/fusion"
	"encoding/json"
	"github.com/leonelquinteros/gotext"
	"sort"
	"strconv"
	"strings"
)

func GetTotalRegisterNum() int64 {
	db := fusion.GetBaseGormDB("db_global")
	var count int64
	db.Table("t_accounts").Count(&count)
	return count
}

func GetTotalPayPriceNum(gsidx int) (countPrice int64, webCountPrice int64, count int64, webCount int64) {
	db := fusion.GetBaseGormDB("db_global")
	tableNames := []string{}
	db.Raw("SELECT table_name FROM information_schema.TABLES WHERE TABLE_SCHEMA='mmorpg_global' ").Scan(&tableNames)
	for _, v := range tableNames {
		var tempCount, tempCountWeb, count1, count2 int64
		index := strings.Index(v, "t_orders_")
		if index != -1 {
			if gsidx == 0 {
				db.Table(v).Select("sum(buyPrice)").Where("payTime != NULL").Scan(&tempCount)
				db.Table(v).Distinct("playerId").Where("payTime != NULL").Count(&count1)
			} else {
				db.Table(v).Select("sum(buyPrice)").Where("gsId=?", gsidx).Where("payTime != NULL").Scan(&tempCount)
				db.Table(v).Select("Count(DISTINCT playerId)").Where("gsId=?", gsidx).Where("payTime != NULL").Scan(&count1)
			}
			countPrice += tempCount
			count += count1
		}
		index = strings.Index(v, "t_web_orders_")
		if index != -1 {
			{
				if gsidx == 0 {
					db.Table(v).Select("sum(buyPrice)").Scan(&tempCountWeb)
					db.Table(v).Distinct("playerId").Count(&count2)
				} else {
					db.Table(v).Select("sum(buyPrice)").Where("gsId=?", gsidx).Scan(&tempCountWeb)
					db.Table(v).Where("gsId=?", gsidx).Distinct("playerId").Count(&count2)
				}
			}
			webCountPrice += tempCountWeb
			webCount += count2
		}
	}
	return countPrice, webCountPrice, count, webCount
}
func GetOnlineNumData(date string, gsIdx int) (int, int, []float64) {
	if date == "" {
		return 0, 0, nil
	}
	var max, min int
	//timeStart, timeEnd := fusion.GetAimDayStartEnd(1)
	go_manager := fusion.GetBaseGormDB("go_manager")
	onlineNumRec := []def.OnlineNumRec{}
	go_manager.Table("player_online_num_record").
		Where("logTime LIKE ?", "%"+date+"%").
		//Where("logTime<=?", timeEnd).
		Order("logTime ASC").
		Scan(&onlineNumRec)
	data := map[int]int{}
	var data1 = make([]float64, 24)
	for _, v := range onlineNumRec {
		b := fusion.String2Bytes(v.Record)
		err := json.Unmarshal(b, &data)
		if err != nil {
			continue
		}
		hour, _ := strconv.Atoi(v.LogTime.Hour())
		data1[hour] = float64(data[gsIdx])
		if max <= data[gsIdx] {
			max = data[gsIdx]
		}
		if min >= data[gsIdx] {
			min = data[gsIdx]
		}
	}
	return min, max, data1
}
func GetQuestWastageNumData(date string, gsIdx int) (int, []string, []float64) {
	if date == "" {
		return 0, nil, nil
	}
	var total int
	var questWastageData []float64
	var fieldsQuest []string
	go_manager := fusion.GetBaseGormDB("go_manager")
	questWastageRecord := def.QuestWastageRecord{}
	go_manager.Table(questWastageRecord.TableName()).
		Where("logTime LIKE ?", "%"+date+"%").
		Where("gsId = ?", gsIdx).
		Scan(&questWastageRecord)
	data := map[int]int{}
	b := fusion.String2Bytes(questWastageRecord.Record)
	err := json.Unmarshal(b, &data)
	if err != nil {
		return 0, nil, nil
	}
	fieldsQuest = []string{}
	questWastageData = []float64{}

	keys := make([]int, 0)
	for k, v := range data {
		keys = append(keys, k)
		total += v
	}
	sort.Ints(keys)
	for _, key := range keys {
		temp := strconv.FormatFloat(float64(data[key])/float64(total)*100, 'f', 2, 32)
		field := strconv.Itoa(key) + "  (" + temp + "%)"
		fieldsQuest = append(fieldsQuest, field)
		questWastageData = append(questWastageData, float64(data[key]))
	}
	return total, fieldsQuest, questWastageData
}

func GetLevelWastageNumData(date string, gsIdx int) (int, []string, []float64) {
	if date == "" {
		return 0, nil, nil
	}

	var total int
	var levelWastageData []float64
	var fieldsLevel []string
	//timeStart, timeEnd := fusion.GetAimDayStartEnd(1)
	go_manager := fusion.GetBaseGormDB("go_manager")
	onlineNumRec := def.LevelWastageRecord{}
	go_manager.Table(onlineNumRec.TableName()).
		Where("logTime LIKE ?", "%"+date+"%").
		Where("gsId = ?", gsIdx).
		Scan(&onlineNumRec)
	data := map[int]int{}
	b := fusion.String2Bytes(onlineNumRec.Record)
	err := json.Unmarshal(b, &data)
	if err != nil {
		return 0, nil, nil
	}
	fieldsLevel = []string{}
	levelWastageData = []float64{}
	keys := make([]int, 0)
	for k := range data {
		keys = append(keys, k)
	}
	sort.Ints(keys)
	for _, key := range keys {
		levelWastageData = append(levelWastageData, float64(data[key]))
		total += data[key]
	}
	for _, key := range keys {
		temp := strconv.FormatFloat(float64(data[key])/float64(total)*100, 'f', 2, 32)
		field := strconv.Itoa(key) + "  (" + temp + "%)"
		fieldsLevel = append(fieldsLevel, field)
	}
	return total, fieldsLevel, levelWastageData
}

func GetPlayerDataCount(date string, gsIdx int) (float32, int64, int64, []string, []string, []float64, []float64) {
	var joinQuestPer float32
	if date == "" {
		timeStart, timeEnd := fusion.GetAimDayStartEnd(1)
		go_manager := fusion.GetBaseGormDB("go_manager")
		countPlayerData := def.CountPlayerData{}
		go_manager.Table(countPlayerData.TableName()).
			Where("logTime>=?", timeStart).
			Where("logTime<=?", timeEnd).
			Where("gsId = ?", gsIdx).
			Scan(&countPlayerData)
		total := 0
		data := map[int]int{}
		b := fusion.String2Bytes(countPlayerData.FamilyQuestCount)
		err := json.Unmarshal(b, &data)
		if err != nil {
			return 0, 0, 0, nil, nil, nil, nil
		}
		for _, v := range data {
			total += v
		}
		joinQuestPer = 0
		if total != 0 {
			joinQuestPer = float32(total) / float32(len(data)*10) * float32(100)
		}

		return joinQuestPer, countPlayerData.GuildCount, countPlayerData.OnlineDurationCount, nil, nil, nil, nil
	}

	var playerLevelData, familyNumData []float64
	var fieldsPlayer, fieldsFamily []string

	go_manager := fusion.GetBaseGormDB("go_manager")
	countPlayerData := def.CountPlayerData{}
	go_manager.Table(countPlayerData.TableName()).
		Where("logTime LIKE ?", "%"+date+"%").
		Where("gsId = ?", gsIdx).
		Scan(&countPlayerData)

	data := map[int]int{}
	b := fusion.String2Bytes(countPlayerData.LevelCount)
	err := json.Unmarshal(b, &data)
	if err != nil {
		return 0, 0, 0, nil, nil, nil, nil
	}
	fieldsPlayer = []string{}
	playerLevelData = []float64{}
	keys := make([]int, 0)
	for k := range data {
		keys = append(keys, k)
	}
	sort.Ints(keys)
	total := 0
	for _, key := range keys {
		playerLevelData = append(playerLevelData, float64(data[key]))
		total += data[key]
	}
	for _, key := range keys {
		temp := strconv.FormatFloat(float64(data[key])/float64(total)*100, 'f', 2, 32)
		field := strconv.Itoa(key) + "  (" + temp + "%)"
		fieldsPlayer = append(fieldsPlayer, field)
	}

	data = map[int]int{}
	b = fusion.String2Bytes(countPlayerData.FamilyCount)
	err = json.Unmarshal(b, &data)
	if err != nil {
		return 0, 0, 0, nil, nil, nil, nil
	}
	fieldsFamily = []string{}
	familyNumData = []float64{}
	keys = make([]int, 0)
	for k := range data {
		keys = append(keys, k)
	}
	sort.Ints(keys)
	total = 0
	for _, key := range keys {
		familyNumData = append(familyNumData, float64(data[key]))
		total += data[key]
	}
	FamilyStrings := fusion.GetFamilyStrings()
	for _, key := range keys {
		temp := strconv.FormatFloat(float64(float64(data[key])/float64(total)*100), 'f', 2, 32)
		field := FamilyStrings[key] + "  (" + temp + "%)"
		fieldsFamily = append(fieldsFamily, field)
	}
	total = 0

	data = map[int]int{}
	b = fusion.String2Bytes(countPlayerData.FamilyQuestCount)
	err = json.Unmarshal(b, &data)
	if err != nil {
		return 0, 0, 0, nil, nil, nil, nil
	}
	for _, v := range data {
		total += v
	}
	joinQuestPer = 0
	if total != 0 {
		joinQuestPer = float32(total) / float32(len(data)*10) * float32(100)
	}

	return joinQuestPer, countPlayerData.GuildCount, countPlayerData.OnlineDurationCount,
		fieldsPlayer, fieldsFamily, playerLevelData, familyNumData
}
func GetShopItemBuyCount(date string, gsIdx int) ([]string, []float64) {
	if date == "" {
		return nil, nil
	}

	var fieldsItemBuy []string
	var ItemBuyData []float64
	//timeStart, timeEnd := fusion.GetAimDayStartEnd(1)
	go_manager := fusion.GetBaseGormDB("go_manager")
	buyShopItemRecord := def.BuyShopItemRecord{}
	go_manager.Table(buyShopItemRecord.TableName()).
		Where("logTime LIKE ?", "%"+date+"%").
		//Where("logTime<=?", timeEnd).
		Where("gsId = ?", gsIdx).
		Scan(&buyShopItemRecord)

	data := map[int]int{}
	b := fusion.String2Bytes(buyShopItemRecord.Record)
	err := json.Unmarshal(b, &data)
	if err != nil {
		return nil, nil
	}
	fieldsItemBuy = []string{}
	ItemBuyData = []float64{}
	keys := make([]int, 0)
	for k := range data {
		keys = append(keys, k)
	}
	sort.Ints(keys)
	total := 0
	for _, key := range keys {
		ItemBuyData = append(ItemBuyData, float64(data[key]))
		total += data[key]
	}
	for _, key := range keys {
		temp := strconv.FormatFloat(float64(data[key])/float64(total)*100, 'f', 2, 32)
		field := strconv.Itoa(key) + "  (" + temp + "%)"
		fieldsItemBuy = append(fieldsItemBuy, field)
	}
	return fieldsItemBuy, ItemBuyData
}
func GetDungeonsCount(date string, gsIdx int) (int64, int64, []string, []float64, []float64) {
	var joinNum, clearNum int64
	var fieldsDungeonsCunt []string
	var joinData, clearData []float64

	timeStart, timeEnd := fusion.GetAimDayStartEnd(1)
	go_manager := fusion.GetBaseGormDB("go_manager")
	dungeonsRecord := def.DungeonsRecord{}
	if date == "" {
		go_manager.Table(dungeonsRecord.TableName()).
			Where("logTime>=?", timeStart).
			Where("logTime<=?", timeEnd).
			Where("gsId = ?", gsIdx).
			Scan(&dungeonsRecord)

		data := map[int]map[string]int{}
		b := fusion.String2Bytes(dungeonsRecord.Record)
		err := json.Unmarshal(b, &data)
		if err != nil {
			return 0, 0, nil, nil, nil
		}
		fieldsDungeonsCunt = []string{}
		joinData = []float64{}
		clearData = []float64{}
		keys := make([]int, 0)
		for k := range data {
			keys = append(keys, k)
		}
		sort.Ints(keys)

		for _, key := range keys {
			joinData = append(joinData, float64(data[key]["join"]))
			clearData = append(clearData, float64(data[key]["clear"]))
			joinNum += int64(data[key]["join"])
			clearNum += int64(data[key]["clear"])
		}
		return joinNum, clearNum, nil, nil, nil
	}
	go_manager.Table(dungeonsRecord.TableName()).
		Where("logTime LIKE ?", "%"+date+"%").
		//Where("logTime<=?", timeEnd).
		Where("gsId = ?", gsIdx).
		Scan(&dungeonsRecord)

	data := map[int]map[string]int{}
	b := fusion.String2Bytes(dungeonsRecord.Record)
	err := json.Unmarshal(b, &data)
	if err != nil {
		return 0, 0, nil, nil, nil
	}
	fieldsDungeonsCunt = []string{}
	joinData = []float64{}
	clearData = []float64{}
	keys := make([]int, 0)
	for k := range data {
		keys = append(keys, k)
	}
	sort.Ints(keys)

	for _, key := range keys {
		joinData = append(joinData, float64(data[key]["join"]))
		clearData = append(clearData, float64(data[key]["clear"]))
		joinNum += int64(data[key]["join"])
		clearNum += int64(data[key]["clear"])
	}
	for _, key := range keys {
		temp := strconv.FormatFloat(float64(data[key]["join"])/float64(joinNum)*100, 'f', 2, 32)
		field := strconv.Itoa(key) + "  (" + temp + "%)"
		fieldsDungeonsCunt = append(fieldsDungeonsCunt, field)
	}

	return joinNum, clearNum, fieldsDungeonsCunt, joinData, clearData
}

func GetOrderCount(date string, platFormID, gsIdx int) (int64, int64, []string, []float64) {

	go_manager := fusion.GetBaseGormDB("go_manager")
	orderRecord := def.OrderRecord{}
	go_manager.Table(orderRecord.TableName()).
		Where("logTime like ?", "%"+date+"%").
		//Where("logTime<=?", timeEnd).
		Where("gsId = ?", gsIdx).
		Scan(&orderRecord)
	if date == "" {
		return int64(orderRecord.TotalPayNum), int64(orderRecord.TotalPrice), nil, nil
	}
	data := map[int]int{}
	var str = ""
	switch platFormID {
	case 1:
		str = orderRecord.AndroidRecord
		break
	case 2:
		str = orderRecord.WebRecord
		break
	case 3:
		str = orderRecord.IOSRecord
		break
	default:
		str = orderRecord.AndroidRecord
		break
	}
	b := fusion.String2Bytes(str)
	err := json.Unmarshal(b, &data)
	if err != nil {
		return 0, 0, nil, nil
	}
	var purchase = make([]map[string]interface{}, 0)
	db_world := fusion.GetBaseGormDB("db_world")
	db_world.Table("auto_purchase").
		Select("ID,buyPriceShow").
		Where("platformID = ?", platFormID).
		Scan(&purchase)
	keys := make([]int, 0)
	for k := range purchase {
		keys = append(keys, k)
	}
	var gear = make([]string, len(purchase))
	var num = make([]float64, len(purchase))
	sort.Ints(keys)
	total := 0
	for _, key := range keys {
		if value, ok := data[key]; ok {
			total += value
		}
	}
	if total == 0 {
		total = 1
	}
	for _, key := range keys {
		payID, _ := strconv.Atoi(fusion.Strval(purchase[key]["ID"]))
		payName := fusion.Strval(purchase[key]["buyPriceShow"])
		value1 := 0
		if value, ok := data[payID]; ok {
			value1 = value
			num[key] = float64(value)
		}
		temp := strconv.FormatFloat(float64(value1)/float64(total)*100, 'f', 2, 32)
		field := payName + "  (" + temp + "%)"
		gear[key] = field
	}
	return int64(orderRecord.TotalRechargeAcct), int64(orderRecord.TotalPrice*1000 + orderRecord.TotalPriceWeb), gear, num
}

func GetDataDetailString() (data map[string]string) {
	data = make(map[string]string)
	data["NewPlayer"] = gotext.Get("新增玩家")
	data["ActivePlayer"] = gotext.Get("活跃玩家")
	data["RechargePrice"] = gotext.Get("充值金额(元)")
	data["RechargeNum"] = gotext.Get("充值人数")

	data["ARPPU"] = gotext.Get("ARPPU")
	data["ARPU"] = gotext.Get("ARPU")
	data["Permeability"] = gotext.Get("充值渗透率(%)")

	data["Retained2"] = gotext.Get("次日留存")
	data["Retained3"] = gotext.Get("三日留存")
	data["Retained7"] = gotext.Get("七日留存")

	data["RetainedPer2"] = gotext.Get("次日留存(%)")
	data["RetainedPer3"] = gotext.Get("三日留存(%)")
	data["RetainedPer7"] = gotext.Get("七日留存(%)")
	return data
}
