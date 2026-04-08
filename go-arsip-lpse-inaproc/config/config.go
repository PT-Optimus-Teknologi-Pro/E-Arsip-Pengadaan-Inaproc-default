package config

import (
	"os"
	"strconv"
	"strings"
	"time"

	"io"
	"github.com/gofiber/fiber/v2/log"
	"github.com/joho/godotenv"
)

var Logfile *os.File

func init() {
	_ = godotenv.Load(".env")
	var err error
	Logfile, err = os.OpenFile("web-app.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening log file: %v", err)
	}
	log.SetOutput(io.MultiWriter(os.Stdout, Logfile))
	log.Info("config init...")
}

func Port() string {
	return os.Getenv("HTTP_PORT")
}

func IsModeDev() bool {
	mode := os.Getenv("APP_MODE")
	if mode == ""{
		return false
	}
	return !strings.EqualFold(mode, "prod")
}

func GetDbUrl() string {
	return os.Getenv("DB_URL")
}

func UploadPath() string {
	return os.Getenv("UPLOAD_PATH")
}

func TahunStart() int {
	tahun, err := strconv.Atoi(os.Getenv("TA_START"))
	if err != nil {
		log.Error("failed convert TA_START")
	}
	return tahun
}

func CronJob() string {
	cron := os.Getenv("CRONJOB")
	if len(cron) == 0 {
		cron = "0 0 * * *" // default 00:00 everyday
	}
	return cron
}

func GetIsbService(api string) string {
	return os.Getenv(api)
}

func GetTahunList() []int {
	tahunNow := time.Now().Year()
	tahunStart, _ := strconv.Atoi(os.Getenv("TA_START"))
	if tahunStart == 0 {
		tahunStart = tahunNow
	}
	var result []int
	for i :=tahunStart;i<=tahunNow;i++{
		result = append(result, i)
	}
	return result
}

func GetKdKldi() string {
	kd := os.Getenv("KD_KLDI")
	if kd == "" {
		kd = os.Getenv("KODE_KLPD")
	}
	if kd == "" {
		return "D291"
	}
	return kd
}
