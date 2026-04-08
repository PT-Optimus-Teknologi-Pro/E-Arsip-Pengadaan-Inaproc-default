package services

import (
	"fmt"
	"log/slog"
	"sync-inaproc/config"
	"sync-inaproc/models"
)

type ApiResponseRupAnggaranSwakelola struct {
	Success		bool							`json:"success"`
	Data		[]models.RupSwakelolaAnggaran	`json:"data"`
	Meta 		Meta							`json:"meta"`
}

func syncRupAnggaranSwakelola() {
	slog.Info("sync Anggaran RUP Swakelola")
	tahunList := config.GetTahunList()
	for _, tahun := range tahunList {
		url := fmt.Sprintf("%s?kode_klpd=%s&tahun=%d", RUP_PAKET_ANGGARAN_SWAKELOLA, config.GetKodeKlpd(), tahun)
		var response ApiResponseRupAnggaranSwakelola
		var result []models.RupSwakelolaAnggaran
		err := fetchApi(url, &response)
		if err != nil {
			slog.Error(err.Error())
		} else if response.Success && len(response.Data) > 0 {
			result = append(result, response.Data...)
			cursor := response.Meta.Cursor
			hasmore := response.Meta.HasMore
			for hasmore && len(cursor) > 0 {
				var responseNext ApiResponseRupAnggaranSwakelola
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
		slog.Info("data rup anggaran swakelola : ", "size", len(result), "tahun", tahun)
		models.SaveRupAnggaranSwakelola(&result, tahun)
	}
}

type ApiRupSwakelola struct {
	Success		bool					`json:"success"`
	Data		[]models.RupSwakelola	`json:"data"`
	Meta 		Meta					`json:"meta"`
}

func syncRupSwakelola() {
	slog.Info("sync RUP Swakelola")
	tahunList := config.GetTahunList()
	for _, tahun := range tahunList {
		url := fmt.Sprintf("%s?kode_klpd=%s&tahun=%d", RUP_PAKET_SWAKELOLA_TERUMUMKAN, config.GetKodeKlpd(), tahun)
		var response ApiRupSwakelola
		result := make(map[uint]models.RupSwakelola)
		err := fetchApi(url, &response)
		if err != nil {
			slog.Error(err.Error())
		} else if response.Success && len(response.Data) > 0 {
			for _, obj := range response.Data {
				result[obj.KdRup] = obj
			}
			cursor := response.Meta.Cursor
			hasmore := response.Meta.HasMore
			for hasmore && len(cursor) > 0 {
				var responseNext ApiRupSwakelola
		 		err = fetchApiCursor(url, cursor, &responseNext)
				if err != nil {
					slog.Error(err.Error())
					hasmore = false
				}else if responseNext.Success && len(responseNext.Data) > 0 {
					for _, obj := range responseNext.Data {
						result[obj.KdRup] = obj
					}
					hasmore = responseNext.Meta.HasMore
					cursor = responseNext.Meta.Cursor
				}
			}
		}
		slog.Info("data rup swakelola : ", "size", len(result), "tahun", tahun)
		models.SaveRupSwakelola(&result, tahun)
	}
}
