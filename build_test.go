package staticfile_test

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"

	. "github.com/onsi/gomega"
	"github.com/paketo-buildpacks/packit/v2"
	"github.com/paketo-buildpacks/packit/v2/chronos"
	"github.com/paketo-buildpacks/packit/v2/scribe"
	"github.com/paketo-community/staticfile"
	"github.com/paketo-community/staticfile/fakes"
	"github.com/sclevine/spec"
)

func testBuild(t *testing.T, when spec.G, it spec.S) {
	var (
		Expect = NewWithT(t).Expect

		layersDir  string
		workingDir string
		cnbDir     string
		buffer     *bytes.Buffer
		timeStamp  time.Time

		clock chronos.Clock

		installProcess *fakes.InstallProcess
		bpYMLParser    *fakes.BpYMLParser
		scriptWriter   *fakes.ScriptWriter
		entryResolver  *fakes.EntryResolver

		build packit.BuildFunc
	)

	it.Before(func() {
		var err error
		layersDir, err = os.MkdirTemp("", "layers")
		Expect(err).NotTo(HaveOccurred())

		workingDir, err = os.MkdirTemp("", "working-dir")
		Expect(err).NotTo(HaveOccurred())

		installProcess = &fakes.InstallProcess{}
		bpYMLParser = &fakes.BpYMLParser{}
		scriptWriter = &fakes.ScriptWriter{}

		entryResolver = &fakes.EntryResolver{}
		entryResolver.ResolveCall.Returns.BuildpackPlanEntry = packit.BuildpackPlanEntry{
			Name: "staticfile",
			Metadata: map[string]interface{}{
				"launch": true,
			},
		}

		buffer = bytes.NewBuffer(nil)

		timeStamp = time.Now()
		clock = chronos.NewClock(func() time.Time {
			return timeStamp
		})

		build = staticfile.Build(
			installProcess,
			bpYMLParser,
			scriptWriter,
			entryResolver,
			scribe.NewEmitter(buffer),
			clock,
		)
	})

	it.After(func() {
		Expect(os.RemoveAll(layersDir)).To(Succeed())
		Expect(os.RemoveAll(cnbDir)).To(Succeed())
		Expect(os.RemoveAll(workingDir)).To(Succeed())
	})

	when("the user specifies an nginx sever", func() {
		it("returns a result that writes an nginx conf file", func() {
			config := staticfile.Config{
				Nginx: &staticfile.Nginx{
					LocationInclude: "some-location",
				},
			}

			bpYMLParser.ParseCall.Returns.Config = config

			buildContext := packit.BuildContext{
				WorkingDir: workingDir,
				Stack:      "some-stack",
				BuildpackInfo: packit.BuildpackInfo{
					Name:    "Some Buildpack",
					Version: "some-version",
				},
				Plan: packit.BuildpackPlan{
					Entries: []packit.BuildpackPlanEntry{
						{
							Name: "staticfile",
							Metadata: map[string]interface{}{
								"launch": true,
							},
						},
					},
				},
				Layers: packit.Layers{Path: layersDir},
			}

			result, err := build(buildContext)

			Expect(err).NotTo(HaveOccurred())
			Expect(result).To(Equal(packit.BuildResult{
				Layers: []packit.Layer{
					{
						Name:             "staticfile",
						Path:             filepath.Join(layersDir, "staticfile"),
						LaunchEnv:        packit.Environment{},
						ProcessLaunchEnv: map[string]packit.Environment{},
						BuildEnv:         packit.Environment{},
						SharedEnv: packit.Environment{
							"APP_ROOT.default": workingDir,
						},
						Build:  false,
						Launch: true,
						Cache:  false,
					},
				},
			}))

			Expect(entryResolver.ResolveCall.Receives.Name).To(Equal("staticfile"))
			Expect(entryResolver.ResolveCall.Receives.Entries).To(Equal([]packit.BuildpackPlanEntry{
				{
					Name: "staticfile",
					Metadata: map[string]interface{}{
						"launch": true,
					},
				},
			}))

			Expect(scriptWriter.WriteInitScriptCall.CallCount).To(Equal(1))
			Expect(scriptWriter.WriteInitScriptCall.Receives.ProfileDPath).To(Equal(filepath.Join(layersDir, "staticfile", "profile.d")))

			Expect(scriptWriter.WriteStartLoggingScriptCall.CallCount).To(Equal(1))
			Expect(scriptWriter.WriteStartLoggingScriptCall.Receives.ProfileDPath).To(Equal(filepath.Join(layersDir, "staticfile", "profile.d")))

			Expect(installProcess.ExecuteCall.CallCount).To(Equal(1))
			Expect(installProcess.ExecuteCall.Receives.Context).To(Equal(buildContext))
			Expect(installProcess.ExecuteCall.Receives.TemplConfig).To(Equal(config))

			Expect(buffer.String()).To(ContainSubstring("Some Buildpack some-version"))
			Expect(buffer.String()).To(ContainSubstring("Parsing buildpack.yml for nginx config"))
			Expect(buffer.String()).To(ContainSubstring("Writing profile.d scripts"))
			Expect(buffer.String()).To(ContainSubstring("Executing build process"))
			Expect(buffer.String()).To(ContainSubstring("Configuring build environment"))
			Expect(buffer.String()).To(ContainSubstring("Configuring launch environment"))
		})
	})

	when("the layers directory cannot be read", func() {
		it.Before(func() {
			Expect(os.Chmod(layersDir, 0000)).To(Succeed())
		})

		it.After(func() {
			Expect(os.Chmod(layersDir, os.ModePerm)).To(Succeed())
		})

		it("returns an error", func() {
			_, err := build(packit.BuildContext{
				WorkingDir: workingDir,
				CNBPath:    cnbDir,
				Stack:      "some-stack",
				BuildpackInfo: packit.BuildpackInfo{
					Name:    "Some Buildpack",
					Version: "some-version",
				},
				Plan: packit.BuildpackPlan{
					Entries: []packit.BuildpackPlanEntry{
						{
							Name: "staticfile",
						},
					},
				},
				Layers: packit.Layers{Path: layersDir},
			})
			Expect(err).To(MatchError(ContainSubstring("failed to get layer")))
		})
	})

	when("parsing the builpack.yml fails", func() {
		it("errors", func() {
			bpYMLParser.ParseCall.Returns.Err = errors.New("some-error")

			buildContext := packit.BuildContext{
				WorkingDir: workingDir,
				Stack:      "some-stack",
				BuildpackInfo: packit.BuildpackInfo{
					Name:    "Some Buildpack",
					Version: "some-version",
				},
				Plan: packit.BuildpackPlan{
					Entries: []packit.BuildpackPlanEntry{
						{
							Name: "staticfile",
						},
					},
				},
				Layers: packit.Layers{Path: layersDir},
			}

			_, err := build(buildContext)
			Expect(err).To(HaveOccurred())
			Expect(err).To(MatchError(fmt.Errorf("failed to parse buildpack.yml: some-error")))
		})
	})

	when("writing the init script fails", func() {
		it("errors", func() {
			config := staticfile.Config{
				Nginx: &staticfile.Nginx{
					LocationInclude: "some-location",
				},
			}

			bpYMLParser.ParseCall.Returns.Config = config

			scriptWriter.WriteInitScriptCall.Returns.Error = errors.New("some-error")

			buildContext := packit.BuildContext{
				WorkingDir: workingDir,
				Stack:      "some-stack",
				BuildpackInfo: packit.BuildpackInfo{
					Name:    "Some Buildpack",
					Version: "some-version",
				},
				Plan: packit.BuildpackPlan{
					Entries: []packit.BuildpackPlanEntry{
						{
							Name: "staticfile",
						},
					},
				},
				Layers: packit.Layers{Path: layersDir},
			}

			_, err := build(buildContext)
			Expect(err).To(MatchError(fmt.Errorf("failed to write init script: some-error")))

		})
	})

	when("writing the start logging script fails", func() {
		it("errors", func() {
			config := staticfile.Config{
				Nginx: &staticfile.Nginx{
					LocationInclude: "some-location",
				},
			}

			bpYMLParser.ParseCall.Returns.Config = config

			scriptWriter.WriteStartLoggingScriptCall.Returns.Error = errors.New("some-error")

			buildContext := packit.BuildContext{
				WorkingDir: workingDir,
				Stack:      "some-stack",
				BuildpackInfo: packit.BuildpackInfo{
					Name:    "Some Buildpack",
					Version: "some-version",
				},
				Plan: packit.BuildpackPlan{
					Entries: []packit.BuildpackPlanEntry{
						{
							Name: "staticfile",
						},
					},
				},
				Layers: packit.Layers{Path: layersDir},
			}

			_, err := build(buildContext)
			Expect(err).To(MatchError(fmt.Errorf("failed to write start_logging script: some-error")))

		})
	})

	when("installing the config file fails", func() {
		it("errors", func() {
			installProcess.ExecuteCall.Returns.Error = errors.New("some-error")

			buildContext := packit.BuildContext{
				WorkingDir: workingDir,
				Stack:      "some-stack",
				BuildpackInfo: packit.BuildpackInfo{
					Name:    "Some Buildpack",
					Version: "some-version",
				},
				Plan: packit.BuildpackPlan{
					Entries: []packit.BuildpackPlanEntry{
						{
							Name: "staticfile",
						},
					},
				},
				Layers: packit.Layers{Path: layersDir},
			}

			_, err := build(buildContext)
			Expect(err).To(MatchError(fmt.Errorf("failed to install config: some-error")))
		})
	})
}
