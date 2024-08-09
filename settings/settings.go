package settings

import (
	"GoShort/logging"
	"fmt"
	"github.com/joho/godotenv"
	"os"
	"strconv"
	"strings"
)

var (
	settingsLogger = logging.NewLogger("settings")
)

var Settings struct {
	BaseURL              string
	Debug                bool
	TimeZone             string
	MaxShortlinksPerUser uint16
	SessionValidityDays  int

	CORS struct {
		AllowedOrigins []string
	}

	Database struct {
		Name string
		User string
		Pass string
		Host string
		Port string
		Dsn  string
	}
}

func LoadSettings() {
	if err := godotenv.Load(".env"); err != nil {
		settingsLogger.Warning("No .env! We'll try without it. Please make sure your environment is set up correctly!")
	} else {
		settingsLogger.Success(".env file loaded!")
	}

	Settings.BaseURL = os.Getenv("BASE_URL")
	Settings.Debug = os.Getenv("DEBUG") == "true"
	Settings.TimeZone = os.Getenv("TIMEZONE")
	Settings.CORS.AllowedOrigins = strings.Split(os.Getenv("CORS_ALLOWED_ORIGINS"), ",")
	maxShortLinksPerUserInt, err := strconv.Atoi(os.Getenv("MAX_SHORTLINKS_PER_USER"))
	if err != nil {
		settingsLogger.Error("Error loading MAX_SHORTLINKS_PER_USER")
	}
	Settings.MaxShortlinksPerUser = uint16(maxShortLinksPerUserInt)

	sessionValidityDaysInt, err := strconv.Atoi(os.Getenv("SESSION_VALIDITY_DAYS"))
	if err != nil {
		settingsLogger.Error("Error loading SESSION_VALIDITY_DAYS")
	}
	Settings.SessionValidityDays = sessionValidityDaysInt

	Settings.Database.Name = os.Getenv("DB_NAME")
	Settings.Database.User = os.Getenv("DB_USER")
	Settings.Database.Pass = os.Getenv("DB_PASS")
	Settings.Database.Host = os.Getenv("DB_HOST")
	Settings.Database.Port = os.Getenv("DB_PORT")

	Settings.Database.Dsn = fmt.Sprintf(
		"host=%s port=%s user=%s dbname=%s password=%s sslmode=disable timezone=%s client_encoding=UTF8",
		Settings.Database.Host,
		Settings.Database.Port,
		Settings.Database.User,
		Settings.Database.Name,
		Settings.Database.Pass,
		Settings.TimeZone,
	)
}
