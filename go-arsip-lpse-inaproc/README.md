# go-arsip-lpse
arsip dokumen lpse menggunakan framework gofiber

# requirement dev
1. golang
2. gofiber
3. postgresql
4. dokumen design

   [DESAIN MODUL APLIKASI BANTU MODUL 1-2 lanjut 18 Juli jam 12.22 (1).xlsx](https://github.com/user-attachments/files/18483987/DESAIN.MODUL.APLIKASI.BANTU.MODUL.1-2.lanjut.18.Juli.jam.12.22.1.xlsx)

   [FORMULIR PERMOHONAN PAKPA & PPK.docx](https://github.com/user-attachments/files/18483990/FORMULIR.PERMOHONAN.PAKPA.PPK.docx)

5. ITKP
   Ref : https://onedrive.live.com/personal/d97a4ab3989db227/_layouts/15/Doc.aspx?sourcedoc=%7Bfe00b8e2-6328-4d87-8c0c-9f59160ef227%7D&action=default&redeem=aHR0cHM6Ly8xZHJ2Lm1zL3gvYy9kOTdhNGFiMzk4OWRiMjI3L0VlSzRBUDRvWTRkTmpBeWZXUllPOGljQkhFc1ZZWVA3ZEhHdkVoa3gtZE5mNkE_ZT1rd2RNVGI&slrid=2c8cbfa1-e01e-a000-0106-9eb321d213b3&originalPath=aHR0cHM6Ly8xZHJ2Lm1zL3gvYy9kOTdhNGFiMzk4OWRiMjI3L0VlSzRBUDRvWTRkTmpBeWZXUllPOGljQkhFc1ZZWVA3ZEhHdkVoa3gtZE5mNkE_cnRpbWU9WWNYa042WGszVWc&CID=249147d4-eaab-4afb-a52a-ccb6717a7b03&_SRM=0:G:66

### isb credential

user : requester-kota-banjarmasin
pass : 5284GCEbcgZQ~{={

u: requester-prov-malut 
p: 6579QPPHMfDg({@=

### how to build
linux 64 bit

1. GOOS=linux GOARCH=amd64 go build -o arsip-inaproc main.go
  on windows :
    $Env:GOOS = "linux"; $Env:GOARCH = "amd64"; go build -o arsip main.go
2. mkdir build
3. cp -r views build
4. cp -r public build
5. cp ./arsip build
6. zip -r arsip.zip build
