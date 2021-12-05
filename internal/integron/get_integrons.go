package integron

import (
	"fmt"

	"github.com/ganvoa/biopipe-tools/internal"
	"github.com/ganvoa/biopipe-tools/internal/fasta"
)

type GetIntegrons struct {
	cleaner         IntegronResultCleaner
	finder          IntegronFinder
	parser          IntegronParser
	fastaRepository fasta.FastaRepository
	logger          internal.Logger
}

func NewGetIntegron(finder IntegronFinder, cleaner IntegronResultCleaner, parser IntegronParser, repository fasta.FastaRepository, logger internal.Logger) GetIntegrons {
	gi := GetIntegrons{}
	gi.finder = finder
	gi.cleaner = cleaner
	gi.parser = parser
	gi.fastaRepository = repository
	gi.logger = logger
	return gi
}

func (gi GetIntegrons) Run(downloadDir string, outputDir string) error {

	succeeded := 0
	failed := 0
	currentFrom := 10000000000

	for {

		gi.logger.Debug("finding strains without integron finder")
		strains, err := gi.fastaRepository.FindWithoutIntegronResult(currentFrom)

		if err != nil {
			return err
		}

		gi.logger.Infof("running integron finder on %d strains", len(strains))
		if len(strains) == 0 {
			break
		}

		for _, strain := range strains {

			fastaFile := fmt.Sprintf("%d.fasta", strain.AssemblyId)

			gi.logger.Infof("proccessing fasta file %s", fastaFile)

			resultFolder, err := gi.finder.Run(downloadDir, fastaFile)
			if err != nil {

				failed = failed + 1
				gi.logger.Warnf("couldnt run integron finder on strain %d", strain.AssemblyId)
				continue
			}

			gi.logger.Infof("cleaning results %s", fastaFile)
			gi.cleaner.Clean(resultFolder)

			gi.logger.Infof("parsing results %s", fastaFile)
			integrons, err := gi.parser.Parse(resultFolder)

			if err != nil {
				failed = failed + 1
				gi.logger.Warnf("couldnt parse integrons on strain %d", strain.AssemblyId)
				continue
			}

			gi.logger.Infof("integrons found %d", len(integrons))
			err = gi.fastaRepository.AddIntegrons(strain.Id, integrons)

			if err != nil {
				failed = failed + 1
				gi.logger.Warnf("couldnt save integrons found on strain %d", strain.AssemblyId)
				continue
			}
			gi.logger.Infof("integrons saved %d", len(integrons))

			succeeded = succeeded + 1

		}
	}

	return nil
}
