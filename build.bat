@echo off
echo ===================================================
echo Building admin.exe for Windows (amd64)...
echo ===================================================

:: Set Go environment variables for Windows build
set CGO_ENABLED=0
set GOOS=windows
set GOARCH=amd64

:: Perform the build
go build -o admin.exe

if %ERRORLEVEL% equ 0 (
    echo.
    echo ===================================================
    echo BUILD SUCCESSFUL! Created admin.exe
    echo ===================================================
) else (
    echo.
    echo ===================================================
    echo BUILD FAILED! Please check the errors above.
    echo ===================================================
)

pause