package main

import (
	"fmt"
	"log/slog"
	_ "sync-inaproc/config"
	_ "sync-inaproc/models"
	"sync-inaproc/services"
	"time"
)
func main() {
	start := time.Now()
	slog.Info(fmt.Sprintf("sync data started at %s", start.Format("2006-01-02 15:04:05")))
	services.Sync()
	selesai := time.Now()
	durasi := selesai.Sub(start)
	slog.Info(fmt.Sprintf("sync data stop at %s", selesai.Format("2006-01-02 15:04:05")))
	slog.Info(fmt.Sprintf("total executed %f detik", durasi.Seconds()))
}
