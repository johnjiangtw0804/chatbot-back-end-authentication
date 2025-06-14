package env

import (
	"log"

	"github.com/spf13/viper"
)

// https://github.com/spf13/viper
// squash tag allows flattening embedded structs
type Configuration struct {
	// App config
	AppName        string `mapstructure:"APP_NAME"`
	AppPort        string `mapstructure:"APP_PORT"`
	AppEnv         string `mapstructure:"APP_Env"`
	AppTimeZone    string `mapstructure:"APP_TIMEZONE"`
	AppFrontendURL string `mapstructure:"APP_FRONTEND_URL"`

	// DB config
	DBHost     string `mapstructure:"DB_HOST"`
	DBUser     string `mapstructure:"DB_USER"`
	DBPassword string `mapstructure:"DB_PASSWORD"`
	DBName     string `mapstructure:"DB_NAME"`
	DBPort     string `mapstructure:"DB_PORT"`

	// JWT
	JWTSecret string `mapstructure:"JWT_SECRET"`
}

func LoadConfig() (*Configuration, error) {
	var config Configuration
	viper.SetConfigName("config")
	viper.SetConfigType("env")
	viper.AddConfigPath(".")
	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Error reading config file, %s", err)
		return nil, err
	}

	err := viper.Unmarshal(&config)
	if err != nil {
		log.Fatalf("error to decode, %v", err)
		return nil, err
	}

	return &config, nil
}
