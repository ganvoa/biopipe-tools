package services

import "github.com/ganvoa/biopipe-tools/internal/integron"

type GetIntegrons struct {
	cleaner integron.IntegronResultCleaner
	finder  integron.IntegronFinder
	parser  integron.IntegronParser
}

func NewGetIntegron(finder integron.IntegronFinder, cleaner integron.IntegronResultCleaner, parser integron.IntegronParser) GetIntegrons {
	gi := GetIntegrons{}
	gi.finder = finder
	gi.cleaner = cleaner
	gi.parser = parser
	return gi
}

func (gi GetIntegrons) Run(fastaPath string) error {

	resultFolder, err := gi.finder.Run(fastaPath)
	if err != nil {
		return err
	}

	gi.cleaner.Clean(resultFolder)

	return nil
}
