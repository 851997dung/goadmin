package AutoPages

import (
	"admin/common/def"
	"admin/common/def/MyTime"
	"admin/common/def/operationAudit"
	"admin/fusion"
	"admin/mgr/pageMgr"
	"encoding/json"
	"errors"
	"strconv"
	"time"

	"github.com/GoAdminGroup/go-admin/context"
	"github.com/GoAdminGroup/go-admin/modules/auth"
	"github.com/GoAdminGroup/go-admin/modules/db"
	form2 "github.com/GoAdminGroup/go-admin/plugins/admin/modules/form"
	"github.com/GoAdminGroup/go-admin/plugins/admin/modules/table"
	"github.com/GoAdminGroup/go-admin/template/types"
	"github.com/GoAdminGroup/go-admin/template/types/form"
	"github.com/leonelquinteros/gotext"
)

func GetOperationauditTable(ctx *context.Context) table.Table {

	operationaudit := table.NewDefaultTable(table.Config{
		Driver:     "mysql",
		Connection: "go_manager",
		CanAdd:     false,
		Editable:   true,
		Deletable:  false,
		Exportable: true,
		PrimaryKey: table.PrimaryKey{
			Type: db.Int,
			Name: "id",
		},
	})

	authNameStr := fusion.GetAllAuthName()
	//serverNameStr := fusion.GetAllServerNameStrings()
	//payShopStr := fusion.GetAllRechargeName()

	info := operationaudit.GetInfo().HideDeleteButton().HideDetailButton()

	info.AddField("Id", "id", db.Int)
	info.AddField(gotext.Get("操作者"), "operatorId", db.Int).
		FieldDisplay(func(model types.FieldModel) interface{} {
			fromIdx, _ := strconv.Atoi(model.Value)
			if authNameStr[fromIdx] != "" {
				return authNameStr[fromIdx]
			}
			return model.Value
		})
	info.AddField(gotext.Get("操作类型"), "operationType", db.Int).
		FieldDisplay(func(model types.FieldModel) interface{} {
			temp, _ := strconv.Atoi(model.Value)
			switch temp {
			//case operationAudit.SimRecharge:
			//	return "模拟充值"
			//case operationAudit.ActivityCompensation:
			//	return "活动补档"
			case operationAudit.SendMail:
				return gotext.Get("发送邮件")
			}
			return gotext.Get(model.Value)
		})
	info.AddField(gotext.Get("操作数据"), "operationalData", db.Mediumtext).
		FieldHide()
	info.AddField(gotext.Get("审核状态"), "auditStatus", db.Int).
		FieldDisplay(func(model types.FieldModel) interface{} {
			temp, _ := strconv.Atoi(model.Value)
			switch temp {
			case operationAudit.AuditStatus_Wait:
				return gotext.Get("等待审批")
			case operationAudit.AuditStatus_Pass:
				return gotext.Get("通过")
			case operationAudit.AuditStatus_Refuse:
				return gotext.Get("拒绝")
			}
			return model.Value
		})
	info.AddField(gotext.Get("审核者"), "auditorId", db.Int).
		FieldDisplay(func(model types.FieldModel) interface{} {
			fromIdx, _ := strconv.Atoi(model.Value)
			if authNameStr[fromIdx] != "" {
				return authNameStr[fromIdx]
			}
			return model.Value
		})
	info.AddField(gotext.Get("提交时间"), "submissionTime", db.Datetime)
	info.AddField(gotext.Get("审核时间"), "auditTime", db.Datetime)
	info.AddField(gotext.Get("处理状态"), "finishStatus", db.Int).
		FieldDisplay(func(model types.FieldModel) interface{} {
			value, _ := strconv.Atoi(model.Value)
			switch value {
			case 0:
				return gotext.Get("等待操作")
			case 1:
				return gotext.Get("操作完成")
			}
			return model.Value
		})

	info.SetTable("operationaudit").SetTitle(gotext.Get("操作审核表"))

	formList := operationaudit.GetForm()
	formList.AddField("Id", "id", db.Int, form.Default).
		FieldDisableWhenCreate().FieldDisplayButCanNotEditWhenUpdate()
	formList.AddField(gotext.Get("操作者"), "operatorId", db.Int, form.Number).
		FieldDisplayButCanNotEditWhenUpdate().
		FieldDisplay(func(model types.FieldModel) interface{} {
			fromIdx, _ := strconv.Atoi(model.Value)
			if authNameStr[fromIdx] != "" {
				return authNameStr[fromIdx]
			}
			return model.Value
		})
	formList.AddField(gotext.Get("操作类型"), "operationType", db.Int, form.Number).
		FieldDisplayButCanNotEditWhenUpdate().
		FieldDisplay(func(model types.FieldModel) interface{} {
			temp, _ := strconv.Atoi(model.Value)
			switch temp {
			//case operationAudit.SimRecharge:
			//	return "模拟充值"
			//case operationAudit.ActivityCompensation:
			//	return "活动补档"
			case operationAudit.SendMail:
				return gotext.Get("发送邮件")
			}
			return gotext.Get(model.Value)
		})
	formList.AddField(gotext.Get("操作数据"), "operationalData", db.Mediumtext, form.TextArea).
		FieldDisplayButCanNotEditWhenUpdate().
		FieldDisplay(func(model types.FieldModel) interface{} {
			temp := model.Row["operationType"]
			opType, _ := strconv.Atoi(fusion.Strval(temp))
			switch opType {
			//case operationAudit.SimRecharge:
			//	var temp1 = fusion.String2Bytes(model.Value)
			//	simRechargeInfo := def.SimRecharge{}
			//	json.Unmarshal(temp1, &simRechargeInfo)
			//	return "服务器id：" + serverNameStr[simRechargeInfo.GsId] + "\r\n" +
			//		"玩家id：" + strconv.Itoa(simRechargeInfo.PlayerId) + "\r\n" +
			//		"模拟充值id：" + simRechargeInfo.PayShopId + "\r\n"
			case operationAudit.SendMail:
				var temp1 = fusion.String2Bytes(model.Value)
				mailAuditInfo := def.MailTransit{}
				err := json.Unmarshal(temp1, &mailAuditInfo)
				if err != nil {
					return model.Value
				}
				return fusion.GetFormatStr4TMailData(&mailAuditInfo)
			}
			return model.Value
		})

	formList.AddField(gotext.Get("审核状态"), "auditStatus", db.Int, form.SelectSingle).
		FieldOptions(types.FieldOptions{
			{Text: gotext.Get("等待审批"), Value: strconv.Itoa(operationAudit.AuditStatus_Wait)},
			{Text: gotext.Get("通过"), Value: strconv.Itoa(operationAudit.AuditStatus_Pass)},
			{Text: gotext.Get("拒绝"), Value: strconv.Itoa(operationAudit.AuditStatus_Refuse)},
		})
	formList.AddField(gotext.Get("审核者"), "auditorId", db.Int, form.Number).
		FieldDisplay(func(model types.FieldModel) interface{} {
			fromIdx, _ := strconv.Atoi(model.Value)
			if authNameStr[fromIdx] != "" {
				return authNameStr[fromIdx]
			}
			return model.Value
		}).FieldHideWhenUpdate()
	formList.AddField(gotext.Get("提交时间"), "submissionTime", db.Datetime, form.Datetime).
		FieldDisplayButCanNotEditWhenUpdate()
	formList.AddField(gotext.Get("审核时间"), "auditTime", db.Datetime, form.Datetime).
		FieldHideWhenUpdate().FieldHideWhenCreate()
	formList.AddField(gotext.Get("处理状态"), "finishStatus", db.Int, form.Text).
		FieldDisplayButCanNotEditWhenUpdate().
		FieldHideWhenCreate().
		FieldHideWhenUpdate().
		FieldDisplay(func(model types.FieldModel) interface{} {
			value, _ := strconv.Atoi(model.Value)
			switch value {
			case 0:
				return gotext.Get("等待操作")
			case 1:
				return gotext.Get("操作完成")
			}
			return model.Value
		})
	formList.SetTable("operationaudit").SetTitle(gotext.Get("操作审核表"))

	formList.SetUpdateFn(func(values form2.Values) error {
		id := values.Get("id")
		go_manager := fusion.GetBaseGormDB("go_manager")

		var finishStatus int
		go_manager.Table("operationaudit").Select("finishStatus").Where("id = ?", id).
			Scan(&finishStatus)
		if finishStatus != 0 {
			return errors.New("不可重复操作！")
		}
		auditStatus, _ := strconv.Atoi(values.Get("auditStatus"))

		var auditArgs string

		switch auditStatus {
		case operationAudit.AuditStatus_Wait:
			return errors.New("请进行审批操作，选择审核状态")
		case operationAudit.AuditStatus_Refuse:
			go_manager.Table("operationaudit").
				Where("id = ?", id).
				Update("auditStatus", auditStatus).
				Update("auditorId", int(auth.Auth(ctx).Id)).
				Update("auditTime", MyTime.MyTime(time.Now())).
				Update("finishStatus", 1)
			return nil
		case operationAudit.AuditStatus_Pass:
			//err := fusion.GoPoolRunWithBlock(func(ctx context2.Context) error {
			var idx, _ = strconv.Atoi(id)
			err, s := doOperationAfterAudit(idx)
			auditArgs = s
			if err != nil {
				return err
			}
			//})
			//if err != nil {
			//	return err
			//}
			go_manager.Table("operationaudit").
				Where("id = ?", id).
				Update("auditStatus", auditStatus).
				Update("auditorId", int(auth.Auth(ctx).Id)).
				Update("auditTime", MyTime.MyTime(time.Now())).
				Update("finishStatus", 1).
				Update("finishArgs", auditArgs)
			return nil
		}
		return nil
	})

	return operationaudit
}

func doOperationAfterAudit(id int) (error, string) {
	var operationauditInfo = operationAudit.Operationaudit{}
	go_manager := fusion.GetBaseGormDB("go_manager")
	go_manager.Table("operationaudit").
		Select("*").
		Where("id = ?", id).
		Scan(&operationauditInfo)

	if operationauditInfo.Id == 0 || operationauditInfo.FinishStatus != 0 {
		return errors.New("数据错误"), ""
	}
	var mailRes string
	switch operationauditInfo.OperationType {
	case operationAudit.SendMail:
		{
			var temp1 = fusion.String2Bytes(operationauditInfo.OperationalData)
			Transit := def.MailTransit{}
			err := json.Unmarshal(temp1, &Transit)
			if err != nil {
				return errors.New("数据有误！请联系管理员"), "数据有误！请联系管理员"
			}
			err = pageMgr.SaveTimingMail2DB(&Transit)
			if err != nil {
				return errors.New("发送失败！请联系管理员"), "发送失败！请联系管理员"
			}
			break
		}
	}
	return nil, mailRes
}
