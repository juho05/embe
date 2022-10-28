package main

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
)

const docsURL = "https://github.com/Bananenpro/embe/blob/main/docs/documentation.md"

func docs() {
	var err error
	switch runtime.GOOS {
	case "linux":
		err = exec.Command("xdg-open", docsURL).Start()
	case "windows":
		err = exec.Command("rundll32", "url.dll,FileProtocolHandler", docsURL).Start()
	case "darwin":
		err = exec.Command("open", docsURL).Start()
	default:
		err = fmt.Errorf("unsupported platform")
	}
	if err != nil {
		printError(err, nil, nil)
		os.Exit(1)
	}
}
