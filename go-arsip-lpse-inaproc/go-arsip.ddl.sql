--
-- PostgreSQL database dump
--

-- Dumped from database version 15.1
-- Dumped by pg_dump version 15.1

SET statement_timeout = 0;
SET lock_timeout = 0;
SET idle_in_transaction_session_timeout = 0;
SET client_encoding = 'UTF8';
SET standard_conforming_strings = on;
SELECT pg_catalog.set_config('search_path', '', false);
SET check_function_bodies = false;
SET xmloption = content;
SET client_min_messages = warning;
SET row_security = off;

--
-- Name: seq_agency; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.seq_agency
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.seq_agency OWNER TO postgres;

SET default_tablespace = '';

SET default_table_access_method = heap;

--
-- Name: agency; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.agency (
    id numeric(19,0) DEFAULT nextval('public.seq_agency'::regclass) NOT NULL,
    updated_by character varying(100) DEFAULT 'ADMIN'::character varying,
    updated_at timestamp without time zone DEFAULT now(),
    agc_nama text NOT NULL,
    agc_alamat text NOT NULL,
    agc_telepon text NOT NULL,
    agc_fax text,
    agc_website text,
    agc_tgl_daftar timestamp without time zone NOT NULL,
    instansi_id text,
    deleted_at timestamp without time zone,
    created_at timestamp without time zone
);


ALTER TABLE public.agency OWNER TO postgres;

--
-- Name: anggarans; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.anggarans (
    id bigint NOT NULL,
    created_at timestamp with time zone,
    updated_at timestamp with time zone,
    deleted_at timestamp with time zone,
    kode_rekening text,
    nilai numeric,
    tahun bigint,
    uraian text,
    sumber text,
    stk_id text
);


ALTER TABLE public.anggarans OWNER TO postgres;

--
-- Name: anggarans_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.anggarans_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.anggarans_id_seq OWNER TO postgres;

--
-- Name: anggarans_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.anggarans_id_seq OWNED BY public.anggarans.id;


--
-- Name: anggota_panitia; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.anggota_panitia (
    id bigint NOT NULL,
    created_at timestamp with time zone,
    updated_at timestamp with time zone,
    deleted_at timestamp with time zone,
    peg_id bigint NOT NULL,
    pnt_id bigint NOT NULL
);


ALTER TABLE public.anggota_panitia OWNER TO postgres;

--
-- Name: anggota_panitia_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.anggota_panitia_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.anggota_panitia_id_seq OWNER TO postgres;

--
-- Name: anggota_panitia_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.anggota_panitia_id_seq OWNED BY public.anggota_panitia.id;


--
-- Name: berita_acara; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.berita_acara (
    id bigint NOT NULL,
    created_at timestamp with time zone,
    updated_at timestamp with time zone,
    deleted_at timestamp with time zone,
    nomor text,
    jenis text,
    uraian text,
    tanggal timestamp with time zone,
    dok_id bigint
);


ALTER TABLE public.berita_acara OWNER TO postgres;

--
-- Name: berita_acara_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.berita_acara_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.berita_acara_id_seq OWNER TO postgres;

--
-- Name: berita_acara_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.berita_acara_id_seq OWNED BY public.berita_acara.id;


--
-- Name: buku_tamu; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.buku_tamu (
    id bigint NOT NULL,
    created_at timestamp with time zone,
    updated_at timestamp with time zone,
    deleted_at timestamp with time zone,
    nama text,
    nama_perusahaan text,
    keperluan text,
    kode_tender bigint,
    feedback text,
    kepuasan text,
    status bigint,
    email text,
    phone text,
    jabatan text,
    dok_id bigint
);


ALTER TABLE public.buku_tamu OWNER TO postgres;

--
-- Name: buku_tamu_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.buku_tamu_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.buku_tamu_id_seq OWNER TO postgres;

--
-- Name: buku_tamu_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.buku_tamu_id_seq OWNED BY public.buku_tamu.id;


--
-- Name: checklist; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.checklist (
    id bigint NOT NULL,
    created_at timestamp with time zone,
    updated_at timestamp with time zone,
    deleted_at timestamp with time zone,
    jenis bigint NOT NULL,
    keterangan text,
    dok_template bigint,
    required boolean
);


ALTER TABLE public.checklist OWNER TO postgres;

--
-- Name: checklist_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.checklist_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.checklist_id_seq OWNER TO postgres;

--
-- Name: checklist_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.checklist_id_seq OWNED BY public.checklist.id;


--
-- Name: checklist_paket; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.checklist_paket (
    id bigint NOT NULL,
    created_at timestamp with time zone,
    updated_at timestamp with time zone,
    deleted_at timestamp with time zone,
    pkt_id bigint,
    chk_id bigint,
    created_by bigint,
    dok_id bigint,
    status bigint,
    tgl_ajukan timestamp with time zone,
    tgl_approve timestamp with time zone,
    tgl_revisi timestamp with time zone,
    tgl_tolak timestamp with time zone
);


ALTER TABLE public.checklist_paket OWNER TO postgres;

--
-- Name: checklist_paket_history; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.checklist_paket_history (
    id bigint NOT NULL,
    created_at timestamp with time zone,
    updated_at timestamp with time zone,
    deleted_at timestamp with time zone,
    created_by bigint,
    check_id bigint,
    dok_id bigint
);


ALTER TABLE public.checklist_paket_history OWNER TO postgres;

--
-- Name: checklist_paket_history_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.checklist_paket_history_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.checklist_paket_history_id_seq OWNER TO postgres;

--
-- Name: checklist_paket_history_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.checklist_paket_history_id_seq OWNED BY public.checklist_paket_history.id;


--
-- Name: checklist_paket_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.checklist_paket_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.checklist_paket_id_seq OWNER TO postgres;

--
-- Name: checklist_paket_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.checklist_paket_id_seq OWNED BY public.checklist_paket.id;


--
-- Name: seq_document; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.seq_document
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.seq_document OWNER TO postgres;

--
-- Name: document; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.document (
    id numeric(19,0) DEFAULT nextval('public.seq_document'::regclass) NOT NULL,
    versi numeric(6,0) DEFAULT 0 NOT NULL,
    updated_by character varying(100) DEFAULT 'ADMIN'::character varying NOT NULL,
    updated_at timestamp without time zone DEFAULT now() NOT NULL,
    filename text,
    filesize numeric(19,0),
    filedate timestamp without time zone NOT NULL,
    filepath text,
    filehash character varying(50),
    deleted_at timestamp without time zone,
    created_at timestamp without time zone,
    jenis character varying NOT NULL,
    peg_id integer NOT NULL
);


ALTER TABLE public.document OWNER TO postgres;

--
-- Name: dok_persiapan; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.dok_persiapan (
    id bigint NOT NULL,
    created_at timestamp with time zone,
    updated_at timestamp with time zone,
    deleted_at timestamp with time zone,
    pkt_id bigint,
    chk_id bigint,
    created_by bigint,
    dok_id bigint
);


ALTER TABLE public.dok_persiapan OWNER TO postgres;

--
-- Name: dok_persiapan_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.dok_persiapan_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.dok_persiapan_id_seq OWNER TO postgres;

--
-- Name: dok_persiapan_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.dok_persiapan_id_seq OWNED BY public.dok_persiapan.id;


--
-- Name: dok_templates; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.dok_templates (
    id bigint NOT NULL,
    created_at timestamp with time zone,
    updated_at timestamp with time zone,
    deleted_at timestamp with time zone,
    jenis text NOT NULL,
    periode_awal timestamp with time zone,
    periode_akhir timestamp with time zone,
    dok_id bigint,
    keterangan text
);


ALTER TABLE public.dok_templates OWNER TO postgres;

--
-- Name: dok_templates_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.dok_templates_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.dok_templates_id_seq OWNER TO postgres;

--
-- Name: dok_templates_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.dok_templates_id_seq OWNED BY public.dok_templates.id;


--
-- Name: feedback; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.feedback (
    id bigint NOT NULL,
    created_at timestamp with time zone,
    updated_at timestamp with time zone,
    deleted_at timestamp with time zone,
    nama text,
    nama_perusahaan text,
    feedback text,
    kepuasan bigint
);


ALTER TABLE public.feedback OWNER TO postgres;

--
-- Name: feedback_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.feedback_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.feedback_id_seq OWNER TO postgres;

--
-- Name: feedback_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.feedback_id_seq OWNED BY public.feedback.id;


--
-- Name: hak_akses; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.hak_akses (
    id bigint NOT NULL,
    created_at timestamp with time zone,
    updated_at timestamp with time zone,
    deleted_at timestamp with time zone,
    peg_id bigint,
    periode_awal timestamp with time zone,
    periode_akhir timestamp with time zone,
    usrgroup text
);


ALTER TABLE public.hak_akses OWNER TO postgres;

--
-- Name: hak_akses_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.hak_akses_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.hak_akses_id_seq OWNER TO postgres;

--
-- Name: hak_akses_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.hak_akses_id_seq OWNED BY public.hak_akses.id;


--
-- Name: inbox; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.inbox (
    id bigint NOT NULL,
    created_at timestamp with time zone,
    updated_at timestamp with time zone,
    deleted_at timestamp with time zone,
    enqueue_date timestamp with time zone,
    subject text,
    content text,
    status text,
    peg_id bigint
);


ALTER TABLE public.inbox OWNER TO postgres;

--
-- Name: inboxes_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.inboxes_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.inboxes_id_seq OWNER TO postgres;

--
-- Name: inboxes_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.inboxes_id_seq OWNED BY public.inbox.id;


--
-- Name: isb_kabupaten; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.isb_kabupaten (
    kd_kabupaten numeric NOT NULL,
    kd_provinsi numeric NOT NULL,
    nama_kabupaten character varying NOT NULL
);


ALTER TABLE public.isb_kabupaten OWNER TO postgres;

--
-- Name: isb_nontender; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.isb_nontender (
    tahun_anggaran bigint,
    kd_klpd text,
    nama_klpd text,
    jenis_klpd text,
    kd_satker bigint,
    kd_satker_str text,
    nama_satker text,
    kd_lpse bigint,
    nama_lpse text,
    kd_nontender bigint NOT NULL,
    kd_pkt_dce bigint,
    kd_rup text,
    nama_paket text,
    pagu numeric,
    hps numeric,
    sumber_dana text,
    mak text,
    kualifikasi_paket text,
    jenis_pengadaan text,
    mtd_pemilihan text,
    kontrak_pembayaran text,
    status_nontender text,
    versi_nontender bigint,
    ket_diulang text,
    ket_ditutup text,
    tgl_buat_paket timestamp without time zone,
    tgl_kolektif_kolegial timestamp without time zone,
    tgl_pengumuman_nontender timestamp without time zone,
    nip_nama_ppk text,
    nip_nama_pp text,
    nip_nama_pokja text,
    lokasi_pekerjaan text,
    url_lpse text,
    realisasi jsonb
);


ALTER TABLE public.isb_nontender OWNER TO postgres;

--
-- Name: isb_pencatatan; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.isb_pencatatan (
    tahun_anggaran bigint,
    kd_klpd text,
    nama_klpd text,
    jenis_klpd text,
    kd_satker bigint,
    kd_satker_str text,
    nama_satker text,
    kd_lpse bigint,
    kd_nontender_pct bigint NOT NULL,
    kd_pkt_dce bigint,
    kd_rup text,
    nama_paket text,
    pagu numeric,
    total_realisasi numeric,
    nilai_pdn_pct numeric,
    nilai_umk_pct numeric,
    sumber_dana text,
    uraian_pekerjaan text,
    informasi_lainnya text,
    kategori_pengadaan text,
    mtd_pemilihan text,
    bukti_pembayaran text,
    status_nontender_pct text,
    status_nontender_pct_ket text,
    alasan_pembatalan text,
    nip_ppk text,
    nama_ppk text,
    tgl_buat_paket timestamp without time zone,
    tgl_mulai_paket timestamp without time zone,
    tgl_selesai_paket timestamp without time zone,
    realisasi jsonb
);


ALTER TABLE public.isb_pencatatan OWNER TO postgres;

--
-- Name: isb_provinsi; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.isb_provinsi (
    kd_provinsi numeric NOT NULL,
    nama_provinsi character varying NOT NULL
);


ALTER TABLE public.isb_provinsi OWNER TO postgres;

--
-- Name: isb_purchase; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.isb_purchase (
    tahun_anggaran integer,
    kd_klpd character varying,
    satker_id numeric,
    nama_satker character varying,
    alamat_satker character varying,
    npwp_satker character varying,
    kd_paket numeric,
    no_paket character varying,
    nama_paket character varying,
    kd_rup numeric,
    nama_sumber_dana character varying,
    kode_anggaran character varying,
    kd_komoditas numeric,
    kd_produk numeric,
    kd_penyedia numeric,
    kd_penyedia_distributor numeric,
    jml_jenis_produk integer,
    kuantitas numeric,
    harga_satuan numeric,
    ongkos_kirim numeric,
    total_harga numeric,
    kd_user_pokja numeric,
    no_telp_user_pokja character varying,
    email_user_pokja character varying,
    kd_user_ppk numeric,
    ppk_nip character varying,
    jabatan_ppk character varying,
    tanggal_buat_paket timestamp without time zone,
    tanggal_edit_paket timestamp without time zone,
    deskripsi character varying,
    status_paket character varying,
    paket_status_str character varying,
    catatan_produk character varying,
    kd_provinsi_wilayah_harga numeric,
    kd_kabupaten_wilayah_harga numeric,
    id integer NOT NULL
);


ALTER TABLE public.isb_purchase OWNER TO postgres;

--
-- Name: isb_purchase_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.isb_purchase_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.isb_purchase_id_seq OWNER TO postgres;

--
-- Name: isb_purchase_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.isb_purchase_id_seq OWNED BY public.isb_purchase.id;


--
-- Name: isb_rup_kegiatan; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.isb_rup_kegiatan (
    tahun_anggaran integer NOT NULL,
    kd_klpd character varying NOT NULL,
    nama_klpd character varying,
    jenis_klpd character varying,
    kd_satker numeric NOT NULL,
    kd_program character varying NOT NULL,
    kd_kegiatan character varying NOT NULL,
    kd_kegiatan_str character varying,
    nama_kegiatan character varying,
    pagu_kegiatan numeric,
    kd_kegiatan_lokal character varying,
    is_deleted boolean
);


ALTER TABLE public.isb_rup_kegiatan OWNER TO postgres;

--
-- Name: isb_rup_paket; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.isb_rup_paket (
    tahun_anggaran integer,
    kd_klpd character varying,
    nama_klpd character varying,
    jenis_klpd character varying,
    kd_satker numeric,
    kd_satker_str character varying,
    nama_satker character varying,
    kd_rup character varying NOT NULL,
    nama_paket character varying,
    pagu numeric,
    kd_metode_pengadaan character varying,
    metode_pengadaan character varying,
    kd_jenis_pengadaan character varying,
    jenis_pengadaan character varying,
    status_pradipa character varying,
    status_pdn character varying,
    status_ukm character varying,
    alasan_non_ukm character varying,
    status_konsolidasi character varying,
    tipe_paket character varying,
    kd_rup_swakelola character varying,
    kd_rup_lokal character varying,
    volume_pekerjaan character varying,
    urarian_pekerjaan character varying,
    spesifikasi_pekerjaan character varying,
    tgl_awal_pemilihan timestamp without time zone,
    tgl_akhir_pemilihan timestamp without time zone,
    tgl_awal_kontrak timestamp without time zone,
    tgl_akhir_kontrak timestamp without time zone,
    tgl_awal_pemanfaatan timestamp without time zone,
    tgl_akhir_pemanfaatan timestamp without time zone,
    tgl_buat_paket timestamp without time zone,
    tgl_pengumuman_paket timestamp without time zone,
    nip_ppk character varying,
    nama_ppk character varying,
    username_ppk character varying,
    status_aktif_rup boolean,
    status_delete_rup boolean,
    status_umumkan_rup character varying,
    status_dikecualikan boolean,
    alasan_dikecualikan character varying,
    tahun_pertama character varying,
    kd_rup_tahun_pertama character varying,
    nomor_kontrak character varying,
    spp_aspek_ekonomi boolean,
    spp_aspek_sosial boolean,
    spp_aspek_lingkungan boolean
);


ALTER TABLE public.isb_rup_paket OWNER TO postgres;

--
-- Name: isb_rup_program; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.isb_rup_program (
    tahun_anggaran integer NOT NULL,
    kd_klpd character varying NOT NULL,
    nama_klpd character varying NOT NULL,
    jenis_klpd character varying,
    kd_satker numeric,
    kd_program character varying NOT NULL,
    kd_program_str character varying,
    nama_program character varying,
    pagu_program numeric,
    kd_program_lokal numeric,
    is_deleted boolean
);


ALTER TABLE public.isb_rup_program OWNER TO postgres;

--
-- Name: isb_rup_swakelola; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.isb_rup_swakelola (
    tahun_anggaran integer,
    kd_klpd character varying,
    nama_klpd character varying,
    jenis_klpd character varying,
    kd_satker numeric,
    kd_satker_str character varying,
    nama_satker character varying,
    kd_rup numeric NOT NULL,
    nama_paket character varying,
    pagu numeric,
    tipe_swakelola smallint,
    volume_pekerjaan character varying,
    uraian_pekerjaan character varying,
    kd_klpd_penyelenggara character varying,
    nama_klpd_penyelenggara character varying,
    nama_satker_penyelenggara character varying,
    tgl_awal_pelaksanaan_kontrak timestamp without time zone,
    tgl_akhir_pelaksanaan_kontrak timestamp without time zone,
    tgl_buat_paket timestamp without time zone,
    tgl_pengumuman_paket timestamp without time zone,
    nip_ppk character varying,
    nama_ppk character varying,
    username_ppk character varying,
    kd_rup_lokal character varying,
    status_aktif_rup boolean,
    status_delete_rup boolean,
    status_umumkan_rup character varying
);


ALTER TABLE public.isb_rup_swakelola OWNER TO postgres;

--
-- Name: isb_satker; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.isb_satker (
    kd_satker numeric NOT NULL,
    kd_satker_str text,
    nama_satker text,
    alamat text,
    telepon text,
    fax text,
    kodepos text,
    status_satker text,
    ket_satker text,
    jenis_satker text,
    kd_klpd text,
    nama_klpd text,
    jenis_klpd text,
    kode_eselon text
);


ALTER TABLE public.isb_satker OWNER TO postgres;

--
-- Name: isb_struktur_anggaran; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.isb_struktur_anggaran (
    tahun_anggaran bigint,
    kd_klpd text,
    nama_klpd text,
    kd_satker bigint,
    kd_satker_str text,
    nama_satker text,
    belanja_operasi numeric,
    belanja_modal numeric,
    belanja_btt numeric,
    belanja_non_pengadaan numeric,
    belanja_pengadaan numeric,
    total_belanja numeric,
    id bigint NOT NULL
);


ALTER TABLE public.isb_struktur_anggaran OWNER TO postgres;

--
-- Name: isb_struktur_anggaran_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.isb_struktur_anggaran_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.isb_struktur_anggaran_id_seq OWNER TO postgres;

--
-- Name: isb_struktur_anggaran_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.isb_struktur_anggaran_id_seq OWNED BY public.isb_struktur_anggaran.id;


--
-- Name: isb_swakelola; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.isb_swakelola (
    tahun_anggaran bigint,
    kd_klpd text,
    nama_klpd text,
    jenis_klpd text,
    kd_satker text,
    kd_satker_str text,
    nama_satker text,
    kd_lpse bigint,
    kd_swakelola_pct bigint NOT NULL,
    kd_pkt_dce bigint,
    kd_rup text,
    nama_paket text,
    pagu numeric,
    total_realisasi numeric,
    nilai_pdn_pct numeric,
    nilai_umk_pct numeric,
    sumber_dana text,
    uraian_pekerjaan text,
    informasi_lainnya text,
    tipe_swakelola bigint,
    tipe_swakelola_nama text,
    status_swakelola_pct text,
    status_swakelola_pct_ket text,
    alasan_pembatalan text,
    nip_ppk text,
    nama_ppk text,
    tgl_buat_paket timestamp without time zone,
    tgl_mulai_paket timestamp without time zone,
    tgl_selesai_paket timestamp without time zone,
    realisasi jsonb
);


ALTER TABLE public.isb_swakelola OWNER TO postgres;

--
-- Name: isb_tahap_nontender; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.isb_tahap_nontender (
    tahun_anggaran integer,
    kd_klpd character varying,
    kd_satker numeric,
    kd_satker_str character varying,
    kd_lpse integer,
    kd_nontender numeric NOT NULL,
    kd_tahapan numeric,
    nama_tahapan character varying,
    kd_akt numeric NOT NULL,
    nama_akt character varying,
    tgl_awal timestamp without time zone,
    tgl_akhir timestamp without time zone
);


ALTER TABLE public.isb_tahap_nontender OWNER TO postgres;

--
-- Name: isb_tahap_tender; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.isb_tahap_tender (
    tahun_anggaran integer,
    kd_klpd character varying,
    kd_satker numeric,
    kd_satker_str character varying,
    kd_lpse integer,
    kd_tender numeric NOT NULL,
    kd_tahapan numeric,
    nama_tahapan character varying,
    kd_akt numeric NOT NULL,
    nama_akt character varying,
    tgl_awal timestamp without time zone,
    tgl_akhir timestamp without time zone
);


ALTER TABLE public.isb_tahap_tender OWNER TO postgres;

--
-- Name: isb_tender; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.isb_tender (
    tahun_anggaran bigint,
    kd_klpd text,
    nama_klpd text,
    jenis_klpd text,
    kd_satker bigint,
    kd_satker_str text,
    nama_satker text,
    kd_lpse bigint,
    nama_lpse text,
    kd_tender bigint NOT NULL,
    kd_pkt_dce bigint,
    kd_rup text,
    nama_paket text,
    pagu numeric,
    hps numeric,
    sumber_dana text,
    kualifikasi_paket text,
    jenis_pengadaan text,
    mtd_pemilihan text,
    mtd_evaluasi text,
    mtd_kualifikasi text,
    kontrak_pembayaran text,
    status_tender text,
    tanggal_status timestamp without time zone,
    versi_tender bigint,
    ket_ditutup text,
    ket_diulang text,
    tgl_buat_paket timestamp without time zone,
    tgl_kolektif_kolegial timestamp without time zone,
    tgl_pengumuman_tender timestamp without time zone,
    nip_ppk text,
    nama_ppk text,
    nip_pokja text,
    nama_pokja text,
    lokasi_pekerjaan text,
    url_lpse text,
    realisasi jsonb
);


ALTER TABLE public.isb_tender OWNER TO postgres;

--
-- Name: isb_tokodaring; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.isb_tokodaring (
    kd_klpd character varying,
    nama_klpd character varying,
    kd_satker character varying,
    nama_satker character varying,
    order_id character varying NOT NULL,
    order_desc character varying,
    valuasi numeric,
    kategori character varying,
    metode_bayar character varying,
    tanggal_transaksi character varying,
    marketplace character varying,
    nama_merchant character varying,
    jenis_transaksi character varying,
    kota_kab character varying,
    provinsi character varying,
    nama_pemesan character varying,
    status_verif character varying,
    sumber_data character varying,
    status_konfirmasi_ppmse character varying,
    keterangan_ppmse character varying
);


ALTER TABLE public.isb_tokodaring OWNER TO postgres;

--
-- Name: kabupaten; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.kabupaten (
    id bigint NOT NULL,
    prp_id bigint,
    nama text
);


ALTER TABLE public.kabupaten OWNER TO postgres;

--
-- Name: kaji_ulang; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.kaji_ulang (
    id bigint NOT NULL,
    created_at timestamp with time zone,
    updated_at timestamp with time zone,
    deleted_at timestamp with time zone,
    parent_id bigint,
    pkt_id bigint,
    peg_id bigint,
    ppk_id bigint,
    pnt_id bigint,
    pp_id bigint,
    role text,
    dokumen text,
    uraian text,
    dok_id bigint
);


ALTER TABLE public.kaji_ulang OWNER TO postgres;

--
-- Name: kaji_ulang_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.kaji_ulang_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.kaji_ulang_id_seq OWNER TO postgres;

--
-- Name: kaji_ulang_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.kaji_ulang_id_seq OWNED BY public.kaji_ulang.id;


--
-- Name: metode_pengadaan; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.metode_pengadaan (
    id bigint NOT NULL,
    metode text
);


ALTER TABLE public.metode_pengadaan OWNER TO postgres;

--
-- Name: paket; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.paket (
    id bigint NOT NULL,
    created_at timestamp with time zone,
    updated_at timestamp with time zone,
    deleted_at timestamp with time zone,
    nama text,
    created_by bigint,
    updated_by bigint,
    pagu numeric,
    hps numeric,
    ukpbj_id bigint,
    tgl_assign_ukpbj timestamp with time zone,
    pnt_id bigint,
    tgl_assign_pokja timestamp with time zone,
    pp_id bigint,
    tgl_assign_pp timestamp with time zone,
    status bigint,
    kgr_id bigint,
    kode_tender bigint,
    ppk_id bigint,
    rup_id bigint,
    satker_id bigint,
    metode bigint,
    tgl_disetujui timestamp with time zone,
    tgl_ditolak timestamp with time zone,
    alasan_ditolak text,
    prioritas boolean
);


ALTER TABLE public.paket OWNER TO postgres;

--
-- Name: paket_anggaran; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.paket_anggaran (
    id bigint NOT NULL,
    created_at timestamp with time zone,
    updated_at timestamp with time zone,
    deleted_at timestamp with time zone,
    pkt_id bigint,
    ang_id bigint,
    ppk_id bigint,
    rup_id bigint
);


ALTER TABLE public.paket_anggaran OWNER TO postgres;

--
-- Name: paket_anggaran_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.paket_anggaran_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.paket_anggaran_id_seq OWNER TO postgres;

--
-- Name: paket_anggaran_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.paket_anggaran_id_seq OWNED BY public.paket_anggaran.id;


--
-- Name: paket_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.paket_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.paket_id_seq OWNER TO postgres;

--
-- Name: paket_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.paket_id_seq OWNED BY public.paket.id;


--
-- Name: paket_lokasi; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.paket_lokasi (
    id bigint NOT NULL,
    created_at timestamp with time zone,
    updated_at timestamp with time zone,
    deleted_at timestamp with time zone,
    p_kt_id bigint,
    kbp_id bigint,
    lokasi text
);


ALTER TABLE public.paket_lokasi OWNER TO postgres;

--
-- Name: paket_lokasi_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.paket_lokasi_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.paket_lokasi_id_seq OWNER TO postgres;

--
-- Name: paket_lokasi_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.paket_lokasi_id_seq OWNED BY public.paket_lokasi.id;


--
-- Name: paket_satker; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.paket_satker (
    id bigint NOT NULL,
    created_at timestamp with time zone,
    updated_at timestamp with time zone,
    deleted_at timestamp with time zone,
    pkt_id bigint,
    stk_id bigint,
    rup_id bigint
);


ALTER TABLE public.paket_satker OWNER TO postgres;

--
-- Name: paket_satker_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.paket_satker_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.paket_satker_id_seq OWNER TO postgres;

--
-- Name: paket_satker_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.paket_satker_id_seq OWNED BY public.paket_satker.id;


--
-- Name: paket_sirup; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.paket_sirup (
    id bigint NOT NULL,
    nama text,
    paket_lokasi_json jsonb,
    volume text,
    keterangan text,
    spesifikasi text,
    is_tkdn boolean,
    is_pradipa boolean,
    paket_anggaran_json jsonb,
    pagu numeric,
    paket_jenis_json jsonb,
    metode_pengadaan bigint,
    tanggal_kebutuhan timestamp with time zone,
    tanggal_awal_pengadaan timestamp with time zone,
    tanggal_akhir_pengadaan timestamp with time zone,
    tanggal_awal_pekerjaan timestamp with time zone,
    tanggal_akhir_pekerjaan timestamp with time zone,
    tanggal_pengumuman timestamp with time zone,
    id_swakelola bigint,
    id_ppk bigint,
    umkm boolean,
    kode_kldi text,
    id_satker bigint,
    encrypted_username_ppk text,
    paket_aktif boolean,
    paket_terhapus boolean,
    status_paket bigint,
    paket_terumumkan boolean,
    tahun bigint,
    jenis_paket bigint
);


ALTER TABLE public.paket_sirup OWNER TO postgres;

--
-- Name: panitia; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.panitia (
    id bigint NOT NULL,
    created_at timestamp with time zone,
    updated_at timestamp with time zone,
    deleted_at timestamp with time zone,
    nama text NOT NULL,
    tahun bigint NOT NULL
);


ALTER TABLE public.panitia OWNER TO postgres;

--
-- Name: panitia_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.panitia_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.panitia_id_seq OWNER TO postgres;

--
-- Name: panitia_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.panitia_id_seq OWNED BY public.panitia.id;


--
-- Name: seq_pegawai; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.seq_pegawai
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.seq_pegawai OWNER TO postgres;

--
-- Name: pegawai; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.pegawai (
    peg_nip text NOT NULL,
    created_at timestamp with time zone,
    updated_at timestamp with time zone,
    peg_nama text NOT NULL,
    peg_alamat text,
    peg_telepon text,
    peg_mobile text,
    peg_email text,
    peg_golongan text,
    peg_pangkat text,
    peg_jabatan text,
    peg_isactive bigint,
    peg_namauser text,
    peg_no_sk text,
    peg_masa_berlaku timestamp without time zone,
    id numeric(19,0) DEFAULT nextval('public.seq_pegawai'::regclass) NOT NULL,
    agc_id bigint,
    peg_no_pbj text,
    passw text,
    reset_password text,
    ukpbj_id bigint,
    satker_id numeric(18,2) DEFAULT NULL::numeric,
    peg_nik text,
    usrgroup text,
    peg_tipe_sertifikat bigint,
    last_change_passw timestamp without time zone,
    deleted_at timestamp without time zone,
    created_by bigint,
    updated_by bigint,
    peg_status bigint,
    peg_catatan text,
    tgl_approve timestamp without time zone,
    tgl_reject timestamp without time zone
);


ALTER TABLE public.pegawai OWNER TO postgres;

--
-- Name: COLUMN pegawai.peg_status; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.pegawai.peg_status IS '0 = verifikasi, 1 = aktif, 2 = revisi';


--
-- Name: provinsi; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.provinsi (
    id bigint NOT NULL,
    nama text
);


ALTER TABLE public.provinsi OWNER TO postgres;

--
-- Name: reviu; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.reviu (
    id bigint NOT NULL,
    created_at timestamp with time zone,
    updated_at timestamp with time zone,
    deleted_at timestamp with time zone,
    bidang text,
    content text,
    opsi1 text,
    opsi2 text
);


ALTER TABLE public.reviu OWNER TO postgres;

--
-- Name: reviu_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.reviu_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.reviu_id_seq OWNER TO postgres;

--
-- Name: reviu_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.reviu_id_seq OWNED BY public.reviu.id;


--
-- Name: reviu_paket; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.reviu_paket (
    id bigint NOT NULL,
    created_at timestamp with time zone,
    updated_at timestamp with time zone,
    deleted_at timestamp with time zone,
    pkt_id bigint,
    rev_id bigint,
    catatan_khusus text,
    keterangan text,
    peg_id bigint
);


ALTER TABLE public.reviu_paket OWNER TO postgres;

--
-- Name: reviu_paket_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.reviu_paket_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.reviu_paket_id_seq OWNER TO postgres;

--
-- Name: reviu_paket_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.reviu_paket_id_seq OWNED BY public.reviu_paket.id;


--
-- Name: satker; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.satker (
    id bigint NOT NULL,
    id_satker text,
    id_kldi text,
    is_deleted boolean,
    nama text,
    auditupdate timestamp with time zone,
    tahun_aktif text,
    blu boolean,
    jenis_satker_id bigint
);


ALTER TABLE public.satker OWNER TO postgres;

--
-- Name: seq_ukpbj; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.seq_ukpbj
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.seq_ukpbj OWNER TO postgres;

--
-- Name: swakelola_sirup; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.swakelola_sirup (
    id bigint NOT NULL,
    nama text,
    paket_lokasi_json jsonb,
    lls_volume text,
    keterangan text,
    is_pradiap boolean,
    paket_anggaran_json jsonb,
    jumlah_pagu numeric,
    tanggal_awal_pekerjaan timestamp with time zone,
    tanggal_akhir_pekerjaan timestamp with time zone,
    tanggal_pengumuman timestamp with time zone,
    id_ppk bigint,
    kode_kldi text,
    id_satker bigint,
    tipe_swakelola bigint,
    satker_lain bigint,
    nama_satker_lain text,
    kld_lain text,
    nama_kld_lain text,
    aktif boolean,
    umumkan boolean,
    is_deleted boolean,
    status bigint,
    tahun bigint
);


ALTER TABLE public.swakelola_sirup OWNER TO postgres;

--
-- Name: ukpbj; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.ukpbj (
    id numeric(19,0) DEFAULT nextval('public.seq_ukpbj'::regclass) NOT NULL,
    agc_id numeric(19,0),
    peg_id numeric(19,0),
    updated_by character varying(100) DEFAULT 'ADMIN'::character varying,
    updated_at timestamp without time zone DEFAULT now(),
    nama text NOT NULL,
    alamat text NOT NULL,
    telepon text,
    fax text,
    tgl_daftar timestamp without time zone,
    is_active boolean,
    deleted_at timestamp without time zone,
    created_at timestamp without time zone
);


ALTER TABLE public.ukpbj OWNER TO postgres;

--
-- Name: COLUMN ukpbj.peg_id; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.ukpbj.peg_id IS 'ID pegawai admin ukpbj';


--
-- Name: anggarans id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.anggarans ALTER COLUMN id SET DEFAULT nextval('public.anggarans_id_seq'::regclass);


--
-- Name: anggota_panitia id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.anggota_panitia ALTER COLUMN id SET DEFAULT nextval('public.anggota_panitia_id_seq'::regclass);


--
-- Name: berita_acara id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.berita_acara ALTER COLUMN id SET DEFAULT nextval('public.berita_acara_id_seq'::regclass);


--
-- Name: buku_tamu id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.buku_tamu ALTER COLUMN id SET DEFAULT nextval('public.buku_tamu_id_seq'::regclass);


--
-- Name: checklist id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.checklist ALTER COLUMN id SET DEFAULT nextval('public.checklist_id_seq'::regclass);


--
-- Name: checklist_paket id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.checklist_paket ALTER COLUMN id SET DEFAULT nextval('public.checklist_paket_id_seq'::regclass);


--
-- Name: checklist_paket_history id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.checklist_paket_history ALTER COLUMN id SET DEFAULT nextval('public.checklist_paket_history_id_seq'::regclass);


--
-- Name: dok_persiapan id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.dok_persiapan ALTER COLUMN id SET DEFAULT nextval('public.dok_persiapan_id_seq'::regclass);


--
-- Name: dok_templates id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.dok_templates ALTER COLUMN id SET DEFAULT nextval('public.dok_templates_id_seq'::regclass);


--
-- Name: feedback id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.feedback ALTER COLUMN id SET DEFAULT nextval('public.feedback_id_seq'::regclass);


--
-- Name: hak_akses id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.hak_akses ALTER COLUMN id SET DEFAULT nextval('public.hak_akses_id_seq'::regclass);


--
-- Name: inbox id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.inbox ALTER COLUMN id SET DEFAULT nextval('public.inboxes_id_seq'::regclass);


--
-- Name: isb_purchase id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.isb_purchase ALTER COLUMN id SET DEFAULT nextval('public.isb_purchase_id_seq'::regclass);


--
-- Name: isb_struktur_anggaran id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.isb_struktur_anggaran ALTER COLUMN id SET DEFAULT nextval('public.isb_struktur_anggaran_id_seq'::regclass);


--
-- Name: kaji_ulang id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.kaji_ulang ALTER COLUMN id SET DEFAULT nextval('public.kaji_ulang_id_seq'::regclass);


--
-- Name: paket id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.paket ALTER COLUMN id SET DEFAULT nextval('public.paket_id_seq'::regclass);


--
-- Name: paket_anggaran id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.paket_anggaran ALTER COLUMN id SET DEFAULT nextval('public.paket_anggaran_id_seq'::regclass);


--
-- Name: paket_lokasi id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.paket_lokasi ALTER COLUMN id SET DEFAULT nextval('public.paket_lokasi_id_seq'::regclass);


--
-- Name: paket_satker id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.paket_satker ALTER COLUMN id SET DEFAULT nextval('public.paket_satker_id_seq'::regclass);


--
-- Name: panitia id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.panitia ALTER COLUMN id SET DEFAULT nextval('public.panitia_id_seq'::regclass);


--
-- Name: reviu id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.reviu ALTER COLUMN id SET DEFAULT nextval('public.reviu_id_seq'::regclass);


--
-- Name: reviu_paket id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.reviu_paket ALTER COLUMN id SET DEFAULT nextval('public.reviu_paket_id_seq'::regclass);


--
-- Name: agency agency_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.agency
    ADD CONSTRAINT agency_pkey PRIMARY KEY (id);


--
-- Name: anggarans anggarans_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.anggarans
    ADD CONSTRAINT anggarans_pkey PRIMARY KEY (id);


--
-- Name: anggota_panitia anggota_panitia_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.anggota_panitia
    ADD CONSTRAINT anggota_panitia_pkey PRIMARY KEY (id);


--
-- Name: berita_acara berita_acara_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.berita_acara
    ADD CONSTRAINT berita_acara_pkey PRIMARY KEY (id);


--
-- Name: buku_tamu buku_tamu_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.buku_tamu
    ADD CONSTRAINT buku_tamu_pkey PRIMARY KEY (id);


--
-- Name: checklist_paket_history checklist_paket_history_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.checklist_paket_history
    ADD CONSTRAINT checklist_paket_history_pkey PRIMARY KEY (id);


--
-- Name: checklist_paket checklist_paket_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.checklist_paket
    ADD CONSTRAINT checklist_paket_pkey PRIMARY KEY (id);


--
-- Name: checklist checklist_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.checklist
    ADD CONSTRAINT checklist_pkey PRIMARY KEY (id);


--
-- Name: document documents_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.document
    ADD CONSTRAINT documents_pkey PRIMARY KEY (id, versi);


--
-- Name: dok_persiapan dok_persiapan_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.dok_persiapan
    ADD CONSTRAINT dok_persiapan_pkey PRIMARY KEY (id);


--
-- Name: dok_templates dok_templates_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.dok_templates
    ADD CONSTRAINT dok_templates_pkey PRIMARY KEY (id);


--
-- Name: feedback feedback_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.feedback
    ADD CONSTRAINT feedback_pkey PRIMARY KEY (id);


--
-- Name: hak_akses hak_akses_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.hak_akses
    ADD CONSTRAINT hak_akses_pkey PRIMARY KEY (id);


--
-- Name: inbox inboxes_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.inbox
    ADD CONSTRAINT inboxes_pkey PRIMARY KEY (id);


--
-- Name: isb_tahap_tender isb_jadwal_tender_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.isb_tahap_tender
    ADD CONSTRAINT isb_jadwal_tender_pkey PRIMARY KEY (kd_tender, kd_akt);


--
-- Name: isb_kabupaten isb_kabupaten_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.isb_kabupaten
    ADD CONSTRAINT isb_kabupaten_pkey PRIMARY KEY (kd_kabupaten);


--
-- Name: isb_nontender isb_nontender_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.isb_nontender
    ADD CONSTRAINT isb_nontender_pkey PRIMARY KEY (kd_nontender);


--
-- Name: isb_pencatatan isb_pencatatan_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.isb_pencatatan
    ADD CONSTRAINT isb_pencatatan_pkey PRIMARY KEY (kd_nontender_pct);


--
-- Name: isb_provinsi isb_provinsi_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.isb_provinsi
    ADD CONSTRAINT isb_provinsi_pkey PRIMARY KEY (kd_provinsi);


--
-- Name: isb_purchase isb_purchase_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.isb_purchase
    ADD CONSTRAINT isb_purchase_pkey PRIMARY KEY (id);


--
-- Name: isb_rup_kegiatan isb_rup_kegiatan_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.isb_rup_kegiatan
    ADD CONSTRAINT isb_rup_kegiatan_pkey PRIMARY KEY (kd_kegiatan);


--
-- Name: isb_rup_paket isb_rup_paket_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.isb_rup_paket
    ADD CONSTRAINT isb_rup_paket_pkey PRIMARY KEY (kd_rup);


--
-- Name: isb_rup_program isb_rup_program_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.isb_rup_program
    ADD CONSTRAINT isb_rup_program_pkey PRIMARY KEY (kd_program);


--
-- Name: isb_rup_swakelola isb_rup_swakelola_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.isb_rup_swakelola
    ADD CONSTRAINT isb_rup_swakelola_pkey PRIMARY KEY (kd_rup);


--
-- Name: isb_satker isb_satker_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.isb_satker
    ADD CONSTRAINT isb_satker_pkey PRIMARY KEY (kd_satker);


--
-- Name: isb_struktur_anggaran isb_struktur_anggaran_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.isb_struktur_anggaran
    ADD CONSTRAINT isb_struktur_anggaran_pkey PRIMARY KEY (id);


--
-- Name: isb_swakelola isb_swakelola_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.isb_swakelola
    ADD CONSTRAINT isb_swakelola_pkey PRIMARY KEY (kd_swakelola_pct);


--
-- Name: isb_tahap_nontender isb_tahap_nontender_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.isb_tahap_nontender
    ADD CONSTRAINT isb_tahap_nontender_pkey PRIMARY KEY (kd_nontender, kd_akt);


--
-- Name: isb_tender isb_tender_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.isb_tender
    ADD CONSTRAINT isb_tender_pkey PRIMARY KEY (kd_tender);


--
-- Name: isb_tokodaring isb_tokodarin_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.isb_tokodaring
    ADD CONSTRAINT isb_tokodarin_pkey PRIMARY KEY (order_id);


--
-- Name: kabupaten kabupaten_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.kabupaten
    ADD CONSTRAINT kabupaten_pkey PRIMARY KEY (id);


--
-- Name: kaji_ulang kaji_ulang_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.kaji_ulang
    ADD CONSTRAINT kaji_ulang_pkey PRIMARY KEY (id);


--
-- Name: metode_pengadaan metode_pengadaan_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.metode_pengadaan
    ADD CONSTRAINT metode_pengadaan_pkey PRIMARY KEY (id);


--
-- Name: paket_anggaran paket_anggaran_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.paket_anggaran
    ADD CONSTRAINT paket_anggaran_pkey PRIMARY KEY (id);


--
-- Name: paket_lokasi paket_lokasi_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.paket_lokasi
    ADD CONSTRAINT paket_lokasi_pkey PRIMARY KEY (id);


--
-- Name: paket paket_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.paket
    ADD CONSTRAINT paket_pkey PRIMARY KEY (id);


--
-- Name: paket_satker paket_satker_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.paket_satker
    ADD CONSTRAINT paket_satker_pkey PRIMARY KEY (id);


--
-- Name: paket_sirup paket_sirup_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.paket_sirup
    ADD CONSTRAINT paket_sirup_pkey PRIMARY KEY (id);


--
-- Name: panitia panitia_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.panitia
    ADD CONSTRAINT panitia_pkey PRIMARY KEY (id);


--
-- Name: pegawai pegawai_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.pegawai
    ADD CONSTRAINT pegawai_pkey PRIMARY KEY (id);


--
-- Name: provinsi provinsi_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.provinsi
    ADD CONSTRAINT provinsi_pkey PRIMARY KEY (id);


--
-- Name: reviu_paket reviu_paket_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.reviu_paket
    ADD CONSTRAINT reviu_paket_pkey PRIMARY KEY (id);


--
-- Name: reviu reviu_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.reviu
    ADD CONSTRAINT reviu_pkey PRIMARY KEY (id);


--
-- Name: satker satker_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.satker
    ADD CONSTRAINT satker_pkey PRIMARY KEY (id);


--
-- Name: swakelola_sirup swakelola_sirup_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.swakelola_sirup
    ADD CONSTRAINT swakelola_sirup_pkey PRIMARY KEY (id);


--
-- Name: ukpbj ukpbj_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.ukpbj
    ADD CONSTRAINT ukpbj_pkey PRIMARY KEY (id);


--
-- Name: idx_anggarans_deleted_at; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX idx_anggarans_deleted_at ON public.anggarans USING btree (deleted_at);


--
-- Name: idx_anggota_panitia_deleted_at; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX idx_anggota_panitia_deleted_at ON public.anggota_panitia USING btree (deleted_at);


--
-- Name: idx_berita_acara_deleted_at; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX idx_berita_acara_deleted_at ON public.berita_acara USING btree (deleted_at);


--
-- Name: idx_buku_tamu_deleted_at; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX idx_buku_tamu_deleted_at ON public.buku_tamu USING btree (deleted_at);


--
-- Name: idx_checklist_deleted_at; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX idx_checklist_deleted_at ON public.checklist USING btree (deleted_at);


--
-- Name: idx_checklist_paket_deleted_at; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX idx_checklist_paket_deleted_at ON public.checklist_paket USING btree (deleted_at);


--
-- Name: idx_checklist_paket_history_deleted_at; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX idx_checklist_paket_history_deleted_at ON public.checklist_paket_history USING btree (deleted_at);


--
-- Name: idx_dok_persiapan_deleted_at; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX idx_dok_persiapan_deleted_at ON public.dok_persiapan USING btree (deleted_at);


--
-- Name: idx_dok_templates_deleted_at; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX idx_dok_templates_deleted_at ON public.dok_templates USING btree (deleted_at);


--
-- Name: idx_feedback_deleted_at; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX idx_feedback_deleted_at ON public.feedback USING btree (deleted_at);


--
-- Name: idx_hak_akses_deleted_at; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX idx_hak_akses_deleted_at ON public.hak_akses USING btree (deleted_at);


--
-- Name: idx_inbox_deleted_at; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX idx_inbox_deleted_at ON public.inbox USING btree (deleted_at);


--
-- Name: idx_inboxes_deleted_at; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX idx_inboxes_deleted_at ON public.inbox USING btree (deleted_at);


--
-- Name: idx_kaji_ulang_deleted_at; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX idx_kaji_ulang_deleted_at ON public.kaji_ulang USING btree (deleted_at);


--
-- Name: idx_paket_anggaran_deleted_at; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX idx_paket_anggaran_deleted_at ON public.paket_anggaran USING btree (deleted_at);


--
-- Name: idx_paket_deleted_at; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX idx_paket_deleted_at ON public.paket USING btree (deleted_at);


--
-- Name: idx_paket_lokasi_deleted_at; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX idx_paket_lokasi_deleted_at ON public.paket_lokasi USING btree (deleted_at);


--
-- Name: idx_paket_satker_deleted_at; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX idx_paket_satker_deleted_at ON public.paket_satker USING btree (deleted_at);


--
-- Name: idx_panitia_deleted_at; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX idx_panitia_deleted_at ON public.panitia USING btree (deleted_at);


--
-- Name: idx_pegawai_deleted_at; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX idx_pegawai_deleted_at ON public.pegawai USING btree (deleted_at);


--
-- Name: idx_reviu_deleted_at; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX idx_reviu_deleted_at ON public.reviu USING btree (deleted_at);


--
-- Name: idx_reviu_paket_deleted_at; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX idx_reviu_paket_deleted_at ON public.reviu_paket USING btree (deleted_at);


--
-- Name: ind_blob_table_nama_file; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX ind_blob_table_nama_file ON public.document USING btree (filename);


--
-- Name: ind_pegawai_deleted; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX ind_pegawai_deleted ON public.pegawai USING btree (deleted_at);


--
-- Name: ind_pegawai_namauser; Type: INDEX; Schema: public; Owner: postgres
--

CREATE UNIQUE INDEX ind_pegawai_namauser ON public.pegawai USING btree (peg_namauser);


--
-- Name: ind_ukpbj_nama; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX ind_ukpbj_nama ON public.ukpbj USING btree (nama);


--
-- Name: pegawai pegawai_agc_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.pegawai
    ADD CONSTRAINT pegawai_agc_id_fkey FOREIGN KEY (agc_id) REFERENCES public.agency(id);


--
-- Name: ukpbj ukpbj_peg_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.ukpbj
    ADD CONSTRAINT ukpbj_peg_id_fkey FOREIGN KEY (peg_id) REFERENCES public.pegawai(id) ON UPDATE CASCADE ON DELETE RESTRICT;


--
-- PostgreSQL database dump complete
--

