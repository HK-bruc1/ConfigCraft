@echo off
echo Building DHF Configuration Manager with TDM-GCC...
echo Current directory: %CD%

REM Change to project root directory
cd /d "%~dp0.."
echo Project directory: %CD%

REM Set TDM-GCC path first
set "ORIGINAL_PATH=%PATH%"
set "PATH=E:\WindowsAppFile\TDM-GCC-64\bin;%PATH%"

REM Check compiler version
echo Current GCC version:
gcc --version

REM Set proxy settings (if needed)
set https_proxy=http://127.0.0.1:7897
set http_proxy=http://127.0.0.1:7897
set all_proxy=socks5://127.0.0.1:7897
echo Proxy settings configured

REM Set build parameters
set CGO_ENABLED=1
set GOOS=windows
set GOARCH=amd64
set GOMAXPROCS=4

echo.
echo Starting build...
go build -v -ldflags "-s -w -H windowsgui" -o build\dhf-config-manager.exe main.go

if %ERRORLEVEL% == 0 (
    echo.
    echo Build successful!
    if exist "build\dhf-config-manager.exe" (
        echo File created: build\dhf-config-manager.exe
        dir "build\dhf-config-manager.exe"
    )
) else (
    echo.
    echo Build failed with error code: %ERRORLEVEL%
)

REM Restore original PATH
set "PATH=%ORIGINAL_PATH%"
echo.
echo Press any key to exit...
pause >nul