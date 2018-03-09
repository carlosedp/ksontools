package yaml

import (
	"bufio"
	"bytes"
	"io"
	"os"

	"github.com/pkg/errors"
	"github.com/spf13/afero"
)

const (
	docSeparator = "---"
)

// Decode decodes YAML into one or more readers.
func Decode(fs afero.Fs, source string) ([]io.Reader, error) {
	if err := checkSource(fs, source); err != nil {
		return nil, errors.Wrap(err, "check source")
	}

	f, err := fs.Open(source)
	if err != nil {
		return nil, errors.Wrap(err, "open source")
	}
	defer f.Close()

	bufs := make([]bytes.Buffer, 1)

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		t := scanner.Text()
		if t == docSeparator {
			bufs = append(bufs, bytes.Buffer{})
			continue
		}

		bufs[len(bufs)-1].WriteString(t)
		bufs[len(bufs)-1].WriteByte('\n')
	}

	var readers []io.Reader
	for i := range bufs {
		readers = append(readers, &bufs[i])
	}

	return readers, nil
}

func checkSource(fs afero.Fs, source string) error {
	if source == "" {
		return errors.New("source is empty")
	}

	if _, err := fs.Stat(source); err != nil {
		if os.IsNotExist(err) {
			return errors.Errorf("%q does not exist", source)
		}

		return errors.Wrap(err, "could not stat source")
	}

	return nil
}
