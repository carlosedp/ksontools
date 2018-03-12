package integration

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

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

	ExpectWithOffset(1, output).To(Equal(string(b)),
		"expected output to be:\n%s\nit was:\n%s\n",
		string(b), output)

}

func assertExitStatus(co *cmdOutput, status int) {
	ExpectWithOffset(1, co.exitCode).To(Equal(status),
		"expected exit status to be %d but was %d\nstdout:\n%s\nstderr:\n%s\nargs:%s\npath:%s",
		status, co.exitCode, co.stdout, co.stderr, strings.Join(co.args, " "), co.cmdName)
}
