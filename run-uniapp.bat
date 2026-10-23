@echo off
echo ========================================================
echo   WallArt Multi-Terminal Launcher (Local Node/Go Mode)
echo ========================================================
echo.
echo [1/2] Starting Backend (Go API on :8080)...
start "WallArt Backend" cmd /c "cd /d %~dp0backend && .\wallpaper-server.exe"

timeout /t 2 /nobreak >nul

echo [2/2] Starting Uni-app Frontend (H5/Web on :5173)...
start "WallArt Uni-App" cmd /c "cd /d %~dp0uni-app && npm run dev:h5"

echo.
echo ========================================================
echo Services started!
echo.
echo  * Frontend URL: http://localhost:5173
echo  * Backend API:  http://localhost:8080/api/health
echo.
echo  Tip: To run with Docker instead, double-click docker-start.bat
echo ========================================================
echo.
pause
