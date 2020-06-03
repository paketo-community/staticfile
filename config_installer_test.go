package staticfile_test

import (
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/paketo-buildpacks/packit"
	"github.com/paketo-community/staticfile"

	"github.com/sclevine/spec"

	. "github.com/onsi/gomega"
)

func testConfigInstaller(t *testing.T, when spec.G, it spec.S) {
	var (
		Expect          = NewWithT(t).Expect
		configInstaller staticfile.ConfigInstaller
		buildpackDir    string
		workingDir      string
		config          staticfile.Config
	)
	it.Before(func() {
		var err error
		buildpackDir, err = ioutil.TempDir("", "buildpackDir")
		Expect(err).NotTo(HaveOccurred())

		workingDir, err = ioutil.TempDir("", "workingDir")
		Expect(err).NotTo(HaveOccurred())

		config = staticfile.Config{
			Nginx: &staticfile.Nginx{
				LocationInclude: "some-location",
			},
		}

		configInstaller = staticfile.NewConfigInstaller()
	})

	it.After(func() {
		err := os.RemoveAll(workingDir)
		Expect(err).NotTo(HaveOccurred())

		err = os.RemoveAll(buildpackDir)
		Expect(err).NotTo(HaveOccurred())
	})

	when("when executing config installer", func() {
		it.Before(func() {
			Expect(os.MkdirAll(filepath.Join(buildpackDir, "server_configs"), os.ModePerm)).To(Succeed())

			confPath := filepath.Join(buildpackDir, "server_configs", "nginx.conf")

			err := ioutil.WriteFile(confPath, []byte(`$(( .LocationInclude ))`), os.ModePerm)
			Expect(err).NotTo(HaveOccurred())
		})

		it("it provides a nginx.conf file", func() {
			err := configInstaller.Execute(packit.BuildContext{
				CNBPath:    buildpackDir,
				WorkingDir: workingDir,
				Plan: packit.BuildpackPlan{
					Entries: []packit.BuildpackPlanEntry{
						{
							Name: "staticfile",
						},
					},
				},
			},
				config,
			)

			Expect(err).NotTo(HaveOccurred())
			appConfPath := filepath.Join(workingDir, "nginx.conf")
			Expect(appConfPath).To(BeARegularFile())

			contents, err := ioutil.ReadFile(appConfPath)
			Expect(err).NotTo(HaveOccurred())

			Expect(string(contents)).To(Equal("some-location"))
		})
	})

	it("parses the nginx conf without erroring", func() {
		Expect(os.MkdirAll(filepath.Join(buildpackDir, "server_configs"), os.ModePerm)).To(Succeed())
		confPath := filepath.Join(buildpackDir, "server_configs", "nginx.conf")

		Copy("server_configs/nginx.conf", confPath)

		err := configInstaller.Execute(packit.BuildContext{
			CNBPath:    buildpackDir,
			WorkingDir: workingDir,
			Plan: packit.BuildpackPlan{
				Entries: []packit.BuildpackPlanEntry{
					{
						Name: "staticfile",
					},
				},
			},
		},
			config,
		)

		Expect(err).NotTo(HaveOccurred())

		appConfPath := filepath.Join(workingDir, "nginx.conf")
		Expect(appConfPath).To(BeARegularFile())
	})

	when("error cases", func() {
		when("the template file can not be found", func() {
			it.Before(func() {
				Expect(os.MkdirAll(filepath.Join(buildpackDir, "server_configs"), os.ModePerm)).To(Succeed())

				confPath := filepath.Join(buildpackDir, "server_configs", "does-not-exist")

				err := ioutil.WriteFile(confPath, []byte(`$(( .LocationInclude ))`), os.ModePerm)
				Expect(err).NotTo(HaveOccurred())
			})

			it("errors", func() {
				err := configInstaller.Execute(packit.BuildContext{
					CNBPath:    buildpackDir,
					WorkingDir: workingDir,
					Plan: packit.BuildpackPlan{
						Entries: []packit.BuildpackPlanEntry{
							{
								Name: "staticfile",
							},
						},
					},
				},
					config,
				)

				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("could not find template file:"))
			})
		})

		when("the filled out config file can not be writtewn to its destination", func() {
			it.Before(func() {
				Expect(os.MkdirAll(filepath.Join(buildpackDir, "server_configs"), os.ModePerm)).To(Succeed())

				confPath := filepath.Join(buildpackDir, "server_configs", "nginx.conf")

				err := ioutil.WriteFile(confPath, []byte(`$(( .LocationInclude ))`), os.ModePerm)
				Expect(err).NotTo(HaveOccurred())

				err = os.RemoveAll(workingDir)
				Expect(err).NotTo(HaveOccurred())
			})

			it("errors", func() {
				err := configInstaller.Execute(packit.BuildContext{
					CNBPath:    buildpackDir,
					WorkingDir: workingDir,
					Plan: packit.BuildpackPlan{
						Entries: []packit.BuildpackPlanEntry{
							{
								Name: "staticfile",
							},
						},
					},
				},
					config,
				)

				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("could not write nginx conf file to working dir:"))
			})
		})
	})

}

func Copy(src, dst string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()

	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, in)
	if err != nil {
		return err
	}
	return out.Close()
}
