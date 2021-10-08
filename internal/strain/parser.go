package strain

import (
	"encoding/json"
	"io/ioutil"
)

type strainParser struct {
	filePath string
}

func NewStrainParser(filePath string) strainParser {
	sp := strainParser{}
	sp.filePath = filePath
	return sp
}

func (sp *strainParser) GetStrains() ([]Strain, error) {
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
