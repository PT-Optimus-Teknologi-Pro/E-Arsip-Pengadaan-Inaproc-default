package services

import (
	"fmt"
	"log/slog"
	"sync-inaproc/config"
	"sync-inaproc/models"
)

type ApiResponseProgram struct {
	Success		bool					`json:"success"`
	Data		[]models.Program		`json:"data"`
	Meta 		Meta					`json:"meta"`
}

func syncProgram() {
	slog.Info("sync Program")
	tahunList := config.GetTahunList()
	for _, tahun := range tahunList {
		url := fmt.Sprintf("%s?kode_klpd=%s&tahun=%d", RUP_PROGRAM, config.GetKodeKlpd(), tahun)
		var response ApiResponseProgram
		result := make(map[uint]models.Program)
		err := fetchApi(url, &response)
		if err != nil {
			slog.Error(err.Error())
		} else if response.Success && len(response.Data) > 0 {
			for _, obj := range response.Data {
				result[obj.KdProgram] = obj
			}
			cursor := response.Meta.Cursor
			hasmore := response.Meta.HasMore
			for hasmore && len(cursor) > 0 {
				var responseNext ApiResponseProgram
			 	err = fetchApiCursor(url, cursor, &responseNext)
	    		if err != nil {
					slog.Error(err.Error())
					hasmore = false
				} else if responseNext.Success && len(responseNext.Data) > 0 {
					for _, obj := range responseNext.Data {
						result[obj.KdProgram] = obj
					}
					hasmore = responseNext.Meta.HasMore
					cursor = responseNext.Meta.Cursor
				}
			}
		}
		slog.Info("data program : ", "size", len(result), "tahun", tahun)
		models.SaveProgram(&result, tahun)
	}
}
