CREATE OR REPLACE FUNCTION public.jenis_pengadaan(id integer)
 RETURNS character varying
 LANGUAGE plpgsql
AS $function$
DECLARE
  result VARCHAR ;
BEGIN
    CASE id
    	WHEN 0 THEN result = 'Belum Ditentukan';
    	WHEN 1 THEN result = 'Barang';
    	WHEN 2 THEN result = 'Pekerjaan Konstruksi';
    	WHEN 3 THEN result = 'Jasa Konsultansi';
    	WHEN 4 THEN result = 'Jasa Lainnya';
    	WHEN 5 THEN result = 'Terintegrasi';
    END CASE;
    RETURN result;
END;
$function$

CREATE OR REPLACE FUNCTION public.metode(id integer)
 RETURNS character varying
 LANGUAGE plpgsql
AS $function$
DECLARE
  result VARCHAR ;
BEGIN
    CASE id
    	WHEN 0 THEN result = 'Belum Ditentukan';
    	WHEN 1 THEN result = 'Lelang Umum';
    	WHEN 2 THEN result = 'Lelang Sederhana';
    	WHEN 3 THEN result = 'Lelang Terbatas';
    	WHEN 4 THEN result = 'Seleksi Umum';
    	WHEN 5 THEN result = 'Seleksi Sederhana';
		WHEN 6 THEN result = 'Pemilihan Langsung';
    	WHEN 7 THEN result = 'Penunjukan Langsung';
    	WHEN 8 THEN result = 'Pengadaan Langsung';
    	WHEN 9 THEN result = 'e-Purchasing';
    	WHEN 10 THEN result = 'Sayembara';
    	WHEN 11 THEN result = 'Kontes';
		WHEN 12 THEN result = 'Lelang Cepat';
		WHEN 13 THEN result = 'Tender';
		WHEN 14 THEN result = 'Tender Cepat';
		WHEN 15 THEN result = 'Seleksi';
		WHEN 16 THEN result = 'Dikecualikan';
    END CASE;
    RETURN result;
END;
$function$
