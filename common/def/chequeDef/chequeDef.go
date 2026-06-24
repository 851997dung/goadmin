package chequeDef

const (
	None             = 0
	Exp              = 1  // 经验
	VipExp                // VIP经验
	Liveness              // 活跃度
	ExpModulus            // 系数经验
	Gold             = 11 // 金币
	BindGold              // 绑金
	Diamond               // 钻石
	BindDiamond           // 绑钻
	PotentialPoint        // 潜力点
	SlaughterValue        // 红名值
	CampScore             // 家族名声
	GuildContrib          // 战盟贡献
	GuildFund             // 战盟基金
	MilitaryExploits      // 功勋
	Max
)

func CurrencyToCheque(id int) int {
	return id + 10
}
