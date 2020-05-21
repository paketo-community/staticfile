package main_test

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"testing"

	"github.com/paketo-buildpacks/packit"
	main "github.com/paketo-buildpacks/staticfile"
	"github.com/paketo-buildpacks/staticfile/fakes"
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
		workingDir, err = ioutil.TempDir("", "working-dir")
		Expect(err).NotTo(HaveOccurred())

		bpYMLParser = &fakes.BpYMLParser{}

		detect = main.Detect(bpYMLParser)
	})

	it.After(func() {
		Expect(os.RemoveAll(workingDir)).To(Succeed())
	})

	context("when the buildpack.yml says to use the Staticfile buildpack", func() {
		it.Before(func() {
			bpYMLParser.ParseCall.Returns.Server = "nginx"
		})

		it("detects", func() {
			result, err := detect(packit.DetectContext{
				WorkingDir: workingDir,
			})
			Expect(err).NotTo(HaveOccurred())
			Expect(result.Plan).To(Equal(packit.BuildPlan{
				Requires: []packit.BuildPlanRequirement{
					{
						Name: "nginx",
					},
				},
			}))
		})
	})

	context("when the buildpack.yml specifies an invalid server", func() {
		it.Before(func() {
			bpYMLParser.ParseCall.Returns.Server = "invalid-server"
		})

		it("detects", func() {
			_, err := detect(packit.DetectContext{
				WorkingDir: workingDir,
			})
			Expect(err).To(MatchError(fmt.Errorf(`"invalid-server" is not a supported server: supported servers are: [%s]`, strings.Join(main.SupportedServers, ","))))
		})
	})

	context("when the buildpack.yml does not specify a server", func() {
		it.Before(func() {
			bpYMLParser.ParseCall.Returns.Server = ""
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
				bpYMLParser.ParseCall.Returns.Server = ""
				bpYMLParser.ParseCall.Returns.Err = errors.New("big bad error")

				_, err := detect(packit.DetectContext{
					WorkingDir: workingDir,
				})
				Expect(err).To(MatchError(`unable to parse buildpack.yml: "big bad error"`))

			})
		})
	})
}
