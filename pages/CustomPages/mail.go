package CustomPages

import (
	"admin/common/def"
	"admin/common/def/MyTime"
	"admin/common/def/chequeDef"
	"admin/common/def/currencyDef"
	"admin/common/def/operationAudit"
	"admin/fusion"
	"encoding/json"
	"fmt"
	"html/template"
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
	template2 "github.com/GoAdminGroup/go-admin/template"
	"github.com/GoAdminGroup/go-admin/template/icon"
	"github.com/GoAdminGroup/go-admin/template/types"
	"github.com/GoAdminGroup/go-admin/template/types/form"
	"github.com/GoAdminGroup/themes/adminlte/components/infobox"
	"github.com/leonelquinteros/gotext"
)

func GetMailPage(ctx *context.Context) (types.Panel, error) {
	components := template2.Get(config.GetTheme())
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
	var formList = types.NewFormPanel()

	myOP := fusion.GetServerList()
	formList.AddField(gotext.Get("操作类型"), "operation", db.Int, form.SelectSingle).
		FieldDefault("0").
		FieldOptions(types.FieldOptions{
			{Text: gotext.Get("全服邮件"), Value: "0", Selected: true},
			{Text: gotext.Get("个人邮件"), Value: "1"},
		}).
		FieldOnChooseHide("0", "serverID", "receiver").
		FieldOnChooseShow("0", "IpcServerID", "noServerID", "mailTargetTime").
		FieldOnChooseShow("1", "serverID", "receiver")

	formList.AddField(gotext.Get("选中的服务器"), "IpcServerID", db.Int, form.SelectBox).FieldOptions(myOP).
		FieldHelpMsg(template.HTML(gotext.Get("在此项中选择发送邮件的服务器,若不选则默认为向所有服务器发送。")))
	formList.AddField(gotext.Get("排除的服务器"), "noServerID", db.Int, form.SelectBox).FieldOptions(myOP).
		FieldHelpMsg(template.HTML(gotext.Get("在此项中选择不发送邮件的服务器。")))

	formList.AddField(gotext.Get("邮件新功能"), "mailTargetTime", db.Int, form.SelectSingle).
		FieldOptions(types.FieldOptions{
			{Value: "0", Text: gotext.Get("无限制")},
			{Value: "1", Text: gotext.Get("在此条邮件发送时间之前注册的用户可收到")},
			{Value: "2", Text: gotext.Get("在此条邮件发送时间之后注册的用户可收到")},
		}).FieldDefault("0")

	formList.AddField(gotext.Get("服务器ID"), "serverID", db.Int, form.SelectSingle).
		FieldOptions(myOP)
	formList.AddField(gotext.Get("玩家ID"), "receiver", db.Bigint, form.TextArea).
		FieldHelpMsg(template.HTML(gotext.Get("<font color=\"#FF0000\">如有多个，请用“,”分隔。</font>")))

	formList.AddField(gotext.Get("邮件类型"), "MailType", db.Int, form.SelectSingle).
		FieldOptions(types.FieldOptions{
			{Text: gotext.Get("普通"), Value: "0"},
			{Text: gotext.Get("系统"), Value: "1"},
			{Text: "GM", Value: "2"},
		}).FieldMust()
	formList.AddField(gotext.Get("发送时间"), "deliverTime", db.Datetime, form.Datetime).FieldMust()
	formList.AddField(gotext.Get("到期时间"), "expireTime", db.Datetime, form.Datetime).FieldMust()
	formList.AddField(gotext.Get("邮件标题"), "mailTitle", db.Varchar, form.Text).FieldMust()
	formList.AddField(gotext.Get("邮件正文"), "mailBody", db.Varchar, form.TextArea).FieldMust()
	itemOP := fusion.GetItemNameOptionAndID()
	formList.AddTable(gotext.Get("道具"), "Item", func(panel *types.FormPanel) {
		panel.AddField(gotext.Get("道具ID"), "ItemID", db.Int, form.SelectSingle).FieldHide().FieldOptions(itemOP)
		panel.AddField(gotext.Get("道具数量"), "ItemNum", db.Int, form.Number).FieldHide().FieldValue("0")
		panel.AddField(gotext.Get("道具有效期"), "ItemActiveLastTime", db.Int, form.Datetime).FieldHide().FieldValue("0")
	})
	currencyOP := currencyDef.GetCurrencyOps()
	formList.AddTable(gotext.Get("票据"), "cheque", func(panel *types.FormPanel) {
		panel.AddField(gotext.Get("票据ID"), "chequeID", db.Int, form.SelectSingle).FieldHide().FieldOptions(currencyOP)
		panel.AddField(gotext.Get("票据数量"), "chequeNum", db.Int, form.Number).FieldHide().FieldValue("0")
	})

	formList.SetTabGroups(types.TabGroups{
		{"operation", "IpcServerID", "noServerID", "mailTargetTime", "serverID", "receiver", "MailType", "deliverTime", "expireTime", "mailTitle", "mailBody", "Item", "cheque"},
	})
	formList.SetTabHeaders(gotext.Get("邮件"))

	field, headers := formList.GroupField()

	aform := components.Form().
		SetTabHeaders(headers).
		SetTabContents(field).
		SetPrefix(config.PrefixFixSlash()).
		SetUrl("/admin/SendMail").
		SetHiddenFields(map[string]string{
			form2.PreviousKey: "/admin",
		}).
		SetMethod("GET").
		SetOperationFooter(col1 + col2)

	return types.Panel{
		Content: components.Box().
			SetHeader(aform.GetDefaultBoxHeader(true)).
			WithHeadBorder().
			SetHeader(template.HTML(gotext.Get("发送邮件"))).
			SetBody(aform.GetContent() + formList.FooterHtml).
			GetContent(),
		Title:       template.HTML(gotext.Get("邮件管理")),
		Description: template.HTML(gotext.Get("发送邮件")),
		CSS:         `.modal.fade.in{z-index:10002}`,
		Callbacks:   nil,
	}, nil
}

func SendMail(ctx *context.Context) (types.Panel, error) {
	var mailInfo = def.TMails{}

	mailInfo.MailSubject = ctx.Request.FormValue("mailTitle")
	mailInfo.MailBody = ctx.Request.FormValue("mailBody")
	mailInfo.MailDeliverTime = ctx.Request.FormValue("deliverTime")
	mailInfo.MailExpireTime = ctx.Request.FormValue("expireTime")
	mailType, _ := strconv.Atoi(ctx.Request.FormValue("MailType"))
	operation, _ := strconv.Atoi(ctx.Request.FormValue("operation"))
	mailTargetTime, _ := strconv.Atoi(ctx.Request.FormValue("mailTargetTime"))

	components := template2.Get(config.GetTheme())
	infobox1 := infobox.New()
	col1 := components.Col()
	var size = map[string]string{"md": "3", "sm": "6", "xs": "12"}
	str := gotext.Get("返回上级")
	preUrl := ctx.Request.Header.Get("Referer")
	index := strings.Index(preUrl, "admin")
	preUrl = preUrl[index:]

	ReportError := func(msg string) (types.Panel, error) {
		infoboxCol1 := col1.SetSize(size).SetContent(infobox1.SetText(template.HTML(msg)).GetContent()).GetContent()
		return types.Panel{
			Content:     infoboxCol1,
			Title:       template.HTML(gotext.Get("邮件管理")),
			Description: template.HTML(fmt.Sprintf(`<a href="/%s">%s</a>`, preUrl, str)),
			Callbacks:   nil,
		}, nil
	}

	serverID := 0
	var playerIDs []int
	//var playerNickNames []string

	var serverIDs []string
	var noServerIDs []string

	if operation == 0 {
		formValue := ctx.Request.Form["IpcServerID[]"]
		for _, v := range formValue {
			serverIDs = append(serverIDs, v)
		}
		for _, v := range serverIDs {
			mailInfo.GsIds += string(v) + ","
		}
		mailInfo.GsIds = strings.TrimSuffix(mailInfo.GsIds, ",")

		formValue = ctx.Request.Form["noServerID[]"]
		for _, v := range formValue {
			noServerIDs = append(noServerIDs, v)
		}

		for _, v := range noServerIDs {
			mailInfo.GsNotIds += string(v) + ","
		}
		mailInfo.GsNotIds = strings.TrimSuffix(mailInfo.GsNotIds, ",")

	} else {
		serverID, _ = strconv.Atoi(ctx.Request.FormValue("serverID"))
		//serverID, _ = strconv.Atoi(ctx.Request.FormValue("serverID4Single"))
		formValue := ctx.FormValue("receiver")
		var tempS1 = strings.Split(formValue, ",")
		for _, v := range tempS1 {
			playerId, _ := strconv.Atoi(v)
			if playerId != 0 {
				playerIDs = append(playerIDs, playerId)
			}
		}
		if len(playerIDs) == 0 {
			infoboxCol1 := col1.SetSize(size).SetContent(infobox1.SetText(template.HTML(gotext.Get("请填写正确的玩家id"))).GetContent()).GetContent()
			return types.Panel{
				Content:     infoboxCol1,
				Title:       template.HTML(gotext.Get("邮件管理")),
				Description: template.HTML(fmt.Sprintf(`<a href="/%s">%s</a>`, preUrl, str)),
				Callbacks:   nil,
			}, nil
		}
	}
	mailInfo.MailTargetTime = uint8(mailTargetTime)
	mailInfo.MailType = uint32(mailType)
	formValue1 := ctx.Request.Form["ItemID"]
	formValue2 := ctx.Request.Form["ItemNum"]
	formValue3 := ctx.Request.Form["ItemActiveLastTime"]

	mailAttachItem := make([]string, 0)
	mailAttachCheque := make([]string, 0)
	for k, _ := range formValue1 {
		itemId, _ := strconv.Atoi(formValue1[k])
		itemNum, _ := strconv.Atoi(formValue2[k])
		itemActive := formValue3[k]

		var timeDuration = 0
		dateTime, err := time.ParseInLocation("2006-01-02 15:04:05", itemActive, time.Local)
		if err == nil {
			timeDuration = int(dateTime.Unix())
		}

		if itemId != 0 && itemNum != 0 {
			tempItem := "0," + strconv.Itoa(itemId) + "," + strconv.Itoa(itemNum) + ",0,1," + strconv.Itoa(timeDuration)
			mailAttachItem = append(mailAttachItem, tempItem)
		}

	}
	temp, _ := json.Marshal(mailAttachItem)
	mailInfo.MailItems = string(temp)
	formValue1 = ctx.Request.Form["chequeID"]
	formValue2 = ctx.Request.Form["chequeNum"]
	for k, _ := range formValue1 {
		chequeID, _ := strconv.Atoi(formValue1[k])
		chequeNum, _ := strconv.Atoi(formValue2[k])
		if chequeNum != 0 {
			Cheque := chequeDef.CurrencyToCheque(chequeID)
			mailAttachCheque = append(mailAttachCheque, strconv.Itoa(Cheque)+","+strconv.Itoa(chequeNum))
		}
	}
	temp, _ = json.Marshal(mailAttachCheque)
	mailInfo.MailCheques = string(temp)

	//now := time.Now()
	timingTMail := def.TimingTMail{}
	timingTMail.Mail = mailInfo
	timingTMail.PlayerIDs = playerIDs
	timingTMail.GsId = serverID
	temp1, _ := json.Marshal(timingTMail)

	Transit := def.MailTransit{}
	Transit.ExpireTime = mailInfo.MailExpireTime
	Transit.DeliverTime = mailInfo.MailDeliverTime
	Transit.Record = string(temp1)
	Transit.IsSend = 0

	var operationAuditInfo = operationAudit.Operationaudit{}
	operationAuditInfo.OperatorId = int(auth.Auth(ctx).Id)
	operationAuditInfo.SubmissionTime = MyTime.MyTime(time.Now())
	operationAuditInfo.AuditStatus = operationAudit.AuditStatus_Wait
	operationAuditInfo.OperationType = operationAudit.SendMail
	temp2, _ := json.Marshal(Transit)
	operationAuditInfo.OperationalData = fusion.Bytes2String(temp2)

	db_manager := fusion.GetBaseGormDB("go_manager")
	db_manager.Create(&operationAuditInfo)
	logger.Info("End SendMail")
	//err := pageMgr.SaveTimingMail2DB(&mailInfo, playerIDs, serverID)
	//msg := gotext.Get("操作成功")
	//if err != nil {
	//	msg = err.Error()
	//}
	return ReportError("已提交审核")

	//err := pageMgr.SaveTimingMail2DB(&mailInfo, playerIDs, serverID)
	//msg := gotext.Get("操作成功")
	//if err != nil {
	//	msg = err.Error()
	//}
	//
	//infoboxCol1 := col1.SetSize(size).SetContent(infobox1.SetText(template.HTML(msg)).GetContent()).GetContent()
	//return types.Panel{
	//	Content:     infoboxCol1,
	//	Title:       template.HTML(gotext.Get("邮件管理")),
	//	Description: template.HTML(fmt.Sprintf(`<a href="/%s">%s</a>`, preUrl, str)),
	//	Callbacks:   nil,
	//}, nil
}

func GetCompensationMail(ctx *context.Context) (types.Panel, error) {
	AcctId := ctx.Request.FormValue("AcctId")
	ServerID, _ := strconv.Atoi(ctx.Request.FormValue("IpcServerID"))
	PlayerId := ctx.Request.FormValue("PlayerId")
	PlayerName := ctx.Request.FormValue("PlayerName")
	LogTime := ctx.Request.FormValue("LogTime")
	AddFailedItems := ctx.Request.FormValue("AddFailedItems")
	Id := ctx.Request.FormValue("Id")

	components := template2.Get(config.GetTheme())
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

	serverStr := fusion.GetServerNameStrings()
	var formList = types.NewFormPanel()

	formList.AddField("id", "Id", db.Varchar, form.Text).
		FieldValue(Id).FieldHide()
	formList.AddField(gotext.Get("账号ID"), "AcctId", db.Varchar, form.Text).
		FieldValue(AcctId).FieldDisplayButCanNotEditWhenUpdate()
	formList.AddField(gotext.Get("角色Id"), "PlayerId", db.Varchar, form.Text).
		FieldValue(PlayerId).FieldDisplayButCanNotEditWhenUpdate()
	formList.AddField(gotext.Get("玩家昵称"), "PlayerName", db.Varchar, form.Text).
		FieldValue(PlayerName).FieldDisplayButCanNotEditWhenUpdate()
	formList.AddField(gotext.Get("服务器Id"), "ServerId", db.Int, form.Text).
		FieldValue(strconv.Itoa(ServerID)).
		FieldDisplay(func(model types.FieldModel) interface{} {
			temp, _ := strconv.Atoi(model.Value)
			if serverStr[temp] != "" {
				return serverStr[temp]
			}
			return ""
		}).FieldDisplayButCanNotEditWhenUpdate()
	formList.AddField(gotext.Get("记录时间"), "LogTime", db.Varchar, form.Text).
		FieldValue(LogTime).FieldDisplayButCanNotEditWhenUpdate()

	formList.AddField(gotext.Get("发送时间"), "deliverTime", db.Datetime, form.Datetime).FieldMust()
	formList.AddField(gotext.Get("到期时间"), "expireTime", db.Datetime, form.Datetime).FieldMust()

	formList.AddField(gotext.Get("未成功添加的道具"), "AddFailedItems", db.Text, form.TextArea).
		FieldValue(AddFailedItems).FieldDisplayButCanNotEditWhenUpdate()

	itemOP := fusion.GetItemNameOption()
	formList.AddTable(gotext.Get("额外补发道具"), "Item", func(panel *types.FormPanel) {
		panel.AddField(gotext.Get("道具ID"), "ItemID", db.Int, form.SelectSingle).FieldHide().FieldOptions(itemOP)
		panel.AddField(gotext.Get("道具数量"), "ItemNum", db.Int, form.Number).FieldHide().FieldValue("0")
	})
	currencyOP := currencyDef.GetCurrencyOps()
	formList.AddTable(gotext.Get("额外补发票据"), "cheque", func(panel *types.FormPanel) {
		panel.AddField(gotext.Get("票据ID"), "chequeID", db.Int, form.SelectSingle).FieldHide().FieldOptions(currencyOP)
		panel.AddField(gotext.Get("票据数量"), "chequeNum", db.Int, form.Number).FieldHide().FieldValue("0")
	})

	formList.AddField(gotext.Get("邮件标题"), "mailTitle", db.Varchar, form.Text).FieldMust().FieldValue("道具补偿邮件")

	formList.AddField(gotext.Get("邮件正文"), "mailBody", db.Varchar, form.TextArea).FieldMust().
		FieldValue("经核实，您于" + LogTime + "获得的道具数量存在异常，现已补发给您。我们对此次问题深表歉意")

	aform := components.Form().
		SetHeader(formList.Header).
		SetContent(formList.FieldList).
		SetPrefix(config.PrefixFixSlash()).
		SetUrl("/admin/SendCompensationMail").
		SetHiddenFields(map[string]string{
			form2.PreviousKey: "/admin",
		}).
		SetMethod("GET").
		SetOperationFooter(col1 + col2)

	return types.Panel{
		Content: components.Box().
			SetHeader(aform.GetDefaultBoxHeader(true)).
			WithHeadBorder().
			SetHeader(template.HTML(gotext.Get("发送邮件"))).
			SetBody(aform.GetContent() + formList.FooterHtml).
			GetContent(),
		Title:       template.HTML(gotext.Get("邮件管理")),
		Description: template.HTML(gotext.Get("发送邮件")),
		CSS:         `.modal.fade.in{z-index:10002}`,
		Callbacks:   nil,
	}, nil
}

func SendCompensationMail(ctx *context.Context) (types.Panel, error) {
	tMail := def.TMails{}
	tMail.MailType = 1
	tMail.MailSubject = ctx.Request.FormValue("mailTitle")
	tMail.MailBody = ctx.Request.FormValue("mailBody")
	tMail.MailDeliverTime = ctx.Request.FormValue("deliverTime")
	tMail.MailExpireTime = ctx.Request.FormValue("expireTime")
	//AddFailedItems := ctx.Request.FormValue("AddFailedItems")
	playerId, _ := strconv.Atoi(ctx.Request.FormValue("PlayerId"))
	ServerId, _ := strconv.Atoi(ctx.Request.FormValue("ServerId"))
	Id, _ := strconv.Atoi(ctx.Request.FormValue("Id"))

	msg := gotext.Get("操作成功")
	db_log := fusion.GetServerGormDB("db_log", int64(ServerId))
	if db_log != nil {
		mailAttachItem := make([]string, 0)
		mailAttachCheque := make([]string, 0)
		formValue1 := ctx.Request.Form["ItemID"]
		formValue2 := ctx.Request.Form["ItemNum"]
		arr := make([]string, 0)
		type errorCut struct {
			IfSolve        int8   `gorm:"column:ifSolve"`
			AddFailedItems string `gorm:"column:addFailedItems"`
		}

		cut := errorCut{}
		db_log.Table("log_player_item_add_error").
			Select("ifSolve,addFailedItems").
			Where("Id = ? ", Id/100000).
			Scan(&cut)
		if cut.IfSolve == 0 {
			b := fusion.String2Bytes(cut.AddFailedItems)
			json.Unmarshal(b, &arr)

			mailAttachItem = append(mailAttachItem, arr...)

			for k, _ := range formValue1 {
				itemId, _ := strconv.Atoi(formValue1[k])
				itemNum, _ := strconv.Atoi(formValue2[k])
				if itemId != 0 && itemNum != 0 {
					tempItem := "0," + strconv.Itoa(itemId) + "," + strconv.Itoa(itemNum) + ",0,0"
					mailAttachItem = append(mailAttachItem, tempItem)
				}
			}

			temp1, _ := json.Marshal(mailAttachItem)
			tMail.MailItems = string(temp1)

			formValue1 = ctx.Request.Form["chequeID"]
			formValue2 = ctx.Request.Form["chequeNum"]
			for k, _ := range formValue1 {
				chequeID, _ := strconv.Atoi(formValue1[k])
				chequeNum, _ := strconv.Atoi(formValue2[k])
				if chequeNum != 0 {
					Cheque := chequeDef.CurrencyToCheque(chequeID)
					mailAttachCheque = append(mailAttachCheque, strconv.Itoa(Cheque)+","+strconv.Itoa(chequeNum))
				}
			}
			temp1, _ = json.Marshal(mailAttachCheque)
			tMail.MailCheques = string(temp1)

			//now := time.Now()
			timingTMail := def.TimingTMail{}
			timingTMail.Mail = tMail
			temp1, _ = json.Marshal(timingTMail)

			timingTMail.PlayerIDs = append(timingTMail.PlayerIDs, playerId)
			timingTMail.GsId = ServerId
			Transit := def.MailTransit{}
			Transit.ExpireTime = tMail.MailExpireTime
			Transit.DeliverTime = tMail.MailDeliverTime
			Transit.Record = string(temp1)
			Transit.IsSend = 0

			var operationAuditInfo = operationAudit.Operationaudit{}
			operationAuditInfo.OperatorId = int(auth.Auth(ctx).Id)
			operationAuditInfo.SubmissionTime = MyTime.MyTime(time.Now())
			operationAuditInfo.AuditStatus = operationAudit.AuditStatus_Wait
			operationAuditInfo.OperationType = operationAudit.SendMail
			temp2, _ := json.Marshal(Transit)
			operationAuditInfo.OperationalData = fusion.Bytes2String(temp2)

			db_manager := fusion.GetBaseGormDB("go_manager")
			db_manager.Create(&operationAuditInfo)
			logger.Info("End SendMail")
			//err := pageMgr.SaveTimingMail2DB(&mailInfo, playerIDs, serverID)
			//msg := gotext.Get("操作成功")
			//if err != nil {
			//	msg = err.Error()
			//}
			msg = gotext.Get("已提交审核！")
			//err := pageMgr.SaveTimingMail2DB(&tMail, []int{playerId}, ServerId)
			//if err != nil {
			//	msg = err.Error()
			//}
			db_log.Table("log_player_item_add_error").Where("id= ?", Id/100000).Update("ifSolve", 1)
		} else {
			msg = gotext.Get("操作失败！该条已处理！")
		}
	} else {
		msg = gotext.Get("操作失败！")
	}

	components := template2.Get(config.GetTheme())
	infobox1 := infobox.New()
	col1 := components.Col()
	var size = map[string]string{"md": "3", "sm": "6", "xs": "12"}
	infoboxCol1 := col1.SetSize(size).SetContent(infobox1.SetText(template.HTML(msg)).GetContent()).GetContent()
	return types.Panel{
		Content:   infoboxCol1,
		Title:     template.HTML(gotext.Get("发送补偿邮件")),
		Callbacks: nil,
	}, nil
}
