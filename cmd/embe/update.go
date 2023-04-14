package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"mime"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/adrg/xdg"
)

// versionCheck returns true if a new version is available
func versionCheck(printWarning, ignoreCache bool) bool {
	if version == "dev" {
		return false
	}

	latest, err := getLatestVersion(ignoreCache)
	if err != nil {
		return false
	}

	if version != latest {
		if printWarning {
			fmt.Fprintf(stderr, "\x1b[33mWARNING:\x1b[0m A new version of embe is available. Run 'embe update' to install the latest version.\n")
		}
		return true
	}
	return false
}

func update() {
	if version == "dev" {
		printError(errors.New("Cannot update dev version."), nil, nil)
		os.Exit(1)
	}

	if !versionCheck(false, true) {
		fmt.Println("Embe is already up-to-date.")
		return
	}

	switch runtime.GOOS {
	case "windows":
		updateWindows()
	case "darwin", "linux":
		updateUnix()
	default:
		printError(errors.New("Automatic updates are not supported for your operating system."), nil, nil)
		os.Exit(1)
	}
}

func updateWindows() {
	fmt.Println("Downloading latest installer...")
	cmd := exec.Command("Powershell.exe", "-Command", "iwr -useb https://raw.githubusercontent.com/juho05/embe/main/install.ps1 | iex")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		printError(errors.New("Failed to update embe."), nil, nil)
		os.Exit(1)
	}
}

func updateUnix() {
	fmt.Println("Downloading latest installer...")

	var installCmd string
	if _, err := exec.LookPath("wget"); err == nil {
		installCmd = "wget -q --show-progress https://raw.githubusercontent.com/juho05/embe/main/install.sh -O- | bash"
	} else if _, err := exec.LookPath("curl"); err == nil {
		installCmd = "curl -L https://raw.githubusercontent.com/juho05/embe/main/install.sh | bash"
	} else {
		printError(errors.New("Please install either wget or curl."), nil, nil)
		os.Exit(1)
	}

	cmd := exec.Command("bash", "-c", installCmd)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		printError(errors.New("Failed to update embe."), nil, nil)
		os.Exit(1)
	}
}

func getLatestVersion(ignoreCache bool) (string, error) {
	cacheDir := filepath.Join(xdg.CacheHome, "embe")
	os.MkdirAll(cacheDir, 0o755)

	if !ignoreCache {
		content, err := os.ReadFile(filepath.Join(cacheDir, "latest_version"))
		if err == nil {
			parts := strings.Split(string(content), "\n")
			if len(parts) >= 2 {
				cacheTime, err := strconv.Atoi(parts[0])
				if err == nil && time.Now().Unix()-int64(cacheTime) <= 60*60*3 {
					return parts[1], nil
				}
			}
		}
	}

	return fetchLatestGithubTag("juho05", "embe")
}

func fetchLatestGithubTag(owner, repo string) (string, error) {
	res, err := http.Get(fmt.Sprintf("https://api.github.com/repos/%s/%s/tags", owner, repo))
	if err != nil || res.StatusCode != http.StatusOK || !hasContentType(res.Header, "application/json") {
		return "", fmt.Errorf("failed to access git tags from 'github.com/%s/%s'", owner, repo)
	}
	defer res.Body.Close()
	type response []struct {
		Name string `json:"name"`
	}
	var data response
	err = json.NewDecoder(res.Body).Decode(&data)
	if err != nil {
		return "", errors.New("failed to decode git tag data")
	}

	cacheDir := filepath.Join(xdg.CacheHome, "embe")
	os.MkdirAll(cacheDir, 0o755)
	os.WriteFile(filepath.Join(cacheDir, "latest_version"), []byte(fmt.Sprintf("%d\n%s", time.Now().Unix(), data[0].Name)), 0o644)

	return data[0].Name, nil
}

func hasContentType(h http.Header, mimetype string) bool {
	contentType := h.Get("content-type")
	if contentType == "" {
		return mimetype == "application/octet-stream"
	}

	for _, v := range strings.Split(contentType, ",") {
		t, _, err := mime.ParseMediaType(v)
		if err != nil {
			break
		}
		if t == mimetype {
			return true
		}
	}
	return false
}
