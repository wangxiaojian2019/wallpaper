@echo off
chcp 65001 >nul
echo ========================================================
echo        壁纸商城系统一键启动器 (WallArt Store)
echo ========================================================
echo.
echo [1/2] 正在启动后端服务 (Go + SQLite on :8080)...
start "WallArt Backend" cmd /c "cd /d %~dp0backend && .\wallpaper-server.exe"

timeout /t 2 /nobreak >nul

echo [2/2] 正在启动前端开发服务 (Vite on :5173)...
start "WallArt Frontend" cmd /c "cd /d %~dp0frontend && npm run dev"

echo.
echo ========================================================
echo 系统启动完成！
echo.
echo 前端访问地址: http://localhost:5173
echo 后端接口地址: http://localhost:8080/api
echo.
echo 默认管理员账号: admin@wallpaper.com
echo 默认管理员密码: admin123456
echo ========================================================
echo.
pause
