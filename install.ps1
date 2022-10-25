if ((Get-ItemProperty 'HKCU:\Software\Microsoft\Windows\CurrentVersion\Internet Settings').ProxyEnable) {
    $proxy = (Get-ItemProperty 'HKCU:\Software\Microsoft\Windows\CurrentVersion\Internet Settings').ProxyServer
    $env:HTTP_PROXY = $proxy
    $env:HTTPS_PROXY = $proxy
	Write-Host "Using proxy: $proxy"
}

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
Invoke-WebRequest -Uri https://github.com/Bananenpro/embe/releases/latest/download/embe-ls-windows-amd64.zip -OutFile embe-ls.zip
Expand-Archive -LiteralPath .\embe-ls.zip -DestinationPath $EmbeDir
Rename-Item -Path $EmbeDir\README.md -NewName README-embe-ls.md
Rename-Item -Path $EmbeDir\LICENSE -NewName LICENSE-embe-ls
rm embe-ls.zip

if (Get-Command code -ErrorAction SilentlyContinue) {
	Write-Host "Installing vscode-embe..."
	Invoke-WebRequest -Uri https://github.com/Bananenpro/vscode-embe/releases/latest/download/embe.vsix -OutFile .\embe.vsix
	code --uninstall-extension bananenpro.embe | Out-Null
	code --install-extension embe.vsix
	rm embe.vsix
}

Write-Host "Refreshing environment variables..."
$HWND_BROADCAST = [intptr]0xffff;
$WM_SETTINGCHANGE = 0x1a;
$result = [uintptr]::zero
if (-not ("win32.nativemethods" -As [type])) {
	Add-Type -Namespace Win32 -Name NativeMethods -MemberDefinition @"
[DllImport("user32.dll", SetLastError = true, CharSet = CharSet.Auto)]
public static extern IntPtr SendMessageTimeout(
IntPtr hWnd, uint Msg, UIntPtr wParam, string lParam,
uint fuFlags, uint uTimeout, out UIntPtr lpdwResult);
"@
}
[void]([win32.nativemethods]::SendMessageTimeout($HWND_BROADCAST, $WM_SETTINGCHANGE, [uintptr]::Zero, "Environment", 2, 5000, [ref]$result))
$env:Path = [System.Environment]::GetEnvironmentVariable("Path","Machine") + ";" + [System.Environment]::GetEnvironmentVariable("Path","User")

Write-Host "Done." -ForegroundColor Green
