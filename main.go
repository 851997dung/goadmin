package main

import (
	"admin/common"
	"admin/common/menuTranslate"
	"admin/fusion"
	"admin/mgr"
	"admin/mgr/pageMgr"
	"admin/pages"
	_ "github.com/GoAdminGroup/go-admin/adapter/gin" // web framework adapter
	"github.com/GoAdminGroup/go-admin/engine"
	"github.com/GoAdminGroup/go-admin/modules/config"
	_ "github.com/GoAdminGroup/go-admin/modules/db/drivers/mysql" // sql driver
	"github.com/GoAdminGroup/go-admin/template"
	"github.com/GoAdminGroup/go-admin/template/chartjs"
	_ "github.com/GoAdminGroup/themes/adminlte" // ui theme
	_ "github.com/GoAdminGroup/themes/sword"    // ui theme
	"github.com/gin-gonic/gin"
	"github.com/leonelquinteros/gotext"
	"io/ioutil"
	"log"
	_ "net/http/pprof"
	"os"
	"os/signal"
)

var Inst mainService

type mainService struct {
	fusion.ServiceBase
}

func main() {
	startServer()
}

func startServer() {
	gin.SetMode(gin.ReleaseMode)
	gin.DefaultWriter = ioutil.Discard
	common.GinEngine = gin.Default()
	template.AddComp(chartjs.NewChart())
	common.AdminEngine = engine.Default()
	if err := common.AdminEngine.AddConfigFromYAML("./config.yml").
		AddGenerators(pages.Generators).
		Use(common.GinEngine); err != nil {
		panic(err)
	}
	common.GinEngine.Static("/uploads", "./uploads")

	common.Cfg_yml = config.ReadFromYaml("./config.yml")
	fusion.InitBaseDataBase()
	fusion.StartService(&Inst.ServiceBase, "main")
	fusion.InitConfig(&common.MyCfg)
	pages.RegisterPath()
	pageMgr.LoadTimingMailFromDB()
	mgr.StartAllTimer()
	menuTranslate.LoadTranslate()
	gotext.Configure("path/to/locales", common.MyCfg.Language, common.MyCfg.Language)
	fusion.InitPprofMonitor()
	_ = common.GinEngine.Run(":8081")

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	fusion.MainRunAndWaitQuit()
	<-quit
	log.Print("closing database connection")
	common.AdminEngine.MysqlConnection().Close()
}
