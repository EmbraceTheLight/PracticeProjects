package adaptor

import (
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
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
