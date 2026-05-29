package adaptor

import (
	"fmt"
	"github.com/redis/go-redis/v9"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	gormLog "gorm.io/gorm/logger"
	"mall/config"
	auth "mall/utils/jwt"
)

type IAdaptor interface {
	GetConfig() *config.Config
	GetDB() *gorm.DB
	GetRedis() *redis.Client
	GetJwtAuth() *auth.JWTAuth
}
type Adaptor struct {
	conf    *config.Config
	db      *gorm.DB
	redis   *redis.Client
	jwtAuth *auth.JWTAuth
}

// NewAdaptor 创建并返回一个新的Adaptor实例
//
// 参数:
//   - conf: 配置信息指针，用于初始化适配器
//   - db: 数据库连接指针，用于数据库操作
//   - redis: Redis客户端指针，用于缓存操作
//
// 返回值:
//   - *Adaptor: 返回一个初始化好的Adaptor结构体指针
func NewAdaptor(conf *config.Config, db *gorm.DB, redis *redis.Client, jwtAuth *auth.JWTAuth) *Adaptor {
	return &Adaptor{
		conf:    conf,    // 将传入的配置信息赋值给 Adaptor 的 conf 字段
		db:      db,      // 将传入的数据库连接赋值给 Adaptor 的 db 字段
		redis:   redis,   // 将传入的 Redis 客户端赋值给 Adaptor 的 redis 字段
		jwtAuth: jwtAuth, // 将传入的 JWTAuth 实例赋值给 Adaptor 的 jwtAuth 字段
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

func (a *Adaptor) GetConfig() *config.Config {
	return a.conf
}

func (a *Adaptor) GetDB() *gorm.DB {
	return a.db
}

func (a *Adaptor) GetRedis() *redis.Client {
	return a.redis
}

func (a *Adaptor) GetJwtAuth() *auth.JWTAuth {
	return a.jwtAuth
}
