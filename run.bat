@echo off
echo ===================================================
echo Starting admin.exe...
echo ===================================================

if not exist admin.exe (
    echo [ERROR] admin.exe not found! Please run build.bat first.
    echo.
    pause
    exit /b 1
)

echo Starting application...
echo You can access the Admin panel at: http://localhost:8081/admin
echo Press Ctrl+C in this window to stop the server.
echo ===================================================
echo.

admin.exe

pause
