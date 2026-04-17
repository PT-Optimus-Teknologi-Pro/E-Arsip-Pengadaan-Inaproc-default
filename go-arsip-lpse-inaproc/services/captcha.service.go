package services

import (
	"encoding/json"
	"fmt"
	"log"

	"arsip/cache"
	"arsip/utils"
	"github.com/golang/freetype/truetype"
	"github.com/wenlng/go-captcha-assets/resources/fonts/fzshengsksjw"
	"github.com/wenlng/go-captcha-assets/resources/images"
	"github.com/wenlng/go-captcha-assets/resources/shapes"
	"github.com/wenlng/go-captcha/v2/base/option"
	"github.com/wenlng/go-captcha/v2/click"
)

var clickCaptchaBuilder click.Builder

func InitCaptchaConfigs() {
	log.Println("Initialize Captcha Builder v2")
	font, err := fzshengsksjw.GetFont()
	if err != nil {
		log.Fatalf("failed to load font: %v", err)
	}

	imgs, err := images.GetImages()
	if err != nil {
		log.Fatalf("failed to load base images: %v", err)
	}

	shapeMaps, err := shapes.GetShapes()
	if err != nil {
		log.Fatalf("failed to load shapes: %v", err)
	}

	clickCaptchaBuilder = click.NewBuilder(
		click.WithRangeLen(option.RangeVal{Min: 4, Max: 6}),
		click.WithRangeVerifyLen(option.RangeVal{Min: 2, Max: 4}),
		click.WithImageSize(option.Size{Width: 320, Height: 160}),
	)

	clickCaptchaBuilder.SetResources(
		click.WithFonts([]*truetype.Font{font}),
		click.WithBackgrounds(imgs),
		click.WithShapes(shapeMaps),
	)
}

func GenerateCaptchaV2() (masterBase64 string, thumbBase64 string, token string, err error) {
	if clickCaptchaBuilder == nil {
		InitCaptchaConfigs()
	}
	capt := clickCaptchaBuilder.MakeShape()
	if capt == nil {
		return "", "", "", fmt.Errorf("failed to make capture")
	}
	captData, err := capt.Generate()
	if err != nil {
		return "", "", "", err
	}

	clickData := captData.(*click.CaptData)
	masterBase64, _ = clickData.GetMasterImage().ToBase64()
	thumbBase64, _ = clickData.GetThumbImage().ToBase64()

	token = utils.RandomStr(16)

	dots := clickData.GetData()
	// Marshal dots into json to store in cache
	jsonData, _ := json.Marshal(dots)
	cache.Set("captv2_"+token, string(jsonData))

	return masterBase64, thumbBase64, token, nil
}

// VerifyCaptchaV2 checks if the provided coordinates list (e.g. "x1,y1,x2,y2") matches the cache
func VerifyCaptchaV2(token string, dotsStr string) error {
	key := "captv2_" + token
	cachedData, found := cache.Get(key)
	if !found {
		return fmt.Errorf("captcha kadaluarsa atau token tidak valid")
	}
	cache.Delete(key)

	var dots map[int]*click.Dot
	strData, ok := cachedData.(string)
	if !ok {
		return fmt.Errorf("data captcha rusak")
	}
	err := json.Unmarshal([]byte(strData), &dots)
	if err != nil {
		return fmt.Errorf("data captcha tidak valid")
	}

	// We expect user dots in the format: "x,y,x,y..."
	// Need to check with click.CheckPoint
	// Parse dotsStr
	var userDots []int64
	var currStr string
	for _, char := range dotsStr {
		if char == ',' {
			val := utils.StringToInt(currStr)
			userDots = append(userDots, int64(val))
			currStr = ""
		} else {
			currStr += string(char)
		}
	}
	if currStr != "" {
		val := utils.StringToInt(currStr)
		userDots = append(userDots, int64(val))
	}

	if len(userDots)%2 != 0 {
		return fmt.Errorf("koordinat klik tidak lengkap")
	}

	userDotCount := len(userDots) / 2
	if userDotCount != len(dots) {
		return fmt.Errorf("jumlah klik tidak sesuai")
	}

	for i := 0; i < len(dots); i++ {
		sysDot := dots[i]
		ux := userDots[i*2]
		uy := userDots[i*2+1]

		// The verification padding allows small pixel drift.
		if !click.Validate(int(ux), int(uy), sysDot.X, sysDot.Y, sysDot.Width, sysDot.Height, 35) {
			return fmt.Errorf("posisi klik tidak akurat")
		}
	}

	return nil
}
