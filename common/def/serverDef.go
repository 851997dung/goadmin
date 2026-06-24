package def

import (
	"github.com/GoAdminGroup/go-admin/template/types"
	"github.com/leonelquinteros/gotext"
	"strconv"
	"time"
)

const (
	ServerOpenStatus_Normal = iota
	ServerOpenStatus_Maintain
	ServerOpenStatus_Hide
)

const (
	ServerSpecialFlag_Hot = 1 << iota
	ServerSpecialFlag_Full
)
const (
	ServerSpecialFlag_New = 1 << (8 + iota)
	ServerSpecialFlag_Recommand
)

func GetServerOpenStatusOpString() map[int]string {
	opStrs := make(map[int]string)
	opStrs[ServerOpenStatus_Normal] = gotext.Get("普通")
	opStrs[ServerOpenStatus_Maintain] = gotext.Get("维护")
	opStrs[ServerOpenStatus_Hide] = gotext.Get("隐藏")
	return opStrs
}
func GetServerSpecialOpString() map[int]string {
	opStrs := make(map[int]string)
	opStrs[ServerSpecialFlag_Hot] = gotext.Get("火爆")
	opStrs[ServerSpecialFlag_Full] = gotext.Get("爆满")
	opStrs[ServerSpecialFlag_New] = gotext.Get("新服")
	opStrs[ServerSpecialFlag_Recommand] = gotext.Get("推荐")
	return opStrs
}
func GetServerStatusOP(genre int) types.FieldOptions {
	myOps := types.FieldOptions{}
	opSrts := make(map[int]string)
	if genre == 1 {
		opSrts = GetServerOpenStatusOpString()
	}
	if genre == 2 {
		opSrts = GetServerSpecialOpString()
	}
	for k, v := range opSrts {
		op := types.FieldOption{}
		op.Value = strconv.Itoa(k)
		op.Text = v
		myOps = append(myOps, op)
	}
	return myOps
}

type ConfigTable struct {
	Key        string `gorm:"primaryKey;column:key" db:"Key"`
	Value      string `gorm:"column:value" db:"Value"`
	UpdateTime MyTime `gorm:"column:updateTime;autoCreateTime" json:"updateTime" db:"updateTime"`
}

func (ConfigTable) TableName() string {
	return "go_config"
}

type GameServerInfo struct {
	CustomOrder       uint32    `gorm:"column:customOrder"`
	ID                uint32    `gorm:"primaryKey;column:Id"`
	ExternalIP        string    `gorm:"column:externalIP"`
	ExternalPort      uint16    `gorm:"column:externalPort"`
	InternalIP        string    `gorm:"column:internalIP"`
	InternalName      string    `gorm:"column:internalName"`
	LogicID           uint32    `gorm:"column:logicId"`
	LogicName         string    `gorm:"column:logicName"`
	LogicOpenTime     time.Time `gorm:"column:logicOpenTime" json:"logicOpenTime"`
	LogicOpenStatus   uint8     `gorm:"column:logicOpenStatus"`
	LogicSpecialFlags uint32    `gorm:"column:logicSpecialFlags"`
}

func (GameServerInfo) TableName() string {
	return "t_game_servers"
}
