package AutoPages

import (
	"admin/common"
	"admin/common/def/excelConfigDef/expconfigdef"
	"admin/fusion"
	"errors"
	"github.com/GoAdminGroup/go-admin/context"
	"github.com/GoAdminGroup/go-admin/modules/db"
	form2 "github.com/GoAdminGroup/go-admin/plugins/admin/modules/form"
	"github.com/GoAdminGroup/go-admin/plugins/admin/modules/table"
	"github.com/GoAdminGroup/go-admin/template/types"
	"github.com/GoAdminGroup/go-admin/template/types/form"
	"github.com/leonelquinteros/gotext"
	"net/url"
	"os"
	"strconv"
)

func GetTExcelConfigTable(ctx *context.Context) table.Table {

	tExcelConfig := table.NewDefaultTable(table.Config{
		Driver:     "mysql",
		Connection: "go_manager",
		CanAdd:     true,
		Editable:   true,
		Deletable:  true,
		Exportable: true,
		PrimaryKey: table.PrimaryKey{
			Type: db.Int,
			Name: "gsId",
		},
	})

	info := tExcelConfig.GetInfo().HideEditButton().HideDetailButton().HideDeleteButton()

	ops := expconfigdef.GetExcelConfigTypeOps()
	strS := expconfigdef.GetExcelConfigTypeStr()
	serverOps := fusion.GetServerList()
	serverName := fusion.GetServerNameStrings()

	info.AddField(gotext.Get("服务器ID"), "gsId", db.Int).
		FieldFilterable(types.FilterType{FormType: form.SelectSingle}).
		FieldFilterOptions(serverOps).
		FieldDisplay(func(model types.FieldModel) interface{} {
			temp, _ := strconv.Atoi(model.Value)
			str, err := serverName[temp]
			if !err {
				return model.Value
			}
			return str
		})
	info.AddField(gotext.Get("配置文件类型"), "configType", db.Int).
		FieldDisplay(func(model types.FieldModel) interface{} {
			temp, _ := strconv.Atoi(model.Value)
			str, err := strS[temp]
			if !err {
				return model.Value
			}
			return str
		}).
		FieldFilterable(types.FilterType{FormType: form.SelectSingle}).
		FieldFilterOptions(ops)
	info.AddField(gotext.Get("文件MD5"), "cfgFileMd5", db.Varchar)
	info.AddField(gotext.Get("文件更新时间"), "cfgFileUpdateTime", db.Datetime)
	info.SetTable("t_excel_config").SetTitle(gotext.Get("配置文件相关"))

	//formList := tExcelConfig.GetForm()
	//formList.AddField(gotext.Get("服务器id"), "gsId", db.Int, form.Number).
	//	FieldDisplay(func(model types.FieldModel) interface{} {
	//		temp, _ := strconv.Atoi(model.Value)
	//		str, err := serverName[temp]
	//		if !err {
	//			return model.Value
	//		}
	//		return str
	//	})
	//
	//formList.AddField(gotext.Get("配置类型"), "configType", db.Int, form.SelectSingle).
	//	FieldOptions(ops)
	//formList.AddField(gotext.Get("上传配置"), "FileName", db.Varchar, form.File)
	//formList.SetUpdateFn(func(values form2.Values) error {
	//	excelConfig := expconfigdef.ExcelConfig{}
	//	excelConfig.GsId, _ = strconv.Atoi(values.Get("gsId"))
	//	excelConfig.ConfigType, _ = strconv.Atoi(values.Get("configType"))
	//	FileName := values.Get("FileName")
	//	excelConfig.CfgFileMd5 = fusion.GetFileMd5(common.Cfg_yml.Store.Path + "/" + FileName)
	//	go_manager := fusion.GetBaseGormDB("go_manager")
	//	go_manager.Save(excelConfig)
	//	return nil
	//})
	//formList.SetTable("t_excel_config").SetTitle(gotext.Get("配置文件相关"))

	newFormList := tExcelConfig.GetNewForm()
	newFormList.AddField(gotext.Get("配置类型"), "configType", db.Int, form.SelectSingle).
		FieldOptions(ops).FieldMust()
	newFormList.AddField(gotext.Get("上传配置"), "FileName", db.Varchar, form.File)
	newFormList.AddField(gotext.Get("服务器id"), "gsId", db.Int, form.SelectBox).
		FieldOptions(serverOps).FieldMust()
	newFormList.SetInsertFn(func(values form2.Values) error {
		gsIds := values["gsId[]"]
		excelConfigs := make([]expconfigdef.ExcelConfig, 0)
		ConfigType, _ := strconv.Atoi(values.Get("configType"))
		FileName := values.Get("FileName")
		md5 := fusion.GetFileMd5(common.Cfg_yml.Store.Path + "/" + FileName)
		err, temp := fusion.ReadExcel(common.Cfg_yml.Store.Path+"/"+FileName, ConfigType)

		defer func() {
			os.Remove(common.Cfg_yml.Store.Path + "/" + FileName)
		}()

		ids := make([]int, 0)
		if err != nil {
			return err
		}
		for _, v := range gsIds {
			gsId, _ := strconv.Atoi(v)
			excelConfig := expconfigdef.ExcelConfig{}
			excelConfig.GsId = gsId
			excelConfig.ConfigType = ConfigType
			excelConfig.CfgFileMd5 = md5
			excelConfigs = append(excelConfigs, excelConfig)
			ids = append(ids, gsId)
		}

		if ConfigType == expconfigdef.EnumExpConfig {
			expConfigs := make([]expconfigdef.ExpConfig, 0)
			for _, v := range gsIds {
				gsId, _ := strconv.Atoi(v)
				for _, data := range temp {
					config := expconfigdef.ExpConfig{}
					config.GsId = gsId
					config.CreatureId = data.CreatureId
					config.CreatureName = data.CreatureName
					config.Coefficient = data.Coefficient
					expConfigs = append(expConfigs, config)
				}
			}
			db_golbal := fusion.GetBaseGormDB("db_global")
			db_golbal.Table("t_exp_cfg").Where("gsId in ?", gsIds).Delete(expconfigdef.ExpConfig{})
			db_golbal.Create(&expConfigs)
		}
		go_manager := fusion.GetBaseGormDB("go_manager")
		go_manager.Table("t_excel_config").Where("configType = ?", ConfigType).Where("gsId in ?", gsIds).Delete(expconfigdef.ExcelConfig{})
		go_manager.Create(&excelConfigs)

		var result string

		for _, v := range gsIds {
			vars := &url.Values{}
			vars.Add("gsId", v)
			ctx1 := &context.Context{}
			err, res := fusion.CallToCenter(vars, common.MyCfg.Api["PushExtraConfig"], "GET", ctx1)
			if err != nil {
				result += err.Error() + "\n"
			}
			if res != "" && res != "All Is OK!" {
				result += res + "\n"
			}
		}
		if len(result) != 0 {
			return errors.New(result)
		}
		return nil
	})
	newFormList.SetTable("t_excel_config").SetTitle(gotext.Get("配置文件相关"))

	return tExcelConfig
}
