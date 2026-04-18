@echo off
echo Building folder_star.exe...
go build -o folder_star.exe main.go
if %ERRORLEVEL% EQU 0 (
    echo Build successful: folder_star.exe
) else (
    echo Build failed!
)
