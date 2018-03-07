package integration

import (
	"io/ioutil"
	"os"
	"path/filepath"

	. "github.com/onsi/gomega"
)

func assertFileExists(path string) {
	_, err := os.Stat(path)
	if err != nil {
		ExpectWithOffset(1, err).To(Not(HaveOccurred()))
	}
}

func assertOutput(name, output string) {
	path := filepath.Join("testdata", "output", name)
	ExpectWithOffset(1, path).To(BeAnExistingFile())

	b, err := ioutil.ReadFile(path)
	ExpectWithOffset(1, err).To(Not(HaveOccurred()))

	Expect(output).To(Equal(string(b)), "output did not match")
}

func assertExitStatus(co *cmdOutput, status int) {
	ExpectWithOffset(1, co.exitCode).To(Equal(status))
}
