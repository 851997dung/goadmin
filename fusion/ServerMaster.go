package fusion

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"sync"
	"time"

	"admin/fusion/base"

	_ "github.com/go-sql-driver/mysql"
	"github.com/panjf2000/ants/v2"
)

var IsStopService bool
var SigStopService chan bool
var AllWebServers sync.Map // map[*http.Server]interface{}
var AllClosers sync.Map    // map[io.Closer]interface{}

var GoPool *ants.Pool
var GoPool10 *ants.Pool
var MainService *ServiceBase

//var allDBs []*sql.DB
//var allRedis []*redis.Client

type MysqlConfig struct {
	Host     string `json:"HOST"`
	Port     uint16 `json:"PORT"`
	User     string `json:"USER"`
	Password string `json:"PASSWORD"`
	Database string `json:"DATABASE"`
}

type RedisConfig struct {
	Addr     string `json:"ADDR"`
	Password string `json:"PASSWORD"`
	DB       int    `json:"DB"`
}

func init() {
	SigStopService = make(chan bool)
}

func MainRunAndWaitQuit() {
	var appName = filepath.Base(os.Args[0])
	log.Printf("%s start successfully.\n", appName)

	for !IsStopService {
		var c string
		if _, err := fmt.Scan(&c); err == nil {
			switch c {
			case "q", "Q":
				IsStopService = true
				close(SigStopService)
			}
		}
		time.Sleep(time.Millisecond * 10)
	}

	AllWebServers.Range(func(srv, _ interface{}) bool {
		srv.(*http.Server).Shutdown(nil)
		return true
	})
	AllClosers.Range(func(cls, _ interface{}) bool {
		cls.(io.Closer).Close()
		return true
	})

	if GoPool != nil {
		GoPool.Release()
		for GoPool.Running() != 0 {
			time.Sleep(100 * time.Millisecond)
		}
	}
	if GoPool10 != nil {
		GoPool10.Release()
		for GoPool10.Running() != 0 {
			time.Sleep(100 * time.Millisecond)
		}
	}

	base.WaitHandler()
	//for _, db := range allDBs {
	//	db.Close()
	//}
	//for _, rds := range allRedis {
	//	rds.Close()
	//}

	log.Printf("%s shutdown gracefully!\n", appName)
}

//func LoadUserConfig(filename string, cfg interface{}) {
//	cfgBytes, err := ioutil.ReadFile(filename)
//	if err != nil {
//		log.Fatalf("read user config error, %s\n", err)
//	}
//	err = json.Unmarshal(cfgBytes, cfg)
//	if err != nil {
//		log.Fatalf("parse user config error, %s\n", err)
//	}
//	log.Printf("load user config successfully.\n")
//}

//func OpenMysql(cfg *MysqlConfig, maxOpen, maxIdle int, isGlobal bool) *sql.DB {
//	dataSourceName := fmt.Sprintf(
//		"%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=true",
//		cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.Database)
//	db, err := sql.Open("mysql", dataSourceName)
//	if err != nil {
//		log.Fatalf("open mysql `%s` is failed, %s\n", dataSourceName, err)
//	}
//	if err = db.Ping(); err != nil {
//		log.Fatalf("ping mysql `%s` is failed, %s\n", dataSourceName, err)
//	}
//	log.Printf("connect mysql database `%s` successfully.\n", dataSourceName)
//	db.SetMaxOpenConns(maxOpen)
//	db.SetMaxIdleConns(maxIdle)
//	//if isGlobal {
//	//	allDBs = append(allDBs, db)
//	//}
//	return db
//}
//
//func OpenRedis(cfg *RedisConfig, isGlobal bool) *redis.Client {
//	rds := redis.NewClient(&redis.Options{
//		Addr:     cfg.Addr,
//		Password: cfg.Password,
//		DB:       cfg.DB,
//	})
//	if _, err := rds.Ping(context.TODO()).Result(); err != nil {
//		log.Fatalf("ping redis `%s`,`%s`,`%d` is failed, %s\n",
//			cfg.Addr, cfg.Password, cfg.DB, err)
//	}
//	log.Printf("open redis `%s`,`%s`,`%d` successfully.\n",
//		cfg.Addr, cfg.Password, cfg.DB)
//	//if isGlobal {
//	//	allRedis = append(allRedis, rds)
//	//}
//	return rds
//}

func InitGoPool(size int) {
	var err error
	GoPool, err = ants.NewPool(size)
	if err != nil {
		log.Fatalf("new go poll `%d` is failed, %s\n", size, err)
	}
	log.Printf("init go poll `%d` successfully.\n", size)
}

func InitGoPool10(size int) {
	var err error
	GoPool10, err = ants.NewPool(size)
	if err != nil {
		log.Fatalf("new go poll `%d` is failed, %s\n", size, err)
	}
	log.Printf("init go poll `%d` successfully.\n", size)
}

func StartService(s *ServiceBase, name string) {
	if name == "main" {
		MainService = s
	}
	s.Init(name)
	s.Start()
	InitGoPool(1000)
	InitGoPool10(7)
}
