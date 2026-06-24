package pageMgr

import (
	"admin/common"
	"admin/common/def"
	"admin/fusion"
	"admin/fusion/timer"
	"encoding/json"
	"github.com/GoAdminGroup/go-admin/context"
	"net/url"
	"strconv"
	"strings"
	"time"
)

type StatisticsTimer struct {
	timer.WheelTriggerOwner
}

func NewMailTimer() *StatisticsTimer {
	recorder := new(StatisticsTimer)
	recorder.InitVT(recorder)
	return recorder
}

func (recorder *StatisticsTimer) Impl_GetWheelTimerMgr() *timer.WheelTimerMgr {
	return fusion.MainService.TimerMgr
}

func SaveMailToDBAndNotifyToCenter(mailInfo *def.TMails, playerIDs []int, serverID int) error {
	var successUserID []int
	if playerIDs == nil {
		db := fusion.GetBaseGormDB("db_global")
		rst := db.Create(&mailInfo)
		if rst.Error != nil {
			return rst.Error
		}
		vals := url.Values{}
		vals.Add("mailID", strconv.Itoa(int(mailInfo.MailID)))
		ctx := &context.Context{}
		fusion.CallToCenter(&vals, common.MyCfg.Api["pushServerMail"], "GET", ctx)
		return nil
	} else {
		for _, v := range playerIDs {
			var userMail = def.InstMail{}
			userMail.MailType = mailInfo.MailType
			userMail.MailFlags = 0
			userMail.MailSender = 0
			userMail.MailReceiver = uint32(v)
			temp, _ := time.ParseInLocation("2006-01-02 15:04:05", mailInfo.MailDeliverTime, time.Local)
			userMail.MailDeliverTime = int(temp.Unix())
			temp, _ = time.ParseInLocation("2006-01-02 15:04:05", mailInfo.MailExpireTime, time.Local)
			userMail.MailExpireTime = int(temp.Unix())
			userMail.MailSubject = mailInfo.MailSubject
			userMail.MailBody = mailInfo.MailBody
			userMail.MailCheques = mailInfo.MailCheques
			userMail.MailItems = mailInfo.MailItems
			userMail.IsGetAttachment = 0
			userMail.IsViewDetail = 0
			db := fusion.GetServerGormDB("db_char", int64(serverID))
			if db == nil {
				return nil
			}
			rst := db.Create(&userMail)
			if rst.Error != nil {
				return rst.Error
			}
			successUserID = append(successUserID, v)
		}
		jsonData, _ := json.Marshal(playerIDs)
		vals := url.Values{}
		vals.Add("playerIds", string(jsonData))
		vals.Add("gsId", strconv.Itoa(serverID))
		ctx := &context.Context{}
		err, _ := fusion.CallToCenter(&vals, common.MyCfg.Api["noticeNewMail"], "GET", ctx)
		return err
	}
}

func SaveTimingMail2DB(Transit *def.MailTransit) error {
	now := time.Now()
	deliver, _ := time.ParseInLocation("2006-01-02 15:04:05", Transit.DeliverTime, time.Local)
	go_manager := fusion.GetBaseGormDB("go_manager")
	go_manager.Create(&Transit)
	mailTimer := NewMailTimer()
	interval := deliver.Unix() - time.Now().Unix()
	if interval <= 1 {
		interval = 1
	}

	if now.After(deliver) {
		SendTimingMail(Transit)
	} else {
		mailTimer.CreateTriggerX4tp(uint64(interval), uint64(deliver.Unix()),
			func() {
				SendTimingMail(Transit)
			}, 1)
	}
	return nil
}

func SendTimingMail(transit *def.MailTransit) {
	now := time.Now()
	Expire, _ := time.ParseInLocation("2006-01-02 15:04:05 ", transit.ExpireTime, time.Local)
	var isSend = 0
	defer func() {
		go_manager := fusion.GetBaseGormDB("go_manager")
		if transit.Id != 0 {
			go_manager.Table(def.MailTransit{}.TableName()).Where("Id", transit.Id).Update("isSend", isSend)
		}
	}()

	if Expire.Before(now) {
		isSend = 2
		return
	}
	temp := fusion.String2Bytes(transit.Record)
	timingTMail := def.TimingTMail{}
	err := json.Unmarshal(temp, &timingTMail)
	if err != nil {
		isSend = 2
		return
	}
	err = SaveMailToDBAndNotifyToCenter(&timingTMail.Mail, timingTMail.PlayerIDs, timingTMail.GsId)
	if err != nil {
		isSend = 2
		return
	}
	isSend = 1
}

func LoadTimingMailFromDB() {
	mailTransits := []def.MailTransit{}
	go_manager := fusion.GetBaseGormDB("go_manager")
	go_manager.Where("isSend = ?", 0).Find(&mailTransits)
	mailTimer := NewMailTimer()
	now := time.Now()
	for _, v := range mailTransits {
		expire, _ := time.ParseInLocation("2006-01-02T15:04:05Z", v.ExpireTime, time.Local)
		deliver, _ := time.ParseInLocation("2006-01-02T15:04:05Z", v.DeliverTime, time.Local)
		v.ExpireTime = strings.Replace(v.ExpireTime, "T", " ", -1)
		v.ExpireTime = strings.Replace(v.ExpireTime, "Z", " ", -1)
		v.DeliverTime = strings.Replace(v.DeliverTime, "T", " ", -1)
		v.DeliverTime = strings.Replace(v.DeliverTime, "Z", " ", -1)
		interval := deliver.Unix() - time.Now().Unix()
		if interval <= 1 {
			interval = 1
		}
		if now.Before(expire) {
			if now.After(deliver) {
				SendTimingMail(&v)
			} else {
				mailTimer.CreateTriggerX4tp(uint64(interval), uint64(deliver.Unix()),
					func() {
						SendTimingMail(&v)
					}, 1)
			}
		}
	}
}
