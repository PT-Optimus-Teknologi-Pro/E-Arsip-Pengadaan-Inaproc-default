UPDATE isb_rup_paket SET tahun_pertama=NULL WHERE tahun_pertama='';
UPDATE isb_rup_paket SET kd_rup_tahun_pertama=NULL WHERE kd_rup_tahun_pertama='';
ALTER TABLE isb_rup_paket ALTER COLUMN tahun_pertama TYPE numeric USING tahun_pertama::numeric;
ALTER TABLE isb_rup_paket ALTER COLUMN kd_rup_tahun_pertama TYPE bigint USING kd_rup_tahun_pertama::bigint;
ALTER TABLE isb_purchase6 ALTER COLUMN jml_produk TYPE numeric(19,2) USING jml_produk::numeric;
ALTER TABLE checklist ADD COLUMN metode int not null default 0;
ALTER TABLE checklist DROP COLUMN keterangan;
ALTER TABLE checklist DROP COLUMN dok_template;
ALTER TABLE checklist DROP COLUMN required;

UPDATE checklist_paket p SET jenis = d.jenis, dok_template = d.id FROM dok_templates d LEFT JOIN checklist_dok c ON c.dok_id=d.id WHERE c.id = p.chk_id;
ALTER TABLE pejabat_pengadaan RENAME COLUMN "group" TO groups;
ALTER TABLE itkp ALTER COLUMN kd_satker TYPE NUMERIC;
