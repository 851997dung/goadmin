package def

import (
	"database/sql/driver"
	"errors"
	"fmt"
	"github.com/leonelquinteros/gotext"
	"strings"
	"time"
)

const (
	None        = iota
	MainLine    // 主线
	BranchLine  // 支线
	Daily       // 日常
	Guide       // 指引
	Transfer    // 转职
	Active      // 活跃
	ActivityMap //活动副本
	Family      // 家族任务
)

const (
	Accept = iota
	Finish
	Failed
	Cancel
	Submit
)

type OnlineNumRec struct {
	Id      int    `gorm:"primaryKey;column:Id" db:"id"`
	LogTime MyTime `gorm:"column:logTime;autoCreateTime" db:"logTime"`
	Record  string `gorm:"column:record" db:"record"`
}
type OnlineRecord struct {
	GsId      int `json:"gsId" db:"gsId"`
	OnlineNum int `json:"onlineNum" db:"onlineNum"`
}

func (s *OnlineNumRec) TableName() string {
	return "player_online_num_record"
}

type OnlineNumRecMin struct {
	Id      int    `gorm:"primaryKey;column:Id" db:"id"`
	LogTime MyTime `gorm:"column:logTime;autoCreateTime" db:"logTime"`
	Record  string `gorm:"column:record" db:"record"`
}
type OnlineRecordMin struct {
	GsId      int `json:"gsId" db:"gsId"`
	OnlineNum int `json:"onlineNum" db:"onlineNum"`
}

func (s *OnlineNumRecMin) TableName() string {
	return "player_online_num_record_min"
}

type LevelWastageRecord struct {
	Id      int    `gorm:"primaryKey;column:Id" db:"id"`
	LogTime MyTime `gorm:"column:logTime;autoCreateTime" db:"logTime"`
	GsId    int    `gorm:"column:gsId" db:"gsId"`
	Record  string `gorm:"column:record" db:"record"`
}

func (s *LevelWastageRecord) TableName() string {
	return "level_wastage_record"
}

type QuestWastageRecord struct {
	Id      int    `gorm:"primaryKey;column:Id" db:"id"`
	LogTime MyTime `gorm:"column:logTime;autoCreateTime" db:"logTime"`
	GsId    int    `gorm:"column:gsId" db:"gsId"`
	Record  string `gorm:"column:record" db:"record"`
}

func (s *QuestWastageRecord) TableName() string {
	return "quest_wastage_record"
}

type LevelRecord struct {
	PlayerID int `gorm:"primaryKey;column:ipcInstID";json:"playerId" db:"playerID"`
	Level    int `gorm:"column:ipcLevel";json:"level" db:"level"`
}

type QuestRecord struct {
	PlayerId    int `gorm:"column:playerId";json:"playerId" db:"playerId"`
	QuestTypeId int `gorm:"column:questTypeID";json:"questTypeId" db:"questTypeId"`
}
type CountPlayerData struct {
	Id                   int    `gorm:"primaryKey;column:Id" db:"id"`
	GsId                 int    `gorm:"column:gsId" db:"gsId"`
	LogTime              MyTime `gorm:"column:logTime;autoCreateTime" db:"logTime"`
	LevelCount           string `gorm:"column:levelCount" db:"levelCount"`
	FamilyCount          string `gorm:"column:familyCount" db:"familyCount"`
	FamilyQuestCount     string `gorm:"column:familyQuestCount" db:"familyQuestCount"`
	GuildCount           int64  `gorm:"column:guildCount" db:"guildCount"`
	OnlineDurationCount  int64  `gorm:"column:onlineDurationCount" db:"onlineDurationCount"`
	OnlinePlayerCount    int64  `gorm:"column:onlinePlayerCount" db:"onlinePlayerCount"`
	LevelRankDaily       string `gorm:"column:levelRankDaily" db:"levelRankDaily"`
	RechargeRankDaily    string `gorm:"column:rechargeRankDaily" db:"rechargeRankDaily"`
	WebRechargeRankDaily string `gorm:"column:webRechargeRankDaily" db:"webRechargeRankDaily"`
	TotalAcctNum         int64  `gorm:"column:totalAcctNum" db:"totalAcctNum"`
}

func (s *CountPlayerData) TableName() string {
	return "player_data_count"
}

type PlayerInfoCut struct {
	PlayerId  int    `gorm:"primaryKey;column:ipcInstID" db:"playerId"`
	Level     int    `gorm:"column:ipcLevel" db:"level"`
	Camp      int    `gorm:"column:ipcCamp" db:"camp"`
	S64Values string `gorm:"column:ipcS64Values" db:"s64Values"`
}

type ActiveCut struct {
	AcctId    int    `gorm:"column:ipcAcctID" db:"acctId"`
	S64Values string `gorm:"column:ipcS64Values" db:"s64Values"`
}

type BuyShopItemRecord struct {
	Id      int    `gorm:"primaryKey;column:Id" db:"id"`
	LogTime MyTime `gorm:"column:logTime;autoCreateTime" db:"logTime"`
	GsId    int    `gorm:"column:gsId" db:"gsId"`
	Record  string `gorm:"column:record" db:"record"`
}

type ItemCut struct {
	ItemTypeId int `gorm:"column:itemTypeId" db:"itemTypeId"`
	ItemNum    int `gorm:"column:itemNum" db:"itemNum"`
}

type LevelUpsCut struct {
	AcctId      int `gorm:"column:acctId" db:"acctId"`
	PlayerId    int `gorm:"column:playerId" db:"playerId"`
	PlayerLevel int `gorm:"column:playerLevel" db:"playerLevel"`
}

type TotalPriceCut struct {
	PlayerId   int `gorm:"column:playerId" db:"playerId"`
	TotalPrice int `gorm:"column:totalPrice" db:"totalPrice"`
}

func (s *BuyShopItemRecord) TableName() string {
	return "buy_shop_item_record"
}

type OrderRecord struct {
	Id                int    `gorm:"primaryKey;column:Id" db:"id"`
	LogTime           MyTime `gorm:"column:logTime;autoCreateTime" db:"logTime"`
	GsId              int    `gorm:"column:gsId" db:"gsId"`
	TotalPayNum       int    `gorm:"column:totalPayNum" db:"totalPayNum"`
	AndroidRecord     string `gorm:"column:AndroidRecord" db:"AndroidRecord"`
	IOSRecord         string `gorm:"column:IOSRecord" db:"IOSRecord"`
	WebRecord         string `gorm:"column:webRecord" db:"webRecord"`
	TotalPrice        int    `gorm:"column:totalPrice" db:"totalPrice"`
	TotalPriceWeb     int    `gorm:"column:totalPriceWeb" db:"totalPriceWeb"`
	TotalRechargeAcct int    `gorm:"column:totalRechargeAcct" db:"totalRechargeAcct"`
}

type OrderCut struct {
	PlayerID int `gorm:"column:playerId" db:"playerId"`
	PayId    int `gorm:"column:payId" db:"payId"`
	BuyPrice int `gorm:"column:buyPrice" db:"buyPrice"`
}

func (s *OrderRecord) TableName() string {
	return "order_record"
}

type DungeonsRecord struct {
	Id      int    `gorm:"primaryKey;column:Id" db:"id"`
	LogTime MyTime `gorm:"column:logTime;autoCreateTime" db:"logTime"`
	GsId    int    `gorm:"column:gsId" db:"gsId"`
	Record  string `gorm:"column:record" db:"record"`
}

type DungeonsCut struct {
	PlayerID  int `gorm:"column:playerID" db:"playerID"`
	DungeonId int `gorm:"column:dungeonId" db:"dungeonId"`
	Progress  int `gorm:"column:progress" db:"progress"`
}

func (s *DungeonsRecord) TableName() string {
	return "dungeons_record"
}

//type PlayerRetainedCharacters struct {
//	Id           int     `gorm:"primaryKey;column:id"`
//	DayNo        string  `gorm:"column:dayNo"`
//	ServerId     int     `gorm:"column:gsId"`
//	New          int64   `gorm:"column:new"`
//	Retained2    int64   `gorm:"column:retained2"`
//	Retained3    int64   `gorm:"column:retained3"`
//	Retained7    int64   `gorm:"column:retained7"`
//	RetainedPer2 float64 `gorm:"column:retainedPer2"`
//	RetainedPer3 float64 `gorm:"column:retainedPer3"`
//	RetainedPer7 float64 `gorm:"column:retainedPer7"`
//	CreateTime   MyTime  `gorm:"column:createTime"`
//	UpdateTime   MyTime  `gorm:"column:updateTime"`
//}
//
//func (s *PlayerRetainedCharacters) TableName() string {
//	return "player_retained_characters"
//}

type PlayerRetainedAccount struct {
	Id           int     `gorm:"primaryKey;column:id"`
	DayNo        string  `gorm:"column:dayNo"`
	ServerId     int     `gorm:"column:gsId"`
	New          int64   `gorm:"column:new"`
	Retained2    int64   `gorm:"column:retained2"`
	Retained3    int64   `gorm:"column:retained3"`
	Retained7    int64   `gorm:"column:retained7"`
	RetainedPer2 float64 `gorm:"column:retainedPer2"`
	RetainedPer3 float64 `gorm:"column:retainedPer3"`
	RetainedPer7 float64 `gorm:"column:retainedPer7"`
	CreateTime   MyTime  `gorm:"column:createTime"`
	UpdateTime   MyTime  `gorm:"column:updateTime"`
}

func (s *PlayerRetainedAccount) TableName() string {
	return "player_retained_account"
}

type AccountCountData struct {
	Id               int     `gorm:"primaryKey;column:id"`
	ServerId         int     `gorm:"column:serverId"`
	Offline7Day      int     `gorm:"column:offline7Day"`
	Online48Hour     int     `gorm:"column:online48Hour"`
	Online48HourPaid int     `gorm:"column:online48HourPaid"`
	PaidPer          float64 `gorm:"column:paidPer"`
	LogTime          MyTime  `gorm:"column:LogTime"`
}

func (s *AccountCountData) TableName() string {
	return "account_data_count"
}

//type PlayerRetainedIP struct {
//	Id           int     `gorm:"primaryKey;column:id"`
//	DayNo        string  `gorm:"column:dayNo"`
//	ServerId     int     `gorm:"column:gsId"`
//	New          int64   `gorm:"column:new"`
//	Retained2    int64   `gorm:"column:retained2"`
//	Retained3    int64   `gorm:"column:retained3"`
//	Retained7    int64   `gorm:"column:retained7"`
//	RetainedPer2 float64 `gorm:"column:retainedPer2"`
//	RetainedPer3 float64 `gorm:"column:retainedPer3"`
//	RetainedPer7 float64 `gorm:"column:retainedPer7"`
//	CreateTime   MyTime  `gorm:"column:createTime"`
//	UpdateTime   MyTime  `gorm:"column:updateTime"`
//}
//
//func (s *PlayerRetainedIP) TableName() string {
//	return "player_retained_ip"
//}

//MyTime 自定义时间
type MyTime time.Time

func (t *MyTime) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		return nil
	}
	var err error
	//前端接收的时间字符串
	str := string(data)
	//去除接收的str收尾多余的"
	timeStr := strings.Trim(str, "\"")
	t1, err := time.Parse("2006-01-02 15:04:05", timeStr)
	*t = MyTime(t1)
	return err
}

func (t MyTime) MarshalJSON() ([]byte, error) {
	formatted := fmt.Sprintf("\"%v\"", time.Time(t).Format("2006-01-02 15:04:05"))
	return []byte(formatted), nil
}

func (t MyTime) Value() (driver.Value, error) {
	// MyTime 转换成 time.Time 类型
	tTime := time.Time(t)
	return tTime.Format("2006-01-02 15:04:05"), nil
}

func (t *MyTime) Scan(v interface{}) error {
	switch vt := v.(type) {
	case time.Time:
		// 字符串转成 time.Time 类型
		*t = MyTime(vt)
	default:
		return errors.New(gotext.Get("类型处理错误"))
	}
	return nil
}

func (t *MyTime) String() string {
	return fmt.Sprintf("hhh:%s", time.Time(*t).String())
}

func (t *MyTime) Hour() string {
	return time.Time(*t).Format("15")
}

type OnlineRecord4Career struct {
	GsId       int         `json:"gsId" db:"gsId"`
	OnlineData map[int]int `json:"onlineData" db:"onlineData"`
}

type OnlineNumRec4Career struct {
	Id      int    `gorm:"primaryKey;column:Id" db:"id"`
	LogTime MyTime `gorm:"column:logTime;autoCreateTime" db:"logTime"`
	Record  string `gorm:"column:record" db:"record"`
}

func (s *OnlineNumRec4Career) TableName() string {
	return "player_online_num_record_4_Career"
}
