package CustomPages

import (
	"admin/common"
	"admin/common/def"
	"admin/common/def/chequeDef"
	"admin/common/def/currencyDef"
	"admin/fusion"
	"admin/mgr/pageMgr"
	"encoding/json"
	"fmt"
	"html/template"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/GoAdminGroup/go-admin/context"
	"github.com/GoAdminGroup/go-admin/modules/auth"
	"github.com/GoAdminGroup/go-admin/modules/config"
	"github.com/GoAdminGroup/go-admin/modules/db"
	"github.com/GoAdminGroup/go-admin/modules/language"
	"github.com/GoAdminGroup/go-admin/modules/logger"
	form2 "github.com/GoAdminGroup/go-admin/plugins/admin/modules/form"
	"github.com/GoAdminGroup/go-admin/plugins/admin/modules/parameter"
	"github.com/GoAdminGroup/go-admin/plugins/admin/modules/table"
	template2 "github.com/GoAdminGroup/go-admin/template"
	"github.com/GoAdminGroup/go-admin/template/icon"
	"github.com/GoAdminGroup/go-admin/template/types"
	"github.com/GoAdminGroup/go-admin/template/types/form"
	"github.com/GoAdminGroup/themes/adminlte/components/infobox"
	"github.com/leonelquinteros/gotext"
)

func GetPlayerDetail(ctx *context.Context) (types.Panel, error) {
	components := template2.Get(config.GetTheme())

	playerInfo := pageMgr.GetPlayerInfo(ctx)

	col3 := components.Col().SetContent(template.HTML(gotext.Get("服务器ID") + ":" + fusion.Strval(playerInfo["IpcServerID"]))).GetContent()
	col4 := components.Col().SetContent(template.HTML(gotext.Get("玩家ID") + ":" + fusion.Strval(playerInfo["IpcInstID"]))).GetContent()
	col5 := components.Col().SetContent(template.HTML(gotext.Get("账号ID") + ":" + fusion.Strval(playerInfo["IpcAcctID"]))).GetContent()
	col6 := components.Col().SetContent(template.HTML(gotext.Get("玩家昵称") + ":" + fusion.Strval(playerInfo["IpcNickName"]))).GetContent()
	col7 := components.Col().SetContent(template.HTML(gotext.Get("玩家等级") + ":" + fusion.Strval(playerInfo["IpcLevel"]))).GetContent()

	displayUnix, _ := strconv.Atoi(fusion.Strval(playerInfo["IpcCreateTime"]))
	displayTime := time.Unix(int64(displayUnix), 0)
	col8 := components.Col().SetContent(template.HTML(gotext.Get("创建时间") + ":" + displayTime.Format("2006-01-02 15:04:05"))).GetContent()

	displayUnix, _ = strconv.Atoi(fusion.Strval(playerInfo["IpcLastLoginTime"]))
	displayTime = time.Unix(int64(displayUnix), 0)
	col9 := components.Col().SetContent(template.HTML(gotext.Get("最后登录时间") + ":" + displayTime.Format("2006-01-02 15:04:05"))).GetContent()

	col13 := components.Col().SetContent(template.HTML(gotext.Get("首次充值时间") + ":" + fusion.Strval(playerInfo["MinPayTime"]))).GetContent()
	col14 := components.Col().SetContent(template.HTML(gotext.Get("最近充值时间") + ":" + fusion.Strval(playerInfo["MaxPayTime"]))).GetContent()

	displayUnix, _ = strconv.Atoi(fusion.Strval(playerInfo["CreateTime"]))
	displayTime = time.Unix(int64(displayUnix), 0)
	col15 := components.Col().SetContent(template.HTML(gotext.Get("账号创建时间") + ":" + displayTime.Format("2006-01-02 15:04:05"))).GetContent()
	col16 := components.Col().SetContent(template.HTML(gotext.Get("最近登录IP") + ":" + fusion.Strval(playerInfo["LastLoginIP"]))).GetContent()
	col17 := components.Col().SetContent(template.HTML(gotext.Get("账号名") + ":" + fusion.Strval(playerInfo["UserName"]))).GetContent()

	cardboard2 := infobox.New().
		SetContent(col3 + col4 + col5 + col17 + col6 + col7 + col8 + col9 + col13 + col14 + col15 + col16).SetText(template.HTML(gotext.Get("玩家信息")))
	row2 := components.Row().SetContent(cardboard2.GetContent()).GetContent()

	var gold, bwing, wing = 0, 0, 0
	statValue := def.StatValue{}
	var temp = fusion.String2Bytes(fusion.Strval(playerInfo["IpcStatValue"]))
	json.Unmarshal(temp, &statValue)
	if len(statValue.BankCurrencies) != 0 {
		gold = statValue.BankCurrencies[currencyDef.Gold-1]
		//bwing = statValue.BankCurrencies[currencyDef.BindDiamond-1]
	}
	temp = fusion.String2Bytes(fusion.Strval(playerInfo["IpcCurrencies"]))
	Currencies := []int{}
	json.Unmarshal(temp, &Currencies)
	if len(Currencies) != 0 {
		gold += Currencies[currencyDef.Gold-1]
		bwing += Currencies[currencyDef.BindDiamond-1]
		wing += Currencies[currencyDef.Diamond-1]
	}
	col10 := components.Col().SetContent(template.HTML(gotext.Get("金币数量") + ":" + strconv.Itoa(gold))).GetContent()
	col11 := components.Col().SetContent(template.HTML(gotext.Get("绑钻数量") + ":" + strconv.Itoa(bwing))).GetContent()
	col12 := components.Col().SetContent(template.HTML(gotext.Get("钻石数量") + ":" + strconv.Itoa(wing))).GetContent()
	cardboard3 := infobox.New().SetContent(col10 + col11 + col12).SetText(template.HTML(gotext.Get("资产统计")))
	row3 := components.Row().SetContent(cardboard3.GetContent()).GetContent()

	itemStr := fusion.Strval(playerInfo["IpcStorageItems"])
	panel := pageMgr.GetBagDetailInfo(ctx, &itemStr)
	params := parameter.GetParamFromURL(ctx.Request.URL.String(), 10, "asc", panel.GetPrimaryKey().Name)
	body := fusion.ShowTable(ctx, params, panel, "/admin/userMgr/PlayerDetail", false)
	data, _ := panel.GetData(params)
	FieldOptions := []types.FieldOption{}

	for _, v := range data.InfoList {
		temp1 := types.FieldOption{
			Text:  string(v["ItemTypeID"].Content),
			Value: v["ItemTypeID"].Value,
		}
		FieldOptions = append(FieldOptions, temp1)
	}

	ServerID := fusion.Strval(playerInfo["IpcServerID"])
	PlayerID := fusion.Strval(playerInfo["IpcInstID"])
	AcctID := fusion.Strval(playerInfo["IpcAcctID"])

	btn1 := components.Button().SetType("submit").
		SetContent(template.HTML(gotext.Get("提交"))).
		SetThemePrimary().
		SetOrientationRight().
		GetContent()
	myops := fusion.GetPurchaseOp()
	tpPointOptWithId := fusion.GetTpPointOptionAndID()
	var form1 = types.NewFormPanel()
	user := auth.Auth(ctx)
	if !user.CheckRole("playerManager") {
		form1.AddField(gotext.Get("角色ID"), "IpcInstID", db.Int, form.Text).FieldDefault(PlayerID).FieldHide()
		form1.AddField(gotext.Get("账号ID"), "AcctID", db.Int, form.Text).FieldDefault(AcctID).FieldHide()
		form1.AddField(gotext.Get("服务器ID"), "IpcServerID", db.Int, form.Text).FieldDefault(ServerID).FieldHide()
		form1.AddField(gotext.Get("操作类型"), "OperationType", db.Int, form.SelectSingle).
			FieldDefault("1").
			FieldOptions(types.FieldOptions{
				{Text: gotext.Get("设置等级"), Value: "1", Selected: true},
				{Text: gotext.Get("设置货币"), Value: "2"},
				{Text: gotext.Get("移除玩家道具"), Value: "3"},
				{Text: gotext.Get("模拟充值"), Value: "6"},
				{Text: gotext.Get("移动角色"), Value: "8"},
				{Text: gotext.Get("禁言"), Value: "4"},
				{Text: gotext.Get("解除禁言"), Value: "9"},
				{Text: gotext.Get("封禁角色"), Value: "5"},
				{Text: gotext.Get("解除角色封禁"), Value: "10"},
				{Text: gotext.Get("封禁账号"), Value: "7"},
				{Text: gotext.Get("解除账号封禁"), Value: "11"},
				{Text: gotext.Get("删除角色"), Value: "12"},
				{Text: gotext.Get("转服"), Value: "13"},
			}).FieldOnChooseHide("1", "CurrencyID", "chequeNum", "ItemTypeID", "ItemCount", "BanChatTime", "BanLoginTime", "RechargeType", "BanLoginTimeAcc", "MovePlayer").
			FieldOnChooseShow("1", "setLevel").
			FieldOnChooseShow("2", "CurrencyID", "chequeNum").
			FieldOnChooseShow("3", "ItemTypeID", "ItemCount").
			FieldOnChooseShow("4", "BanChatTime").
			FieldOnChooseShow("5", "BanLoginTime").
			FieldOnChooseShow("6", "RechargeType").
			FieldOnChooseShow("7", "BanLoginTimeAcc").
			FieldOnChooseShow("8", "MovePlayer").
			FieldOnChooseShow("9", "ZhuanFu")

		form1.AddField(gotext.Get("设置等级"), "setLevel", db.Int, form.Text)
		form1.AddField(gotext.Get("选择货币"), "CurrencyID", db.Int, form.SelectSingle).FieldOptions(currencyDef.GetCurrencyOps())
		form1.AddField(gotext.Get("货币数量"), "chequeNum", db.Int, form.Number)
		form1.AddField(gotext.Get("选择移除道具"), "ItemTypeID", db.Int, form.SelectSingle).FieldOptions(FieldOptions)
		form1.AddField(gotext.Get("道具数量"), "ItemCount", db.Int, form.Number)
		form1.AddField(gotext.Get("禁言时间"), "BanChatTime", db.Int, form.Datetime)
		form1.AddField(gotext.Get("角色封禁时间"), "BanLoginTime", db.Int, form.Datetime)
		form1.AddField(gotext.Get("模拟充值类型"), "RechargeType", db.Int, form.SelectSingle).FieldOptions(myops) //TODO:验证数据是否超过最大值
		form1.AddField(gotext.Get("账号封禁时间"), "BanLoginTimeAcc", db.Int, form.Datetime)
		form1.AddField(gotext.Get("移动角色到"), "MovePlayer", db.Int, form.SelectSingle).FieldOptions(tpPointOptWithId)
		form1.AddField(gotext.Get("转服角色到"), "ZhuanFu", db.Int, form.Text)

		form1.SetTabGroups(types.TabGroups{
			{"AcctID", "IpcInstID", "IpcServerID", "OperationType", "setLevel", "CurrencyID", "chequeNum",
				"ItemTypeID", "ItemCount", "BanChatTime", "BanLoginTime", "RechargeType", "BanLoginTimeAcc", "MovePlayer", "ZhuanFu"},
		})
	} else {
		form1.AddField(gotext.Get("角色ID"), "IpcInstID", db.Int, form.Text).FieldDefault(PlayerID).FieldHide()
		form1.AddField(gotext.Get("账号ID"), "AcctID", db.Int, form.Text).FieldDefault(AcctID).FieldHide()
		form1.AddField(gotext.Get("服务器ID"), "IpcServerID", db.Int, form.Text).FieldDefault(ServerID).FieldHide()
		form1.AddField(gotext.Get("操作类型"), "OperationType", db.Int, form.SelectSingle).
			FieldDefault("4").
			FieldOptions(types.FieldOptions{
				{Text: gotext.Get("禁言"), Value: "4"},
				{Text: gotext.Get("禁言"), Value: "9"},
				{Text: gotext.Get("封禁角色"), Value: "5"},
				{Text: gotext.Get("解除角色封禁"), Value: "10"},
				{Text: gotext.Get("封禁账号"), Value: "7"},
				{Text: gotext.Get("解除账号封禁"), Value: "11"},
			}).
			FieldOnChooseHide("4", "BanChatTime", "BanLoginTime", "BanLoginTimeAcc").
			FieldOnChooseShow("4", "BanChatTime").
			FieldOnChooseShow("5", "BanLoginTime").
			FieldOnChooseShow("7", "BanLoginTimeAcc")

		form1.AddField(gotext.Get("禁言时间"), "BanChatTime", db.Int, form.Datetime)
		form1.AddField(gotext.Get("角色封禁时间"), "BanLoginTime", db.Int, form.Datetime)
		form1.AddField(gotext.Get("账号封禁时间"), "BanLoginTimeAcc", db.Int, form.Datetime)

		form1.SetTabGroups(types.TabGroups{
			{"AcctID", "IpcInstID", "IpcServerID", "OperationType", "BanChatTime", "BanLoginTime", "BanLoginTimeAcc"},
		})
	}
	form1.SetTabHeaders(gotext.Get("操作"))

	field, headers := form1.GroupField()

	dataForm1 := components.Form().
		SetTabHeaders(headers).
		SetTabContents(field).
		SetPrefix(config.PrefixFixSlash()).
		SetUrl("/admin/userMgr/SubmitOperation").
		SetHiddenFields(map[string]string{
			form2.PreviousKey: "/admin",
		}).
		SetMethod("GET").
		SetOperationFooter(btn1).
		GetContent()

	return types.Panel{
		Content: components.Box().
			SetHeader(template.HTML(gotext.Get("玩家详情"))).
			SetBody(row2 + row3 + dataForm1 + form1.FooterHtml + body).
			GetContent(),
		Title: template.HTML(gotext.Get("玩家信息")),
	}, nil
}

func SubmitOperation(ctx *context.Context) (types.Panel, error) {
	components := template2.Get(config.GetTheme())
	setLevel := ctx.Request.FormValue("setLevel")
	currencyID := ctx.Request.FormValue("CurrencyID")
	chequeNum := ctx.Request.FormValue("chequeNum")
	ItemTypeID := ctx.Request.FormValue("ItemTypeID")
	ItemCount := ctx.Request.FormValue("ItemCount")
	AcctID := ctx.Request.FormValue("AcctID")
	PlayerID := ctx.Request.FormValue("IpcInstID")
	ServerID := ctx.Request.FormValue("IpcServerID")
	BanChatTime := ctx.Request.FormValue("BanChatTime")
	BanLoginTime := ctx.Request.FormValue("BanLoginTime")
	BanLoginTimeAcc := ctx.Request.FormValue("BanLoginTimeAcc")
	OperationType := ctx.Request.FormValue("OperationType")
	MovePlayer := ctx.Request.FormValue("MovePlayer")
	ZhuanFu := ctx.Request.FormValue("ZhuanFu")
	rechargeType, _ := strconv.Atoi(ctx.Request.FormValue("RechargeType"))
	operationType, _ := strconv.Atoi(OperationType)

	msg := ""
	vars := url.Values{}
	vars.Add("gsId", ServerID)
	args := PlayerID + ","

	//vars.Add("name", "")//TODO:lobbyServer
	if PlayerID != "" && ServerID != "" {
		switch operationType {
		case 1:
			if setLevel != "" {
				vars.Add("cmd", "ModifyPlayer")
				args += "SetLevel" + "," + setLevel
			} else {
				msg = gotext.Get("请输入正确的值 ")
			}
			break
		case 2:
			if chequeNum != "" && currencyID != "" {
				vars.Add("cmd", "ModifyPlayer")
				id, _ := strconv.Atoi(currencyID)
				chequeID := strconv.Itoa(chequeDef.CurrencyToCheque(id))
				args += "SetCheque" + "," + chequeID + "," + chequeNum
			} else {
				msg = gotext.Get("请输入正确的数量 ")
			}
			break
		case 3:
			if ItemTypeID != "" && ItemCount != "" {
				vars.Add("cmd", "ModifyPlayer")
				args += "RemoveItem" + "," + ItemTypeID + "," + ItemCount
			} else {
				msg = gotext.Get("请输入正确的数量 ")
			}
			break
		case 4:
			if BanChatTime != "" {
				vars.Add("cmd", "BanChatTime")
				temp, _ := time.ParseInLocation("2006-01-02 15:04:05", BanChatTime, time.Local)
				BanChatTimeUnix := int(temp.Unix())
				args += strconv.Itoa(BanChatTimeUnix)
			} else {
				msg = gotext.Get("请检查输入 ")
			}
			break
		case 9:
			vars.Add("cmd", "BanChatTime")
			args += strconv.Itoa(0)
			break
		case 5:
			if BanLoginTime != "" {
				vars.Add("cmd", "BanLoginTime")
				temp, _ := time.ParseInLocation("2006-01-02 15:04:05", BanLoginTime, time.Local)
				BanLoginTimeUnix := int(temp.Unix())
				args += strconv.Itoa(BanLoginTimeUnix)
			} else {
				msg = gotext.Get("请检查输入 ")
			}
			break
		case 10:
			vars.Add("cmd", "BanLoginTime")
			args += strconv.Itoa(0)
			break
		case 7:
			if BanLoginTimeAcc != "" {
				vars.Add("cmd", "KickPlayer")
				temp, _ := time.ParseInLocation("2006-01-02 15:04:05", BanLoginTimeAcc, time.Local)
				if time.Now().Before(temp) {
					db_global := fusion.GetBaseGormDB("db_global")
					db_global.Table("t_accounts").Where("Id = ?", AcctID).Update("banExpireTime", temp.Unix())
				}
			} else {
				msg = gotext.Get("请检查输入 ")
			}
			break
		case 11:
			db_global := fusion.GetBaseGormDB("db_global")
			db_global.Table("t_accounts").Where("Id = ?", AcctID).Update("banExpireTime", 0)
			break
		case 8:
			if MovePlayer != "" {
				args += MovePlayer
				vars.Add("cmd", "MoveTo")
			} else {
				msg = gotext.Get("请检查输入 ")
			}
		case 13:
			if ZhuanFu != "" {
				
			} else {
				msg = gotext.Get("请检查输入 ")
			}
		}

		api := ""
		if operationType == 12 {
			vars = url.Values{}
			vars.Add("playerId", PlayerID)
			api = common.MyCfg.Api["DeleteCharacter"]
		} else if operationType != 6 {
			vars.Add("args", args)
			api = common.MyCfg.Api["GM2GS"]
		} else {
			if rechargeType != 0 && rechargeType != 11 {
				Purchase := fusion.GetPurchase()
				op, ok := Purchase[rechargeType].(def.PurchaseCut)
				if !ok {
					msg = gotext.Get("操作失败") + "Purchase is error "
				} else {
					now := time.Now()
					str := "MN" + strconv.Itoa(int(now.UnixNano())) + "_" + strconv.Itoa(rechargeType)
					vars.Add("orderId", str)
					vars.Add("money", strconv.Itoa(op.BuyPrice))
					vars.Add("payId", strconv.Itoa(rechargeType))
					vars.Add("playerId", PlayerID)
					vars.Add("gsId", ServerID)
					vars.Add("token", fusion.VerifyWebArgs(str, ServerID))
					api = common.MyCfg.Api["webPurchase"]
				}
			} else {
				msg = gotext.Get("请检查输入 ")
			}
		}
		ctx1 := &context.Context{}
		if api == "" {
			msg = gotext.Get("操作失败") + " api = " + api
		}
		err, res := fusion.CallToCenter(&vars, api, "GET", ctx1)
		if err != nil {
			msg = gotext.Get("操作失败")
		}
		if res != "" {
			if res != "All Is OK!" {
				msg = res
			} else {
				msg = gotext.Get("操作成功")
			}
		}
	} else {
		msg = gotext.Get("请输入正确的数量 ")
	}
	infobox1 := infobox.New().SetText(template.HTML(msg))
	col1 := components.Col()
	var size = map[string]string{"md": "3", "sm": "6", "xs": "12"}
	infoboxCol1 := col1.SetSize(size).SetContent(infobox1.GetContent()).GetContent()
	str := gotext.Get("返回上级")
	preUrl := ctx.Request.Header.Get("Referer")
	index := strings.Index(preUrl, "admin")
	preUrl = preUrl[index:]

	return types.Panel{
		Content:     infoboxCol1,
		Title:       template.HTML(gotext.Get("修改结果")),
		Description: template.HTML(fmt.Sprintf(`<a href="/%s">%s</a>`, preUrl, str)),
		Callbacks:   nil,
	}, nil
}

func GetTransformCharacterPage(ctx *context.Context) (types.Panel, error) {
	components := template2.Get(config.GetTheme())
	panel := types.NewFormPanel()
	myop := fusion.GetServerList()
	panel.AddField(gotext.Get("发起账号ID"), "fromAccountID", db.Int, form.Text).FieldMust()
	panel.AddField(gotext.Get("发起服务器ID"), "fromServerID", db.Int, form.SelectSingle).FieldOptions(myop).FieldMust()
	panel.AddField(gotext.Get("转移角色ID"), "characterID", db.Int, form.Text).FieldMust()
	panel.AddField(gotext.Get("接收账号ID"), "toAccountID", db.Int, form.Text).FieldMust()

	col1 := components.Col().GetContent()
	btn1 := components.Button().SetType("submit").
		SetContent(language.GetFromHtml("Save")).
		SetThemePrimary().
		SetOrientationRight().
		SetLoadingText(icon.Icon("fa-spinner fa-spin", 2) + `Save`).
		GetContent()
	btn2 := components.Button().SetType("reset").
		SetContent(language.GetFromHtml("Reset")).
		SetThemeWarning().
		SetOrientationLeft().
		GetContent()
	col2 := components.Col().SetSize(types.SizeMD(8)).
		SetContent(btn1 + btn2).GetContent()
	form1 := components.Form().
		SetHeader(panel.Header).
		SetContent(panel.FieldList).
		SetPrefix(config.PrefixFixSlash()).
		SetUrl("/admin/SubmitTransfer").
		SetHiddenFields(map[string]string{
			form2.PreviousKey: "/admin",
		}).
		SetMethod("GET").
		SetOperationFooter(col1 + col2)

	return types.Panel{
		Content: components.Box().
			SetHeader(form1.GetDefaultBoxHeader(true)).
			WithHeadBorder().
			SetBody(form1.GetContent() + panel.FooterHtml /*+ popup.GetContent()*/).
			GetContent(),
		Title:     template.HTML(gotext.Get("同服跨账号角色转移")),
		Callbacks: nil,
	}, nil
}

func GetSubmitTransferPage(ctx *context.Context) (types.Panel, error) {
	fromAccountID := ctx.FormValue("fromAccountID")
	fromServerID := ctx.FormValue("fromServerID")
	characterID := ctx.FormValue("characterID")
	toAccountID := ctx.FormValue("toAccountID")

	ctx1 := &context.Context{}
	vars := url.Values{}
	vars.Add("accountId", fromAccountID)
	vars.Add("playerId", characterID)
	vars.Add("serverId", fromServerID)
	vars.Add("newAccountId", toAccountID)
	err, res := fusion.CallToCenter(&vars, common.MyCfg.Api["TransferToOtherAccount"], "GET", ctx1)
	if err != nil {
		logger.Info("CountPlayerOnlineNum Failed, Because CenterServer has some error")
	}

	components := template2.Get(config.GetTheme())
	infobox1 := infobox.New().SetText(template.HTML(res))
	str := gotext.Get("返回上级")
	preUrl := ctx.Request.Header.Get("Referer")
	index := strings.Index(preUrl, "admin")
	preUrl = preUrl[index:]
	col1 := components.Col()
	var size = map[string]string{"md": "3", "sm": "6", "xs": "12"}
	infoboxCol1 := col1.SetSize(size).SetContent(infobox1.GetContent()).GetContent()
	return types.Panel{
		Content:     infoboxCol1,
		Title:       template.HTML(gotext.Get("转移结果")),
		Description: template.HTML(fmt.Sprintf(`<a href="/%s">%s</a>`, preUrl, str)),
		Callbacks:   nil,
	}, nil
}

func GetTransformServerPage(ctx *context.Context) (types.Panel, error) {
	myop := fusion.GetServerList()
	components := template2.Get(config.GetTheme())
	panel := types.NewFormPanel()
	panel.AddField(gotext.Get("发起服务器ID"), "fromServerID", db.Int, form.SelectSingle).FieldOptions(myop).FieldMust()
	panel.AddField(gotext.Get("接收服务器ID"), "toServerID", db.Int, form.SelectSingle).FieldOptions(myop).FieldMust()
	panel.AddField(gotext.Get("发起账号ID"), "fromAccountID", db.Int, form.Text).FieldMust()
	panel.AddField(gotext.Get("转移角色ID"), "characterID", db.Int, form.Text).FieldMust()

	col1 := components.Col().GetContent()
	btn1 := components.Button().SetType("submit").
		SetContent(language.GetFromHtml("submit")).
		SetThemePrimary().
		SetOrientationRight().
		SetLoadingText(icon.Icon("fa-spinner fa-spin", 2) + `Save`).
		GetContent()
	btn2 := components.Button().SetType("reset").
		SetContent(language.GetFromHtml("Reset")).
		SetThemeWarning().
		SetOrientationLeft().
		GetContent()
	col2 := components.Col().SetSize(types.SizeMD(8)).
		SetContent(btn1 + btn2).GetContent()
	form1 := components.Form().
		SetHeader(panel.Header).
		SetContent(panel.FieldList).
		SetPrefix(config.PrefixFixSlash()).
		SetUrl("/admin/SubmitTransferServer").
		SetHiddenFields(map[string]string{
			form2.PreviousKey: "/admin",
		}).
		SetMethod("GET").
		SetOperationFooter(col1 + col2)

	return types.Panel{
		Content: components.Box().
			SetHeader(form1.GetDefaultBoxHeader(true)).
			WithHeadBorder().
			SetBody(form1.GetContent() + panel.FooterHtml /*+ popup.GetContent()*/).
			GetContent(),
		Title:     template.HTML(gotext.Get("跨服同账号转移角色")),
		Callbacks: nil,
	}, nil
}

func GetSubmitTransferServerPage(ctx *context.Context) (types.Panel, error) {
	fromServerID := ctx.FormValue("fromServerID")
	toServerID := ctx.FormValue("toServerID")
	fromAccountID := ctx.FormValue("fromAccountID")
	characterID := ctx.FormValue("characterID")

	var resultStr string

	if fromServerID == toServerID {
		resultStr = gotext.Get("请选择不同的服务器")
	} else {
		ctx1 := &context.Context{}
		vars := url.Values{}
		vars.Add("accountId", fromAccountID)
		vars.Add("playerId", characterID)
		vars.Add("serverId", fromServerID)
		vars.Add("newAccountId", fromAccountID)
		vars.Add("newServerId", toServerID)
		err, res := fusion.CallToCenter(&vars, common.MyCfg.Api["TransferToOtherServer"], "GET", ctx1)
		if err != nil {
			logger.Info("CountPlayerOnlineNum Failed, Because CenterServer has some error")
		}
		resultStr = res
	}

	components := template2.Get(config.GetTheme())
	infobox1 := infobox.New().SetText(template.HTML(resultStr))
	str := gotext.Get("返回上级")
	preUrl := ctx.Request.Header.Get("Referer")
	index := strings.Index(preUrl, "admin")
	preUrl = preUrl[index:]
	col1 := components.Col()
	var size = map[string]string{"md": "3", "sm": "6", "xs": "12"}
	infoboxCol1 := col1.SetSize(size).SetContent(infobox1.GetContent()).GetContent()
	return types.Panel{
		Content:     infoboxCol1,
		Title:       template.HTML(gotext.Get("转移结果")),
		Description: template.HTML(fmt.Sprintf(`<a href="/%s">%s</a>`, preUrl, str)),
		Callbacks:   nil,
	}, nil
}

func GetDeletePlayerPage(ctx *context.Context) (types.Panel, error) {
	components := template2.Get(config.GetTheme())
	panel := types.NewFormPanel()
	panel.AddField(gotext.Get("删除角色ID"), "characterID", db.Int, form.Text).FieldMust()

	col1 := components.Col().GetContent()
	btn1 := components.Button().SetType("submit").
		SetContent(language.GetFromHtml("submit")).
		SetThemePrimary().
		SetOrientationRight().
		SetLoadingText(icon.Icon("fa-spinner fa-spin", 2) + `Save`).
		GetContent()
	btn2 := components.Button().SetType("reset").
		SetContent(language.GetFromHtml("Reset")).
		SetThemeWarning().
		SetOrientationLeft().
		GetContent()
	col2 := components.Col().SetSize(types.SizeMD(8)).
		SetContent(btn1 + btn2).GetContent()
	form1 := components.Form().
		SetHeader(panel.Header).
		SetContent(panel.FieldList).
		SetPrefix(config.PrefixFixSlash()).
		SetUrl("/admin/SubmitDeletePlayer").
		SetHiddenFields(map[string]string{
			form2.PreviousKey: "/admin",
		}).
		SetMethod("GET").
		SetOperationFooter(col1 + col2)

	return types.Panel{
		Content: components.Box().
			SetHeader(form1.GetDefaultBoxHeader(true)).
			WithHeadBorder().
			SetBody(form1.GetContent() + panel.FooterHtml /*+ popup.GetContent()*/).
			GetContent(),
		Title:     template.HTML(gotext.Get("删除角色")),
		Callbacks: nil,
	}, nil
}

func GetSubmitDeletePlayerPage(ctx *context.Context) (types.Panel, error) {
	characterID := ctx.FormValue("characterID")

	ctx1 := &context.Context{}
	vars := url.Values{}
	vars.Add("playerId", characterID)
	err, res := fusion.CallToCenter(&vars, common.MyCfg.Api["DeleteCharacter"], "GET", ctx1)
	if err != nil {
		logger.Info("CountPlayerOnlineNum Failed, Because CenterServer has some error")
	}

	components := template2.Get(config.GetTheme())
	infobox1 := infobox.New().SetText(template.HTML(res))
	str := gotext.Get("返回上级")
	preUrl := ctx.Request.Header.Get("Referer")
	index := strings.Index(preUrl, "admin")
	preUrl = preUrl[index:]
	col1 := components.Col()
	var size = map[string]string{"md": "3", "sm": "6", "xs": "12"}
	infoboxCol1 := col1.SetSize(size).SetContent(infobox1.GetContent()).GetContent()
	return types.Panel{
		Content:     infoboxCol1,
		Title:       template.HTML(gotext.Get("操作结果")),
		Description: template.HTML(fmt.Sprintf(`<a href="/%s">%s</a>`, preUrl, str)),
		Callbacks:   nil,
	}, nil
}

func GetRetrievePWDPage(ctx *context.Context) table.Table {
	retrievePWD := table.NewDefaultTable(table.Config{
		Driver:     "mysql",
		Connection: "",
		CanAdd:     true,
		Editable:   true,
		Deletable:  true,
		Exportable: true,
		PrimaryKey: table.PrimaryKey{
			Type: db.Int,
			Name: "PlayerId",
		},
	})

	serverOps := fusion.GetServerList()
	info := retrievePWD.GetInfo()
	info.AddField(gotext.Get("服务器ID"), "ServerId", db.Int).
		FieldFilterable(types.FilterType{FormType: form.SelectSingle}).
		FieldFilterOptions(serverOps).FieldHide()

	info.AddField(gotext.Get("角色ID"), "PlayerId", db.Int).FieldFilterable(types.FilterType{HelpMsg: "角色ID或者 角色名字 不能为空 否则无法查询到结果"})
	info.AddField(gotext.Get("角色名字"), "PlayerName", db.Text).FieldFilterable()
	info.AddField(gotext.Get("角色密码"), "pwd", db.Text)

	info.SetGetDataFn(func(param parameter.Parameters) (data []map[string]interface{}, size int) {
		PlayerName := ctx.Request.FormValue("PlayerName")
		PlayerId, _ := strconv.Atoi(ctx.Request.FormValue("PlayerId"))
		ServerId, _ := strconv.Atoi(ctx.Request.FormValue("ServerId"))
		if ServerId == 0 || (PlayerId == 0 && PlayerName == "") {
			return nil, 0
		}

		db_char := fusion.GetServerGormDB("db_char", int64(ServerId))
		if db_char == nil {
			logger.Error("can`t find db_char serverId = " + strconv.Itoa(ServerId))
			return nil, 0
		}
		if PlayerId == 0 && PlayerName == "" {
			return nil, 0
		}
		sql := "SELECT SUBSTRING_INDEX(SUBSTRING_INDEX(ipcComponentStatus->'$.RoleOperationLock', ',', 1), ':', -1) as pwd,ipcInstID as PlayerId,ipcNickName as PlayerName FROM inst_player_char"
		where := ""
		if PlayerId != 0 {
			where += "ipcInstID = " + strconv.Itoa(PlayerId) + " and "
		}
		if PlayerName != "" {
			where += "ipcNickName like '%" + PlayerName + "%' "
		}
		if len(where) != 0 {
			sql += " where " + where + ";"
		}
		db_char.Raw(sql).Scan(&data)
		return data, len(data)
	})
	info.HideNewButton().HideEditButton().HideDeleteButton().HideDetailButton().HideRowSelector().HideExportButton()
	info.SetTitle(gotext.Get("查找角色密码"))

	return retrievePWD
}
