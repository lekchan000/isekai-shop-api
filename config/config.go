package config

import (
	"strings"
	"sync"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/spf13/viper"
)

type (
	Config struct {
		Server   *Server   `mapstructure:"server" Validate:"required"`
		OAuth2   *OAuth2   `mapstructure:"oauth2" Validate:"required"`
		State    *State    `mapstructure:"state" Validate:"required"`
		Database *Database `mapstructure:"database" Validate:"required"`
	}

	Server struct {
		Port         int           `mapstructure:"port" Validate:"required"`
		AllowOrigins []string      `mapstructure:"allowOrigins" Validate:"required"`
		BodyLimit    string        `mapstructure:"bodyLimit" Validate:"required"`
		TimeOut      time.Duration `mapstructure:"timeout" Validate:"required"`
	}

	OAuth2 struct {
		PlayerRedirectUrl string   `mapstructure:"playerRedirectUrl" Validate:"required"`
		AdminRedirectUrl  string   `mapstructures:"adminRedirectUrl" Validate:"required"`
		ClientID          string   `mapstructure:"clientId" Validate:"required"`
		ClientSecret      string   `mapstructure:"clientSecret" Validate:"required"`
		EndPoints         endpoint `mapstructure:"endpoints" Validate:"required"`
		Scopes            []string `mapstructure:"scopes" Validate:"required"`
		UserInfoUrl       string   `mapstructure:"userInfoUrl" Validate:"required"`
		RevokeUrl         string   `mapstructure:"revokeUrl" Validate:"required"`
	}

	endpoint struct {
		AuthUrl       string `mapstructure:"authUrl" Validate:"required"`
		TokenUrl      string `mapstructure:"tokenUrl" Validate:"required"`
		DeviceAuthUrl string `mapstructure:"deviceAuthUrl" Validate:"required"`
	}

	State struct {
		Secret   string        `mapstructure:"secret" Validate:"required"`
		ExpireAt time.Duration `mapstructure:"expiresAt" Validate:"required"`
		Issuer   string        `mapstructure:"issuer" Validate:"required"`
	}

	Database struct {
		Host     string `mapstructure:"host" Validate:"required"`
		Port     int    `mapstructure:"port" Validate:"required"`
		User     string `mapstructure:"user" Validate:"required"`
		Password string `mapstructure:"password" Validate:"required"`
		DBName   string `mapstructure:"dbname" Validate:"required"`
		SSLMode  string `mapstructure:"sslmode" Validate:"required"`
		Schema   string `mapstructure:"schema" Validate:"required"`
	}
)

var (
	once           sync.Once
	configInstance *Config
)

func ConfigGetting() *Config {
	once.Do(func() {
		viper.SetConfigName("config")
		viper.SetConfigType("yaml")
		viper.AddConfigPath("./config")
		viper.AutomaticEnv()
		viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

		if err := viper.ReadInConfig(); err != nil {
			panic(err)
		}

		if err := viper.Unmarshal(&configInstance); err != nil {
			panic(err)
		}

		validating := validator.New()

		if err := validating.Struct(configInstance); err != nil {
			panic(err)
		}
	})
	return configInstance
}
