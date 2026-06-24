package AutoPages

import (
	"admin/common"
	"admin/common/def"
	"admin/fusion"
	"github.com/GoAdminGroup/go-admin/context"
	"github.com/GoAdminGroup/go-admin/modules/db"
	form2 "github.com/GoAdminGroup/go-admin/plugins/admin/modules/form"
	"github.com/GoAdminGroup/go-admin/plugins/admin/modules/table"
	"github.com/GoAdminGroup/go-admin/template/types"
	"github.com/GoAdminGroup/go-admin/template/types/action"
	"github.com/GoAdminGroup/go-admin/template/types/form"
	"github.com/leonelquinteros/gotext"
	"html/template"
	"net/url"
	"strconv"
	"strings"
	"time"
)

func GetTNoticesTable(ctx *context.Context) table.Table {

	tNotices := table.NewDefaultTable(table.Config{
		Driver:     "mysql",
		Connection: "db_global",
		CanAdd:     true,
		Editable:   true,
		Deletable:  true,
		Exportable: true,
		PrimaryKey: table.PrimaryKey{
			Type: db.Int,
			Name: "Id",
		},
	})
	info := tNotices.GetInfo()

	info.HideDetailButton().HideRowSelector().HideFilterArea()
	myOP := fusion.GetServerList()

	info.AddField("Id", "Id", db.Int)
	info.AddField(gotext.Get("开始时间"), "startTime", db.Datetime)
	info.AddField(gotext.Get("结束时间"), "endTime", db.Datetime)
	info.AddField(gotext.Get("间隔")+"("+gotext.Get("单位")+":"+gotext.Get("秒")+")", "interval", db.Int)
	info.AddField(gotext.Get("内容"), "msg", db.Text)
	info.AddField("MsgFlags", "msgFlags", db.Int).
		FieldDisplay(func(model types.FieldModel) interface{} {
			temp, _ := strconv.Atoi(model.Value)
			return def.GetOpStringsMsgFlag()[temp]
		})
	info.AddField(gotext.Get("发送频道"), "toChannels", db.Int).
		FieldDisplay(func(model types.FieldModel) interface{} {
			temp, _ := strconv.Atoi(model.Value)
			return def.GetOpStringsChannel()[temp]
		})
	info.AddField(gotext.Get("发往指定服务器"), "gsIds", db.Text)
	info.AddField(gotext.Get("不发往指定服务器"), "gsNotIds", db.Text)
	info.AddField(gotext.Get("是否循环"), "isOnce", db.Tinyint).
		FieldFilterable(types.FilterType{FormType: form.SelectSingle}).
		FieldDisplay(func(model types.FieldModel) interface{} {
			if model.Value == "1" {
				return gotext.Get("是")
			}
			return gotext.Get("否")
		}).FieldFilterOptions(types.FieldOptions{
		{Text: gotext.Get("是"), Value: "1"},
		{Text: gotext.Get("否"), Value: "0"}}).FieldHide()
	info.SetTable("t_notices").SetTitle(gotext.Get("广播列表"))
	info.SetDeleteHook(func(ids []string) error {
		for _, v := range ids {
			if v != "" {
				vars := &url.Values{}
				vars.Add("noticeId", v)
				vars.Add("isCancelNotice", "true")
				ctx1 := &context.Context{}
				err, res := fusion.CallToCenter(vars, common.MyCfg.Api["pushServerNotice"], "GET", ctx1)
				if err != nil {
					return err
				}
				if res != "" {
					return err
				}
			}
		}
		return nil
	})

	info.AddButton(template.HTML(gotext.Get("将选中条目同步到服务器")), "",
		action.PopUp("/notice/send", gotext.Get("操作结果"), func(ctx *context.Context) (success bool, msg string, data interface{}) {
			failed := []string{}
			ids := ctx.Request.FormValue("ids")
			noticeIds := strings.Split(ids, ",")
			ifEmpty := true
			for _, v := range noticeIds {
				if v != "" {
					ifEmpty = false
					vars := &url.Values{}
					vars.Add("noticeId", v)
					vars.Add("isCancelNotice", "false")
					ctx1 := &context.Context{}
					err, res := fusion.CallToCenter(vars, common.MyCfg.Api["pushServerNotice"], "GET", ctx1)
					if err != nil {
						failed = append(failed, "<h5>"+v+":"+err.Error()+"</h5>")
					}
					if res != "" {
						failed = append(failed, "<h5>"+res+"</h5>")
					}
				}
			}
			if ifEmpty {
				return false, gotext.Get("请选择条目"), ""
			}
			return true, "", failed
		}))

	info.AddButton(template.HTML(gotext.Get("向服务器请求取消选中条目")), "",
		action.PopUp("/notice/cancel", gotext.Get("操作结果"), func(ctx *context.Context) (success bool, msg string, data interface{}) {
			failed := []string{}
			ids := ctx.Request.FormValue("ids")
			noticeIds := strings.Split(ids, ",")
			ifEmpty := true
			for _, v := range noticeIds {
				if v != "" {
					ifEmpty = false
					vars := &url.Values{}
					vars.Add("noticeId", v)
					vars.Add("isCancelNotice", "true")
					ctx1 := &context.Context{}
					err, res := fusion.CallToCenter(vars, common.MyCfg.Api["pushServerNotice"], "GET", ctx1)
					if err != nil {
						failed = append(failed, "<h5>"+v+":"+err.Error()+"</h5>")
					}
					if res != "" {
						failed = append(failed, "<h5>"+res+"</h5>")
					}
				}
			}
			if ifEmpty {
				return false, gotext.Get("请选择条目"), ""
			}
			return true, "", failed
		}))

	formList := tNotices.GetForm()
	formList.AddField("Id", "Id", db.Int, form.Number).
		FieldDisableWhenCreate()
	formList.AddField(gotext.Get("开始时间"), "startTime", db.Datetime, form.Datetime).FieldMust()
	formList.AddField(gotext.Get("结束时间"), "endTime", db.Datetime, form.Datetime).FieldMust()
	formList.AddField(gotext.Get("间隔")+"("+gotext.Get("单位")+":"+gotext.Get("秒")+")",
		"interval", db.Int, form.Number).FieldDefault("0").FieldMust()
	formList.AddField(gotext.Get("内容"), "msg", db.Text, form.TextArea).FieldMust()
	formList.AddField("MsgFlags", "msgFlags", db.Int, form.SelectSingle).FieldDefault("0").
		FieldMust().FieldOptions(def.GetMsgFlagOptions())
	formList.AddField(gotext.Get("发送频道"), "toChannels", db.Int, form.SelectSingle).FieldDefault("0").
		FieldMust().FieldOptions(def.GetChannelOptions()).FieldDefault("6")
	formList.AddField(gotext.Get("发往指定服务器"), "gsIds", db.Text, form.SelectBox).
		FieldOptions(myOP).FieldHelpMsg(template.HTML(gotext.Get("在此项选择要发往的服务器,若不填则向所有服务器发送")))

	formList.AddField(gotext.Get("不发往指定服务器"), "gsNotIds", db.Text, form.SelectBox).
		FieldOptions(myOP).
		FieldHelpMsg(template.HTML(gotext.Get("在此项选择不发往的服务器")))
	formList.AddField(gotext.Get("是否循环"), "isOnce", db.Tinyint, form.SelectSingle).FieldOptions(types.FieldOptions{
		{Text: gotext.Get("是"), Value: "1"},
		{Text: gotext.Get("否"), Value: "0"}}).FieldHide()
	formList.SetTable("t_notices").SetTitle(gotext.Get("添加广播")).HideContinueNewCheckBox().HideContinueEditCheckBox()

	formList.SetInsertFn(func(values form2.Values) error {
		notice := def.NoticeData{}
		startTime := values["startTime"][0]
		endTime := values["endTime"][0]
		start, _ := time.ParseInLocation("2006-01-02 15:04:05", startTime, time.Local)
		end, _ := time.ParseInLocation("2006-01-02 15:04:05", endTime, time.Local)
		notice.StartTime = def.MyTime(start)
		notice.EndTime = def.MyTime(end)
		notice.Msg = values["msg"][0]
		temp, _ := strconv.Atoi(values["msgFlags"][0])
		notice.MsgFlags = uint32(temp)
		temp, _ = strconv.Atoi(values["interval"][0])
		notice.Interval = uint32(temp)
		temp, _ = strconv.Atoi(values["isOnce"][0])
		notice.IsOnce = temp
		temp, _ = strconv.Atoi(values["toChannels"][0])
		notice.ToChannels = uint32(temp)

		temp1 := values["gsIds[]"]
		for _, v := range temp1 {
			notice.GsIds += string(v) + ","
		}
		notice.GsIds = strings.TrimSuffix(notice.GsIds, ",")

		temp1 = values["gsNotIds[]"]
		for _, v := range temp1 {
			notice.GsNotIds += string(v) + ","
		}
		notice.GsNotIds = strings.TrimSuffix(notice.GsNotIds, ",")
		db_global := fusion.GetBaseGormDB("db_global")
		db_global.Create(&notice)
		return nil
	})

	return tNotices
}
