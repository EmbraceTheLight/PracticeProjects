package main

import (
	"context"
	"errors"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
	adaptor "mall/adaptor"
	"mall/config"
	"mall/router"
	"mall/utils/logger"
	"strconv"
	"time"
)

func main() {
	conf := config.InitConfig()
	logger.SetLevel(conf.HttpServer.LogLevel)

	dbClient, err := initMySQL(conf.Mysql)
	handleErr(err)
	logger.Debug("Mysql init success")

	rdbClient, err := initRedis(conf.Redis)
	handleErr(err)

	logger.Debug("Redis init success")

	startServer(conf, dbClient, rdbClient).Run()
}

func startServer(conf *config.Config, db *gorm.DB, redis *redis.Client) *router.App {
	adpt := adaptor.NewAdaptor(conf, db, redis)
	return router.NewApp(conf.HttpServer.HttpPort, router.NewRouter(adpt, conf, func() error {
		pingDB, err := db.DB()
		handleErr(err)
		err = pingDB.Ping()
		if err != nil {
			return errors.New("mysql connect failed")
		}
		ctx := context.TODO()
		return redis.Ping(ctx).Err()
	}))
}

func initMySQL(mysqlConf *config.Mysql) (*gorm.DB, error) {
	return adaptor.NewMysqlData(mysqlConf)
}

func initRedis(redisConf *config.Redis) (*redis.Client, error) {
	client := redis.NewClient(&redis.Options{
		Addr:         redisConf.Host + ":" + strconv.Itoa(int(redisConf.Port)),
		Password:     redisConf.Password,
		PoolSize:     int(redisConf.PoolSize),
		MinIdleConns: int(redisConf.MinIdleConns),
		MaxRetries:   int(redisConf.MaxRetries),
		DialTimeout:  redisConf.DialTimeout * time.Second,
		ReadTimeout:  redisConf.ReadTimeout * time.Second,
		WriteTimeout: redisConf.WriteTimeout * time.Second,
		PoolTimeout:  redisConf.PoolTimeout * time.Second,
	})
	return client, nil
}

func handleErr(err error) {
	if err != nil {
		panic(err)
	}
}
