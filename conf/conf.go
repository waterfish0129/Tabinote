package conf

import (
	"fmt"
	"github.com/joho/godotenv"
	"github.com/spf13/viper"
	"os"
	"strings"
)

func IntiConfig() {
	e := godotenv.Load()
	if e != nil {
		fmt.Println(fmt.Sprintf("error loading .env file: %s", e))
	}

	viper.SetConfigName("config")
	viper.SetConfigType("yml")
	viper.AddConfigPath("./conf/")
	err := viper.ReadInConfig()
	//允許viper讀取環境變數
	viper.AutomaticEnv()
	//對應環境變數的命名規則
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	if err != nil {
		panic(fmt.Sprintf("Fatal error config file: %s \n", err))
	}
}

func GetServerPort() string {
	//給Railway讀的
	if port := os.Getenv("SERVER_PORT"); port != "" {
		return port
	}

	// fallback to config
	if port := viper.GetString("server.port"); port != "" {
		return port
	}

	return "8087"
}
