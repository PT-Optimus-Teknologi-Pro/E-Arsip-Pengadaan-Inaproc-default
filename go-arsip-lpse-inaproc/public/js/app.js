const rupiah = (number)=>{
    return new Intl.NumberFormat("id-ID", {
      style: "currency",
      currency: "IDR"
    }).format(number);
}

const statusPaket = (number) => {
  if (number == 0) {
    return "Draft";
  }
  if (number == 1) {
    return  "Pengajuan";
  }
  if (number == 2) {
    return "Disetujui";
  }
  if (number == 3) {
    return "Tolak";
  }
  if (number == 4) {
    return "Kaji Ulang";
  }
  if (number == 5) {
    return "Proses";
  }
  if (number == 6) {
    return "Selesai";
  }
}

const metode = (number) => {
  if (number == 0) {
    return "Belum Ditentukan";
  }
  if (number == 1) {
    return  "Lelang Umum";
  }
  if (number == 2) {
    return "Lelang Sederhana";
  }
  if (number == 3) {
    return "Lelang Terbatas";
  }
  if (number == 4) {
    return "Seleksi Umum";
  }
  if (number == 5) {
    return "Seleksi Sederhana";
  }
  if (number == 6) {
    return "Pemilihan Langsung";
  }
  if (number == 7) {
    return "Penunjukan Langsung";
  }
  if (number == 8) {
    return "Pengadaan Langsung";
  }
  if (number == 9) {
    return "e-Purchasing";
  }
  if (number == 10) {
    return "Sayembara";
  }
  if (number == 11) {
    return "Kontes";
  }
  if (number == 12) {
    return "Lelang Cepat";
  }
  if (number == 13) {
    return "Tender";
  }
  if (number == 14) {
    return "Tender Cepat";
  }
  if (number == 15) {
    return "Seleksi";
  }
  if (number == 16) {
    return "Dikecualikan";
  }
}

const feedback = (number) => {
  if (number == 1) {
    return "Tidak Puas";
  }
  if (number == 2) {
    return "Puas";
  }
  if (number == 3) {
    return "Sangat Puas";
  }
}

const statusPerubahan = (number) => {
  if (number == 0) {
    return "Draft";
  }
  if (number == 1) {
    return  "Proses";
  }
  if (number == 2) {
    return "Selesai";
  }
}

const jenisPengadaan = (number) => {
  if (number == 0) {
    return "Belum Ditentukan";
  }
  if (number == 1) {
    return  "Barang";
  }
  if (number == 2) {
    return "Pekerjaan Konstruksi";
  }
  if (number == 3) {
    return "Jasa Konsultansi";
  }
  if (number == 4) {
    return "Jasa Lainnya";
  }
  if (number == 5) {
    return "Terintegrasi";
  }
}
