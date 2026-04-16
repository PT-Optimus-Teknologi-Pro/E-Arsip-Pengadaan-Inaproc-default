package services

const (
	DOMAIN_URL = "https://data.inaproc.id"
	RUP_SATKER = "/api/v1/rup/master-satker"
	RUP_PROGRAM = "/api/v1/rup/program-master"
	RUP_PAKET_ANGGARAN_PENYEDIA = "/api/v1/rup/paket-anggaran-penyedia"
	RUP_PAKET_ANGGARAN_SWAKELOLA = "/api/v1/rup/paket-anggaran-swakelola"
	RUP_PAKET_PENYEDIA_TERUMUMKAN = "/api/v1/rup/paket-penyedia-terumumkan"
	RUP_PAKET_SWAKELOLA_TERUMUMKAN = "/api/v1/rup/paket-swakelola-terumumkan"
	TENDER_PENGUMUMAN = "/api/v1/tender/pengumuman"
	TENDER_JADWAL = "/api/v1/tender/jadwal-tahapan-tender"
	TENDER_PESERTA = "/api/v1/tender/peserta-tender"
	TENDER_KONTRAK = "/api/v1/tender/tender-ekontrak-kontrak"
	TENDER_SELESAI_NILAI = "/api/v1/tender/tender-selesai-nilai"
	NONTENDER_JADWAL = "/api/v1/tender/jadwal-tahapan-non-tender"
	NONTENDER_PENGUMUMAN= "/api/v1/tender/non-tender-pengumuman"
	NONTENDER_SELESAI = "/api/v1/tender/non-tender-selesai"
	NONTENDER_KONTRAK ="/api/v1/tender/non-tender-ekontrak-kontrak"
	PENCATATAN = "/api/v1/tender/pencatatan-non-tender"
	PENCATATAN_REALISASI = "/api/v1/tender/pencatatan-non-tender-realisasi"
	SWAKELOLA ="/api/v1/tender/pencatatan-swakelola"
	SWAKELOLA_REALISASI = "/api/v1/tender/pencatatan-swakelola-realisasi"
	KATALOG_PURCHASING = "/api/v1/ekatalog/paket-e-purchasing" // v6 katalog
	KATALOG_PURCHASING_V5 ="/api/v1/ekatalog-archive/paket-e-purchasing" // v5 katalog
	KATALOG_INSTANSI_SATKER = "/api/v1/ekatalog-archive/instansi-satker"
	KATALOG_KOMODITAS = "/api/v1/ekatalog-archive/komoditas-detail"
	KATALOG_DISRIBUTOR = "/api/v1/ekatalog-archive/penyedia-distributor-detail"
	KATALOG_PENYEDIA_v5 = "/api/v1/ekatalog-archive/penyedia-detail"
	KATALOG_PENYEDIA = "/api/v1/ekatalog/penyedia-detail"
)

func Sync() {
	// syncSatker()
	// syncProgram()
	syncRupAnggaran()
	syncRup()
	syncRupAnggaranSwakelola()
	syncRupSwakelola()
	syncKatalog()
	syncPenyediaKatalog()
	syncKatalogArchive()
	syncPenyediaArchive()
	syncNontender()
	syncNontenderSelesai()
	syncJadwalNontender()
	syncPencatatan()
	syncPencatatanRealisasi()
	syncSwakelola()
	syncSwakelolaRealisasi()
	syncTender()
	syncTenderSelesai()
	syncJadwal()
	syncPeserta()
	syncKontrak()
	syncKontrakNontender()
}
