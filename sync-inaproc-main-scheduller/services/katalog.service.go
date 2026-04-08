package services

import (
	"fmt"
	"log/slog"
	"sync-inaproc/config"
	"sync-inaproc/models"
)

type ApiResponseKatalog struct {
	Success		bool					`json:"success"`
	Data		[]models.Katalog		`json:"data"`
	Meta 		Meta					`json:"meta"`
}

func syncKatalog() {
	slog.Info("sync Katalog")
	tahunList := config.GetTahunList()
	for _, tahun := range tahunList {
		url := fmt.Sprintf("%s?kode_klpd=%s&tahun=%d", KATALOG_PURCHASING, config.GetKodeKlpd(), tahun)
		var response ApiResponseKatalog
		result := make(map[string]models.Katalog)
		err := fetchApi(url, &response)
		if err != nil {
			slog.Error(err.Error())
		} else if response.Success && len(response.Data) > 0 {
			for _, obj := range response.Data {
				result[obj.OrderId] = obj
			}
			cursor := response.Meta.Cursor
			hasmore := response.Meta.HasMore
			for hasmore && len(cursor) > 0 {
				var responseNext ApiResponseKatalog
				err = fetchApiCursor(url, cursor, &responseNext)
				if err != nil {
					slog.Error(err.Error())
					hasmore = false
				}else if responseNext.Success && len(responseNext.Data) > 0 {
					for _, obj := range responseNext.Data {
						result[obj.OrderId] = obj
					}
					hasmore = responseNext.Meta.HasMore
					cursor = responseNext.Meta.Cursor
				}
			}
		}
		slog.Info("data katalog : ", "size", len(result), "tahun", tahun)
		models.SaveKatalog(&result, tahun)
	}
}


type ApiResponsePenyediaKatalog struct {
	Success		bool					`json:"success"`
	Data		[]models.Penyedia			`json:"data"`
	Meta 		Meta					`json:"meta"`
}

func syncPenyediaKatalog() {
	slog.Info("sync Penyedia Katalog")
	list := models.GetListPenyedia()
	var result []models.Penyedia
	for _, obj := range list {
		// slog.Info("penyedia katalog :", "kode", obj)
		url := fmt.Sprintf("%s?kode_penyedia=%s", KATALOG_PENYEDIA, obj)
		var response ApiResponsePenyediaKatalog
		err := fetch(url, &response)
		if err != nil {
			slog.Error(err.Error())
		} else if response.Success && len(response.Data) > 0 {
			result = append(result, response.Data...)
		}
	}
	models.SavePenyedia(&result)
}
