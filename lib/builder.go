package gin

import (
	"fmt"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
)

type Builder interface {
	Build() error
	Binary() string
	Errors() string
}

type builder struct {
	dir      string
	binary   string
	errors   string
	useGodep bool
	wd       string
	mainPath string
}

func NewBuilder(dir string, bin string, useGodep bool, wd string, mainPath string) Builder {
	if len(bin) == 0 {
		bin = "bin"
	}

	// does not work on Windows without the ".exe" extension
	if runtime.GOOS == "windows" {
		if !strings.HasSuffix(bin, ".exe") { // check if it already has the .exe extension
			bin += ".exe"
		}
	}

	return &builder{dir: dir, binary: bin, useGodep: useGodep, wd: wd, mainPath: mainPath}
}

func (b *builder) Binary() string {
	return b.binary
}

func (b *builder) Errors() string {
	return b.errors
}

func (b *builder) Build() error {
	var command *exec.Cmd
	if b.useGodep {
		command = exec.Command("godep", "go", "build", "-o", filepath.Join(b.wd, b.binary), b.mainPath)
	} else {
		command = exec.Command("go", "build", "-o", filepath.Join(b.wd, b.binary), b.mainPath)
	}
	command.Dir = b.dir

	output, err := command.CombinedOutput()

	if command.ProcessState.Success() {
		b.errors = ""
	} else {
		b.errors = string(output)
	}

	if len(b.errors) > 0 {
		return fmt.Errorf(b.errors)
	}

	return err
}
