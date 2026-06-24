package AutoPages

import (
	"github.com/GoAdminGroup/go-admin/context"
	"github.com/GoAdminGroup/go-admin/modules/db"
	"github.com/GoAdminGroup/go-admin/plugins/admin/modules/table"
	"github.com/GoAdminGroup/go-admin/template/types/form"
)

func GetLogPersonalShopTable(ctx *context.Context) table.Table {

	logPersonalShop := table.NewDefaultTable(table.DefaultConfigWithDriverAndConnection("mysql", "db_log"))

	info := logPersonalShop.GetInfo()

	info.AddField("Id", "Id", db.Int)
	info.AddField("PlayerId", "PlayerId", db.Int)
	info.AddField("TradPersonName", "tradPersonName", db.Text)
	info.AddField("BundleNum", "bundleNum", db.Int)
	info.AddField("SoulGemNum", "soulGemNum", db.Int)
	info.AddField("BlessingGemNum", "blessingGemNum", db.Int)
	info.AddField("DiamondsNum", "DiamondsNum", db.Bigint)
	info.AddField("Price", "price", db.Bigint)
	info.AddField("IsBuy", "isBuy", db.Tinyint)
	info.AddField("ItemQualitys", "itemQualitys", db.JSON)
	info.AddField("ItemLevels", "itemLevels", db.JSON)
	info.AddField("ItemTypeIDs", "itemTypeIDs", db.JSON)
	info.AddField("ItemCounts", "itemCounts", db.JSON)
	info.AddField("ItemProps", "itemProps", db.JSON)
	info.AddField("CreateTime", "createTime", db.Bigint)

	info.SetTable("log_personal_shop").SetTitle("LogPersonalShop").SetDescription("LogPersonalShop")

	formList := logPersonalShop.GetForm()
	formList.AddField("Id", "Id", db.Int, form.Number)
	formList.AddField("PlayerId", "PlayerId", db.Int, form.Number)
	formList.AddField("TradPersonName", "tradPersonName", db.Text, form.RichText)
	formList.AddField("BundleNum", "bundleNum", db.Int, form.Number)
	formList.AddField("SoulGemNum", "soulGemNum", db.Int, form.Number)
	formList.AddField("BlessingGemNum", "blessingGemNum", db.Int, form.Number)
	formList.AddField("DiamondsNum", "DiamondsNum", db.Bigint, form.Number)
	formList.AddField("Price", "price", db.Bigint, form.Number)
	formList.AddField("IsBuy", "isBuy", db.Tinyint, form.Number)
	formList.AddField("ItemQualitys", "itemQualitys", db.JSON, form.Text)
	formList.AddField("ItemLevels", "itemLevels", db.JSON, form.Text)
	formList.AddField("ItemTypeIDs", "itemTypeIDs", db.JSON, form.Text)
	formList.AddField("ItemCounts", "itemCounts", db.JSON, form.Text)
	formList.AddField("ItemProps", "itemProps", db.JSON, form.Text)
	formList.AddField("CreateTime", "createTime", db.Bigint, form.Number)

	formList.SetTable("log_personal_shop").SetTitle("LogPersonalShop").SetDescription("LogPersonalShop")

	return logPersonalShop
}
