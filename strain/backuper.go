package strain

import "fmt"

type StrainBakuper struct {
	repository *StrainRepositoryElasticSearch
	parser     *StrainParser
}

func NewStrainBackuper(repository *StrainRepositoryElasticSearch, parser *StrainParser) *StrainBakuper {
	s := new(StrainBakuper)
	s.repository = repository
	s.parser = parser
	return s
}

func (sb StrainBakuper) Backup() error {

	strains, err := sb.parser.GetStrains()
	if err != nil {
		return err
	}

	fmt.Printf("Se encontraron %v strains\n", len(strains))

	sb.repository.SaveAll(strains)

	return nil
}
