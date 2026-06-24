package pageMgr

import (
	"admin/common/def"
	"admin/fusion"
	"strconv"
	"strings"
	"time"

	"github.com/GoAdminGroup/go-admin/context"
	"github.com/GoAdminGroup/go-admin/modules/db"
	"github.com/GoAdminGroup/go-admin/modules/logger"
	"github.com/GoAdminGroup/go-admin/plugins/admin/modules/parameter"
	"github.com/GoAdminGroup/go-admin/plugins/admin/modules/table"
	"github.com/GoAdminGroup/go-admin/template/types"
	"github.com/fatih/structs"
	"github.com/leonelquinteros/gotext"
)

func GetBagDetailInfo(ctx *context.Context, items *string) table.Table {
	BagDetailInfo := table.NewDefaultTable(table.Config{
		Driver:     "mysql",
		Connection: "",
		CanAdd:     true,
		Editable:   true,
		Deletable:  true,
		Exportable: true,
		PrimaryKey: table.PrimaryKey{
			Type: db.Int,
			Name: "ItemGuid",
		},
	})
	itemNameString := fusion.GetItemNameStrings()
	panelList := BagDetailInfo.GetInfo()
	panelList.AddField(gotext.Get("道具格子"), "Slot", db.Int).
		FieldHide()
	//FieldSortable()
	panelList.AddField(gotext.Get("道具GUID"), "ItemGuid", db.Int).
		FieldHide()
	//FieldHide()
	panelList.AddField(gotext.Get("道具ID"), "ItemTypeID", db.Int).
		FieldSortable().
		FieldDisplay(func(model types.FieldModel) interface{} {
			temp, _ := strconv.Atoi(model.Value)
			if itemNameString[temp] != "" {
				return itemNameString[temp] + "(" + model.Value + ")"
			}
			return model.Value
		})
	panelList.AddField(gotext.Get("道具数量"), "ItemCount", db.Int)
	//FieldSortable()
	panelList.AddField(gotext.Get("是/否绑定"), "ItemOwner", db.Int).
		FieldDisplay(func(model types.FieldModel) interface{} {
			temp, _ := strconv.Atoi(model.Value)
			if temp == 1 {
				return gotext.Get("是")
			}
			return gotext.Get("否")
		})
	panelList.AddField("道具状态", "ItemStatus", db.Int).
		FieldDisplay(func(model types.FieldModel) interface{} {
			temp, _ := strconv.Atoi(model.Value)
			myOps := def.GetOpStrings()
			return myOps[temp]
		})

	panelList.HideNewButton().HideEditButton().HideDeleteButton().
		HideDetailButton().HideRowSelector().HideFilterArea().HideFilterButton().HideCheckBoxColumn()

	panelList.SetGetDataFn(
		func(param parameter.Parameters) (data []map[string]interface{}, size int) {
			arr := strings.Split(*items, ";")
			var list = fusion.ReadItemListFromByte(arr, true)
			temp := []map[string]interface{}{}
			for _, v := range list {
				m3 := structs.Map(&v)
				temp = append(temp, m3)
			}
			//max := math.Min(float64(len(temp)), float64((param.PageInt)*param.PageSizeInt))
			//return temp[(param.PageInt-1)*param.PageSizeInt : int(max)], len(temp)
			return temp, len(temp)
		})
	return BagDetailInfo
}

func GetPlayerInfo(ctx *context.Context) map[string]interface{} {
	ServerID := ctx.Request.FormValue("IpcServerID")
	AccountID := ctx.Request.FormValue("IpcAcctID")
	PlayerID := ctx.Request.FormValue("IpcInstID")
	totalStart := time.Now()
	defer func() {
		logger.Infof("[GetPlayerInfo] total elapsed=%v serverID=%s acctID=%s playerID=%s",
			time.Since(totalStart), ServerID, AccountID, PlayerID)
	}()

	ServerIDint, _ := strconv.Atoi(ServerID)

	t0 := time.Now()
	db_char := fusion.GetServerGormDB("db_char", int64(ServerIDint))
	logger.Infof("[GetPlayerInfo] GetServerGormDB elapsed=%v serverID=%s", time.Since(t0), ServerID)
	if db_char == nil {
		return nil
	}
	db_global := fusion.GetBaseGormDB("db_global")

	type accountCut struct {
		UserName    string `gorm:"column:username" db:"username"`
		Id          int    `gorm:"column:Id" db:"Id"`
		CreateTime  int    `gorm:"column:createTime" db:"createTime"`
		LastLoginIP string `gorm:"column:lastLoginIP" db:"lastLoginIP"`
	}
	type orderCut1 struct {
		MinPayTime time.Time `gorm:"column:minPayTime" db:"minPayTime"`
		MaxPayTime time.Time `gorm:"column:maxPayTime" db:"maxPayTime"`
		PlayerId   int       `gorm:"column:playerId" db:"playerId"`
	}
	type tempCut struct {
		PlayerId   int       `gorm:"column:playerId" db:"playerId"`
		CreateTime time.Time `gorm:"column:createTime" db:"createTime"`
	}

	fields := fusion.GetFields2Struct(def.PlayerInfo{})
	conditions := []string{}
	args := []interface{}{}
	if ServerID != "" {
		conditions = append(conditions, "ipcServerID = ?")
		args = append(args, ServerID)
	}
	if AccountID != "" {
		conditions = append(conditions, "ipcAcctID = ?")
		args = append(args, AccountID)
	}
	if PlayerID != "" {
		conditions = append(conditions, "ipcInstID = ?")
		args = append(args, PlayerID)
	}
	whereClause := ""
	if len(conditions) > 0 {
		whereClause = " where " + strings.Join(conditions, " and ") + " and ipcDeleteTime = 0"
	}
	sql1 := "select " + fields + " from inst_player_char" + whereClause

	pid, _ := strconv.Atoi(PlayerID)
	if pid > 0 {
		ch1 := make(chan def.PlayerInfo, 1)
		ch2 := make(chan accountCut, 1)
		ch3 := make(chan orderCut1, 1)

		go func() {
			t := time.Now()
			info := def.PlayerInfo{}
			db_char.Raw(sql1, args...).Scan(&info)
			logger.Infof("[GetPlayerInfo] Q1-char elapsed=%v pid=%d", time.Since(t), pid)
			ch1 <- info
		}()
		go func() {
			t := time.Now()
			info := accountCut{}
			db_global.Raw("select username,Id,createTime, lastLoginIP from t_accounts "+
				"where id = (select accountId from t_account_characters where characterId = ?)\n",
				pid).Scan(&info)
			logger.Infof("[GetPlayerInfo] Q2-account elapsed=%v pid=%d", time.Since(t), pid)
			ch2 <- info
		}()
		go func() {
			t := time.Now()
			info := orderCut1{}
			tUnion := time.Now()
			sql2, _ := fusion.CreateUnionTableSqlByDays("", "", "t_web_orders_", db_global,
				tempCut{}, "", true)
			logger.Infof("[GetPlayerInfo] Q3-unionBuilder elapsed=%v pid=%d", time.Since(tUnion), pid)
			sql2 = "SELECT\n" +
				"playerId,\n" +
				"MIN(a.createTime)AS minPayTime,\n" +
				"Max(a.createTime)AS maxPayTime\n" +
				"FROM\n" +
				"(" + sql2 + ")AS a\n" +
				"WHERE a.playerId = ? \n" +
				"GROUP BY playerId "
			db_global.Raw(sql2, pid).Scan(&info)
			logger.Infof("[GetPlayerInfo] Q3-order elapsed=%v pid=%d", time.Since(t), pid)
			ch3 <- info
		}()

		playerInfo := <-ch1
		cut1 := <-ch2
		order1 := <-ch3

		m1 := structs.Map(playerInfo)
		m2 := structs.Map(&cut1)
		m3 := structs.Map(&order1)
		m4 := fusion.MergeMap(m1, m2)
		return fusion.MergeMap(m4, m3)
	}

	t1 := time.Now()
	playerInfo := def.PlayerInfo{}
	db_char.Raw(sql1, args...).Scan(&playerInfo)
	logger.Infof("[GetPlayerInfo] Q1-char(seq) elapsed=%v pid=%d", time.Since(t1), playerInfo.IpcInstID)
	m1 := structs.Map(playerInfo)
	pid = playerInfo.IpcInstID

	t2 := time.Now()
	cut1 := accountCut{}
	db_global.Raw("select username,Id,createTime, lastLoginIP from t_accounts "+
		"where id = (select accountId from t_account_characters where characterId = ?)\n",
		pid).Scan(&cut1)
	logger.Infof("[GetPlayerInfo] Q2-account(seq) elapsed=%v pid=%d", time.Since(t2), pid)
	m2 := structs.Map(&cut1)

	t3 := time.Now()
	order1 := orderCut1{}
	sql2, _ := fusion.CreateUnionTableSqlByDays("", "", "t_web_orders_", db_global,
		tempCut{}, "", true)
	sql2 = "SELECT\n" +
		"playerId,\n" +
		"MIN(a.createTime)AS minPayTime,\n" +
		"Max(a.createTime)AS maxPayTime\n" +
		"FROM\n" +
		"(" + sql2 + ")AS a\n" +
		"WHERE a.playerId = ? \n" +
		"GROUP BY playerId "
	db_global.Raw(sql2, pid).Scan(&order1)
	logger.Infof("[GetPlayerInfo] Q3-order(seq) elapsed=%v pid=%d", time.Since(t3), pid)
	m3 := structs.Map(&order1)
	m4 := fusion.MergeMap(m1, m2)
	return fusion.MergeMap(m4, m3)
}
