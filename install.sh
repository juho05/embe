#!/bin/bash

echo "Installing embe and embe-ls..."

cd /tmp

rm -f "embe.tar.gz"
rm -f "embe-ls.tar.gz"

os=$(uname)
arch=$(uname -m)

download () {
	if hash wget 2>/dev/null; then
		wget -q --show-progress https://github.com/Bananenpro/embe/releases/latest/download/embe-$1-$2.tar.gz -O embe.tar.gz || exit 1
		wget -q --show-progress https://github.com/Bananenpro/embe/releases/latest/download/embe-ls-$1-$2.tar.gz -O embe-ls.tar.gz || exit 1
	elif hash curl 2>/dev/null; then
		curl -L https://github.com/Bananenpro/embe/releases/latest/download/embe-$1-$2.tar.gz > embe.tar.gz || exit 1
		curl -L https://github.com/Bananenpro/embe/releases/latest/download/embe-ls-$1-$2.tar.gz > embe-ls.tar.gz || exit 1
	else
		echo "Please install either wget or curl."
		exit 1
	fi
}

shopt -s nocasematch

if [[ $os == *"linux"* ]]; then
	if [[ $arch == *"x86"* ]]; then
		echo "Detected OS: Linux x86_64"
		download "linux" "amd64"
	elif [[ $arch == *"aarch64"* ]]; then
		echo "Detected OS: Linux ARM64"
		download "linux" "arm64"
	elif [[ $arch == *"arm"* ]]; then
		echo "Detected OS: Linux ARM64"
		download "linux" "arm64"
	else
		echo "Detected OS: $os $arch"
		echo "Your architecture is not supported by this installer."
		exit 1
	fi
elif [[ $os == *"darwin"* ]]; then
	export PATH="$PATH:/Applications/Visual Studio Code.app/Contents/Resources/app/bin"
	if [[ $arch == *"x86"* ]]; then
		echo "Detected OS: macOS x86_64"
		download "darwin" "amd64"
	elif [[ $arch == *"aarch64"* ]]; then
		echo "Detected OS: macOS ARM64"
		download "darwin" "arm64"
	elif [[ $arch == *"arm"* ]]; then
		echo "Detected OS: macOS ARM64"
		download "darwin" "arm64"
	else
		echo "Detected OS: $os $arch"
		echo "Your architecture is not supported by this installer."
		exit 1
	fi
else
	echo "Detected OS: $os $arch"
	echo "Your OS is not supported by this installer."
	exit 1
fi

if [[ :$PATH: == *:"$HOME/.local/bin":* ]] ; then
	echo "Installing binaries into ~/.local/bin..."
	mkdir -p $HOME/.local/bin || exit 1
	if test -f /usr/local/bin/embe; then
		echo "Removing old version in /usr/local/bin..."
		sudo rm -f /usr/local/bin/embe
		sudo rm -f /usr/local/bin/embe-ls
	fi
	tar -xzf embe.tar.gz embe && mv embe $HOME/.local/bin || exit 1
	tar -xzf embe-ls.tar.gz embe-ls && mv embe-ls $HOME/.local/bin || exit 1
else
	echo "Installing binaries into /usr/local/bin..."
	sudo mkdir -p /usr/local/bin || exit 1
	if test -f $HOME/.local/bin/embe; then
		echo "Removing old version in ~/.local/bin..."
		rm -f $HOME/.local/bin/embe
		rm -f $HOME/.local/bin/embe-ls
	fi
	tar -xzf embe.tar.gz embe && sudo mv embe /usr/local/bin || exit 1
	tar -xzf embe-ls.tar.gz embe-ls && sudo mv embe-ls /usr/local/bin || exit 1
fi

rm embe.tar.gz
rm embe-ls.tar.gz

if hash code 2>/dev/null; then
	echo "Installing embe VS Code extension..."
	if hash wget 2>/dev/null; then
		wget -q --show-progress https://github.com/Bananenpro/vscode-embe/releases/latest/download/embe.vsix -O embe.vsix || exit 1
	else
		curl -L https://github.com/Bananenpro/vscode-embe/releases/latest/download/embe.vsix > embe.vsix || exit 1
	fi
	code --uninstall-extension bananenpro.embe &>/dev/null
	code --install-extension embe.vsix || exit 1
	rm embe.vsix
fi

echo "Done."
