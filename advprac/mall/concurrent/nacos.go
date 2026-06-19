package concurrent

import (
	"time"

	"github.com/nacos-group/nacos-sdk-go/v2/vo"
	"go.uber.org/zap"
	"gopkg.in/yaml.v3"

	"mall/config"
	"mall/utils/logger"
)

func Go() {
	go ListenConfig()
}

func ListenConfig() {
	if !config.IsNacosEnabled() {
		logger.Warn("skip nacos config listener: nacos env is not configured")
		return
	}
	if err := config.InitNacosClient(); err != nil {
		logger.Error("init nacos client error", zap.Error(err))
		return
	}

	configParam := vo.ConfigParam{
		DataId: config.DataID,
		Group:  config.NacosGroup,
		OnChange: func(namespace, group, dataId, data string) {
			var mallConfig config.Config
			if err := yaml.Unmarshal([]byte(data), &mallConfig); err != nil {
				logger.Error("unmarshal nacos config error", zap.Error(err))
				return
			}
			config.GlobalConfig = &mallConfig
			logger.Info("nacos config changed",
				zap.String("namespace", namespace),
				zap.String("group", group),
				zap.String("data_id", dataId),
			)
		},
	}

	var lastErr error
	for i := 0; i < config.NacosRetryTimes; i++ {
		err := config.NacosConfigClient.ListenConfig(configParam)
		if err == nil {
			logger.Info("listen nacos config success")
			return
		}
		lastErr = err
		logger.Warn("listen nacos config failed, will retry",
			zap.Int("attempt", i+1),
			zap.Int("max_attempts", config.NacosRetryTimes),
			zap.Error(err),
		)
		time.Sleep(time.Second)
	}
	logger.Error("listen nacos config failed", zap.Error(lastErr))
}
