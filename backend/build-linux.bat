@echo off
set GOOS=linux
set GOARCH=amd64
set CGO_ENABLED=0
set GOCACHE=d:\WALLPAPER\backend\.gocache
set GOTMPDIR=d:\WALLPAPER\backend\.gotmp
go build -ldflags "-s -w" -o wallpaper-server-linux main.go
if %ERRORLEVEL% EQU 0 (
    echo [SUCCESS] wallpaper-server-linux compiled successfully!
) else (
    echo [FAILED] Compilation failed with error %ERRORLEVEL%
)
