package integron_test

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/ganvoa/biopipe-tools/internal/integron"
	"github.com/ganvoa/biopipe-tools/internal/platform"
	"github.com/otiai10/copy"
	"github.com/stretchr/testify/assert"
)

func TestWhenFolderHasFilesShouldPreserveFilesWithIntegronResultsOnly(t *testing.T) {
	baseFolder := filepath.Join(".", "cleaner_test_data", "Results_Integron_Finder_488112")
	newFolder := filepath.Join(".", "cleaner_test_data", "Results_Integron_Finder_488112_test")

	err := copy.Copy(baseFolder, newFolder)
	assert.NoError(t, err)
	logger := &platform.FakeLogger{}
	cleaner := integron.NewIntegronResultCleaner(logger)
	err = cleaner.Clean(newFolder)
	assert.NoError(t, err)

	files, _ := ioutil.ReadDir(newFolder)
	assert.Equal(t, 3, len(files))

	err = os.RemoveAll(newFolder)
	assert.NoError(t, err)
}

func TestWhenFoundDirectoryIsNotFoundShouldReturnError(t *testing.T) {
	folder := filepath.Join(".", "cleaner_test_data", "Results_Integron_Finder_Invalid")
	logger := &platform.FakeLogger{}
	cleaner := integron.NewIntegronResultCleaner(logger)
	err := cleaner.Clean(folder)
	assert.Error(t, err)
}
