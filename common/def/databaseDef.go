package def

type DatabaseCfg struct {
	ServerId     int64  `gorm:"primaryKey;column:serverId" db:"serverId"`
	LogDataBase  string `gorm:"column:logDataBase" db:"logDataBase"`
	CharDataBase string `gorm:"column:charDataBase" db:"charDataBase"`
	RedisSn      int    `gorm:"column:redisSn" db:"redisSn"`
	DeployHost   string `gorm:"column:deployHost" db:"deployHost"`
	DeployPort   string `gorm:"column:deployPort" db:"deployPort"`
	CreateTime   MyTime `gorm:"column:createTime" db:"createTime"`
	UpdateTime   MyTime `gorm:"column:updateTime" db:"updateTime"`
}

func (DatabaseCfg) TableName() string {
	return "databasecfg"
}

type DataBase struct {
	Host         string `json:"host" db:"host"`
	Port         string `json:"port" db:"port"`
	User         string `json:"user" db:"user"`
	Pwd          string `json:"pwd" db:"pwd"`
	DataBaseName string `json:"dataBaseName" db:"dataBaseName"`
}
