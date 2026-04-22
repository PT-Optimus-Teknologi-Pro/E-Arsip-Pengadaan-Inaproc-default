package models

import (
	"database/sql"
	"gorm.io/gorm"
)

var BidangList = []string{
	"Anggaran",
	"Rencana Umum Pengadaan (RUP)",
	"Spesifikasi",
	"Gambar(Barang)/ Kerja (Konstruksi)",
	"KAK",
	"Metode pengadaan",
	"HPS",
	"Rancangan Kontrak",
	"Analisa Pasar",
	"jaminan penawaran untuk Jasa konstruksi dengan nilai total HPS >10M",
	"Sertifikat garansi/kartu-jaminan/garansi purna jual untuk pengadaan barang",
	"Sertifikat/dokumen dalam rangka pengadaan barang impor untuk pengadaan barang",
}

type ReviuBidang struct {
	gorm.Model
	Nama string `gorm:"uniqueIndex" json:"nama" form:"nama"`
}

func (ReviuBidang) TableName() string {
	return "reviu_bidang"
}

func GetReviuBidang(id uint) ReviuBidang {
	var res ReviuBidang
	db.First(&res, id)
	return res
}

func GetAllReviuBidang() []ReviuBidang {
	var results []ReviuBidang
	db.Order("nama asc").Find(&results)
	return results
}

func SaveReviuBidang(obj *ReviuBidang) error {
	return db.Save(obj).Error
}

func DeleteReviuBidang(obj *ReviuBidang) error {
	return db.Delete(obj).Error
}

type BeritaAcara struct {
	gorm.Model
	PktId            uint         `gorm:"pkt_id" json:"pkt_id"` // Link to Paket
	Nomor            string       `gorm:"nomor"`
	Jenis            string       `form:"jenis"`
	Hari             string       `form:"hari"`
	Tanggal          sql.NullTime `json:"tanggal"`
	Tempat           string       `form:"tempat"`
	Waktu            string       `form:"waktu"`
	SubKeg           string       `form:"sub_keg"`
	Pengadaan        string       `form:"pengadaan"`
	Uraian           string       `form:"uraian"`
	DokId            uint         `form:"dok_id"`
	DasarPelaksanaan string       `form:"dasar_pelaksanaan" gorm:"type:text"` // Isi BA - Dasar
	Pembahasan       string       `form:"pembahasan" gorm:"type:text"`        // Isi BA - Pembahasan
	Kesimpulan       string       `form:"kesimpulan" gorm:"type:text"`        // Isi BA - Kesimpulan
}

func (BeritaAcara) TableName() string {
	return "berita_acara"
}

func GetBeritaAcara(id uint) BeritaAcara {
	var result BeritaAcara
	db.First(&result, id)
	return result
}

type Reviu struct {
	gorm.Model
	Bidang 		string		`json:"bidang"`
	Content		string 		`json:"content"`
	Opsi1		string 		`json:"opsi1"`
	Opsi2		string 		`json:"opsi2"`
}

func (Reviu) TableName() string {
	return "reviu"
}

func GetReviu(id uint) Reviu {
	var reviu Reviu
	db.First(&reviu, id)
	return reviu
}

func SaveReviu(reviu *Reviu) error {
	return db.Save(reviu).Error
}

func DeleteReviu(reviu *Reviu) error {
	return db.Delete(reviu).Error
}

func GetAllReviu() []Reviu {
	var results []Reviu
	db.Find(&results)
	return results
}

type ReviuPaket struct {
	gorm.Model
	PktId			uint 		`json:"pkt_id"`
	RevId			uint		`json:"rev_id"`
	Status          int         `json:"status"` // 0: Kosong, 1: Sesuai/Tersedia, 2: Tidak Sesuai/Tidak Tersedia
	Keterangan		string		`json:"Keteranga"`
	CatatanKhusus	string		`json:"catatan_khusus"`
	PegId			uint 		`json:"peg_id"`
}

func (ReviuPaket) TableName() string {
	return "reviu_paket"
}

func GetTemplateByVariable(variable string) Templates {
	var t Templates
	db.Where("variable = ?", variable).First(&t)
	return t
}

func AutoCreateReviuMaster() {
	// 0. Seed BA Reviu default templates
	defaultTemplates := []Templates{
		{
			Nama:     "Dasar Pelaksanaan BA Reviu",
			Variable: "reviu_dasar",
			Content:  "<p>Berdasarkan Peraturan Presiden Nomor 16 Tahun 2018 tentang Pengadaan Barang/Jasa Pemerintah beserta perubahannya, dan Peraturan LKPP Nomor 12 Tahun 2021 tentang Pedoman Pelaksanaan Pengadaan Barang/Jasa Pemerintah melalui Penyedia, maka dilaksanakan Reviu Dokumen Persiapan Pengadaan.</p>",
		},
		{
			Nama:     "Pembahasan BA Reviu",
			Variable: "reviu_pembahasan",
			Content:  "<p>Pokja Pemilihan/Pejabat Pengadaan bersama Pejabat Pembuat Komitmen (PPK) telah melakukan reviu terhadap dokumen persiapan pengadaan yang meliputi: Dokumen Anggaran, Rencana Umum Pengadaan (RUP), Spesifikasi Teknis/KAK, HPS, dan Rancangan Kontrak.</p>",
		},
		{
			Nama:     "Kesimpulan BA Reviu",
			Variable: "reviu_kesimpulan",
			Content:  "<p>Berdasarkan hasil reviu terhadap dokumen persiapan pengadaan sebagaimana tersebut di atas, maka dokumen persiapan pengadaan dinyatakan <strong>[SESUAI / BELUM SESUAI]</strong> dan proses pengadaan dapat <strong>[DILANJUTKAN / DIPERBAIKI TERLEBIH DAHULU]</strong>.</p>",
		},
	}
	for _, tpl := range defaultTemplates {
		var count int64
		db.Model(&Templates{}).Where("variable = ?", tpl.Variable).Count(&count)
		if count == 0 {
			db.Create(&tpl)
		}
	}

	// 1. Seed Bidang/Kategori
	var countBidang int64
	db.Model(&ReviuBidang{}).Count(&countBidang)
	if countBidang == 0 {
		for _, name := range BidangList {
			db.Create(&ReviuBidang{Nama: name})
		}
	}

	// 2. Seed Reviu Items
	var count int64
	db.Model(&Reviu{}).Count(&count)
	if count > 0 {
		return
	}

	items := []Reviu{
		// 1. Anggaran
		{Bidang: "Anggaran", Content: "Dokumen Anggaran Belanja (DPA atau RKA yang telah ditetapkan)", Opsi1: "Tersedia", Opsi2: "Tidak Tersedia"},
		{Bidang: "Anggaran", Content: "Apakah Mata Anggaran Kegiatan (MAK) sudah sesuai?", Opsi1: "Sesuai", Opsi2: "Tidak Sesuai"},
		{Bidang: "Anggaran", Content: "Kode akun yang tercantum di Dokumen Pelaksanaan Anggaran (DPA) telah sesuai dengan jenis peruntukan pengeluarannya", Opsi1: "Sesuai", Opsi2: "Tidak Sesuai"},
		{Bidang: "Anggaran", Content: "Perkiraan jumlah anggaran yang tersedia untuk paket pekerjaan dalam dokumen anggaran mencukupi kebutuhan pelaksanaan pekerjaan", Opsi1: "Tersedia", Opsi2: "Tidak Tersedia"},

		// 2. RUP
		{Bidang: "Rencana Umum Pengadaan (RUP)", Content: "Pokja Pemilihan mengingatkan kepada PA, dalam melakukan pemaketan pengadaan barang mengacu pada Peraturan Presiden Nomor 16 tahun 2018 beserta perubahan dan turunannnya", Opsi1: "Sesuai", Opsi2: "Tidak Sesuai"},
		{Bidang: "Rencana Umum Pengadaan (RUP)", Content: "Pengguna Anggaran telah mengumumkan Rencana Umum Pengadaan", Opsi1: "Sesuai", Opsi2: "Tidak Sesuai"},
		{Bidang: "Rencana Umum Pengadaan (RUP)", Content: "PPK dalam merencanakan pemilihan penyedia pada SiRUP agar memperhatikan pasal 38 tentang metode pemilihan penyedia barang/pekerjaan konstruksi/jasa_lainnya terdiri atas E-Purchasing, Pengadaan Langsung, Penunjukan Langsung, Tender Cepat dan Tender", Opsi1: "Sesuai", Opsi2: "Tidak Sesuai"},
		{Bidang: "Rencana Umum Pengadaan (RUP)", Content: "Rencana waktu penggunaan barang/jasa", Opsi1: "Sesuai", Opsi2: "Tidak Sesuai"},

		// 3. Spesifikasi
		{Bidang: "Spesifikasi", Content: "Untuk Persyaratan spesifikasi teknis yang telah dibuat oleh PPK, agar sesuai Peraturan Presiden Nomor 12 tahun 2021 pasal 19 ayat (2)", Opsi1: "Sesuai", Opsi2: "Tidak Sesuai"},
		{Bidang: "Spesifikasi", Content: "Dalam penyusunan spesifikasi teknis/KAK dimungkinkan penyebutan merek terhadap: a. Komponen barang/jasa; b. suku cadang; c. bagian dari satu sistem yang sudah ada; atau d. barang/jasa dalam katalog elektronik atau Toko Daring.", Opsi1: "Sesuai", Opsi2: "Tidak Sesuai"},
		{Bidang: "Spesifikasi", Content: "Untuk Spesifikasi Tender Cepat dilaksanakan dalam hal : 1. Spesifikasi dan volume pekerjaannya sudah dapat ditentukan secara rinci; dan 2. Pelaku Usaha telah terkualifikasi dalam Sistem Informasi Kinerja Penyedia. 3. Penyusunan spesifikasi teknis dapat menyebutkan merek terhadap: 1) suku cadang; atau 2) bagian dari satu sistem yang sudah ada.", Opsi1: "Sesuai", Opsi2: "Tidak Sesuai"},
		{Bidang: "Spesifikasi", Content: "Dalam menyusun Spesifikasi Bahan konstruksi ketersediaan barang/jasa perlu memperhatikan: 1) Tingkat Komponen Dalam Negeri (TKDN) yang mengacu pada daftar inventarisasi barang/jasa produksi dalam negeri (tkdn.kemenperin.go.id); 2) memenuhi Standar Nasional Indonesia(SNI); 3) produk usaha mikro dan kecil serta koperasi dari hasil produksi dalam negeri; dan 4) produk ramah lingkungan hidup", Opsi1: "Sesuai", Opsi2: "Tidak Sesuai"},

		// 4. Gambar
		{Bidang: "Gambar(Barang)/ Kerja (Konstruksi)", Content: "Pengadaan Barang/ jasa Lainnya: Gambar apakah sudah sesuai spesifikasi teknis pekerjaan", Opsi1: "Sesuai", Opsi2: "Tidak Sesuai"},
		{Bidang: "Gambar(Barang)/ Kerja (Konstruksi)", Content: "Pekerjaan Konstruksi : Gambar-gambar untuk pelaksanaan pekerjaan harus ditetapkan oleh Pejabat Pembuat Komitmen (PPK) secara terinci, lengkap dan jelas, antara lain : 1. Peta Lokasi 2. Lay out 3. Potongan memanjang 4. Potongan melintang 5. Detail-detail konstruksi", Opsi1: "Sesuai", Opsi2: "Tidak Sesuai"},
		{Bidang: "Gambar(Barang)/ Kerja (Konstruksi)", Content: "Apakah Gambar Kerja sudah sesuai dengan Spesifikasi serta volume RAB", Opsi1: "Sesuai", Opsi2: "Tidak Sesuai"},

		// 5. KAK
		{Bidang: "KAK", Content: "Apakah KAK sudah menjelaskan latar belakang kegiatan/pengadaan yang akan dilaksanakan", Opsi1: "Sesuai", Opsi2: "Tidak Sesuai"},
		{Bidang: "KAK", Content: "Apakah KAK sudah menjelaskan maksud dan tujuan kegiatan/pengadaan yang akan dilaksanakan", Opsi1: "Sesuai", Opsi2: "Tidak Sesuai"},
		{Bidang: "KAK", Content: "Fasilitas peralatan", Opsi1: "Sesuai", Opsi2: "Tidak Sesuai"},

		// 6. Metode pengadaan
		{Bidang: "Metode pengadaan", Content: "Memenuhi kebutuhan pengguna akhir", Opsi1: "Sesuai", Opsi2: "Tidak Sesuai"},

		// 7. HPS
		{Bidang: "HPS", Content: "HPS benar-benar didasarkan pada Harga Pasar setempat. HPS yang telah disusun perlu dicermati agar sesuai dengan ketentuan dan kebutuhan barang. Disamping itu dalam pembuatan HPS agar memperbandingkan harga pasar harus lebih 1 (Satu) sumber yang berbeda Berdasarkan pasal 2 6 Peraturan Presiden Nomor 16 tahun 2018 beserta perubahan dan turunannnya.", Opsi1: "Sesuai", Opsi2: "Tidak Sesuai"},
		{Bidang: "HPS", Content: "Apakah sesuai dengan spesifikasi teknis/KAK dan ruang lingkup pekerjaan?", Opsi1: "Sesuai", Opsi2: "Tidak Sesuai"},
		{Bidang: "HPS", Content: "Apakah sudah memperhitungkan kewajiban (perpajakan/cukai/asuransi/SMK3 atau biaya lain yang dipersyaratkan dalam pelaksanaan pekerjaan", Opsi1: "Sesuai", Opsi2: "Tidak Sesuai"},
		{Bidang: "HPS", Content: "Biaya penyelenggaraan K3", Opsi1: "Sesuai", Opsi2: "Tidak Sesuai"},
		{Bidang: "HPS", Content: "Apakah sudah membuat perhitungan Nilai TKDN", Opsi1: "Sesuai", Opsi2: "Tidak Sesuai"},

		// 8. Rancangan Kontrak
		{Bidang: "Rancangan Kontrak", Content: "Apakah sesuai dengan ruang lingkup pekerjaan", Opsi1: "Sesuai", Opsi2: "Tidak Sesuai"},
		{Bidang: "Rancangan Kontrak", Content: "Syarat-Syarat Umum Kontrak (SSUK)", Opsi1: "Sesuai", Opsi2: "Tidak Sesuai"},
		{Bidang: "Rancangan Kontrak", Content: "Syarat-Syarat Khusus Kontrak (SSKK)", Opsi1: "Sesuai", Opsi2: "Tidak Sesuai"},
		{Bidang: "Rancangan Kontrak", Content: "Surat Perjanjian", Opsi1: "Sesuai", Opsi2: "Tidak Sesuai"},

		// 9. Analisa Pasar
		{Bidang: "Analisa Pasar", Content: "Ketersediaan barang/jasa di pasar", Opsi1: "Tersedia", Opsi2: "Tidak Tersedia"},
		{Bidang: "Analisa Pasar", Content: "Ketersediaan penyedia barang / jasa dalam negeri", Opsi1: "Tersedia", Opsi2: "Tidak Tersedia"},

		// 10. Jaminan Penawaran
		{Bidang: "jaminan penawaran untuk Jasa konstruksi dengan nilai total HPS >10M", Content: "Memakai jaminan penawaran 1% hingga 3% dari nilai total HPS", Opsi1: "Sesuai", Opsi2: "Tidak Sesuai"},

		// 11. Sertifikat Garansi
		{Bidang: "Sertifikat garansi/kartu-jaminan/garansi purna jual untuk pengadaan barang", Content: "Sesuai dengan KAK", Opsi1: "Sesuai", Opsi2: "Tidak Sesuai"},

		// 12. Sertifikat Impor
		{Bidang: "Sertifikat/dokumen dalam rangka pengadaan barang impor untuk pengadaan barang", Content: "Apakah sudah mencakup : a. Supporting Letter/Letter of Intent/Letter of Agreement dari pabrikan/prinsipal di negara asal; b. Surat Keterangan Asal (Certificate of Origin); dan c. Sertifikat Produksi.", Opsi1: "Sesuai", Opsi2: "Tidak Sesuai"},
	}

	for _, item := range items {
		db.Create(&item)
	}
}
