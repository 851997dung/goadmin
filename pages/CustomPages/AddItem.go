package CustomPages

import (
	"admin/common"
	"admin/fusion"
	"fmt"
	"html/template"
	"log"
	"net/url"
	"strconv"
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

// GetAddItemPage 返回添加物品页面
func GetAddItemPage(ctx *context.Context) (types.Panel, error) {
	components := template2.Get(config.GetTheme())

	// 构建表单
	formList := buildItemForm()

	// 获取表单字段和标签头
	fields, headers := formList.GroupField()

	// 构建操作按钮 - 修复：使用正确的函数名
	buttons := buildItemFormButtonsAddItem(components)

	// 构建完整表单
	aform := components.Form().
		SetTabHeaders(headers).
		SetTabContents(fields).
		SetPrefix(config.PrefixFixSlash()).
		SetUrl("/admin/AddItem").
		SetHiddenFields(map[string]string{
			form2.PreviousKey: "/admin",
		}).
		SetMethod("POST"). // 改为POST方法以支持更多数据
		SetOperationFooter(template.HTML(buttons))

	// 返回面板
	return types.Panel{
		Content: components.Box().
			SetHeader(aform.GetDefaultBoxHeader(true)).
			WithHeadBorder().
			SetHeader(template.HTML(gotext.Get("添加物品"))).
			SetBody(aform.GetContent() + formList.FooterHtml).
			GetContent(),
		Title:     template.HTML(gotext.Get("添加物品")),
		CSS:       `.modal.fade.in{z-index:10002}`,
		Callbacks: nil,
	}, nil
}

func buildItemForm() *types.FormPanel {
	formList := types.NewFormPanel()

	formList.AddField(gotext.Get("playerId"), "playerId", db.Int, form.Number).FieldMust()
	formList.AddField(gotext.Get("itemTypeID"), "itemTypeID", db.Int, form.Number).FieldMust()
	formList.AddField(gotext.Get("itemCount"), "itemCount", db.Int, form.Number).FieldMust()
	formList.AddField(gotext.Get("itemFlags"), "itemFlags", db.Int, form.Number).FieldDefault("0")
	formList.AddField(gotext.Get("itemForgeLevel"), "itemForgeLevel", db.Int, form.Number).FieldDefault("0")
	formList.AddField(gotext.Get("itemExpireTime"), "itemExpireTime", db.Int, form.Number).FieldDefault("0")
	formList.AddField(gotext.Get("exattrIdNum"), "exattrIdNum", db.Int, form.Number).FieldDefault("0")
	formList.AddField(gotext.Get("gsid"), "gsid", db.Int, form.Text).FieldMust()

	// 添加额外属性对字段
	formList.AddField(gotext.Get("额外属性组数"), "extraAttrCount", db.Int, form.Number).
		FieldDefault("0").
		FieldHelpMsg(template.HTML("输入要添加的额外属性组数（每组包含index和lv两个参数）"))

	// 动态添加额外属性字段的JavaScript
	formList.FooterHtml += template.HTML(`
<script>
document.addEventListener('DOMContentLoaded', function() {
    const extraAttrCountField = document.querySelector('[name="extraAttrCount"]');
    const formContainer = document.querySelector('.box-body');
    
    if (extraAttrCountField && formContainer) {
        extraAttrCountField.addEventListener('change', updateExtraAttrFields);
        // 初始化
        updateExtraAttrFields();
    }
    
    function updateExtraAttrFields() {
        const count = parseInt(extraAttrCountField.value) || 0;
        
        // 移除现有的额外属性字段
        const existingFields = document.querySelectorAll('.extra-attr-field');
        existingFields.forEach(function(field) {
            field.remove();
        });
        
        // 添加新的额外属性字段
        for (let i = 0; i < count; i++) {
            const groupDiv = document.createElement('div');
            groupDiv.className = 'form-group extra-attr-field';
            
            groupDiv.innerHTML = 
                '<label for="extraIndex' + i + '" class="col-sm-2 control-label">额外属性' + (i+1) + '</label>' +
                '<div class="col-sm-4">' +
                '<div class="input-group">' +
                '<span class="input-group-addon">Index</span>' +
                '<input type="number" class="form-control" name="extraIndex' + i + '" placeholder="索引">' +
                '</div>' +
                '</div>' +
                '<div class="col-sm-4">' +
                '<div class="input-group">' +
                '<span class="input-group-addon">Level</span>' +
                '<input type="number" class="form-control" name="extraLv' + i + '" placeholder="等级">' +
                '</div>' +
                '</div>';
            
            // 在extraAttrCount字段后插入
            extraAttrCountField.closest('.form-group').after(groupDiv);
        }
    }
});
</script>
`)

	formList.SetTabHeaders(gotext.Get("设置"))
	formList.SetTabGroups(types.TabGroups{
		{"playerId", "itemTypeID", "itemCount", "itemFlags", "itemForgeLevel", "itemExpireTime", "exattrIdNum", "gsid", "extraAttrCount"},
	})

	return formList
}

// 修复：使用正确的函数名
func buildItemFormButtonsAddItem(components template2.Template) string {
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

func SubmitAddItem(ctx *context.Context) (types.Panel, error) {
	components := template2.Get(config.GetTheme())

	// 获取基础表单数据
	playerId := ctx.Request.FormValue("playerId")
	itemTypeID := ctx.Request.FormValue("itemTypeID")
	itemCount := ctx.Request.FormValue("itemCount")
	itemFlags := ctx.Request.FormValue("itemFlags")
	itemForgeLevel := ctx.Request.FormValue("itemForgeLevel")
	itemExpireTime := ctx.Request.FormValue("itemExpireTime")
	exattrIdNum := ctx.Request.FormValue("exattrIdNum")
	gsid := ctx.Request.FormValue("gsid")
	extraAttrCount := ctx.Request.FormValue("extraAttrCount")

	// 验证输入
	if playerId == "" || itemTypeID == "" || itemCount == "" || gsid == "" {
		return showItemErrorResult(components, ctx, gotext.Get("必填字段不能为空"))
	}

	// 构建参数
	args, err := buildItemArgs(playerId, itemTypeID, itemCount, itemFlags, itemForgeLevel,
		itemExpireTime, exattrIdNum, extraAttrCount, ctx.Request.Form)
	if err != nil {
		return showItemErrorResult(components, ctx, err.Error())
	}

	// 修复：使用正确的函数名
	result, err := callAddItemAPIForItem(args, gsid)
	if err != nil {
		return showItemErrorResult(components, ctx, err.Error())
	}

	return showItemSuccessResult(components, ctx, result)
}

// buildItemArgs 构建物品参数
func buildItemArgs(playerId, itemTypeID, itemCount, itemFlags, itemForgeLevel,
	itemExpireTime, exattrIdNum, extraAttrCount string, form url.Values) (string, error) {

	// 基础参数
	args := []string{
		playerId,
		itemTypeID,
		itemCount,
		itemFlags,
		itemForgeLevel,
		itemExpireTime,
		exattrIdNum,
	}

	// 处理额外属性
	count, err := strconv.Atoi(extraAttrCount)
	if err != nil {
		count = 0
	}

	// 添加额外属性组数
	args = append(args, strconv.Itoa(count))

	// 添加额外属性对
	for i := 0; i < count; i++ {
		indexKey := fmt.Sprintf("extraIndex%d", i)
		lvKey := fmt.Sprintf("extraLv%d", i)

		index := form.Get(indexKey)
		lv := form.Get(lvKey)

		if index == "" || lv == "" {
			return "", fmt.Errorf("额外属性组 %d 的 index 或 lv 不能为空", i+1)
		}

		args = append(args, index, lv)
	}

	// 用逗号分隔所有参数
	return strings.Join(args, ","), nil
}

// 修复：使用正确的函数名
func callAddItemAPIForItem(args, gsid string) (string, error) {
	vars := &url.Values{}
	vars.Add("cmd", "AddItemToPlayer")
	vars.Add("gsId", gsid)
	vars.Add("args", args)
	log.Println(args)
	ctx := &context.Context{}
	err, res := fusion.CallToCenter(vars, common.MyCfg.Api["GM2GS"], "GET", ctx)

	if err != nil {
		return "", fmt.Errorf("%s: %v", common.MyCfg.Api["MainServer"], err)
	}

	if res != "" {
		return res, nil
	}

	return gotext.Get("添加物品成功"), nil
}

// showItemSuccessResult 显示成功结果
func showItemSuccessResult(components template2.Template, ctx *context.Context, message string) (types.Panel, error) {
	infobox1 := infobox.New().SetText(template.HTML("<h5 class='text-success'>" + message + "</h5>"))
	col := components.Col().SetSize(map[string]string{"md": "3", "sm": "6", "xs": "12"}).
		SetContent(infobox1.GetContent())

	return types.Panel{
		Content:     col.GetContent(),
		Title:       template.HTML(gotext.Get("提交结果")),
		Description: buildItemBackLink(ctx),
		Callbacks:   nil,
	}, nil
}

// showItemErrorResult 显示错误结果
func showItemErrorResult(components template2.Template, ctx *context.Context, errorMsg string) (types.Panel, error) {
	infobox1 := infobox.New().SetText(template.HTML("<h5 class='text-danger'>" + errorMsg + "</h5>"))
	col := components.Col().SetSize(map[string]string{"md": "3", "sm": "6", "xs": "12"}).
		SetContent(infobox1.GetContent())

	return types.Panel{
		Content:     col.GetContent(),
		Title:       template.HTML(gotext.Get("提交结果")),
		Description: buildItemBackLink(ctx),
		Callbacks:   nil,
	}, nil
}

// buildItemBackLink 构建返回链接
func buildItemBackLink(ctx *context.Context) template.HTML {
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
