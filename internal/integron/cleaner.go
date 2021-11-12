package integron

import (
	"bufio"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

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
			fileToRemove := filepath.Join(cleaner.resultDir, file.Name())
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

func (cleaner integronResultCleaner) fileHasResult(filePath string) (bool, error) {
	fullPath := filepath.Join(cleaner.resultDir, filePath)
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
	// optionally, resize scanner's capacity for lines over 64K, see next example
	line, _, err := reader.ReadLine()

	if err != nil {
		return false, err
	}

	if strings.Contains(string(line), "No Integron found") {
		return false, nil
	}

	return true, nil
}
