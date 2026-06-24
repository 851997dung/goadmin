package pages

import (
	"admin/common"
	"admin/pages/AutoPages"
	"admin/pages/CustomPages"
)

func RegisterPath() {
	common.AdminEngine.HTML("GET", "/admin", CustomPages.GetIndex)
	common.AdminEngine.HTML("GET", "/admin/SaveConfig", CustomPages.SaveConfig)

	common.AdminEngine.HTML("GET", "/admin/MailPage", CustomPages.GetMailPage)
	common.AdminEngine.HTML("GET", "/admin/CompensationMail", CustomPages.GetCompensationMail)
	common.AdminEngine.HTML("GET", "/admin/SendMail", CustomPages.SendMail)
	common.AdminEngine.HTML("GET", "/admin/SendCompensationMail", CustomPages.SendCompensationMail)

	common.AdminEngine.HTML("GET", "/admin/ServerDetail", CustomPages.GetServerDetail)
	common.AdminEngine.HTML("GET", "/admin/EditServerDatabaseCfg", CustomPages.EditServerDatabaseCfg)
	common.AdminEngine.HTML("GET", "/admin/UpdateServerDatabaseCfg", CustomPages.UpdateServerDatabaseCfg)
	common.AdminEngine.HTML("GET", "/admin/EditServerDatabaseBackupCfg", CustomPages.EditServerDatabaseBackupCfg)
	common.AdminEngine.HTML("GET", "/admin/UpdateServerDatabaseBackupCfg", CustomPages.UpdateServerDatabaseBackupCfg)

	common.AdminEngine.HTML("GET", "/admin/userMgr/PlayerDetail", CustomPages.GetPlayerDetail)
	common.AdminEngine.HTML("GET", "/admin/userMgr/SubmitOperation", CustomPages.SubmitOperation)

	common.AdminEngine.HTML("GET", "/admin/DataDetail", CustomPages.GetTDataChart)
	common.AdminEngine.HTML("GET", "/admin/OnlinePopulationTrendChart", CustomPages.GetOnlinePopulationTrendChart)
	common.AdminEngine.Data("POST", "/admin/ChartData", CustomPages.GetChartData)

	common.AdminEngine.Data("POST", "/admin/popup/updateServerSetting", CustomPages.UpdateServerSetting)
	common.AdminEngine.Data("POST", "/admin/popup/HotfixServer", AutoPages.HotfixServer)
	common.AdminEngine.Data("POST", "/admin/popup/UpdateServerSettingPro", AutoPages.UpdateServerSettingPro)

	common.AdminEngine.Data("GET", "/admin/SendComplaint", CustomPages.SendComplaint, true)

	common.AdminEngine.HTML("GET", "/admin/EmailSettingPage", CustomPages.GetEmailSettingPage)
	common.AdminEngine.HTML("GET", "/admin/UpdateEmailSetting", CustomPages.UpdateEmailSetting)

	common.AdminEngine.HTML("GET", "/admin/BattleGrouping", CustomPages.GetBattleGroupingPage)
	common.AdminEngine.HTML("GET", "/admin/SubmitBattleInfo", CustomPages.SubmitBattleInfo)

	common.AdminEngine.HTML("GET", "/admin/2V2PkConfig", CustomPages.Get2V2PkConfigPage)
	common.AdminEngine.HTML("GET", "/admin/Submit2V2PkInfo", CustomPages.Submit2V2PkInfo)

	common.AdminEngine.Data("POST", "/admin/popup/OperationResFilesPro", AutoPages.OperationResFilesPro)

	//common.AdminEngine.HTML("GET", "/admin/TransformCharacter", CustomPages.GetTransformCharacterPage)
	//common.AdminEngine.HTML("GET", "/admin/SubmitTransfer", CustomPages.GetSubmitTransferPage)
	//
	common.AdminEngine.HTML("GET", "/admin/TransformServer", CustomPages.GetTransformServerPage)
	common.AdminEngine.HTML("GET", "/admin/SubmitTransferServer", CustomPages.GetSubmitTransferServerPage)
	//
	common.AdminEngine.HTML("GET", "/admin/DeletePlayer", CustomPages.GetDeletePlayerPage)
	common.AdminEngine.HTML("GET", "/admin/SubmitDeletePlayer", CustomPages.GetSubmitDeletePlayerPage)

	common.AdminEngine.HTML("GET", "/admin/GetSetRebornPage", CustomPages.GetSetRebornPage)
	common.AdminEngine.HTML("GET", "/admin/SubmitRebornNum", CustomPages.SubmitRebornNum)
	common.AdminEngine.HTML("GET", "/admin/GetAddItemPage", CustomPages.GetAddItemPage)
	common.AdminEngine.HTML("GET", "/admin/SubmitAddItem", CustomPages.SubmitAddItem)

	common.AdminEngine.HTML("GET", "/admin/DailySignIn", CustomPages.GetDailySignInPage)
	common.AdminEngine.HTML("GET", "/admin/SubmitDailySignIn", CustomPages.SubmitDailySignIn)

	common.AdminEngine.HTML("GET", "/admin/MonsterAttr", CustomPages.GetMonsterAttrConfigPage)
	common.AdminEngine.HTML("GET", "/admin/SubmitMonsterAttrConfig", CustomPages.SubmitMonsterAttrConfig)

	common.AdminEngine.HTML("GET", "/admin/HotfixDataTables", CustomPages.GetHotfixDataTablesPage)
	common.AdminEngine.Data("POST", "/admin/hotfix/start", CustomPages.StartHotfixPipeline)
	common.AdminEngine.Data("GET", "/admin/hotfix/progress", CustomPages.GetHotfixProgress)
	common.AdminEngine.Data("POST", "/admin/hotfix/stop", CustomPages.StopHotfixPipeline)

	common.AdminEngine.Data("POST", "/admin/popup/HotfixDBTables", AutoPages.HotfixDBTables)

}
