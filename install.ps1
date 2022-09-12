$EmbeDir = "$HOME\AppData\Local\Programs\embe"
if (!($EmbeDir | Test-Path)) {
	New-Item -ItemType "directory" -Path $EmbeDir
	[System.Environment]::SetEnvironmentVariable("PATH",[System.Environment]::GetEnvironmentVariable("PATH","USER") + ";" + $EmbeDir,"USER")
} else {
	rm -r $EmbeDir
	New-Item -ItemType "directory" -Path $EmbeDir
}

$TempDir = [System.IO.Path]::GetTempPath()
cd $TempDir

Write-Host "Installing embe..."
Invoke-WebRequest -Uri https://github.com/Bananenpro/embe/releases/latest/download/embe-windows-amd64.zip -OutFile embe.zip
Expand-Archive -LiteralPath embe.zip -DestinationPath $EmbeDir
Rename-Item -Path $EmbeDir\README.md -NewName README-embe.md
Rename-Item -Path $EmbeDir\LICENSE -NewName LICENSE-embe
rm embe.zip

Write-Host "Installing embe-ls..."
Invoke-WebRequest -Uri https://github.com/Bananenpro/embe-ls/releases/latest/download/embe-ls-windows-amd64.zip -OutFile embe-ls.zip
Expand-Archive -LiteralPath .\embe-ls.zip -DestinationPath $EmbeDir
Rename-Item -Path $EmbeDir\README.md -NewName README-embe-ls.md
Rename-Item -Path $EmbeDir\LICENSE -NewName LICENSE-embe-ls
rm embe-ls.zip

if (Get-Command code -ErrorAction SilentlyContinue) { 
	Write-Host "Installing vscode-embe..."
	Invoke-WebRequest -Uri https://github.com/Bananenpro/vscode-embe/releases/latest/download/embe.vsix -OutFile .\embe.vsix
	code --uninstall-extension bananenpro.embe
	code --install-extension embe.vsix
	rm embe.vsix
}

Write-Host "Done."
Write-Host "Please reboot for the installation to take effect." -ForegroundColor Yellow
