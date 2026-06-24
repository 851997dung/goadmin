package CustomPages

import (
	"admin/common"
	"admin/fusion"
	"fmt"
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
	"html/template"
	"math"
	"net/url"
	"sort"
	"strconv"
	"strings"
)

type GuildSignInfo struct {
	GuildId    int
	GuildPoint int
	GuildName  string
	LogTime    uint64
}

func GetBattleGroupingPage(ctx *context.Context) (types.Panel, error) {
	components := template2.Get(config.GetTheme())
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

	var guildOpts = GetGuildOpts()
	var formList = types.NewFormPanel()
	formList.AddField(gotext.Get("第一组："), "group1", db.Int, form.Text).
		FieldDisplayButCanNotEditWhenUpdate().FieldDisplayButCanNotEditWhenCreate()
	formList.AddField(gotext.Get("工会1"), "Guild1_1", db.Int, form.SelectSingle).FieldOptions(guildOpts)
	formList.AddField(gotext.Get("工会2"), "Guild1_2", db.Int, form.SelectSingle).FieldOptions(guildOpts)

	formList.AddField(gotext.Get("第二组："), "group2", db.Int, form.Text).
		FieldDisplayButCanNotEditWhenUpdate().FieldDisplayButCanNotEditWhenCreate()
	formList.AddField(gotext.Get("工会1"), "Guild2_1", db.Int, form.SelectSingle).FieldOptions(guildOpts)
	formList.AddField(gotext.Get("工会2"), "Guild2_2", db.Int, form.SelectSingle).FieldOptions(guildOpts)

	formList.AddField(gotext.Get("第三组："), "group3", db.Int, form.Text).
		FieldDisplayButCanNotEditWhenUpdate().FieldDisplayButCanNotEditWhenCreate()
	formList.AddField(gotext.Get("工会1"), "Guild3_1", db.Int, form.SelectSingle).FieldOptions(guildOpts)
	formList.AddField(gotext.Get("工会2"), "Guild3_2", db.Int, form.SelectSingle).FieldOptions(guildOpts)

	formList.AddField(gotext.Get("第四组："), "group4", db.Int, form.Text).
		FieldDisplayButCanNotEditWhenUpdate().FieldDisplayButCanNotEditWhenCreate()
	formList.AddField(gotext.Get("工会1"), "Guild4_1", db.Int, form.SelectSingle).FieldOptions(guildOpts)
	formList.AddField(gotext.Get("工会2"), "Guild4_2", db.Int, form.SelectSingle).FieldOptions(guildOpts)

	formList.SetTabGroups(types.TabGroups{
		{"group1", "Guild1_1", "Guild1_2", "group2", "Guild2_1", "Guild2_2", "group3", "Guild3_1", "Guild3_2", "group4", "Guild4_1", "Guild4_2"},
	})
	formList.SetTabHeaders(gotext.Get("战盟争霸分组"))

	field, headers := formList.GroupField()

	aform := components.Form().
		SetTabHeaders(headers).
		SetTabContents(field).
		SetPrefix(config.PrefixFixSlash()).
		SetUrl("/admin/SubmitBattleInfo").
		SetHiddenFields(map[string]string{
			form2.PreviousKey: "/admin",
		}).
		SetMethod("GET").
		SetOperationFooter(col1 + col2)

	return types.Panel{
		Content: components.Box().
			SetHeader(aform.GetDefaultBoxHeader(true)).
			WithHeadBorder().
			SetHeader(template.HTML(gotext.Get("战盟争霸分组"))).
			SetBody(aform.GetContent() + formList.FooterHtml).
			GetContent(),
		Title:     template.HTML(gotext.Get("战盟争霸分组")),
		CSS:       `.modal.fade.in{z-index:10002}`,
		Callbacks: nil,
	}, nil
}

// 通过两重循环过滤重复元素
func RemoveRepByLoop(slc []GuildSignInfo) []GuildSignInfo {
	result := []GuildSignInfo{} // 存放结果
	for i := range slc {
		flag := true
		for j := range result {
			if slc[i].GuildId == result[j].GuildId {
				flag = false // 存在重复元素，标识为false
				break
			}
		}
		if flag { // 标识为false，不添加进结果
			result = append(result, slc[i])
		}
	}
	return result
}

func GetSignGuilds() []GuildSignInfo {
	var serverIdList = fusion.GetActiveServerIds()
	//var guildInfos4Server = make(map[int][]GuildSignInfo, 0)
	var guildInfos = make([]GuildSignInfo, 0)
	for _, v := range serverIdList {
		var db_char = fusion.GetServerGormDB("db_char", int64(v))
		if db_char == nil {
			continue
		}
		var res string
		db_char.Table("inst_configure").Select("cfgValue").Where("cfgID = ?", 33).Scan(&res)

		var unPacker = fusion.NewTextUnpacker(res)
		var tempSize = unPacker.UnpackUint()
		for i := 0; i < int(tempSize)*4; i++ {
			unPacker.UnpackUint()
		}
		tempSize = unPacker.UnpackUint()
		for i := 0; i < int(tempSize); i++ {
			var info = GuildSignInfo{}
			info.GuildId = int(unPacker.UnpackUint())
			info.GuildPoint = int(unPacker.UnpackInt())
			info.GuildName = unPacker.UnpackString()
			info.LogTime = unPacker.UnpackUint()
			guildInfos = append(guildInfos, info)
		}

		guildInfos = RemoveRepByLoop(guildInfos)

		sort.Slice(guildInfos, func(i, j int) bool {
			if guildInfos[i].GuildPoint > guildInfos[j].GuildPoint {
				return true
			} else if guildInfos[i].GuildPoint == guildInfos[j].GuildPoint {
				if guildInfos[i].LogTime < guildInfos[j].LogTime {
					return true
				}
			}
			return false
		})
	}
	var endMark = math.Min(float64(len(guildInfos)), 8)
	return guildInfos[:int(endMark)]
}

func GetGuildOpts() (opts types.FieldOptions) {
	var guildInfos = GetSignGuilds()
	for _, info := range guildInfos {
		opts = append(opts, types.FieldOption{
			Value: strconv.Itoa(info.GuildId),
			Text:
			//"server:" + strconv.Itoa(k) +
			"id:" + strconv.Itoa(info.GuildId) +
				"Name:" + info.GuildName +
				"point:" + strconv.Itoa(info.GuildPoint)})
	}
	return opts
}

func SubmitBattleInfo(ctx *context.Context) (types.Panel, error) {

	components := template2.Get(config.GetTheme())
	infobox1 := infobox.New()
	col1 := components.Col()
	var size = map[string]string{"md": "3", "sm": "6", "xs": "12"}
	str := gotext.Get("返回上级")
	preUrl := ctx.Request.Header.Get("Referer")
	index := strings.Index(preUrl, "admin")
	preUrl = preUrl[index:]

	//{"group1", "Guild1_1", "Guild1_2", "group2", "Guild2_1", "Guild2_2", "group3", "Guild3_1", "Guild3_2", "group4", "Guild4_1", "Guild4_2"},
	var guildList = make([]int, 0)
	var temp int
	temp, _ = strconv.Atoi(ctx.Request.FormValue("Guild1_1"))
	guildList = append(guildList, temp)
	temp, _ = strconv.Atoi(ctx.Request.FormValue("Guild1_2"))
	guildList = append(guildList, temp)
	temp, _ = strconv.Atoi(ctx.Request.FormValue("Guild2_1"))
	guildList = append(guildList, temp)
	temp, _ = strconv.Atoi(ctx.Request.FormValue("Guild2_2"))
	guildList = append(guildList, temp)
	temp, _ = strconv.Atoi(ctx.Request.FormValue("Guild3_1"))
	guildList = append(guildList, temp)
	temp, _ = strconv.Atoi(ctx.Request.FormValue("Guild3_2"))
	guildList = append(guildList, temp)
	temp, _ = strconv.Atoi(ctx.Request.FormValue("Guild4_1"))
	guildList = append(guildList, temp)
	temp, _ = strconv.Atoi(ctx.Request.FormValue("Guild4_2"))
	guildList = append(guildList, temp)
	if len(guildList) == 0 {
		infoboxCol1 := col1.SetSize(size).SetContent(infobox1.SetText(template.HTML(gotext.Get("请输入工会id"))).GetContent()).GetContent()
		return types.Panel{
			Content:     infoboxCol1,
			Title:       template.HTML(gotext.Get("提交结果")),
			Description: template.HTML(fmt.Sprintf(`<a href="/%s">%s</a>`, preUrl, str)),
			Callbacks:   nil,
		}, nil
	}

	type CrossServerGroup struct {
		MainServer    int    `gorm:"column:id" db:"id"`
		MemberServers string `gorm:"column:serverId" db:"serverId"`
	}

	//var hasSameGuild = false
	//for _, v := range guildList {
	//	for _, v1 := range guildList {
	//		if v == v1 {
	//			hasSameGuild = true
	//		}
	//	}
	//}
	//
	//if hasSameGuild {
	//	infoboxCol1 := col1.SetSize(size).SetContent(infobox1.SetText(template.HTML(gotext.Get("有重复的工会ID"))).GetContent()).GetContent()
	//	return types.Panel{
	//		Content:     infoboxCol1,
	//		Title:       template.HTML(gotext.Get("提交结果")),
	//		Description: template.HTML(fmt.Sprintf(`<a href="/%s">%s</a>`, preUrl, str)),
	//		Callbacks:   nil,
	//	}, nil`
	//}

	var groupInfo = make([]CrossServerGroup, 0)
	db_world := fusion.GetBaseGormDB("db_world")
	db_world.Table("auto_battle_alliance_group").Select("*").Scan(&groupInfo)

	var ServerIds = make([]int, 0)
	db_global := fusion.GetBaseGormDB("db_global")
	db_global.Table("t_guilds").Select("groupId").Where("Id in ?", guildList).Scan(&ServerIds)
	var mainServer int
	var ifFalse bool

	for _, v := range ServerIds {

		for _, v1 := range groupInfo {

			for _, v2 := range strings.Split(v1.MemberServers, ",") {
				tempserver, err := strconv.Atoi(v2)
				if err != nil || tempserver == 0 {
					continue
				}
				if tempserver == v {
					if mainServer == 0 {
						mainServer = v1.MainServer
					}
					if mainServer != v1.MainServer {
						ifFalse = true
						break
					}
				}

			}
		}
	}
	if ifFalse || mainServer == 0 {
		infoboxCol1 := col1.SetSize(size).SetContent(infobox1.SetText(template.HTML(gotext.Get("提交的工会不在同一个跨服分组"))).GetContent()).GetContent()
		return types.Panel{
			Content:     infoboxCol1,
			Title:       template.HTML(gotext.Get("提交结果")),
			Description: template.HTML(fmt.Sprintf(`<a href="/%s">%s</a>`, preUrl, str)),
			Callbacks:   nil,
		}, nil
	}
	var guildIdsStr string
	for _, v := range guildList {
		guildIdsStr += strconv.Itoa(v) + ","
	}

	if len(guildIdsStr) != 0 {
		guildIdsStr = strings.TrimSuffix(guildIdsStr, ",")
	}

	var faild string
	vars := &url.Values{}
	vars.Add("args", guildIdsStr)
	vars.Add("gsId", strconv.Itoa(mainServer))
	vars.Add("cmd", "BattleAllianceSetGuild")
	ctx1 := &context.Context{}
	err, res := fusion.CallToCenter(vars, common.MyCfg.Api["GM2GS"], "GET", ctx1)
	if err != nil {
		faild = "<h5>" + strconv.Itoa(mainServer) + ":" + err.Error() + "</h5>"
	}
	if res != "" {
		faild = "<h5>" + res + "</h5>"
	}

	infoboxCol1 := col1.SetSize(size).SetContent(infobox1.SetText(template.HTML(faild)).GetContent()).GetContent()
	return types.Panel{
		Content:     infoboxCol1,
		Title:       template.HTML(gotext.Get("提交结果")),
		Description: template.HTML(fmt.Sprintf(`<a href="/%s">%s</a>`, preUrl, str)),
		Callbacks:   nil,
	}, nil
}
