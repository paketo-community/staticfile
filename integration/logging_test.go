package integration

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/paketo-buildpacks/occam"
	"github.com/sclevine/spec"

	. "github.com/onsi/gomega"
	. "github.com/paketo-buildpacks/occam/matchers"
)

func testLogging(t *testing.T, when spec.G, it spec.S) {
	var (
		Expect = NewWithT(t).Expect
		pack   occam.Pack
		docker occam.Docker
	)

	it.Before(func() {
		pack = occam.NewPack().WithVerbose()
		docker = occam.NewDocker()
	})

	when("when the buildpack is run with pack build", func() {
		var (
			image occam.Image

			name   string
			source string
		)

		it.Before(func() {
			var err error
			name, err = occam.RandomName()
			Expect(err).NotTo(HaveOccurred())
		})

		it.After(func() {
			Expect(docker.Image.Remove.Execute(image.ID)).To(Succeed())
			Expect(docker.Volume.Remove.Execute(occam.CacheVolumeNames(name))).To(Succeed())
			Expect(os.RemoveAll(source)).To(Succeed())
		})

		it("logs useful information for the user", func() {
			var err error
			source, err = occam.Source(filepath.Join("testdata", "nginx_helloworld"))
			Expect(err).NotTo(HaveOccurred())

			var logs fmt.Stringer
			image, logs, err = pack.WithNoColor().Build.
				WithPullPolicy("never").
				WithBuildpacks(nginxBuildpack, buildpack).
				Execute(name, source)
			Expect(err).NotTo(HaveOccurred(), logs.String)

			buildpackVersion, err := GetGitVersion()
			Expect(err).NotTo(HaveOccurred())

			Expect(logs).To(ContainLines(
				fmt.Sprintf("Staticfile Buildpack %s", buildpackVersion),
				"  Parsing buildpack.yml for nginx config",
				"  Writing profile.d scripts",
				"  Executing build process",
				`    Filling out nginx.conf template`,
				MatchRegexp(`      Completed in \d+\.?\d*`),
				"",
				"  Configuring environment",
				`    APP_ROOT -> "/workspace"`,
			))
		})
	})
}
