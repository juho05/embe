@echo off
echo Downloading latest installer...
Powershell.exe -Command "iwr -useb https://raw.githubusercontent.com/juho05/embe/main/install.ps1 | iex"
pause
