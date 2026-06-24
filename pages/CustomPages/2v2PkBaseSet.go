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
	"github.com/GoAdminGroup/go-admin/modules/language"
	form2 "github.com/GoAdminGroup/go-admin/plugins/admin/modules/form"
	template2 "github.com/GoAdminGroup/go-admin/template"
	"github.com/GoAdminGroup/go-admin/template/icon"
	"github.com/GoAdminGroup/go-admin/template/types"
	"github.com/GoAdminGroup/go-admin/template/types/form"
	"github.com/GoAdminGroup/themes/adminlte/components/infobox"
	"github.com/leonelquinteros/gotext"
)

type signInfo struct {
	charName      string
	signTime      uint64
	playerId      uint32
	serverId      uint32
	state         uint32
	submitItemNum uint32
	teamId        uint32
}

type teamInfo struct {
	id        uint32
	teamName  string
	winTimes  uint32
	loseTimes uint32
}

type matchState struct {
	index          uint32
	winTeam        uint32
	teamA          uint32
	teamB          uint32
	matchStartTime uint64
	playerState    uint32
}

func Get2V2PkConfigPage(ctx *context.Context) (types.Panel, error) {
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
	MainServer, _ := strconv.Atoi(common.MyCfg.Api["MainServer"])
	gormDB := fusion.GetServerGormDB("db_char", int64(MainServer))
	var cfgValue string
	gormDB.Table("inst_configure").Select("cfgValue").Where("cfgID = ?", 37).Scan(&cfgValue)
	unpacker := fusion.NewTextUnpacker(cfgValue)
	var state = unpacker.UnpackUint()
	var signStartTime = unpacker.UnpackUint()
	var signEndTime = unpacker.UnpackUint()
	var preMatchStartTime = unpacker.UnpackUint()
	var preMatchEndTime = unpacker.UnpackUint()
	var size = 0
	signDatas := make([]signInfo, 0)
	size = int(unpacker.UnpackUint())
	for i := 0; i < size; i++ {
		signData := signInfo{}
		signData.playerId = uint32(unpacker.UnpackUint())
		signData.charName = unpacker.UnpackString()
		signData.signTime = unpacker.UnpackUint()
		signData.state = uint32(unpacker.UnpackUint())
		signData.serverId = uint32(unpacker.UnpackUint())
		signData.submitItemNum = uint32(unpacker.UnpackUint())
		signData.teamId = uint32(unpacker.UnpackUint())
		signDatas = append(signDatas, signData)
	}
	size = int(unpacker.UnpackUint())
	teamDatas := make([]teamInfo, 0)
	teamId2Name := make(map[int]string, 0)
	for i := 0; i < size; i++ {
		teamData := teamInfo{}
		teamData.id = uint32(unpacker.UnpackUint())
		teamData.teamName = unpacker.UnpackString()
		teamData.winTimes = uint32(unpacker.UnpackUint())
		teamData.loseTimes = uint32(unpacker.UnpackUint())
		memberSize := uint32(unpacker.UnpackUint())
		for j := 0; j < int(memberSize); j++ {
			_ = uint32(unpacker.UnpackUint())
		}
		teamDatas = append(teamDatas, teamData)
		teamId2Name[int(teamData.id)] = teamData.teamName
	}
	size = int(unpacker.UnpackUint())
	matchDatas := make(map[int]matchState, 0)
	for i := 0; i < size; i++ {
		matchData := matchState{}
		matchData.index = uint32(unpacker.UnpackUint())
		matchData.matchStartTime = unpacker.UnpackUint()
		matchData.winTeam = uint32(unpacker.UnpackUint())
		matchData.teamA = uint32(unpacker.UnpackUint())
		matchData.teamB = uint32(unpacker.UnpackUint())
		matchData.playerState = uint32(unpacker.UnpackUint())
		matchDatas[int(matchData.index)] = matchData
	}

	var teamOpts = make(types.FieldOptions, 0)
	for _, v := range teamDatas {
		teamOpts = append(teamOpts, types.FieldOption{Text: v.teamName, Value: strconv.Itoa(int(v.id))})
	}

	var formList = types.NewFormPanel()
	formList.AddField(gotext.Get("当前阶段"), "state", db.Int, form.Text).
		FieldDisplay(func(model types.FieldModel) interface{} {
			value, _ := strconv.Atoi(model.Value)
			switch value {
			case 0:
				return gotext.Get("初始状态")
			case 1:
				return gotext.Get("报名状态")
			case 2:
				return gotext.Get("报名结束状态")
			case 3:
				return gotext.Get("预组队状态")
			case 4:
				return gotext.Get("预组队结束状态")
			case 5:
				return gotext.Get("比赛状态")
			case 6:
				return gotext.Get("结束状态")
			}
			return model.Value
		}).
		FieldDefault(strconv.Itoa(int(state))).
		FieldDisplayButCanNotEditWhenUpdate().
		FieldDisplayButCanNotEditWhenCreate()

	formList.AddField(gotext.Get("报名开始时间"), "signStartTime", db.Int, form.Datetime).FieldDefault(strconv.Itoa(int(signStartTime))).
		FieldDisplay(func(model types.FieldModel) interface{} {
			cTime, _ := strconv.Atoi(model.Value)
			if cTime == 0 {
				return ""
			}
			displayTime := time.Unix(int64(cTime), 0)
			return displayTime.Format("2006-01-02 15:04:05")
		})

	formList.AddField(gotext.Get("报名结束时间"), "signEndTime", db.Int, form.Datetime).FieldDefault(strconv.Itoa(int(signEndTime))).
		FieldDisplay(func(model types.FieldModel) interface{} {
			cTime, _ := strconv.Atoi(model.Value)
			if cTime == 0 {
				return ""
			}
			displayTime := time.Unix(int64(cTime), 0)
			return displayTime.Format("2006-01-02 15:04:05")
		})

	formList.AddField(gotext.Get("预组队开始时间"), "preMatchStartTime", db.Int, form.Datetime).FieldDefault(strconv.Itoa(int(preMatchStartTime))).
		FieldDisplay(func(model types.FieldModel) interface{} {
			cTime, _ := strconv.Atoi(model.Value)
			if cTime == 0 {
				return nil
			}
			displayTime := time.Unix(int64(cTime), 0)
			return displayTime.Format("2006-01-02 15:04:05")
		})
	formList.AddField(gotext.Get("预组队结束时间"), "preMatchEndTime", db.Int, form.Datetime).FieldDefault(strconv.Itoa(int(preMatchEndTime))).
		FieldDisplay(func(model types.FieldModel) interface{} {
			cTime, _ := strconv.Atoi(model.Value)
			if cTime == 0 {
				return nil
			}
			displayTime := time.Unix(int64(cTime), 0)
			return displayTime.Format("2006-01-02 15:04:05")
		})

	startTime1 := 0
	teamA1 := 0
	teamB1 := 0
	singleMatchData, isok := matchDatas[1]
	if isok {
		startTime1 = int(singleMatchData.matchStartTime)
		teamA1 = int(singleMatchData.teamA)
		teamB1 = int(singleMatchData.teamB)
	}
	formList.AddField(gotext.Get("第一组："), "group1", db.Int, form.Text).
		FieldDisplayButCanNotEditWhenUpdate().FieldDisplayButCanNotEditWhenCreate().FieldDefault("1")
	formList.AddField(gotext.Get("队伍A"), "TeamA_1", db.Int, form.SelectSingle).FieldOptions(teamOpts).FieldDefault(strconv.Itoa(teamA1)).
		FieldDisplay(func(model types.FieldModel) interface{} {
			teamId, _ := strconv.Atoi(model.Value)
			teamName, isOk := teamId2Name[teamId]
			if !isOk {
				return model.Value
			}
			return teamName
		})
	formList.AddField(gotext.Get("队伍B"), "TeamB_1", db.Int, form.SelectSingle).FieldOptions(teamOpts).FieldDefault(strconv.Itoa(teamB1)).
		FieldDisplay(func(model types.FieldModel) interface{} {
			teamId, _ := strconv.Atoi(model.Value)
			teamName, isOk := teamId2Name[teamId]
			if !isOk {
				return model.Value
			}
			return teamName
		})
	formList.AddField(gotext.Get("对局开始时间"), "matchStartTime_1", db.Int, form.Datetime).
		FieldDefault(strconv.Itoa(startTime1)).
		FieldDisplay(func(model types.FieldModel) interface{} {
			cTime, _ := strconv.Atoi(model.Value)
			if cTime == 0 {
				return ""
			}
			displayTime := time.Unix(int64(cTime), 0)
			return displayTime.Format("2006-01-02 15:04:05")
		})

	startTime2 := 0
	teamA2 := 0
	teamB2 := 0
	singleMatchData, isok = matchDatas[2]
	if isok {
		startTime2 = int(singleMatchData.matchStartTime)
		teamA2 = int(singleMatchData.teamA)
		teamB2 = int(singleMatchData.teamB)
	}
	formList.AddField(gotext.Get("第二组："), "group2", db.Int, form.Text).
		FieldDisplayButCanNotEditWhenUpdate().FieldDisplayButCanNotEditWhenCreate().FieldDefault("2")
	formList.AddField(gotext.Get("队伍A"), "TeamA_2", db.Int, form.SelectSingle).FieldOptions(teamOpts).FieldDefault(strconv.Itoa(teamA2)).
		FieldDisplay(func(model types.FieldModel) interface{} {
			teamId, _ := strconv.Atoi(model.Value)
			teamName, isOk := teamId2Name[teamId]
			if !isOk {
				return model.Value
			}
			return teamName
		})
	formList.AddField(gotext.Get("队伍B"), "TeamB_2", db.Int, form.SelectSingle).FieldOptions(teamOpts).FieldDefault(strconv.Itoa(teamB2)).
		FieldDisplay(func(model types.FieldModel) interface{} {
			teamId, _ := strconv.Atoi(model.Value)
			teamName, isOk := teamId2Name[teamId]
			if !isOk {
				return model.Value
			}
			return teamName
		})

	formList.AddField(gotext.Get("对局开始时间"), "matchStartTime_2", db.Int, form.Datetime).
		FieldDefault(strconv.Itoa(startTime2)).
		FieldDisplay(func(model types.FieldModel) interface{} {
			cTime, _ := strconv.Atoi(model.Value)
			if cTime == 0 {
				return ""
			}
			displayTime := time.Unix(int64(cTime), 0)
			return displayTime.Format("2006-01-02 15:04:05")
		})

	startTime3 := 0
	teamA3 := 0
	teamB3 := 0
	singleMatchData, isok = matchDatas[3]
	if isok {
		startTime3 = int(singleMatchData.matchStartTime)
		teamA3 = int(singleMatchData.teamA)
		teamB3 = int(singleMatchData.teamB)
	}
	formList.AddField(gotext.Get("第三组："), "group3", db.Int, form.Text).
		FieldDisplayButCanNotEditWhenUpdate().FieldDisplayButCanNotEditWhenCreate().FieldDefault("3")
	formList.AddField(gotext.Get("队伍A"), "TeamA_3", db.Int, form.SelectSingle).FieldOptions(teamOpts).FieldDefault(strconv.Itoa(teamA3)).
		FieldDisplay(func(model types.FieldModel) interface{} {
			teamId, _ := strconv.Atoi(model.Value)
			teamName, isOk := teamId2Name[teamId]
			if !isOk {
				return model.Value
			}
			return teamName
		})
	formList.AddField(gotext.Get("队伍B"), "TeamB_3", db.Int, form.SelectSingle).FieldOptions(teamOpts).FieldDefault(strconv.Itoa(teamB3)).
		FieldDisplay(func(model types.FieldModel) interface{} {
			teamId, _ := strconv.Atoi(model.Value)
			teamName, isOk := teamId2Name[teamId]
			if !isOk {
				return model.Value
			}
			return teamName
		})
	formList.AddField(gotext.Get("对局开始时间"), "matchStartTime_3", db.Int, form.Datetime).
		FieldDefault(strconv.Itoa(startTime3)).
		FieldDisplay(func(model types.FieldModel) interface{} {
			cTime, _ := strconv.Atoi(model.Value)
			if cTime == 0 {
				return ""
			}
			displayTime := time.Unix(int64(cTime), 0)
			return displayTime.Format("2006-01-02 15:04:05")
		})

	startTime4 := 0
	teamA4 := 0
	teamB4 := 0
	singleMatchData, isok = matchDatas[4]
	if isok {
		startTime4 = int(singleMatchData.matchStartTime)
		teamA4 = int(singleMatchData.teamA)
		teamB4 = int(singleMatchData.teamB)
	}
	formList.AddField(gotext.Get("第四组："), "group4", db.Int, form.Text).
		FieldDisplayButCanNotEditWhenUpdate().FieldDisplayButCanNotEditWhenCreate().FieldDefault("4")
	formList.AddField(gotext.Get("队伍A"), "TeamA_4", db.Int, form.SelectSingle).FieldOptions(teamOpts).FieldDefault(strconv.Itoa(teamA4)).
		FieldDisplay(func(model types.FieldModel) interface{} {
			teamId, _ := strconv.Atoi(model.Value)
			teamName, isOk := teamId2Name[teamId]
			if !isOk {
				return model.Value
			}
			return teamName
		})
	formList.AddField(gotext.Get("队伍B"), "TeamB_4", db.Int, form.SelectSingle).FieldOptions(teamOpts).FieldDefault(strconv.Itoa(teamB4)).
		FieldDisplay(func(model types.FieldModel) interface{} {
			teamId, _ := strconv.Atoi(model.Value)
			teamName, isOk := teamId2Name[teamId]
			if !isOk {
				return model.Value
			}
			return teamName
		})
	formList.AddField(gotext.Get("对局开始时间"), "matchStartTime_4", db.Int, form.Datetime).
		FieldDefault(strconv.Itoa(startTime4)).
		FieldDisplay(func(model types.FieldModel) interface{} {
			cTime, _ := strconv.Atoi(model.Value)
			if cTime == 0 {
				return ""
			}
			displayTime := time.Unix(int64(cTime), 0)
			return displayTime.Format("2006-01-02 15:04:05")
		})

	startTime5 := 0
	teamA5 := 0
	teamB5 := 0
	singleMatchData, isok = matchDatas[5]
	if isok {
		startTime5 = int(singleMatchData.matchStartTime)
		teamA5 = int(singleMatchData.teamA)
		teamB5 = int(singleMatchData.teamB)
	}
	formList.AddField(gotext.Get("第五组："), "group5", db.Int, form.Text).
		FieldDisplayButCanNotEditWhenUpdate().FieldDisplayButCanNotEditWhenCreate().FieldDefault("5")
	formList.AddField(gotext.Get("队伍A"), "TeamA_5", db.Int, form.SelectSingle).FieldOptions(teamOpts).FieldDefault(strconv.Itoa(teamA5)).
		FieldDisplay(func(model types.FieldModel) interface{} {
			teamId, _ := strconv.Atoi(model.Value)
			teamName, isOk := teamId2Name[teamId]
			if !isOk {
				return model.Value
			}
			return teamName
		})
	formList.AddField(gotext.Get("队伍B"), "TeamB_5", db.Int, form.SelectSingle).FieldOptions(teamOpts).FieldDefault(strconv.Itoa(teamB5)).
		FieldDisplay(func(model types.FieldModel) interface{} {
			teamId, _ := strconv.Atoi(model.Value)
			teamName, isOk := teamId2Name[teamId]
			if !isOk {
				return model.Value
			}
			return teamName
		})
	formList.AddField(gotext.Get("对局开始时间"), "matchStartTime_5", db.Int, form.Datetime).
		FieldDefault(strconv.Itoa(startTime5)).
		FieldDisplay(func(model types.FieldModel) interface{} {
			cTime, _ := strconv.Atoi(model.Value)
			if cTime == 0 {
				return ""
			}
			displayTime := time.Unix(int64(cTime), 0)
			return displayTime.Format("2006-01-02 15:04:05")
		})

	startTime6 := 0
	teamA6 := 0
	teamB6 := 0
	singleMatchData, isok = matchDatas[6]
	if isok {
		startTime6 = int(singleMatchData.matchStartTime)
		teamA6 = int(singleMatchData.teamA)
		teamB6 = int(singleMatchData.teamB)
	}
	formList.AddField(gotext.Get("第六组："), "group6", db.Int, form.Text).
		FieldDisplayButCanNotEditWhenUpdate().FieldDisplayButCanNotEditWhenCreate().FieldDefault("6")
	formList.AddField(gotext.Get("队伍A"), "TeamA_6", db.Int, form.SelectSingle).FieldOptions(teamOpts).FieldDefault(strconv.Itoa(teamA6)).
		FieldDisplay(func(model types.FieldModel) interface{} {
			teamId, _ := strconv.Atoi(model.Value)
			teamName, isOk := teamId2Name[teamId]
			if !isOk {
				return model.Value
			}
			return teamName
		})
	formList.AddField(gotext.Get("队伍B"), "TeamB_6", db.Int, form.SelectSingle).FieldOptions(teamOpts).FieldDefault(strconv.Itoa(teamB6)).
		FieldDisplay(func(model types.FieldModel) interface{} {
			teamId, _ := strconv.Atoi(model.Value)
			teamName, isOk := teamId2Name[teamId]
			if !isOk {
				return model.Value
			}
			return teamName
		})
	formList.AddField(gotext.Get("对局开始时间"), "matchStartTime_6", db.Int, form.Datetime).
		FieldDefault(strconv.Itoa(startTime6)).
		FieldDisplay(func(model types.FieldModel) interface{} {
			cTime, _ := strconv.Atoi(model.Value)
			if cTime == 0 {
				return ""
			}
			displayTime := time.Unix(int64(cTime), 0)
			return displayTime.Format("2006-01-02 15:04:05")
		})

	formList.SetTabGroups(types.TabGroups{
		{"state", "signStartTime", "signEndTime", "preMatchStartTime", "preMatchEndTime",
			"group1", "TeamA_1", "TeamB_1", "matchStartTime_1", "group2", "TeamA_2", "TeamB_2", "matchStartTime_2", "group3", "TeamA_3", "TeamB_3", "matchStartTime_3",
			"group4", "TeamA_4", "TeamB_4", "matchStartTime_4", "group5", "TeamA_5", "TeamB_5", "matchStartTime_5", "group6", "TeamA_6", "TeamB_6", "matchStartTime_6"}})
	formList.SetTabHeaders(gotext.Get("2v2Pk分组"))
	field, headers := formList.GroupField()

	formList.SetTabHeaders(gotext.Get("2v2Pk分组"))
	aform := components.Form().
		SetTabHeaders(headers).
		SetTabContents(field).
		SetPrefix(config.PrefixFixSlash()).
		SetUrl("/admin/Submit2V2PkInfo").
		SetHiddenFields(map[string]string{
			form2.PreviousKey: "/admin",
		}).
		SetMethod("GET").
		SetOperationFooter(col1 + col2)

	return types.Panel{
		Content: components.Box().
			SetHeader(aform.GetDefaultBoxHeader(true)).
			WithHeadBorder().
			SetHeader(template.HTML(gotext.Get("2v2Pk分组"))).
			SetBody(aform.GetContent() + formList.FooterHtml).
			GetContent(),
		Title:     template.HTML(gotext.Get("2v2Pk分组")),
		CSS:       `.modal.fade.in{z-index:10002}`,
		Callbacks: nil,
	}, nil
}

func Submit2V2PkInfo(ctx *context.Context) (types.Panel, error) {
	components := template2.Get(config.GetTheme())
	infobox1 := infobox.New()
	col1 := components.Col()
	var size = map[string]string{"md": "3", "sm": "6", "xs": "12"}
	str := gotext.Get("返回上级")
	preUrl := ctx.Request.Header.Get("Referer")
	index := strings.Index(preUrl, "admin")
	if index > -1 {
		preUrl = preUrl[index:]
	} else {
		preUrl = ""
	}

	signStartTimeStr := ctx.Request.FormValue("signStartTime")
	signStartTimeUnix := "0"
	if signStartTimeStr != "" && signStartTimeStr != "0" {
		signStartTime, _ := time.ParseInLocation("2006-01-02 15:04:05", signStartTimeStr, time.Local)
		signStartTimeUnix = strconv.Itoa(int(signStartTime.Unix()))
	}
	signEndTimeStr := ctx.Request.FormValue("signEndTime")
	signEndTimeUnix := "0"
	if signEndTimeStr != "" && signEndTimeStr != "0" {
		signEndTime, _ := time.ParseInLocation("2006-01-02 15:04:05", signEndTimeStr, time.Local)
		signEndTimeUnix = strconv.Itoa(int(signEndTime.Unix()))
	}
	preMatchStartTimeStr := ctx.Request.FormValue("preMatchStartTime")
	preMatchStartTimeUnix := "0"
	if preMatchStartTimeStr != "" && preMatchStartTimeStr != "0" {
		preMatchStartTime, _ := time.ParseInLocation("2006-01-02 15:04:05", preMatchStartTimeStr, time.Local)
		preMatchStartTimeUnix = strconv.Itoa(int(preMatchStartTime.Unix()))
	}

	preMatchEndTimeStr := ctx.Request.FormValue("preMatchEndTime")
	preMatchEndTimeUnix := "0"
	if preMatchEndTimeStr != "" && preMatchEndTimeStr != "0" {
		preMatchEndTime, _ := time.ParseInLocation("2006-01-02 15:04:05", preMatchEndTimeStr, time.Local)
		preMatchEndTimeUnix = strconv.Itoa(int(preMatchEndTime.Unix()))
	}

	type matchInfo struct {
		index          string
		teamA          string
		teamB          string
		matchStartTime string
	}

	matchDatas := make([]matchInfo, 0)
	for i := 1; i <= 6; i++ {
		group, _ := strconv.Atoi(ctx.Request.FormValue("group" + strconv.Itoa(i)))
		teamA, _ := strconv.Atoi(ctx.Request.FormValue("TeamA_" + strconv.Itoa(i)))
		teamB, _ := strconv.Atoi(ctx.Request.FormValue("TeamB_" + strconv.Itoa(i)))

		matchStartTimeStr := ctx.Request.FormValue("matchStartTime_" + strconv.Itoa(i))
		matchStartTimeUnix := "0"
		if matchStartTimeStr != "" && matchStartTimeStr != "0" {
			matchStartTime, _ := time.ParseInLocation("2006-01-02 15:04:05", matchStartTimeStr, time.Local)
			matchStartTimeUnix = strconv.Itoa(int(matchStartTime.Unix()))
		}
		matchDatas = append(matchDatas, matchInfo{strconv.Itoa(group), strconv.Itoa(teamA),
			strconv.Itoa(teamB), matchStartTimeUnix})
	}

	var packer = signEndTimeUnix + "," + signStartTimeUnix + "," + preMatchStartTimeUnix + "," + preMatchEndTimeUnix + ","
	for _, v := range matchDatas {
		packer += v.index + "," + v.teamA + "," + v.teamB + "," + v.matchStartTime + ","
	}
	packer = strings.TrimRight(packer, ",")

	vars := &url.Values{}
	var faild string
	vars.Add("cmd", "ModifyAnyOnAnyMatchInfo")
	vars.Add("gsId", "0")
	vars.Add("args", packer)

	ctx1 := &context.Context{}
	err, res := fusion.CallToCenter(vars, common.MyCfg.Api["GM2GS"], "GET", ctx1)
	if err != nil {
		faild = "<h5>" + common.MyCfg.Api["MainServer"] + ":" + err.Error() + "</h5>"
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
