package db

import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
	"time"
	"user_growth/conf"
	"xorm.io/xorm"
	xlog "xorm.io/xorm/log"
)

var (
	DbEngine *xorm.Engine
	DSN      = fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s",
		conf.GlobalConfig.Db.Username,
		conf.GlobalConfig.Db.Password,
		conf.GlobalConfig.Db.Host,
		conf.GlobalConfig.Db.Port,
		conf.GlobalConfig.Db.Database,
		conf.GlobalConfig.Db.Charset)
	Db *gorm.DB
)

func InitDb() {
	if DbEngine != nil {
		return
	}
	if engine, err := xorm.NewEngine(conf.GlobalConfig.Db.Engine, DSN); err != nil {
		log.Fatalf("dbhelper.initDb(%s) error%s\n", DSN, err.Error())
		return
	} else {
		DbEngine = engine
	}
	if err := DbEngine.Ping(); err != nil {
		log.Fatalf("dbhelper.initDb(%s) ping error=%s\n", DSN, err.Error())
		return
	}
	logger := xlog.NewSimpleLogger(log.Writer())
	logger.ShowSQL(conf.GlobalConfig.Db.ShowSql)
	DbEngine.SetLogger(logger)
	DbEngine.SetMaxIdleConns(conf.GlobalConfig.Db.MaxIdleConns)
	DbEngine.SetMaxOpenConns(conf.GlobalConfig.Db.MaxOpenConns)
	DbEngine.SetConnMaxLifetime(time.Duration(conf.GlobalConfig.Db.ConnMaxLifetime) * time.Minute)
}

func InitGormDb() {
	if Db != nil {
		return
	}
	var err error
	Db, err = gorm.Open(mysql.Open(DSN), &gorm.Config{
		SkipDefaultTransaction: true, // 禁用默认事务
	})
	if err != nil {
		log.Fatalf("dbhelper.InitGormDb(%s) error=%s\n", DSN, err.Error())
		return
	}
	sqlDB, _ := Db.DB()
	sqlDB.SetMaxIdleConns(conf.GlobalConfig.Db.MaxIdleConns)
	sqlDB.SetMaxOpenConns(conf.GlobalConfig.Db.MaxOpenConns)
	sqlDB.SetConnMaxLifetime(time.Duration(conf.GlobalConfig.Db.ConnMaxLifetime) * time.Minute)
	log.Printf("dbhelper.InitGormDb(%s) success\n", DSN)
}

func CloseDb() {
	if DbEngine != nil {
		if err := DbEngine.Close(); err != nil {
			log.Printf("dbhelper.CloseDb error=%s\n", err.Error())
		} else {
			log.Printf("dbhelper.CloseDb success\n")
		}
		DbEngine = nil
	}
	if Db != nil {
		sqlDB, _ := Db.DB()
		if err := sqlDB.Close(); err != nil {
			log.Printf("dbhelper.CloseGormDb error=%s\n", err.Error())
		} else {
			log.Printf("dbhelper.CloseGormDb success\n")
		}
		Db = nil
	}
}

func init() {
	// 初始化数据库连接
	InitDb()
	// 初始化gorm连接
	//InitGormDb()
}
