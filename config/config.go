package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	ChatGPTToken  string
	ChargeBackUrl string
	ServiceId     string
}

var config *Config

func Load() *Config {

	err := godotenv.Load("dev.env")
	if err == nil {
		fmt.Println("Load dev.env file for local dev")
	}

	if config == nil {
		if os.Getenv("CHATGPT_TOKEN") == "" { //other env value might not set as well
			_ = fmt.Errorf("CHATGPT_TOKEN is not set:")
		}

		config = &Config{
			ChatGPTToken:  os.Getenv("CHATGPT_TOKEN"),
			ChargeBackUrl: os.Getenv("CHARGE_BACK_URL"),
			ServiceId:     os.Getenv("SERVICE_ID"),
		}
	}
	return config
}
