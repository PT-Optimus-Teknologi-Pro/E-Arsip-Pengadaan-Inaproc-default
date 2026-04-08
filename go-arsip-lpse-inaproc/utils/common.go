package utils

import (
	"crypto/md5"
	"crypto/sha512"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
	"github.com/gofiber/template/django/v3"
)

var BULAN = [12]string{"Januari", "Februari", "Maret", "April", "Mei", "Juni", "Juli", "Agustus", "September", "Oktober", "November", "Desember"}

func Setup(engine *django.Engine) {
	engine.AddFunc("rupiah", FormatRupiah)
	engine.AddFunc("number", FormatNumber)
	engine.AddFunc("formatDate", FormatDate)
	engine.AddFunc("formatDateTime", FormatDateTime)
	engine.AddFunc("StartWith", StartWith)
	engine.AddFunc("Prosentase", Prosentase)
	engine.AddFunc("bulanLabel", Bulan)
	engine.AddFunc("len", Len)
	engine.AddFunc("rating", Rating)
}

func StringToUint(s string) uint {
    return uint(StringToInt(s))
}

func UintToString(s uint) string  {
	return fmt.Sprintf("%d", s)
}

func StringToInt(s string) int {
    i, _ := strconv.Atoi(s)
    return i
}

func IntToString(s int) string  {
	return fmt.Sprintf("%d", s)
}

func StartWith(word string, prefix string) bool {
	return strings.HasPrefix(word, prefix)
}

func HashFile(filepath string) string {
	file, err := os.Open(filepath)

	if err != nil {
		panic(err)
	}

	defer file.Close()

	hash := md5.New()
	_, err = io.Copy(hash, file)

	if err != nil {
		panic(err)
	}
	return hex.EncodeToString(hash.Sum(nil))
}



func HashPassword(password string) string {
	var sha512Hasher = sha512.New()
	sha512Hasher.Write([]byte(password))
	var hashedPasswordBytes = sha512Hasher.Sum(nil)
	var hashedPasswordHex = hex.EncodeToString(hashedPasswordBytes)
	return hashedPasswordHex
}


func Bulan(bulan int) string {
	return BULAN[bulan - 1]
}


func Rating(star int, name ...string) string {
	var builder strings.Builder
	builder.WriteString("<select id=\"rating-default\"")
	// name="kualitas[]" required
	if len(name) > 0 {
		builder.WriteString(fmt.Sprintf("name=\"%s\" required", name[0]))
	} else {
			builder.WriteString("disabled")
	}
	builder.WriteString(">")
	if star == 0 {
    	builder.WriteString("<option>Select a rating</option>")
	} else {
    	builder.WriteString("<option selected>Select a rating</option>")
    }
    if star == 5 {
    	builder.WriteString("<option value=\"5\" selected>Excellent</option>")
    } else {
    	builder.WriteString("<option value=\"5\">Excellent</option>")
    }
    if star == 4 {
    	builder.WriteString("<option value=\"4\" selected>Very Good</option>")
    } else {
    	builder.WriteString("<option value=\"4\">Very Good</option>")
    }
    if star == 3 {
    	builder.WriteString("<option value=\"3\" selected>Average</option>")
    } else {
    	builder.WriteString("<option value=\"3\">Average</option>")
    }
    if star == 2 {
    	builder.WriteString("<option value=\"2\" selected>Poor</option>")
    } else {
    	builder.WriteString("<option value=\"2\">Poor</option>")
    }
    if star == 1 {
    	builder.WriteString("<option value=\"1\" selected>Terrible</option>")
    } else {
    	builder.WriteString("<option value=\"1\">Terrible</option>")
    }
    builder.WriteString("</select>")
	return builder.String()
}
func ToWebPath(p string) string {
	if p == "" {
		return ""
	}
	// Remove file:/// or other protocol prefixes if present
	if strings.HasPrefix(p, "file:///") {
		p = strings.TrimPrefix(p, "file:///")
	} else if strings.HasPrefix(p, "file://") {
		p = strings.TrimPrefix(p, "file://")
	}

	// Normalisasi file upload path
	if strings.Contains(p, "fileupload") {
		idx := strings.Index(p, "fileupload")
		res := "/" + strings.ReplaceAll(p[idx:], "\\", "/")
		return strings.ReplaceAll(res, "//", "/")
	}

	// Handle absolute paths windows
	if len(p) > 2 && p[1] == ':' && p[2] == '/' {
		// Possibly C:/... without fileupload? Return it normalized anyway
		return strings.ReplaceAll(p, "\\", "/")
	}

	return p
}
