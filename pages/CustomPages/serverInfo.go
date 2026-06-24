package CustomPages

import (
	"admin/common"
	"admin/common/def"
	"admin/fusion"
	"admin/mgr/pageMgr"
	"encoding/json"
	"fmt"
	"github.com/GoAdminGroup/go-admin/context"
	"github.com/GoAdminGroup/go-admin/modules/config"
	"github.com/GoAdminGroup/go-admin/modules/db"
	"github.com/GoAdminGroup/go-admin/plugins/admin/modules/constant"
	form2 "github.com/GoAdminGroup/go-admin/plugins/admin/modules/form"
	"github.com/GoAdminGroup/go-admin/plugins/admin/modules/parameter"
	"github.com/GoAdminGroup/go-admin/plugins/admin/modules/response"
	"github.com/GoAdminGroup/go-admin/plugins/admin/modules/table"
	template2 "github.com/GoAdminGroup/go-admin/template"
	"github.com/GoAdminGroup/go-admin/template/chartjs"
	"github.com/GoAdminGroup/go-admin/template/icon"
	"github.com/GoAdminGroup/go-admin/template/types"
	"github.com/GoAdminGroup/go-admin/template/types/form"
	"github.com/GoAdminGroup/themes/adminlte/components/infobox"
	"github.com/GoAdminGroup/themes/sword/components/card"
	"github.com/leonelquinteros/gotext"
	"html/template"
	"math"
	"net/url"
	"sort"
	"strconv"
	"strings"
	"time"
)

func GetServerDetail(ctx *context.Context) (types.Panel, error) {
	components := template2.Get(config.GetTheme())
	serverID := ctx.Request.FormValue("Id")
	serverName := ctx.Request.FormValue("logicName")
	logicOpenStatus := ctx.Request.FormValue("logicOpenStatus")
	ctx1 := &context.Context{}
	vals := url.Values{}
	vals.Add("gsId", serverID)
	err, res := fusion.CallToCenter(&vals, common.MyCfg.Api["getGSOnlineNumber"], "GET", ctx1)
	if err != nil {
		return types.Panel{}, err
	}
	var gsOnlineNum int
	gsOnlineNum, err = strconv.Atoi(res)
	if err != nil {
		gsOnlineNum = 0
	} else {
		res = "开启中"
	}
	gsidx, _ := strconv.Atoi(serverID)
	col1 := components.Col().SetContent(template.HTML("ID:" + serverID)).GetContent()
	col2 := components.Col().SetContent(template.HTML(gotext.Get("服务器标签") + ":" + logicOpenStatus)).GetContent()
	col4 := components.Col().SetContent(template.HTML(gotext.Get("服务器状态") + ":" + res)).GetContent()
	col5 := components.Col().SetContent(template.HTML(gotext.Get("总注册账号数") + ":" + strconv.
		Itoa(int(pageMgr.GetTotalRegisterNum())))).GetContent()
	row := components.Row().SetContent(col1 + col2 + col4 + col5).GetContent()
	countPrice, webCountPrice, count, webCount := pageMgr.GetTotalPayPriceNum(gsidx)
	col6 := components.Col().SetContent(template.HTML(gotext.Get("总充值游戏内金额数") + ":" + strconv.Itoa(int(countPrice)))).GetContent()
	col7 := components.Col().SetContent(template.HTML(gotext.Get("总网页充值金额") + ":" + strconv.Itoa(int(webCountPrice)))).GetContent()
	col8 := components.Col().SetContent(template.HTML(gotext.Get("总游戏内充值人数") + ":" + strconv.Itoa(int(count)))).GetContent()
	col9 := components.Col().SetContent(template.HTML(gotext.Get("总网页充值人数") + ":" + strconv.Itoa(int(webCount)))).GetContent()
	row2 := components.Row().SetContent(col6 + col7 + col8 + col9).GetContent()
	TotalPayNum, TotalPrice, gear, num := pageMgr.GetOrderCount("", 1, gsidx)
	col27 := components.Col().SetContent(template.HTML(gotext.Get("昨日总充值人数") + ":" + strconv.Itoa(int(TotalPayNum)))).GetContent()
	col28 := components.Col().SetContent(template.HTML(gotext.Get("昨日总充值金额") + ":" + strconv.Itoa(int(TotalPrice)))).GetContent()
	row13 := components.Row().SetContent(col27 + col28).GetContent()

	successjs := "var charts = $(\"div.chart\");\n" +
		"$(charts[%d]).children().remove();\n " +
		"$(charts[%d]).append($(\"<canvas />\", {id: \"%s\", style: \"height: 300px\"}));\n" +
		"var temp = new Chart(document.getElementById('%s')," +
		" {\"type\":\"%s\"," +
		"\"data\"" +
		":{\"labels\":data.data.Labels," +
		"\"datasets\":[{\"label\":\"%s\",\"data\":data.data.DsData," +
		"\"type\":\"%s\",\"borderColor\":\"rgb(210, 214, 222)\",\"fill\":false,\"lineTension\":0.1}]}});\n" +
		"console.log(data.data.DsData)"

	btn1 := components.Button().SetType("submit").
		SetContent(template.HTML(gotext.Get("提交"))).
		SetThemePrimary().
		AddClass("submit").
		GetContent()
	info := types.NewFormPanel()
	info.AddField(gotext.Get("ChartID"), "ChartID", db.Int, form.Text).FieldHide()
	info.AddField(gotext.Get("服务器ID"), "ServerID", db.Int, form.Text).FieldValue(serverID).FieldHide()
	info.AddField(gotext.Get("平台"), "PlatFormID", db.Date, form.SelectSingle).FieldWidth(360).
		FieldOptions(types.FieldOptions{
			{Text: gotext.Get("安卓"), Value: "1"},
			{Text: gotext.Get("网页"), Value: "2"},
			{Text: gotext.Get("IOS"), Value: "3"},
		}).FieldValue("1")
	info.AddField(gotext.Get("日期"), "LogTime", db.Date, form.Date).FieldWidth(360).FieldFoot(btn1)

	chart8 := chartjs.Bar()
	chart8.SetTitle(template.HTML(gotext.Get("充值分布"))).
		SetID("OrderDistribution").
		SetHeight(180).
		SetLabels(gear).
		AddDataSet(gotext.Get("人次")).
		DSData(num).
		DSBorderColor("rgb(210, 214, 222)")
	col29 := components.Col().SetContent(chart8.GetContent()).SetSize(types.SizeMD(8)).GetContent()

	_, _, chart1data := pageMgr.GetOnlineNumData("", gsidx)
	col3 := components.Col().SetContent(template.HTML(gotext.Get("当前在线人数") + ":" + strconv.Itoa(gsOnlineNum))).GetContent()
	//col10 := components.Col().SetContent(template.HTML(gotext.Get("昨日最大在线人数") + ":" + strconv.Itoa(max))).GetContent()
	//col11 := components.Col().SetContent(template.HTML(gotext.Get("昨日最小在线人数") + ":" + strconv.Itoa(min))).GetContent()
	chart1 := chartjs.Line()
	chart1.SetTitle(template.HTML(gotext.Get("整点在线人数"))).
		SetID("onlineNum").SetHeight(300).
		SetLabels([]string{"00", "01", "02", "03", "04", "05", "06", "07", "08", "09", "10", "11", "12", "13", "14", "15", "16", "17", "18", "19", "20", "21", "22", "23"}).
		AddDataSet(gotext.Get("人数")).
		DSData(chart1data).
		DSFill(false).
		DSBorderColor("rgb(210, 214, 222)"). // 线边框颜色
		DSLineTension(0.1)                   // 设置压力度
	col12 := components.Col().SetContent(chart1.GetContent()).SetSize(types.SizeMD(8)).GetContent()
	row4 := components.Row().SetContent(col3 /*+ col10 + col11*/).GetContent()

	_, fieldsLevel, levelWastageData := pageMgr.GetLevelWastageNumData("", gsidx)
	_, fieldsQuest, questWastageData := pageMgr.GetQuestWastageNumData("", gsidx)

	//col13 := components.Col().SetContent(template.HTML(gotext.Get("总流失人数") + ":" + strconv.Itoa(total))).GetContent()

	chart2 := chartjs.Bar()
	chart2.SetTitle(template.HTML(gotext.Get("等级疑似流失"))).
		SetID("levelWastage").
		SetHeight(300).
		SetLabels(fieldsLevel).
		AddDataSet(gotext.Get("人数")).
		DSData(levelWastageData).
		DSBorderColor("rgb(210, 214, 222)")

	chart3 := chartjs.Bar()
	chart3.SetTitle(template.HTML(gotext.Get("任务疑似流失"))).
		SetID("QuestWastage").
		SetHeight(300).
		SetLabels(fieldsQuest).
		AddDataSet(gotext.Get("人数")).
		DSData(questWastageData).
		DSBorderColor("rgb(210, 214, 222)")
	col14 := components.Col().SetContent(chart2.GetContent()).SetSize(types.SizeMD(8)).GetContent()
	col17 := components.Col().SetContent(chart3.GetContent()).SetSize(types.SizeMD(8)).GetContent()
	//row9 := components.Row().SetContent(col13).GetContent()

	joinQuestPer, GuildCount, OnlineDurationCount, fieldsPlayer, fieldsFamily, playerLevelData, familyNumData := pageMgr.GetPlayerDataCount("", gsidx)

	col18 := components.Col().SetContent(template.HTML(gotext.Get("战盟人数统计") + ":" + strconv.Itoa(int(GuildCount)))).GetContent()
	col19 := components.Col().SetContent(template.HTML(gotext.Get("家族任务完成率") + ":" + strconv.Itoa(int(joinQuestPer)) + "%")).GetContent()
	col20 := components.Col().SetContent(template.HTML(gotext.Get("玩家在线总时长") + ":" + strconv.Itoa(int(OnlineDurationCount)))).GetContent()
	row11 := components.Row().SetContent(col18 + col19 + col20).GetContent()

	chart4 := chartjs.Bar()
	chart4.SetTitle(template.HTML(gotext.Get("玩家等级分布"))).
		SetID("LevelDistribution").
		SetHeight(300).
		SetLabels(fieldsPlayer).
		AddDataSet(gotext.Get("人数")).
		DSData(playerLevelData).
		DSBorderColor("rgb(210, 214, 222)")
	col15 := components.Col().SetContent(chart4.GetContent()).SetSize(types.SizeMD(8)).GetContent()

	chart5 := chartjs.Bar()
	chart5.SetTitle(template.HTML(gotext.Get("玩家家族分布"))).
		SetID("FamilyDistribution").
		SetHeight(300).
		SetLabels(fieldsFamily).
		AddDataSet("人数").
		DSData(familyNumData).
		DSBorderColor("rgb(210, 214, 222)")
	col16 := components.Col().SetSize(types.SizeMD(2)).SetContent(chart5.GetContent()).GetContent()

	field, data := pageMgr.GetShopItemBuyCount("", gsidx)

	chart6 := chartjs.Bar()
	chart6.SetTitle(template.HTML(gotext.Get("商店道具购买分布"))).
		SetID("ShopItemBuyDistribution").
		SetHeight(300).
		SetLabels(field).
		AddDataSet(gotext.Get("次数")).
		DSData(data).
		DSBorderColor("rgb(210, 214, 222)")
	col22 := components.Col().SetContent(chart6.GetContent()).SetSize(types.SizeMD(8)).GetContent()

	joinNum, clearNum, fieldsDungeonsCunt, joinData, clearData := pageMgr.GetDungeonsCount("", gsidx)

	col23 := components.Col().SetContent(template.HTML(gotext.Get("副本参与人数") + ":" + strconv.Itoa(int(joinNum)))).GetContent()
	col24 := components.Col().SetContent(template.HTML(gotext.Get("副本完成人数") + ":" + strconv.Itoa(int(clearNum)))).GetContent()
	temp := "0"
	if joinNum != 0 {
		temp = strconv.FormatFloat(float64(clearNum)/float64(joinNum)*100, 'f', 2, 32)
	}
	col25 := components.Col().SetContent(template.HTML(gotext.Get("副本完成度") + ":" + temp + "%")).GetContent()
	row12 := components.Row().SetContent(col23 + col24 + col25).GetContent()

	chart7 := chartjs.Bar()
	chart7.SetTitle(template.HTML(gotext.Get("副本参与记录"))).
		SetID("DungeonsCount").
		SetHeight(300).
		SetLabels(fieldsDungeonsCunt).
		AddDataSet(gotext.Get("参与人数")).
		DSData(joinData).
		DSBorderColor("rgb(210, 214, 222)").
		AddDataSet(gotext.Get("通关人数")).
		DSData(clearData).
		DSBorderColor("rgb(54, 162, 235)")
	col26 := components.Col().SetContent(chart7.GetContent()).SetSize(types.SizeMD(8)).GetContent()

	boxInternalRow := components.Row().SetContent(col12).GetContent()
	boxInternalRow2 := components.Row().SetContent(col14).GetContent()
	boxInternalRow3 := components.Row().SetContent(col17).GetContent()
	boxInternalRow4 := components.Row().SetContent(col15).GetContent()
	boxInternalRow5 := components.Row().SetContent(col16).GetContent()
	boxInternalRow6 := components.Row().SetContent(col22).GetContent()
	boxInternalRow7 := components.Row().SetContent(col26).GetContent()
	boxInternalRow8 := components.Row().SetContent(col29).GetContent()

	cardboard1 := card.New().
		SetTitle(serverName + "(S" + serverID + ")").
		SetContent(row)

	info.FieldList[0].Value = "0"
	OrderDistribution := fmt.Sprintf(successjs, 0, 0, "OrderDistribution", "OrderDistribution", "bar", gotext.Get("人次"), "bar")
	aform := components.Form().
		SetHeader(info.Header).
		SetContent(info.FieldList).
		SetPrefix(config.PrefixFixSlash()).
		SetUrl("/admin/ChartData").
		SetId("OrderDistributionForm").
		SetMethod("POST").
		SetAjax(template.JS(OrderDistribution), "alert(\"Failed\")").
		SetHiddenFields(map[string]string{
			form2.NoAnimationKey: "true",
		}).
		GetContent()
	box6 := components.Box().SetTheme("danger").WithHeadBorder().SetHeader(template.HTML(gotext.Get("充值"))).
		SetBody(row2 + row13 + aform + boxInternalRow8).
		GetContent()

	info.FieldList[0].Value = "1"
	info.FieldList[2].Hide = true
	onlineNum := fmt.Sprintf(successjs, 1, 1, "onlineNum", "onlineNum", "line", gotext.Get("人数"), "line")
	aform1 := components.Form().
		SetHeader(info.Header).
		SetContent(info.FieldList).
		SetPrefix(config.PrefixFixSlash()).
		SetUrl("/admin/ChartData").
		SetId("onlineNumForm").
		SetMethod("POST").
		SetAjax(template.JS(onlineNum), "alert(\"Failed\")").
		SetHiddenFields(map[string]string{
			form2.NoAnimationKey: "true",
		}).
		GetContent()
	box1 := components.Box().SetTheme("danger").WithHeadBorder().SetHeader(template.HTML(gotext.Get("在线人数统计"))).
		SetBody(row4 + aform1 + boxInternalRow).
		GetContent()

	info.FieldList[0].Value = "2"
	info.FieldList[2].Hide = true
	levelWastage := fmt.Sprintf(successjs, 2, 2, "levelWastage", "levelWastage", "bar", gotext.Get("人数"), "bar")
	aform2 := components.Form().
		SetHeader(info.Header).
		SetContent(info.FieldList).
		SetPrefix(config.PrefixFixSlash()).
		SetUrl("/admin/ChartData").
		SetId("levelWastageForm").
		SetMethod("POST").
		SetAjax(template.JS(levelWastage), "alert(\"Failed\")").
		SetHiddenFields(map[string]string{
			form2.NoAnimationKey: "true",
		}).
		GetContent()

	info.FieldList[0].Value = "3"
	info.FieldList[2].Hide = true
	questWastage := fmt.Sprintf(successjs, 3, 3, "questWastage", "questWastage", "bar", gotext.Get("人数"), "bar")
	aform3 := components.Form().
		SetHeader(info.Header).
		SetContent(info.FieldList).
		SetPrefix(config.PrefixFixSlash()).
		SetUrl("/admin/ChartData").
		SetId("questWastageForm").
		SetMethod("POST").
		SetAjax(template.JS(questWastage), "alert(\"Failed\")").
		SetHiddenFields(map[string]string{
			form2.NoAnimationKey: "true",
		}).
		GetContent()
	box2 := components.Box().SetTheme("danger").WithHeadBorder().SetHeader(template.HTML(gotext.Get("流失统计"))).
		SetBody( /*row9+*/ aform2 + boxInternalRow2 + aform3 + boxInternalRow3).
		GetContent()

	info.FieldList[0].Value = "4"
	info.FieldList[2].Hide = true
	LevelDistribution := fmt.Sprintf(successjs, 4, 4, "LevelDistribution", "LevelDistribution", "bar", gotext.Get("人数"), "bar")
	aform4 := components.Form().
		SetHeader(info.Header).
		SetContent(info.FieldList).
		SetPrefix(config.PrefixFixSlash()).
		SetUrl("/admin/ChartData").
		SetId("LevelDistributionForm").
		SetMethod("POST").
		SetAjax(template.JS(LevelDistribution), "alert(\"Failed\")").
		SetHiddenFields(map[string]string{
			form2.NoAnimationKey: "true",
		}).
		GetContent()
	info.FieldList[0].Value = "5"
	info.FieldList[2].Hide = true
	FamilyDistribution := fmt.Sprintf(successjs, 5, 5, "FamilyDistribution", "FamilyDistribution", "bar", gotext.Get("人数"), "bar")
	aform5 := components.Form().
		SetHeader(info.Header).
		SetContent(info.FieldList).
		SetPrefix(config.PrefixFixSlash()).
		SetUrl("/admin/ChartData").
		SetId("FamilyDistributionForm").
		SetMethod("POST").
		SetAjax(template.JS(FamilyDistribution), "alert(\"Failed\")").
		SetHiddenFields(map[string]string{
			form2.NoAnimationKey: "true",
		}).
		GetContent()
	box3 := components.Box().SetTheme("danger").WithHeadBorder().SetHeader(template.HTML(gotext.Get("玩家信息统计"))).
		SetBody(row11 + aform4 + boxInternalRow4 + aform5 + boxInternalRow5).
		GetContent()

	info.FieldList[0].Value = "6"
	info.FieldList[2].Hide = true
	ShopItemBuyDistribution := fmt.Sprintf(successjs, 6, 6, "ShopItemBuyDistribution", "ShopItemBuyDistribution", "bar", gotext.Get("次数"), "bar")
	aform6 := components.Form().
		SetHeader(info.Header).
		SetContent(info.FieldList).
		SetPrefix(config.PrefixFixSlash()).
		SetUrl("/admin/ChartData").
		SetId("ShopItemBuyDistributionForm").
		SetMethod("POST").
		SetAjax(template.JS(ShopItemBuyDistribution), "alert(\"Failed\")").
		SetHiddenFields(map[string]string{
			form2.NoAnimationKey: "true",
		}).
		GetContent()
	box4 := components.Box().SetTheme("danger").WithHeadBorder().SetHeader(template.HTML(gotext.Get("商城"))).
		SetBody(aform6 + boxInternalRow6).
		GetContent()

	info.FieldList[0].Value = "7"
	info.FieldList[2].Hide = true
	DungeonsCount := fmt.Sprintf("var charts = $(\"div.chart\");\n"+
		"$(charts[%d]).children().remove();\n "+
		"$(charts[%d]).append($(\"<canvas />\", {id: \"%s\", style: \"height: 300px\"}));\n"+
		"new Chart(document.getElementById('%s'), "+
		"{\"type\":\"bar\",\"data\":\n{\"labels\":data.data.Labels,"+
		"\"datasets\":[{\"label\":\"%s\",\"data\":data.data.DsData,\"type\":\"bar\",\"borderColor\":\"rgb(210, 214, 222)\",\"backgroundColor\":\"rgb(210, 214, 222)\"},"+
		"{\"label\":\"%s\",\"data\":data.data.DsData1,\"type\":\"bar\",\"borderColor\":\"rgb(255,127,80)\",\"backgroundColor\":\"rgb(255,127,80)\"}]}});"+
		"console.log(data.data.DsData)", 7, 7, "DungeonsCount", "DungeonsCount", gotext.Get("参与人数"), gotext.Get("通关人数"))
	aform7 := components.Form().
		SetHeader(info.Header).
		SetContent(info.FieldList).
		SetPrefix(config.PrefixFixSlash()).
		SetUrl("/admin/ChartData").
		SetId("DungeonsCountForm").
		SetMethod("POST").
		SetAjax(template.JS(DungeonsCount), "alert(\"Failed\")").
		SetHiddenFields(map[string]string{
			form2.NoAnimationKey: "true",
		}).
		GetContent()
	box5 := components.Box().SetTheme("danger").WithHeadBorder().SetHeader(template.HTML(gotext.Get("副本"))).
		SetBody(row12 + aform7 + boxInternalRow7).
		GetContent()

	return types.Panel{
		Content: cardboard1.GetContent() +
			box6 + box1 + box2 + box3 + box4 + box5,
	}, nil
}

func EditServerDatabaseCfg(ctx *context.Context) (types.Panel, error) {
	components := template2.Get(config.GetTheme())
	col1 := components.Col().GetContent()
	btn1 := components.Button().SetType("submit").
		SetContent(template.HTML(gotext.Get("提交"))).
		SetThemePrimary().
		SetOrientationRight().
		SetLoadingText(icon.Icon("fa-spinner fa-spin", 2) + `Save`).
		GetContent()
	btn2 := components.Button().SetType("reset").
		SetContent(template.HTML(gotext.Get("重置"))).
		SetThemeWarning().
		SetOrientationLeft().
		GetContent()
	col2 := components.Col().SetSize(types.SizeMD(8)).
		SetContent(btn1 + btn2).GetContent()
	var formList = types.NewFormPanel()
	Id := ctx.Request.FormValue("Id")
	LogicName := ctx.Request.FormValue("logicName")

	databaseCfg := def.DatabaseCfg{}
	go_manager := fusion.GetBaseGormDB("go_manager")
	go_manager.Table("databasecfg").Where("ServerID = ?", Id).Scan(&databaseCfg)
	dataBase := def.DataBase{}
	var temp = fusion.String2Bytes(databaseCfg.LogDataBase)
	json.Unmarshal(temp, &dataBase)

	formList.AddField(gotext.Get("服务器ID"), "Id", db.Int, form.Text).FieldValue(Id).FieldDisplayButCanNotEditWhenUpdate()
	formList.AddField(gotext.Get("服务器名"), "LogicName", db.Int, form.Text).FieldValue(LogicName).FieldDisplayButCanNotEditWhenUpdate()

	formList.AddField(gotext.Get("数据库"), "LogDatabase", db.Int, form.Text).FieldValue(gotext.Get("LOG数据库")).FieldDisplayButCanNotEditWhenUpdate().
		FieldHelpMsg(template.HTML(gotext.Get("此处填写后台服务器对应读取数据库配置")))
	formList.AddField("host", "host_log", db.Varchar, form.Ip).FieldMust().FieldValue(dataBase.Host)
	formList.AddField(gotext.Get("端口"), "port_log", db.Varchar, form.Text).FieldMust().FieldValue(dataBase.Port)
	formList.AddField(gotext.Get("用户名"), "user_log", db.Varchar, form.Text).FieldMust().FieldValue(dataBase.User)
	formList.AddField(gotext.Get("密码"), "pwd_log", db.Varchar, form.Text).FieldMust().FieldValue(dataBase.Pwd)
	var database string
	if dataBase.DataBaseName != "" {
		database = dataBase.DataBaseName
	} else {
		database = "mmorpg_log_s" + Id
	}
	formList.AddField(gotext.Get("数据库名"), "dataBaseName_log", db.Varchar, form.Text).FieldMust().FieldValue(database)

	dataBase = def.DataBase{}
	var temp1 = fusion.String2Bytes(databaseCfg.CharDataBase)
	json.Unmarshal(temp1, &dataBase)
	formList.AddField(gotext.Get("数据库"), "CharDatabase", db.Int, form.Text).FieldValue(gotext.Get("CHAR数据库")).FieldDisplayButCanNotEditWhenUpdate().
		FieldHelpMsg(template.HTML(gotext.Get("此处填写后台服务器对应读取数据库配置")))
	formList.AddField("host", "host_char", db.Varchar, form.Ip).FieldMust().FieldValue(dataBase.Host)
	formList.AddField(gotext.Get("端口"), "port_char", db.Varchar, form.Text).FieldMust().FieldValue(dataBase.Port)
	formList.AddField(gotext.Get("用户名"), "user_char", db.Varchar, form.Text).FieldMust().FieldValue(dataBase.User)
	formList.AddField(gotext.Get("密码"), "pwd_char", db.Varchar, form.Text).FieldMust().FieldValue(dataBase.Pwd)
	var database1 string
	if dataBase.DataBaseName != "" {
		database1 = dataBase.DataBaseName
	} else {
		database1 = "mmorpg_char_s" + Id
	}
	formList.AddField(gotext.Get("数据库名"), "dataBaseName_char", db.Varchar, form.Text).FieldMust().FieldValue(database1)
	sn := 2
	if databaseCfg.RedisSn != 0 {
		sn = databaseCfg.RedisSn
	}
	formList.AddField(gotext.Get("Redis序号"), "redisSn", db.Varchar, form.Text).
		FieldMust().FieldValue(strconv.Itoa(sn)).
		FieldHelpMsg(template.HTML(gotext.Get("填写服务器配置文件对应的redis序号,默认为2")))
	formList.AddField(gotext.Get("deployServerHost"), "deployHost", db.Varchar, form.Text).
		FieldMust().FieldValue(databaseCfg.DeployHost).
		FieldHelpMsg(template.HTML(gotext.Get("请填写对应的服务器所部署的deploy服务端的地址")))
	formList.AddField(gotext.Get("deployServerPort"), "deployPort", db.Varchar, form.Text).
		FieldMust().FieldValue(databaseCfg.DeployPort).
		FieldHelpMsg(template.HTML(gotext.Get("请填写对应的服务器所部署的deploy服务端的端口")))

	aform := components.Form().
		SetHeader(formList.Header).
		SetContent(formList.FieldList).
		SetPrefix(config.PrefixFixSlash()).
		SetUrl("/admin/UpdateServerDatabaseCfg").
		SetHiddenFields(map[string]string{
			form2.PreviousKey: "/admin",
		}).
		SetMethod("GET").
		SetOperationFooter(col1 + col2)

	return types.Panel{
		Content: components.Box().
			SetHeader(aform.GetDefaultBoxHeader(true)).
			WithHeadBorder().
			SetBody(aform.GetContent() + formList.FooterHtml).
			GetContent(),
		Title: template.HTML(gotext.Get("编辑后台服务器读取数据库配置")),
	}, nil
}

func UpdateServerDatabaseCfg(ctx *context.Context) (types.Panel, error) {
	components := template2.Get(config.GetTheme())
	databaseCfg := def.DatabaseCfg{}

	Id := ctx.Request.FormValue("Id")
	databaseCfg.ServerId, _ = strconv.ParseInt(Id, 10, 0)
	dataBase := def.DataBase{}
	dataBase.Host = ctx.Request.FormValue("host_log")
	dataBase.Port = ctx.Request.FormValue("port_log")
	dataBase.User = ctx.Request.FormValue("user_log")
	dataBase.Pwd = ctx.Request.FormValue("pwd_log")
	dataBase.DataBaseName = ctx.Request.FormValue("dataBaseName_log")
	temp, _ := json.Marshal(dataBase)
	databaseCfg.LogDataBase = fusion.Bytes2String(temp)

	dataBase = def.DataBase{}
	dataBase.Host = ctx.Request.FormValue("host_char")
	dataBase.Port = ctx.Request.FormValue("port_char")
	dataBase.User = ctx.Request.FormValue("user_char")
	dataBase.Pwd = ctx.Request.FormValue("pwd_char")
	dataBase.DataBaseName = ctx.Request.FormValue("dataBaseName_char")
	temp = []byte{}
	temp, _ = json.Marshal(dataBase)
	databaseCfg.CharDataBase = fusion.Bytes2String(temp)
	databaseCfg.RedisSn, _ = strconv.Atoi(ctx.Request.FormValue("redisSn"))
	databaseCfg.DeployHost = ctx.Request.FormValue("deployHost")
	databaseCfg.DeployPort = ctx.Request.FormValue("deployPort")

	go_manager := fusion.GetBaseGormDB("go_manager")
	id := 0
	go_manager.Table("databasecfg").Select("serverId").Where("serverId = ?", databaseCfg.ServerId).Scan(&id)
	if id != 0 {
		databaseCfg.UpdateTime = def.MyTime(time.Now())
		go_manager.Updates(&databaseCfg)
	} else {
		databaseCfg.CreateTime = def.MyTime(time.Now())
		go_manager.Create(&databaseCfg)
	}
	fusion.UpdateGormDBConn(&databaseCfg)

	//go_backup := fusion.GetBaseGormDB("go_backup")
	//id = 0
	//go_backup.Table("databasecfg").Select("serverId").Where("serverId = ?", databaseCfg.ServerId).Scan(&id)
	//if id != 0 {
	//	databaseCfg.UpdateTime = def.MyTime(time.Now())
	//	go_backup.Updates(&databaseCfg)
	//} else {
	//	databaseCfg.CreateTime = def.MyTime(time.Now())
	//	go_backup.Create(&databaseCfg)
	//}

	var size = map[string]string{"md": "3", "sm": "6", "xs": "12"}
	col1 := components.Col()
	infobox1 := infobox.New().SetText(template.HTML(gotext.Get("操作成功"))).GetContent()
	infoboxCol1 := col1.SetSize(size).SetContent(infobox1).GetContent()
	return types.Panel{
		Content:     infoboxCol1,
		Title:       template.HTML(gotext.Get("服务器管理")),
		Description: template.HTML(gotext.Get("更新服务器数据库配置")),
	}, nil
}

func GetTDataDetail(ctx *context.Context) table.Table {
	tDataDetail := table.NewDefaultTable(table.Config{
		Driver:     "mysql",
		Connection: "db_global",
		CanAdd:     true,
		Editable:   true,
		Deletable:  true,
		Exportable: true,
		PrimaryKey: table.PrimaryKey{
			Type: db.Varchar,
			Name: "DayNo",
		},
	})
	info := tDataDetail.GetInfo()

	info.HideDetailButton().HideRowSelector().HideNewButton().HideEditButton()

	myOP := fusion.GetServerList()
	myStrs := fusion.GetServerNameStrings()
	info.AddField(gotext.Get("服务器"), "ServerId", db.Int).
		FieldFilterable(types.FilterType{FormType: form.SelectSingle}).
		FieldFilterOptions(myOP).
		FieldDisplay(func(model types.FieldModel) interface{} {
			temp, _ := strconv.Atoi(model.Value)
			return myStrs[temp]
		}).
		FieldSortable()
	info.AddField(gotext.Get("日期"), "DayNo", db.Varchar).
		FieldFilterable(types.FilterType{FormType: form.DateRange}).
		FieldSortable().
		FieldWidth(95)
	info.AddField(gotext.Get("新增玩家"), "NewPlayer", db.Int).
		FieldSortable()
	info.AddField(gotext.Get("活跃玩家"), "ActivePlayer", db.Int).
		FieldSortable()
	info.AddField(gotext.Get("充值玩家"), "RechargeNum", db.Int).
		FieldSortable()
	info.AddField(gotext.Get("充值金额(元)"), "RechargePrice", db.Int).
		FieldSortable()
	info.AddField("ARPPU", "ARPPU", db.Float).
		FieldSortable().
		FieldDisplay(func(model types.FieldModel) interface{} {
			temp, err := strconv.ParseFloat(model.Value, 64)
			if err != nil {
				return "0"
			}
			return strconv.FormatFloat(temp, 'f', 2, 64)
		})
	info.AddField("ARPU", "ARPU", db.Float).
		FieldSortable().
		FieldDisplay(func(model types.FieldModel) interface{} {
			temp, err := strconv.ParseFloat(model.Value, 64)
			if err != nil {
				return "0"
			}
			return strconv.FormatFloat(temp, 'f', 2, 64)
		})
	info.AddField(gotext.Get("充值渗透率(%)"), "Permeability", db.Float).
		FieldSortable().
		FieldDisplay(func(model types.FieldModel) interface{} {
			temp, err := strconv.ParseFloat(model.Value, 64)
			if err != nil {
				return "0"
			}
			return strconv.FormatFloat(temp, 'f', 2, 64)
		})
	info.AddField(gotext.Get("次日留存"), "Retained2", db.Int).
		FieldSortable()
	info.AddField(gotext.Get("三日留存"), "Retained3", db.Int).
		FieldSortable()
	info.AddField(gotext.Get("七日留存"), "Retained7", db.Int).
		FieldSortable()
	info.AddField(gotext.Get("次日留存(%)"), "RetainedPer2", db.Float).
		FieldSortable().
		FieldDisplay(func(model types.FieldModel) interface{} {
			temp, err := strconv.ParseFloat(model.Value, 64)
			if err != nil {
				return "0"
			}
			return strconv.FormatFloat(temp, 'f', 2, 64)
		})
	info.AddField(gotext.Get("三日留存(%)"), "RetainedPer3", db.Float).
		FieldSortable().
		FieldDisplay(func(model types.FieldModel) interface{} {
			temp, err := strconv.ParseFloat(model.Value, 64)
			if err != nil {
				return "0"
			}
			return strconv.FormatFloat(temp, 'f', 2, 64)
		})
	info.AddField(gotext.Get("七日留存(%)"), "RetainedPer7", db.Float).
		FieldSortable().
		FieldDisplay(func(model types.FieldModel) interface{} {
			temp, err := strconv.ParseFloat(model.Value, 64)
			if err != nil {
				return "0"
			}
			return strconv.FormatFloat(temp, 'f', 2, 64)
		})
	info.SetTitle(gotext.Get("数据明细")).HideDeleteButton().HideDetailButton()

	info.SetGetDataFn(
		func(param parameter.Parameters) (data []map[string]interface{}, size int) {
			serverId := ctx.Request.FormValue("ServerId")
			DayNoStart := ctx.Request.FormValue("DayNo_start__goadmin")
			DayNoEnd := ctx.Request.FormValue("DayNo_end__goadmin")
			startTime, _ := time.ParseInLocation("2006-01-02", DayNoStart, time.Local)
			endTime, _ := time.ParseInLocation("2006-01-02", DayNoEnd, time.Local)
			tempTime := endTime.AddDate(0, 0, 1)
			tempdate := tempTime.Format("2006-01-02")

			if serverId == "" || DayNoStart == "" || DayNoEnd == "" {
				return nil, 0
			}
			go_manager := fusion.GetBaseGormDB("go_manager")

			countOnlineData := make([]def.CountPlayerData, 0)
			go_manager.Table("player_data_count").
				Select("*").
				Where("logTime >= ?", DayNoStart).
				Where("logTime <= ?", tempdate).
				Where("gsId = ?", serverId).
				Scan(&countOnlineData)
			map1 := make(map[string]def.CountPlayerData)
			for _, v := range countOnlineData {
				dayNo := time.Time(v.LogTime).Format("20060102")
				map1[dayNo] = v
			}

			orders := make([]def.OrderRecord, 0)
			go_manager.Table("order_record").
				Select("*").
				Where("logTime >= ?", DayNoStart).
				Where("logTime <= ?", tempdate).
				Where("gsId = ?", serverId).
				Scan(&orders)
			map2 := make(map[string]def.OrderRecord)

			for _, v := range orders {
				dayNo := time.Time(v.LogTime).Format("20060102")
				map2[dayNo] = v
			}
			DayNoStart = strings.Replace(DayNoStart, "-", "", -1)
			DayNoEnd = strings.Replace(DayNoEnd, "-", "", -1)

			retainData := make([]def.PlayerRetainedAccount, 0)
			go_manager.Table("player_retained_account").
				Where("dayNo >= ?", DayNoStart).
				Where("dayNo <= ?", DayNoEnd).
				Where("gsId = ?", serverId).
				Scan(&retainData)
			map3 := make(map[string]def.PlayerRetainedAccount)
			for _, v := range retainData {
				map3[v.DayNo] = v
			}

			for startTime.Before(endTime) {
				tempMap := make(map[string]interface{})
				dayNo := startTime.Format("20060102")
				tempMap["DayNo"] = dayNo
				tempMap["ServerId"] = serverId
				retain, err3 := map3[dayNo]
				if !err3 {
					tempMap["NewPlayer"] = 0
					tempMap["Retained2"] = 0
					tempMap["Retained3"] = 0
					tempMap["Retained7"] = 0
					tempMap["RetainedPer2"] = 0.00
					tempMap["RetainedPer3"] = 0.00
					tempMap["RetainedPer7"] = 0.00
				} else {
					tempMap["NewPlayer"] = retain.New
					tempMap["Retained2"] = retain.Retained2
					tempMap["Retained3"] = retain.Retained3
					tempMap["Retained7"] = retain.Retained7
					tempMap["RetainedPer2"] = retain.RetainedPer2
					tempMap["RetainedPer3"] = retain.RetainedPer3
					tempMap["RetainedPer7"] = retain.RetainedPer7
				}
				tempTime1 := startTime.AddDate(0, 0, 1)
				orderTime := tempTime1.Format("20060102")
				count, err1 := map1[orderTime]
				if !err1 {
					tempMap["ActivePlayer"] = 0
				} else {
					tempMap["ActivePlayer"] = count.OnlinePlayerCount
				}
				order, err2 := map2[orderTime]
				if !err2 {
					tempMap["RechargeNum"] = 0
					tempMap["RechargePrice"] = 0
					tempMap["ARPPU"] = 0.00
				} else {
					tempMap["RechargeNum"] = order.TotalPayNum
					tempMap["RechargePrice"] = order.TotalPrice + order.TotalPriceWeb*1000
					if order.TotalPayNum == 0 {
						tempMap["ARPPU"] = 0.00
					} else {
						tempMap["ARPPU"] = float64(order.TotalPrice+order.TotalPriceWeb*1000) / math.Max(float64(order.TotalRechargeAcct), float64(1))
					}
				}
				if !err1 || !err2 {
					tempMap["ARPU"] = 0.00
					tempMap["Permeability"] = 0.00
				} else {
					if count.TotalAcctNum == 0 {
						tempMap["ARPU"] = 0.00
					} else {
						tempMap["ARPU"] = float64(float64(order.TotalPrice+order.TotalPriceWeb*1000) / math.Max(float64(count.TotalAcctNum), float64(1)))
					}
					if count.OnlinePlayerCount == 0 {
						tempMap["Permeability"] = 0.00
					} else {
						tempMap["Permeability"] = (float64(order.TotalPayNum) / math.Max(float64(count.OnlinePlayerCount), float64(1))) * 100
					}
				}
				data = append(data, tempMap)
				startTime = startTime.AddDate(0, 0, 1)
			}
			if param.IsAll() {
				return data, len(data)
			}
			max := math.Min(float64(len(data)), float64((param.PageInt)*param.PageSizeInt))
			return data[(param.PageInt-1)*param.PageSizeInt : int(max)], len(data)
		})

	return tDataDetail
}

//func GetRetained4Account(ctx *context.Context) table.Table {
//	tDataDetail := table.NewDefaultTable(table.Config{
//		Driver:     "mysql",
//		Connection: "db_global",
//		CanAdd:     true,
//		Editable:   true,
//		Deletable:  true,
//		Exportable: true,
//		PrimaryKey: table.PrimaryKey{
//			Type: db.Varchar,
//			Name: "DayNo",
//		},
//	})
//	info := tDataDetail.GetInfo()
//
//	info.HideDetailButton().HideRowSelector().HideNewButton().HideEditButton()
//
//	myOP := fusion.GetServerList()
//	myStrs := fusion.GetServerNameStrings()
//	info.AddField(gotext.Get("服务器"), "ServerId", db.Varchar).
//		FieldFilterable(types.FilterType{FormType: form.SelectSingle}).
//		FieldFilterOptions(myOP).
//		FieldDisplay(func(model types.FieldModel) interface{} {
//			temp, _ := strconv.Atoi(model.Value)
//			return myStrs[temp]
//		}).
//		FieldSortable()
//	info.AddField(gotext.Get("日期"), "DayNo", db.Varchar).
//		FieldFilterable(types.FilterType{FormType: form.DateRange}).
//		FieldSortable().
//		FieldWidth(95)
//	info.AddField(gotext.Get("新增玩家"), "NewPlayer", db.Int).
//		FieldSortable()
//	info.AddField(gotext.Get("次日留存"), "Retained2", db.Int).
//		FieldSortable()
//	info.AddField(gotext.Get("三日留存"), "Retained3", db.Int).
//		FieldSortable()
//	info.AddField(gotext.Get("七日留存"), "Retained7", db.Int).
//		FieldSortable()
//	info.AddField(gotext.Get("次日留存(%)"), "RetainedPer2", db.Float).
//		FieldSortable().
//		FieldDisplay(func(model types.FieldModel) interface{} {
//			temp, err := strconv.ParseFloat(model.Value, 64)
//			if err != nil {
//				return "0"
//			}
//			return strconv.FormatFloat(temp, 'f', 2, 64)
//		})
//	info.AddField(gotext.Get("三日留存(%)"), "RetainedPer3", db.Float).
//		FieldSortable().
//		FieldDisplay(func(model types.FieldModel) interface{} {
//			temp, err := strconv.ParseFloat(model.Value, 64)
//			if err != nil {
//				return "0"
//			}
//			return strconv.FormatFloat(temp, 'f', 2, 64)
//		})
//	info.AddField(gotext.Get("七日留存(%)"), "RetainedPer7", db.Float).
//		FieldSortable().
//		FieldDisplay(func(model types.FieldModel) interface{} {
//			temp, err := strconv.ParseFloat(model.Value, 64)
//			if err != nil {
//				return "0"
//			}
//			return strconv.FormatFloat(temp, 'f', 2, 64)
//		})
//	info.SetTitle("数据明细").HideDeleteButton().HideDetailButton()
//
//	info.SetGetDataFn(
//		func(param parameter.Parameters) (data []map[string]interface{}, size int) {
//			serverId := ctx.Request.FormValue("ServerId")
//			DayNoStart := ctx.Request.FormValue("DayNo_start__goadmin")
//			DayNoEnd := ctx.Request.FormValue("DayNo_end__goadmin")
//			startTime, _ := time.ParseInLocation("2006-01-02", DayNoStart, time.Local)
//			endTime, _ := time.ParseInLocation("2006-01-02", DayNoEnd, time.Local)
//
//			if serverId == "" || DayNoStart == "" || DayNoEnd == "" {
//				return nil, 0
//			}
//			go_manager := fusion.GetBaseGormDB("go_manager")
//
//			DayNoStart = strings.Replace(DayNoStart, "-", "", -1)
//			DayNoEnd = strings.Replace(DayNoEnd, "-", "", -1)
//
//			retainData := make([]def.PlayerRetainedCharacters, 0)
//			go_manager.Table("player_retained_account").
//				Where("dayNo >= ?", DayNoStart).
//				Where("dayNo <= ?", DayNoEnd).
//				Where("gsId = ?", serverId).
//				Scan(&retainData)
//			map3 := make(map[string]def.PlayerRetainedCharacters)
//			for _, v := range retainData {
//				map3[v.DayNo] = v
//			}
//
//			for startTime.Before(endTime) {
//				tempMap := make(map[string]interface{})
//				dayNo := startTime.Format("20060102")
//				tempMap["DayNo"] = dayNo
//				tempMap["ServerId"] = serverId
//				retain, err3 := map3[dayNo]
//				if !err3 {
//					tempMap["NewPlayer"] = 0
//					tempMap["Retained2"] = 0
//					tempMap["Retained3"] = 0
//					tempMap["Retained7"] = 0
//					tempMap["RetainedPer2"] = 0.0
//					tempMap["RetainedPer3"] = 0.0
//					tempMap["RetainedPer7"] = 0.0
//				} else {
//					tempMap["NewPlayer"] = retain.New
//					tempMap["Retained2"] = retain.Retained2
//					tempMap["Retained3"] = retain.Retained3
//					tempMap["Retained7"] = retain.Retained7
//					tempMap["RetainedPer2"] = retain.RetainedPer2
//					tempMap["RetainedPer3"] = retain.RetainedPer3
//					tempMap["RetainedPer7"] = retain.RetainedPer7
//				}
//				data = append(data, tempMap)
//				startTime = startTime.AddDate(0, 0, 1)
//			}
//			if param.IsAll() {
//				return data, len(data)
//			}
//			max := math.Min(float64(len(data)), float64((param.PageInt)*param.PageSizeInt))
//			return data[(param.PageInt-1)*param.PageSizeInt : int(max)], len(data)
//		})
//	return tDataDetail
//}
//
//func GetRetained4IP(ctx *context.Context) table.Table {
//	tDataDetail := table.NewDefaultTable(table.Config{
//		Driver:     "mysql",
//		Connection: "db_global",
//		CanAdd:     true,
//		Editable:   true,
//		Deletable:  true,
//		Exportable: true,
//		PrimaryKey: table.PrimaryKey{
//			Type: db.Varchar,
//			Name: "DayNo",
//		},
//	})
//	info := tDataDetail.GetInfo()
//
//	info.HideDetailButton().HideRowSelector().HideNewButton().HideEditButton()
//
//	myOP := fusion.GetServerList()
//	myStrs := fusion.GetServerNameStrings()
//	info.AddField(gotext.Get("服务器"), "ServerId", db.Int).
//		FieldFilterable(types.FilterType{FormType: form.SelectSingle}).
//		FieldFilterOptions(myOP).
//		FieldDisplay(func(model types.FieldModel) interface{} {
//			temp, _ := strconv.Atoi(model.Value)
//			return myStrs[temp]
//		}).
//		FieldSortable()
//	info.AddField(gotext.Get("日期"), "DayNo", db.Varchar).
//		FieldFilterable(types.FilterType{FormType: form.DateRange}).
//		FieldSortable().
//		FieldWidth(95)
//	info.AddField(gotext.Get("新增玩家"), "NewPlayer", db.Int).
//		FieldSortable()
//	info.AddField(gotext.Get("次日留存"), "Retained2", db.Int).
//		FieldSortable()
//	info.AddField(gotext.Get("三日留存"), "Retained3", db.Int).
//		FieldSortable()
//	info.AddField(gotext.Get("七日留存"), "Retained7", db.Int).
//		FieldSortable()
//	info.AddField(gotext.Get("次日留存(%)"), "RetainedPer2", db.Float).
//		FieldSortable().
//		FieldDisplay(func(model types.FieldModel) interface{} {
//			temp, err := strconv.ParseFloat(model.Value, 64)
//			if err != nil {
//				return "0"
//			}
//			return strconv.FormatFloat(temp, 'f', 2, 64)
//		})
//	info.AddField(gotext.Get("三日留存(%)"), "RetainedPer3", db.Float).
//		FieldSortable().
//		FieldDisplay(func(model types.FieldModel) interface{} {
//			temp, err := strconv.ParseFloat(model.Value, 64)
//			if err != nil {
//				return "0"
//			}
//			return strconv.FormatFloat(temp, 'f', 2, 64)
//		})
//	info.AddField(gotext.Get("七日留存(%)"), "RetainedPer7", db.Float).
//		FieldSortable().
//		FieldDisplay(func(model types.FieldModel) interface{} {
//			temp, err := strconv.ParseFloat(model.Value, 64)
//			if err != nil {
//				return "0"
//			}
//			return strconv.FormatFloat(temp, 'f', 2, 64)
//		})
//	info.SetTitle("数据明细").HideDeleteButton().HideDetailButton()
//
//	info.SetGetDataFn(
//		func(param parameter.Parameters) (data []map[string]interface{}, size int) {
//			serverId := ctx.Request.FormValue("ServerId")
//			DayNoStart := ctx.Request.FormValue("DayNo_start__goadmin")
//			DayNoEnd := ctx.Request.FormValue("DayNo_end__goadmin")
//			startTime, _ := time.ParseInLocation("2006-01-02", DayNoStart, time.Local)
//			endTime, _ := time.ParseInLocation("2006-01-02", DayNoEnd, time.Local)
//
//			if serverId == "" || DayNoStart == "" || DayNoEnd == "" {
//				return nil, 0
//			}
//			go_manager := fusion.GetBaseGormDB("go_manager")
//
//			DayNoStart = strings.Replace(DayNoStart, "-", "", -1)
//			DayNoEnd = strings.Replace(DayNoEnd, "-", "", -1)
//
//			retainData := make([]def.PlayerRetainedCharacters, 0)
//			go_manager.Table("player_retained_ip").
//				Where("dayNo >= ?", DayNoStart).
//				Where("dayNo <= ?", DayNoEnd).
//				Where("gsId = ?", serverId).
//				Scan(&retainData)
//			map3 := make(map[string]def.PlayerRetainedCharacters)
//			for _, v := range retainData {
//				map3[v.DayNo] = v
//			}
//
//			for startTime.Before(endTime) {
//				tempMap := make(map[string]interface{})
//				dayNo := startTime.Format("20060102")
//				tempMap["DayNo"] = dayNo
//				tempMap["ServerId"] = serverId
//				retain, err3 := map3[dayNo]
//				if !err3 {
//					tempMap["NewPlayer"] = 0
//					tempMap["Retained2"] = 0
//					tempMap["Retained3"] = 0
//					tempMap["Retained7"] = 0
//					tempMap["RetainedPer2"] = 0.0
//					tempMap["RetainedPer3"] = 0.0
//					tempMap["RetainedPer7"] = 0.0
//				} else {
//					tempMap["NewPlayer"] = retain.New
//					tempMap["Retained2"] = retain.Retained2
//					tempMap["Retained3"] = retain.Retained3
//					tempMap["Retained7"] = retain.Retained7
//					tempMap["RetainedPer2"] = retain.RetainedPer2
//					tempMap["RetainedPer3"] = retain.RetainedPer3
//					tempMap["RetainedPer7"] = retain.RetainedPer7
//				}
//				data = append(data, tempMap)
//				startTime = startTime.AddDate(0, 0, 1)
//			}
//			if param.IsAll() {
//				return data, len(data)
//			}
//			max := math.Min(float64(len(data)), float64((param.PageInt)*param.PageSizeInt))
//			return data[(param.PageInt-1)*param.PageSizeInt : int(max)], len(data)
//		})
//	return tDataDetail
//}

func GetTDataChart(ctx *context.Context) (types.Panel, error) {
	ServerID, _ := strconv.Atoi(ctx.Request.FormValue("ServerId"))
	DayNoStart := ctx.Request.FormValue("DayNo_start__goadmin")
	DayNoEnd := ctx.Request.FormValue("DayNo_end__goadmin")

	startTime, _ := time.ParseInLocation("2006-01-02", DayNoStart, time.Local)
	endTime, _ := time.ParseInLocation("2006-01-02", DayNoEnd, time.Local)

	DayNoStart = strings.Replace(DayNoStart, "-", "", -1)
	DayNoEnd = strings.Replace(DayNoEnd, "-", "", -1)

	components := template2.Get(config.GetTheme())
	btn1 := components.Button().SetType("submit").
		AddClass("submit").
		SetContent(icon.Icon(icon.Search, 2) + template.HTML(gotext.Get("查找"))).
		SetThemePrimary().
		SetSmallSize().
		SetOrientationLeft().
		SetLoadingText(icon.Icon(icon.Spinner, 1) + template.HTML(gotext.Get("查找"))).
		SetID("submit1").
		GetContent()

	var form1 = types.NewFormPanel()
	myOP := fusion.GetServerList()
	//myStr := fusion.GetServerNameStrings()
	form1.AddField(gotext.Get("服务器ID"), "ServerId", db.Int, form.SelectSingle).
		FieldOptions(myOP).
		FieldMust()

	form1.AddField(gotext.Get("日期"), "DayNo", db.Varchar, form.DateRange).
		FieldMust()

	params := parameter.GetParam(ctx.Request.URL, 10)
	params.SortField = ""

	labels := []string{}

	for startTime.Before(endTime) {
		dayNo := startTime.Format("20060102")
		labels = append(labels, dayNo)
		startTime = startTime.AddDate(0, 0, 1)
	}

	data, _ := GetTDataDetail(ctx).GetData(params.WithIsAll(true))

	allInfoItem := map[string][]float64{}
	for _, v := range data.InfoList {
		for key, DayNo := range labels {
			if DayNo == v["DayNo"].Value {
				temp := v["ActivePlayer"].Value
				tempF, _ := strconv.ParseFloat(temp, 64)
				if allInfoItem["ActivePlayer"] == nil {
					allInfoItem["ActivePlayer"] = make([]float64, len(labels))
				}
				allInfoItem["ActivePlayer"][key] = tempF

				temp = v["NewPlayer"].Value
				tempF, _ = strconv.ParseFloat(temp, 64)
				if allInfoItem["NewPlayer"] == nil {
					allInfoItem["NewPlayer"] = make([]float64, len(labels))
				}
				allInfoItem["NewPlayer"][key] = tempF

				temp = v["RechargePrice"].Value
				tempF, _ = strconv.ParseFloat(temp, 64)
				if allInfoItem["RechargePrice"] == nil {
					allInfoItem["RechargePrice"] = make([]float64, len(labels))
				}
				allInfoItem["RechargePrice"][key] = tempF

				temp = v["RetainedPer2"].Value
				tempF, _ = strconv.ParseFloat(temp, 64)
				if allInfoItem["RetainedPer2"] == nil {
					allInfoItem["RetainedPer2"] = make([]float64, len(labels))
				}
				allInfoItem["RetainedPer2"][key] = tempF

				temp = v["RetainedPer3"].Value
				tempF, _ = strconv.ParseFloat(temp, 64)
				if allInfoItem["RetainedPer3"] == nil {
					allInfoItem["RetainedPer3"] = make([]float64, len(labels))
				}
				allInfoItem["RetainedPer3"][key] = tempF

				temp = v["RetainedPer7"].Value
				tempF, _ = strconv.ParseFloat(temp, 64)
				if allInfoItem["RetainedPer7"] == nil {
					allInfoItem["RetainedPer7"] = make([]float64, len(labels))
				}
				allInfoItem["RetainedPer7"][key] = tempF

				temp = v["Retained2"].Value
				tempF, _ = strconv.ParseFloat(temp, 64)
				if allInfoItem["Retained2"] == nil {
					allInfoItem["Retained2"] = make([]float64, len(labels))
				}
				allInfoItem["Retained2"][key] = tempF

				temp = v["Retained3"].Value
				tempF, _ = strconv.ParseFloat(temp, 64)
				if allInfoItem["Retained3"] == nil {
					allInfoItem["Retained3"] = make([]float64, len(labels))
				}
				allInfoItem["Retained3"][key] = tempF

				temp = v["Retained7"].Value
				tempF, _ = strconv.ParseFloat(temp, 64)
				if allInfoItem["Retained7"] == nil {
					allInfoItem["Retained7"] = make([]float64, len(labels))
				}
				allInfoItem["Retained7"][key] = tempF

				temp = v["RechargeNum"].Value
				tempF, _ = strconv.ParseFloat(temp, 64)
				if allInfoItem["RechargeNum"] == nil {
					allInfoItem["RechargeNum"] = make([]float64, len(labels))
				}
				allInfoItem["RechargeNum"][key] = tempF

				temp = v["ARPPU"].Value
				tempF, _ = strconv.ParseFloat(temp, 64)
				if allInfoItem["ARPPU"] == nil {
					allInfoItem["ARPPU"] = make([]float64, len(labels))
				}
				allInfoItem["ARPPU"][key] = tempF

				temp = v["ARPU"].Value
				tempF, _ = strconv.ParseFloat(temp, 64)
				if allInfoItem["ARPU"] == nil {
					allInfoItem["ARPU"] = make([]float64, len(labels))
				}
				allInfoItem["ARPU"][key] = tempF

				temp = v["Permeability"].Value
				tempF, _ = strconv.ParseFloat(temp, 64)
				if allInfoItem["Permeability"] == nil {
					allInfoItem["Permeability"] = make([]float64, len(labels))
				}
				allInfoItem["Permeability"][key] = tempF
			}
		}
	}

	keys := []string{}
	for k, _ := range allInfoItem {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	i := 0
	stringMap := pageMgr.GetDataDetailString()
	mytabs := make([]map[string]template.HTML, len(allInfoItem))
	for _, key := range keys {
		chart := chartjs.Line()
		chart.SetID(key).SetHeight(180).
			SetLabels(labels).
			AddDataSet(stringMap[key]).
			DSData(allInfoItem[key]).
			DSFill(false).
			DSBorderColor("rgb(210, 214, 222)").
			DSLineTension(0.1)

		col := components.Col().SetContent(chart.GetContent()).SetSize(types.SizeMD(8)).GetContent()
		boxInternalRow := components.Row().SetContent(col).GetContent()
		box := components.Box().SetTheme("danger").WithHeadBorder().SetHeader(template.HTML(stringMap[key])).
			SetBody(boxInternalRow).SetIframeStyle(true).
			GetContent()
		mytabs[i] = map[string]template.HTML{"title": template.HTML(stringMap[key]), "content": box}
		i++
	}
	tabs := components.Tabs().SetData(mytabs).GetContent()
	myStr := fusion.GetServerNameStrings()

	col3 := components.Col().SetContent(template.HTML(gotext.Get("服务器ID") + ":" + myStr[ServerID])).GetContent()
	col10 := components.Col().SetContent(template.HTML(gotext.Get("记录时间") + ":" + DayNoStart + "-" + DayNoEnd)).GetContent()
	row4 := components.Row().SetContent(col3 + col10).GetContent()

	isNotIframe := ctx.Query(constant.IframeKey) != "true"
	boxModel := components.Box().
		SetNoPadding().
		WithHeadBorder().
		SetIframeStyle(!isNotIframe).
		SetHeader(template.HTML(gotext.Get("数据明细")))
	if ServerID != 0 {
		boxModel.SetBody(row4 + tabs)
	}

	boxModel = boxModel.SetSecondHeaderClass("filter-area").
		SetSecondHeader(
			components.Form().
				SetHeader(form1.Header).
				SetContent(form1.FieldList).
				SetPrefix(config.PrefixFixSlash()).
				SetUrl("/admin/DataDetail").
				SetId("form1").
				SetMethod("GET").
				SetOperationFooter(btn1).
				SetHiddenFields(map[string]string{
					form2.NoAnimationKey: "true",
				}).
				GetContent())

	return types.Panel{Content: boxModel.GetContent() + form1.FooterHtml}, nil
}

func GetChartData(ctx *context.Context) {
	chartId, _ := strconv.Atoi(ctx.Request.FormValue("ChartID"))
	gsidx, _ := strconv.Atoi(ctx.Request.FormValue("ServerID"))
	platFormID, _ := strconv.Atoi(ctx.Request.FormValue("PlatFormID"))
	date := ctx.Request.FormValue("LogTime")
	temp := map[string]interface{}{}
	switch chartId {
	case 0:
		_, _, gear, num := pageMgr.GetOrderCount(date, platFormID, gsidx)
		temp["DsData"] = num
		temp["Labels"] = gear
		break
	case 1:
		_, _, chart1data := pageMgr.GetOnlineNumData(date, gsidx)
		temp["DsData"] = chart1data
		temp["Labels"] = []string{"00", "01", "02", "03", "04", "05", "06", "07", "08", "09", "10", "11", "12", "13", "14", "15", "16", "17", "18", "19", "20", "21", "22", "23"}
		break
	case 2:
		_, fieldsLevel, levelWastageData := pageMgr.GetLevelWastageNumData(date, gsidx)
		temp["DsData"] = levelWastageData
		temp["Labels"] = fieldsLevel
		break
	case 3:
		_, fieldsQuest, questWastageData := pageMgr.GetQuestWastageNumData(date, gsidx)
		temp["DsData"] = questWastageData
		temp["Labels"] = fieldsQuest
		break
	case 4:
		_, _, _, fieldsPlayer, _, playerLevelData, _ := pageMgr.GetPlayerDataCount(date, gsidx)
		temp["DsData"] = playerLevelData
		temp["Labels"] = fieldsPlayer
		break
	case 5:
		_, _, _, _, fieldsFamily, _, familyNumData := pageMgr.GetPlayerDataCount(date, gsidx)
		temp["DsData"] = familyNumData
		temp["Labels"] = fieldsFamily
		break
	case 6:
		fieldsItemBuy, ItemBuyData := pageMgr.GetShopItemBuyCount(date, gsidx)
		temp["DsData"] = ItemBuyData
		temp["Labels"] = fieldsItemBuy
		break
	case 7:
		_, _, fieldsDungeonsCunt, joinData, clearData := pageMgr.GetDungeonsCount(date, gsidx)
		temp["DsData"] = joinData
		temp["DsData1"] = clearData
		temp["Labels"] = fieldsDungeonsCunt
		break

	}
	response.OkWithData(ctx, temp)
}

func UpdateServerSetting(ctx *context.Context) {
	OPType, _ := strconv.Atoi(ctx.Request.FormValue("OPType"))
	formValue := ctx.Request.Form["Id[]"]
	var LogicIds []int
	for _, v := range formValue {
		id, err := strconv.Atoi(v)
		if err != nil {
			continue
		}
		LogicIds = append(LogicIds, id)
	}
	db_global := fusion.GetBaseGormDB("db_global")
	// 合服修复：将选中服务器的同IP同端口（合服）服务器一并纳入更新范围
	if len(LogicIds) > 0 {
		var relatedIds []int
		db_global.Raw(`SELECT Id FROM t_game_servers WHERE CONCAT(externalIP, ':', externalPort) IN (SELECT CONCAT(externalIP, ':', externalPort) FROM t_game_servers WHERE Id IN ?)`, LogicIds).Scan(&relatedIds)
		LogicIds = relatedIds
	}
	if OPType == 1 {
		logicSpecialFlags := ctx.Request.FormValue("logicSpecialFlags")
		if logicSpecialFlags != "" {
			status, _ := strconv.Atoi(logicSpecialFlags)
			db_global.Table("t_game_servers").
				Where("Id in ?", LogicIds).
				Update("logicSpecialFlags", status)
		}
	} else if OPType == 2 {
		logicOpenStatus := ctx.Request.FormValue("logicOpenStatus")
		if logicOpenStatus != "" {
			symbol, _ := strconv.Atoi(logicOpenStatus)
			db_global.Table("t_game_servers").
				Where("Id in ?", LogicIds).
				Update("logicOpenStatus", symbol)
		}
	}
	response.Ok(ctx)
}

func EditServerDatabaseBackupCfg(ctx *context.Context) (types.Panel, error) {
	components := template2.Get(config.GetTheme())
	col1 := components.Col().GetContent()
	btn1 := components.Button().SetType("submit").
		SetContent(template.HTML(gotext.Get("提交"))).
		SetThemePrimary().
		SetOrientationRight().
		SetLoadingText(icon.Icon("fa-spinner fa-spin", 2) + `Save`).
		GetContent()
	btn2 := components.Button().SetType("reset").
		SetContent(template.HTML(gotext.Get("重置"))).
		SetThemeWarning().
		SetOrientationLeft().
		GetContent()
	col2 := components.Col().SetSize(types.SizeMD(8)).
		SetContent(btn1 + btn2).GetContent()
	var formList = types.NewFormPanel()
	Id := ctx.Request.FormValue("Id")
	LogicName := ctx.Request.FormValue("logicName")

	databaseCfg := def.DatabaseCfg{}
	go_backup := fusion.GetBaseGormDB("go_backup")
	go_backup.Table("databasecfg").Where("ServerID = ?", Id).Scan(&databaseCfg)
	dataBase := def.DataBase{}
	var temp = fusion.String2Bytes(databaseCfg.LogDataBase)
	json.Unmarshal(temp, &dataBase)

	formList.AddField(gotext.Get("服务器ID"), "Id", db.Int, form.Text).FieldValue(Id).FieldDisplayButCanNotEditWhenUpdate()
	formList.AddField(gotext.Get("服务器名"), "LogicName", db.Int, form.Text).FieldValue(LogicName).FieldDisplayButCanNotEditWhenUpdate()

	formList.AddField(gotext.Get("数据库"), "LogDatabase", db.Int, form.Text).FieldValue(gotext.Get("LOG数据库")).FieldDisplayButCanNotEditWhenUpdate().
		FieldHelpMsg(template.HTML(gotext.Get("此处填写备份服务器对应读取数据库配置")))
	formList.AddField("host", "host_log", db.Varchar, form.Ip).FieldMust().FieldValue(dataBase.Host)
	formList.AddField(gotext.Get("端口"), "port_log", db.Varchar, form.Text).FieldMust().FieldValue(dataBase.Port)
	formList.AddField(gotext.Get("用户名"), "user_log", db.Varchar, form.Text).FieldMust().FieldValue(dataBase.User)
	formList.AddField(gotext.Get("密码"), "pwd_log", db.Varchar, form.Text).FieldMust().FieldValue(dataBase.Pwd)
	var database string
	if dataBase.DataBaseName != "" {
		database = dataBase.DataBaseName
	} else {
		database = "mmorpg_log_s" + Id
	}
	formList.AddField(gotext.Get("数据库名"), "dataBaseName_log", db.Varchar, form.Text).FieldMust().FieldValue(database)

	dataBase = def.DataBase{}
	var temp1 = fusion.String2Bytes(databaseCfg.CharDataBase)
	json.Unmarshal(temp1, &dataBase)
	formList.AddField(gotext.Get("数据库"), "CharDatabase", db.Int, form.Text).FieldValue(gotext.Get("CHAR数据库")).FieldDisplayButCanNotEditWhenUpdate().
		FieldHelpMsg(template.HTML(gotext.Get("此处填写备份服务器对应读取数据库配置")))
	formList.AddField("host", "host_char", db.Varchar, form.Ip).FieldMust().FieldValue(dataBase.Host)
	formList.AddField(gotext.Get("端口"), "port_char", db.Varchar, form.Text).FieldMust().FieldValue(dataBase.Port)
	formList.AddField(gotext.Get("用户名"), "user_char", db.Varchar, form.Text).FieldMust().FieldValue(dataBase.User)
	formList.AddField(gotext.Get("密码"), "pwd_char", db.Varchar, form.Text).FieldMust().FieldValue(dataBase.Pwd)
	var database1 string
	if dataBase.DataBaseName != "" {
		database1 = dataBase.DataBaseName
	} else {
		database1 = "mmorpg_char_s" + Id
	}
	formList.AddField(gotext.Get("数据库名"), "dataBaseName_char", db.Varchar, form.Text).FieldMust().FieldValue(database1)

	aform := components.Form().
		SetHeader(formList.Header).
		SetContent(formList.FieldList).
		SetPrefix(config.PrefixFixSlash()).
		SetUrl("/admin/UpdateServerDatabaseBackupCfg").
		SetHiddenFields(map[string]string{
			form2.PreviousKey: "/admin",
		}).
		SetMethod("GET").
		SetOperationFooter(col1 + col2)

	return types.Panel{
		Content: components.Box().
			SetHeader(aform.GetDefaultBoxHeader(true)).
			WithHeadBorder().
			SetBody(aform.GetContent() + formList.FooterHtml).
			GetContent(),
		Title: template.HTML(gotext.Get("编辑备份服务器读取数据库配置")),
	}, nil

}
func UpdateServerDatabaseBackupCfg(ctx *context.Context) (types.Panel, error) {
	components := template2.Get(config.GetTheme())
	databaseCfg := def.DatabaseCfg{}

	Id := ctx.Request.FormValue("Id")
	databaseCfg.ServerId, _ = strconv.ParseInt(Id, 10, 0)
	dataBase := def.DataBase{}
	dataBase.Host = ctx.Request.FormValue("host_log")
	dataBase.Port = ctx.Request.FormValue("port_log")
	dataBase.User = ctx.Request.FormValue("user_log")
	dataBase.Pwd = ctx.Request.FormValue("pwd_log")
	dataBase.DataBaseName = ctx.Request.FormValue("dataBaseName_log")
	temp, _ := json.Marshal(dataBase)
	databaseCfg.LogDataBase = fusion.Bytes2String(temp)

	dataBase = def.DataBase{}
	dataBase.Host = ctx.Request.FormValue("host_char")
	dataBase.Port = ctx.Request.FormValue("port_char")
	dataBase.User = ctx.Request.FormValue("user_char")
	dataBase.Pwd = ctx.Request.FormValue("pwd_char")
	dataBase.DataBaseName = ctx.Request.FormValue("dataBaseName_char")
	temp = []byte{}
	temp, _ = json.Marshal(dataBase)
	databaseCfg.CharDataBase = fusion.Bytes2String(temp)

	go_backup := fusion.GetBaseGormDB("go_backup")
	id := 0
	go_backup.Table("databasecfg").Select("serverId").Where("serverId = ?", databaseCfg.ServerId).Scan(&id)
	if id != 0 {
		databaseCfg.UpdateTime = def.MyTime(time.Now())
		go_backup.Updates(&databaseCfg)
	} else {
		databaseCfg.CreateTime = def.MyTime(time.Now())
		go_backup.Create(&databaseCfg)
	}

	var size = map[string]string{"md": "3", "sm": "6", "xs": "12"}
	col1 := components.Col()
	infobox1 := infobox.New().SetText(template.HTML(gotext.Get("操作成功"))).GetContent()
	infoboxCol1 := col1.SetSize(size).SetContent(infobox1).GetContent()
	return types.Panel{
		Content:     infoboxCol1,
		Title:       template.HTML(gotext.Get("服务器管理")),
		Description: template.HTML(gotext.Get("更新服务器数据库配置")),
	}, nil
}
