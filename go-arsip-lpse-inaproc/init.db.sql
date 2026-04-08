INSERT INTO "public"."metode_pengadaan" ("id", "metode") VALUES
(0, 'Belum Ditentukan'),
(1, 'Lelang Umum'),
(2, 'Lelang Sederhana'),
(3, 'Lelang Terbatas'),
(4, 'Seleksi Umum'),
(5, 'Seleksi Sederhana'),
(6, 'Pemilihan Langsung'),
(7, 'Penunjukan Langsung'),
(8, 'Pengadaan Langsung'),
(9, 'e-Purchasing'),
(10, 'Sayembara'),
(11, 'Kontes'),
(12, 'Lelang Cepat'),
(13, 'Tender'),
(14, 'Tender Cepat'),
(15, 'Seleksi'),
(16, 'Dikecualikan');

DROP TABLE isb_tender;
DROP TABLE isb_nontender;
DROP TABLE isb_swakelola;
DROP TABLE isb_pencatatan;
DROP TABLE isb_tokodaring;
DROP TABLE isb_tahap_tender;
DROP TABLE isb_tahap_nontender;
DROP TABLE isb_pegawai_purchase;
DROP TABLE isb_purchase_penyedia;
DROP TABLE isb_purchase_penyedia6;
DROP TABLE isb_purchase6;
DROP TABLE isb_purchase;
DROP TABLE isb_rup_program;
DROP TABLE isb_rup_kegiatan;
DROP TABLE isb_rup_swakelola;
DROP TABLE isb_rup_paket;
DROP TABLE isb_struktur_anggaran;
DROP TABLE isb_kabupaten;
DROP TABLE isb_provinsi;
ALTER TABLE itkp DROP CONSTRAINT "fk_itkp_satker";
DROP TABLE isb_satker;
