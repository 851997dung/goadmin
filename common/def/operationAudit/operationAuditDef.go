package operationAudit

import (
	"admin/common/def/MyTime"
)

type Operationaudit struct {
	Id              int           `gorm:"primaryKey;column:id" db:"id"`
	OperatorId      int           `gorm:"column:operatorId" db:"operatorId"`
	OperationType   int           `gorm:"column:operationType" db:"operationType"`
	OperationalData string        `gorm:"column:operationalData" db:"operationalData"`
	AuditStatus     int           `gorm:"column:auditStatus" db:"auditStatus"`
	AuditorId       int           `gorm:"column:auditorId" db:"auditorId"`
	SubmissionTime  MyTime.MyTime `gorm:"column:submissionTime" db:"submissionTime"`
	AuditTime       MyTime.MyTime `gorm:"column:auditTime" db:"auditTime"`
	FinishStatus    int           `gorm:"column:finishStatus" db:"finishStatus"`
	FinishArgs      string        `gorm:"column:finishArgs" db:"finishArgs"`
}

func (Operationaudit) TableName() string {
	return "operationaudit"
}

const (
	//SimRecharge = iota
	//ActivityCompensation
	SendMail = iota
)

const (
	AuditStatus_Wait = iota
	AuditStatus_Pass
	AuditStatus_Refuse
)
