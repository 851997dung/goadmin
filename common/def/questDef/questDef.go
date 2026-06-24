package questDef

import (
	"github.com/GoAdminGroup/go-admin/template/types"
	"github.com/leonelquinteros/gotext"
	"strconv"
)

const (
	None        = iota
	MainLine    // 主线
	BranchLine  // 支线
	Daily       // 日常
	Guide       // 指引
	Transfer    // 转职
	Active      // 活跃
	ActivityMap //活动副本
	Family      // 家族任务

	Count
)
const (
	DailyQuest    = iota + 1 // 每日任务
	DailyActive              //每日活跃
	VolunteerArmy            // 义军任务
	TreasureMap              // 藏宝图任务
	Bandit                   // 江湖大盗任务
)
const (
	Accept = iota
	Finish
	Failed
	Cancel
	Submit
)

func GetQuestString() map[int]string {
	myStrs := make(map[int]string)
	myStrs[MainLine] = gotext.Get("主线")
	myStrs[BranchLine] = gotext.Get("支线")
	myStrs[Daily] = gotext.Get("日常")
	myStrs[Guide] = gotext.Get("指引")
	myStrs[Transfer] = gotext.Get("转职")
	myStrs[Active] = gotext.Get("活跃")
	myStrs[ActivityMap] = gotext.Get("活动副本")
	myStrs[MainLine] = gotext.Get("主线")
	myStrs[Family] = gotext.Get("家族任务")
	return myStrs
}
func GetQuestOptions() types.FieldOptions {
	myOps := types.FieldOptions{}
	myStr := GetQuestString()
	for k, v := range myStr {
		myOps = append(myOps, types.FieldOption{Text: v, Value: strconv.Itoa(k)})
	}
	return myOps
}
func GetQuestSubString() map[int]string {
	myStrs := make(map[int]string)
	myStrs[DailyQuest] = gotext.Get("每日任务")
	myStrs[DailyActive] = gotext.Get("每日活跃")
	myStrs[VolunteerArmy] = gotext.Get("义军任务")
	myStrs[TreasureMap] = gotext.Get("藏宝图任务")
	myStrs[Bandit] = gotext.Get("江湖大盗任务")
	return myStrs
}
func GetQuestSubOptions() types.FieldOptions {
	myOps := types.FieldOptions{}
	myStr := GetQuestSubString()
	for k, v := range myStr {
		myOps = append(myOps, types.FieldOption{Text: v, Value: strconv.Itoa(k)})
	}
	return myOps
}
func GetQuestWhenTypeString() map[int]string {
	myStrs := make(map[int]string)
	myStrs[Accept] = gotext.Get("接受")
	myStrs[Finish] = gotext.Get("结束")
	myStrs[Failed] = gotext.Get("失败")
	myStrs[Cancel] = gotext.Get("取消")
	myStrs[Submit] = gotext.Get("提交")
	return myStrs
}
func GetQuestWhenTypeOptions() types.FieldOptions {
	myOps := types.FieldOptions{}
	myStr := GetQuestWhenTypeString()
	for k, v := range myStr {
		myOps = append(myOps, types.FieldOption{Text: v, Value: strconv.Itoa(k)})
	}
	return myOps
}
