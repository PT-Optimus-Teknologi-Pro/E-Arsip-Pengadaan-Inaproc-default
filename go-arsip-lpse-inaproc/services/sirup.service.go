package services

import (
	"arsip/config"
	"arsip/models"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/gofiber/fiber/v2/log"
)

const (
	// sirup_provinsi = "https://sirup.inaproc.id/sirup/service2/getJsonProv"
	// sirup_kabupaten = "https://sirup.inaproc.id/sirup/service2/getJsonKabKot?prpId=%d"
	// sirup_daftar_kldi = "https://sirup.inaproc.id/sirup/service2/daftarKLDI"
	sirup_daftar_satker = "https://sirup.inaproc.id/sirup/service2/daftarSatkerByKLDI?kldi=%s&tahunAnggaran=%d"
	sirup_daftar_paket = "https://sirup.inaproc.id/sirup/service2/paketPenyediaPerSatkerTampilAllStatus?idSatker=%d&tahunAnggaran=%d"
	sirup_daftar_swakelola = "https://sirup.inaproc.id/sirup/service2/paketSwakelolaPerSatkerTampilAllStatus?idSatker=%d&tahunAnggaran=%d"
	sirup_struktur_anggaran = "https://sirup.inaproc.id/sirup/service2/strukturAnggaranDaerahPerSatker2?kodeKldi=%s&tahunAnggaran=%d"
	// sirup_daftar_program = "https://sirup.inaproc.id/sirup/service2/daftarProgramByKLDI?idKldi=%s&tahunAnggaran=%d"
	// sirup_daftar_akun = "https://sirup.inaproc.id/sirup/service2/daftarakunBySatker?idSatker=%s&tahunAnggaran=%d"
)

func updateSatkerSirup(tahunAnggaran int) {
	kdKldi := config.GetKdKldi()
	url := fmt.Sprintf(sirup_daftar_satker, kdKldi, tahunAnggaran)
	log.Info("kd kldi : ",kdKldi)
	// log.Info("get API satker tahun ", tahunAnggaran, " : ", url)
	content, err := getContent(url)
	if err != nil {
		log.Error(err)
		return
	}
	if len(content) > 0 {
		var satker []models.SatkerSirup
		if err := json.Unmarshal(content, &satker); err != nil {
			log.Error(err)
			return
		}
		if len(satker) > 0 {
	 		log.Info("save satker sirup...", len(satker))
			models.SaveAllSatkerSirup(&satker)
		}
	}
}

func updatePaketSirup(satkerId uint, tahunAnggaran int) {
	url := fmt.Sprintf(sirup_daftar_paket, satkerId, tahunAnggaran)
	content, err := getContent(url)
	if err != nil {
		log.Info("get API paket : ", url)
		log.Error(err)
		return
	}
	if len(content) > 0 {
		var paket []models.PaketSirup
		if err := json.Unmarshal(content, &paket); err != nil {
			log.Info("get API paket : ", url)
			log.Error(err)
			return
		}
		if len(paket) > 0 {
			for i := range paket {
				paket[i].Tahun = tahunAnggaran
				jenis := paket[i].Jenis()
				if len(jenis) > 0 {
					paket[i].JenisPaket = jenis[0].Jenisid
				}
			}
			err = models.SavePaketSirupTransaction(&paket)
			if err != nil {
				log.Error(err)
				return
			}
		}
	}
}

func updateSwakelolaSirup(satkerId uint, tahunAnggaran int) {
	url := fmt.Sprintf(sirup_daftar_swakelola, satkerId, tahunAnggaran)
	content, err := getContent(url)
	if err != nil {
		log.Error("get API swakelola : ", url)
		log.Error(err)
		return
	}
	if len(content) > 0 {
		var paket []models.SwakelolaSirup
		if err := json.Unmarshal(content, &paket); err != nil {
			log.Error(err)
			return
		}
		if len(paket) > 0 {
			for i := range paket {
				paket[i].Tahun = tahunAnggaran
			}
			log.Info("save swakelola sirup...")
			err = models.SaveSwakelolaSirupTransaction(&paket)
			if err != nil {
				log.Error(err)
				return
			}
		}
	}
}


func SyncSirup() {
	log.Info("Sync Data Sirup")
	tahunlist := GetTahunRupList()
	for _, tahun := range tahunlist {
		updateStrukturAnggaran(tahun)
		updateSatkerSirup(tahun)
		satkers := models.GetAllSatkerSirup(tahun)
		for _, satker := range satkers {
			updatePaketSirup(satker.ID, tahun)
			updateSwakelolaSirup(satker.ID, tahun)
		}
	}
}

func isPaketSirupEmpty() bool {
	return models.GetCountPaketSirup() == 0
}

func isSwakelolaSirupEmpty() bool {
	return models.GetCountSwakelolaSirup() == 0
}

func IsEmptyProvinsi() bool {
	return models.GetCountProvinsi() == 0
}

// func updateKabupaten(prpId uint) {
// 	url := fmt.Sprintf(sirup_kabupaten, prpId)
// 	log.Info("get API kabupaten : ", url)
// 	content, err := getContent(url)
// 	if err != nil {
// 		log.Error(err)
// 		return
// 	}
// 	if len(content) > 0 {
// 		var kabupaten []models.Kabupaten
// 		if err := json.Unmarshal(content, &kabupaten); err != nil {
// 			log.Error(err)
// 			return
// 		}
// 		if len(kabupaten) > 0 {
// 			log.Info("save kabupaten ...", len(kabupaten), "ID prod ", prpId)
// 			models.SaveAllKabupaten(&kabupaten)
// 		}
// 	}
// }

// func SyncProvinsiKabupaten(forceUpdate bool) {
// 	if models.GetCountProvinsi() == 0  || forceUpdate {
// 		url := fmt.Sprintf(sirup_provinsi)
// 		log.Info("get API propinsi : ", url)
// 		content, err := getContent(url)
// 		if err != nil {
// 			log.Error(err)
// 			return
// 		}
// 		if len(content) > 0 {
// 			var provinsi []models.Provinsi
// 			if err := json.Unmarshal(content, &provinsi); err != nil {
// 				log.Error(err)
// 				return
// 			}
// 			if len(provinsi) > 0 {
// 		 		log.Info("save provinsi...", len(provinsi))
// 		   		err = models.SaveAllPropinsi(&provinsi)
// 				if err != nil {
// 					log.Error(err)
// 					return
// 				}
// 				for _, obj := range provinsi {
// 					go updateKabupaten(obj.ID)
// 				}
// 			}
// 		}
// 	}
// }

func GetPaketSirup(id uint) models.PaketSirup {
	return models.GetPaketSirup(id)
}

func GetSwakelolaSirup(id uint) models.SwakelolaSirup {
	return models.GetSwakelolaSirup(id)
}

func GetSatkerAPI(tahun int) []models.APISatkerSirup {
	return models.GetSatkerAPI(tahun)
}


func getContent(url string) ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		return []byte{}, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return []byte{}, fmt.Errorf("Status error: %v", resp.StatusCode)
	}
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return []byte{}, fmt.Errorf("Read body: %v", err)
	}
	if bytes.HasPrefix(data, []byte("<")){
		return nil, fmt.Errorf("Error format response : %s", data)
	}
	return data, nil
}

func GetSatkerSirup(id uint) models.SatkerSirup {
	return models.GetSatkerSirup(id)
}

func updateStrukturAnggaran(tahunAnggaran int) {
	kdKldi := config.GetKdKldi()
	url := fmt.Sprintf(sirup_struktur_anggaran, kdKldi, tahunAnggaran)
	// log.Info("kd kldi : ",kdKldi)
	// log.Info("get API struktur anggaran tahun ", tahunAnggaran, " : ", url)
	content, err := getContent(url)
	if err != nil {
		log.Error(err)
		return
	}
	if len(content) > 0 {
		var satker []models.StrukturAnggaran
		if err := json.Unmarshal(content, &satker); err != nil {
			log.Error(err)
			return
		}
		if len(satker) > 0 {
	 		log.Info("save struktur anggaran sirup...", len(satker))
			err := models.SaveAllStrukturAnggaran(&satker)
			if err != nil {
				log.Error("Failed to save struktur anggaran: ", err)
			}
		}
	}
}
