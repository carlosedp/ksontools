package integration

import (
	"bytes"
	"os/exec"
	"path/filepath"
	"syscall"

	"github.com/bryanl/ksonnet/plugin"
	yaml "gopkg.in/yaml.v2"

	. "github.com/onsi/gomega"
)

func buildPluginConfig(te *testEnv) ([]byte, error) {
	config := plugin.Config{
		Name:        te.id,
		Version:     "0.1.0",
		Description: "component prototype",
		Command:     "$KS_PLUGIN_DIR/plugin",
	}

	return yaml.Marshal(&config)
}

type cmdOutput struct {
	stdout   string
	stderr   string
	exitCode int
}

func runWithOutput(cmd *exec.Cmd) (*cmdOutput, error) {
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	if err := cmd.Start(); err != nil {
		return nil, err
	}

	var exitCode int
	if err := cmd.Wait(); err != nil {
		switch t := err.(type) {
		default:
			return nil, err
		case *exec.ExitError:
			status, ok := t.Sys().(syscall.WaitStatus)
			if !ok {
				return nil, t
			}
			exitCode = status.ExitStatus()
		}
	}

	co := &cmdOutput{
		stdout:   stdout.String(),
		stderr:   stderr.String(),
		exitCode: exitCode,
	}

	return co, nil
}

func createDefaultComponents(te *testEnv, appDir string) {
	dataPath, err := filepath.Abs("testdata")
	ExpectWithOffset(1, err).To(Not(HaveOccurred()))

	dir := filepath.Join(dataPath, "rbac")
	co := te.runInApp(appDir, "import", "-f", dir)
	assertExitStatus(co, 0)
}
