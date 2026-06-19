package config

import (
	"fmt"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/nacos-group/nacos-sdk-go/v2/clients"
	"github.com/nacos-group/nacos-sdk-go/v2/clients/config_client"
	"github.com/nacos-group/nacos-sdk-go/v2/clients/naming_client"
	"github.com/nacos-group/nacos-sdk-go/v2/common/constant"
	"github.com/nacos-group/nacos-sdk-go/v2/vo"
	"go.uber.org/zap"

	"mall/utils/logger"
)

const nacosRetryInterval = time.Second

var (
	NacosNamingClient naming_client.INamingClient
	NacosConfigClient config_client.IConfigClient

	nacosInitOnce sync.Once
	nacosInitErr  error
)

func InitNacosClient() error {
	nacosInitOnce.Do(func() {
		nacosInitErr = initNacosClient()
	})
	return nacosInitErr
}

func NacosRegister() {
	if !IsNacosEnabled() {
		logger.Warn("skip nacos register: nacos env is not configured")
		return
	}
	if err := InitNacosClient(); err != nil {
		logger.Error("init nacos client error", zap.Error(err))
		return
	}

	servicePort, err := strconv.Atoi(ServicePort)
	if err != nil {
		logger.Error("invalid SERVICE_PORT", zap.String("service_port", ServicePort), zap.Error(err))
		return
	}

	registerParam := vo.RegisterInstanceParam{
		Ip:          ServiceIP,
		Port:        uint64(servicePort),
		ServiceName: ServerName,
		GroupName:   NacosGroup,
		Weight:      10,
		Enable:      true,
		Healthy:     true,
		Ephemeral:   true,
		Metadata:    map[string]string{},
	}
	if strings.TrimSpace(registerParam.Ip) == "" {
		registerParam.Ip = "127.0.0.1"
	}

	var lastErr error
	for i := 0; i < NacosRetryTimes; i++ {
		success, err := NacosNamingClient.RegisterInstance(registerParam)
		if err == nil && success {
			logger.Info("mall service registered to nacos",
				zap.String("service_name", registerParam.ServiceName),
				zap.String("group", registerParam.GroupName),
				zap.String("ip", registerParam.Ip),
				zap.Uint64("port", registerParam.Port),
			)
			return
		}
		if err != nil {
			lastErr = err
			logger.Warn("register nacos instance failed, will retry",
				zap.Int("attempt", i+1),
				zap.Int("max_attempts", NacosRetryTimes),
				zap.Error(err),
			)
		} else {
			lastErr = fmt.Errorf("nacos register returned success=false")
			logger.Warn("register nacos instance returned false, will retry",
				zap.Int("attempt", i+1),
				zap.Int("max_attempts", NacosRetryTimes),
				zap.String("service_name", registerParam.ServiceName),
				zap.String("group", registerParam.GroupName),
				zap.String("ip", registerParam.Ip),
				zap.Uint64("port", registerParam.Port),
				zap.String("namespace_id", NamespaceID),
			)
		}
		time.Sleep(nacosRetryInterval)
	}
	logger.Error("register nacos instance failed", zap.Error(lastErr))
}

func initNacosClient() error {
	if !IsNacosEnabled() {
		return fmt.Errorf("nacos env is not configured")
	}

	port, err := strconv.Atoi(NacosPort)
	if err != nil {
		return fmt.Errorf("invalid NACOS_PORT %q: %w", NacosPort, err)
	}
	grpcPort := 0
	if strings.TrimSpace(NacosGrpcPort) != "" {
		grpcPort, err = strconv.Atoi(NacosGrpcPort)
		if err != nil {
			return fmt.Errorf("invalid NACOS_GRPC_PORT %q: %w", NacosGrpcPort, err)
		}
	}

	logger.Info("initializing nacos client",
		zap.String("host", NacosIP),
		zap.Int("http_port", port),
		zap.Int("grpc_port", grpcPort),
		zap.String("namespace_id", NamespaceID),
		zap.String("group", NacosGroup),
		zap.String("data_id", DataID),
		zap.Bool("has_username", strings.TrimSpace(NacosUsername) != ""),
	)
	clientConfig := constant.ClientConfig{
		NamespaceId:         NamespaceID,
		Username:            NacosUsername,
		Password:            NacosPassword,
		TimeoutMs:           5000,
		NotLoadCacheAtStart: true,
		LogDir:              "log/nacos",
		CacheDir:            "log/nacos",
		LogLevel:            "debug",
	}
	serverConfigs := []constant.ServerConfig{
		{
			IpAddr:      NacosIP,
			ContextPath: "/nacos",
			Port:        uint64(port),
			GrpcPort:    uint64(grpcPort),
			Scheme:      "http",
		},
	}

	namingClient, err := clients.CreateNamingClient(map[string]interface{}{
		"serverConfigs": serverConfigs,
		"clientConfig":  clientConfig,
	})
	if err != nil {
		return err
	}

	configClient, err := clients.NewConfigClient(vo.NacosClientParam{
		ClientConfig:  &clientConfig,
		ServerConfigs: serverConfigs,
	})
	if err != nil {
		return err
	}

	NacosNamingClient = namingClient
	NacosConfigClient = configClient
	return nil
}

func IsNacosEnabled() bool {
	return strings.TrimSpace(NacosIP) != "" &&
		strings.TrimSpace(NacosPort) != "" &&
		strings.TrimSpace(NacosGroup) != "" &&
		strings.TrimSpace(NamespaceID) != "" &&
		strings.TrimSpace(DataID) != ""
}
