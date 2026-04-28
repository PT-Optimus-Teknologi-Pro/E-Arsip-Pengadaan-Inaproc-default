package models

import (
	"arsip/utils"
	"encoding/json"
)

// definsikan model2 terkait dashboard
type RekapSatker struct {
	KdSatker 	uint
	NamaSatker	string
}

func GetRekapSatker(tahun int) []RekapSatker {
	var result []RekapSatker
	db.Raw("SELECT DISTINCT kd_satker, nama_satker FROM rup WHERE tahun_anggaran=?", tahun).Scan(&result)
	return result
}

type RupProgress struct {
	KdSatker		uint
	Nama			string
	Pagu			float64
	PaguPds			float64
	PaguSwakelola	float64
	Belanja 		float64
	TotalBelanja	float64
}

func (c RupProgress) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]interface{}{
		"kd_satker" : c.KdSatker,
		"nama_satker" : c.Nama,
		"pagu" : utils.FormatRupiah(c.Pagu),
		"pagu_pds" : utils.FormatRupiah(c.PaguPds),
		"pagu_swakelola" : utils.FormatRupiah(c.PaguSwakelola),
		"belanja" : utils.FormatRupiah(c.Belanja),
		"total_belanja" : utils.FormatRupiah(c.TotalBelanja),
		"prosentase" : utils.Prosentase(c.Pagu+c.PaguPds+c.PaguSwakelola, c.TotalBelanja),
	})
}

type PaketPrioritas struct {
	ID 		uint
	Nama 	string
	PegNama	string
	Status	int
	Pagu 	float64
	Metode 	int
}
func (obj PaketPrioritas) StatusLabel() string {
	return statusPaket[obj.Status]
}
func (c PaketPrioritas) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]interface{}{
		"id" : c.ID,
		"nama_ppk" : c.PegNama,
		"nama_paket" : c.Nama,
		"status" : c.StatusLabel(),
		"pagu" : c.Pagu,
		"metode" : metodePengadaan[c.Metode],
	})
}

type BebanPersonil struct {
	ID 			uint
	PegNama		string
	PaketPokja	int
	PaketPp		int
	PaguPokja	float64
	PaguPp		float64
}
func (c BebanPersonil) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]interface{}{
		"nama_pegawai" : c.PegNama,
		"jumlah" : c.Jumlah(),
		"pagu" : utils.FormatRupiah(c.Pagu()),
	})
}

func (obj BebanPersonil) Jumlah() int {
	return obj.PaketPokja + obj.PaketPp
}

func (obj BebanPersonil) Pagu() float64 {
	return obj.PaguPokja + obj.PaguPp
}

type RekapPaketSatker struct {
	KdSatker		uint   `gorm:"column:kd_satker" json:"kd_satker"`
	NamaSatker		string `gorm:"column:nama_satker" json:"nama_satker"`
	Tender			int    `gorm:"column:tender" json:"tender"`
	Nontender		int    `gorm:"column:nontender" json:"nontender"`
	Pencatatan		int    `gorm:"column:pencatatan" json:"pencatatan"`
	Katalog			int    `gorm:"column:katalog" json:"katalog"`
}

func (c RekapPaketSatker) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]interface{}{
		"nama":       c.NamaSatker,
		"tender":     c.Tender,
		"nontender":  c.Nontender,
		"pencatatan": c.Pencatatan,
		"katalog":    c.Katalog,
	})
}

type RekapFeedback struct {
	ID		int
	Jan		int
	Feb 	int
	Mar 	int
	Apr 	int
	Mei 	int
	Jun		int
	Jul 	int
	Agu		int
	Sep		int
	Okt		int
	Nov		int
	Des		int
}

func (obj RekapFeedback) Label() string {
	switch obj.ID {
	case 1:
		return "Tidak Puas"
	case 2:
		return "Puas"
	case 3:
		return "Sangat Puas"
	}
	return ""
}

type RekapFeedbackDashboard struct {
	Tahun 			int
	TidakPuas		[]int
	Puas			[]int
	SangatPuas		[]int
}

type RekapFeedbackBulanan struct {
	Count		int
	Bulan 		int
	Kepuasan	int
}

func GetRekapFeedbackBulanan(tahun int) []RekapFeedbackBulanan {
	var rekap []RekapFeedbackBulanan
	db.Raw(`SELECT count(id), jenis, EXTRACT(MONTH FROM created_at) bulan FROM feedback
			WHERE EXTRACT(YEAR FROM created_at) = ? GROUP BY jenis, EXTRACT(MONTH FROM created_at)`, tahun).Scan(&rekap)
	return rekap
}

type RekapSatkerDashboard struct {
	KdSatker		uint 		`json:"kd_satker"`
	NamaSatker		string		`json:"nama_satker"`
	Tahun			int			`json:"tahun"`
	PaketRencana	[]int		`json:"paket_rencana"`
	Rencana 		[]float64	`json:"rencana"`
	PaketRealisasi	[]int		`json:"paket_realisasi"`
	Realisasi 		[]float64	`json:"realisasi"`
	PaketPurchase	[]int		`json:"paket_purchase"`
	Purchase 		[]float64	`json:"purchase"`
}

func (c RekapSatkerDashboard) TotalPaketRealisasi() int {
	result := 0
    for i := 0; i < len(c.PaketRealisasi); i++ {
        result += c.PaketRealisasi[i]
    }
    return result
}

func (c RekapSatkerDashboard) TotalPaketPurchase() int {
	result := 0
    for i := 0; i < len(c.PaketPurchase); i++ {
        result += c.PaketPurchase[i]
    }
    return result
}

// func (c RekapSatkerDashboard) MarshalJSON() ([]byte, error) {
// 	return json.Marshal(map[string]interface{}{
// 		"nama_satker" : c.NamaSatker,
// 		"realisasi" : c.TotalPaketRealisasi(),
// 		"purchase" : c.TotalPaketPurchase(),
// 	})
// }

type RekapRencanaSatkerBulanan struct {
	KdSatker	uint
	Bulan 		int
	Paket 		int
	Pagu  		float64
}

type RekapRealisasiSatkerBulanan struct {
	KdSatker	uint
	Bulan 		int
	Paket 		int
	Realisasi	float64
}

type DetailPaketBulanan struct {
	Kode 			uint
	NamaPaket		string
	TglRealisasi	Datetime
	Realisasi		float64
}


func GetRupProgress(tahun int) []RupProgress {
	var result []RupProgress
	db.Raw(`SELECT a.id_satker, s.nama, belanja_pengadaan as belanja, total_belanja, 
			p.pagu, pds.pagu as pagu_pds, w.pagu as pagu_swakelola FROM struktur_anggaran a
			LEFT JOIN (SELECT id_satker, sum(pagu) as pagu FROM paket_sirup WHERE tahun = ? and paket_aktif='TRUE' and paket_terumumkan='TRUE' and (id_swakelola = 0 OR id_swakelola IS NULL) GROUP BY id_satker) p On a.id_satker = p.id_satker
			LEFT JOIN (SELECT id_satker, sum(pagu) as pagu FROM paket_sirup WHERE tahun = ? and paket_aktif='TRUE' and paket_terumumkan='TRUE' and id_swakelola > 0 GROUP BY id_satker) pds On a.id_satker = pds.id_satker
			LEFT JOIN (SELECT id_satker, sum(jumlah_pagu) as pagu FROM swakelola_sirup WHERE tahun = ? and aktif='TRUE' and umumkan='TRUE' GROUP BY id_satker) w On a.id_satker = w.id_satker
			LEFT JOIN satker s ON a.id_satker = s.id WHERE a.tahun_anggaran = ?`, tahun, tahun, tahun, tahun).Scan(&result)
	return result
}

func GetPaketPrioritas(tahun int) []PaketPrioritas {
	var res []PaketPrioritas
	db.Raw(`SELECT p.id, p.nama, p.status, p.pagu, p.metode, peg.peg_nama FROM paket p, pegawai peg, paket_sirup s WHERE p.rup_id=s.id AND p.ppk_id = peg.id AND prioritas = ? AND s.tahun=?`, true, tahun).Scan(&res)
	return res
}

func GetBebanPersonil(tahun int) []BebanPersonil {
	var res []BebanPersonil
	db.Raw(`SELECT id, peg_nama, pokja.jml paket_pokja, pp.jml paket_pp, pokja.pagu pagu_pokja, pp.pagu pagu_pp FROM pegawai p
LEFT JOIN (SELECT peg_id, p.pnt_id, count(p.id) jml, sum(pagu) pagu FROM panitia pa, anggota_panitia ap, paket p
WHERE pa.id = ap.pnt_id AND ap.pnt_id = p.pnt_id AND tahun = ?  GROUP BY p.pnt_id, peg_id) pokja ON pokja.peg_id=p.id
LEFT JOIN (SELECT pp_id, count(id) jml, sum(pagu) pagu FROM paket WHERE pp_id > 0 AND EXTRACT(YEAR FROM tgl_disetujui)= ? GROUP BY pp_id) pp ON pp.pp_id=p.id
WHERE usrgroup in ('POKJA', 'PP')`, tahun, tahun).Scan(&res)
	return res
}

func GetRekapPaketSatker(tahun int) []RekapPaketSatker {
	var res []RekapPaketSatker
	db.Raw(`SELECT 
		kd_satker::bigint AS kd_satker, 
		MAX(nama_satker) AS nama_satker, 
		SUM(tender)::int AS tender, 
		SUM(nontender)::int AS nontender, 
		SUM(pencatatan)::int AS pencatatan, 
		SUM(katalog)::int AS katalog 
	FROM (
		SELECT kd_satker, nama_satker, count(*) AS tender, 0 AS nontender, 0 AS pencatatan, 0 AS katalog FROM tender WHERE tahun_anggaran=? GROUP BY kd_satker, nama_satker
		UNION ALL
		SELECT kd_satker, nama_satker, 0 AS tender, count(*) AS nontender, 0 AS pencatatan, 0 AS katalog FROM nontender WHERE tahun_anggaran=? GROUP BY kd_satker, nama_satker
		UNION ALL
		SELECT kd_satker, nama_satker, 0 AS tender, 0 AS nontender, count(*) AS pencatatan, 0 AS katalog FROM pencatatan WHERE tahun_anggaran=? GROUP BY kd_satker, nama_satker
		UNION ALL
		SELECT satker_id::text AS kd_satker, nama_satker, 0 AS tender, 0 AS nontender, 0 AS pencatatan, count(*) AS katalog FROM katalog_archive WHERE tahun_anggaran=? GROUP BY satker_id, nama_satker
	) combined
	WHERE kd_satker IS NOT NULL
	GROUP BY kd_satker
	ORDER BY nama_satker`, tahun, tahun, tahun, tahun).Scan(&res)
	return res
}

func GetRekapSatkerDashboard(tahun int) []RekapSatker {
	var result []RekapSatker
	db.Raw(`SELECT DISTINCT kd_satker, nama_satker  FROM rup WHERE tahun_anggaran = ?`, tahun).Scan(&result)
	return result
}

func GetRekapRencanaSatkerBulanan(tahun int) []RekapRencanaSatkerBulanan {
	var rekap []RekapRencanaSatkerBulanan
	db.Raw(`SELECT kd_satker, count(kd_rup) paket, sum(pagu) pagu,  EXTRACT(MONTH FROM tgl_awal_pemilihan) as bulan
    		FROM rup WHERE tahun_anggaran = ? group by  kd_satker, EXTRACT(MONTH FROM tgl_awal_pemilihan)
      		ORDER BY  kd_satker, EXTRACT(MONTH FROM tgl_awal_pemilihan)`, tahun).Scan(&rekap)
	return rekap
}

func GetRencanaPaketSatkerBulanan(tahun int, satkerId uint, bulan int) []Rup {
	var rekap []Rup
	db.Find(&rekap, `tahun_anggaran=? AND kd_satker=? AND  EXTRACT(MONTH FROM tgl_awal_pemilihan) = ?`, tahun, satkerId, bulan)
	return rekap
}

func GetRealisasiPaketTenderBulanan(tahun int, satkerId uint, bulan int) []DetailPaketBulanan {
	var rekap []DetailPaketBulanan
	db.Raw(`SELECT p.kd_tender kode, p.nama_paket, t.tgl_awal tgl_realisasi, p.realisasi -> 0 -> 'nilai_kontrak' realisasi
			FROM tender p LEFT JOIN rup rup ON rup.kd_rup::text = ANY(string_to_array(p.kd_rup, ';')) LEFT JOIN jadwal t ON p.kd_tender=t.kd_tender
			WHERE p.tahun_anggaran = ? AND p.realisasi IS NOT NULL AND t.kd_tahapan = 18803 AND EXTRACT(MONTH FROM t.tgl_awal) = ? and rup.kd_satker=? `,
			tahun, bulan, satkerId).Scan(&rekap)
	return rekap
}

func GetRealisasiPaketNontenderBulanan(tahun int, satkerId uint, bulan int) []DetailPaketBulanan {
	var rekap []DetailPaketBulanan
	db.Raw(`SELECT p.kd_nontender kode, p.nama_paket, t.tgl_awal tgl_realisasi, p.realisasi -> 0 -> 'nilai_kontrak' realisasi
			FROM nontender p LEFT JOIN rup rup ON rup.kd_rup::text = ANY(string_to_array(p.kd_rup, ';')) LEFT JOIN jadwal_nontender t ON p.kd_nontender=t.kd_nontender
			WHERE p.tahun_anggaran = ? AND p.realisasi IS NOT NULL AND t.kd_tahapan = 18841 and EXTRACT(MONTH FROM t.tgl_awal) = ? and  rup.kd_satker=?`,
			tahun, bulan, satkerId).Scan(&rekap)
	return rekap
}

func GetRealisasiPencatatanBulanan(tahun int, satkerId uint, bulan int) []DetailPaketBulanan {
	var rekap []DetailPaketBulanan
	db.Raw(`SELECT kd_nontender_pct kode, p.nama_paket, tgl_selesai_paket tgl_realisasi, total_realisasi realisasi
			FROM pencatatan p LEFT JOIN rup rup ON rup.kd_rup::text = ANY(string_to_array(p.kd_rup, ';'))
			WHERE p.tahun_anggaran = ? AND p.realisasi IS NOT NULL AND EXTRACT(MONTH FROM tgl_mulai_paket) = ?  and  rup.kd_satker=?`,
			tahun, bulan, satkerId).Scan(&rekap)
	return rekap
}

func GetRealisasiPurchaseBulanan(tahun int, satkerId uint, bulan int) []DetailPaketBulanan {
	var rekap []DetailPaketBulanan
	db.Raw(`SELECT kd_paket kode, p.nama_paket, tanggal_edit_paket tgl_realisasi, total_harga as realisasi FROM katalog_archive p
			WHERE tahun_anggaran = ? and status_paket IN ('proses_kontrak_ppk', 'melakukan_pengiriman_dan_penerimaan') AND EXTRACT(MONTH FROM tanggal_edit_paket) = ? AND p.satker_id=?`,
			tahun, bulan, satkerId).Scan(&rekap)
	return rekap
}

func GetRealisasiNontenderPerSkpd(tahun int) []RekapRealisasiSatkerBulanan {
	var rekap []RekapRealisasiSatkerBulanan
	db.Raw(`SELECT rup.kd_satker, count(p.kd_nontender) paket, sum((p.realisasi -> 0 -> 'nilai_kontrak')::numeric) realisasi, EXTRACT(MONTH FROM t.tgl_awal) as bulan
		FROM nontender p LEFT JOIN rup rup ON rup.kd_rup::text = ANY(string_to_array(p.kd_rup, ';')) LEFT JOIN jadwal_nontender t ON p.kd_nontender=t.kd_nontender
		WHERE p.tahun_anggaran = ? AND p.realisasi IS NOT NULL AND t.kd_tahapan = 18841 Group by rup.kd_satker, EXTRACT(MONTH FROM t.tgl_awal)`, tahun).Scan(&rekap)
	return rekap
}

func GetRealisasiTenderPerSkpd(tahun int) []RekapRealisasiSatkerBulanan {
	var rekap []RekapRealisasiSatkerBulanan
	db.Raw(`SELECT rup.kd_satker, count(p.kd_tender) paket, sum((p.realisasi -> 0 -> 'nilai_kontrak')::numeric) realisasi, EXTRACT(MONTH FROM t.tgl_awal) as bulan
		FROM tender p LEFT JOIN rup rup ON rup.kd_rup::text = ANY(string_to_array(p.kd_rup, ';')) LEFT JOIN jadwal t ON p.kd_tender=t.kd_tender
		WHERE p.tahun_anggaran = ? AND p.realisasi IS NOT NULL AND t.kd_tahapan = 18803 Group by rup.kd_satker, EXTRACT(MONTH FROM t.tgl_awal)`, tahun).Scan(&rekap)
	return rekap
}

func GetRealisasiPencatatanPerSkpd(tahun int) []RekapRealisasiSatkerBulanan {
	var rekap []RekapRealisasiSatkerBulanan
	db.Raw(`SELECT rup.kd_satker, count(kd_nontender_pct) paket, sum(total_realisasi) realisasi, EXTRACT(MONTH FROM tgl_mulai_paket) as bulan
    		FROM pencatatan p, rup rup WHERE rup.kd_rup::text = ANY(string_to_array(p.kd_rup, ';')) AND p.tahun_anggaran = ?
      		GROUP BY rup.kd_satker, EXTRACT( MONTH FROM tgl_mulai_paket)`, tahun).Scan(&rekap)
	return rekap
}

func GetRealisasiPurchasePerSkpd(tahun int) []RekapRealisasiSatkerBulanan {
	var rekap []RekapRealisasiSatkerBulanan
	db.Raw(`SELECT satker_id, count(*) paket, sum(total_harga) realisasi, EXTRACT(MONTH FROM tanggal_edit_paket) as bulan  FROM katalog_archive
    		WHERE tahun_anggaran = ?  and status_paket IN ('proses_kontrak_ppk', 'melakukan_pengiriman_dan_penerimaan')
      		GROUP BY satker_id, EXTRACT(MONTH FROM tanggal_edit_paket)`, tahun).Scan(&rekap)
	return rekap
}

type PPKRekap struct {
	NipPpk 			string
	NamaPpk			string
	Count 			int
	CountPurchase	int
}

func (c PPKRekap) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]interface{}{
		"nip" : c.NipPpk,
		"nama_ppk" : c.NamaPpk,
		"spse" : c.Count,
		"katalog" : c.CountPurchase,
	})
}

func GetPaketTenderPpk(tahun int) []PPKRekap {
	var rekap []PPKRekap
	db.Raw(`SELECT nip_ppk, nama_ppk, count(*) FROM tender where tahun_anggaran = ? GROUP BY nip_ppk,  nama_ppk`, tahun).Scan(&rekap)
	return rekap
}

func GetDetilPaketTenderPpk(nip string, tahun int) []Tender {
	var rekap []Tender
	db.Raw(`SELECT * FROM tender where tahun_anggaran = ? AND nip_ppk=?`, tahun, nip).Scan(&rekap)
	return rekap
}

func GetPaketNontenderPpk(tahun int) []PPKRekap {
	var rekap []PPKRekap
	db.Raw(`SELECT split_part(nip_nama_ppk, ' - ', 1) nip_ppk, split_part(nip_nama_ppk, ' - ', 2) nama_ppk, count(*) FROM nontender where tahun_anggaran = ? GROUP BY nip_nama_ppk`, tahun).Scan(&rekap)
	return rekap
}

func GetDetilPaketNontenderPpk(nip string, tahun int) []Nontender {
	var rekap []Nontender
	db.Raw(`SELECT * FROM nontender where tahun_anggaran = ? AND split_part(nip_nama_ppk, ' - ', 1) = ?`, tahun, nip).Scan(&rekap)
	return rekap
}

func GetPaketPencatatanPpk(tahun int) []PPKRekap {
	var rekap []PPKRekap
	db.Raw(`SELECT nip_ppk,  nama_ppk, count(*) FROM pencatatan where tahun_anggaran = ? GROUP BY nip_ppk,  nama_ppk`, tahun).Scan(&rekap)
	return rekap
}

func GetDetilPaketPencatatanPpk(nip string, tahun int) []Pencatatan {
	var rekap []Pencatatan
	db.Raw(`SELECT * FROM pencatatan where tahun_anggaran = ? AND nip_ppk=?`, tahun, nip).Scan(&rekap)
	return rekap
}

func GetPaketSwakelolaPpk(tahun int) []PPKRekap {
	var rekap []PPKRekap
	db.Raw(`SELECT nip_ppk, nama_ppk, count(*) FROM swakelola where tahun_anggaran = ? GROUP BY nip_ppk,  nama_ppk`, tahun).Scan(&rekap)
	return rekap
}

func GetPaketPurchasePpk(tahun int) []PPKRekap {
	var rekap []PPKRekap
	db.Raw(`SELECT nip_pegawai nip_ppk, nama_lengkap nama_ppk, count(id)
    FROM isb_pegawai_purchase p, katalog_archive ip WHERE p.kd_user_pegawai = ip.kd_user_ppk AND tahun_anggaran = ? GROUP BY nip_pegawai, nama_lengkap`, tahun).Scan(&rekap)
	return rekap
}

func GetDetilPaketPurchasePpk(nip string, tahun int) []Katalog {
	var rekap []Katalog
	db.Raw(`SELECT ip.* FROM isb_pegawai_purchase p, katalog_archive ip WHERE p.kd_user_pegawai = ip.kd_user_ppk AND tahun_anggaran = ? AND nip_pegawai=?`, tahun, nip).Scan(&rekap)
	return rekap
}
