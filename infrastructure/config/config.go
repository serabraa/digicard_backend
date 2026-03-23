package config

import (
	"os"
)

type Config struct {
	Port                string
	AssetsDir           string
	Organization        string
	LogoText            string
	PassTypeID          string
	TeamID              string
	BackgroundRGB       string
	LabelColorRGB       string
	ForegroundRGB       string
	PassCertP12Path     string
	PassCertP12Password string
	WWDRCertPath        string
}

func Load() Config {
	return Config{
		Port:          getEnv("PORT", "8080"),
		AssetsDir:     getEnv("ASSETS_DIR", "assets/wallet"),
		Organization:  getEnv("ORGANIZATION_NAME", "UstaDev"),
		LogoText:      getEnv("LOGO_TEXT", "DigiCard"),
		PassTypeID:    getEnv("PASS_TYPE_ID", ""),
		TeamID:        getEnv("TEAM_ID", ""),
		BackgroundRGB: getEnv("BACKGROUND_COLOR", "rgb(255,255,255)"),
		LabelColorRGB: getEnv("LABEL_COLOR", "rgb(100,100,100)"),
		ForegroundRGB: getEnv("FOREGROUND_COLOR", "rgb(0,0,0)"),
	}
}

func getEnv(key, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}
	return value
}
