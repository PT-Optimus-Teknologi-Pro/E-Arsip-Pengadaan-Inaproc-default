package services

import (
	"fmt"
	"log/slog"
	"sync-inaproc/config"
	"sync-inaproc/models"
)

type ApiResponsePencatatan struct {
	Success		bool					`json:"success"`
	Data		[]models.Pencatatan		`json:"data"`
	Meta 		Meta					`json:"meta"`
}

func syncPencatatan() {
	slog.Info("sync pencatatan")
	tahunList := config.GetTahunList()
	for _, tahun := range tahunList {
		url := fmt.Sprintf("%s?kode_klpd=%s&tahun=%d", PENCATATAN, config.GetKodeKlpd(), tahun)
		var response ApiResponsePencatatan
		result := make(map[uint]models.Pencatatan)
		err := fetchApi(url, &response)
		if err != nil {
			slog.Error(err.Error())
		} else if response.Success && len(response.Data) > 0 {
			for _, obj := range response.Data {
				result[obj.KdNontenderPct] = obj
			}
			cursor := response.Meta.Cursor
			hasmore := response.Meta.HasMore
			for hasmore && len(cursor) > 0 {
				var responseNext ApiResponsePencatatan
				err = fetchApiCursor(url, cursor, &responseNext)
				if err != nil {
					slog.Error(err.Error())
					hasmore = false
				}else if responseNext.Success && len(responseNext.Data) > 0 {
					for _, obj := range responseNext.Data {
						result[obj.KdNontenderPct] = obj
					}
					hasmore = responseNext.Meta.HasMore
					cursor = responseNext.Meta.Cursor
				}
			}
		}
		slog.Info("data pencatatan : ", "size", len(result), "tahun", tahun)
		models.SavePencatatan(&result, tahun)
	}
}

type ApiResponsePencatatanRealisasi struct {
	Success		bool							`json:"success"`
	Data		[]models.PencatatanRealisasi	`json:"data"`
	Meta 		Meta							`json:"meta"`
}

func syncPencatatanRealisasi() {
	slog.Info("sync pencatatan realisasi")
	tahunList := config.GetTahunList()
	for _, tahun := range tahunList {
		url := fmt.Sprintf("%s?kode_klpd=%s&tahun=%d", PENCATATAN_REALISASI, config.GetKodeKlpd(), tahun)
		var response ApiResponsePencatatanRealisasi
		result := make(map[uint]models.PencatatanRealisasi)
		err := fetchApi(url, &response)
		if err != nil {
			slog.Error(err.Error())
		} else if response.Success && len(response.Data) > 0 {
			for _, obj := range response.Data {
				result[obj.KdNontenderPct] = obj
			}
			cursor := response.Meta.Cursor
			hasmore := response.Meta.HasMore
			for hasmore && len(cursor) > 0 {
				var responseNext ApiResponsePencatatanRealisasi
				err = fetchApiCursor(url, cursor, &responseNext)
				if err != nil {
					slog.Error(err.Error())
					hasmore = false
				}else if responseNext.Success && len(responseNext.Data) > 0 {
					for _, obj := range responseNext.Data {
						result[obj.KdNontenderPct] = obj
					}
					hasmore = responseNext.Meta.HasMore
					cursor = responseNext.Meta.Cursor
				}
			}
		}
		slog.Info("data realisasi pencatatan : ", "size", len(result), "tahun", tahun)
		models.SavePencatatanRealisasi(&result, tahun)
	}
}
