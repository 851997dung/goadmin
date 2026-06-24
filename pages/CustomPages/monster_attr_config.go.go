package CustomPages

import (
	"admin/common"
	"admin/fusion"
	"fmt"
	"html/template"
	"net/url"
	"strconv"
	"strings"

	"github.com/GoAdminGroup/go-admin/context"
	"github.com/GoAdminGroup/go-admin/modules/config"
	"github.com/GoAdminGroup/go-admin/modules/db"
	form2 "github.com/GoAdminGroup/go-admin/plugins/admin/modules/form"
	template2 "github.com/GoAdminGroup/go-admin/template"
	"github.com/GoAdminGroup/go-admin/template/icon"
	"github.com/GoAdminGroup/go-admin/template/types"
	"github.com/GoAdminGroup/go-admin/template/types/form"
	"github.com/GoAdminGroup/themes/adminlte/components/infobox"
	"github.com/leonelquinteros/gotext"
	"go.uber.org/zap"
)

// ========== 结构体定义 ==========
type MonsterAttr struct {
	MonsterId int        `json:"monster_id"`
	AttrNum   int        `json:"attr_num"`
	Attrs     []AttrItem `json:"attrs"`
}

type AttrItem struct {
	AttrId  int     `json:"attr_id"`
	AttrVal float64 `json:"attr_val"` // float64对应double类型
	Day     int     `json:"day"`      // 【新增】天数（int）
	IsStop  int     `json:"is_stop"`  // 【新增】是否停止（int，如0=否，1=是）
}

type BossAttrConfig struct {
	ServerDay  int           `json:"server_day"`
	MonsterNum int           `json:"monster_num"`
	Monsters   []MonsterAttr `json:"monsters"`
}

// ========== 核心页面 ==========
func GetMonsterAttrConfigPage(ctx *context.Context) (types.Panel, error) {
	components := template2.Get(config.GetTheme())

	// 1. 按钮组件
	loadBtn := template.HTML(`<button type="button" id="loadConfigBtn" class="btn btn-info">` + gotext.Get("加载配置") + `</button>`)
	addMonsterBtn := template.HTML(`<button type="button" id="addMonsterBtn" class="btn btn-success">` + gotext.Get("新增怪物") + `</button>`)

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
			components.Col().SetSize(types.SizeMD(2)).SetContent(addMonsterBtn).GetContent() +
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
	var bossConfig BossAttrConfig
	serverId, err := strconv.Atoi(selectedServerId)
	if err == nil {
		bossConfig = loadBossAttrConfig(serverId)
	} else {
		bossConfig = BossAttrConfig{ServerDay: 0, MonsterNum: 0, Monsters: []MonsterAttr{}}
		zap.L().Error("解析服务器ID失败", zap.String("serverId", selectedServerId), zap.Error(err))
	}

	// 5. 构建表单
	formList := types.NewFormPanel()
	formList.AddField(gotext.Get("选择服务器"), "serverId", db.Int, form.SelectSingle).
		FieldMust().
		FieldOptions(serverOptions).
		FieldDefault(selectedServerId).
		FieldHelpMsg(template.HTML(gotext.Get("选择后点击【加载配置】读取该服务器的怪物属性配置")))

	formList.AddField(gotext.Get("服务器天数"), "serverDay", db.Int, form.Number).
		FieldMust().
		FieldDefault(strconv.Itoa(bossConfig.ServerDay)).
		FieldHelpMsg(template.HTML(gotext.Get("配置生效的服务器天数（如5=开服第5天）")))

	monsterListHtml := renderMonsterList(bossConfig.Monsters)
	formList.AddField(gotext.Get("怪物配置列表"), "monsterList", db.Text, form.Text).
		FieldDefault(monsterListHtml).
		FieldDisplay(func(model types.FieldModel) interface{} {
			return template.HTML(model.Value)
		}).
		FieldDisplayButCanNotEditWhenUpdate().
		FieldDisplayButCanNotEditWhenCreate()

	formList.SetTabGroups(types.TabGroups{
		{"serverId", "serverDay", "monsterList"},
	})
	formList.SetTabHeaders(gotext.Get("怪物属性配置"))
	field, headers := formList.GroupField()

	// 7. 前端JS
	monsterCount := len(bossConfig.Monsters)
	monsterTpl := getMonsterTemplate(999, MonsterAttr{})
	attrTpl := getAttrTemplate(999, 999, AttrItem{})

	jsCode := template.HTML(`
	<script>
	// 怪物计数
	let monsterCount = ` + strconv.Itoa(monsterCount) + `;
	
	// 加载配置：页面跳转，不提交表单
	document.getElementById('loadConfigBtn').addEventListener('click', function() {
		let serverSel = document.querySelector('select[name="serverId"]');
		let serverId = serverSel.value;
		window.location.href = window.location.pathname + '?serverId=' + serverId;
	});

	// 新增怪物
	document.getElementById('addMonsterBtn').addEventListener('click', function() {
		monsterCount++;
		let monsterContainer = document.getElementById('monsterContainer');
		let monsterHtml = ` + "`" + monsterTpl + "`" + `.replace(/999/g, monsterCount);
		monsterContainer.insertAdjacentHTML('beforeend', monsterHtml);
	});

	// 新增属性
	function addAttr(monsterIdx) {
		let attrContainer = document.getElementById('attrContainer_' + monsterIdx);
		let attrCount = attrContainer.querySelectorAll('.attr-item').length + 1;
		let attrHtml = ` + "`" + attrTpl + "`" + `.replace(/999/g, monsterIdx).replace(/888/g, attrCount);
		attrContainer.insertAdjacentHTML('beforeend', attrHtml);
	}

	// 删除怪物/属性
	function deleteMonster(monsterIdx) {
		let el = document.getElementById('monsterItem_' + monsterIdx);
		if(el) el.remove();
		monsterCount = document.querySelectorAll('.monster-item').length;
	}
	function deleteAttr(monsterIdx, attrIdx) {
		let el = document.getElementById('attrItem_' + monsterIdx + '_' + attrIdx);
		if(el) el.remove();
	}
	</script>
	`)

	// 8. 构建表单UI
	aform := components.Form().
		SetTabHeaders(headers).
		SetTabContents(field).
		SetPrefix(config.PrefixFixSlash()).
		SetUrl("/admin/SubmitMonsterAttrConfig").
		SetHiddenFields(map[string]string{
			form2.PreviousKey: "/admin",
		}).
		SetMethod("GET").
		SetOperationFooter(btnRow)

	// 9. 页面内容
	pageContent := components.Box().
		SetHeader(template.HTML(gotext.Get("怪物属性配置（BossAttrByDay）"))).
		WithHeadBorder().
		SetBody(aform.GetContent() + formList.FooterHtml + jsCode).
		GetContent()

	return types.Panel{
		Content:   pageContent,
		Title:     template.HTML(gotext.Get("怪物属性配置")),
		CSS:       `.modal.fade.in{z-index:10002}; .monster-item {border:1px solid #eee; padding:10px; margin:10px 0; border-radius:5px;} .attr-item {margin:5px 0; padding:5px; border-left:3px solid #666;}`,
		Callbacks: nil,
	}, nil
}

// ========== 提交处理函数 ==========
func SubmitMonsterAttrConfig(ctx *context.Context) (types.Panel, error) {
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

	// 2. 解析服务器天数
	serverDay, _ := strconv.Atoi(ctx.Request.FormValue("serverDay"))

	// 3. 收集动态怪物配置
	var monsters []MonsterAttr
	monsterIdx := 1
	for {
		monsterIdStr := ctx.Request.FormValue(fmt.Sprintf("monster_%d_id", monsterIdx))
		if monsterIdStr == "" {
			break
		}
		monsterId, _ := strconv.Atoi(monsterIdStr)
		if monsterId <= 0 {
			monsterIdx++
			continue
		}

		var attrs []AttrItem
		attrIdx := 1
		for {
			// 【新增】读取天数、是否停止参数（与AttrId/AttrVal逻辑一致）
			attrIdStr := ctx.Request.FormValue(fmt.Sprintf("monster_%d_attr_%d_id", monsterIdx, attrIdx))
			attrValStr := ctx.Request.FormValue(fmt.Sprintf("monster_%d_attr_%d_val", monsterIdx, attrIdx))
			dayStr := ctx.Request.FormValue(fmt.Sprintf("monster_%d_attr_%d_day", monsterIdx, attrIdx))        // 天数
			isStopStr := ctx.Request.FormValue(fmt.Sprintf("monster_%d_attr_%d_is_stop", monsterIdx, attrIdx)) // 是否停止

			if attrIdStr == "" {
				break
			}

			// 解析基础字段
			attrId, _ := strconv.Atoi(attrIdStr)
			attrVal, _ := strconv.ParseFloat(attrValStr, 64)
			// 【新增】解析天数、是否停止（int类型）
			day, _ := strconv.Atoi(dayStr)
			isStop, _ := strconv.Atoi(isStopStr)

			if attrId > 0 {
				attrs = append(attrs, AttrItem{
					AttrId:  attrId,
					AttrVal: attrVal,
					Day:     day,    // 赋值
					IsStop:  isStop, // 赋值
				})
			}
			attrIdx++
		}

		monsters = append(monsters, MonsterAttr{
			MonsterId: monsterId,
			AttrNum:   len(attrs),
			Attrs:     attrs,
		})
		monsterIdx++
	}

	// 4. 打包配置字符串
	packerStr := packBossAttrConfig(serverDay, monsters)

	// 5. 调用接口提交
	vars := &url.Values{}
	var faild string
	vars.Add("cmd", "BossAttrByDay")
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

	// 6. 返回结果
	return buildResultPanel(col1, size, faild, preUrl, str), nil
}

// ========== 辅助函数 ==========
func loadBossAttrConfig(serverId int) BossAttrConfig {
	gormDB := fusion.GetServerGormDB("db_char", int64(serverId))
	if gormDB == nil {
		zap.L().Error("获取数据库连接失败", zap.Int("serverId", serverId))
		return BossAttrConfig{ServerDay: 0, MonsterNum: 0, Monsters: []MonsterAttr{}}
	}

	var cfgValue string
	err := gormDB.Table("inst_configure").Select("cfgValue").Where("cfgID = ?", 137).Scan(&cfgValue).Error
	if err != nil {
		zap.L().Error("读取配置失败", zap.Int("cfgID", 137), zap.Int("serverId", serverId), zap.Error(err))
		return BossAttrConfig{ServerDay: 0, MonsterNum: 0, Monsters: []MonsterAttr{}}
	}

	return parseBossAttrConfig(cfgValue)
}

func parseBossAttrConfig(cfgValue string) BossAttrConfig {
	config := BossAttrConfig{ServerDay: 0, MonsterNum: 0, Monsters: []MonsterAttr{}}
	if cfgValue == "" {
		return config
	}

	parts := strings.Split(cfgValue, ",")
	idx := 0
	lenParts := len(parts)

	if idx < lenParts {
		config.ServerDay, _ = strconv.Atoi(parts[idx])
		idx++
	}

	if idx < lenParts {
		config.MonsterNum, _ = strconv.Atoi(parts[idx])
		idx++
	}

	for i := 0; i < config.MonsterNum && idx < lenParts; i++ {
		monster := MonsterAttr{Attrs: []AttrItem{}}

		if idx < lenParts {
			monster.MonsterId, _ = strconv.Atoi(parts[idx])
			idx++
		}

		if idx < lenParts {
			monster.AttrNum, _ = strconv.Atoi(parts[idx])
			idx++
		}

		// 解析每个属性：原有2个字段→新增为4个字段（AttrId/AttrVal/Day/IsStop）
		for j := 0; j < monster.AttrNum && idx+3 < lenParts; j++ { // 【修改】idx+3（需同时有Id/Val/Day/IsStop）
			attrId, _ := strconv.Atoi(parts[idx])
			attrVal, _ := strconv.ParseFloat(parts[idx+1], 64)
			// 【新增】解析天数、是否停止
			day, _ := strconv.Atoi(parts[idx+2])
			isStop, _ := strconv.Atoi(parts[idx+3])

			monster.Attrs = append(monster.Attrs, AttrItem{
				AttrId:  attrId,
				AttrVal: attrVal,
				Day:     day,    // 赋值
				IsStop:  isStop, // 赋值
			})
			idx += 4 // 【修改】步长从2→4（4个字段：Id/Val/Day/IsStop）
		}

		config.Monsters = append(config.Monsters, monster)
	}

	return config
}

func packBossAttrConfig(serverDay int, monsters []MonsterAttr) string {
	var packer strings.Builder
	packer.WriteString(strconv.Itoa(serverDay))
	packer.WriteString(",")
	packer.WriteString(strconv.Itoa(len(monsters)))
	packer.WriteString(",")

	for _, monster := range monsters {
		packer.WriteString(strconv.Itoa(monster.MonsterId))
		packer.WriteString(",")
		packer.WriteString(strconv.Itoa(len(monster.Attrs)))
		packer.WriteString(",")

		// 打包每个属性：原有2个字段→新增为4个字段
		for _, attr := range monster.Attrs {
			packer.WriteString(strconv.Itoa(attr.AttrId))
			packer.WriteString(",")
			packer.WriteString(fmt.Sprintf("%f", attr.AttrVal))
			packer.WriteString(",")
			// 【新增】打包天数、是否停止（与其他字段逻辑一致）
			packer.WriteString(strconv.Itoa(attr.Day))
			packer.WriteString(",")
			packer.WriteString(strconv.Itoa(attr.IsStop))
			packer.WriteString(",")
		}
	}

	return strings.TrimRight(packer.String(), ",")
}

func renderMonsterList(monsters []MonsterAttr) string {
	if len(monsters) == 0 {
		return `<div id="monsterContainer"><div class="alert alert-info">暂无怪物配置，点击【新增怪物】添加</div></div>`
	}

	var html strings.Builder
	html.WriteString(`<div id="monsterContainer">`)
	for idx, monster := range monsters {
		html.WriteString(getMonsterTemplate(idx+1, monster))
	}
	html.WriteString(`</div>`)
	return html.String()
}

func getMonsterTemplate(monsterIdx int, monster MonsterAttr) string {
	var html strings.Builder
	html.WriteString(fmt.Sprintf(`
	<div class="monster-item" id="monsterItem_%d">
		<h5>怪物 %d</h5>
		<div class="form-group">
			<label>怪物ID：</label>
			<input type="number" name="monster_%d_id" value="%d" class="form-control" style="width:200px; display:inline-block;" required>
			<button type="button" class="btn btn-sm btn-success" onclick="addAttr(%d)">新增属性</button>
			<button type="button" class="btn btn-sm btn-danger" onclick="deleteMonster(%d)">删除</button>
		</div>
		<div id="attrContainer_%d" class="attr-container" style="margin-left:20px;">`,
		monsterIdx, monsterIdx, monsterIdx, monster.MonsterId, monsterIdx, monsterIdx, monsterIdx))

	if len(monster.Attrs) == 0 {
		html.WriteString(getAttrTemplate(monsterIdx, 1, AttrItem{}))
	} else {
		for attrIdx, attr := range monster.Attrs {
			html.WriteString(getAttrTemplate(monsterIdx, attrIdx+1, attr))
		}
	}

	html.WriteString(`
		</div>
		<hr>
	</div>`)
	return html.String()
}

func getAttrTemplate(monsterIdx, attrIdx int, attr AttrItem) string {
	attrIdxPlaceholder := 888
	if attrIdx == 999 {
		attrIdxPlaceholder = 888
	} else {
		attrIdxPlaceholder = attrIdx
	}
	return fmt.Sprintf(`
	<div class="attr-item" id="attrItem_%d_%d">
		<label>属性%d ID：</label>
		<input type="number" name="monster_%d_attr_%d_id" value="%d" class="form-control" style="width:150px; display:inline-block;">
		
		<label style="margin-left:10px;">属性值：</label>
		<input type="number" step="any" name="monster_%d_attr_%d_val" value="%f" class="form-control" style="width:150px; display:inline-block;">
		
		<!-- 【新增】天数输入框（与其他字段样式/命名规则一致） -->
		<label style="margin-left:10px;">天数：</label>
		<input type="number" name="monster_%d_attr_%d_day" value="%d" class="form-control" style="width:150px; display:inline-block;">
		
		<!-- 【新增】是否停止输入框（与其他字段样式/命名规则一致） -->
		<label style="margin-left:10px;">是否停止：</label>
		<input type="number" name="monster_%d_attr_%d_is_stop" value="%d" class="form-control" style="width:150px; display:inline-block;" placeholder="0=否,1=是">
		
		<button type="button" class="btn btn-sm btn-danger" onclick="deleteAttr(%d, %d)">删除</button>
	</div>`,
		// 原有字段占位符
		monsterIdx, attrIdxPlaceholder, attrIdxPlaceholder,
		monsterIdx, attrIdxPlaceholder, attr.AttrId,
		monsterIdx, attrIdxPlaceholder, attr.AttrVal,
		// 【新增】天数占位符+值
		monsterIdx, attrIdxPlaceholder, attr.Day,
		// 【新增】是否停止占位符+值
		monsterIdx, attrIdxPlaceholder, attr.IsStop,
		// 删除按钮参数
		monsterIdx, attrIdxPlaceholder)
}

func buildResultPanel(col1 types.ColAttribute, size map[string]string, faild, preUrl, str string) types.Panel {
	infoboxComp := infobox.New()
	infoboxCol1 := col1.SetSize(size).SetContent(infoboxComp.SetText(template.HTML(faild)).GetContent()).GetContent()
	return types.Panel{
		Content:     infoboxCol1,
		Title:       template.HTML(gotext.Get("提交结果")),
		Description: template.HTML(fmt.Sprintf(`<a href="/%s">%s</a>`, preUrl, str)),
		Callbacks:   nil,
	}
}
