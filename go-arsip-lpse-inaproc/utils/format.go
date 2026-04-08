package utils

import (
	"fmt"
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


func FormatRupiah(amount float64) string {
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
	persen := (value/total) * 100
	return fmt.Sprintf("%.2f", persen)+"%"
}


func Len(obj []interface{}) int {
	return len(obj)
}

func FormatNumber(amount float64) string {
	humanizeValue := humanize.CommafWithDigits(amount, 0)
	return strings.Replace(humanizeValue, ",", ".", -1)
}
