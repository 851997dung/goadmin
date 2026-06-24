package CustomPages

import (
	"admin/common/def"
	"admin/fusion"
	"github.com/GoAdminGroup/go-admin/context"
	"github.com/GoAdminGroup/go-admin/plugins/admin/modules/response"
	"strconv"
)

func SendComplaint(ctx *context.Context) {
	cType, _ := strconv.Atoi(ctx.Request.FormValue("cType"))
	cFrom, _ := strconv.Atoi(ctx.Request.FormValue("cFrom"))
	gsId, _ := strconv.Atoi(ctx.Request.FormValue("gsId"))
	content := ctx.Request.FormValue("content")
	playerName := ctx.Request.FormValue("playerName")
	uid, _ := strconv.Atoi(ctx.Request.FormValue("uid"))

	complaint := def.Complaint{
		CType:      cType,
		CFrom:      cFrom,
		Content:    content,
		PlayerName: playerName,
		Status:     def.Pending,
		GsId:       gsId,
		Uid:        uid,
	}
	go_manager := fusion.GetBaseGormDB("go_manager")
	go_manager.Create(&complaint)
	response.Ok(ctx)
}
