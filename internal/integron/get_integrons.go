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

			gi.logger.Infof("proccessing fasta %s", fastaFile)
			resultFolder, err := gi.finder.Run(downloadDir, fastaFile)
			if err != nil {

				failed = failed + 1
				gi.logger.Warnf("error %v", err)
				continue
			}

			gi.logger.Info("cleaning results")
			err = gi.cleaner.Clean(resultFolder)
			if err != nil {

				failed = failed + 1
				gi.logger.Warnf("error %v", err)
				continue
			}

			gi.logger.Infof("parsing results %s", fastaFile)
			integrons, err := gi.parser.Parse(resultFolder)

			if err != nil {
				failed = failed + 1
				gi.logger.Warnf("error %v", err)
				continue
			}
			gi.logger.Infof("integrons found %d", len(integrons))
			err = gi.fastaRepository.AddIntegrons(strain.Id, integrons)

			if err != nil {
				failed = failed + 1
				gi.logger.Warnf("error %v", err)
				continue
			}
			gi.logger.Infof("integrons saved %d", len(integrons))

			gi.logger.Info("removing folder")
			err = gi.cleaner.Remove(resultFolder)
			if err != nil {

				gi.logger.Warnf("error %v", err)
			}

			succeeded = succeeded + 1

		}
	}

	return nil
}
