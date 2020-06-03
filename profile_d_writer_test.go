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

func testProfileDWriter(t *testing.T, when spec.G, it spec.S) {
	var (
		Expect = NewWithT(t).Expect

		profileDWriter staticfile.ProfileDWriter
		layersDir      string
	)

	it.Before(func() {
		var err error
		layersDir, err = ioutil.TempDir("", "layers")
		Expect(err).NotTo(HaveOccurred())

		profileDWriter = staticfile.NewProfileDWriter()

	})

	it.After(func() {
		err := os.RemoveAll(layersDir)
		Expect(err).NotTo(HaveOccurred())
	})

	when("WriteInitScript", func() {
		it("writes the init script to the profile.d directory", func() {
			profileDDest := filepath.Join(layersDir, "profile.d")
			err := profileDWriter.WriteInitScript(profileDDest)
			Expect(err).NotTo(HaveOccurred())

			contents, err := ioutil.ReadFile(filepath.Join(profileDDest, "00_staticfile.sh"))
			Expect(err).NotTo(HaveOccurred())

			Expect(string(contents)).To(Equal(staticfile.InitScriptContents))

		})
	})
	when("WriteStartoggingScript", func() {
		it("writes the start logging script to the profile.d directory", func() {
			profileDDest := filepath.Join(layersDir, "profile.d")
			err := profileDWriter.WriteStartLoggingScript(profileDDest)
			Expect(err).NotTo(HaveOccurred())

			contents, err := ioutil.ReadFile(filepath.Join(profileDDest, "01_start_logging.sh"))
			Expect(err).NotTo(HaveOccurred())

			Expect(string(contents)).To(Equal(staticfile.StartLoggingContents))
		})

	})
	when("error cases", func() {
		when("it can not create the profile.d directory", func() {
			it("errors", func() {
				profileDDest := filepath.Join(layersDir, "profile.d")
				err := ioutil.WriteFile(profileDDest, []byte(``), 0000)
				Expect(err).NotTo(HaveOccurred())

				err = profileDWriter.WriteStartLoggingScript(profileDDest)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("could not create the profile.d layer directory"))
			})
		})

		when("it can not create the start_logging.sh script", func() {
			it("errors", func() {
				profileDDest := filepath.Join(layersDir, "profile.d")
				err := os.MkdirAll(profileDDest, 0744)
				Expect(err).NotTo(HaveOccurred())

				err = ioutil.WriteFile(filepath.Join(profileDDest, "01_start_logging.sh"), []byte(``), 0000)
				Expect(err).NotTo(HaveOccurred())

				err = profileDWriter.WriteStartLoggingScript(profileDDest)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("could not write the 01_start_logging.sh script"))
			})
		})
	})

}
