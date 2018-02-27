package ksutil

import (
	"path/filepath"

	"github.com/ksonnet/ksonnet/metadata/app"
	"github.com/spf13/afero"
)

// App represents a ksonnet app.
type App struct {
	fs     afero.Fs
	appDir string
}

// NewApp creates an instance of App.
func NewApp(fs afero.Fs, appDir string) *App {
	return &App{
		fs:     fs,
		appDir: appDir,
	}
}

// AppSpec returns the specifcation for a ksonnet application.
func (a *App) AppSpec() (*app.Spec, error) {
	path := filepath.Join(a.appDir, "app.yaml")
	b, err := afero.ReadFile(a.fs, path)
	if err != nil {
		return nil, err
	}

	schema, err := app.Unmarshal(b)
	if err != nil {
		return nil, err
	}

	if schema.Contributors == nil {
		schema.Contributors = app.ContributorSpecs{}
	}

	if schema.Registries == nil {
		schema.Registries = app.RegistryRefSpecs{}
	}

	if schema.Libraries == nil {
		schema.Libraries = app.LibraryRefSpecs{}
	}

	if schema.Environments == nil {
		schema.Environments = app.EnvironmentSpecs{}
	}

	return schema, nil
}
