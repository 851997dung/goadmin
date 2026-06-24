package CustomPages

import (
	"admin/common"
	"admin/fusion"
	"fmt"
	"html/template"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/GoAdminGroup/go-admin/context"
	"github.com/GoAdminGroup/go-admin/modules/config"
	"github.com/GoAdminGroup/go-admin/modules/db"
	form2 "github.com/GoAdminGroup/go-admin/plugins/admin/modules/form"
	template2 "github.com/GoAdminGroup/go-admin/template"
	"github.com/GoAdminGroup/go-admin/template/icon"
	"github.com/GoAdminGroup/go-admin/template/types"
	"github.com/GoAdminGroup/go-admin/template/types/form"
	"github.com/leonelquinteros/gotext"
	"go.uber.org/zap"
)

type DailySignInConfig struct {
	StartTime int64 `json:"start_time"`
	EndTime   int64 `json:"end_time"`
	Period    int   `json:"period"`
}

func GetDailySignInPage(ctx *context.Context) (types.Panel, error) {
	components := template2.Get(config.GetTheme())

	// 1. 按钮组件
	loadBtn := template.HTML(`<button type="button" id="loadConfigBtn" class="btn btn-info">` + gotext.Get("加载配置") + `</button>`)

	submitBtn := components.Button().SetType("submit").
		SetContent(template.HTML(gotext.Get("提交"))).
		SetThemePrimary().
		SetLoadingText(icon.Icon("fa-spinner fa-spin", 2) + `提交中`).
		GetContent()
	resetBtn := components.Button().SetType("reset").
		SetContent(template.HTML(gotext.Get("重置"))).
		SetThemeWarning().
		GetContent()

	// 2. 布局：按钮行
	btnRow := components.Row().SetContent(
		components.Col().SetSize(types.SizeMD(2)).SetContent(loadBtn).GetContent() +
			components.Col().SetSize(types.SizeMD(6)).SetContent(submitBtn+" "+resetBtn).GetContent(),
	).GetContent()

	// 3. 服务器选择下拉框
	serverOptions := fusion.GetServerList()
	defaultServerId := "0"
	if len(serverOptions) > 0 {
		defaultServerId = serverOptions[0].Value
	}
	selectedServerId := ctx.Request.FormValue("serverId")
	if selectedServerId == "" {
		selectedServerId = defaultServerId
	}

	// 4. 读取并解析配置
	var signInConfig DailySignInConfig
	serverId, err := strconv.Atoi(selectedServerId)
	if err == nil {
		signInConfig = loadDailySignInConfig(serverId)
	} else {
		signInConfig = DailySignInConfig{StartTime: 0, EndTime: 0, Period: 0}
		zap.L().Error("解析服务器ID失败", zap.String("serverId", selectedServerId), zap.Error(err))
	}

	// 5. 构建表单
	formList := types.NewFormPanel()
	formList.AddField(gotext.Get("选择服务器"), "serverId", db.Int, form.SelectSingle).
		FieldMust().
		FieldOptions(serverOptions).
		FieldDefault(selectedServerId).
		FieldHelpMsg(template.HTML(gotext.Get("选择后点击【加载配置】读取该服务器的每日签到配置")))

	formList.AddField(gotext.Get("开始时间"), "startTime", db.Int, form.Datetime).
		FieldMust().
		FieldDefault(strconv.FormatInt(signInConfig.StartTime, 10)).
		FieldDisplay(func(model types.FieldModel) interface{} {
			cTime, _ := strconv.ParseInt(model.Value, 10, 64)
			if cTime == 0 {
				return ""
			}
			return time.Unix(cTime, 0).Format("2006-01-02 15:04:05")
		}).
		FieldHelpMsg(template.HTML(gotext.Get("签到活动开始时间")))

	formList.AddField(gotext.Get("结束时间"), "endTime", db.Int, form.Datetime).
		FieldMust().
		FieldDefault(strconv.FormatInt(signInConfig.EndTime, 10)).
		FieldDisplay(func(model types.FieldModel) interface{} {
			cTime, _ := strconv.ParseInt(model.Value, 10, 64)
			if cTime == 0 {
				return ""
			}
			return time.Unix(cTime, 0).Format("2006-01-02 15:04:05")
		}).
		FieldHelpMsg(template.HTML(gotext.Get("签到活动结束时间")))

	formList.AddField(gotext.Get("期数"), "period", db.Int, form.Number).
		FieldMust().
		FieldDefault(strconv.Itoa(signInConfig.Period)).
		FieldHelpMsg(template.HTML(gotext.Get("签到活动期数")))

	formList.SetTabGroups(types.TabGroups{
		{"serverId", "startTime", "endTime", "period"},
	})
	formList.SetTabHeaders(gotext.Get("每日签到配置"))
	field, headers := formList.GroupField()

	// 6. 前端JS
	jsCode := template.HTML(`
	<script>
	document.getElementById('loadConfigBtn').addEventListener('click', function() {
		let serverSel = document.querySelector('select[name="serverId"]');
		let serverId = serverSel.value;
		window.location.href = window.location.pathname + '?serverId=' + serverId;
	});
	</script>
	`)

	// 7. 构建表单UI
	aform := components.Form().
		SetTabHeaders(headers).
		SetTabContents(field).
		SetPrefix(config.PrefixFixSlash()).
		SetUrl("/admin/SubmitDailySignIn").
		SetHiddenFields(map[string]string{
			form2.PreviousKey: "/admin",
		}).
		SetMethod("GET").
		SetOperationFooter(btnRow)

	// 8. 页面内容
	pageContent := components.Box().
		SetHeader(template.HTML(gotext.Get("每日签到设置（DailySignIn）"))).
		WithHeadBorder().
		SetBody(aform.GetContent() + formList.FooterHtml + jsCode).
		GetContent()

	return types.Panel{
		Content:   pageContent,
		Title:     template.HTML(gotext.Get("每日签到设置")),
		CSS:       `.modal.fade.in{z-index:10002}`,
		Callbacks: nil,
	}, nil
}

func SubmitDailySignIn(ctx *context.Context) (types.Panel, error) {
	components := template2.Get(config.GetTheme())
	col1 := components.Col()
	size := map[string]string{"md": "3", "sm": "6", "xs": "12"}
	str := gotext.Get("返回上级")
	preUrl := ctx.Request.Header.Get("Referer")
	if idx := strings.Index(preUrl, "admin"); idx > -1 {
		preUrl = preUrl[idx:]
	} else {
		preUrl = ""
	}

	// 1. 解析服务器ID
	serverIdStr := ctx.Request.FormValue("serverId")
	serverId, err := strconv.Atoi(serverIdStr)
	if err != nil {
		faild := fmt.Sprintf("<h5 style='color:red'>解析服务器ID失败：%s</h5>", err.Error())
		return buildResultPanel(col1, size, faild, preUrl, str), nil
	}

	// 2. 解析表单字段
	startTimeStr := ctx.Request.FormValue("startTime")
	endTimeStr := ctx.Request.FormValue("endTime")
	period, _ := strconv.Atoi(ctx.Request.FormValue("period"))

	var startTime int64
	var endTime int64
	if startTimeStr != "" && startTimeStr != "0" {
		t, _ := time.ParseInLocation("2006-01-02 15:04:05", startTimeStr, time.Local)
		startTime = t.Unix()
	}
	if endTimeStr != "" && endTimeStr != "0" {
		t, _ := time.ParseInLocation("2006-01-02 15:04:05", endTimeStr, time.Local)
		endTime = t.Unix()
	}

	// 3. 打包配置字符串
	packerStr := packDailySignInConfig(startTime, endTime, period)

	// 4. 调用接口提交
	vars := &url.Values{}
	var faild string
	vars.Add("cmd", "DailySignIn")
	vars.Add("gsId", serverIdStr)
	vars.Add("args", packerStr)

	ctx1 := &context.Context{}
	err, res := fusion.CallToCenter(vars, common.MyCfg.Api["GM2GS"], "GET", ctx1)
	if err != nil {
		faild = fmt.Sprintf("<h5 style='color:red'>服务器[%d]：%s</h5>", serverId, err.Error())
	} else if res != "" {
		faild = fmt.Sprintf("<h5 style='color:red'>服务器[%d]：%s</h5>", serverId, res)
	} else {
		faild = fmt.Sprintf("<h5 style='color:green'>服务器[%d]配置提交成功！</h5>", serverId)
	}

	// 5. 返回结果
	return buildResultPanel(col1, size, faild, preUrl, str), nil
}

func loadDailySignInConfig(serverId int) DailySignInConfig {
	gormDB := fusion.GetServerGormDB("db_char", int64(serverId))
	if gormDB == nil {
		zap.L().Error("获取数据库连接失败", zap.Int("serverId", serverId))
		return DailySignInConfig{StartTime: 0, EndTime: 0, Period: 0}
	}

	var cfgValue string
	err := gormDB.Table("inst_configure").Select("cfgValue").Where("cfgID = ?", 138).Scan(&cfgValue).Error
	if err != nil {
		zap.L().Error("读取配置失败", zap.Int("cfgID", 138), zap.Int("serverId", serverId), zap.Error(err))
		return DailySignInConfig{StartTime: 0, EndTime: 0, Period: 0}
	}

	return parseDailySignInConfig(cfgValue)
}

func parseDailySignInConfig(cfgValue string) DailySignInConfig {
	config := DailySignInConfig{StartTime: 0, EndTime: 0, Period: 0}
	if cfgValue == "" {
		return config
	}

	parts := strings.Split(cfgValue, ",")
	if len(parts) >= 1 {
		config.StartTime, _ = strconv.ParseInt(parts[0], 10, 64)
	}
	if len(parts) >= 2 {
		config.EndTime, _ = strconv.ParseInt(parts[1], 10, 64)
	}
	if len(parts) >= 3 {
		config.Period, _ = strconv.Atoi(parts[2])
	}

	return config
}

func packDailySignInConfig(startTime, endTime int64, period int) string {
	var packer strings.Builder
	packer.WriteString(strconv.FormatInt(startTime, 10))
	packer.WriteString(",")
	packer.WriteString(strconv.FormatInt(endTime, 10))
	packer.WriteString(",")
	packer.WriteString(strconv.Itoa(period))
	return packer.String()
}
