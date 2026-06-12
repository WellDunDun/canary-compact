@echo off
setlocal
set "HOOK=%~dp0windows-amd64\canary-compact-hook.exe"
if not exist "%HOOK%" (
  echo canary-compact: missing helper binary: %HOOK% 1>&2
  exit /b 1
)
"%HOOK%" %*
