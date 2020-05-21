package main

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/paketo-buildpacks/packit"
)

//go:generate faux --interface BpYMLParser --output fakes/bp_yml_parser.go
type BpYMLParser interface {
	Parse(path string) (server string, err error)
}

var SupportedServers = []string{"nginx"}

func Detect(bpYMLParser BpYMLParser) packit.DetectFunc {
	return func(context packit.DetectContext) (packit.DetectResult, error) {

		server, err := bpYMLParser.Parse(filepath.Join(context.WorkingDir, "buildpack.yml"))
		if err != nil {
			return packit.DetectResult{}, fmt.Errorf("unable to parse buildpack.yml: %q", err)
		}

		if server == "" {
			return packit.DetectResult{}, packit.Fail
		} else if !serverSupported(server) {
			return packit.DetectResult{}, fmt.Errorf("%q is not a supported server: supported servers are: [%s]", server, strings.Join(SupportedServers, ","))
		}

		return packit.DetectResult{
			Plan: packit.BuildPlan{
				Requires: []packit.BuildPlanRequirement{
					{
						Name: server,
					},
				},
			},
		}, nil
	}
}

func serverSupported(server string) bool {
	for _, supportedServer := range SupportedServers {
		if server == supportedServer {
			return true
		}
	}
	return false
}
