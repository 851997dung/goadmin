package currencyDef

import (
	"strconv"

	"github.com/GoAdminGroup/go-admin/template/types"
	"github.com/leonelquinteros/gotext"
)

const (
	None             = iota
	Gold             // 金币
	BindGold         // 绑金
	Diamond          // 钻石
	BindDiamond      // 绑钻
	PotentialPoint   // 潜力点
	SlaughterValue   // 红名值
	CampScore        // 家族名声
	GuildContrib     // 战盟贡献
	GuildFund        // 战盟基金
	MilitaryExploits // 功勋
	Max
)

func ChequeToCurrency(id int) int {
	return id - 10
}
func GetCurrencyStrings() map[int]string {
	CurrencyStrings := make(map[int]string)
	CurrencyStrings[Gold] = gotext.Get("金币")
	//CurrencyStrings[BindGold] = gotext.Get("绑定金币")
	CurrencyStrings[Diamond] = gotext.Get("钻石")
	CurrencyStrings[BindDiamond] = gotext.Get("绑定钻石")
	return CurrencyStrings
}
func GetCurrencyOps() types.FieldOptions {
	myop := types.FieldOptions{}
	mystrings := GetCurrencyStrings()
	for k, v := range mystrings {
		temp := types.FieldOption{}
		temp.Text = v
		temp.Value = strconv.Itoa(k)
		myop = append(myop, temp)
	}
	return myop
}
