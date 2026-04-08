package services

import (
	"fmt"
	"log/slog"
	"sync-inaproc/config"
	"sync-inaproc/models"
)

type ApiResponseRupAnggaran struct {
	Success		bool							`json:"success"`
	Data		[]models.RupAnggaran			`json:"data"`
	Meta 		Meta							`json:"meta"`
}

func syncRupAnggaran() {
	slog.Info("sync Anggaran RUP")
	tahunList := config.GetTahunList()
	for _, tahun := range tahunList {
		url := fmt.Sprintf("%s?kode_klpd=%s&tahun=%d", RUP_PAKET_ANGGARAN_PENYEDIA, config.GetKodeKlpd(), tahun)
		var response ApiResponseRupAnggaran
		var result []models.RupAnggaran
		err := fetchApi(url, &response)
		if err != nil {
			slog.Error(err.Error())
		} else if response.Success && len(response.Data) > 0 {
			result = append(result, response.Data...)
			cursor := response.Meta.Cursor
			hasmore := response.Meta.HasMore
			for hasmore && len(cursor) > 0 {
				var responseNext ApiResponseRupAnggaran
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
		slog.Info("data rup anggaran : ", "size", len(result), "tahun", tahun)
		models.SaveRupAnggaran(&result, tahun)
	}
}

type ApiRupTerumumkan struct {
	Success		bool					`json:"success"`
	Data		[]models.Rup			`json:"data"`
	Meta 		Meta					`json:"meta"`
}

func syncRup() {
	slog.Info("sync RUP")
	tahunList := config.GetTahunList()
	for _, tahun := range tahunList {
		url := fmt.Sprintf("%s?kode_klpd=%s&tahun=%d", RUP_PAKET_PENYEDIA_TERUMUMKAN, config.GetKodeKlpd(), tahun)
		var response ApiRupTerumumkan
		result := make(map[uint]models.Rup)
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
				var responseNext ApiRupTerumumkan
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
		slog.Info("data rup : ", "size", len(result), "tahun", tahun)
		models.SaveRup(&result, tahun)
	}
}
