package hotfix

import (
	"admin/common"
	"admin/fusion/myConfig"
	"net"
)

type Config struct {
	ProjectDir    string
	ServerAddr    string
	BuildGoModule string
	CenterHost    string
	CenterPort    string
	DeployPort    string
}

func NewConfig(cfg myConfig.MyConfig) Config {
	host, port, _ := net.SplitHostPort(cfg.Center.Host)

	c := Config{
		ProjectDir:    cfg.Hotfix.ProjectDir,
		ServerAddr:    cfg.Hotfix.ServerAddr,
		BuildGoModule: cfg.Hotfix.BuildGoModule,
		CenterHost:    host,
		CenterPort:    cfg.Center.Port,
		DeployPort:    "10005",
	}
	if port != "" {
		c.CenterPort = port
	}
	if c.BuildGoModule == "" {
		c.BuildGoModule = "mytest.com/runXML"
	}
	if c.DeployPort == "" {
		c.DeployPort = "10005"
	}
	return c
}

func GetConfig() Config {
	return NewConfig(common.MyCfg)
}
