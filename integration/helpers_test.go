package integration

import (
	"os"

	. "github.com/onsi/gomega"
)

func assertFileExists(path string) {
	_, err := os.Stat(path)
	if err != nil {
		ExpectWithOffset(1, err).To(Not(HaveOccurred()))
	}

}
