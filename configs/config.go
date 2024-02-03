package configs

import (
	"github.com/spf13/viper"
)

type conf struct {
	BdType               string `mapstructure:"BD_TYPE"`
	BdAddr               string `mapstructure:"BD_ADDR"`
	BdPort               string `mapstructure:"BD_PORT"`
	IPMaxNumberAccess    string `mapstructure:"IP_MAX_NUMBER_ACCESS"`
	IPTimeLimit          string `mapstructure:"IP_TIME_LIMIT"`
	IPTimeBlock          string `mapstructure:"IP_TIME_BLOCK"`
	TokenMaxNumberAccess string `mapstructure:"TOKEN_MAX_NUMBER_ACCESS"`
	TokenTimeLimit       string `mapstructure:"TOKEN_TIME_LIMIT"`
	TokenTimeBlock       string `mapstructure:"TOKEN_TIME_BLOCK"`
	TokenFileLimits      string `mapstructure:"TOKEN_FILE_LIMITS"`
}

func LoadConfig(path string) (*conf, error) {
	var cfg *conf

	viper.SetConfigName("app_config")
	viper.SetConfigType("env")
	viper.AddConfigPath(path)
	viper.SetConfigFile(".env")
	viper.AutomaticEnv()
	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}
	err = viper.Unmarshal(&cfg)
	if err != nil {
		panic(err)
	}

	//	cfg.TokenAuth = jwtauth.New("HS256", []byte(cfg.JWTSecret), nil)

	return cfg, err
}
