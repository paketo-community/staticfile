package main_test

import (
	"io/ioutil"
	"os"
	"testing"

	main "github.com/paketo-buildpacks/staticfile"
	"github.com/sclevine/spec"

	. "github.com/onsi/gomega"
)

func testBuildpackYAMLParser(t *testing.T, when spec.G, it spec.S) {
	var (
		Expect             = NewWithT(t).Expect
		buildpackYMLParser main.BuildpackYMLParser
		path               string
	)

	it.Before(func() {
		buildpackYMLParser = main.NewBuildpackYMLParser()
	})

	when("buildpack.yml exists", func() {
		when("staticfile field is present", func() {
			it.Before(func() {
				file, err := ioutil.TempFile("", "buildpack.yml")
				Expect(err).NotTo(HaveOccurred())
				defer file.Close()

				path = file.Name()

				err = ioutil.WriteFile(path, []byte(`---
staticfile:
   server: nginx
`), os.ModePerm)
				Expect(err).NotTo(HaveOccurred())

			})
			it("parses", func() {
				server, err := buildpackYMLParser.Parse(path)
				Expect(err).NotTo(HaveOccurred())

				Expect(server).To(Equal("nginx"))
			})
		})
	})

}
