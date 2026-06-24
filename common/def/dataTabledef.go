package def

type PlayerActiveCount struct {
	Id            int     `gorm:"primaryKey;column:Id" db:"Id"`
	GsId          int     `gorm:"column:gsId" db:"gsId"`
	ActiveA       int     `gorm:"column:activeA" db:"activeA"`
	ActiveB       int     `gorm:"column:activeB" db:"activeB"`
	NewActivePer  float64 `gorm:"column:newActivePer" db:"newActivePer"`
	LoginManTimes int     `gorm:"column:loginManTimes" db:"loginManTimes"`
	LogTime       MyTime  `gorm:"autoCreateTime;column:logTime" db:"logTime"`
}

func (PlayerActiveCount) TableName() string {
	return "player_active_count"
}

type SingleServerDataCount struct {
	Id              int    `gorm:"primaryKey;column:Id" db:"Id"`
	GsId            int    `gorm:"column:gsId" db:"gsId"`
	MaxOnline       int    `gorm:"column:maxOnline" db:"maxOnline"`
	MinOnline       int    `gorm:"column:minOnline" db:"minOnline"`
	AverageOnline   int    `gorm:"column:averageOnline" db:"averageOnline"`
	NewAccountNum   int    `gorm:"column:newAccountNum" db:"newAccountNum"`
	AllAccountNum   int    `gorm:"column:allAccountNum" db:"allAccountNum"`
	NewRNum         int    `gorm:"column:newRNum" db:"newRNum"`
	CurDayRNum      int    `gorm:"column:curDayRNum" db:"curDayRNum"`
	CountRNum       int    `gorm:"column:countRNum" db:"countRNum"`
	CurDayManTimesR int    `gorm:"column:curDayManTimesR" db:"curDayManTimesR"`
	LogTime         string `gorm:"column:logTime" db:"logTime"`
}

func (SingleServerDataCount) TableName() string {
	return "single_server_count"
}
