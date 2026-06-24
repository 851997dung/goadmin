package def

type ActivationCodeDef struct {
	Code      string `gorm:"primaryKey;column:code" db:"code"`
	StartTime MyTime `gorm:"column:startTime" db:"startTime"`
	EndTime   MyTime `gorm:"column:endTime" db:"endTime"`
	//UseTime   sql.NullTime `gorm:"column:useTime" db:"useTime"`
	IsGlobal int    `gorm:"column:isGlobal" db:"isGlobal"`
	BatchNum int    `gorm:"column:batchNum" db:"batchNum"`
	GiftId   int    `gorm:"column:giftId" db:"giftId"`
	GsIds    string `gorm:"column:gsIds" db:"gsIds"`
	GsNotIds string `gorm:"column:gsNotIds" db:"gsNotIds"`
}

func (ActivationCodeDef) TableName() string {
	return "t_activation_code"
}
