package def

import (
	"github.com/GoAdminGroup/go-admin/template/types"
	"github.com/leonelquinteros/gotext"
	"strconv"
)

type Complaint struct {
	Id           int    `gorm:"primaryKey;column:Id" db:"Id"`
	GsId         int    `gorm:"column:gsId" db:"gsId"`
	CType        int    `gorm:"column:cType" db:"cType"`
	CFrom        int    `gorm:"column:cFrom" db:"cFrom"`
	Content      string `gorm:"column:content" db:"content"`
	PlayerName   string `gorm:"column:playerName" db:"playerName"`
	LogTime      MyTime `gorm:"autoCreateTime;column:logTime" db:"logTime"`
	Status       int    `gorm:"column:status" db:"status"`
	OperatorName string `gorm:"column:operatorName" db:"operatorName"`
	Uid          int    `gorm:"column:uid" db:"uid"`
}

func (Complaint) TableName() string {
	return "complaint_management"
}

const (
	Suggestion = iota //建议
	BUG               //BUG
	Recharge          //充值问题
)

const (
	FromOperation = iota //运营转交
	BBS                  //论坛
	WebSite              //BUG
	FromGame             //游戏内提交
)

const (
	Pending = iota
	Reply
	Ignore
)

func GetComplaintTypeString() map[int]string {
	myStr := make(map[int]string)
	myStr[Suggestion] = gotext.Get("建议")
	myStr[BUG] = gotext.Get("BUG")
	myStr[Recharge] = gotext.Get("充值问题")
	return myStr
}

func GetComplaintFromString() map[int]string {
	myStr := make(map[int]string)
	myStr[FromOperation] = gotext.Get("运营转交")
	myStr[BBS] = gotext.Get("论坛")
	myStr[WebSite] = gotext.Get("官网")
	myStr[FromGame] = gotext.Get("游戏内提交")
	return myStr
}

func GetComplaintStatusString() map[int]string {
	myStr := make(map[int]string)
	myStr[Pending] = gotext.Get("待处理")
	myStr[Reply] = gotext.Get("已回复")
	myStr[Ignore] = gotext.Get("已忽略")
	return myStr
}

func GetComplaintOps(cType int) types.FieldOptions {
	myStr := make(map[int]string)
	if cType == 0 {
		myStr = GetComplaintTypeString()
	} else if cType == 1 {
		myStr = GetComplaintFromString()
	} else if cType == 2 {
		myStr = GetComplaintStatusString()
	}
	myOP := types.FieldOptions{}
	for k, v := range myStr {
		myOP = append(myOP, types.FieldOption{Value: strconv.Itoa(k), Text: v})
	}
	return myOP
}
