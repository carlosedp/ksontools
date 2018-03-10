// +build integration

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
			appDir = te.initApp("")
		})

		AfterEach(func() {
			err = os.RemoveAll(appDir)
			Expect(err).To(Not(HaveOccurred()))
		})

		Context("env", func() {
			Context("describe", func() {
				It("describes an environment", func() {
					co := te.runInApp(appDir, "env", "describe", "default")
					assertExitStatus(co, 0)
					assertOutput("env_describe.txt", co.stdout)
				})
			})

			Context("targets", func() {
				It("updates targets", func() {
					co := te.runInApp(appDir, "env", "describe", "default")
					assertExitStatus(co, 0)
					assertOutput("env_describe.txt", co.stdout)

					co = te.runInApp(appDir, "env", "targets", "default", "--component", "/")
					assertExitStatus(co, 0)

					co = te.runInApp(appDir, "env", "describe", "default")
					assertExitStatus(co, 0)
					assertOutput("env_describe_with_target.txt", co.stdout)

					co = te.runInApp(appDir, "env", "targets", "default")
					assertExitStatus(co, 0)

					co = te.runInApp(appDir, "env", "describe", "default")
					assertExitStatus(co, 0)
					assertOutput("env_describe.txt", co.stdout)
				})
			})
		})

		Context("component", func() {
			Context("list", func() {
				Context("with a namespace which has components", func() {
					It("lists the components in a namespace", func() {
						createDefaultComponents(te, appDir)

						co := te.runInApp(appDir, "component", "list")
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

				Context("in wide format", func() {
					Context("with a namespace which has components", func() {
						It("lists the components in a namespace", func() {
							createDefaultComponents(te, appDir)

							co := te.runInApp(appDir, "component", "list", "-o", "wide")
							assertExitStatus(co, 0)
							assertOutput("component_list_wide.txt", co.stdout)
						})
					})
					Context("with a namespace which has no components", func() {
						It("returns an empty list", func() {
							co := te.runInApp(appDir, "component", "list", "-o", "wide")
							assertExitStatus(co, 0)

							assertOutput("component_list_wide_empty.txt", co.stdout)
						})
					})
					Context("with an invalid namespace", func() {
						It("returns an empty list", func() {
							co := te.runInApp(appDir, "component", "list", "--ns", "invalid", "-o", "wide")
							assertExitStatus(co, 1)

							assertOutput("component_list_invalid.txt", co.stderr)
						})
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

		Context("ns", func() {
			Context("list", func() {
				It("it lists namespaces", func() {
					co := te.runInApp(appDir, "ns", "list")
					assertExitStatus(co, 0)
					assertOutput("ns_list.txt", co.stdout)
				})
			})
		})
		Context("param", func() {
			BeforeEach(func() {
				createDefaultComponents(te, appDir)
			})

			Context("delete", func() {
				Context("an existing parameter", func() {
					BeforeEach(func() {
						co := te.runInApp(appDir, "param", "set",
							"deployment", "metadata.name", "cert-manager2")
						assertExitStatus(co, 0)
					})

					It("deletes a parameter for a namespace", func() {
						co := te.runInApp(appDir, "param", "delete",
							"deployment", "metadata.name")
						assertExitStatus(co, 0)

						co = te.runInApp(appDir, "param", "list")
						assertOutput("param_list_empty.txt", co.stdout)
					})
				})
				Context("without an existing parameter", func() {
					It("deletes a parameter for a namespace", func() {
						co := te.runInApp(appDir, "param", "delete",
							"deployment", "metadata.name")
						assertExitStatus(co, 1)
					})
				})

				Context("an item with an index", func() {
					BeforeEach(func() {
						co := te.runInApp(appDir, "param", "set",
							"rbac", "metadata.name", "cert-manager2", "-i", "1")
						assertExitStatus(co, 0)
					})

					It("deletes a parameter", func() {
						co := te.runInApp(appDir, "param", "delete",
							"rbac", "metadata.name", "-i", "1")
						assertExitStatus(co, 0)

						co = te.runInApp(appDir, "param", "list")
						assertOutput("param_list_empty.txt", co.stdout)
					})
				})
			})

			Context("list", func() {
				It("lists parameters for a namespace", func() {
					co := te.runInApp(appDir, "param", "list")
					assertOutput("param_list_empty.txt", co.stdout)
				})

				Context("with custom parameters set", func() {
					It("lists parameters for a namespace", func() {
						co := te.runInApp(appDir, "param", "set",
							"deployment", "metadata.labels", `{"session":"session-a"}`)
						assertExitStatus(co, 0)

						co = te.runInApp(appDir, "param", "list")
						assertOutput("param_list_map_key.txt", co.stdout)
					})
				})
			})

			Context("set", func() {
				Context("map", func() {
					It("sets a map value", func() {
						co := te.runInApp(appDir, "param", "set",
							"deployment", "metadata.labels", `{"session":"session-a"}`)
						assertExitStatus(co, 0)

						co = te.runInApp(appDir, "show", "default", "-c", "deployment")
						assertOutput("param_set_map.txt", co.stdout)
					})
				})
				Context("array", func() {
					It("sets an array value", func() {
						co := te.runInApp(appDir, "param", "set",
							"deployment", "metadata.array", "[1,2,3,4]")
						assertExitStatus(co, 0)

						co = te.runInApp(appDir, "show", "default", "-c", "deployment")
						assertOutput("param_set_array.txt", co.stdout)
					})
				})
				Context("literal", func() {
					It("sets a literal value", func() {
						co := te.runInApp(appDir, "param", "set",
							"deployment", "metadata.name", "cert-manager2")
						assertExitStatus(co, 0)

						co = te.runInApp(appDir, "show", "default", "-c", "deployment")
						assertOutput("param_set_literal.txt", co.stdout)
					})
				})
			})

			Context("set-global", func() {
				Context("map", func() {
					It("sets a map value", func() {
						co := te.runInApp(appDir, "param", "set-global",
							"metadata.labels", `{"session":"session-a"}`)
						assertExitStatus(co, 0)

						co = te.runInApp(appDir, "show", "default", "-c", "deployment")
						assertOutput("param_set_map.txt", co.stdout)
					})
				})
				Context("array", func() {
					It("sets an array value", func() {
						co := te.runInApp(appDir, "param", "set-global",
							"metadata.array", "[1,2,3,4]")
						assertExitStatus(co, 0)

						co = te.runInApp(appDir, "show", "default", "-c", "deployment")
						assertOutput("param_set_array.txt", co.stdout)
					})
				})
				Context("literal", func() {
					It("sets a literal value", func() {
						co := te.runInApp(appDir, "param", "set-global",
							"metadata.name", "cert-manager2")
						assertExitStatus(co, 0)

						co = te.runInApp(appDir, "show", "default", "-c", "deployment")
						assertOutput("param_set_literal.txt", co.stdout)
					})
				})
			})

		})

		Context("show", func() {
			It("shows an environment", func() {
				dir := filepath.Join(dataPath, "rbac")
				co := te.runInApp(appDir, "import", "-f", dir)
				assertExitStatus(co, 0)

				co = te.runInApp(appDir, "show", "default")
				assertExitStatus(co, 0)
				assertOutput("show.txt", co.stdout)
			})
		})
	})
})
