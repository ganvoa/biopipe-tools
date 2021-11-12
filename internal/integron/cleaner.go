package integron

import (
	"io/ioutil"
	"path/filepath"

	"github.com/ganvoa/biopipe-tools/internal"
)

type integronResultCleaner struct {
	resultDir string
	logger    internal.Logger
}

func NewIntegronResultCleaner(resultDir string, logger internal.Logger) integronResultCleaner {
	cleaner := integronResultCleaner{}
	cleaner.resultDir = resultDir
	cleaner.logger = logger
	return cleaner
}

func (cleaner integronResultCleaner) Clean() error {

	cleaner.logger.Info("looking for non empty integron files")

	files, err := ioutil.ReadDir(cleaner.resultDir)
	if err != nil {
		return err
	}

	for _, file := range files {

		if file.IsDir() {
			continue
		}

		hasResult, err := cleaner.fileHasResult(file.Name())
		if err != nil {
			return err
		}

		if !hasResult {
			cleaner.logger.Infof("delete file %s", file.Name())
			continue
		}

		cleaner.logger.Infof("file with results %s", file.Name())
	}

	return nil
}

func (cleaner integronResultCleaner) fileHasResult(filePath string) (bool, error) {
	fullPath := filepath.Join(cleaner.resultDir, filePath)
	extension := filepath.Ext(fullPath)

	if extension != ".integrons" {
		return false, nil
	}

	return true, nil
}
