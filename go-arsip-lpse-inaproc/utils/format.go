package utils

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/dustin/go-humanize"
)


func FormatDate(obj time.Time) string {
	if obj.IsZero() {
		return ""
	}
	return obj.Format("02-01-2006")
}

func FormatDateTime(obj time.Time) string{
	if obj.IsZero() {
		return ""
	}
	return obj.Format("02-01-2006 15:04")
}


func toFloat(i interface{}) float64 {
	switch v := i.(type) {
	case int:
		return float64(v)
	case int32:
		return float64(v)
	case int64:
		return float64(v)
	case float32:
		return float64(v)
	case float64:
		return v
	case string:
		f, _ := strconv.ParseFloat(v, 64)
		return f
	default:
		return 0
	}
}

func FormatRupiah(i interface{}) string {
	amount := toFloat(i)
	if amount == 0 {
		return ""
	}
	humanizeValue := humanize.CommafWithDigits(amount, 0)
	stringValue := strings.Replace(humanizeValue, ",", ".", -1)
	return "Rp " + stringValue
}

func Prosentase(value float64, total float64) string {
	if total == 0 {
		return ""
	}
	persen := (value / total) * 100
	return fmt.Sprintf("%.2f", persen) + "%"
}

func Len(obj []interface{}) int {
	return len(obj)
}

func FormatNumber(i interface{}) string {
	amount := toFloat(i)
	humanizeValue := humanize.CommafWithDigits(amount, 0)
	return strings.Replace(humanizeValue, ",", ".", -1)
}
