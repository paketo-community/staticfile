package main

import "github.com/paketo-buildpacks/packit"

func main() {
	parser := NewBuildpackYMLParser()
	packit.Run(Detect(parser), Build())
}
