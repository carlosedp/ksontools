package integration

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"syscall"

	"github.com/bryanl/ksonnet/plugin"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	yaml "gopkg.in/yaml.v2"

	. "github.com/onsi/gomega"
)

type testEnv struct {
	dir string
	id  string

	home string
}

func newTestEnv(dir, id string) (*testEnv, error) {
	// TODO: make this run on windows by using the proper paths

	home := os.Getenv("HOME")
	if home == "" {
		return nil, errors.Errorf("unable to find home directory")

	}

	return &testEnv{
		dir:  dir,
		id:   id,
		home: home,
	}, nil
}

func (te *testEnv) run(options ...string) *cmdOutput {
	cmd := exec.Command("ks", append([]string{te.id}, options...)...)
	co, err := runWithOutput(cmd)
	ExpectWithOffset(1, err).ToNot(HaveOccurred())
	return co
}

func (te *testEnv) runInApp(appDir string, options ...string) *cmdOutput {
	ExpectWithOffset(1, appDir).To(BeADirectory())
	cmd := exec.Command("ks", append([]string{te.id, "--"}, options...)...)
	cmd.Dir = appDir

	co, err := runWithOutput(cmd)
	ExpectWithOffset(1, err).ToNot(HaveOccurred())

	return co
}

func (te *testEnv) pluginDir() string {
	return filepath.Join(te.home, ".config", "ksonnet", "plugins", te.id)
}

func createTestDir() (string, error) {
	return ioutil.TempDir("", "")
}

func deleteTestDir(name string) error {
	return os.RemoveAll(name)
}

func (te *testEnv) compilePlugin() error {
	if err := os.MkdirAll(te.pluginDir(), 0755); err != nil {
		return err
	}

	buildPath := filepath.Join(te.pluginDir(), "plugin")

	options := []string{
		"build",
		"-o",
		buildPath,
		"github.com/bryanl/woowoo/cmd/kscomp",
	}

	cmd := exec.Command("go", options...)
	if b, err := cmd.CombinedOutput(); err != nil {
		logrus.Error(string(b))
		return err
	}

	data, err := buildPluginConfig(te)
	if err != nil {
		return err
	}

	configPath := filepath.Join(te.pluginDir(), "plugin.yaml")
	if err := ioutil.WriteFile(configPath, data, 0644); err != nil {
		return err
	}

	return nil
}

func (te *testEnv) initApp() string {
	appID := randString(6)
	appDir := filepath.Join(te.dir, appID)
	options := []string{
		"init",
		appID,
		"--dir",
		appDir,
	}
	cmd := exec.Command("ks", options...)

	co, err := runWithOutput(cmd)
	Expect(err).ToNot(HaveOccurred())

	msg := fmt.Sprintf("exitCode: %d\nstdout:\n%s\nstderr:\n%s",
		co.exitCode, co.stdout, co.stderr)
	Expect(co.exitCode).To(Equal(0), msg)

	return appDir
}

func (te *testEnv) cleanup() error {
	if err := deleteTestDir(te.dir); err != nil {
		return err
	}

	return nil
}

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
