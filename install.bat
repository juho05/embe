@echo off
echo Downloading latest installer...
Powershell.exe -Command "iwr -useb https://raw.githubusercontent.com/Bananenpro/embe/main/install.ps1 | iex"
pause
