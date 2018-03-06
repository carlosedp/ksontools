package integration

import (
	"os"
	"path/filepath"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Integration", func() {
	var err error
	var dataPath string

	BeforeEach(func() {
		dataPath, err = filepath.Abs("testdata")
		Expect(err).To(Not(HaveOccurred()))
	})

	Context("with a ks app", func() {
		var appDir string

		BeforeEach(func() {
			appDir = te.initApp()
		})

		AfterEach(func() {
			err = os.RemoveAll(appDir)
			Expect(err).To(Not(HaveOccurred()))
		})

		Context("import", func() {
			Context("file", func() {
				It("should create the component", func() {
					file := filepath.Join(dataPath, "rbac", "certificate-crd.yaml")
					te.runInApp(appDir, "import", "--", "-f", file)

					c := filepath.Join(appDir, "components", "certificate-crd.yaml")
					assertFileExists(c)
				})
			})
			Context("directory", func() {
				It("should create components for all the files in the directory", func() {
					dir := filepath.Join(dataPath, "rbac")
					te.runInApp(appDir, "import", "--", "-f", dir)

					expected := []string{"certificate-crd.yaml", "clusterissuer-crd.yaml",
						"deployment.yaml", "issuer-crd.yaml", "rbac.yaml", "serviceaccount.yaml"}

					for _, f := range expected {
						c := filepath.Join(appDir, "components", f)
						assertFileExists(c)
					}
				})
			})
		})
	})
})
