package config

import (
	"log"
	"os"
	"path/filepath"

	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

type Config struct {
	App struct {
		Name string `yaml:"name"`
		Port string `yaml:"port"`
	} `yaml:"app"`
	Database struct {
		Host     string `yaml:"host"`
		Port     string `yaml:"port"`
		Username string `yaml:"username"`
		Password string `yaml:"password"`
		Name     string `yaml:"name"`
	} `yaml:"database"`
	Redis struct {
		Host     string `yaml:"host"`
		Port     string `yaml:"port"`
		Password string `yaml:"password"`
		DB       int    `yaml:"db"`
	} `yaml:"redis"`
}

var AppConfig *Config

func tryLoadEnv(paths []string) bool {
	for _, p := range paths {
		clean := filepath.Clean(p)
		if fi, err := os.Stat(clean); err == nil && !fi.IsDir() {
			if err := godotenv.Load(clean); err == nil {
				return true
			}
		}
	}
	return false
}

func loadDotEnv() {
	cwd, _ := os.Getwd()
	exePath, _ := os.Executable()
	exeDir := filepath.Dir(exePath)

	candidates := []string{
		filepath.Join(cwd, ".env"),
		filepath.Join(cwd, "..", ".env"),
		filepath.Join(cwd, "..", "..", ".env"),
		filepath.Join(exeDir, ".env"),
		filepath.Join(exeDir, "..", ".env"),
		filepath.Join(exeDir, "..", "..", ".env"),
	}
	_ = tryLoadEnv(candidates)
}

func InitConfig() {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./config")

	// 读取配置文件
	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Error reading config file: %v", err)
	}

	// 解析配置
	AppConfig = &Config{}
	if err := viper.Unmarshal(AppConfig); err != nil {
		log.Fatalf("Fatal error config file: %v", err)
	}

	// 加载环境变量（探测工作目录与可执行目录）
	loadDotEnv()
	// 初始化数据库
	InitDB()
	InitRedis()
	InitOSS()
}
