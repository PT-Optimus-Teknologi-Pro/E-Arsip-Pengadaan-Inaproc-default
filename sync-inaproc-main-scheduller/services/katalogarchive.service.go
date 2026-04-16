package services

import (
	"fmt"
	"log/slog"
	"sync-inaproc/config"
	"sync-inaproc/models"
)

type ApiResponseKatalogArchive struct {
	Success		bool					`json:"success"`
	Data		[]models.KatalogArchive	`json:"data"`
	Meta 		Meta					`json:"meta"`
}

func syncKatalogArchive() {
	slog.Info("sync katalog archive")
	tahunList := config.GetTahunList()
	for _, tahun := range tahunList {
		url := fmt.Sprintf("%s?kode_klpd=%s&tahun=%d", KATALOG_PURCHASING_V5, config.GetKodeKlpd(), tahun)
		var response ApiResponseKatalogArchive
		var result []models.KatalogArchive
		err := fetchApi(url, &response)
		if err != nil {
			slog.Error(err.Error())
		} else if response.Success && len(response.Data) > 0 {
			result = append(result, response.Data...)
			cursor := response.Meta.Cursor
			hasmore := response.Meta.HasMore
			for hasmore && len(cursor) > 0 {
				var responseNext ApiResponseKatalogArchive
				err = fetchApiCursor(url, cursor, &responseNext)
				if err != nil {
					slog.Error(err.Error())
					hasmore = false
				} else if responseNext.Success && len(responseNext.Data) > 0 {
					result = append(result, responseNext.Data...)
					hasmore = responseNext.Meta.HasMore
					cursor = responseNext.Meta.Cursor
				}
			}
		}
		slog.Info("data katalog archive : ", "size", len(result), "tahun", tahun)
		models.SaveKatalogArchive(&result, tahun)
	}
}


type ApiResponsePenyediaArchive struct {
	Success		bool						`json:"success"`
	Data		[]models.PenyediaArchive	`json:"data"`
	Meta 		Meta						`json:"meta"`
}

func syncPenyediaArchive() {
	slog.Info("sync penyedia katalog archive")
	list := models.GetListPenyediaArchive()
	var result []models.PenyediaArchive
	for _, obj := range list {
		url := fmt.Sprintf("%s?kode_penyedia=%s", KATALOG_PENYEDIA_v5, obj)
		var response ApiResponsePenyediaArchive
		err := fetchApi(url, &response)
		if err != nil {
			slog.Error(err.Error())
		} else if response.Success && len(response.Data) > 0 {
			result = append(result, response.Data...)
		}
	}
	if len(result) > 0 {
		models.SavePenyediaArchive(&result)
	}
}
