package strain

import (
	"github.com/ganvoa/biopipe-tools/internal"
)

type StrainBakuper struct {
	repository StrainRepository
	parser     strainParser
	logger     internal.Logger
}

func NewStrainBackuper(repository StrainRepository, parser strainParser, logger internal.Logger) *StrainBakuper {
	s := new(StrainBakuper)
	s.repository = repository
	s.parser = parser
	s.logger = logger
	return s
}

func (sb StrainBakuper) Backup() error {

	strains, err := sb.parser.GetStrains()
	if err != nil {
		return err
	}

	sb.logger.Infof("%d strains found", len(strains))
	sb.repository.SaveAll(strains)
	sb.logger.Infof("strains saved", len(strains))

	return nil
}
