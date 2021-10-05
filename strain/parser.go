package strain

import (
	"encoding/json"
	"io/ioutil"
)

type StrainParser struct {
	filePath string
}

func NewStrainParser(filePath string) *StrainParser {
	sp := new(StrainParser)
	sp.filePath = filePath
	return sp
}

func (sp *StrainParser) GetStrains() ([]Strain, error) {
	strainJson, err := ioutil.ReadFile(sp.filePath)
	if err != nil {
		return nil, err
	}

	strains := struct {
		Strains []Strain
	}{}

	err = json.Unmarshal(strainJson, &strains)

	if err != nil {
		return nil, err
	}

	return strains.Strains, nil
}
