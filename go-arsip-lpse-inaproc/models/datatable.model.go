package models

import (
	"log/slog"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"gorm.io/gorm"
)

type Datatable[T any] struct {
	Draw            int   `json:"draw"`
	RecordsTotal    int64 `json:"recordsTotal"`
	RecordsFiltered int64 `json:"recordsFiltered"`
	Data            []T   `json:"data"`
}

func (data *Datatable[T]) Populate(db *gorm.DB, param map[string]string, columns ...string) {
	start, _ := strconv.Atoi(param["start"])
	length, _ := strconv.Atoi(param["length"])
	queryCount := db.Model(&data.Data)
	queryCount.Count(&data.RecordsTotal)
	// queryDB := db
	search := param["search[value]"]
	if search != "" {
		search = "%" + search + "%"
		searchable := false
		first := true
		var filter strings.Builder
		var paramFilter []interface{}
		for i, column := range columns {
			searchable, _ = strconv.ParseBool(param["columns["+strconv.Itoa(i)+"][searchable]"])
			if searchable {
				if first {
					filter.WriteString(" LOWER("+column+"::varchar) like ?")
					// queryCount.Where(" LOWER("+column+"::varchar) like ?", search)
					// querydb.Where(" LOWER("+column+"::varchar) like ?", search)
					first = false
				} else {
					filter.WriteString(" OR LOWER("+column+"::varchar) like ?")
					// queryCount.Or(" LOWER("+column+"::varchar) like ?", search)
					// querydb.Or(" LOWER("+column+"::varchar) like ?", search)
				}
				paramFilter = append(paramFilter, search)
			}
		}
		log.Info("filter.Len() : ", filter.Len())
		if filter.Len() > 0 {
			queryCount.Where("("+filter.String()+")", paramFilter...)
			db = db.Where("("+filter.String()+")", paramFilter...)
			log.Info("filter dtable : ", filter.String())
		}
	}
	queryCount.Count(&data.RecordsFiltered)
	// ORDERING
	var requestColumn string
	columnIdx := 0
	orderable := false
	for i, column := range columns {
		requestColumn = param["order["+strconv.Itoa(i)+"][column]"]
		if len(requestColumn) == 0 {
			continue
		}
		columnIdx, _ = strconv.Atoi(requestColumn)
		orderable, _ = strconv.ParseBool(param["columns["+requestColumn+"][orderable]"])
		if orderable && (columnIdx >= 0 && columnIdx < len(column)) {
			order_status := param["order["+strconv.Itoa(i)+"][dir]"]
			db = db.Order(columns[columnIdx] + " " + order_status)
		}
	}
	db.Limit(length).Offset(start).Find(&data.Data)
	data.Draw, _ = strconv.Atoi(param["draw"])
}

func GetDataTableAgency(c *fiber.Ctx) error {
	return populate(db.Model(&Agency{}), c, &Agency{}, "id", "agc_nama", "agc_alamat", "agc_tgl_daftar")
}

func GetDataTableUkpbj(c *fiber.Ctx) error {
	var datas []Ukpbj
	return populate(db.Model(&Ukpbj{}), c, &datas, "id", "nama", "alamat", "tgl_daftar")
}

func GetDataTablePerubahanData(c *fiber.Ctx, usrsession UserSession) error {
	var datas []PerubahanData
	orm := db.Model(&PerubahanData{}).Preload("Pegawai")
	if usrsession.IsPp() || usrsession.IsPpk() || usrsession.IsPokja() {
		orm.Where("dok_id > 0 AND peg_id=?", usrsession.Id)
	} else {
		orm.Where("dok_id > 0")
	}
	return populate(orm, c, &datas,  "id", "nomor", "perihal", "created_at", "status", "pegawai.peg_nama")
}

func GetDataTableVerifikasi(c *fiber.Ctx) error {
	var datas []Pegawai
	orm := db.Model(&Pegawai{}).Where("usrgroup IN ('PPK', 'PP', 'POKJA', 'PEGAWAI', 'ARSIPARIS', '') and deleted_at IS NULL")
	statusFilter := c.Query("status")
	if statusFilter != "" && statusFilter != "all" {
		orm = orm.Where("peg_status = ?", statusFilter)
	}
	return populate(orm, c, &datas, "id", "peg_nama", "peg_nip", "peg_email", "peg_namauser", "peg_status")
}

func GetDataTablePaketSirup(c *fiber.Ctx, tahun int, satker string, metode string, jenis string) error {
	var datas []PaketSirup
	orm := db.Model(&PaketSirup{})
	orm = orm.Where("paket_terhapus = ?", false)
	if tahun > 0 {
		orm = orm.Where("tahun = ?", tahun)
	}
	if len(satker) > 0 {
		orm = orm.Where("id_satker = ?", satker)
	}
	if len(metode) > 0 {
		var metodeId int;
		for i, v := range metodePengadaan {
			if v == metode {
				metodeId = i;
				break;
			}
		}
		orm = orm.Where("metode_pengadaan = ?", metodeId)
	}
	if len(jenis) > 0 {
		var jenisId int;
		for i, v := range jenisPengadaan {
			if v == jenis {
				jenisId = i;
				break;
			}
		}
		orm = orm.Where("jenis_paket = ?", jenisId)
	}
	return populate(orm, c, &datas,  "id", "nama", "pagu", "tahun", "kode_kldi")
}

func GetDataTableSwakelolaSirup(c *fiber.Ctx) error {
	orm := db.Model(&SwakelolaSirup{})
	var datas []SwakelolaSirup
	return populate(orm, c, &datas,  "id", "nama", "pagu", "tahun", "kode_kldi")
}

func GetDataTablePegawai(c *fiber.Ctx, usrgroup string) error {
	orm := db.Model(&Pegawai{})
	if usrgroup == ADMIN {
		orm.Where("peg_status IN (1, 2) AND usrgroup NOT IN ('ADMIN')")
	} else {
		orm.Where("peg_status IN (1, 2) AND usrgroup IN ('PPK', 'POKJA', 'PP', 'ARSIPARIS')")
	}
	var datas []Pegawai
	return populate(orm, c, &datas,  "id", "peg_nama", "peg_nip", "peg_namauser")
}

func GetDataTablePaket(c *fiber.Ctx, id uint, isPPK, isUkpbj, isPokja, isPp bool) error {
	orm := db.Model(&Paket{})
	pegawai := GetPegawai(id)
	if isPPK && pegawai.IsApprove(){
		orm = orm.Where("ppk_id = ?", id)
	} else if isUkpbj {
		orm = orm.Where("ukpbj_id <> 0 OR status >= 1")
	} else if isPokja && pegawai.IsApprove() {
		orm = orm.Where("pnt_id IN (SELECT pnt_id FROM anggota_panitia WHERE peg_id=? and deleted_at IS NULL)", id)
	} else if isPp && pegawai.IsApprove() {
		orm = orm.Where("pp_id = ?", id)
	} else {
		return populateEmpty(c)
	}
	metode := c.Query("metode")
	if metode != "" && metode != "all" {
		orm = orm.Where("metode = ?", metode)
	}
	var datas []Paket
	return populate(orm, c, &datas,  "id", "nama", "pagu", "hps", "Created_at", "created_by", "status")
}

func GetDataTableTemplates(c *fiber.Ctx) error {
	orm :=db.Model(&Templates{})
	var datas []Templates
	return populate(orm, c, &datas,  "id", "nama", "content")
}

func GetDataTableReviu(c *fiber.Ctx) error {
	orm := db.Model(&Reviu{})
	var datas []Reviu
	return populate(orm, c, &datas,  "id", "bidang", "content", "opsi1", "opsi2")
}

func GetDataTablePanitia(c *fiber.Ctx) error {
	orm := db.Model(&Panitia{})
	var datas []Panitia
	return populate(orm, c, &datas,  "id", "nama", "tahun")
}

func GetDataTablePp(c *fiber.Ctx) error {
	orm := db.Model(&PejabatPengadaan{})
	var datas []PejabatPengadaan
	return populate(orm, c, &datas,  "id", "groups", "tahun", "no_sk")
}

func GetDataTableInbox(c *fiber.Ctx, id uint) error {
	orm := db.Model(&Inbox{}).Where("peg_id=?", id)
	var datas []Inbox
	return populate(orm, c, &datas,  "id", "subject", "enqueue_date", "status")
}

func GetDataTableDocTemplate(c *fiber.Ctx) error {
	orm := db.Model(&DokTemplate{})
	var datas []DokTemplate
	return populate(orm, c, &datas,  "id", "jenis", "periode_awal", "periode_akhir")
}

func GetDataTableBukuTamu(c *fiber.Ctx, isUkpbj bool) error {
	orm := db.Model(&BukuTamu{})
	// if isUkpbj {
	// 	orm.Where("kategori = 'non-pengadaan'")
	// } else  {
	// 	orm.Where("kategori = 'pengadaan'")
	// }
	var datas []BukuTamu
	return populate(orm, c, &datas, "id", "nama", "nama_perusahaan", "email", "keperluan")
}

func GetDataTableFeedback(c *fiber.Ctx) error {
	orm := db.Model(&Feedback{})
	var datas []Feedback
	return populate(orm, c, &datas, "id", "nama", "nama_perusahaan", "feedback", "kepuasan")
}

func GetDataTableDocument(c *fiber.Ctx) error {
	orm := db.Model(&Document{})
	var datas []Document
	return populate(orm, c, &datas, "id", "filename", "filesize", "filedate")
}

func GetDataTableAdminDocument(c *fiber.Ctx) error {
	var datas []Document
	orm := db.Model(&Document{}).Where("jenis = ?", ADMIN_DOK)
	return populate(orm, c, &datas, "id", "filename", "filesize", "filedate")
}

func GetDataTableChecklist(c *fiber.Ctx) error {
	orm := db.Model(&Checklist{})
	var datas []Checklist
	return populate(orm, c, &datas, "id", "jenis", "metode", "periode_awal", "periode_akhir")
}

func populate(db *gorm.DB, c *fiber.Ctx, result interface{}, columns ...string) error {
	start := c.QueryInt("start", 0)
	length := c.QueryInt("length", 10)
	var total int64
	db.Session(&gorm.Session{}).Count(&total)
	search := c.Query("search[value]")
	if search != "" {
		search = "%" + search + "%"
		searchable := false
		first := true
		var filter strings.Builder
		var paramFilter []interface{}
		for i, column := range columns {
			searchable = c.QueryBool("columns["+strconv.Itoa(i)+"][searchable]")
			if searchable {
				if first {
					filter.WriteString(" LOWER("+column+"::varchar) like ?")
					first = false
				} else {
					filter.WriteString(" OR LOWER("+column+"::varchar) like ?")
				}
				paramFilter = append(paramFilter, search)
			}
		}
		if filter.Len() > 0 {
			db.Where(filter.String(), paramFilter...)
		}
	}
	var filterCount int64
	db.Session(&gorm.Session{}).Count(&filterCount)
	// ORDERING
	var requestColumn string
	columnIdx := 0
	orderable := false
	for i, column := range columns {
		requestColumn = c.Query("order["+strconv.Itoa(i)+"][column]")
		if len(requestColumn) == 0 {
			continue
		}
		columnIdx, _ = strconv.Atoi(requestColumn)
		orderable  = c.QueryBool("columns["+requestColumn+"][orderable]")
		if orderable && (columnIdx >= 0 && columnIdx < len(column)) {
			order_status := c.Query("order["+strconv.Itoa(i)+"][dir]")
			db = db.Order(columns[columnIdx] + " " + order_status)
		}
	}
	db.Debug().Limit(length).Offset(start).Find(result)
	draw := c.QueryInt("draw")
	responseData := fiber.Map{
		"draw" : draw,
		"recordsTotal" : total,
		"recordsFiltered" : filterCount,
		"data" : result,
	}
	return c.JSON(responseData)
}

func populateEmpty(c *fiber.Ctx) error {
	draw := c.QueryInt("draw")
	responseData := fiber.Map{
		"draw" : draw,
		"recordsTotal" : 0,
		"recordsFiltered" : 0,
		"data" : []interface{}{},
	}
	return c.JSON(responseData)
}

/**
 * pre requirement : kolom pk harus di awal select
 */
func GetDataTable(c *fiber.Ctx, result interface{}, columns []string, queryFrom string, param ...interface{}) error {
	start := c.QueryInt("start", 1)
	length := c.QueryInt("length", 10)
	total := Count("SELECT COUNT ("+columns[0]+") "+queryFrom, param...)
	// dbOrm := db
	search := c.Query("search[value]")
	if search != "" {
		var filter strings.Builder
		search = "%" + search + "%"
		searchable := false
		first := true
		for i, column := range columns {
			searchable = c.QueryBool("columns["+strconv.Itoa(i)+"][searchable]")
			if searchable {
				if first {
					filter.WriteString(" LOWER("+column+"::varchar) like ?")
					first = false
				} else {
					filter.WriteString(" OR LOWER("+column+"::varchar) like ?")
				}
				param = append(param, search)
			}
		}
		if filter.Len() > 0 {
			if strings.Contains(strings.ToUpper(queryFrom), " WHERE ") {
				queryFrom += " AND "
			} else {
				queryFrom += " WHERE "
			}
			queryFrom += " ("+filter.String()+")"
		}
	}
	filterCount := Count("SELECT COUNT ("+columns[0]+") "+queryFrom, param...)

	// ORDERING
	var requestColumn string
	columnIdx := 0
	orderable := false
	var filter strings.Builder
	for i, column := range columns {
		requestColumn = c.Query("order["+strconv.Itoa(i)+"][column]")
		if len(requestColumn) == 0 {
			continue
		}
		columnIdx, _ = strconv.Atoi(requestColumn)
		orderable  = c.QueryBool("columns["+requestColumn+"][orderable]")
		if orderable && (columnIdx >= 0 && columnIdx < len(column)) {
			order_status := c.Query("order["+strconv.Itoa(i)+"][dir]")
			if(filter.Len() > 1) {
				filter.WriteString(",")
			}
			filter.WriteString(columns[columnIdx] + " " + order_status);
		}
	}
	if filter.Len() > 0 {
		queryFrom += " ORDER BY "+filter.String()
	}
	// no change needed if removing it, but let's just comment it out or remove it
	// param = append(param) // Removed redundant append
	if length > 0 {
		queryFrom += " LIMIT "+strconv.Itoa(length)
	}
	if start > 0 {
		queryFrom += " OFFSET "+strconv.Itoa(start)
	}
	// slog.Info("limit ", "limit", length, "offset", start, "params", param);
	err := db.Raw("SELECT "+strings.Join(columns, ",")+" "+queryFrom, param...).Scan(&result).Error
	if err != nil {
		slog.Error("query datatable failed ", "err", err)
	}
	draw := c.QueryInt("draw")
	responseData := fiber.Map{
		"draw" : draw,
		"recordsTotal" : total,
		"recordsFiltered" : filterCount,
		"data" : result,
	}
	return c.JSON(responseData)
}
