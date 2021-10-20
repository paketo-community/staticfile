package main

import (
	"os"

	"github.com/paketo-buildpacks/packit"
	"github.com/paketo-buildpacks/packit/chronos"
	"github.com/paketo-buildpacks/packit/draft"
	"github.com/paketo-buildpacks/packit/scribe"
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
