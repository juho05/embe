package main

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/adrg/xdg"
)

func areYouSure(question string) bool {
	fmt.Print(question + " [y/N] ")
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	return strings.TrimSpace(strings.ToLower(scanner.Text())) == "y"
}

func uninstall() {
	if !areYouSure("Are you sure you want to uninstall embe?") {
		fmt.Println("Canceled.")
		return
	}

	switch runtime.GOOS {
	case "windows":
		uninstallWindows()
	case "darwin", "linux":
		uninstallUnix()
	default:
		printError(errors.New("Uninstallation is not supported for your operating system."), nil)
		os.Exit(1)
	}

	os.RemoveAll(filepath.Join(xdg.CacheHome, "embe"))

	uninstallVSCodeExt()
}

func uninstallWindows() {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		printError(fmt.Errorf("Failed to get the user home directory: %s", err), nil)
		os.Exit(1)
	}
	homeDir = filepath.Clean(homeDir)

	embeDir := homeDir + "\\AppData\\Local\\Programs\\embe"

	fmt.Println("Uninstalling embe and embe-ls...")
	cmd := exec.Command("Powershell.exe", "-Command", "[System.Environment]::SetEnvironmentVariable(\"PATH\", [System.Environment]::GetEnvironmentVariable(\"PATH\",\"USER\") -replace \";"+strings.ReplaceAll(embeDir, "\\", "\\\\")+"\",\"USER\")")
	cmd.Stderr = os.Stderr
	err = cmd.Run()
	if err != nil {
		printError(fmt.Errorf("Failed to remove %s from PATH: %s", embeDir, err), nil)
	}

	err = os.RemoveAll(embeDir)
	if err != nil {
		printError(fmt.Errorf("Failed to uninstall embe: %s", err), nil)
	}
}

func uninstallUnix() {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		printError(fmt.Errorf("Failed to get the user home directory: %s", err), nil)
		os.Exit(1)
	}
	homeDir = filepath.Clean(homeDir)

	fmt.Println("Uninstalling embe...")
	if _, err := os.Stat("/usr/local/bin/embe"); !os.IsNotExist(err) {
		cmd := exec.Command("bash", "-c", "sudo rm /usr/local/bin/embe")
		cmd.Stderr = os.Stderr
		err := cmd.Run()
		if err != nil {
			printError(fmt.Errorf("Failed to remove /usr/local/bin/embe: %s", err), nil)
		}
	}
	if _, err := os.Stat(homeDir + "/.local/bin/embe"); !os.IsNotExist(err) {
		err := os.Remove(homeDir + "/.local/bin/embe")
		if err != nil {
			printError(fmt.Errorf("Failed to remove %s: %s", homeDir+"/.local/bin/embe", err), nil)
		}
	}

	fmt.Println("Uninstalling embe-ls...")
	if _, err := os.Stat("/usr/local/bin/embe-ls"); !os.IsNotExist(err) {
		cmd := exec.Command("bash", "-c", "sudo rm /usr/local/bin/embe-ls")
		cmd.Stderr = os.Stderr
		err := cmd.Run()
		if err != nil {
			printError(fmt.Errorf("Failed to remove /usr/local/bin/embe-ls: %s", err), nil)
		}
	}
	if _, err := os.Stat(homeDir + "/.local/bin/embe-ls"); !os.IsNotExist(err) {
		err := os.Remove(homeDir + "/.local/bin/embe-ls")
		if err != nil {
			printError(fmt.Errorf("Failed to remove %s: %s", homeDir+"/.local/bin/embe-ls", err), nil)
		}
	}
}

func uninstallVSCodeExt() {
	if runtime.GOOS == "darwin" {
		os.Setenv("PATH", os.Getenv("PATH")+":/Applications/Visual Studio Code.app/Contents/Resources/app/bin")
	}
	if _, err := exec.LookPath("code"); err != nil {
		return
	}

	buf := &bytes.Buffer{}
	cmd := exec.Command("code", "--list-extensions")
	cmd.Stdout = buf
	err := cmd.Run()
	if err != nil || !strings.Contains(buf.String(), "bananenpro.embe") {
		return
	}

	fmt.Println("Uninstalling vscode-embe...")
	cmd = exec.Command("code", "--uninstall-extension", "bananenpro.embe")
	cmd.Stderr = os.Stderr
	err = cmd.Run()
	if err != nil {
		printError(fmt.Errorf("Failed to uninstall vscode-embe: %s.", err), nil)
	}
}
