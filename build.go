package staticfile

import (
	"fmt"
	"path/filepath"
	"time"

	"github.com/paketo-buildpacks/packit/v2"
	"github.com/paketo-buildpacks/packit/v2/chronos"
	"github.com/paketo-buildpacks/packit/v2/scribe"
)

//go:generate faux --interface InstallProcess --output fakes/install_process.go
type InstallProcess interface {
	Execute(context packit.BuildContext, templConfig Config) error
}

//go:generate faux --interface EntryResolver --output fakes/entry_resolver.go
type EntryResolver interface {
	Resolve(name string, entries []packit.BuildpackPlanEntry, priorities []interface{}) (packit.BuildpackPlanEntry, []packit.BuildpackPlanEntry)
}

//go:generate faux --interface ScriptWriter --output fakes/script_writer.go
type ScriptWriter interface {
	WriteInitScript(profileDPath string) error
	WriteStartLoggingScript(profileDPath string) error
}

func Build(
	installProcess InstallProcess,
	bpYMLParser BpYMLParser,
	scriptWriter ScriptWriter,
	entryResolver EntryResolver,
	logger scribe.Emitter,
	clock chronos.Clock,
) packit.BuildFunc {
	return func(context packit.BuildContext) (packit.BuildResult, error) {
		logger.Title("%s %s", context.BuildpackInfo.Name, context.BuildpackInfo.Version)

		logger.Process("Parsing buildpack.yml for nginx config")
		config, err := bpYMLParser.Parse(filepath.Join(context.WorkingDir, "buildpack.yml"))
		if err != nil {
			return packit.BuildResult{}, fmt.Errorf("failed to parse buildpack.yml: %v", err)
		}

		layer, err := context.Layers.Get(LayerNameStaticfile)
		if err != nil {
			return packit.BuildResult{}, fmt.Errorf("failed to get layer: %v", err)
		}

		entry, _ := entryResolver.Resolve(StaticfileDependency, context.Plan.Entries, nil)

		layer, err = layer.Reset()
		if err != nil {
			return packit.BuildResult{}, err
		}

		layer.Launch = entry.Metadata["launch"] == true

		logger.Process("Writing profile.d scripts")
		err = scriptWriter.WriteInitScript(filepath.Join(layer.Path, "profile.d"))
		if err != nil {
			return packit.BuildResult{}, fmt.Errorf("failed to write init script: %v", err)
		}

		err = scriptWriter.WriteStartLoggingScript(filepath.Join(layer.Path, "profile.d"))
		if err != nil {
			return packit.BuildResult{}, fmt.Errorf("failed to write start_logging script: %v", err)
		}

		logger.Process("Executing build process")
		logger.Subprocess("Filling out nginx.conf template")

		duration, err := clock.Measure(func() error {
			return installProcess.Execute(context, config)
		})
		if err != nil {
			return packit.BuildResult{}, fmt.Errorf("failed to install config: %v", err)
		}

		logger.Action("Completed in %s", duration.Round(time.Millisecond))
		logger.Break()

		layer.SharedEnv.Default("APP_ROOT", context.WorkingDir)
		logger.EnvironmentVariables(layer)

		return packit.BuildResult{
			Layers: []packit.Layer{layer},
		}, nil
	}
}
