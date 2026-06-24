package pageMgr

import (
	"admin/common/def/flowDef"
	"github.com/GoAdminGroup/go-admin/template/types"
	"strconv"
)

func GetFlowTypeOptions() types.FieldOptions {
	myOptions := types.FieldOptions{}
	str := flowDef.GetFlowTypeStrings()
	for k, v := range str {
		Option := types.FieldOption{}
		Option.Value = strconv.Itoa(k)
		Option.Text = v
		myOptions = append(myOptions, Option)
	}
	return myOptions
}
