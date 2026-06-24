package def

type PlayerBriefInfo struct {
	IpcInstID        int    `gorm:"primaryKey;column:ipcInstID" db:"ipcInstID"`
	IpcLevel         int    `gorm:"column:ipcLevel" db:"ipcLevel"`
	IpcAcctID        int    `gorm:"column:ipcAcctID" db:"ipcAcctID"`
	IpcNickName      string `gorm:"column:ipcNickName" db:"ipcNickName"`
	IpcServerID      uint32 `gorm:"column:ipcServerID" db:"ipcServerID"`
	IpcCreateTime    int    `gorm:"column:ipcCreateTime" db:"ipcCreateTime"`
	IpcLastLoginTime int    `gorm:"column:ipcLastLoginTime" db:"ipcLastLoginTime"`
}

type PlayerInfo struct {
	IpcInstID        int    `gorm:"primaryKey;column:ipcInstID" db:"ipcInstID"`
	IpcLevel         int    `gorm:"column:ipcLevel" db:"ipcLevel"`
	IpcAcctID        int    `gorm:"column:ipcAcctID" db:"ipcAcctID"`
	IpcNickName      string `gorm:"column:ipcNickName" db:"ipcNickName"`
	IpcServerID      uint32 `gorm:"column:ipcServerID" db:"ipcServerID"`
	IpcCreateTime    int    `gorm:"column:ipcCreateTime" db:"ipcCreateTime"`
	IpcLastLoginTime int    `gorm:"column:ipcLastLoginTime" db:"ipcLastLoginTime"`
	IpcStorageItems  string `gorm:"column:ipcStorageItems" db:"ipcStorageItems"`
	IpcCurrencies    string `gorm:"column:ipcCurrencies" db:"ipcCurrencies"`
	IpcStatValue     string `gorm:"column:ipcStatValue" db:"ipcStatValue"`
}

func (PlayerInfo) TableName() string {
	return "inst_player_char"
}

type StatValue struct {
	Potentials     []int         `json:"potentials,omitempty" db:"potentials"`
	BagCapacity    int           `json:"bagCapacity" db:"bagCapacity"`
	BankCapacity   int           `json:"bankCapacity" db:"bankCapacity"`
	BagExCapacity  []interface{} `json:"bagExCapacity,omitempty" db:"bagExCapacity"`
	BankCurrencies []int         `json:"bankCurrencies,omitempty" db:"bankCurrencies"`
}
type PurchaseCut struct {
	BuyPrice     int
	BuyPriceShow string
}
