package services

import (
	"fmt"
	"log/slog"
	"sync-inaproc/config"
	"sync-inaproc/models"
)


type ApiResponseSatker struct {
	Success		bool					`json:"success"`
	Data		[]models.Satker			`json:"data"`
	Meta 		Meta					`json:"meta"`
}

func syncSatker() {
	slog.Info("sync satker")
	tahunList := config.GetTahunList()
	for _, tahun := range tahunList {
		url := fmt.Sprintf("%s?kode_klpd=%s&tahun=%d", RUP_SATKER, config.GetKodeKlpd(), tahun)
		var response ApiResponseSatker
		result := make(map[uint]models.Satker)
		err := fetchApi(url, &response)
		if err != nil {
			slog.Error(err.Error())
		} else if response.Success && len(response.Data) > 0 {
			for _, obj := range response.Data {
				result[obj.KdSatker] = obj
			}
			cursor := response.Meta.Cursor
			hasmore := response.Meta.HasMore
			for hasmore && len(cursor) > 0 {
				var responseNext ApiResponseSatker
			 	err = fetchApiCursor(url, cursor, &responseNext)
				if err != nil {
					slog.Error(err.Error())
					hasmore = false
				} else if responseNext.Success && len(responseNext.Data) > 0 {
					for _, obj := range responseNext.Data {
						result[obj.KdSatker] = obj
					}
					hasmore = responseNext.Meta.HasMore
					cursor = responseNext.Meta.Cursor
				}
			}
		}
		slog.Info("data satker : ", "size", len(result), "tahun", tahun)
		models.SaveSatker(&result)
	}
}
