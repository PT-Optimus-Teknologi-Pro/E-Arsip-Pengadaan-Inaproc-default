package services

import (
	"arsip/models"
	"errors"
	"time"

	"github.com/gofiber/fiber/v2/log"
)

func GetUkpbj(id uint) models.Ukpbj {
	return models.GetUkpbj(id)
}

func CreateUkpbj(ukpbj models.Ukpbj) error {
	return models.CreateUkpbj(ukpbj)
}

func SaveUkpbj(ukpbj models.Ukpbj) error {
	return 	models.SaveUkpbj(ukpbj)
}

func DeleteUkpbj(ukpbj models.Ukpbj) error {
	return models.DeleteUkpbj(ukpbj)
}

func GetAgency(id uint) models.Agency {
	return models.GetAgency(id)
}

func CreateAgency(agency models.Agency) error {
	return models.CreateAgency(agency)
}

func SaveAgency(agency models.Agency) error {
	return models.SaveAgency(agency)
}

func DeleteAgency(agency models.Agency) error {
	return models.DeleteAgency(agency)
}

func GetCountAgency() int64 {
	return models.GetCountAgency()
}

func GetCountUkpbj() int64 {
	return models.GetCountUkpbj()
}

func GetCountPegawai() int64 {
	return models.GetCountPegawai()
}


func GetPanitias() []models.Panitia {
	return models.GetPanitias()
}

func GetPanitia(id uint) models.Panitia {
	return models.GetPanitia(id)
}

func GetAnggotaPokja(id uint) []models.AnggotaPanitia {
	return models.GetAnggotaPokja(id)
}

func CreatePanitia(panitia models.PanitiaDTO) error {
	objPanitia := models.Panitia{Nama: panitia.Nama, Tahun: panitia.Tahun}
	err := models.SavePanitia(&objPanitia)
	if err != nil {
		log.Error("error saving panitia ", err)
		return errors.New("Gagal Simpan panitia")
	}
	err = objPanitia.SaveAnggotaPokja(panitia)
	if err != nil {
		log.Error("error saving anggota panitia ", err)
		return errors.New("Gagal Simpan anggota panitia")
	}
	return nil
}

func UpdatePanitia(id uint, panitia models.PanitiaDTO) error {
	objPanitia := GetPanitia(id)
	if objPanitia.ID == 0 {
		return errors.New("Panitia tidak ditemukan")
	}
	objPanitia.Nama = panitia.Nama
	objPanitia.Tahun = panitia.Tahun
	err := models.SavePanitia(&objPanitia)
	if err != nil {
		return errors.New("Gagal Simpan panitia")
	}
	err = objPanitia.SaveAnggotaPokja(panitia)
	if err != nil {
		return errors.New("Gagal Simpan anggota panitia")
	}
	return nil
}

func DeletePanitia(id uint) error {
	panitia := GetPanitia(id)
	if panitia.ID == 0 {
		return errors.New("Panitia tidak ditemukan")
	}
	err := models.DeleteAnggotaPokja(panitia.ID)
	if err != nil {
		return err
	}
	models.DeleteAnggotaPokja(panitia.ID)
	return models.DeletePanitia(&panitia)
}

func GetPPs(satkerid uint) []models.Pegawai {
	return models.GetPPs(satkerid)
}

func GetPPKs() []models.Pegawai {
	return models.GetPpks()
}

func GetListAnggotaPokja() []models.Pegawai {
	return models.GetListAnggotaPokja()
}

func GetListAnggotaPp() []models.Pegawai {
	return models.GetListAnggotaPp()
}

func GetPejabatPengadaan(id uint) models.PejabatPengadaan {
	return models.GetPejabatPengadaan(id)
}

func SavePejabatPengadaan(dto *models.PejabatPengadaanDTO) error {
	log.Info("pp dto : ", dto.Pegawai)
	var objPP models.PejabatPengadaan
	if dto.ID > 0 {
		objPP = GetPejabatPengadaan(dto.ID)
	}
	
	objPP.Groups = dto.Groups
	objPP.NoSk = dto.NoSk
	objPP.TempatSk = dto.TempatSk
	objPP.Tahun = dto.Tahun
	objPP.UkpbjId = dto.UkpbjId
	
	// Parse dates
	if dto.PeriodeAwal != "" {
		objPP.PeriodeAwal, _ = time.Parse("2006-01-02", dto.PeriodeAwal)
	}
	if dto.PeriodeAkhir != "" {
		objPP.PeriodeAkhir, _ = time.Parse("2006-01-02", dto.PeriodeAkhir)
	}
	if dto.TglSk != "" {
		objPP.TglSk, _ = time.Parse("2006-01-02", dto.TglSk)
	}

	err := models.SavePejabatPengadaan(&objPP)
	if err != nil {
		log.Error("error saving Pejabat Pengadaan ", err)
		return err
	}
	var satkers []models.PejabatPengadaanSatker
	for _, v := range dto.Satker {
		objSatker := models.PejabatPengadaanSatker {
			PpId: objPP.ID,
			SatkerId: v,
		}
		satkers = append(satkers, objSatker)
	}
	err = objPP.SavePejabatPengadaanSatker(&satkers)
	if err != nil {
		log.Error("error saving Pejabat Pengadaan satker", err)
		return err
	}
	err = objPP.SavePejabatPengadaanPegawai(dto.Pegawai)
	if err != nil {
		log.Error("error saving Pejabat Pengadaan pegawai", err)
		return err
	}
	return nil
}

func DeletePejabatPengadaan(id uint) error {
	pp := GetPejabatPengadaan(id)
	if pp.ID == 0 {
		return errors.New("Pejabat Pengadaan tidak ditemukan")
	}
	err := models.DeletePejabatPengadaanSatker(pp.ID)
	if err != nil {
		return err
	}
	err = models.DeletePejabatPengadaanPegawai(pp.ID)
	if err != nil {
		return err
	}
	return models.DeletePejabatPengadaan(&pp)
}
