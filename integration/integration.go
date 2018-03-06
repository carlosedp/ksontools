package integration

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"

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

func (te *testEnv) run(options ...string) (string, string, error) {
	cmd := exec.Command("ks", append([]string{te.id}, options...)...)
	return runWithOutput(cmd)
}

func (te *testEnv) runInApp(appDir string, options ...string) {
	ExpectWithOffset(1, appDir).To(BeADirectory())
	cmd := exec.Command("ks", append([]string{te.id}, options...)...)
	cmd.Dir = appDir

	stdout, stderr, err := runWithOutput(cmd)
	msg := fmt.Sprintf("stdout:\n%s\nstderr:\n%s", stdout, stderr)
	ExpectWithOffset(0, err).ToNot(HaveOccurred(), msg)
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

	stdout, stderr, err := runWithOutput(cmd)
	msg := fmt.Sprintf("stdout:\n%s\nstderr:\n%s", stdout, stderr)
	Expect(err).ToNot(HaveOccurred(), msg)
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

func runWithOutput(cmd *exec.Cmd) (string, string, error) {
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		return "", "", err
	}

	return stdout.String(), stderr.String(), nil
}
