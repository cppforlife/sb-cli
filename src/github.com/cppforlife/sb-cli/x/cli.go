package x

import (
	"bytes"
	"fmt"
	"io"
	"os/exec"

	boshui "github.com/cloudfoundry/bosh-cli/ui"
)

type CLIImpl struct {
	ui boshui.UI
}

func NewCLIImpl(ui boshui.UI) CLIImpl {
	return CLIImpl{ui}
}

func (d CLIImpl) Execute(args []string, stdin io.Reader) ([]byte, error) {
	cmd := exec.Command("bosh", args...) // todo use cmdRunner

	var outBuf, errBuf bytes.Buffer
	cmd.Stdout = &outBuf
	cmd.Stderr = &errBuf
	cmd.Stdin = stdin

	err := cmd.Run()
	if err != nil {
		return nil, fmt.Errorf("executing bosh: %s (stderr: %s)", err, errBuf.String())
	}

	return outBuf.Bytes(), nil
}
