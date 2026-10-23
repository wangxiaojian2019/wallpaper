@echo off
echo ========================================================
echo   WallArt Docker One-Click Launcher
echo ========================================================
echo.
echo Starting all services (Backend + Uni-App Client + Admin)...
docker compose up -d --build
echo.
echo ========================================================
echo All services are running!
echo.
echo  * Client Web/H5 : http://localhost:5173  (or http://localhost)
echo  * Admin Panel   : http://localhost:5174
echo  * Backend API   : http://localhost:8080/api/health
echo ========================================================
echo.
pause
