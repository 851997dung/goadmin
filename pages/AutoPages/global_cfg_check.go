package AutoPages

import (
	"admin/common/def"
	"admin/common/def/currencyDef"
	"admin/fusion"
	"errors"
	"strconv"
	"time"

	"github.com/GoAdminGroup/go-admin/context"
	"github.com/GoAdminGroup/go-admin/modules/db"
	form2 "github.com/GoAdminGroup/go-admin/plugins/admin/modules/form"
	"github.com/GoAdminGroup/go-admin/plugins/admin/modules/table"
	"github.com/GoAdminGroup/go-admin/template/types"
	"github.com/GoAdminGroup/go-admin/template/types/form"
	"github.com/leonelquinteros/gotext"
)

func GetGlobalCfgCheckTable(ctx *context.Context) table.Table {

	globalCfgCheck := table.NewDefaultTable(table.Config{
		Driver:     "mysql",
		Connection: "go_manager",
		CanAdd:     true,
		Editable:   true,
		Deletable:  true,
		Exportable: true,
		PrimaryKey: table.PrimaryKey{
			Type: db.Int,
			Name: "Id",
		},
	})

	info := globalCfgCheck.GetInfo()

	cfgOps := def.GetGlobalCfgOps()
	cfgStrS := def.GetGlobalCfgStr()

	cheatTypeS := def.GetCheatTypeStr()
	cheatOps := def.GetCheatTypeOps()

	info.AddField("Id", "Id", db.Int)
	info.AddField(gotext.Get("类型"), "type", db.Int).
		FieldDisplay(func(model types.FieldModel) interface{} {
			temp, _ := strconv.Atoi(model.Value)
			return cfgStrS[temp]
		})
	info.AddField("Key", "key", db.Int).
		FieldDisplay(func(model types.FieldModel) interface{} {
			Type, _ := strconv.Atoi(fusion.Strval(model.Row["type"]))
			tempStr := make(map[int]string)
			switch Type {
			case 0:
				tempStr = currencyDef.GetCurrencyStrings()
				break
			case 1:
				tempStr = fusion.GetItemNameStrings()
				break
			case 2:
				tempStr = currencyDef.GetCurrencyStrings()
				break
			case 3:
				tempStr = fusion.GetItemNameStrings()
				break
			case 5:
				tempStr = cheatTypeS
				break
			}
			temp, _ := strconv.Atoi(model.Value)
			if tempStr[temp] != "" {
				return tempStr[temp]

			}
			return model.Value
		})

	info.AddField("Value", "value", db.Int)
	info.AddField(gotext.Get("更新时间"), "updateTime", db.Datetime)

	info.SetTable("global_cfg_check").SetTitle(gotext.Get("全局检测配置"))

	formList1 := globalCfgCheck.GetNewForm()

	itemOP := fusion.GetItemNameOptionAndID()
	currencyOP := currencyDef.GetCurrencyOps()
	formList1.AddTable(gotext.Get("票据获取总上限"), "cheque", func(panel *types.FormPanel) {
		panel.AddField(gotext.Get("票据ID"), "chequeID", db.Int, form.SelectSingle).FieldHide().FieldOptions(currencyOP)
		panel.AddField(gotext.Get("票据数量"), "chequeNum", db.Int, form.Number).FieldHide().FieldValue("0")
	}).FieldDisableWhenUpdate()
	formList1.AddTable(gotext.Get("道具获取总上限"), "Item", func(panel *types.FormPanel) {
		panel.AddField(gotext.Get("道具ID"), "ItemID", db.Int, form.SelectSingle).FieldHide().FieldOptions(itemOP)
		panel.AddField(gotext.Get("道具数量"), "ItemNum", db.Int, form.Number).FieldHide().FieldValue("0")
	}).FieldDisableWhenUpdate()

	formList1.AddTable(gotext.Get("道具获取日上限"), "ItemDaily", func(panel *types.FormPanel) {
		panel.AddField(gotext.Get("道具ID"), "ItemIDDaily", db.Int, form.SelectSingle).FieldHide().FieldOptions(itemOP)
		panel.AddField(gotext.Get("道具数量"), "ItemNumDaily", db.Int, form.Number).FieldHide().FieldValue("0")
	}).FieldDisableWhenUpdate()

	formList1.AddTable(gotext.Get("票据获取日上限"), "chequeDaily", func(panel *types.FormPanel) {
		panel.AddField(gotext.Get("票据ID"), "chequeIDDaily", db.Int, form.SelectSingle).FieldHide().FieldOptions(currencyOP)
		panel.AddField(gotext.Get("票据数量"), "chequeNumDaily", db.Int, form.Number).FieldHide().FieldValue("0")
	}).FieldDisableWhenUpdate()
	formList1.AddTable(gotext.Get("其他"), "Other", func(panel *types.FormPanel) {
		panel.AddField(gotext.Get("Key"), "OtherID", db.Int, form.SelectSingle).FieldHide().FieldOptions(cheatOps)
		panel.AddField(gotext.Get("Value"), "OtherNum", db.Int, form.Number).FieldHide().FieldValue("0")
	}).FieldDisableWhenUpdate()

	formList := globalCfgCheck.GetForm()
	formList.AddField("Id", "Id", db.Int, form.Number).
		FieldDisableWhenCreate().FieldDisplayButCanNotEditWhenUpdate()
	formList.AddField(gotext.Get("类型"), "type", db.Int, form.SelectSingle).
		FieldOptions(cfgOps).
		FieldMust().
		FieldOnChooseHide("0", "KeyItem", "ValueItem", "KeyOther", "ValueOther",
			"KeyChequeDaily", "ValueChequeDaily", "KeyItemDaily", "ValueItemDaily").
		FieldOnChooseHide("1", "KeyCheque", "ValueCheque", "KeyOther", "ValueOther",
			"KeyChequeDaily", "ValueChequeDaily", "KeyItemDaily", "ValueItemDaily").
		FieldOnChooseHide("2", "KeyCheque", "ValueCheque", "KeyOther", "ValueOther",
			"KeyItemDaily", "ValueItemDaily", "KeyItem", "ValueItem").
		FieldOnChooseHide("3", "KeyCheque", "ValueCheque", "KeyOther", "ValueOther",
			"KeyItem", "ValueItem", "KeyChequeDaily", "ValueChequeDaily").
		FieldOnChooseHide("4", "KeyItem", "ValueItem", "KeyCheque", "ValueCheque",
			"KeyChequeDaily", "ValueChequeDaily", "KeyItemDaily", "ValueItemDaily").
		FieldOnChooseShow("0", "KeyCheque", "ValueCheque").
		FieldOnChooseShow("1", "KeyItem", "ValueItem").
		FieldOnChooseShow("2", "KeyChequeDaily", "ValueChequeDaily").
		FieldOnChooseShow("3", "KeyItemDaily", "ValueItemDaily").
		FieldOnChooseShow("4", "KeyOther", "ValueOther").
		FieldDisableWhenCreate()

	formList.AddRow(func(pa *types.FormPanel) {
		myOp := currencyDef.GetCurrencyOps()
		formList.AddField(gotext.Get("票据Id"), "KeyCheque", db.Char, form.SelectSingle).
			FieldOptions(myOp).
			FieldRowWidth(3).FieldHeadWidth(2).FieldInputWidth(4).FieldDisableWhenCreate()
		formList.AddField(gotext.Get("票据获取总上限"), "ValueCheque", db.Char, form.Number).
			FieldDefault("0").
			FieldRowWidth(3).FieldHeadWidth(3).FieldInputWidth(6).FieldDisableWhenCreate()
	})
	formList.AddRow(func(pa *types.FormPanel) {
		myOp := fusion.GetItemNameOption()
		formList.AddField(gotext.Get("道具Id"), "KeyItem", db.Char, form.SelectSingle).
			FieldOptions(myOp).
			FieldRowWidth(3).FieldHeadWidth(2).FieldInputWidth(4).FieldDisableWhenCreate()
		formList.AddField(gotext.Get("道具获取总上限"), "ValueItem", db.Char, form.Number).
			FieldDefault("0").
			FieldRowWidth(3).FieldHeadWidth(3).FieldInputWidth(6).FieldDisableWhenCreate()
	})
	formList.AddRow(func(pa *types.FormPanel) {
		myOp := currencyDef.GetCurrencyOps()
		formList.AddField(gotext.Get("票据Id"), "KeyChequeDaily", db.Char, form.SelectSingle).
			FieldOptions(myOp).
			FieldRowWidth(3).FieldHeadWidth(2).FieldInputWidth(4).FieldDisableWhenCreate()
		formList.AddField(gotext.Get("票据获取日上限"), "ValueChequeDaily", db.Char, form.Number).
			FieldDefault("0").
			FieldRowWidth(3).FieldHeadWidth(3).FieldInputWidth(6).FieldDisableWhenCreate()
	})
	formList.AddRow(func(pa *types.FormPanel) {
		myOp := fusion.GetItemNameOption()
		formList.AddField(gotext.Get("道具Id"), "KeyItemDaily", db.Char, form.SelectSingle).
			FieldOptions(myOp).
			FieldRowWidth(3).FieldHeadWidth(2).FieldInputWidth(4).FieldDisableWhenCreate()
		formList.AddField(gotext.Get("道具获取日上限"), "ValueItemDaily", db.Char, form.Number).
			FieldDefault("0").
			FieldRowWidth(3).FieldHeadWidth(3).FieldInputWidth(6).FieldDisableWhenCreate()
	})
	formList.AddRow(func(pa *types.FormPanel) {
		formList.AddField(gotext.Get("key"), "KeyOther", db.Char, form.SelectSingle).
			FieldOptions(cheatOps).
			FieldRowWidth(3).FieldHeadWidth(2).FieldInputWidth(4).FieldDisableWhenCreate()
		formList.AddField(gotext.Get("Value"), "ValueOther", db.Char, form.Number).
			FieldDefault("0").
			FieldRowWidth(3).FieldHeadWidth(3).FieldInputWidth(6).FieldDisableWhenCreate()
	})

	formList.AddField(gotext.Get("更新时间"), "updateTime", db.Datetime, form.Datetime).FieldHide()

	formList.SetTable("global_cfg_check").SetTitle(gotext.Get("全局检测配置"))

	formList1.SetInsertFn(func(values form2.Values) error {

		type obj struct {
			Type       int        `gorm:"column:type" db:"type"`
			Key        int        `gorm:"column:key" db:"key"`
			Value      int        `gorm:"column:value" db:"value"`
			UpdateTime def.MyTime `gorm:"column:updateTime" db:"updateTime"`
		}
		objs := make([]obj, 0)
		formValue1 := make([]string, 0)
		formValue2 := make([]string, 0)

		formValue1 = ctx.Request.Form["chequeID"]
		formValue2 = ctx.Request.Form["chequeNum"]
		for k, _ := range formValue1 {
			Key, _ := strconv.Atoi(formValue1[k])
			Value, _ := strconv.Atoi(formValue2[k])
			if Key != 0 && Value != 0 {
				objs = append(objs, obj{
					Key:        Key,
					Value:      Value,
					Type:       0,
					UpdateTime: def.MyTime(time.Now()),
				})
			}
		}

		formValue1 = ctx.Request.Form["ItemID"]
		formValue2 = ctx.Request.Form["ItemNum"]
		for k, _ := range formValue1 {
			Key, _ := strconv.Atoi(formValue1[k])
			Value, _ := strconv.Atoi(formValue2[k])
			if Key != 0 && Value != 0 {
				objs = append(objs, obj{
					Key:        Key,
					Value:      Value,
					Type:       1,
					UpdateTime: def.MyTime(time.Now()),
				})
			}
		}

		formValue1 = ctx.Request.Form["chequeIDDaily"]
		formValue2 = ctx.Request.Form["chequeNumDaily"]
		for k, _ := range formValue1 {
			Key, _ := strconv.Atoi(formValue1[k])
			Value, _ := strconv.Atoi(formValue2[k])
			if Key != 0 && Value != 0 {
				objs = append(objs, obj{
					Key:        Key,
					Value:      Value,
					Type:       2,
					UpdateTime: def.MyTime(time.Now()),
				})
			}
		}

		formValue1 = ctx.Request.Form["ItemIDDaily"]
		formValue2 = ctx.Request.Form["ItemNumDaily"]
		for k, _ := range formValue1 {
			Key, _ := strconv.Atoi(formValue1[k])
			Value, _ := strconv.Atoi(formValue2[k])
			if Key != 0 && Value != 0 {
				objs = append(objs, obj{
					Key:        Key,
					Value:      Value,
					Type:       3,
					UpdateTime: def.MyTime(time.Now()),
				})
			}
		}

		formValue1 = ctx.Request.Form["OtherID"]
		formValue2 = ctx.Request.Form["OtherNum"]
		for k, _ := range formValue1 {
			Key, _ := strconv.Atoi(formValue1[k])
			Value, _ := strconv.Atoi(formValue2[k])
			if Key != 0 && Value != 0 {
				objs = append(objs, obj{
					Key:        Key,
					Value:      Value,
					Type:       4,
					UpdateTime: def.MyTime(time.Now()),
				})
			}
		}

		go_manager := fusion.GetBaseGormDB("go_manager")
		go_manager.Table("global_cfg_check").Create(objs)
		return nil
	})
	formList.SetUpdateFn(func(values form2.Values) error {
		Type, _ := strconv.Atoi(values.Get("type"))
		Id := values.Get("Id")
		var key, value string
		switch Type {
		case 0:
			key = values.Get("KeyCheque")
			value = values.Get("ValueCheque")
			break
		case 1:
			key = values.Get("KeyItem")
			value = values.Get("ValueItem")
			break
		case 2:
			key = values.Get("KeyOther")
			value = values.Get("ValueOther")
			break
		case 3:
			key = values.Get("KeyChequeDaily")
			value = values.Get("ValueChequeDaily")
			break
		case 4:
			key = values.Get("KeyItemDaily")
			value = values.Get("ValueItemDaily")
			break
		default:
			return errors.New(gotext.Get("something wrong"))
		}
		go_manager := fusion.GetBaseGormDB("go_manager")
		go_manager.Table("global_cfg_check").
			Where("id = ?", Id).
			Update("type", Type).
			Update("key", key).
			Update("value", value).
			Update("updateTime", def.MyTime(time.Now()))
		return nil
	})

	return globalCfgCheck
}
