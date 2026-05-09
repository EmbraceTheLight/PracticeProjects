package config

import (
	"flag"
	"fmt"
	"github.com/gogf/gf/util/gconv"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v3"
	"os"
	"time"
)

type Config struct {
	HttpServer *Server `yaml:"server"`
	Mysql      *Mysql  `yaml:"mysql"`
	Redis      *Redis  `yaml:"redis"`
}
type Server struct {
	HttpPort        int    `yaml:"http_port"`
	Env             string `yaml:"env"`
	EnableFullPPROF bool   `yaml:"enable_full_pprof"`
	LogLevel        string `yaml:"log_level"`
}
type Mysql struct {
	Host     string `yaml:"host"`
	Port     string `yaml:"port"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
	Database string `yaml:"database"`
	Charset  string `yaml:"charset"`
	MaxIdle  int    `yaml:"max_idle"`
	MaxOpen  int    `yaml:"max_open"`
	LogLevel int    `yaml:"log_level"`
}

type Redis struct {
	Host         string        `yaml:"host"`
	Port         int           `yaml:"port"`
	Password     string        `yaml:"password"`
	DialTimeout  time.Duration `yaml:"dial_timeout"`
	ReadTimeout  time.Duration `yaml:"read_timeout"`
	WriteTimeout time.Duration `yaml:"write_timeout"`
	PoolTimeout  time.Duration `yaml:"pool_timeout"`
	PoolSize     int32         `yaml:"pool_size"`
	MaxRetries   int32         `yaml:"max_retries"`
	MinIdleConns int32         `yaml:"min_idle_conns"`
	MaxIdle      int32         `yaml:"max_idle"`
	MaxOpen      int32         `yaml:"max_open"`
}

const (
	ServerName     string = "mall"
	ServerFullName        = "edu.mall"
)

var (
	etcdKey         = fmt.Sprintf("/config/%s/systerm", ServerFullName)
	etcdAddr        string
	localConfigPath string
	GlobalConfig    Config
)

func init() {
	flag.StringVar(&localConfigPath, "c", ServerName+"_local.yaml", "default config path")
	flag.StringVar(&etcdAddr, "r", os.Getenv("ETCD_ADDR"), "default consul address")
}

func InitConfig() *Config {
	var (
		err      error
		tempConf = &Config{}
		vipConf  = viper.New()
	)
	flag.Parse()

	if etcdAddr != "" {
		tempConf, err = getFromRemoteAndWatchUpdate(vipConf)
		if err != nil {
			panic(err)
		}
		return tempConf
	}
	tempConf, err = getFromLocal()
	if err != nil {
		panic(err)
	}
	return tempConf
}

func getFromRemoteAndWatchUpdate(v *viper.Viper) (*Config, error) {
	tempConf := Config{}
	if err := v.AddRemoteProvider("etcd3", etcdAddr, etcdKey); err != nil {
		return nil, err
	}
	if err := v.ReadRemoteConfig(); err != nil {
		return nil, err
	}
	if err := v.Unmarshal(&tempConf); err != nil {
		return nil, err
	}
	go func() {
		for {
			time.Sleep(1 * time.Minute)
			if err := v.WatchRemoteConfig(); err != nil {
				_ = v.Unmarshal(&GlobalConfig)
				fmt.Println(">>> etcd config hot-update ", gconv.String(&GlobalConfig))
			}
		}
	}()
	return &tempConf, nil
}

func getFromLocal() (*Config, error) {
	tempConf := Config{}
	if _, err := os.Stat(localConfigPath); err != nil {
		return nil, fmt.Errorf("local config not found. file name: %s", localConfigPath)
	}

	bytes, err := os.ReadFile(localConfigPath)
	if err != nil {
		return nil, fmt.Errorf("local config open failed. file name: %s", localConfigPath)
	}
	err = yaml.Unmarshal(bytes, &tempConf)
	if err != nil {
		return nil, err
	}

	return &tempConf, nil
}
