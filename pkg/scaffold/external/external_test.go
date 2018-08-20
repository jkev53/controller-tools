package external

import (
	"path/filepath"

	"fmt"

	// "strings"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"sigs.k8s.io/controller-tools/pkg/scaffold/input"
	"sigs.k8s.io/controller-tools/pkg/scaffold/resource"
	"sigs.k8s.io/controller-tools/pkg/scaffold/scaffoldtest"
)

var _ = Describe("Resource", func() {

	exts := []*External{
		{Resource: &resource.Resource{Group: "starfleet", Version: "v1", Kind: "NumberOne", Namespaced: true, CreateExampleReconcileBody: false}},
	}

	for i := range exts {
		e := exts[i]
		Describe(fmt.Sprintf("scaffolding API %s", e.Resource.Kind), func() {
			files := []struct {
				instance input.File
				file     string
			}{
				{
					file: filepath.Join("pkg", "apis",
						fmt.Sprintf("addtoscheme_%s_%s.go", e.Resource.Group, e.Resource.Version)),
					instance: &AddToScheme{External: e},
				},
			}

			for j := range files {
				f := files[j]
				Context(f.file, func() {
					It("should write a file matching the golden file", func() {
						s, result := scaffoldtest.NewTestScaffold(f.file, f.file)
						Expect(s.Execute(scaffoldtest.Options(), f.instance)).To(Succeed())
						Expect(result.Actual.String()).To(Equal(result.Golden), result.Actual.String())
					})
				})
			}
		})
	}
})
