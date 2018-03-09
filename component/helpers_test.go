package component

import (
	"path/filepath"

	"github.com/stretchr/testify/mock"

	"github.com/bryanl/woowoo/ksutil/mocks"
	"github.com/spf13/afero"
)

func appMock(root string) (*mocks.SuperApp, afero.Fs) {
	fs := afero.NewMemMapFs()
	app := &mocks.SuperApp{}
	app.On("Fs").Return(fs)
	app.On("Root").Return(root)
	app.On("LibPath", mock.AnythingOfType("string")).Return(filepath.Join(root, "lib", "v1.8.7"), nil)

	return app, fs

}
