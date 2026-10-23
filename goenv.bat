@echo off
set "PATH=D:\Program Files\Go\bin;%PATH%"
set "GOPROXY=https://goproxy.io,direct"
set "GONOSUMDB=*"
set "GOFLAGS=-mod=mod"
set "GOPATH=C:\gomod"
set "GOMODCACHE=C:\gomod\pkg\mod"
set "GOROOT=D:\goproject\go"
cd /d e:\WALLPAPER\backend
go get github.com/gin-gonic/gin
go get github.com/gin-contrib/cors
go get gorm.io/gorm
go get gorm.io/driver/sqlite
go get github.com/golang-jwt/jwt/v4
go get golang.org/x/crypto/bcrypt
go get github.com/google/uuid
echo ALL DONE
