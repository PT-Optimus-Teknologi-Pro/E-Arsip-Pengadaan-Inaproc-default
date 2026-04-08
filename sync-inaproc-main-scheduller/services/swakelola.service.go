package services

import (
	"fmt"
	"log/slog"
	"sync-inaproc/config"
	"sync-inaproc/models"
)

type ApiResponseSwakelola struct {
	Success		bool					`json:"success"`
	Data		[]models.Swakelola		`json:"data"`
	Meta 		Meta					`json:"meta"`
}

func syncSwakelola() {
	slog.Info("sync swakelola")
	tahunList := config.GetTahunList()
	for _, tahun := range tahunList {
		url := fmt.Sprintf("%s?kode_klpd=%s&tahun=%d", SWAKELOLA, config.GetKodeKlpd(), tahun)
		var response ApiResponseSwakelola
		result := make(map[uint]models.Swakelola)
		err := fetchApi(url, &response)
		if err != nil {
			slog.Error(err.Error())
		} else if response.Success && len(response.Data) > 0 {
			for _, obj := range response.Data {
				result[obj.KdSwakelolaPct] = obj
			}
			cursor := response.Meta.Cursor
			hasmore := response.Meta.HasMore
			for hasmore && len(cursor) > 0 {
				var responseNext ApiResponseSwakelola
				err = fetchApiCursor(url, cursor, &responseNext)
				if err != nil {
					slog.Error(err.Error())
					hasmore = false
				}else if responseNext.Success && len(responseNext.Data) > 0 {
					for _, obj := range responseNext.Data {
						result[obj.KdSwakelolaPct] = obj
					}
					hasmore = responseNext.Meta.HasMore
					cursor = responseNext.Meta.Cursor
				}
			}
		}
		slog.Info("data swakelola : ", "size", len(result), "tahun", tahun)
		models.SaveSwakelola(&result, tahun)
	}
}

type ApiResponseSwakelolaRealisasi struct {
	Success		bool							`json:"success"`
	Data		[]models.SwakelolaRealisasi		`json:"data"`
	Meta 		Meta							`json:"meta"`
}

func syncSwakelolaRealisasi() {
	slog.Info("sync swakelola realisasi")
	tahunList := config.GetTahunList()
	for _, tahun := range tahunList {
		url := fmt.Sprintf("%s?kode_klpd=%s&tahun=%d", SWAKELOLA_REALISASI, config.GetKodeKlpd(), tahun)
		var response ApiResponseSwakelolaRealisasi
		result := make(map[uint]models.SwakelolaRealisasi)
		err := fetchApi(url, &response)
		if err != nil {
			slog.Error(err.Error())
		} else if response.Success && len(response.Data) > 0 {
			for _, obj := range response.Data {
				result[obj.KdSwakelolaPct] = obj
			}
			cursor := response.Meta.Cursor
			hasmore := response.Meta.HasMore
			for hasmore && len(cursor) > 0 {
				var responseNext ApiResponseSwakelolaRealisasi
				err = fetchApiCursor(url, cursor, &responseNext)
				if err != nil {
					slog.Error(err.Error())
					hasmore = false
				}else if responseNext.Success && len(responseNext.Data) > 0 {
					for _, obj := range responseNext.Data {
						result[obj.KdSwakelolaPct] = obj
					}
					hasmore = responseNext.Meta.HasMore
					cursor = responseNext.Meta.Cursor
				}
			}
		}
		slog.Info("data swakelola selesai : ", "size", len(result), "tahun", tahun)
		models.SaveSwakelolaRealisasi(&result, tahun)
	}
}
