package services

import (
	"fmt"
	"log/slog"
	"sync-inaproc/config"
	"sync-inaproc/models"
)

type ApiResponseKontrak struct {
	Success		bool					`json:"success"`
	Data		[]models.Kontrak		`json:"data"`
	Meta 		Meta					`json:"meta"`
}

func syncKontrak() {
	slog.Info("sync kontrak")
	tahunList := config.GetTahunList()
	for _, tahun := range tahunList {
		url := fmt.Sprintf("%s?kode_klpd=%s&tahun=%d", TENDER_KONTRAK, config.GetKodeKlpd(), tahun)
		var response ApiResponseKontrak
		result := make(map[string]models.Kontrak)
		err := fetchApi(url, &response)
		if err != nil {
			slog.Error(err.Error())
		} else if response.Success && len(response.Data) > 0 {
			for _, obj := range response.Data {
				result[obj.NoKontrak] = obj
			}
			cursor := response.Meta.Cursor
			hasmore := response.Meta.HasMore
			for hasmore && len(cursor) > 0 {
				var responseNext ApiResponseKontrak
				err = fetchApiCursor(url, cursor, &responseNext)
				if err != nil {
					slog.Error(err.Error())
					hasmore = false
				}else if responseNext.Success && len(responseNext.Data) > 0 {
					for _, obj := range responseNext.Data {
						result[obj.NoKontrak] = obj
					}
					hasmore = responseNext.Meta.HasMore
					cursor = responseNext.Meta.Cursor
				}
			}
		}
		slog.Info("data kontrak tender : ", "size", len(result), "tahun", tahun)
		models.SaveKontrak(&result, tahun)
	}
}

type ApiResponseKontrakNontender struct {
	Success		bool						`json:"success"`
	Data		[]models.KontrakNontender	`json:"data"`
	Meta 		Meta						`json:"meta"`
}

func syncKontrakNontender() {
	slog.Info("sync kontrak nontender")
	tahunList := config.GetTahunList()
	for _, tahun := range tahunList {
		url := fmt.Sprintf("%s?kode_klpd=%s&tahun=%d", NONTENDER_KONTRAK, config.GetKodeKlpd(), tahun)
		var response ApiResponseKontrakNontender
		result := make(map[string]models.KontrakNontender)
		err := fetchApi(url, &response)
		if err != nil {
			slog.Error(err.Error())
		} else if response.Success && len(response.Data) > 0 {
			for _, obj := range response.Data {
				result[obj.NoKontrak] = obj
			}
			cursor := response.Meta.Cursor
			hasmore := response.Meta.HasMore
			for hasmore && len(cursor) > 0 {
				var responseNext ApiResponseKontrakNontender
				err = fetchApiCursor(url, cursor, &responseNext)
				if err != nil {
					slog.Error(err.Error())
					hasmore = false
				}else if responseNext.Success && len(responseNext.Data) > 0 {
					for _, obj := range responseNext.Data {
						result[obj.NoKontrak] = obj
					}
					hasmore = responseNext.Meta.HasMore
					cursor = responseNext.Meta.Cursor
				}
			}
		}
		slog.Info("data kontrak nontender : ", "size", len(result), "tahun", tahun)
		models.SaveKontrakNontender(&result, tahun)
	}
}
