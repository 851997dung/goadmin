package AutoPages

import (
	"admin/common/def"
	"admin/common/def/flowDef"
	"admin/fusion"
	"encoding/json"
	"github.com/GoAdminGroup/go-admin/context"
	"github.com/GoAdminGroup/go-admin/modules/db"
	"github.com/GoAdminGroup/go-admin/plugins/admin/modules/parameter"
	"github.com/GoAdminGroup/go-admin/plugins/admin/modules/table"
	"github.com/GoAdminGroup/go-admin/template"
	"github.com/GoAdminGroup/go-admin/template/types"
	"github.com/GoAdminGroup/go-admin/template/types/action"
	"github.com/GoAdminGroup/go-admin/template/types/form"
	"github.com/fatih/structs"
	"github.com/leonelquinteros/gotext"
	"math"
	"strconv"
	"strings"
)

func GetLogPlayerItemAddError(ctx *context.Context) table.Table {

	LogPlayerItemAddError := table.NewDefaultTable(table.Config{
		Driver:     "mysql",
		Connection: "db_log",
		CanAdd:     false,
		Editable:   false,
		Deletable:  true,
		Exportable: false,
		PrimaryKey: table.PrimaryKey{
			Type: db.Int,
			Name: "Id",
		},
	})
	itemNameString := fusion.GetItemNameStrings()

	info := LogPlayerItemAddError.GetInfo()

	flowStr := flowDef.GetFlowTypeStrings()
	serverStr := fusion.GetServerNameStrings()
	serverOP := fusion.GetServerList()
	info.HideNewButton().HideEditButton().HideDetailButton().HideRowSelector().HideExportButton()

	if strings.Index(ctx.Request.RequestURI, "/export/") != -1 {
		info.AddField(gotext.Get("服务器名"), "ServerName", db.Varchar)
	}
	info.AddField(gotext.Get("服务器ID"), "IpcServerID", db.Int).
		FieldFilterable(types.FilterType{FormType: form.SelectSingle}).
		FieldFilterOptions(serverOP).
		FieldDisplay(func(model types.FieldModel) interface{} {
			temp, _ := strconv.Atoi(model.Value)
			if serverStr[temp] != "" {
				return serverStr[temp]
			}
			return ""
		})

	info.AddField("Id", "Id", db.Bigint)
	info.AddField(gotext.Get("账号ID"), "AcctId", db.Int).FieldHide()
	info.AddField(gotext.Get("角色ID"), "PlayerId", db.Int).FieldHide()
	info.AddField(gotext.Get("角色昵称"), "PlayerName", db.Varchar)
	info.AddField(gotext.Get("角色等级"), "PlayerLevel", db.Int)
	info.AddField(gotext.Get("VIP等级"), "PlayerVipLevel", db.Int).
		FieldHide()
	info.AddField("PlayerFightValue", "PlayerFightValue", db.Bigint).
		FieldHide()
	info.AddField("PlayerRichValue", "PlayerRichValue", db.Double).
		FieldHide()
	info.AddField(gotext.Get("未成功添加的道具"), "AddFailedItems", db.JSON).
		FieldDisplay(func(model types.FieldModel) interface{} {
			arr := make([]string, 0)
			b := fusion.String2Bytes(model.Value)
			json.Unmarshal(b, &arr)
			var list = fusion.ReadItemListFromByte(arr, false)
			return fusion.GetShowItemFormatStr(list, itemNameString)
		})

	info.AddField(gotext.Get("操作来源"), "FlowType", db.Int).
		FieldDisplay(func(model types.FieldModel) interface{} {
			temp, _ := strconv.Atoi(model.Value)
			if flowStr[temp] == "" {
				return model.Value
			}
			return flowStr[temp]
		})
	info.AddField("FlowParams", "FlowParams", db.JSON).
		FieldHide()
	info.AddField(gotext.Get("记录时间"), "LogTime", db.Datetime).
		FieldSortable()
	info.AddField(gotext.Get("是否处理"), "IfSolve", db.Boolean)
	info.SetTable("log_player_items").SetTitle(gotext.Get("未成功添加道具日志"))

	info.AddActionButton(template.HTML(gotext.Get("发送补偿邮件")), action.Jump("/admin/CompensationMail?Id={{.Id}}"+
		"&AcctId={{(index .Value \"AcctId\").Value}}"+
		"&IpcServerID={{(index .Value \"IpcServerID\").Value}}"+
		"&PlayerId={{(index .Value \"PlayerId\").Value}}"+
		"&PlayerName={{(index .Value \"PlayerName\").Value}}"+
		"&LogTime={{(index .Value \"LogTime\").Value}}"+
		"&AddFailedItems={{(index .Value \"AddFailedItems\")}}"))

	info.SetGetDataFn(
		func(param parameter.Parameters) (data []map[string]interface{}, size int) {
			serverId, _ := strconv.Atoi(ctx.Request.FormValue("IpcServerID"))
			var gsIds = make([]int, 0)
			if serverId == 0 {
				gsIds = fusion.GetActiveServerIds()
			} else {
				gsIds = append(gsIds, serverId)
			}
			for _, v := range gsIds {
				db_log := fusion.GetServerGormDB("db_log", int64(v))
				if db_log == nil {
					continue
				}
				ItemAddError := make([]def.PlayerItemAddError, 0)
				fields := fusion.GetFields2Struct(def.PlayerItemAddError{})
				db_log.Table("log_player_item_add_error").Select(fields).
					Scan(&ItemAddError)
				for i, _ := range ItemAddError {
					ItemAddError[i].Id = ItemAddError[i].Id*100000 + v
					ItemAddError[i].LogTime = strings.Replace(ItemAddError[i].LogTime, "T", " ", -1)
					ItemAddError[i].LogTime = strings.Replace(ItemAddError[i].LogTime, "Z", " ", -1)
				}
				for _, item := range ItemAddError {
					m3 := structs.Map(item)
					m3["IpcServerID"] = v
					data = append(data, m3)
				}
			}
			min := math.Min(float64(param.PageInt*param.PageSizeInt), float64(len(data)))
			return data[((param.PageInt - 1) * param.PageSizeInt):int(min)], len(data)
		})

	info.SetDeleteFn(func(idArr []string) error {
		idList := make(map[int][]int)
		for _, v := range idArr {
			Id, _ := strconv.Atoi(v)
			if idList[Id%100000] == nil {
				idList[Id%100000] = make([]int, 0)
			}
			idList[Id%100000] = append(idList[Id%100000], Id/100000)
		}

		for serverId, ids := range idList {
			db_log := fusion.GetServerGormDB("db_log", int64(serverId))
			db_log.Delete(def.PlayerItemAddError{}, "id in ?", ids)
		}
		return nil
	})

	return LogPlayerItemAddError
}
