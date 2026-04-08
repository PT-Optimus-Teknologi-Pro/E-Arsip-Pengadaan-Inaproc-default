package models

import (
	"log/slog"

	"golang.org/x/exp/maps"
)

type Program struct {
	IsDeleted		bool 		`json:"is_deleted"`
    JenisKlpd		string 		`json:"jenis_klpd"`
    KdKlpd			string 		`json:"kd_klpd"`
    KdProgram		uint 		`gorm:"primaryKey;autoIncrement:false" json:"kd_program"`
    KdProgramLokal	uint		`json:"kd_program_lokal"`
    KdProgramStr	string 		`json:"kd_program_str"`
    KdSatker		uint 		`json:"kd_satker"`
    NamaKlpd		string 		`json:"nama_klpd"`
    NamaProgram		string 		`json:"nama_program"`
    PaguProgram		float64		`json:"pagu_program"`
    TahunAnggaran	int 		`json:"tahun_anggaran"`
}

func (c Program) TableName() string {
	return "program"
}

func SaveProgram(datas *map[uint]Program, tahun int) {
	if len(*datas) > 0 {
		values := maps.Values(*datas)
		DB.Unscoped().Where("tahun_anggaran=?", tahun).Delete(&Program{})
		err := DB.CreateInBatches(values, MAX_DATA).Error
		if err != nil {
			slog.Error(err.Error())
		}
	}
}
