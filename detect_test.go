package staticfile_test

import (
	"errors"
	"os"
	"testing"

	"github.com/paketo-buildpacks/packit/v2"
	"github.com/paketo-community/staticfile"
	"github.com/paketo-community/staticfile/fakes"
	"github.com/sclevine/spec"

	. "github.com/onsi/gomega"
)

func testDetect(t *testing.T, context spec.G, it spec.S) {
	var (
		Expect = NewWithT(t).Expect

		bpYMLParser *fakes.BpYMLParser
		workingDir  string
		detect      packit.DetectFunc
	)

	it.Before(func() {
		var err error
		workingDir, err = os.MkdirTemp("", "working-dir")
		Expect(err).NotTo(HaveOccurred())

		bpYMLParser = &fakes.BpYMLParser{}

		detect = staticfile.Detect(bpYMLParser)
	})

	it.After(func() {
		Expect(os.RemoveAll(workingDir)).To(Succeed())
	})

	context("when the buildpack.yml indicates it wants an nginx config", func() {
		it.Before(func() {
			bpYMLParser.ValidConfigCall.Returns.Valid = true
		})

		it("detects", func() {
			result, err := detect(packit.DetectContext{
				WorkingDir: workingDir,
			})
			Expect(err).NotTo(HaveOccurred())
			Expect(result.Plan).To(Equal(packit.BuildPlan{
				Provides: []packit.BuildPlanProvision{
					{
						Name: "staticfile",
					},
				},
				Requires: []packit.BuildPlanRequirement{
					{
						Name: "nginx",
						Metadata: staticfile.BuildPlanMetadata{
							Launch: true,
						},
					},
					{
						Name: "staticfile",
						Metadata: staticfile.BuildPlanMetadata{
							Launch: true,
						},
					},
				},
			}))
		})
	})

	context("when the buildpack.yml does not indicate it wants an nginx config", func() {
		it.Before(func() {
			bpYMLParser.ValidConfigCall.Returns.Valid = false
		})
		it("detect fails with an error", func() {
			_, err := detect(packit.DetectContext{
				WorkingDir: workingDir,
			})
			Expect(err).To(MatchError(packit.Fail))

		})
	})

	context("error cases", func() {
		context("when unable to parse buildpack.yml", func() {
			it("returns an error", func() {
				bpYMLParser.ValidConfigCall.Returns.Err = errors.New("some-error")

				_, err := detect(packit.DetectContext{
					WorkingDir: workingDir,
				})
				Expect(err).To(MatchError(`unable to parse buildpack.yml: "some-error"`))

			})
		})
	})
}
