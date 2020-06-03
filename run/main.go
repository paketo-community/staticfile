package main

import (
	"os"
	"time"

	"github.com/paketo-buildpacks/packit"
	"github.com/paketo-community/staticfile"
)

func main() {
	parser := staticfile.NewBuildpackYMLParser()
	configInstaller := staticfile.NewConfigInstaller()
	profileDWriter := staticfile.NewProfileDWriter()
	logEmitter := staticfile.NewLogEmitter(os.Stdout)
	clock := staticfile.NewClock(time.Now)

	packit.Run(
		staticfile.Detect(parser),
		staticfile.Build(
			configInstaller,
			parser,
			profileDWriter,
			logEmitter,
			clock,
		),
	)
}
