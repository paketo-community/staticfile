package main

import (
	"os"

	"github.com/paketo-buildpacks/packit"
	"github.com/paketo-buildpacks/packit/chronos"
	"github.com/paketo-community/staticfile"
)

func main() {
	parser := staticfile.NewBuildpackYMLParser()
	configInstaller := staticfile.NewConfigInstaller()
	profileDWriter := staticfile.NewProfileDWriter()
	logEmitter := staticfile.NewLogEmitter(os.Stdout)
	planEntryResolver := staticfile.NewPlanEntryResolver(logEmitter)

	packit.Run(
		staticfile.Detect(parser),
		staticfile.Build(
			configInstaller,
			parser,
			profileDWriter,
			planEntryResolver,
			logEmitter,
			chronos.DefaultClock,
		),
	)
}
