package staticfile

import (
	"fmt"
	"path/filepath"
	"time"

	"github.com/paketo-buildpacks/packit"
)

//go:generate faux --interface InstallProcess --output fakes/install_process.go
type InstallProcess interface {
	Execute(context packit.BuildContext, templConfig Config) error
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
	logger LogEmitter,
	clock Clock,
) packit.BuildFunc {
	return func(context packit.BuildContext) (packit.BuildResult, error) {
		logger.Title("%s %s", context.BuildpackInfo.Name, context.BuildpackInfo.Version)

		logger.Process("Parsing buildpack.yml for nginx config")
		config, err := bpYMLParser.Parse(filepath.Join(context.WorkingDir, "buildpack.yml"))
		if err != nil {
			return packit.BuildResult{}, fmt.Errorf("failed to parse buildpack.yml: %v", err)
		}

		staticfileLayer, err := context.Layers.Get(LayerNameStaticfile, packit.LaunchLayer)
		if err != nil {
			return packit.BuildResult{}, fmt.Errorf("failed to get layer: %v", err)
		}

		logger.Process("Writing profile.d scripts")
		err = scriptWriter.WriteInitScript(filepath.Join(staticfileLayer.Path, "profile.d"))
		if err != nil {
			return packit.BuildResult{}, fmt.Errorf("failed to write init script: %v", err)
		}

		err = scriptWriter.WriteStartLoggingScript(filepath.Join(staticfileLayer.Path, "profile.d"))
		if err != nil {
			return packit.BuildResult{}, fmt.Errorf("failed to write start_logging script: %v", err)
		}

		logger.Process("Executing build process")
		logger.Subprocess("Filling out nginx.conf template")
		then := clock.Now()
		err = installProcess.Execute(context, config)
		if err != nil {
			return packit.BuildResult{}, fmt.Errorf("failed to install config: %v", err)
		}
		logger.Action("Completed in %s", time.Since(then).Round(time.Millisecond))
		logger.Break()

		staticfileLayer.SharedEnv.Default("APP_ROOT", context.WorkingDir)
		logger.Environment(staticfileLayer.SharedEnv)

		return packit.BuildResult{
			Plan:   context.Plan,
			Layers: []packit.Layer{staticfileLayer},
		}, nil
	}
}
