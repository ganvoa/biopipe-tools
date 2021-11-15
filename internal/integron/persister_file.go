package integron

import (
	"encoding/json"
	"io/ioutil"

	"github.com/ganvoa/biopipe-tools/internal"
)

type FilePersister struct {
	outputPath string
	logger     internal.Logger
}

func NewFilePersister(outputPath string, logger internal.Logger) FilePersister {
	ipf := FilePersister{
		outputPath: outputPath,
		logger:     logger,
	}
	return ipf
}

func (ipf FilePersister) Save(integrons []Integron) error {
	output, err := json.Marshal(integrons)

	if err != nil {
		return err
	}

	ipf.logger.Infof("saving to file %s", ipf.outputPath)
	err = ioutil.WriteFile(ipf.outputPath, output, 0664)

	if err != nil {
		return err
	}

	return nil
}
