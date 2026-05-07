package adaptor

import (
	"fmt"
	"github.com/redis/go-redis/v9"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	gormLog "gorm.io/gorm/logger"
	"mall/config"
)

type Adaptor struct {
	conf  *config.Config
	db    *gorm.DB
	redis *redis.Client
}

func NewAdaptor(conf *config.Config, db *gorm.DB, redis *redis.Client) *Adaptor {
	return &Adaptor{
		conf:  conf,
		db:    db,
		redis: redis,
	}
}

func NewMysqlData(mysqlConf *config.Mysql) (*gorm.DB, error) {
	host, port, username, password, dbname := mysqlConf.Host, mysqlConf.Port, mysqlConf.Username, mysqlConf.Password, mysqlConf.Database
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", username, password, host, port, dbname)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		DisableForeignKeyConstraintWhenMigrating: true,
	})
	if err != nil {
		return nil, err
	}
	db.Logger.LogMode(gormLog.LogLevel(mysqlConf.LogLevel))
	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}
	sqlDB.SetMaxIdleConns(int(mysqlConf.MaxIdle))
	sqlDB.SetMaxOpenConns(int(mysqlConf.MaxOpen))
	return db, nil
}

func (a *Adaptor) GetDB() *gorm.DB {
	return a.db
}

func (a *Adaptor) GetRedis() *redis.Client {
	return a.redis
}
