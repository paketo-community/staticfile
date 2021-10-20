package staticfile

import (
	"fmt"
	"path/filepath"

	"github.com/paketo-buildpacks/packit"
)

//go:generate faux --interface BpYMLParser --output fakes/bp_yml_parser.go
type BpYMLParser interface {
	ValidConfig(path string) (valid bool, err error)
	Parse(path string) (config Config, err error)
}

type BuildPlanMetadata struct {
	Launch bool `toml:"launch"`
}

func Detect(bpYMLParser BpYMLParser) packit.DetectFunc {
	return func(context packit.DetectContext) (packit.DetectResult, error) {

		valid, err := bpYMLParser.ValidConfig(filepath.Join(context.WorkingDir, "buildpack.yml"))
		if err != nil {
			return packit.DetectResult{}, fmt.Errorf("unable to parse buildpack.yml: %q", err)
		}

		if !valid {
			return packit.DetectResult{}, packit.Fail
		}

		return packit.DetectResult{
			Plan: packit.BuildPlan{
				Provides: []packit.BuildPlanProvision{
					{Name: StaticfileDependency},
				},
				Requires: []packit.BuildPlanRequirement{
					{
						Name: NginxDependency,
						Metadata: BuildPlanMetadata{
							Launch: true,
						},
					},
					{
						Name: StaticfileDependency,
						Metadata: BuildPlanMetadata{
							Launch: true,
						},
					},
				},
			},
		}, nil
	}
}
