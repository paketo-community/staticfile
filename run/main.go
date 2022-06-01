package main

import (
	"os"

	"github.com/paketo-buildpacks/packit/v2"
	"github.com/paketo-buildpacks/packit/v2/chronos"
	"github.com/paketo-buildpacks/packit/v2/draft"
	"github.com/paketo-buildpacks/packit/v2/scribe"
	"github.com/paketo-community/staticfile"
)

func main() {
	parser := staticfile.NewBuildpackYMLParser()

	packit.Run(
		staticfile.Detect(parser),
		staticfile.Build(
			staticfile.NewConfigInstaller(),
			parser,
			staticfile.NewProfileDWriter(),
			draft.NewPlanner(),
			scribe.NewEmitter(os.Stdout),
			chronos.DefaultClock,
		),
	)
}
