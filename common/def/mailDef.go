package def

type InstMail struct {
	MailID          uint32 `gorm:"primaryKey;column:mailID" db:"mailID"`
	MailType        uint32 `gorm:"column:mailType" db:"mailType"`
	MailFlags       uint32 `gorm:"column:mailFlags" db:"mailFlags"`
	MailSender      uint32 `gorm:"column:mailSender" db:"mailSender"`
	MailReceiver    uint32 `gorm:"column:mailReceiver" db:"mailReceiver"`
	MailDeliverTime int    `gorm:"column:mailDeliverTime" db:"mailDeliverTime"`
	MailExpireTime  int    `gorm:"column:mailExpireTime" db:"mailExpireTime"`
	MailSubject     string `gorm:"column:mailSubject" db:"mailSubject"`
	MailBody        string `gorm:"column:mailBody" db:"mailBody"`
	MailCheques     string `gorm:"column:mailCheques" db:"mailCheques"`
	MailItems       string `gorm:"column:mailItems" db:"mailItems"`
	IsGetAttachment int    `gorm:"column:isGetAttachment" db:"isGetAttachment"`
	IsViewDetail    int    `gorm:"column:isViewDetail" db:"isViewDetail"`
}

func (InstMail) TableName() string {
	return "inst_mail"
}

type TMails struct {
	MailID          uint32 `gorm:"primaryKey;column:mailID";json:"MailID" db:"mailID"`
	MailType        uint32 `gorm:"column:mailType";json:"MailType" db:"mailType"`
	MailDeliverTime string `gorm:"column:mailDeliverTime";json:"MailDeliverTime" db:"mailDeliverTime"`
	MailExpireTime  string `gorm:"column:mailExpireTime";json:"MailExpireTime" db:"mailExpireTime"`
	MailSubject     string `gorm:"column:mailSubject";json:"MailSubject" db:"mailSubject"`
	MailBody        string `gorm:"column:mailBody";json:"MailBody" db:"mailBody"`
	MailCheques     string `gorm:"column:mailCheques";json:"MailCheques" db:"mailCheques"`
	MailItems       string `gorm:"column:mailItems";json:"MailItems" db:"mailItems"`
	GsIds           string `gorm:"column:gsIds";json:"GsIds" db:"gsIds"`
	GsNotIds        string `gorm:"column:gsNotIds";json:"GsNotIds" db:"gsNotIds"`
	MailTargetTime  uint8  `gorm:"column:mailTargetTime";json:"MailTargetTime" db:"mailTargetTime"`
	//CancelTime      sql.NullTime `gorm:"column:cancelTime";json:"CancelTime"`
	CreateTime MyTime `gorm:"column:createTime;autoCreateTime";json:"IpcCreateTime" db:"createTime"`
}

func (TMails) TableName() string {
	return "t_mails"
}

type TimingTMail struct {
	Mail      TMails `json:"Mail" db:"mail"`
	PlayerIDs []int  `json:"PlayerIDs" db:"playerIDs"`
	GsId      int    `json:"GsId" db:"gsId"`
	ID        int    `json:"ID" db:"ID"`
}
type MailTransit struct {
	Id          int    `gorm:"primaryKey;column:Id" db:"id"`
	Record      string `gorm:"column:record" db:"record"`
	CreateTime  MyTime `gorm:"column:createTime;autoCreateTime" db:"createTime"`
	DeliverTime string `gorm:"column:deliverTime" db:"deliverTime"`
	ExpireTime  string `gorm:"column:expireTime" db:"expireTime"`
	IsSend      int    `gorm:"column:isSend" db:"isSend"`
}

func (MailTransit) TableName() string {
	return "mail_transit"
}
