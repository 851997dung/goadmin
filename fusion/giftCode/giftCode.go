package giftCode

import "github.com/pochard/commons/randstr"

var AwardCodeDict = []rune{'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'J',
	'K', 'M', 'N', 'P', 'Q', 'R', 'S', 'T', 'U',
	'V', 'W', 'X', 'Y', 'Z',
	'1', '2', '3', '4', '5', '6', '7', '8', '9'}

func GetActivationCode(createNum int) []string {
	temp := make([]string, 0)
	for i := 0; i < createNum; i++ {
		str := randstr.Random(8, string(AwardCodeDict))
		if !ifHaveThis(temp, str) {
			temp = append(temp, str)
		} else {
			i--
		}
	}
	return temp
}

func ifHaveThis(sl []string, str string) bool {
	temp := ConvertStrSlice2Map(sl)
	return InMap(temp, str)
}

// ConvertStrSlice2Map 将字符串 slice 转为 map[string]struct{}。
func ConvertStrSlice2Map(sl []string) map[string]struct{} {
	set := make(map[string]struct{}, len(sl))
	for _, v := range sl {
		set[v] = struct{}{}
	}
	return set
}

// InMap 判断字符串是否在 map 中。
func InMap(m map[string]struct{}, s string) bool {
	_, ok := m[s]
	return ok
}
