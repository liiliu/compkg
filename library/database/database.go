package database

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	stdGorm "gorm.io/gorm"
	stdLogger "gorm.io/gorm/logger"
	"log"
	"os"
	"time"
	"weihu_server/library/common"
	"weihu_server/library/config"
	"weihu_server/library/logger"
)

var MasterDb *stdGorm.DB
var BackendReadDb *stdGorm.DB

func Initial() {

	var logLevel stdLogger.LogLevel
	if config.GetBool("server.debug") {
		logLevel = stdLogger.Info
	} else {
		logLevel = stdLogger.Warn
	}

	slowLogger := stdLogger.New(
		//将标准输出作为Writer
		log.New(os.Stdout, "\r\n", log.LstdFlags),
		stdLogger.Config{
			//设定慢查询时间阈值为200ms
			SlowThreshold: 200 * time.Millisecond,
			//设置日志级别，只有Warn和Info级别会输出慢查询日志
			LogLevel: logLevel,
			//设置是否彩色打印
			Colorful: true,
			//忽略RecordNotFound错误
			IgnoreRecordNotFoundError: true,
		},
	)

	logger.Info(common.LogTagDb, "初始化数据库 ...")
	{
		dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
			config.GetString("postgresMaster.username"), config.GetString("postgresMaster.password"),
			config.GetString("postgresMaster.host"), config.GetInt("postgresMaster.port"), config.GetString("postgresMaster.database"))
		db, err := stdGorm.Open(mysql.Open(dsn), &stdGorm.Config{
			Logger: slowLogger,
		})
		if err != nil {
			panic("failed to connect database")
		} else {
			sqlDB, err := db.DB()
			if err != nil {
				panic("failed to connect database")
			}
			sqlDB.SetMaxIdleConns(config.GetInt("postgresMaster.maxIdleConns"))
			sqlDB.SetMaxOpenConns(config.GetInt("postgresMaster.maxOpenConns"))

			if sqlDB.Ping() != nil {
				panic("failed to connect database")
			}
			MasterDb = db
		}
	}
	{
		dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
			config.GetString("postgresBackend.username"), config.GetString("postgresBackend.password"),
			config.GetString("postgresBackend.host"), config.GetInt("postgresBackend.port"), config.GetString("postgresBackend.database"))
		db, err := stdGorm.Open(mysql.Open(dsn), &stdGorm.Config{
			Logger: slowLogger,
		})
		if err != nil {
			panic("failed to connect database")
		} else {
			sqlDB, err := db.DB()
			if err != nil {
				panic("failed to connect database")
			}
			sqlDB.SetMaxIdleConns(config.GetInt("postgresBackend.maxIdleConns"))
			sqlDB.SetMaxOpenConns(config.GetInt("postgresBackend.maxOpenConns"))

			if sqlDB.Ping() != nil {
				panic("failed to connect database")
			}
			BackendReadDb = db
		}
	}
	logger.Info(common.LogTagDb, "数据库初始化完成")

}

func InitialPostgresql() {

	var logLevel stdLogger.LogLevel
	if config.GetBool("server.debug") {
		logLevel = stdLogger.Info
	} else {
		logLevel = stdLogger.Warn
	}

	slowLogger := stdLogger.New(
		//将标准输出作为Writer
		log.New(os.Stdout, "\r\n", log.LstdFlags),
		stdLogger.Config{
			//设定慢查询时间阈值为200ms
			SlowThreshold: 200 * time.Millisecond,
			//设置日志级别，只有Warn和Info级别会输出慢查询日志
			LogLevel: logLevel,
			//设置是否彩色打印
			Colorful: true,
			//忽略RecordNotFound错误
			IgnoreRecordNotFoundError: true,
		},
	)

	logger.Info(common.LogTagDb, "初始化数据库 ...")

	{
		dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Shanghai",
			config.GetString("postgresMaster.host"), config.GetString("postgresMaster.username"), config.GetString("postgresMaster.password"), config.GetString("postgresMaster.database"), config.GetString("postgresMaster.port"))

		db, err := stdGorm.Open(postgres.Open(dsn), &stdGorm.Config{Logger: slowLogger})
		if err != nil {
			panic("failed to connect database")
		} else {
			sqlDB, err := db.DB()
			if err != nil {
				panic("failed to connect database")
			}
			sqlDB.SetMaxIdleConns(config.GetInt("postgresMaster.maxIdleConns"))
			sqlDB.SetMaxOpenConns(config.GetInt("postgresMaster.maxOpenConns"))

			if sqlDB.Ping() != nil {
				panic("failed to connect database")
			}
			MasterDb = db
		}
	}
	{
		dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Shanghai",
			config.GetString("postgresBackend.host"), config.GetString("postgresBackend.username"), config.GetString("postgresBackend.password"), config.GetString("postgresBackend.database"), config.GetString("postgresBackend.port"))

		db, err := stdGorm.Open(postgres.Open(dsn), &stdGorm.Config{Logger: slowLogger})
		if err != nil {
			panic("failed to connect database")
		} else {
			sqlDB, err := db.DB()
			if err != nil {
				panic("failed to connect database")
			}
			sqlDB.SetMaxIdleConns(config.GetInt("postgresBackend.maxIdleConns"))
			sqlDB.SetMaxOpenConns(config.GetInt("postgresBackend.maxOpenConns"))

			if sqlDB.Ping() != nil {
				panic("failed to connect database")
			}
			BackendReadDb = db
		}
	}
	logger.Info(common.LogTagDb, "数据库初始化完成")

}
