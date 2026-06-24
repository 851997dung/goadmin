package CustomPages

import (
	"admin/common"
	"admin/fusion"
	"admin/hotfix"
	"fmt"
	"html/template"
	"net/url"
	"strconv"
	"strings"

	"github.com/GoAdminGroup/go-admin/context"
	"github.com/GoAdminGroup/go-admin/modules/db"
	"github.com/GoAdminGroup/go-admin/plugins/admin/modules/parameter"
	"github.com/GoAdminGroup/go-admin/plugins/admin/modules/table"
	"github.com/GoAdminGroup/go-admin/template/types"
	"github.com/leonelquinteros/gotext"
)

func GetHotfixDataTablesPage(ctx *context.Context) (types.Panel, error) {
	content := template.HTML(hotfixDataTablesHeader()) + executeHotfixGM(ctx)
	return types.Panel{
		Content: template.HTML(content),
		Title:   template.HTML(gotext.Get("热更数据表")),
	}, nil
}

func hotfixDataTablesHeader() string {
	return `<div class="panel panel-default">
<div class="panel-heading">` + gotext.Get("热更数据表") + `</div>
<div class="panel-body">
<form method="post" action="/admin/hotfix/start">
<div class="form-group">
<label>` + gotext.Get("目标服务器") + `</label>
<select class="form-control" name="gsIds" multiple>` +
		hotfixServerOptions() +
		`</select>
</div>
<div class="form-group">
<label>` + gotext.Get("数据表名") + `</label>
<input type="text" class="form-control" name="tableNames" placeholder="` + gotext.Get("表名以逗号分隔") + `">
</div>
<button type="submit" class="btn btn-primary">` + gotext.Get("开始热更") + `</button>
</form>
<hr>
<h5>` + gotext.Get("特殊表名说明") + `</h5>
<table class="table table-bordered table-striped">
<thead><tr><th>` + gotext.Get("表名") + `</th><th>` + gotext.Get("触发的 GM 命令") + `</th></tr></thead>
<tbody>
<tr><td><code>$Spell</code></td><td>MapServer HotfixSpellRelation + ClearSpellEffectArgs</td></tr>
<tr><td><code>$Loot</code></td><td>MapServer HotfixLootRelation</td></tr>
<tr><td><code>QuestPrototype</code></td><td>MapServer HotfixQuests</td></tr>
<tr><td><code>ItemPrototype</code></td><td>MapServer HotfixItemPrototypes</td></tr>
<tr><td><code>MapSpecial</code></td><td>MapServer HotfixMapSpecial</td></tr>
<tr><td><code>operating_activities</code></td><td>SocialServer ReloadActivity</td></tr>
</tbody>
</table>
</div>`
}

func hotfixServerOptions() string {
	var opts string
	for _, s := range fusion.GetServerList() {
		opts += fmt.Sprintf(`<option value="%s">%s</option>`, s.Value, s.Text)
	}
	return opts
}

func executeHotfixGM(ctx *context.Context) template.HTML {
	gsIDStrs := strings.Split(ctx.Request.FormValue("gsIds"), ",")
	if len(gsIDStrs) == 0 || gsIDStrs[0] == "" {
		return template.HTML(`<div class="alert alert-danger">` + gotext.Get("请选择服务器") + `</div>`)
	}

	var gsIDs []uint32
	for _, s := range gsIDStrs {
		id, err := strconv.ParseUint(strings.TrimSpace(s), 10, 32)
		if err != nil {
			return template.HTML(fmt.Sprintf(`<div class="alert alert-danger">invalid gsId: %s</div>`, s))
		}
		gsIDs = append(gsIDs, uint32(id))
	}

	tableNamesStr := ctx.Request.FormValue("tableNames")
	if tableNamesStr == "" {
		return template.HTML(`<div class="alert alert-danger">` + gotext.Get("请输入数据表名") + `</div>`)
	}
	tableNames := strings.FieldsFunc(tableNamesStr, func(r rune) bool { return r == ',' || r == ' ' || r == '\t' || r == '\n' || r == '\r' })

	// Resolve special table commands
	tableNamesForGM, specialCmds := resolveSpecialTableCmdsPage(tableNames)

	var rows []string
	for _, gsID := range gsIDs {

		// Step 1: AdminServer HotfixDBFile
		row, _ := doGMCommandPage(gsID, "AdminServer", "HotfixDBFile", "")
		rows = append(rows, row)

		// Step 3: AdminServer HotfixDBTables
		row, _ = doGMCommandPage(gsID, "AdminServer", "HotfixDBTables", strings.Join(tableNamesForGM, ","))
		rows = append(rows, row)

		// Steps 4-10: Conditional commands
		for _, cmd := range specialCmds {
			row, _ = doGMCommandPage(gsID, cmd.server, cmd.cmd, "")
			rows = append(rows, row)
		}

		rows = append(rows, fmt.Sprintf(`<tr class="success"><td>%d</td><td colspan="2">%s</td></tr>`, gsID, gotext.Get("操作完成")))
	}

	header := `<h4>` + gotext.Get("操作结果") + `</h4>`
	table := `<table class="table table-hover"><thead><tr><th>GS ID</th><th>Command</th><th>Result</th></tr></thead><tbody>` +
		strings.Join(rows, "") + `</tbody></table>`
	return template.HTML(header + table)
}

type specialCmdPage struct {
	server string
	cmd    string
}

func resolveSpecialTableCmdsPage(tableNames []string) ([]string, []specialCmdPage) {
	var cmds []specialCmdPage
	remaining := make([]string, 0, len(tableNames))

	hasSpell := false
	hasLoot := false
	for _, name := range tableNames {
		if name == "$Spell" || hotfix.IsSpellTable(name) {
			hasSpell = true
		} else if name == "$Loot" || hotfix.IsLootTable(name) {
			hasLoot = true
		} else {
			remaining = append(remaining, name)
		}
	}
	if hasSpell {
		cmds = append(cmds, specialCmdPage{"MapServer", "HotfixSpellRelation"})
		cmds = append(cmds, specialCmdPage{"MapServer", "ClearSpellEffectArgs"})
	}
	if hasLoot {
		cmds = append(cmds, specialCmdPage{"MapServer", "HotfixLootRelation"})
	}

	specialMap := map[string]specialCmdPage{
		"QuestPrototype":       {"MapServer", "HotfixQuests"},
		"ItemPrototype":        {"MapServer", "HotfixItemPrototypes"},
		"MapSpecial":           {"MapServer", "HotfixMapSpecial"},
		"operating_activities": {"SocialServer", "ReloadActivity"},
	}
	for i := 0; i < len(remaining); i++ {
		if cmd, ok := specialMap[remaining[i]]; ok {
			cmds = append(cmds, cmd)
			remaining = append(remaining[:i], remaining[i+1:]...)
			i--
		}
	}
	return remaining, cmds
}

func doGMCommandPage(gsID uint32, server, cmd, args string) (string, bool) {
	apis := []string{
		common.MyCfg.Api["GM2S"],
		common.MyCfg.Api["GM2MS"],
		common.MyCfg.Api["GM2GS"],
		common.MyCfg.Api["GM2GATE"],
		common.MyCfg.Api["GM2SOCIAL"],
		common.MyCfg.Api["GM2DBP"],
	}
	var result string
	for _, api := range apis {
		err, res := fusion.CallToCenter(&url.Values{
			"name": {server},
			"cmd":  {cmd},
			"args": {args},
			"gsId": {strconv.FormatUint(uint64(gsID), 10)},
		}, api, "GET", &context.Context{})
		if err != nil {
			result = fmt.Sprintf(`<td class="danger">%s</td>`, err.Error())
			break
		}
		result = fmt.Sprintf(`<td>%s</td>`, res)
	}
	cell := fmt.Sprintf(`<tr><td>%d</td><td>%s %s</td>%s</tr>`, gsID, server, cmd, result)
	return cell, !strings.Contains(result, "danger")
}

func GetHotfixDataTablesTable(ctx *context.Context) table.Table {
	tbl := table.NewDefaultTable(table.Config{
		Driver:     "mysql",
		Connection: "default",
		CanAdd:     false,
		Editable:   false,
		Deletable:  false,
		Exportable: false,
		PrimaryKey: table.PrimaryKey{Type: db.Int, Name: "id"},
	})
	info := tbl.GetInfo()
	info.AddField("ID", "id", db.Int)
	info.AddField(gotext.Get("表名"), "table_name", db.Text)
	info.AddField(gotext.Get("状态"), "status", db.Text)
	info.SetTable("hotfix_log").SetTitle(gotext.Get("热更记录"))
	info.HideNewButton().HideEditButton().HideDeleteButton().HideDetailButton().HideRowSelector().HideExportButton()
	info.SetGetDataFn(func(param parameter.Parameters) (data []map[string]interface{}, size int) {
		// empty table for now
		return nil, 0
	})
	return tbl
}

func StartHotfixPipeline(ctx *context.Context) {
	gsIDStrs := strings.Split(ctx.Request.FormValue("gsIds"), ",")
	var gsIDs []uint32
	for _, s := range gsIDStrs {
		id, _ := strconv.ParseUint(strings.TrimSpace(s), 10, 32)
		gsIDs = append(gsIDs, uint32(id))
	}

	tableNamesStr := ctx.Request.FormValue("tableNames")
	tableNames := strings.FieldsFunc(tableNamesStr, func(r rune) bool { return r == ',' || r == ' ' || r == '\t' || r == '\n' || r == '\r' })

	cfg := hotfix.GetConfig()
	session, err := hotfix.StartSession("hotfix_datatables")
	if err != nil {
			ctx.JSON(200, map[string]interface{}{"status": "error", "message": err.Error()})
			return
		}
		go hotfix.RunFullPipeline(session, cfg, gsIDs, tableNames, hotfix.SkipFlags{IsSkipBuildTool: true, IsSkipConfigs: false})

	ctx.JSON(200, map[string]interface{}{
		"sessionId": session.ID,
		"status":    "started",
	})
}

func GetHotfixProgress(ctx *context.Context) {
	sessionID := ctx.Request.FormValue("sessionId")
	session := hotfix.GetSession(sessionID)
	if session == nil {
		ctx.JSON(404, map[string]interface{}{"error": "session not found"})
		return
	}
	ctx.JSON(200, map[string]interface{}{
		"status":   session.Status,
		"messages": session.Messages,
	})
}

func StopHotfixPipeline(ctx *context.Context) {
	sessionID := ctx.Request.FormValue("sessionId")
	session := hotfix.GetSession(sessionID)
	if session == nil {
		ctx.JSON(404, map[string]interface{}{"error": "session not found"})
		return
	}
	session.Cancel()
	ctx.JSON(200, map[string]interface{}{"status": "cancelled"})
}
