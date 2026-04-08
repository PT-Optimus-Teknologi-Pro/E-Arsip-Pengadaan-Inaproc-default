package services

import (
	"fmt"
	"log/slog"
	"sync-inaproc/config"
	"sync-inaproc/models"
)

type ApiResponseTender struct {
	Success		bool					`json:"success"`
	Data		[]models.Tender			`json:"data"`
	Meta 		Meta					`json:"meta"`
}

func syncTender() {
	slog.Info("sync tender")
	tahunList := config.GetTahunList()
	for _, tahun := range tahunList {
		url := fmt.Sprintf("%s?kode_klpd=%s&tahun=%d", TENDER_PENGUMUMAN, config.GetKodeKlpd(), tahun)
		var response ApiResponseTender
		result := make(map[uint]models.Tender)
		err := fetchApi(url, &response)
		if err != nil {
			slog.Error(err.Error())
		} else if response.Success && len(response.Data) > 0 {
			for _, obj := range response.Data {
				result[obj.KdTender] = obj
			}
			cursor := response.Meta.Cursor
			hasmore := response.Meta.HasMore
			for hasmore && len(cursor) > 0 {
				var responseNext ApiResponseTender
				err = fetchApiCursor(url, cursor, &responseNext)
				if err != nil {
					slog.Error(err.Error())
					hasmore = false
				}else if responseNext.Success && len(responseNext.Data) > 0 {
					for _, obj := range responseNext.Data {
						result[obj.KdTender] = obj
					}
					hasmore = responseNext.Meta.HasMore
					cursor = responseNext.Meta.Cursor
				}
			}
		}
		slog.Info("data tender : ", "size", len(result), "tahun", tahun)
		models.SaveTender(&result, tahun)
	}
}

type ApiResponseTenderSelesai struct {
	Success		bool						`json:"success"`
	Data		[]models.TenderSelesai		`json:"data"`
	Meta 		Meta						`json:"meta"`
}

func syncTenderSelesai() {
	slog.Info("sync tender selesai")
	tahunList := config.GetTahunList()
	for _, tahun := range tahunList {
		url := fmt.Sprintf("%s?kode_klpd=%s&tahun=%d", TENDER_SELESAI_NILAI, config.GetKodeKlpd(), tahun)
		var response ApiResponseTenderSelesai
		result := make(map[uint]models.TenderSelesai)
		err := fetchApi(url, &response)
		if err != nil {
			slog.Error(err.Error())
		} else if response.Success && len(response.Data) > 0 {
			for _, obj := range response.Data {
				result[obj.KdTender] = obj
			}
			cursor := response.Meta.Cursor
			hasmore := response.Meta.HasMore
			for hasmore && len(cursor) > 0 {
				var responseNext ApiResponseTenderSelesai
				err = fetchApiCursor(url, cursor, &responseNext)
				if err != nil {
					slog.Error(err.Error())
					hasmore = false
				}else if responseNext.Success && len(responseNext.Data) > 0 {
					for _, obj := range responseNext.Data {
						result[obj.KdTender] = obj
					}
					hasmore = responseNext.Meta.HasMore
					cursor = responseNext.Meta.Cursor
				}
			}
		}
		slog.Info("data tender selesai : ", "size", len(result), "tahun", tahun)
		models.SaveTenderSelesai(&result, tahun)
	}
}


type ApiResponseJadwal struct {
	Success		bool						`json:"success"`
	Data		[]models.Jadwal				`json:"data"`
	Meta 		Meta						`json:"meta"`
}

func syncJadwal() {
	slog.Info("sync jadwal")
	tahunList := config.GetTahunList()
	for _, tahun := range tahunList {
		url := fmt.Sprintf("%s?kode_klpd=%s&tahun=%d", TENDER_JADWAL, config.GetKodeKlpd(), tahun)
		var response ApiResponseJadwal
		var result []models.Jadwal
		err := fetchApi(url, &response)
		if err != nil {
			slog.Error(err.Error())
		} else if response.Success && len(response.Data) > 0 {
			result = append(result, response.Data...)
			cursor := response.Meta.Cursor
			hasmore := response.Meta.HasMore
			for hasmore && len(cursor) > 0 {
				var responseNext ApiResponseJadwal
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
		slog.Info("data jadwal tender : ", "size", len(result), "tahun", tahun)
		models.SaveJadwal(&result, tahun)
	}
}

type ApiResponsePeserta struct {
	Success		bool						`json:"success"`
	Data		[]models.Peserta			`json:"data"`
	Meta 		Meta						`json:"meta"`
}

func syncPeserta() {
	slog.Info("sync peserta")
	tahunList := config.GetTahunList()
	for _, tahun := range tahunList {
		url := fmt.Sprintf("%s?kode_klpd=%s&tahun=%d", TENDER_PESERTA, config.GetKodeKlpd(), tahun)
		var response ApiResponsePeserta
		result := make(map[uint]models.Peserta)
		err := fetchApi(url, &response)
		if err != nil {
			slog.Error(err.Error())
		} else if response.Success && len(response.Data) > 0 {
			for _, obj := range response.Data {
				result[obj.KdPeserta] = obj
			}
			cursor := response.Meta.Cursor
			hasmore := response.Meta.HasMore
			for hasmore && len(cursor) > 0 {
				var responseNext ApiResponsePeserta
				err = fetchApiCursor(url, cursor, &responseNext)
				if err != nil {
					slog.Error(err.Error())
					hasmore = false
				}else if responseNext.Success && len(responseNext.Data) > 0 {
					for _, obj := range responseNext.Data {
						result[obj.KdPeserta] = obj
					}
					hasmore = responseNext.Meta.HasMore
					cursor = responseNext.Meta.Cursor
				}
			}
		}
		slog.Info("data peserta : ", "size", len(result), "tahun", tahun)
		models.SavePeserta(&result, tahun)
	}
}
