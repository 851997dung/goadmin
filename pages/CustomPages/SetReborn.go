package CustomPages

import (
	"admin/common"
	"admin/fusion"
	"fmt"
	"html/template"
	"net/url"
	"strings"

	"github.com/GoAdminGroup/go-admin/context"
	"github.com/GoAdminGroup/go-admin/modules/config"
	"github.com/GoAdminGroup/go-admin/modules/db"
	"github.com/GoAdminGroup/go-admin/modules/language"
	form2 "github.com/GoAdminGroup/go-admin/plugins/admin/modules/form"
	template2 "github.com/GoAdminGroup/go-admin/template"
	"github.com/GoAdminGroup/go-admin/template/icon"
	"github.com/GoAdminGroup/go-admin/template/types"
	"github.com/GoAdminGroup/go-admin/template/types/form"
	"github.com/GoAdminGroup/themes/adminlte/components/infobox"
	"github.com/leonelquinteros/gotext"
)

// GetSetRebornPage 返回设置转生页面
func GetSetRebornPage(ctx *context.Context) (types.Panel, error) {
	components := template2.Get(config.GetTheme())

	// 构建表单
	formList := buildRebornForm()

	// 获取表单字段和标签头
	fields, headers := formList.GroupField()

	// 构建操作按钮
	buttons := buildFormButtons(components)

	// 构建完整表单
	aform := components.Form().
		SetTabHeaders(headers).
		SetTabContents(fields).
		SetPrefix(config.PrefixFixSlash()).
		SetUrl("/admin/SubmitRebornNum").
		SetHiddenFields(map[string]string{
			form2.PreviousKey: "/admin",
		}).
		SetMethod("GET").
		SetOperationFooter(template.HTML(buttons))

	// 返回面板
	return types.Panel{
		Content: components.Box().
			SetHeader(aform.GetDefaultBoxHeader(true)).
			WithHeadBorder().
			SetHeader(template.HTML(gotext.Get("设置转生"))).
			SetBody(aform.GetContent() + formList.FooterHtml).
			GetContent(),
		Title:     template.HTML(gotext.Get("设置转生")),
		CSS:       `.modal.fade.in{z-index:10002}`,
		Callbacks: nil,
	}, nil
}

// buildRebornForm 构建转生表单
func buildRebornForm() *types.FormPanel {
	formList := types.NewFormPanel()

	// 修复：移除 FieldHide() 让 RebornNum 显示
	formList.AddField(gotext.Get("RebornNum"), "Reborn3_1", db.Int, form.Number)
	// FieldHelpMsg 需要 HTML 类型，这里直接使用字符串会报错
	// 可以在前端通过其他方式添加帮助信息

	formList.AddField(gotext.Get("gsid"), "gsid", db.Int, form.Text)

	formList.SetTabHeaders(gotext.Get("设置"))
	formList.SetTabGroups(types.TabGroups{
		{"Reborn3_1", "gsid"},
	})

	return formList
}

// buildFormButtons 构建表单按钮
func buildFormButtons(components template2.Template) string {
	submitBtn := components.Button().
		SetType("submit").
		SetContent(language.GetFromHtml("submit")).
		SetThemePrimary().
		SetOrientationRight().
		SetLoadingText(icon.Icon("fa-spinner fa-spin", 2) + `Save`).
		GetContent()

	resetBtn := components.Button().
		SetType("reset").
		SetContent(language.GetFromHtml("Reset")).
		SetThemeWarning().
		SetOrientationLeft().
		GetContent()

	col1 := components.Col().GetContent()
	col2 := components.Col().SetSize(types.SizeMD(8)).
		SetContent(submitBtn + resetBtn).GetContent()

	return string(col1 + col2)
}

// SubmitRebornNum 处理转生设置提交
func SubmitRebornNum(ctx *context.Context) (types.Panel, error) {
	components := template2.Get(config.GetTheme())

	// 获取表单数据
	rebornNum := ctx.Request.FormValue("Reborn3_1")
	gsid := ctx.Request.FormValue("gsid")

	// 验证输入
	if rebornNum == "" || gsid == "" {
		return showErrorResult(components, ctx, gotext.Get("转生次数和服务器ID不能为空"))
	}

	// 调用 API
	result, err := callRebornAPI(rebornNum, gsid)
	if err != nil {
		return showErrorResult(components, ctx, err.Error())
	}

	return showSuccessResult(components, ctx, result)
}

// callRebornAPI 调用转生API
func callRebornAPI(rebornNum, gsid string) (string, error) {
	vars := &url.Values{}
	vars.Add("cmd", "AddMaxRebornCount")
	vars.Add("gsId", gsid)
	vars.Add("args", rebornNum)

	ctx := &context.Context{}
	err, res := fusion.CallToCenter(vars, common.MyCfg.Api["GM2GS"], "GET", ctx)

	if err != nil {
		return "", fmt.Errorf("%s: %v", common.MyCfg.Api["MainServer"], err)
	}

	if res != "" {
		return res, nil
	}

	return gotext.Get("设置成功"), nil
}

// showSuccessResult 显示成功结果
func showSuccessResult(components template2.Template, ctx *context.Context, message string) (types.Panel, error) {
	infobox1 := infobox.New().SetText(template.HTML("<h5 class='text-success'>" + message + "</h5>"))
	col := components.Col().SetSize(map[string]string{"md": "3", "sm": "6", "xs": "12"}).
		SetContent(infobox1.GetContent())

	return types.Panel{
		Content:     col.GetContent(),
		Title:       template.HTML(gotext.Get("提交结果")),
		Description: buildBackLink(ctx),
		Callbacks:   nil,
	}, nil
}

// showErrorResult 显示错误结果
func showErrorResult(components template2.Template, ctx *context.Context, errorMsg string) (types.Panel, error) {
	infobox1 := infobox.New().SetText(template.HTML("<h5 class='text-danger'>" + errorMsg + "</h5>"))
	col := components.Col().SetSize(map[string]string{"md": "3", "sm": "6", "xs": "12"}).
		SetContent(infobox1.GetContent())

	return types.Panel{
		Content:     col.GetContent(),
		Title:       template.HTML(gotext.Get("提交结果")),
		Description: buildBackLink(ctx),
		Callbacks:   nil,
	}, nil
}

// buildBackLink 构建返回链接
func buildBackLink(ctx *context.Context) template.HTML {
	preUrl := ctx.Request.Header.Get("Referer")
	if preUrl == "" {
		preUrl = "/admin"
	} else {
		index := strings.Index(preUrl, "admin")
		if index >= 0 {
			preUrl = preUrl[index:]
		}
	}

	str := gotext.Get("返回上级")
	return template.HTML(fmt.Sprintf(`<a href="/%s" class="btn btn-default">%s</a>`, preUrl, str))
}
