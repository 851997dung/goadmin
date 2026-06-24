package expconfigdef

import (
	"admin/common/def"
	"github.com/GoAdminGroup/go-admin/template/types"
	"github.com/leonelquinteros/gotext"
	"strconv"
)

type ExpConfig struct {
	GsId         int    `gorm:"column:gsId" db:"gsId"`
	CreatureId   int    `gorm:"primaryKey;column:CreatureId" db:"CreatureId"`
	CreatureName string `gorm:"column:CreatureName" db:"CreatureName"`
	Coefficient  int    `gorm:"column:Coefficient" db:"Coefficient"`
}

func (ExpConfig) TableName() string {
	return "t_exp_Cfg"
}

type ExcelConfig struct {
	GsId              int        `gorm:"column:gsId" db:"gsId"`
	ConfigType        int        `gorm:"column:configType" db:"configType"`
	CfgFileMd5        string     `gorm:"column:cfgFileMd5" db:"expCfgFileMd5"`
	CfgFileUpdateTime def.MyTime `gorm:"column:cfgFileUpdateTime;autoUpdateTime" db:"CfgFileUpdateTime"`
}

func (ExcelConfig) TableName() string {
	return "t_excel_config"
}

const (
	None = iota
	EnumExpConfig
	Max
)

func GetExcelConfigTypeStr() map[int]string {
	ExpConfigTypeStr := make(map[int]string)
	ExpConfigTypeStr[EnumExpConfig] = gotext.Get("经验配置")
	return ExpConfigTypeStr
}

func GetExcelConfigTypeOps() types.FieldOptions {
	myop := types.FieldOptions{}
	mystrings := GetExcelConfigTypeStr()
	for k, v := range mystrings {
		temp := types.FieldOption{}
		temp.Text = v
		temp.Value = strconv.Itoa(k)
		myop = append(myop, temp)
	}
	return myop
}

func GetTypeStr() map[int]string {
	ConfigTypeStr := make(map[int]string)
	ConfigTypeStr[EnumExpConfig] = "expConfig"
	return ConfigTypeStr
}
