package util

import (
	"github.com/spf13/viper"
)

type Config struct {
	DbHost           string `mapstructure:"MYSQL_HOST"`
	DbPort           string `mapstructure:"MYSQL_PORT"`
	DbUsername       string `mapstructure:"MYSQL_USERNAME"`
	DbPassword       string `mapstructure:"MYSQL_PASSWORD"`
	DbName           string `mapstructure:"MYSQL_DATABASE_NAME"`
	SecretKey        string `mapstructure:"SECRET_KEY"`
	SecretKeyAnother string `mapstructure:"ANOTHER_SECRET_KEY"`
	RefreshKey       string `mapstructure:"REFRESH_TOKEN_KEY"`
}

func LoadConfig(path string, filename string) (config Config, err error) {
	println(path)
	println(filename)
	viper.AddConfigPath(path)
	viper.SetConfigName(filename)
	viper.SetConfigType("env")
	viper.AllowEmptyEnv(true)
	// viper.SetEnvPrefix("c_")
	viper.AutomaticEnv()

	err = viper.ReadInConfig()

	if err != nil {
		return
	}
	err = viper.Unmarshal(&config)

	return
}
