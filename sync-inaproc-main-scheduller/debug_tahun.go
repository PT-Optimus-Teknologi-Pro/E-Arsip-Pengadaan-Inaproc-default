package main

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

func DebugTahun() {
	godotenv.Load()
	tahunNow := time.Now().Year()
	tahunStart, _ := strconv.Atoi(os.Getenv("TA_START"))
	if tahunStart == 0 {
		tahunStart = tahunNow
	}
	var result []int
	for i := tahunStart; i <= tahunNow; i++ {
		result = append(result, i)
	}
	fmt.Printf("TA_START from env: %s\n", os.Getenv("TA_START"))
	fmt.Printf("Tahun Now: %d\n", tahunNow)
	fmt.Printf("Tahun List: %v\n", result)
}
