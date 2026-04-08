package services

import (
	"fmt"
	"log/slog"
	"sync-inaproc/config"
	"sync-inaproc/models"
)

type ApiResponseNontender struct {
	Success		bool					`json:"success"`
	Data		[]models.Nontender		`json:"data"`
	Meta 		Meta					`json:"meta"`
}

func syncNontender() {
	slog.Info("sync nontender")
	tahunList := config.GetTahunList()
	for _, tahun := range tahunList {
		url := fmt.Sprintf("%s?kode_klpd=%s&tahun=%d", NONTENDER_PENGUMUMAN, config.GetKodeKlpd(), tahun)
		var response ApiResponseNontender
		result := make(map[uint]models.Nontender)
		err := fetchApi(url, &response)
		if err != nil {
			slog.Error(err.Error())
		} else if response.Success && len(response.Data) > 0 {
			for _, obj := range response.Data {
				result[obj.KdNontender] = obj
			}
			cursor := response.Meta.Cursor
			hasmore := response.Meta.HasMore
			for hasmore && len(cursor) > 0 {
				var responseNext ApiResponseNontender
				err = fetchApiCursor(url, cursor, &responseNext)
				if err != nil {
					slog.Error(err.Error())
					hasmore = false
				}else if responseNext.Success && len(responseNext.Data) > 0 {
					for _, obj := range responseNext.Data {
						result[obj.KdNontender] = obj
					}
					hasmore = responseNext.Meta.HasMore
					cursor = responseNext.Meta.Cursor
				}
			}
		}
		slog.Info("data nontender : ", "size", len(result), "tahun", tahun)
		models.SaveNontender(&result, tahun)
	}
}

type ApiResponseNontenderSelesai struct {
	Success		bool						`json:"success"`
	Data		[]models.NontenderSelesai	`json:"data"`
	Meta 		Meta						`json:"meta"`
}

func syncNontenderSelesai() {
	slog.Info("sync nontender selesai")
	tahunList := config.GetTahunList()
	for _, tahun := range tahunList {
		url := fmt.Sprintf("%s?kode_klpd=%s&tahun=%d", NONTENDER_SELESAI, config.GetKodeKlpd(), tahun)
		var response ApiResponseNontenderSelesai
		result := make(map[uint]models.NontenderSelesai)
		err := fetchApi(url, &response)
		if err != nil {
			slog.Error(err.Error())
		} else if response.Success && len(response.Data) > 0 {
			for _, obj := range response.Data {
				result[obj.KdNontender] = obj
			}
			cursor := response.Meta.Cursor
			hasmore := response.Meta.HasMore
			for hasmore && len(cursor) > 0 {
				var responseNext ApiResponseNontenderSelesai
				err = fetchApiCursor(url, cursor, &responseNext)
				if err != nil {
					slog.Error(err.Error())
					hasmore = false
				}else if responseNext.Success && len(responseNext.Data) > 0 {
					for _, obj := range responseNext.Data {
						result[obj.KdNontender] = obj
					}
					hasmore = responseNext.Meta.HasMore
					cursor = responseNext.Meta.Cursor
				}
			}
		}
		slog.Info("data nontender selesai : ", "size", len(result), "tahun", tahun)
		models.SaveNontenderSelesai(&result, tahun)
	}
}


type ApiResponseJadwalNontender struct {
	Success		bool						`json:"success"`
	Data		[]models.JadwalNontender	`json:"data"`
	Meta 		Meta						`json:"meta"`
}

func syncJadwalNontender() {
	slog.Info("sync jadwal nontender")
	tahunList := config.GetTahunList()
	for _, tahun := range tahunList {
		url := fmt.Sprintf("%s?kode_klpd=%s&tahun=%d", NONTENDER_JADWAL, config.GetKodeKlpd(), tahun)
		var response ApiResponseJadwalNontender
		var result []models.JadwalNontender
		err := fetchApi(url, &response)
		if err != nil {
			slog.Error(err.Error())
		} else if response.Success && len(response.Data) > 0 {
			result = append(result, response.Data...)
			cursor := response.Meta.Cursor
			hasmore := response.Meta.HasMore
			for hasmore && len(cursor) > 0 {
				var responseNext ApiResponseJadwalNontender
				err = fetchApiCursor(url, cursor, &responseNext)
				if err != nil {
					slog.Error(err.Error())
					hasmore = false
				}else if responseNext.Success && len(responseNext.Data) > 0 {
					result = append(result, responseNext.Data...)
					hasmore = responseNext.Meta.HasMore
					cursor = responseNext.Meta.Cursor
				}
			}
		}
		slog.Info("data jadwal nontender : ", "size", len(result), "tahun", tahun)
		models.SaveJadwalNontender(&result, tahun)
	}
}
