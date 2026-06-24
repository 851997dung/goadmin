package CustomPages

import (
	"admin/common/def"
	"admin/fusion"
	"encoding/json"
	"github.com/GoAdminGroup/go-admin/context"
	"github.com/GoAdminGroup/go-admin/modules/config"
	"github.com/GoAdminGroup/go-admin/modules/db"
	"github.com/GoAdminGroup/go-admin/plugins/admin/modules/constant"
	form2 "github.com/GoAdminGroup/go-admin/plugins/admin/modules/form"
	"github.com/GoAdminGroup/go-admin/plugins/admin/modules/parameter"
	"github.com/GoAdminGroup/go-admin/plugins/admin/modules/table"
	template2 "github.com/GoAdminGroup/go-admin/template"
	"github.com/GoAdminGroup/go-admin/template/chartjs"
	"github.com/GoAdminGroup/go-admin/template/icon"
	"github.com/GoAdminGroup/go-admin/template/types"
	"github.com/GoAdminGroup/go-admin/template/types/form"
	"github.com/leonelquinteros/gotext"
	"html/template"
	"sort"
	"strconv"
	"time"
)

func GetOnlinePopulationTrendChart(ctx *context.Context) (types.Panel, error) {
	ServerID, _ := strconv.Atoi(ctx.Request.FormValue("ServerId"))

	components := template2.Get(config.GetTheme())
	btn1 := components.Button().SetType("submit").
		AddClass("submit").
		SetContent(icon.Icon(icon.Search, 2) + template.HTML(gotext.Get("查找"))).
		SetThemePrimary().
		SetSmallSize().
		SetOrientationLeft().
		SetLoadingText(icon.Icon(icon.Spinner, 1) + template.HTML(gotext.Get("查找"))).
		SetID("submit1").
		GetContent()

	var form1 = types.NewFormPanel()
	myOP := fusion.GetServerList()

	form1.AddField(gotext.Get("服务器ID"), "ServerId", db.Int, form.SelectSingle).
		FieldOptions(myOP).
		FieldMust()

	params := parameter.GetParam(ctx.Request.URL, 10)
	params.SortField = ""

	data, _ := GetOnlineNumByMin(ctx).GetData(params)

	dsData := make([]float64, 0)

	labels := []string{}
	for _, v := range data.InfoList {
		if ctx.Request.FormValue("ServerId") != v["serverId"].Value {
			continue
		}
		labels = append(labels, v["label"].Value)
		temp := v["num"].Value
		tempF, _ := strconv.ParseFloat(temp, 64)
		dsData = append(dsData, tempF)
	}

	chart := chartjs.Line()
	chart.SetID("online4Min").SetHeight(300).
		SetLabels(labels).
		AddDataSet(gotext.Get("人数")).
		DSData(dsData).
		DSFill(false).
		DSBorderColor("rgb(210, 214, 222)").
		DSLineTension(0.1)

	col := components.Col().SetContent(chart.GetContent()).SetSize(types.SizeMD(8)).GetContent()
	boxInternalRow := components.Row().SetContent(col).GetContent()
	box := components.Box().SetTheme("danger").WithHeadBorder().SetHeader(template.HTML(gotext.Get("人数"))).
		SetBody(boxInternalRow).SetIframeStyle(true).
		GetContent()

	myStr := fusion.GetServerNameStrings()
	col3 := components.Col().SetContent(template.HTML(gotext.Get("服务器ID") + ":" + myStr[ServerID])).GetContent()
	//col10 := components.Col().SetContent(template.HTML(gotext.Get("记录时间") + ":" + DayNoStart + "-" + DayNoEnd)).GetContent()
	row4 := components.Row().SetContent(col3 /*+ col10*/).GetContent()

	isNotIframe := ctx.Query(constant.IframeKey) != "true"
	boxModel := components.Box().
		SetNoPadding().
		WithHeadBorder().
		SetIframeStyle(!isNotIframe).
		SetHeader(template.HTML(gotext.Get("数据明细")))
	if ServerID != 0 {
		boxModel.SetBody(row4 + box)
	}

	boxModel = boxModel.SetSecondHeaderClass("filter-area").
		SetSecondHeader(
			components.Form().
				SetHeader(form1.Header).
				SetContent(form1.FieldList).
				SetPrefix(config.PrefixFixSlash()).
				SetUrl("/admin/OnlinePopulationTrendChart").
				SetId("form1").
				SetMethod("GET").
				SetOperationFooter(btn1).
				SetHiddenFields(map[string]string{
					form2.NoAnimationKey: "true",
				}).
				GetContent())

	return types.Panel{Content: boxModel.GetContent() + form1.FooterHtml}, nil
}

func GetOnlineNumByMin(ctx *context.Context) table.Table {
	OnlineNumByHour := table.NewDefaultTable(table.Config{
		Driver:     "mysql",
		Connection: "db_global",
		CanAdd:     true,
		Editable:   true,
		Deletable:  true,
		Exportable: true,
		PrimaryKey: table.PrimaryKey{
			Type: db.Varchar,
			Name: "label",
		},
	})
	info := OnlineNumByHour.GetInfo().HideEditButton().HideNewButton().HideRowSelector().HideDetailButton().HideDeleteButton()

	myOp := fusion.GetServerList()
	info.AddField(gotext.Get("服务器ID"), "serverId", db.Int).
		FieldFilterable(types.FilterType{FormType: form.SelectSingle}).
		FieldFilterOptions(myOp)
	info.AddField("label", "label", db.Varchar).
		FieldFilterable(types.FilterType{FormType: form.Text})
	info.AddField("num", "num", db.Int).
		FieldFilterable(types.FilterType{FormType: form.Text})

	info.SetTitle(gotext.Get("分时在线"))
	info.SetGetDataFn(
		func(param parameter.Parameters) (data []map[string]interface{}, size int) {
			serverId, _ := strconv.Atoi(ctx.Request.FormValue("ServerId"))
			if serverId == 0 {
				return nil, 0
			}

			go_manager := fusion.GetBaseGormDB("go_manager")
			onlineNumRec := []def.OnlineNumRecMin{}
			go_manager.Table("player_online_num_record_min").
				Order("logTime DESC").
				Scan(&onlineNumRec)
			sort.Slice(onlineNumRec, func(i, j int) bool {
				return time.Time(onlineNumRec[i].LogTime).After(time.Time(onlineNumRec[j].LogTime))
			})

			for i := len(onlineNumRec) - 1; i >= 0; i-- {
				jsonData := map[int]int{}
				temp := make(map[string]interface{})
				temp["serverId"] = serverId
				logTime := time.Time(onlineNumRec[i].LogTime)
				temp["label"] = logTime.Format("15:04")
				json.Unmarshal(fusion.String2Bytes(onlineNumRec[i].Record), &jsonData)
				temp["num"] = jsonData[serverId]
				data = append(data, temp)
			}
			return data, len(data)
		})

	info.SetExportProcessFn(func(param parameter.Parameters) (types.PanelInfo, error) {
		panelInfo, _ := OnlineNumByHour.GetData(param.WithIsAll(param.IsAll()))
		return types.PanelInfo{Thead: panelInfo.Thead, InfoList: panelInfo.InfoList}, nil
	})

	return OnlineNumByHour
}
