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

		Context("component", func() {
			Context("list", func() {
				Context("with a namespace which has components", func() {
					It("lists the components in a namespace", func() {
						dir := filepath.Join(dataPath, "rbac")
						co := te.runInApp(appDir, "import", "-f", dir)
						assertExitStatus(co, 0)

						co = te.runInApp(appDir, "component", "list")
						assertExitStatus(co, 0)
						assertOutput("component_list.txt", co.stdout)
					})
				})
				Context("with a namespace which has no components", func() {
					It("returns an empty list", func() {
						co := te.runInApp(appDir, "component", "list")
						assertExitStatus(co, 0)

						assertOutput("component_list_empty.txt", co.stdout)
					})
				})
				Context("with an invalid namespace", func() {
					It("returns an empty list", func() {
						co := te.runInApp(appDir, "component", "list", "--ns", "invalid")
						assertExitStatus(co, 1)

						assertOutput("component_list_invalid.txt", co.stderr)
					})
				})
			})
		})

		Context("import", func() {
			Context("file", func() {
				It("should create the component", func() {
					file := filepath.Join(dataPath, "rbac", "certificate-crd.yaml")
					co := te.runInApp(appDir, "import", "-f", file)
					assertExitStatus(co, 0)

					c := filepath.Join(appDir, "components", "certificate-crd.yaml")
					assertFileExists(c)
				})
			})
			Context("directory", func() {
				It("should create components for all the files in the directory", func() {
					dir := filepath.Join(dataPath, "rbac")
					co := te.runInApp(appDir, "import", "-f", dir)
					assertExitStatus(co, 0)

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
