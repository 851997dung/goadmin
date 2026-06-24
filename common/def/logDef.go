package def

import (
	"admin/common/def/vipSeviceType"
	"strconv"

	"github.com/GoAdminGroup/go-admin/template/types"
	"github.com/leonelquinteros/gotext"
)

type ChatLog struct {
	Id          int64  `gorm:"primaryKey;column:Id" db:"Id"`
	AcctId      int    `gorm:"column:acctId" db:"acctId"`
	PlayerId    int    `gorm:"column:characterId" db:"characterId"`
	PlayerName  string `gorm:"column:characterName" db:"characterName"`
	ChannelType int8   `gorm:"column:channelType" db:"channelType"`
	ToPlayerId  int    `gorm:"column:toPlayerId" db:"toPlayerId"`
	StrMsg      string `gorm:"column:strMsg" db:"strMsg"`
	LogTime     string `gorm:"column:logTime" db:"logTime"`
}

func (ChatLog) TableName() string {
	return "Log_character_chats"
}

type OnlinesLog struct {
	Id               int64   `gorm:"primaryKey;column:Id" db:"id"`
	AcctId           int     `gorm:"column:acctId" db:"acctId"`
	PlayerId         int     `gorm:"column:characterId" db:"characterId"`
	PlayerName       string  `gorm:"column:characterName" db:"characterName"`
	PlayerLevel      int     `gorm:"column:characterLevel" db:"characterLevel"`
	PlayerVipLevel   int     `gorm:"column:characterVipLevel" db:"characterVipLevel"`
	PlayerFightValue int     `gorm:"column:characterFightValue" db:"characterFightValue"`
	PlayerRichValue  float64 `gorm:"column:characterRichValue" db:"characterRichValue"`
	IsOnline         bool    `gorm:"column:isOnline" db:"isOnline"`
	CharacterLoginIP string  `gorm:"column:characterLoginIP" db:"characterLoginIP"`
	LogTime          string  `gorm:"column:logTime" db:"logTime"`
}

func (OnlinesLog) TableName() string {
	return "Log_character_onlines"
}

type LevelUpLog struct {
	Id               int64   `gorm:"primaryKey;column:Id" db:"id"`
	AcctId           int     `gorm:"column:acctId" db:"acctId"`
	PlayerId         int     `gorm:"column:playerId" db:"playerId"`
	PlayerName       string  `gorm:"column:playerName" db:"playerName"`
	PlayerLevel      int     `gorm:"column:playerLevel" db:"playerLevel"`
	PlayerVipLevel   int     `gorm:"column:playerVipLevel" db:"playerVipLevel"`
	PlayerFightValue int     `gorm:"column:playerFightValue" db:"playerFightValue"`
	PlayerRichValue  float64 `gorm:"column:playerRichValue" db:"playerRichValue"`
	LogTime          string  `gorm:"column:logTime" db:"logTime"`
}

func (LevelUpLog) TableName() string {
	return "Log_player_levelups"
}

type QuestLog struct {
	Id               int64   `gorm:"primaryKey;column:Id" db:"id"`
	AcctId           int     `gorm:"column:acctId" db:"acctId"`
	PlayerId         int     `gorm:"column:playerId" db:"playerId"`
	PlayerName       string  `gorm:"column:playerName" db:"playerName"`
	PlayerLevel      int     `gorm:"column:playerLevel" db:"playerLevel"`
	PlayerVipLevel   int     `gorm:"column:playerVipLevel" db:"playerVipLevel"`
	PlayerFightValue int     `gorm:"column:playerFightValue" db:"playerFightValue"`
	PlayerRichValue  float64 `gorm:"column:playerRichValue" db:"playerRichValue"`
	QuestTypeID      int     `gorm:"column:questTypeID" db:"questTypeID"`
	QuestName        string  `gorm:"column:questName" db:"questName"`
	QuestClass       int     `gorm:"column:questClass" db:"questClass"`
	QuestSubClass    int     `gorm:"column:questSubClass" db:"questSubClass"`
	WhenType         int     `gorm:"column:whenType" db:"whenType"`
	LogTime          string  `gorm:"column:logTime" db:"logTime"`
}

func (QuestLog) TableName() string {
	return "Log_player_quests"
}

type ExpLog struct {
	Id               int64   `gorm:"primaryKey;column:Id" db:"id"`
	AcctId           int     `gorm:"column:acctId" db:"acctId"`
	PlayerId         int     `gorm:"column:playerId" db:"playerId"`
	PlayerName       string  `gorm:"column:playerName" db:"playerName"`
	PlayerLevel      int     `gorm:"column:playerLevel" db:"playerLevel"`
	PlayerVipLevel   int     `gorm:"column:playerVipLevel" db:"playerVipLevel"`
	PlayerFightValue int     `gorm:"column:playerFightValue" db:"playerFightValue"`
	PlayerRichValue  float64 `gorm:"column:playerRichValue" db:"playerRichValue"`
	Exp              int     `gorm:"column:exp" db:"exp"`
	Count            int     `gorm:"column:count" db:"count"`
	FlowType         int     `gorm:"column:flowType" db:"flowType"`
	FlowParams       string  `gorm:"column:flowParams" db:"flowParams"`
	//LogStartTime     string  `gorm:"column:logStartTime"`
	LogTime string `gorm:"column:logStopTime" db:"logStopTime"`
}

func (ExpLog) TableName() string {
	return "log_player_exps_"
}

type MoneyLog struct {
	Id               int64   `gorm:"primaryKey;column:Id" db:"id"`
	AcctId           int     `gorm:"column:acctId" db:"acctId"`
	PlayerId         int     `gorm:"column:playerId" db:"playerId"`
	PlayerName       string  `gorm:"column:playerName" db:"playerName"`
	PlayerLevel      int     `gorm:"column:playerLevel" db:"playerLevel"`
	PlayerVipLevel   int     `gorm:"column:playerVipLevel" db:"playerVipLevel"`
	PlayerFightValue int     `gorm:"column:playerFightValue" db:"playerFightValue"`
	PlayerRichValue  float64 `gorm:"column:playerRichValue" db:"playerRichValue"`
	MoneyType        int     `gorm:"column:moneyType" db:"moneyType"`
	MoneyName        string  `gorm:"column:moneyName" db:"moneyName"`
	MoneyValue       int64   `gorm:"column:moneyValue" db:"moneyValue"`
	MoneyCurValue    int64   `gorm:"column:moneyCurValue" db:"moneyCurValue"`
	FlowType         int     `gorm:"column:flowType" db:"flowType"`
	FlowParams       string  `gorm:"column:flowParams" db:"flowParams"`
	LogTime          string  `gorm:"column:logTime" db:"logTime"`
}

func (MoneyLog) TableName() string {
	return "log_player_moneys_"
}

type ItemLog struct {
	Id               int64   `gorm:"primaryKey;column:Id" db:"id"`
	AcctId           int     `gorm:"column:acctId" db:"acctId"`
	PlayerId         int     `gorm:"column:playerId" db:"playerId"`
	PlayerName       string  `gorm:"column:playerName" db:"playerName"`
	PlayerLevel      int     `gorm:"column:playerLevel" db:"playerLevel"`
	PlayerVipLevel   int     `gorm:"column:playerVipLevel" db:"playerVipLevel"`
	PlayerFightValue int     `gorm:"column:playerFightValue" db:"playerFightValue"`
	PlayerRichValue  float64 `gorm:"column:playerRichValue" db:"playerRichValue"`
	ItemSlotType     int     `gorm:"column:itemSlotType" db:"itemSlotType"`
	ItemSlot         int     `gorm:"column:itemSlot" db:"itemSlot"`
	ItemConstGuid    int     `gorm:"column:itemConstGuid" db:"itemConstGuid"`
	ItemTypeId       int     `gorm:"column:itemTypeId" db:"itemTypeId"`
	ItemName         string  `gorm:"column:itemName" db:"itemName"`
	ItemNum          float64 `gorm:"column:itemNum" db:"itemNum"`
	ItemCurNum       float64 `gorm:"column:itemCurNum" db:"itemCurNum"`
	FlowType         int     `gorm:"column:flowType" db:"flowType"`
	FlowParams       string  `gorm:"column:flowParams" db:"flowParams"`
	LogType          int     `gorm:"column:logType" db:"logType"`
	LogTime          string  `gorm:"column:logTime" db:"logTime"`
}

func (ItemLog) TableName() string {
	return "log_player_items_"
}
func (ItemLog) GetOpStrings() map[int]string {
	op := make(map[int]string)
	op[NewItem] = gotext.Get("新道具")
	op[DeleteItem] = gotext.Get("移除道具")
	op[AddItem] = gotext.Get("添加道具")
	op[SubItem] = gotext.Get("分割道具")
	op[SwapItem] = gotext.Get("交换道具")
	op[166] = gotext.Get("个人商店")
	return op
}
func (s ItemLog) GetOptions() types.FieldOptions {
	myOps := types.FieldOptions{}
	logMap := s.GetOpStrings()
	for k, v := range logMap {
		op := types.FieldOption{}
		op.Value = strconv.Itoa(k)
		op.Text = v
		myOps = append(myOps, op)
	}
	return myOps
}

type DungeonsLog struct {
	Id               int64   `gorm:"primaryKey;column:Id" db:"id"`
	AcctId           int     `gorm:"column:acctId" db:"acctId"`
	PlayerId         int     `gorm:"column:playerId" db:"playerId"`
	PlayerName       string  `gorm:"column:playerName" db:"playerName"`
	PlayerLevel      int     `gorm:"column:playerLevel" db:"playerLevel"`
	PlayerVipLevel   int     `gorm:"column:playerVipLevel" db:"playerVipLevel"`
	PlayerFightValue int     `gorm:"column:playerFightValue" db:"playerFightValue"`
	PlayerRichValue  float64 `gorm:"column:playerRichValue" db:"playerRichValue"`
	DungeonId        int     `gorm:"column:dungeonId" db:"dungeonId"`
	Progress         int     `gorm:"column:progress" db:"progress"`
	LogTime          string  `gorm:"column:logTime" db:"logTime"`
}

func (DungeonsLog) TableName() string {
	return "Log_player_dungeons"
}

func (DungeonsLog) GetOpStrings() map[int]string {
	op := make(map[int]string)
	op[vipSeviceType.DemonSquare] = gotext.Get("恶魔广场")
	op[vipSeviceType.BloodyCastle] = gotext.Get("血色城堡")
	op[vipSeviceType.RedFortress] = gotext.Get("赤色要塞")
	op[vipSeviceType.DailyQuest] = gotext.Get("每日任务")
	op[vipSeviceType.PhantomTemple] = gotext.Get("幻影寺院")
	op[vipSeviceType.SoulSquare] = gotext.Get("生魂广场")
	op[vipSeviceType.DailyActive] = gotext.Get("每日活跃")
	op[vipSeviceType.WolfSpiritSortress] = gotext.Get("狼魂要塞")
	op[vipSeviceType.FamilyQuest] = gotext.Get("家族任务")
	op[vipSeviceType.ImperialFortress] = gotext.Get("帝国要塞")
	op[vipSeviceType.RefineTower] = gotext.Get("提炼之塔")
	op[vipSeviceType.Maze] = gotext.Get("迷宫")
	op[vipSeviceType.Siege] = gotext.Get("攻城战")
	return op
}

type StoneLog struct {
	Id               int64   `gorm:"primaryKey;column:Id" db:"id"`
	AcctId           int     `gorm:"column:acctId" db:"acctId"`
	PlayerId         int     `gorm:"column:playerId" db:"playerId"`
	PlayerName       string  `gorm:"column:playerName" db:"playerName"`
	PlayerLevel      int     `gorm:"column:playerLevel" db:"playerLevel"`
	PlayerVipLevel   int     `gorm:"column:playerVipLevel" db:"playerVipLevel"`
	PlayerFightValue int     `gorm:"column:playerFightValue" db:"playerFightValue"`
	PlayerRichValue  float64 `gorm:"column:playerRichValue" db:"playerRichValue"`
	LogType          int     `gorm:"column:logType" db:"logType"`
	GemItemTypeID    int     `gorm:"column:gemItemTypeID" db:"gemItemTypeID"`
	HoleID           int     `gorm:"column:holeID" db:"holeID"`
	GemsLucky        string  `gorm:"column:gemsLucky" db:"gemsLucky"`
	GemLevel         int     `gorm:"column:gemLevel" db:"gemLevel"`
	LogTime          string  `gorm:"column:logTime" db:"logTime"`
}

func (StoneLog) TableName() string {
	return "log_player_stone_"
}

func (StoneLog) GetOpStrings() map[int]string {
	op := make(map[int]string)
	op[Inlay] = gotext.Get("镶嵌")
	op[Synthetic] = gotext.Get("合成")
	op[Enhancement] = gotext.Get("强化")
	op[Remove] = gotext.Get("移除")
	return op
}
func (s StoneLog) GetOptions() types.FieldOptions {
	myOps := types.FieldOptions{}
	logMap := s.GetOpStrings()
	for k, v := range logMap {
		op := types.FieldOption{}
		op.Value = strconv.Itoa(k)
		op.Text = v
		myOps = append(myOps, op)
	}
	return myOps
}

type EquipLog struct {
	Id               int64   `gorm:"primaryKey;column:Id" db:"id"`
	AcctId           int     `gorm:"column:acctId" db:"acctId"`
	PlayerId         int     `gorm:"column:playerId" db:"playerId"`
	PlayerName       string  `gorm:"column:playerName" db:"playerName"`
	PlayerLevel      int     `gorm:"column:playerLevel" db:"playerLevel"`
	PlayerVipLevel   int     `gorm:"column:playerVipLevel" db:"playerVipLevel"`
	PlayerFightValue int     `gorm:"column:playerFightValue" db:"playerFightValue"`
	PlayerRichValue  float64 `gorm:"column:playerRichValue" db:"playerRichValue"`
	LogType          int     `gorm:"column:logType" db:"logType"`
	FixEquipType     int     `gorm:"column:fixEquipType" db:"fixEquipType"`
	//OldLevel         int     `gorm:"column:oldLevel" db:"oldLevel"`
	AfterOpLevel  int    `gorm:"column:afterOpLevel" db:"afterOpLevel"`
	ChangeAttrIdx int    `gorm:"column:changeAttrIdx" db:"changeAttrIdx"`
	ChangeAttrVal int    `gorm:"column:changeAttrVal" db:"changeAttrVal"`
	LogTime       string `gorm:"column:logTime" db:"logTime"`
}

func (EquipLog) TableName() string {
	return "log_player_equip_"
}

func (EquipLog) GetOpStrings() map[int]string {
	op := make(map[int]string)
	op[Fix] = gotext.Get("修理")
	op[Forge] = gotext.Get("强化")
	op[Append] = gotext.Get("追加")
	op[Regenerate] = gotext.Get("重生")
	op[Evolution] = gotext.Get("进化")
	return op
}
func (s EquipLog) GetOptions() types.FieldOptions {
	myOps := types.FieldOptions{}
	logMap := s.GetOpStrings()
	for k, v := range logMap {
		op := types.FieldOption{}
		op.Value = strconv.Itoa(k)
		op.Text = v
		myOps = append(myOps, op)
	}
	return myOps
}

type AmuletLog struct {
	Id               int64   `gorm:"primaryKey;column:Id" db:"id"`
	AcctId           int     `gorm:"column:acctId" db:"acctId"`
	PlayerId         int     `gorm:"column:playerId" db:"playerId"`
	PlayerName       string  `gorm:"column:playerName" db:"playerName"`
	PlayerLevel      int     `gorm:"column:playerLevel" db:"playerLevel"`
	PlayerVipLevel   int     `gorm:"column:playerVipLevel" db:"playerVipLevel"`
	PlayerFightValue int     `gorm:"column:playerFightValue" db:"playerFightValue"`
	PlayerRichValue  float64 `gorm:"column:playerRichValue" db:"playerRichValue"`
	LogType          int     `gorm:"column:logType" db:"logType"`
	ItemErtl         string  `gorm:"column:itemErtl" db:"itemErtl"`
	LogTime          string  `gorm:"column:logTime" db:"logTime"`
	ErtlTypeID       string  `gorm:"column:ertlTypeID" db:"ertlTypeID"`
}

func (AmuletLog) TableName() string {
	return "log_player_amulet_"
}

func (AmuletLog) GetOpStrings() map[int]string {
	op := make(map[int]string)
	op[AmuletInlay] = gotext.Get("镶嵌")
	op[AmuletRemove] = gotext.Get("移除")
	return op
}
func (s AmuletLog) GetOptions() types.FieldOptions {
	myOps := types.FieldOptions{}
	logMap := s.GetOpStrings()
	for k, v := range logMap {
		op := types.FieldOption{}
		op.Value = strconv.Itoa(k)
		op.Text = v
		myOps = append(myOps, op)
	}
	return myOps
}

type ErtlLog struct {
	Id               int64   `gorm:"primaryKey;column:Id" db:"id"`
	AcctId           int     `gorm:"column:acctId" db:"acctId"`
	PlayerId         int     `gorm:"column:playerId" db:"playerId"`
	PlayerName       string  `gorm:"column:playerName" db:"playerName"`
	PlayerLevel      int     `gorm:"column:playerLevel" db:"playerLevel"`
	PlayerVipLevel   int     `gorm:"column:playerVipLevel" db:"playerVipLevel"`
	PlayerFightValue int     `gorm:"column:playerFightValue" db:"playerFightValue"`
	PlayerRichValue  float64 `gorm:"column:playerRichValue" db:"playerRichValue"`
	LogType          int     `gorm:"column:logType" db:"logType"`
	ErtlCurLevel     int     `gorm:"column:ertlCurLevel" db:"ertlCurLevel"`
	ErtlCurRating    int     `gorm:"column:ertlCurRating" db:"ertlCurRating"`
	AwardAttr        string  `gorm:"column:awardAttr" db:"awardAttr"`
	LogTime          string  `gorm:"column:logTime" db:"logTime"`
	ErtlTypeID       int     `gorm:"column:ertlTypeID" db:"ertlTypeID"`
}

func (ErtlLog) TableName() string {
	return "log_player_Ertl_"
}
func (ErtlLog) GetOpStrings() map[int]string {
	op := make(map[int]string)
	op[upgrade] = gotext.Get("升阶")
	op[level] = gotext.Get("升级")
	return op
}
func (s ErtlLog) GetOptions() types.FieldOptions {
	myOps := types.FieldOptions{}
	logMap := s.GetOpStrings()
	for k, v := range logMap {
		op := types.FieldOption{}
		op.Value = strconv.Itoa(k)
		op.Text = v
		myOps = append(myOps, op)
	}
	return myOps
}

const (
	tInvalid = -1
	tWhisper = 0 // 私聊
	//tFriend  = 1 //好友
	tSpeak  = 1 // 区域
	tTeam   = 2 // 队伍
	tGuild  = 3 // 工会
	tCamp   = 4 // 阵营
	tWorld  = 5 // 世界
	tServer = 6 // 系统
	//tMarquee = 6 //跑马灯
	//tInstant = 7 //即时公告
	//tSpecial = 8 //特殊提示
)

// 0=私聊，1=好友，2=附近，3=队伍，4=家族，5=门派，6=世界，7=系统，8=跑马灯，9=即时公告，10=特殊提示
func GetOpStringsChannel() map[int]string {
	op := make(map[int]string)
	op[tWhisper] = gotext.Get("私聊")
	//op[tFriend] = gotext.Get("好友")
	op[tSpeak] = gotext.Get("区域")
	op[tTeam] = gotext.Get("队伍")
	op[tGuild] = gotext.Get("工会")
	op[tCamp] = gotext.Get("阵营")
	op[tWorld] = gotext.Get("世界")
	op[tServer] = gotext.Get("系统")
	//op[tMarquee] = gotext.Get("跑马灯")
	//op[tInstant] = gotext.Get("即时公告")
	//op[tSpecial] = gotext.Get("特殊提示")
	return op
}

func GetChannelOptions() types.FieldOptions {
	myOps := types.FieldOptions{}
	strs := GetOpStringsChannel()
	for k, v := range strs {
		if k != tServer {
			continue
		}
		op := types.FieldOption{}
		op.Value = strconv.Itoa(k)
		op.Text = v
		myOps = append(myOps, op)
	}
	return myOps
}

const (
	Marquee     = 1 << 0  // 跑马灯
	Notice      = 1 << 1  // 通知
	Instant     = 1 << 2  // 即时
	Bubble      = 1 << 3  // 冒泡
	Window      = 1 << 4  // 弹窗
	FormatTwice = 1 << 16 // 格式化两次
)

func GetOpStringsMsgFlag() map[int]string {
	op := make(map[int]string)
	op[Marquee] = gotext.Get("跑马灯")
	op[Notice] = gotext.Get("通知")
	op[Instant] = gotext.Get("即时")
	op[Bubble] = gotext.Get("冒泡")
	op[Window] = gotext.Get("弹窗")
	//op[FormatTwice] = gotext.Get("格式化两次")
	return op
}

func GetMsgFlagOptions() types.FieldOptions {
	myOps := types.FieldOptions{}
	strs := GetOpStringsMsgFlag()
	for k, v := range strs {
		op := types.FieldOption{}
		op.Value = strconv.Itoa(k)
		op.Text = v
		myOps = append(myOps, op)
	}
	return myOps
}

const (
	Log_character_chats = iota
	Log_character_onlines
	Log_player_amulet
	Log_player_dungeons
	Log_player_equip
	Log_player_ertl
	Log_player_exps
	Log_player_items
	Log_player_levelups
	Log_player_moneys
	Log_player_quests
	Log_player_stone
	Log_player_login
)

func GetTableNameStrings() map[int]string {
	NameStrings := make(map[int]string)
	NameStrings[Log_character_chats] = gotext.Get("聊天日志")
	NameStrings[Log_character_onlines] = gotext.Get("登录日志表格")
	NameStrings[Log_player_amulet] = gotext.Get("配饰镶嵌记录")
	NameStrings[Log_player_dungeons] = gotext.Get("副本日志")
	NameStrings[Log_player_equip] = gotext.Get("装备日志")
	NameStrings[Log_player_ertl] = gotext.Get("艾尔特日志")
	NameStrings[Log_player_exps] = gotext.Get("经验获取日志")
	NameStrings[Log_player_items] = gotext.Get("道具日志")
	NameStrings[Log_player_levelups] = gotext.Get("升级日志")
	NameStrings[Log_player_moneys] = gotext.Get("货币日志")
	NameStrings[Log_player_quests] = gotext.Get("任务日志")
	NameStrings[Log_player_stone] = gotext.Get("荧之石日志")
	NameStrings[Log_player_login] = gotext.Get("日志")
	return NameStrings
}

func GetTableNameOp() types.FieldOptions {
	myOptions := types.FieldOptions{}
	temp := GetTableNameStrings()
	for k, v := range temp {
		tempOp := types.FieldOption{}
		tempOp.Value = strconv.Itoa(k)
		tempOp.Text = v
		myOptions = append(myOptions, tempOp)
	}
	return myOptions
}

const (
	Cfg_type_Cheque = iota
	Cfg_type_Item
	Cfg_type_chequeDaily
	Cfg_type_ItemDaily
	Cfg_type_Other
)

func GetGlobalCfgStr() map[int]string {
	str := make(map[int]string)
	str[0] = gotext.Get("票据获取总上限")
	str[1] = gotext.Get("道具获取总上限")
	str[2] = gotext.Get("票据获取日上限")
	str[3] = gotext.Get("道具获取日上限")
	str[4] = gotext.Get("其他")
	return str
}
func GetGlobalCfgOps() types.FieldOptions {
	myOptions := types.FieldOptions{}
	temp := GetGlobalCfgStr()
	for k, v := range temp {
		tempOp := types.FieldOption{}
		tempOp.Value = strconv.Itoa(k)
		tempOp.Text = v
		myOptions = append(myOptions, tempOp)
	}
	return myOptions
}

func GetCheatTypeStr() map[int]string {
	str := make(map[int]string)
	str[0] = gotext.Get("攻击速度")
	str[1] = gotext.Get("移动速度")
	return str
}
func GetCheatTypeOps() types.FieldOptions {
	myOptions := types.FieldOptions{}
	temp := GetCheatTypeStr()
	for k, v := range temp {
		tempOp := types.FieldOption{}
		tempOp.Value = strconv.Itoa(k)
		tempOp.Text = v
		myOptions = append(myOptions, tempOp)
	}
	return myOptions
}

type SuspiciousDataRecord struct {
	Id         int    `gorm:"primaryKey;column:Id" db:"id"`
	GsId       int    `gorm:"column:gsId" db:"gsId"`
	AccountId  int    `gorm:"column:accountId" db:"accountId"`
	PlayerId   int    `gorm:"column:playerId" db:"playerId"`
	PlayerName string `gorm:"column:playerName" db:"playerName"`
	Type       int    `gorm:"column:type" db:"type"`
	Key        int    `gorm:"column:key" db:"key"`
	Value      int    `gorm:"column:value" db:"value"`
	LogTime    MyTime `gorm:"column:logTime" db:"logTime"`
}

func (SuspiciousDataRecord) TableName() string {
	return "suspicious_data_record"
}

type LoginLog struct {
	DeviceUniqueIdentifier string `gorm:"primaryKey;column:deviceUniqueIdentifier" db:"deviceUniqueIdentifier"` //nolint:unused
	DeviceModel            string `gorm:"column:deviceModel" db:"deviceModel"`
	LastTime               string `gorm:"column:lastTime" db:"lastTime"`
	LastIP                 string `gorm:"column:lastIP" db:"lastIP"`
	AcctId                 int    `gorm:"column:accountId" db:"accountId"`
	IpcServerID            int    `gorm:"column:serverId" db:"serverId"` //nolint:unused                           //nolint:unused
	CurStep                string `gorm:"column:lastStep" db:"lastStep"` //nolint:unused
	FinishStep             string `gorm:"column:finalStep" db:"finalStep"`
	IsPerfectPlay          int    `gorm:"column:isPerfectPlay" db:"isPerfectPlay"`   //nolint:unused
	DayFirstLoginTime      string `gorm:"column:firstStartTime" db:"firstStartTime"` //nolint:unused
	PackageStatus          string `gorm:"column:packageStatus" db:"packageStatus"`   //nolint:unused
	UpdateState            int    `gorm:"column:updateStatus" db:"updateStatus"`
	ResourceLoad           int    `gorm:"column:resLoadStatus" db:"resLoadStatus"` //nolint:unused
	SDKAcct                int    `gorm:"column:sdkLoginStatus" db:"sdkLoginStatus"`
	ServerListStatus       int    `gorm:"column:serverListStatus" db:"serverListStatus"` //nolint:unused
	ServerConnectStatus    int    `gorm:"column:serverConnectStatus" db:"serverConnectStatus"`
	//nolint:unused
	EnterGameStatus    int `gorm:"column:enterGameStatus" db:"enterGameStatus"`
	SdkLoginFailCount  int `gorm:"column:sdkLoginFailCount" db:"sdkLoginFailCount"`           //nolint:unused
	GetSerListFaildNum int `gorm:"column:serverListFailCount" db:"serverListFailCount"`       //nolint:unused
	ConnSerFaildNum    int `gorm:"column:connectServerFailCount" db:"connectServerFailCount"` //nolint:unused

}

func (LoginLog) TableName() string {
	return "t_device_logs_"
}

func (LoginLog) GetOpStrings() map[string]string {
	op := make(map[string]string)
	op["SystemInfo"] = gotext.Get("获取系统信息")
	op["Package Status"] = gotext.Get("获取包状态")
	op["GetConfig from"] = gotext.Get("获取配置信息")
	op["GetConfig succeeded"] = gotext.Get("获取配置成功")
	op["No hotfix resource"] = gotext.Get("无需热更")
	op["Load hotfix dll"] = gotext.Get("下载热更资源")
	op["Start game"] = gotext.Get("开始进入游戏步骤")
	op["Start play hotfixVersion"] = gotext.Get("开始热更版本")
	op["Start LoadDB"] = gotext.Get("开始加载DB")
	op["LoadDB succeeded"] = gotext.Get("加载DB成功")
	op["Open LoginForm succeeded"] = gotext.Get("打开登录面板成功")
	op["Start Login"] = gotext.Get("开始登录")
	op["Login succeeded"] = gotext.Get("登录成功")
	op["Connect Server"] = gotext.Get("连接服务器")
	op["Connect succeeded"] = gotext.Get("连接服务器成功")
	op["Send MmorpgLogin"] = gotext.Get("发送登录请求")
	op["Send GetCharacterList"] = gotext.Get("获取角色列表")
	op["Open RoleSelectorForm succeeded"] = gotext.Get("打开选角界面成功")
	op["CharacterListResp Count"] = gotext.Get("获取角色列表成功")
	op["Login Character"] = gotext.Get("登录角色")
	op["LoadMap Id"] = gotext.Get("加载地图")
	op["Load Hero Packet"] = gotext.Get("加载角色信息")
	op["Load Hero succeeded"] = gotext.Get("加载角色信息成功")
	op["LoadMap succeeded"] = gotext.Get("加载地图成功")
	op["Perfect Play"] = gotext.Get("成功游玩游戏")
	return op
}
func (LoginLog) GetStatusStrings() map[int]string {
	op := make(map[int]string)
	op[0] = gotext.Get("失败")
	op[1] = gotext.Get("成功")
	return op
}
