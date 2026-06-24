package def

type NoticeData struct {
	Id         uint32 `gorm:"primaryKey;column:Id" db:"id"`
	StartTime  MyTime `gorm:"column:startTime" db:"startTime"`
	EndTime    MyTime `gorm:"column:endTime" db:"endTime"`
	Interval   uint32 `gorm:"column:interval" db:"interval"`
	Msg        string `gorm:"column:msg" db:"msg"`
	MsgFlags   uint32 `gorm:"column:msgFlags" db:"msgFlags"`
	ToChannels uint32 `gorm:"column:toChannels" db:"toChannels"`
	GsIds      string `gorm:"column:gsIds" db:"gsIds"`
	GsNotIds   string `gorm:"column:gsNotIds" db:"gsNotIds"`
	IsOnce     int    `gorm:"column:isOnce" db:"isOnce"`
}

func (NoticeData) TableName() string {
	return "t_notices"
}
