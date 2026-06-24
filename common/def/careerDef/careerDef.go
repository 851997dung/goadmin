package careerDef

const (
	None        = iota
	Magician    // 魔法师
	Swordsman   // 剑士
	Archer      // 弓箭手
	Summon      // 召唤术士
	Gladiator   // 格斗家
	Spellsword  // 魔剑士
	HolyTeacher // 圣导师
	GrowLancer  // 梦幻骑士
	RuneMage    // 符文法师
	HighWnd     // 疾风
	Gunner      // 枪手
	Max
)

var Careers = []int{Magician, Swordsman, Archer, Summon, Gladiator,
	Spellsword, HolyTeacher, GrowLancer, RuneMage, HighWnd, Gunner}
