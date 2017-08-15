package cmd

import (
	"code.cloudfoundry.org/clock"
	boshui "github.com/cloudfoundry/bosh-cli/ui"
	boshcmd "github.com/cloudfoundry/bosh-utils/fileutil"
	boshlog "github.com/cloudfoundry/bosh-utils/logger"
	boshsys "github.com/cloudfoundry/bosh-utils/system"
	boshuuid "github.com/cloudfoundry/bosh-utils/uuid"
)

type BasicDeps struct {
	FS     boshsys.FileSystem
	UI     *boshui.ConfUI
	Logger boshlog.Logger

	UUIDGen    boshuuid.Generator
	CmdRunner  boshsys.CmdRunner
	Compressor boshcmd.Compressor

	Time clock.Clock
}

func NewBasicDeps(ui *boshui.ConfUI, logger boshlog.Logger) BasicDeps {
	return NewBasicDepsWithFS(ui, boshsys.NewOsFileSystemWithStrictTempRoot(logger), logger)
}

func NewBasicDepsWithFS(ui *boshui.ConfUI, fs boshsys.FileSystem, logger boshlog.Logger) BasicDeps {
	cmdRunner := boshsys.NewExecCmdRunner(logger)

	return BasicDeps{
		FS:     fs,
		UI:     ui,
		Logger: logger,

		UUIDGen:    boshuuid.NewGenerator(),
		CmdRunner:  cmdRunner,
		Compressor: boshcmd.NewTarballCompressor(cmdRunner, fs),

		Time: clock.NewClock(),
	}
}
