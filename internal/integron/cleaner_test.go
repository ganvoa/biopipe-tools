package integron_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/ganvoa/biopipe-tools/internal/integron"
	"github.com/ganvoa/biopipe-tools/internal/platform"
	"github.com/otiai10/copy"
	"github.com/stretchr/testify/assert"
)

func TestRemoveFilesWithoutIntegronExtension(t *testing.T) {
	baseFolder := filepath.Join(".", "cleaner_test_data", "Results_Integron_Finder_488112")
	newFolder := filepath.Join(".", "cleaner_test_data", "Results_Integron_Finder_488112_test")
	err := copy.Copy(baseFolder, newFolder)
	assert.NoError(t, err)
	logger := &platform.FakeLogger{}
	cleaner := integron.NewIntegronResultCleaner(newFolder, logger)
	err = cleaner.Clean()
	assert.NoError(t, err)
	err = os.RemoveAll(newFolder)
	assert.NoError(t, err)
}

func TestErrorOnNotFoundDirectory(t *testing.T) {
	folder := filepath.Join(".", "cleaner_test_data", "Results_Integron_Finder_Invalid")
	logger := &platform.FakeLogger{}
	cleaner := integron.NewIntegronResultCleaner(folder, logger)
	err := cleaner.Clean()
	assert.Error(t, err)
}
