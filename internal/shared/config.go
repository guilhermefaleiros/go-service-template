package shared

import (
	"fmt"
	"github.com/spf13/viper"
	"os"
	"strings"
)

type Config struct {
	App struct {
		Name        string `mapstructure:"name"`
		Environment string `mapstructure:"environment"`
	} `mapstructure:"app"`

	API struct {
		Port        int    `mapstructure:"port"`
		HealthCheck string `mapstructure:"healthCheck"`
	} `mapstructure:"api"`

	Observability struct {
		Jaeger struct {
			Endpoint string `mapstructure:"endpoint"`
		} `mapstructure:"jaeger"`
	} `mapstructure:"observability"`

	Postgres struct {
		Host          string `mapstructure:"host"`
		Port          int    `mapstructure:"port"`
		User          string `mapstructure:"user"`
		Password      string `mapstructure:"password"`
		Name          string `mapstructure:"name"`
		MaxConnection int32  `mapstructure:"max_connections"`
		MinConnection int32  `mapstructure:"min_connections"`
		MaxIdleTime   int32  `mapstructure:"max_idle_time"`
		MaxLifeTime   int32  `mapstructure:"max_conn_lifetime"`
	} `mapstructure:"postgres"`
}

func LoadConfig(environment string) (*Config, error) {
	v := viper.New()

	v.SetConfigType("yaml")

	basePath, err := os.Getwd()
	if err != nil {
		return nil, fmt.Errorf("failed to get current working directory: %w", err)
	}

	v.AddConfigPath(fmt.Sprintf("%s/config", basePath))
	v.SetConfigName(environment)

	v.AutomaticEnv()
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	if err := v.MergeInConfig(); err != nil {
		return nil, fmt.Errorf("error reading %s config: %w", environment, err)
	}

	for _, key := range v.AllKeys() {
		value := v.GetString(key)
		if strings.Contains(value, "${") {
			v.Set(key, os.ExpandEnv(value))
		}
	}

	var config Config
	if err := v.Unmarshal(&config); err != nil {
		return nil, fmt.Errorf("error unmarshaling config: %w", err)
	}

	return &config, nil
}
