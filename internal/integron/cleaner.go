package integron

import (
	"bufio"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/ganvoa/biopipe-tools/internal"
)

type IntegronResultCleaner struct {
	logger internal.Logger
}

func NewIntegronResultCleaner(logger internal.Logger) IntegronResultCleaner {
	cleaner := IntegronResultCleaner{}
	cleaner.logger = logger
	return cleaner
}

func (cleaner IntegronResultCleaner) Clean(resultFolder string) error {

	cleaner.logger.Info("looking for non empty integron files")

	files, err := ioutil.ReadDir(resultFolder)
	if err != nil {
		return err
	}

	for _, file := range files {

		if file.IsDir() {
			continue
		}

		hasResult, err := cleaner.fileHasResult(resultFolder, file.Name())
		if err != nil {
			return err
		}

		if !hasResult {
			cleaner.logger.Infof("delete file %s", file.Name())
			fileToRemove := filepath.Join(resultFolder, file.Name())
			err = os.Remove(fileToRemove)
			if err != nil {
				return err
			}
			continue
		}

		cleaner.logger.Infof("file with results %s", file.Name())
	}

	return nil
}

func (cleaner IntegronResultCleaner) fileHasResult(resultFolder string, filePath string) (bool, error) {
	fullPath := filepath.Join(resultFolder, filePath)
	extension := filepath.Ext(fullPath)

	if extension != ".integrons" {
		return false, nil
	}

	file, err := os.Open(fullPath)
	if err != nil {
		return false, err
	}

	defer file.Close()

	reader := bufio.NewReader(file)
	line, _, err := reader.ReadLine()

	if err != nil {
		return false, err
	}

	if strings.Contains(string(line), "No Integron found") {
		return false, nil
	}

	return true, nil
}
