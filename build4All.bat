@echo off
setlocal enabledelayedexpansion

set SRC=porkbunDDNS.go
set OUTDIR=builds
mkdir %OUTDIR%

REM List of all supported OS
set GOOS_LIST=android darwin freebsd linux netbsd openbsd solaris windows

REM List of all supported architectures
set GOARCH_LIST=386 amd64 amd64p32 arm arm64 arm64be armbe

for %%o in (%GOOS_LIST%) do (
    for %%a in (%GOARCH_LIST%) do (
        set OUTFILE=%OUTDIR%\porkbunDDNS-%%o-%%a
        if "%%o"=="windows" set OUTFILE=!OUTFILE!.exe
        echo Building %SRC% for %%o/%%a ...
        REM Inline environment variables for this build only
        cmd /C "set GOOS=%%o&& set GOARCH=%%a&& go build -o !OUTFILE! %SRC%"
    )
)

echo All builds attempted!
pause
