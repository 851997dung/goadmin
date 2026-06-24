package def

import "github.com/leonelquinteros/gotext"

const (
	Expirable = iota
	Consumable
	Transaction
	Equipment
	Ertl
	Amulet
	Wing
	Mount
	PetConverter
	Pet_Guard
	Count
)

//type ItemDetail struct {
//	Slot                       uint32
//	ItemGuid                   uint32
//	ItemTypeID                 uint32
//	ItemCount                  uint32
//	ItemOwner                  uint32
//	ItemStatus                 int32
//	ItemExpireTime             int64
//	ItemConsumableUseCount     uint32
//	ItemTransactionCount       uint32
//	ItemEquipIsHighestGrade    bool
//	ItemEquipIsUnbindEnable    bool
//	ItemEquipForgeLevel        uint32
//	ItemEquipAdditionLevel     uint32
//	ItemEquipRegenerateAttrIdx uint8
//	ItemEquipRegenerateAttrVal float32
//	ItemDurability             uint32
//	ItemWearDuration           uint32
//	ItemEquipApplyWeight       uint32
//	ItemEquipSpellIdxs         []uint8
//	ItemEquipExcellentAttrIdxs []uint8
//}

type ItemInfo struct {
	Slot                       uint32
	ItemGuid                   uint32
	ItemTypeID                 uint32
	ItemCount                  uint32
	ItemOwner                  uint32
	ItemStatus                 int32
	ItemExpireTime             int64
	ItemConsumableUseCount     uint32
	ItemTransactionCount       uint32
	ItemEquipIsHighestGrade    bool
	ItemEquipIsUnbindEnable    bool
	ItemEquipForgeLevel        uint32
	ItemEquipAdditionLevel     uint32
	ItemEquipRegenerateAttrIdx uint8
	ItemEquipRegenerateAttrVal float32
	ItemDurability             uint32
	ItemWearDuration           uint32
	ItemEquipApplyWeight       uint32
	ItemEquipSpellIdxs         []uint8
	ItemEquipExcellentAttrIdxs []uint8

}

func GetOpStrings() map[int]string {
	opStings := make(map[int]string)
	opStings[-1] = gotext.Get("普通")
	opStings[Expirable] = gotext.Get("有期限")
	opStings[Consumable] = gotext.Get("可使用")
	opStings[Transaction] = gotext.Get("可交易的")
	opStings[Equipment] = gotext.Get("装备")
	opStings[Ertl] = gotext.Get("艾尔特")
	opStings[Amulet] = gotext.Get("配饰")
	opStings[Wing] = gotext.Get("翅膀")
	opStings[Mount] = gotext.Get("坐骑")
	opStings[PetConverter] = gotext.Get("宠物转换器")
	opStings[Pet_Guard] = gotext.Get("精灵")
	return opStings
}

const (
	NewItem    = iota //新增
	DeleteItem        //删除
	AddItem           //添加
	SubItem           //分割
	SwapItem          //交换
)

const (
	Fix = iota
	Forge
	Append
	Regenerate
	Evolution
)

const (
	Inlay = iota
	Synthetic
	Enhancement
	Remove
)
const (
	upgrade = iota
	level
)
const (
	AmuletInlay = iota
	AmuletRemove
)

type PlayerItemAddError struct {
	Id               int     `gorm:"primaryKey;column:Id" db:"Id"`
	AcctId           int     `gorm:"column:acctId" db:"acctId"`
	PlayerId         int     `gorm:"column:playerId" db:"playerId"`
	PlayerName       string  `gorm:"column:playerName" db:"playerName"`
	PlayerLevel      int     `gorm:"column:playerLevel" db:"playerLevel"`
	PlayerVipLevel   int     `gorm:"column:playerVipLevel" db:"playerVipLevel"`
	PlayerFightValue int     `gorm:"column:playerFightValue" db:"playerFightValue"`
	PlayerRichValue  float64 `gorm:"column:playerRichValue" db:"playerRichValue"`
	AddFailedItems   string  `gorm:"column:addFailedItems" db:"addFailedItems"`
	FlowType         int     `gorm:"column:flowType" db:"flowType"`
	FlowParams       string  `gorm:"column:flowParams" db:"flowParams"`
	IfSolve          bool    `gorm:"column:ifSolve" db:"ifSolve"`
	LogTime          string  `gorm:"column:logTime" db:"logTime"`
}

func (PlayerItemAddError) TableName() string {
	return "log_player_item_add_error"
}
