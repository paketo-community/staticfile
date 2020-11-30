package staticfile_test

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/paketo-community/staticfile"
	"github.com/sclevine/spec"

	. "github.com/onsi/gomega"
)

func testBuildpackYAMLParser(t *testing.T, when spec.G, it spec.S) {
	var (
		Expect             = NewWithT(t).Expect
		buildpackYMLParser staticfile.BuildpackYMLParser
		path               string
	)

	it.Before(func() {
		buildpackYMLParser = staticfile.NewBuildpackYMLParser()
	})

	when("Parse", func() {
		when("buildpack.yml exists", func() {

			it.Before(func() {
				file, err := ioutil.TempFile("", "buildpack.yml")
				Expect(err).NotTo(HaveOccurred())
				defer file.Close()

				path = file.Name()

				err = ioutil.WriteFile(path, []byte(`---
staticfile:
  nginx:
    root: some-root-dir
    host_dot_files: true
    location_include: some-location
    directory: true
    ssi: true
    pushstate: true
    http_strict_transport_security: true
    http_strict_transport_security_include_subdomains: true
    http_strict_transport_security_preload: true
    force_https: true
    basic_auth: true
    status_codes:
      some-staus: some-code
`), os.ModePerm)
				Expect(err).NotTo(HaveOccurred())

			})
			it("parses", func() {
				configData, err := buildpackYMLParser.Parse(path)
				Expect(err).NotTo(HaveOccurred())

				Expect(configData.Nginx.RootDir).To(Equal("some-root-dir"))
				Expect(configData.Nginx.HostDotFiles).To(BeTrue())
				Expect(configData.Nginx.LocationInclude).To(Equal("some-location"))
				Expect(configData.Nginx.DirectoryIndex).To(BeTrue())
				Expect(configData.Nginx.SSI).To(BeTrue())
				Expect(configData.Nginx.PushState).To(BeTrue())
				Expect(configData.Nginx.HSTSIncludeSubDomains).To(BeTrue())
				Expect(configData.Nginx.HSTSPreload).To(BeTrue())
				Expect(configData.Nginx.ForceHTTPS).To(BeTrue())
				Expect(configData.Nginx.BasicAuth).To(BeTrue())
				Expect(configData.Nginx.StatusCodes).To(Equal(map[string]string{
					"some-staus": "some-code",
				}))
			})

			when("the root dir is not specified", func() {
				it.Before(func() {
					err := ioutil.WriteFile(path, []byte(`---
staticfile:
  nginx: {}
`), os.ModePerm)
					Expect(err).NotTo(HaveOccurred())
				})
				it("sets RootDir to 'public'", func() {
					configData, err := buildpackYMLParser.Parse(path)
					Expect(err).NotTo(HaveOccurred())

					Expect(configData.Nginx.RootDir).To(Equal("public"))
				})
			})
		})

	})

	when("ValidConfig", func() {
		when("the buildpack.yml indicates the user wants to generate an nginx config", func() {
			it.Before(func() {
				file, err := ioutil.TempFile("", "buildpack.yml")
				Expect(err).NotTo(HaveOccurred())
				defer file.Close()

				path = file.Name()

				err = ioutil.WriteFile(path, []byte(`---
staticfile:
  nginx: {}
`), os.ModePerm)
				Expect(err).NotTo(HaveOccurred())

			})
			it("returns true", func() {
				valid, err := buildpackYMLParser.ValidConfig(path)
				Expect(err).NotTo(HaveOccurred())

				Expect(valid).To(BeTrue())
			})
		})

		when("the buildpack.yml indicates the user does not want to generate an nginx config", func() {
			it.Before(func() {
				file, err := ioutil.TempFile("", "buildpack.yml")
				Expect(err).NotTo(HaveOccurred())
				defer file.Close()

				path = file.Name()

				err = ioutil.WriteFile(path, []byte(`---
staticfile:
`), os.ModePerm)
				Expect(err).NotTo(HaveOccurred())

			})
			it("returns false", func() {
				valid, err := buildpackYMLParser.ValidConfig(path)
				Expect(err).NotTo(HaveOccurred())

				Expect(valid).To(BeFalse())
			})
		})

		when("buildpack.yml does not exist", func() {
			it("returns false", func() {
				path = filepath.Join("does-not-exist", "buildpack.yml")
				valid, err := buildpackYMLParser.ValidConfig(path)
				Expect(err).NotTo(HaveOccurred())

				Expect(valid).To(BeFalse())
			})
		})
	})

	when("error cases", func() {
		when("buildpack.yml is unable to be opened", func() {
			var workingDir string
			it.Before(func() {
				var err error
				workingDir, err = ioutil.TempDir("", "workingDir")
				Expect(err).NotTo(HaveOccurred())

				path = filepath.Join(workingDir, "buildpack.yml")
				err = ioutil.WriteFile(path, []byte(``), os.ModePerm)
				Expect(err).NotTo(HaveOccurred())

				Expect(os.Chmod(workingDir, 0000)).To(Succeed())
			})
			it.After(func() {
				Expect(os.Chmod(workingDir, os.ModePerm)).To(Succeed())
			})

			it("returns an error", func() {
				_, err := buildpackYMLParser.Parse(path)
				Expect(err).To(MatchError(ContainSubstring("unable to open buildpack.yml: ")))
			})
		})
		when("buildpack.yml is malformed", func() {
			it.Before(func() {
				file, err := ioutil.TempFile("", "buildpack.yml")
				Expect(err).NotTo(HaveOccurred())
				path = file.Name()
			})

			it("returns an error", func() {
				_, err := buildpackYMLParser.Parse(path)
				Expect(err).To(MatchError(ContainSubstring("unable to parse buildpack.yml: ")))
			})

		})
	})

}
