Aplikasi ini digunakan untuk sync data inaproc di backend server. Project ini merupakan project CLI menggunakan golang

### requirement development
- golang minimum 1.25
- postgresql minimum 15


### requirement modul
- go get github.com/joho/godotenv
- go get golang.org/x/exp
- go get gorm.io/driver/postgres
- go get gorm.io/gorm


### Deploly on linux
1. chmod +x sync-inaproc
2. edit .env
3. Open the Crontab Editor
4. 0 * * * * sync-inaproc >> app.log 2>&1


### how to build
linux 64 bit :
  GOOS=linux GOARCH=amd64 go build -o schduler main.go
windows :
  $Env:GOOS = "linux"; $Env:GOARCH = "amd64"; go build -o arsip main.go
