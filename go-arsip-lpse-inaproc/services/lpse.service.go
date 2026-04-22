package services

import (
	"arsip/config"
	"arsip/models"
	"encoding/json"
	"fmt"
	"net/http"
	"io"
	"bytes"
	"github.com/gofiber/fiber/v2/log"
)

const (
	isb_tender             = "https://isb.inaproc.id/isb/api/getJsonPaketTender?token=%s&lpseId=%s"
	isb_tender_selesai     = "https://isb.inaproc.id/isb/api/getJsonPaketTenderSelesai?token=%s&lpseId=%s"
	isb_nontender          = "https://isb.inaproc.id/isb/api/getJsonPaketNonTender?token=%s&lpseId=%s"
	isb_nontender_selesai  = "https://isb.inaproc.id/isb/api/getJsonPaketNonTenderSelesai?token=%s&lpseId=%s"
	isb_purchasing         = "https://isb.inaproc.id/isb/api/getJsonPaketPurchasing?token=%s&lpseId=%s"
)

func SyncLpse() {
	log.Info("Sync Data LPSE (Inaproc ISB)")
	token := config.GetToken()
	lpseId := config.GetLpseId()

	if token == "" || lpseId == "" {
		log.Warn("LPSE Token or LPSE_ID is not configured in .env")
		return
	}

	// 1. Sync Tender
	syncTender(token, lpseId)
	// 2. Sync Tender Selesai
	syncTenderSelesai(token, lpseId)
	// 3. Sync NonTender
	syncNontender(token, lpseId)
	// 4. Sync NonTender Selesai
	syncNontenderSelesai(token, lpseId)
	// 5. Sync Purchasing
	syncPurchasing(token, lpseId)
}

func syncTender(token, lpseId string) {
	url := fmt.Sprintf(isb_tender, token, lpseId)
	content, err := getLpseContent(url)
	if err == nil {
		var data []models.Tender
		if err := json.Unmarshal(content, &data); err == nil {
			models.SaveTender(&data)
			log.Info("Synced Tender: ", len(data))
		}
	}
}

func syncTenderSelesai(token, lpseId string) {
	url := fmt.Sprintf(isb_tender_selesai, token, lpseId)
	content, err := getLpseContent(url)
	if err == nil {
		var data []models.TenderSelesai
		if err := json.Unmarshal(content, &data); err == nil {
			models.SaveTenderSelesai(&data)
			log.Info("Synced Tender Selesai: ", len(data))
		}
	}
}

func syncNontender(token, lpseId string) {
	url := fmt.Sprintf(isb_nontender, token, lpseId)
	content, err := getLpseContent(url)
	if err == nil {
		var data []models.Nontender
		if err := json.Unmarshal(content, &data); err == nil {
			models.SaveNontender(&data)
			log.Info("Synced NonTender: ", len(data))
		}
	}
}

func syncNontenderSelesai(token, lpseId string) {
	url := fmt.Sprintf(isb_nontender_selesai, token, lpseId)
	content, err := getLpseContent(url)
	if err == nil {
		var data []models.NontenderSelesai
		if err := json.Unmarshal(content, &data); err == nil {
			models.SaveNontenderSelesai(&data)
			log.Info("Synced NonTender Selesai: ", len(data))
		}
	}
}

func syncPurchasing(token, lpseId string) {
	url := fmt.Sprintf(isb_purchasing, token, lpseId)
	content, err := getLpseContent(url)
	if err == nil {
		var data []models.Katalog
		if err := json.Unmarshal(content, &data); err == nil {
			models.SaveKatalogArchive(&data)
			log.Info("Synced Purchasing: ", len(data))
		}
	}
}

func FetchTenderByCode(kode uint) {
	token := config.GetToken()
	lpseId := config.GetLpseId()
	if token == "" || lpseId == "" {
		return
	}
	// Currently ISB doesn't support filtering by code in URL, so we fetch all and filter via Save/DB logic
	// In a real scenario, you might want to only fetch and filter in memory if the API is huge
	syncTender(token, lpseId)
	syncTenderSelesai(token, lpseId)
}

func FetchNontenderByCode(kode uint) {
	token := config.GetToken()
	lpseId := config.GetLpseId()
	if token == "" || lpseId == "" {
		return
	}
	syncNontender(token, lpseId)
	syncNontenderSelesai(token, lpseId)
}

func getLpseContent(url string) ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Status error: %v", resp.StatusCode)
	}
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if bytes.HasPrefix(data, []byte("<")) {
		return nil, fmt.Errorf("Error format response : %s", data)
	}
	return data, nil
}
