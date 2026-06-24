package AutoPages

import (
	"admin/common/def"
	"admin/fusion"
	"admin/fusion/giftCode"
	"fmt"
	"github.com/GoAdminGroup/go-admin/context"
	"github.com/GoAdminGroup/go-admin/modules/db"
	form2 "github.com/GoAdminGroup/go-admin/plugins/admin/modules/form"
	"github.com/GoAdminGroup/go-admin/plugins/admin/modules/table"
	"github.com/GoAdminGroup/go-admin/template/types"
	"github.com/GoAdminGroup/go-admin/template/types/form"
	"github.com/leonelquinteros/gotext"
	"html/template"
	"strconv"
	"strings"
	"time"
)

func GetTActivationCodeTable(ctx *context.Context) table.Table {

	tActivationCode := table.NewDefaultTable(table.Config{
		Driver:     "mysql",
		Connection: "db_global",
		CanAdd:     true,
		Editable:   true,
		Deletable:  true,
		Exportable: true,
		PrimaryKey: table.PrimaryKey{
			Type: db.Varchar,
			Name: "code",
		},
	})

	info := tActivationCode.GetInfo().HideDetailButton().HideEditButton().HideRowSelector()

	myOp := fusion.GetBatchNumOption()
	info.AddField(gotext.Get("礼包码"), "code", db.Char)
	info.AddField(gotext.Get("发往指定服务器"), "gsIds", db.Text)
	info.AddField(gotext.Get("不发往指定服务器"), "gsNotIds", db.Text)
	info.AddField(gotext.Get("礼包码类型"), "isGlobal", db.Tinyint).
		FieldFilterable(types.FilterType{FormType: form.SelectSingle}).
		FieldFilterOptions(types.FieldOptions{
			{Text: gotext.Get("全服"), Value: "1"},
			{Text: gotext.Get("个人"), Value: "0"}}).
		FieldDisplay(func(model types.FieldModel) interface{} {
			temp, _ := strconv.Atoi(model.Value)
			if temp == 1 {
				return gotext.Get("全服")
			} else if temp == 0 {
				return gotext.Get("个人")
			}
			return def.AmuletLog{}.GetOpStrings()[temp]
		})
	info.AddField(gotext.Get("批号"), "batchNum", db.Int).
		FieldFilterable(types.FilterType{FormType: form.SelectSingle}).
		FieldFilterOptions(myOp)
	info.AddField(gotext.Get("礼物掉落Id"), "giftId", db.Int) //对应掉落表
	info.AddField(gotext.Get("开始时间"), "startTime", db.Datetime)
	info.AddField(gotext.Get("结束时间"), "endTime", db.Datetime)
	info.AddField(gotext.Get("使用时间"), "useTime", db.Datetime)

	info.SetTable("t_activation_code").SetTitle(gotext.Get("礼包码管理"))

	ServerList := fusion.GetServerList()
	formList := tActivationCode.GetForm()
	formList.AddField(gotext.Get("生成方式"), "createType", db.Int, form.SelectSingle).
		FieldOptions(types.FieldOptions{
			{Text: gotext.Get("批量"), Value: "1"},
			{Text: gotext.Get("单个"), Value: "0"}}).
		FieldMust().
		FieldOnChooseHide("1", "code").
		FieldOnChooseShow("1", "createNum").
		FieldOnChooseShow("0", "code").
		FieldDefault("1")
	formList.AddField(gotext.Get("生成数量"), "createNum", db.Int, form.Number)
	formList.AddField(gotext.Get("礼包码"), "code", db.Char, form.Text).
		FieldHelpMsg(template.HTML(gotext.Get("礼包码必须为8位，请尽量使用大写字母和数字，尽量避免使用'I','1','l','O','0'等不易区分的字符。")))

	formList.AddField(gotext.Get("礼包码类型"), "isGlobal", db.Tinyint, form.SelectSingle).
		FieldOptions(types.FieldOptions{
			{Text: gotext.Get("全服"), Value: "1"},
			{Text: gotext.Get("个人"), Value: "0"}}).
		FieldMust()
	formList.AddField(gotext.Get("批号"), "batchNum", db.Int, form.Text).
		FieldHelpMsg(template.HTML(gotext.Get("同一个玩家只能使用一个同一批号的礼包码"))).
		FieldMust()
	formList.AddField(gotext.Get("礼物掉落Id"), "giftId", db.Int, form.Text).
		FieldHelpMsg(template.HTML(gotext.Get("礼物掉落ID对应掉落表中的掉落ID"))).
		FieldMust()
	formList.AddField(gotext.Get("开始时间"), "startTime", db.Datetime, form.Datetime).
		FieldMust()
	formList.AddField(gotext.Get("结束时间"), "endTime", db.Datetime, form.Datetime).
		FieldMust()
	formList.AddField(gotext.Get("发往指定服务器"), "gsIds", db.Text, form.SelectBox).
		FieldOptions(ServerList).
		FieldHelpMsg(template.HTML(gotext.Get("在此项选择要发往的服务器,若不填则向所有服务器发送")))
	formList.AddField(gotext.Get("不发往指定服务器"), "gsNotIds", db.Text, form.SelectBox).
		FieldOptions(ServerList).
		FieldHelpMsg(template.HTML(gotext.Get("在此项选择不发往的服务器")))

	formList.SetTable("t_activation_code").SetTitle(gotext.Get("礼包码管理")).HideContinueNewCheckBox().HideContinueEditCheckBox()

	formList.SetPostValidator(func(values form2.Values) error {
		createType, _ := strconv.Atoi(values.Get("createType"))
		if createType == 0 {
			code := values.Get("code")
			if strings.Index(code, "I") != -1 {
				return fmt.Errorf("error info")
			} else if strings.Index(code, "1") != -1 {
				return fmt.Errorf("error info")
			} else if strings.Index(code, "O") != -1 {
				return fmt.Errorf("error info")
			} else if strings.Index(code, "0") != -1 {
				return fmt.Errorf("error info")
			}
		}
		return nil
	})

	formList.SetInsertFn(func(values form2.Values) error {
		activationCode := def.ActivationCodeDef{}

		startTime := values["startTime"][0]
		endTime := values["endTime"][0]

		start, _ := time.ParseInLocation("2006-01-02 15:04:05", startTime, time.Local)
		activationCode.StartTime = def.MyTime(start)
		end, _ := time.ParseInLocation("2006-01-02 15:04:05", endTime, time.Local)
		activationCode.EndTime = def.MyTime(end)

		createType, _ := strconv.Atoi(values["createType"][0])

		createNum := 0
		if createType == 0 {
			activationCode.Code = values["code"][0]
		} else if createType == 1 {
			createNum, _ = strconv.Atoi(values["createNum"][0])
			if createNum <= 0 {
				return fmt.Errorf("wrong createNum : %d", createNum)
			}
		} else {
			return fmt.Errorf("wrong createType : %d", createType)
		}

		activationCode.BatchNum, _ = strconv.Atoi(values["batchNum"][0])
		activationCode.GiftId, _ = strconv.Atoi(values["giftId"][0])
		activationCode.IsGlobal, _ = strconv.Atoi(values["isGlobal"][0])

		temp1 := values["gsIds[]"]
		for _, v := range temp1 {
			activationCode.GsIds += string(v) + ","
		}
		activationCode.GsIds = strings.TrimSuffix(activationCode.GsIds, ",")

		temp1 = values["gsNotIds[]"]
		for _, v := range temp1 {
			activationCode.GsNotIds += string(v) + ","
		}
		activationCode.GsNotIds = strings.TrimSuffix(activationCode.GsNotIds, ",")
		db_global := fusion.GetBaseGormDB("db_global")
		if createType == 0 {
			db_global.Create(&activationCode)
		} else if createType == 1 {
			ActivationCodes := make([]def.ActivationCodeDef, 0)
			codes := giftCode.GetActivationCode(createNum)
			for _, v := range codes {
				temp := def.ActivationCodeDef{}
				temp = activationCode
				temp.Code = v
				ActivationCodes = append(ActivationCodes, temp)
			}
			db_global.Create(&ActivationCodes)
		}
		return nil
	})

	return tActivationCode
}
