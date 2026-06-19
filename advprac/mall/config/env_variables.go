package config

import (
	"os"
	"strconv"
)

var (
	NacosIP       = os.Getenv("NACOS_IP")
	NacosPort     = os.Getenv("NACOS_PORT")
	NacosGrpcPort = os.Getenv("NACOS_GRPC_PORT")
	NacosGroup    = os.Getenv("NACOS_GROUP")
	NacosUsername = os.Getenv("NACOS_USERNAME")
	NacosPassword = os.Getenv("NACOS_PASSWORD")
	NamespaceID   = os.Getenv("NAMESPACE_ID")
)

var (
	ServiceIP   = os.Getenv("SERVICE_IP")
	ServicePort = os.Getenv("SERVICE_PORT")
	DataID      = os.Getenv("DATA_ID")
)

var (
	NacosRetryTimes = getEnvInt("NACOS_RETRY_TIMES", 30)
)

func getEnvInt(key string, defaultValue int) int {
	raw := os.Getenv(key)
	if raw == "" {
		return defaultValue
	}
	value, err := strconv.Atoi(raw)
	if err != nil || value <= 0 {
		return defaultValue
	}
	return value
}
