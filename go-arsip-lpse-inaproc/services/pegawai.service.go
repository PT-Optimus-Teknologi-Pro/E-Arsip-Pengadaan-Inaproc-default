package services

import (
	"arsip/models"
	"arsip/utils"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
)

func Otentikasi(userid, password string) (models.Pegawai, error) {
	var user models.Pegawai
	if userid == "" || password == "" {
		return user, errors.New("Invalid password or userid!")
	}
	hashedPassword := utils.HashPassword(password)
	userid = strings.ToUpper(userid)
	user = models.GetPegawaiByUserPass(userid, hashedPassword)
	if user.ID == 0 {
		return user, errors.New("Invalid password or userid!")
	}
	return user, nil
}


func GetPegawai(id uint) models.Pegawai {
	return models.GetPegawai(id)
}

func GetPegawaiPengadaan() []models.Pegawai {
	return models.GetPegawaiPengadaan()
}

func GetPegawaiNonPengadaan() []models.Pegawai {
	return models.GetPegawaiNonPengadaan()
}

func DeletePegawai(id uint) error {
	user := GetPegawai(id)
	if user.ID == 0 {
		return errors.New("Pegawai tidak ditemukan")
	}
	err := models.DeletePegawai(&user)
	if err != nil {
		log.Error(err)
		return errors.New("Hapus Pegawai Gagal")
	}
	return nil
}

func CreatePegawai(pegawai *models.Pegawai) error {
	pegawai.PegNamauser = strings.ToUpper(pegawai.PegNamauser)
	pegawai.Passw = utils.HashPassword(pegawai.Passw)
	return models.SavePegawai(pegawai)
}

func UpdatePegawai(pegawai models.Pegawai) error {
	if(len(pegawai.Passw) == 0 && pegawai.ID > 0){
		pegExist := GetPegawai(pegawai.ID)
		pegawai.Passw = pegExist.Passw
	}
	if len(pegawai.Passw) != 128 {
		pegawai.Passw = utils.HashPassword(pegawai.Passw)
	}
	log.Info("pegawai isActive ", pegawai.PegIsactive)
	return models.SavePegawai(&pegawai)
}

func VerifikasiAkun(pegawai models.Pegawai, action string, usrgroup string) error {
	switch action {
	case "approve":
		switch usrgroup {
		case models.UKPBJ:
			pegawai.PegStatus = models.APPROVED_UKPBJ
		case models.ADMIN:
			pegawai.PegStatus = models.APPROVED
		}
		pegawai.TglApprove = time.Now()
		pegawai.PegIsactive = 1
	case "reject":
		pegawai.PegStatus = models.REJECT
		pegawai.TglReject = time.Now()
		pegawai.PegIsactive = 0
	}
	return models.SavePegawai(&pegawai)
}

func GetPerubahanData(id uint) models.PerubahanData {
	return models.GetPerubahanData(id)
}

func CreatePerubahanData(c *fiber.Ctx, perubahan *models.PerubahanData, name string) error {
	perubahan.DokId , _ = models.SaveDocument(c, perubahan.PegId, models.PERUBAHAN_DATA, name)
	perubahan.Status = 0;
	return models.SavePerubahanData(perubahan)
}

func UpdatePerubahanData(perubahan models.PerubahanData) error {
	err := models.SavePerubahanData(&perubahan)
	if err == nil {
		var email models.Inbox
		if(perubahan.IsProses()){
			email = models.Inbox{
				Subject: "Perubahan Data Akun",
				PegId: perubahan.PegId,
				Content: `Permohonan Perubahan data akun sedang di proses oleh UKPBJ`,
				Status: "inbox",
				EnqueueDate: time.Now(),
			}
		} else if(perubahan.IsSelesai()){
			email = models.Inbox{
				Subject: "Perubahan Data Akun",
				PegId: perubahan.PegId,
				Content: `Permohonan Perubahan data akun sudah dilakukan. Harap segera check di menu profile`,
				Status: "inbox",
				EnqueueDate: time.Now(),
			}
		}
		err = models.SaveInbox(&email)
		// TODO send email
	}
	return err
}

func DeletePerubahanData(id uint) error {
	perubahan := GetPerubahanData(id)
	if perubahan.ID == 0 {
		return errors.New("Perubahan tidak ditemukan")
	}
	err := models.DeletePerubahanData(&perubahan)
	if err != nil {
		log.Error(err)
		return errors.New("Hapus Perubahan Gagal")
	}
	return nil
}

func AutoCreateAdminIfNoExist() {
	models.AutoCreateAdmin()
}

func GetInbox(id uint) models.Inbox {
	return models.GetInbox(id)
}

func SendNotifikasi(user *models.Pegawai, plainPassw string) error {
	content := `Proses Registrasi Akun berhasil dilakukan. Informasi akun : <br />
 					User ID : %s <br/>
				    Password : %s <br />
				Silahkan login dengan menggunakan Akun.`
	email := models.Inbox{
		Subject: "Notifikasi Registrasi Akun",
		PegId: user.ID,
		Content: fmt.Sprintf(content, user.PegNamauser, plainPassw),
		Status: "inbox",
		EnqueueDate: time.Now(),
	}
	return models.SaveInbox(&email)
}

func Registrasi(c *fiber.Ctx, pegawai *models.Pegawai, plainPassw string) error {
	pegawai.PegNamauser = strings.ToUpper(pegawai.PegNamauser)
	pegawai.Passw = utils.HashPassword(pegawai.Passw)
	err := models.SavePegawai(pegawai)
	if err != nil {
		log.Error(err)
		return errors.New("Registrasi Akun Gagal")
	}
	_, err = models.SaveDocument(c, pegawai.ID, models.KTP, "ktp")
	_, err =models.SaveDocument(c, pegawai.ID, models.SK, "sk")
	_, err =models.SaveDocument(c, pegawai.ID, models.SERTIFIKAT, "sertifikat")
	err = SaveTTD(c, pegawai.ID)
	if err != nil {
		log.Error(err)
		return errors.New("Registrasi Akun Gagal")
	}
	err = SendNotifikasi(pegawai, plainPassw)
	if err != nil {
		log.Error(err)
		return errors.New("Registrasi Akun Gagal")
	}
	return nil
}
