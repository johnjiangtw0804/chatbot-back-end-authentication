package env

import "github.com/spf13/viper"

// squash tag allows flattening embedded structs
type AppConfig struct {
	Name     string `mapstructure:"name"`
	Port     string `mapstructure:"port"`
	Env      string `mapstructure:"env"`
	TimeZone string `mapstructure:"timezone"`
}

type DatabaseConfig struct {
	Host     string `mapstructure:"host"`
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
	Name     string `mapstructure:"name"`
}

type JWT struct {
	Secret string `mapstructure:"secret"`
}

type Configuration struct {
	App      AppConfig      `mapstructure:"app"`
	Database DatabaseConfig `mapstructure:"database"`
	Jwt      JWT            `mapstructure:"jwt"`
}

func LoadConfig() (*Configuration, error) {
	var config Configuration
	viper.SetConfigName("config")
	viper.AddConfigPath(".")
	if err := viper.ReadInConfig(); err != nil {
		// logger.Errorf("Error reading config file, %s", err)
		return nil, err
	}

	err := viper.Unmarshal(&config)
	if err != nil {
		// logger.Errorf("error to decode, %v", err)
		return nil, err
	}

	return &config, nil
}
