package config

import (
	"github.com/spf13/viper"
	"log"
)

type Config struct {
	App struct {
		Name string
		Port string
	}

	Database struct {
		Dsn          string
		MaxIdleConns int
		MaxOpenConns int
	}

	Redis struct {
		Addr     string
		DB       int
		Password string
	}
}

var AppConig *Config

func InitConfig() {
	viper.SetConfigName("config")   //  设置配置文件名（不带拓展）
	viper.SetConfigType("yml")      // 明确配置的文件格式
	viper.AddConfigPath("./config") // 配置文件路径（当前目录）

	// 读取配置文件
	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Error reading config file, %v", err) // %v占位符: 相应值的默认格式
	}

	AppConig = &Config{}

	if err := viper.Unmarshal(AppConig); err != nil {
		log.Fatalf("Error unmarshalling config, %v", err)
	}

	initDB()

	InitRedis()
}
