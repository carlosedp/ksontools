package integration

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var (
	te *testEnv
)

var _ = BeforeSuite(func() {
	testDir, err := createTestDir()
	Expect(err).NotTo(HaveOccurred())

	testID := randString(6)

	te, err = newTestEnv(testDir, testID)
	Expect(err).NotTo(HaveOccurred())

	err = te.compilePlugin()
	Expect(err).NotTo(HaveOccurred())
})

var _ = AfterSuite(func() {
	err := te.cleanup()
	Expect(err).NotTo(HaveOccurred())
})

func TestIntegration(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "plugin integration")
}
