package common

import (
	"admin/fusion/myConfig"
	"github.com/GoAdminGroup/go-admin/engine"
	"github.com/GoAdminGroup/go-admin/modules/config"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"gorm.io/gorm"
)

var (
	GormDBList   map[string]*gorm.DB
	ServerDBList map[int64]map[string]*gorm.DB
	MyCfg        myConfig.MyConfig
	Cfg_yml      config.Config
	GinEngine    *gin.Engine
	AdminEngine  *engine.Engine
)
