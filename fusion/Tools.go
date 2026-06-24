package fusion

import (
	"admin/common"
	"admin/common/def"
	"admin/common/def/chequeDef"
	"admin/common/def/currencyDef"
	"admin/common/def/excelConfigDef/expconfigdef"
	"admin/common/def/flowDef"
	"admin/fusion/myConfig"
	"crypto/md5"
	"crypto/tls"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"html/template"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"reflect"
	"strconv"
	"strings"
	"sync"
	"time"
	"unsafe"

	"github.com/360EntSecGroup-Skylar/excelize"
	"github.com/GoAdminGroup/go-admin/context"
	"github.com/GoAdminGroup/go-admin/modules/auth"
	"github.com/GoAdminGroup/go-admin/modules/config"
	"github.com/GoAdminGroup/go-admin/modules/db"
	"github.com/GoAdminGroup/go-admin/modules/language"
	"github.com/GoAdminGroup/go-admin/modules/logger"
	"github.com/GoAdminGroup/go-admin/plugins/admin/modules/constant"
	"github.com/GoAdminGroup/go-admin/plugins/admin/modules/form"
	"github.com/GoAdminGroup/go-admin/plugins/admin/modules/parameter"
	"github.com/GoAdminGroup/go-admin/plugins/admin/modules/table"
	template2 "github.com/GoAdminGroup/go-admin/template"
	"github.com/GoAdminGroup/go-admin/template/icon"
	"github.com/GoAdminGroup/go-admin/template/types"
	"github.com/GoAdminGroup/html"
	"github.com/leonelquinteros/gotext"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"
)

var combinedServerCache sync.Map // int -> int

var centerHTTPClient = &http.Client{
	Transport: &http.Transport{
		TLSClientConfig:   &tls.Config{InsecureSkipVerify: true},
		MaxIdleConns:      20,
		IdleConnTimeout:   90 * time.Second,
		DisableKeepAlives: false,
	},
	Timeout: 30 * time.Second,
}

var (
	itemNameCache     map[int]string
	itemNameCacheOnce sync.Once
)

type tableCacheEntry struct {
	tables    []string
	expiresAt time.Time
}

var (
	tableCacheMu sync.RWMutex
	tableCache   = map[string]*tableCacheEntry{}
)

func getCachedTableNames(gormDB *gorm.DB, tablePrefix string) []string {
	var dbName string
	gormDB.Raw("SELECT DATABASE();").Scan(&dbName)
	if dbName == "" {
		return nil
	}
	cacheKey := dbName + "_" + tablePrefix
	tableCacheMu.RLock()
	entry, ok := tableCache[cacheKey]
	tableCacheMu.RUnlock()
	if ok && time.Now().Before(entry.expiresAt) {
		return entry.tables
	}
	tableNames := []string{}
	sql := "SELECT table_name FROM information_schema.TABLES WHERE" +
		" TABLE_SCHEMA='" + dbName + "' AND table_name LIKE '%" + tablePrefix + "%' "
	gormDB.Raw(sql).Scan(&tableNames)
	tableCacheMu.Lock()
	tableCache[cacheKey] = &tableCacheEntry{tables: tableNames, expiresAt: time.Now().Add(5 * time.Minute)}
	tableCacheMu.Unlock()
	return tableNames
}

func CallToDeploy1(vals url.Values, path string, method string, ctx context.Context, host, port string) (error, string) {
	sUrl := &url.URL{
		Scheme:   "http",
		Host:     host + ":" + port,
		Path:     path,
		Fragment: "anchor",
	}
	sUrl.RawQuery = vals.Encode()
	ctx.Request, _ = http.NewRequest(method, sUrl.String(), strings.NewReader(vals.Encode()))
	tr := &http.Transport{
		TLSClientConfig:   &tls.Config{InsecureSkipVerify: true},
		DisableKeepAlives: true,
	}
	client := &http.Client{Transport: tr}
	var err error
	ctx.Response, err = client.Do(ctx.Request)
	if err != nil {
		fmt.Printf("Fail exec http client Do,err:%s\n", err.Error())
		return err, ""
	}
	defer func() {
		ctx.Response.Body.Close()
	}()
	returnBody, _ := ParseResponse(ctx.Response)
	res := (*string)(unsafe.Pointer(&returnBody))
	return nil, *res
}

func CallToDeploy(vals *url.Values, path string, method string, ctx *context.Context, host, port string) (error, string) {
	sUrl := &url.URL{
		Scheme:   "http",
		Host:     host + ":" + port,
		Path:     path,
		Fragment: "anchor",
	}
	sUrl.RawQuery = vals.Encode()
	ctx.Request, _ = http.NewRequest(method, sUrl.String(), strings.NewReader(vals.Encode()))
	tr := &http.Transport{
		TLSClientConfig:   &tls.Config{InsecureSkipVerify: true},
		DisableKeepAlives: true,
	}
	client := &http.Client{Transport: tr}
	var err error
	ctx.Response, err = client.Do(ctx.Request)
	if err != nil {
		fmt.Printf("Fail exec http client Do,err:%s\n", err.Error())
		return err, ""
	}
	defer func() {
		ctx.Response.Body.Close()
	}()
	returnBody, err := ParseResponse(ctx.Response)
	res := (*string)(unsafe.Pointer(&returnBody))
	return nil, *res
}

func CallToBackup(vals *url.Values, path string, method string, ctx *context.Context, host, port string) (error, string) {
	sUrl := &url.URL{
		Scheme:   "http",
		Host:     host + ":" + port,
		Path:     path,
		Fragment: "anchor",
	}
	sUrl.RawQuery = vals.Encode()
	ctx.Request, _ = http.NewRequest(method, sUrl.String(), strings.NewReader(vals.Encode()))
	tr := &http.Transport{
		TLSClientConfig:   &tls.Config{InsecureSkipVerify: true},
		DisableKeepAlives: true,
	}
	client := &http.Client{Transport: tr}
	var err error
	ctx.Response, err = client.Do(ctx.Request)
	if err != nil {
		fmt.Printf("Fail exec http client Do,err:%s\n", err.Error())
		return err, ""
	}
	defer func() {
		ctx.Response.Body.Close()
	}()
	returnBody, err := ParseResponse(ctx.Response)
	res := (*string)(unsafe.Pointer(&returnBody))
	return nil, *res
}

func CallToCenter(vals *url.Values, path string, method string, ctx *context.Context) (error, string) {
	sUrl := &url.URL{
		Scheme:   "https",
		Host:     common.MyCfg.Center.Host + ":" + common.MyCfg.Center.Port,
		Path:     path,
		Fragment: "anchor",
	}
	sUrl.RawQuery = vals.Encode()
	ctx.Request, _ = http.NewRequest(method, sUrl.String(), strings.NewReader(vals.Encode()))
	var err error
	ctx.Response, err = centerHTTPClient.Do(ctx.Request)
	if err != nil {
		fmt.Printf("Fail exec http client Do,err:%s\n", err.Error())
		return err, ""
	}
	defer func() {
		ctx.Response.Body.Close()
	}()
	returnBody, err := ParseResponse(ctx.Response)
	res := (*string)(unsafe.Pointer(&returnBody))
	return nil, *res
}

func ParseResponse(response *http.Response) ([]byte, error) {
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}
	return body, err
}

func String2Bytes(s string) []byte {
	sh := (*reflect.StringHeader)(unsafe.Pointer(&s))
	bh := reflect.SliceHeader{
		Data: sh.Data,
		Len:  sh.Len,
		Cap:  sh.Len,
	}
	return *(*[]byte)(unsafe.Pointer(&bh))
}

func Bytes2String(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}
func CreateInfoList(list *[]map[string]interface{}) (infolist types.InfoList) {
	for _, v := range *list {
		temp := map[string]types.InfoItem{}
		t := reflect.TypeOf(v)
		if t.Kind() == reflect.Ptr {
			t = t.Elem()
		}
		for key, value := range v {
			tempItem := types.InfoItem{}
			tempItem.Value = Strval(value)
			tempItem.Content = template.HTML(Strval(value))
			temp[key] = tempItem
		}
		infolist = append(infolist, temp)
	}
	return infolist
}
func Strval(value interface{}) string {
	// interface 转 string
	var key string
	if value == nil {
		return key
	}
	switch value.(type) {
	case float64:
		ft := value.(float64)
		key = strconv.FormatFloat(ft, 'f', -1, 64)
	case float32:
		ft := value.(float32)
		key = strconv.FormatFloat(float64(ft), 'f', -1, 64)
	case int:
		it := value.(int)
		key = strconv.Itoa(it)
	case uint:
		it := value.(uint)
		key = strconv.Itoa(int(it))
	case int8:
		it := value.(int8)
		key = strconv.Itoa(int(it))
	case uint8:
		it := value.(uint8)
		key = strconv.Itoa(int(it))
	case int16:
		it := value.(int16)
		key = strconv.Itoa(int(it))
	case uint16:
		it := value.(uint16)
		key = strconv.Itoa(int(it))
	case int32:
		it := value.(int32)
		key = strconv.Itoa(int(it))
	case uint32:
		it := value.(uint32)
		key = strconv.Itoa(int(it))
	case int64:
		it := value.(int64)
		key = strconv.FormatInt(it, 10)
	case uint64:
		it := value.(uint64)
		key = strconv.FormatUint(it, 10)
	case string:
		key = value.(string)
	case []byte:
		key = string(value.([]byte))
	case time.Time:
		if value.(time.Time).Unix() > 0 {
			key = value.(time.Time).Format("2006-01-02 15:04:05")
		} else {
			key = "0000-00-00 00:00:00"
		}
	default:
		newValue, _ := json.Marshal(value)
		key = string(newValue)
	}
	return key
}

func filterFormFooter(infoUrl string) template.HTML {
	components := template2.Get(config.GetTheme())
	col1 := components.Col().SetSize(types.SizeMD(2)).GetContent()
	btn1 := components.Button().SetType("submit").
		AddClass("submit").
		SetContent(icon.Icon(icon.Search, 2) + language.GetFromHtml("search")).
		SetThemePrimary().
		SetSmallSize().
		SetOrientationLeft().
		SetLoadingText(icon.Icon(icon.Spinner, 1) + language.GetFromHtml("search")).
		GetContent()
	btn2 := components.Button().SetType("reset").
		AddClass("reset").
		SetContent(icon.Icon(icon.Undo, 2) + language.GetFromHtml("reset")).
		SetThemeDefault().
		SetOrientationLeft().
		SetSmallSize().
		SetHref(infoUrl).
		SetMarginLeft(12).
		GetContent()
	col2 := components.Col().SetSize(types.SizeMD(8)).
		SetContent(btn1 + btn2).GetContent()
	return col1 + col2
}

func ShowTable(ctx *context.Context, params parameter.Parameters, panel table.Table, infoURL string, ifFoot bool) template.HTML {
	registerPath(panel.GetInfo())
	authHandler := auth.Middleware(db.GetConnection(common.AdminEngine.Services))
	for _, cb := range panel.GetInfo().Callbacks {
		if cb.Value[constant.ContextNodeNeedAuth] == 1 {
			common.AdminEngine.AdminPlugin().GetAddOperationFn()(context.Node{
				Path:     cb.Path,
				Method:   cb.Method,
				Handlers: append([]context.Handler{authHandler}, cb.Handlers...),
			})
		} else {
			common.AdminEngine.AdminPlugin().GetAddOperationFn()(context.Node{Path: cb.Path, Method: cb.Method, Handlers: cb.Handlers})
		}
	}
	panelInfo, _ := panel.GetData(params.WithIsAll(false))
	components := template2.Get(config.GetTheme())
	var (
		actionJs  template.JS
		body      template.HTML
		dataTable types.DataTableAttribute

		user          = auth.Auth(ctx)
		info          = panel.GetInfo()
		actionBtns    = info.Action
		allActionBtns = info.ActionButtons.CheckPermissionWhenURLAndMethodNotEmpty(user)
	)

	if actionBtns == template2.HTML("") && len(allActionBtns) > 0 {
		if info.ActionButtonFold {
			var content template.HTML
			content, actionJs = allActionBtns.Content()
			actionBtns = html.Div(html.Div(
				html.A(icon.Icon(icon.EllipsisV),
					html.M{"color": "#676565"},
					html.M{"href": "#"},
				), html.M{"cursor": "pointer", "width": "100%"}, html.M{"class": "dropdown-toggle", "data-toggle": "dropdown"})+
				html.Ul(content,
					html.M{"min-width": "20px !important", "left": "-32px", "overflow": "hidden"},
					html.M{"class": "dropdown-menu", "role": "menu", "aria-labelledby": "dLabel"}),

				html.M{"text-align": "center"}, html.M{"class": "dropdown"})
		} else {
			actionBtns, actionJs = allActionBtns.Content()
		}
	} else {
		info.ActionButtonFold = false
	}

	btns, btnsJs := info.Buttons.CheckPermissionWhenURLAndMethodNotEmpty(user).Content()

	if info.TabGroups.Valid() {

		dataTable = components.DataTable().
			SetThead(panelInfo.Thead)

		var (
			tabsHtml    = make([]map[string]template.HTML, len(info.TabHeaders))
			infoListArr = panelInfo.InfoList.GroupBy(info.TabGroups)
			theadArr    = panelInfo.Thead.GroupBy(info.TabGroups)
		)
		for key, header := range info.TabHeaders {
			tabsHtml[key] = map[string]template.HTML{
				"title": template2.HTML(header),
				"content": components.DataTable().
					SetInfoList(infoListArr[key]).
					SetButtons(btns).
					SetActionJs(btnsJs + actionJs).
					SetHasFilter(len(panelInfo.FilterFormData) > 0).
					SetAction(actionBtns).
					SetActionFold(info.ActionButtonFold).
					SetIsTab(key != 0).
					SetPrimaryKey(panel.GetPrimaryKey().Name).
					SetThead(theadArr[key]).
					SetHideRowSelector(info.IsHideRowSelector).
					SetLayout(info.TableLayout).
					SetSortUrl(params.GetFixedParamStrWithoutSort()).
					SetInfoUrl(infoURL).
					GetContent(),
			}
		}
		body = components.Tabs().SetData(tabsHtml).GetContent()
	} else {
		dataTable = components.DataTable().
			SetInfoList(panelInfo.InfoList).
			SetButtons(btns).
			SetLayout(info.TableLayout).
			SetActionJs(btnsJs + actionJs).
			SetAction(actionBtns).
			SetHasFilter(len(panelInfo.FilterFormData) > 0).
			SetPrimaryKey(panel.GetPrimaryKey().Name).
			SetThead(panelInfo.Thead).
			SetActionFold(info.ActionButtonFold).
			SetHideRowSelector(info.IsHideRowSelector).
			SetHideFilterArea(info.IsHideFilterArea).
			SetInfoUrl(infoURL).
			SetSortUrl(params.GetFixedParamStrWithoutSort())
		body = dataTable.GetContent()
	}

	isNotIframe := ctx.Query(constant.IframeKey) != "true"
	paginator := panelInfo.Paginator

	if !isNotIframe {
		paginator = paginator.SetEntriesInfo("")
	}

	boxModel := components.Box().
		SetBody(body).
		SetNoPadding().
		SetHeader(dataTable.GetDataTableHeader() + info.HeaderHtml).
		WithHeadBorder().
		SetIframeStyle(!isNotIframe)
	if ifFoot {
		boxModel.SetFooter(paginator.GetContent() + info.FooterHtml)
	}
	if len(panelInfo.FilterFormData) > 0 {
		boxModel = boxModel.SetSecondHeaderClass("filter-area").
			SetSecondHeader(components.Form().
				SetContent(panelInfo.FilterFormData).
				SetInputWidth(info.FilterFormInputWidth).
				SetHeadWidth(info.FilterFormHeadWidth).
				SetMethod("get").
				SetLayout(info.FilterFormLayout).
				SetUrl(infoURL). //  + params.GetFixedParamStrWithoutColumnsAndPage()
				SetHiddenFields(map[string]string{
					form.NoAnimationKey: "true",
				}).
				SetOperationFooter(filterFormFooter(infoURL)).
				GetContent())
	}

	content := boxModel.GetContent()

	if info.Wrapper != nil {
		content = info.Wrapper(content)
	}

	interval := make([]int, 0)
	autoRefresh := info.AutoRefresh != uint(0)
	if autoRefresh {
		interval = append(interval, int(info.AutoRefresh))
	}
	return content
}

func registerPath(InfoPanel *types.InfoPanel) {
	authHandler := auth.Middleware(db.GetConnection(common.AdminEngine.Services))
	for _, cb := range InfoPanel.Callbacks {
		if cb.Value[constant.ContextNodeNeedAuth] == 1 {
			common.AdminEngine.AdminPlugin().GetAddOperationFn()(context.Node{
				Path:     cb.Path,
				Method:   cb.Method,
				Handlers: append([]context.Handler{authHandler}, cb.Handlers...),
			})
		} else {
			common.AdminEngine.AdminPlugin().GetAddOperationFn()(context.Node{Path: cb.Path, Method: cb.Method, Handlers: cb.Handlers})
		}
	}
}

// 获取两个时间相差的天数，0表同一天，正数表t1>t2，负数表t1<t2
func GetDiffDays(t1, t2 time.Time) int {
	t1 = time.Date(t1.Year(), t1.Month(), t1.Day(), 0, 0, 0, 0, time.Local)
	t2 = time.Date(t2.Year(), t2.Month(), t2.Day(), 0, 0, 0, 0, time.Local)

	return int(t1.Sub(t2).Hours() / 24)
}

// 获取t1和t2的相差天数，单位：秒，0表同一天，正数表t1>t2，负数表t1<t2
func GetDiffDaysBySecond(t1, t2 int64) int {
	time1 := time.Unix(t1, 0)
	time2 := time.Unix(t2, 0)
	// 调用上面的函数
	return GetDiffDays(time1, time2)
}

func CreateUnionTableSqlByDays(startTime, endTime string, tableName string,
	gormDB *gorm.DB, structs interface{}, where string, all bool) (sql1, sql2 string) {
	if !all && startTime == "" && endTime == "" {
		return "", ""
	}
	var dbName string
	gormDB.Raw("SELECT DATABASE();").Scan(&dbName)
	if dbName == "" {
		return "", ""
	}
	logTimeStart, _ := time.ParseInLocation("2006-01-02 15:04:05", startTime, time.Local)
	logTimeEnd, _ := time.ParseInLocation("2006-01-02 15:04:05", endTime, time.Local)
	tableNames := getCachedTableNames(gormDB, tableName)
	////select * from table1 union all select * from table2
	fileds := GetFields2Struct(structs)
	for _, v := range tableNames {
		arr := strings.Split(v, "_")
		t, _ := time.ParseInLocation("20060102", arr[len(arr)-1], time.Local)
		if t.After(logTimeStart) && t.Before(logTimeEnd) || all || t.Equal(logTimeStart) || t.Equal(logTimeEnd) ||
			DiffNatureDays(t.Unix(), logTimeStart.Unix()) == 0 || DiffNatureDays(t.Unix(), logTimeEnd.Unix()) == 0 {
			sql1 += " select " + fileds + " from " + v + where + " UNION ALL "
			sql2 += " select count(*) from " + v + where + " UNION ALL "
		}
	}
	if len(sql1) > len(" UNION ALL ") {
		sql1 = strings.TrimSuffix(sql1, " UNION ALL ")
	}
	if len(sql2) > len(" UNION ALL ") {
		sql2 = strings.TrimSuffix(sql2, " UNION ALL ")
	}
	return sql1, sql2
}

func CreateUnionTableSqlByDays4Money(startTime, endTime string, tableName string,
	gormDB *gorm.DB, structs interface{}, where string, all bool) (sql1, sql2 string) {
	if !all && startTime == "" && endTime == "" {
		return "", ""
	}
	var dbName string
	gormDB.Raw("SELECT DATABASE();").Scan(&dbName)
	if dbName == "" {
		return "", ""
	}
	logTimeStart, _ := time.ParseInLocation("2006-01-02 15:04:05", startTime, time.Local)
	logTimeEnd, _ := time.ParseInLocation("2006-01-02 15:04:05", endTime, time.Local)
	tableNames := getCachedTableNames(gormDB, tableName)
	////select * from table1 union all select * from table2
	fileds := GetFields2Struct(structs)
	for _, v := range tableNames {
		arr := strings.Split(v, "_")
		t, _ := time.ParseInLocation("20060102", arr[len(arr)-1], time.Local)
		if t.After(logTimeStart) && t.Before(logTimeEnd) || all || t.Equal(logTimeStart) || t.Equal(logTimeEnd) {
			sql1 += " select " + fileds + " from " + v + where + " UNION ALL "
			sql2 += " select count(*) from " + v + where + " UNION ALL "
		}
	}
	if len(sql1) > len(" UNION ALL ") {
		sql1 = strings.TrimSuffix(sql1, " UNION ALL ")
	}
	if len(sql2) > len(" UNION ALL ") {
		sql2 = strings.TrimSuffix(sql2, " UNION ALL ")
	}
	temp4Money := GetFields2StructOnly4Money(structs)

	sql1 = "select " + temp4Money + " from ( " + sql1 + ") as A group by A.FlowType , A.PlayerId"

	return sql1, sql2
}

func GetFields2Struct(structs interface{}) (fields string) {
	ret := reflect.TypeOf(structs)
	if ret.Kind() != reflect.Struct {
		return ""
	}
	for i := 0; i < ret.NumField(); i++ {
		fields += ret.Field(i).Tag.Get("db") + ","
	}
	if len(fields) > len(",") {
		fields = strings.TrimSuffix(fields, ",")
	}
	return fields
}

func GetFields2StructOnly4Money(structs interface{}) (fields string) {
	ret := reflect.TypeOf(structs)
	if ret.Kind() != reflect.Struct {
		return ""
	}
	for i := 0; i < ret.NumField(); i++ {
		if ret.Field(i).Tag.Get("db") == "moneyValue" {
			fields += "SUM(" + ret.Field(i).Tag.Get("db") + ") as " + ret.Field(i).Tag.Get("db") + ","
		} else if ret.Field(i).Tag.Get("db") == "flowType" {
			fields += ret.Field(i).Tag.Get("db") + ","
		} else if ret.Field(i).Tag.Get("db") == "playerId" {
			fields += ret.Field(i).Tag.Get("db") + ","
		} else {
			fields += "any_value(" + ret.Field(i).Tag.Get("db") + ") as " + ret.Field(i).Tag.Get("db") + ","
		}
	}
	if len(fields) > len(",") {
		fields = strings.TrimSuffix(fields, ",")
	}
	return fields
}

func ReadItemListFromByte(items []string, ifBag bool) (list []def.ItemInfo) {
	for _, v := range items {
		itemInfo := def.ItemInfo{}
		var unPacker = NewTextUnpacker(v)
		if !unPacker.IsEmpty() {
			if ifBag {
				itemInfo.Slot = uint32(unPacker.UnpackUint())
			}
			loadInstItem(&itemInfo, unPacker)
			//itemInfo.ItemGuid = uint32(unPacker.UnpackUint())
			//itemInfo.ItemTypeID = uint32(unPacker.UnpackUint())
			//itemInfo.ItemCount = uint32(unPacker.UnpackUint())
			//itemInfo.ItemOwner = uint32(unPacker.UnpackUint())
			//temp := unPacker.UnpackUint()
			//itemInfo.ItemStatus = int32(GetStatus(uint32(temp)))
			list = append(list, itemInfo)
		}
	}
	return list
}

func loadInstItem(itemInfo *def.ItemInfo, textUnPacker *TextUnpacker) {
	itemInfo.ItemGuid = uint32(textUnPacker.UnpackUint())
	itemInfo.ItemTypeID = uint32(textUnPacker.UnpackUint())
	itemInfo.ItemCount = uint32(textUnPacker.UnpackUint())
	itemInfo.ItemOwner = uint32(textUnPacker.UnpackUint())
	temp := textUnPacker.UnpackUint()
	itemInfo.ItemStatus = int32(GetStatus(uint32(temp)))
	if BIT_ISSET(uint32(temp), def.Expirable) {
		itemInfo.ItemExpireTime = textUnPacker.UnpackInt()
	}
	if BIT_ISSET(uint32(temp), def.Consumable) {
		itemInfo.ItemConsumableUseCount = uint32(textUnPacker.UnpackUint())
	}
	if BIT_ISSET(uint32(temp), def.Transaction) {
		itemInfo.ItemTransactionCount = uint32(textUnPacker.UnpackUint())
	}
	if BIT_ISSET(uint32(temp), def.Equipment) {
		itemInfo.ItemEquipIsHighestGrade = textUnPacker.UnpackBool()
		itemInfo.ItemEquipIsUnbindEnable = textUnPacker.UnpackBool()
		itemInfo.ItemEquipForgeLevel = uint32(textUnPacker.UnpackUint())
		itemInfo.ItemEquipAdditionLevel = uint32(textUnPacker.UnpackUint())
		itemInfo.ItemEquipRegenerateAttrIdx = uint8(textUnPacker.UnpackUint())
		itemInfo.ItemEquipRegenerateAttrVal = float32(textUnPacker.UnpackFloat())
		itemInfo.ItemDurability = uint32(textUnPacker.UnpackUint())
		itemInfo.ItemWearDuration = uint32(textUnPacker.UnpackUint())
		itemInfo.ItemEquipApplyWeight = uint32(textUnPacker.UnpackUint())
		var size uint8
		size = uint8(uint32(textUnPacker.UnpackUint()))
		for i := uint8(0); i < size; i++ {
			itemInfo.ItemEquipSpellIdxs = append(itemInfo.ItemEquipSpellIdxs, uint8(textUnPacker.UnpackUint()))
		}
		size = 0
		size = uint8(uint32(textUnPacker.UnpackUint()))
		for i := uint8(0); i < size; i++ {
			itemInfo.ItemEquipExcellentAttrIdxs = append(itemInfo.ItemEquipExcellentAttrIdxs, uint8(textUnPacker.UnpackUint()))
		}
	}
}

func GetEligibleEquipCount(attrCount int, items []string) (count int) {
	for _, v := range items {
		itemInfo := def.ItemInfo{}
		var unPacker = NewTextUnpacker(v)
		if !unPacker.IsEmpty() {
			itemInfo.Slot = uint32(unPacker.UnpackUint())
			if itemInfo.Slot == 31 {
				print("fxxk")
			}

			loadInstItem(&itemInfo, unPacker)
			if len(itemInfo.ItemEquipExcellentAttrIdxs) >= attrCount {
				count++
			}
		}
	}
	return count
}

func GetAimDayStartEnd(dist int) (string, string) {
	date := time.Now()
	date = date.Add(-(time.Hour * time.Duration(24*dist)))
	dateString := date.Format("2006-01-02")

	//获取当前时区
	//loc, _ := time.LoadLocation("Local")

	//日期当天0点时间戳(拼接字符串)
	startDate := dateString + " 00:00:00"
	//startTime, _ := time.ParseInLocation("2006-01-02 15:04:05", startDate, loc)

	//日期当天23时59分时间戳
	endDate := dateString + " 23:59:59"
	//end, _ := time.ParseInLocation("2006-01-02 15:04:05", endDate, loc)

	return startDate, endDate
}

func GetAimDayStartEndTime(dist int) (time.Time, time.Time) {
	start, end := GetAimDayStartEnd(dist)
	startTime, _ := time.ParseInLocation("2006-01-02 15:04:05", start, time.Local)
	endTime, _ := time.ParseInLocation("2006-01-02 15:04:05", end, time.Local)
	return startTime, endTime
}

func InitBaseDataBase() {
	common.GormDBList = map[string]*gorm.DB{}
	cfg := common.Cfg_yml.Databases["go_manager"]
	goManager := OpenDB(&def.DataBase{
		User:         cfg.User,
		Pwd:          cfg.Pwd,
		Host:         cfg.Host,
		Port:         cfg.Port,
		DataBaseName: cfg.Name})
	common.GormDBList["go_manager"] = goManager

	cfg = common.Cfg_yml.Databases["db_global"]
	dbGlobal := OpenDB(&def.DataBase{
		User:         cfg.User,
		Pwd:          cfg.Pwd,
		Host:         cfg.Host,
		Port:         cfg.Port,
		DataBaseName: cfg.Name})
	common.GormDBList["db_global"] = dbGlobal

	cfg = common.Cfg_yml.Databases["db_world"]
	dbWorld := OpenDB(&def.DataBase{
		User:         cfg.User,
		Pwd:          cfg.Pwd,
		Host:         cfg.Host,
		Port:         cfg.Port,
		DataBaseName: cfg.Name})
	common.GormDBList["db_world"] = dbWorld

	cfg = common.Cfg_yml.Databases["default"]
	dbDefault := OpenDB(&def.DataBase{
		User:         cfg.User,
		Pwd:          cfg.Pwd,
		Host:         cfg.Host,
		Port:         cfg.Port,
		DataBaseName: cfg.Name})
	common.GormDBList["default"] = dbDefault

	cfg = common.Cfg_yml.Databases["go_backup"]
	go_backup := OpenDB(&def.DataBase{
		User:         cfg.User,
		Pwd:          cfg.Pwd,
		Host:         cfg.Host,
		Port:         cfg.Port,
		DataBaseName: cfg.Name})
	common.GormDBList["go_backup"] = go_backup

	cfg = common.Cfg_yml.Databases["db_login_log"]
	db_login_log := OpenDB(&def.DataBase{
		User:         cfg.User,
		Pwd:          cfg.Pwd,
		Host:         cfg.Host,
		Port:         cfg.Port,
		DataBaseName: cfg.Name})
	common.GormDBList["db_login_log"] = db_login_log

	common.ServerDBList = make(map[int64]map[string]*gorm.DB, 0)
}

func OpenDB(pconfig *def.DataBase) *gorm.DB {
	var dbURI string
	var dialector gorm.Dialector
	dbURI = fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=true",
		pconfig.User,
		pconfig.Pwd,
		pconfig.Host,
		pconfig.Port,
		pconfig.DataBaseName)
	dialector = mysql.New(mysql.Config{
		DSN:                       dbURI, // data source name
		DefaultStringSize:         256,   // default size for string fields
		DisableDatetimePrecision:  true,  // disable datetime precision, which not supported before MySQL 5.6
		DontSupportRenameIndex:    true,  // drop & create when rename index, rename index not supported before MySQL 5.7, MariaDB
		DontSupportRenameColumn:   true,  // `change` when rename column, rename column not supported before MySQL 8, MariaDB
		SkipInitializeWithVersion: false, // auto configure based on currently MySQL version
	})
	newLogger := gormlogger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags),
		gormlogger.Config{
			SlowThreshold:             200 * time.Millisecond,
			LogLevel:                  gormlogger.Warn,
			IgnoreRecordNotFoundError: true,
			Colorful:                  false,
		},
	)
	conn, err := gorm.Open(dialector, &gorm.Config{Logger: newLogger})
	if err != nil {
		log.Print(err.Error())
		return nil
	}
	sqlDB, err := conn.DB()
	if err != nil {
		log.Print("connect db server failed.")
		return nil
	}
	sqlDB.SetMaxIdleConns(10)                   // SetMaxIdleConns sets the maximum number of connections in the idle connection pool.
	sqlDB.SetMaxOpenConns(100)                  // SetMaxOpenConns sets the maximum number of open connections to the database.
	sqlDB.SetConnMaxLifetime(time.Second * 600) // SetConnMaxLifetime sets the maximum amount of time a connection may be reused.
	return conn
}

func InitServerDataBases() {
	databaseCfgs := []def.DatabaseCfg{}
	go_manager := GetBaseGormDB("go_manager")
	go_manager.Select("*").Table("databasecfg").Scan(&databaseCfgs)
	common.ServerDBList = map[int64]map[string]*gorm.DB{}
	for _, v := range databaseCfgs {
		tempMap := map[string]*gorm.DB{}
		dataBase := def.DataBase{}
		var temp = String2Bytes(v.LogDataBase)
		json.Unmarshal(temp, &dataBase)
		db_log := OpenDB(&dataBase)
		if db_log != nil {
			tempMap["db_log"] = db_log
		}

		dataBase = def.DataBase{}
		var temp1 = String2Bytes(v.CharDataBase)
		json.Unmarshal(temp1, &dataBase)
		db_char := OpenDB(&dataBase)
		if db_char != nil {
			tempMap["db_char"] = db_char
		}
		common.ServerDBList[v.ServerId] = tempMap
	}
}

func GetBaseGormDB(dbName string) *gorm.DB {
	sqlDB, err := common.GormDBList[dbName].DB()
	if err != nil {
		fmt.Errorf("connect db server failed")
		InitBaseDataBase()
	}
	if err := sqlDB.Ping(); err != nil {
		sqlDB.Close()
		InitBaseDataBase()
	}
	return common.GormDBList[dbName]
}

func GetDeployTargetByServerId(serverId string) (host string, port string, name string) {
	port = "10005"
	tmp := make(map[string]interface{})
	db_global := GetBaseGormDB("db_global")
	db_global.Table("t_game_servers").Where("Id = ?", serverId).Select("internalIP", "internalName").Scan(&tmp)
	host = Strval(tmp["internalIP"])
	name = Strval(tmp["internalName"])
	return
}

func GetMinCombinedServerID(gsId int) (minId int, err error) {
	if cached, ok := combinedServerCache.Load(gsId); ok {
		return cached.(int), nil
	}

	t := time.Now()
	defer func() {
		logger.Infof("[GetMinCombinedServerID] cache=miss gsId=%d elapsed=%v", gsId, time.Since(t))
	}()

	type serverCut struct {
		ExternalIP   string `gorm:"column:externalIP" db:"externalIP"`
		ExternalPort int    `gorm:"column:externalPort" db:"externalPort"`
	}
	var cut1 = serverCut{}
	db_global := GetBaseGormDB("db_global")
	db_global.Table("t_game_servers").Select("externalIP,externalPort").Where("Id = ?", gsId).Scan(&cut1)
	if cut1.ExternalIP == "" || cut1.ExternalPort == 0 {
		return 0, errors.New("server not exist")
	}

	db_global.Table("t_game_servers").Select("min(id)").
		Where("externalIP = ? and externalPort = ?", cut1.ExternalIP, cut1.ExternalPort).
		Group("externalIP,externalPort").Scan(&minId)
	if minId == 0 {
		return 0, errors.New("server not exist")
	}
	combinedServerCache.Store(gsId, minId)
	return minId, nil
}

func GetServerGormDB(db_name string, serverID int64) *gorm.DB {
	temp, err := GetMinCombinedServerID(int(serverID))
	if err != nil {
		return nil
	}
	logger.Info(fmt.Sprintf(" !!!! tempserverIdid :%d > serverId: %d", temp, serverID))
	if int64(temp) > serverID {
		logger.Error(fmt.Sprintf("something wrong!!!! CombinedId :%d > serverId: %d", temp, serverID))
	}
	serverID = int64(temp)

	var dataBaseName = ""
	if db_name == "db_char" {
		dataBaseName = common.MyCfg.Api["BassCharName"]
	} else if db_name == "db_log" {
		dataBaseName = common.MyCfg.Api["BassLogName"]
	}

	if common.Cfg_yml.Debug {
		serverID = 0
	} else {
		dataBaseName += "_s" + strconv.Itoa(int(serverID))
	}
	logger.Info(fmt.Sprintf(" !!!! dbname :%s > serverId: %d", dataBaseName, serverID))
	var serverDBMap, isOk = common.ServerDBList[serverID]
	if !isOk {
		common.ServerDBList[serverID] = make(map[string]*gorm.DB, 0)
	}
	var hadInit = false
	if !isOk || serverDBMap == nil {
		hadInit = true
	}
	if serverDBMap != nil {
		gorm_db, isOk := serverDBMap[db_name]
		if !isOk || gorm_db == nil {
			hadInit = true
		}
	}
	if hadInit {
		if db_name == "db_char" {
			common.ServerDBList[serverID][db_name] = OpenDB(&def.DataBase{
				User:         common.MyCfg.Api["BaseCharUser"],
				Pwd:          common.MyCfg.Api["BaseCharPass"],
				Host:         common.MyCfg.Api["BaseCharHost"],
				Port:         common.MyCfg.Api["BaseCharProt"],
				DataBaseName: dataBaseName})
		} else if db_name == "db_log" {
			common.ServerDBList[serverID][db_name] = OpenDB(&def.DataBase{
				User:         common.MyCfg.Api["BaseLogUser"],
				Pwd:          common.MyCfg.Api["BaseLogPass"],
				Host:         common.MyCfg.Api["BaseLogHost"],
				Port:         common.MyCfg.Api["BaseLogProt"],
				DataBaseName: dataBaseName})
		}
	}
	return common.ServerDBList[serverID][db_name]
}

func UpdateGormDBConn(cfg *def.DatabaseCfg) {
	tempMap := map[string]*gorm.DB{}
	dataBase := def.DataBase{}
	var temp = String2Bytes(cfg.LogDataBase)
	json.Unmarshal(temp, &dataBase)
	db_log := OpenDB(&dataBase)
	tempMap["db_log"] = db_log
	dataBase = def.DataBase{}
	var temp1 = String2Bytes(cfg.CharDataBase)
	json.Unmarshal(temp1, &dataBase)
	db_char := OpenDB(&dataBase)
	tempMap["db_char"] = db_char
	common.ServerDBList[cfg.ServerId] = tempMap
}

func BIT_SET(value uint32, pos int) uint32 {
	value |= (1) << pos
	return value
}

func BIT_ISSET(value uint32, pos int) bool {
	return (value & (1 << pos)) != 0
}

func GetStatus(temp uint32) int {
	for i := 0; i < def.Count; i++ {
		if BIT_ISSET(temp, i) {
			return i
		}
	}
	return -1 //-1为普通道具
}

func GetServerNameStrings() map[int]string {
	type ServerCut struct {
		Id        int    `gorm:"primaryKey;column:Id"`
		LogicName string `gorm:"column:logicName"`
	}
	db_global := GetBaseGormDB("db_global")
	temp := []ServerCut{}
	db_global.Raw("SELECT Id,logicName FROM t_game_servers WHERE id IN (SELECT MIN(id) FROM t_game_servers GROUP BY externalIP,externalPort)").Scan(&temp)
	opStrings := make(map[int]string)
	for _, v := range temp {
		opStrings[v.Id] = v.LogicName
	}
	return opStrings
}

func GetServerListWithAll() types.FieldOptions {
	temp := GetServerNameStrings()
	myOptions := types.FieldOptions{}
	myOptions = append(myOptions,
		types.FieldOption{Text: gotext.Get("全服"), Value: "0"})
	for k, v := range temp {
		tempOp := types.FieldOption{}
		tempOp.Value = strconv.Itoa(k)
		tempOp.Text = "[" + tempOp.Value + "]" + v
		myOptions = append(myOptions, tempOp)
	}
	return myOptions
}

func GetServerList() types.FieldOptions {
	temp := GetServerNameStrings()
	myOptions := types.FieldOptions{}
	for k, v := range temp {
		tempOp := types.FieldOption{}
		tempOp.Value = strconv.Itoa(k)
		tempOp.Text = "[" + tempOp.Value + "]" + v
		myOptions = append(myOptions, tempOp)
	}
	return myOptions
}

func GetItemNameStrings() map[int]string {
	itemNameCacheOnce.Do(func() {
		itemNameCache = make(map[int]string)

		filed := "stringCN"
		if common.MyCfg.Language == "vi" {
			filed = "stringVN"
		}
		itemIDs := make([]int, 0)
		db_world := GetBaseGormDB("db_world")
		db_world.Table("item_prototype").
			Select("itemTypeID").Scan(&itemIDs)
		stringIds := make([]int, 0)
		for _, v := range itemIDs {
			stringIds = append(stringIds, 20<<32+v)
		}
		stringS := make([]map[string]interface{}, 0)
		db_world.Table("string_text_list").
			Select(filed+",stringID").
			Where("stringID in ?", stringIds).
			Scan(&stringS)
		for _, v := range stringS {
			itemID, _ := strconv.Atoi(Strval(v["stringID"]))
			itemString := Strval(v[filed])
			itemID -= 20 << 32
			if itemString == "" {
				itemString = strconv.Itoa(itemID)
			}
			itemNameCache[itemID] = itemString
		}
	})
	return itemNameCache
}

func GetItemNameOption() types.FieldOptions {
	ItemNameOP := types.FieldOptions{}
	ItemNameString := GetItemNameStrings()
	for k, v := range ItemNameString {
		ItemNameOP = append(ItemNameOP, types.FieldOption{Text: v, Value: strconv.Itoa(k)})
	}
	return ItemNameOP
}

func GetItemNameOptionAndID() types.FieldOptions {
	ItemNameOP := types.FieldOptions{}
	ItemNameString := GetItemNameStrings()
	for k, v := range ItemNameString {
		ItemNameOP = append(ItemNameOP, types.FieldOption{Text: v + "(" + strconv.Itoa(k) + ")", Value: strconv.Itoa(k)})
	}
	return ItemNameOP
}

func GetTpPointStrings() map[int]string {
	transmitStrings := make(map[int]string)

	filed := "stringCN"
	if common.MyCfg.Language == "vi" {
		filed = "stringVN"
	}
	transmitIDs := make([]int, 0)
	db_world := GetBaseGormDB("db_world")
	db_world.Table("auto_map_transmit").
		Select("Id").Scan(&transmitIDs)
	stringIds := make([]int, 0)
	for _, v := range transmitIDs {
		stringIds = append(stringIds, 1016<<32+v)
	}

	stringS := make([]map[string]interface{}, 0)
	db_world.Table("string_text_list").
		Select(filed+",stringID").
		Where("stringID in ?", stringIds).
		Scan(&stringS)
	for _, v := range stringS {
		transmitID, _ := strconv.Atoi(Strval(v["stringID"]))
		transmitString := Strval(v[filed])
		transmitID -= 1016 << 32
		if transmitString == "" {
			transmitString = strconv.Itoa(transmitID)
		}
		transmitStrings[transmitID] = transmitString
	}
	return transmitStrings
}

func GetTpPointOptionAndID() types.FieldOptions {
	transmitOP := types.FieldOptions{}
	transmitString := GetTpPointStrings()
	for k, v := range transmitString {
		transmitOP = append(transmitOP, types.FieldOption{Text: v + "(" + strconv.Itoa(k) + ")", Value: strconv.Itoa(k)})
	}
	return transmitOP
}

func GetFamilyStrings() map[int]string {
	FamilyStrings := make(map[int]string)
	FamilyStrings[0] = gotext.Get("无")
	FamilyStrings[1] = gotext.Get("多普瑞恩家族")
	FamilyStrings[2] = gotext.Get("巴内尔特家族")
	return FamilyStrings
}

func InitConfig(cfg *myConfig.MyConfig) {
	temp := []def.ConfigTable{}
	go_manager := GetBaseGormDB("go_manager")
	go_manager.Table(def.ConfigTable{}.TableName()).Select("*").Scan(&temp)
	cfg.Api = make(map[string]string)
	for _, v := range temp {
		switch v.Key {
		case "centerHost":
			cfg.Center.Host = v.Value
		case "centerPort":
			cfg.Center.Port = v.Value
		case "backupHost":
			cfg.BackUp.Host = v.Value
		case "backupPort":
			cfg.BackUp.Port = v.Value
		case "language":
			cfg.Language = v.Value
		case "EmailAccount":
			cfg.Email.EmailAccount = v.Value
		case "EmailAuthorizationCode":
			cfg.Email.EmailAuthorizationCode = v.Value
		case "EmailServerHost":
			cfg.Email.EmailServerHost = v.Value
		case "EmailServerPort":
			cfg.Email.EmailServerPort = v.Value
		case "SecretKey":
			cfg.SecretKey = v.Value
		default:
			cfg.Api[v.Key] = v.Value
		}
	}
}

func GetBatchNumOption() types.FieldOptions {
	myOp := types.FieldOptions{}
	db_global := GetBaseGormDB("db_global")
	var batchNums = make([]int, 0)
	db_global.Table("t_activation_code").
		Select("DISTINCT batchNum").
		Scan(&batchNums)
	for _, v := range batchNums {
		myOp = append(myOp, types.FieldOption{Text: strconv.Itoa(v), Value: strconv.Itoa(v)})
	}
	return myOp
}

func GetPurchase() map[int]interface{} {
	var purchase = make([]map[string]interface{}, 0)
	db_world := GetBaseGormDB("db_world")
	db_world.Table("auto_purchase").
		Select("ID,buyPrice,buyPriceShow").
		Where("platformID = ?", 2).
		Scan(&purchase)
	var data = make(map[int]interface{})

	for _, v := range purchase {
		id, _ := strconv.Atoi(Strval(v["ID"]))
		buyPrice, _ := strconv.Atoi(Strval(v["buyPrice"]))
		buyPriceShow := Strval(v["buyPriceShow"])
		data[id] = def.PurchaseCut{BuyPrice: buyPrice, BuyPriceShow: buyPriceShow}
	}
	return data
}
func GetPurchaseOp() types.FieldOptions {
	myOps := types.FieldOptions{}
	Purchase := GetPurchase()
	for k, v := range Purchase {
		op, ok := v.(def.PurchaseCut)
		if !ok {
			continue
		}
		myOps = append(myOps, types.FieldOption{Text: op.BuyPriceShow, Value: strconv.Itoa(k)})
	}
	return myOps
}

func MergeMap(m1, m2 map[string]interface{}) (m3 map[string]interface{}) {
	m3 = make(map[string]interface{})
	for i, v := range m1 {
		for j, w := range m2 {
			if i == j {
				m3[i] = w
			} else {
				if _, ok := m3[i]; !ok {
					m3[i] = v
				}
				if _, ok := m3[j]; !ok {
					m3[j] = w
				}
			}
		}
	}
	return m3
}

func GetShowItemFormatStr(list []def.ItemInfo, strS map[int]string) string {
	str := ""
	for _, v := range list {
		if strS[int(v.ItemTypeID)] != "" {
			str += strS[int(v.ItemTypeID)] + "  x" + strconv.Itoa(int(v.ItemCount)) + "\n"
		} else {
			str += strconv.Itoa(int(v.ItemTypeID)) + "  x" + strconv.Itoa(int(v.ItemCount)) + "\n"
		}
	}
	return str
}

func GetActiveServerIds() (ids []int) {
	db_global := GetBaseGormDB("db_global")
	db_global.Raw("SELECT MIN(id) FROM t_game_servers GROUP BY externalIP,externalPort").
		Scan(&ids)
	return ids
}

func VerifyWebArgs(args ...string) string {
	h := md5.New()
	for _, arg := range args {
		io.WriteString(h, arg)
	}
	io.WriteString(h, common.MyCfg.SecretKey)
	return hex.EncodeToString(h.Sum(nil))
}

func ReadExcel(filepath string, configType int) (error, []expconfigdef.ExpConfig) {
	str := expconfigdef.GetTypeStr()

	xlsx, err := excelize.OpenFile(filepath)
	if err != nil {
		fmt.Println(err)
		return err, nil
	}

	rows := xlsx.GetRows("Sheet1")
	expcfgs := make([]expconfigdef.ExpConfig, 0)
	for i, row := range rows {
		if i == 0 {
			for j, info := range row {
				if j == 0 {
					if str[configType] != info {
						return errors.New("ErrConfigType"), nil
					}
				}
			}
			continue
		}
		if i == 1 {
			continue
		}
		expCfg := expconfigdef.ExpConfig{}
		for j, info := range row {
			if j == 0 {
				expCfg.CreatureId, err = strconv.Atoi(info)
				if err != nil {
					return err, nil
				}
			}
			if j == 1 {
				expCfg.CreatureName = info
			}
			if j == 2 {
				expCfg.Coefficient, err = strconv.Atoi(info)
				if err != nil {
					return err, nil
				}
			}
		}
		expcfgs = append(expcfgs, expCfg)
	}
	return nil, expcfgs
}

func GetFileMd5(filename string) string {
	// 文件全路径名
	path := fmt.Sprintf("./%s", filename)
	pFile, err := os.Open(path)
	if err != nil {
		fmt.Errorf("打开文件失败，filename=%v, err=%v", filename, err)
		return ""
	}
	defer pFile.Close()
	md5h := md5.New()
	io.Copy(md5h, pFile)

	return hex.EncodeToString(md5h.Sum(nil))
}

func DiffNatureDays(t1, t2 int64) int {
	if t1 == t2 {
		return -1
	}
	if t1 > t2 {
		t1, t2 = t2, t1
	}

	diffDays := 0
	secDiff := t2 - t1
	if secDiff > (60 * 60 * 24) {
		tmpDays := int(secDiff / (60 * 60 * 24))
		t1 += int64(tmpDays) * (60 * 60 * 24)
		diffDays += tmpDays
	}

	st := time.Unix(t1, 0)
	et := time.Unix(t2, 0)
	dateFormatTpl := "20060102"
	if st.Format(dateFormatTpl) != et.Format(dateFormatTpl) {
		diffDays += 1
	}

	return diffDays
}

func ParseFlowParams4FlowType(flowType uint32, flowParams string, serverID int64) string {
	if flowType < flowDef.LFT_NONE {
		return flowParams
	}

	var params = make([]int, 0)
	err := json.Unmarshal(String2Bytes(flowParams), &params)
	if err != nil {
		return flowParams
	}

	switch flowType {
	case flowDef.LFT_MAIL:
		err, showStr := parseMailId(params, serverID)
		if err != nil {
			print(err.Error())
			return flowParams
		}
		return showStr
	case flowDef.LFT_LOOT:
		err, showStr := parseLoot(params)
		if err != nil {
			print(err.Error())
			return flowParams
		}
		return showStr
	case flowDef.LFT_ACTIVITY:
		err, showStr := parseActivity(params)
		if err != nil {
			print(err.Error())
			return flowParams
		}
		return showStr
	case flowDef.LFT_TRANSACTION_ITEM:
		err, showStr := parseTransaction(params)
		if err != nil {
			print(err.Error())
			return flowParams
		}
		return showStr
	case flowDef.LFT_PERSONAL_SHOP:
		err, showStr := parsePersonalShop(params, serverID)
		if err != nil {
			print(err.Error())
			return flowParams
		}
		return showStr
	default:
		return flowParams
	}
}

func GetMailTypeString() map[int]string {
	var temp = make(map[int]string)
	temp[0] = gotext.Get("普通")
	temp[1] = gotext.Get("系統")
	temp[2] = "GM"
	return temp
}

func parseMailId(flowParams []int, serverID int64) (error, string) {
	if len(flowParams) < 2 {
		return errors.New(gotext.Get("参数有误")), ""
	}
	var temp = make(map[string]interface{})
	if flowParams[0] == 0 {
		temp[gotext.Get("操作")] = gotext.Get("接收")
	} else if flowParams[0] == 1 {
		temp[gotext.Get("操作")] = gotext.Get("发送")
	}
	db_char := GetServerGormDB("db_char", serverID)

	type mailCut struct {
		MailID          int    `gorm:"column:mailID"`
		MailType        int    `gorm:"column:mailType"`
		MailSender      string `gorm:"column:mailSender"`
		MailReceiver    string `gorm:"column:mailReceiver"`
		MailDeliverTime int64  `gorm:"column:mailDeliverTime"`
	}
	var Cut mailCut
	db_char.Table("inst_mail_backup").Select("mailID,mailType,mailSender,mailReceiver,mailDeliverTime").
		Where("mailID = ?", flowParams[1]).Scan(&Cut)
	if db_char.Error != nil {
		return db_char.Error, ""
	}

	temp[gotext.Get("邮件ID")] = Cut.MailID
	temp[gotext.Get("邮件类型")] = GetMailTypeString()[Cut.MailType]
	temp[gotext.Get("发件人ID")] = Cut.MailSender
	temp[gotext.Get("收件人ID")] = Cut.MailReceiver
	temp[gotext.Get("发送时间")] = def.MyTime(time.Unix(Cut.MailDeliverTime, 0))

	marshal, err := json.Marshal(temp)
	if err != nil {
		return err, ""
	}
	return nil, Strval(marshal)
}

func parseLoot(flowParams []int) (error, string) {
	var temp = make(map[string]interface{})
	if len(flowParams) <= 1 {
		return errors.New(gotext.Get("参数有误")), ""
	}
	if len(flowParams) < 3 {
		if flowParams[0] != 0 {
			temp[gotext.Get("怪物生成ID")] = flowParams[0]
		}
		if flowParams[1] != 0 {
			temp[gotext.Get("怪物类型ID")] = flowParams[1]
		}

	} else {
		if flowParams[1] != 0 {
			temp[gotext.Get("怪物生成ID")] = flowParams[1]
		}
		if flowParams[2] != 0 {
			temp[gotext.Get("怪物类型ID")] = flowParams[2]
		}
		if len(flowParams) > 3 {
			switch flowParams[3] {
			case 1:
				if len(flowParams) >= 6 {
					temp[gotext.Get("使用者的剩余数量")] = flowParams[4]
					temp[gotext.Get("道具ID")] = flowParams[5]
				} else {
					return errors.New(gotext.Get("参数有误")), ""
				}
				break
			case 2:
				temp[gotext.Get("来源")] = gotext.Get("血色城堡")
				break
			case 3:
				temp[gotext.Get("来源")] = gotext.Get("幻影寺院")
				break
			case 4:
				temp[gotext.Get("来源")] = gotext.Get("赤色要塞")
				break
			case 5:
				temp[gotext.Get("来源")] = gotext.Get("生魂广场")
				break
			case 6:
				if len(flowParams) >= 5 {
					temp[gotext.Get("来源")] = gotext.Get("丢弃")
					temp[gotext.Get("丢弃者")] = flowParams[4]
				} else {
					return errors.New(gotext.Get("参数有误")), ""
				}
				break
			case 7:
				temp[gotext.Get("来源")] = gotext.Get("死亡掉落")
				break
			}
		}
	}
	marshal, err := json.Marshal(temp)
	if err != nil {
		return err, ""
	}
	return nil, Strval(marshal)
}

func parseActivity(flowParams []int) (error, string) {
	if len(flowParams) < 3 {
		return errors.New(gotext.Get("参数有误")), ""
	}
	var temp = make(map[string]interface{})
	temp[gotext.Get("活动ID")] = flowParams[0]
	temp[gotext.Get("难度ID")] = flowParams[1]
	temp[gotext.Get("活动完成状态ID")] = flowParams[2]
	marshal, err := json.Marshal(temp)
	if err != nil {
		return err, ""
	}
	return nil, Strval(marshal)
}

func parseTransaction(flowParams []int) (error, string) {
	if len(flowParams) < 1 {
		return errors.New(gotext.Get("参数有误")), ""
	}
	var temp = make(map[string]interface{})
	temp[gotext.Get("交易对象playerID")] = flowParams[0]

	marshal, err := json.Marshal(temp)
	if err != nil {
		return err, ""
	}
	return nil, Strval(marshal)
}

func GetPersonShopOpStr() map[int]string {
	var temp = make(map[int]string)
	temp[0] = gotext.Get("上架")
	temp[1] = gotext.Get("下架")
	temp[2] = gotext.Get("购买消耗")
	temp[3] = gotext.Get("上架失败")
	temp[4] = gotext.Get("购买失败")
	return temp
}

func parsePersonalShop(flowParams []int, serverId int64) (error, string) {
	if len(flowParams) < 1 {
		return errors.New(gotext.Get("参数有误")), ""
	}
	var temp = make(map[string]interface{})
	type auctionCut struct {
		Id          int    `gorm:"column:Id"`
		PlayerId    int    `gorm:"column:playerId"`
		ItemTypeIDs string `gorm:"column:itemTypeIDs"`
		ItemCounts  string `gorm:"column:itemCounts"`
		CreateTime  int64  `gorm:"column:createTime"`
	}
	var Cut auctionCut
	switch flowParams[0] {
	case 0:
		temp[gotext.Get("背包格子位置")] = flowParams[1]
		temp[gotext.Get("上架后位置")] = flowParams[2]
		break
	case 1:
		temp[gotext.Get("个人商店记录id")] = flowParams[1]
		break
	case 2:
		temp[gotext.Get("个人商店记录id")] = flowParams[1]
		db_log := GetServerGormDB("db_log", serverId)
		db_log.Table("log_personal_shop").Select("Id,playerId,itemTypeIDs,itemCounts,createTime").
			Where("Id = ? and isBuy = 1", flowParams[1]).Scan(&Cut)
		if db_log.Error != nil {
			return db_log.Error, ""
		}
		temp[gotext.Get("购买玩家ID")] = Cut.PlayerId
		var items = make(map[int]int)
		var tempID = make([]int, 0)
		json.Unmarshal(String2Bytes(Cut.ItemTypeIDs), &tempID)
		var tempNum = make([]int, 0)
		json.Unmarshal(String2Bytes(Cut.ItemCounts), &tempNum)
		for i, v := range tempID {
			items[v] = tempNum[i]
		}
		marshal, _ := json.Marshal(items)
		temp[gotext.Get("道具数量")] = Strval(marshal)
		temp[gotext.Get("操作时间")] = def.MyTime(time.Unix(Cut.CreateTime, 0))
		break
	default:
		temp[gotext.Get("操作")] = GetPersonShopOpStr()[flowParams[0]]
		break
	}

	marshal, err := json.Marshal(temp)
	if err != nil {
		return err, ""
	}
	return nil, Strval(marshal)
}

func InitPprofMonitor() {
	go func() {
		err := http.ListenAndServe(":6060", nil)
		if err != nil {
			logger.Error("funcRetErr=http.ListenAndServe||err=%s", err.Error())
		}
	}()
}

func CreateUnionTableSql(tableName string, gormDB *gorm.DB, args ...string) (sql1, sql2 string) {

	var dbName string
	gormDB.Raw("SELECT DATABASE();").Scan(&dbName)
	if dbName == "" {
		return "", ""
	}

	tableNames := []string{}
	sql := "SELECT table_name FROM information_schema.TABLES WHERE" +
		" TABLE_SCHEMA='" + dbName + "' AND table_name LIKE '%" + tableName + "%' "
	gormDB.Raw(sql).Scan(&tableNames)
	////select * from table1 union all select * from table2

	sql = ""
	for _, v := range tableNames {
		sql += " select playerId,buyPrice,gsId,payId from " + v + " UNION ALL "
	}
	if len(sql) > len(" UNION ALL ") {
		sql = strings.TrimSuffix(sql, " UNION ALL ")
	}
	sql = "( " + sql + " ) as a "
	if len(args) != 0 {
		for _, v := range args {
			sql += v
		}
	}
	//where playerId in ? and gsId = ? group by playerId,gsId

	sql1 = "select gsId,playerId,sum(buyPrice)as buyPrice,payId from" + sql
	sql2 = "select gsId,playerId,count(buyPrice)as rechargeTimes from" + sql
	return sql1, sql2
}

func GetSql4ExcellentAttr(tableName string, gormDB *gorm.DB, args ...string) (sql1 string) {

	var dbName string
	gormDB.Raw("SELECT DATABASE();").Scan(&dbName)
	if dbName == "" {
		return ""
	}

	tableNames := []string{}
	sql := "SELECT table_name FROM information_schema.TABLES WHERE" +
		" TABLE_SCHEMA='" + dbName + "' AND table_name LIKE '%" + tableName + "%' "
	gormDB.Raw(sql).Scan(&tableNames)
	////select * from table1 union all select * from table2

	sql = ""
	for _, v := range tableNames {
		sql += " select playerId,buyPrice,gsId,createTime from " + v + " UNION ALL "
	}
	if len(sql) > len(" UNION ALL ") {
		sql = strings.TrimSuffix(sql, " UNION ALL ")
	}
	sql = "( " + sql + " ) as a "
	if len(args) != 0 {
		for _, v := range args {
			sql += v
		}
	}
	//where playerId in ? and gsId = ? group by playerId,gsId

	sql1 = "select gsId,playerId,sum(buyPrice)as buyPrice ,max(createTime) as createTime from" + sql
	return sql1
}

func GetSql4AccountCount(tableName string, gormDB *gorm.DB, args ...string) (sql1 string) {

	var dbName string
	gormDB.Raw("SELECT DATABASE();").Scan(&dbName)
	if dbName == "" {
		return ""
	}

	tableNames := []string{}
	sql := "SELECT table_name FROM information_schema.TABLES WHERE" +
		" TABLE_SCHEMA='" + dbName + "' AND table_name LIKE '%" + tableName + "%' "
	gormDB.Raw(sql).Scan(&tableNames)
	////select * from table1 union all select * from table2

	sql = ""
	for _, v := range tableNames {
		sql += " select playerId,gsId from " + v + " UNION ALL "
	}
	if len(sql) > len(" UNION ALL ") {
		sql = strings.TrimSuffix(sql, " UNION ALL ")
	}
	sql = "( " + sql + " ) as a "
	if len(args) != 0 {
		for _, v := range args {
			sql += v
		}
	}
	//where playerId in ? and gsId = ? group by playerId,gsId

	sql1 = "select playerId from" + sql
	return sql1
}

func GetDataQuarter(time2 time.Time) int {
	month, _ := strconv.Atoi(time2.Format("01"))
	switch month {
	case 1:
	case 2:
	case 3:
		return 1
	case 4:
	case 5:
	case 6:
		return 2
	case 7:
	case 8:
	case 9:
		return 3
	case 10:
	case 11:
	case 12:
		return 4
	}
	return 0
}

func GetAllAuthName() map[int]string {
	type AuthCut struct {
		Id       int    `gorm:"primaryKey;column:id"`
		UserName string `gorm:"column:username"`
	}
	db_account := GetBaseGormDB("default")
	temp := []AuthCut{}
	db_account.Raw("SELECT id,username FROM goadmin_users").Scan(&temp)
	opStrings := make(map[int]string)
	for _, v := range temp {
		opStrings[v.Id] = v.UserName
	}
	return opStrings
}

func GetFormatStr4TMailData(data *def.MailTransit) string {
	var temp string

	var currencyStr = currencyDef.GetCurrencyStrings()
	var itemWithIds = GetItemNameStrings()

	tempData := String2Bytes(data.Record)
	timingTMail := def.TimingTMail{}
	err := json.Unmarshal(tempData, &timingTMail)
	if err != nil {
		return "解析失败！，请联系管理员"
	}

	var chequeStr string
	if timingTMail.Mail.MailCheques != "" {
		mailAttachCheque := make([]string, 0)
		tempData := String2Bytes(timingTMail.Mail.MailCheques)
		err := json.Unmarshal(tempData, &mailAttachCheque)
		if err != nil {
			return "解析失败！，请联系管理员"
		}

		for _, v := range mailAttachCheque {
			var splitStr = strings.Split(v, ",")
			if len(splitStr) < 2 {
				continue
			}
			var chequeType, _ = strconv.Atoi(splitStr[0])
			if currencyStr[chequeDef.CurrencyToCheque(chequeType)] != "" {
				chequeStr += "货币类型：" + currencyStr[chequeDef.CurrencyToCheque(chequeType)] + "，货币值：" + splitStr[1]
			} else {
				chequeStr += "货币类型：" + splitStr[0] + "，货币值：" + splitStr[1]
			}
		}
	}

	var itemStr string
	if timingTMail.Mail.MailItems != "" {
		mailAttachItem := make([]string, 0)
		tempData := String2Bytes(timingTMail.Mail.MailItems)
		err := json.Unmarshal(tempData, &mailAttachItem)
		if err != nil {
			return "解析失败！，请联系管理员"
		}
		for _, v := range mailAttachItem {
			itemInfo := def.ItemInfo{}
			var unPacker = NewTextUnpacker(v)
			loadInstItem(&itemInfo, unPacker)
			if itemWithIds[int(itemInfo.ItemTypeID)] != "" {
				itemStr += "道具：" + itemWithIds[int(itemInfo.ItemTypeID)] + "，数量：" + strconv.Itoa(int(itemInfo.ItemCount))
			} else {
				itemStr += "道具：" + strconv.Itoa(int(itemInfo.ItemTypeID)) + ",数量：" + strconv.Itoa(int(itemInfo.ItemCount))
			}
		}
	}

	temp = "邮件标题：" + timingTMail.Mail.MailSubject + "\r\n" +
		"邮件正文：" + timingTMail.Mail.MailBody + "\r\n" +
		"发送时间：" + data.DeliverTime + "\r\n" +
		"过期时间：" + data.ExpireTime + "\r\n" +
		"赠送货币：" + chequeStr + "\r\n" +
		"赠送道具：" + itemStr + "\r\n"
	return temp
}
