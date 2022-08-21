package generator

import (
	"archive/zip"
	"encoding/json"
	"io"
	"text/template"
	"time"

	_ "embed"

	"github.com/Bananenpro/embe/blocks"
)

//go:embed assets/project.json.tmpl
var projectTemplate string

//go:embed assets/mblock5.tmpl
var mblock5Template string

//go:embed assets/06d70cb3d65abe36615f0d51e08c3404.svg
var svg1 []byte

//go:embed assets/cd21514d0531fdffb22204e0ec5ed84a.svg
var svg2 []byte

//go:embed assets/83a9787d4cb6f3b7632b4ddfebf74367.wav
var wav []byte

//go:embed assets/mscratch.json
var mscratch []byte

func Package(writer io.Writer, blocks map[string]*blocks.Block) error {
	w := zip.NewWriter(writer)
	defer w.Close()

	err := createProject(w, blocks)
	if err != nil {
		return err
	}

	err = createMBlock5(w)
	if err != nil {
		return err
	}

	err = createAssets(w)
	if err != nil {
		return err
	}

	err = createMScratch(w)
	if err != nil {
		return err
	}

	return nil
}

func createProject(zw *zip.Writer, blockMap map[string]*blocks.Block) error {
	w, err := zw.Create("project.json")
	if err != nil {
		return err
	}

	tmpl, err := template.New("project.json").Parse(projectTemplate)
	if err != nil {
		return err
	}

	blockJSON, err := json.Marshal(blockMap)
	if err != nil {
		return err
	}

	tmpl.Execute(w, struct {
		Blocks string
	}{
		Blocks: string(blockJSON),
	})
	return nil
}

func createMBlock5(zw *zip.Writer) error {
	w, err := zw.Create("mblock5")
	if err != nil {
		return err
	}
	tmpl, err := template.New("mblock5").Parse(mblock5Template)
	if err != nil {
		return err
	}
	return tmpl.Execute(w, struct {
		CreatedAt int64 `json:"createdAt"`
	}{
		CreatedAt: time.Now().UnixMilli(),
	})
}

func createAssets(zw *zip.Writer) error {
	w, err := zw.Create("06d70cb3d65abe36615f0d51e08c3404.svg")
	if err != nil {
		return err
	}
	w.Write(svg1)

	w, err = zw.Create("cd21514d0531fdffb22204e0ec5ed84a.svg")
	if err != nil {
		return err
	}
	w.Write(svg2)

	w, err = zw.Create("83a9787d4cb6f3b7632b4ddfebf74367.wav")
	if err != nil {
		return err
	}
	w.Write(wav)
	return nil
}

func createMScratch(zw *zip.Writer) error {
	w, err := zw.Create("mscratch.json")
	if err != nil {
		return err
	}
	w.Write(mscratch)
	return nil
}
