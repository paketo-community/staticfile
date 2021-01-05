package staticfile_test

import (
	"bytes"
	"testing"

	"github.com/paketo-buildpacks/packit"
	"github.com/paketo-community/staticfile"
	"github.com/sclevine/spec"

	. "github.com/onsi/gomega"
)

func testPlanEntryResolver(t *testing.T, context spec.G, it spec.S) {
	var (
		Expect = NewWithT(t).Expect

		buffer   *bytes.Buffer
		resolver staticfile.PlanEntryResolver
	)

	it.Before(func() {
		buffer = bytes.NewBuffer(nil)
		resolver = staticfile.NewPlanEntryResolver(staticfile.NewLogEmitter(buffer))
	})

	context("when entry flags differ", func() {
		context("OR's them together on best plan entry", func() {
			it("has all flags", func() {
				entry := resolver.Resolve([]packit.BuildpackPlanEntry{
					{
						Name: "staticfile",
						Metadata: map[string]interface{}{
							"launch": true,
						},
					},
					{
						Name: "staticfile",
						Metadata: map[string]interface{}{
							"build": true,
						},
					},
				})
				Expect(entry).To(Equal(packit.BuildpackPlanEntry{
					Name: "staticfile",
					Metadata: map[string]interface{}{
						"build":  true,
						"launch": true,
					},
				}))
			})
		})
	})
}
