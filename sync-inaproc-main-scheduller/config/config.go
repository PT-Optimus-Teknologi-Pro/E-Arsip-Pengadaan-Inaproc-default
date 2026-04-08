package config

import (
	"log"
	"log/slog"
	"os"
	"strconv"
	"time"

	"io"
	"github.com/joho/godotenv"
)

func init() {
	_ = godotenv.Load()
	file, err := os.OpenFile("scheduler.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening log file: %v", err)
	}

	// Create a TextHandler with a specific minimum log level (e.g., LevelDebug)
	handler := slog.NewTextHandler(io.MultiWriter(os.Stdout, file), &slog.HandlerOptions{
		Level: slog.LevelDebug, // This enables logs from Debug to Error
	})
	// Set this handler as the default for the entire application
	slog.SetDefault(slog.New(handler))
}

func GetDbUrl() string {
	return os.Getenv("DB_URL")
}

func GetToken() string {
	return os.Getenv("TOKEN")
}

func GetKodeKlpd() string {
	return os.Getenv("KODE_KLPD")
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

func GetDelay() int {
	delay, _ := strconv.Atoi(os.Getenv("DELAY"))
	if delay == 0 {
		delay = 3 // default 3 second
	}
	return delay
}
